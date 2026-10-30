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

package rest

import "net/http"

// securityHeaders are added to every response served by the REST server.
// The values match the ones already set by the docker-compose nginx config,
// so deployments behind it do not receive conflicting duplicates.
var securityHeaders = []struct {
	key   string
	value string
}{
	// The dashboard must not be embedded in a frame on another origin.
	{key: "X-Frame-Options", value: "DENY"},
	// The server only listens with TLS; browsers must not fall back to plain HTTP.
	{key: "Strict-Transport-Security", value: "max-age=31536000; includeSubDomains"},
	// Scripts and stylesheets are only used when served with the declared MIME type.
	{key: "X-Content-Type-Options", value: "nosniff"},
	// Only the origin is sent as referrer to other origins, nothing on HTTPS->HTTP.
	{key: "Referrer-Policy", value: "strict-origin-when-cross-origin"},
}

// securityHeadersMiddleware sets the browser security headers before the
// request reaches the application handlers, so every response (including
// error responses written by later middlewares or handlers) carries them.
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		for _, h := range securityHeaders {
			header.Set(h.key, h.value)
		}
		next.ServeHTTP(w, r)
	})
}
