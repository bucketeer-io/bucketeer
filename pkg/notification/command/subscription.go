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

type subscriptionCommandHandler struct {
	editor               *eventproto.Editor
	subscription         *domain.Subscription
	environmentNamespace string
	events               []*eventproto.Event
}

func NewSubscriptionCommandHandler(
	editor *eventproto.Editor,
	subscription *domain.Subscription,
	environmentNamespace string) Handler {
	return &subscriptionCommandHandler{
		editor:               editor,
		subscription:         subscription,
		environmentNamespace: environmentNamespace,
		events:               []*eventproto.Event{},
	}
}

// for unit test
func NewEmptySubscriptionCommandHandler() Handler {
	return &subscriptionCommandHandler{}
}

func (h *subscriptionCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateSubscriptionCommand:
		return h.create(ctx, c)
	case *proto.DeleteSubscriptionCommand:
		return h.delete(ctx, c)
	case *proto.AddSourceTypesCommand:
		return h.addSourceTypes(ctx, c)
	case *proto.DeleteSourceTypesCommand:
		return h.deleteSourceTypes(ctx, c)
	case *proto.EnableSubscriptionCommand:
		return h.enable(ctx, c)
	case *proto.DisableSubscriptionCommand:
		return h.disable(ctx, c)
	case *proto.RenameSubscriptionCommand:
		return h.rename(ctx, c)
	}
	return errUnknownCommand
}

func (h *subscriptionCommandHandler) create(ctx context.Context, cmd *proto.CreateSubscriptionCommand) error {
	return h.createEvent(ctx, eventproto.Event_SUBSCRIPTION_CREATED, &eventproto.SubscriptionCreatedEvent{
		SourceTypes: h.subscription.SourceTypes,
		Recipient:   h.subscription.Recipient,
		Name:        h.subscription.Name,
	})
}

func (h *subscriptionCommandHandler) delete(ctx context.Context, cmd *proto.DeleteSubscriptionCommand) error {
	return h.createEvent(ctx, eventproto.Event_SUBSCRIPTION_DELETED, &eventproto.SubscriptionDeletedEvent{})
}

func (h *subscriptionCommandHandler) addSourceTypes(ctx context.Context, cmd *proto.AddSourceTypesCommand) error {
	err := h.subscription.AddSourceTypes(cmd.SourceTypes)
	if err != nil {
		return err
	}
	return h.createEvent(
		ctx,
		eventproto.Event_SUBSCRIPTION_SOURCE_TYPE_ADDED,
		&eventproto.SubscriptionSourceTypesAddedEvent{
			SourceTypes: cmd.SourceTypes,
		},
	)
}

func (h *subscriptionCommandHandler) deleteSourceTypes(ctx context.Context, cmd *proto.DeleteSourceTypesCommand) error {
	err := h.subscription.DeleteSourceTypes(cmd.SourceTypes)
	if err != nil {
		return err
	}
	return h.createEvent(
		ctx,
		eventproto.Event_SUBSCRIPTION_SOURCE_TYPE_DELETED,
		&eventproto.SubscriptionSourceTypesDeletedEvent{
			SourceTypes: cmd.SourceTypes,
		},
	)
}

func (h *subscriptionCommandHandler) enable(ctx context.Context, cmd *proto.EnableSubscriptionCommand) error {
	if err := h.subscription.Enable(); err != nil {
		return err
	}
	return h.createEvent(ctx, eventproto.Event_SUBSCRIPTION_ENABLED, &eventproto.SubscriptionEnabledEvent{})
}

func (h *subscriptionCommandHandler) disable(ctx context.Context, cmd *proto.DisableSubscriptionCommand) error {
	err := h.subscription.Disable()
	if err != nil {
		return err
	}
	return h.createEvent(ctx, eventproto.Event_SUBSCRIPTION_DISABLED, &eventproto.SubscriptionDisabledEvent{})
}

func (h *subscriptionCommandHandler) rename(ctx context.Context, cmd *proto.RenameSubscriptionCommand) error {
	err := h.subscription.Rename(cmd.Name)
	if err != nil {
		return err
	}
	return h.createEvent(ctx, eventproto.Event_SUBSCRIPTION_RENAMED, &eventproto.SubscriptionRenamedEvent{Name: cmd.Name})
}

func (h *subscriptionCommandHandler) createEvent(
	ctx context.Context,
	eventType eventproto.Event_Type,
	event pb.Message,
) error {
	e, err := domainevent.NewEvent(
		h.editor,
		eventproto.Event_SUBSCRIPTION,
		h.subscription.Id,
		eventType,
		event,
		h.environmentNamespace,
	)
	if err != nil {
		return err
	}
	h.events = append(h.events, e)
	return nil
}

func (h *subscriptionCommandHandler) Events() []*eventproto.Event {
	return h.events
}
