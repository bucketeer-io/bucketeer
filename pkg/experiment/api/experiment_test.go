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

package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	v2es "github.com/bucketeer-io/bucketeer/pkg/experiment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	storagetesting "github.com/bucketeer-io/bucketeer/pkg/storage/testing"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

func TestGetExperimentMySQL(t *testing.T) {
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
			expectedErr:          createError(statusExperimentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
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
		req := &experimentproto.GetExperimentRequest{Id: p.id, EnvironmentNamespace: p.environmentNamespace}
		_, err := service.GetExperiment(ctx, req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestListExperimentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.ListExperimentsRequest
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
			req:         &experimentproto.ListExperimentsRequest{FeatureId: "id-0", EnvironmentNamespace: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.ListExperiments(createContextWithTokenRoleUnassigned(), p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestCreateExperimentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(s *experimentService)
		input       *experimentproto.CreateExperimentRequest
		expectedErr error
	}{
		{
			setup: func(s *experimentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &experimentproto.CreateExperimentRequest{
				Command: &experimentproto.CreateExperimentCommand{
					FeatureId: "fid",
					GoalIds:   []string{"goalId"},
					StartAt:   1,
					StopAt:    10,
				},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	ctx := createContextWithToken()
	for _, p := range patterns {
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.CreateExperiment(ctx, p.input)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestValidateCreateExperimentRequest(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
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
		in       *experimentproto.CreateExperimentRequest
		expected error
	}{
		{
			in: &experimentproto.CreateExperimentRequest{
				Command: &experimentproto.CreateExperimentCommand{
					FeatureId: "fid",
					GoalIds:   []string{"gid"},
					StartAt:   1,
					StopAt:    10,
				},
				EnvironmentNamespace: "ns0",
			},
			expected: nil,
		},
		{
			in: &experimentproto.CreateExperimentRequest{
				Command: &experimentproto.CreateExperimentCommand{
					FeatureId: "",
					GoalIds:   []string{"gid"},
				},
				EnvironmentNamespace: "ns0",
			},
			expected: createError(statusFeatureIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id")),
		},
		{
			in: &experimentproto.CreateExperimentRequest{
				Command: &experimentproto.CreateExperimentCommand{
					FeatureId: "fid",
					GoalIds:   nil,
				},
				EnvironmentNamespace: "ns0",
			},
			expected: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			in: &experimentproto.CreateExperimentRequest{
				Command: &experimentproto.CreateExperimentCommand{
					FeatureId: "fid",
					GoalIds:   []string{""},
				},
				EnvironmentNamespace: "ns0",
			},
			expected: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			in: &experimentproto.CreateExperimentRequest{
				Command: &experimentproto.CreateExperimentCommand{
					FeatureId: "fid",
					GoalIds:   []string{"gid", ""},
				},
				EnvironmentNamespace: "ns0",
			},
			expected: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			in: &experimentproto.CreateExperimentRequest{
				Command: &experimentproto.CreateExperimentCommand{
					FeatureId: "fid",
					GoalIds:   []string{"gid0", "gid1"},
					StartAt:   1,
					StopAt:    30*24*60*60 + 2,
				},
				EnvironmentNamespace: "ns0",
			},
			expected: createError(statusPeriodTooLong, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "period")),
		},
		{
			in: &experimentproto.CreateExperimentRequest{
				Command: &experimentproto.CreateExperimentCommand{
					FeatureId: "fid",
					GoalIds:   []string{"gid0", "gid1"},
					StartAt:   1,
					StopAt:    10,
				},
				EnvironmentNamespace: "ns0",
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := validateCreateExperimentRequest(p.in, localizer)
		assert.Equal(t, p.expected, err)
	}
}

func TestUpdateExperimentMySQL(t *testing.T) {
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
		req         *experimentproto.UpdateExperimentRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.UpdateExperimentRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusExperimentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			setup: nil,
			req: &experimentproto.UpdateExperimentRequest{
				Id: "id-1",
				ChangeExperimentPeriodCommand: &experimentproto.ChangeExperimentPeriodCommand{
					StartAt: time.Now().Unix(),
					StopAt:  time.Now().AddDate(0, 0, 31).Unix(),
				},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusPeriodTooLong, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "period")),
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.UpdateExperimentRequest{
				Id:                   "id-0",
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
			req: &experimentproto.UpdateExperimentRequest{
				Id:                   "id-1",
				ChangeNameCommand:    &experimentproto.ChangeExperimentNameCommand{Name: "test-name"},
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
		_, err := service.UpdateExperiment(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestStartExperimentMySQL(t *testing.T) {
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
		desc        string
		setup       func(*experimentService)
		req         *experimentproto.StartExperimentRequest
		expectedErr error
	}{
		{
			desc:  "error id required",
			setup: nil,
			req: &experimentproto.StartExperimentRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusExperimentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error no command",
			setup: nil,
			req: &experimentproto.StartExperimentRequest{
				Id:                   "eid",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "error not found",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.StartExperimentRequest{
				Id:                   "noop",
				Command:              &experimentproto.StartExperimentCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.StartExperimentRequest{
				Id:                   "eid",
				Command:              &experimentproto.StartExperimentCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createExperimentService(mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.StartExperiment(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestFinishExperimentMySQL(t *testing.T) {
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
		desc        string
		setup       func(*experimentService)
		req         *experimentproto.FinishExperimentRequest
		expectedErr error
	}{
		{
			desc:  "error id required",
			setup: nil,
			req: &experimentproto.FinishExperimentRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusExperimentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error no command",
			setup: nil,
			req: &experimentproto.FinishExperimentRequest{
				Id:                   "eid",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "error not found",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.FinishExperimentRequest{
				Id:                   "noop",
				Command:              &experimentproto.FinishExperimentCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.FinishExperimentRequest{
				Id:                   "eid",
				Command:              &experimentproto.FinishExperimentCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createExperimentService(mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.FinishExperiment(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestStopExperimentMySQL(t *testing.T) {
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
		req         *experimentproto.StopExperimentRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.StopExperimentRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusExperimentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			setup: nil,
			req: &experimentproto.StopExperimentRequest{
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
				).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.StopExperimentRequest{
				Id:                   "id-0",
				Command:              &experimentproto.StopExperimentCommand{},
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
			req: &experimentproto.StopExperimentRequest{
				Id:                   "id-1",
				Command:              &experimentproto.StopExperimentCommand{},
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
		_, err := service.StopExperiment(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestArchiveExperimentMySQL(t *testing.T) {
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
		req         *experimentproto.ArchiveExperimentRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.ArchiveExperimentRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusExperimentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "experiment_id")),
		},
		{
			setup: nil,
			req: &experimentproto.ArchiveExperimentRequest{
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
				).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.ArchiveExperimentRequest{
				Id:                   "id-0",
				Command:              &experimentproto.ArchiveExperimentCommand{},
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
			req: &experimentproto.ArchiveExperimentRequest{
				Id:                   "id-1",
				Command:              &experimentproto.ArchiveExperimentCommand{},
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
		_, err := service.ArchiveExperiment(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestDeleteExperimentMySQL(t *testing.T) {
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
		req         *experimentproto.DeleteExperimentRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.DeleteExperimentRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusExperimentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			setup: nil,
			req: &experimentproto.DeleteExperimentRequest{
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
				).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.DeleteExperimentRequest{
				Id:                   "id-0",
				Command:              &experimentproto.DeleteExperimentCommand{},
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
			req: &experimentproto.DeleteExperimentRequest{
				Id:                   "id-1",
				Command:              &experimentproto.DeleteExperimentCommand{},
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
		_, err := service.DeleteExperiment(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestExperimentPermissionDenied(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleUnassigned()
	s := storagetesting.NewInMemoryStorage()
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
			desc: "CreateExperiment",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.CreateExperiment(ctx, &experimentproto.CreateExperimentRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "UpdateExperiment",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "StopExperiment",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.StopExperiment(ctx, &experimentproto.StopExperimentRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "DeleteExperiment",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.DeleteExperiment(ctx, &experimentproto.DeleteExperimentRequest{})
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
