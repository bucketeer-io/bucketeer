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

package command

import (
	"context"

	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func (h *FeatureCommandHandler) RenameFeature(ctx context.Context, cmd *proto.RenameFeatureCommand) error {
	err := h.feature.Rename(cmd.Name)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_RENAMED, &eventproto.FeatureRenamedEvent{
		Id:   h.feature.Id,
		Name: cmd.Name,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeDescription(ctx context.Context, cmd *proto.ChangeDescriptionCommand) error {
	err := h.feature.ChangeDescription(cmd.Description)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_FEATURE_DESCRIPTION_CHANGED,
		&eventproto.FeatureDescriptionChangedEvent{
			Id:          h.feature.Id,
			Description: cmd.Description,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) AddTag(ctx context.Context, cmd *proto.AddTagCommand) error {
	err := h.feature.AddTag(cmd.Tag)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_TAG_ADDED, &eventproto.FeatureTagAddedEvent{
		Id:  h.feature.Id,
		Tag: cmd.Tag,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) RemoveTag(ctx context.Context, cmd *proto.RemoveTagCommand) error {
	err := h.feature.RemoveTag(cmd.Tag)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_TAG_REMOVED, &eventproto.FeatureTagRemovedEvent{
		Id:  h.feature.Id,
		Tag: cmd.Tag,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}
