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
	"time"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

type EnvironmentV2 struct {
	*proto.EnvironmentV2
}

const (
	defaultAutoArchiveUnusedDays    int32 = 90
	defaultAutoArchiveCheckCodeRefs       = true
)

func NewEnvironmentV2(
	name,
	urlCode,
	description,
	projectID,
	organizationID string,
	requireComment bool,
	logger *zap.Logger,
) (*EnvironmentV2, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		logger.Error("failed to generate uuid", zap.Error(err))
		return nil, err
	}
	now := time.Now().Unix()
	return &EnvironmentV2{&proto.EnvironmentV2{
		Id:                       uid.String(),
		Name:                     name,
		UrlCode:                  urlCode,
		Description:              description,
		ProjectId:                projectID,
		OrganizationId:           organizationID,
		Archived:                 false,
		RequireComment:           requireComment,
		CreatedAt:                now,
		UpdatedAt:                now,
		AutoArchiveEnabled:       false,
		AutoArchiveUnusedDays:    defaultAutoArchiveUnusedDays,
		AutoArchiveCheckCodeRefs: defaultAutoArchiveCheckCodeRefs,
	}}, nil
}

func (e *EnvironmentV2) Update(
	name *wrapperspb.StringValue,
	description *wrapperspb.StringValue,
	requireComment *wrapperspb.BoolValue,
	archived *wrapperspb.BoolValue,
	autoArchiveEnabled *wrapperspb.BoolValue,
	autoArchiveUnusedDays *wrapperspb.Int32Value,
	autoArchiveCheckCodeRefs *wrapperspb.BoolValue,
) (*EnvironmentV2, error) {
	updated := &EnvironmentV2{}
	if err := copier.Copy(updated, e); err != nil {
		return nil, err
	}

	if name != nil {
		updated.Name = name.Value
	}
	if description != nil {
		updated.Description = description.Value
	}
	if requireComment != nil {
		updated.RequireComment = requireComment.Value
	}
	if archived != nil {
		updated.Archived = archived.Value
	}
	if autoArchiveEnabled != nil {
		updated.AutoArchiveEnabled = autoArchiveEnabled.Value
	}
	if autoArchiveUnusedDays != nil {
		updated.AutoArchiveUnusedDays = autoArchiveUnusedDays.Value
	}
	if autoArchiveCheckCodeRefs != nil {
		updated.AutoArchiveCheckCodeRefs = autoArchiveCheckCodeRefs.Value
	}

	updated.UpdatedAt = time.Now().Unix()
	return updated, nil
}

func (e *EnvironmentV2) Rename(name string) {
	e.Name = name
	e.UpdatedAt = time.Now().Unix()
}

func (e *EnvironmentV2) ChangeDescription(description string) {
	e.Description = description
	e.UpdatedAt = time.Now().Unix()
}

func (e *EnvironmentV2) ChangeRequireComment(state bool) {
	e.RequireComment = state
	e.UpdatedAt = time.Now().Unix()
}

func (e *EnvironmentV2) SetArchived() {
	e.Archived = true
	e.UpdatedAt = time.Now().Unix()
}

func (e *EnvironmentV2) SetUnarchived() {
	e.Archived = false
	e.UpdatedAt = time.Now().Unix()
}
