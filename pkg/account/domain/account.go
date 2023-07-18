// Copyright 2023 The Bucketeer Authors.
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
	"time"

	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

type Account struct {
	*proto.Account
}

func NewAccount(email string, role proto.Account_Role) (*Account, error) {
	now := time.Now().Unix()
	return &Account{&proto.Account{
		Id:        email,
		Email:     email,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}}, nil
}

func (a *Account) ChangeRole(role proto.Account_Role) error {
	a.Account.Role = role
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *Account) Enable() error {
	a.Account.Disabled = false
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *Account) Disable() error {
	a.Account.Disabled = true
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *Account) Delete() error {
	a.Account.Deleted = true
	a.UpdatedAt = time.Now().Unix()
	return nil
}
