package file

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const TempDownloadExtension = "download"

// Temporary downloads live in a namespace reserved for this client: a path
// element that starts with tempDownloadNamespace and ends with the temporary
// extension is never uploaded or downloaded, in either direction, whatever the
// caller's ignore and include rules say. Nothing else may occupy a temporary
// path, so an unfinished download can only ever be continued by the file it
// belongs to.
const (
	tempDownloadNamespace = ".~files-cli"

	// tempDownloadPrefix starts the temporary name of a file whose name fits
	// within a path element: ".~files-cli.<name>.download".
	tempDownloadPrefix = tempDownloadNamespace + "."

	// tempDownloadEncodedPrefix starts the temporary name of a file that needs
	// its name encoded: ".~files-cli~<uniqueness>.<digest>.<shortened>.download".
	// The character that tells the two apart is part of this prefix, so a file
	// name, which only ever follows it, cannot spell the other form.
	tempDownloadEncodedPrefix = tempDownloadNamespace + "~"

	// maxTempDownloadElementBytes is the path element length the supported
	// filesystems allow. A generated temporary name stays within it so every
	// file name that could be written can also be staged.
	maxTempDownloadElementBytes = 255

	// tempDownloadDigestBytes of the file name's SHA-256 identify it inside an
	// encoded name. The digest is a fixed width in a fixed position, so two
	// different names cannot produce one encoded name.
	tempDownloadDigestBytes = 16
)

// tmpDownloadBase is the path temporary file names derive from: the final file
// itself, or the final file's name directly inside an external temp directory.
func tmpDownloadBase(final destinationPath, tempPath string) destinationPath {
	if tempPath != "" {
		return destinationPath{dir: tempPath, name: final.base()}
	}
	return final
}

// existingTmpDownloadPath returns the canonical temporary download path for
// final when it already exists.
func existingTmpDownloadPath(final destinationPath, tempPath string) (destinationPath, bool) {
	return existingTmpDownloadFile(final, tmpDownloadCandidate(tmpDownloadBase(final, tempPath), ""))
}

// tmpDownloadPath generates a unique temporary download path for final inside
// the reserved namespace and, if necessary, with additional identifiers to
// avoid name conflicts.
func tmpDownloadPath(final destinationPath, tempPath string) (destinationPath, error) {
	var index int
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	base := tmpDownloadBase(final, tempPath)

	for {
		var uniqueness string
		if index > 25 {
			return destinationPath{}, fmt.Errorf("unable to create a unique temporary path after 25 attempts, consider deleting existing %v*.%v files; attempted path: %v", tempDownloadNamespace, TempDownloadExtension, base)
		} else if index > 10 {
			for i := 0; i < 4; i++ {
				uniqueness += string(rune(randGenerator.Intn(26) + 'a'))
			}
		} else if index > 0 {
			uniqueness = fmt.Sprintf("%v", index)
		}
		candidate := tmpDownloadCandidate(base, uniqueness)

		if _, err := candidate.lstat(); errors.Is(err, fs.ErrNotExist) {
			return tmpDownloadPathOnNotExist(final, candidate)
		}
		index++
	}
}

// tmpDownloadCandidate is the temporary path beside base for one uniqueness
// token. Only the last element of the name becomes temporary, so a file staged
// for a folder download stays next to where it is going.
func tmpDownloadCandidate(base destinationPath, uniqueness string) destinationPath {
	dir, name := filepath.Split(base.name)
	return base.withName(filepath.Join(dir, tmpDownloadElement(name, uniqueness)))
}

// tmpDownloadElement names the temporary file or folder of one file name.
//
// A name that fits is kept as it is, so an unfinished download is recognizable
// and the same name always produces the same temporary path. A name that is too
// long, or that needs a uniqueness token, is encoded after a different prefix
// character and identified by a digest of the whole name. The shortened name in
// an encoded element is only there to be read: the digest, not the shortened
// name, decides which file the element belongs to, so no other file name can
// take over an encoded temporary path.
func tmpDownloadElement(name string, uniqueness string) string {
	plain := tempDownloadPrefix + name + "." + TempDownloadExtension
	if uniqueness == "" && len(plain) <= maxTempDownloadElementBytes {
		return plain
	}

	digest := sha256.Sum256([]byte(name))
	encoded := tempDownloadEncodedPrefix + uniqueness + "." + hex.EncodeToString(digest[:tempDownloadDigestBytes]) + "."
	shortened := truncateOnRuneBoundary(name, maxTempDownloadElementBytes-len(encoded)-len(TempDownloadExtension)-1)
	return encoded + shortened + "." + TempDownloadExtension
}

