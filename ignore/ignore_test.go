package ignore_test

import (
	"testing"

	"github.com/Files-com/files-sdk-go/v3/ignore"
)

func TestNew(t *testing.T) {
	ig, err := ignore.New()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ig == nil {
		t.Fatal("expected non-nil GitIgnore instance")
	}
}

type testCase struct {
	path    string
	ignored bool
}

// only testing things that are in the common.gitignore to avoid
// CI issues with OS specific gitignore files
func TestNewWithNoOverrides(t *testing.T) {
	ig, err := ignore.New()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ig == nil {
		t.Fatal("expected non-nil GitIgnore instance")
	}

	// *.crdownload, *.part, and *.download are in the common.gitignore file, so
	// they should be ignored by default in this test case.
	testCases := []testCase{
		{"error.log", false},
		{"temp/file.txt", false},
		{"file.txt", false},
		{"xyz.crdownload", true},
		{"xyz.part", true},
		{"xyz.download", true},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			matched := ig.MatchesPath(tc.path)
			if matched != tc.ignored {
				t.Errorf("expected MatchesPath(%q) to be %v, got %v", tc.path, tc.ignored, matched)
			}
		})
	}
}

func TestNewWithOverrides(t *testing.T) {
	overrides := []string{"*.log", "temp/"}
	ig, err := ignore.New(overrides...)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ig == nil {
		t.Fatal("expected non-nil GitIgnore instance")
	}

	// only things ending in .log or starting with temp/ should be ignored in this test case.
	testCases := []testCase{
		{"error.log", true},
		{"temp/file.txt", true},
		{"file.txt", false},
		{"xyz.crdownload", false},
		{"xyz.part", false},
		{"xyz.download", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			matched := ig.MatchesPath(tc.path)
			if matched != tc.ignored {
				t.Errorf("expected MatchesPath(%q) to be %v, got %v", tc.path, tc.ignored, matched)
			}
		})
	}
}

func TestNewWithAllowed(t *testing.T) {
	allowed := []string{"*.crdownload", "temp/"}
	ig, err := ignore.NewWithAllowList(allowed...)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ig == nil {
		t.Fatal("expected non-nil GitIgnore instance")
	}

	// *.crdownload, *.part, and *.download are in the common.gitignore file, but *.crdownload and
	// temp/ are in the allowed list, so only *.part and *.download should be ignored in this test case.
	testCases := []testCase{
		{"error.log", false},
		{"temp/file.txt", false},
		{"file.txt", false},
		{"xyz.crdownload", false},
		{"xyz.part", true},
		{"xyz.download", true},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			matched := ig.MatchesPath(tc.path)
			if matched != tc.ignored {
				t.Errorf("expected MatchesPath(%q) to be %v, got %v", tc.path, tc.ignored, matched)
			}
		})
	}
}
