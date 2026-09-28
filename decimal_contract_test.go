package files_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// decimalInput sets the float64 and exact decimal fields shared by
// RemoteMountBackend Create and Update, so each case runs against both.
type decimalInput struct {
	cpu, mem               float64
	cpuDecimal, memDecimal *string
	fall                   int64 // an unrelated field, omitted when zero
}

// requests builds both request types with their own copies of the decimal
// text, so a test can detect writes through the pointers.
func (in decimalInput) requests() map[string]interface{} {
	cpuDecimal, memDecimal := copyText(in.cpuDecimal), copyText(in.memDecimal)
	return map[string]interface{}{
		"create": RemoteMountBackendCreateParams{
			MinFreeCpu: in.cpu, MinFreeCpuDecimal: cpuDecimal,
			MinFreeMem: in.mem, MinFreeMemDecimal: memDecimal,
			Fall: in.fall, Priority: 5, CanaryFilePath: "canary.txt", RemoteServerMountId: 2, RemoteServerId: 3,
		},
		"update": RemoteMountBackendUpdateParams{
			Id:         7,
			MinFreeCpu: in.cpu, MinFreeCpuDecimal: cpuDecimal,
			MinFreeMem: in.mem, MinFreeMemDecimal: memDecimal,
			Fall: in.fall, Priority: 5,
		},
	}
}

// The request fields other than the decimals, as lib.Params sends them.
var unrelatedWireFields = map[string]map[string]interface{}{
	"create": {"priority": json.Number("5"), "canary_file_path": "canary.txt", "remote_server_mount_id": json.Number("2"), "remote_server_id": json.Number("3")},
	"update": {"priority": json.Number("5")},
}

func TestRemoteMountBackendParams_SendFloatsAsNumbersAndDecimalsAsExactText(t *testing.T) {
	for _, test := range []struct {
		name  string
		input decimalInput
		// A string is exact decimal text; a json.Number is a float64 value.
		want map[string]interface{}
	}{
		{"long exact fraction", decimalInput{cpuDecimal: lib.String("1.0049999999999999999999999999")}, map[string]interface{}{"min_free_cpu": "1.0049999999999999999999999999"}},
		{"exponent", decimalInput{memDecimal: lib.String("1e-30")}, map[string]interface{}{"min_free_mem": "1e-30"}},
		{"exact zero is sent", decimalInput{cpuDecimal: lib.String("0"), memDecimal: lib.String("-0.00")}, map[string]interface{}{"min_free_cpu": "0", "min_free_mem": "-0.00"}},
		{"float zero is omitted", decimalInput{cpu: 0, mem: math.Copysign(0, -1)}, map[string]interface{}{}},
		{"float is a number", decimalInput{cpu: 1.005, mem: 12.5}, map[string]interface{}{"min_free_cpu": json.Number("1.005"), "min_free_mem": json.Number("12.5")}},
		{"float and exact fields mix", decimalInput{cpu: 12.5, memDecimal: lib.String("33.333333333333333333")}, map[string]interface{}{"min_free_cpu": json.Number("12.5"), "min_free_mem": "33.333333333333333333"}},
	} {
		for kind, request := range test.input.requests() {
			t.Run(test.name+"/"+kind, func(t *testing.T) {
				// Direct marshaling behaves the same for values and pointers.
				valueJSON, err := json.Marshal(request)
				require.NoError(t, err)
				pointerJSON, err := json.Marshal(pointerTo(request))
				require.NoError(t, err)
				assert.Equal(t, string(valueJSON), string(pointerJSON))
				direct := decodeJSONObject(t, valueJSON)
				for _, key := range []string{"min_free_cpu", "min_free_mem"} {
					assert.LessOrEqual(t, strings.Count(string(valueJSON), `"`+key+`"`), 1, "one %s key", key)
					assert.Equal(t, test.want[key], direct[key], key)
				}
				assert.Equal(t, json.Number("5"), direct["priority"])

				// lib.Params builds the request body and query values.
				want := map[string]interface{}{}
				wantValues := url.Values{}
				for _, fields := range []map[string]interface{}{unrelatedWireFields[kind], test.want} {
					for key, value := range fields {
						want[key] = value
						wantValues.Set(key, stringValue(value))
					}
				}
				body, err := lib.Params{Params: request}.ToJSON()
				require.NoError(t, err)
				bodyJSON, err := io.ReadAll(body)
				require.NoError(t, err)
				assert.Equal(t, want, decodeJSONObject(t, bodyJSON))
				values, err := lib.Params{Params: request}.ToValues()
				require.NoError(t, err)
				assert.Equal(t, wantValues, values)

				assert.Equal(t, test.input.requests()[kind], request, "encoding must not change the request")
			})
		}
	}
}

