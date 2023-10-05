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

package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	v2es "github.com/bucketeer-io/bucketeer/pkg/experiment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	storeclient "github.com/bucketeer-io/bucketeer/pkg/storage/testing"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

func TestGetGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleUnassigned()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		setup                func(*experimentService)
		id                   string
		environmentNamespace string
		expectedErr          error
	}{
		{
			setup:                nil,
			id:                   "",
			environmentNamespace: "ns0",
			expectedErr:          createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			setup: func(s *experimentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:                   "id-0",
			environmentNamespace: "ns0",
			expectedErr:          createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			setup: func(s *experimentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:                   "id-1",
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		req := &experimentproto.GetGoalRequest{Id: p.id, EnvironmentNamespace: p.environmentNamespace}
		_, err := service.GetGoal(ctx, req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestListGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.ListGoalsRequest
		expectedErr error
	}{
		{
			setup: func(s *experimentService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req:         &experimentproto.ListGoalsRequest{EnvironmentNamespace: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.ListGoals(createContextWithTokenRoleUnassigned(), p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestCreateGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		setup       func(s *experimentService)
		req         *experimentproto.CreateGoalRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Command:              nil,
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Command:              &experimentproto.CreateGoalCommand{Id: ""},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Command:              &experimentproto.CreateGoalCommand{Id: "bucketeer_goal_id?"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInvalidGoalID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "goal_id")),
		},
		{
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Command:              &experimentproto.CreateGoalCommand{Id: "Bucketeer-id-2019", Name: ""},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusGoalNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalAlreadyExists)
			},
			req: &experimentproto.CreateGoalRequest{
				Command:              &experimentproto.CreateGoalCommand{Id: "Bucketeer-id-2019", Name: "name-0"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.CreateGoalRequest{
				Command:              &experimentproto.CreateGoalCommand{Id: "Bucketeer-id-2020", Name: "name-1"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.CreateGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestUpdateGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.UpdateGoalRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.UpdateGoalRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			setup: nil,
			req: &experimentproto.UpdateGoalRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalNotFound)
			},
			req: &experimentproto.UpdateGoalRequest{
				Id:                   "id-0",
				RenameCommand:        &experimentproto.RenameGoalCommand{Name: "name-0"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.UpdateGoalRequest{
				Id:                   "id-1",
				RenameCommand:        &experimentproto.RenameGoalCommand{Name: "name-1"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.UpdateGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestArchiveGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.ArchiveGoalRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.ArchiveGoalRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			setup: nil,
			req: &experimentproto.ArchiveGoalRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalNotFound)
			},
			req: &experimentproto.ArchiveGoalRequest{
				Id:                   "id-0",
				Command:              &experimentproto.ArchiveGoalCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.ArchiveGoalRequest{
				Id:                   "id-1",
				Command:              &experimentproto.ArchiveGoalCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.ArchiveGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestDeleteGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.DeleteGoalRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.DeleteGoalRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			setup: nil,
			req: &experimentproto.DeleteGoalRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalNotFound)
			},
			req: &experimentproto.DeleteGoalRequest{
				Id:                   "id-0",
				Command:              &experimentproto.DeleteGoalCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.DeleteGoalRequest{
				Id:                   "id-1",
				Command:              &experimentproto.DeleteGoalCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.DeleteGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestGoalPermissionDenied(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleUnassigned()
	s := storeclient.NewInMemoryStorage()
	service := createExperimentService(mockController, s)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []struct {
		desc     string
		action   func(context.Context, *experimentService) error
		expected error
	}{
		{
			desc: "CreateGoal",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.CreateGoal(ctx, &experimentproto.CreateGoalRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "UpdateGoal",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.UpdateGoal(ctx, &experimentproto.UpdateGoalRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "DeleteGoal",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.DeleteGoal(ctx, &experimentproto.DeleteGoalRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
	}
	for _, p := range patterns {
		actual := p.action(ctx, service)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
	}
}
