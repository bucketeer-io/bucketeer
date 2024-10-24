// Copyright 2024 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/ui/dashboard"
	webv2 "github.com/bucketeer-io/bucketeer/ui/web-v2"
)

type spaFileSystem struct {
	root http.FileSystem
}

func (fs *spaFileSystem) Open(name string) (http.File, error) {
	f, err := fs.root.Open(name)
	if os.IsNotExist(err) {
		return fs.root.Open("index.html")
	}
	return f, err
}

// webConsoleHandler returns a http.Handler for the old web console UI.
func webConsoleHandler() http.Handler {
	return http.FileServer(&spaFileSystem{http.FS(webv2.FS)})
}

// dashboardHandler returns a http.Handler for the new dashboard UI.
func dashboardHandler() http.Handler {
	return http.FileServer(&spaFileSystem{http.FS(dashboard.FS)})
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
	mux.HandleFunc("/v3", dashboardHandler().ServeHTTP)
	mux.HandleFunc("/static/js/",
		http.StripPrefix("/static/js/", webConsoleEnvJSHandler(c.consoleEnvJSPath)).ServeHTTP)
}
