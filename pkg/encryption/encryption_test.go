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

package encryption

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

const (
	testKeyID       = "key-1"
	testServiceAcct = `{"type":"service_account","project_id":"p","private_key":"-----BEGIN PRIVATE KEY-----\nMIIB\n"}`
	testWebhookURL  = "https://hooks.slack.com/services/T000/B000/secret-token"
)

func testKey(t *testing.T, seed byte) []byte {
	t.Helper()
	key := make([]byte, aeadKeySize)
	for i := range key {
		key[i] = seed + byte(i)
	}
	return key
}

func writeKeyFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "key")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestAEADLocalRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	enc, err := NewAEADLocal(ctx, testKeyID, testKey(t, 1))
	require.NoError(t, err)

	patterns := []struct {
		desc      string
		plaintext string
	}{
		{desc: "fcm service account", plaintext: testServiceAcct},
		{desc: "slack webhook url", plaintext: testWebhookURL},
		{desc: "unicode", plaintext: "日本語 ✓"},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sealed, err := enc.Encrypt(ctx, p.plaintext)
			require.NoError(t, err)
			assert.True(t, IsEncrypted(sealed))
			assert.True(t, strings.HasPrefix(sealed, "enc:v1:"+testKeyID+":"))
			assert.NotContains(t, sealed, p.plaintext)
			assert.NotContains(t, sealed, "secret-token")

			opened, err := enc.Decrypt(ctx, sealed)
			require.NoError(t, err)
			assert.Equal(t, p.plaintext, opened)

			// A fresh nonce is used for every call.
			again, err := enc.Encrypt(ctx, p.plaintext)
			require.NoError(t, err)
			assert.NotEqual(t, sealed, again)
		})
	}
}

func TestAEADLocalEmptyPlaintextStaysEmpty(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	enc, err := NewAEADLocal(ctx, testKeyID, testKey(t, 1))
	require.NoError(t, err)
	sealed, err := enc.Encrypt(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, "", sealed)
	opened, err := enc.Decrypt(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, "", opened)
}

func TestAEADLocalDecryptLegacyPlaintext(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	enc, err := NewAEADLocal(ctx, testKeyID, testKey(t, 1))
	require.NoError(t, err)
	for _, legacy := range []string{testServiceAcct, testWebhookURL} {
		opened, err := enc.Decrypt(ctx, legacy)
		require.NoError(t, err)
		assert.Equal(t, legacy, opened)
	}
}

func TestAEADLocalDecryptErrors(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	enc, err := NewAEADLocal(ctx, testKeyID, testKey(t, 1))
	require.NoError(t, err)
	sealed, err := enc.Encrypt(ctx, testWebhookURL)
	require.NoError(t, err)

	otherKey, err := NewAEADLocal(ctx, testKeyID, testKey(t, 2))
	require.NoError(t, err)
	otherKeyID, err := NewAEADLocal(ctx, "key-2", testKey(t, 1))
	require.NoError(t, err)

	// Flip one bit of the ciphertext inside the envelope payload.
	keyID, payload, err := decodeEnvelope(sealed)
	require.NoError(t, err)
	blob := &wrapping.BlobInfo{}
	require.NoError(t, proto.Unmarshal(payload, blob))
	blob.Ciphertext[len(blob.Ciphertext)-1] ^= 0x01
	tamperedPayload, err := proto.Marshal(blob)
	require.NoError(t, err)
	tampered := encodeEnvelope(keyID, tamperedPayload)

	patterns := []struct {
		desc        string
		encryptor   Encryptor
		value       string
		expectedErr error
	}{
		{
			desc:        "wrong key",
			encryptor:   otherKey,
			value:       sealed,
			expectedErr: nil, // authentication failure from the cipher
		},
		{
			desc:        "unknown key id",
			encryptor:   otherKeyID,
			value:       sealed,
			expectedErr: ErrUnknownKeyID,
		},
		{
			desc:        "tampered ciphertext",
			encryptor:   enc,
			value:       tampered,
			expectedErr: nil,
		},
		{
			desc:        "unsupported envelope version",
			encryptor:   enc,
			value:       "enc:v2:" + testKeyID + ":AAAA",
			expectedErr: ErrMalformedValue,
		},
		{
			desc:        "missing payload",
			encryptor:   enc,
			value:       "enc:v1:" + testKeyID + ":",
			expectedErr: ErrMalformedValue,
		},
		{
			desc:        "missing key id",
			encryptor:   enc,
			value:       "enc:v1::AAAA",
			expectedErr: ErrMalformedValue,
		},
		{
			desc:        "payload is not base64url",
			encryptor:   enc,
			value:       "enc:v1:" + testKeyID + ":not*base64",
			expectedErr: ErrMalformedValue,
		},
		{
			desc:        "payload is not a blob",
			encryptor:   enc,
			value:       "enc:v1:" + testKeyID + ":" + base64.RawURLEncoding.EncodeToString([]byte{0xff, 0xff}),
			expectedErr: ErrMalformedValue,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			opened, err := p.encryptor.Decrypt(ctx, p.value)
			require.Error(t, err)
			assert.Equal(t, "", opened)
			if p.expectedErr != nil {
				assert.ErrorIs(t, err, p.expectedErr)
			}
		})
	}
}

