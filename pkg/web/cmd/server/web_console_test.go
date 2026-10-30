// Copyright 2026 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/ui/dashboard"
)

func TestSpaFileSystem_Open_invalidUTF8Path(t *testing.T) {
	t.Parallel()

	sfs := &spaFileSystem{root: http.FS(dashboard.FS)}
	_, err := sfs.Open(".env.local\x85")
	require.ErrorIs(t, err, fs.ErrNotExist)
}

func TestDashboardHandler_invalidUTF8PathReturns404(t *testing.T) {
	t.Parallel()

	handler := dashboardHandler()
	req := httptest.NewRequest(http.MethodGet, "/.env.local%85", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// newStaticFileRoots returns filesystems shaped like the two static routes:
// an in-memory root standing in for the embedded assets and an on-disk root
// standing in for the env JS directory. Both contain a file with
// pre-compressed variants and a sub directory.
func newStaticFileRoots(t *testing.T) (embedded, disk http.FileSystem) {
	t.Helper()

	embedded = http.FS(fstest.MapFS{
		"index-abc.js":      {Data: []byte("plain")},
		"index-abc.js.br":   {Data: []byte("brotli")},
		"index-abc.js.gz":   {Data: []byte("gzip")},
		"fonts/sofia.woff2": {Data: []byte("font")},
	})

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "env.js"), []byte("plain"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "env.js.br"), []byte("brotli"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "env.js.gz"), []byte("gzip"), 0o600))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "nested"), 0o700))
	disk = http.Dir(dir)
	return embedded, disk
}

func TestFileOnlyFileSystem_Open(t *testing.T) {
	t.Parallel()

	embedded, disk := newStaticFileRoots(t)
	cases := []struct {
		name    string
		root    http.FileSystem
		path    string
		wantErr error
	}{
		{name: "embedded root directory", root: embedded, path: "/", wantErr: fs.ErrNotExist},
		{name: "embedded sub directory", root: embedded, path: "/fonts", wantErr: fs.ErrNotExist},
		{name: "embedded missing file", root: embedded, path: "/missing.js", wantErr: fs.ErrNotExist},
		{name: "embedded file", root: embedded, path: "/index-abc.js"},
		{name: "embedded file in sub directory", root: embedded, path: "/fonts/sofia.woff2"},
		{name: "disk root directory", root: disk, path: "/", wantErr: fs.ErrNotExist},
		{name: "disk sub directory", root: disk, path: "/nested", wantErr: fs.ErrNotExist},
		{name: "disk missing file", root: disk, path: "/missing.js", wantErr: fs.ErrNotExist},
		{name: "disk file", root: disk, path: "/env.js"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ffs := &fileOnlyFileSystem{root: tc.root}
			f, err := ffs.Open(tc.path)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			defer f.Close()
			stat, err := f.Stat()
			require.NoError(t, err)
			assert.False(t, stat.IsDir())
		})
	}
}

func TestCompressedFileServer_directoryPathsReturn404(t *testing.T) {
	t.Parallel()

	embedded, disk := newStaticFileRoots(t)
	roots := []struct {
		name   string
		prefix string
		root   http.FileSystem
		file   string
		subdir string
	}{
		{name: "embedded assets", prefix: "/assets/", root: embedded, file: "index-abc.js", subdir: "fonts"},
		{name: "env js directory", prefix: "/static/js/", root: disk, file: "env.js", subdir: "nested"},
	}
	for _, root := range roots {
		// Mount the handler the same way DashboardService.Register does.
		mux := http.NewServeMux()
		mux.Handle(root.prefix, http.StripPrefix(root.prefix, compressedFileServer(root.root)))

		cases := []struct {
			name           string
			path           string
			acceptEncoding string
			wantCode       int
			wantBody       string
			wantEncoding   string
		}{
			{
				name:     "root directory returns 404",
				path:     root.prefix,
				wantCode: http.StatusNotFound,
			},
			{
				name:     "sub directory returns 404",
				path:     root.prefix + root.subdir,
				wantCode: http.StatusNotFound,
			},
			{
				name:     "sub directory with trailing slash returns 404",
				path:     root.prefix + root.subdir + "/",
				wantCode: http.StatusNotFound,
			},
			{
				name:     "missing file returns 404",
				path:     root.prefix + "missing.js",
				wantCode: http.StatusNotFound,
			},
			{
				name:     "file is served",
				path:     root.prefix + root.file,
				wantCode: http.StatusOK,
				wantBody: "plain",
			},
			{
				name:           "file is served pre-compressed with brotli",
				path:           root.prefix + root.file,
				acceptEncoding: "gzip, deflate, br",
				wantCode:       http.StatusOK,
				wantBody:       "brotli",
				wantEncoding:   "br",
			},
			{
				name:           "file is served pre-compressed with gzip",
				path:           root.prefix + root.file,
				acceptEncoding: "gzip",
				wantCode:       http.StatusOK,
				wantBody:       "gzip",
				wantEncoding:   "gzip",
			},
		}
		for _, tc := range cases {
			t.Run(root.name+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				req := httptest.NewRequest(http.MethodGet, tc.path, nil)
				if tc.acceptEncoding != "" {
					req.Header.Set("Accept-Encoding", tc.acceptEncoding)
				}
				rec := httptest.NewRecorder()
				mux.ServeHTTP(rec, req)
				assert.Equal(t, tc.wantCode, rec.Code)
				assert.Equal(t, tc.wantEncoding, rec.Header().Get("Content-Encoding"))
				if tc.wantCode == http.StatusOK {
					assert.Equal(t, tc.wantBody, rec.Body.String())
				}
			})
		}
	}
}

func TestDashboardService_Register_envJSDirectoryReturns404(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "env.js"), []byte("window.env = {};"), 0o600))
	mux := http.NewServeMux()
	NewDashboardService(dir).Register(mux)

	cases := []struct {
		name     string
		path     string
		wantCode int
	}{
		{name: "env js directory returns 404", path: "/static/js/", wantCode: http.StatusNotFound},
		{name: "env js file is served", path: "/static/js/env.js", wantCode: http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			assert.Equal(t, tc.wantCode, rec.Code)
		})
	}
}
