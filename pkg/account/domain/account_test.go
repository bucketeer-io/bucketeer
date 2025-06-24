// Copyright 2025 The Bucketeer Authors.
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
	"google.golang.org/protobuf/types/known/wrapperspb"

	proto "github.com/bucketeer-io/bucketeer/proto/account"
	"github.com/bucketeer-io/bucketeer/proto/common"
)

func TestNewAccountV2(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"John",
		"Doe",
		"en",
		"avatarImageURL",
		[]string{"tag"},
		[]string{"team"},
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	assert.Equal(t, "email", a.Email)
	assert.Equal(t, "name", a.Name)
	assert.Equal(t, "John", a.FirstName)
	assert.Equal(t, "Doe", a.LastName)
	assert.Equal(t, "en", a.Language)
	assert.Equal(t, "avatarImageURL", a.AvatarImageUrl)
	assert.Equal(t, []string{"tag"}, a.Tags)
	assert.Equal(t, "organizationID", a.OrganizationId)
	assert.Equal(t, proto.AccountV2_Role_Organization_MEMBER, a.OrganizationRole)
	assert.Equal(t, []*proto.AccountV2_EnvironmentRole{}, a.EnvironmentRoles)
}

func TestChangeFirstName(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"fname",
		"lname",
		"en",
		"avatarImageURL",
		[]string{"tag"},
		[]string{"team"},
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	a.ChangeFirstName("newName")
	assert.Equal(t, "newName", a.FirstName)
}

func TestChangeLastName(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"fname",
		"lname",
		"en",
		"avatarImageURL",
		[]string{"tag"},
		[]string{"team"},
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	a.ChangeLastName("newLastName")
	assert.Equal(t, "newLastName", a.LastName)
}

func TestChangeAvatarImageURL(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"fname",
		"lname",
		"en",
		"avatarImageURL",
		[]string{"tag-1"},
		[]string{"team"},
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	a.ChangeAvatarImageURL("newURL")
	assert.Equal(t, "newURL", a.AvatarImageUrl)
}

func TestChangeTags(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"fname",
		"lname",
		"en",
		"avatarImageURL",
		[]string{"tag-1"},
		[]string{"team"},
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{},
	)
	a.ChangeFirstName("newName")
	assert.Equal(t, "newName", a.FirstName)
}

func TestChangeOrganizationRole(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"fname",
		"lname",
		"en",
		"avatarImageURL",
		[]string{"tag"},
		[]string{"team"},
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
		"fname",
		"lname",
		"en",
		"avatarImageURL",
		[]string{"tag"},
		[]string{"team"},
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
				"fname",
				"lname",
				"en",
				"avatarImageURL",
				[]string{"tag"},
				[]string{"team"},
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
				})
			a.PatchEnvironmentRole(p.envRoles)
			assert.Equal(t, p.expected, a.EnvironmentRoles)
		})
	}
}

func TestEnableV2(t *testing.T) {
	a := NewAccountV2(
		"email",
		"name",
		"fname",
		"lname",
		"en",
		"avatarImageURL",
		[]string{"tag"},
		[]string{"team"},
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
		"fname",
		"lname",
		"en",
		"avatarImageURL",
		[]string{"tag"},
		[]string{"team"},
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
		FirstName:        "John",
		LastName:         "Doe",
		AvatarImageUrl:   "avatarImageURL",
		Tags:             []string{"tag"},
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
					DefaultFilter:    true,
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
		{
			desc: "add same targetType and environmentID filters with default filter true",
			expectedFilters: []*proto.SearchFilter{
				{
					Name:             "name0",
					Query:            "query0",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID",
					DefaultFilter:    true,
				},
				{
					Name:             "name1",
					Query:            "query1",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID",
					DefaultFilter:    true,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			a := NewAccountV2(
				account.Email,
				account.Name,
				account.FirstName,
				account.LastName,
				account.Language,
				account.AvatarImageUrl,
				account.Tags,
				[]string{"team"},
				account.OrganizationId,
				account.OrganizationRole, account.EnvironmentRoles)
			for _, f := range p.expectedFilters {
				_, err := a.AddSearchFilter(f.Name, f.Query, f.FilterTargetType, f.EnvironmentId, f.DefaultFilter)
				assert.Nil(t, err)
			}
			// account has not changed.
			assert.Equal(t, account.Name, a.Name)
			assert.Equal(t, account.FirstName, a.FirstName)
			assert.Equal(t, account.LastName, a.LastName)
			assert.Equal(t, account.Language, a.Language)
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
			}

			// If target and EnvID are the same, only one DefaultFilter can exist
			for srcCnt, f := range a.SearchFilters {
				if f.DefaultFilter {
					for dctCnt, ff := range a.SearchFilters {
						if srcCnt != dctCnt && ff.DefaultFilter && ff.FilterTargetType == f.FilterTargetType && ff.EnvironmentId == f.EnvironmentId {
							assert.New(t).Fail("multiple default filters")
						}
					}
				}
			}
		})
	}
}

