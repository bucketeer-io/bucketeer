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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

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
	h, err := NewExperimentCommandHandler(
		&eventproto.Editor{
			Email: "email",
		},
		experiment,
		publisher,
		"ns0",
	)
	if err != nil {
		t.Fatal(err)
	}
	return h
}
