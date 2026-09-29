package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/test"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createClient(fixture string) (client *Client, r *recorder.Recorder, err error) {
	client = &Client{}
	client.Config, r, err = test.CreateConfig(fixture)

	return client, r, err
}

func TestClient_Create(t *testing.T) {
	client, r, err := createClient("TestClient_Create")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)

	user, err := findOrCreateUser(client, files_sdk.UserCreateParams{Username: "TestMo"})
	assert.NoError(err)

	assert.Equal("TestMo", user.Username)
}

func TestClient_Update(t *testing.T) {
	client, r, err := createClient("TestClient_Update")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)

	user, err := findOrCreateUser(client, files_sdk.UserCreateParams{Username: "TestMo"})
	assert.NoError(err)

	assert.Equal(lib.Bool(true), user.SftpPermission)

	user, err = client.Update(
		files_sdk.UserUpdateParams{
			Id:             user.Id,
			SftpPermission: lib.Bool(false),
		},
	)

	assert.NoError(err)

	assert.Equal(lib.Bool(false), user.SftpPermission)
}

func TestClient_List(t *testing.T) {
	client, r, err := createClient("TestClient_List")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)

	_, err = findOrCreateUser(client, files_sdk.UserCreateParams{Username: "test-list-user"})
	assert.NoError(err)

	it, err := client.List(files_sdk.UserListParams{})
	assert.NoError(err)
	var users []files_sdk.User
	for it.Next() {
		users = append(users, it.User())
		loaderUser, err := it.LoadResource(it.User().Identifier())
		assert.NoError(err)
		assert.Equal(loaderUser, it.User())
	}
	assert.NoError(it.Err())
	assert.Len(users, 1)
}

func TestClient_List_ReloadRepeatsListingWithItsOwnCursor(t *testing.T) {
	pages := map[string]struct {
		ids  []int64
		next string
	}{
		"":       {[]int64{1, 2}, "page-2"},
		"page-2": {[]int64{3}, "page-3"},
		"page-3": {[]int64{4}, ""},
	}
	var mu sync.Mutex
	var requests []*http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests = append(requests, r.Clone(context.Background()))
		mu.Unlock()
		page := pages[r.URL.Query().Get("cursor")]
		writeUsers(t, w, page.next, page.ids...)
	}))
	defer server.Close()
	client := newListClient(server.URL)

	source, err := client.List(
		files_sdk.UserListParams{
			Search:     "ops",
			Filter:     map[string]interface{}{"not_site_admin": "true"},
			ListParams: files_sdk.ListParams{PerPage: 2},
		},
		headerOption("X-Test-Listing", "users"),
		headerOption("X-Test-Traversal", "source"),
	)
	require.NoError(t, err)
	// Configuration made before the first page belongs to the listing.
	source.GetParams().MaxPages = 2
	assert.Equal(t, int64(1), nextUserID(t, source))

	reloaded, ok := source.Reload(headerOption("X-Test-Traversal", "reloaded")).(*Iter)
	require.True(t, ok, "Reload keeps the generated iterator type")
	assert.Same(t, client, reloaded.Client)

	assert.Equal(t, int64(1), nextUserID(t, reloaded))
	assert.Equal(t, int64(2), nextUserID(t, source))
	assert.Equal(t, int64(3), nextUserID(t, source))
	assert.Equal(t, int64(2), nextUserID(t, reloaded))
	assert.Equal(t, int64(3), nextUserID(t, reloaded))
	assert.False(t, source.Next())
	assert.False(t, reloaded.Next())

	assert.Equal(t, []int64{1, 2, 3}, userIDs(reloaded.Reload()))

	type traversalRequest struct{ Traversal, Cursor string }
	var seen []traversalRequest
	mu.Lock()
	defer mu.Unlock()
	for _, r := range requests {
		query := r.URL.Query()
		seen = append(seen, traversalRequest{r.Header.Get("X-Test-Traversal"), query.Get("cursor")})
		assert.Equal(t, "ops", query.Get("search"))
		assert.Equal(t, "true", query.Get("filter[not_site_admin]"))
		assert.Equal(t, "2", query.Get("per_page"))
		assert.Equal(t, "list-key", r.Header.Get("X-FilesAPI-Key"))
		assert.Equal(t, "42", r.Header.Get("X-Files-Workspace-Id"))
		assert.Equal(t, "users", r.Header.Get("X-Test-Listing"))
	}
	assert.Equal(t, []traversalRequest{
		{"source", ""},
		{"reloaded", ""},
		{"source", "page-2"},
		{"reloaded", "page-2"},
		{"reloaded", ""},
		{"reloaded", "page-2"},
	}, seen)
}

