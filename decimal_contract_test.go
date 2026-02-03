package files_sdk

import (
	"encoding/json"
	"strings"
	"testing"
)

type decimalContractPayload struct {
	Amount string  `json:"amount"`
	Ratio  float64 `json:"ratio"`
}

func TestDecimalContract_JSONMarshalling(t *testing.T) {
	payload := decimalContractPayload{Amount: "1.23", Ratio: 1.23}

	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	s := string(b)

	if !strings.Contains(s, `"amount":"1.23"`) {
		t.Fatalf("expected decimal to serialize as JSON string, got: %s", s)
	}

	if !strings.Contains(s, `"ratio":1.23`) {
		t.Fatalf("expected double to serialize as JSON number, got: %s", s)
	}
}