func TestNoneEncryptor(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	enc := NewNone()

	sealed, err := enc.Encrypt(ctx, testServiceAcct)
	require.NoError(t, err)
	assert.Equal(t, testServiceAcct, sealed)
	assert.False(t, IsEncrypted(sealed))

	opened, err := enc.Decrypt(ctx, testWebhookURL)
	require.NoError(t, err)
	assert.Equal(t, testWebhookURL, opened)

	// A value encrypted while encryption was enabled is never returned as a credential.
	opened, err = enc.Decrypt(ctx, "enc:v1:"+testKeyID+":AAAA")
	assert.ErrorIs(t, err, ErrNotConfigured)
	assert.Equal(t, "", opened)
}

func TestIsEncrypted(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		value    string
		expected bool
	}{
		{desc: "envelope v1", value: "enc:v1:k:AAAA", expected: true},
		{desc: "future envelope version", value: "enc:v2:k:AAAA", expected: true},
		{desc: "json", value: testServiceAcct, expected: false},
		{desc: "url", value: testWebhookURL, expected: false},
		{desc: "empty", value: "", expected: false},
		{desc: "prefix only in the middle", value: "https://enc:v1", expected: false},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, IsEncrypted(p.value))
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	validKey := base64.StdEncoding.EncodeToString(testKey(t, 1))

	patterns := []struct {
		desc        string
		cfg         Config
		expectedErr error
		check       func(t *testing.T, enc Encryptor)
	}{
		{
			desc: "none",
			cfg:  Config{Type: TypeNone},
			check: func(t *testing.T, enc Encryptor) {
				sealed, err := enc.Encrypt(ctx, testWebhookURL)
				require.NoError(t, err)
				assert.Equal(t, testWebhookURL, sealed)
			},
		},
		{
			desc: "aead-local with a key file including a trailing newline",
			cfg: Config{
				Type:    TypeAEADLocal,
				KeyID:   testKeyID,
				KeyPath: writeKeyFile(t, validKey+"\n"),
			},
			check: func(t *testing.T, enc Encryptor) {
				sealed, err := enc.Encrypt(ctx, testWebhookURL)
				require.NoError(t, err)
				assert.True(t, IsEncrypted(sealed))
				// The same key loaded again decrypts the value.
				same, err := NewAEADLocal(ctx, testKeyID, testKey(t, 1))
				require.NoError(t, err)
				opened, err := same.Decrypt(ctx, sealed)
				require.NoError(t, err)
				assert.Equal(t, testWebhookURL, opened)
			},
		},
		{
			desc:        "unsupported type",
			cfg:         Config{Type: "vault"},
			expectedErr: ErrUnsupportedType,
		},
		{
			desc:        "empty type",
			cfg:         Config{Type: ""},
			expectedErr: ErrUnsupportedType,
		},
		{
			desc:        "aead-local without key path",
			cfg:         Config{Type: TypeAEADLocal, KeyID: testKeyID},
			expectedErr: errKeyPathRequired,
		},
		{
			desc:        "aead-local without key id",
			cfg:         Config{Type: TypeAEADLocal, KeyPath: writeKeyFile(t, validKey)},
			expectedErr: ErrInvalidKeyID,
		},
		{
			desc:        "aead-local with a key id containing the separator",
			cfg:         Config{Type: TypeAEADLocal, KeyID: "a:b", KeyPath: writeKeyFile(t, validKey)},
			expectedErr: ErrInvalidKeyID,
		},
		{
			desc:        "aead-local with a key of the wrong length",
			cfg:         Config{Type: TypeAEADLocal, KeyID: testKeyID, KeyPath: writeKeyFile(t, "c2hvcnQ=")},
			expectedErr: ErrInvalidKey,
		},
		{
			desc:        "aead-local with a key that is not base64",
			cfg:         Config{Type: TypeAEADLocal, KeyID: testKeyID, KeyPath: writeKeyFile(t, "not base64!")},
			expectedErr: ErrInvalidKey,
		},
		{
			desc: "aead-local with a missing key file",
			cfg: Config{
				Type:    TypeAEADLocal,
				KeyID:   testKeyID,
				KeyPath: filepath.Join(t.TempDir(), "missing"),
			},
			expectedErr: os.ErrNotExist,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			enc, err := New(ctx, p.cfg)
			if p.expectedErr != nil {
				assert.ErrorIs(t, err, p.expectedErr)
				assert.Nil(t, enc)
				return
			}
			require.NoError(t, err)
			p.check(t, enc)
		})
	}
}
