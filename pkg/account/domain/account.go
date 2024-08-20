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
	"errors"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

var (
	errSearchFilterNotFound = errors.New("account: search filter not found")
)

type AccountV2 struct {
	*proto.AccountV2
}

type AccountWithOrganization struct {
	*proto.AccountV2
	*environmentproto.Organization
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
		SearchFilters:    nil,
	}}
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

func (a *AccountV2) PatchEnvironmentRole(patchRoles []*proto.AccountV2_EnvironmentRole) error {
	for _, p := range patchRoles {
		e := getEnvironmentRole(a.AccountV2.EnvironmentRoles, p.EnvironmentId)
		if e == nil {
			a.AccountV2.EnvironmentRoles = append(a.AccountV2.EnvironmentRoles, p)
			continue
		}
		e.Role = p.Role
	}
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func getEnvironmentRole(roles []*proto.AccountV2_EnvironmentRole, envID string) *proto.AccountV2_EnvironmentRole {
	for _, r := range roles {
		if r.EnvironmentId == envID {
			return r
		}
	}
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

func (a *AccountV2) AddSearchFilter(
	name string,
	query string,
	targetType proto.FilterTargetType,
	environmentID string, defaultFilter bool) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	searchFilter := &proto.SearchFilter{
		Id:               id.String(),
		Name:             name,
		Query:            query,
		FilterTargetType: targetType,
		EnvironmentId:    environmentID,
		DefaultFilter:    defaultFilter,
	}
	a.AccountV2.SearchFilters = append(a.AccountV2.SearchFilters, searchFilter)
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) DeleteSearchFilter(id string) error {
	for i, f := range a.AccountV2.SearchFilters {
		if f.Id == id {
			a.AccountV2.SearchFilters = append(a.AccountV2.SearchFilters[:i], a.AccountV2.SearchFilters[i+1:]...)
			if len(a.AccountV2.SearchFilters) == 0 {
				a.AccountV2.SearchFilters = nil
			}
			a.UpdatedAt = time.Now().Unix()
			return nil
		}
	}
	return errSearchFilterNotFound
}

func (a *AccountV2) UpdateSearchFilter(searchFilter *proto.SearchFilter) error {
	for i, f := range a.AccountV2.SearchFilters {
		if f.Id == searchFilter.Id {
			a.AccountV2.SearchFilters[i] = searchFilter
			a.UpdatedAt = time.Now().Unix()
			return nil
		}
	}
	return errSearchFilterNotFound
}
