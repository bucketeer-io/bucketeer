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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	experimentdomain "github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestChangePeriod(t *testing.T) {
	now := time.Now()
	startAt := now.Unix()
	stopAt := now.Local().Add(time.Hour * 1).Unix()

	mockController := gomock.NewController(t)
	defer mockController.Finish()
	m := publishermock.NewMockPublisher(mockController)
	e := newExperiment(startAt, stopAt)
	h := newExperimentCommandHandler(t, m, e)
	patterns := []*struct {
		startAt     int64
		stopAt      int64
		expectedErr error
	}{
		{
			startAt:     startAt + 10,
			stopAt:      stopAt + 10,
			expectedErr: nil,
		},
		{
			startAt:     stopAt + 10,
			stopAt:      startAt + 10,
			expectedErr: experimentdomain.ErrExperimentStartIsAfterStop,
		},
		{
			startAt:     startAt - 100,
			stopAt:      startAt - 10,
			expectedErr: experimentdomain.ErrExperimentStopIsBeforeNow,
		},
	}
	for i, p := range patterns {
		if p.expectedErr == nil {
			m.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
		}
		cmd := &experimentproto.ChangeExperimentPeriodCommand{StartAt: p.startAt, StopAt: p.stopAt}
		err := h.Handle(context.Background(), cmd)
		des := fmt.Sprintf("index: %d", i)
		assert.Equal(t, p.expectedErr, err, des)
		if err == nil {
			assert.Equal(t, p.startAt, e.Experiment.StartAt, des)
			assert.Equal(t, p.stopAt, e.Experiment.StopAt, des)
		}
	}
}

func TestHandleRenameChangeCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	e := newExperiment(0, 0)
	h := newExperimentCommandHandler(t, publisher, e)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newName := "newGName"
	cmd := &experimentproto.ChangeExperimentNameCommand{Name: newName}
	err := h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newName, e.Name)
}

func TestHandleChangeDescriptionExperimentCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	e := newExperiment(0, 0)
	h := newExperimentCommandHandler(t, publisher, e)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newDesc := "newGDesc"
	cmd := &experimentproto.ChangeExperimentDescriptionCommand{Description: newDesc}
	err := h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newDesc, e.Description)
}

func TestHandleArchiveExperimentCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	e := newExperiment(0, 0)
	h := newExperimentCommandHandler(t, publisher, e)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &experimentproto.ArchiveExperimentCommand{}
	err := h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.True(t, e.Archived)
}

func TestHandleDeleteExperimentCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	e := newExperiment(0, 0)
	h := newExperimentCommandHandler(t, publisher, e)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &experimentproto.DeleteExperimentCommand{}
	err := h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.True(t, e.Deleted)
}

func newExperiment(startAt int64, stopAt int64) *experimentdomain.Experiment {
	return &experimentdomain.Experiment{
		Experiment: &experimentproto.Experiment{
			Id:             "experiment-id",
			GoalId:         "goal-id",
			FeatureId:      "feature-id",
			FeatureVersion: 1,
			Variations: []*featureproto.Variation{
				{
					Id:          "variation-A",
					Value:       "A",
					Name:        "Variation A",
					Description: "Thing does A",
				},
				{
					Id:          "variation-B",
					Value:       "B",
					Name:        "Variation B",
					Description: "Thing does B",
				},
			},
			StartAt:   startAt,
			StopAt:    stopAt,
			CreatedAt: time.Now().Unix(),
		},
	}
}

func newExperimentCommandHandler(t *testing.T, publisher publisher.Publisher, experiment *experimentdomain.Experiment) Handler {
	t.Helper()
	return NewExperimentCommandHandler(
		&eventproto.Editor{
			Email: "email",
		},
		experiment,
		publisher,
		"ns0",
	)
}
