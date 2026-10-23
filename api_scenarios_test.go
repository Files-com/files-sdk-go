package files_sdk_test

import (
	"bytes"
	"cmp"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	file "github.com/Files-com/files-sdk-go/v3/file"
	file_status "github.com/Files-com/files-sdk-go/v3/file/status"
	file_migration "github.com/Files-com/files-sdk-go/v3/filemigration"
	folder "github.com/Files-com/files-sdk-go/v3/folder"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Generated from verification/scenarios/api.json (SHA-256 5d9b8eff0ef0724a68a458db77dc85784bb7d24436f104b98c2a2fe67086bbe0), the Files.com API
// scenarios every client shares. Each subtest runs one scenario's API operations and client actions through this
// SDK against the mock server's simulator, in a folder of its own. Steps only a Files.com site can answer are left
// out, and a scenario that needs a client action this binding does not implement yet is skipped, naming those
// actions. With FILES_API_SCENARIOS=required, as the Go CI job runs them, a test that would be skipped fails.
func TestAPIScenarios(t *testing.T) {
	t.Run("files.owned-folder-namespace", func(t *testing.T) {
		run := startScenario(t)
		// Step 1: FileActionFind (GET /file_actions/metadata/{path}), 422 bad-request/request-params-required
		run.step(t, "01 FileActionFind", "1 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: "/"}, scenarioObserve(&responses))
			s.failedWith(err, responses, 422, "bad-request/request-params-required", files_sdk.ErrRequestParamsRequired)
		})

		// Step 2: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "02 FolderListForPath", "2 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: "/", ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			_ = records
		})

		// Step 3 (BehaviorListForPath) is for a Files.com site only.

		// Step 4: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "04 FileActionFind", "4 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.namespace}, scenarioObserve(&responses))
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 5: PostFoldersPath (POST /folders/{path}), 201
		run.step(t, "05 PostFoldersPath", "5 (PostFoldersPath, POST /folders/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.namespace, MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 201)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.namespace)
			s.noValue(body, "size")
			s.madeNamespace(entry, err, responses)
		})

		// Step 6 (FileActionFind) is for a Files.com site only.

		// Step 7: PostFoldersPath (POST /folders/{path}), 201
		run.step(t, "07 PostFoldersPath", "7 (PostFoldersPath, POST /folders/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.path("alpha"), MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 201)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("alpha"))
			s.noValue(body, "size")
		})

		// Step 8: PostFoldersPath (POST /folders/{path}), 201
		run.step(t, "08 PostFoldersPath", "8 (PostFoldersPath, POST /folders/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.path("bravo"), MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 201)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("bravo"))
			s.noValue(body, "size")
		})

		// Step 9: PostFoldersPath (POST /folders/{path}), 201
		run.step(t, "09 PostFoldersPath", "9 (PostFoldersPath, POST /folders/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.path("charlie"), MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 201)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("charlie"))
			s.noValue(body, "size")
		})

		// Step 10: PostFoldersPath (POST /folders/{path}), 201
		run.step(t, "10 PostFoldersPath", "10 (PostFoldersPath, POST /folders/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.path("alpha/inner"), MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 201)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("alpha/inner"))
			s.noValue(body, "size")
		})

		// Step 11: PostFoldersPath (POST /folders/{path}), 422 processing-failure/destination-exists
		run.step(t, "11 PostFoldersPath", "11 (PostFoldersPath, POST /folders/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			_, err := (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.path("alpha"), MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
			s.failedWith(err, responses, 422, "processing-failure/destination-exists", files_sdk.ErrDestinationExists)
		})

		// Step 12: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "12 FileActionFind", "12 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("alpha")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("alpha"))
			s.noValue(body, "size")
		})

		// Step 13: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "13 FileActionFind", "13 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("alpha/inner")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("alpha/inner"))
			s.noValue(body, "size")
		})

		// Step 14: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "14 FolderListForPath", "14 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: run.namespace, ListParams: files_sdk.ListParams{PerPage: 1, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1, 16, 0)
			s.cursorChain(before, responses)
			listed := make([]string, 0, len(records))
			for _, record := range records {
				listed = append(listed, record.Type+" "+record.Path)
			}
			s.sameEntries(listed, []string{"directory " + run.path("alpha"), "directory " + run.path("bravo"), "directory " + run.path("charlie")})
		})

		// Step 15: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "15 FolderListForPath", "15 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: run.path("alpha"), ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			listed := make([]string, 0, len(records))
			for _, record := range records {
				listed = append(listed, record.Type+" "+record.Path)
			}
			s.sameEntries(listed, []string{"directory " + run.path("alpha/inner")})
		})

		// Step 16: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "16 DeleteFilesPath", "16 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("alpha/inner"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 17: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "17 FileActionFind", "17 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("alpha/inner")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 18: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "18 DeleteFilesPath", "18 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("alpha"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 19: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "19 FileActionFind", "19 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("alpha")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 20: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "20 DeleteFilesPath", "20 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("bravo"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 21: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "21 FileActionFind", "21 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("bravo")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 22: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "22 DeleteFilesPath", "22 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("charlie"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 23: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "23 FileActionFind", "23 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("charlie")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 24: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "24 DeleteFilesPath", "24 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.namespace, Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 25: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "25 FileActionFind", "25 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.namespace}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})
	})
	t.Run("files.owned-binary-roundtrip", func(t *testing.T) {
		run := startScenario(t)
		// Step 1: FileActionFind (GET /file_actions/metadata/{path}), 422 bad-request/request-params-required
		run.step(t, "01 FileActionFind", "1 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: "/"}, scenarioObserve(&responses))
			s.failedWith(err, responses, 422, "bad-request/request-params-required", files_sdk.ErrRequestParamsRequired)
		})

		// Step 2: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "02 FolderListForPath", "2 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: "/", ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			_ = records
		})

		// Step 3 (BehaviorListForPath) is for a Files.com site only.

		// Step 4: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "04 FileActionFind", "4 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.namespace}, scenarioObserve(&responses))
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 5: PostFoldersPath (POST /folders/{path}), 201
		run.step(t, "05 PostFoldersPath", "5 (PostFoldersPath, POST /folders/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.namespace, MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 201)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.namespace)
			s.noValue(body, "size")
			s.madeNamespace(entry, err, responses)
		})

		// Step 6 (FileActionFind) is for a Files.com site only.

		// Step 7: upload binary.bin
		run.step(t, "07 upload binary.bin", "7 (upload binary.bin)", func(s *scenarioStep) {
			s.upload(scenarioUpload{path: run.path("binary.bin"), recipe: scenarioRecipe{number: 1, multiplier: 131, offset: 17, modulus: 251, length: 65539, sha256: "b65094d36ac1009dac9c73c27b48638a47eb7f23d219a80d1b4d05b58f514499"}, source: "local-file", providedMtime: s.instant("2026-09-24T00:00:00Z")})
		})

		// Step 8: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "08 FileActionFind", "8 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("binary.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("binary.bin"))
			scenarioEqual(s, "size", entry.Size, int64(65539))
			s.integer(body, "size", 65539)
			scenarioTimeOneOf(s, "provided_mtime", entry.ProvidedMtime, s.instant("2026-09-24T00:00:00Z"), nil)
		})

		// Step 9: upload empty.bin
		run.step(t, "09 upload empty.bin", "9 (upload empty.bin)", func(s *scenarioStep) {
			s.upload(scenarioUpload{path: run.path("empty.bin"), recipe: scenarioRecipe{number: 2, multiplier: 1, offset: 0, modulus: 256, length: 0, sha256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}, source: "seekable-stream", providedMtime: s.instant("2026-09-24T00:00:00Z")})
		})

		// Step 10: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "10 FileActionFind", "10 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("empty.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("empty.bin"))
			scenarioEqual(s, "size", entry.Size, int64(0))
			s.integer(body, "size", 0)
			scenarioTimeOneOf(s, "provided_mtime", entry.ProvidedMtime, s.instant("2026-09-24T00:00:00Z"), nil)
		})

		// Step 11: upload short-read.bin
		run.step(t, "11 upload short-read.bin", "11 (upload short-read.bin)", func(s *scenarioStep) {
			s.upload(scenarioUpload{path: run.path("short-read.bin"), recipe: scenarioRecipe{number: 3, multiplier: 29, offset: 101, modulus: 256, length: 65537, sha256: "ed5a86d5b614d5385a85c170c15f25fe75c3faab6a8b29abbb470a6c51faa7f9"}, source: "forward-only-stream", maxRead: 3, providedMtime: s.instant("2026-09-24T00:00:00Z")})
		})

		// Step 12: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "12 FileActionFind", "12 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("short-read.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("short-read.bin"))
			scenarioEqual(s, "size", entry.Size, int64(65537))
			s.integer(body, "size", 65537)
			scenarioTimeOneOf(s, "provided_mtime", entry.ProvidedMtime, s.instant("2026-09-24T00:00:00Z"), nil)
		})

		// Step 13: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "13 FolderListForPath", "13 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: run.namespace, ListParams: files_sdk.ListParams{PerPage: 1, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1, 16, 0)
			s.cursorChain(before, responses)
			listed := make([]string, 0, len(records))
			for _, record := range records {
				listed = append(listed, record.Type+" "+record.Path)
			}
			s.sameEntries(listed, []string{"file " + run.path("binary.bin"), "file " + run.path("empty.bin"), "file " + run.path("short-read.bin")})
		})

		// Step 14: download binary.bin
		run.step(t, "14 download binary.bin", "14 (download binary.bin)", func(s *scenarioStep) {
			s.download(run.path("binary.bin"), scenarioRecipe{number: 1, multiplier: 131, offset: 17, modulus: 251, length: 65539, sha256: "b65094d36ac1009dac9c73c27b48638a47eb7f23d219a80d1b4d05b58f514499"}, false)
		})

		// Step 15: download empty.bin
		run.step(t, "15 download empty.bin", "15 (download empty.bin)", func(s *scenarioStep) {
			s.download(run.path("empty.bin"), scenarioRecipe{number: 2, multiplier: 1, offset: 0, modulus: 256, length: 0, sha256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}, false)
		})

		// Step 16: download short-read.bin
		run.step(t, "16 download short-read.bin", "16 (download short-read.bin)", func(s *scenarioStep) {
			s.download(run.path("short-read.bin"), scenarioRecipe{number: 3, multiplier: 29, offset: 101, modulus: 256, length: 65537, sha256: "ed5a86d5b614d5385a85c170c15f25fe75c3faab6a8b29abbb470a6c51faa7f9"}, false)
		})

		// Step 17: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "17 DeleteFilesPath", "17 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("binary.bin"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 18: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "18 FileActionFind", "18 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("binary.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 19: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "19 DeleteFilesPath", "19 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("empty.bin"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 20: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "20 FileActionFind", "20 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("empty.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 21: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "21 DeleteFilesPath", "21 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("short-read.bin"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 22: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "22 FileActionFind", "22 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("short-read.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 23: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "23 DeleteFilesPath", "23 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.namespace, Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 24: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "24 FileActionFind", "24 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.namespace}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})
	})
	t.Run("files.owned-namespace-extended", func(t *testing.T) {
		run := startScenario(t)
		// Step 1: FileActionFind (GET /file_actions/metadata/{path}), 422 bad-request/request-params-required
		run.step(t, "01 FileActionFind", "1 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: "/"}, scenarioObserve(&responses))
			s.failedWith(err, responses, 422, "bad-request/request-params-required", files_sdk.ErrRequestParamsRequired)
		})

		// Step 2: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "02 FolderListForPath", "2 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: "/", ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			_ = records
		})

		// Step 3 (BehaviorListForPath) is for a Files.com site only.

		// Step 4: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "04 FileActionFind", "4 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.namespace}, scenarioObserve(&responses))
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 5: PostFoldersPath (POST /folders/{path}), 201
		run.step(t, "05 PostFoldersPath", "5 (PostFoldersPath, POST /folders/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.namespace, MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 201)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.namespace)
			s.noValue(body, "size")
			s.madeNamespace(entry, err, responses)
		})

		// Step 6 (FileActionFind) is for a Files.com site only.

		// Step 7: upload hash-name
		run.step(t, "07 upload hash-name", "7 (upload hash-name)", func(s *scenarioStep) {
			s.upload(scenarioUpload{path: run.path("a#b.txt"), recipe: scenarioRecipe{number: 4, multiplier: 37, offset: 5, modulus: 256, length: 1031, sha256: "29f483b8739ae2de994d4794133a7a858a0f0416caa03e5a96165af41ebdf3e8"}, source: "local-file"})
		})

		// Step 8: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "08 FileActionFind", "8 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("a#b.txt")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("a#b.txt"))
			scenarioEqual(s, "size", entry.Size, int64(1031))
			s.integer(body, "size", 1031)
		})

		// Step 9: upload percent-name
		run.step(t, "09 upload percent-name", "9 (upload percent-name)", func(s *scenarioStep) {
			s.upload(scenarioUpload{path: run.path("100%.txt"), recipe: scenarioRecipe{number: 5, multiplier: 41, offset: 11, modulus: 256, length: 2053, sha256: "83298104fe01d905ee9ce8fae2f23da08007b5e8ab356d5bfc064b49c2664b98"}, source: "local-file"})
		})

		// Step 10: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "10 FileActionFind", "10 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("100%.txt")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("100%.txt"))
			scenarioEqual(s, "size", entry.Size, int64(2053))
			s.integer(body, "size", 2053)
		})

		// Step 11: upload unicode-name
		run.step(t, "11 upload unicode-name", "11 (upload unicode-name)", func(s *scenarioStep) {
			s.upload(scenarioUpload{path: run.path("é +.txt"), recipe: scenarioRecipe{number: 6, multiplier: 43, offset: 3, modulus: 256, length: 517, sha256: "b664357b3151e48ca4c5c6c452dd20d58293897f887536f6e0e0627c0d8cd25b"}, source: "local-file"})
		})

		// Step 12: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "12 FileActionFind", "12 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("é +.txt")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("é +.txt"))
			scenarioEqual(s, "size", entry.Size, int64(517))
			s.integer(body, "size", 517)
		})

		// Step 13: PostFoldersPath (POST /folders/{path}), 201
		run.step(t, "13 PostFoldersPath", "13 (PostFoldersPath, POST /folders/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.path("nested"), MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 201)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("nested"))
			s.noValue(body, "size")
		})

		// Step 14: upload inner.bin
		run.step(t, "14 upload inner.bin", "14 (upload inner.bin)", func(s *scenarioStep) {
			s.upload(scenarioUpload{path: run.path("nested/inner.bin"), recipe: scenarioRecipe{number: 7, multiplier: 47, offset: 19, modulus: 256, length: 4099, sha256: "9f251477458622d9da868e31b8708ec638d8ba9a01c9b575bb8c0a0764734521"}, source: "seekable-stream"})
		})

		// Step 15: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "15 FileActionFind", "15 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("nested/inner.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("nested/inner.bin"))
			scenarioEqual(s, "size", entry.Size, int64(4099))
			s.integer(body, "size", 4099)
		})

		// Step 16: upload large.bin
		run.step(t, "16 upload large.bin", "16 (upload large.bin)", func(s *scenarioStep) {
			s.upload(scenarioUpload{path: run.path("large.bin"), recipe: scenarioRecipe{number: 8, multiplier: 53, offset: 23, modulus: 251, length: 6291459, sha256: "4ab1615925cb4b1fa9b9ddc3345a920200888d8ba7d3ae1b880aa6c0c19f87fa"}, source: "local-file", minParts: 2, progressFinal: true})
		})

		// Step 17: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "17 FileActionFind", "17 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("large.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("large.bin"))
			scenarioEqual(s, "size", entry.Size, int64(6291459))
			s.integer(body, "size", 6291459)
		})

		// Step 18: copy copied.bin
		run.step(t, "18 copy copied.bin", "18 (copy copied.bin)", func(s *scenarioStep) {
			action, err := (&file.Client{Config: run.config}).Copy(files_sdk.FileCopyParams{Path: run.path("nested/inner.bin"), Destination: run.path("copied.bin")})
			if err != nil {
				s.fatalf("the copy was refused: %v", err)
			}
			migration, err := (&file_migration.Client{Config: run.config}).Wait(action, func(files_sdk.FileMigration) {}, s.deadline(scenarioMigrationTimeout))
			s.completed("copy", migration.Status, err)
		})

		// Step 19: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "19 FileActionFind", "19 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("copied.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("copied.bin"))
			scenarioEqual(s, "size", entry.Size, int64(4099))
			s.integer(body, "size", 4099)
		})

		// Step 20: move moved.bin
		run.step(t, "20 move moved.bin", "20 (move moved.bin)", func(s *scenarioStep) {
			action, err := (&file.Client{Config: run.config}).Move(files_sdk.FileMoveParams{Path: run.path("copied.bin"), Destination: run.path("nested/moved.bin")})
			if err != nil {
				s.fatalf("the move was refused: %v", err)
			}
			migration, err := (&file_migration.Client{Config: run.config}).Wait(action, func(files_sdk.FileMigration) {}, s.deadline(scenarioMigrationTimeout))
			s.completed("move", migration.Status, err)
		})

		// Step 21: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "21 FileActionFind", "21 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("copied.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 22: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "22 FileActionFind", "22 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("nested/moved.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("nested/moved.bin"))
			scenarioEqual(s, "size", entry.Size, int64(4099))
			s.integer(body, "size", 4099)
		})

		// Step 23: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "23 FolderListForPath", "23 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: run.namespace, ListParams: files_sdk.ListParams{PerPage: 1, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1, 16, 0)
			s.cursorChain(before, responses)
			listed := make([]string, 0, len(records))
			for _, record := range records {
				listed = append(listed, record.Type+" "+record.Path)
				run.listed[record.Path] = record
			}
			s.sameEntries(listed, []string{"file " + run.path("a#b.txt"), "file " + run.path("100%.txt"), "file " + run.path("é +.txt"), "directory " + run.path("nested"), "file " + run.path("large.bin")})
		})

		// Step 24: download hash-name
		run.step(t, "24 download hash-name", "24 (download hash-name)", func(s *scenarioStep) {
			s.download(run.path("a#b.txt"), scenarioRecipe{number: 4, multiplier: 37, offset: 5, modulus: 256, length: 1031, sha256: "29f483b8739ae2de994d4794133a7a858a0f0416caa03e5a96165af41ebdf3e8"}, true)
		})

		// Step 25: download percent-name
		run.step(t, "25 download percent-name", "25 (download percent-name)", func(s *scenarioStep) {
			s.download(run.path("100%.txt"), scenarioRecipe{number: 5, multiplier: 41, offset: 11, modulus: 256, length: 2053, sha256: "83298104fe01d905ee9ce8fae2f23da08007b5e8ab356d5bfc064b49c2664b98"}, true)
		})

		// Step 26: download unicode-name
		run.step(t, "26 download unicode-name", "26 (download unicode-name)", func(s *scenarioStep) {
			s.download(run.path("é +.txt"), scenarioRecipe{number: 6, multiplier: 43, offset: 3, modulus: 256, length: 517, sha256: "b664357b3151e48ca4c5c6c452dd20d58293897f887536f6e0e0627c0d8cd25b"}, true)
		})

		// Step 27: download large.bin
		run.step(t, "27 download large.bin", "27 (download large.bin)", func(s *scenarioStep) {
			s.download(run.path("large.bin"), scenarioRecipe{number: 8, multiplier: 53, offset: 23, modulus: 251, length: 6291459, sha256: "4ab1615925cb4b1fa9b9ddc3345a920200888d8ba7d3ae1b880aa6c0c19f87fa"}, true)
		})

		// Step 28: download inner.bin
		run.step(t, "28 download inner.bin", "28 (download inner.bin)", func(s *scenarioStep) {
			s.download(run.path("nested/inner.bin"), scenarioRecipe{number: 7, multiplier: 47, offset: 19, modulus: 256, length: 4099, sha256: "9f251477458622d9da868e31b8708ec638d8ba9a01c9b575bb8c0a0764734521"}, false)
		})

		// Step 29: download moved.bin
		run.step(t, "29 download moved.bin", "29 (download moved.bin)", func(s *scenarioStep) {
			s.download(run.path("nested/moved.bin"), scenarioRecipe{number: 7, multiplier: 47, offset: 19, modulus: 256, length: 4099, sha256: "9f251477458622d9da868e31b8708ec638d8ba9a01c9b575bb8c0a0764734521"}, false)
		})

		// Step 30: download large.bin
		run.step(t, "30 download large.bin", "30 (download large.bin)", func(s *scenarioStep) {
			s.download(run.path("large.bin"), scenarioRecipe{number: 8, multiplier: 53, offset: 23, modulus: 251, length: 6291459, sha256: "4ab1615925cb4b1fa9b9ddc3345a920200888d8ba7d3ae1b880aa6c0c19f87fa"}, false)
		})

		// Step 31: upload cancelled.bin
		run.step(t, "31 upload cancelled.bin", "31 (upload cancelled.bin)", func(s *scenarioStep) {
			s.upload(scenarioUpload{path: run.path("cancelled.bin"), recipe: scenarioRecipe{number: 9, multiplier: 59, offset: 29, modulus: 251, length: 6291461, sha256: "d7a71574a24faecb2fea4915b8096a86a43b19db98d440f16730caac0a66cfd4"}, source: "local-file", cancelAfterFirstPart: true})
		})

		// Step 32: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "32 FileActionFind", "32 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("cancelled.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 33: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "33 DeleteFilesPath", "33 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("nested/moved.bin"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 34: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "34 FileActionFind", "34 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("nested/moved.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 35: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "35 DeleteFilesPath", "35 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("nested/inner.bin"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 36: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "36 FileActionFind", "36 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("nested/inner.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 37: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "37 DeleteFilesPath", "37 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("nested"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 38: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "38 FileActionFind", "38 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("nested")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 39: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "39 DeleteFilesPath", "39 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("large.bin"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 40: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "40 FileActionFind", "40 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("large.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 41: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "41 DeleteFilesPath", "41 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("a#b.txt"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 42: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "42 FileActionFind", "42 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("a#b.txt")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 43: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "43 DeleteFilesPath", "43 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("100%.txt"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 44: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "44 FileActionFind", "44 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("100%.txt")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 45: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "45 DeleteFilesPath", "45 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("é +.txt"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 46: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "46 FileActionFind", "46 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("é +.txt")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 47: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "47 DeleteFilesPath", "47 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.namespace, Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 48: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "48 FileActionFind", "48 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.namespace}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})
	})
	t.Run("files.owned-directory-helpers", func(t *testing.T) {
		run := startScenario(t)
		// Step 1: FileActionFind (GET /file_actions/metadata/{path}), 422 bad-request/request-params-required
		run.step(t, "01 FileActionFind", "1 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: "/"}, scenarioObserve(&responses))
			s.failedWith(err, responses, 422, "bad-request/request-params-required", files_sdk.ErrRequestParamsRequired)
		})

		// Step 2: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "02 FolderListForPath", "2 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: "/", ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			_ = records
		})

		// Step 3 (BehaviorListForPath) is for a Files.com site only.

		// Step 4: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "04 FileActionFind", "4 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.namespace}, scenarioObserve(&responses))
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 5: PostFoldersPath (POST /folders/{path}), 201
		run.step(t, "05 PostFoldersPath", "5 (PostFoldersPath, POST /folders/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.namespace, MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 201)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.namespace)
			s.noValue(body, "size")
			s.madeNamespace(entry, err, responses)
		})

		// Step 6 (FileActionFind) is for a Files.com site only.

		// Step 7: upload-directory up
		run.step(t, "07 upload-directory up", "7 (upload-directory up)", func(s *scenarioStep) {
			s.uploadDirectory(run.path("up"), []scenarioTreeEntry{
				{relative: "a.bin", recipe: &scenarioRecipe{number: 10, multiplier: 61, offset: 7, modulus: 256, length: 1543, sha256: "983f31356cea338492c47d6c9c979deb3a45914fe6aa8eafa96cd44e02f08e9d"}},
				{relative: "sub"},
				{relative: "sub/b.bin", recipe: &scenarioRecipe{number: 11, multiplier: 67, offset: 13, modulus: 256, length: 2089, sha256: "03e528690696066e61d8a65fa0fac38974bb95cb0d31874240287f325484fddb"}},
				{relative: "empty"},
			}, []string{"a.bin", "sub/b.bin"}, 2, 3632, []string{"empty", "sub"}, []string{"empty"}, 1)
		})

		// Step 8: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "08 FileActionFind", "8 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("up")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("up"))
			s.noValue(body, "size")
		})

		// Step 9: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "09 FileActionFind", "9 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("up/a.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("up/a.bin"))
			scenarioEqual(s, "size", entry.Size, int64(1543))
			s.integer(body, "size", 1543)
		})

		// Step 10: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "10 FileActionFind", "10 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("up/sub")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("up/sub"))
			s.noValue(body, "size")
		})

		// Step 11: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "11 FileActionFind", "11 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("up/sub/b.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("up/sub/b.bin"))
			scenarioEqual(s, "size", entry.Size, int64(2089))
			s.integer(body, "size", 2089)
		})

		// Step 12: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "12 FileActionFind", "12 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("up/empty")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("up/empty"))
			s.noValue(body, "size")
		})

		// Step 13: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "13 FolderListForPath", "13 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: run.path("up"), ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			listed := make([]string, 0, len(records))
			for _, record := range records {
				listed = append(listed, record.Type+" "+record.Path)
			}
			s.sameEntries(listed, []string{"file " + run.path("up/a.bin"), "directory " + run.path("up/sub"), "directory " + run.path("up/empty")})
		})

		// Step 14: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "14 FolderListForPath", "14 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: run.path("up/sub"), ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			listed := make([]string, 0, len(records))
			for _, record := range records {
				listed = append(listed, record.Type+" "+record.Path)
			}
			s.sameEntries(listed, []string{"file " + run.path("up/sub/b.bin")})
		})

		// Step 15: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "15 FolderListForPath", "15 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: run.path("up/empty"), ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			listed := make([]string, 0, len(records))
			for _, record := range records {
				listed = append(listed, record.Type+" "+record.Path)
			}
			s.sameEntries(listed, []string{})
		})

		// Step 16: download-directory up
		run.step(t, "16 download-directory up", "16 (download-directory up)", func(s *scenarioStep) {
			s.downloadDirectory(run.path("up"), 2, 1, 3632, []scenarioTreeEntry{
				{relative: "a.bin", recipe: &scenarioRecipe{number: 10, multiplier: 61, offset: 7, modulus: 256, length: 1543, sha256: "983f31356cea338492c47d6c9c979deb3a45914fe6aa8eafa96cd44e02f08e9d"}},
				{relative: "sub/b.bin", recipe: &scenarioRecipe{number: 11, multiplier: 67, offset: 13, modulus: 256, length: 2089, sha256: "03e528690696066e61d8a65fa0fac38974bb95cb0d31874240287f325484fddb"}},
			})
		})

		// Step 17: copy copy
		run.step(t, "17 copy copy", "17 (copy copy)", func(s *scenarioStep) {
			action, err := (&file.Client{Config: run.config}).Copy(files_sdk.FileCopyParams{Path: run.path("up/sub"), Destination: run.path("copy")})
			if err != nil {
				s.fatalf("the copy was refused: %v", err)
			}
			migration, err := (&file_migration.Client{Config: run.config}).Wait(action, func(files_sdk.FileMigration) {}, s.deadline(scenarioMigrationTimeout))
			s.completed("copy", migration.Status, err)
		})

		// Step 18: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "18 FileActionFind", "18 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("copy")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("copy"))
			s.noValue(body, "size")
		})

		// Step 19: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "19 FileActionFind", "19 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("copy/b.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("copy/b.bin"))
			scenarioEqual(s, "size", entry.Size, int64(2089))
			s.integer(body, "size", 2089)
		})

		// Step 20: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "20 FolderListForPath", "20 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: run.path("copy"), ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			listed := make([]string, 0, len(records))
			for _, record := range records {
				listed = append(listed, record.Type+" "+record.Path)
			}
			s.sameEntries(listed, []string{"file " + run.path("copy/b.bin")})
		})

		// Step 21: move moved
		run.step(t, "21 move moved", "21 (move moved)", func(s *scenarioStep) {
			action, err := (&file.Client{Config: run.config}).Move(files_sdk.FileMoveParams{Path: run.path("copy"), Destination: run.path("moved")})
			if err != nil {
				s.fatalf("the move was refused: %v", err)
			}
			migration, err := (&file_migration.Client{Config: run.config}).Wait(action, func(files_sdk.FileMigration) {}, s.deadline(scenarioMigrationTimeout))
			s.completed("move", migration.Status, err)
		})

		// Step 22: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "22 FileActionFind", "22 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("copy/b.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 23: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "23 FileActionFind", "23 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("copy")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 24: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "24 FileActionFind", "24 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("moved")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "directory")
			scenarioEqual(s, "path", entry.Path, run.path("moved"))
			s.noValue(body, "size")
		})

		// Step 25: FileActionFind (GET /file_actions/metadata/{path}), 200
		run.step(t, "25 FileActionFind", "25 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			var entry files_sdk.File
			var err error
			entry, err = (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("moved/b.bin")}, scenarioObserve(&responses))
			body := s.succeeded(err, responses, 200)
			scenarioEqual(s, "type", entry.Type, "file")
			scenarioEqual(s, "path", entry.Path, run.path("moved/b.bin"))
			scenarioEqual(s, "size", entry.Size, int64(2089))
			s.integer(body, "size", 2089)
		})

		// Step 26: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "26 FolderListForPath", "26 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: run.path("moved"), ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			listed := make([]string, 0, len(records))
			for _, record := range records {
				listed = append(listed, record.Type+" "+record.Path)
			}
			s.sameEntries(listed, []string{"file " + run.path("moved/b.bin")})
		})

		// Step 27: FolderListForPath (GET /folders/{path}), 200
		run.step(t, "27 FolderListForPath", "27 (FolderListForPath, GET /folders/{path})", func(s *scenarioStep) {
			before := run.lastSeq()
			var responses []scenarioResponse
			listing, err := (&folder.Client{Config: run.config}).ListFor(files_sdk.FolderListForParams{Path: run.namespace, ListParams: files_sdk.ListParams{PerPage: 1000, MaxPages: 16}}, scenarioObserve(&responses))
			if err != nil {
				s.fatalf("the listing could not start: %v", err)
			}
			records := scenarioPages[files_sdk.File](s, listing.Iter, &responses, 200, 1000, 16, 0)
			s.cursorChain(before, responses)
			listed := make([]string, 0, len(records))
			for _, record := range records {
				listed = append(listed, record.Type+" "+record.Path)
			}
			s.sameEntries(listed, []string{"directory " + run.path("up"), "directory " + run.path("moved")})
		})

		// Step 28: download moved-b
		run.step(t, "28 download moved-b", "28 (download moved-b)", func(s *scenarioStep) {
			s.download(run.path("moved/b.bin"), scenarioRecipe{number: 11, multiplier: 67, offset: 13, modulus: 256, length: 2089, sha256: "03e528690696066e61d8a65fa0fac38974bb95cb0d31874240287f325484fddb"}, false)
		})

		// Step 29: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "29 DeleteFilesPath", "29 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("moved/b.bin"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 30: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "30 FileActionFind", "30 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("moved/b.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 31: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "31 DeleteFilesPath", "31 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("moved"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 32: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "32 FileActionFind", "32 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("moved")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 33: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "33 DeleteFilesPath", "33 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("up/sub/b.bin"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 34: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "34 FileActionFind", "34 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("up/sub/b.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 35: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "35 DeleteFilesPath", "35 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("up/sub"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 36: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "36 FileActionFind", "36 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("up/sub")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 37: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "37 DeleteFilesPath", "37 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("up/empty"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 38: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "38 FileActionFind", "38 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("up/empty")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 39: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "39 DeleteFilesPath", "39 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("up/a.bin"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 40: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "40 FileActionFind", "40 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("up/a.bin")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 41: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "41 DeleteFilesPath", "41 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.path("up"), Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 42: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "42 FileActionFind", "42 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.path("up")}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})

		// Step 43: DeleteFilesPath (DELETE /files/{path}), 204
		run.step(t, "43 DeleteFilesPath", "43 (DeleteFilesPath, DELETE /files/{path})", func(s *scenarioStep) {
			var responses []scenarioResponse
			err := (&file.Client{Config: run.config}).Delete(files_sdk.FileDeleteParams{Path: run.namespace, Recursive: lib.Bool(false)}, scenarioObserve(&responses))
			s.noContent(s.succeeded(err, responses, 204))
		})

		// Step 44: FileActionFind (GET /file_actions/metadata/{path}), 404 not-found
		run.step(t, "44 FileActionFind", "44 (FileActionFind, GET /file_actions/metadata/{path})", func(s *scenarioStep) {
			responses, err := s.eventually(5, 1*time.Second, files_sdk.ErrNotFound, func(options ...files_sdk.RequestResponseOption) error {
				_, err := (&file.Client{Config: run.config}).Find(files_sdk.FileFindParams{Path: run.namespace}, options...)
				return err
			})
			s.failedWith(err, responses, 404, "not-found", files_sdk.ErrNotFound)
		})
	})
}

