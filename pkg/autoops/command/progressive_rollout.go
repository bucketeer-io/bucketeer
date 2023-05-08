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

package command

import (
	"context"

	pb "github.com/golang/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type progressiveRolloutCommandHandler struct {
	editor               *eventproto.Editor
	progressiveRollout   *domain.ProgressiveRollout
	publisher            publisher.Publisher
	environmentNamespace string
}

func NewProgressiveRolloutCommandHandler(
	editor *eventproto.Editor,
	progressiveRollout *domain.ProgressiveRollout,
	publisher publisher.Publisher,
	environmentNamespace string,
) Handler {
	return &progressiveRolloutCommandHandler{
		editor:               editor,
		progressiveRollout:   progressiveRollout,
		publisher:            publisher,
		environmentNamespace: environmentNamespace,
	}
}

func (h *progressiveRolloutCommandHandler) Handle(
	ctx context.Context,
	cmd Command,
) error {
	switch c := cmd.(type) {
	case *autoopsproto.CreateProgressiveRolloutCommand:
		return h.create(ctx, c)
	case *autoopsproto.DeleteProgressiveRolloutCommand:
		return h.delete(ctx, c)
	case *autoopsproto.ChangeProgressiveRolloutScheduleTriggeredAtCommand:
		return h.changeTriggeredAt(ctx, c)
	}
	return errUnknownCommand
}

func (h *progressiveRolloutCommandHandler) create(
	ctx context.Context,
	c *autoopsproto.CreateProgressiveRolloutCommand,
) error {
	return h.send(
		ctx,
		eventproto.Event_PROGRESSIVE_ROLLOUT_CREATED,
		&eventproto.ProgressiveRolloutCreatedEvent{
			Id:        h.progressiveRollout.Id,
			FeatureId: h.progressiveRollout.FeatureId,
			Clause:    h.progressiveRollout.Clause,
			CreatedAt: h.progressiveRollout.CreatedAt,
			UpdatedAt: h.progressiveRollout.UpdatedAt,
			Type:      h.progressiveRollout.Type,
		},
	)
}

func (h *progressiveRolloutCommandHandler) delete(
	ctx context.Context,
	c *autoopsproto.DeleteProgressiveRolloutCommand,
) error {
	return h.send(
		ctx,
		eventproto.Event_PROGRESSIVE_ROLLOUT_DELETED,
		&eventproto.ProgressiveRolloutDeletedEvent{
			Id: h.progressiveRollout.Id,
		},
	)
}

func (h *progressiveRolloutCommandHandler) changeTriggeredAt(
	ctx context.Context,
	c *autoopsproto.ChangeProgressiveRolloutScheduleTriggeredAtCommand,
) error {
	if err := h.progressiveRollout.SetTriggeredAt(c.ScheduleId); err != nil {
		return err
	}
	return h.send(
		ctx,
		eventproto.Event_PROGRESSIVE_ROLLOUT_SCHEDULE_TRIGGERED_AT_CHANGED,
		&eventproto.ProgressiveRolloutScheduleTriggeredAtChangedEvent{
			ScheduleId: c.ScheduleId,
		},
	)
}

func (h *progressiveRolloutCommandHandler) send(
	ctx context.Context,
	eventType eventproto.Event_Type,
	event pb.Message,
) error {
	e, err := domainevent.NewEvent(
		h.editor,
		eventproto.Event_AUTOOPS_RULE,
		h.progressiveRollout.Id,
		eventType,
		event,
		h.environmentNamespace,
	)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
