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

// Package testing provides the signing keys that tests of the token package
// and of its callers need.
package testing

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

// keyBits is the RSA key size. 2048 bits is enough for a key that only lives
// for one test, and it is much faster to generate than a 4096-bit key.
const keyBits = 2048

// WriteKeyPair generates an RSA key pair, writes both keys as PEM files under
// t.TempDir(), and returns their paths. The private key is for token.NewSigner
// and the public key is for token.NewVerifier, so a token signed with the
// first path verifies with the second.
//
// Every call returns a new pair. A test that signs and verifies must take
// both paths from the same call. Do not commit a key to the repository
// instead, not even one that is only used by tests.
func WriteKeyPair(t testing.TB) (privateKeyPath, publicKeyPath string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		t.Fatalf("failed to generate the RSA key: %v", err)
	}
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal the public key: %v", err)
	}
	dir := t.TempDir()
	privateKeyPath = filepath.Join(dir, "private.pem")
	publicKeyPath = filepath.Join(dir, "public.pem")
	writePEM(t, privateKeyPath, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key), 0o600)
	writePEM(t, publicKeyPath, "PUBLIC KEY", publicKeyDER, 0o644)
	return privateKeyPath, publicKeyPath
}

func writePEM(t testing.TB, path, blockType string, der []byte, perm os.FileMode) {
	t.Helper()
	data := pem.EncodeToMemory(&pem.Block{Type: blockType, Bytes: der})
	if err := os.WriteFile(path, data, perm); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}
