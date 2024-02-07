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
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

type Project struct {
	*proto.Project
}

func NewProject(name, urlCode, description, creatorEmail, organizationID string, trial bool) (*Project, error) {
	now := time.Now().Unix()
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Project{&proto.Project{
		Id:             uid.String(),
		Name:           name,
		UrlCode:        urlCode,
		Description:    description,
		Disabled:       false,
		Trial:          trial,
		CreatorEmail:   creatorEmail,
		OrganizationId: organizationID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}}, nil
}

func (p *Project) ChangeDescription(description string) {
	p.Project.Description = description
	p.Project.UpdatedAt = time.Now().Unix()
}

func (p *Project) Rename(name string) {
	p.Project.Name = name
	p.Project.UpdatedAt = time.Now().Unix()
}

func (p *Project) Enable() {
	p.Project.Disabled = false
	p.Project.UpdatedAt = time.Now().Unix()
}

func (p *Project) Disable() {
	p.Project.Disabled = true
	p.Project.UpdatedAt = time.Now().Unix()
}

func (p *Project) ConvertTrial() {
	p.Project.Trial = false
	p.Project.UpdatedAt = time.Now().Unix()
}
