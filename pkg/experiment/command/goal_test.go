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
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

func TestHandleRenameGoalCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	g, err := domain.NewGoal("gId", "gName", "gDesc")
	assert.NoError(t, err)

	h := newGoalCommandHandler(t, publisher, g)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newName := "newGName"
	cmd := &experimentproto.RenameGoalCommand{Name: newName}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newName, g.Name)
}

func TestHandleChangeDescriptionGoalCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	g, err := domain.NewGoal("gId", "gName", "gDesc")
	assert.NoError(t, err)

	h := newGoalCommandHandler(t, publisher, g)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newDesc := "newGDesc"
	cmd := &experimentproto.ChangeDescriptionGoalCommand{Description: newDesc}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newDesc, g.Description)
}

func TestHandleArchiveGoalCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	g, err := domain.NewGoal("gId", "gName", "gDesc")
	assert.NoError(t, err)

	h := newGoalCommandHandler(t, publisher, g)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &experimentproto.ArchiveGoalCommand{}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.True(t, g.Archived)
}

func TestHandleDeleteGoalCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	g, err := domain.NewGoal("gId", "gName", "gDesc")
	assert.NoError(t, err)

	h := newGoalCommandHandler(t, publisher, g)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &experimentproto.DeleteGoalCommand{}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.True(t, g.Deleted)
}

func newGoalCommandHandler(t *testing.T, publisher publisher.Publisher, goal *domain.Goal) Handler {
	t.Helper()
	return NewGoalCommandHandler(
		&eventproto.Editor{
			Email: "email",
		},
		goal,
		publisher,
		"ns0",
	)
}
