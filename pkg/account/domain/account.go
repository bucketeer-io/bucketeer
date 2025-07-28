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
	"regexp"
	"slices"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
	"github.com/bucketeer-io/bucketeer/proto/common"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

// nolint:lll
var (
	maxAccountNameLength = 250
	// nolint:lll
	emailRegex                 = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	ErrSearchFilterNotFound    = pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "search filter not found", nil, "search_filter")
	ErrTeamNotFound            = pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "team not found", nil, "team")
	ErrMissingOrganizationID   = pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "organization id must be specified", pkgErr.InvalidTypeEmpty, nil, "organization_id")
	ErrEmailIsEmpty            = pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "email is empty", pkgErr.InvalidTypeEmpty, nil, "email")
	ErrEmailInvalidFormat      = pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "invalid email format", pkgErr.InvalidTypeNotMatchFormat, nil, "email")
	ErrFullNameIsEmpty         = pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "full name is empty", pkgErr.InvalidTypeEmpty, nil, "full_name")
	ErrFirstNameInvalidFormat  = pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "invalid first name format", pkgErr.InvalidTypeNotMatchFormat, nil, "first name")
	ErrLastNameInvalidFormat   = pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "invalid last name format", pkgErr.InvalidTypeNotMatchFormat, nil, "last_name")
	ErrLanguageIsEmpty         = pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "language is empty", pkgErr.InvalidTypeEmpty, nil, "language")
	ErrOrganizationRoleInvalid = pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "invalid organization role", pkgErr.InvalidTypeEmpty, nil, "organization_role")
)

type AccountV2 struct {
	*proto.AccountV2
}

type AccountWithOrganization struct {
	*proto.AccountV2
	*environmentproto.Organization
}

func NewAccountV2(
	email, name, firstName, lastName, language, avatarImageURL string,
	tags []string,
	teams []string,
	organizationID string,
	organizationRole proto.AccountV2_Role_Organization,
	environmentRoles []*proto.AccountV2_EnvironmentRole,
) *AccountV2 {
	now := time.Now().Unix()
	if organizationRole == proto.AccountV2_Role_Organization_ADMIN {
		environmentRoles = []*proto.AccountV2_EnvironmentRole{}
	}
	return &AccountV2{&proto.AccountV2{
		Email:            email,
		Name:             name,
		AvatarImageUrl:   avatarImageURL,
		Tags:             tags,
		Teams:            teams,
		OrganizationId:   organizationID,
		OrganizationRole: organizationRole,
		EnvironmentRoles: environmentRoles,
		Disabled:         false,
		CreatedAt:        now,
		UpdatedAt:        now,
		SearchFilters:    nil,
		FirstName:        firstName,
		LastName:         lastName,
		Language:         language,
	}}
}

