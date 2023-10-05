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

	pb "github.com/golang/protobuf/proto" // nolint:staticcheck

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type environmentCommandHandler struct {
	editor      *eventproto.Editor
	environment *domain.Environment
	publisher   publisher.Publisher
}

func NewEnvironmentCommandHandler(
	editor *eventproto.Editor,
	environment *domain.Environment,
	p publisher.Publisher,
) Handler {
	return &environmentCommandHandler{
		editor:      editor,
		environment: environment,
		publisher:   p,
	}
}

func (h *environmentCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateEnvironmentCommand:
		return h.create(ctx, c)
	case *proto.RenameEnvironmentCommand:
		return h.rename(ctx, c)
	case *proto.ChangeDescriptionEnvironmentCommand:
		return h.changeDescription(ctx, c)
	case *proto.DeleteEnvironmentCommand:
		return h.delete(ctx, c)
	default:
		return errUnknownCommand
	}
}

func (h *environmentCommandHandler) create(ctx context.Context, cmd *proto.CreateEnvironmentCommand) error {
	return h.send(ctx, eventproto.Event_ENVIRONMENT_CREATED, &eventproto.EnvironmentCreatedEvent{
		Id:          h.environment.Id,
		Namespace:   h.environment.Namespace,
		Name:        h.environment.Name,
		Description: h.environment.Description,
		Deleted:     h.environment.Deleted,
		CreatedAt:   h.environment.CreatedAt,
		UpdatedAt:   h.environment.UpdatedAt,
		ProjectId:   h.environment.ProjectId,
	})
}

func (h *environmentCommandHandler) rename(ctx context.Context, cmd *proto.RenameEnvironmentCommand) error {
	h.environment.Rename(cmd.Name)
	return h.send(ctx, eventproto.Event_ENVIRONMENT_RENAMED, &eventproto.EnvironmentRenamedEvent{
		Id:   h.environment.Id,
		Name: cmd.Name,
	})
}

func (h *environmentCommandHandler) changeDescription(
	ctx context.Context,
	cmd *proto.ChangeDescriptionEnvironmentCommand,
) error {
	h.environment.ChangeDescription(cmd.Description)
	return h.send(ctx, eventproto.Event_ENVIRONMENT_DESCRIPTION_CHANGED, &eventproto.EnvironmentDescriptionChangedEvent{
		Id:          h.environment.Id,
		Description: cmd.Description,
	})
}

func (h *environmentCommandHandler) delete(ctx context.Context, cmd *proto.DeleteEnvironmentCommand) error {
	h.environment.SetDeleted()
	return h.send(ctx, eventproto.Event_ENVIRONMENT_DELETED, &eventproto.EnvironmentDeletedEvent{
		Id:        h.environment.Id,
		Namespace: h.environment.Namespace,
	})
}

func (h *environmentCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event pb.Message) error {
	e, err := domainevent.NewAdminEvent(h.editor, eventproto.Event_ENVIRONMENT, h.environment.Id, eventType, event)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
