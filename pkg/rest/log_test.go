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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogServerMiddleware(t *testing.T) {
	t.Parallel()

	handlerFor := func(logger *zap.Logger) http.Handler {
		return LogServerMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/redirect":
				w.WriteHeader(http.StatusMovedPermanently)
			case "/not-found":
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusOK)
			}
		}))
	}

	patterns := []struct {
		desc          string
		path          string
		expectedCount int
	}{
		{
			desc:          "success: skip logging for 200",
			path:          "/ok",
			expectedCount: 0,
		},
		{
			desc:          "success: skip logging for 301 redirect",
			path:          "/redirect",
			expectedCount: 0,
		},
		{
			desc:          "warn: log 404",
			path:          "/not-found",
			expectedCount: 1,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(&buf),
				zapcore.WarnLevel,
			)
			handler := handlerFor(zap.New(core))

			req := httptest.NewRequest(http.MethodGet, p.path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			out := strings.TrimSpace(buf.String())
			var logLines []string
			if out != "" {
				logLines = strings.Split(out, "\n")
			}
			require.Len(t, logLines, p.expectedCount)
			if p.expectedCount == 0 {
				return
			}
			var logEntry map[string]any
			require.NoError(t, json.Unmarshal([]byte(logLines[0]), &logEntry))
			assert.Equal(t, float64(http.StatusNotFound), logEntry["statusCode"])
		})
	}
}

func TestDecodeBody(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		body        io.Reader
		expected    interface{}
		expectedErr bool
	}{
		{
			desc:        "err: not json",
			body:        strings.NewReader(`{tag: "ios", user: {id: "pingdom", data: {foo: "bar"}}}`),
			expected:    nil,
			expectedErr: true,
		},
		{
			desc:        "success: nil",
			body:        bytes.NewReader(nil),
			expected:    nil,
			expectedErr: false,
		},
		{
			desc: "success: json",
			body: strings.NewReader(`{"tag":"ios","user":{"id":"pingdom","data":{"foo":"bar"}}}`),
			expected: map[string]interface{}{
				"tag": "ios",
				"user": map[string]interface{}{
					"id": "pingdom",
					"data": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			expectedErr: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			decoded, err := decodeBody(p.body)
			assert.Equal(t, p.expected, decoded)
			assert.Equal(t, p.expectedErr, err != nil)
		})
	}
}
