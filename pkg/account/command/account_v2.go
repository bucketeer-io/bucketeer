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

package command

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type accountV2CommandHandler struct {
	editor          *eventproto.Editor
	account         *domain.AccountV2
	previousAccount *domain.AccountV2
	publisher       publisher.Publisher
	organizationID  string
}

func NewAccountV2CommandHandler(
	editor *eventproto.Editor,
	account *domain.AccountV2,
	p publisher.Publisher,
	organizationID string,
) (Handler, error) {
	prev := &domain.AccountV2{}
	if err := copier.Copy(prev, account); err != nil {
		return nil, err
	}
	return &accountV2CommandHandler{
		editor:          editor,
		account:         account,
		previousAccount: prev,
		publisher:       p,
		organizationID:  organizationID,
	}, nil
}

func (h *accountV2CommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *accountproto.CreateAccountV2Command:
		return h.create(ctx, c)
	case *accountproto.ChangeAccountV2NameCommand:
		return h.changeName(ctx, c)
	case *accountproto.ChangeAccountV2FirstNameCommand:
		return h.changeFirstName(ctx, c)
	case *accountproto.ChangeAccountV2LastNameCommand:
		return h.changeLastName(ctx, c)
	case *accountproto.ChangeAccountV2LanguageCommand:
		return h.changeLanguage(ctx, c)
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
	case *accountproto.CreateSearchFilterCommand:
		return h.createSearchFilter(ctx, c)
	case *accountproto.ChangeSearchFilterNameCommand:
		return h.changeSearchFilerName(ctx, c)
	case *accountproto.ChangeSearchFilterQueryCommand:
		return h.changeSearchFilterQuery(ctx, c)
	case *accountproto.ChangeDefaultSearchFilterCommand:
		return h.changeDefaultSearchFilter(ctx, c)
	case *accountproto.DeleteSearchFilterCommand:
		return h.deleteSearchFiler(ctx, c)
	case *accountproto.ChangeAccountV2LastSeenCommand:
		return h.changeLastSeen(ctx, c)
	case *accountproto.ChangeAccountV2AvatarCommand:
		return h.changeAvatar(ctx, c)
	default:
		return ErrBadCommand
	}
}

func (h *accountV2CommandHandler) create(ctx context.Context, cmd *accountproto.CreateAccountV2Command) error {
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_CREATED, &eventproto.AccountV2CreatedEvent{
		Email:            h.account.Email,
		FirstName:        h.account.FirstName,
		LastName:         h.account.LastName,
		Language:         h.account.Language,
		AvatarImageUrl:   h.account.AvatarImageUrl,
		OrganizationId:   h.account.OrganizationId,
		OrganizationRole: h.account.OrganizationRole,
		EnvironmentRoles: h.account.EnvironmentRoles,
		Disabled:         h.account.Disabled,
		CreatedAt:        h.account.CreatedAt,
		UpdatedAt:        h.account.UpdatedAt,
	})
}

func (h *accountV2CommandHandler) changeName(
	ctx context.Context,
	cmd *accountproto.ChangeAccountV2NameCommand,
) error {
	if err := h.account.ChangeName(cmd.Name); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_NAME_CHANGED, &eventproto.AccountV2NameChangedEvent{
		Email: h.account.Email,
		Name:  cmd.Name,
	})
}

func (h *accountV2CommandHandler) changeFirstName(
	ctx context.Context,
	cmd *accountproto.ChangeAccountV2FirstNameCommand,
) error {
	if err := h.account.ChangeFirstName(cmd.FirstName); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_FIRST_NAME_CHANGED, &eventproto.AccountV2FirstNameChangedEvent{
		Email:     h.account.Email,
		FirstName: cmd.FirstName,
	})
}

func (h *accountV2CommandHandler) changeLastName(
	ctx context.Context,
	cmd *accountproto.ChangeAccountV2LastNameCommand,
) error {
	if err := h.account.ChangeLastName(cmd.LastName); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_LAST_NAME_CHANGED, &eventproto.AccountV2LastNameChangedEvent{
		Email:    h.account.Email,
		LastName: cmd.LastName,
	})
}

