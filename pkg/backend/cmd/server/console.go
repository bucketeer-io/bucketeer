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
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var console embed.FS

func ConsoleHandler() (http.Handler, error) {
	consoleFiles, err := fs.Sub(console, "dist")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(consoleFiles)), nil
}
