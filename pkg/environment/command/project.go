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

	pb "github.com/golang/protobuf/proto" // nolint:staticcheck

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type projectCommandHandler struct {
	editor    *eventproto.Editor
	project   *domain.Project
	publisher publisher.Publisher
}

func NewProjectCommandHandler(
	editor *eventproto.Editor,
	project *domain.Project,
	p publisher.Publisher,
) Handler {
	return &projectCommandHandler{
		editor:    editor,
		project:   project,
		publisher: p,
	}
}

func (h *projectCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateProjectCommand:
		return h.create(ctx, c)
	case *proto.CreateTrialProjectCommand:
		return h.createTrial(ctx, c)
	case *proto.ChangeDescriptionProjectCommand:
		return h.changeDescription(ctx, c)
	case *proto.RenameProjectCommand:
		return h.rename(ctx, c)
	case *proto.EnableProjectCommand:
		return h.enable(ctx, c)
	case *proto.DisableProjectCommand:
		return h.disable(ctx, c)
	case *proto.ConvertTrialProjectCommand:
		return h.convertTrial(ctx, c)
	default:
		return errUnknownCommand
	}
}

func (h *projectCommandHandler) create(ctx context.Context, cmd *proto.CreateProjectCommand) error {
	return h.send(ctx, eventproto.Event_PROJECT_CREATED, &eventproto.ProjectCreatedEvent{
		Id:           h.project.Id,
		Name:         h.project.Name,
		UrlCode:      h.project.UrlCode,
		Description:  h.project.Description,
		Disabled:     h.project.Disabled,
		Trial:        h.project.Trial,
		CreatorEmail: h.project.CreatorEmail,
		CreatedAt:    h.project.CreatedAt,
		UpdatedAt:    h.project.UpdatedAt,
	})
}

func (h *projectCommandHandler) createTrial(ctx context.Context, cmd *proto.CreateTrialProjectCommand) error {
	return h.send(ctx, eventproto.Event_PROJECT_TRIAL_CREATED, &eventproto.ProjectTrialCreatedEvent{
		Id:           h.project.Id,
		Name:         h.project.Name,
		UrlCode:      h.project.UrlCode,
		Description:  h.project.Description,
		Disabled:     h.project.Disabled,
		Trial:        h.project.Trial,
		CreatorEmail: h.project.CreatorEmail,
		CreatedAt:    h.project.CreatedAt,
		UpdatedAt:    h.project.UpdatedAt,
	})
}

func (h *projectCommandHandler) changeDescription(
	ctx context.Context,
	cmd *proto.ChangeDescriptionProjectCommand,
) error {
	h.project.ChangeDescription(cmd.Description)
	return h.send(ctx, eventproto.Event_PROJECT_DESCRIPTION_CHANGED, &eventproto.ProjectDescriptionChangedEvent{
		Id:          h.project.Id,
		Description: cmd.Description,
	})
}

func (h *projectCommandHandler) rename(ctx context.Context, cmd *proto.RenameProjectCommand) error {
	newName := strings.TrimSpace(cmd.Name)
	h.project.Rename(newName)
	return h.send(ctx, eventproto.Event_PROJECT_RENAMED, &eventproto.ProjectRenamedEvent{
		Id:   h.project.Id,
		Name: newName,
	})
}

func (h *projectCommandHandler) enable(ctx context.Context, cmd *proto.EnableProjectCommand) error {
	h.project.Enable()
	return h.send(ctx, eventproto.Event_PROJECT_ENABLED, &eventproto.ProjectEnabledEvent{
		Id: h.project.Id,
	})
}

func (h *projectCommandHandler) disable(ctx context.Context, cmd *proto.DisableProjectCommand) error {
	h.project.Disable()
	return h.send(ctx, eventproto.Event_PROJECT_DISABLED, &eventproto.ProjectDisabledEvent{
		Id: h.project.Id,
	})
}

func (h *projectCommandHandler) convertTrial(ctx context.Context, cmd *proto.ConvertTrialProjectCommand) error {
	h.project.ConvertTrial()
	return h.send(ctx, eventproto.Event_PROJECT_TRIAL_CONVERTED, &eventproto.ProjectTrialConvertedEvent{
		Id: h.project.Id,
	})
}

func (h *projectCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event pb.Message) error {
	e, err := domainevent.NewAdminEvent(h.editor, eventproto.Event_PROJECT, h.project.Id, eventType, event)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
