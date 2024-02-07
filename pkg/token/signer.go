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
	"encoding/pem"
	"fmt"
	"io/ioutil"

	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Signer interface {
	Sign(*IDToken) (string, error)
}

type signer struct {
	sig jose.Signer
}

func NewSigner(keyPath string) (Signer, error) {
	data, err := ioutil.ReadFile(keyPath)
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

func (s *signer) Sign(token *IDToken) (string, error) {
	return jwt.Signed(s.sig).Claims(token).CompactSerialize()
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