const (
	// How long the simulator has to answer its readiness report.
	scenarioReadyTimeout = 30 * time.Second
	// Not a real API key: the simulator accepts any key and journals it only by number.
	scenarioAPIKey = "synthetic-mock-server-api-key"
	// The folder each run's own folder is made in; it is left in place.
	scenarioParent = "go_api_scenarios"
	// How long a copy or move may take to complete.
	scenarioMigrationTimeout = time.Minute
	// How many lookups, a second apart, the runner's own cleanup waits for a deleted folder to be absent.
	scenarioCleanupLookups = 5
)

// scenarioRun is one scenario's run: an SDK configuration that sends every request to the simulator, the owned
// folder its fixtures are made in (which must not exist yet), the entries its listings returned and a local
// folder for the files it uploads and downloads.
type scenarioRun struct {
	t         *testing.T
	origin    string
	control   *http.Client
	transport *scenarioTransport
	config    files_sdk.Config
	folder    string
	namespace string
	// Whether the run made its owned folder: the scenario's step creating it answered 201 Created with exactly that
	// folder (madeNamespace). Only then does cleanup delete it.
	namespaceMade bool
	scratch       string
	listed        map[string]files_sdk.File
	recipes       map[int][]byte
	files         int
}

// startScenario checks that FILES_MOCK_SERVER_HOST and FILES_MOCK_SERVER_PORT name a loopback address where the
// simulator whose schema is FILES_MOCK_SCHEMA_SHA256 answers its readiness report (the schema_sha256 in its
// lib/simulation/generation.json), and makes the run's folder there. Without FILES_MOCK_SCHEMA_SHA256 the test is
// skipped (scenarioSkip); with it, anything else at that address fails the test before the SDK sends anything. The run
// owns its folder, and registers the cleanup that deletes it, only once the folder's create answered 201 Created with
// exactly that folder (scenarioMade); any other answer, a success included, fails the test and leaves the folder alone.
func startScenario(t *testing.T) *scenarioRun {
	t.Helper()
	schema := strings.TrimSpace(os.Getenv("FILES_MOCK_SCHEMA_SHA256"))
	if schema == "" {
		scenarioSkip(t, "Set FILES_MOCK_SERVER_HOST, FILES_MOCK_SERVER_PORT and FILES_MOCK_SCHEMA_SHA256 for the Files.com mock server's simulator (targets/mock-server/with-scenario.rb sets them); these tests never use Files.com.")
	}
	host := strings.TrimSpace(os.Getenv("FILES_MOCK_SERVER_HOST"))
	port := strings.TrimSpace(os.Getenv("FILES_MOCK_SERVER_PORT"))
	if !slices.Contains([]string{"127.0.0.1", "::1", "localhost"}, host) {
		t.Fatalf("FILES_MOCK_SERVER_HOST must be a loopback address (127.0.0.1, ::1, localhost), not %q.", host)
	}
	if number, err := strconv.Atoi(port); err != nil || len(port) > 5 || number < 1 || number > 65535 {
		t.Fatalf("FILES_MOCK_SERVER_PORT must be a port number, not %q.", port)
	}
	base := http.DefaultTransport.(*http.Transport).Clone()
	t.Cleanup(base.CloseIdleConnections)
	run := &scenarioRun{
		t:         t,
		origin:    "http://" + net.JoinHostPort(host, port),
		control:   &http.Client{Transport: base, Timeout: 10 * time.Second},
		transport: &scenarioTransport{base: base},
		scratch:   t.TempDir(),
		listed:    map[string]files_sdk.File{},
		recipes:   map[int][]byte{},
	}
	run.ready(schema)
	after := run.lastSeq()
	key := cmp.Or(strings.TrimSpace(os.Getenv("FILES_MOCK_API_KEY")), scenarioAPIKey)
	run.config = files_sdk.Config{APIKey: key, EndpointOverride: run.origin}.Init().SetCustomClient(&http.Client{Transport: run.transport})

	folders := &folder.Client{Config: run.config}
	if _, err := folders.Create(files_sdk.FolderCreateParams{Path: scenarioParent, MkdirParents: lib.Bool(false)}); err != nil && !files_sdk.IsExist(err) {
		t.Fatalf("Could not make %s: %v", scenarioParent, err)
	}
	run.folder = scenarioParent + "/" + scenarioRandom(t)
	var responses []scenarioResponse
	made, err := folders.Create(files_sdk.FolderCreateParams{Path: run.folder, MkdirParents: lib.Bool(false)}, scenarioObserve(&responses))
	if err = scenarioMade(run.folder, made, err, responses); err != nil {
		t.Fatalf("Could not make the run's folder %s: %v", run.folder, err)
	}
	run.namespace = run.folder + "/scenario"
	t.Cleanup(func() { run.cleanup(after) })
	return run
}

