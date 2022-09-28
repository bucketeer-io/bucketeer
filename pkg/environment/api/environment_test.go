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

	v2es "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestGetEnvironmentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		id          string
		expectedErr error
	}{
		"err: ErrEnvironmentIDRequired": {
			setup:       nil,
			id:          "",
			expectedErr: localizedError(statusEnvironmentIDRequired, locale.JaJP),
		},
		"err: ErrEnvironmentNotFound": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "id-0",
			expectedErr: localizedError(statusEnvironmentNotFound, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "id-1",
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
			id:          "id-3",
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			req := &proto.GetEnvironmentRequest{Id: p.id}
			resp, err := s.GetEnvironment(createContextWithToken(t), req)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestGetEnvironmentByNamespaceMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		namespace   string
		expectedErr error
	}{
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			namespace:   "ns-0",
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"err: ErrEnvironmentNotFound": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			namespace:   "ns-1",
			expectedErr: localizedError(statusEnvironmentNotFound, locale.JaJP),
		},
		"success": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			namespace:   "ns-2",
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			req := &proto.GetEnvironmentByNamespaceRequest{Namespace: p.namespace}
			resp, err := s.GetEnvironmentByNamespace(createContextWithToken(t), req)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestListEnvironmentsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		input       *proto.ListEnvironmentsRequest
		expected    *proto.ListEnvironmentsResponse
		expectedErr error
	}{
		"err: ErrInvalidCursor": {
			setup:       nil,
			input:       &proto.ListEnvironmentsRequest{Cursor: "XXX"},
			expected:    nil,
			expectedErr: localizedError(statusInvalidCursor, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       &proto.ListEnvironmentsRequest{},
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
			input:       &proto.ListEnvironmentsRequest{PageSize: 2, Cursor: ""},
			expected:    &proto.ListEnvironmentsResponse{Environments: []*proto.Environment{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListEnvironments(createContextWithToken(t), p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestCreateEnvironmentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		req         *proto.CreateEnvironmentRequest
		expectedErr error
	}{
		"err: ErrNoCommand": {
			setup: nil,
			req: &proto.CreateEnvironmentRequest{
				Command: nil,
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrInvalidEnvironmentID: empty id": {
			setup: nil,
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: ""},
			},
			expectedErr: localizedError(statusInvalidEnvironmentID, locale.JaJP),
		},
		"err: ErrInvalidEnvironmentID: can't use uppercase": {
			setup: nil,
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: "NS-1"},
			},
			expectedErr: localizedError(statusInvalidEnvironmentID, locale.JaJP),
		},
		"err: ErrInvalidEnvironmentID: max id length exceeded": {
			setup: nil,
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: strings.Repeat("a", 51)},
			},
			expectedErr: localizedError(statusInvalidEnvironmentID, locale.JaJP),
		},
		"err: ErrProjectIDRequired": {
			setup: nil,
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: "ns-0", ProjectId: ""},
			},
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
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: "ns-0", ProjectId: "project-id-0"},
			},
			expectedErr: localizedError(statusProjectNotFound, locale.JaJP),
		},
		"err: ErrEnvironmentAlreadyExists": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentAlreadyExists)
			},
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: "ns-0", ProjectId: "project-id-0"},
			},
			expectedErr: localizedError(statusEnvironmentAlreadyExists, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: "ns-1", ProjectId: "project-id-0"},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *EnvironmentService) {
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
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: "ns-2", ProjectId: "project-id-0"},
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
			_, err := service.CreateEnvironment(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateEnvironmentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		req         *proto.UpdateEnvironmentRequest
		expectedErr error
	}{
		"err: ErrNoCommand": {
			setup: nil,
			req: &proto.UpdateEnvironmentRequest{
				Id: "ns0",
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrEnvironmentIDRequired": {
			setup: nil,
			req: &proto.UpdateEnvironmentRequest{
				RenameCommand: &proto.RenameEnvironmentCommand{Name: "name-0"},
			},
			expectedErr: localizedError(statusEnvironmentIDRequired, locale.JaJP),
		},
		"err: ErrEnvironmentNotFound": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentNotFound)
			},
			req: &proto.UpdateEnvironmentRequest{
				Id:            "ns0",
				RenameCommand: &proto.RenameEnvironmentCommand{Name: "name-0"},
			},
			expectedErr: localizedError(statusEnvironmentNotFound, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.UpdateEnvironmentRequest{
				Id:            "ns1",
				RenameCommand: &proto.RenameEnvironmentCommand{Name: "name-1"},
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
			req: &proto.UpdateEnvironmentRequest{
				Id:                       "ns1",
				RenameCommand:            &proto.RenameEnvironmentCommand{Name: "name-1"},
				ChangeDescriptionCommand: &proto.ChangeDescriptionEnvironmentCommand{Description: "desc-1"},
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
			_, err := service.UpdateEnvironment(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteEnvironmentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*EnvironmentService)
		req         *proto.DeleteEnvironmentRequest
		expectedErr error
	}{
		"err: ErrNoCommand": {
			setup:       nil,
			req:         &proto.DeleteEnvironmentRequest{},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrEnvironmentIDRequired": {
			setup: nil,
			req: &proto.DeleteEnvironmentRequest{
				Command: &proto.DeleteEnvironmentCommand{},
			},
			expectedErr: localizedError(statusEnvironmentIDRequired, locale.JaJP),
		},
		"err: ErrEnvironmentNotFound": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentNotFound)
			},
			req: &proto.DeleteEnvironmentRequest{
				Id:      "ns0",
				Command: &proto.DeleteEnvironmentCommand{},
			},
			expectedErr: localizedError(statusEnvironmentNotFound, locale.JaJP),
		},
		"err: ErrInternal": {
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.DeleteEnvironmentRequest{
				Id:      "ns1",
				Command: &proto.DeleteEnvironmentCommand{},
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
			req: &proto.DeleteEnvironmentRequest{
				Id:      "ns1",
				Command: &proto.DeleteEnvironmentCommand{},
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
			_, err := service.DeleteEnvironment(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestEnvironmentPermissionDeniedMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		action   func(context.Context, *EnvironmentService) error
		expected error
	}{
		"CreateEnvironment": {
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.CreateEnvironment(ctx, &proto.CreateEnvironmentRequest{})
				return err
			},
			expected: localizedError(statusPermissionDenied, locale.JaJP),
		},
		"UpdateEnvironment": {
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.UpdateEnvironment(ctx, &proto.UpdateEnvironmentRequest{})
				return err
			},
			expected: localizedError(statusPermissionDenied, locale.JaJP),
		},
		"DeleteEnvironment": {
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.DeleteEnvironment(ctx, &proto.DeleteEnvironmentRequest{})
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