func TestClient_List_ReloadContextReplacesListingContext(t *testing.T) {
	var servedID atomic.Int64
	servedID.Store(1)
	heldRequestArrived := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Hold") != "" {
			close(heldRequestArrived)
			<-r.Context().Done() // The client abandoned the request.
			return
		}
		writeUsers(t, w, "", servedID.Load())
	}))
	defer server.Close()

	listCtx, cancelList := context.WithCancel(context.Background())
	it, err := newListClient(server.URL).List(files_sdk.UserListParams{}, files_sdk.WithContext(listCtx))
	require.NoError(t, err)
	assert.Equal(t, []int64{1}, userIDs(it))
	cancelList()

	stale := it.Reload()
	assert.Empty(t, userIDs(stale))
	assert.ErrorIs(t, stale.Err(), context.Canceled, "without a new context, the listing's context still applies")

	servedID.Store(2)
	live := it.Reload(files_sdk.WithContext(context.Background()))
	assert.Equal(t, []int64{2}, userIDs(live))
	assert.NoError(t, live.Err())

	reloadCtx, cancelReload := context.WithCancel(context.Background())
	held := it.Reload(files_sdk.WithContext(reloadCtx), headerOption("X-Test-Hold", "true"))
	go func() {
		<-heldRequestArrived
		cancelReload()
	}()
	assert.Empty(t, userIDs(held))
	assert.ErrorIs(t, held.Err(), context.Canceled)
}

func TestClient_List_ReloadWhileCanceledRequestIsReturning(t *testing.T) {
	var requests atomic.Int32
	firstRequestArrived := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "ops", r.URL.Query().Get("search"))
		if requests.Add(1) == 1 {
			close(firstRequestArrived)
			<-r.Context().Done() // The client abandoned the request.
			return
		}
		writeUsers(t, w, "", 2)
	}))
	defer server.Close()

	listCtx, cancelList := context.WithCancel(context.Background())
	source, err := newListClient(server.URL).List(files_sdk.UserListParams{Search: "ops"}, files_sdk.WithContext(listCtx))
	require.NoError(t, err)
	sourceNext := make(chan bool)
	go func() { sourceNext <- source.Next() }()

	// Like the CLI's Backspace while a list loads: cancel the traversal and
	// reload it without waiting for the canceled Next to return.
	<-firstRequestArrived
	cancelList()
	reloaded := source.Reload(files_sdk.WithContext(context.Background()))

	assert.False(t, <-sourceNext)
	assert.ErrorIs(t, source.Err(), context.Canceled)
	assert.Equal(t, []int64{2}, userIDs(reloaded))
	assert.NoError(t, reloaded.Err())
}

