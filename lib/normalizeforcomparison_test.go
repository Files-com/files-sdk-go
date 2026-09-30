package lib

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"
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

func TestNormalizeForComparisonAppliesSharedMapping(t *testing.T) {
	data, err := os.ReadFile("../shared/path_comparison.json")
	if err != nil {
		t.Fatal(err)
	}
	var comparison struct {
		Mapping map[string]string `json:"mapping"`
	}
	if err := json.Unmarshal(data, &comparison); err != nil {
		t.Fatal(err)
	}
	mapping := make(map[rune]string, len(comparison.Mapping))
	for hex, replacement := range comparison.Mapping {
		scalar, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			t.Fatal(err)
		}
		mapping[rune(scalar)] = replacement
	}

	// Printable ASCII is lowercased directly. Every other scalar is replaced
	// by its shared mapping entry, or kept when it has none.
	for scalar := rune(0); scalar <= utf8.MaxRune; scalar++ {
		if !utf8.ValidRune(scalar) || (scalar >= ' ' && scalar <= '~') {
			continue
		}
		want, ok := mapping[scalar]
		if !ok {
			want = string(scalar)
		}
		if actual := NormalizeForComparison(string(scalar)); actual != want {
			t.Errorf("U+%04X: got %q, want %q", scalar, actual, want)
		}
	}
}

// inFreshProcess reruns the calling test alone in a new process, where the
// comparison table has not been loaded yet, and fails if that run fails. It
// returns true inside the new process.
func inFreshProcess(t *testing.T) bool {
	t.Helper()
	const env = "NORMALIZE_FOR_COMPARISON_FRESH_PROCESS"
	if os.Getenv(env) == t.Name() {
		return true
	}
	cmd := exec.Command(os.Args[0], "-test.run=^"+t.Name()+"$", "-test.count=1", "-test.v")
	cmd.Env = append(os.Environ(), env+"="+t.Name())
	out, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "--- PASS: "+t.Name()) {
		t.Fatalf("fresh process run failed: %v\n%s", err, out)
	}
	return false
}

func panicValue(f func()) (recovered any) {
	defer func() { recovered = recover() }()
	f()
	return nil
}

func TestNormalizeForComparisonLoadsMappingOnFirstLookup(t *testing.T) {
	if !inFreshProcess(t) {
		return
	}
	// Truncating the embedded table after startup is only noticed once a
	// lookup outside printable ASCII needs it, and then on every lookup.
	pathComparisonGzip = pathComparisonGzip[:len(pathComparisonGzip)/2]

	if actual := NormalizeForComparison("Folder\\Report.TXT"); actual != "folder/report.txt" {
		t.Errorf("got %q, want %q", actual, "folder/report.txt")
	}
	for range 2 {
		if panicValue(func() { NormalizeForComparison("Résumé.txt") }) == nil {
			t.Error("lookup with a corrupt table did not panic")
		}
	}
}

func TestNormalizeForComparisonConcurrentFirstLookups(t *testing.T) {
	if !inFreshProcess(t) {
		return
	}
	start := make(chan struct{})
	results := make([]string, 32)
	var wg sync.WaitGroup
	for i := range results {
		wg.Go(func() {
			<-start
			results[i] = NormalizeForComparison("FÎŁĘÑÂMÉ.TXT")
		})
	}
	close(start)
	wg.Wait()

	for i, actual := range results {
		if actual != "filename.txt" {
			t.Errorf("caller %d: got %q, want %q", i, actual, "filename.txt")
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
