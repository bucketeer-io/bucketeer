// Copyright 2024 The Bucketeer Authors.
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
	"io/ioutil"
	"time"

	jose "gopkg.in/square/go-jose.v2"
)

type Verifier interface {
	Verify(string) (*IDToken, error)
}

type verifier struct {
	issuer    string
	clientID  string
	algorithm jose.SignatureAlgorithm
	pubKey    *rsa.PublicKey
}

func NewVerifier(keyPath, issuer, clientID string) (Verifier, error) {
	data, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	key, err := parseRSAPublicKey(data)
	if err != nil {
		return nil, err
	}
	return &verifier{
		issuer:    issuer,
		clientID:  clientID,
		algorithm: jose.RS256,
		pubKey:    key,
	}, nil
}

func (v *verifier) Verify(rawIDToken string) (*IDToken, error) {
	jws, err := jose.ParseSigned(rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("malformed jwt: %v", err)
	}
	payload, err := jws.Verify(v.pubKey)
	if err != nil {
		return nil, fmt.Errorf("invalid jwt: %v", err)
	}
	t := &IDToken{}
	if err := json.Unmarshal(payload, t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %v", err)
	}
	if t.Issuer != v.issuer {
		return nil, fmt.Errorf("id token issued by a different provider, expected %q got %q", v.issuer, t.Issuer)
	}
	if t.Audience != v.clientID {
		return nil, fmt.Errorf("expected audience %q got %q", v.clientID, t.Audience)
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
