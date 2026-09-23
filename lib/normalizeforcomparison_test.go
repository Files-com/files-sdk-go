package lib

import (
	"encoding/json"
	"os"
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
		// Component removal happens before Unicode folding.
		{"Ä/../É.txt", "a/e.txt"},
		{"a/．/b", "a/./b"},
		{"a/./B.txt \t", "a/b.txt \t"},
		{"a/. /b", "a/. /b"},
		{". ", ". "},
		{"a/b/c.txt  ", "a/b/c.txt  "},
		{"a/b/c.txt\t", "a/b/c.txt\t"},
		{"a/b/c.txt\n", "a/b/c.txt\n"},
		{"a/b/c.txt\r", "a/b/c.txt\r"},
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

func TestNormalizeForComparisonSharedData(t *testing.T) {
	for _, filename := range []string{"normalization_for_comparison_test_data.json", "comparison_examples.json"} {
		data, err := os.ReadFile("../shared/" + filename)
		if err != nil {
			t.Fatal(err)
		}
		var pairs [][2]string
		if err := json.Unmarshal(data, &pairs); err != nil {
			t.Fatal(err)
		}
		for _, pair := range pairs {
			if actual := NormalizeForComparison(pair[0]); actual != pair[1] {
				t.Errorf("%q: got %q, want %q", pair[0], actual, pair[1])
			}
			if filename == "normalization_for_comparison_test_data.json" {
				if actual := NormalizeForComparison(pair[1]); actual != pair[1] {
					t.Errorf("comparison is not idempotent for %q: %q", pair[1], actual)
				}
			}
		}
	}
}

func BenchmarkNormalizeForComparison(b *testing.B) {
	for _, path := range []string{"folder/subfolder/report.txt", "Folder/Subfolder/Report.txt", "Résumé/q́/カ/𐐀.txt"} {
		b.Run(path, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				NormalizeForComparison(path)
			}
		})
	}
}
