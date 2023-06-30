// Copyright 2022 The Bucketeer Authors.
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
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestHandleCreateEnvironmentCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	env := domain.NewEnvironment("env-id", "env desc", "project-id")

	h := newEnvironmentCommandHandler(t, publisher, env)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.CreateEnvironmentCommand{Id: env.Id, Description: env.Description}
	err := h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}

func TestHandleRenameEnvironmentCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	env := domain.NewEnvironment("env-id", "env desc", "project-id")

	h := newEnvironmentCommandHandler(t, publisher, env)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newName := "new-env-name"
	cmd := &environmentproto.RenameEnvironmentCommand{Name: newName}
	err := h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newName, env.Name)
}

func TestHandleChangeDescriptionEnvironmentCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	env := domain.NewEnvironment("env-id", "env desc", "project-id")

	h := newEnvironmentCommandHandler(t, publisher, env)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newDesc := "new env desc"
	cmd := &environmentproto.ChangeDescriptionEnvironmentCommand{Description: newDesc}
	err := h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newDesc, env.Description)
}

func TestHandleDeleteEnvironmentCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	env := domain.NewEnvironment("env-id", "env desc", "project-id")

	h := newEnvironmentCommandHandler(t, publisher, env)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.DeleteEnvironmentCommand{}
	err := h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.True(t, env.Deleted)
}

func newEnvironmentCommandHandler(t *testing.T, publisher publisher.Publisher, env *domain.Environment) Handler {
	t.Helper()
	return NewEnvironmentCommandHandler(
		&eventproto.Editor{
			Email: "email",
			Role:  accountproto.Account_EDITOR,
		},
		env,
		publisher,
	)
}
