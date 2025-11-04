package io_test

import (
	"io"
	"testing"

	fsio "github.com/Files-com/files-sdk-go/v3/fsmount/internal/io"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
)

type writeAtOffset struct {
	Offset int64
	Data   string
}

func TestOutOfOrderWrites(t *testing.T) {
	op, err := fsio.NewOrderedPipe("/test/path", fsio.WithLogger(&log.NoOpLogger{}))
	if err != nil {
		t.Fatalf("Error creating ordered pipe: %v", err)
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

	// Start reading in a goroutine (simulating the upload goroutine)
	// Read will block waiting for contiguous data as writes come in
	readDone := make(chan error, 1)
	var sortedData []byte
	go func() {
		defer close(readDone)
		var err error
		sortedData, err = io.ReadAll(op.Out)
		if err != nil {
			readDone <- err
		}
	}()

	// Meanwhile, write data in the main goroutine (simulating FUSE writes)
	for _, w := range writeOffsets {
		_, err := op.WriteAt([]byte(w.Data), w.Offset)
		if err != nil {
			t.Fatalf("Error writing to sorted pipe: %v", err)
		}
	}

	// Signal that all writes are complete (simulates closeWriter())
	op.Close()

	// Wait for reading to complete and check for errors
	if err := <-readDone; err != nil {
		t.Fatalf("Error during reads: %v", err)
	}
	got := string(sortedData)
	want := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do"
	if got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}

func TestReaderAt(t *testing.T) {
	op, err := fsio.NewOrderedPipe("/test/path", fsio.WithLogger(&log.NoOpLogger{}))
	if err != nil {
		t.Fatalf("Error creating ordered pipe: %v", err)
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

	// Start reading in a goroutine (simulating the upload goroutine)
	readDone := make(chan error, 1)
	var sortedData []byte
	go func() {
		defer close(readDone)
		var err error
		sortedData, err = io.ReadAll(op.Out)
		if err != nil {
			readDone <- err
		}
	}()

	// Write data and test ReadAt during the process
	for idx, w := range writeOffsets {
		_, err := op.WriteAt([]byte(w.Data), w.Offset)
		if err != nil {
			t.Fatalf("Error writing to sorted pipe: %v", err)
		}
		if idx == 3 {
			// At this point, the data in the pipe should be:
			// "Lorem ipsum dolor sit " based on completing the first four writes.
			readBuff := make([]byte, 6)
			// Call Read with a start offset of 6, which should read "ipsum " and pass
			// a buffer of 6 bytes. The result should be "ipsum ".
			n := op.ReadAt(readBuff, 6)
			if n != 6 {
				t.Errorf("Expected to read 6 bytes, but got %d", n)
			}
			if string(readBuff) != "ipsum " {
				t.Errorf("Expected to read 'ipsum ', but got '%s'", string(readBuff))
			}
		}
		if idx == 7 {
			// At this point, the data in the pipe should be:
			// "Lorem ipsum dolor sit amet, consectetur "
			// based on completing the first eight writes.
			readBuff := make([]byte, 12)
			// Call Read with a start offset of 28, which should read "consectetur"...
			// and pass a buffer of 12 bytes. The result should be "consectetur ".
			n := op.ReadAt(readBuff, 28)
			if n != 12 {
				t.Errorf("Expected to read 12 bytes, but got %d", n)
			}
			if string(readBuff) != "consectetur " {
				t.Errorf("Expected to read 'consectetur ', but got '%s'", string(readBuff))
			}
		}
	}

	// Signal that all writes are complete (simulates closeWriter())
	op.Close()

	// Wait for reading to complete and check for errors
	if err := <-readDone; err != nil {
		t.Fatalf("Error during reads: %v", err)
	}
	got := string(sortedData)
	want := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do"
	if got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}
