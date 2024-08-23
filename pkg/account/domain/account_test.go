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
	"time"

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

func TestAddSearchFilter(t *testing.T) {
	account := proto.AccountV2{
		Email:            "email",
		Name:             "name",
		AvatarImageUrl:   "avatarImageURL",
		OrganizationId:   "organizationID",
		OrganizationRole: proto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
			{
				EnvironmentId: "environmentID",
				Role:          proto.AccountV2_Role_Environment_VIEWER,
			},
		},
		UpdatedAt: time.Now().Unix(),
	}

	patterns := []struct {
		desc            string
		expectedFilters []*proto.SearchFilter
	}{
		{
			desc: "add one filter",
			expectedFilters: []*proto.SearchFilter{
				{
					Name:             "name",
					Query:            "query",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID",
					DefaultFilter:    false,
				},
			},
		},
		{
			desc: "add some filters",
			expectedFilters: []*proto.SearchFilter{
				{
					Name:             "name0",
					Query:            "query0",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID0",
					DefaultFilter:    false,
				},
				{
					Name:             "name1",
					Query:            "query1",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID1",
					DefaultFilter:    false,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			a := NewAccountV2(
				account.Email,
				account.Name,
				account.AvatarImageUrl,
				account.OrganizationId,
				account.OrganizationRole,
				account.EnvironmentRoles,
			)
			for _, f := range p.expectedFilters {
				_, err := a.AddSearchFilter(f.Name, f.Query, f.FilterTargetType, f.EnvironmentId, f.DefaultFilter)
				assert.Nil(t, err)
			}
			// account has not changed.
			assert.Equal(t, account.Name, a.Name)
			assert.Equal(t, account.Email, a.Email)
			assert.Equal(t, account.AvatarImageUrl, a.AvatarImageUrl)
			assert.Equal(t, account.OrganizationId, a.OrganizationId)
			assert.Equal(t, account.OrganizationRole, a.OrganizationRole)
			assert.Equal(t, account.EnvironmentRoles, a.EnvironmentRoles)
			assert.Equal(t, account.UpdatedAt, a.UpdatedAt)

			assert.Equal(t, len(p.expectedFilters), len(a.SearchFilters))
			for i, f := range p.expectedFilters {
				filter := a.SearchFilters[i]
				assert.NotNil(t, filter.Id)
				assert.Equal(t, f.Name, filter.Name)
				assert.Equal(t, f.Query, filter.Query)
				assert.Equal(t, f.FilterTargetType, filter.FilterTargetType)
				assert.Equal(t, f.EnvironmentId, filter.EnvironmentId)
				assert.Equal(t, f.DefaultFilter, filter.DefaultFilter)
			}
		})
	}
}

func TestUpdateSearchFilter(t *testing.T) {
	account := proto.AccountV2{
		Email:            "email",
		Name:             "name",
		AvatarImageUrl:   "avatarImageURL",
		OrganizationId:   "organizationID",
		OrganizationRole: proto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
			{
				EnvironmentId: "environmentID",
				Role:          proto.AccountV2_Role_Environment_VIEWER,
			},
		},
		UpdatedAt: time.Now().Unix(),
	}

	patterns := []struct {
		desc            string
		existingFilters []*proto.SearchFilter
		updateFilter    *proto.SearchFilter
		expectedFilters []*proto.SearchFilter
		error           error
	}{
		{
			desc:            "don't have a filter",
			existingFilters: nil,
			updateFilter: &proto.SearchFilter{
				Name:             "update-name",
				Query:            "update-query",
				FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
				EnvironmentId:    "environmentID",
				DefaultFilter:    false,
			},
			expectedFilters: []*proto.SearchFilter{},
			error:           errSearchFilterNotFound,
		},
		{
			desc: "have a filter",
			existingFilters: []*proto.SearchFilter{
				{
					Name:             "name0",
					Query:            "query0",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID0",
					DefaultFilter:    false,
				},
			},
			updateFilter: &proto.SearchFilter{
				Name:             "update-name",
				Query:            "update-query",
				FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
				EnvironmentId:    "environmentID",
				DefaultFilter:    true,
			},
			expectedFilters: []*proto.SearchFilter{
				{
					Name:             "update-name",
					Query:            "update-query",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID",
					DefaultFilter:    true,
				},
			},
			error: nil,
		},
		{
			desc: "have some filters",
			existingFilters: []*proto.SearchFilter{
				{
					Name:             "name0",
					Query:            "query0",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID0",
					DefaultFilter:    false,
				},
				{
					Name:             "name1",
					Query:            "query1",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID1",
					DefaultFilter:    false,
				},
				{
					Name:             "name2",
					Query:            "query2",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID2",
					DefaultFilter:    false,
				},
			},
			updateFilter: &proto.SearchFilter{
				Name:             "update-name",
				Query:            "update-query",
				FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
				EnvironmentId:    "environmentID",
				DefaultFilter:    true,
			},
			expectedFilters: []*proto.SearchFilter{
				{
					Name:             "name0",
					Query:            "query0",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID0",
					DefaultFilter:    false,
				},
				{
					Name:             "update-name",
					Query:            "update-query",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID",
					DefaultFilter:    true,
				},
				{
					Name:             "name2",
					Query:            "query2",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID2",
					DefaultFilter:    false,
				},
			},
			error: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			a := NewAccountV2(
				account.Email,
				account.Name,
				account.AvatarImageUrl,
				account.OrganizationId,
				account.OrganizationRole,
				account.EnvironmentRoles,
			)
			for _, f := range p.existingFilters {
				_, err := a.AddSearchFilter(f.Name, f.Query, f.FilterTargetType, f.EnvironmentId, f.DefaultFilter)
				assert.Nil(t, err)
			}
			updateFilterId := "update-filter-id"
			if len(a.SearchFilters) > 0 {
				updateFilterId = a.SearchFilters[(len(a.SearchFilters) / 2)].Id
			}
			updateFilter := &proto.SearchFilter{
				Id:               updateFilterId,
				Name:             p.updateFilter.Name,
				Query:            p.updateFilter.Query,
				FilterTargetType: p.updateFilter.FilterTargetType,
				EnvironmentId:    p.updateFilter.EnvironmentId,
				DefaultFilter:    p.updateFilter.DefaultFilter,
			}
			err := a.UpdateSearchFilter(updateFilter)
			assert.Equal(t, err, p.error)

			// account has not changed.
			assert.Equal(t, account.Name, a.Name)
			assert.Equal(t, account.Email, a.Email)
			assert.Equal(t, account.AvatarImageUrl, a.AvatarImageUrl)
			assert.Equal(t, account.OrganizationId, a.OrganizationId)
			assert.Equal(t, account.OrganizationRole, a.OrganizationRole)
			assert.Equal(t, account.EnvironmentRoles, a.EnvironmentRoles)
			assert.Equal(t, account.UpdatedAt, a.UpdatedAt)

			assert.Equal(t, len(p.expectedFilters), len(a.SearchFilters))
			for i, f := range p.expectedFilters {
				assert.Equal(t, f.Name, a.SearchFilters[i].Name)
				assert.Equal(t, f.Query, a.SearchFilters[i].Query)
				assert.Equal(t, f.FilterTargetType, a.SearchFilters[i].FilterTargetType)
				assert.Equal(t, f.EnvironmentId, a.SearchFilters[i].EnvironmentId)
				assert.Equal(t, f.DefaultFilter, a.SearchFilters[i].DefaultFilter)
			}
		})
	}
}

