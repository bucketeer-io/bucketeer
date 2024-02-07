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

	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestHandleCreateProjectCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	project, err := domain.NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		false)
	assert.NoError(t, err)

	h := newProjectCommandHandler(t, publisher, project)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.CreateProjectCommand{Name: project.Name, Description: project.Description}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}

func TestHandleCreateTrialProjectCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	project, err := domain.NewProject(
		"project-name",
		"project-code",
		"",
		"test@example.com",
		"organization-id",
		true)
	assert.NoError(t, err)

	h := newProjectCommandHandler(t, publisher, project)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.CreateTrialProjectCommand{Name: project.Name, Email: project.CreatorEmail}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}

func TestHandleChangeDescriptionProjectCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	project, err := domain.NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		false)
	assert.NoError(t, err)

	h := newProjectCommandHandler(t, publisher, project)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newDesc := "new project desc"
	cmd := &environmentproto.ChangeDescriptionProjectCommand{Description: newDesc}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newDesc, project.Description)
}

func TestHandleChangeNameProjectCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	project, err := domain.NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		false,
	)
	assert.NoError(t, err)

	h := newProjectCommandHandler(t, publisher, project)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newName := "new-project-name"
	cmd := &environmentproto.RenameProjectCommand{Name: newName}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newName, project.Name)
}

func TestHandleEnableProjectCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	project, err := domain.NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		false,
	)
	assert.NoError(t, err)
	project.Disabled = true

	h := newProjectCommandHandler(t, publisher, project)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.EnableProjectCommand{}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.False(t, project.Disabled)
}

func TestHandleDisableProjectCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	project, err := domain.NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		false,
	)
	assert.NoError(t, err)

	h := newProjectCommandHandler(t, publisher, project)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.DisableProjectCommand{}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.True(t, project.Disabled)
}

func TestHandleConvertTrialProjectCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	project, err := domain.NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		true,
	)
	assert.NoError(t, err)

	h := newProjectCommandHandler(t, publisher, project)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.ConvertTrialProjectCommand{}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.False(t, project.Trial)
}

func newProjectCommandHandler(t *testing.T, publisher publisher.Publisher, project *domain.Project) Handler {
	t.Helper()
	return NewProjectCommandHandler(
		&eventproto.Editor{
			Email: "email",
		},
		project,
		publisher,
	)
}
