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
		{"/../../remote\\path//./to/file.txt", "remote/path/to/file.txt"},
		{"remote/../path/to/file.txt", "remote/path/to/file.txt"},
		// Exact "." and ".." components are dropped wherever they appear, without
		// collapsing the component before "..".
		{"./a/b.txt", "a/b.txt"},
		{"a/./b.txt", "a/b.txt"},
		{"a/b/.", "a/b"},
		{"../a/b.txt", "a/b.txt"},
		{"a/../b.txt", "a/b.txt"},
		{"a/b/..", "a/b"},
		{"a/b/c/../../hello", "a/b/c/hello"},
		{".", ""},
		{"..", ""},
		{"./", ""},
		{"/../.", ""},
		{".\\", ""},
		// Names that merely contain dots are preserved.
		{".hidden", ".hidden"},
		{"...", "..."},
		{"a..b", "a..b"},
		{"file.", "file."},
		{".git/config", ".git/config"},
		// Null removal and backslash conversion happen before component removal.
		{"a\\.\x00\\b.txt", "a/b.txt"},
		{"a/.\x00./b", "a/b"},
		// Component removal happens before Unicode folding and whitespace trimming.
		{"Ä/../É.txt", "a/e.txt"},
		{"a/．/b", "a/./b"},
		{"a/./B.txt \t", "a/b.txt"},
		{"a/. /b", "a/. /b"},
		{". ", "."},
		{"a/b/c.txt  ", "a/b/c.txt"},
		{"a/b/c.txt\t", "a/b/c.txt"},
		{"a/b/c.txt\n", "a/b/c.txt"},
		{"a/b/c.txt\r", "a/b/c.txt"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			output := NormalizeForComparison(tc.input)
			if output != tc.expected {
				t.Errorf("Expected %s but got %s", tc.expected, output)
			}
		})
	}
}
