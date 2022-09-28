// Copyright 2022 The Bucketeer Authors.
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

	"github.com/golang/protobuf/proto" // nolint:staticcheck

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type adminAccountCommandHandler struct {
	editor    *eventproto.Editor
	account   *domain.Account
	publisher publisher.Publisher
}

func NewAdminAccountCommandHandler(
	editor *eventproto.Editor,
	account *domain.Account,
	p publisher.Publisher,
) Handler {
	return &adminAccountCommandHandler{
		editor:    editor,
		account:   account,
		publisher: p,
	}
}

func (h *adminAccountCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *accountproto.CreateAdminAccountCommand:
		return h.create(ctx, c)
	case *accountproto.EnableAdminAccountCommand:
		return h.enable(ctx, c)
	case *accountproto.DisableAdminAccountCommand:
		return h.disable(ctx, c)
	default:
		return ErrBadCommand
	}
}

func (h *adminAccountCommandHandler) create(ctx context.Context, cmd *accountproto.CreateAdminAccountCommand) error {
	return h.send(ctx, eventproto.Event_ADMIN_ACCOUNT_CREATED, &eventproto.AdminAccountCreatedEvent{
		Id:        h.account.Id,
		Email:     h.account.Email,
		Name:      h.account.Name,
		Role:      h.account.Role,
		Disabled:  h.account.Disabled,
		CreatedAt: h.account.CreatedAt,
		UpdatedAt: h.account.UpdatedAt,
	})
}

func (h *adminAccountCommandHandler) enable(ctx context.Context, cmd *accountproto.EnableAdminAccountCommand) error {
	if err := h.account.Enable(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ADMIN_ACCOUNT_ENABLED, &eventproto.AdminAccountEnabledEvent{
		Id: h.account.Id,
	})
}

func (h *adminAccountCommandHandler) disable(ctx context.Context, cmd *accountproto.DisableAdminAccountCommand) error {
	if err := h.account.Disable(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ADMIN_ACCOUNT_DISABLED, &eventproto.AdminAccountDisabledEvent{
		Id: h.account.Id,
	})
}

func (h *adminAccountCommandHandler) send(
	ctx context.Context,
	eventType eventproto.Event_Type,
	event proto.Message,
) error {
	e, err := domainevent.NewAdminEvent(h.editor, eventproto.Event_ADMIN_ACCOUNT, h.account.Id, eventType, event)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
