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

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type segmentCommandHandler struct {
	editor               *eventproto.Editor
	segment              *domain.Segment
	publisher            publisher.Publisher
	environmentNamespace string
}

func NewSegmentCommandHandler(
	editor *eventproto.Editor,
	segment *domain.Segment,
	publisher publisher.Publisher,
	environmentNamespace string,
) Handler {
	return &segmentCommandHandler{
		editor:               editor,
		segment:              segment,
		publisher:            publisher,
		environmentNamespace: environmentNamespace,
	}
}

func (h *segmentCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *featureproto.CreateSegmentCommand:
		return h.CreateSegment(ctx, c)
	case *featureproto.DeleteSegmentCommand:
		return h.DeleteSegment(ctx)
	case *featureproto.ChangeSegmentNameCommand:
		return h.ChangeName(ctx, c)
	case *featureproto.ChangeSegmentDescriptionCommand:
		return h.ChangeDescription(ctx, c)
	case *featureproto.AddRuleCommand:
		return h.AddRule(ctx, c)
	case *featureproto.DeleteRuleCommand:
		return h.DeleteRule(ctx, c)
	case *featureproto.AddClauseCommand:
		return h.AddClause(ctx, c)
	case *featureproto.DeleteClauseCommand:
		return h.DeleteClause(ctx, c)
	case *featureproto.ChangeClauseAttributeCommand:
		return h.ChangeClauseAttribute(ctx, c)
	case *featureproto.ChangeClauseOperatorCommand:
		return h.ChangeClauseOperator(ctx, c)
	case *featureproto.AddClauseValueCommand:
		return h.AddClauseValue(ctx, c)
	case *featureproto.RemoveClauseValueCommand:
		return h.RemoveClauseValue(ctx, c)
	case *featureproto.AddSegmentUserCommand:
		return h.AddSegmentUser(ctx, c)
	case *featureproto.DeleteSegmentUserCommand:
		return h.DeleteSegmentUser(ctx, c)
	case *featureproto.BulkUploadSegmentUsersCommand:
		return h.BulkUploadSegmentUsers(ctx, c)
	case *featureproto.ChangeBulkUploadSegmentUsersStatusCommand:
		return h.ChangeBulkUploadSegmentUsersStatus(ctx, c)
	default:
		return errBadCommand
	}
}

func (h *segmentCommandHandler) CreateSegment(ctx context.Context, cmd *featureproto.CreateSegmentCommand) error {
	return h.send(ctx, eventproto.Event_SEGMENT_CREATED, &eventproto.SegmentCreatedEvent{
		Id:          h.segment.Id,
		Name:        h.segment.Name,
		Description: h.segment.Description,
	})
}

func (h *segmentCommandHandler) DeleteSegment(ctx context.Context) error {
	if err := h.segment.SetDeleted(); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_DELETED, &eventproto.SegmentDeletedEvent{
		Id: h.segment.Id,
	})
}

func (h *segmentCommandHandler) ChangeName(ctx context.Context, cmd *featureproto.ChangeSegmentNameCommand) error {
	if err := h.segment.ChangeName(cmd.Name); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_NAME_CHANGED, &eventproto.SegmentNameChangedEvent{
		Id:   h.segment.Id,
		Name: h.segment.Name,
	})
}

func (h *segmentCommandHandler) ChangeDescription(
	ctx context.Context,
	cmd *featureproto.ChangeSegmentDescriptionCommand,
) error {
	if err := h.segment.ChangeDescription(cmd.Description); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_DESCRIPTION_CHANGED, &eventproto.SegmentDescriptionChangedEvent{
		Id:          h.segment.Id,
		Description: h.segment.Description,
	})
}

func (h *segmentCommandHandler) AddRule(ctx context.Context, cmd *featureproto.AddRuleCommand) error {
	for _, clause := range cmd.Rule.Clauses {
		id, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		clause.Id = id.String()
	}
	if err := h.segment.AddRule(cmd.Rule); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_RULE_ADDED, &eventproto.SegmentRuleAddedEvent{
		Id:   h.segment.Id,
		Rule: cmd.Rule,
	})
}

func (h *segmentCommandHandler) DeleteRule(ctx context.Context, cmd *featureproto.DeleteRuleCommand) error {
	if err := h.segment.DeleteRule(cmd.Id); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_RULE_DELETED, &eventproto.SegmentRuleDeletedEvent{
		Id:     h.segment.Id,
		RuleId: cmd.Id,
	})
}

func (h *segmentCommandHandler) AddClause(ctx context.Context, cmd *featureproto.AddClauseCommand) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	cmd.Clause.Id = id.String()
	if err := h.segment.AddClause(cmd.RuleId, cmd.Clause); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_RULE_CLAUSE_ADDED, &eventproto.SegmentRuleClauseAddedEvent{
		SegmentId: h.segment.Id,
		RuleId:    cmd.RuleId,
		Clause:    cmd.Clause,
	})
}

