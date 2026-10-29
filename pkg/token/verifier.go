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
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	jose "github.com/go-jose/go-jose/v4"
)

// devOAuthPublicKeyFingerprint is the SHA-256 (hex) of the PKIX/DER encoding of
// the development OAuth public key that used to be committed to the public
// repository (tools/dev/cert/oauth-public.pem). The key pair is now generated
// locally and gitignored, but the old key remains in git history, so its private
// half is public forever. NewVerifier refuses this specific key for any non-local
// issuer so it can never be reused to forge system-admin tokens in a real
// deployment; local development (docker-compose, the dev container and the
// localenv cluster) uses the https://localhost issuer and is unaffected.
const devOAuthPublicKeyFingerprint = "fad5196defc62660595507a2f5496ec52237c9fb0bfb2122e13fda9f0e755acb"

type Verifier interface {
	VerifyAccessToken(string) (*AccessToken, error)
	VerifyRefreshToken(string) (*RefreshToken, error)
	VerifyDemoCreationToken(string) (*DemoCreationToken, error)
}

type verifier struct {
	issuer    string
	audience  string
	algorithm jose.SignatureAlgorithm
	pubKey    *rsa.PublicKey
}

func NewVerifier(keyPath, issuer, audience string) (Verifier, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	key, err := parseRSAPublicKey(data)
	if err != nil {
		return nil, err
	}
	if err := ensureNotInsecureDevKey(key, issuer); err != nil {
		return nil, err
	}
	return &verifier{
		issuer:    issuer,
		audience:  audience,
		algorithm: jose.RS256,
		pubKey:    key,
	}, nil
}

// ensureNotInsecureDevKey rejects the committed development OAuth public key for
// any non-local issuer. Local development uses the https://localhost issuer and
// is therefore unaffected; any other issuer paired with this key means the
// public development key leaked into a real deployment.
func ensureNotInsecureDevKey(pubKey *rsa.PublicKey, issuer string) error {
	fingerprint, err := publicKeyFingerprint(pubKey)
	if err != nil {
		return err
	}
	if fingerprint == devOAuthPublicKeyFingerprint && !isLocalIssuer(issuer) {
		return fmt.Errorf(
			"token: refusing to verify tokens with the leaked development OAuth key "+
				"(previously committed as tools/dev/cert/oauth-public.pem) for non-local issuer %q; "+
				"generate a dedicated OAuth key pair for this deployment",
			issuer,
		)
	}
	return nil
}

func publicKeyFingerprint(pubKey *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(der)
	return hex.EncodeToString(sum[:]), nil
}

// isLocalIssuer reports whether the issuer points at a local development host.
func isLocalIssuer(issuer string) bool {
	if issuer == "" {
		return false
	}
	host := issuer
	if u, err := url.Parse(issuer); err == nil && u.Host != "" {
		host = u.Hostname()
	} else if u, err := url.Parse("//" + issuer); err == nil && u.Host != "" {
		host = u.Hostname()
	}
	switch host {
	case "localhost", "127.0.0.1", "::1":
		return true
	}
	return strings.HasSuffix(host, ".localhost")
}

func (v *verifier) VerifyAccessToken(rawAccessToken string) (*AccessToken, error) {
	jws, err := jose.ParseSigned(rawAccessToken, []jose.SignatureAlgorithm{v.algorithm})
	if err != nil {
		return nil, fmt.Errorf("malformed jwt: %w", err)
	}
	payload, err := jws.Verify(v.pubKey)
	if err != nil {
		return nil, fmt.Errorf("invalid jwt: %w", err)
	}
	t := &AccessToken{}
	if err := json.Unmarshal(payload, t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}
	if t.Issuer != v.issuer {
		return nil, fmt.Errorf("id token issued by a different provider, expected %q got %q", v.issuer, t.Issuer)
	}
	if t.Audience != v.audience {
		return nil, fmt.Errorf("expected audience %q got %q", v.audience, t.Audience)
	}
	if t.Expiry.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired (Token Expiry: %v)", t.Expiry)
	}
	if t.Email == "" {
		return nil, fmt.Errorf("email must be not empty")
	}
	return t, nil
}

func (v *verifier) VerifyRefreshToken(rawRefreshToken string) (*RefreshToken, error) {
	jws, err := jose.ParseSigned(rawRefreshToken, []jose.SignatureAlgorithm{v.algorithm})
	if err != nil {
		return nil, fmt.Errorf("malformed jwt: %w", err)
	}
	payload, err := jws.Verify(v.pubKey)
	if err != nil {
		return nil, fmt.Errorf("invalid jwt: %w", err)
	}
	t := &RefreshToken{}
	if err := json.Unmarshal(payload, t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}
	if t.Expiry.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired (Token Expiry: %v)", t.Expiry)
	}
	if t.Email == "" {
		return nil, fmt.Errorf("email must be not empty")
	}
	return t, nil
}

func (v *verifier) VerifyDemoCreationToken(rawDemoToken string) (*DemoCreationToken, error) {
	jws, err := jose.ParseSigned(rawDemoToken, []jose.SignatureAlgorithm{v.algorithm})
	if err != nil {
		return nil, fmt.Errorf("malformed jwt: %w", err)
	}
	payload, err := jws.Verify(v.pubKey)
	if err != nil {
		return nil, fmt.Errorf("invalid jwt: %w", err)
	}
	t := &DemoCreationToken{}
	if err := json.Unmarshal(payload, t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}
	if t.Issuer != v.issuer {
		return nil, fmt.Errorf("demo token issued by a different provider, expected %q got %q", v.issuer, t.Issuer)
	}
	if t.Audience != v.audience {
		return nil, fmt.Errorf("expected audience %q got %q", v.audience, t.Audience)
	}
	if t.Expiry.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired (Token Expiry: %v)", t.Expiry)
	}
	if t.Email == "" {
		return nil, fmt.Errorf("email must be not empty")
	}
	return t, nil
}

func parseRSAPublicKey(data []byte) (*rsa.PublicKey, error) {
	input := data
	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}
	parsedKey, err := x509.ParsePKIXPublicKey(input)
	if err != nil {
		cert, err := x509.ParseCertificate(input)
		if err != nil {
			return nil, err
		}
		parsedKey = cert.PublicKey
	}
	pubKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not a valid RSA public key")
	}
	return pubKey, nil
}
