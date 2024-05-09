package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	fs2 "io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/folder"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/test"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateClient(fixture string) (client *Client, r *recorder.Recorder, err error) {
	client = &Client{}
	client.Config, r, err = test.CreateConfig(fixture)

	return client, r, err
}

func deletePath(client *Client, path string) {
	err := client.Delete(files_sdk.FileDeleteParams{Path: path})
	var responseError files_sdk.ResponseError
	ok := errors.As(err, &responseError)
	if ok && responseError.Type == "not-found" {
	} else if ok && responseError.Type == "processing-failure/folder-not-empty" {
		err = client.Delete(files_sdk.FileDeleteParams{Path: path, Recursive: lib.Bool(true)})
		ok = errors.As(err, &responseError)
		if ok && responseError.Type == "not-found" {
			//noop
		} else if ok {
			panic(err)
		}
	} else if ok {
		panic(err)
	}
}

func ignoreSomeErrors(err error) {
	if err != nil && !files_sdk.IsExist(err) {
		panic(err)
	}
}

func buildScenario(base string, client *Client) {
	folderClient := folder.Client{Config: client.Config}

	_, err := folderClient.Create(files_sdk.FolderCreateParams{Path: base})
	ignoreSomeErrors(err)
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: lib.UrlJoinNoEscape(base, "nested_1")})
	ignoreSomeErrors(err)
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: lib.UrlJoinNoEscape(base, "nested_1", "nested_2")})
	ignoreSomeErrors(err)
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: lib.UrlJoinNoEscape(base, "nested_1", "nested_2", "nested_3")})
	ignoreSomeErrors(err)

	client.Upload()

	err = client.Upload(
		UploadWithSize(9),
		UploadWithDestinationPath(lib.UrlJoinNoEscape(base, "nested_1", "nested_2", "3.text")),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("testing 3")),
	)

	ignoreSomeErrors(err)

	err = client.Upload(
		UploadWithSize(9),
		UploadWithDestinationPath(lib.UrlJoinNoEscape(base, "nested_1", "nested_2", "nested_3", "4.text")),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("testing 3")),
	)
	ignoreSomeErrors(err)
}

func runDownloadScenario(path string, destination string, client *Client) map[string][]JobFile {
	m := &sync.Mutex{}
	results := make(map[string][]JobFile)

	reporter := func(r JobFile) {
		m.Lock()
		results[r.File.Path] = append(results[r.File.Path], r)
		m.Unlock()
	}

	job := client.Downloader(
		DownloaderParams{RemotePath: path, LocalPath: destination, EventsReporter: CreateReporter(reporter)},
	)

	job.Start()
	job.Wait()

	return results
}

func CreateReporter(callback Reporter) EventsReporter {
	return CreateFileEvents(callback, append(status.Excluded, append(status.Included, OnBytesChange(status.Uploading))...)...)
}

func TestClient_UploadFolder(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	resultsMapMutex := sync.RWMutex{}
	results := make(map[string][]ReporterCall)

	job := client.Uploader(
		UploaderParams{
			LocalPath:  "../lib",
			RemotePath: "golib",
			EventsReporter: CreateReporter(func(status JobFile) {
				resultsMapMutex.Lock()
				results[status.RemotePath] = append(results[status.RemotePath], ReporterCall{JobFile: status, err: status.Err})
				resultsMapMutex.Unlock()
			}),
			Manager: manager.Default(),
		},
	)

	job.Start()
	job.Wait()
	files, err := os.ReadDir("../lib")
	assert.NoError(err)
	gitIgnore, err := ignore.New()
	assert.NoError(err)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if gitIgnore.MatchesPath(f.Name()) {
			continue
		}
		remotePath := fmt.Sprintf("golib/lib/%v", f.Name())
		assert.Contains(results, remotePath)
		lastStatuses, ok := results[remotePath]
		if !ok {
			continue
		}
		lastStatus := lastStatuses[len(lastStatuses)-1]
		if lastStatus.Err != nil && strings.Contains(lastStatus.Err.Error(), "Requested interaction not found") {
			assert.Equal(status.Errored, lastStatus.Status)
		} else {
			assert.Equal(status.Complete, lastStatus.Status)
			assert.NoError(lastStatus.Err)
		}
	}
	deletePath(client, "golib")
}

