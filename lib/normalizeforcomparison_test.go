package lib

import (
	"testing"
)

func TestNormalizeForComparison(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"filename.txt", "filename.txt"},
		{"FiLeNaMe.TxT", "filename.txt"},
		{"FILENAME.TXT", "filename.txt"},
		{"FÎŁĘÑÂMÉ.TXT", "filename.txt"},
		{"Fïłèńämê.Txt", "filename.txt"},
		{"a/b/c.txt", "a/b/c.txt"},
		{"A\\B\\C.TXT", "a/b/c.txt"},
		{"A/B\\C.TXT", "a/b/c.txt"},
		{"//a/b//c.txt", "a/b/c.txt"},
		{"a/b/c.txt  ", "a/b/c.txt"},
		{"a/b/c.txt\t", "a/b/c.txt"},
		{"a/b/c.txt\n", "a/b/c.txt"},
		{"a/b/c.txt\r", "a/b/c.txt"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			output := normalizeForComparison(tc.input)
			if output != tc.expected {
				t.Errorf("Expected %s but got %s", tc.expected, output)
			}
		})
	}
}
