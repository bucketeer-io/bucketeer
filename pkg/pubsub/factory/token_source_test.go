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

package factory

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

type stubTokenSource struct {
	tok *oauth2.Token
	err error
}

func (s *stubTokenSource) Token() (*oauth2.Token, error) {
	return s.tok, s.err
}

func TestResilientTokenSource(t *testing.T) {
	t.Parallel()
	errRefused := errors.New("connection refused")

	validToken := func(access string, expiry time.Duration) *oauth2.Token {
		return &oauth2.Token{AccessToken: access, Expiry: time.Now().Add(expiry)}
	}

	tests := []struct {
		desc    string
		setup   func(base *stubTokenSource, rts *resilientTokenSource)
		baseTok *oauth2.Token
		baseErr error
		wantTok string
		wantErr bool
	}{
		{
			desc:    "delegates to base on success",
			baseTok: validToken("good", time.Hour),
			wantTok: "good",
		},
		{
			desc: "falls back to cached token when base fails",
			setup: func(base *stubTokenSource, rts *resilientTokenSource) {
				base.tok = validToken("cached", time.Hour)
				_, _ = rts.Token()
			},
			baseErr: errRefused,
			wantTok: "cached",
		},
		{
			desc: "fails when cached token is expired",
			setup: func(base *stubTokenSource, rts *resilientTokenSource) {
				base.tok = validToken("old", -time.Minute)
				_, _ = rts.Token()
			},
			baseErr: errRefused,
			wantErr: true,
		},
		{
			desc: "updates cache and returns latest on fallback",
			setup: func(base *stubTokenSource, rts *resilientTokenSource) {
				base.tok = validToken("v1", time.Hour)
				_, _ = rts.Token()
				base.tok = validToken("v2", 2*time.Hour)
				_, _ = rts.Token()
			},
			baseErr: errRefused,
			wantTok: "v2",
		},
		{
			desc:    "no cache falls through with error",
			baseErr: errors.New("no credentials"),
			wantErr: true,
		},
		{
			desc: "nil token from base does not overwrite cache",
			setup: func(base *stubTokenSource, rts *resilientTokenSource) {
				base.tok = validToken("cached", time.Hour)
				_, _ = rts.Token()
			},
			// base returns (nil, nil)
			wantTok: "cached",
		},
		{
			desc: "nil token from base with no cache returns error",
			// base returns (nil, nil), no prior cache
			wantErr: true,
		},
		{
			desc: "zero expiry treated as never expires",
			setup: func(base *stubTokenSource, rts *resilientTokenSource) {
				base.tok = &oauth2.Token{AccessToken: "forever"}
				_, _ = rts.Token()
			},
			baseErr: errRefused,
			wantTok: "forever",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			base := &stubTokenSource{}
			rts := &resilientTokenSource{base: base}

			if tt.setup != nil {
				tt.setup(base, rts)
			}
			base.tok = tt.baseTok
			base.err = tt.baseErr

			got, err := rts.Token()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantTok, got.AccessToken)
		})
	}
}
