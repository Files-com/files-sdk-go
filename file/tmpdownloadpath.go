package file

import (
	"context"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
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
//
// A temporary download has two names. While a transfer may still write to it,
// it is under its in-flight name (tmpDownloadPath). When a transfer stops in an
// orderly way and has ended the file at the bytes it received, the file is
// renamed to its paused name (tempDownloadPausedToken), and only a file under
// that name is continued from, or published without asking for content. The
// rename is the proof of what the file holds: parts are written at their own
// offsets, so a file still under its in-flight name may have been left by a
// process that stopped mid-write, or by an earlier version of this client, and
// its length says nothing about how much of it was received. Such a file is
// discarded and the download starts over.
const (
	tempDownloadNamespace = ".~files-cli"

	// tempDownloadPausedToken starts the token of a paused temporary download's
	// encoded name, ahead of the uniqueness token the file staged under.
	// Uniqueness tokens are empty, digits or letters, so the "-" keeps a paused
	// token from ever spelling one.
	tempDownloadPausedToken = "p-"

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

	// tempDownloadDigestBytes of a SHA-256 identify a file inside an encoded
	// name. The digest is a fixed width in a fixed position, so two different
	// inputs cannot produce one encoded name.
	tempDownloadDigestBytes = 16

	// tempDownloadPathDigestTag marks a digest of a file's whole local path, as
	// opposed to a digest of its name alone. It begins with a character that is
	// not a hexadecimal digit, so an encoded name identified by a path is never
	// spelled like one identified by a name: an external stage from before path
	// identities were introduced cannot become the canonical stage of any file.
	tempDownloadPathDigestTag = "path-"
)

// tmpDownloadNaming is what the temporary names of one file are made from.
type tmpDownloadNaming struct {
	// base is the path a temporary name is formed beside: the final file
	// itself, or a stand-in carrying the final file's name directly inside an
	// external temp directory.
	base destinationPath

	// pathDigest identifies the final file by its whole local path. It is set
	// only inside an external temp directory: files from every folder, and
	// from every destination the directory is shared with, are staged side by
	// side there, so a name alone does not say which file a temporary download
	// belongs to. Beside the final file, the directory it is in says that.
	pathDigest string
}

// tmpDownloadNamingFor derives the temporary naming of final: beside final
// itself, or inside the external temp directory tempPath when it is set.
func tmpDownloadNamingFor(final destinationPath, tempPath string) (tmpDownloadNaming, error) {
	if tempPath == "" {
		return tmpDownloadNaming{base: final}, nil
	}
	digest, err := tmpDownloadPathDigest(final)
	if err != nil {
		return tmpDownloadNaming{}, err
	}
	return tmpDownloadNaming{base: destinationPath{dir: tempPath, name: final.base()}, pathDigest: digest}, nil
}

// tmpDownloadPathDigest identifies final by its absolute, cleaned local path,
// so one file reached through a relative path, a redundant separator, or a
// different split between directory and name has one temporary download, and
// two files whose paths differ anywhere have two. It fails only when the
// working directory a relative path depends on cannot be determined; the
// temporary download is then not named at all, rather than named in a way
// another file could share.
func tmpDownloadPathDigest(final destinationPath) (string, error) {
	absolute, err := filepath.Abs(final.String())
	if err != nil {
		return "", fmt.Errorf("download: unable to identify the temporary file of %v: %w", final, err)
	}
	digest := sha256.Sum256([]byte(absolute))
	return tempDownloadPathDigestTag + hex.EncodeToString(digest[:tempDownloadDigestBytes]), nil
}

// candidate is the temporary path for one uniqueness token. Only the last
// element of the name becomes temporary, so a file staged for a folder
// download stays next to where it is going.
func (n tmpDownloadNaming) candidate(uniqueness string) destinationPath {
	dir, name := filepath.Split(n.base.name)
	return n.base.withName(filepath.Join(dir, tmpDownloadElement(name, n.pathDigest, uniqueness)))
}

// existingTmpDownloadPath returns the canonical in-flight temporary download of
// final when it exists: the file a transfer of final writes, or one a transfer
// left without pausing. It is never continued from; see pausedTmpDownloadPath.
func existingTmpDownloadPath(final destinationPath, tempPath string) (destinationPath, bool) {
	naming, err := tmpDownloadNamingFor(final, tempPath)
	if err != nil {
		return destinationPath{}, false
	}
	return existingTmpDownloadFile(final, naming.candidate(""))
}

// pausedTmpDownloadPath returns the paused temporary download of final under
// the canonical paused name: what a pause of a transfer that staged under the
// canonical in-flight name leaves, found without listing the directory. A
// transfer that staged under a uniqueness token (the canonical in-flight name
// was occupied) pauses under a tokenized paused name; see
// Job.tokenizedPausedTmpDownload for how a later job finds that one.
func pausedTmpDownloadPath(final destinationPath, tempPath string) (destinationPath, bool) {
	naming, err := tmpDownloadNamingFor(final, tempPath)
	if err != nil {
		return destinationPath{}, false
	}
	return existingTmpDownloadFile(final, naming.candidate(tempDownloadPausedToken))
}

// tokenizedPausedTmpDownload finds a paused temporary download of final that a
// transfer left under a tokenized paused name. The stage directory is listed
// once per job and directory, on the first file that needs it, and the
// listing is kept for the job (pausedListings): a folder download of many
// files then costs one listing per directory rather than one per file. Files
// paused by another process after the listing was taken are not seen by this
// job; a later job sees them. When several such files exist for one
// destination, the most recently modified one is continued from; each is a
// verified prefix of the same file.
func (r *Job) tokenizedPausedTmpDownload(final destinationPath, tempPath string) (destinationPath, bool) {
	naming, err := tmpDownloadNamingFor(final, tempPath)
	if err != nil {
		return destinationPath{}, false
	}
	canonical := naming.candidate(tempDownloadPausedToken)
	_, digest, _, ok := encodedTmpDownloadParts(filepath.Base(canonical.name))
	if !ok {
		return destinationPath{}, false
	}
	dir := canonical.withName(filepath.Dir(canonical.name))
	var found destinationPath
	var foundAt time.Time
	for _, name := range r.pausedListing(dir)[digest] {
		candidate := dir.withName(filepath.Join(dir.name, name))
		file, ok := existingTmpDownloadFile(final, candidate)
		if !ok {
			continue
		}
		info, err := file.lstat()
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		if found.name == "" || info.ModTime().After(foundAt) {
			found, foundAt = file, info.ModTime()
		}
	}
	return found, found.name != ""
}

// pausedListing returns the tokenized paused entries of dir by file digest,
// listing dir the first time the job asks about it.
func (r *Job) pausedListing(dir destinationPath) map[string][]string {
	key := filepath.Clean(dir.String())
	r.pausedListingsMutex.Lock()
	defer r.pausedListingsMutex.Unlock()
	if listing, ok := r.pausedListings[key]; ok {
		return listing
	}
	listing := map[string][]string{}
	entries, err := inRoot(dir, func(root *os.Root) ([]fs.DirEntry, error) {
		return fs.ReadDir(root.FS(), filepath.ToSlash(dir.name))
	})
	if err == nil {
		for _, entry := range entries {
			token, digest, _, ok := encodedTmpDownloadParts(entry.Name())
			if ok && strings.HasPrefix(token, tempDownloadPausedToken) && token != tempDownloadPausedToken {
				listing[digest] = append(listing[digest], entry.Name())
			}
		}
	}
	if r.pausedListings == nil {
		r.pausedListings = map[string]map[string][]string{}
	}
	r.pausedListings[key] = listing
	return listing
}

// plainTmpDownloadName returns the file name inside a plain in-flight element,
// ".~files-cli.<name>.download".
func plainTmpDownloadName(element string) (string, bool) {
	suffix := "." + TempDownloadExtension
	if !strings.HasPrefix(element, tempDownloadPrefix) || !strings.HasSuffix(element, suffix) {
		return "", false
	}
	name := element[len(tempDownloadPrefix) : len(element)-len(suffix)]
	return name, name != ""
}

// encodedTmpDownloadParts splits an encoded element,
// ".~files-cli~<token>.<digest>.<shortened>.download", into its parts.
func encodedTmpDownloadParts(element string) (token string, digest string, shortened string, ok bool) {
	suffix := "." + TempDownloadExtension
	if !strings.HasPrefix(element, tempDownloadEncodedPrefix) || !strings.HasSuffix(element, suffix) {
		return "", "", "", false
	}
	body := strings.TrimSuffix(element[len(tempDownloadEncodedPrefix):], suffix)
	token, rest, found := strings.Cut(body, ".")
	if !found {
		return "", "", "", false
	}
	digest, shortened, found = strings.Cut(rest, ".")
	if !found || digest == "" {
		return "", "", "", false
	}
	return token, digest, shortened, true
}

// encodedTmpDownloadElement spells an encoded element and keeps it within the
// path element limit, shortening the readable part as tmpDownloadElement does.
func encodedTmpDownloadElement(token string, digest string, shortened string) string {
	encoded := tempDownloadEncodedPrefix + token + "." + digest + "."
	shortened = truncateOnRuneBoundary(shortened, maxTempDownloadElementBytes-len(encoded)-len(TempDownloadExtension)-1)
	return encoded + shortened + "." + TempDownloadExtension
}

// tmpDownloadToken returns the uniqueness token of a temporary download's
// reserved element: "" for the plain form, otherwise what its encoded name
// carries, which for a paused file starts with tempDownloadPausedToken.
func tmpDownloadToken(element string) (string, bool) {
	if _, ok := plainTmpDownloadName(element); ok {
		return "", true
	}
	token, _, _, ok := encodedTmpDownloadParts(element)
	return token, ok
}

// tmpDownloadElementWithToken respells a reserved element under another
// token. A plain element is encoded first; an encoded one keeps its digest and
// readable part.
func tmpDownloadElementWithToken(element string, token string) (string, bool) {
	if name, ok := plainTmpDownloadName(element); ok {
		return tmpDownloadElement(name, "", token), true
	}
	_, digest, shortened, ok := encodedTmpDownloadParts(element)
	if !ok {
		return "", false
	}
	return encodedTmpDownloadElement(token, digest, shortened), true
}

func isPausedTmpDownloadElement(element string) bool {
	token, ok := tmpDownloadToken(element)
	return ok && strings.HasPrefix(token, tempDownloadPausedToken)
}

// pausedTmpDownloadElement is the paused name of a temporary download's
// reserved element (the file, or on macOS the folder): the same name with the
// paused token in front of its uniqueness token. Each in-flight name has
// exactly one paused name and no two attempts share one, so a pause never
// writes where another attempt's paused file could be. A paused element is
// its own.
func pausedTmpDownloadElement(element string) (string, bool) {
	token, ok := tmpDownloadToken(element)
	if !ok {
		return "", false
	}
	if strings.HasPrefix(token, tempDownloadPausedToken) {
		return element, true
	}
	return tmpDownloadElementWithToken(element, tempDownloadPausedToken+token)
}

// pausedTmpDownloadOf is the path the temporary download tmp has once it is
// parked; it is not renamed here.
func pausedTmpDownloadOf(tmp destinationPath) (destinationPath, bool) {
	entry := tmpDownloadEntry(tmp)
	element, ok := pausedTmpDownloadElement(entry.base())
	if !ok {
		return destinationPath{}, false
	}
	return tmpDownloadInEntry(entry.withName(filepath.Join(filepath.Dir(entry.name), element)), tmp.base()), true
}

// parkTmpDownload marks the temporary download tmp as continued from later.
// Every writer of the file must have closed it. The file is first claimed under
// a name nothing else can reach and checked there to be expected, the file this
// transfer wrote (claimTmpDownload): a filesystem can reach one file through
// more than one name, so another transfer can have put its own file at tmp,
// and that file must not become a checkpoint. Only the verified file is made
// durable and given the paused name, so nothing that can be discovered under a
// paused name was ever unverified. As at publication, a replacement is removed
// and reported. It reports where the file now is.
func parkTmpDownload(tmp destinationPath, expected fs.FileInfo) (destinationPath, error) {
	entry := tmpDownloadEntry(tmp)
	element, ok := pausedTmpDownloadElement(entry.base())
	if !ok {
		return destinationPath{}, fmt.Errorf("download: %v is not a temporary download and cannot be kept for a resume", tmp)
	}
	if element == entry.base() {
		return tmp, nil
	}
	claimed, err := claimTmpDownload(tmp, expected)
	if err != nil {
		return destinationPath{}, err
	}
	if err := syncTmpDownload(claimed); err != nil {
		return destinationPath{}, errors.Join(err, restoreTmpDownload(claimed, tmp))
	}
	paused, err := publishPausedTmpDownload(claimed, tmp, element)
	if err != nil {
		return destinationPath{}, errors.Join(err, restoreTmpDownload(claimed, tmp))
	}
	return paused, nil
}

// refuseOccupied reports fs.ErrExist when something is already under target.
// It is a check, not an atomic no-replace rename: a file put there between the
// check and the rename would be replaced. The names checked this way are
// specific to one attempt (its uniqueness token, or the paused form of it), so
// in the SDK's own use only this attempt writes them; see
// renameTmpDownloadEntry for the one exception.
func refuseOccupied(target destinationPath) error {
	if _, err := target.lstat(); err == nil {
		return &fs.PathError{Op: "rename", Path: target.String(), Err: fs.ErrExist}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

// activateTmpDownload moves the paused temporary download of final back under
// an in-flight name before a transfer writes to it again, so a process that
// stops before the next pause leaves nothing that looks continued from. The
// name is the one the file paused from when that is free, otherwise an unused
// one of the same kind: another attempt may be staging under the first, and
// its file is left alone. A checkpointed file elsewhere goes back to its own
// in-flight name only.
func activateTmpDownload(paused destinationPath, final destinationPath, tempPath string) (destinationPath, error) {
	entry := tmpDownloadEntry(paused)
	token, ok := tmpDownloadToken(entry.base())
	if !ok || !strings.HasPrefix(token, tempDownloadPausedToken) {
		return destinationPath{}, fmt.Errorf("download: %v is not a paused temporary download", paused)
	}
	token = strings.TrimPrefix(token, tempDownloadPausedToken)
	naming, err := tmpDownloadNamingFor(final, tempPath)
	inPlace := err == nil && sameDirectory(naming.candidate(""), entry)
	if !inPlace {
		element, ok := tmpDownloadElementWithToken(entry.base(), token)
		if !ok {
			return destinationPath{}, fmt.Errorf("download: %v is not a paused temporary download", paused)
		}
		return renameTmpDownloadEntry(paused, element)
	}
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	for attempt := 0; attempt <= 25; attempt++ {
		candidate := naming.candidate(token)
		if _, err := candidate.lstat(); errors.Is(err, fs.ErrNotExist) {
			return renameTmpDownloadEntry(paused, filepath.Base(candidate.name))
		}
		token = tmpDownloadUniqueness(attempt+1, randGenerator)
	}
	return destinationPath{}, fmt.Errorf("download: no unused temporary name to continue %v under", paused)
}

// tmpDownloadUniqueness is the uniqueness token of the index'th alternate
// temporary name: a small number first, then letters.
func tmpDownloadUniqueness(index int, randGenerator *rand.Rand) string {
	if index > 10 {
		var uniqueness string
		for i := 0; i < 4; i++ {
			uniqueness += string(rune(randGenerator.Intn(26) + 'a'))
		}
		return uniqueness
	}
	return fmt.Sprintf("%v", index)
}

func sameDirectory(a destinationPath, b destinationPath) bool {
	return filepath.Clean(filepath.Join(a.dir, filepath.Dir(a.name))) == filepath.Clean(filepath.Join(b.dir, filepath.Dir(b.name)))
}

// renameTmpDownloadEntry moves the temporary download tmp to another reserved
// name in the same directory. It is used to take a paused file back into
// flight, under a name that was unused a moment ago. The check and the rename
// are two steps, not an atomic no-replace rename: an attempt that starts in
// between and allocates that very name (tmpDownloadPath checks and creates in
// two steps as well, a pre-existing race) could have its file replaced, and
// the identity checks at pause and publication only detect a file that was
// replaced, not two writers holding the same file. Two transfers of one file
// into one destination at the same time are not a supported arrangement, and
// this change neither introduces nor resolves that.
func renameTmpDownloadEntry(tmp destinationPath, element string) (destinationPath, error) {
	entry := tmpDownloadEntry(tmp)
	target := entry.withName(filepath.Join(filepath.Dir(entry.name), element))
	if err := refuseOccupied(target); err != nil {
		return destinationPath{}, err
	}
	if err := entry.rename(target); err != nil {
		return destinationPath{}, err
	}
	return tmpDownloadInEntry(target, tmp.base()), nil
}

// syncTmpDownload makes the bytes of tmp durable before its name says they
// were received.
func syncTmpDownload(tmp destinationPath) error {
	file, err := tmp.openFile(os.O_RDWR, 0)
	if err != nil {
		return err
	}
	return errors.Join(file.Sync(), file.Close())
}

// discardUnpausedCheckpoint removes the file the caller's checkpoint names when
// it is still under an in-flight name: it was left by a run that stopped while
// writing, and its length says nothing about what was received. The path is
// resolved exactly as the lookup resolves it, inside the caller-selected root;
// a path that root refuses is not touched. Files under the canonical in-flight
// name that this transfer did not checkpoint are left alone: another attempt
// may be writing them.
func (d *DownloadStatus) discardUnpausedCheckpoint() {
	tmp, ok := d.checkpointedTmpDownloadPath()
	if !ok || isPausedTmpDownloadElement(tmpDownloadEntry(tmp).base()) {
		return
	}
	if _, err := tmp.lstat(); err != nil {
		return
	}
	logger := d.Job().Logger
	if err := removeTmpDownload(tmp); err != nil && !errors.Is(err, fs.ErrNotExist) {
		logger.Printf("could not remove the unfinished temporary download %v, starting over beside it: %v", tmp, err)
		return
	}
	logger.Printf("discarding the unfinished temporary download %v: it was not paused, so its bytes cannot be continued from", tmp)
}

// tmpDownloadPath generates a unique temporary download path for final inside
// the reserved namespace and, if necessary, with additional identifiers to
// avoid name conflicts.
func tmpDownloadPath(final destinationPath, tempPath string) (destinationPath, error) {
	var index int
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	naming, err := tmpDownloadNamingFor(final, tempPath)
	if err != nil {
		return destinationPath{}, err
	}

	for {
		var uniqueness string
		if index > 25 {
			return destinationPath{}, fmt.Errorf("unable to create a unique temporary path after 25 attempts, consider deleting existing %v*.%v files; attempted path: %v", tempDownloadNamespace, TempDownloadExtension, naming.base)
		} else if index > 0 {
			uniqueness = tmpDownloadUniqueness(index, randGenerator)
		}
		candidate := naming.candidate(uniqueness)

		if _, err := candidate.lstat(); errors.Is(err, fs.ErrNotExist) {
			return tmpDownloadPathOnNotExist(final, candidate)
		}
		index++
	}
}

// tmpDownloadElement names the temporary file or folder of one file name.
//
// Beside the final file, a name that fits is kept as it is, so an unfinished
// download is recognizable and the same name always produces the same
// temporary path. A name that is too long, or that needs a uniqueness token,
// is encoded after a different prefix character and identified by a digest of
// the whole name. Inside an external temp directory the name is always encoded
// and identified by pathDigest, the digest of the file's whole local path,
// because files from every folder are staged side by side there. The shortened
// name in an encoded element is only there to be read: the digest, not the
// shortened name, decides which file the element belongs to, so no other file
// can take over an encoded temporary path.
func tmpDownloadElement(name string, pathDigest string, uniqueness string) string {
	plain := tempDownloadPrefix + name + "." + TempDownloadExtension
	if pathDigest == "" && uniqueness == "" && len(plain) <= maxTempDownloadElementBytes {
		return plain
	}

	digest := pathDigest
	if digest == "" {
		sum := sha256.Sum256([]byte(name))
		digest = hex.EncodeToString(sum[:tempDownloadDigestBytes])
	}
	encoded := tempDownloadEncodedPrefix + uniqueness + "." + digest + "."
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

// localPathReachesTempDownload reports whether a local path is, or is inside, a
// temporary download, including through another name the filesystem resolves to
// the same file. It reports an error rather than an answer when a path that
// exists cannot be resolved, so a transfer is refused instead of being allowed
// on the strength of a check that did not run.
func localPathReachesTempDownload(path string) (bool, error) {
	if isReservedTempDownloadPath(path) {
		return true, nil
	}
	canonical, err := canonicalLocalPath(path)
	if err != nil {
		return false, err
	}
	return isReservedTempDownloadPath(canonical), nil
}

// recordTmpDownloadIdentity remembers which file the open temporary download
// is, so publishing can tell this transfer's file from any other file that
// reaches the same path while the transfer runs. A resumed transfer has already
// recorded the file it found; the file it then opened must be that same file,
// or the bytes it is about to append would go onto another file's content.
func (d *DownloadStatus) recordTmpDownloadIdentity(file *os.File) error {
	got, err := file.Stat()
	if err != nil {
		return err
	}
	if d.tmpIdentity != nil && !os.SameFile(d.tmpIdentity, got) {
		return fmt.Errorf("download: %w: %v", errTempDownloadReplaced, file.Name())
	}
	d.tmpIdentity = got
	return nil
}

// tmpDownloadClaimContainer is the folder a finished download is claimed in
// before it is published. Its name is in the reserved namespace, so whatever is
// left in it if this process stops is never transferred. When the temporary
// file already sits in such a folder, as a macOS temporary download does, that
// folder is used. Otherwise, including for a temporary file the caller named
// explicitly (a checkpoint's TmpPath), a fresh folder is created beside the
// temporary file, which is where this transfer has already been writing: the
// caller-selected directory above it may not be writable. The claim inside it
// is a random 8.3 name, so NTFS derives no second name for the claimed file.
func tmpDownloadClaimContainer(tmp destinationPath) (destinationPath, error) {
	parent := tmp.parent()
	if isReservedTempDownloadElement(filepath.Base(parent.name)) {
		return parent, nil
	}
	container := parent.withName(filepath.Join(parent.name, tempDownloadEncodedPrefix+cryptorand.Text()+"."+TempDownloadExtension))
	if err := container.mkdir(); err != nil {
		return destinationPath{}, err
	}
	return container, nil
}

// finalizeTmpDownload publishes the finished temporary download as final. It
// first claims the file under a name nothing else can reach and checks that it
// is the file this transfer wrote (claimTmpDownload), then moves it into place.
// If the move fails, for example because the job was paused, the file is put
// back under its temporary name so a later run continues from it; a failure to
// put it back is reported with the original error, because the download then
// has to start again.
func finalizeTmpDownload(ctx context.Context, tmp destinationPath, final destinationPath, expected fs.FileInfo) error {
	claimed, err := claimTmpDownload(tmp, expected)
	if err != nil {
		return err
	}
	if err := claimed.moveTo(ctx, final); err != nil {
		restoreErr := restoreTmpDownload(claimed, tmp)
		if restoreErr == nil && errors.Is(context.Cause(ctx), ErrJobPaused) {
			// The file is complete and verified; a later run may publish it
			// without asking for content, so it is kept under its paused name.
			_, restoreErr = parkTmpDownload(tmp, expected)
		}
		return errors.Join(err, restoreErr)
	}
	return errors.Join(removeTmpDownloadClaimContainer(claimed, tmp), removeTmpDownloadFolder(tmp))
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
// checkpointedTmpDownloadPath resolves the path the caller's checkpoint names
// the way every operation on it must: inside the caller-selected directory or
// the external temp directory when it lies below one of them, so the operating
// system refuses a path that leaves them; otherwise as an explicit path.
func (d *DownloadStatus) checkpointedTmpDownloadPath() (destinationPath, bool) {
	checkpointed := d.tmpDownloadPath()
	if checkpointed == "" {
		return destinationPath{}, false
	}
	tmp := explicitTmpDownload(checkpointed)
	for _, dir := range []string{d.destination.dir, d.tempPath} {
		if dir == "" {
			continue
		}
		if rel, err := filepath.Rel(dir, checkpointed); err == nil && rel != "." && filepath.IsLocal(rel) {
			tmp = destinationPath{dir: dir, name: rel}
			break
		}
	}
	return tmp, true
}

// checkpointedTmpDownload returns the paused temporary download the caller's
// checkpoint names, when there is one. The checkpoint may name the file under
// its in-flight name, as it was while the transfer ran; the paused file it
// became counts as well. A file still under an in-flight name does not.
func (d *DownloadStatus) checkpointedTmpDownload() (destinationPath, bool) {
	tmp, ok := d.checkpointedTmpDownloadPath()
	if !ok {
		return destinationPath{}, false
	}
	paused, ok := pausedTmpDownloadOf(tmp)
	if !ok {
		return destinationPath{}, false
	}
	if _, err := paused.stat(); err != nil {
		return destinationPath{}, false
	}
	return paused, true
}

// tmpDownloadToResume returns the paused temporary file a download continues
// from: the checkpointed one, otherwise the canonical one, otherwise one this
// job finds under a tokenized paused name.
func (d *DownloadStatus) tmpDownloadToResume() (destinationPath, bool) {
	if tmp, ok := d.checkpointedTmpDownload(); ok {
		return tmp, true
	}
	if tmp, ok := pausedTmpDownloadPath(d.destination, d.tempPath); ok {
		return tmp, true
	}
	if d.job == nil {
		return destinationPath{}, false
	}
	return d.job.tokenizedPausedTmpDownload(d.destination, d.tempPath)
}