func (a *AccountV2) Update(
	name, firstName, lastName, language, avatarImageURL *wrapperspb.StringValue,
	avatar *proto.UpdateAccountV2Request_AccountV2Avatar,
	tags *common.StringListValue,
	teamChanges []*proto.TeamChange,
	organizationRole *proto.UpdateAccountV2Request_OrganizationRoleValue,
	environmentRoles []*proto.AccountV2_EnvironmentRole,
	isDisabled *wrapperspb.BoolValue,
) (*AccountV2, error) {
	updated := &AccountV2{}
	if err := copier.Copy(updated, a); err != nil {
		return nil, err
	}

	if name != nil {
		updated.Name = name.Value
	}
	if firstName != nil {
		updated.FirstName = firstName.Value
	}
	if lastName != nil {
		updated.LastName = lastName.Value
	}
	if language != nil {
		updated.Language = language.Value
	}
	if avatarImageURL != nil {
		updated.AvatarImageUrl = avatarImageURL.Value
	}
	if avatar != nil {
		updated.AvatarImage = avatar.AvatarImage
		updated.AvatarFileType = avatar.AvatarFileType
	}
	if tags != nil {
		updated.Tags = tags.Values
	}
	for _, teamChange := range teamChanges {
		switch teamChange.ChangeType {
		case proto.ChangeType_CREATE, proto.ChangeType_UPDATE:
			if err := updated.AddTeam(teamChange.Team); err != nil {
				return nil, err
			}
		case proto.ChangeType_DELETE:
			if err := updated.RemoveTeam(teamChange.Team); err != nil {
				return nil, err
			}
		}
	}
	if organizationRole != nil {
		updated.OrganizationRole = organizationRole.Role
	}
	if len(environmentRoles) > 0 {
		updated.EnvironmentRoles = environmentRoles
	}
	if updated.OrganizationRole == proto.AccountV2_Role_Organization_ADMIN {
		updated.EnvironmentRoles = []*proto.AccountV2_EnvironmentRole{}
	}
	if isDisabled != nil {
		updated.Disabled = isDisabled.Value
	}
	updated.UpdatedAt = time.Now().Unix()
	if err := validate(updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (a *AccountV2) AddTeam(team string) error {
	if slices.Contains(a.Teams, team) {
		// output info log
		return nil
	}
	a.Teams = append(a.Teams, team)
	return nil
}

func (a *AccountV2) RemoveTeam(team string) error {
	idx := slices.Index(a.Teams, team)
	if idx == -1 {
		return ErrTeamNotFound
	}
	a.Teams = slices.Delete(a.Teams, idx, idx+1)
	return nil
}

func (a *AccountV2) ChangeName(newName string) error {
	a.AccountV2.Name = newName
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) ChangeFirstName(newFirstName string) error {
	a.AccountV2.FirstName = newFirstName
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) ChangeLastName(newLastName string) error {
	a.AccountV2.LastName = newLastName
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) ChangeLanguage(newLanguage string) error {
	a.AccountV2.Language = newLanguage
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) ChangeAvatarImageURL(url string) error {
	a.AccountV2.AvatarImageUrl = url
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) ChangeAvatar(image []byte, fileType string) error {
	a.AccountV2.AvatarImage = image
	a.AccountV2.AvatarFileType = fileType
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) ChangeTags(tags []string) error {
	a.AccountV2.Tags = tags
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AccountV2) ChangeOrganizationRole(role proto.AccountV2_Role_Organization) error {
	a.AccountV2.OrganizationRole = role
	if role == proto.AccountV2_Role_Organization_ADMIN {
		a.AccountV2.EnvironmentRoles = []*proto.AccountV2_EnvironmentRole{}
	}
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

func (a *AccountV2) ChangeLastSeen(lastSeen int64) error {
	a.AccountV2.LastSeen = lastSeen
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
	environmentID string, defaultFilter bool) (*proto.SearchFilter, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	// Since there is only one default setting for a filter target, set the existing default to OFF.
	if defaultFilter {
		a.resetDefaultFilter(targetType, environmentID)
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
	return searchFilter, nil
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
	return ErrSearchFilterNotFound
}

func (a *AccountV2) ChangeSearchFilterName(id string, name string) error {
	for _, f := range a.AccountV2.SearchFilters {
		if f.Id == id {
			f.Name = name
			a.UpdatedAt = time.Now().Unix()
			return nil
		}
	}
	return ErrSearchFilterNotFound
}

func (a *AccountV2) ChangeSearchFilterQuery(id string, query string) error {
	for _, f := range a.AccountV2.SearchFilters {
		if f.Id == id {
			f.Query = query
			a.UpdatedAt = time.Now().Unix()
			return nil
		}
	}
	return ErrSearchFilterNotFound
}

func (a *AccountV2) ChangeDefaultSearchFilter(id string, defaultFilter bool) error {
	for _, f := range a.AccountV2.SearchFilters {
		if f.Id == id {
			// Since there is only one default setting for a filter target, set the existing default to OFF.
			if defaultFilter {
				a.resetDefaultFilter(f.FilterTargetType, f.EnvironmentId)
			}

			f.DefaultFilter = defaultFilter
			a.UpdatedAt = time.Now().Unix()
			return nil
		}
	}
	return ErrSearchFilterNotFound
}

func (a *AccountV2) resetDefaultFilter(targetFilter proto.FilterTargetType, environmentID string) {
	for _, f := range a.AccountV2.SearchFilters {
		if f.DefaultFilter &&
			targetFilter == f.FilterTargetType &&
			environmentID == f.EnvironmentId {
			f.DefaultFilter = false
		}
	}
}

func (a *AccountV2) GetAccountFullName() string {
	if a.FirstName == "" && a.LastName == "" {
		return a.Name
	}
	if a.FirstName == "" {
		return a.LastName
	}
	if a.LastName == "" {
		return a.FirstName
	}
	return a.FirstName + " " + a.LastName
}

func validate(a *AccountV2) error {
	if a.OrganizationId == "" {
		return ErrMissingOrganizationID
	}
	if a.Email == "" {
		return ErrEmailIsEmpty
	}
	if !emailRegex.MatchString(a.Email) {
		return ErrEmailInvalidFormat
	}
	// If both first name and last name are empty, the name field must not be empty
	if a.FirstName == "" && a.LastName == "" {
		// TODO: This should be removed after the new console is released and the migration is completed
		if a.Name == "" {
			return ErrFullNameIsEmpty
		}
	}
	// Validate first name length if it's provided
	if a.FirstName != "" && len(a.FirstName) > maxAccountNameLength {
		return ErrFirstNameInvalidFormat
	}
	// Validate last name length if it's provided
	if a.LastName != "" && len(a.LastName) > maxAccountNameLength {
		return ErrLastNameInvalidFormat
	}
	if a.Language == "" {
		return ErrLanguageIsEmpty
	}
	if a.OrganizationRole == proto.AccountV2_Role_Organization_UNASSIGNED {
		return ErrOrganizationRoleInvalid
	}
	return nil
}
