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
	"strings"

	pb "github.com/golang/protobuf/proto"

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type environmentV2CommandHandler struct {
	editor      *eventproto.Editor
	environment *domain.EnvironmentV2
	publisher   publisher.Publisher
}

func NewEnvironmentV2CommandHandler(
	editor *eventproto.Editor,
	environment *domain.EnvironmentV2,
	p publisher.Publisher,
) Handler {
	return &environmentV2CommandHandler{
		editor:      editor,
		environment: environment,
		publisher:   p,
	}
}

func (h *environmentV2CommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateEnvironmentV2Command:
		return h.create(ctx, c)
	case *proto.RenameEnvironmentV2Command:
		return h.rename(ctx, c)
	case *proto.ChangeDescriptionEnvironmentV2Command:
		return h.changeDescription(ctx, c)
	case *proto.ArchiveEnvironmentV2Command:
		return h.archive(ctx, c)
	case *proto.UnarchiveEnvironmentV2Command:
		return h.unarchive(ctx, c)
	default:
		return errUnknownCommand
	}
}

func (h *environmentV2CommandHandler) create(ctx context.Context, _ *proto.CreateEnvironmentV2Command) error {
	return h.send(ctx, eventproto.Event_ENVIRONMENT_V2_CREATED, &eventproto.EnvironmentV2CreatedEvent{
		Id:          h.environment.Id,
		Name:        h.environment.Name,
		UrlCode:     h.environment.UrlCode,
		Description: h.environment.Description,
		ProjectId:   h.environment.ProjectId,
		Archived:    h.environment.Archived,
		CreatedAt:   h.environment.CreatedAt,
		UpdatedAt:   h.environment.UpdatedAt,
	})
}

func (h *environmentV2CommandHandler) rename(ctx context.Context, cmd *proto.RenameEnvironmentV2Command) error {
	oldName := h.environment.Name
	newName := strings.TrimSpace(cmd.Name)
	h.environment.Rename(cmd.Name)
	return h.send(ctx, eventproto.Event_ENVIRONMENT_V2_RENAMED, &eventproto.EnvironmentV2RenamedEvent{
		Id:        h.environment.Id,
		OldName:   oldName,
		NewName:   newName,
		ProjectId: h.environment.ProjectId,
	})
}

func (h *environmentV2CommandHandler) changeDescription(
	ctx context.Context,
	cmd *proto.ChangeDescriptionEnvironmentV2Command,
) error {
	oldDescription := h.environment.Description
	h.environment.ChangeDescription(cmd.Description)
	return h.send(
		ctx,
		eventproto.Event_ENVIRONMENT_V2_DESCRIPTION_CHANGED,
		&eventproto.EnvironmentV2DescriptionChangedEvent{
			Id:             h.environment.Id,
			Name:           h.environment.Name,
			ProjectId:      h.environment.ProjectId,
			OldDescription: oldDescription,
			NewDescription: cmd.Description,
		})
}

func (h *environmentV2CommandHandler) archive(ctx context.Context, _ *proto.ArchiveEnvironmentV2Command) error {
	h.environment.SetArchived()
	return h.send(ctx, eventproto.Event_ENVIRONMENT_V2_ARCHIVED, &eventproto.EnvironmentV2ArchivedEvent{
		Id:        h.environment.Id,
		Name:      h.environment.Name,
		ProjectId: h.environment.ProjectId,
	})
}

func (h *environmentV2CommandHandler) unarchive(ctx context.Context, _ *proto.UnarchiveEnvironmentV2Command) error {
	h.environment.SetUnarchived()
	return h.send(ctx, eventproto.Event_ENVIRONMENT_V2_UNARCHIVED, &eventproto.EnvironmentV2UnarchivedEvent{
		Id:        h.environment.Id,
		Name:      h.environment.Name,
		ProjectId: h.environment.ProjectId,
	})
}

func (h *environmentV2CommandHandler) send(
	ctx context.Context,
	eventType eventproto.Event_Type,
	event pb.Message,
) error {
	e, err := domainevent.NewAdminEvent(h.editor, eventproto.Event_ENVIRONMENT, h.environment.Id, eventType, event)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
