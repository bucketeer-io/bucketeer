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
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

type Account struct {
	*proto.Account
}

type AccountV2 struct {
	*proto.AccountV2
}

type AccountWithOrganization struct {
	*proto.AccountV2
	*environmentproto.Organization
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

func NewAccountV2(
	email, name, avatarImageURL, organizationID string,
	organizationRole proto.AccountV2_Role_Organization,
	environmentRoles []*proto.AccountV2_EnvironmentRole,
) *AccountV2 {
	now := time.Now().Unix()
	return &AccountV2{&proto.AccountV2{
		Email:            email,
		Name:             name,
		AvatarImageUrl:   avatarImageURL,
		OrganizationId:   organizationID,
		OrganizationRole: organizationRole,
		EnvironmentRoles: environmentRoles,
		Disabled:         false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}}
}

// TODO remove this function after migration
func ConvertAccountV2(account *Account, environmentID, organizationID string) *AccountV2 {
	envRole := convertEnvironmentRoleForV2(environmentID, account.Role)
	return &AccountV2{&proto.AccountV2{
		Email:            account.Email,
		Name:             account.Name,
		AvatarImageUrl:   "",
		OrganizationId:   organizationID,
		OrganizationRole: convertOrganizationRoleForV2(account.Role),
		EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{envRole},
		Disabled:         account.Disabled,
		CreatedAt:        account.CreatedAt,
		UpdatedAt:        account.UpdatedAt,
	}}
}

// TODO remove this function after migration
func convertEnvironmentRoleForV2(environmentID string, role proto.Account_Role) *proto.AccountV2_EnvironmentRole {
	v2Role := proto.AccountV2_Role_Environment_VIEWER
	switch role {
	case proto.Account_VIEWER:
		v2Role = proto.AccountV2_Role_Environment_VIEWER
	case proto.Account_EDITOR, proto.Account_OWNER:
		v2Role = proto.AccountV2_Role_Environment_EDITOR
	}
	return &proto.AccountV2_EnvironmentRole{
		EnvironmentId: environmentID,
		Role:          v2Role,
	}
}

// TODO remove this function after migration
func convertOrganizationRoleForV2(role proto.Account_Role) proto.AccountV2_Role_Organization {
	switch role {
	case proto.Account_VIEWER, proto.Account_EDITOR:
		return proto.AccountV2_Role_Organization_MEMBER
	case proto.Account_OWNER:
		return proto.AccountV2_Role_Organization_ADMIN
	}
	return proto.AccountV2_Role_Organization_MEMBER
}

// TODO remove this function after migration
func (a *AccountV2) PatchAccountV2EnvironmentRoles(environmentID string, newRole proto.Account_Role) {
	newEnvironmentRole := convertEnvironmentRoleForV2(environmentID, newRole)
	for _, e := range a.AccountV2.EnvironmentRoles {
		if e.EnvironmentId == newEnvironmentRole.EnvironmentId {
			e.Role = newEnvironmentRole.Role
			a.UpdatedAt = time.Now().Unix()
			return
		}
	}
	a.AccountV2.EnvironmentRoles = append(a.AccountV2.EnvironmentRoles, newEnvironmentRole)
	a.UpdatedAt = time.Now().Unix()
}

// TODO remove this function after migration
func (a *AccountV2) RemoveAccountV2EnvironmentRole(rmEnvironmentID string) {
	roles := make([]*proto.AccountV2_EnvironmentRole, 0, len(a.AccountV2.EnvironmentRoles))
	for _, r := range a.AccountV2.EnvironmentRoles {
		if r.EnvironmentId != rmEnvironmentID {
			roles = append(roles, r)
		}
	}
	a.AccountV2.EnvironmentRoles = roles
	a.UpdatedAt = time.Now().Unix()
}

func (a *AccountV2) ChangeName(newName string) error {
	a.AccountV2.Name = newName
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) ChangeAvatarImageURL(url string) error {
	a.AccountV2.AvatarImageUrl = url
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) ChangeOrganizationRole(role proto.AccountV2_Role_Organization) error {
	a.AccountV2.OrganizationRole = role
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) ChangeEnvironmentRole(roles []*proto.AccountV2_EnvironmentRole) error {
	a.AccountV2.EnvironmentRoles = roles
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) Enable() error {
	a.AccountV2.Disabled = false
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) Disable() error {
	a.AccountV2.Disabled = true
	a.UpdatedAt = time.Now().Unix()
	return nil
}
