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

type writeAtOffset struct {
	Offset int64
	Data   string
}

func TestOrderdedPipe(t *testing.T) {
	f := &Fakefs{
		orderdPipe: newOrderedPipe("/test/path", 123, logger),
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
		defer f.orderdPipe.done()
		defer close(errchan)
		for _, w := range writeOffsets {
			_, err := f.Write([]byte(w.Data), int(w.Offset))
			if err != nil {
				errchan <- fmt.Errorf("Error writing to sorted pipe: %v", err)
			}
		}
	}()

	var sortedData []byte
	var err error

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