func TestRemoteMountBackendParams_RejectAmbiguousOrInvalidInputBeforeEncoding(t *testing.T) {
	for _, test := range []struct {
		name, wantErr string
		input         decimalInput
	}{
		{"float and decimal for one field", "MinFreeCpu and MinFreeCpuDecimal are both set", decimalInput{cpu: 1.5, cpuDecimal: lib.String("1.5")}},
		{"empty decimal", `MinFreeMemDecimal must be a decimal number such as "1.5"`, decimalInput{memDecimal: lib.String("")}},
		{"hexadecimal", "MinFreeCpuDecimal must be a decimal number", decimalInput{cpuDecimal: lib.String("0x1p0")}},
		{"underscores", "MinFreeMemDecimal must be a decimal number", decimalInput{memDecimal: lib.String("1_0.5")}},
		{"surrounding space", "MinFreeCpuDecimal must be a decimal number", decimalInput{cpuDecimal: lib.String(" 2.5 ")}},
		{"infinity", "MinFreeCpuDecimal must be a decimal number", decimalInput{cpuDecimal: lib.String("Infinity")}},
		{"non-finite float", "unsupported value: NaN", decimalInput{mem: math.NaN()}},
	} {
		for kind, request := range test.input.requests() {
			t.Run(test.name+"/"+kind, func(t *testing.T) {
				_, valueErr := json.Marshal(request)
				_, pointerErr := json.Marshal(pointerTo(request))
				_, bodyErr := lib.Params{Params: request}.ToJSON()
				_, valuesErr := lib.Params{Params: request}.ToValues()
				for _, err := range []error{valueErr, pointerErr, bodyErr, valuesErr} {
					require.Error(t, err)
					assert.Contains(t, err.Error(), test.wantErr)
					for _, text := range []*string{test.input.cpuDecimal, test.input.memDecimal} {
						if text != nil && *text != "" {
							assert.NotContains(t, err.Error(), *text, "errors must not repeat the input")
						}
					}
				}
			})
		}
	}
}

func TestRemoteMountBackendParams_NilPointerStillEncodesAsNull(t *testing.T) {
	body, err := lib.Params{Params: (*RemoteMountBackendUpdateParams)(nil)}.ToJSON()
	require.NoError(t, err)
	bodyJSON, err := io.ReadAll(body)
	require.NoError(t, err)
	assert.Equal(t, "null", string(bodyJSON))
}

