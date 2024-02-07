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

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	proto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type autoOpsRuleCommandHandler struct {
	editor               *eventproto.Editor
	autoOpsRule          *domain.AutoOpsRule
	publisher            publisher.Publisher
	environmentNamespace string
}

func NewAutoOpsCommandHandler(
	editor *eventproto.Editor,
	autoOpsRule *domain.AutoOpsRule,
	p publisher.Publisher,
	environmentNamespace string,
) Handler {
	return &autoOpsRuleCommandHandler{
		editor:               editor,
		autoOpsRule:          autoOpsRule,
		publisher:            p,
		environmentNamespace: environmentNamespace,
	}
}

func (h *autoOpsRuleCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateAutoOpsRuleCommand:
		return h.create(ctx, c)
	case *proto.ChangeAutoOpsRuleOpsTypeCommand:
		return h.changeOpsType(ctx, c)
	case *proto.DeleteAutoOpsRuleCommand:
		return h.delete(ctx, c)
	case *proto.ChangeAutoOpsRuleTriggeredAtCommand:
		return h.changeTriggeredAt(ctx, c)
	case *proto.AddOpsEventRateClauseCommand:
		return h.addOpsEventRateClause(ctx, c)
	case *proto.ChangeOpsEventRateClauseCommand:
		return h.changeOpsEventRateClause(ctx, c)
	case *proto.DeleteClauseCommand:
		return h.deleteClause(ctx, c)
	case *proto.AddDatetimeClauseCommand:
		return h.addDatetimeClause(ctx, c)
	case *proto.ChangeDatetimeClauseCommand:
		return h.changeDatetimeClause(ctx, c)
	}
	return errUnknownCommand
}

func (h *autoOpsRuleCommandHandler) create(ctx context.Context, cmd *proto.CreateAutoOpsRuleCommand) error {
	return h.send(ctx, eventproto.Event_AUTOOPS_RULE_CREATED, &eventproto.AutoOpsRuleCreatedEvent{
		FeatureId:   h.autoOpsRule.FeatureId,
		OpsType:     h.autoOpsRule.OpsType,
		Clauses:     h.autoOpsRule.Clauses,
		TriggeredAt: h.autoOpsRule.TriggeredAt,
		CreatedAt:   h.autoOpsRule.CreatedAt,
		UpdatedAt:   h.autoOpsRule.UpdatedAt,
	})
}

func (h *autoOpsRuleCommandHandler) changeOpsType(
	ctx context.Context,
	cmd *proto.ChangeAutoOpsRuleOpsTypeCommand,
) error {
	h.autoOpsRule.SetOpsType(cmd.OpsType)
	return h.send(ctx, eventproto.Event_AUTOOPS_RULE_OPS_TYPE_CHANGED, &eventproto.AutoOpsRuleOpsTypeChangedEvent{
		OpsType: h.autoOpsRule.OpsType,
	})
}

func (h *autoOpsRuleCommandHandler) delete(ctx context.Context, cmd *proto.DeleteAutoOpsRuleCommand) error {
	h.autoOpsRule.SetDeleted()
	return h.send(ctx, eventproto.Event_AUTOOPS_RULE_DELETED, &eventproto.AutoOpsRuleDeletedEvent{})
}

func (h *autoOpsRuleCommandHandler) changeTriggeredAt(
	ctx context.Context,
	cmd *proto.ChangeAutoOpsRuleTriggeredAtCommand,
) error {
	h.autoOpsRule.SetTriggeredAt()
	return h.send(
		ctx,
		eventproto.Event_AUTOOPS_RULE_TRIGGERED_AT_CHANGED,
		&eventproto.AutoOpsRuleTriggeredAtChangedEvent{},
	)
}

func (h *autoOpsRuleCommandHandler) addOpsEventRateClause(
	ctx context.Context,
	cmd *proto.AddOpsEventRateClauseCommand,
) error {
	clause, err := h.autoOpsRule.AddOpsEventRateClause(cmd.OpsEventRateClause)
	if err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_OPS_EVENT_RATE_CLAUSE_ADDED, &eventproto.OpsEventRateClauseAddedEvent{
		ClauseId:           clause.Id,
		OpsEventRateClause: cmd.OpsEventRateClause,
	})
}

func (h *autoOpsRuleCommandHandler) changeOpsEventRateClause(
	ctx context.Context,
	cmd *proto.ChangeOpsEventRateClauseCommand,
) error {
	if err := h.autoOpsRule.ChangeOpsEventRateClause(cmd.Id, cmd.OpsEventRateClause); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_OPS_EVENT_RATE_CLAUSE_CHANGED, &eventproto.OpsEventRateClauseChangedEvent{
		ClauseId:           cmd.Id,
		OpsEventRateClause: cmd.OpsEventRateClause,
	})
}

func (h *autoOpsRuleCommandHandler) deleteClause(ctx context.Context, cmd *proto.DeleteClauseCommand) error {
	if err := h.autoOpsRule.DeleteClause(cmd.Id); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_AUTOOPS_RULE_CLAUSE_DELETED, &eventproto.AutoOpsRuleClauseDeletedEvent{
		ClauseId: cmd.Id,
	})
}

func (h *autoOpsRuleCommandHandler) addDatetimeClause(ctx context.Context, cmd *proto.AddDatetimeClauseCommand) error {
	clause, err := h.autoOpsRule.AddDatetimeClause(cmd.DatetimeClause)
	if err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_DATETIME_CLAUSE_ADDED, &eventproto.DatetimeClauseAddedEvent{
		ClauseId:       clause.Id,
		DatetimeClause: cmd.DatetimeClause,
	})
}

func (h *autoOpsRuleCommandHandler) changeDatetimeClause(
	ctx context.Context,
	cmd *proto.ChangeDatetimeClauseCommand,
) error {
	if err := h.autoOpsRule.ChangeDatetimeClause(cmd.Id, cmd.DatetimeClause); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_DATETIME_CLAUSE_CHANGED, &eventproto.DatetimeClauseChangedEvent{
		ClauseId:       cmd.Id,
		DatetimeClause: cmd.DatetimeClause,
	})
}

func (h *autoOpsRuleCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event pb.Message) error {
	e, err := domainevent.NewEvent(
		h.editor,
		eventproto.Event_AUTOOPS_RULE,
		h.autoOpsRule.Id,
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