// scenarioSkip skips t, saying why, unless FILES_API_SCENARIOS is "required", where these tests must run
// (the Go CI job) and a test that would be skipped fails instead.
func scenarioSkip(t *testing.T, why string) {
	t.Helper()
	if os.Getenv("FILES_API_SCENARIOS") == "required" {
		t.Fatalf("FILES_API_SCENARIOS is required, so this test may not be skipped: %s", why)
	}
	t.Skip(why)
}

func scenarioRandom(t *testing.T) string {
	random := make([]byte, 16)
	if _, err := rand.Read(random); err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(random)
}

// path is a fixture's path inside the owned folder.
func (r *scenarioRun) path(relative string) string {
	return r.namespace + "/" + relative
}

// ready waits up to scenarioReadyTimeout for something to listen at the origin. Whatever answers first decides: a
// legacy mock server or any other server has no readiness report.
func (r *scenarioRun) ready(schema string) {
	deadline := time.Now().Add(scenarioReadyTimeout)
	for {
		response, err := r.control.Get(r.origin + "/__files_mock/v1/ready")
		if err != nil {
			if time.Now().After(deadline) {
				r.t.Fatalf("Nothing answered at %s within %s: %v", r.origin, scenarioReadyTimeout, err)
			}
			time.Sleep(200 * time.Millisecond)
			continue
		}
		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil || response.StatusCode != http.StatusOK {
			r.t.Fatalf("%s answered %d for its readiness report, so it is not the simulator (%v).", r.origin, response.StatusCode, err)
		}
		var reported struct {
			Status          string `json:"status"`
			Mode            string `json:"mode"`
			ContractVersion int    `json:"contract_version"`
			SchemaSHA256    string `json:"schema_sha256"`
		}
		if err := json.Unmarshal(body, &reported); err != nil {
			r.t.Fatalf("%s answered its readiness report with something other than JSON, so it is not the simulator.", r.origin)
		}
		if reported.Status != "ready" || reported.Mode != "simulation" || reported.ContractVersion != 3 || reported.SchemaSHA256 != schema {
			r.t.Fatalf("%s is not the expected simulator: it reports %+v, not ready, simulation, contract 3 and schema %s.", r.origin, reported, schema)
		}
		return
	}
}

