package ignore_test

import (
	"testing"

	"github.com/Files-com/files-sdk-go/v3/ignore"
)

// only testing things in the common.gitignore to avoid CI issues with OS specific gitignore files

type testCase struct {
	path    string
	ignored bool
}

func TestNew(t *testing.T) {
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
		{path: "error.log", ignored: false},
		{path: "temp/file.txt", ignored: false},
		{path: "file.txt", ignored: false},
		{path: "xyz.crdownload", ignored: true},
		{path: "xyz.part", ignored: true},
		{path: "xyz.download", ignored: true},
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

func TestIgnoreNothing(t *testing.T) {
	// pass an empty slice to ignore nothing
	ig, err := ignore.New([]string{}...)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ig == nil {
		t.Fatal("expected non-nil GitIgnore instance")
	}

	// nothing should be ignored in this test case.
	testCases := []testCase{
		{path: "error.log", ignored: false},
		{path: "temp/file.txt", ignored: false},
		{path: "file.txt", ignored: false},
		{path: "xyz.crdownload", ignored: false},
		{path: "xyz.part", ignored: false},
		{path: "xyz.download", ignored: false},
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
		{path: "error.log", ignored: true},
		{path: "temp/file.txt", ignored: true},
		{path: "file.txt", ignored: false},
		{path: "xyz.crdownload", ignored: false},
		{path: "xyz.part", ignored: false},
		{path: "xyz.download", ignored: false},
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
		{path: "error.log", ignored: false},
		{path: "temp/file.txt", ignored: false},
		{path: "file.txt", ignored: false},
		{path: "xyz.crdownload", ignored: false},
		{path: "xyz.part", ignored: true},
		{path: "xyz.download", ignored: true},
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

func TestNewWithDenied(t *testing.T) {
	denied := []string{
		"*.no-such-extension",
		"tests/",
		".~*",
		"~*",
	}
	ig, err := ignore.NewWithDenyList(denied...)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ig == nil {
		t.Fatal("expected non-nil GitIgnore instance")
	}

	// anything in the common gitignore file, plus the following:
	// file.no-such-extension,
	// tests/file.txt,
	// .~WRD0000,
	// ~WRL0001
	// should be ignored in this test case.
	testCases := []testCase{
		{path: "error.log", ignored: false},
		{path: "temp/file.txt", ignored: false},
		{path: "file.txt", ignored: false},
		{path: "xyz.crdownload", ignored: true},
		{path: "xyz.part", ignored: true},
		{path: "xyz.download", ignored: true},
		{path: "file.no-such-extension", ignored: true},
		{path: "tests/file.txt", ignored: true},
		{path: "/Skills.docx.sb-c76ae02f-s1IOeW/.~WRD0000", ignored: true},
		{path: "/Skills-1.docx.sb-c76ae02f-f69gHw/~WRL0001", ignored: true},
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
