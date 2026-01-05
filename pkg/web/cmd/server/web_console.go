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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bucketeer-io/bucketeer/v2/ui/dashboard"
)

type spaFileSystem struct {
	root   http.FileSystem
	prefix string
}

// Open method for spaFileSystem
func (fs *spaFileSystem) Open(name string) (http.File, error) {
	// First try with original path
	f, err := fs.root.Open(name)
	if !os.IsNotExist(err) {
		return f, err
	}

	// If file not found and we have a prefix, try stripping it
	if fs.prefix != "" && strings.HasPrefix(name, fs.prefix) {
		strippedName := strings.TrimPrefix(name, fs.prefix)
		f, err = fs.root.Open(strippedName)
		if !os.IsNotExist(err) {
			return f, err
		}
	}

	// If still not found, return index.html
	return fs.root.Open("index.html")
}

// fontCacheHandler adds cache headers for font files
func cacheHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add cache headers for font files
		if isFontFile(r.URL.Path) {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			w.Header().Set("Vary", "Accept-Encoding")
		}
		next.ServeHTTP(w, r)
	})
}

func isFontFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".woff2" || ext == ".woff" || ext == ".ttf" || ext == ".otf"
}

// compressedFileServer serves pre-compressed Brotli/Gzip assets when available.
func compressedFileServer(root http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		ae := r.Header.Get("Accept-Encoding")
		// Brotli
		if strings.Contains(ae, "br") {
			if f, err := root.Open(path + ".br"); err == nil {
				defer f.Close()
				w.Header().Set("Content-Encoding", "br")
				w.Header().Set("Vary", "Accept-Encoding")
				http.ServeContent(w, r, path, time.Time{}, f)
				return
			}
		}
		// Gzip
		if strings.Contains(ae, "gzip") {
			if f, err := root.Open(path + ".gz"); err == nil {
				defer f.Close()
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Vary", "Accept-Encoding")
				http.ServeContent(w, r, path, time.Time{}, f)
				return
			}
		}
		// Fallback
		http.FileServer(root).ServeHTTP(w, r)
	})
}

// dashboardHandler returns a http.Handler for the new dashboard UI.
func dashboardHandler() http.Handler {
	fs := http.FS(dashboard.FS)
	return http.FileServer(&spaFileSystem{root: fs})
}

type DashboardService struct {
	consoleEnvJSPath string
}

func NewDashboardService(consoleEnvJSPath string) DashboardService {
	return DashboardService{consoleEnvJSPath: consoleEnvJSPath}
}

// Register sets up handlers for assets, fonts, the SPA, and env-JS.
func (d DashboardService) Register(mux *http.ServeMux) {
	// Subtree for embedded assets
	embedded := dashboard.FS
	assetsSub, err := fs.Sub(embedded, "assets")
	if err != nil {
		panic("failed to get embedded dashboard assets: " + err.Error())
	}
	assetsFS := http.FS(assetsSub)

	// Serve pre-compressed static assets (JS/CSS/images/fonts) with font caching
	mux.Handle("/assets/", http.StripPrefix(
		"/assets/",
		cacheHandler(
			compressedFileServer(assetsFS),
		),
	))

	// Serve dynamic env JS with compression
	mux.Handle("/static/js/", http.StripPrefix(
		"/static/js/",
		compressedFileServer(http.Dir(d.consoleEnvJSPath)),
	))

	// SPA entry point with index.html fallback
	mux.Handle("/", dashboardHandler())
}
