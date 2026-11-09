package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
)

// failureCase is one case of the shared failure-case matrix (the mock server's scenarios/failure-cases.json).
type failureCase struct {
	ID     string           `json:"id"`
	Call   string           `json:"call"`
	UserID int64            `json:"user_id"`
	Params map[string]any   `json:"params"`
	Expect []map[string]any `json:"expect"`
}

// TestFailureCases runs the shared failure-case matrix through this package's public Client: Find, and UpdateWithMap,
// whose map carries the definition's invalid values as they are (UserUpdateParams types authenticate_until as a
// time.Time, which cannot hold an invalid date-time). simulator-tests.sh runs it through the mock server's launcher,
// which resets the simulator to the scenario's users and adds its fault rules; Config.RetryMax is the definition's
// max_network_retries. Each call's result or error is compared with the case's expectation: a ResponseError's status,
// type and message, a transport failure once the retries are spent, or the user found. FILES_MOCK_OUT/failure-cases.json
// records each case's outcome, and the launcher's failure cases check then compares the journal and the fault rules'
// decisions with each case's expected wire outcome. Without that launcher the test is skipped, unless
// FILES_API_SCENARIOS is "required".
func TestFailureCases(t *testing.T) {
	if os.Getenv("FILES_MOCK_SCENARIO") != "failure-cases" {
		why := "Run this test through targets/mock-server/with-scenario.rb --scenario failure-cases (simulator-tests.sh)."
		if os.Getenv("FILES_API_SCENARIOS") == "required" {
			t.Fatalf("FILES_API_SCENARIOS is required, so this test may not be skipped: %s", why)
		}
		t.Skip(why)
	}
	origin := os.Getenv("FILES_MOCK_SERVER_URL")
	if !strings.HasPrefix(origin, "http://127.0.0.1:") {
		t.Fatalf("FILES_MOCK_SERVER_URL is %q, not the loopback simulator with-scenario.rb starts.", origin)
	}
	awaitSimulator(t, &http.Client{Timeout: 10 * time.Second}, origin, os.Getenv("FILES_MOCK_SCHEMA_SHA256"), os.Getenv("FILES_MOCK_INSTANCE"))

	scenario, err := os.ReadFile(os.Getenv("FILES_MOCK_SCENARIO_FILE"))
	if err != nil {
		t.Fatal(err)
	}
	var definition struct {
		FailureCases struct {
			MaxNetworkRetries *int          `json:"max_network_retries"`
			Cases             []failureCase `json:"cases"`
		} `json:"failure_cases"`
	}
	if err := json.Unmarshal(scenario, &definition); err != nil || definition.FailureCases.MaxNetworkRetries == nil || len(definition.FailureCases.Cases) == 0 {
		t.Fatalf("FILES_MOCK_SCENARIO_FILE defines no failure_cases with max_network_retries (%v)", err)
	}
	config := files_sdk.Config{APIKey: os.Getenv("FILES_MOCK_API_KEY"), EndpointOverride: origin}.Init()
	config.RetryMax = *definition.FailureCases.MaxNetworkRetries
	client := &Client{Config: config}

	type result struct {
		ID      string  `json:"id"`
		Outcome string  `json:"outcome"`
		Failure *string `json:"failure"`
	}
	var results []result
	for _, testCase := range definition.FailureCases.Cases {
		var problems []string
		for index, expected := range testCase.Expect {
			actual := failureCaseCall(client, testCase)
			if !reflect.DeepEqual(actual, expected) {
				got, _ := json.Marshal(actual)
				want, _ := json.Marshal(expected)
				problems = append(problems, fmt.Sprintf("call %d gave %s, not %s", index+1, got, want))
			}
		}
		outcome := result{ID: testCase.ID, Outcome: "passed"}
		if len(problems) > 0 {
			failure := strings.Join(problems, "; ")
			outcome.Outcome, outcome.Failure = "failed", &failure
			t.Errorf("%s: %s", testCase.ID, failure)
		}
		results = append(results, outcome)
	}

	report, err := json.MarshalIndent(map[string]any{
		"format": "files-failure-cases/1", "sdk": "go", "runtime": runtime.Version(), "cases": results,
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(os.Getenv("FILES_MOCK_OUT"), "failure-cases.json"), append(report, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

// failureCaseCall makes one call of testCase and returns what it did in the expectation's terms, as JSON decodes them.
func failureCaseCall(client *Client, testCase failureCase) map[string]any {
	var actual map[string]any
	switch testCase.Call {
	case "users.update":
		params := map[string]any{"id": testCase.UserID}
		for name, value := range testCase.Params {
			params[name] = value
		}
		_, err := client.UpdateWithMap(params)
		actual = failureCaseOutcome(err, map[string]any{"updated": true})
	case "users.find":
		user, err := client.Find(files_sdk.UserFindParams{Id: testCase.UserID})
		actual = failureCaseOutcome(err, map[string]any{"user": map[string]any{"id": user.Id, "username": user.Username}})
	default:
		actual = map[string]any{"unknown call": testCase.Call}
	}
	encoded, _ := json.Marshal(actual)
	var decoded map[string]any
	_ = json.Unmarshal(encoded, &decoded)
	return decoded
}

func failureCaseOutcome(err error, success map[string]any) map[string]any {
	if err == nil {
		return success
	}
	var apiError files_sdk.ResponseError
	if errors.As(err, &apiError) {
		return map[string]any{"error": map[string]any{"status": apiError.HttpCode, "type": apiError.Type, "message": apiError.ErrorMessage}}
	}
	// A request that got no answer, after its retries: the transport's own error, which the SDK returns as it is.
	var transport *url.Error
	if errors.As(err, &transport) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return map[string]any{"error": map[string]any{"network": true}}
	}
	return map[string]any{"unexpected": fmt.Sprintf("%T: %v", err, err)}
}