type scenarioJournalEntry struct {
	Seq       int64  `json:"seq"`
	Path      string `json:"path"`
	Operation string `json:"operation"`
	Status    int    `json:"status"`
	UserAgent string `json:"user_agent"`
	// A listing page's cursors as the simulator journals them: the SHA-256 of the cursor the request sent and of
	// the one its response returned, never the cursors themselves.
	CursorSHA256     string `json:"cursor_sha256"`
	NextCursorSHA256 string `json:"next_cursor_sha256"`
	Credentials      struct {
		APIKey  any `json:"api_key"`
		Session any `json:"session"`
	} `json:"credentials"`
}

func (r *scenarioRun) journal() []scenarioJournalEntry {
	response, err := r.control.Get(r.origin + "/__files_mock/v1/journal")
	if err != nil {
		r.t.Fatalf("Could not read the simulator's journal: %v", err)
	}
	defer response.Body.Close()
	var journal struct {
		Complete bool                   `json:"complete"`
		Entries  []scenarioJournalEntry `json:"entries"`
	}
	if err := json.NewDecoder(response.Body).Decode(&journal); err != nil || response.StatusCode != http.StatusOK {
		r.t.Fatalf("The simulator's journal could not be read (%d, %v).", response.StatusCode, err)
	}
	if !journal.Complete {
		r.t.Fatal("The simulator's journal is incomplete, so it cannot show which requests the SDK sent.")
	}
	return journal.Entries
}

