package file

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Files-com/files-sdk-go/v2/ignore"

	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib/test"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/folder"
	"github.com/Files-com/files-sdk-go/v2/lib"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func CreateClient(fixture string) (client *Client, r *recorder.Recorder, err error) {
	client = &Client{}
	client.Config, r, err = test.CreateConfig(fixture)

	return client, r, err
}

func deletePath(client *Client, path string) {
	err := client.Delete(context.Background(), files_sdk.FileDeleteParams{Path: path})
	responseError, ok := err.(files_sdk.ResponseError)
	if ok && responseError.Type == "not-found" {
	} else if ok && responseError.Type == "processing-failure/folder-not-empty" {
		err = client.Delete(context.Background(), files_sdk.FileDeleteParams{Path: path, Recursive: lib.Bool(true)})
		responseError, ok = err.(files_sdk.ResponseError)
		if ok && responseError.Type == "not-found" {

		} else if ok {
			panic(err)
		}
	} else if ok {
		panic(err)
	}
}

func ignoreSomeErrors(err error) {
	if err != nil && !files_sdk.IsDestinationExistsError(err) {
		panic(err)
	}
}

func buildScenario(base string, client *Client) {
	folderClient := folder.Client{Config: client.Config}

	_, err := folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: base})
	ignoreSomeErrors(err)
	_, err = folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: filepath.Join(base, "nested_1")})
	ignoreSomeErrors(err)
	_, err = folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: filepath.Join(base, "nested_1", "nested_2")})
	ignoreSomeErrors(err)
	_, err = folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: filepath.Join(base, "nested_1", "nested_2", "nested_3")})
	ignoreSomeErrors(err)

	_, _, _, _, err = client.UploadIO(
		context.Background(),
		UploadIOParams{
			Path:   filepath.Join(base, "nested_1", "nested_2", "3.text"),
			Reader: strings.NewReader("testing 3"), Size: int64(9),
			ProvidedMtime: time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	)
	ignoreSomeErrors(err)
	_, _, _, _, err = client.UploadIO(
		context.Background(),
		UploadIOParams{
			Path:   filepath.Join(base, "nested_1", "nested_2", "nested_3", "4.text"),
			Reader: strings.NewReader("testing 3"), Size: int64(9),
			ProvidedMtime: time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	)
	ignoreSomeErrors(err)
}

func runDownloadScenario(path string, destination string, client *Client) map[string][]status.File {
	m := &sync.Mutex{}
	results := make(map[string][]status.File)

	reporter := func(r status.File) {
		m.Lock()
		results[r.File.Path] = append(results[r.File.Path], r)
		m.Unlock()
	}

	job := client.Downloader(
		context.Background(),
		DownloaderParams{RemotePath: path, LocalPath: destination, EventsReporter: Reporter(reporter)},
	)

	job.Start()
	job.Wait()

	return results
}

func Reporter(callback status.Reporter) status.EventsReporter {
	events := make(status.EventsReporter)

	for _, s := range status.Included {
		events[s] = append(events[s], callback)
	}

	for _, s := range status.Excluded {
		events[s] = append(events[s], callback)
	}

	return events
}

func TestClient_UploadFolder_path_spec(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_path_spec")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	t.Run("files", func(t *testing.T) {
		for _, tt := range pathSpec() {
			t.Run(tt.name, func(t *testing.T) {
				filesDest := fmt.Sprintf(
					"files-dest-%v",
					strings.Replace(
						strings.Replace(tt.name, " ", "_", -1),
						"/", "_", -1,
					),
				)
				os.MkdirAll(filesDest, 0755)
				folderClient := folder.Client{Config: client.Config}
				_, err := folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: filesDest})
				ignoreSomeErrors(err)
				localSrc, err := os.MkdirTemp("", "local-src")
				assert.NoError(t, err)

				for _, e := range tt.src {
					if e.dir {
						err := os.MkdirAll(filepath.Join(localSrc, e.path), 0750)
						assert.NoError(t, err)
					} else {
						file, err := os.Create(filepath.Join(localSrc, e.path))
						file.Write([]byte(e.path))
						assert.NoError(t, err)
					}
				}
				for _, e := range tt.dest {
					if !e.preexisting {
						continue
					}
					if e.dir {
						_, err := folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: filepath.Join(filesDest, e.path)})
						ignoreSomeErrors(err)
					} else {
						_, _, _, _, err = client.UploadIO(
							context.Background(),
							UploadIOParams{
								Path:   filepath.Join(filesDest, e.path),
								Reader: strings.NewReader(e.path), Size: int64(len(e.path)),
								ProvidedMtime: time.Now(),
							},
						)
						assert.NoError(t, err)
					}
				}
				params := UploaderParams{
					LocalPath:  strings.Join([]string{localSrc, tt.args.src}, string(os.PathSeparator)),
					RemotePath: filepath.Join(filesDest, tt.args.dest),
				}
				if tt.args.dest == "" {
					params.RemotePath = ""
				}
				job := client.Uploader(context.Background(), params)

				t.Logf("RemotePath: %v, LocalPath: %v", params.RemotePath, params.LocalPath)

				job.Start()
				job.Wait()
				assert.Len(t, job.Statuses, 1)
				assert.NoError(t, job.Statuses[0].Err())
				it, err := client.ListForRecursive(context.Background(), files_sdk.FolderListForParams{Path: filesDest})
				assert.NoError(t, err)
				for it.Next() {
					fmt.Println(it.Current().(files_sdk.File).Path)
					assert.NoError(t, it.Err())
				}
				var filePath string
				for _, e := range tt.dest {
					if tt.args.dest == "" {
						filePath = e.path
					} else {
						filePath = filepath.Join(filesDest, e.path)
					}
					file, err := client.Find(context.Background(), files_sdk.FileFindParams{Path: filePath})
					assert.NoError(t, err)
					assert.Equal(t, e.dir, file.Type == "directory", e.path)
				}

				assert.NoError(t, os.RemoveAll(localSrc))
				assert.NoError(t, client.Delete(context.Background(), files_sdk.FileDeleteParams{Path: filesDest, Recursive: lib.Bool(true)}))
				if tt.args.dest == "" {
					assert.NoError(t, client.Delete(context.Background(), files_sdk.FileDeleteParams{Path: filePath, Recursive: lib.Bool(true)}))
				}
			})
		}
	})
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
		context.Background(),
		UploaderParams{
			LocalPath:  "../lib",
			RemotePath: "golib",
			EventsReporter: Reporter(func(status status.File) {
				resultsMapMutex.Lock()
				results[status.RemotePath] = append(results[status.RemotePath], ReporterCall{File: status, err: status.Err})
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
		lastStatuses, ok := results[fmt.Sprintf("golib/%v", f.Name())]
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

		assert.Contains(results, fmt.Sprintf("golib/%v", f.Name()))
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

	f1, err := os.Create(filepath.Join(tmpDir, "dot/1.text"))
	f1.Write([]byte("hello 1"))
	f1.Close()

	f2, err := os.Create(filepath.Join(tmpDir, "dot/2.text"))
	f2.Write([]byte("hello 2"))
	f2.Close()

	f3, err := os.Create(filepath.Join(tmpDir, "dot/3.text"))
	f3.Write([]byte("hello 3"))
	f3.Close()

	currentPwd, err := os.Getwd()
	err = os.Chdir(filepath.Join(tmpDir, "dot"))
	defer func() {
		os.Chdir(currentPwd)
		os.RemoveAll(filepath.Join(tmpDir, "/dot"))
	}()
	assert.NoError(err)

	ctx, _ := context.WithTimeout(context.Background(), 90*time.Second)
	job := client.Uploader(
		ctx,
		UploaderParams{
			LocalPath:  "./",
			RemotePath: "go-from-dot",
			EventsReporter: Reporter(func(s status.File) {
				resultsMapMutex.Lock()
				if s.Err != nil {
					panic(s.Err)
				}

				results[s.File.Path] = append(results[s.File.Path], s.TransferBytes)
				resultsMapMutex.Unlock()
			}),
		})
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

	f1, err := os.Create(filepath.Join(tmpDir, "relative/1.text"))
	f1.Write([]byte("hello 1"))
	f1.Close()

	f2, err := os.Create(filepath.Join(tmpDir, "relative/2.text"))
	f2.Write([]byte("hello 2"))
	f2.Close()

	f3, err := os.Create(filepath.Join(tmpDir, "relative/3.text"))
	f3.Write([]byte("hello 3"))
	f3.Close()

	currentPwd, err := os.Getwd()
	err = os.Chdir(tmpDir)
	defer os.Chdir(currentPwd)
	assert.NoError(err)

	job := client.Uploader(
		context.Background(),
		UploaderParams{
			LocalPath:  "relative/",
			RemotePath: "relative",
			EventsReporter: Reporter(func(status status.File) {
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

func TestClient_UploadFile(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	uploadPath := "../LICENSE"
	job := client.Uploader(context.Background(), UploaderParams{LocalPath: uploadPath})
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
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(context.Background(), files_sdk.FileDownloadParams{Path: "LICENSE"}, tempFile.Name())
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
	_, err = folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: "docs"})
	defer deletePath(client, "docs")

	assert.NoError(err)
	uploadPath := "../LICENSE"
	job := client.Uploader(context.Background(), UploaderParams{LocalPath: uploadPath, RemotePath: "docs"})
	job.Start()
	job.Wait()
	tempFile, err := os.CreateTemp(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(context.Background(), files_sdk.FileDownloadParams{Path: "docs/LICENSE"}, tempFile.Name())
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
	uploadPath := "../LICENSE"
	job := client.Uploader(context.Background(), UploaderParams{LocalPath: uploadPath, RemotePath: "taco"})
	defer deletePath(client, "taco")
	job.Start()
	job.Wait()
	tempFile, err := ioutil.TempFile(tmpDir, "LICENSE")
	downloadPath, err := filepath.Abs(filepath.Dir(tempFile.Name()))
	if err != nil {
		panic(err)
	}
	downloadPath = path.Join(downloadPath, tempFile.Name())
	file, err := client.DownloadToFile(context.Background(), files_sdk.FileDownloadParams{Path: "taco"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal("taco", file.DisplayName, "because the docs did not exist as a folder it becomes the file")

	downloadData, err := ioutil.ReadFile(tempFile.Name())
	assert.NoError(err)
	localData, err := ioutil.ReadFile(uploadPath)
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
	uploadPath := "../LICENSE"
	deletePath(client, "docs")
	job := client.Uploader(context.Background(), UploaderParams{LocalPath: uploadPath, RemotePath: "docs/"})
	defer deletePath(client, "docs")
	job.Start()
	job.Wait()
	tempFile, err := os.CreateTemp(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(context.Background(), files_sdk.FileDownloadParams{Path: "docs/LICENSE"}, tempFile.Name())
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

	uploadPath := "../LICENSE"
	job := client.Uploader(context.Background(), UploaderParams{LocalPath: uploadPath})
	job.Start()
	job.Wait()

	assert.Equal(int64(1102), job.TransferBytes())
	assert.Equal(int64(1102), job.TotalBytes())

	tempFile, err := ioutil.TempFile(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(context.Background(), files_sdk.FileDownloadParams{Path: "LICENSE"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal(file.DisplayName, "LICENSE")

	downloadData, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		panic(err)
	}
	localData, err := ioutil.ReadFile(uploadPath)
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

	f1, err := os.Create(filepath.Join(tmpDir, "test/1.text"))
	f1.Write([]byte("hello 1"))
	f1.Close()

	f2, err := os.Create(filepath.Join(tmpDir, "test/2.text"))
	f2.Write([]byte("hello 2"))
	f2.Close()

	f3, err := os.Create(filepath.Join(tmpDir, "test/3.text"))
	f3.Write([]byte("hello 3"))
	f3.Close()
	job := client.Uploader(context.Background(), UploaderParams{LocalPath: filepath.Join(tmpDir, "test"), RemotePath: "/test", Manager: manager.New(1, 1)})
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

	f1, err := os.Create(filepath.Join(tmpDir, "zero_byte_folder/zero-byte-file.text"))
	f1.Close()

	job := client.Uploader(context.Background(), UploaderParams{LocalPath: filepath.Join(tmpDir, "zero_byte_folder"), RemotePath: "", Manager: manager.New(1, 1)})
	job.Start()
	job.Wait()
	require.Len(t, job.Statuses, 1)
	assert.NoError(t, job.Statuses[0].Err())
	assert.Equal(t, "zero_byte_folder/zero-byte-file.text", job.Statuses[0].RemotePath())

	job = client.Downloader(context.Background(), DownloaderParams{RemotePath: "zero_byte_folder", LocalPath: tmpDir + "/"})
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
	folderClient := folder.Client{Config: client.Config}

	it, err := folderClient.ListFor(context.Background(), files_sdk.FolderListForParams{
		ListParams: lib.ListParams{PerPage: 1},
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

	results := runDownloadScenario(filepath.Join("TestClient_DownloadFolder_Smart", "nested_1", "nested_2", "3.text"), "download/", client)

	assert.Len(t, results["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"], 3)
	assert.Equal(t, int64(9), results["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"][2].TransferBytes)

	results2 := runDownloadScenario(filepath.Join("TestClient_DownloadFolder_Smart", "nested_1", "nested_2")+"/", "download", client)

	assert.NoError(t, err)

	path, err := os.Getwd()
	assert.NoError(t, err)

	require.Len(t, results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"], 3)
	assert.Equal(t, int64(9), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"][2].TransferBytes)
	assert.Equal(t, filepath.Join(path, "download/3.text"), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"][2].LocalPath)
	fileInfo, err := os.Stat(filepath.Join(path, "download/3.text"))
	require.NoError(t, err)
	assert.Equal(t, int64(9), fileInfo.Size())
	require.Len(t, results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/nested_3/4.text"], 3)
	assert.Equal(t, int64(9), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/nested_3/4.text"][2].TransferBytes)
	assert.Equal(t, filepath.Join(path, "download/nested_3/4.text"), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/nested_3/4.text"][2].LocalPath)
}

func TestClient_DownloadFolder_file_to_file(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_to_file")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	buildScenario("TestClient_DownloadFolder_file_to_file", client)
	assert := assert.New(t)

	tmpDir, err := ioutil.TempDir(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	results := runDownloadScenario(filepath.Join("TestClient_DownloadFolder_file_to_file", "nested_1", "nested_2", "3.text"), filepath.Join(tmpDir, "3.text"), client)
	assert.NoError(err)

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
	results := runDownloadScenario(filepath.Join("file_to_implicit", "nested_1", "nested_2", "3.text"), "", client)
	assert.NoError(err)

	assert.Equal(int64(9), results["file_to_implicit/nested_1/nested_2/3.text"][2].TransferBytes)
	os.RemoveAll("3.text")
}

func TestClient_DownloadFolder_file_only(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_only")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	_, _, _, _, err = client.UploadIO(
		context.Background(),
		UploadIOParams{
			Path:          filepath.Join("i am at the root.text"),
			Reader:        strings.NewReader("hello"),
			Size:          int64(5),
			ProvidedMtime: time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	)
	require.NoError(t, err)

	assert := assert.New(t)
	results := runDownloadScenario("i am at the root.text", "", client)
	assert.NoError(err)
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

	folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: "test-delete-source"})
	_, _, _, _, err = client.UploadIO(
		context.Background(),
		UploadIOParams{
			Path:          filepath.Join("test-delete-source", "test.text"),
			Reader:        strings.NewReader("testing 3"),
			Size:          int64(9),
			ProvidedMtime: time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	)
	assert.NoError(err)
	job := client.Downloader(
		context.Background(),
		DownloaderParams{RemotePath: filepath.Join("test-delete-source"), LocalPath: ""},
	)
	var fi status.File
	var log status.Log
	job.RegisterFileEvent(func(f status.File) {
		fi = f
		log, err = DeleteSource{Config: client.Config, Direction: job.Direction}.Call(context.Background(), f)
	}, status.Complete)
	job.Start()
	<-job.Finished.Subscribe()
	assert.NoError(err)
	assert.Equal("delete source", log.Action)
	assert.Equal(fi.RemotePath, log.Path)

	_, err = client.Find(context.Background(), files_sdk.FileFindParams{Path: filepath.Join("test-delete-source", "test.text")})
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

	folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: "test-move-source"})
	_, _, _, _, err = client.UploadIO(
		context.Background(),
		UploadIOParams{
			Path:          filepath.Join("test-move-source", "test.text"),
			Reader:        strings.NewReader("testing 3"),
			Size:          int64(9),
			ProvidedMtime: time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	)
	assert.NoError(err)
	job := client.Downloader(
		context.Background(),
		DownloaderParams{RemotePath: filepath.Join("test-move-source"), LocalPath: ""},
	)
	var log status.Log
	job.RegisterFileEvent(func(f status.File) {
		log, err = MoveSource{Config: client.Config, Direction: job.Direction, Path: "test-moved-source"}.Call(context.Background(), f)
	}, status.Complete)
	job.Start()
	job.Wait()

	require.NoError(t, err)
	assert.Equal("move source", log.Action)
	assert.Equal("test-moved-source/test.text", log.Path)

	_, err = client.Find(context.Background(), files_sdk.FileFindParams{Path: filepath.Join("test-move-source", "test.text")})
	assert.Equal("Not Found - `Not Found`", err.Error())
	_, err = client.Find(context.Background(), files_sdk.FileFindParams{Path: filepath.Join("test-moved-source", "test.text")})
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

	tempFile, err := os.Create(filepath.Join(tmpDir, "move-source/upload-move-source.text"))
	assert.NoError(err)
	tempFile.Write([]byte("testing"))
	job := client.Uploader(context.Background(), UploaderParams{LocalPath: tempFile.Name()})
	job.RegisterFileEvent(func(f status.File) {
		log, err = MoveSource{Config: client.Config, Direction: job.Direction, Path: filepath.Join(tmpDir, "move-source/test-moved-source.text")}.Call(context.Background(), f)
	}, status.Complete)
	job.Start()
	job.Wait()
	assert.Equal(false, job.Any(status.Errored))
	assert.NoError(err)
	assert.Equal("move source", log.Action)
	assert.Equal(filepath.Join(tmpDir, "move-source/test-moved-source.text"), log.Path)
	stat, err := os.Stat(log.Path)
	assert.NoError(err)
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
	err = os.MkdirAll(filepath.Join(tmpDir, "/move-source-dir"), 0755)
	assert.NoError(err)
	tempFile, err := os.Create(filepath.Join(tmpDir, "/move-source-dir/upload-move-source.text"))
	assert.NoError(err)
	tempFile.Write([]byte("testing"))
	tempFile.Close()
	job := client.Uploader(context.Background(), UploaderParams{LocalPath: filepath.Join(tmpDir, "move-source-dir") + "/"})
	job.RegisterFileEvent(func(f status.File) {
		log, err = MoveSource{Config: client.Config, Direction: job.Direction, Path: filepath.Join(tmpDir, "moved-source-dir")}.Call(context.Background(), f)
	}, status.Complete)
	job.Start()
	job.Wait()
	erroredFile, ok := job.Find(status.Errored)
	if ok {
		require.NoError(t, erroredFile.Err(), erroredFile.LocalPath())
	}
	require.Equal(t, false, job.Any(status.Errored), "")
	assert.NoError(err)
	assert.Equal("move source", log.Action)
	assert.Equal(filepath.Join(tmpDir, "moved-source-dir/upload-move-source.text"), log.Path)
	stat, err := os.Stat(log.Path)
	assert.NoError(err)
	assert.Equal("upload-move-source.text", stat.Name())
	assert.Equal(false, stat.IsDir())

	_, err = os.Stat(filepath.Join(tmpDir, "move-source-dir/upload-move-source.text"))
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
		context.Background(),
		DownloaderParams{
			RemotePath: "TestClient_Downloader_Move_Source_Missing_Dir",
			LocalPath:  filepath.Join(tmpDir, "TestClient_Downloader_Move_Source_Missing_Dir") + "/",
		},
	)
	job.RegisterFileEvent(func(f status.File) {
		moveLog, err := MoveSource{Config: client.Config, Direction: job.Direction, Path: "TestClient_Downloader_Move_Source_Missing_Dir-moved"}.Call(context.Background(), f)
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
	assert.Equal(files[log1.Path], log1.Path)
	assert.NoError(<-errChan)

	log2 := <-logChan
	assert.Equal("move source", log2.Action)
	assert.Equal(files[log2.Path], log2.Path)
	assert.NoError(<-errChan)

	job.Wait()

	assert.Equal(false, job.Any(status.Errored))

	movedDir, err := client.Find(context.Background(), files_sdk.FileFindParams{Path: "TestClient_Downloader_Move_Source_Missing_Dir-moved"})
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
	tempFile.Write([]byte("testing"))
	var fi status.File
	job := client.Uploader(context.Background(), UploaderParams{LocalPath: tempFile.Name()})
	job.RegisterFileEvent(func(f status.File) {
		fi = f
		log, err = DeleteSource{Config: client.Config, Direction: job.Direction}.Call(context.Background(), f)
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

func TestClient_ListForRecursive(t *testing.T) {
	client, r, err := CreateClient("TestClient_ListForRecursive")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	buildScenario("TestClient_ListForRecursive", client)

	it, err := client.ListForRecursive(context.Background(), files_sdk.FolderListForParams{Path: "TestClient_ListForRecursive"})
	var files []files_sdk.File
	for it.Next() {
		files = append(files, it.Current().(files_sdk.File))
	}

	require.Equal(t, len(files), 6)
	assert.Equal(files[0].Path, "TestClient_ListForRecursive")
	assert.Equal(files[1].Path, "TestClient_ListForRecursive/nested_1")
	assert.Equal(files[2].Path, "TestClient_ListForRecursive/nested_1/nested_2")
	assert.Equal(files[3].Path, "TestClient_ListForRecursive/nested_1/nested_2/3.text")
	assert.Equal(files[4].Path, "TestClient_ListForRecursive/nested_1/nested_2/nested_3")
	assert.Equal(files[5].Path, "TestClient_ListForRecursive/nested_1/nested_2/nested_3/4.text")
}

func TestClient_ListForRecursive_Error(t *testing.T) {
	client, r, err := CreateClient("TestClient_ListForRecursive_Error")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	it, err := client.ListForRecursive(context.Background(), files_sdk.FolderListForParams{Path: "TestClient_ListForRecursive-Not-Found"})
	var files []files_sdk.File
	for it.Next() {
		files = append(files, it.Current().(files_sdk.File))
	}

	assert.Equal(len(files), 0)
	assert.Equal(it.Err().Error(), "open : Authentication Required - `Unauthorized. The API key or Session token is either missing or invalid.`")
}

func TestClient_ListForRecursive_Root(t *testing.T) {
	client, r, err := CreateClient("TestClient_ListForRecursive_Root")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	it, err := client.ListForRecursive(context.Background(), files_sdk.FolderListForParams{Path: ""})
	var files []files_sdk.File
	for it.Next() {
		files = append(files, it.Current().(files_sdk.File))
		assert.NotEqual(it.Current().(files_sdk.File).Path, "")
	}
}
