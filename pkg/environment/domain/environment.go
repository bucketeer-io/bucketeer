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
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

type EnvironmentV2 struct {
	*proto.EnvironmentV2
}

func NewEnvironmentV2(name, urlCode, description, projectID string, logger *zap.Logger) (*EnvironmentV2, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		logger.Error("failed to generate uuid", zap.Error(err))
		return nil, err
	}
	now := time.Now().Unix()
	return &EnvironmentV2{&proto.EnvironmentV2{
		Id:          uid.String(),
		Name:        name,
		UrlCode:     urlCode,
		Description: description,
		ProjectId:   projectID,
		Archived:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}}, nil
}

func (e *EnvironmentV2) Rename(name string) {
	e.EnvironmentV2.Name = name
	e.EnvironmentV2.UpdatedAt = time.Now().Unix()
}

func (e *EnvironmentV2) ChangeDescription(description string) {
	e.EnvironmentV2.Description = description
	e.EnvironmentV2.UpdatedAt = time.Now().Unix()
}

func (e *EnvironmentV2) SetArchived() {
	e.EnvironmentV2.Archived = true
	e.EnvironmentV2.UpdatedAt = time.Now().Unix()
}

func (e *EnvironmentV2) SetUnarchived() {
	e.EnvironmentV2.Archived = false
	e.EnvironmentV2.UpdatedAt = time.Now().Unix()
}
