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

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	autoopspb "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventpb "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type webhookCommandHandler struct {
	editor               *eventpb.Editor
	publisher            publisher.Publisher
	webhook              *domain.Webhook
	environmentNamespace string
}

func NewWebhookCommandHandler(
	editor *eventpb.Editor,
	p publisher.Publisher,
	webhook *domain.Webhook,
	environmentNamespace string,
) Handler {
	return &webhookCommandHandler{
		editor:               editor,
		publisher:            p,
		webhook:              webhook,
		environmentNamespace: environmentNamespace,
	}
}

func (h *webhookCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *autoopspb.CreateWebhookCommand:
		return h.CreateWebhook(ctx, c)
	case *autoopspb.DeleteWebhookCommand:
		return h.DeleteWebhook(ctx, c)
	case *autoopspb.ChangeWebhookNameCommand:
		return h.ChangeWebhookName(ctx, c)
	case *autoopspb.ChangeWebhookDescriptionCommand:
		return h.ChangeWebhookDescription(ctx, c)
	default:
		return errUnknownCommand
	}
}

func (h *webhookCommandHandler) CreateWebhook(ctx context.Context, cmd *autoopspb.CreateWebhookCommand) error {
	event, err := domainevent.NewEvent(
		h.editor,
		eventpb.Event_WEBHOOK,
		h.webhook.Id,
		eventpb.Event_WEBHOOK_CREATED,
		&eventpb.WebhookCreatedEvent{
			Id:          h.webhook.Id,
			Name:        h.webhook.Name,
			Description: h.webhook.Description,
			CreatedAt:   h.webhook.CreatedAt,
			UpdatedAt:   h.webhook.UpdatedAt,
		},
		h.environmentNamespace,
	)
	if err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}

func (h *webhookCommandHandler) DeleteWebhook(ctx context.Context, cmd *autoopspb.DeleteWebhookCommand) error {
	event, err := domainevent.NewEvent(
		h.editor,
		eventpb.Event_WEBHOOK,
		h.webhook.Id,
		eventpb.Event_WEBHOOK_DELETED,
		&eventpb.WebhookDeletedEvent{
			Id: h.webhook.Id,
		},
		h.environmentNamespace,
	)
	if err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}

func (h *webhookCommandHandler) ChangeWebhookName(ctx context.Context, cmd *autoopspb.ChangeWebhookNameCommand) error {
	if err := h.webhook.ChangeName(cmd.Name); err != nil {
		return err
	}
	event, err := domainevent.NewEvent(
		h.editor,
		eventpb.Event_WEBHOOK,
		h.webhook.Id,
		eventpb.Event_WEBHOOK_NAME_CHANGED,
		&eventpb.WebhookNameChangedEvent{
			Id:   h.webhook.Id,
			Name: cmd.Name,
		},
		h.environmentNamespace,
	)
	if err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}

func (h *webhookCommandHandler) ChangeWebhookDescription(
	ctx context.Context,
	cmd *autoopspb.ChangeWebhookDescriptionCommand,
) error {
	if err := h.webhook.ChangeDescription(cmd.Description); err != nil {
		return err
	}
	event, err := domainevent.NewEvent(
		h.editor,
		eventpb.Event_WEBHOOK,
		h.webhook.Id,
		eventpb.Event_WEBHOOK_DESCRIPTION_CHANGED,
		&eventpb.WebhookDescriptionChangedEvent{
			Id:          h.webhook.Id,
			Description: cmd.Description,
		},
		h.environmentNamespace,
	)
	if err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}
