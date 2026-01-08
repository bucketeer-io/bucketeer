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

type goalCommandHandler struct {
	editor        *eventproto.Editor
	goal          *domain.Goal
	previousGoal  *domain.Goal
	publisher     publisher.Publisher
	environmentId string
}

func NewGoalCommandHandler(
	editor *eventproto.Editor,
	goal *domain.Goal,
	p publisher.Publisher,
	environmentId string,
) (Handler, error) {
	prev := &domain.Goal{}
	if err := copier.Copy(prev, goal); err != nil {
		return nil, err
	}
	return &goalCommandHandler{
		editor:        editor,
		goal:          goal,
		previousGoal:  prev,
		publisher:     p,
		environmentId: environmentId,
	}, nil
}

func (h *goalCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.ArchiveGoalCommand:
		return h.archive(ctx, c)
	default:
		return ErrUnknownCommand
	}
}

func (h *goalCommandHandler) archive(ctx context.Context, cmd *proto.ArchiveGoalCommand) error {
	if err := h.goal.SetArchived(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_GOAL_ARCHIVED, &eventproto.GoalArchivedEvent{
		Id: h.goal.Id,
	})
}

func (h *goalCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event pb.Message) error {
	var prev *proto.Goal
	if h.previousGoal != nil && h.previousGoal.Goal != nil {
		prev = h.previousGoal.Goal
	}
	e, err := domainevent.NewEvent(
		h.editor,
		eventproto.Event_GOAL,
		h.goal.Id,
		eventType,
		event,
		h.environmentId,
		h.goal.Goal,
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