func TestClient_List_OverlappingReloadsKeepTheirOwnContexts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeUsers(t, w, "", 1)
	}))
	defer server.Close()

	// The first request waits while its options are applied until another
	// reload's request has finished, so the two requests overlap.
	firstRequestPaused := make(chan struct{})
	resumeFirstRequest := make(chan struct{})
	var paused atomic.Bool
	pauseFirstRequest := files_sdk.RequestOption(func(*http.Request) error {
		if paused.CompareAndSwap(false, true) {
			close(firstRequestPaused)
			<-resumeFirstRequest
		}
		return nil
	})
	listOpts := make([]files_sdk.RequestResponseOption, 0, 2) // spare capacity
	listOpts = append(listOpts, pauseFirstRequest)
	it, err := newListClient(server.URL).List(files_sdk.UserListParams{}, listOpts...)
	require.NoError(t, err)

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	live := it.Reload(files_sdk.WithContext(context.Background()))
	canceled := it.Reload(files_sdk.WithContext(canceledCtx))

	liveIDs := make(chan []int64)
	go func() { liveIDs <- userIDs(live) }()
	<-firstRequestPaused
	assert.Empty(t, userIDs(canceled))
	assert.ErrorIs(t, canceled.Err(), context.Canceled)
	close(resumeFirstRequest)

	assert.Equal(t, []int64{1}, <-liveIDs)
	assert.NoError(t, live.Err())
}

func TestClient_List_SourceAndReloadCanRequestConcurrently(t *testing.T) {
	var arrivals atomic.Int32
	bothArrived := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reply once both traversals are waiting, so their responses decode together.
		if arrivals.Add(1) == 2 {
			close(bothArrived)
		}
		select {
		case <-bothArrived:
			writeUsers(t, w, "", 1, 2)
		case <-time.After(10 * time.Second):
			http.Error(w, "the other traversal did not send a request", http.StatusRequestTimeout)
		}
	}))
	defer server.Close()

	source, err := newListClient(server.URL).List(files_sdk.UserListParams{})
	require.NoError(t, err)
	traversals := []files_sdk.IterI{source, source.Reload()}

	ids := make([][]int64, len(traversals))
	var done sync.WaitGroup
	for n, traversal := range traversals {
		done.Add(1)
		go func() {
			defer done.Done()
			ids[n] = userIDs(traversal)
		}()
	}
	done.Wait()

	for n, traversal := range traversals {
		assert.Equal(t, []int64{1, 2}, ids[n])
		assert.NoError(t, traversal.Err())
	}
}

func newListClient(serverURL string) *Client {
	return &Client{Config: files_sdk.Config{APIKey: "list-key", WorkspaceId: 42, EndpointOverride: serverURL}.Init()}
}

func headerOption(name, value string) files_sdk.RequestResponseOption {
	return files_sdk.RequestHeadersOption(&http.Header{name: {value}})
}

func writeUsers(t *testing.T, w http.ResponseWriter, nextCursor string, ids ...int64) {
	users := make([]map[string]int64, len(ids))
	for n, id := range ids {
		users[n] = map[string]int64{"id": id}
	}
	w.Header().Set("Content-Type", "application/json")
	if nextCursor != "" {
		w.Header().Set("X-Files-Cursor", nextCursor)
	}
	assert.NoError(t, json.NewEncoder(w).Encode(users))
}

func nextUserID(t *testing.T, it *Iter) int64 {
	t.Helper()
	require.True(t, it.Next(), "expected another user, err: %v", it.Err())
	return it.User().Id
}

func userIDs(it files_sdk.IterI) (ids []int64) {
	for it.Next() {
		ids = append(ids, it.Current().(files_sdk.User).Id)
	}
	return ids
}

func findOrCreateUser(client *Client, params files_sdk.UserCreateParams) (files_sdk.User, error) {
	user, err := findUser(client, params)
	if err != nil && err.Error() == "user not found" {
		return client.Create(
			params,
		)
	}
	return user, err
}

func findUser(client *Client, params files_sdk.UserCreateParams) (files_sdk.User, error) {
	it, err := client.List(
		files_sdk.UserListParams{},
	)

	if err != nil {
		return files_sdk.User{}, err
	}
	var user *files_sdk.User
	for it.Next() {
		if it.User().Username == params.Username {
			u := it.User()
			user = &u
			continue
		}
	}
	if user == nil {
		return files_sdk.User{}, fmt.Errorf("user not found")
	}

	return *user, nil
}
