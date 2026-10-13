package manager

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/http/httptrace"
	"sync"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

func TestAdaptiveUploadV2ConcurrentFilePartsOverrideCapsDefault(t *testing.T) {
	t.Cleanup(func() {
		SetAdaptiveUploadV2ConcurrentFileParts(0)
	})

	SetAdaptiveUploadV2ConcurrentFileParts(7)

	if got := EffectiveAdaptiveUploadV2ConcurrentFileParts(1024); got != 7 {
		t.Fatalf("EffectiveAdaptiveUploadV2ConcurrentFileParts(1024) = %d, want 7", got)
	}
	if got := EffectiveAdaptiveUploadV2ConcurrentFileParts(4); got != 4 {
		t.Fatalf("EffectiveAdaptiveUploadV2ConcurrentFileParts(4) = %d, want 4", got)
	}
}

func TestAdaptiveUploadV2ConcurrentFilePartsOverrideCanReset(t *testing.T) {
	t.Cleanup(func() {
		SetAdaptiveUploadV2ConcurrentFileParts(0)
	})

	SetAdaptiveUploadV2ConcurrentFileParts(7)
	SetAdaptiveUploadV2ConcurrentFileParts(0)

	if got := EffectiveAdaptiveUploadV2ConcurrentFileParts(1024); got != 1024 {
		t.Fatalf("EffectiveAdaptiveUploadV2ConcurrentFileParts(1024) = %d, want 1024", got)
	}
}

func TestAdaptiveDownloadV2ConcurrentFilePartsOverrideCapsDefault(t *testing.T) {
	t.Cleanup(func() {
		SetAdaptiveDownloadV2ConcurrentFileParts(0)
	})

	SetAdaptiveDownloadV2ConcurrentFileParts(7)

	if got := EffectiveAdaptiveDownloadV2ConcurrentFileParts(); got != 7 {
		t.Fatalf("EffectiveAdaptiveDownloadV2ConcurrentFileParts() = %d, want 7", got)
	}
}

func TestAdaptiveDownloadV2ConcurrentFilePartsOverrideCanReset(t *testing.T) {
	t.Cleanup(func() {
		SetAdaptiveDownloadV2ConcurrentFileParts(0)
	})

	SetAdaptiveDownloadV2ConcurrentFileParts(7)
	SetAdaptiveDownloadV2ConcurrentFileParts(0)

	if got := EffectiveAdaptiveDownloadV2ConcurrentFileParts(); got != AdaptiveDownloadV2ConcurrentFileParts {
		t.Fatalf("EffectiveAdaptiveDownloadV2ConcurrentFileParts() = %d, want %d", got, AdaptiveDownloadV2ConcurrentFileParts)
	}
}

func TestManagerTransferStats(t *testing.T) {
	manager := New(2, 3, 1)
	manager.FilePartsManager.Wait()
	manager.FilePartsManager.Wait()
	t.Cleanup(func() {
		manager.FilePartsManager.Done()
		manager.FilePartsManager.Done()
	})

	stats := manager.TransferStats()

	if stats.Active != 2 {
		t.Fatalf("TransferStats().Active = %d, want 2", stats.Active)
	}
	if stats.Max != 3 {
		t.Fatalf("TransferStats().Max = %d, want 3", stats.Max)
	}
}

func TestNilManagerTransferStats(t *testing.T) {
	var manager *Manager

	stats := manager.TransferStats()

	if stats.Active != 0 {
		t.Fatalf("TransferStats().Active = %d, want 0", stats.Active)
	}
	if stats.Max != 0 {
		t.Fatalf("TransferStats().Max = %d, want 0", stats.Max)
	}
}