// truncateOnRuneBoundary shortens value to at most maxBytes without splitting a
// character.
func truncateOnRuneBoundary(value string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len(value) <= maxBytes {
		return value
	}
	for maxBytes > 0 && !utf8.RuneStart(value[maxBytes]) {
		maxBytes--
	}
	return value[:maxBytes]
}

// isReservedTempDownloadPath reports whether path is, or is inside, a temporary
// download. Every element is checked, because a macOS temporary download is a
// folder holding the file being downloaded.
func isReservedTempDownloadPath(path string) bool {
	for _, element := range strings.FieldsFunc(path, isPathSeparator) {
		if isReservedTempDownloadElement(element) {
			return true
		}
	}
	return false
}

// isReservedTempDownloadElement reports whether one path element is a name this
// client generates for a temporary download. Only the two prefixes it actually
// generates are reserved, so ordinary names that merely start the same way,
// such as ".~files-cli-notes.download", stay transferable.
func isReservedTempDownloadElement(element string) bool {
	if !hasSuffixFold(element, "."+TempDownloadExtension) {
		return false
	}
	return hasPrefixFold(element, tempDownloadPrefix) || hasPrefixFold(element, tempDownloadEncodedPrefix)
}

// hasPrefixFold reports whether s begins with prefix, counting characters a
// case-insensitive filesystem resolves to the same name as the same character.
// That includes ones written with a different number of bytes, such as "ſ" for
// "s", which macOS resolves to one file, so the comparison walks characters
// rather than bytes.
func hasPrefixFold(s string, prefix string) bool {
	for _, want := range prefix {
		got, size := utf8.DecodeRuneInString(s)
		if size == 0 || !runeEqualFold(got, want) {
			return false
		}
		s = s[size:]
	}
	return true
}

// hasSuffixFold is hasPrefixFold from the end of s.
func hasSuffixFold(s string, suffix string) bool {
	for len(suffix) > 0 {
		want, wantSize := utf8.DecodeLastRuneInString(suffix)
		got, gotSize := utf8.DecodeLastRuneInString(s)
		if gotSize == 0 || !runeEqualFold(got, want) {
			return false
		}
		s, suffix = s[:len(s)-gotSize], suffix[:len(suffix)-wantSize]
	}
	return true
}

// runeEqualFold is what strings.EqualFold compares, for one character.
func runeEqualFold(a rune, b rune) bool {
	if a == b {
		return true
	}
	for folded := unicode.SimpleFold(a); folded != a; folded = unicode.SimpleFold(folded) {
		if folded == b {
			return true
		}
	}
	return false
}

// isPathSeparator accepts both separators, because the paths checked come from
// the local filesystem and from server responses.
func isPathSeparator(r rune) bool {
	return r == '/' || r == filepath.Separator
}

// checkpointedTmpDownload returns the temporary file recorded by a prior
// paused session when it still exists. A recorded path below the destination
// directory or the temp directory is resolved inside that directory, so the
// server-derived part of it stays confined. Any other path is the caller's
// explicit choice and is used as given, so its partial data is kept.
func (d *DownloadStatus) checkpointedTmpDownload() (destinationPath, bool) {
	if d.TmpPath == "" {
		return destinationPath{}, false
	}
	tmp := explicitTmpDownload(d.TmpPath)
	for _, dir := range []string{d.destination.dir, d.tempPath} {
		if dir == "" {
			continue
		}
		if rel, err := filepath.Rel(dir, d.TmpPath); err == nil && rel != "." && filepath.IsLocal(rel) {
			tmp = destinationPath{dir: dir, name: rel}
			break
		}
	}
	if _, err := tmp.stat(); err != nil {
		return destinationPath{}, false
	}
	return tmp, true
}

// tmpDownloadToResume returns the temporary file a download continues from:
// the checkpointed one, otherwise the canonical one when it exists.
func (d *DownloadStatus) tmpDownloadToResume() (destinationPath, bool) {
	if tmp, ok := d.checkpointedTmpDownload(); ok {
		return tmp, true
	}
	return existingTmpDownloadPath(d.destination, d.tempPath)
}
