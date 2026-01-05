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
	"encoding/json"
	"time"

	pb "github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	"github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

type Option func(*domain.Options)

func WithComment(c string) Option {
	return func(opts *domain.Options) {
		opts.Comment = c
	}
}

func WithNewVersion(ver int32) Option {
	return func(opts *domain.Options) {
		opts.NewVersion = ver
	}
}

func NewEvent(
	editor *domain.Editor,
	entityType domain.Event_EntityType,
	entityID string,
	eventType domain.Event_Type,
	event pb.Message,
	environmentID string,
	entityData, previousEntityData interface{},
	opts ...Option,
) (*domain.Event, error) {
	return newEvent(
		editor, entityType, entityID, eventType, event,
		environmentID, false, entityData, previousEntityData, opts...)
}

func NewAdminEvent(
	editor *domain.Editor,
	entityType domain.Event_EntityType,
	entityID string,
	eventType domain.Event_Type,
	event pb.Message,
	entityData, previousEntityData interface{},
	opts ...Option,
) (*domain.Event, error) {
	return newEvent(
		editor, entityType, entityID, eventType, event,
		storage.AdminEnvironmentID, true, entityData, previousEntityData, opts...)
}

func newEvent(
	editor *domain.Editor,
	entityType domain.Event_EntityType,
	entityID string,
	eventType domain.Event_Type,
	event pb.Message,
	environmentID string,
	isAdminEvent bool,
	entity, previousEntity interface{},
	opts ...Option,
) (*domain.Event, error) {
	options := domain.Options{
		Comment:    "",
		NewVersion: 1,
	}
	for _, opt := range opts {
		opt(&options)
	}
	buf, err := ptypes.MarshalAny(event)
	if err != nil {
		return nil, err
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	var entityData, prevEntityData []byte
	if entity != nil {
		entityData, err = json.MarshalIndent(entity, "", "  ")
		if err != nil {
			return nil, err
		}
	}
	if previousEntity != nil {
		prevEntityData, err = json.MarshalIndent(previousEntity, "", "  ")
		if err != nil {
			return nil, err
		}
	}
	return &domain.Event{
		Id:                 id.String(),
		Timestamp:          time.Now().Unix(),
		EntityType:         entityType,
		EntityId:           entityID,
		Type:               eventType,
		Editor:             editor,
		Data:               buf,
		EnvironmentId:      environmentID,
		IsAdminEvent:       isAdminEvent,
		EntityData:         string(entityData),
		PreviousEntityData: string(prevEntityData),
		Options:            &options,
	}, nil
}