// A matched client never writes the client it is made from, whose transport other requests may be using, and keeps
// its own cap: making a client with another cap from the same original, while requests run on both, does not raise it.
// Each client reuses its own connections once their responses are read.
func TestCreateMatchingClientLeavesTheOriginalInUseAndKeepsEachCap(t *testing.T) {
	held := []string{"original", "cap1-first", "cap1-second", "cap2-first", "cap2-second"}
	gates := make(map[string]chan struct{}, len(held))
	for _, id := range held {
		gates[id] = make(chan struct{})
	}
	done := make(chan struct{})
	arrivals := make(chan string, 16)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		arrivals <- id
		if gate, ok := gates[id]; ok {
			select {
			case <-gate:
			case <-done:
			}
		}
		_, _ = io.WriteString(w, id)
	}))

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	errRedirect := errors.New("redirect checked")
	transport := lib.DefaultPooledTransport()
	original := &http.Client{
		Transport:     transport,
		Timeout:       time.Minute,
		Jar:           jar,
		CheckRedirect: func(*http.Request, []*http.Request) error { return errRedirect },
	}
	var cap1, cap2 *http.Client
	var requests sync.WaitGroup
	t.Cleanup(func() {
		close(done)
		requests.Wait()
		for _, client := range []*http.Client{original, cap1, cap2} {
			if client != nil {
				client.CloseIdleConnections()
			}
		}
		server.Close()
	})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	// get sends the request id with client and reads its whole response, reporting whether the connection was reused.
	get := func(client *http.Client, id string) (reused bool, err error) {
		trace := &httptrace.ClientTrace{GotConn: func(info httptrace.GotConnInfo) { reused = info.Reused }}
		req, err := http.NewRequestWithContext(httptrace.WithClientTrace(ctx, trace), http.MethodGet, server.URL+"/?id="+id, nil)
		if err != nil {
			return false, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		_, err = io.Copy(io.Discard, resp.Body)
		return reused, err
	}
	start := func(client *http.Client, id string) <-chan struct{} {
		finished := make(chan struct{})
		requests.Add(1)
		go func() {
			defer requests.Done()
			defer close(finished)
			if _, err := get(client, id); err != nil {
				t.Errorf("request %s: %v", id, err)
			}
		}()
		return finished
	}
	arrived := map[string]bool{}
	await := func(ids ...string) {
		t.Helper()
		for _, id := range ids {
			for !arrived[id] {
				select {
				case got := <-arrivals:
					arrived[got] = true
				case <-ctx.Done():
					t.Fatalf("request %s never reached the server", id)
				}
			}
		}
	}

	// The matched clients are made while a request is running on the original.
	originalDone := start(original, "original")
	await("original")
	cap1 = Build(1, 1).CreateMatchingClient(original)
	cap2 = Build(2, 1).CreateMatchingClient(original)

	if original.Transport != transport || transport.MaxConnsPerHost != 75 {
		t.Fatalf("the original client's transport changed: %T allowing %d connections per host", original.Transport, transport.MaxConnsPerHost)
	}
	for _, matched := range []struct {
		client *http.Client
		cap    int
	}{{cap1, 1}, {cap2, 2}} {
		matchedTransport, ok := matched.client.Transport.(*lib.Transport)
		if !ok || matched.client == original || matchedTransport == transport {
			t.Fatalf("cap %d: got the original client or a %T transport, want a new SDK transport", matched.cap, matched.client.Transport)
		}
		if matchedTransport.MaxConnsPerHost != matched.cap {
			t.Fatalf("cap %d: the matched transport allows %d connections per host", matched.cap, matchedTransport.MaxConnsPerHost)
		}
		if matchedTransport.MaxIdleConns != transport.MaxIdleConns || matchedTransport.MaxIdleConnsPerHost != transport.MaxIdleConnsPerHost || matchedTransport.IdleConnTimeout != transport.IdleConnTimeout {
			t.Fatalf("cap %d: idle limits %d/%d/%v, want the original's %d/%d/%v", matched.cap, matchedTransport.MaxIdleConns, matchedTransport.MaxIdleConnsPerHost, matchedTransport.IdleConnTimeout, transport.MaxIdleConns, transport.MaxIdleConnsPerHost, transport.IdleConnTimeout)
		}
		if matched.client.Timeout != original.Timeout || matched.client.Jar != original.Jar || matched.client.CheckRedirect == nil || matched.client.CheckRedirect(nil, nil) != errRedirect {
			t.Fatalf("cap %d: the client's timeout, cookie jar or redirect policy was not kept", matched.cap)
		}
	}
	if again := Build(1, 1).CreateMatchingClient(cap1); again != cap1 {
		t.Fatal("a client whose cap already matches was not returned as it is")
	}

	// cap1's second request waits for its first, even once cap2's two requests are running.
	cap1First := start(cap1, "cap1-first")
	await("cap1-first")
	cap1Second := start(cap1, "cap1-second")
	cap2First := start(cap2, "cap2-first")
	cap2Second := start(cap2, "cap2-second")
	await("cap2-first", "cap2-second")
	for drained := false; !drained; {
		select {
		case got := <-arrivals:
			arrived[got] = true
		default:
			drained = true
		}
	}
	if arrived["cap1-second"] {
		t.Fatal("the client matched to one connection per host sent two requests at once")
	}
	close(gates["cap1-first"])
	await("cap1-second")
	for _, id := range []string{"original", "cap1-second", "cap2-first", "cap2-second"} {
		close(gates[id])
	}
	for _, finished := range []<-chan struct{}{originalDone, cap1First, cap1Second, cap2First, cap2Second} {
		select {
		case <-finished:
		case <-ctx.Done():
			t.Fatal("a request never finished")
		}
	}

	for _, client := range []struct {
		name   string
		client *http.Client
	}{{"original", original}, {"cap1", cap1}} {
		reused, err := get(client.client, client.name+"-again")
		if err != nil {
			t.Fatal(err)
		}
		if !reused {
			t.Fatalf("%s did not reuse its own connection", client.name)
		}
	}
}
