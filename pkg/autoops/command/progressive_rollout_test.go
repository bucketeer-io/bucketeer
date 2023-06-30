// Copyright 2023 The Bucketeer Authors.
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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestProgressiveRolloutDelete(t *testing.T) {
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
		a := createProgressiveRollout(t)
		h := newProgressiveRolloutCommandHandler(m, a)
		if p.expected == nil {
			m.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
		}
		cmd := &autoopsproto.DeleteProgressiveRolloutCommand{}
		err := h.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
}

func createProgressiveRollout(t *testing.T) *domain.ProgressiveRollout {
	p, err := domain.NewProgressiveRollout(
		"feature-id",
		nil,
		&autoopsproto.ProgressiveRolloutTemplateScheduleClause{
			Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
				{
					ScheduleId: "schedule-id-0",
					ExecuteAt:  time.Now().Unix(),
					Weight:     0,
				},
				{
					ScheduleId: "schedule-id-1",
					ExecuteAt:  time.Now().AddDate(1, 0, 0).Unix(),
					Weight:     20,
				},
				{
					ScheduleId: "schedule-id-2",
					ExecuteAt:  time.Now().AddDate(2, 0, 0).Unix(),
					Weight:     40,
				},
				{
					ScheduleId: "schedule-id-3",
					ExecuteAt:  time.Now().AddDate(3, 0, 0).Unix(),
					Weight:     60,
				},
				{
					ScheduleId: "schedule-id-4",
					ExecuteAt:  time.Now().AddDate(4, 0, 0).Unix(),
					Weight:     80,
				},
				{
					ScheduleId: "schedule-id-5",
					ExecuteAt:  time.Now().AddDate(5, 0, 0).Unix(),
					Weight:     100,
				},
			},
			Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
			Increments:  20,
			VariationId: "vid-1",
		},
	)
	require.NoError(t, err)
	return p
}

func newProgressiveRolloutCommandHandler(publisher publisher.Publisher, progressiveRollout *domain.ProgressiveRollout) Handler {
	return NewProgressiveRolloutCommandHandler(
		&eventproto.Editor{
			Email: "email",
			Role:  accountproto.Account_EDITOR,
		},
		progressiveRollout,
		publisher,
		"ns0",
	)
}
