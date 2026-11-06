package file

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	boundPausedDownloadToken    = "b1-"
	boundPausedDownloadData     = "DATA.BIN"
	boundPausedDownloadRecord   = "BINDING"
	maxPausedDownloadProofBytes = 1 << 20
	pausedDownloadProofMagic    = "FSDK-B1\x00"
)

var errDownloadPauseProof = errors.New("download: paused representation proof could not be settled")

func pauseProofError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Join(errDownloadPauseProof, err)
}

func isBoundPausedDownloadElement(element string) bool {
	token, _, _, ok := encodedTmpDownloadParts(element)
	return ok && strings.HasPrefix(token, boundPausedDownloadToken)
}

func boundPausedDownloadContainer(tmp destinationPath) (destinationPath, bool) {
	parent := tmp.parent()
	return parent, tmp.base() == boundPausedDownloadData && isBoundPausedDownloadElement(parent.base())
}

type pausedProofFile interface {
	io.ReadWriteCloser
	Stat() (fs.FileInfo, error)
	Sync() error
}

// A narrow file boundary permits real filesystem tests to return independent
// sync/close errors without replacing the artifact or transfer implementation.
var openPausedProofFile = func(root *os.Root, name string, flag int, perm fs.FileMode) (pausedProofFile, error) {
	return root.OpenFile(name, flag, perm)
}

type boundPausedDownload struct {
	data, container                                destinationPath
	binding                                        DownloadBinding
	dataIdentity, proofIdentity, containerIdentity fs.FileInfo
}

type pausedDownloadCustodyError struct {
	cause    error
	terminal bool
}

func (e pausedDownloadCustodyError) Error() string {
	return fmt.Sprintf("download: recovered paused artifact changed: %v", e.cause)
}
func (e pausedDownloadCustodyError) Unwrap() []error {
	return []error{e.cause, errDownloadPauseProof, errTempDownloadReplaced}
}

func observePausedIdentity(path destinationPath, expected fs.FileInfo) (bool, error) {
	got, err := path.lstat()
	if err != nil {
		return errors.Is(err, fs.ErrNotExist), err
	}
	if expected == nil || !os.SameFile(expected, got) || expected.Mode().Type() != got.Mode().Type() {
		return true, fmt.Errorf("%w: %v", errTempDownloadReplaced, path)
	}
	return false, nil
}

// An observed invalid b1 must not be rediscovered and accepted on a retry.
// Quarantine only a container whose original identity still matches. All child
// bytes are retained; an unfamiliar child or lost container makes this file
// terminal instead of moving somebody else's directory.
func quarantinePausedCustody(a *boundPausedDownload, cause error) error {
	changed, err := observePausedIdentity(a.container, a.containerIdentity)
	if err != nil {
		return pausedDownloadCustodyError{cause: errors.Join(cause, err), terminal: true}
	}
	if changed {
		return pausedDownloadCustodyError{cause: cause, terminal: true}
	}
	err = a.container.do(func(parent *os.Root) error {
		root, err := parent.OpenRoot(a.container.name)
		if err != nil {
			return err
		}
		entries, readErr := fs.ReadDir(root.FS(), ".")
		if readErr == nil && len(entries) != 2 {
			readErr = errors.New("changed paused artifact has an unfamiliar child")
		}
		for _, entry := range entries {
			if entry.Name() != boundPausedDownloadData && entry.Name() != boundPausedDownloadRecord {
				readErr = errors.Join(readErr, errors.New("changed paused artifact has an unfamiliar child"))
			}
		}
		return errors.Join(readErr, root.Close())
	})
	if err != nil {
		return pausedDownloadCustodyError{cause: errors.Join(cause, err), terminal: true}
	}
	claim := a.container.withName(filepath.Join(filepath.Dir(a.container.name), tempDownloadPrefix+cryptorand.Text()+"."+TempDownloadExtension))
	_, err = movePausedOwnedPath(a.container, claim, a.containerIdentity)
	return pausedDownloadCustodyError{cause: errors.Join(cause, err), terminal: err != nil}
}