func TestClient_UploadFolder_Dot(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_Dot")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	resultsMapMutex := sync.RWMutex{}
	results := make(map[string][]int64)
	err = os.MkdirAll(filepath.Join(tmpDir, "dot"), 0755)
	assert.NoError(err)

	f1, _ := os.Create(filepath.Join(tmpDir, "dot/1.text"))
	f1.Write([]byte("hello 1"))
	f1.Close()

	f2, _ := os.Create(filepath.Join(tmpDir, "dot/2.text"))
	f2.Write([]byte("hello 2"))
	f2.Close()

	f3, _ := os.Create(filepath.Join(tmpDir, "dot/3.text"))
	f3.Write([]byte("hello 3"))
	f3.Close()

	currentPwd, _ := os.Getwd()
	err = os.Chdir(filepath.Join(tmpDir, "dot"))
	defer func() {
		os.Chdir(currentPwd)
		os.RemoveAll(filepath.Join(tmpDir, "/dot"))
	}()
	assert.NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	job := client.Uploader(
		UploaderParams{
			LocalPath:  "." + string(os.PathSeparator),
			RemotePath: "go-from-dot",
			EventsReporter: CreateReporter(func(s JobFile) {
				resultsMapMutex.Lock()
				require.NoError(t, s.Err)

				results[s.File.Path] = append(results[s.File.Path], s.TransferBytes)
				resultsMapMutex.Unlock()
			}),
		}, files_sdk.WithContext(ctx))
	job.Start()
	job.Wait()
	assert.Contains(results, "go-from-dot/1.text")
	assert.Contains(results, "go-from-dot/2.text")
	assert.Contains(results, "go-from-dot/3.text")
	assert.Equal(int64(7), job.Statuses[0].TransferBytes())
	assert.Equal(int64(7), job.Statuses[1].TransferBytes())
	assert.Equal(int64(7), job.Statuses[2].TransferBytes())

	deletePath(client, "go-from-dot")
}

func TestClient_UploadFolder_Relative(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_Relative")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	resultsMapMutex := sync.RWMutex{}
	results := make(map[string][]int64)

	err = os.MkdirAll(filepath.Join(tmpDir, "relative"), 0755)
	assert.NoError(err)

	f1, _ := os.Create(filepath.Join(tmpDir, "relative", "1.text"))
	f1.Write([]byte("hello 1"))
	f1.Close()

	f2, _ := os.Create(filepath.Join(tmpDir, "relative", "2.text"))
	f2.Write([]byte("hello 2"))
	f2.Close()

	f3, _ := os.Create(filepath.Join(tmpDir, "relative", "3.text"))
	f3.Write([]byte("hello 3"))
	f3.Close()

	currentPwd, _ := os.Getwd()
	err = os.Chdir(tmpDir)
	defer os.Chdir(currentPwd)
	assert.NoError(err)

	job := client.Uploader(
		UploaderParams{
			LocalPath:  "relative" + string(os.PathSeparator),
			RemotePath: "relative",
			EventsReporter: CreateReporter(func(status JobFile) {
				resultsMapMutex.Lock()
				results[status.File.Path] = append(results[status.File.Path], status.TransferBytes)
				resultsMapMutex.Unlock()
			}),
		})
	job.Start()
	job.Wait()
	assert.Contains(results, "relative/1.text")
	assert.Contains(results, "relative/2.text")
	assert.Contains(results, "relative/3.text")
	assert.Equal(int64(7), job.Statuses[0].TransferBytes())
	assert.Equal(int64(7), job.Statuses[1].TransferBytes())
	assert.Equal(int64(7), job.Statuses[2].TransferBytes())
	assert.Equal(int64(21), job.TotalBytes(status.Valid...))
	assert.Equal(true, job.All(status.Ended...))

	deletePath(client, "relative")
}

