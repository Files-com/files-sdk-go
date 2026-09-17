package lib

import (
	"context"
	"errors"
	"io/fs"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWalkExitsWhenContextCanceledWhileWaitingForWorker(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	walk := (&Walk[string]{
		FS:                 fstest.MapFS{"file.txt": &fstest.MapFile{Data: []byte("fixture")}},
		ConcurrencyManager: NewConstrainedWorkGroup(1),
		WalkFile: func(_ fs.DirEntry, path string, _ error) (string, error) {
			close(entered)
			<-release
			return path, nil
		},
	}).Walk(ctx)

	<-entered
	done := make(chan bool, 1)
	go func() {
		done <- walk.Next()
	}()

	cancel()
	select {
	case ok := <-done:
		assert.False(t, ok)
	case <-time.After(250 * time.Millisecond):
		t.Fatal("Walk did not exit after parent context cancellation")
	}
	close(release)
}

func TestWalkSendReturnsWhenContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	iter := (&IterChan[string]{}).Init(context.Background())
	defer iter.Stop()

	walk := &Walk[string]{
		WalkFile: func(_ fs.DirEntry, path string, _ error) (string, error) {
			return path, nil
		},
	}

	err := walk.send(ctx, nil, "file.txt", iter, nil)

	assert.True(t, errors.Is(err, context.Canceled))
}

func TestWalkContinuesWhenQueueBecomesNonEmptyBeforeDoneWaitReturnsFalse(t *testing.T) {
	walkState := &Walk[string]{
		FS: fstest.MapFS{
			"root":       &fstest.MapFile{Mode: fs.ModeDir},
			"queued.txt": &fstest.MapFile{Data: []byte("queued")},
		},
		Root: "root",
		WalkFile: func(_ fs.DirEntry, path string, _ error) (string, error) {
			return path, nil
		},
	}
	walkState.ConcurrencyManager = &queueOnFalseManager{
		ConstrainedWorkGroup: NewConstrainedWorkGroup(1),
		queue:                &walkState.Queue,
		path:                 "queued.txt",
	}

	iter := walkState.Walk(context.Background())

	path, ok := nextWalkString(t, iter)
	assert.True(t, ok)
	assert.Equal(t, "queued.txt", path)

	_, ok = nextWalkString(t, iter)
	assert.False(t, ok)
}

func TestWalkDirReturnsContextCanceledWhenCanceledBeforeEntry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	iter := (&IterChan[string]{}).Init(context.Background())
	defer iter.Stop()
	called := false

	walk := &Walk[string]{
		FS: fstest.MapFS{"file.txt": &fstest.MapFile{Data: []byte("fixture")}},
		WalkFile: func(_ fs.DirEntry, path string, _ error) (string, error) {
			called = true
			return path, nil
		},
	}

	err := walk.walkDir(ctx, ".", iter)

	assert.True(t, errors.Is(err, context.Canceled))
	assert.False(t, called)
}

func TestWalkDeadlineExceededIsQuietShutdown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	entered := make(chan struct{})
	walk := (&Walk[string]{
		FS:                 fstest.MapFS{"file.txt": &fstest.MapFile{Data: []byte("fixture")}},
		ConcurrencyManager: NewConstrainedWorkGroup(1),
		WalkFile: func(_ fs.DirEntry, path string, _ error) (string, error) {
			close(entered)
			<-ctx.Done()
			return "", ctx.Err()
		},
	}).Walk(ctx)

	select {
	case <-entered:
	case <-time.After(250 * time.Millisecond):
		t.Fatal("WalkFile was not called")
	}
	_, ok := nextWalkString(t, walk)

	assert.False(t, ok)
	assert.NoError(t, walk.Err())
}

func TestWalkRootEmission(t *testing.T) {
	fileSystem := fstest.MapFS{
		"top.txt":       &fstest.MapFile{Data: []byte("top")},
		"dir/a.txt":     &fstest.MapFile{Data: []byte("a")},
		"dir/empty":     &fstest.MapFile{Mode: fs.ModeDir},
		"dir/sub/b.txt": &fstest.MapFile{Data: []byte("b")},
	}
	testCases := []struct {
		name            string
		root            string
		listDirectories bool
		expected        []string
	}{
		{
			name:            "empty root lists the contents without the root itself",
			root:            "",
			listDirectories: true,
			expected:        []string{"dir", "dir/a.txt", "dir/empty", "dir/sub", "dir/sub/b.txt", "top.txt"},
		},
		{
			name:            "dot root includes the root directory",
			root:            ".",
			listDirectories: true,
			expected:        []string{".", "dir", "dir/a.txt", "dir/empty", "dir/sub", "dir/sub/b.txt", "top.txt"},
		},
		{
			name:            "named root includes the root directory",
			root:            "dir",
			listDirectories: true,
			expected:        []string{"dir", "dir/a.txt", "dir/empty", "dir/sub", "dir/sub/b.txt"},
		},
		{
			name:            "empty root without directories lists only files",
			root:            "",
			listDirectories: false,
			expected:        []string{"dir/a.txt", "dir/sub/b.txt", "top.txt"},
		},
		{
			name:            "dot root without directories lists only files",
			root:            ".",
			listDirectories: false,
			expected:        []string{"dir/a.txt", "dir/sub/b.txt", "top.txt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iter := (&Walk[string]{
				FS:                 fileSystem,
				Root:               tc.root,
				ListDirectories:    tc.listDirectories,
				ConcurrencyManager: NewConstrainedWorkGroup(1),
				WalkFile: func(_ fs.DirEntry, path string, _ error) (string, error) {
					return path, nil
				},
			}).Walk(context.Background())

			var paths []string
			for {
				path, ok := nextWalkString(t, iter)
				if !ok {
					break
				}
				paths = append(paths, path)
			}

			assert.NoError(t, iter.Err())
			assert.ElementsMatch(t, tc.expected, paths)
		})
	}
}

type queueOnFalseManager struct {
	*ConstrainedWorkGroup
	queue *Queue[string]
	path  string
	once  sync.Once
}

func (m *queueOnFalseManager) NewSubWorker() ConcurrencyManager {
	return &queueOnFalseSubWorker{
		ConcurrencyManager: m.ConstrainedWorkGroup.NewSubWorker(),
		queue:              m.queue,
		path:               m.path,
		once:               &m.once,
	}
}

type queueOnFalseSubWorker struct {
	ConcurrencyManager
	queue *Queue[string]
	path  string
	once  *sync.Once
}

func (s *queueOnFalseSubWorker) WaitForADoneWithContext(ctx context.Context) bool {
	if s.ConcurrencyManager.WaitForADoneWithContext(ctx) {
		return true
	}
	s.once.Do(func() {
		s.queue.Push(s.path)
	})
	return false
}

func nextWalkString(t *testing.T, iter *IterChan[string]) (string, bool) {
	t.Helper()
	type result struct {
		path string
		ok   bool
	}
	done := make(chan result, 1)
	go func() {
		ok := iter.Next()
		path := ""
		if ok {
			path = iter.Resource()
		}
		done <- result{path: path, ok: ok}
	}()

	select {
	case result := <-done:
		return result.path, result.ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("iterator did not return")
		return "", false
	}
}
