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

	"github.com/golang/protobuf/proto" // nolint:staticcheck

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type apiKeyCommandHandler struct {
	editor               *eventproto.Editor
	apiKey               *domain.APIKey
	publisher            publisher.Publisher
	environmentNamespace string
}

func NewAPIKeyCommandHandler(
	editor *eventproto.Editor,
	apiKey *domain.APIKey,
	p publisher.Publisher,
	environmentNamespace string,
) Handler {
	return &apiKeyCommandHandler{
		editor:               editor,
		apiKey:               apiKey,
		publisher:            p,
		environmentNamespace: environmentNamespace,
	}
}

func (h *apiKeyCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *accountproto.CreateAPIKeyCommand:
		return h.create(ctx, c)
	case *accountproto.ChangeAPIKeyNameCommand:
		return h.rename(ctx, c)
	case *accountproto.EnableAPIKeyCommand:
		return h.enable(ctx, c)
	case *accountproto.DisableAPIKeyCommand:
		return h.disable(ctx, c)
	default:
		return ErrBadCommand
	}
}

func (h *apiKeyCommandHandler) create(ctx context.Context, cmd *accountproto.CreateAPIKeyCommand) error {
	return h.send(ctx, eventproto.Event_APIKEY_CREATED, &eventproto.APIKeyCreatedEvent{
		Id:        h.apiKey.Id,
		Name:      h.apiKey.Name,
		Role:      h.apiKey.Role,
		Disabled:  h.apiKey.Disabled,
		CreatedAt: h.apiKey.CreatedAt,
		UpdatedAt: h.apiKey.UpdatedAt,
	})
}

func (h *apiKeyCommandHandler) rename(ctx context.Context, cmd *accountproto.ChangeAPIKeyNameCommand) error {
	if err := h.apiKey.Rename(cmd.Name); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_APIKEY_NAME_CHANGED, &eventproto.APIKeyNameChangedEvent{
		Id:   h.apiKey.Id,
		Name: h.apiKey.Name,
	})
}

func (h *apiKeyCommandHandler) enable(ctx context.Context, cmd *accountproto.EnableAPIKeyCommand) error {
	if err := h.apiKey.Enable(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_APIKEY_ENABLED, &eventproto.APIKeyEnabledEvent{
		Id: h.apiKey.Id,
	})
}

func (h *apiKeyCommandHandler) disable(ctx context.Context, cmd *accountproto.DisableAPIKeyCommand) error {
	if err := h.apiKey.Disable(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_APIKEY_DISABLED, &eventproto.APIKeyDisabledEvent{
		Id: h.apiKey.Id,
	})
}

func (h *apiKeyCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event proto.Message) error {
	e, err := domainevent.NewEvent(
		h.editor,
		eventproto.Event_APIKEY,
		h.apiKey.Id,
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