func (h *accountV2CommandHandler) changeLanguage(
	ctx context.Context,
	cmd *accountproto.ChangeAccountV2LanguageCommand,
) error {
	if err := h.account.ChangeLanguage(cmd.Language); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_LANGUAGE_CHANGED, &eventproto.AccountV2LanguageChangedEvent{
		Email:    h.account.Email,
		Language: cmd.Language,
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

func (h *accountV2CommandHandler) changeLastSeen(
	ctx context.Context,
	cmd *accountproto.ChangeAccountV2LastSeenCommand,
) error {
	return h.account.ChangeLastSeen(cmd.LastSeen)
}

func (h *accountV2CommandHandler) changeAvatar(
	ctx context.Context,
	cmd *accountproto.ChangeAccountV2AvatarCommand,
) error {
	return h.account.ChangeAvatar(cmd.AvatarImage, cmd.AvatarFileType)
}

func (h *accountV2CommandHandler) createSearchFilter(
	ctx context.Context,
	cmd *accountproto.CreateSearchFilterCommand) error {
	searchFilter, err := h.account.AddSearchFilter(
		cmd.Name,
		cmd.Query,
		cmd.FilterTargetType,
		cmd.EnvironmentId,
		cmd.DefaultFilter)
	if err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_CREATED_SEARCH_FILTER, &eventproto.SearchFilterCreatedEvent{
		Name:          searchFilter.Name,
		Query:         searchFilter.Query,
		TargetType:    searchFilter.FilterTargetType,
		EnvironmentId: searchFilter.EnvironmentId,
		DefaultFilter: searchFilter.DefaultFilter,
	})
}

func (h *accountV2CommandHandler) changeSearchFilerName(
	ctx context.Context,
	cmd *accountproto.ChangeSearchFilterNameCommand) error {
	if err := h.account.ChangeSearchFilterName(cmd.Id, cmd.Name); err != nil {
		return err
	}
	return h.send(
		ctx,
		eventproto.Event_ACCOUNT_V2_SEARCH_FILTER_NANE_CHANGED,
		&eventproto.SearchFilterNameChangedEvent{
			Id:   cmd.Id,
			Name: cmd.Name,
		},
	)
}

func (h *accountV2CommandHandler) changeSearchFilterQuery(
	ctx context.Context,
	cmd *accountproto.ChangeSearchFilterQueryCommand) error {
	if err := h.account.ChangeSearchFilterQuery(cmd.Id, cmd.Query); err != nil {
		return err
	}
	return h.send(
		ctx,
		eventproto.Event_ACCOUNT_V2_SEARCH_FILTER_QUERY_CHANGED,
		&eventproto.SearchFilterQueryChangedEvent{
			Id:    cmd.Id,
			Query: cmd.Query,
		},
	)
}

func (h *accountV2CommandHandler) changeDefaultSearchFilter(
	ctx context.Context,
	cmd *accountproto.ChangeDefaultSearchFilterCommand) error {
	if err := h.account.ChangeDefaultSearchFilter(cmd.Id, cmd.DefaultFilter); err != nil {
		return err
	}
	return h.send(
		ctx,
		eventproto.Event_ACCOUNT_V2_SEARCH_FILTER_DEFAULT_CHANGED,
		&eventproto.SearchFilterDefaultChangedEvent{
			Id:            cmd.Id,
			DefaultFilter: cmd.DefaultFilter,
		},
	)
}

func (h *accountV2CommandHandler) deleteSearchFiler(
	ctx context.Context,
	cmd *accountproto.DeleteSearchFilterCommand) error {
	if err := h.account.DeleteSearchFilter(cmd.Id); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_SEARCH_FILTER_DELETED, &eventproto.SearchFilterDeletedEvent{
		Id: cmd.Id,
	})
}

func (h *accountV2CommandHandler) send(
	ctx context.Context,
	eventType eventproto.Event_Type,
	event proto.Message,
) error {
	var prev *accountproto.AccountV2
	if h.previousAccount != nil && h.previousAccount.AccountV2 != nil {
		prev = h.previousAccount.AccountV2
	}
	e, err := domainevent.NewAdminEvent(
		h.editor,
		eventproto.Event_ACCOUNT,
		h.account.Email,
		eventType,
		event,
		h.account.AccountV2,
		prev,
	)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
