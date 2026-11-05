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

package gateway

import "net/http"

// securityHeaders are added to every response served by the gateway. Where
// the REST server or the docker-compose nginx config already sets a header,
// the value matches it, so deployments behind them do not receive
// conflicting duplicates.
var securityHeaders = []struct {
	key   string
	value string
}{
	// The gateway is reached over TLS; browsers must not fall back to plain HTTP.
	{key: "Strict-Transport-Security", value: "max-age=31536000; includeSubDomains"},
	// Responses are JSON; browsers must not sniff them into another type.
	{key: "X-Content-Type-Options", value: "nosniff"},
	// Responses are authenticated and specific to the caller, so neither a
	// browser nor a shared cache may store them.
	{key: "Cache-Control", value: "no-store"},
}

// securityHeadersMiddleware sets the security headers before the request
// reaches the gateway mux, so every response (including error responses
// written by the mux or by a handler) carries them.
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		for _, h := range securityHeaders {
			header.Set(h.key, h.value)
		}
		next.ServeHTTP(w, r)
	})
}
