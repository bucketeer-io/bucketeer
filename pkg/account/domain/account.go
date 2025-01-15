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
	"errors"
	"regexp"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

var (
	maxAccountNameLength = 250
	// nolint:lll
	emailRegex                  = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	ErrSearchFilterNotFound     = errors.New("account: search filter not found")
	statusMissingOrganizationID = gstatus.New(
		codes.InvalidArgument,
		"account: organization id must be specified",
	)
	statusEmailIsEmpty            = gstatus.New(codes.InvalidArgument, "account: email is empty")
	statusInvalidEmail            = gstatus.New(codes.InvalidArgument, "account: invalid email format")
	statusFirstNameIsEmpty        = gstatus.New(codes.InvalidArgument, "account: first name is empty")
	statusInvalidFirstName        = gstatus.New(codes.InvalidArgument, "account: invalid first name format")
	statusLastNameIsEmpty         = gstatus.New(codes.InvalidArgument, "account: last name is empty")
	statusInvalidLastName         = gstatus.New(codes.InvalidArgument, "account: invalid last name format")
	statusLanguageIsEmpty         = gstatus.New(codes.InvalidArgument, "account: language is empty")
	statusInvalidOrganizationRole = gstatus.New(codes.InvalidArgument, "account: invalid organization role")
)

type AccountV2 struct {
	*proto.AccountV2
}

type AccountWithOrganization struct {
	*proto.AccountV2
	*environmentproto.Organization
}

func NewAccountV2(
	email, name, firstName, lastName, language, avatarImageURL, organizationID string,
	organizationRole proto.AccountV2_Role_Organization,
	environmentRoles []*proto.AccountV2_EnvironmentRole,
) *AccountV2 {
	now := time.Now().Unix()
	return &AccountV2{&proto.AccountV2{
		Email:            email,
		Name:             name,
		AvatarImageUrl:   avatarImageURL,
		Tags:             []string{}, // TODO: Implement tags
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
	if organizationRole != nil {
		updated.OrganizationRole = organizationRole.Role
	}
	if len(updated.EnvironmentRoles) > 0 {
		updated.EnvironmentRoles = environmentRoles
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
		return statusMissingOrganizationID.Err()
	}
	if a.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if !emailRegex.MatchString(a.Email) {
		return statusInvalidEmail.Err()
	}
	if a.FirstName == "" {
		return statusFirstNameIsEmpty.Err()
	}
	if len(a.FirstName) > maxAccountNameLength {
		return statusInvalidFirstName.Err()
	}
	if a.LastName == "" {
		return statusLastNameIsEmpty.Err()
	}
	if len(a.LastName) > maxAccountNameLength {
		return statusInvalidLastName.Err()
	}
	if a.Language == "" {
		return statusLanguageIsEmpty.Err()
	}
	if a.OrganizationRole == proto.AccountV2_Role_Organization_UNASSIGNED {
		return statusInvalidOrganizationRole.Err()
	}
	return nil
}
