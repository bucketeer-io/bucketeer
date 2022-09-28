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

package api

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	acmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	v2es "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestGetProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		id          string
		expectedErr error
	}{
		"err: ErrProjectIDRequired": {
			setup:       nil,
			id:          "",
			expectedErr: localizedError(statusProjectIDRequired, locale.JaJP),
		},
		"err: ErrProjectNotFound": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "err-id-0",
			expectedErr: localizedError(statusProjectNotFound, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "err-id-1",
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "success-id-0",
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			req := &proto.GetProjectRequest{Id: p.id}
			resp, err := s.GetProject(createContextWithToken(t), req)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestListProjectsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		input       *proto.ListProjectsRequest
		expected    *proto.ListProjectsResponse
		expectedErr error
	}{
		"err: ErrInvalidCursor": {
			setup:       nil,
			input:       &proto.ListProjectsRequest{Cursor: "XXX"},
			expected:    nil,
			expectedErr: localizedError(statusInvalidCursor, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       &proto.ListProjectsRequest{},
			expected:    nil,
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *EnvironmentService) {
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
			input:       &proto.ListProjectsRequest{PageSize: 2, Cursor: ""},
			expected:    &proto.ListProjectsResponse{Projects: []*proto.Project{}, Cursor: "0", TotalCount: 0},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListProjects(createContextWithToken(t), p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestCreateProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		req         *proto.CreateProjectRequest
		expectedErr error
	}{
		"err: ErrNoCommand": {
			setup: nil,
			req: &proto.CreateProjectRequest{
				Command: nil,
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrInvalidProjectID: empty id": {
			setup: nil,
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Id: ""},
			},
			expectedErr: localizedError(statusInvalidProjectID, locale.JaJP),
		},
		"err: ErrInvalidProjectID: can't use uppercase": {
			setup: nil,
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Id: "ID-1"},
			},
			expectedErr: localizedError(statusInvalidProjectID, locale.JaJP),
		},
		"err: ErrInvalidProjectID: max id length exceeded": {
			setup: nil,
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Id: strings.Repeat("a", 51)},
			},
			expectedErr: localizedError(statusInvalidProjectID, locale.JaJP),
		},
		"err: ErrProjectAlreadyExists: duplicate id": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectAlreadyExists)
			},
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Id: "id-0"},
			},
			expectedErr: localizedError(statusProjectAlreadyExists, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Id: "id-1"},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Id: "id-2"},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t)
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateTrialProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		req         *proto.CreateTrialProjectRequest
		expectedErr error
	}{
		"err: ErrNoCommand": {
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: nil,
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrInvalidProjectID: empty id": {
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Id: ""},
			},
			expectedErr: localizedError(statusInvalidProjectID, locale.JaJP),
		},
		"err: ErrInvalidProjectID: can't use uppercase": {
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Id: "ID-1"},
			},
			expectedErr: localizedError(statusInvalidProjectID, locale.JaJP),
		},
		"err: ErrInvalidProjectID: max id length exceeded": {
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Id: strings.Repeat("a", 51)},
			},
			expectedErr: localizedError(statusInvalidProjectID, locale.JaJP),
		},
		"err: ErrInvalidProjectCreatorEmail": {
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Id: "id-0", Email: "email"},
			},
			expectedErr: localizedError(statusInvalidProjectCreatorEmail, locale.JaJP),
		},
		"err: ErrProjectAlreadyExists: trial exists": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Id: "id-0", Email: "test@example.com"},
			},
			expectedErr: localizedError(statusProjectAlreadyExists, locale.JaJP),
		},
		"err: ErrProjectAlreadyExists: duplicated id": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectAlreadyExists)
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Id: "id-0", Email: "test@example.com"},
			},
			expectedErr: localizedError(statusProjectAlreadyExists, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Id: "id-1", Email: "test@example.com"},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil).Times(4)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil).Times(4)
				s.accountClient.(*acmock.MockClient).EXPECT().GetAdminAccount(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "not found"))
				s.accountClient.(*acmock.MockClient).EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(
					&accountproto.CreateAccountResponse{}, nil).Times(3)
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Id: "id-2", Email: "test@example.com"},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t)
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateTrialProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		req         *proto.UpdateProjectRequest
		expectedErr error
	}{
		"err: ErrNoCommand": {
			setup: nil,
			req: &proto.UpdateProjectRequest{
				Id: "id-0",
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrProjectIDRequired": {
			setup: nil,
			req: &proto.UpdateProjectRequest{
				ChangeDescriptionCommand: &proto.ChangeDescriptionProjectCommand{Description: "desc"},
			},
			expectedErr: localizedError(statusProjectIDRequired, locale.JaJP),
		},
		"err: ErrProjectNotFound": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectNotFound)
			},
			req: &proto.UpdateProjectRequest{
				Id:                       "id-0",
				ChangeDescriptionCommand: &proto.ChangeDescriptionProjectCommand{Description: "desc"},
			},
			expectedErr: localizedError(statusProjectNotFound, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.UpdateProjectRequest{
				Id:                       "id-1",
				ChangeDescriptionCommand: &proto.ChangeDescriptionProjectCommand{Description: "desc"},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.UpdateProjectRequest{
				Id:                       "id-1",
				ChangeDescriptionCommand: &proto.ChangeDescriptionProjectCommand{Description: "desc"},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t)
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdateProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestEnableProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		req         *proto.EnableProjectRequest
		expectedErr error
	}{
		"err: ErrNoCommand": {
			setup: nil,
			req: &proto.EnableProjectRequest{
				Id: "id-0",
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrProjectIDRequired": {
			setup: nil,
			req: &proto.EnableProjectRequest{
				Command: &proto.EnableProjectCommand{},
			},
			expectedErr: localizedError(statusProjectIDRequired, locale.JaJP),
		},
		"err: ErrProjectNotFound": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectNotFound)
			},
			req: &proto.EnableProjectRequest{
				Id:      "id-0",
				Command: &proto.EnableProjectCommand{},
			},
			expectedErr: localizedError(statusProjectNotFound, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.EnableProjectRequest{
				Id:      "id-1",
				Command: &proto.EnableProjectCommand{},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.EnableProjectRequest{
				Id:      "id-1",
				Command: &proto.EnableProjectCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t)
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.EnableProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDisableProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		req         *proto.DisableProjectRequest
		expectedErr error
	}{
		"err: ErrNoCommand": {
			setup: nil,
			req: &proto.DisableProjectRequest{
				Id: "id-0",
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrProjectIDRequired": {
			setup: nil,
			req: &proto.DisableProjectRequest{
				Command: &proto.DisableProjectCommand{},
			},
			expectedErr: localizedError(statusProjectIDRequired, locale.JaJP),
		},
		"err: ErrProjectNotFound": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectNotFound)
			},
			req: &proto.DisableProjectRequest{
				Id:      "id-0",
				Command: &proto.DisableProjectCommand{},
			},
			expectedErr: localizedError(statusProjectNotFound, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.DisableProjectRequest{
				Id:      "id-1",
				Command: &proto.DisableProjectCommand{},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.DisableProjectRequest{
				Id:      "id-1",
				Command: &proto.DisableProjectCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t)
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DisableProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestConvertTrialProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		req         *proto.ConvertTrialProjectRequest
		expectedErr error
	}{
		"err: ErrNoCommand": {
			setup: nil,
			req: &proto.ConvertTrialProjectRequest{
				Id: "id-0",
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrProjectIDRequired": {
			setup: nil,
			req: &proto.ConvertTrialProjectRequest{
				Command: &proto.ConvertTrialProjectCommand{},
			},
			expectedErr: localizedError(statusProjectIDRequired, locale.JaJP),
		},
		"err: ErrProjectNotFound": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectNotFound)
			},
			req: &proto.ConvertTrialProjectRequest{
				Id:      "id-0",
				Command: &proto.ConvertTrialProjectCommand{},
			},
			expectedErr: localizedError(statusProjectNotFound, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.ConvertTrialProjectRequest{
				Id:      "id-1",
				Command: &proto.ConvertTrialProjectCommand{},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.ConvertTrialProjectRequest{
				Id:      "id-1",
				Command: &proto.ConvertTrialProjectCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t)
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ConvertTrialProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestProjectPermissionDeniedMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		action   func(context.Context, *EnvironmentService) error
		expected error
	}{
		"CreateProject": {
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.CreateProject(ctx, &proto.CreateProjectRequest{})
				return err
			},
			expected: localizedError(statusPermissionDenied, locale.JaJP),
		},
		"CreateTrialProject": {
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.CreateTrialProject(ctx, &proto.CreateTrialProjectRequest{})
				return err
			},
			expected: localizedError(statusPermissionDenied, locale.JaJP),
		},
		"UpdateProject": {
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.UpdateProject(ctx, &proto.UpdateProjectRequest{})
				return err
			},
			expected: localizedError(statusPermissionDenied, locale.JaJP),
		},
		"EnableProject": {
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.EnableProject(ctx, &proto.EnableProjectRequest{})
				return err
			},
			expected: localizedError(statusPermissionDenied, locale.JaJP),
		},
		"DisableProject": {
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.DisableProject(ctx, &proto.DisableProjectRequest{})
				return err
			},
			expected: localizedError(statusPermissionDenied, locale.JaJP),
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithTokenRoleUnassigned(t)
			service := newEnvironmentService(t, mockController, nil)
			actual := p.action(ctx, service)
			assert.Equal(t, p.expected, actual)
		})
	}
}
