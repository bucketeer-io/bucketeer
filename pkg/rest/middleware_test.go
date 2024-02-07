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
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const firstKey = "first"
const secondKey = "second"
const dummyURL = "http://example.com"

func TestHandle(t *testing.T) {
	t.Parallel()
	var firstRun, secondRun, handlerRun bool
	first := func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add(firstKey, firstKey)
				firstRun = true
				next.ServeHTTP(w, r)
			},
		)
	}
	second := func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, w.Header().Get(firstKey), firstKey)
				w.Header().Add(secondKey, secondKey)
				secondRun = true
				next.ServeHTTP(w, r)
			},
		)
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerRun = true
		require.Equal(t, w.Header().Get(firstKey), firstKey)
		require.Equal(t, w.Header().Get(secondKey), secondKey)
	}
	mws := newMiddleWares()
	mws.Append(first)
	mws.Append(second)
	handlers := mws.Handle(http.HandlerFunc(handler))
	req := httptest.NewRequest(http.MethodGet, dummyURL, nil)
	w := httptest.NewRecorder()
	handlers.ServeHTTP(w, req)
	assert.True(t, firstRun)
	assert.True(t, secondRun)
	assert.True(t, handlerRun)
}

func TestSplitURLPath(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc,
		input,
		expectedApiVersion,
		expectedServiceName,
		expectedApiName string
	}{
		{
			desc:                "error: wrong path format",
			input:               "scheme://host/api_version/service_name/api_name/api/",
			expectedApiVersion:  "unknown",
			expectedServiceName: "unknown",
			expectedApiName:     "unknown",
		},
		{
			desc:                "error: using slash in the end of the path",
			input:               "scheme://host/api_version/service_name/api_name/",
			expectedApiVersion:  "unknown",
			expectedServiceName: "unknown",
			expectedApiName:     "unknown",
		},
		{
			desc:                "sucess",
			input:               "scheme://host/api_version/service_name/api_name",
			expectedApiVersion:  "api_version",
			expectedServiceName: "service_name",
			expectedApiName:     "api_name",
		},
	}
	for _, p := range patterns {
		url, err := url.Parse(p.input)
		assert.NoError(t, err)
		apiVersion, serviceName, apiName := splitURLPath(url.Path)
		assert.Equal(t, apiVersion, p.expectedApiVersion, p.desc)
		assert.Equal(t, serviceName, p.expectedServiceName, p.desc)
		assert.Equal(t, apiName, p.expectedApiName, p.desc)
	}
}
