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
	"github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

type experimentCommandHandler struct {
	editor               *eventproto.Editor
	experiment           *domain.Experiment
	publisher            publisher.Publisher
	environmentNamespace string
}

func NewExperimentCommandHandler(
	editor *eventproto.Editor,
	experiment *domain.Experiment,
	p publisher.Publisher,
	environmentNamespace string,
) Handler {
	return &experimentCommandHandler{
		editor:               editor,
		experiment:           experiment,
		publisher:            p,
		environmentNamespace: environmentNamespace,
	}
}

func (h *experimentCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateExperimentCommand:
		return h.create(ctx, c)
	case *proto.ChangeExperimentPeriodCommand:
		return h.changePeriod(ctx, c)
	case *proto.ChangeExperimentNameCommand:
		return h.changeName(ctx, c)
	case *proto.ChangeExperimentDescriptionCommand:
		return h.changeDescription(ctx, c)
	case *proto.StopExperimentCommand:
		return h.stop(ctx, c)
	case *proto.StartExperimentCommand:
		return h.start(ctx, c)
	case *proto.FinishExperimentCommand:
		return h.finish(ctx, c)
	case *proto.ArchiveExperimentCommand:
		return h.archive(ctx, c)
	case *proto.DeleteExperimentCommand:
		return h.delete(ctx, c)
	default:
		return ErrUnknownCommand
	}
}

func (h *experimentCommandHandler) create(ctx context.Context, cmd *proto.CreateExperimentCommand) error {
	return h.send(ctx, eventproto.Event_EXPERIMENT_CREATED, &eventproto.ExperimentCreatedEvent{
		Id:              h.experiment.Id,
		FeatureId:       h.experiment.FeatureId,
		FeatureVersion:  h.experiment.FeatureVersion,
		Variations:      h.experiment.Variations,
		GoalId:          h.experiment.GoalId,
		GoalIds:         h.experiment.GoalIds,
		StartAt:         h.experiment.StartAt,
		StopAt:          h.experiment.StopAt,
		Stopped:         h.experiment.Stopped,
		StoppedAt:       h.experiment.StoppedAt,
		CreatedAt:       h.experiment.CreatedAt,
		UpdatedAt:       h.experiment.UpdatedAt,
		Name:            h.experiment.Name,
		Description:     h.experiment.Description,
		BaseVariationId: h.experiment.BaseVariationId,
	})
}

func (h *experimentCommandHandler) changePeriod(ctx context.Context, cmd *proto.ChangeExperimentPeriodCommand) error {
	if err := h.experiment.ChangePeriod(cmd.StartAt, cmd.StopAt); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_EXPERIMENT_PERIOD_CHANGED, &eventproto.ExperimentPeriodChangedEvent{
		Id:      h.experiment.Id,
		StartAt: cmd.StartAt,
		StopAt:  cmd.StopAt,
	})
}

func (h *experimentCommandHandler) changeName(ctx context.Context, cmd *proto.ChangeExperimentNameCommand) error {
	if err := h.experiment.ChangeName(cmd.Name); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_EXPERIMENT_NAME_CHANGED, &eventproto.ExperimentNameChangedEvent{
		Id:   h.experiment.Id,
		Name: h.experiment.Name,
	})
}

func (h *experimentCommandHandler) changeDescription(
	ctx context.Context,
	cmd *proto.ChangeExperimentDescriptionCommand,
) error {
	if err := h.experiment.ChangeDescription(cmd.Description); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_EXPERIMENT_DESCRIPTION_CHANGED, &eventproto.ExperimentDescriptionChangedEvent{
		Id:          h.experiment.Id,
		Description: h.experiment.Description,
	})
}

func (h *experimentCommandHandler) stop(ctx context.Context, cmd *proto.StopExperimentCommand) error {
	if err := h.experiment.Stop(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_EXPERIMENT_STOPPED, &eventproto.ExperimentStoppedEvent{
		Id:        h.experiment.Id,
		StoppedAt: h.experiment.StoppedAt,
	})
}

func (h *experimentCommandHandler) archive(ctx context.Context, cmd *proto.ArchiveExperimentCommand) error {
	if err := h.experiment.SetArchived(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_EXPERIMENT_ARCHIVED, &eventproto.ExperimentArchivedEvent{
		Id: h.experiment.Id,
	})
}

func (h *experimentCommandHandler) delete(ctx context.Context, cmd *proto.DeleteExperimentCommand) error {
	if err := h.experiment.SetDeleted(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_EXPERIMENT_DELETED, &eventproto.ExperimentDeletedEvent{
		Id: h.experiment.Id,
	})
}

func (h *experimentCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event pb.Message) error {
	e, err := domainevent.NewEvent(
		h.editor,
		eventproto.Event_EXPERIMENT,
		h.experiment.Id,
		eventType,
		event,
		h.environmentNamespace,
	)
	if err != nil {
		return err
	}
	// TODO: more reliable
	// TODO: add metrics
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}

func (h *experimentCommandHandler) start(ctx context.Context, cmd *proto.StartExperimentCommand) error {
	if err := h.experiment.Start(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_EXPERIMENT_STARTED, &eventproto.ExperimentStartedEvent{})
}

func (h *experimentCommandHandler) finish(ctx context.Context, cmd *proto.FinishExperimentCommand) error {
	if err := h.experiment.Finish(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_EXPERIMENT_FINISHED, &eventproto.ExperimentFinishedEvent{})
}