func checkRecoveredPausedIdentity(a *boundPausedDownload) error {
	changed, err := observePausedIdentity(a.container, a.containerIdentity)
	if err != nil {
		if changed {
			return pausedDownloadCustodyError{cause: err, terminal: true}
		}
		return pauseProofError(err)
	}
	for _, child := range []struct {
		path     destinationPath
		identity fs.FileInfo
	}{
		{a.data, a.dataIdentity},
		{a.container.withName(filepath.Join(a.container.name, boundPausedDownloadRecord)), a.proofIdentity},
	} {
		changed, err := observePausedIdentity(child.path, child.identity)
		if err != nil {
			if changed {
				return quarantinePausedCustody(a, err)
			}
			return pauseProofError(err)
		}
	}
	return nil
}

// readBoundPausedDownload validates one bounded record inside a pinned root.
// Malformed/legacy artifacts are not trusted. IO and close failures remain
// independent errors, never an implicit permission to retry or discard data.
func readBoundPausedDownload(tmp destinationPath) (artifact *boundPausedDownload, err error) {
	container, ok := boundPausedDownloadContainer(tmp)
	if !ok {
		return nil, nil
	}
	containerInfo, statErr := container.lstat()
	if errors.Is(statErr, fs.ErrNotExist) {
		return nil, nil
	}
	if statErr != nil {
		return nil, pauseProofError(statErr)
	}
	if !containerInfo.IsDir() || containerInfo.Mode()&os.ModeSymlink != 0 {
		return nil, nil
	}
	root, openErr := inRoot(container, func(parent *os.Root) (*os.Root, error) { return parent.OpenRoot(container.name) })
	if openErr != nil {
		return nil, pauseProofError(openErr)
	}
	defer func() { err = errors.Join(err, pauseProofError(root.Close())) }()
	rootInfo, statErr := root.Stat(".")
	if statErr != nil {
		return nil, pauseProofError(statErr)
	}
	if !os.SameFile(containerInfo, rootInfo) {
		return nil, pauseProofError(errTempDownloadReplaced)
	}
	dataInfo, dataErr := root.Lstat(boundPausedDownloadData)
	proofInfo, proofErr := root.Lstat(boundPausedDownloadRecord)
	if dataErr != nil && !errors.Is(dataErr, fs.ErrNotExist) || proofErr != nil && !errors.Is(proofErr, fs.ErrNotExist) {
		return nil, pauseProofError(errors.Join(dataErr, proofErr))
	}
	if dataErr != nil || proofErr != nil {
		return nil, nil
	}
	if !dataInfo.Mode().IsRegular() || !proofInfo.Mode().IsRegular() || proofInfo.Size() > maxPausedDownloadProofBytes {
		return nil, nil
	}
	file, openErr := openPausedProofFile(root, boundPausedDownloadRecord, os.O_RDONLY, 0)
	if openErr != nil {
		return nil, pauseProofError(openErr)
	}
	openedInfo, infoErr := file.Stat()
	if infoErr != nil || !os.SameFile(openedInfo, proofInfo) {
		return nil, pauseProofError(errors.Join(infoErr, errTempDownloadReplaced, file.Close()))
	}
	encoded, readErr := io.ReadAll(io.LimitReader(file, maxPausedDownloadProofBytes+1))
	if readErr = errors.Join(readErr, file.Close()); readErr != nil {
		return nil, pauseProofError(readErr)
	}
	binding, decodeErr := decodePausedDownloadProof(encoded)
	if decodeErr != nil || binding.TmpPath != checkpointLocalPath(tmp.String()) || binding.Prefix != dataInfo.Size() || binding.Source == "" || binding.ETag == "" || strongDownloadETag(binding.ETag) != binding.ETag {
		return nil, nil
	}
	return &boundPausedDownload{data: tmp, container: container, binding: binding, dataIdentity: dataInfo, proofIdentity: proofInfo, containerIdentity: containerInfo}, nil
}

func boundPausedDataIn(container destinationPath) destinationPath {
	return container.withName(filepath.Join(container.name, boundPausedDownloadData))
}

func existingBoundPausedData(container destinationPath) (destinationPath, bool) {
	info, err := container.lstat()
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return destinationPath{}, false
	}
	data := boundPausedDataIn(container)
	info, err = data.lstat()
	return data, err == nil && info.Mode().IsRegular()
}

func verifyPausedIdentity(path destinationPath, expected fs.FileInfo) error {
	got, err := path.lstat()
	if err != nil {
		return err
	}
	if expected == nil || !os.SameFile(expected, got) || expected.Mode().Type() != got.Mode().Type() {
		return fmt.Errorf("%w: %v", errTempDownloadReplaced, path)
	}
	return nil
}

