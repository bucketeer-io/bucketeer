// Copyright 2025 The Bucketeer Authors.
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
	"google.golang.org/protobuf/types/known/wrapperspb"

	autoopsclientmock "github.com/bucketeer-io/bucketeer/pkg/autoops/client/mock"
	v2es "github.com/bucketeer-io/bucketeer/pkg/experiment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
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
		desc          string
		setup         func(*experimentService)
		orgRole       *accountproto.AccountV2_Role_Organization
		envRole       *accountproto.AccountV2_Role_Environment
		id            string
		environmentId string
		expectedErr   error
	}{
		{
			desc:          "error: ErrRequiredFieldTemplate",
			setup:         nil,
			id:            "",
			environmentId: "ns0",
			expectedErr:   createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			desc: "error: ErrNotFound",
			setup: func(s *experimentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "id-0",
			environmentId: "ns0",
			expectedErr:   createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc:          "error: ErrPermissionDenied",
			orgRole:       toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:       toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			id:            "id-1",
			environmentId: "ns0",
			expectedErr:   createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:    "success",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *experimentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.autoOpsClient.(*autoopsclientmock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(), gomock.Any(),
				).Return(&autoopsproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*autoopsproto.AutoOpsRule{},
				}, nil)
			},
			id:            "id-1",
			environmentId: "ns0",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createExperimentService(mockController, nil, p.orgRole, p.envRole)
			if p.setup != nil {
				p.setup(service)
			}
			req := &experimentproto.GetGoalRequest{Id: p.id, EnvironmentId: p.environmentId}
			_, err := service.GetGoal(ctx, req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListGoalMySQL(t *testing.T) {
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
		desc        string
		orgRole     *accountproto.AccountV2_Role_Organization
		envRole     *accountproto.AccountV2_Role_Environment
		setup       func(*experimentService)
		req         *experimentproto.ListGoalsRequest
		expectedErr error
	}{
		{
			desc:        "error: ErrPermissionDenied",
			orgRole:     toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:     toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			req:         &experimentproto.ListGoalsRequest{EnvironmentId: "ns0"},
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:    "success",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
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
				s.autoOpsClient.(*autoopsclientmock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(), gomock.Any(),
				).Return(&autoopsproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*autoopsproto.AutoOpsRule{},
				}, nil)
			},
			req:         &experimentproto.ListGoalsRequest{EnvironmentId: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createExperimentService(mockController, nil, p.orgRole, p.envRole)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ListGoals(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
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
				Command:       &experimentproto.CreateGoalCommand{Id: ""},
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Command:       &experimentproto.CreateGoalCommand{Id: "bucketeer_goal_id?"},
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusInvalidGoalID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "goal_id")),
		},
		{
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Command:       &experimentproto.CreateGoalCommand{Id: "Bucketeer-id-2019", Name: ""},
				EnvironmentId: "ns0",
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
				Command:       &experimentproto.CreateGoalCommand{Id: "Bucketeer-id-2019", Name: "name-0"},
				EnvironmentId: "ns0",
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
				Command:       &experimentproto.CreateGoalCommand{Id: "Bucketeer-id-2020", Name: "name-1"},
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.CreateGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestCreateGoalNoCommandMySQL(t *testing.T) {
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
		setup       func(s *experimentService)
		req         *experimentproto.CreateGoalRequest
		expectedErr error
	}{
		{
			desc:  "error: missing Id",
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Id:            "",
				Name:          "name-0",
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			desc:  "error: invalid Id",
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Id:            "bucketeer_goal_id?",
				Name:          "name-0",
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusInvalidGoalID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "goal_id")),
		},
		{
			desc:  "error: missing Name",
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				EnvironmentId: "ns0",
				Id:            "Bucketeer-id-2019",
				Name:          "",
			},
			expectedErr: createError(statusGoalNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc: "error: ErrGoalAlreadyExists",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalAlreadyExists)
			},
			req: &experimentproto.CreateGoalRequest{
				Id:            "Bucketeer-id-2019",
				Name:          "name-0",
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "success",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.CreateGoalRequest{
				Id:            "Bucketeer-id-2020",
				Name:          "name-1",
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.CreateGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err, p.desc)
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
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalNotFound)
			},
			req: &experimentproto.UpdateGoalRequest{
				Id:            "id-0",
				RenameCommand: &experimentproto.RenameGoalCommand{Name: "name-0"},
				EnvironmentId: "ns0",
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
				Id:            "id-1",
				RenameCommand: &experimentproto.RenameGoalCommand{Name: "name-1"},
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.UpdateGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestUpdateGoalNoCommandMySQL(t *testing.T) {
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
		req         *experimentproto.UpdateGoalRequest
		expectedErr error
	}{
		{
			desc:  "error: missing Id",
			setup: nil,
			req: &experimentproto.UpdateGoalRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			desc:  "error: name empty",
			setup: nil,
			req: &experimentproto.UpdateGoalRequest{
				Id:            "id-0",
				Name:          wrapperspb.String(""),
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusGoalNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc: "error: not found",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalNotFound)
			},
			req: &experimentproto.UpdateGoalRequest{
				Id:            "id-0",
				Name:          wrapperspb.String("name-0"),
				EnvironmentId: "ns0",
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
			req: &experimentproto.UpdateGoalRequest{
				Id:            "id-1",
				Name:          wrapperspb.String("name-0"),
				Description:   wrapperspb.String("description-0"),
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
		{
			desc: "success: archived goal",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.UpdateGoalRequest{
				Id:            "id-1",
				Archived:      wrapperspb.Bool(true),
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
		{
			desc: "success: delete goal",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.UpdateGoalRequest{
				Id:            "id-1",
				Deleted:       wrapperspb.Bool(true),
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
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
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			setup: nil,
			req: &experimentproto.ArchiveGoalRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
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
				Id:            "id-0",
				Command:       &experimentproto.ArchiveGoalCommand{},
				EnvironmentId: "ns0",
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
				Id:            "id-1",
				Command:       &experimentproto.ArchiveGoalCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
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
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			setup: nil,
			req: &experimentproto.DeleteGoalRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
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
				Id:            "id-0",
				Command:       &experimentproto.DeleteGoalCommand{},
				EnvironmentId: "ns0",
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
				Id:            "id-1",
				Command:       &experimentproto.DeleteGoalCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
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
	service := createExperimentService(mockController, nil, nil, nil)
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
