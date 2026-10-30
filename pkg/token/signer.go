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
	"encoding/pem"
	"fmt"
	"os"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type Signer interface {
	SignAccessToken(*AccessToken) (string, error)
	SignRefreshToken(*RefreshToken) (string, error)
	SignDemoCreationToken(*DemoCreationToken) (string, error)
}

type signer struct {
	sig jose.Signer
}

func NewSigner(keyPath string) (Signer, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	key, err := parseRSAPrivateKey(data)
	if err != nil {
		return nil, err
	}
	// TODO: Currently, we are using RSA algorithm to be compatible with istio envoy.
	// https://github.com/istio/proxy/tree/master/src/envoy/auth
	// But in the future, we should consider to move to HMAC for a better performance.
	return NewSignerWithPrivateKey(key)
}

func NewSignerWithPrivateKey(privateKey *rsa.PrivateKey) (Signer, error) {
	signingKey := jose.SigningKey{
		Key:       privateKey,
		Algorithm: jose.RS256,
	}
	sig, err := jose.NewSigner(signingKey, &jose.SignerOptions{})
	if err != nil {
		return nil, err
	}
	return &signer{sig: sig}, nil
}

func (s *signer) SignAccessToken(token *AccessToken) (string, error) {
	// token_use is stamped here, not taken from the caller, so an access token
	// cannot be signed as a refresh or demo-creation token.
	claims := *token
	claims.TokenUse = TokenUseAccess
	return jwt.Signed(s.sig).Claims(claims).Serialize()
}

func (s *signer) SignRefreshToken(token *RefreshToken) (string, error) {
	claims := *token
	claims.TokenUse = TokenUseRefresh
	// iss and aud are not written. The previous access and demo verifiers
	// accept a token that carries them and do not read token_use, so a
	// refresh token with those claims would verify as an access or demo
	// token on a pod that has not been updated, for the refresh lifetime.
	claims.Issuer = ""
	claims.Audience = ""
	return jwt.Signed(s.sig).Claims(claims).Serialize()
}

func (s *signer) SignDemoCreationToken(token *DemoCreationToken) (string, error) {
	claims := *token
	claims.TokenUse = TokenUseDemo
	return jwt.Signed(s.sig).Claims(claims).Serialize()
}

func parseRSAPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	input := data
	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}
	var parsedKey interface{}
	parsedKey, err := x509.ParsePKCS1PrivateKey(input)
	if err != nil {
		parsedKey, err = x509.ParsePKCS8PrivateKey(input)
		if err != nil {
			return nil, err
		}
	}
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not a valid RSA private key")
	}
	return rsaKey, nil
}
