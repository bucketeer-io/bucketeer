// Copyright 2025 The Bucketeer Authors.
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
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bucketeer-io/bucketeer/ui/dashboard"
	webv2 "github.com/bucketeer-io/bucketeer/ui/web-v2"
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

// webConsoleHandler returns a http.Handler for the old web console UI.
func webConsoleHandler() http.Handler {
	return http.FileServer(&spaFileSystem{root: http.FS(webv2.FS), prefix: "/legacy/"})
}

// fontCacheHandler adds cache headers for font files
func fontCacheHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add cache headers for font files
		if isFontFile(r.URL.Path) {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			w.Header().Set("ETag", `"font-v1"`)
		}
		next.ServeHTTP(w, r)
	})
}

func isFontFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".woff2" || ext == ".woff" || ext == ".ttf" || ext == ".otf"
}

// dashboardHandler returns a http.Handler for the new dashboard UI.
func dashboardHandler() http.Handler {
	return fontCacheHandler(http.FileServer(&spaFileSystem{root: http.FS(dashboard.FS)}))
}

func webConsoleEnvJSHandler(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}

type WebConsoleService struct {
	consoleEnvJSPath string
}

func NewWebConsoleService(consoleEnvJSPath string) WebConsoleService {
	return WebConsoleService{consoleEnvJSPath: consoleEnvJSPath}
}

func (c WebConsoleService) Register(mux *http.ServeMux) {
	mux.HandleFunc("/", webConsoleHandler().ServeHTTP)
	mux.HandleFunc("/legacy/static/js/",
		http.StripPrefix("/legacy/static/js/", webConsoleEnvJSHandler(c.consoleEnvJSPath)).ServeHTTP)
}

type DashboardService struct {
	consoleEnvJSPath string
}

func NewDashboardService(consoleEnvJSPath string) DashboardService {
	return DashboardService{consoleEnvJSPath: consoleEnvJSPath}
}

func (d DashboardService) Register(mux *http.ServeMux) {
	mux.HandleFunc("/", dashboardHandler().ServeHTTP)
	mux.HandleFunc("/static/js/",
		http.StripPrefix("/static/js/", webConsoleEnvJSHandler(d.consoleEnvJSPath)).ServeHTTP)
}
