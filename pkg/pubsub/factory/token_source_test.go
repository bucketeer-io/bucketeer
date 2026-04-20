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

func TestResilientTokenSource_DelegatesOnSuccess(t *testing.T) {
	want := &oauth2.Token{
		AccessToken: "good",
		Expiry:      time.Now().Add(time.Hour),
	}
	rts := &resilientTokenSource{base: &stubTokenSource{tok: want}}

	got, err := rts.Token()
	require.NoError(t, err)
	assert.Equal(t, want.AccessToken, got.AccessToken)
}

func TestResilientTokenSource_FallsBackToCachedToken(t *testing.T) {
	good := &oauth2.Token{
		AccessToken: "cached",
		Expiry:      time.Now().Add(time.Hour),
	}
	base := &stubTokenSource{tok: good}
	rts := &resilientTokenSource{base: base}

	// First call succeeds and caches.
	_, err := rts.Token()
	require.NoError(t, err)

	// Base starts failing (metadata server gone).
	base.tok = nil
	base.err = errors.New("connection refused")

	got, err := rts.Token()
	require.NoError(t, err)
	assert.Equal(t, "cached", got.AccessToken)
}

func TestResilientTokenSource_FailsWhenCachedTokenExpired(t *testing.T) {
	expired := &oauth2.Token{
		AccessToken: "old",
		Expiry:      time.Now().Add(-time.Minute),
	}
	base := &stubTokenSource{tok: expired}
	rts := &resilientTokenSource{base: base}

	// Cache an already-expired token.
	_, _ = rts.Token()

	// Base starts failing.
	base.tok = nil
	base.err = errors.New("connection refused")

	_, err := rts.Token()
	assert.Error(t, err)
}

func TestResilientTokenSource_UpdatesCacheOnNewToken(t *testing.T) {
	tok1 := &oauth2.Token{
		AccessToken: "v1",
		Expiry:      time.Now().Add(time.Hour),
	}
	tok2 := &oauth2.Token{
		AccessToken: "v2",
		Expiry:      time.Now().Add(2 * time.Hour),
	}

	base := &stubTokenSource{tok: tok1}
	rts := &resilientTokenSource{base: base}

	got1, _ := rts.Token()
	assert.Equal(t, "v1", got1.AccessToken)

	base.tok = tok2
	got2, _ := rts.Token()
	assert.Equal(t, "v2", got2.AccessToken)

	// Fail base — should return v2, not v1.
	base.tok = nil
	base.err = errors.New("fail")
	got3, err := rts.Token()
	require.NoError(t, err)
	assert.Equal(t, "v2", got3.AccessToken)
}

func TestResilientTokenSource_NoCacheFallsThrough(t *testing.T) {
	base := &stubTokenSource{err: errors.New("no credentials")}
	rts := &resilientTokenSource{base: base}

	_, err := rts.Token()
	assert.Error(t, err)
}
