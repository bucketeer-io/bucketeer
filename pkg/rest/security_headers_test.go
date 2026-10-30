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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedSecurityHeaders = map[string]string{
	"X-Frame-Options":           "DENY",
	"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	"X-Content-Type-Options":    "nosniff",
	"Referrer-Policy":           "strict-origin-when-cross-origin",
	"Content-Security-Policy":   contentSecurityPolicy,
}

func assertSecurityHeaders(t *testing.T, header http.Header) {
	t.Helper()
	for key, value := range expectedSecurityHeaders {
		assert.Equal(t, []string{value}, header.Values(key), key)
	}
	// The policy is enforced; it must not be sent as report-only as well,
	// which would make browsers evaluate and report it twice.
	assert.Empty(t, header.Values("Content-Security-Policy-Report-Only"))
}

// parseContentSecurityPolicy splits a serialized policy into its directives,
// keyed by directive name with the source list as value.
func parseContentSecurityPolicy(t *testing.T, policy string) map[string]string {
	t.Helper()
	directives := map[string]string{}
	for _, directive := range strings.Split(policy, ";") {
		fields := strings.Fields(directive)
		require.NotEmpty(t, fields, "empty directive in %q", policy)
		name := fields[0]
		_, duplicated := directives[name]
		require.False(t, duplicated, "directive %q appears twice in %q", name, policy)
		directives[name] = strings.Join(fields[1:], " ")
	}
	return directives
}

func TestContentSecurityPolicyDirectives(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc      string
		directive string
		want      string
	}{
		{
			desc:      "everything not listed below is restricted to the dashboard origin",
			directive: "default-src",
			want:      "'self'",
		},
		{
			desc:      "scripts and workers only come from the dashboard origin",
			directive: "script-src",
			want:      "'self'",
		},
		{
			desc:      "inline style attributes and injected style elements are allowed",
			directive: "style-src",
			want:      "'self' 'unsafe-inline'",
		},
		{
			desc:      "inlined assets, uploaded avatars and identity provider pictures load",
			directive: "img-src",
			want:      "'self' data: https:",
		},
		{
			desc:      "the same-origin API, SSE stream and inlined avatar images are reachable",
			directive: "connect-src",
			want:      "'self' data:",
		},
		{
			desc:      "the dashboard cannot be framed",
			directive: "frame-ancestors",
			want:      "'none'",
		},
		{
			desc:      "plugins are not allowed",
			directive: "object-src",
			want:      "'none'",
		},
		{
			desc:      "the document base URL cannot be changed",
			directive: "base-uri",
			want:      "'self'",
		},
	}
	directives := parseContentSecurityPolicy(t, contentSecurityPolicy)
	assert.Len(t, directives, len(patterns), "unexpected directive in %q", contentSecurityPolicy)
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			value, ok := directives[p.directive]
			require.True(t, ok, "directive %q is missing", p.directive)
			assert.Equal(t, p.want, value)
		})
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc            string
		handler         http.HandlerFunc
		wantStatus      int
		wantContentType string
	}{
		{
			desc: "success response keeps the content type set by the handler",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("<!doctype html>"))
			},
			wantStatus:      http.StatusOK,
			wantContentType: "text/html; charset=utf-8",
		},
		{
			desc: "error response written with http.Error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			},
			wantStatus:      http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			desc: "response written without an explicit status",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("ok"))
			},
			wantStatus:      http.StatusOK,
			wantContentType: "text/plain; charset=utf-8",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, dummyURL, nil)
			rec := httptest.NewRecorder()
			securityHeadersMiddleware(p.handler).ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, p.wantStatus, res.StatusCode)
			assert.Equal(t, p.wantContentType, res.Header.Get("Content-Type"))
			assertSecurityHeaders(t, res.Header)
		})
	}
}

func TestServerResponsesIncludeSecurityHeaders(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc       string
		path       string
		wantStatus int
	}{
		{
			desc:       "unknown path answered by the mux",
			path:       "/",
			wantStatus: http.StatusNotFound,
		},
		{
			desc:       "path rejected by the request path middleware",
			path:       "/.env",
			wantStatus: http.StatusNotFound,
		},
		{
			desc:       "malformed path rejected by the request path middleware",
			path:       "/foo%85",
			wantStatus: http.StatusBadRequest,
		},
	}
	client := newHTTPClient()
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			resp, err := client.Get(fmt.Sprintf("https://localhost:%d%s", port, p.path))
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, p.wantStatus, resp.StatusCode)
			assertSecurityHeaders(t, resp.Header)
		})
	}
}
