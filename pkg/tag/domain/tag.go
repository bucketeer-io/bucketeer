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

	proto "github.com/bucketeer-io/bucketeer/proto/tag"
)

type Tag struct {
	*proto.Tag
}

func NewTag(id, environmentID string, entityType proto.Tag_EntityType) *Tag {
	now := time.Now().Unix()
	return &Tag{
		Tag: &proto.Tag{
			Id:            id,
			CreatedAt:     now,
			UpdatedAt:     now,
			EntityType:    entityType,
			EnvironmentId: environmentID,
		},
	}
}
