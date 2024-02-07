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

	pb "github.com/golang/protobuf/proto" // nolint:staticcheck

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

type goalCommandHandler struct {
	editor               *eventproto.Editor
	goal                 *domain.Goal
	publisher            publisher.Publisher
	environmentNamespace string
}

func NewGoalCommandHandler(
	editor *eventproto.Editor,
	goal *domain.Goal,
	p publisher.Publisher,
	environmentNamespace string,
) Handler {
	return &goalCommandHandler{
		editor:               editor,
		goal:                 goal,
		publisher:            p,
		environmentNamespace: environmentNamespace,
	}
}

func (h *goalCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateGoalCommand:
		return h.create(ctx, c)
	case *proto.RenameGoalCommand:
		return h.rename(ctx, c)
	case *proto.ChangeDescriptionGoalCommand:
		return h.changeDescription(ctx, c)
	case *proto.ArchiveGoalCommand:
		return h.archive(ctx, c)
	case *proto.DeleteGoalCommand:
		return h.delete(ctx, c)
	default:
		return ErrUnknownCommand
	}
}

func (h *goalCommandHandler) create(ctx context.Context, cmd *proto.CreateGoalCommand) error {
	return h.send(ctx, eventproto.Event_GOAL_CREATED, &eventproto.GoalCreatedEvent{
		Id:          h.goal.Id,
		Name:        h.goal.Name,
		Description: h.goal.Description,
		Deleted:     h.goal.Deleted,
		CreatedAt:   h.goal.CreatedAt,
		UpdatedAt:   h.goal.UpdatedAt,
	})
}

func (h *goalCommandHandler) rename(ctx context.Context, cmd *proto.RenameGoalCommand) error {
	if err := h.goal.Rename(cmd.Name); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_GOAL_RENAMED, &eventproto.GoalRenamedEvent{
		Id:   h.goal.Id,
		Name: cmd.Name,
	})
}

func (h *goalCommandHandler) changeDescription(ctx context.Context, cmd *proto.ChangeDescriptionGoalCommand) error {
	if err := h.goal.ChangeDescription(cmd.Description); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_GOAL_DESCRIPTION_CHANGED, &eventproto.GoalDescriptionChangedEvent{
		Id:          h.goal.Id,
		Description: cmd.Description,
	})
}

func (h *goalCommandHandler) archive(ctx context.Context, cmd *proto.ArchiveGoalCommand) error {
	if err := h.goal.SetArchived(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_GOAL_ARCHIVED, &eventproto.GoalArchivedEvent{
		Id: h.goal.Id,
	})
}

func (h *goalCommandHandler) delete(ctx context.Context, cmd *proto.DeleteGoalCommand) error {
	if err := h.goal.SetDeleted(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_GOAL_DELETED, &eventproto.GoalDeletedEvent{
		Id: h.goal.Id,
	})
}

func (h *goalCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event pb.Message) error {
	e, err := domainevent.NewEvent(h.editor, eventproto.Event_GOAL, h.goal.Id, eventType, event, h.environmentNamespace)
	if err != nil {
		return err
	}
	// TODO: more reliable
	// TODO: add metrics
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