func TestClient_Uploader(t *testing.T) {
	client, r, err := CreateClient("TestClient_Uploader")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	uploadPath := ".." + string(os.PathSeparator) + "LICENSE"
	job := client.Uploader(UploaderParams{LocalPath: uploadPath})
	job.Start()
	job.Wait()
	assert.Equal(true, job.Started.Called())
	assert.Equal(true, job.Scanning.Called())
	assert.Equal(true, job.EndScanning.Called())
	assert.Equal(true, job.Finished.Called())
	assert.Equal(false, job.Canceled.Called())
	assert.Equal("LICENSE", job.Files()[0].DisplayName)
	assert.Equal(1, job.Count())
	assert.Equal(int64(1102), job.TotalBytes())
	assert.Equal(true, job.All(status.Ended...))

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tempFile, err := os.CreateTemp(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "LICENSE"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal(file.DisplayName, "LICENSE")

	downloadData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		panic(err)
	}
	localData, err := os.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFile_To_Existing_Dir(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_To_Existing_Dir")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	folderClient := folder.Client{Config: client.Config}
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: "docs"})
	defer deletePath(client, "docs")

	assert.NoError(err)
	uploadPath := ".." + string(os.PathSeparator) + "LICENSE"
	job := client.Uploader(UploaderParams{LocalPath: uploadPath, RemotePath: "docs"})
	job.Start()
	job.Wait()
	tempFile, err := os.CreateTemp(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "docs/LICENSE"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal(file.DisplayName, "LICENSE")

	downloadData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		panic(err)
	}
	localData, err := os.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFile_To_NonExistingPath(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_To_NonExistingPath")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	deletePath(client, "taco")
	uploadPath := ".." + string(os.PathSeparator) + "LICENSE"
	job := client.Uploader(UploaderParams{LocalPath: uploadPath, RemotePath: "taco"})
	defer deletePath(client, "taco")
	job.Start()
	job.Wait()
	tempFile, _ := os.CreateTemp(tmpDir, "LICENSE")
	_, err = filepath.Abs(filepath.Dir(tempFile.Name()))
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "taco"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal("taco", file.DisplayName, "because the docs did not exist as a folder it becomes the file")

	downloadData, err := os.ReadFile(tempFile.Name())
	assert.NoError(err)
	localData, err := os.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFile_To_NonExistingPath_WithSlash(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_To_NonExistingPath_WithSlash")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	assert.NoError(err)
	uploadPath := ".." + string(os.PathSeparator) + "LICENSE"
	deletePath(client, "docs")
	job := client.Uploader(UploaderParams{LocalPath: uploadPath, RemotePath: "docs/"})
	defer deletePath(client, "docs")
	job.Start()
	job.Wait()
	tempFile, err := os.CreateTemp(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "docs/LICENSE"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal("file", file.Type)
	assert.Equal("LICENSE", file.DisplayName)

	downloadData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		panic(err)
	}
	localData, err := os.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFolder_as_file2(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_as_file2")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	uploadPath := ".." + string(os.PathSeparator) + "LICENSE"
	job := client.Uploader(UploaderParams{LocalPath: uploadPath})
	job.Start()
	job.Wait()

	assert.Equal(int64(1102), job.TransferBytes())
	assert.Equal(int64(1102), job.TotalBytes())

	tempFile, err := os.CreateTemp(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "LICENSE"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal(file.DisplayName, "LICENSE")

	downloadData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		panic(err)
	}
	localData, err := os.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFolder_RemotePathWithStartingSlash(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_RemotePathWithStartingSlash")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	err = os.MkdirAll(filepath.Join(tmpDir, "test"), 0755)
	assert.NoError(t, err)

	f1, _ := os.Create(filepath.Join(tmpDir, "test", "1.text"))
	f1.Write([]byte("hello 1"))
	f1.Close()

	f2, _ := os.Create(filepath.Join(tmpDir, "test", "2.text"))
	f2.Write([]byte("hello 2"))
	f2.Close()

	f3, _ := os.Create(filepath.Join(tmpDir, "test", "3.text"))
	f3.Write([]byte("hello 3"))
	f3.Close()
	job := client.Uploader(UploaderParams{LocalPath: filepath.Join(tmpDir, "test"), RemotePath: "/test", Manager: manager.Sync()})
	job.Start()
	job.Wait()
	assert.NoError(t, job.Statuses[0].Err())
	assert.Len(t, job.Statuses, 3)
	dir, _ := filepath.Split(job.Statuses[0].RemotePath())
	assert.Equal(t, "test/test/", dir)
}

func TestClient_UploadFolder_ZeroByteFile(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_ZeroByteFile")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	err = os.MkdirAll(filepath.Join(tmpDir, "zero_byte_folder"), 0755)
	assert.NoError(t, err)

	f1, _ := os.Create(filepath.Join(tmpDir, "zero_byte_folder", "zero-byte-file.text"))
	f1.Close()

	job := client.Uploader(UploaderParams{LocalPath: filepath.Join(tmpDir, "zero_byte_folder"), RemotePath: "", Manager: manager.Sync()})
	job.Start()
	job.Wait()
	require.Len(t, job.Statuses, 1)
	assert.NoError(t, job.Statuses[0].Err())
	assert.Equal(t, "zero_byte_folder/zero-byte-file.text", job.Statuses[0].RemotePath())

	job = client.Downloader(DownloaderParams{RemotePath: "zero_byte_folder", LocalPath: tmpDir + string(os.PathSeparator)})
	job.Start()
	job.Wait()
	require.Len(t, job.Statuses, 1)
	assert.NoError(t, job.Statuses[0].Err())
}

func TestClient_DownloadFolder(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	buildScenario("TestClient_DownloadFolder", client)

	assert := assert.New(t)

	it, err := client.ListFor(files_sdk.FolderListForParams{
		ListParams: files_sdk.ListParams{PerPage: 1},
		Path:       "TestClient_DownloadFolder/nested_1/nested_2",
	})

	assert.NoError(err)
	files := files_sdk.FileCollection{}
	for it.Next() {
		files = append(files, it.File())
	}

	assert.Len(files, 2, "something is wrong with cursor")

	results := runDownloadScenario("TestClient_DownloadFolder", "download/", client)
	assert.NoError(err)
	assert.Equal(int64(9), results["TestClient_DownloadFolder/nested_1/nested_2/3.text"][2].TransferBytes)
	assert.Equal(int64(9), results["TestClient_DownloadFolder/nested_1/nested_2/nested_3/4.text"][2].TransferBytes)
	os.RemoveAll("download")
}

func TestClient_DownloadFolder_Smart(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_Smart")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	defer os.RemoveAll("download")

	buildScenario("TestClient_DownloadFolder_Smart", client)

	results := runDownloadScenario(lib.UrlJoinNoEscape("TestClient_DownloadFolder_Smart", "nested_1", "nested_2", "3.text"), "download"+string(os.PathSeparator), client)
	for _, result := range results {
		for _, f := range result {
			require.NoError(t, f.Err)
		}
	}

	assert.Len(t, results["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"], 3)
	assert.Equal(t, int64(9), results["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"][2].TransferBytes)

	results2 := runDownloadScenario(lib.UrlJoinNoEscape("TestClient_DownloadFolder_Smart", "nested_1", "nested_2")+"/", "download", client)

	assert.NoError(t, err)

	path, err := os.Getwd()
	assert.NoError(t, err)

	require.Len(t, results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"], 3)
	assert.Equal(t, int64(9), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"][2].TransferBytes)
	assert.Equal(t, filepath.Join(path, "download", "3.text"), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"][2].LocalPath)
	fileInfo, err := os.Stat(filepath.Join(path, "download", "3.text"))
	require.NoError(t, err)
	assert.Equal(t, int64(9), fileInfo.Size())
	require.Len(t, results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/nested_3/4.text"], 3)
	assert.Equal(t, int64(9), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/nested_3/4.text"][2].TransferBytes)
	assert.Equal(t, filepath.Join(path, "download", "nested_3", "4.text"), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/nested_3/4.text"][2].LocalPath)
}

func TestClient_DownloadFolder_file_to_file(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_to_file")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	buildScenario("TestClient_DownloadFolder_file_to_file", client)
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	results := runDownloadScenario(lib.UrlJoinNoEscape("TestClient_DownloadFolder_file_to_file", "nested_1", "nested_2", "3.text"), filepath.Join(tmpDir, "3.text"), client)
	assert.NoError(err)
	for _, result := range results {
		for _, f := range result {
			require.NoError(t, f.Err)
		}
	}
	assert.Equal(int64(9), results["TestClient_DownloadFolder_file_to_file/nested_1/nested_2/3.text"][2].TransferBytes)
}

func TestClient_DownloadFolder_file_to_implicit(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_to_implicit")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	buildScenario("file_to_implicit", client)
	assert := assert.New(t)
	results := runDownloadScenario(lib.UrlJoinNoEscape("file_to_implicit", "nested_1", "nested_2", "3.text"), "", client)
	assert.NoError(err)
	for _, result := range results {
		for _, f := range result {
			require.NoError(t, f.Err)
		}
	}
	assert.Equal(int64(9), results["file_to_implicit/nested_1/nested_2/3.text"][2].TransferBytes)
	os.RemoveAll("3.text")
}

func TestClient_DownloadFolder_file_only(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_only")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	err = client.Upload(
		UploadWithSize(5),
		UploadWithDestinationPath("i am at the root.text"),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("hello")),
	)

	require.NoError(t, err)

	assert := assert.New(t)
	results := runDownloadScenario("i am at the root.text", "", client)
	assert.NoError(err)
	for _, result := range results {
		for _, f := range result {
			require.NoError(t, f.Err)
		}
	}
	require.Len(t, results["i am at the root.text"], 3)
	assert.Equal(int64(5), results["i am at the root.text"][2].TransferBytes)
	os.RemoveAll("i am at the root.text")
}

func TestClient_Downloader_Delete_Source(t *testing.T) {
	client, r, err := CreateClient("TestClient_Downloader_Delete_Source")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	folderClient := folder.Client{Config: client.Config}

	folderClient.Create(files_sdk.FolderCreateParams{Path: "test-delete-source"})

	err = client.Upload(
		UploadWithSize(9),
		UploadWithDestinationPath(lib.UrlJoinNoEscape("test-delete-source", "test.text")),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("testing 3")),
	)

	require.NoError(t, err)
	localPath, err := os.MkdirTemp("", "TestClient_Downloader_Delete_Source")
	require.NoError(t, err)

	job := client.Downloader(
		DownloaderParams{RemotePath: "test-delete-source", LocalPath: localPath},
	)
	var fi JobFile
	var log status.Log
	job.RegisterFileEvent(func(f JobFile) {
		fi = f
		log, err = DeleteSource{Config: client.Config, Direction: job.Direction}.Call(f)
	}, status.Complete)
	job.Start()
	<-job.Finished.C
	assert.NoError(err)
	assert.Equal("delete source", log.Action)
	assert.Equal(fi.RemotePath, log.Path)

	_, err = client.Find(files_sdk.FileFindParams{Path: lib.UrlJoinNoEscape("test-delete-source", "test.text")})
	require.NotNil(t, err)
	assert.Equal("Not Found - `Not Found`", err.Error())
	os.RemoveAll("test.text")
}

func TestClient_Downloader_Move_Source(t *testing.T) {
	client, r, err := CreateClient("TestClient_Downloader_Move_Source")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	folderClient := folder.Client{Config: client.Config}

	folderClient.Create(files_sdk.FolderCreateParams{Path: "test-move-source"})

	err = client.Upload(
		UploadWithSize(9),
		UploadWithDestinationPath(lib.UrlJoinNoEscape("test-move-source", "test.text")),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("testing 3")),
	)
	require.NoError(t, err)
	localPath, err := os.MkdirTemp("", "TestClient_Downloader_Move_Source")
	require.NoError(t, err)
	job := client.Downloader(
		DownloaderParams{RemotePath: "test-move-source", LocalPath: localPath},
	)
	var log status.Log
	job.RegisterFileEvent(func(f JobFile) {
		log, err = MoveSource{Config: client.Config, Direction: job.Direction, Path: "test-moved-source"}.Call(f)
	}, status.Complete)
	job.Start()
	job.Wait()

	require.NoError(t, err)
	assert.Equal("move source", log.Action)
	assert.Equal(filepath.Join("test-moved-source", "test.text"), log.Path)

	_, err = client.Find(files_sdk.FileFindParams{Path: lib.UrlJoinNoEscape("test-move-source", "test.text")})
	assert.Equal("Not Found - `Not Found`", err.Error())
	_, err = client.Find(files_sdk.FileFindParams{Path: lib.UrlJoinNoEscape("test-moved-source", "test.text")})
	assert.NoError(err)
	os.RemoveAll("test.text")
}

func TestClient_UploadFolder_Move_Source(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_Move_Source")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	var log status.Log

	defer os.RemoveAll(tmpDir)
	err = os.MkdirAll(filepath.Join(tmpDir, "move-source"), 0755)
	assert.NoError(err)

	tempFile, err := os.Create(filepath.Join(tmpDir, "move-source", "upload-move-source.text"))
	assert.NoError(err)
	tempFile.Write([]byte("testing"))
	require.NoError(t, tempFile.Close())
	job := client.Uploader(UploaderParams{LocalPath: tempFile.Name()})
	job.RegisterFileEvent(func(f JobFile) {
		fmt.Println("RegisterFileEvent")
		log, err = MoveSource{Config: client.Config, Direction: job.Direction, Path: filepath.Join(tmpDir, "move-source", "test-moved-source.text")}.Call(f)
	}, status.Complete)
	job.Start()
	job.Wait()
	assert.Equal(false, job.Any(status.Errored))
	assert.NoError(err)
	assert.Equal("move source", log.Action)
	assert.Equal(filepath.Join(tmpDir, "move-source", "test-moved-source.text"), log.Path)
	stat, err := os.Stat(log.Path)
	require.NoError(t, err)
	assert.Equal("test-moved-source.text", stat.Name())
	tempFile.Close()
}

func TestClient_UploadFolder_Move_Source_Missing_Dir(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_Move_Source_Missing_Dir")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	var log status.Log
	err = os.MkdirAll(filepath.Join(tmpDir, "move-source-dir"), 0755)
	assert.NoError(err)
	tempFile, err := os.Create(filepath.Join(tmpDir, "move-source-dir", "upload-move-source.text"))
	assert.NoError(err)
	tempFile.Write([]byte("testing"))
	tempFile.Close()
	job := client.Uploader(UploaderParams{LocalPath: filepath.Join(tmpDir, "move-source-dir") + string(os.PathSeparator)})
	job.RegisterFileEvent(func(f JobFile) {
		log, err = MoveSource{Config: client.Config, Direction: job.Direction, Path: filepath.Join(tmpDir, "moved-source-dir")}.Call(f)
	}, status.Complete)
	job.Start()
	job.Wait()
	erroredFile, ok := job.Find(status.Errored)
	if ok {
		assert.NoError(erroredFile.Err(), erroredFile.LocalPath())
	}
	require.Equal(t, false, job.Any(status.Errored), "")
	assert.NoError(err)
	assert.Equal("move source", log.Action)
	assert.Equal(filepath.Join(tmpDir, "moved-source-dir", "upload-move-source.text"), log.Path)
	stat, err := os.Stat(log.Path)
	assert.NoError(err)
	assert.Equal("upload-move-source.text", stat.Name())
	assert.Equal(false, stat.IsDir())

	_, err = os.Stat(filepath.Join(tmpDir, "move-source-dir", "upload-move-source.text"))
	assert.Equal(true, os.IsNotExist(err))
}

func TestClient_Downloader_Move_Source_Missing_Dir(t *testing.T) {
	client, r, err := CreateClient("TestClient_Downloader_Move_Source_Missing_Dir")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	buildScenario("TestClient_Downloader_Move_Source_Missing_Dir", client)
	os.MkdirAll(filepath.Join(tmpDir, "TestClient_Downloader_Move_Source_Missing_Dir"), 0755)

	logChan := make(chan status.Log)
	errChan := make(chan error)
	job := client.Downloader(
		DownloaderParams{
			RemotePath: "TestClient_Downloader_Move_Source_Missing_Dir",
			LocalPath:  filepath.Join(tmpDir, "TestClient_Downloader_Move_Source_Missing_Dir") + string(os.PathSeparator),
		},
	)
	job.RegisterFileEvent(func(f JobFile) {
		moveLog, err := MoveSource{Config: client.Config, Direction: job.Direction, Path: "TestClient_Downloader_Move_Source_Missing_Dir-moved"}.Call(f)
		logChan <- moveLog
		errChan <- err
	}, status.Complete)
	job.Start()

	fileValues := []string{
		"TestClient_Downloader_Move_Source_Missing_Dir-moved/nested_1/nested_2/3.text",
		"TestClient_Downloader_Move_Source_Missing_Dir-moved/nested_1/nested_2/nested_3/4.text",
	}
	files := make(map[string]string)
	files[fileValues[0]] = fileValues[0]
	files[fileValues[1]] = fileValues[1]

	log1 := <-logChan
	assert.Equal("move source", log1.Action)
	assert.Equal(files[lib.Path{Path: log1.Path}.NormalizePathSystemForAPI().String()], lib.Path{Path: log1.Path}.NormalizePathSystemForAPI().String())
	assert.NoError(<-errChan)

	log2 := <-logChan
	assert.Equal("move source", log2.Action)
	assert.Equal(files[lib.Path{Path: log2.Path}.NormalizePathSystemForAPI().String()], lib.Path{Path: log2.Path}.NormalizePathSystemForAPI().String())
	assert.NoError(<-errChan)

	job.Wait()

	assert.Equal(false, job.Any(status.Errored))

	movedDir, err := client.Find(files_sdk.FileFindParams{Path: "TestClient_Downloader_Move_Source_Missing_Dir-moved"})
	assert.NoError(err)
	assert.Equal("TestClient_Downloader_Move_Source_Missing_Dir-moved", movedDir.Path)
	assert.Equal("directory", movedDir.Type)
}

func TestClient_UploadFile_Delete_Source(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_Delete_Source")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	var log status.Log
	tempFile, err := os.Create(filepath.Join(tmpDir, "upload-delete-source.text"))
	assert.NoError(err)
	tempFile.Write([]byte("testing"))
	require.NoError(t, tempFile.Close())
	var fi JobFile
	job := client.Uploader(UploaderParams{LocalPath: tempFile.Name()})
	job.RegisterFileEvent(func(f JobFile) {
		fi = f
		log, err = DeleteSource{Config: client.Config, Direction: job.Direction}.Call(f)
	}, status.Complete)
	job.Start()
	job.Wait()
	assert.Equal(false, job.Any(status.Errored))
	assert.NoError(err)
	assert.Equal("delete source", log.Action)
	assert.Equal(fi.LocalPath, log.Path)
	tempFile.Close()
	os.Remove(tempFile.Name())
}

func TestClient_Uploader_Files(t *testing.T) {
	client, r, err := CreateClient("TestClient_Uploader_Files")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir := t.TempDir()

	filesAndStatus := []struct {
		name   string
		status string
		size   int
	}{{name: "1 (1).text", status: "complete", size: 24}, {name: "2.text", status: "complete", size: 24}, {name: "3.pdf", status: "ignored"}}
	var filePaths []string
	for _, file := range filesAndStatus {
		f, err := os.Create(filepath.Join(tmpDir, file.name))
		assert.NoError(err)
		f.Write([]byte("hello how are you doing?"))
		f.Close()
		if file.status == "complete" {
			filePaths = append(filePaths, f.Name())
		}
	}

	job := client.Uploader(UploaderParams{LocalPath: tmpDir + string(os.PathSeparator), LocalPaths: filePaths})
	job.Start()
	job.Wait()

	assert.Len(job.Statuses, 2)
	assert.Equal(filePaths[0], job.Statuses[0].LocalPath())
	assert.Equal(status.Complete, job.Statuses[0].Status())
	assert.NoError(job.Statuses[0].Err())

	assert.Equal(filePaths[1], job.Statuses[1].LocalPath())
	assert.Equal(status.Complete, job.Statuses[1].Status())
	assert.NoError(job.Statuses[1].Err())
}

func TestClient_Uploader_Directories(t *testing.T) {
	client, r, err := CreateClient("TestClient_Uploader_Directories")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir := t.TempDir()

	filesAndStatus := []struct {
		name   string
		status string
		size   int
	}{{name: "A/1.text", status: "complete", size: 24}, {name: "B/2.text", status: "complete", size: 24}, {name: "B/Z/4.text", status: "complete", size: 24}, {name: "3.text", status: "complete", size: 24}}
	for index, file := range filesAndStatus {
		file.name = filepath.Join(tmpDir, file.name)
		require.NoError(t, os.MkdirAll(filepath.Dir(file.name), 0750))
		filesAndStatus[index] = file
		f, err := os.Create(file.name)
		assert.NoError(err)
		f.Write([]byte("hello how are you doing?"))
		f.Close()
	}

	job := client.Uploader(
		UploaderParams{
			LocalPath: tmpDir + string(os.PathSeparator),
			LocalPaths: []string{
				filepath.Join(tmpDir, "A") + string(os.PathSeparator),
				filepath.Join(tmpDir, "B") + string(os.PathSeparator),
				filepath.Join(tmpDir, "3.text"),
			},
		},
	)
	job.Start()
	job.Wait()

	require.Equal(t, 4, len(job.Statuses), "the right number of files did not upload")
	assert.Equal(filesAndStatus[0].name, job.Statuses[0].LocalPath())
	assert.Equal("1.text", job.Statuses[0].RemotePath())
	assert.Equal(status.Complete, job.Statuses[0].Status())
	assert.NoError(job.Statuses[0].Err())

	assert.Equal(filesAndStatus[1].name, job.Statuses[1].LocalPath())
	assert.Equal("2.text", job.Statuses[1].RemotePath())
	assert.Equal(status.Complete, job.Statuses[1].Status())
	assert.NoError(job.Statuses[1].Err())

	assert.Equal(filesAndStatus[2].name, job.Statuses[2].LocalPath())
	assert.Equal("Z/4.text", job.Statuses[2].RemotePath())
	assert.Equal(status.Complete, job.Statuses[2].Status())
	assert.NoError(job.Statuses[2].Err())

	assert.Equal(filesAndStatus[3].name, job.Statuses[3].LocalPath())
	assert.Equal("3.text", job.Statuses[3].RemotePath())
	assert.Equal(status.Complete, job.Statuses[3].Status())
	assert.NoError(job.Statuses[3].Err())
}

func TestClient_ListForRecursive(t *testing.T) {
	client, r, err := CreateClient("TestClient_ListForRecursive")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	buildScenario("TestClient_ListForRecursive", client)

	it, _ := client.ListForRecursive(files_sdk.FolderListForParams{Path: "/TestClient_ListForRecursive"})
	var files []RecursiveItem
	for it.Next() {
		files = append(files, it.Resource())
	}

	paths := lo.Map[RecursiveItem, string](files, func(item RecursiveItem, index int) string {
		return item.Path
	})
	assert.Contains(paths, "TestClient_ListForRecursive")
	assert.Contains(paths, "TestClient_ListForRecursive/nested_1")
	assert.Contains(paths, "TestClient_ListForRecursive/nested_1/nested_2")
	assert.Contains(paths, "TestClient_ListForRecursive/nested_1/nested_2/nested_3")
	assert.Contains(paths, "TestClient_ListForRecursive/nested_1/nested_2/nested_3/4.text")
	assert.Contains(paths, "TestClient_ListForRecursive/nested_1/nested_2/3.text")
}

func TestClient_ListForRecursiveInsensitive(t *testing.T) {
	client, r, err := CreateClient("TestClient_ListForRecursiveInsensitive")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	buildScenario("TestClient_ListForRecursiveInsensitive", client)

	it, _ := client.ListForRecursive(files_sdk.FolderListForParams{Path: "/TestcLient_listforrecursiveinseNsitive"})
	var files []RecursiveItem
	for it.Next() {
		files = append(files, it.Resource())
	}

	require.Equal(t, 6, len(files))
	paths := lo.Map[RecursiveItem, string](files, func(item RecursiveItem, index int) string {
		return item.Path
	})
	assert.Contains(paths, "TestClient_ListForRecursiveInsensitive")
	assert.Contains(paths, "TestClient_ListForRecursiveInsensitive/nested_1")
	assert.Contains(paths, "TestClient_ListForRecursiveInsensitive/nested_1/nested_2")
	assert.Contains(paths, "TestClient_ListForRecursiveInsensitive/nested_1/nested_2/nested_3")
	assert.Contains(paths, "TestClient_ListForRecursiveInsensitive/nested_1/nested_2/nested_3/4.text")
	assert.Contains(paths, "TestClient_ListForRecursiveInsensitive/nested_1/nested_2/3.text")
}

func TestClient_ListForRecursive_Error(t *testing.T) {
	client, r, err := CreateClient("TestClient_ListForRecursive_Error")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	it, err := client.ListForRecursive(files_sdk.FolderListForParams{Path: "TestClient_ListForRecursive-Not-Found"})
	var files []interface{}
	if err == nil {
		for it.Next() {
			files = append(files, it.Current())
		}
	}

	assert.Equal(len(files), 0)
	assert.Equal("open TestClient_ListForRecursive-Not-Found: Authentication Required - `Unauthorized. The API key or Session token is either missing or invalid.`", err.Error())
}

func TestClient_ListForRecursive_Root(t *testing.T) {
	client, r, err := CreateClient("TestClient_ListForRecursive_Root")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	it, _ := client.ListForRecursive(files_sdk.FolderListForParams{Path: ""}, files_sdk.WithContext(ctx))
	recursiveItems := make([]RecursiveItem, 0)
	for it.Next() {
		if it.Err() != nil {
			assert.NoError(it.Err())
			continue
		}

		recursiveItems = append(recursiveItems, it.Resource())
		assert.NotEqual(it.Resource().Path, "")
	}
	paths := lo.Map[RecursiveItem, string](recursiveItems, func(item RecursiveItem, index int) string {
		return item.Path
	})
	assert.Len(paths, 4)
	errs := lo.Map[RecursiveItem, error](recursiveItems, func(item RecursiveItem, index int) error {
		return item.Err()
	})
	errs = lo.Reject[error](errs, func(err error, index int) bool {
		return err == nil
	})
	assert.ElementsMatch([]string{"azure", "aws-sftp"}, []string{errs[0].(*fs2.PathError).Path, errs[1].(*fs2.PathError).Path})
}

func TestClient_UploadFile(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile")
	if err != nil {
		t.Fatal(err)
	}
	fileName := filepath.Join(t.TempDir(), "anything")
	file, err := os.Create(fileName)
	require.NoError(t, err)
	file.Write([]byte("anything"))
	file.Close()

	err = client.UploadFile(fileName, "test/anything")
	require.NoError(t, err)

	sdkFile, err := client.Find(files_sdk.FileFindParams{Path: "test/anything"})
	require.NoError(t, err)
	require.Equal(t, sdkFile.DisplayName, "anything")
	fs := (&FS{}).Init(client.Config, false).WithContext(context.Background()).(fs2.FS)
	fsFile, err := fs.Open("test/anything")
	require.NoError(t, err)
	fileBytes, err := io.ReadAll(fsFile)
	require.NoError(t, err)
	require.Equal(t, "anything", string(fileBytes))
	r.Stop()
}

func TestClient_ListFor(t *testing.T) {
	client, r, err := CreateClient("TestClient_ListFor")
	if err != nil {
		t.Fatal(err)
	}

	buildScenario("TestClient_DownloadFolder", client)

	assert := assert.New(t)

	it, err := client.ListFor(files_sdk.FolderListForParams{
		ListParams: files_sdk.ListParams{PerPage: 1},
		Path:       "TestClient_DownloadFolder/nested_1/nested_2",
	})

	assert.NoError(err)
	var files []files_sdk.File
	for it.Next() {
		if it.File().Type == "directory" {
			subIt, _ := it.Iterate(it.File().Identifier())
			for subIt.Next() {
				subFile := subIt.Current().(files_sdk.File)
				files = append(files, subFile)
				loadedFile, err := it.LoadResource(subFile.Identifier())
				assert.NoError(err)
				assert.Equal(subFile, loadedFile, "LoadResource with Identifier matches file from list")
			}
		}
		files = append(files, it.File())

	}
	paths := lo.Map[files_sdk.File, string](files, func(item files_sdk.File, index int) string {
		return item.Path
	})
	assert.Equal([]string{
		"TestClient_DownloadFolder/nested_1/nested_2/nested_3/4.text",
		"TestClient_DownloadFolder/nested_1/nested_2/nested_3",
		"TestClient_DownloadFolder/nested_1/nested_2/3.text",
	}, paths)

	r.Stop()
}
