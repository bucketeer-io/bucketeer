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
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

type Organization struct {
	*proto.Organization
}

var (
	ErrCannotDisableSystemAdmin = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.EnvironmentPackageName,
		"cannot disable system admin",
		"system_admin_organization")
	ErrCannotArchiveSystemAdmin = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.EnvironmentPackageName,
		"cannot archive system admin",
		"system_admin_organization")
	ErrCannotDisableGoogleAuthentication = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.EnvironmentPackageName,
		"cannot disable google authentication",
		"google_authentication_required")
)

func NewOrganization(
	name, urlCode, ownerEmail, description string,
	trial, systemAdmin bool,
	authenticationSettings *proto.AuthenticationSettings,
) (*Organization, error) {
	now := time.Now().Unix()
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// Set default authentication settings if not provided: Google is always enabled
	if authenticationSettings == nil {
		authenticationSettings = &proto.AuthenticationSettings{
			EnabledTypes: []proto.AuthenticationType{proto.AuthenticationType_AUTHENTICATION_TYPE_GOOGLE},
		}
	}

	return &Organization{&proto.Organization{
		Id:                     uid.String(),
		Name:                   name,
		UrlCode:                urlCode,
		OwnerEmail:             ownerEmail,
		Description:            description,
		Disabled:               false,
		Archived:               false,
		Trial:                  trial,
		SystemAdmin:            systemAdmin,
		AuthenticationSettings: authenticationSettings,
		CreatedAt:              now,
		UpdatedAt:              now,
	}}, nil
}

func (p *Organization) Update(
	name *wrapperspb.StringValue,
	description *wrapperspb.StringValue,
	ownerEmail *wrapperspb.StringValue,
	authenticationSettings *proto.AuthenticationSettings,
) (*Organization, error) {
	updated := &Organization{}
	if err := copier.Copy(updated, p); err != nil {
		return nil, err
	}
	if name != nil {
		updated.Name = name.Value
	}
	if description != nil {
		updated.Description = description.Value
	}
	if ownerEmail != nil {
		updated.OwnerEmail = ownerEmail.Value
	}
	if authenticationSettings != nil {
		updated.AuthenticationSettings = authenticationSettings
	}
	updated.UpdatedAt = time.Now().Unix()
	return updated, nil
}

func (p *Organization) ChangeDescription(description string) {
	p.Description = description
	p.UpdatedAt = time.Now().Unix()
}

func (p *Organization) ChangeOwnerEmail(ownerEmail string) {
	p.OwnerEmail = ownerEmail
	p.UpdatedAt = time.Now().Unix()
}

func (p *Organization) ChangeName(name string) {
	p.Name = name
	p.UpdatedAt = time.Now().Unix()
}

func (p *Organization) UpdateAuthenticationSettings(settings *proto.AuthenticationSettings) {
	p.AuthenticationSettings = settings
	p.UpdatedAt = time.Now().Unix()
}

func (p *Organization) EnableAuthenticationType(authType proto.AuthenticationType) {
	if p.AuthenticationSettings == nil {
		// Initialize with Google as the base authentication method
		p.AuthenticationSettings = &proto.AuthenticationSettings{
			EnabledTypes: []proto.AuthenticationType{proto.AuthenticationType_AUTHENTICATION_TYPE_GOOGLE},
		}
	}

	// Check if the authentication type is already enabled
	for _, t := range p.AuthenticationSettings.EnabledTypes {
		if t == authType {
			return // Already enabled, nothing to do
		}
	}

	// Add the authentication type
	p.AuthenticationSettings.EnabledTypes = append(
		p.AuthenticationSettings.EnabledTypes,
		authType,
	)
	p.UpdatedAt = time.Now().Unix()
}

func (p *Organization) DisableAuthenticationType(authType proto.AuthenticationType) error {
	// Google authentication cannot be disabled as it's the required base authentication method
	if authType == proto.AuthenticationType_AUTHENTICATION_TYPE_GOOGLE {
		return ErrCannotDisableGoogleAuthentication
	}

	if p.AuthenticationSettings == nil {
		return nil // Nothing to disable
	}

	// Remove the specified authentication type
	var newTypes []proto.AuthenticationType
	for _, t := range p.AuthenticationSettings.EnabledTypes {
		if t != authType {
			newTypes = append(newTypes, t)
		}
	}

	p.AuthenticationSettings.EnabledTypes = newTypes
	p.UpdatedAt = time.Now().Unix()
	return nil
}

func (p *Organization) IsPasswordAuthenticationEnabled() bool {
	if p.AuthenticationSettings == nil {
		return false
	}

	for _, authType := range p.AuthenticationSettings.EnabledTypes {
		if authType == proto.AuthenticationType_AUTHENTICATION_TYPE_PASSWORD {
			return true
		}
	}
	return false
}

func (p *Organization) Enable() {
	p.Disabled = false
	p.UpdatedAt = time.Now().Unix()
}

func (p *Organization) Disable() error {
	if p.SystemAdmin {
		return ErrCannotDisableSystemAdmin
	}
	p.Disabled = true
	p.UpdatedAt = time.Now().Unix()
	return nil
}

func (p *Organization) Archive() error {
	if p.SystemAdmin {
		return ErrCannotArchiveSystemAdmin
	}
	p.Archived = true
	p.UpdatedAt = time.Now().Unix()
	return nil
}

func (p *Organization) Unarchive() {
	p.Archived = false
	p.UpdatedAt = time.Now().Unix()
}

func (p *Organization) ConvertTrial() {
	p.Trial = false
	p.UpdatedAt = time.Now().Unix()
}
