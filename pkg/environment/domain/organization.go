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
	"errors"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

type Organization struct {
	*proto.Organization
}

var (
	ErrCannotDisableSystemAdmin = errors.New("environment: cannot disable system admin")
	ErrCannotArchiveSystemAdmin = errors.New("environment: cannot archive system admin")
)

func NewOrganization(name, urlCode, description string, trial, systemAdmin bool) (*Organization, error) {
	now := time.Now().Unix()
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Organization{&proto.Organization{
		Id:          uid.String(),
		Name:        name,
		UrlCode:     urlCode,
		Description: description,
		Disabled:    false,
		Archived:    false,
		Trial:       trial,
		SystemAdmin: systemAdmin,
		CreatedAt:   now,
		UpdatedAt:   now,
	}}, nil
}

func (p *Organization) ChangeDescription(description string) {
	p.Organization.Description = description
	p.Organization.UpdatedAt = time.Now().Unix()
}

func (p *Organization) ChangeName(name string) {
	p.Organization.Name = name
	p.Organization.UpdatedAt = time.Now().Unix()
}

func (p *Organization) Enable() {
	p.Organization.Disabled = false
	p.Organization.UpdatedAt = time.Now().Unix()
}

func (p *Organization) Disable() error {
	if p.Organization.SystemAdmin {
		return ErrCannotDisableSystemAdmin
	}
	p.Organization.Disabled = true
	p.Organization.UpdatedAt = time.Now().Unix()
	return nil
}

func (p *Organization) Archive() error {
	if p.Organization.SystemAdmin {
		return ErrCannotArchiveSystemAdmin
	}
	p.Organization.Archived = true
	p.Organization.UpdatedAt = time.Now().Unix()
	return nil
}

func (p *Organization) Unarchive() {
	p.Organization.Archived = false
	p.Organization.UpdatedAt = time.Now().Unix()
}

func (p *Organization) ConvertTrial() {
	p.Organization.Trial = false
	p.Organization.UpdatedAt = time.Now().Unix()
}
