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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// devOAuthPublicKeyPEM is the development OAuth public key that was previously
// committed as tools/dev/cert/oauth-public.pem. It is kept here as a fixture so
// the guard's rejection of that specific (now-leaked) key stays tested; the key
// pair itself is generated locally and gitignored.
const devOAuthPublicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAx7b6Kk1Rf0f8lBAJCkE1
+hpiHbrTJbabSfO1pTB0VLOr32nRCPjBAcEseS7JdUmvBIKfeRGzoE17q9kbTuBI
SiAPxeQwOs3uqsm522GZIHErvkWoQ/rUn6OF0AuQ6Co506nblaPfYvsqdKMpBre7
ZkXu3M0r+upxyJDPnx/fMz1xfwrLdwtoSXRl5B2e68MveFc6olfI3RG1VR2geLWf
0/lUaTD/i2EnvmjHTIpxE4cvYPhNLZLE95zIq9h7WDTwQWXtJdFvHtnMhpzZHVl7
9OPAVHJMBURXXcMHobQwMxghkniDydmz7/ZHYoBY+tZXWuZ4CwC7Vc9UOfYPX3Vu
5asQfqIo0fzkfQowkXHBViII3HQOckjErbU2mGKZjRlDDynEo0TsxJTryoyvKSi0
l7l8mJbR8al/XhCskuIWZxSfCVAX2ewLzqMsRNRm6q+j7xQc234v8TGC6fXYUWB3
pKbiRYzGc5b8iFJmgKw0P7rxAqUWcgAZR+1GpCuxRILS0Xm/zq5ES7ItLlOIR7el
RtG2J1uExQYzItQbswJlYSmGGVjNK+wtLw+JKU8McaIyHmk3j45s+qJyKXlb9yCH
tVNoZD8ZcBmCyQ6Mcj0sfrOIcMXe1DVQxN2KxEmXul5WWIJtbEffCMLcW3HJ+hOv
kW3mHcyUV1j7p8Qxu30xlP0CAwEAAQ==
-----END PUBLIC KEY-----
`

func writeTempKey(t *testing.T, pemData string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "key.pem")
	require.NoError(t, os.WriteFile(path, []byte(pemData), 0o600))
	return path
}

// TestNewVerifierRejectsDevKeyForNonLocalIssuer ensures the committed public
// development OAuth key cannot be used to verify tokens in a real deployment,
// while local development (https://localhost issuer) keeps working.
func TestNewVerifierRejectsDevKeyForNonLocalIssuer(t *testing.T) {
	t.Parallel()

	devKeyPath := writeTempKey(t, devOAuthPublicKeyPEM)

	// A freshly generated key is the "regenerated for production" control case.
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	der, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	require.NoError(t, err)
	otherKeyPath := writeTempKey(t, string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})))

	cases := []struct {
		desc    string
		keyPath string
		issuer  string
		wantErr bool
	}{
		{"dev key + production issuer is rejected", devKeyPath, "https://app.example.com", true},
		{"dev key + empty issuer is rejected", devKeyPath, "", true},
		{"dev key + https localhost is allowed", devKeyPath, "https://localhost", false},
		{"dev key + localhost with port is allowed", devKeyPath, "https://localhost:9000", false},
		{"dev key + bare localhost is allowed", devKeyPath, "localhost", false},
		{"non-dev key + production issuer is allowed", otherKeyPath, "https://app.example.com", false},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			t.Parallel()
			_, err := NewVerifier(c.keyPath, c.issuer, "bucketeer")
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestLeakedDevKeyFingerprintMatchesConstant ensures the fingerprint constant
// used by the guard matches the leaked development key embedded above, so the
// guard keeps rejecting that specific key.
func TestLeakedDevKeyFingerprintMatchesConstant(t *testing.T) {
	t.Parallel()

	key, err := parseRSAPublicKey([]byte(devOAuthPublicKeyPEM))
	require.NoError(t, err)
	fingerprint, err := publicKeyFingerprint(key)
	require.NoError(t, err)
	assert.Equal(t, devOAuthPublicKeyFingerprint, fingerprint,
		"devOAuthPublicKeyFingerprint is out of sync with the embedded leaked dev key")
}

func TestIsLocalIssuer(t *testing.T) {
	t.Parallel()
	for _, issuer := range []string{
		"https://localhost",
		"http://localhost",
		"https://localhost:9000",
		"localhost",
		"https://foo.localhost",
		"127.0.0.1",
	} {
		assert.True(t, isLocalIssuer(issuer), issuer)
	}
	for _, issuer := range []string{
		"",
		"https://app.example.com",
		"https://bucketeer.example.io",
		"https://notlocalhost.com",
	} {
		assert.False(t, isLocalIssuer(issuer), issuer)
	}
}