func (r *scenarioRun) lastSeq() int64 {
	var last int64
	for _, entry := range r.journal() {
		last = max(last, entry.Seq)
	}
	return last
}

// scenarioMade returns nil only when a folder create made path itself: the call returned no error, the last response
// it received was 201 Created, and it returned path as a directory. Any other answer, a success included, does not
// show that this run made the folder, and neither does a lookup that finds it.
func scenarioMade(path string, entry files_sdk.File, err error, responses []scenarioResponse) error {
	switch {
	case err != nil:
		return err
	case len(responses) == 0:
		return errors.New("the create received no response")
	case responses[len(responses)-1].status != http.StatusCreated:
		return fmt.Errorf("the create answered %d, not %d", responses[len(responses)-1].status, http.StatusCreated)
	case entry.Path != path || entry.Type != "directory":
		return fmt.Errorf("the create returned %q as %q, not the directory %q", entry.Path, entry.Type, path)
	}
	return nil
}

// madeNamespace makes the run the owner of its owned folder, which its cleanup then deletes, at the end of the step
// that creates that folder, once every other assertion of the step has passed: the create must have made exactly that
// folder (scenarioMade). A create that failed or answered otherwise, or one that never ran, such as after the
// scenario's check that the folder was absent failed, leaves it unowned.
func (s *scenarioStep) madeNamespace(entry files_sdk.File, err error, responses []scenarioResponse) {
	s.t.Helper()
	if err = scenarioMade(s.namespace, entry, err, responses); err != nil {
		s.fatalf("did not make the run's owned folder %s: %v", s.namespace, err)
	}
	s.namespaceMade = true
}

// cleanup deletes what the run made, and nothing else. A scenario that stopped before its own deletes leaves its owned
// folder, which is deleted with everything in it, but only when the run made it (madeNamespace): a lookup that finds
// the folder does not show that, so a folder whose create failed, answered otherwise or never ran is left as it is.
// Then the run's folder (startScenario made it) goes, without anything in it, and the parent must stay; each accepted
// delete must leave its folder absent. Last, the API requests the run sent each carried the same synthetic API key
// (the journal numbers keys and never records them), no session and this SDK's User-Agent.
func (r *scenarioRun) cleanup(after int64) {
	t := r.t
	files := &file.Client{Config: r.config}
	if r.namespaceMade {
		if _, err := files.Find(files_sdk.FileFindParams{Path: r.namespace}); err == nil {
			if err := files.Delete(files_sdk.FileDeleteParams{Path: r.namespace, Recursive: lib.Bool(true)}); err != nil {
				t.Errorf("Could not delete the owned folder %s the scenario left: %v", r.namespace, err)
			} else {
				r.confirmAbsent(files, r.namespace)
			}
		} else if !files_sdk.IsNotExist(err) {
			t.Errorf("Could not look up the owned folder %s: %v", r.namespace, err)
		}
	}
	if err := files.Delete(files_sdk.FileDeleteParams{Path: r.folder, Recursive: lib.Bool(false)}); err != nil {
		t.Errorf("Could not delete the run's folder %s: %v", r.folder, err)
	} else {
		r.confirmAbsent(files, r.folder)
	}
	if _, err := files.Find(files_sdk.FileFindParams{Path: scenarioParent}); err != nil {
		t.Errorf("Deleting %s also removed %s: %v", r.folder, scenarioParent, err)
	}

	var api []scenarioJournalEntry
	for _, entry := range r.journal() {
		if entry.Seq > after && strings.HasPrefix(entry.Path, files_sdk.APIPath+"/") {
			api = append(api, entry)
		}
	}
	if len(api) == 0 {
		t.Error("No API request reached the simulator.")
		return
	}
	keys := map[string]bool{}
	for _, entry := range api {
		key, _ := json.Marshal(entry.Credentials.APIKey)
		keys[string(key)] = true
		if _, number := entry.Credentials.APIKey.(float64); !number {
			t.Errorf("Request %d carried no journaled API key: %s.", entry.Seq, key)
		}
		if entry.Credentials.Session != nil {
			t.Errorf("Request %d carried a session.", entry.Seq)
		}
		if !strings.HasPrefix(entry.UserAgent, files_sdk.UserAgent+" ") {
			t.Errorf("Request %d came with the User-Agent %q, not the Go SDK's.", entry.Seq, entry.UserAgent)
		}
	}
	if len(keys) != 1 {
		t.Errorf("The SDK's requests carried %d API keys, not one key on every request.", len(keys))
	}
}

func (r *scenarioRun) confirmAbsent(files *file.Client, path string) {
	for lookup := 0; lookup < scenarioCleanupLookups; lookup++ {
		if lookup > 0 {
			time.Sleep(time.Second)
		}
		_, err := files.Find(files_sdk.FileFindParams{Path: path})
		if files_sdk.IsNotExist(err) {
			return
		}
		if err != nil {
			r.t.Errorf("Could not look up %s: %v", path, err)
			return
		}
	}
	r.t.Errorf("Deleting %s left it.", path)
}

// scenarioTransport sends the SDK's requests. It also tells the run when the SDK has finished with a successful
// response to an upload part it sent outside the API (to an upload URL the server issued), so an upload can be
// cancelled once its first part is acknowledged.
type scenarioTransport struct {
	base   *http.Transport
	mutex  sync.Mutex
	onPart func()
}

func (s *scenarioTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := s.base.RoundTrip(request)
	if err != nil || request.Method != http.MethodPut || strings.HasPrefix(request.URL.Path, files_sdk.APIPath+"/") || response.StatusCode < 200 || response.StatusCode > 299 {
		return response, err
	}
	s.mutex.Lock()
	onPart := s.onPart
	s.mutex.Unlock()
	if onPart != nil {
		response.Body = &scenarioAcknowledged{ReadCloser: response.Body, closed: onPart}
	}
	return response, err
}

// afterPart calls call (once) when the SDK closes the next successful upload part response, or stops doing so.
func (s *scenarioTransport) afterPart(call func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if call == nil {
		s.onPart = nil
		return
	}
	var once sync.Once
	s.onPart = func() { once.Do(call) }
}

type scenarioAcknowledged struct {
	io.ReadCloser
	closed func()
}

func (a *scenarioAcknowledged) Close() error {
	err := a.ReadCloser.Close()
	a.closed()
	return err
}

// scenarioResponse is a response as the SDK received it: its HTTP status and its body.
type scenarioResponse struct {
	status int
	body   []byte
	// The request's path as the SDK sent it, which the simulator journals.
	path string
}

// scenarioObserve records each response a call receives, before the SDK reads it, and gives the SDK the same body.
func scenarioObserve(into *[]scenarioResponse) files_sdk.RequestResponseOption {
	return files_sdk.ResponseOption(func(response *http.Response) error {
		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return err
		}
		response.Body = io.NopCloser(bytes.NewReader(body))
		if response.Request == nil {
			return errors.New("the response names no request")
		}
		*into = append(*into, scenarioResponse{status: response.StatusCode, body: body, path: response.Request.URL.EscapedPath()})
		return nil
	})
}

// scenarioRecipe is a fixture's contents: byte i is (i*multiplier + offset) % modulus, for i below length.
type scenarioRecipe struct {
	number, multiplier, offset, modulus, length int
	sha256                                      string
}

// scenarioStep is one step of a scenario, run as a subtest of its own; its failures name it.
type scenarioStep struct {
	*scenarioRun
	t     *testing.T
	label string
}

// step runs body as the subtest name. A step that fails ends the scenario; one that is skipped does not.
func (r *scenarioRun) step(t *testing.T, name, label string, body func(*scenarioStep)) {
	t.Helper()
	if !t.Run(name, func(t *testing.T) { body(&scenarioStep{scenarioRun: r, t: t, label: label}) }) {
		t.FailNow()
	}
}

func (s *scenarioStep) fatalf(format string, args ...any) {
	s.t.Helper()
	s.t.Fatalf("Step "+s.label+": "+format, args...)
}

func (s *scenarioStep) instant(value string) *time.Time {
	instant, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		s.fatalf("%q is not a date-time: %v", value, err)
	}
	return &instant
}

// deadline is a request option whose context ends after limit.
func (s *scenarioStep) deadline(limit time.Duration) files_sdk.RequestResponseOption {
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	s.t.Cleanup(cancel)
	return files_sdk.WithContext(ctx)
}

// bytes are a recipe's bytes, which must have the length and SHA-256 the definition records.
func (s *scenarioStep) bytes(recipe scenarioRecipe) []byte {
	if data, ok := s.recipes[recipe.number]; ok {
		return data
	}
	data := make([]byte, recipe.length)
	for index := range data {
		data[index] = byte((index*recipe.multiplier + recipe.offset) % recipe.modulus)
	}
	if sum := sha256.Sum256(data); hex.EncodeToString(sum[:]) != recipe.sha256 {
		s.fatalf("recipe %d makes bytes whose SHA-256 is %x, not %s", recipe.number, sum, recipe.sha256)
	}
	s.recipes[recipe.number] = data
	return data
}

// sameBytes checks bytes the client received against a recipe: its length, its SHA-256 and every byte.
func (s *scenarioStep) sameBytes(what string, got []byte, recipe scenarioRecipe) {
	sum := sha256.Sum256(got)
	if len(got) != recipe.length || hex.EncodeToString(sum[:]) != recipe.sha256 || !bytes.Equal(got, s.bytes(recipe)) {
		s.fatalf("%s holds %d bytes with SHA-256 %x, not recipe %d's %d bytes with SHA-256 %s", what, len(got), sum, recipe.number, recipe.length, recipe.sha256)
	}
}

