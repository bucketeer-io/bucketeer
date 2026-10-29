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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestAllowedIncomingHeaderMatcher(t *testing.T) {
	t.Parallel()
	matcher := AllowedIncomingHeaderMatcher("authorization", "X-Request-Id")

	patterns := []struct {
		desc        string
		key         string
		expectedKey string
		expectedOK  bool
	}{
		{
			desc:        "allowed header in canonical form",
			key:         "Authorization",
			expectedKey: "authorization",
			expectedOK:  true,
		},
		{
			desc:        "allowed header in lower case",
			key:         "authorization",
			expectedKey: "authorization",
			expectedOK:  true,
		},
		{
			desc:        "allowed header in upper case",
			key:         "X-REQUEST-ID",
			expectedKey: "x-request-id",
			expectedOK:  true,
		},
		{
			desc:       "grpc-metadata prefixed header is not forwarded",
			key:        "Grpc-Metadata-Apikey-Token",
			expectedOK: false,
		},
		{
			desc:       "grpc-metadata prefixed allowed header is not forwarded",
			key:        "Grpc-Metadata-Authorization",
			expectedOK: false,
		},
		{
			desc:       "internal metadata key sent as a plain header is not forwarded",
			key:        "Apikey-Maintainer",
			expectedOK: false,
		},
		{
			desc:       "permanent http header is not forwarded",
			key:        "Cookie",
			expectedOK: false,
		},
		{
			desc:       "empty key is not forwarded",
			key:        "",
			expectedOK: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			key, ok := matcher(p.key)
			assert.Equal(t, p.expectedOK, ok)
			assert.Equal(t, p.expectedKey, key)
		})
	}
}

func TestAllowedIncomingHeaderMatcherAnnotateContext(t *testing.T) {
	t.Parallel()

	newRequest := func() *http.Request {
		req := httptest.NewRequest(http.MethodPost, "/v1/feature", nil)
		req.Header.Set("Authorization", "Bearer token")
		req.Header.Set("Cookie", "session=abc")
		req.Header.Set("Grpc-Metadata-Apikey-Token", "forged-token")
		req.Header.Set("Grpc-Metadata-Apikey-Maintainer", "someone@example.com")
		req.Header.Set("Grpc-Metadata-Apikey-Name", "forged-name")
		req.Header.Set("Apikey-Token", "forged-token")
		return req
	}

	patterns := []struct {
		desc            string
		muxOpts         []runtime.ServeMuxOption
		expectedPresent []string
		expectedAbsent  []string
	}{
		{
			desc:            "default matcher forwards grpc-metadata prefixed headers",
			muxOpts:         nil,
			expectedPresent: []string{"authorization", "apikey-token", "apikey-maintainer", "apikey-name", "grpcgateway-cookie"},
			expectedAbsent:  nil,
		},
		{
			desc: "allowlist matcher forwards only the allowed headers",
			muxOpts: []runtime.ServeMuxOption{
				runtime.WithIncomingHeaderMatcher(AllowedIncomingHeaderMatcher("authorization")),
			},
			expectedPresent: []string{"authorization"},
			expectedAbsent:  []string{"apikey-token", "apikey-maintainer", "apikey-name", "grpcgateway-cookie", "cookie"},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			mux := runtime.NewServeMux(p.muxOpts...)
			ctx, err := runtime.AnnotateContext(
				context.Background(),
				mux,
				newRequest(),
				"/bucketeer.feature.FeatureService/CreateFeature",
			)
			require.NoError(t, err)
			md, ok := metadata.FromOutgoingContext(ctx)
			require.True(t, ok)
			for _, key := range p.expectedPresent {
				assert.NotEmpty(t, md.Get(key), "expected metadata %q to be forwarded", key)
			}
			for _, key := range p.expectedAbsent {
				assert.Empty(t, md.Get(key), "expected metadata %q not to be forwarded", key)
			}
			assert.Contains(t, md.Get("authorization"), "Bearer token")
		})
	}
}
