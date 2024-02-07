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

type organizationCommandHandler struct {
	editor       *eventproto.Editor
	organization *domain.Organization
	publisher    publisher.Publisher
}

func NewOrganizationCommandHandler(
	editor *eventproto.Editor,
	organization *domain.Organization,
	p publisher.Publisher,
) Handler {
	return &organizationCommandHandler{
		editor:       editor,
		organization: organization,
		publisher:    p,
	}
}

func (h *organizationCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateOrganizationCommand:
		return h.create(ctx, c)
	case *proto.ChangeNameOrganizationCommand:
		return h.changeName(ctx, c)
	case *proto.ChangeDescriptionOrganizationCommand:
		return h.changeDescription(ctx, c)
	case *proto.EnableOrganizationCommand:
		return h.enable(ctx, c)
	case *proto.DisableOrganizationCommand:
		return h.disable(ctx, c)
	case *proto.ArchiveOrganizationCommand:
		return h.archive(ctx, c)
	case *proto.UnarchiveOrganizationCommand:
		return h.unarchive(ctx, c)
	case *proto.ConvertTrialOrganizationCommand:
		return h.convertTrial(ctx, c)
	default:
		return errUnknownCommand
	}
}

func (h *organizationCommandHandler) create(ctx context.Context, cmd *proto.CreateOrganizationCommand) error {
	return h.send(ctx, eventproto.Event_ORGANIZATION_CREATED, &eventproto.OrganizationCreatedEvent{
		Id:          h.organization.Id,
		Name:        h.organization.Name,
		UrlCode:     h.organization.UrlCode,
		Description: h.organization.Description,
		Disabled:    h.organization.Disabled,
		Archived:    h.organization.Archived,
		Trial:       h.organization.Trial,
		CreatedAt:   h.organization.CreatedAt,
		UpdatedAt:   h.organization.UpdatedAt,
	})
}

func (h *organizationCommandHandler) changeName(ctx context.Context, cmd *proto.ChangeNameOrganizationCommand) error {
	newName := strings.TrimSpace(cmd.Name)
	h.organization.ChangeName(newName)
	return h.send(ctx, eventproto.Event_ORGANIZATION_NAME_CHANGED, &eventproto.OrganizationNameChangedEvent{
		Id:   h.organization.Id,
		Name: newName,
	})
}

func (h *organizationCommandHandler) changeDescription(
	ctx context.Context,
	cmd *proto.ChangeDescriptionOrganizationCommand,
) error {
	h.organization.ChangeDescription(cmd.Description)
	return h.send(ctx, eventproto.Event_ORGANIZATION_DESCRIPTION_CHANGED, &eventproto.OrganizationDescriptionChangedEvent{
		Id:          h.organization.Id,
		Description: cmd.Description,
	})
}

func (h *organizationCommandHandler) enable(ctx context.Context, cmd *proto.EnableOrganizationCommand) error {
	h.organization.Enable()
	return h.send(ctx, eventproto.Event_ORGANIZATION_ENABLED, &eventproto.OrganizationEnabledEvent{
		Id: h.organization.Id,
	})
}

func (h *organizationCommandHandler) disable(ctx context.Context, cmd *proto.DisableOrganizationCommand) error {
	if err := h.organization.Disable(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ORGANIZATION_DISABLED, &eventproto.OrganizationDisabledEvent{
		Id: h.organization.Id,
	})
}

func (h *organizationCommandHandler) archive(ctx context.Context, cmd *proto.ArchiveOrganizationCommand) error {
	if err := h.organization.Archive(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ORGANIZATION_ARCHIVED, &eventproto.OrganizationArchivedEvent{
		Id: h.organization.Id,
	})
}

func (h *organizationCommandHandler) unarchive(ctx context.Context, cmd *proto.UnarchiveOrganizationCommand) error {
	h.organization.Unarchive()
	return h.send(ctx, eventproto.Event_ORGANIZATION_UNARCHIVED, &eventproto.OrganizationUnarchivedEvent{
		Id: h.organization.Id,
	})
}

func (h *organizationCommandHandler) convertTrial(
	ctx context.Context,
	cmd *proto.ConvertTrialOrganizationCommand,
) error {
	h.organization.ConvertTrial()
	return h.send(ctx, eventproto.Event_ORGANIZATION_TRIAL_CONVERTED, &eventproto.OrganizationTrialConvertedEvent{
		Id: h.organization.Id,
	})
}

func (h *organizationCommandHandler) send(
	ctx context.Context,
	eventType eventproto.Event_Type,
	event pb.Message,
) error {
	e, err := domainevent.NewAdminEvent(h.editor, eventproto.Event_ORGANIZATION, h.organization.Id, eventType, event)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
