// Copyright 2026 The Bucketeer Authors.
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
	"github.com/jinzhu/copier"

	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

type experimentCommandHandler struct {
	editor             *eventproto.Editor
	experiment         *domain.Experiment
	previousExperiment *domain.Experiment
	publisher          publisher.Publisher
	environmentId      string
}

func NewExperimentCommandHandler(
	editor *eventproto.Editor,
	experiment *domain.Experiment,
	p publisher.Publisher,
	environmentId string,
) (Handler, error) {
	prev := &domain.Experiment{}
	if err := copier.Copy(prev, experiment); err != nil {
		return nil, err
	}
	return &experimentCommandHandler{
		editor:             editor,
		experiment:         experiment,
		previousExperiment: prev,
		publisher:          p,
		environmentId:      environmentId,
	}, nil
}

func (h *experimentCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
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
	var prev *proto.Experiment
	if h.previousExperiment != nil && h.previousExperiment.Experiment != nil {
		prev = h.previousExperiment.Experiment
	}
	e, err := domainevent.NewEvent(
		h.editor,
		eventproto.Event_EXPERIMENT,
		h.experiment.Id,
		eventType,
		event,
		h.environmentId,
		h.experiment.Experiment,
		prev,
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
