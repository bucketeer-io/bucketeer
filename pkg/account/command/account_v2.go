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

	"github.com/golang/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type accountV2CommandHandler struct {
	editor         *eventproto.Editor
	account        *domain.AccountV2
	publisher      publisher.Publisher
	organizationID string
}

func NewAccountV2CommandHandler(
	editor *eventproto.Editor,
	account *domain.AccountV2,
	p publisher.Publisher,
	organizationID string,
) Handler {
	return &accountV2CommandHandler{
		editor:         editor,
		account:        account,
		publisher:      p,
		organizationID: organizationID,
	}
}

func (h *accountV2CommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *accountproto.CreateAccountV2Command:
		return h.create(ctx, c)
	case *accountproto.ChangeAccountV2NameCommand:
		return h.changeName(ctx, c)
	case *accountproto.ChangeAccountV2AvatarImageUrlCommand:
		return h.changeAvatarImageURL(ctx, c)
	case *accountproto.ChangeAccountV2OrganizationRoleCommand:
		return h.changeOrganizationRole(ctx, c)
	case *accountproto.ChangeAccountV2EnvironmentRolesCommand:
		return h.changeEnvironmentRoles(ctx, c)
	case *accountproto.EnableAccountV2Command:
		return h.enable(ctx, c)
	case *accountproto.DisableAccountV2Command:
		return h.disable(ctx, c)
	case *accountproto.DeleteAccountV2Command:
		return h.delete(ctx, c)
	default:
		return ErrBadCommand
	}
}

func (h *accountV2CommandHandler) create(ctx context.Context, cmd *accountproto.CreateAccountV2Command) error {
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_CREATED, &eventproto.AccountV2CreatedEvent{
		Email:            h.account.Email,
		Name:             h.account.Name,
		AvatarImageUrl:   h.account.AvatarImageUrl,
		OrganizationId:   h.account.OrganizationId,
		OrganizationRole: h.account.OrganizationRole,
		EnvironmentRoles: h.account.EnvironmentRoles,
		Disabled:         h.account.Disabled,
		CreatedAt:        h.account.CreatedAt,
		UpdatedAt:        h.account.UpdatedAt,
	})
}

func (h *accountV2CommandHandler) changeName(ctx context.Context, cmd *accountproto.ChangeAccountV2NameCommand) error {
	if err := h.account.ChangeName(cmd.Name); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_NAME_CHANGED, &eventproto.AccountV2NameChangedEvent{
		Email: h.account.Email,
		Name:  cmd.Name,
	})
}

func (h *accountV2CommandHandler) changeAvatarImageURL(
	ctx context.Context,
	cmd *accountproto.ChangeAccountV2AvatarImageUrlCommand,
) error {
	if err := h.account.ChangeAvatarImageURL(cmd.AvatarImageUrl); err != nil {
		return err
	}
	return h.send(
		ctx,
		eventproto.Event_ACCOUNT_V2_AVATAR_IMAGE_URL_CHANGED,
		&eventproto.AccountV2AvatarImageURLChangedEvent{
			Email:          h.account.Email,
			AvatarImageUrl: cmd.AvatarImageUrl,
		},
	)
}

func (h *accountV2CommandHandler) changeOrganizationRole(
	ctx context.Context,
	cmd *accountproto.ChangeAccountV2OrganizationRoleCommand,
) error {
	if err := h.account.ChangeOrganizationRole(cmd.Role); err != nil {
		return err
	}
	return h.send(
		ctx,
		eventproto.Event_ACCOUNT_V2_ORGANIZATION_ROLE_CHANGED,
		&eventproto.AccountV2OrganizationRoleChangedEvent{
			Email:            h.account.Email,
			OrganizationRole: cmd.Role,
		},
	)
}

func (h *accountV2CommandHandler) changeEnvironmentRoles(
	ctx context.Context,
	cmd *accountproto.ChangeAccountV2EnvironmentRolesCommand,
) error {
	if cmd.WriteType == accountproto.ChangeAccountV2EnvironmentRolesCommand_WriteType_OVERRIDE {
		if err := h.account.ChangeEnvironmentRole(cmd.Roles); err != nil {
			return err
		}
	} else if cmd.WriteType == accountproto.ChangeAccountV2EnvironmentRolesCommand_WriteType_PATCH {
		if err := h.account.PatchEnvironmentRole(cmd.Roles); err != nil {
			return err
		}
	}
	return h.send(
		ctx,
		eventproto.Event_ACCOUNT_V2_ENVIRONMENT_ROLES_CHANGED,
		&eventproto.AccountV2EnvironmentRolesChangedEvent{
			Email:            h.account.Email,
			EnvironmentRoles: cmd.Roles,
		},
	)
}

func (h *accountV2CommandHandler) enable(ctx context.Context, _ *accountproto.EnableAccountV2Command) error {
	if err := h.account.Enable(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_ENABLED, &eventproto.AccountV2EnabledEvent{
		Email: h.account.Email,
	})
}

func (h *accountV2CommandHandler) disable(ctx context.Context, _ *accountproto.DisableAccountV2Command) error {
	if err := h.account.Disable(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_DISABLED, &eventproto.AccountV2DisabledEvent{
		Email: h.account.Email,
	})
}

func (h *accountV2CommandHandler) delete(ctx context.Context, _ *accountproto.DeleteAccountV2Command) error {
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_DELETED, &eventproto.AccountV2DeletedEvent{
		Email: h.account.Email,
	})
}

func (h *accountV2CommandHandler) send(
	ctx context.Context,
	eventType eventproto.Event_Type,
	event proto.Message,
) error {
	e, err := domainevent.NewAdminEvent(
		h.editor,
		eventproto.Event_ACCOUNT,
		h.account.Email,
		eventType,
		event,
	)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
