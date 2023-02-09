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
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

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
		desc        string
		setup       func(*EnvironmentService)
		id          string
		expectedErr error
	}{
		{
			desc:        "err: ErrEnvironmentIDRequired",
			setup:       nil,
			id:          "",
			expectedErr: createError(statusEnvironmentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrEnvironmentNotFound",
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "id-0",
			expectedErr: createError(statusEnvironmentNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "id-1",
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
		desc        string
		setup       func(*EnvironmentService)
		namespace   string
		expectedErr error
	}{
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			namespace:   "ns-0",
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "err: ErrEnvironmentNotFound",
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			namespace:   "ns-1",
			expectedErr: createError(statusEnvironmentNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
		desc        string
		setup       func(*EnvironmentService)
		input       *proto.ListEnvironmentsRequest
		expected    *proto.ListEnvironmentsResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			setup:       nil,
			input:       &proto.ListEnvironmentsRequest{Cursor: "XXX"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       &proto.ListEnvironmentsRequest{},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.CreateEnvironmentRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.CreateEnvironmentRequest{
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrInvalidEnvironmentID: empty id",
			setup: nil,
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: ""},
			},
			expectedErr: createError(statusInvalidEnvironmentID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "id")),
		},
		{
			desc:  "err: ErrInvalidEnvironmentID: can't use uppercase",
			setup: nil,
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: "NS-1"},
			},
			expectedErr: createError(statusInvalidEnvironmentID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "id")),
		},
		{
			desc:  "err: ErrInvalidEnvironmentID: max id length exceeded",
			setup: nil,
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: strings.Repeat("a", 51)},
			},
			expectedErr: createError(statusInvalidEnvironmentID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "id")),
		},
		{
			desc:  "err: ErrProjectIDRequired",
			setup: nil,
			req: &proto.CreateEnvironmentRequest{
				Command: &proto.CreateEnvironmentCommand{Id: "ns-0", ProjectId: ""},
			},
			expectedErr: createError(statusProjectIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "project_id")),
		},
		{
			desc: "err: ErrProjectNotFound",
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
			expectedErr: createError(statusProjectNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrEnvironmentAlreadyExists",
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
			expectedErr: createError(statusEnvironmentAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "err: ErrInternal",
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
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.UpdateEnvironmentRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.UpdateEnvironmentRequest{
				Id: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrEnvironmentIDRequired",
			setup: nil,
			req: &proto.UpdateEnvironmentRequest{
				RenameCommand: &proto.RenameEnvironmentCommand{Name: "name-0"},
			},
			expectedErr: createError(statusEnvironmentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrEnvironmentNotFound",
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
			expectedErr: createError(statusEnvironmentNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
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
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.DeleteEnvironmentRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrNoCommand",
			setup:       nil,
			req:         &proto.DeleteEnvironmentRequest{},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrEnvironmentIDRequired",
			setup: nil,
			req: &proto.DeleteEnvironmentRequest{
				Command: &proto.DeleteEnvironmentCommand{},
			},
			expectedErr: createError(statusEnvironmentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrEnvironmentNotFound",
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
			expectedErr: createError(statusEnvironmentNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
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
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
		desc     string
		action   func(context.Context, *EnvironmentService) error
		expected error
	}{
		{
			desc: "CreateEnvironment",
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.CreateEnvironment(ctx, &proto.CreateEnvironmentRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "UpdateEnvironment",
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.UpdateEnvironment(ctx, &proto.UpdateEnvironmentRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "DeleteEnvironment",
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.DeleteEnvironment(ctx, &proto.DeleteEnvironmentRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := createContextWithTokenRoleUnassigned(t)
			service := newEnvironmentService(t, mockController, nil)
			actual := p.action(ctx, service)
			assert.Equal(t, p.expected, actual)
		})
	}
}
