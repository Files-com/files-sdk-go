//go:build !windows

package file

import (
	"io"
	"os"
)

// Elsewhere than Windows, reads of one *os.File run concurrently, so the parts of an upload share the file it opened, and
// none reads through a handle of its own (uploadv2_ownedparts_windows.go). These stand-ins hold nothing and do nothing.

type uploadOwnedFile struct{}

func uploadWithOwnedFile(file *os.File) UploadOption {
	return UploadWithReaderAt(file)
}

type uploadV2OwnedFile struct{}

func newUploadV2OwnedFile(*uploadIO) uploadV2OwnedFile {
	return uploadV2OwnedFile{}
}

func (uploadV2OwnedFile) readerFor(readerAt io.ReaderAt) (io.ReaderAt, uploadV2PartSource) {
	return readerAt, uploadV2PartSource{}
}

type uploadV2PartSource struct{}

func (*uploadV2Engine) attachSource(*uploadV2Part) {}

func (*uploadV2Part) releaseSource() {}

type uploadV2OwnedPartStats struct{}

func (*uploadV2OwnedPartStats) addLogAttrs(map[string]any) {}
