package fsmount

import (
	"testing"
)

func TestAdditionalIgnorePatternsMatchTmpFiles(t *testing.T) {
	ig, err := ignoreFromPatterns(nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	fs := &Filescomfs{ignore: ig}

	cases := []struct {
		path           string
		storedRemotely bool
	}{
		// Illustrator-style temp file with various casings should NOT be stored remotely
		{path: "/Illustrator/A9Rdwrczv_1b81nor_aj0.tmp", storedRemotely: false},
		{path: "/Illustrator/A9Rdwrczv_1b81nor_aj0.TMP", storedRemotely: false},
		{path: "/Illustrator/A9Rdwrczv_1b81nor_aj0.Tmp", storedRemotely: false},
		{path: "/Illustrator/A9Rdwrczv_1b81nor_aj0.tMp", storedRemotely: false},
		{path: "/Illustrator/A9Rdwrczv_1b81nor_aj0.tmP", storedRemotely: false},
		{path: "/Illustrator/a9rdwrczv_1b81nor_aj0.tmp", storedRemotely: false},
		{path: "/Illustrator/A9RDWRCZV_1B81NOR_AJ0.TMP", storedRemotely: false},
		{path: "/illustrator/a9rdwrczv_1b81nor_aj0.tmp", storedRemotely: false},

		// bare filename without directory
		{path: "A9Rdwrczv_1b81nor_aj0.tmp", storedRemotely: false},

		// deeply nested path
		{path: "/some/deep/path/Illustrator/A9Rdwrczv_1b81nor_aj0.tmp", storedRemotely: false},

		// other temp file patterns that should also NOT be stored remotely
		{path: "/Illustrator/A9R2gnwq7_ax8nkp_1944.tmp", storedRemotely: false},
		{path: "/Office/~WR1234.tmp", storedRemotely: false},
		{path: "/Photoshop/psAF90.tmp", storedRemotely: false},
		{path: "/AutoCAD/save686566b0.tmp", storedRemotely: false},

		// non-tmp extensions SHOULD be stored remotely
		{path: "/Illustrator/A9Rdwrczv_1b81nor_aj0.txt", storedRemotely: true},
		{path: "/Illustrator/A9Rdwrczv_1b81nor_aj0.psd", storedRemotely: true},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got := fs.isStoredRemotely(tc.path)
			if got != tc.storedRemotely {
				t.Errorf("expected isStoredRemotely(%q) to be %v, got %v", tc.path, tc.storedRemotely, got)
			}
		})
	}
}
