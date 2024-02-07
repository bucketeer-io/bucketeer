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
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestHandleCreateEnvironmentV2Command(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	pub := publishermock.NewMockPublisher(mockController)
	env, err := domain.NewEnvironmentV2(
		"env-name",
		"env-url-code",
		"env-desc",
		"project-id",
		"organization-id",
		zap.NewNop(),
	)
	assert.NoError(t, err)
	h := newEnvironmentV2CommandHandler(t, pub, env)
	pub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.CreateEnvironmentV2Command{
		Name:        env.Name,
		UrlCode:     env.UrlCode,
		Description: env.Description,
		ProjectId:   env.ProjectId,
	}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}

func TestHandleRenameEnvironmentV2Command(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	pub := publishermock.NewMockPublisher(mockController)
	env, err := domain.NewEnvironmentV2(
		"env-name",
		"env-url-code",
		"env-desc",
		"project-id",
		"organization-id",
		zap.NewNop(),
	)
	h := newEnvironmentV2CommandHandler(t, pub, env)
	pub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newName := "new-env-name"
	cmd := &environmentproto.RenameEnvironmentV2Command{Name: newName}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newName, env.Name)
}

func TestHandleChangeDescriptionEnvironmentV2Command(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	pub := publishermock.NewMockPublisher(mockController)
	env, err := domain.NewEnvironmentV2(
		"env-name",
		"env-url-code",
		"env-desc",
		"project-id",
		"organization-id",
		zap.NewNop(),
	)
	h := newEnvironmentV2CommandHandler(t, pub, env)
	pub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newDesc := "new env desc"
	cmd := &environmentproto.ChangeDescriptionEnvironmentV2Command{Description: newDesc}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newDesc, env.Description)
}

func TestHandleArchiveAndUnarchiveEnvironmentV2Command(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	pub := publishermock.NewMockPublisher(mockController)
	env, err := domain.NewEnvironmentV2(
		"env-name",
		"env-url-code",
		"env-desc",
		"project-id",
		"organization-id",
		zap.NewNop(),
	)
	h := newEnvironmentV2CommandHandler(t, pub, env)
	pub.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(2).Return(nil)
	archiveCmd := &environmentproto.ArchiveEnvironmentV2Command{}
	err = h.Handle(context.Background(), archiveCmd)
	assert.NoError(t, err)
	assert.True(t, env.Archived)

	unArchiveCmd := &environmentproto.UnarchiveEnvironmentV2Command{}
	err = h.Handle(context.Background(), unArchiveCmd)
	assert.NoError(t, err)
	assert.False(t, env.Archived)
}

func newEnvironmentV2CommandHandler(t *testing.T, publisher publisher.Publisher, env *domain.EnvironmentV2) Handler {
	t.Helper()
	return NewEnvironmentV2CommandHandler(
		&eventproto.Editor{
			Email: "email",
		},
		env,
		publisher,
	)
}
