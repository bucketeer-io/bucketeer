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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRejectInvalidRequestPathMiddleware(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc       string
		target     string
		wantStatus int
		wantNext   bool
	}{
		{
			desc:       "valid dashboard path",
			target:     "http://example.com/assets/app.js",
			wantStatus: http.StatusOK,
			wantNext:   true,
		},
		{
			desc:       "malformed percent encoding in request URI",
			target:     "http://example.com/.env.local%85",
			wantStatus: http.StatusBadRequest,
			wantNext:   false,
		},
		{
			desc:       "secret probe path",
			target:     "http://example.com/credentials.json",
			wantStatus: http.StatusNotFound,
			wantNext:   false,
		},
		{
			desc:       "env probe path",
			target:     "http://example.com/.env.local",
			wantStatus: http.StatusNotFound,
			wantNext:   false,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			nextCalled := false
			handler := rejectInvalidRequestPathMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			}))
			req := httptest.NewRequest(http.MethodGet, p.target, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, p.wantStatus, rec.Code)
			assert.Equal(t, p.wantNext, nextCalled)
		})
	}
}

func TestInvalidRequestPathStatus(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://example.com/terraform.tfstate%85", nil)
	status, reject := invalidRequestPathStatus(req)
	require.True(t, reject)
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestInvalidRequestPathStatus_incompletePercentEncoding(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.RequestURI = "/foo%2"
	status, reject := invalidRequestPathStatus(req)
	require.True(t, reject)
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestInvalidRequestPathStatus_utf8C1ControlChar(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo%C2%85", nil)
	status, reject := invalidRequestPathStatus(req)
	require.True(t, reject)
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestRESTMiddlewareChain_invalidUTF8PathSegment(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("application handler must not run for invalid path")
	})
	mws := newMiddleWares()
	mws.Append(MetricsServerMiddleware)
	mws.Append(rejectInvalidRequestPathMiddleware)
	handler := mws.Handle(mux)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/v1/foo/%85", nil)
	rec := httptest.NewRecorder()
	require.NotPanics(t, func() {
		handler.ServeHTTP(rec, req)
	})
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
