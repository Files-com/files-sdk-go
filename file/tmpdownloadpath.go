package file

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"path/filepath"
	"time"
)

const TempDownloadExtension = "download"

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
	base := tmpDownloadBase(final, tempPath)
	return existingTmpDownloadFile(final, base.withName(fmt.Sprintf("%v.%v", base.name, TempDownloadExtension)))
}

// tmpDownloadPath generates a unique temporary download path for final by appending a ".download" extension and, if necessary, additional identifiers to avoid name conflicts.
func tmpDownloadPath(final destinationPath, tempPath string) (destinationPath, error) {
	var index int
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	base := tmpDownloadBase(final, tempPath)

	for {
		var candidate destinationPath
		var uniqueness string
		if index == 0 {
			candidate = base.withName(fmt.Sprintf("%v.%v", base.name, TempDownloadExtension))
		} else if index > 25 {
			return destinationPath{}, fmt.Errorf("unable to create a unique temporary path after 25 attempts, consider deleting existing .%v files; attempted path: %v", TempDownloadExtension, base)
		} else {
			if index > 10 {
				for i := 0; i < 4; i++ {
					uniqueness += string(rune(randGenerator.Intn(26) + 'a'))
				}
			} else {
				uniqueness = fmt.Sprintf("%v", index)
			}
			candidate = base.withName(fmt.Sprintf("%v (%v).%v", base.name, uniqueness, TempDownloadExtension))
		}

		if _, err := candidate.lstat(); errors.Is(err, fs.ErrNotExist) {
			return tmpDownloadPathOnNotExist(final, candidate)
		}
		index++
	}
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