// local is a new path in the run's local folder.
func (s *scenarioStep) local(name string) string {
	s.files++
	return filepath.Join(s.scratch, strconv.Itoa(s.files)+"-"+name)
}

// succeeded checks that the call returned no error and that the last response it received had status, and returns
// that response's body as JSON (nil when it had none).
func (s *scenarioStep) succeeded(err error, responses []scenarioResponse, status int) map[string]any {
	s.t.Helper()
	if err != nil {
		s.fatalf("the call failed: %v", err)
	}
	if len(responses) == 0 {
		s.fatalf("the call received no response")
	}
	last := responses[len(responses)-1]
	if last.status != status {
		s.fatalf("the response had status %d, not %d", last.status, status)
	}
	if len(bytes.TrimSpace(last.body)) == 0 {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(last.body))
	decoder.UseNumber()
	var object map[string]any
	if err := decoder.Decode(&object); err != nil {
		s.fatalf("the response body is not a JSON object: %v", err)
	}
	return object
}

// failedWith checks that the call failed with the API error errorType (matching match with errors.Is) and that
// both the error and the response the SDK received had status.
func (s *scenarioStep) failedWith(err error, responses []scenarioResponse, status int, errorType string, match error) {
	s.t.Helper()
	if err == nil {
		s.fatalf("the call succeeded; it should fail with %s", errorType)
	}
	if !errors.Is(err, match) {
		s.fatalf("the call failed with %v, not %s", err, errorType)
	}
	var responseError files_sdk.ResponseError
	var responseErrorPointer *files_sdk.ResponseError
	switch {
	case errors.As(err, &responseError):
	case errors.As(err, &responseErrorPointer) && responseErrorPointer != nil:
		responseError = *responseErrorPointer
	default:
		s.fatalf("the call failed with %T, not a files_sdk.ResponseError", err)
	}
	if responseError.Type != errorType || responseError.HttpCode != status {
		s.fatalf("the call failed with %s (%d), not %s (%d)", responseError.Type, responseError.HttpCode, errorType, status)
	}
	if len(responses) == 0 || responses[len(responses)-1].status != status {
		s.fatalf("the responses had statuses %v, ending with %d", responses, status)
	}
}

// eventually calls call up to attempts times, interval apart, until it fails with match, and returns that error
// and the responses its last call received. Another error stops the test.
func (s *scenarioStep) eventually(attempts int, interval time.Duration, match error, call func(...files_sdk.RequestResponseOption) error) ([]scenarioResponse, error) {
	s.t.Helper()
	var err error
	var responses []scenarioResponse
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			time.Sleep(interval)
		}
		responses = nil
		err = call(scenarioObserve(&responses))
		if errors.Is(err, match) {
			return responses, err
		}
		if err != nil {
			s.fatalf("lookup %d failed with %v, not %v", attempt+1, err, match)
		}
	}
	s.fatalf("%d lookups, %s apart, never answered %v", attempts, interval, match)
	return responses, err
}

func (s *scenarioStep) noContent(body map[string]any) {
	s.t.Helper()
	if body != nil {
		s.fatalf("the response has content %v, not none", body)
	}
}

// noValue checks that the body has no value for name: it is absent or null.
func (s *scenarioStep) noValue(body map[string]any, name string) {
	s.t.Helper()
	if value := body[name]; value != nil {
		s.fatalf("%s is %v, not absent", name, value)
	}
}

// integer checks that the body gives name as the JSON integer want.
func (s *scenarioStep) integer(body map[string]any, name string, want int64) {
	s.t.Helper()
	number, ok := body[name].(json.Number)
	if got, err := number.Int64(); !ok || err != nil || got != want {
		s.fatalf("%s is %v, not the integer %d", name, body[name], want)
	}
}

// valueOneOf checks the body's value for name against allowed, where nil is an absent or null value.
func (s *scenarioStep) valueOneOf(body map[string]any, name string, allowed ...any) {
	s.t.Helper()
	if !slices.Contains(allowed, body[name]) {
		s.fatalf("%s is %#v, not one of %#v", name, body[name], allowed)
	}
}

func scenarioEqual[T comparable](s *scenarioStep, name string, got, want T) {
	s.t.Helper()
	if got != want {
		s.fatalf("%s is %#v, not %#v", name, got, want)
	}
}

func scenarioOneOf[T comparable](s *scenarioStep, name string, got T, allowed ...T) {
	s.t.Helper()
	if !slices.Contains(allowed, got) {
		s.fatalf("%s is %#v, not one of %#v", name, got, allowed)
	}
}

func scenarioNot[T comparable](s *scenarioStep, name string, got, excluded T) {
	s.t.Helper()
	if got == excluded {
		s.fatalf("%s is %#v", name, got)
	}
}

func scenarioIncludes(s *scenarioStep, name string, got string, characters ...string) {
	s.t.Helper()
	for _, character := range characters {
		if !strings.Contains(got, character) {
			s.fatalf("%s is %q, without %q", name, got, character)
		}
	}
}

// scenarioPointerOneOf compares what pointers point to; a nil one is an absent or null value.
func scenarioPointerOneOf[T comparable](s *scenarioStep, name string, got *T, allowed ...*T) {
	s.t.Helper()
	if !slices.ContainsFunc(allowed, func(want *T) bool { return (got == nil) == (want == nil) && (got == nil || *got == *want) }) {
		s.fatalf("%s is %v, not one of the %d allowed values", name, got, len(allowed))
	}
}

func scenarioPointerNot[T comparable](s *scenarioStep, name string, got, excluded *T) {
	s.t.Helper()
	if (got == nil) == (excluded == nil) && (got == nil || *got == *excluded) {
		s.fatalf("%s is %v", name, got)
	}
}

// scenarioTimeOneOf compares instants; a nil time is an absent or null value.
func scenarioTimeOneOf(s *scenarioStep, name string, got *time.Time, allowed ...*time.Time) {
	s.t.Helper()
	if !slices.ContainsFunc(allowed, func(want *time.Time) bool { return (got == nil) == (want == nil) && (got == nil || got.Equal(*want)) }) {
		s.fatalf("%s is %v, not one of %v", name, got, allowed)
	}
}

func scenarioTimeNot(s *scenarioStep, name string, got, excluded *time.Time) {
	s.t.Helper()
	if (got == nil) == (excluded == nil) && (got == nil || got.Equal(*excluded)) {
		s.fatalf("%s is %v", name, got)
	}
}

// scenarioPages reads every page of a listing the way the SDK pages it: each page holds at most perPage records of
// type T, every page's response had status, there are at most maxPages pages (and maxRecords records, unless it is
// 0), and the last page gives no next page.
func scenarioPages[T any](s *scenarioStep, listing *files_sdk.Iter, responses *[]scenarioResponse, status, perPage, maxPages, maxRecords int) []T {
	s.t.Helper()
	var records []T
	for listing.GetPage() {
		page := *listing.Values
		if len(page) > perPage {
			s.fatalf("a page holds %d records, more than %d", len(page), perPage)
		}
		for _, value := range page {
			record, ok := value.(T)
			if !ok {
				s.fatalf("the listing returned a %T, not a %T", value, record)
			}
			records = append(records, record)
		}
	}
	if err := listing.Err(); err != nil {
		s.fatalf("a page failed: %v", err)
	}
	if listing.NextPage() {
		s.fatalf("the listing has more than %d pages", maxPages)
	}
	if len(*responses) == 0 || len(*responses) > maxPages {
		s.fatalf("the listing took %d pages, not 1 to %d", len(*responses), maxPages)
	}
	for page, response := range *responses {
		if response.status != status {
			s.fatalf("page %d had status %d, not %d", page+1, response.status, status)
		}
	}
	if maxRecords > 0 && len(records) > maxRecords {
		s.fatalf("the listing holds %d records, more than %d", len(records), maxRecords)
	}
	return records
}

// cursorChain checks the listing's page requests in the simulator's journal (those after seq before, at the paths
// the SDK's responses came from): one for each response the SDK received, the first sending no cursor, each later
// one sending exactly the cursor the page before it returned, each answered with 200, and the last returning no
// cursor. The journal holds only the cursors' SHA-256s.
func (s *scenarioStep) cursorChain(before int64, responses []scenarioResponse) {
	s.t.Helper()
	paths := map[string]bool{}
	for _, response := range responses {
		paths[response.path] = true
	}
	var pages []scenarioJournalEntry
	for _, entry := range s.journal() {
		if entry.Seq > before && entry.Operation == "folders.list" && paths[entry.Path] {
			pages = append(pages, entry)
		}
	}
	if len(pages) != len(responses) {
		s.fatalf("the simulator journaled %d page requests, not the %d responses the SDK received", len(pages), len(responses))
	}
	for index, page := range pages {
		sent := ""
		if index > 0 {
			sent = pages[index-1].NextCursorSHA256
		}
		if page.CursorSHA256 != sent || page.Status != http.StatusOK {
			s.fatalf("page %d (journal seq %d) answered %d to a cursor that was not the one the page before it returned", index+1, page.Seq, page.Status)
		}
	}
	if last := pages[len(pages)-1]; last.NextCursorSHA256 != "" {
		s.fatalf("the last page (journal seq %d) returned a next cursor", last.Seq)
	}
}

// sameEntries checks that a listing holds each expected entry ("type path") exactly once and nothing else.
func (s *scenarioStep) sameEntries(got, want []string) {
	s.t.Helper()
	got, want = slices.Sorted(slices.Values(got)), slices.Sorted(slices.Values(want))
	if !slices.Equal(got, want) {
		s.fatalf("the listing holds %q, not %q", got, want)
	}
}

func (s *scenarioStep) completed(action, status string, err error) {
	s.t.Helper()
	if err != nil || status != "completed" {
		s.fatalf("the %s ended %q (%v), not completed", action, status, err)
	}
}

// scenarioUpload is an upload intent: a recipe's bytes from a local file (file.Client.Uploader, which reports
// the job's progress and is cancelled with Job.Cancel), a seekable stream (UploadWithReader with a
// *bytes.Reader) or a forward-only stream giving at most maxRead bytes a read and nothing else
// (UploadWithReader), with what the step expects of it.
type scenarioUpload struct {
	path                 string
	recipe               scenarioRecipe
	source               string
	maxRead              int
	providedMtime        *time.Time
	minParts             int
	progressFinal        bool
	cancelAfterFirstPart bool
}

