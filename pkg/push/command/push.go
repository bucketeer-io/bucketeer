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
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/push/domain"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/push"
)

type pushCommandHandler struct {
	editor               *eventproto.Editor
	push                 *domain.Push
	publisher            publisher.Publisher
	environmentNamespace string
}

func NewPushCommandHandler(
	editor *eventproto.Editor,
	push *domain.Push,
	p publisher.Publisher,
	environmentNamespace string,
) Handler {
	return &pushCommandHandler{
		editor:               editor,
		push:                 push,
		publisher:            p,
		environmentNamespace: environmentNamespace,
	}
}

func (h *pushCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreatePushCommand:
		return h.create(ctx, c)
	case *proto.DeletePushCommand:
		return h.delete(ctx, c)
	case *proto.AddPushTagsCommand:
		return h.addTags(ctx, c)
	case *proto.DeletePushTagsCommand:
		return h.deleteTags(ctx, c)
	case *proto.RenamePushCommand:
		return h.rename(ctx, c)
	}
	return errUnknownCommand
}

func (h *pushCommandHandler) create(ctx context.Context, cmd *proto.CreatePushCommand) error {
	return h.send(ctx, eventproto.Event_PUSH_CREATED, &eventproto.PushCreatedEvent{
		Name:      h.push.Name,
		FcmApiKey: h.push.FcmApiKey,
		Tags:      h.push.Tags,
	})
}

func (h *pushCommandHandler) delete(ctx context.Context, cmd *proto.DeletePushCommand) error {
	h.push.SetDeleted()
	return h.send(ctx, eventproto.Event_PUSH_DELETED, &eventproto.PushDeletedEvent{})
}

func (h *pushCommandHandler) addTags(ctx context.Context, cmd *proto.AddPushTagsCommand) error {
	err := h.push.AddTags(cmd.Tags)
	if err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_PUSH_TAGS_ADDED, &eventproto.PushTagsAddedEvent{
		Tags: cmd.Tags,
	})
}

func (h *pushCommandHandler) deleteTags(ctx context.Context, cmd *proto.DeletePushTagsCommand) error {
	err := h.push.DeleteTags(cmd.Tags)
	if err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_PUSH_TAGS_DELETED, &eventproto.PushTagsDeletedEvent{
		Tags: cmd.Tags,
	})
}

func (h pushCommandHandler) rename(ctx context.Context, cmd *proto.RenamePushCommand) error {
	if err := h.push.Rename(cmd.Name); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_PUSH_RENAMED, &eventproto.PushRenamedEvent{
		Name: cmd.Name,
	})
}

func (h *pushCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event pb.Message) error {
	e, err := domainevent.NewEvent(h.editor, eventproto.Event_PUSH, h.push.Id, eventType, event, h.environmentNamespace)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
