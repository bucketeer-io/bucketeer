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

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// expectedSecurityHeaders are the headers every gateway response must carry.
var expectedSecurityHeaders = map[string]string{
	"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	"X-Content-Type-Options":    "nosniff",
	"Cache-Control":             "no-store",
}

func assertSecurityHeaders(t *testing.T, header http.Header) {
	t.Helper()
	for key, value := range expectedSecurityHeaders {
		assert.Equal(t, []string{value}, header.Values(key), key)
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc       string
		handler    http.HandlerFunc
		wantStatus int
	}{
		{
			desc: "success response written by the handler",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			},
			wantStatus: http.StatusOK,
		},
		{
			desc: "error response written with http.Error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			desc: "response written without an explicit status",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("ok"))
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "http://example.com/v1/features", nil)
			rec := httptest.NewRecorder()
			securityHeadersMiddleware(p.handler).ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, p.wantStatus, res.StatusCode)
			assertSecurityHeaders(t, res.Header)
		})
	}
}

func TestGatewayResponsesIncludeSecurityHeaders(t *testing.T) {
	t.Parallel()
	baseURL := startTestGateway(t)

	patterns := []struct {
		desc       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{
			desc:       "handled route",
			method:     http.MethodPost,
			path:       "/v1/echo",
			body:       `{"name":"ok"}`,
			wantStatus: http.StatusOK,
		},
		{
			desc:       "routing error written by the mux",
			method:     http.MethodGet,
			path:       "/v1/unknown",
			wantStatus: http.StatusNotFound,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(p.method, baseURL+p.path, strings.NewReader(p.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, p.wantStatus, resp.StatusCode)
			assertSecurityHeaders(t, resp.Header)
		})
	}
}
