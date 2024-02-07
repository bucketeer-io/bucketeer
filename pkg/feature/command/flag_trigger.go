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
//

package command

import (
	"context"

	pb "github.com/golang/protobuf/proto"

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type flagTriggerCommandHandler struct {
	editor               *eventproto.Editor
	flagTrigger          *domain.FlagTrigger
	publisher            publisher.Publisher
	environmentNamespace string
}

func NewFlagTriggerCommandHandler(
	editor *eventproto.Editor,
	flagTrigger *domain.FlagTrigger,
	publisher publisher.Publisher,
	environmentNamespace string,
) Handler {
	return &flagTriggerCommandHandler{
		editor:               editor,
		flagTrigger:          flagTrigger,
		publisher:            publisher,
		environmentNamespace: environmentNamespace,
	}
}

func (f *flagTriggerCommandHandler) Handle(
	ctx context.Context,
	cmd Command,
) error {
	switch c := cmd.(type) {
	case *proto.CreateFlagTriggerCommand:
		return f.create(ctx, c)
	case *proto.ResetFlagTriggerCommand:
		return f.reset(ctx, c)
	case *proto.ChangeFlagTriggerDescriptionCommand:
		return f.changeDescription(ctx, c)
	case *proto.DisableFlagTriggerCommand:
		return f.disable(ctx, c)
	case *proto.EnableFlagTriggerCommand:
		return f.enable(ctx, c)
	case *proto.DeleteFlagTriggerCommand:
		return f.delete(ctx, c)
	case *proto.UpdateFlagTriggerUsageCommand:
		return f.updateUsage(ctx, c)
	}
	return errBadCommand
}

func (f *flagTriggerCommandHandler) create(
	ctx context.Context,
	cmd *proto.CreateFlagTriggerCommand,
) error {
	if err := f.flagTrigger.GenerateToken(); err != nil {
		return err
	}
	return f.send(ctx, eventproto.Event_FLAG_TRIGGER_CREATED, &eventproto.FlagTriggerCreatedEvent{
		Id:                   f.flagTrigger.Id,
		FeatureId:            f.flagTrigger.FeatureId,
		EnvironmentNamespace: f.flagTrigger.EnvironmentNamespace,
		Type:                 f.flagTrigger.Type,
		Action:               f.flagTrigger.Action,
		Description:          f.flagTrigger.Description,
		CreatedAt:            f.flagTrigger.CreatedAt,
		UpdatedAt:            f.flagTrigger.UpdatedAt,
		Token:                f.flagTrigger.Token,
	})
}

func (f *flagTriggerCommandHandler) reset(
	ctx context.Context,
	cmd *proto.ResetFlagTriggerCommand,
) error {
	if err := f.flagTrigger.GenerateToken(); err != nil {
		return err
	}
	return f.send(ctx, eventproto.Event_FLAG_TRIGGER_RESET, &eventproto.FlagTriggerResetEvent{
		Id:                   f.flagTrigger.Id,
		FeatureId:            f.flagTrigger.FeatureId,
		EnvironmentNamespace: f.flagTrigger.EnvironmentNamespace,
		Token:                f.flagTrigger.Token,
	})
}

func (f *flagTriggerCommandHandler) changeDescription(
	ctx context.Context,
	cmd *proto.ChangeFlagTriggerDescriptionCommand,
) error {
	_ = f.flagTrigger.ChangeDescription(cmd.Description)
	return f.send(ctx,
		eventproto.Event_FLAG_TRIGGER_DESCRIPTION_CHANGED,
		&eventproto.FlagTriggerDescriptionChangedEvent{
			Id:                   f.flagTrigger.Id,
			FeatureId:            f.flagTrigger.FeatureId,
			EnvironmentNamespace: f.flagTrigger.EnvironmentNamespace,
			Description:          f.flagTrigger.Description,
		})
}

func (f *flagTriggerCommandHandler) disable(
	ctx context.Context,
	cmd *proto.DisableFlagTriggerCommand,
) error {
	_ = f.flagTrigger.Disable()
	return f.send(ctx, eventproto.Event_FLAG_TRIGGER_DISABLED, &eventproto.FlagTriggerDisabledEvent{
		Id:                   f.flagTrigger.Id,
		FeatureId:            f.flagTrigger.FeatureId,
		EnvironmentNamespace: f.flagTrigger.EnvironmentNamespace,
	})
}

func (f *flagTriggerCommandHandler) enable(
	ctx context.Context,
	cmd *proto.EnableFlagTriggerCommand,
) error {
	_ = f.flagTrigger.Enable()
	return f.send(ctx, eventproto.Event_FLAG_TRIGGER_ENABLED, &eventproto.FlagTriggerEnabledEvent{
		Id:                   f.flagTrigger.Id,
		FeatureId:            f.flagTrigger.FeatureId,
		EnvironmentNamespace: f.flagTrigger.EnvironmentNamespace,
	})
}

func (f *flagTriggerCommandHandler) updateUsage(
	ctx context.Context,
	c *proto.UpdateFlagTriggerUsageCommand,
) error {
	_ = f.flagTrigger.UpdateTriggerUsage()
	return f.send(ctx, eventproto.Event_FLAG_TRIGGER_USAGE_UPDATED, &eventproto.FlagTriggerUsageUpdatedEvent{
		Id:                   f.flagTrigger.Id,
		FeatureId:            f.flagTrigger.FeatureId,
		EnvironmentNamespace: f.flagTrigger.EnvironmentNamespace,
		LastTriggeredAt:      f.flagTrigger.LastTriggeredAt,
		TriggerTimes:         f.flagTrigger.TriggerCount,
	})
}

func (f *flagTriggerCommandHandler) delete(
	ctx context.Context,
	cmd *proto.DeleteFlagTriggerCommand,
) error {
	return f.send(ctx, eventproto.Event_FLAG_TRIGGER_DELETED, &eventproto.FlagTriggerDeletedEvent{
		Id:                   f.flagTrigger.Id,
		FeatureId:            f.flagTrigger.FeatureId,
		EnvironmentNamespace: f.flagTrigger.EnvironmentNamespace,
	})
}

func (f *flagTriggerCommandHandler) send(
	ctx context.Context,
	eventType eventproto.Event_Type,
	event pb.Message,
) error {
	e, err := domainevent.NewEvent(
		f.editor,
		eventproto.Event_FLAG_TRIGGER,
		f.flagTrigger.Id,
		eventType,
		event,
		f.environmentNamespace,
	)
	if err != nil {
		return err
	}
	if err := f.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
