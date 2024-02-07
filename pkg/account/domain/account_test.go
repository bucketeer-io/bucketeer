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

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

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

func TestChangeEnvironmentRole(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"avatarImageURL",
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	a.ChangeEnvironmentRole([]*proto.AccountV2_EnvironmentRole{
		{
			EnvironmentId: "environmentID",
			Role:          proto.AccountV2_Role_Environment_EDITOR,
		},
	})
	assert.Equal(t, []*proto.AccountV2_EnvironmentRole{
		{
			EnvironmentId: "environmentID",
			Role:          proto.AccountV2_Role_Environment_EDITOR,
		},
	}, a.EnvironmentRoles)
}

func TestPatchEnvironmentRole(t *testing.T) {
	patterns := []struct {
		desc     string
		envRoles []*proto.AccountV2_EnvironmentRole
		expected []*proto.AccountV2_EnvironmentRole
	}{
		{
			desc: "append a new role",
			envRoles: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID3",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
			},
			expected: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID",
					Role:          proto.AccountV2_Role_Environment_VIEWER,
				},
				{
					EnvironmentId: "environmentID2",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
				{

					EnvironmentId: "environmentID3",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
			},
		},
		{
			desc: "replace a role",
			envRoles: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
			},
			expected: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
				{
					EnvironmentId: "environmentID2",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
			},
		},
		{
			desc: "no change",
			envRoles: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID",
					Role:          proto.AccountV2_Role_Environment_VIEWER,
				},
			},
			expected: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID",
					Role:          proto.AccountV2_Role_Environment_VIEWER,
				},
				{
					EnvironmentId: "environmentID2",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
			},
		},
		{
			desc: "mix",
			envRoles: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
				{
					EnvironmentId: "environmentID3",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
			},
			expected: []*proto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "environmentID",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
				{
					EnvironmentId: "environmentID2",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
				{
					EnvironmentId: "environmentID3",
					Role:          proto.AccountV2_Role_Environment_EDITOR,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			a := NewAccountV2(
				"email",
				"name",
				"avatarImageURL",
				"organizationID",
				proto.AccountV2_Role_Organization_MEMBER,
				[]*proto.AccountV2_EnvironmentRole{
					{
						EnvironmentId: "environmentID",
						Role:          proto.AccountV2_Role_Environment_VIEWER,
					},
					{
						EnvironmentId: "environmentID2",
						Role:          proto.AccountV2_Role_Environment_EDITOR,
					},
				},
			)
			a.PatchEnvironmentRole(p.envRoles)
			assert.Equal(t, p.expected, a.EnvironmentRoles)
		})
	}
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