func TestChangeSearchFilterName(t *testing.T) {
	account := proto.AccountV2{
		Email:            "email",
		FirstName:        "John",
		LastName:         "Doe",
		Language:         "en",
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
		desc             string
		existingFilters  []*proto.SearchFilter
		updateFilterName string
		expectedFilters  []*proto.SearchFilter
		error            error
	}{
		{
			desc:             "don't have a filter",
			existingFilters:  nil,
			updateFilterName: "update-name",
			expectedFilters:  []*proto.SearchFilter{},
			error:            ErrSearchFilterNotFound,
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
			updateFilterName: "update-name",
			expectedFilters: []*proto.SearchFilter{
				{
					Name:             "update-name",
					Query:            "query0",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID0",
					DefaultFilter:    false,
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
			updateFilterName: "update-name",
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
			error: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			a := NewAccountV2(
				account.Email,
				account.Name,
				account.FirstName,
				account.LastName,
				account.Language,
				account.AvatarImageUrl,
				account.Tags,
				[]string{"team"},
				account.OrganizationId,
				account.OrganizationRole, account.EnvironmentRoles)
			for _, f := range p.existingFilters {
				_, err := a.AddSearchFilter(f.Name, f.Query, f.FilterTargetType, f.EnvironmentId, f.DefaultFilter)
				assert.Nil(t, err)
			}
			updateFilterId := "update-filter-id"
			if len(a.SearchFilters) > 0 {
				updateFilterId = a.SearchFilters[(len(a.SearchFilters) / 2)].Id
			}
			err := a.ChangeSearchFilterName(updateFilterId, p.updateFilterName)
			assert.Equal(t, err, p.error)

			// account has not changed.
			assert.Equal(t, account.Name, a.Name)
			assert.Equal(t, account.FirstName, a.FirstName)
			assert.Equal(t, account.LastName, a.LastName)
			assert.Equal(t, account.Language, a.Language)
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

func TestChangeSearchFilterQuery(t *testing.T) {
	account := proto.AccountV2{
		Email:            "email",
		FirstName:        "John",
		LastName:         "Doe",
		Language:         "en",
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
		desc              string
		existingFilters   []*proto.SearchFilter
		updateFilterQuery string
		expectedFilters   []*proto.SearchFilter
		error             error
	}{
		{
			desc:              "don't have a filter",
			existingFilters:   nil,
			updateFilterQuery: "update-query",
			expectedFilters:   []*proto.SearchFilter{},
			error:             ErrSearchFilterNotFound,
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
			updateFilterQuery: "update-query",
			expectedFilters: []*proto.SearchFilter{
				{
					Name:             "name0",
					Query:            "update-query",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID0",
					DefaultFilter:    false,
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
			updateFilterQuery: "update-query",
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
					Query:            "update-query",
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
			error: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			a := NewAccountV2(
				account.Email,
				account.Name,
				account.FirstName,
				account.LastName,
				account.Language,
				account.AvatarImageUrl,
				account.Tags,
				[]string{"team"},
				account.OrganizationId,
				account.OrganizationRole, account.EnvironmentRoles)
			for _, f := range p.existingFilters {
				_, err := a.AddSearchFilter(f.Name, f.Query, f.FilterTargetType, f.EnvironmentId, f.DefaultFilter)
				assert.Nil(t, err)
			}
			updateFilterId := "update-filter-id"
			if len(a.SearchFilters) > 0 {
				updateFilterId = a.SearchFilters[(len(a.SearchFilters) / 2)].Id
			}
			err := a.ChangeSearchFilterQuery(updateFilterId, p.updateFilterQuery)
			assert.Equal(t, err, p.error)

			// account has not changed.
			assert.Equal(t, account.Name, a.Name)
			assert.Equal(t, account.FirstName, a.FirstName)
			assert.Equal(t, account.LastName, a.LastName)
			assert.Equal(t, account.Language, a.Language)
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

func TestChangeDefaultSearchFilter(t *testing.T) {
	account := proto.AccountV2{
		Email:            "email",
		FirstName:        "John",
		LastName:         "Doe",
		Language:         "en",
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
		desc                string
		existingFilters     []*proto.SearchFilter
		updateDefaultFilter bool
		expectedFilters     []*proto.SearchFilter
		error               error
	}{
		{
			desc:                "don't have a filter",
			existingFilters:     nil,
			updateDefaultFilter: false,
			expectedFilters:     []*proto.SearchFilter{},
			error:               ErrSearchFilterNotFound,
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
			updateDefaultFilter: true,
			expectedFilters: []*proto.SearchFilter{
				{
					Name:             "name0",
					Query:            "query0",
					FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "environmentID0",
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
			updateDefaultFilter: true,
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
				account.FirstName,
				account.LastName,
				account.Language,
				account.AvatarImageUrl,
				account.Tags,
				[]string{"team"},
				account.OrganizationId,
				account.OrganizationRole, account.EnvironmentRoles)
			for _, f := range p.existingFilters {
				_, err := a.AddSearchFilter(f.Name, f.Query, f.FilterTargetType, f.EnvironmentId, f.DefaultFilter)
				assert.Nil(t, err)
			}
			updateFilterId := "update-filter-id"
			if len(a.SearchFilters) > 0 {
				updateFilterId = a.SearchFilters[(len(a.SearchFilters) / 2)].Id
			}
			updateFilter := &proto.SearchFilter{
				Id:            updateFilterId,
				DefaultFilter: p.updateDefaultFilter,
			}
			err := a.ChangeDefaultSearchFilter(
				updateFilter.Id,
				updateFilter.DefaultFilter)
			assert.Equal(t, err, p.error)

			// account has not changed.
			assert.Equal(t, account.Name, a.Name)
			assert.Equal(t, account.FirstName, a.FirstName)
			assert.Equal(t, account.LastName, a.LastName)
			assert.Equal(t, account.Language, a.Language)
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
		FirstName:        "John",
		LastName:         "Doe",
		Language:         "en",
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
			error:           ErrSearchFilterNotFound,
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
				account.FirstName,
				account.LastName,
				account.Language,
				account.AvatarImageUrl,
				account.Tags,
				[]string{"team"},
				account.OrganizationId,
				account.OrganizationRole, account.EnvironmentRoles)
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
			assert.Equal(t, account.FirstName, a.FirstName)
			assert.Equal(t, account.LastName, a.LastName)
			assert.Equal(t, account.Language, a.Language)
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

func TestAccountV2_Update(t *testing.T) {
	account := NewAccountV2(
		"bucketeer@gmail.com",
		"name",
		"firstName",
		"lastName",
		"en",
		"avatarImageURL",
		[]string{"tag"},
		[]string{"team"},
		"organizationID",
		proto.AccountV2_Role_Organization_MEMBER,
		[]*proto.AccountV2_EnvironmentRole{
			{
				EnvironmentId: "e2e",
			},
		},
	)
	updated, err := account.Update(
		wrapperspb.String("newName"),
		wrapperspb.String("newFirstName"),
		wrapperspb.String("newLastName"),
		wrapperspb.String("ja"),
		wrapperspb.String("newAvatarImageURL"),
		nil,
		&common.StringListValue{Values: []string{"tag-1"}},
		[]*proto.TeamChange{
			{
				ChangeType: proto.ChangeType_CREATE,
				Team:       "team-1",
			},
		},
		&proto.UpdateAccountV2Request_OrganizationRoleValue{
			Role: proto.AccountV2_Role_Organization_ADMIN,
		},
		[]*proto.AccountV2_EnvironmentRole{
			{
				EnvironmentId: "e2e",
			},
			{
				EnvironmentId: "default",
			},
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "newName", updated.Name)
	assert.Equal(t, "newFirstName", updated.FirstName)
	assert.Equal(t, "newLastName", updated.LastName)
	assert.Equal(t, "ja", updated.Language)
	assert.Equal(t, "newAvatarImageURL", updated.AvatarImageUrl)
	assert.Equal(t, []string{"tag-1"}, updated.Tags)
	assert.Equal(t, "organizationID", updated.OrganizationId)
	assert.Equal(t, proto.AccountV2_Role_Organization_ADMIN, updated.OrganizationRole)
	assert.Equal(t, []*proto.AccountV2_EnvironmentRole{
		{
			EnvironmentId: "e2e",
		},
		{
			EnvironmentId: "default",
		},
	}, updated.EnvironmentRoles)
}

func TestAccountV2_GetAccountFullName(t *testing.T) {
	patterns := []struct {
		desc         string
		account      *proto.AccountV2
		expectedName string
	}{
		{
			desc: "no first name",
			account: &proto.AccountV2{
				FirstName: "",
				LastName:  "newLastName",
			},
			expectedName: "newLastName",
		},
		{
			desc: "no last name",
			account: &proto.AccountV2{
				FirstName: "newFirstName",
			},
			expectedName: "newFirstName",
		},
		{
			desc: "both first name and last name",
			account: &proto.AccountV2{
				FirstName: "newFirstName",
				LastName:  "newLastName",
			},
			expectedName: "newFirstName newLastName",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			account := NewAccountV2(
				p.account.Email,
				p.account.Name,
				p.account.FirstName,
				p.account.LastName,
				p.account.Language,
				p.account.AvatarImageUrl,
				p.account.Tags,
				p.account.Teams,
				p.account.OrganizationId,
				p.account.OrganizationRole,
				p.account.EnvironmentRoles,
			)
			assert.Equal(t, p.expectedName, account.GetAccountFullName())
		})
	}
}