func TestRemoteMountBackendParams_DecodeJSONByTokenType(t *testing.T) {
	exact := "1.0049999999999999999999999999"
	for _, test := range []struct {
		name, data  string
		start, want decimalInput
		wantErr     string
	}{
		{"string sets exact text and clears the float", `{"min_free_cpu":"` + exact + `"}`, decimalInput{cpu: 3}, decimalInput{cpuDecimal: &exact}, ""},
		{"number sets the float and clears exact text", `{"min_free_mem":1.25}`, decimalInput{memDecimal: lib.String("2")}, decimalInput{mem: 1.25}, ""},
		{"null clears exact text and keeps the float", `{"min_free_cpu":null,"min_free_mem":null}`, decimalInput{cpu: 4, memDecimal: lib.String("2")}, decimalInput{cpu: 4}, ""},
		{"missing keys keep their values", `{"priority":5}`, decimalInput{cpuDecimal: lib.String("7"), mem: 2}, decimalInput{cpuDecimal: lib.String("7"), mem: 2}, ""},
		// Keys match case-insensitively and the last occurrence wins, as in plain decoding.
		{"the last case variant decides the route", `{"min_free_cpu":"3","MIN_FREE_CPU":2,"min_free_mem":1,"Min_Free_Mem":"2.50"}`, decimalInput{}, decimalInput{cpu: 2, memDecimal: lib.String("2.50")}, ""},
		{"case variants clear stale exact text", `{"MIN_FREE_CPU":2,"Min_Free_Mem":null}`, decimalInput{cpuDecimal: lib.String("1.5"), memDecimal: lib.String("2")}, decimalInput{cpu: 2}, ""},
		{"other repeated keys keep document order", `{"fall":1,"FALL":2,"min_free_cpu":3}`, decimalInput{}, decimalInput{fall: 2, cpu: 3}, ""},
		{"invalid text", `{"priority":9,"min_free_cpu":"0x1p0"}`, decimalInput{mem: 2}, decimalInput{mem: 2}, "MinFreeCpuDecimal must be a decimal number"},
		{"wrong type", `{"priority":9,"min_free_mem":true}`, decimalInput{mem: 2}, decimalInput{mem: 2}, "cannot unmarshal bool"},
		{"number beyond float64", `{"priority":9,"min_free_cpu":1e400}`, decimalInput{mem: 2}, decimalInput{mem: 2}, "cannot unmarshal number 1e400"},
	} {
		for kind, start := range test.start.requests() {
			t.Run(test.name+"/"+kind, func(t *testing.T) {
				target := pointerTo(start)
				err := json.Unmarshal([]byte(test.data), target)
				// A failed decode leaves the whole request unchanged, including priority.
				assert.Equal(t, test.want.requests()[kind], reflect.ValueOf(target).Elem().Interface())
				if test.wantErr == "" {
					require.NoError(t, err)
					return
				}
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.wantErr)
				if strings.Contains(test.data, "1e400") {
					var typeErr *json.UnmarshalTypeError
					assert.True(t, errors.As(err, &typeErr), "an out-of-range number stays a float64 type error")
				}
			})
		}
	}
}

func TestRemoteMountBackendParams_RoundTripThroughJSON(t *testing.T) {
	input := decimalInput{cpu: 12.5, memDecimal: lib.String("33.333333333333333333")}
	for kind, request := range input.requests() {
		t.Run(kind, func(t *testing.T) {
			encoded, err := json.Marshal(request)
			require.NoError(t, err)
			decoded := reflect.New(reflect.TypeOf(request))
			require.NoError(t, json.Unmarshal(encoded, decoded.Interface()))
			assert.Equal(t, request, decoded.Elem().Interface())
			reencoded, err := json.Marshal(decoded.Interface())
			require.NoError(t, err)
			assert.Equal(t, string(encoded), string(reencoded))
		})
	}
}

func copyText(text *string) *string {
	if text == nil {
		return nil
	}
	return lib.String(*text)
}

func pointerTo(v interface{}) interface{} {
	pointer := reflect.New(reflect.TypeOf(v))
	pointer.Elem().Set(reflect.ValueOf(v))
	return pointer.Interface()
}

func stringValue(value interface{}) string {
	if number, ok := value.(json.Number); ok {
		return number.String()
	}
	return value.(string)
}

// decodeJSONObject keeps numbers as json.Number so assertions see the
// encoded token, never a float64 conversion of it.
func decodeJSONObject(t *testing.T, data []byte) map[string]interface{} {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var object map[string]interface{}
	require.NoError(t, decoder.Decode(&object))
	return object
}