func (s *scenarioStep) upload(upload scenarioUpload) {
	s.t.Helper()
	data := s.bytes(upload.recipe)
	files := &file.Client{Config: s.config}
	if upload.source != "local-file" {
		var reader io.Reader = bytes.NewReader(data)
		if upload.source == "forward-only-stream" {
			reader = &scenarioForwardOnly{data: data, maxRead: upload.maxRead}
		}
		options := []file.UploadOption{file.UploadWithReader(reader), file.UploadWithDestinationPath(upload.path)}
		if upload.providedMtime != nil {
			options = append(options, file.UploadWithProvidedMtime(*upload.providedMtime))
		}
		uploaded, err := files.UploadWithResume(options...)
		if err != nil {
			s.fatalf("the upload failed: %v", err)
		}
		if forwardOnly, ok := reader.(*scenarioForwardOnly); ok && forwardOnly.offset != len(data) {
			s.fatalf("the upload read %d of the stream's %d bytes", forwardOnly.offset, len(data))
		}
		scenarioEqual(s, "the uploaded file's path", uploaded.File.Path, upload.path)
		s.parts(uploaded.Parts, upload.minParts)
		return
	}

	local := s.local(filepath.Base(upload.path))
	if err := os.WriteFile(local, data, 0o600); err != nil {
		s.fatalf("could not write %s: %v", local, err)
	}
	if upload.providedMtime != nil {
		// The uploader sends a local file's modification time as its provided mtime.
		if err := os.Chtimes(local, *upload.providedMtime, *upload.providedMtime); err != nil {
			s.fatalf("could not set %s's modification time: %v", local, err)
		}
	}
	job := files.Uploader(file.UploaderParams{LocalPath: local, RemotePath: upload.path, PreserveTimes: upload.providedMtime != nil})
	var acknowledged atomic.Bool
	if upload.cancelAfterFirstPart {
		s.transport.afterPart(func() {
			acknowledged.Store(true)
			job.Cancel()
		})
		defer s.transport.afterPart(nil)
	}
	job.Start()
	job.Wait()
	if len(job.Statuses) != 1 {
		s.fatalf("the upload job has %d files, not 1", len(job.Statuses))
	}
	uploaded := job.Statuses[0]
	if upload.cancelAfterFirstPart {
		s.transport.afterPart(nil)
		if !acknowledged.Load() || !uploaded.Status().Is(file_status.Canceled) {
			s.fatalf("the upload ended %s (%v), a part acknowledged: %v; not cancelled after its first part", uploaded.Status().Name, uploaded.Err(), acknowledged.Load())
		}
		return
	}
	if !uploaded.Status().Is(file_status.Complete) || uploaded.Err() != nil {
		s.fatalf("the upload ended %s (%v), not complete", uploaded.Status().Name, uploaded.Err())
	}
	scenarioEqual(s, "the uploaded file's path", uploaded.RemotePath(), upload.path)
	status, ok := uploaded.(*file.UploadStatus)
	if !ok {
		s.fatalf("the upload job reports a %T, not a *file.UploadStatus", uploaded)
	}
	s.parts(status.UploadResumable.Parts, upload.minParts)
	if upload.progressFinal && (job.TotalBytes() != int64(upload.recipe.length) || job.TransferBytes() != int64(upload.recipe.length)) {
		s.fatalf("the upload's final progress is %d of %d bytes, not %d of %d", job.TransferBytes(), job.TotalBytes(), upload.recipe.length, upload.recipe.length)
	}
}

// parts checks that an upload was stored in at least minParts parts the server acknowledged.
func (s *scenarioStep) parts(parts file.Parts, minParts int) {
	s.t.Helper()
	acknowledged := 0
	for _, part := range parts {
		if part.Successful() {
			acknowledged++
		}
	}
	if acknowledged < minParts {
		s.fatalf("the upload was stored in %d acknowledged parts, fewer than %d", acknowledged, minParts)
	}
}

// scenarioForwardOnly is a stream that can only be read forward, at most maxRead bytes at a time.
type scenarioForwardOnly struct {
	data    []byte
	offset  int
	maxRead int
}

func (f *scenarioForwardOnly) Read(buffer []byte) (int, error) {
	if f.offset == len(f.data) {
		return 0, io.EOF
	}
	count := copy(buffer[:min(len(buffer), f.maxRead)], f.data[f.offset:])
	f.offset += count
	return count, nil
}

// download downloads the file at path, or the entry a listing returned for it, with DownloadToFile and checks its
// bytes against the recipe.
func (s *scenarioStep) download(path string, recipe scenarioRecipe, listedEntry bool) {
	s.t.Helper()
	params := files_sdk.FileDownloadParams{Path: path}
	if listedEntry {
		entry, ok := s.listed[path]
		if !ok {
			s.fatalf("no listing returned an entry for %s", path)
		}
		params = files_sdk.FileDownloadParams{File: entry}
	}
	local := s.local(filepath.Base(path))
	if _, err := (&file.Client{Config: s.config}).DownloadToFile(params, local); err != nil {
		s.fatalf("the download failed: %v", err)
	}
	got, err := os.ReadFile(local)
	if err != nil {
		s.fatalf("could not read %s: %v", local, err)
	}
	s.sameBytes("the downloaded file", got, recipe)
}

// scenarioTreeEntry is a file (with its recipe) or a folder (without one), relative to a folder fixture.
type scenarioTreeEntry struct {
	relative string
	recipe   *scenarioRecipe
}

// uploadDirectory makes the tree in a new local folder and uploads that folder's contents to remote with
// file.Client.Uploader. The job must complete exactly the files completed (relative paths), files of them with
// totalBytes bytes in all, and create exactly the folders directories (every folder the tree has, including those
// above its entries), of which those that hold nothing it uploaded are exactly emptyDirectories, emptyFolders of
// them; and report nothing else.
func (s *scenarioStep) uploadDirectory(remote string, tree []scenarioTreeEntry, completed []string, files int, totalBytes int64, directories, emptyDirectories []string, emptyFolders int) {
	s.t.Helper()
	local := s.local("upload-directory")
	if err := os.Mkdir(local, 0o700); err != nil {
		s.fatalf("could not make %s: %v", local, err)
	}
	for _, entry := range tree {
		target := filepath.Join(local, filepath.FromSlash(entry.relative))
		if entry.recipe == nil {
			if err := os.MkdirAll(target, 0o700); err != nil {
				s.fatalf("could not make %s: %v", target, err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			s.fatalf("could not make %s: %v", filepath.Dir(target), err)
		}
		if err := os.WriteFile(target, s.bytes(*entry.recipe), 0o600); err != nil {
			s.fatalf("could not write %s: %v", target, err)
		}
	}
	job := (&file.Client{Config: s.config}).Uploader(file.UploaderParams{LocalPath: local + string(filepath.Separator), RemotePath: remote})
	job.Start()
	job.Wait()
	var done, created []string
	var gotBytes int64
	for _, uploaded := range job.Statuses {
		relative, err := filepath.Rel(local, uploaded.LocalPath())
		if err != nil {
			s.fatalf("the job uploaded %s, outside %s", uploaded.LocalPath(), local)
		}
		switch {
		case uploaded.Status().Is(file_status.Complete) && uploaded.Err() == nil:
			done = append(done, filepath.ToSlash(relative))
			gotBytes += uploaded.Size()
		case uploaded.Status().Is(file_status.FolderCreated) && uploaded.Err() == nil:
			created = append(created, filepath.ToSlash(relative))
		default:
			s.fatalf("the job ended %s for %s (%v)", uploaded.Status().Name, filepath.ToSlash(relative), uploaded.Err())
		}
	}
	slices.Sort(done)
	if want := slices.Sorted(slices.Values(completed)); !slices.Equal(done, want) || len(done) != files {
		s.fatalf("the job completed %q, not %q (%d files)", done, want, files)
	}
	if gotBytes != totalBytes {
		s.fatalf("the job uploaded %d bytes, not %d", gotBytes, totalBytes)
	}
	slices.Sort(created)
	if want := slices.Sorted(slices.Values(directories)); !slices.Equal(created, want) {
		s.fatalf("the job created the folders %q, not %q", created, want)
	}
	// A folder the job created is empty when nothing it uploaded or created is inside it.
	var empty []string
	for _, folder := range created {
		inside := func(path string) bool { return strings.HasPrefix(path, folder+"/") }
		if !slices.ContainsFunc(done, inside) && !slices.ContainsFunc(created, inside) {
			empty = append(empty, folder)
		}
	}
	if want := slices.Sorted(slices.Values(emptyDirectories)); !slices.Equal(empty, want) || len(empty) != emptyFolders {
		s.fatalf("the job's empty folders are %q, not the %d expected, %q", empty, emptyFolders, want)
	}
}

// downloadDirectory downloads the folder at remote into a new local folder with file.Client.Downloader (its
// contents, as a remote path ending in a slash selects them), which must then hold files files (bytes bytes in
// all) and emptyFolders empty folders, each matched file with its recipe's bytes.
func (s *scenarioStep) downloadDirectory(remote string, files, emptyFolders int, totalBytes int64, matched []scenarioTreeEntry) {
	s.t.Helper()
	local := s.local("download-directory")
	job := (&file.Client{Config: s.config}).Downloader(file.DownloaderParams{RemotePath: remote + "/", LocalPath: local + string(filepath.Separator)})
	job.Start()
	job.Wait()
	for _, downloaded := range job.Statuses {
		if downloaded.Err() != nil || !downloaded.Status().Is(file_status.Complete, file_status.FolderCreated) {
			s.fatalf("the job ended %s for %s (%v)", downloaded.Status().Name, downloaded.RemotePath(), downloaded.Err())
		}
	}
	var gotFiles, gotEmpty int
	var gotBytes int64
	err := filepath.WalkDir(local, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			gotFiles++
			gotBytes += info.Size()
			return nil
		}
		children, err := os.ReadDir(path)
		if err == nil && len(children) == 0 && path != local {
			gotEmpty++
		}
		return err
	})
	if err != nil {
		s.fatalf("could not read the downloaded folder %s: %v", local, err)
	}
	if gotFiles != files || gotEmpty != emptyFolders || gotBytes != totalBytes {
		s.fatalf("the download holds %d files (%d bytes) and %d empty folders, not %d files (%d bytes) and %d empty folders", gotFiles, gotBytes, gotEmpty, files, totalBytes, emptyFolders)
	}
	for _, entry := range matched {
		got, err := os.ReadFile(filepath.Join(local, filepath.FromSlash(entry.relative)))
		if err != nil {
			s.fatalf("the download has no %s: %v", entry.relative, err)
		}
		s.sameBytes("the downloaded "+entry.relative, got, *entry.recipe)
	}
}

