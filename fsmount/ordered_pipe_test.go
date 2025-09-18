package fsmount

import (
	"fmt"
	"io"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

var logger = lib.NewLeveledLogger(lib.NullLogger{})

type Fakefs struct {
	orderdPipe *orderedPipe
}

func (f *Fakefs) Write(p []byte, offset int) (n int, err error) {
	return f.orderdPipe.writeAt(p, int64(offset))
}

func (fs *Fakefs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	return fs.orderdPipe.readAt(buff, ofst)
}

type writeAtOffset struct {
	Offset int64
	Data   string
}

func TestOutOfOrderWrites(t *testing.T) {
	op, err := newOrderedPipe("/test/path", logger)
	if err != nil {
		t.Fatalf("Error creating ordered pipe: %v", err)
	}
	f := &Fakefs{
		orderdPipe: op,
	}

	// This slice is made up of the string
	// "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do"
	// with the offsets where each part of the string should be written.
	writeOffsets := []writeAtOffset{
		{Offset: 6, Data: "ipsum "},
		{Offset: 0, Data: "Lorem "},
		{Offset: 18, Data: "sit "},
		{Offset: 12, Data: "dolor "},
		{Offset: 28, Data: "consectetur "},
		{Offset: 22, Data: "amet, "},
		{Offset: 51, Data: "elit. "},
		{Offset: 40, Data: "adipiscing "},
		{Offset: 61, Data: "do"},
		{Offset: 57, Data: "Sed "},
	}

	// make the channel capable of buffering at least as many errors as there are
	// writes to ensure the test doesn't hang if there are errors
	errChanLen := len(writeOffsets) + 1
	errchan := make(chan error, errChanLen)

	go func() {
		defer close(errchan)
		for _, w := range writeOffsets {
			_, err := f.Write([]byte(w.Data), int(w.Offset))
			if err != nil {
				errchan <- fmt.Errorf("Error writing to sorted pipe: %v", err)
			}
		}
		f.orderdPipe.close()
	}()

	var sortedData []byte

	sortedData, err = io.ReadAll(f.orderdPipe.out)
	if err != nil {
		t.Errorf("Error reading from sorted pipe: %v", err)
	}

	for err := range errchan {
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	}
	got := string(sortedData)
	want := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do"
	if got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}

func TestReaderAt(t *testing.T) {
	op, err := newOrderedPipe("/test/path", logger)
	if err != nil {
		t.Fatalf("Error creating ordered pipe: %v", err)
	}
	f := &Fakefs{
		orderdPipe: op,
	}

	// This slice is made up of the string
	// "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do"
	// with the offsets where each part of the string should be written.
	writeOffsets := []writeAtOffset{
		{Offset: 6, Data: "ipsum "},        // bytes 6-11
		{Offset: 0, Data: "Lorem "},        // bytes 0-5
		{Offset: 18, Data: "sit "},         // bytes 18-21
		{Offset: 12, Data: "dolor "},       // bytes 12-17
		{Offset: 28, Data: "consectetur "}, // bytes 28-39
		{Offset: 22, Data: "amet, "},       // bytes 22-27
		{Offset: 51, Data: "elit. "},       // bytes 51-56
		{Offset: 40, Data: "adipiscing "},  // bytes 40-50
		{Offset: 61, Data: "do"},           // bytes 61-62
		{Offset: 57, Data: "Sed "},         // bytes 57-60
	}

	// make the channel capable of buffering at least as many errors as there are
	// writes to ensure the test doesn't hang if there are errors
	errChanLen := len(writeOffsets) + 1
	errchan := make(chan error, errChanLen)

	go func() {
		defer f.orderdPipe.close()
		defer close(errchan)
		for idx, w := range writeOffsets {
			_, err := f.Write([]byte(w.Data), int(w.Offset))
			if err != nil {
				errchan <- fmt.Errorf("Error writing to sorted pipe: %v", err)
			}
			if idx == 3 {
				// At this point, the data in the pipe should be:
				// "Lorem ipsum dolor sit " based on completing the first four writes.
				readBuff := make([]byte, 6)
				// Call Read with a start offset of 6, which should read "ipsum " and pass
				// a buffer of 6 bytes. The result should be "ipsum ".
				n := f.Read("/test/path", readBuff, 6, 123)
				if n != 6 {
					errchan <- fmt.Errorf("Expected to read 6 bytes, but got %d", n)
				}
				if string(readBuff) != "ipsum " {
					errchan <- fmt.Errorf("Expected to read 'ipsum ', but got '%s'", string(readBuff))
				}
			}
			if idx == 7 {
				// At this point, the data in the pipe should be:
				// "Lorem ipsum dolor sit amet, consectetur "
				// based on completing the first eight writes.
				readBuff := make([]byte, 12)
				// Call Read with a start offset of 28, which should read "consectetur"...
				// and pass a buffer of 12 bytes. The result should be "consectetur ".
				n := f.Read("/test/path", readBuff, 28, 123)
				if n != 12 {
					errchan <- fmt.Errorf("Expected to read 12 bytes, but got %d", n)
				}
				if string(readBuff) != "consectetur " {
					errchan <- fmt.Errorf("Expected to read 'consectetur ', but got '%s'", string(readBuff))
				}
			}
		}
	}()

	var sortedData []byte

	sortedData, err = io.ReadAll(f.orderdPipe.out)
	if err != nil {
		t.Errorf("Error reading from sorted pipe: %v", err)
	}

	for err := range errchan {
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	}
	got := string(sortedData)
	want := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do"
	if got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

}
