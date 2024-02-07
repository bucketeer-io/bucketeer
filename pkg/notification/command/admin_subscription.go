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
	"github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

type adminSubscriptionCommandHandler struct {
	editor       *eventproto.Editor
	subscription *domain.Subscription
	events       []*eventproto.Event
}

func NewAdminSubscriptionCommandHandler(
	editor *eventproto.Editor,
	subscription *domain.Subscription) Handler {
	return &adminSubscriptionCommandHandler{
		editor:       editor,
		subscription: subscription,
		events:       []*eventproto.Event{},
	}
}

// for unit test
func NewEmptyAdminSubscriptionCommandHandler() Handler {
	return &adminSubscriptionCommandHandler{}
}

func (h *adminSubscriptionCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateAdminSubscriptionCommand:
		return h.create(ctx, c)
	case *proto.DeleteAdminSubscriptionCommand:
		return h.delete(ctx, c)
	case *proto.AddAdminSubscriptionSourceTypesCommand:
		return h.addSourceTypes(ctx, c)
	case *proto.DeleteAdminSubscriptionSourceTypesCommand:
		return h.deleteSourceTypes(ctx, c)
	case *proto.EnableAdminSubscriptionCommand:
		return h.enable(ctx, c)
	case *proto.DisableAdminSubscriptionCommand:
		return h.disable(ctx, c)
	case *proto.RenameAdminSubscriptionCommand:
		return h.rename(ctx, c)
	}
	return errUnknownCommand
}

func (h *adminSubscriptionCommandHandler) create(ctx context.Context, cmd *proto.CreateAdminSubscriptionCommand) error {
	return h.createEvent(ctx, eventproto.Event_ADMIN_SUBSCRIPTION_CREATED, &eventproto.AdminSubscriptionCreatedEvent{
		SourceTypes: h.subscription.SourceTypes,
		Recipient:   h.subscription.Recipient,
		Name:        h.subscription.Name,
	})
}

func (h *adminSubscriptionCommandHandler) delete(ctx context.Context, cmd *proto.DeleteAdminSubscriptionCommand) error {
	return h.createEvent(ctx, eventproto.Event_ADMIN_SUBSCRIPTION_DELETED, &eventproto.AdminSubscriptionDeletedEvent{})
}

func (h *adminSubscriptionCommandHandler) addSourceTypes(
	ctx context.Context,
	cmd *proto.AddAdminSubscriptionSourceTypesCommand,
) error {
	err := h.subscription.AddSourceTypes(cmd.SourceTypes)
	if err != nil {
		return err
	}
	return h.createEvent(
		ctx,
		eventproto.Event_ADMIN_SUBSCRIPTION_SOURCE_TYPE_ADDED,
		&eventproto.AdminSubscriptionSourceTypesAddedEvent{
			SourceTypes: cmd.SourceTypes,
		},
	)
}

func (h *adminSubscriptionCommandHandler) deleteSourceTypes(
	ctx context.Context,
	cmd *proto.DeleteAdminSubscriptionSourceTypesCommand,
) error {
	err := h.subscription.DeleteSourceTypes(cmd.SourceTypes)
	if err != nil {
		return err
	}
	return h.createEvent(
		ctx,
		eventproto.Event_ADMIN_SUBSCRIPTION_SOURCE_TYPE_DELETED,
		&eventproto.AdminSubscriptionSourceTypesDeletedEvent{
			SourceTypes: cmd.SourceTypes,
		},
	)
}

func (h *adminSubscriptionCommandHandler) enable(ctx context.Context, cmd *proto.EnableAdminSubscriptionCommand) error {
	if err := h.subscription.Enable(); err != nil {
		return err
	}
	return h.createEvent(ctx, eventproto.Event_ADMIN_SUBSCRIPTION_ENABLED, &eventproto.AdminSubscriptionEnabledEvent{})
}

func (h *adminSubscriptionCommandHandler) disable(
	ctx context.Context,
	cmd *proto.DisableAdminSubscriptionCommand,
) error {
	err := h.subscription.Disable()
	if err != nil {
		return err
	}
	return h.createEvent(ctx, eventproto.Event_ADMIN_SUBSCRIPTION_DISABLED, &eventproto.AdminSubscriptionDisabledEvent{})
}

func (h *adminSubscriptionCommandHandler) rename(ctx context.Context, cmd *proto.RenameAdminSubscriptionCommand) error {
	err := h.subscription.Rename(cmd.Name)
	if err != nil {
		return err
	}
	return h.createEvent(
		ctx,
		eventproto.Event_ADMIN_SUBSCRIPTION_RENAMED,
		&eventproto.AdminSubscriptionRenamedEvent{Name: cmd.Name},
	)
}

func (h *adminSubscriptionCommandHandler) createEvent(
	ctx context.Context,
	eventType eventproto.Event_Type,
	event pb.Message,
) error {
	e, err := domainevent.NewAdminEvent(h.editor, eventproto.Event_ADMIN_SUBSCRIPTION, h.subscription.Id, eventType, event)
	if err != nil {
		return err
	}
	h.events = append(h.events, e)
	return nil
}

func (h *adminSubscriptionCommandHandler) Events() []*eventproto.Event {
	return h.events
}
