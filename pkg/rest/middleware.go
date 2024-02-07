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

package rest

import (
	"bytes"
	"net/http"
	"strings"
)

type middleware func(http.Handler) http.Handler

type middlewares struct {
	mw []middleware
}

func newMiddleWares() *middlewares {
	return &middlewares{}
}

func (ms *middlewares) Append(mw middleware) *middlewares {
	ms.mw = append(ms.mw, mw)
	return ms
}

func (ms *middlewares) Handle(handler http.Handler) http.Handler {
	next := handler
	for i := len(ms.mw) - 1; i >= 0; i-- {
		next = ms.mw[i](next)
	}
	return next
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.body.Write(b)
	return rr.ResponseWriter.Write(b)
}

func splitURLPath(path string) (string, string, string) {
	// format: /api_version/service_name/api_name
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		return "unknown", "unknown", "unknown"
	}
	return parts[1], parts[2], parts[3]
}
