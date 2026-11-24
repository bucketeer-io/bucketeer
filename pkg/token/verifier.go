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
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	jose "github.com/go-jose/go-jose/v4"
)

type Verifier interface {
	VerifyAccessToken(string) (*AccessToken, error)
	VerifyRefreshToken(string) (*RefreshToken, error)
	VerifyServiceToken(string) (*ServiceToken, error)
	VerifyDemoCreationToken(string) (*DemoCreationToken, error)
}

type verifier struct {
	issuer                      string
	audience                    string
	algorithm                   jose.SignatureAlgorithm
	pubKey                      *rsa.PublicKey
	serviceTokenKeys            []*rsa.PublicKey
	serviceTokenNotBefore       map[string]time.Time
	legacyServiceTokenNotBefore time.Time
}

type VerifierOption func(*verifier) error

// WithServiceTokenPublicKey verifies service tokens with the keys in the file
// (one or more PEM blocks) instead of the shared key.
func WithServiceTokenPublicKey(keyPath string) VerifierOption {
	return func(v *verifier) error {
		data, err := os.ReadFile(keyPath)
		if err != nil {
			return err
		}
		keys, err := parseRSAPublicKeys(data)
		if err != nil {
			return err
		}
		v.serviceTokenKeys = keys
		return nil
	}
}

// WithServiceTokenNotBefore rejects the service tokens of each named service issued before its time.
func WithServiceTokenNotBefore(notBefore map[string]time.Time) VerifierOption {
	return func(v *verifier) error {
		v.serviceTokenNotBefore = notBefore
		return nil
	}
}

// WithLegacyServiceTokenNotBefore rejects service tokens without a service name issued before the time.
func WithLegacyServiceTokenNotBefore(notBefore time.Time) VerifierOption {
	return func(v *verifier) error {
		v.legacyServiceTokenNotBefore = notBefore
		return nil
	}
}

func NewVerifier(keyPath, issuer, audience string, opts ...VerifierOption) (Verifier, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	key, err := parseRSAPublicKey(data)
	if err != nil {
		return nil, err
	}
	v := &verifier{
		issuer:    issuer,
		audience:  audience,
		algorithm: jose.RS256,
		pubKey:    key,
	}
	for _, opt := range opts {
		if err := opt(v); err != nil {
			return nil, err
		}
	}
	return v, nil
}

// ServiceTokenVerifierOptions builds the options from the command line values; each is optional.
func ServiceTokenVerifierOptions(
	publicKeyPath string,
	notBefore map[string]string,
	legacyNotBefore string,
) ([]VerifierOption, error) {
	var opts []VerifierOption
	if publicKeyPath != "" {
		opts = append(opts, WithServiceTokenPublicKey(publicKeyPath))
	}
	if len(notBefore) > 0 {
		parsed := make(map[string]time.Time, len(notBefore))
		for service, value := range notBefore {
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return nil, fmt.Errorf("service token not before for %q: %w", service, err)
			}
			parsed[service] = t
		}
		opts = append(opts, WithServiceTokenNotBefore(parsed))
	}
	if legacyNotBefore != "" {
		t, err := time.Parse(time.RFC3339, legacyNotBefore)
		if err != nil {
			return nil, fmt.Errorf("legacy service token not before: %w", err)
		}
		opts = append(opts, WithLegacyServiceTokenNotBefore(t))
	}
	return opts, nil
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
	// Empty token_use is a token issued before the claim existed, including
	// the service token. It still verifies here so the process can roll
	// forward or back without a new token. Refresh and demo verifiers do
	// not accept that token as their own kind.
	if t.TokenUse != "" && t.TokenUse != TokenUseAccess {
		return nil, fmt.Errorf("unexpected token_use %q, expected %q", t.TokenUse, TokenUseAccess)
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
	if t.IsSystemAdmin && t.IsServiceToken && t.IssuedAt.Before(v.legacyServiceTokenNotBefore) {
		return nil, fmt.Errorf("service token issued before %v", v.legacyServiceTokenNotBefore)
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
	switch t.TokenUse {
	case TokenUseRefresh:
		// A refresh token must not carry iss or aud. The previous access and
		// demo verifiers treat those claims as sufficient and ignore token_use.
		if t.Issuer != "" || t.Audience != "" {
			return nil, fmt.Errorf("refresh token must not include iss or aud")
		}
	case "":
		// Issued before token_use. Those refresh tokens also have no iss or aud.
		// An access token always has both, so it cannot take this path.
		if t.Issuer != "" || t.Audience != "" {
			return nil, fmt.Errorf("unexpected token_use %q, expected %q", t.TokenUse, TokenUseRefresh)
		}
	default:
		return nil, fmt.Errorf("unexpected token_use %q, expected %q", t.TokenUse, TokenUseRefresh)
	}
	if t.Expiry.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired (Token Expiry: %v)", t.Expiry)
	}
	if t.Email == "" {
		return nil, fmt.Errorf("email must be not empty")
	}
	return t, nil
}

