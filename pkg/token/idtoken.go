// Copyright 2022 The Bucketeer Authors.
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
	"encoding/base64"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck

	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
)

type IDToken struct {
	Issuer   string    `json:"iss"`
	Subject  string    `json:"sub"`
	Audience string    `json:"aud"`
	Expiry   time.Time `json:"exp"`
	IssuedAt time.Time `json:"iat"`
	Email    string    `json:"email"`
	// Use "role" as json tag to keep compatibility.
	// AdminRole is accountproto.Account_OWNER in the case of AdminAccount.
	// AdminRole is accountproto.Account_UNASSIGNED in the case of Account.
	AdminRole accountproto.Account_Role `json:"role"`
}

func (t *IDToken) IsAdmin() bool {
	return t.AdminRole != accountproto.Account_UNASSIGNED
}

func ExtractUserID(subject string) (string, error) {
	tokenSubject := &authproto.IDTokenSubject{}
	// Q: Why do we need to decode the sub string
	// A: https://github.com/coreos/dex/blob/master/server/internal/codec.go#L20
	data, err := base64.RawURLEncoding.DecodeString(subject)
	if err != nil {
		return "", err
	}
	err = proto.Unmarshal(data, tokenSubject)
	if err != nil {
		return "", err
	}
	return tokenSubject.UserId, nil
}