// Set in the test process TestAPIScenarioRunDeletesOnlyFoldersItMade starts for each of its cases: that process runs
// the case's run (scenarioCleanupRun) instead of the cases.
const scenarioCleanupCase = "FILES_GO_API_SCENARIOS_CLEANUP_CASE"

// TestAPIScenarioRunDeletesOnlyFoldersItMade checks that a run's cleanup deletes the run's folder and its owned folder
// only when the run made them: their creates answered 201 Created with exactly that folder, a directory. Each case's
// run starts (startScenario), creates its owned folder as a scenario's step does (madeNamespace) and cleans up, in a
// test process of its own, since a run that did not make a folder fails as a scenario would. Its SDK requests go to a
// loopback stand-in for the simulator (scenarioStandIn), which answers those creates as the case says, makes each
// folder all the same, so that a lookup finds it, and records the deletes it is asked for. A folder the run did not
// make must be left in place, neither deleted itself nor through the folder holding it.
func TestAPIScenarioRunDeletesOnlyFoldersItMade(t *testing.T) {
	if os.Getenv(scenarioCleanupCase) != "" {
		scenarioCleanupRun(t)
		return
	}
	made := scenarioCreate{status: http.StatusCreated}
	cases := []struct {
		name              string
		folder, namespace scenarioCreate
		// Whether the run made its folder and its owned folder, and so must delete them.
		madeFolder, madeNamespace bool
	}{
		{"made both folders", made, made, true, true},
		{"folder create answered 200", scenarioCreate{status: http.StatusOK}, made, false, false},
		{"folder create answered another path", scenarioCreate{status: http.StatusCreated, elsewhere: true}, made, false, false},
		{"folder create answered a file", scenarioCreate{status: http.StatusCreated, file: true}, made, false, false},
		{"owned folder create answered 200", made, scenarioCreate{status: http.StatusOK}, true, false},
		{"owned folder create answered another path", made, scenarioCreate{status: http.StatusCreated, elsewhere: true}, true, false},
		{"owned folder create answered a file", made, scenarioCreate{status: http.StatusCreated, file: true}, true, false},
		{"owned folder already there", made, scenarioCreate{existing: true}, true, false},
	}
	test := t.Name()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			standIn := &scenarioStandIn{schema: scenarioRandom(t) + scenarioRandom(t), folder: c.folder, namespace: c.namespace, folders: map[string]bool{}}
			server := httptest.NewServer(standIn)
			t.Cleanup(server.Close)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			process := exec.CommandContext(ctx, os.Args[0], "-test.run=^"+test+"$", "-test.v")
			process.Env = append(os.Environ(), scenarioCleanupCase+"="+c.name, "FILES_MOCK_SERVER_HOST=127.0.0.1",
				"FILES_MOCK_SERVER_PORT="+strconv.Itoa(server.Listener.Addr().(*net.TCPAddr).Port), "FILES_MOCK_SCHEMA_SHA256="+standIn.schema)
			output, err := process.CombinedOutput()

			standIn.mutex.Lock()
			defer standIn.mutex.Unlock()
			if standIn.runFolder == "" {
				t.Fatalf("The run never created its folder (%v):\n%s", err, output)
			}
			var failed *exec.ExitError
			if passes := c.madeFolder && c.madeNamespace; passes && err != nil {
				t.Errorf("The run failed (%v):\n%s", err, output)
			} else if !passes && !errors.As(err, &failed) {
				t.Errorf("The run did not fail (%v):\n%s", err, output)
			}
			var want []string
			if c.madeNamespace {
				want = append(want, standIn.ownedFolder)
			}
			if c.madeFolder {
				want = append(want, standIn.runFolder)
			}
			if !slices.Equal(standIn.deleted, want) {
				t.Errorf("The run asked to delete %q, not %q:\n%s", standIn.deleted, want, output)
			}
			for path, madeIt := range map[string]bool{standIn.runFolder: c.madeFolder, standIn.ownedFolder: c.madeNamespace} {
				if path != "" && !madeIt && !standIn.folders[path] {
					t.Errorf("%s, which the run did not make, is gone:\n%s", path, output)
				}
			}
		})
	}
}

// scenarioCleanupRun is one case of TestAPIScenarioRunDeletesOnlyFoldersItMade, in the test process the case started:
// the run starts, creates its owned folder through the SDK as the scenario step creating it does, and cleans up when
// the test ends.
func scenarioCleanupRun(t *testing.T) {
	run := startScenario(t)
	run.step(t, "01 FolderCreate", "1 (create the owned folder)", func(s *scenarioStep) {
		var responses []scenarioResponse
		entry, err := (&folder.Client{Config: run.config}).Create(files_sdk.FolderCreateParams{Path: run.namespace}, scenarioObserve(&responses))
		s.madeNamespace(entry, err, responses)
	})
}

// scenarioStandIn stands in for the simulator, on loopback, in a case of TestAPIScenarioRunDeletesOnlyFoldersItMade.
// It answers the readiness report and the journal, and creates, looks up and deletes folders, deleting one that holds
// others only when the delete asks for its contents too, as the simulator does. It answers the creates of the run's
// folder (the first one in scenarioParent) and of its owned folder (one in the run's folder) as the case says, and
// records the paths it is asked to delete.
type scenarioStandIn struct {
	schema            string
	folder, namespace scenarioCreate
	mutex             sync.Mutex
	folders           map[string]bool
	runFolder         string
	ownedFolder       string
	deleted           []string
	journal           []scenarioJournalEntry
}

// scenarioCreate is how the stand-in answers a folder's create. It makes the folder and answers status with it, at
// another path when elsewhere is set and as a file when file is; an existing folder is already there, so its create
// is refused.
type scenarioCreate struct {
	status          int
	elsewhere, file bool
	existing        bool
}

// The simulator's answer for a path that is not there.
var scenarioNotFound = map[string]any{"error": "Not Found", "http-code": http.StatusNotFound, "title": "Not Found", "type": "not-found"}

func (s *scenarioStandIn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	switch r.URL.Path {
	case "/__files_mock/v1/ready":
		scenarioReply(w, http.StatusOK, map[string]any{"status": "ready", "mode": "simulation", "contract_version": 3, "schema_sha256": s.schema})
		return
	case "/__files_mock/v1/journal":
		scenarioReply(w, http.StatusOK, map[string]any{"complete": true, "entries": s.journal})
		return
	}
	status, body := s.answer(r)
	entry := scenarioJournalEntry{Seq: int64(len(s.journal) + 1), Path: r.URL.EscapedPath(), Status: status, UserAgent: r.UserAgent()}
	// As the simulator does, the journal numbers the credentials a request carried and never records them.
	if r.Header.Get("X-FilesAPI-Key") != "" {
		entry.Credentials.APIKey = 1
	}
	if r.Header.Get("X-FilesAPI-Auth") != "" {
		entry.Credentials.Session = 1
	}
	s.journal = append(s.journal, entry)
	scenarioReply(w, status, body)
}

// answer answers an API request for scenarioParent or a folder in it by its method: POST creates the folder, GET looks
// it up and DELETE deletes it, with everything in it when its recursive query parameter is true.
func (s *scenarioStandIn) answer(r *http.Request) (int, any) {
	start := strings.Index(r.URL.Path, "/"+scenarioParent)
	if start < 0 {
		return http.StatusNotFound, scenarioNotFound
	}
	path := r.URL.Path[start+1:]
	switch r.Method {
	case http.MethodPost:
		return s.create(path)
	case http.MethodGet:
		if !s.folders[path] {
			return http.StatusNotFound, scenarioNotFound
		}
		return http.StatusOK, map[string]any{"path": path, "type": "directory"}
	case http.MethodDelete:
		recursive, _ := strconv.ParseBool(r.URL.Query().Get("recursive"))
		return s.remove(path, recursive)
	}
	return http.StatusMethodNotAllowed, nil
}

func (s *scenarioStandIn) create(path string) (int, any) {
	answer := scenarioCreate{status: http.StatusCreated}
	switch {
	case s.runFolder == "" && strings.HasPrefix(path, scenarioParent+"/"):
		s.runFolder, answer = path, s.folder
	case s.runFolder != "" && strings.HasPrefix(path, s.runFolder+"/"):
		s.ownedFolder, answer = path, s.namespace
	}
	if s.folders[path] || answer.existing {
		s.folders[path] = true
		return http.StatusUnprocessableEntity, map[string]any{"error": "Destination exists", "http-code": http.StatusUnprocessableEntity, "title": "Destination Exists", "type": "processing-failure/destination-exists"}
	}
	s.folders[path] = true
	entry := map[string]any{"path": path, "type": "directory"}
	if answer.elsewhere {
		entry["path"] = path + "-elsewhere"
	}
	if answer.file {
		entry["type"] = "file"
	}
	return answer.status, entry
}

func (s *scenarioStandIn) remove(path string, recursive bool) (int, any) {
	s.deleted = append(s.deleted, path)
	if !s.folders[path] {
		return http.StatusNotFound, scenarioNotFound
	}
	for other := range s.folders {
		if strings.HasPrefix(other, path+"/") && !recursive {
			return http.StatusUnprocessableEntity, map[string]any{"error": "Folder not empty", "http-code": http.StatusUnprocessableEntity, "title": "Folder Not Empty", "type": "processing-failure/folder-not-empty"}
		}
	}
	for other := range s.folders {
		if other == path || strings.HasPrefix(other, path+"/") {
			delete(s.folders, other)
		}
	}
	return http.StatusNoContent, nil
}

// scenarioReply answers status with body as JSON, or with no body when body is nil.
func scenarioReply(w http.ResponseWriter, status int, body any) {
	if body == nil {
		w.WriteHeader(status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