func (v *verifier) VerifyServiceToken(rawServiceToken string) (*ServiceToken, error) {
	keys := v.serviceTokenKeys
	if len(keys) == 0 {
		keys = []*rsa.PublicKey{v.pubKey}
	}
	t, err := v.verifyClaims(rawServiceToken, TokenUseService, keys)
	if err != nil {
		return nil, err
	}
	if t.Subject == "" {
		return nil, fmt.Errorf("sub must be not empty")
	}
	if t.ID == "" {
		return nil, fmt.Errorf("jti must be not empty")
	}
	if notBefore, ok := v.serviceTokenNotBefore[t.Subject]; ok && t.IssuedAt.Before(notBefore) {
		return nil, fmt.Errorf("service token of %q issued before %v", t.Subject, notBefore)
	}
	return (*ServiceToken)(t), nil
}

func (v *verifier) VerifyDemoCreationToken(rawDemoToken string) (*DemoCreationToken, error) {
	t, err := v.verifyClaims(rawDemoToken, TokenUseDemo, []*rsa.PublicKey{v.pubKey})
	if err != nil {
		return nil, err
	}
	if t.Email == "" {
		return nil, fmt.Errorf("email must be not empty")
	}
	return (*DemoCreationToken)(t), nil
}

func (v *verifier) verifyClaims(raw, tokenUse string, keys []*rsa.PublicKey) (*Claims, error) {
	jws, err := jose.ParseSigned(raw, []jose.SignatureAlgorithm{v.algorithm})
	if err != nil {
		return nil, fmt.Errorf("malformed jwt: %w", err)
	}
	var payload []byte
	for _, key := range keys {
		if payload, err = jws.Verify(key); err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("invalid jwt: %w", err)
	}
	t := &Claims{}
	if err := json.Unmarshal(payload, t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}
	if err := checkTokenUse(t.TokenUse, tokenUse); err != nil {
		return nil, err
	}
	if t.Issuer != v.issuer {
		return nil, fmt.Errorf("%s token issued by a different provider, expected %q got %q", tokenUse, v.issuer, t.Issuer)
	}
	if t.Audience != v.audience {
		return nil, fmt.Errorf("expected audience %q got %q", v.audience, t.Audience)
	}
	if t.Expiry.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired (Token Expiry: %v)", t.Expiry)
	}
	return t, nil
}

func checkTokenUse(got, want string) error {
	if got != want {
		return fmt.Errorf("unexpected token_use %q, expected %q", got, want)
	}
	return nil
}

func parseRSAPublicKeys(data []byte) ([]*rsa.PublicKey, error) {
	var keys []*rsa.PublicKey
	for {
		block, rest := pem.Decode(data)
		if block == nil {
			break
		}
		key, err := parseRSAPublicKey(pem.EncodeToMemory(block))
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
		data = rest
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("no public key found")
	}
	return keys, nil
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
