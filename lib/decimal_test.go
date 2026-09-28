package lib

import (
	"strings"
	"testing"
)

func TestIsDecimalText(t *testing.T) {
	accepted := []string{
		"0", "-0", "+0.00", "1.5", "+1.5", "-1.5", ".5", "1.", "01.50", "1.e3", "1e-30", "2.5E+3",
		"1e9999", // syntax only: exact text has no float64 range limit
		"1.0049999999999999999999999999",
		strings.Repeat("9", 4096) + "." + strings.Repeat("1", 4096),
	}
	for _, text := range accepted {
		if !IsDecimalText(text) {
			t.Errorf("IsDecimalText(%.40q) = false, want true", text)
		}
	}

	rejected := []string{
		"", " ", " 1.5", "1.5 ", "1.5\n", ".", "+", "-", "+-1", "1..5", "1.5.2",
		"NaN", "nan", "Inf", "+Infinity", "-Infinity",
		"0x1p0", "0X1.8P+2", "1_0.5", "1/2", "1e", "1e+", "e3", "1e3.5", "1.5x", "١",
	}
	for _, text := range rejected {
		if IsDecimalText(text) {
			t.Errorf("IsDecimalText(%q) = true, want false", text)
		}
	}
}