// Targets are fresh private names, as in the existing claim implementation.
// This refuses an observed occupant; it is not an atomic no-replace operation
// against an external actor inserting a guessed name between check and rename.
func movePausedOwnedPath(from, to destinationPath, expected fs.FileInfo) (bool, error) {
	if err := verifyPausedIdentity(from, expected); err != nil {
		return false, err
	}
	if from.String() == to.String() {
		return false, nil
	}
	if err := refuseOccupied(to); err != nil {
		return false, err
	}
	if err := from.rename(to); err != nil {
		return false, err
	}
	return true, verifyPausedIdentity(to, expected)
}

func removePausedOwnedFile(path destinationPath, expected fs.FileInfo) error {
	if _, err := path.lstat(); errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	claim := path.withName(filepath.Join(filepath.Dir(path.name), tmpDownloadClaimName()))
	if _, err := movePausedOwnedPath(path, claim, expected); err != nil {
		return err // an uncertain replacement stays in the reserved claim
	}
	return claim.remove()
}

func removePausedOwnedContainer(container destinationPath, expected fs.FileInfo) error {
	if err := verifyPausedIdentity(container, expected); err != nil {
		return err
	}
	return container.remove() // only empty containers; never RemoveAll
}

func verifyPausedChildren(a *boundPausedDownload, children map[string]fs.FileInfo) error {
	if err := verifyPausedIdentity(a.container, a.containerIdentity); err != nil {
		return err
	}
	return a.container.do(func(parent *os.Root) error {
		root, err := parent.OpenRoot(a.container.name)
		if err != nil {
			return err
		}
		entries, readErr := fs.ReadDir(root.FS(), ".")
		if readErr == nil && len(entries) != len(children) {
			readErr = errors.New("paused artifact contains a foreign child")
		}
		for _, entry := range entries {
			expected, ok := children[entry.Name()]
			if !ok {
				readErr = errors.Join(readErr, errors.New("paused artifact contains a foreign child"))
				continue
			}
			got, err := root.Lstat(entry.Name())
			if err != nil || expected == nil || !os.SameFile(expected, got) || !got.Mode().IsRegular() {
				readErr = errors.Join(readErr, err, errTempDownloadReplaced)
			}
		}
		return errors.Join(readErr, root.Close())
	})
}

