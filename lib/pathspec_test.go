package lib

import (
	"os/exec"
	"sync"
	"testing"
)

func TestRsync(t *testing.T) {
	_, err := exec.LookPath("rsync")
	if err != nil {
		t.Log(err)
		return
	}

	destinationFs := ReadWriteFs(LocalFileSystem{})
	sourceFs := ReadWriteFs(LocalFileSystem{})

	t.Run("rsync", func(t *testing.T) {
		mutex := &sync.Mutex{}
		for _, tt := range PathSpec(sourceFs.PathSeparator(), destinationFs.PathSeparator()) {
			t.Run(tt.Name, func(t *testing.T) {
				BuildPathSpecTest(t, mutex, tt, sourceFs, destinationFs, func(source, destination string) Cmd {
					return ExeCmd{Cmd: exec.Command("rsync", "-av", source, destination)}
				})
			})
		}
	})
}
