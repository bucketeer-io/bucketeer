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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	domain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestChangeOpsType(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		input    proto.OpsType
		expected error
	}{
		{
			input:    proto.OpsType_DISABLE_FEATURE,
			expected: nil,
		},
	}
	for _, p := range patterns {
		m := publishermock.NewMockPublisher(mockController)
		a := newAutoOpsRule(t)
		h := newAutoOpsRuleCommandHandler(m, a)
		if p.expected == nil {
			m.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
		}
		cmd := &proto.ChangeAutoOpsRuleOpsTypeCommand{OpsType: proto.OpsType_DISABLE_FEATURE}
		err := h.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, p.input, a.OpsType)
	}
}

func TestDelete(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		expected error
	}{
		{
			expected: nil,
		},
	}
	for _, p := range patterns {
		m := publishermock.NewMockPublisher(mockController)
		a := newAutoOpsRule(t)
		h := newAutoOpsRuleCommandHandler(m, a)
		if p.expected == nil {
			m.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
		}
		cmd := &proto.DeleteAutoOpsRuleCommand{}
		err := h.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
}

func TestChangeTriggeredAt(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		expected error
	}{
		{
			expected: nil,
		},
	}
	for _, p := range patterns {
		m := publishermock.NewMockPublisher(mockController)
		a := newAutoOpsRule(t)
		h := newAutoOpsRuleCommandHandler(m, a)
		if p.expected == nil {
			m.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
		}
		cmd := &proto.ChangeAutoOpsRuleTriggeredAtCommand{}
		err := h.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
}

func TestAddOpsEventRateClause(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		input    *proto.OpsEventRateClause
		expected error
	}{
		{
			input:    &proto.OpsEventRateClause{},
			expected: nil,
		},
	}
	for _, p := range patterns {
		m := publishermock.NewMockPublisher(mockController)
		a := newAutoOpsRule(t)
		l := len(a.Clauses)
		h := newAutoOpsRuleCommandHandler(m, a)
		if p.expected == nil {
			m.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
		}
		cmd := &proto.AddOpsEventRateClauseCommand{OpsEventRateClause: p.input}
		err := h.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, l+1, len(a.Clauses))
	}
}

func TestChangeOpsEventRateClause(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		input    *proto.OpsEventRateClause
		expected error
	}{
		{
			input:    &proto.OpsEventRateClause{},
			expected: nil,
		},
	}
	for _, p := range patterns {
		m := publishermock.NewMockPublisher(mockController)
		a := newAutoOpsRule(t)
		h := newAutoOpsRuleCommandHandler(m, a)
		if p.expected == nil {
			m.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
		}
		cmd := &proto.ChangeOpsEventRateClauseCommand{Id: a.Clauses[0].Id, OpsEventRateClause: p.input}
		err := h.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
}

func TestDeleteClause(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		expected error
	}{
		{
			expected: nil,
		},
	}
	for _, p := range patterns {
		m := publishermock.NewMockPublisher(mockController)
		a := newAutoOpsRule(t)
		l := len(a.Clauses)
		h := newAutoOpsRuleCommandHandler(m, a)
		if p.expected == nil {
			m.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
		}
		cmd := &proto.DeleteClauseCommand{Id: a.Clauses[0].Id}
		err := h.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, l-1, len(a.Clauses))
	}
}

func TestAddDatetimeClause(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		input    *proto.DatetimeClause
		expected error
	}{
		{
			input:    &proto.DatetimeClause{},
			expected: nil,
		},
	}
	for _, p := range patterns {
		m := publishermock.NewMockPublisher(mockController)
		a := newAutoOpsRule(t)
		l := len(a.Clauses)
		h := newAutoOpsRuleCommandHandler(m, a)
		if p.expected == nil {
			m.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
		}
		cmd := &proto.AddDatetimeClauseCommand{DatetimeClause: p.input}
		err := h.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, l+1, len(a.Clauses))
	}
}

func TestChangeDatetimeClause(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		input    *proto.DatetimeClause
		expected error
	}{
		{
			input:    &proto.DatetimeClause{},
			expected: nil,
		},
	}
	for _, p := range patterns {
		m := publishermock.NewMockPublisher(mockController)
		a := newAutoOpsRule(t)
		h := newAutoOpsRuleCommandHandler(m, a)
		if p.expected == nil {
			m.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
		}
		cmd := &proto.ChangeDatetimeClauseCommand{Id: a.Clauses[0].Id, DatetimeClause: p.input}
		err := h.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
}

func newAutoOpsRule(t *testing.T) *domain.AutoOpsRule {
	oerc1 := &proto.OpsEventRateClause{
		GoalId:          "gid",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        proto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	oerc2 := &proto.OpsEventRateClause{
		GoalId:          "gid",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        proto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	dc1 := &proto.DatetimeClause{
		Time: 1000000001,
	}
	dc2 := &proto.DatetimeClause{
		Time: 1000000002,
	}
	aor, err := domain.NewAutoOpsRule("fid", proto.OpsType_ENABLE_FEATURE, []*proto.OpsEventRateClause{oerc1, oerc2}, []*proto.DatetimeClause{dc1, dc2})
	require.NoError(t, err)
	return aor
}

func newAutoOpsRuleCommandHandler(publisher publisher.Publisher, autoOpsRule *domain.AutoOpsRule) Handler {
	return NewAutoOpsCommandHandler(
		&eventproto.Editor{
			Email: "email",
		},
		autoOpsRule,
		publisher,
		"ns0",
	)
}
