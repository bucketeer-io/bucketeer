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

package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	wrapperspb "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jinzhu/copier"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

const keyBytes = 32

var (
	ErrLastUsedAtNotUpdated = pkgErr.NewErrorFailedPrecondition(
		pkgErr.AccountPackageName, "last used at not updated")
	ErrInvalidLastUsedAt = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.AccountPackageName, "invalid last used at", "last_used_at")
)

type APIKey struct {
	*proto.APIKey
}

type EnvironmentAPIKey struct {
	*proto.EnvironmentAPIKey
}

func NewAPIKey(
	name string,
	role proto.APIKey_Role,
	maintainer string,
	description string,
) (*APIKey, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	key, err := generateKey()
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()

	return &APIKey{&proto.APIKey{
		Id:          id.String(),
		Name:        name,
		Role:        role,
		CreatedAt:   now,
		UpdatedAt:   now,
		Maintainer:  maintainer,
		ApiKey:      key,
		Description: description,
	}}, nil
}

func (a *APIKey) Rename(name string) error {
	a.Name = name
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *APIKey) Enable() error {
	a.Disabled = false
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *APIKey) Disable() error {
	a.Disabled = true
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func generateKey() (string, error) {
	b := make([]byte, keyBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (a *APIKey) Update(
	name *wrapperspb.StringValue,
	description *wrapperspb.StringValue,
	role proto.APIKey_Role,
	maintainer *wrapperspb.StringValue,
	disabled *wrapperspb.BoolValue,
) (*APIKey, error) {
	updated := &APIKey{}
	if err := copier.Copy(updated, a); err != nil {
		return nil, err
	}
	if name != nil {
		updated.Name = name.Value
	}
	if description != nil {
		updated.Description = description.Value
	}
	if role != proto.APIKey_UNKNOWN {
		updated.Role = role
	}
	if maintainer != nil {
		updated.Maintainer = maintainer.Value
	}
	if disabled != nil {
		updated.Disabled = disabled.Value
	}
	updated.UpdatedAt = time.Now().Unix()
	return updated, nil
}
