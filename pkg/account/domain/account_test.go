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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestNewAccount(t *testing.T) {
	a, err := NewAccount("email", proto.Account_VIEWER)
	assert.NoError(t, err)
	assert.Equal(t, "email", a.Id)
	assert.Equal(t, "email", a.Email)
	assert.Equal(t, proto.Account_VIEWER, a.Role)
}

func TestChangeRole(t *testing.T) {
	a, err := NewAccount("email", proto.Account_VIEWER)
	assert.NoError(t, err)
	a.ChangeRole(proto.Account_EDITOR)
	assert.Equal(t, proto.Account_EDITOR, a.Role)
}

func TestEnable(t *testing.T) {
	a, err := NewAccount("email", proto.Account_VIEWER)
	assert.NoError(t, err)
	a.Disabled = true
	a.Enable()
	assert.Equal(t, false, a.Disabled)
}

func TestDisable(t *testing.T) {
	a, err := NewAccount("email", proto.Account_VIEWER)
	assert.NoError(t, err)
	a.Disable()
	assert.Equal(t, true, a.Disabled)
}

func TestNewAccountV2(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"avatarImageURL",
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	assert.Equal(t, "email", a.Email)
	assert.Equal(t, "name", a.Name)
	assert.Equal(t, "avatarImageURL", a.AvatarImageUrl)
	assert.Equal(t, "organizationID", a.OrganizationId)
	assert.Equal(t, proto.AccountV2_Role_Organization_MEMBER, a.OrganizationRole)
	assert.Equal(t, []*proto.AccountV2_EnvironmentRole{}, a.EnvironmentRoles)
}

func TestChangeName(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"avatarImageURL",
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	a.ChangeName("newName")
	assert.Equal(t, "newName", a.Name)
}

func TestChangeAvatarImageURL(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"avatarImageURL",
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	a.ChangeAvatarImageURL("newURL")
	assert.Equal(t, "newURL", a.AvatarImageUrl)
}

func TestChangeOrganizationRole(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"avatarImageURL",
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	a.ChangeOrganizationRole(proto.AccountV2_Role_Organization_ADMIN)
	assert.Equal(t, proto.AccountV2_Role_Organization_ADMIN, a.OrganizationRole)
}

func TestEnableV2(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"avatarImageURL",
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	a.Disabled = true
	a.Enable()
	assert.Equal(t, false, a.Disabled)
}

func TestDisableV2(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"avatarImageURL",
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	a.Disable()
	assert.Equal(t, true, a.Disabled)
}

func TestConvertAccountV2(t *testing.T) {
	now := time.Now().Unix()
	patterns := []struct {
		desc                     string
		fromRole                 proto.Account_Role
		expectedOrganizationRole proto.AccountV2_Role_Organization
		expectedEnvironmentRole  proto.AccountV2_Role_Environment
	}{
		{
			desc:                     "convert from viewer",
			fromRole:                 proto.Account_VIEWER,
			expectedOrganizationRole: proto.AccountV2_Role_Organization_MEMBER,
			expectedEnvironmentRole:  proto.AccountV2_Role_Environment_VIEWER,
		},
		{
			desc:                     "convert from editor",
			fromRole:                 proto.Account_EDITOR,
			expectedOrganizationRole: proto.AccountV2_Role_Organization_MEMBER,
			expectedEnvironmentRole:  proto.AccountV2_Role_Environment_EDITOR,
		},
		{
			desc:                     "convert from owner",
			fromRole:                 proto.Account_OWNER,
			expectedOrganizationRole: proto.AccountV2_Role_Organization_ADMIN,
			expectedEnvironmentRole:  proto.AccountV2_Role_Environment_EDITOR,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			from := &Account{&proto.Account{
				Id:        "email",
				Email:     "email",
				Role:      proto.Account_UNASSIGNED,
				CreatedAt: now,
				UpdatedAt: now,
			}}
			from.Role = p.fromRole
			expected := &AccountV2{
				&proto.AccountV2{
					Email:            "email",
					Name:             "",
					AvatarImageUrl:   "",
					OrganizationId:   "organizationID",
					OrganizationRole: proto.AccountV2_Role_Organization_UNASSIGNED,
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "environmentID",
							Role:          proto.AccountV2_Role_Environment_UNASSIGNED,
						},
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
			}
			expected.OrganizationRole = p.expectedOrganizationRole
			expected.EnvironmentRoles[0].Role = p.expectedEnvironmentRole
			actual := ConvertAccountV2(from, "environmentID", "organizationID")
			assert.Equal(t, expected, actual)
		})
	}
}

func TestPatchAccountV2EnvironmentRoles(t *testing.T) {
	now := time.Now().Unix()
	patterns := []struct {
		desc          string
		EnvironmentID string
		Role          proto.Account_Role
		expectedRoles []*proto.AccountV2_EnvironmentRole
	}{
		{
			desc:          "append a new role",
			EnvironmentID: "environmentID2",
			Role:          proto.Account_EDITOR,
			expectedRoles: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID1",
					Role:          proto.AccountV2_Role_Environment_VIEWER,
				},
				{
					EnvironmentId: "environmentID2",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
			},
		},
		{
			desc:          "replace a role",
			EnvironmentID: "environmentID1",
			Role:          proto.Account_EDITOR,
			expectedRoles: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID1",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
			},
		},
		{
			desc:          "no change",
			EnvironmentID: "environmentID1",
			Role:          proto.Account_VIEWER,
			expectedRoles: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID1",
					Role:          proto.AccountV2_Role_Environment_VIEWER,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := &AccountV2{
				&proto.AccountV2{
					Email:            "email",
					Name:             "",
					AvatarImageUrl:   "",
					OrganizationId:   "organizationID",
					OrganizationRole: proto.AccountV2_Role_Organization_MEMBER,
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "environmentID1",
							Role:          proto.AccountV2_Role_Environment_VIEWER,
						},
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
			}
			actual.PatchAccountV2EnvironmentRoles(p.EnvironmentID, p.Role)
			assert.Equal(t, p.expectedRoles, actual.EnvironmentRoles)
		})
	}
}

func TestRemoveAccountV2EnvironmentRole(t *testing.T) {
	now := time.Now().Unix()
	patterns := []struct {
		desc          string
		EnvironmentID string
		expectedRoles []*proto.AccountV2_EnvironmentRole
	}{
		{
			desc:          "remove a role",
			EnvironmentID: "environmentID1",
			expectedRoles: []*proto.AccountV2_EnvironmentRole{},
		},
		{
			desc:          "no change",
			EnvironmentID: "environmentID2",
			expectedRoles: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID1",
					Role:          proto.AccountV2_Role_Environment_VIEWER,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := &AccountV2{
				&proto.AccountV2{
					Email:            "email",
					Name:             "",
					AvatarImageUrl:   "",
					OrganizationId:   "organizationID",
					OrganizationRole: proto.AccountV2_Role_Organization_MEMBER,
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "environmentID1",
							Role:          proto.AccountV2_Role_Environment_VIEWER,
						},
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
			}
			actual.RemoveAccountV2EnvironmentRole(p.EnvironmentID)
			assert.Equal(t, p.expectedRoles, actual.EnvironmentRoles)
		})
	}
}