func writePausedProof(container destinationPath, encoded []byte) (identity fs.FileInfo, err error) {
	root, err := inRoot(container, func(parent *os.Root) (*os.Root, error) { return parent.OpenRoot(container.name) })
	if err != nil {
		return nil, err
	}
	defer func() { err = errors.Join(err, root.Close()) }()
	file, err := openPausedProofFile(root, boundPausedDownloadRecord, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	identity, statErr := file.Stat()
	n, writeErr := file.Write(encoded)
	if writeErr == nil && n != len(encoded) {
		writeErr = io.ErrShortWrite
	}
	return identity, errors.Join(statErr, writeErr, file.Sync(), file.Close())
}

// fallbackPausedClaim never restores to the old, potentially occupied name.
// It keeps only verified settled bytes under a fresh legacy p- name; a later
// automatic HTTP download therefore restarts from zero without a proof.
func fallbackPausedClaim(a *boundPausedDownload, template destinationPath, cause error) (destinationPath, error) {
	if err := verifyPausedIdentity(a.container, a.containerIdentity); err != nil {
		return destinationPath{}, pauseProofError(errors.Join(cause, err))
	}
	element, ok := tmpDownloadElementWithToken(tmpDownloadEntry(template).base(), tempDownloadPausedToken+cryptorand.Text())
	if !ok {
		return destinationPath{}, pauseProofError(cause)
	}
	entry := a.container.withName(filepath.Join(filepath.Dir(a.container.name), element))
	legacy, err := preparePausedProofLegacy(entry, filepath.Base(a.binding.LocalPath))
	if err != nil {
		return destinationPath{}, pauseProofError(errors.Join(cause, err))
	}
	if _, err = movePausedOwnedPath(a.data, legacy, a.dataIdentity); err != nil {
		return destinationPath{}, pauseProofError(errors.Join(cause, err))
	}
	var cleanupErr error
	if a.proofIdentity != nil {
		proof := a.container.withName(filepath.Join(a.container.name, boundPausedDownloadRecord))
		cleanupErr = removePausedOwnedFile(proof, a.proofIdentity)
	}
	cleanupErr = errors.Join(cleanupErr, removePausedOwnedContainer(a.container, a.containerIdentity))
	return legacy, pauseProofError(errors.Join(cause, cleanupErr))
}

func publishPausedBoundClaim(a *boundPausedDownload, template destinationPath) (destinationPath, error) {
	if err := syncTmpDownload(a.data); err != nil {
		return destinationPath{}, pauseProofError(err) // do not certify a failed data close
	}
	if a.proofIdentity != nil {
		proof := a.container.withName(filepath.Join(a.container.name, boundPausedDownloadRecord))
		if err := removePausedOwnedFile(proof, a.proofIdentity); err != nil {
			return fallbackPausedClaim(a, template, err)
		}
		a.proofIdentity = nil
	}
	element, ok := tmpDownloadElementWithToken(tmpDownloadEntry(template).base(), boundPausedDownloadToken+cryptorand.Text())
	if !ok {
		return fallbackPausedClaim(a, template, errors.New("invalid paused artifact name"))
	}
	target := a.container.withName(filepath.Join(filepath.Dir(a.container.name), element))
	a.binding.TmpPath = checkpointLocalPath(boundPausedDataIn(target).String())
	encoded, err := encodePausedDownloadProof(a.binding)
	if err != nil {
		return fallbackPausedClaim(a, template, err)
	}
	if err := verifyPausedChildren(a, map[string]fs.FileInfo{a.data.base(): a.dataIdentity}); err != nil {
		return fallbackPausedClaim(a, template, err)
	}
	data := boundPausedDataIn(a.container)
	moved, err := movePausedOwnedPath(a.data, data, a.dataIdentity)
	if moved {
		a.data = data
	}
	if err != nil {
		return fallbackPausedClaim(a, template, err)
	}
	a.proofIdentity, err = writePausedProof(a.container, encoded)
	if err != nil {
		return fallbackPausedClaim(a, template, err)
	}
	if err := verifyPausedChildren(a, map[string]fs.FileInfo{boundPausedDownloadData: a.dataIdentity, boundPausedDownloadRecord: a.proofIdentity}); err != nil {
		return fallbackPausedClaim(a, template, err)
	}
	moved, err = movePausedOwnedPath(a.container, target, a.containerIdentity)
	if moved {
		a.container, a.data = target, boundPausedDataIn(target)
	}
	if err != nil {
		return fallbackPausedClaim(a, template, err)
	}
	return a.data, nil
}

func parkBoundTmpDownload(tmp destinationPath, expected fs.FileInfo, binding DownloadBinding) (destinationPath, error) {
	claimed, err := claimTmpDownload(tmp, expected)
	if err != nil {
		return destinationPath{}, pauseProofError(err)
	}
	container := claimed.parent()
	containerInfo, err := container.lstat()
	if err != nil {
		return destinationPath{}, pauseProofError(err)
	}
	return publishPausedBoundClaim(&boundPausedDownload{data: claimed, container: container, binding: binding, dataIdentity: expected, containerIdentity: containerInfo}, tmp)
}

// claimPausedBoundDownload removes the public paused name before any mutation,
// and verifies both file identities again inside a fresh private container.
func claimPausedBoundDownload(tmp destinationPath, expected fs.FileInfo, recovered ...*boundPausedDownload) (*boundPausedDownload, error) {
	var original *boundPausedDownload
	if len(recovered) > 0 && recovered[0] != nil {
		original = recovered[0]
		if err := checkRecoveredPausedIdentity(original); err != nil {
			return nil, err
		}
	}
	a, err := readBoundPausedDownload(tmp)
	if original != nil {
		if custodyErr := checkRecoveredPausedIdentity(original); custodyErr != nil {
			return nil, errors.Join(err, custodyErr)
		}
	}
	if err != nil || a == nil {
		return nil, pauseProofError(errors.Join(err, errors.New("paused artifact has no matching proof")))
	}
	if expected != nil && !os.SameFile(expected, a.dataIdentity) {
		if original != nil {
			return nil, quarantinePausedCustody(original, errTempDownloadReplaced)
		}
		return nil, pauseProofError(errTempDownloadReplaced)
	}
	if original != nil {
		if !os.SameFile(original.dataIdentity, a.dataIdentity) || !os.SameFile(original.proofIdentity, a.proofIdentity) || !os.SameFile(original.containerIdentity, a.containerIdentity) {
			return nil, quarantinePausedCustody(original, errTempDownloadReplaced)
		}
		if original.binding != a.binding || checkpointLocalPath(original.data.String()) != checkpointLocalPath(tmp.String()) {
			return nil, quarantinePausedCustody(original, errors.New("recovered paused proof no longer matches its original binding"))
		}
		a.dataIdentity, a.proofIdentity, a.containerIdentity = original.dataIdentity, original.proofIdentity, original.containerIdentity
	}
	if err := verifyPausedChildren(a, map[string]fs.FileInfo{boundPausedDownloadData: a.dataIdentity, boundPausedDownloadRecord: a.proofIdentity}); err != nil {
		if original != nil {
			if custodyErr := checkRecoveredPausedIdentity(original); custodyErr != nil {
				return nil, errors.Join(err, custodyErr)
			}
		}
		return nil, pauseProofError(err)
	}
	claim := a.container.withName(filepath.Join(filepath.Dir(a.container.name), tempDownloadPrefix+cryptorand.Text()+"."+TempDownloadExtension))
	moved, err := movePausedOwnedPath(a.container, claim, a.containerIdentity)
	if moved {
		a.container, a.data = claim, boundPausedDataIn(claim)
	}
	if err != nil {
		if original != nil {
			if custodyErr := checkRecoveredPausedIdentity(a); custodyErr != nil {
				return nil, errors.Join(err, custodyErr)
			}
		}
		return nil, pauseProofError(err)
	}
	if err := verifyPausedChildren(a, map[string]fs.FileInfo{boundPausedDownloadData: a.dataIdentity, boundPausedDownloadRecord: a.proofIdentity}); err != nil {
		if original != nil {
			if custodyErr := checkRecoveredPausedIdentity(a); custodyErr != nil {
				return nil, errors.Join(err, custodyErr)
			}
		}
		return nil, pauseProofError(err)
	}
	return a, nil
}

type pausedDownloadPathError struct {
	error
	path destinationPath
}

func (e pausedDownloadPathError) Unwrap() error { return e.error }

func pausedPathFromError(err error) (destinationPath, bool) {
	var paused pausedDownloadPathError
	if errors.As(err, &paused) {
		return paused.path, true
	}
	return destinationPath{}, false
}

func finalizeBoundTmpDownload(ctx context.Context, tmp, final destinationPath, expected fs.FileInfo, recovered ...*boundPausedDownload) error {
	a, err := claimPausedBoundDownload(tmp, expected, recovered...)
	if err != nil {
		return err
	}
	moveErr := context.Cause(ctx)
	if moveErr == nil {
		moveErr = a.data.moveTo(ctx, final)
	}
	if err := moveErr; err != nil {
		if errors.Is(context.Cause(ctx), ErrJobPaused) && stoppedForPauseOnly(err) {
			paused, parkErr := publishPausedBoundClaim(a, tmp)
			joined := errors.Join(err, parkErr)
			if paused.name != "" {
				return pausedDownloadPathError{error: joined, path: paused}
			}
			return joined
		}
		cleanup := errors.Join(removePausedOwnedFile(a.data, a.dataIdentity), removePausedOwnedFile(a.container.withName(filepath.Join(a.container.name, boundPausedDownloadRecord)), a.proofIdentity))
		cleanup = errors.Join(cleanup, removePausedOwnedContainer(a.container, a.containerIdentity))
		if cleanup != nil || !stoppedForPauseOnly(err) {
			return pauseProofError(errors.Join(err, cleanup))
		}
		return err
	}
	proof := a.container.withName(filepath.Join(a.container.name, boundPausedDownloadRecord))
	return pauseProofError(errors.Join(removePausedOwnedFile(proof, a.proofIdentity), removePausedOwnedContainer(a.container, a.containerIdentity)))
}

func activateBoundTmpDownload(tmp destinationPath, final destinationPath, expected fs.FileInfo, recovered ...*boundPausedDownload) (destinationPath, error) {
	a, err := claimPausedBoundDownload(tmp, expected, recovered...)
	if err != nil {
		return destinationPath{}, err
	}
	element, _ := tmpDownloadElementWithToken(tmpDownloadEntry(tmp).base(), cryptorand.Text())
	entry := a.container.withName(filepath.Join(filepath.Dir(a.container.name), element))
	active, err := preparePausedProofLegacy(entry, final.base())
	if err != nil {
		return fallbackPausedClaim(a, tmp, err)
	}
	if _, err := movePausedOwnedPath(a.data, active, a.dataIdentity); err != nil {
		return fallbackPausedClaim(a, tmp, err)
	}
	proof := a.container.withName(filepath.Join(a.container.name, boundPausedDownloadRecord))
	err = errors.Join(removePausedOwnedFile(proof, a.proofIdentity), removePausedOwnedContainer(a.container, a.containerIdentity))
	return active, pauseProofError(err)
}

func removeBoundTmpDownload(tmp destinationPath) error {
	if _, err := tmp.lstat(); errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	a, err := claimPausedBoundDownload(tmp, nil)
	if err != nil {
		return err
	}
	proof := a.container.withName(filepath.Join(a.container.name, boundPausedDownloadRecord))
	err = errors.Join(removePausedOwnedFile(a.data, a.dataIdentity), removePausedOwnedFile(proof, a.proofIdentity))
	return pauseProofError(errors.Join(err, removePausedOwnedContainer(a.container, a.containerIdentity)))
}

// The private record preserves Go strings byte for byte, including file names
// and opaque validators that JSON would otherwise repair as invalid UTF-8.
// Its entire encoded size is checked before allocating a producer buffer.
func encodePausedDownloadProof(binding DownloadBinding) ([]byte, error) {
	fields := []string{binding.Source, binding.RemotePath, binding.LocalPath, binding.TmpPath, binding.ETag}
	size := len(pausedDownloadProofMagic) + 4*len(fields) + 16
	for _, field := range fields {
		if len(field) > maxPausedDownloadProofBytes-size {
			return nil, fmt.Errorf("paused representation proof exceeds %d bytes", maxPausedDownloadProofBytes)
		}
		size += len(field)
	}
	if binding.Prefix < 0 || binding.Total < binding.Prefix {
		return nil, errors.New("invalid paused representation sizes")
	}
	encoded := make([]byte, 0, size)
	encoded = append(encoded, pausedDownloadProofMagic...)
	for _, field := range fields {
		encoded = binary.BigEndian.AppendUint32(encoded, uint32(len(field)))
		encoded = append(encoded, field...)
	}
	encoded = binary.BigEndian.AppendUint64(encoded, uint64(binding.Total))
	encoded = binary.BigEndian.AppendUint64(encoded, uint64(binding.Prefix))
	return encoded, nil
}

func decodePausedDownloadProof(encoded []byte) (DownloadBinding, error) {
	var binding DownloadBinding
	if len(encoded) > maxPausedDownloadProofBytes || !bytes.HasPrefix(encoded, []byte(pausedDownloadProofMagic)) {
		return binding, errors.New("invalid paused representation proof format or size")
	}
	encoded = encoded[len(pausedDownloadProofMagic):]
	fields := []*string{&binding.Source, &binding.RemotePath, &binding.LocalPath, &binding.TmpPath, &binding.ETag}
	for _, field := range fields {
		if len(encoded) < 4 {
			return DownloadBinding{}, errors.New("truncated paused representation proof")
		}
		length := uint64(binary.BigEndian.Uint32(encoded[:4]))
		encoded = encoded[4:]
		if length > uint64(len(encoded)) {
			return DownloadBinding{}, errors.New("invalid paused representation field length")
		}
		*field = string(encoded[:int(length)])
		encoded = encoded[int(length):]
	}
	if len(encoded) != 16 {
		return DownloadBinding{}, errors.New("invalid paused representation proof suffix")
	}
	total, prefix := binary.BigEndian.Uint64(encoded[:8]), binary.BigEndian.Uint64(encoded[8:])
	if total > uint64(^uint64(0)>>1) || prefix > total {
		return DownloadBinding{}, errors.New("invalid paused representation sizes")
	}
	binding.Total, binding.Prefix = int64(total), int64(prefix)
	return binding, nil
}
