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
)

func NewOrganization(
	name, urlCode, ownerEmail, description string,
	trial, systemAdmin, passwordAuthenticationEnabled bool,
) (*Organization, error) {
	now := time.Now().Unix()
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Organization{&proto.Organization{
		Id:                            uid.String(),
		Name:                          name,
		UrlCode:                       urlCode,
		OwnerEmail:                    ownerEmail,
		Description:                   description,
		Disabled:                      false,
		Archived:                      false,
		Trial:                         trial,
		SystemAdmin:                   systemAdmin,
		PasswordAuthenticationEnabled: passwordAuthenticationEnabled,
		CreatedAt:                     now,
		UpdatedAt:                     now,
	}}, nil
}

func (p *Organization) Update(
	name *wrapperspb.StringValue,
	description *wrapperspb.StringValue,
	ownerEmail *wrapperspb.StringValue,
	passwordAuthenticationEnabled *wrapperspb.BoolValue,
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
	if passwordAuthenticationEnabled != nil {
		updated.PasswordAuthenticationEnabled = passwordAuthenticationEnabled.Value
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

func (p *Organization) ChangePasswordAuthentication(enabled bool) {
	p.Organization.PasswordAuthenticationEnabled = enabled
	p.Organization.UpdatedAt = time.Now().Unix()
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