func (h *segmentCommandHandler) DeleteClause(ctx context.Context, cmd *featureproto.DeleteClauseCommand) error {
	if err := h.segment.DeleteClause(cmd.RuleId, cmd.Id); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_RULE_CLAUSE_DELETED, &eventproto.SegmentRuleClauseDeletedEvent{
		SegmentId: h.segment.Id,
		RuleId:    cmd.RuleId,
		ClauseId:  cmd.Id,
	})
}

func (h *segmentCommandHandler) ChangeClauseAttribute(
	ctx context.Context,
	cmd *featureproto.ChangeClauseAttributeCommand,
) error {
	if err := h.segment.ChangeClauseAttribute(cmd.RuleId, cmd.Id, cmd.Attribute); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_CLAUSE_ATTRIBUTE_CHANGED, &eventproto.SegmentClauseAttributeChangedEvent{
		SegmentId: h.segment.Id,
		RuleId:    cmd.RuleId,
		ClauseId:  cmd.Id,
		Attribute: cmd.Attribute,
	})
}

func (h *segmentCommandHandler) ChangeClauseOperator(
	ctx context.Context,
	cmd *featureproto.ChangeClauseOperatorCommand,
) error {
	if err := h.segment.ChangeClauseOperator(cmd.RuleId, cmd.Id, cmd.Operator); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_CLAUSE_OPERATOR_CHANGED, &eventproto.SegmentClauseOperatorChangedEvent{
		SegmentId: h.segment.Id,
		RuleId:    cmd.RuleId,
		ClauseId:  cmd.Id,
		Operator:  cmd.Operator,
	})
}

func (h *segmentCommandHandler) AddClauseValue(ctx context.Context, cmd *featureproto.AddClauseValueCommand) error {
	if err := h.segment.AddClauseValue(cmd.RuleId, cmd.Id, cmd.Value); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_CLAUSE_VALUE_ADDED, &eventproto.SegmentClauseValueAddedEvent{
		SegmentId: h.segment.Id,
		RuleId:    cmd.RuleId,
		ClauseId:  cmd.Id,
		Value:     cmd.Value,
	})
}

func (h *segmentCommandHandler) RemoveClauseValue(
	ctx context.Context,
	cmd *featureproto.RemoveClauseValueCommand,
) error {
	if err := h.segment.RemoveClauseValue(cmd.RuleId, cmd.Id, cmd.Value); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_SEGMENT_CLAUSE_VALUE_REMOVED, &eventproto.SegmentClauseValueRemovedEvent{
		SegmentId: h.segment.Id,
		RuleId:    cmd.RuleId,
		ClauseId:  cmd.Id,
		Value:     cmd.Value,
	})
}

func (h *segmentCommandHandler) AddSegmentUser(ctx context.Context, cmd *featureproto.AddSegmentUserCommand) error {
	count := int64(len(cmd.UserIds))
	switch cmd.State {
	case featureproto.SegmentUser_INCLUDED:
		h.segment.AddIncludedUserCount(count)
	}
	return h.send(ctx, eventproto.Event_SEGMENT_USER_ADDED, &eventproto.SegmentUserAddedEvent{
		SegmentId: h.segment.Id,
		UserIds:   cmd.UserIds,
		State:     cmd.State,
	})
}

func (h *segmentCommandHandler) DeleteSegmentUser(
	ctx context.Context,
	cmd *featureproto.DeleteSegmentUserCommand,
) error {
	count := int64(len(cmd.UserIds))
	switch cmd.State {
	case featureproto.SegmentUser_INCLUDED:
		h.segment.RemoveIncludedUserCount(count)
	}
	return h.send(ctx, eventproto.Event_SEGMENT_USER_DELETED, &eventproto.SegmentUserDeletedEvent{
		SegmentId: h.segment.Id,
		UserIds:   cmd.UserIds,
		State:     cmd.State,
	})
}

func (h *segmentCommandHandler) BulkUploadSegmentUsers(
	ctx context.Context,
	cmd *featureproto.BulkUploadSegmentUsersCommand,
) error {
	h.segment.SetStatus(featureproto.Segment_UPLOADING)
	return h.send(ctx, eventproto.Event_SEGMENT_BULK_UPLOAD_USERS, &eventproto.SegmentBulkUploadUsersEvent{
		SegmentId: h.segment.Id,
		Status:    featureproto.Segment_UPLOADING,
		State:     cmd.State,
	})
}

func (h *segmentCommandHandler) ChangeBulkUploadSegmentUsersStatus(
	ctx context.Context,
	cmd *featureproto.ChangeBulkUploadSegmentUsersStatusCommand,
) error {
	h.segment.SetStatus(cmd.Status)
	switch cmd.State {
	case featureproto.SegmentUser_INCLUDED:
		h.segment.SetIncludedUserCount(cmd.Count)
	}
	return h.send(
		ctx,
		eventproto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED,
		&eventproto.SegmentBulkUploadUsersStatusChangedEvent{
			SegmentId: h.segment.Id,
			Status:    cmd.Status,
			State:     cmd.State,
			Count:     cmd.Count,
		},
	)
}

func (h *segmentCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event proto.Message) error {
	e, err := domainevent.NewEvent(
		h.editor,
		eventproto.Event_SEGMENT,
		h.segment.Id,
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
