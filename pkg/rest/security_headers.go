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

import (
	"net/http"
	"strings"
)

// contentSecurityPolicy lists the sources the dashboard (ui/dashboard) loads
// resources from. Every deployment serves the dashboard, the web gRPC gateway
// (/v1/) and the AI chat SSE stream (/v1/aichat/chat) from the same origin,
// so 'self' covers all API traffic; the sign-in flow is a top-level redirect
// to the identity provider and back, which CSP does not restrict.
var contentSecurityPolicy = strings.Join([]string{
	"default-src 'self'",
	// The build emits only same-origin module scripts and worker files, and
	// /static/js/env.js is served by this server. Nothing turns strings into
	// code at runtime (the JSON Schema validation of variation values
	// interprets the schema instead of generating a validator function), so
	// 'unsafe-eval' is not needed.
	"script-src 'self'",
	// React inline style attributes and the style elements injected at
	// runtime by the code editor and the select components.
	"style-src 'self' 'unsafe-inline'",
	// Small assets are inlined as data: URLs by the build, uploaded avatars are
	// rendered as data: URLs, and identity provider account pictures are
	// loaded from their https URL.
	"img-src 'self' data: https:",
	// The avatar picker converts the inlined avatar images with fetch(),
	// which reads a data: URL.
	"connect-src 'self' data:",
	// Same restriction as X-Frame-Options: DENY.
	"frame-ancestors 'none'",
	"object-src 'none'",
	"base-uri 'self'",
}, "; ")

// securityHeaders are added to every response served by the REST server.
// Where the docker-compose nginx config already sets a header, the value
// matches it, so deployments behind it do not receive conflicting duplicates.
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
	// Browsers block every resource load, connection, inline script and
	// framing attempt the policy does not allow.
	{key: "Content-Security-Policy", value: contentSecurityPolicy},
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