func TestDeleteSearchFilter(t *testing.T) {
	account := proto.AccountV2{
		Email:            "email",
		Name:             "name",
		AvatarImageUrl:   "avatarImageURL",
		OrganizationId:   "organizationID",
		OrganizationRole: proto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
			{
				EnvironmentId: "environmentID",
				Role:          proto.AccountV2_Role_Environment_VIEWER,
			},
		},
		UpdatedAt: time.Now().Unix(),
	}

	patterns := []struct {
		desc            string
		existingFilters []*proto.SearchFilter
		expectedFilters []*proto.SearchFilter
		error           error
	}{
		{
			desc:            "don't have a filter",
			existingFilters: nil,
			expectedFilters: nil,
			error:           errSearchFilterNotFound,
		},
		{
			desc: "have a filter",
			existingFilters: []*proto.SearchFilter{
				{
					Name:             "name0",
					Query:            "query0",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID0",
					DefaultFilter:    false,
				},
			},
			expectedFilters: nil,
			error:           nil,
		},
		{
			desc: "have some filters",
			existingFilters: []*proto.SearchFilter{
				{
					Name:             "name0",
					Query:            "query0",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID0",
					DefaultFilter:    false,
				},
				{
					Name:             "name1",
					Query:            "query1",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID1",
					DefaultFilter:    false,
				},
				{
					Name:             "name2",
					Query:            "query2",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID2",
					DefaultFilter:    false,
				},
			},
			expectedFilters: []*proto.SearchFilter{
				{
					Name:             "name0",
					Query:            "query0",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID0",
					DefaultFilter:    false,
				},
				{
					Name:             "name2",
					Query:            "query2",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID2",
					DefaultFilter:    false,
				},
			},
			error: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			a := NewAccountV2(
				account.Email,
				account.Name,
				account.AvatarImageUrl,
				account.OrganizationId,
				account.OrganizationRole,
				account.EnvironmentRoles,
			)
			for _, f := range p.existingFilters {
				_, err := a.AddSearchFilter(f.Name, f.Query, f.FilterTargetType, f.EnvironmentId, f.DefaultFilter)
				assert.Nil(t, err)
			}
			deleteFilterId := "delete-filter-id"
			if len(a.SearchFilters) > 0 {
				deleteFilterId = a.SearchFilters[(len(a.SearchFilters) / 2)].Id
			}
			err := a.DeleteSearchFilter(deleteFilterId)
			assert.Equal(t, err, p.error)

			// account has not changed.
			assert.Equal(t, account.Name, a.Name)
			assert.Equal(t, account.Email, a.Email)
			assert.Equal(t, account.AvatarImageUrl, a.AvatarImageUrl)
			assert.Equal(t, account.OrganizationId, a.OrganizationId)
			assert.Equal(t, account.OrganizationRole, a.OrganizationRole)
			assert.Equal(t, account.EnvironmentRoles, a.EnvironmentRoles)
			assert.Equal(t, account.UpdatedAt, a.UpdatedAt)

			assert.Equal(t, len(p.expectedFilters), len(a.SearchFilters))
			if len(a.SearchFilters) > 0 {
				for i, expectedFilter := range p.expectedFilters {
					actualFilter := a.SearchFilters[i]
					assert.Equal(t, expectedFilter.Name, actualFilter.Name)
					assert.Equal(t, expectedFilter.Query, actualFilter.Query)
					assert.Equal(t, expectedFilter.FilterTargetType, actualFilter.FilterTargetType)
					assert.Equal(t, expectedFilter.EnvironmentId, actualFilter.EnvironmentId)
					assert.Equal(t, expectedFilter.DefaultFilter, actualFilter.DefaultFilter)
				}
			} else {
				assert.Equal(t, p.expectedFilters, a.SearchFilters)
			}
		})
	}
}
