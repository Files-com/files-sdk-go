package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stageWithContent writes content to the temporary download of name inside dir,
// as an interrupted download of that file would have left it.
func stageWithContent(t *testing.T, dir string, name string, content string) destinationPath {
	t.Helper()
	tmp, err := tmpDownloadPath(explicitDestination(filepath.Join(dir, name)), "")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(tmp.String(), []byte(content), 0600))
	return tmp
}

// resumedContent is what a later run of a download of name would continue from.
func resumedContent(t *testing.T, dir string, name string) (string, bool) {
	t.Helper()
	tmp, ok := existingTmpDownloadPath(explicitDestination(filepath.Join(dir, name)), "")
	if !ok {
		return "", false
	}
	content, err := os.ReadFile(tmp.String())
	require.NoError(t, err)
	return string(content), true
}

// The reported defect: a file named after another file's temporary download must
// not be mistaken for that file's unfinished download. This also covers
// unfinished downloads left by a client from before this naming, which are
// equally indistinguishable from a planted file.
func Test_existingTmpDownloadPath_doesNotAdoptASiblingNamedLikeATemporaryDownload(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "flat.csv.download"), []byte("not mine"), 0600))
	folder := filepath.Join(dir, "folder.csv.download")
	require.NoError(t, os.Mkdir(folder, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(folder, "folder.csv"), []byte("not mine either"), 0600))

	for _, name := range []string{"flat.csv", "folder.csv"} {
		_, ok := existingTmpDownloadPath(explicitDestination(filepath.Join(dir, name)), "")
		assert.Falsef(t, ok, "%v must not resume from a sibling it does not own", name)
	}
}

// No file name may reach another file name's temporary download. A shortened
// temporary name is what invites this: the second name of each pair below is
// spelled the way the first one's temporary name is shortened.
func Test_tmpDownloadPath_keepsEachFileNameToItsOwnTemporaryDownload(t *testing.T) {
	longName := strings.Repeat("x", 244)
	for _, pair := range []struct {
		name  string
		first string
		other string
	}{
		{
			name:  "a name spelled as another name shortened and digested",
			first: longName,
			other: strings.Repeat("x", 221) + "-7d19aed0baf5",
		},
		{
			name:  "a name spelled as the readable part of another file's temporary name",
			first: longName,
			other: strings.TrimSuffix(strings.TrimPrefix(tmpDownloadElement(longName, ""), tempDownloadEncodedPrefix), "."+TempDownloadExtension),
		},
	} {
		t.Run(pair.name, func(t *testing.T) {
			dir := t.TempDir()
			stageWithContent(t, dir, pair.first, "bytes of the first file")
			stageWithContent(t, dir, pair.other, "bytes of the other file")

			first, ok := resumedContent(t, dir, pair.first)
			require.True(t, ok)
			other, ok := resumedContent(t, dir, pair.other)
			require.True(t, ok)

			assert.Equal(t, "bytes of the first file", first)
			assert.Equal(t, "bytes of the other file", other)
		})
	}
}

// Every attempt returns a path nothing else is using, and an alternate one does
// not take over the temporary download of a file that is really named that way.
func Test_tmpDownloadPath_returnsAnUnusedPathAndLeavesOtherFilesAlone(t *testing.T) {
	dir := t.TempDir()
	taken := stageWithContent(t, dir, "notes.txt (1)", "bytes of the file really named that")

	used := map[string]bool{taken.String(): true}
	for attempt := 0; attempt < 12; attempt++ {
		path, err := tmpDownloadPath(explicitDestination(filepath.Join(dir, "notes.txt")), "")
		require.NoError(t, err)
		require.Falsef(t, used[path.String()], "attempt %v reused %v", attempt, path)
		used[path.String()] = true
		require.NoError(t, os.WriteFile(path.String(), []byte("partial"), 0600))
	}

	content, ok := resumedContent(t, dir, "notes.txt (1)")
	require.True(t, ok)
	assert.Equal(t, "bytes of the file really named that", content)
}

// Long names, including ones whose characters are several bytes each, stay
// within what the filesystem accepts and still resume.
func Test_tmpDownloadPath_staysWithinThePathElementLimit(t *testing.T) {
	for _, name := range []string{
		strings.Repeat("a", 234),          // the longest name kept as it is
		strings.Repeat("a", 235),          // the shortest name that is encoded
		strings.Repeat("a", 240) + ".txt", // a name that worked before this change
		strings.Repeat("日", 70) + ".txt",  // three bytes per character, kept as it is
		strings.Repeat("日", 83) + ".txt",  // three bytes per character, shortened mid-character
		strings.Repeat("é", 121) + ".txt", // two bytes per character
	} {
		dir := t.TempDir()
		staged := stageWithContent(t, dir, name, "partial "+name)

		for _, element := range strings.Split(staged.name, string(filepath.Separator)) {
			assert.LessOrEqualf(t, len(element), maxTempDownloadElementBytes, "%q is too long for a path element", element)
		}
		content, ok := resumedContent(t, dir, name)
		require.True(t, ok)
		assert.Equal(t, "partial "+name, content)
	}
}

// Names this client generates never transfer; names that only look similar do.
func Test_isReservedTempDownloadPath(t *testing.T) {
	for path, reserved := range map[string]bool{
		".~files-cli.report.csv.download":                    true,
		"folder/.~files-cli.report.csv.download":             true,
		".~files-cli.report.csv.download/report.csv":         true, // the macOS folder's contents
		".~FILES-CLI.report.csv.download/report.csv":         true, // reaches the same path where case does not matter
		".~files-cli.report.csv.DOWNLOAD":                    true,
		".~file\u017f-cli.report.csv.download":               true, // macOS resolves \u017f and s to one name
		".~file\u017f-cli.report.csv.download/report.csv":    true,
		".~FILE\u017f-CLI.report.csv.DOWNLOAD":               true,
		".~files-cli~1.0123456789abcdef.report.csv.download": true,
		"folder/.~files-cli~.0123456789abcdef.x.download":    true,
		"report.csv":                       false,
		"report.csv.download":              false,
		".~files-cli-notes.download":       false,
		".~files-clients.download":         false,
		".~files-cli.report.csv":           false,
		"notes/.~files-cli-notes.download": false,
		".~file\u017f-cli-notes.download":  false,
		".~file\u017f-clients.download":    false,
	} {
		assert.Equalf(t, reserved, isReservedTempDownloadPath(path), "%q", path)
	}
}

// On a filesystem that resolves different spellings of a name to one file, the
// alias must be excluded too: it physically reaches the same temporary path.
// "\u017f" folds to "s", which macOS resolves but a plain lower-casing does not.
func Test_isReservedTempDownloadPath_coversAnAliasTheFilesystemResolves(t *testing.T) {
	dir := t.TempDir()
	alias := filepath.Join(dir, ".~file\u017f-cli.report.csv.download")
	require.NoError(t, os.Mkdir(alias, 0700))
	aliasInfo, err := os.Stat(alias)
	require.NoError(t, err)

	canonicalInfo, err := os.Stat(filepath.Join(dir, ".~files-cli.report.csv.download"))
	if os.IsNotExist(err) {
		t.Skip("this filesystem keeps the two spellings apart")
	}
	require.NoError(t, err)
	require.True(t, os.SameFile(aliasInfo, canonicalInfo), "the two spellings are expected to be one file here")

	assert.True(t, isReservedTempDownloadPath(alias))
}
