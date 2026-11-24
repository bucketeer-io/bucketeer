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

package token

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tokentesting "github.com/bucketeer-io/bucketeer/v2/pkg/token/testing"
)

const (
	optionsTestIssuer   = "test_issuer"
	optionsTestAudience = "test_audience"
)

func signTestServiceToken(t *testing.T, signer Signer, service string, issuedAt time.Time) string {
	t.Helper()
	raw, err := signer.SignServiceToken(&ServiceToken{
		Issuer:   optionsTestIssuer,
		Subject:  service,
		Audience: optionsTestAudience,
		Expiry:   issuedAt.Add(time.Hour),
		IssuedAt: issuedAt,
		ID:       "token-1",
	})
	require.NoError(t, err)
	return raw
}

func signTestAccessToken(t *testing.T, signer Signer, issuedAt time.Time, service bool) string {
	t.Helper()
	raw, err := signer.SignAccessToken(&AccessToken{
		Issuer:         optionsTestIssuer,
		Audience:       optionsTestAudience,
		Email:          "user@example.com",
		Expiry:         issuedAt.Add(time.Hour),
		IssuedAt:       issuedAt,
		IsSystemAdmin:  service,
		IsServiceToken: service,
	})
	require.NoError(t, err)
	return raw
}

func TestVerifyServiceTokenNotBefore(t *testing.T) {
	t.Parallel()
	privateKeyPath, publicKeyPath := tokentesting.WriteKeyPair(t)
	signer, err := NewSigner(privateKeyPath)
	require.NoError(t, err)
	notBefore := time.Now().Truncate(time.Second)
	verifier, err := NewVerifier(publicKeyPath, optionsTestIssuer, optionsTestAudience,
		WithServiceTokenNotBefore(map[string]time.Time{"batch": notBefore}),
		WithLegacyServiceTokenNotBefore(notBefore),
	)
	require.NoError(t, err)
	patterns := []struct {
		desc        string
		raw         string
		verify      func(string) error
		expectedErr string
	}{
		{
			desc: "success: service token issued at the time",
			raw:  signTestServiceToken(t, signer, "batch", notBefore),
			verify: func(raw string) error {
				_, err := verifier.VerifyServiceToken(raw)
				return err
			},
		},
		{
			desc: "err: service token issued before the time",
			raw:  signTestServiceToken(t, signer, "batch", notBefore.Add(-time.Second)),
			verify: func(raw string) error {
				_, err := verifier.VerifyServiceToken(raw)
				return err
			},
			expectedErr: `service token of "batch" issued before`,
		},
		{
			desc: "success: service token of another service issued before the time",
			raw:  signTestServiceToken(t, signer, "api", notBefore.Add(-time.Second)),
			verify: func(raw string) error {
				_, err := verifier.VerifyServiceToken(raw)
				return err
			},
		},
		{
			desc: "success: legacy service token issued at the time",
			raw:  signTestAccessToken(t, signer, notBefore, true),
			verify: func(raw string) error {
				_, err := verifier.VerifyAccessToken(raw)
				return err
			},
		},
		{
			desc: "err: legacy service token issued before the time",
			raw:  signTestAccessToken(t, signer, notBefore.Add(-time.Second), true),
			verify: func(raw string) error {
				_, err := verifier.VerifyAccessToken(raw)
				return err
			},
			expectedErr: "service token issued before",
		},
		{
			desc: "success: access token of a user issued before the time",
			raw:  signTestAccessToken(t, signer, notBefore.Add(-time.Second), false),
			verify: func(raw string) error {
				_, err := verifier.VerifyAccessToken(raw)
				return err
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			err := p.verify(p.raw)
			if p.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), p.expectedErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestVerifyServiceTokenWithoutNotBefore(t *testing.T) {
	t.Parallel()
	privateKeyPath, publicKeyPath := tokentesting.WriteKeyPair(t)
	signer, err := NewSigner(privateKeyPath)
	require.NoError(t, err)
	verifier, err := NewVerifier(publicKeyPath, optionsTestIssuer, optionsTestAudience)
	require.NoError(t, err)
	old := time.Now().Add(-time.Minute)
	_, err = verifier.VerifyServiceToken(signTestServiceToken(t, signer, "batch", old))
	assert.NoError(t, err)
	_, err = verifier.VerifyAccessToken(signTestAccessToken(t, signer, old, true))
	assert.NoError(t, err)
}

func TestVerifyServiceTokenPublicKey(t *testing.T) {
	t.Parallel()
	sharedPrivateKeyPath, sharedPublicKeyPath := tokentesting.WriteKeyPair(t)
	servicePrivateKeyPath, servicePublicKeyPath := tokentesting.WriteKeyPair(t)
	nextPrivateKeyPath, nextPublicKeyPath := tokentesting.WriteKeyPair(t)
	sharedSigner, err := NewSigner(sharedPrivateKeyPath)
	require.NoError(t, err)
	serviceSigner, err := NewSigner(servicePrivateKeyPath)
	require.NoError(t, err)
	nextSigner, err := NewSigner(nextPrivateKeyPath)
	require.NoError(t, err)
	// One file with the current and the next service token key.
	bothPublicKeysPath := filepath.Join(t.TempDir(), "service-token-public.pem")
	current, err := os.ReadFile(servicePublicKeyPath)
	require.NoError(t, err)
	next, err := os.ReadFile(nextPublicKeyPath)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(bothPublicKeysPath, append(current, next...), 0600))

	now := time.Now()
	patterns := []struct {
		desc        string
		opts        []VerifierOption
		raw         string
		verify      string
		expectedErr string
	}{
		{
			desc:   "success: without a service token key, the shared key verifies a service token",
			raw:    signTestServiceToken(t, sharedSigner, "batch", now),
			verify: TokenUseService,
		},
		{
			desc:   "success: service token signed with the service token key",
			opts:   []VerifierOption{WithServiceTokenPublicKey(servicePublicKeyPath)},
			raw:    signTestServiceToken(t, serviceSigner, "batch", now),
			verify: TokenUseService,
		},
		{
			desc:        "err: service token signed with the shared key is rejected once a service token key is set",
			opts:        []VerifierOption{WithServiceTokenPublicKey(servicePublicKeyPath)},
			raw:         signTestServiceToken(t, sharedSigner, "batch", now),
			verify:      TokenUseService,
			expectedErr: "invalid jwt",
		},
		{
			desc:        "err: access token signed with the service token key is rejected",
			opts:        []VerifierOption{WithServiceTokenPublicKey(servicePublicKeyPath)},
			raw:         signTestAccessToken(t, serviceSigner, now, false),
			verify:      TokenUseAccess,
			expectedErr: "invalid jwt",
		},
		{
			desc:        "err: legacy service token signed with the service token key is rejected",
			opts:        []VerifierOption{WithServiceTokenPublicKey(servicePublicKeyPath)},
			raw:         signTestAccessToken(t, serviceSigner, now, true),
			verify:      TokenUseAccess,
			expectedErr: "invalid jwt",
		},
		{
			desc:   "success: legacy service token is still verified with the shared key",
			opts:   []VerifierOption{WithServiceTokenPublicKey(servicePublicKeyPath)},
			raw:    signTestAccessToken(t, sharedSigner, now, true),
			verify: TokenUseAccess,
		},
		{
			desc:   "success: a file with two keys accepts the current key",
			opts:   []VerifierOption{WithServiceTokenPublicKey(bothPublicKeysPath)},
			raw:    signTestServiceToken(t, serviceSigner, "batch", now),
			verify: TokenUseService,
		},
		{
			desc:   "success: a file with two keys accepts the next key",
			opts:   []VerifierOption{WithServiceTokenPublicKey(bothPublicKeysPath)},
			raw:    signTestServiceToken(t, nextSigner, "batch", now),
			verify: TokenUseService,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			verifier, err := NewVerifier(sharedPublicKeyPath, optionsTestIssuer, optionsTestAudience, p.opts...)
			require.NoError(t, err)
			switch p.verify {
			case TokenUseService:
				_, err = verifier.VerifyServiceToken(p.raw)
			default:
				_, err = verifier.VerifyAccessToken(p.raw)
			}
			if p.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), p.expectedErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewVerifierServiceTokenPublicKeyErrors(t *testing.T) {
	t.Parallel()
	_, publicKeyPath := tokentesting.WriteKeyPair(t)
	emptyPath := filepath.Join(t.TempDir(), "empty.pem")
	require.NoError(t, os.WriteFile(emptyPath, []byte("not a pem\n"), 0600))
	patterns := []struct {
		desc string
		path string
	}{
		{desc: "err: missing file", path: filepath.Join(t.TempDir(), "missing.pem")},
		{desc: "err: no pem block", path: emptyPath},
		{desc: "err: invalid key", path: "testdata/invalid-public.pem"},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			_, err := NewVerifier(publicKeyPath, optionsTestIssuer, optionsTestAudience, WithServiceTokenPublicKey(p.path))
			assert.Error(t, err)
		})
	}
}

func TestServiceTokenVerifierOptions(t *testing.T) {
	t.Parallel()
	_, publicKeyPath := tokentesting.WriteKeyPair(t)
	patterns := []struct {
		desc            string
		publicKeyPath   string
		notBefore       map[string]string
		legacyNotBefore string
		expectedLen     int
		expectedErr     string
	}{
		{
			desc: "success: nothing set",
		},
		{
			desc:            "success: all set",
			publicKeyPath:   publicKeyPath,
			notBefore:       map[string]string{"batch": "2026-12-01T00:00:00Z"},
			legacyNotBefore: "2026-12-01T00:00:00Z",
			expectedLen:     3,
		},
		{
			desc:        "err: not before is not RFC 3339",
			notBefore:   map[string]string{"batch": "2026-12-01"},
			expectedErr: `service token not before for "batch"`,
		},
		{
			desc:            "err: legacy not before is not RFC 3339",
			legacyNotBefore: "yesterday",
			expectedErr:     "legacy service token not before",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			opts, err := ServiceTokenVerifierOptions(p.publicKeyPath, p.notBefore, p.legacyNotBefore)
			if p.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), p.expectedErr)
				return
			}
			require.NoError(t, err)
			assert.Len(t, opts, p.expectedLen)
		})
	}
}
