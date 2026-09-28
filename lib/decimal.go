package lib

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/appscode/go-querystring/query"
)

// decimalText matches finite base-10 numbers. Go's regexp runs in linear
// time, so long input costs no more than reading it.
var decimalText = regexp.MustCompile(`^[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?$`)

// IsDecimalText reports whether text is a finite base-10 number, such as
// "1.5", "-0.25", ".5", "01.50" or "2e-3". It rejects empty text, spaces,
// NaN, infinities, hexadecimal, fractions and underscores. It checks syntax
// only: it sets no range or precision limit, and the API decides which values
// it accepts.
func IsDecimalText(text string) bool {
	return decimalText.MatchString(text)
}

// DecimalOverride pairs a float64 request field, released before the API
// accepted exact decimals for it, with the *string field that sends exact
// decimal text under the same key instead. Generated request types keep the
// float64 field for existing callers and use these helpers to encode and
// decode either.
type DecimalOverride struct {
	Key     string // JSON and query key shared by both fields
	Field   string // Go name of the float64 field; the text field is Field + "Decimal"
	Float   *float64
	Decimal **string
}

// text returns the decimal text to send in place of the float64 value, if
// the caller set one.
func (o DecimalOverride) text() (text string, ok bool, err error) {
	decimal := *o.Decimal
	if decimal == nil {
		return "", false, nil
	}
	// A zero float64 cannot be told apart from an unset one, so only a
	// nonzero value (including NaN) conflicts with the decimal text.
	if *o.Float != 0 {
		return "", false, fmt.Errorf("%s and %sDecimal are both set; leave %s zero to send exact decimal text", o.Field, o.Field, o.Field)
	}
	if !IsDecimalText(*decimal) {
		return "", false, o.invalidTextError()
	}
	return *decimal, true, nil
}

// The message leaves out the rejected text, which may be arbitrary input.
func (o DecimalOverride) invalidTextError() error {
	return fmt.Errorf(`%sDecimal must be a decimal number such as "1.5" or "2e-3"`, o.Field)
}

func decimalTexts(overrides []DecimalOverride) (map[string]string, error) {
	texts := make(map[string]string)
	for _, override := range overrides {
		text, ok, err := override.text()
		if err != nil {
			return nil, err
		}
		if ok {
			texts[override.Key] = text
		}
	}
	return texts, nil
}

// MarshalDecimalOverrides encodes v, a request struct type without its own
// MarshalJSON method, as JSON, sending each decimal text as a JSON string.
// Without decimal text the result is exactly json.Marshal(v).
func MarshalDecimalOverrides(v any, overrides ...DecimalOverride) ([]byte, error) {
	texts, err := decimalTexts(overrides)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(v)
	if err != nil || len(texts) == 0 {
		return data, err
	}

	// Raw values keep every other field's encoded bytes unchanged.
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, err
	}
	for key, text := range texts {
		if object[key], err = json.Marshal(text); err != nil {
			return nil, err
		}
	}
	return json.Marshal(object)
}

// DecimalOverrideValues encodes v like query.Values, sending each decimal
// text as the single value for its key.
func DecimalOverrideValues(v any, overrides ...DecimalOverride) (url.Values, error) {
	texts, err := decimalTexts(overrides)
	if err != nil {
		return nil, err
	}
	values, err := query.Values(v)
	if err != nil {
		return nil, err
	}
	for key, text := range texts {
		values.Set(key, text)
	}
	return values, nil
}

// UnmarshalJSON makes o the decoding target for its key. Generated request
// types decode each override in the same encoding/json pass as their other
// fields, so key matching and repeated keys behave as in plain decoding, and
// the last occurrence of the key decides the route. A JSON string sets the
// decimal text and clears the float64; null clears the text and, as in plain
// decoding, leaves the float64 alone; any other value clears the text and
// decodes into the float64 as plain decoding would.
func (o *DecimalOverride) UnmarshalJSON(data []byte) error {
	switch data[0] {
	case 'n':
		*o.Decimal = nil
		return nil
	case '"':
		var text string
		if err := json.Unmarshal(data, &text); err != nil {
			return err
		}
		if !IsDecimalText(text) {
			return o.invalidTextError()
		}
		*o.Decimal, *o.Float = &text, 0
		return nil
	}
	*o.Decimal = nil
	if err := json.Unmarshal(data, o.Float); err != nil {
		return fmt.Errorf("%s: %w", o.Key, err)
	}
	return nil
}
