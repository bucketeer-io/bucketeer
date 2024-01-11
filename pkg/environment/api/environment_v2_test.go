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
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestGetEnvironmentV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
			req := &proto.GetEnvironmentV2Request{Id: p.id}
			resp, err := s.GetEnvironmentV2(ctx, req)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestListEnvironmentsV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
		input       *proto.ListEnvironmentsV2Request
		expected    *proto.ListEnvironmentsV2Response
		expectedErr error
	}{
		{
			desc:     "err: ErrInvalidCursor",
			setup:    nil,
			input:    &proto.ListEnvironmentsV2Request{Cursor: "XXX"},
			expected: nil,
			expectedErr: createError(
				statusInvalidCursor,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
			),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       &proto.ListEnvironmentsV2Request{},
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
			input:       &proto.ListEnvironmentsV2Request{PageSize: 2, Cursor: ""},
			expected:    &proto.ListEnvironmentsV2Response{Environments: []*proto.EnvironmentV2{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListEnvironmentsV2(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestCreateEnvironmentV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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

	envExpected, err := domain.NewEnvironmentV2(
		"Env Name-dev01",
		"url-code-01",
		"description",
		"project-id01",
		"organization-id01",
		nil,
	)
	envExpected.Archived = false
	require.NoError(t, err)

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.CreateEnvironmentV2Request
		expected    *proto.EnvironmentV2
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: nil,
			},
			expectedErr: createError(
				statusNoCommand,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
			),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: empty name",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: ""},
			},
			expectedErr: createError(
				statusEnvironmentNameRequired,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
			),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: only space",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "    "},
			},
			expectedErr: createError(
				statusEnvironmentNameRequired,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
			),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: max name length exceeded",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: strings.Repeat("a", 51)},
			},
			expectedErr: createError(
				statusInvalidEnvironmentName,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
			),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: empty url code",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: ""},
			},
			expectedErr: createError(
				statusInvalidEnvironmentUrlCode,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
			),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: can't use uppercase",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "URLCODE"},
			},
			expectedErr: createError(
				statusInvalidEnvironmentUrlCode,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
			),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: can't use space",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "url code"},
			},
			expectedErr: createError(
				statusInvalidEnvironmentUrlCode,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
			),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: max url code length exceeded",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: strings.Repeat("a", 51)},
			},
			expectedErr: createError(
				statusInvalidEnvironmentUrlCode,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
			),
		},
		{
			desc:  "err: ErrProjectIDRequired",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "url-code", ProjectId: ""},
			},
			expectedErr: createError(
				statusProjectIDRequired,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "project_id"),
			),
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
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "url-code", ProjectId: "project-id"},
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
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "url-code", ProjectId: "project-id"},
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
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "url-code", ProjectId: "project-id"},
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
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{
					Name:        envExpected.Name,
					UrlCode:     envExpected.UrlCode,
					Description: envExpected.Description,
					ProjectId:   envExpected.ProjectId,
				},
			},
			expected: envExpected.EnvironmentV2,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.CreateEnvironmentV2(ctx, p.req)
			if resp != nil {
				assert.True(t, len(resp.Environment.Id) > 0)
				assert.Equal(t, p.expected.Name, resp.Environment.Name)
				assert.Equal(t, p.expected.UrlCode, resp.Environment.UrlCode)
				assert.Equal(t, p.expected.Description, resp.Environment.Description)
				assert.Equal(t, p.expected.ProjectId, resp.Environment.ProjectId)
				assert.Equal(t, p.expected.Archived, resp.Environment.Archived)
				assert.True(t, resp.Environment.CreatedAt > 0)
				assert.True(t, resp.Environment.UpdatedAt > 0)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateEnvironmentV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
		req         *proto.UpdateEnvironmentV2Request
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.UpdateEnvironmentV2Request{
				Id: "id01",
			},
			expectedErr: createError(
				statusNoCommand,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
			),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: only space",
			setup: nil,
			req: &proto.UpdateEnvironmentV2Request{
				Id:                       "id01",
				RenameCommand:            &proto.RenameEnvironmentV2Command{Name: "  "},
				ChangeDescriptionCommand: &proto.ChangeDescriptionEnvironmentV2Command{Description: "desc-1"},
			},
			expectedErr: createError(
				statusEnvironmentNameRequired,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
			),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: max name length exceeded",
			setup: nil,
			req: &proto.UpdateEnvironmentV2Request{
				Id:                       "id01",
				RenameCommand:            &proto.RenameEnvironmentV2Command{Name: strings.Repeat("a", 51)},
				ChangeDescriptionCommand: &proto.ChangeDescriptionEnvironmentV2Command{Description: "desc-1"},
			},
			expectedErr: createError(
				statusInvalidEnvironmentName,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
			),
		},
		{
			desc: "err: ErrEnvironmentNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentNotFound)
			},
			req: &proto.UpdateEnvironmentV2Request{
				Id:            "id01",
				RenameCommand: &proto.RenameEnvironmentV2Command{Name: "name-0"},
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
			req: &proto.UpdateEnvironmentV2Request{
				Id:            "id02",
				RenameCommand: &proto.RenameEnvironmentV2Command{Name: "name-1"},
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
			req: &proto.UpdateEnvironmentV2Request{
				Id:                       "id01",
				RenameCommand:            &proto.RenameEnvironmentV2Command{Name: "name-1"},
				ChangeDescriptionCommand: &proto.ChangeDescriptionEnvironmentV2Command{Description: "desc-1"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdateEnvironmentV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestArchiveEnvironmentV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
		req         *proto.ArchiveEnvironmentV2Request
		expectedErr error
	}{
		{
			desc: "err: ErrEnvironmentNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentNotFound)
			},
			req: &proto.ArchiveEnvironmentV2Request{Id: "id01", Command: &proto.ArchiveEnvironmentV2Command{}},
			expectedErr: createError(
				statusEnvironmentNotFound,
				localizer.MustLocalize(locale.NotFoundError),
			),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req:         &proto.ArchiveEnvironmentV2Request{Id: "id02", Command: &proto.ArchiveEnvironmentV2Command{}},
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
			req:         &proto.ArchiveEnvironmentV2Request{Id: "id01", Command: &proto.ArchiveEnvironmentV2Command{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ArchiveEnvironmentV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUnarchiveEnvironmentV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
		req         *proto.UnarchiveEnvironmentV2Request
		expectedErr error
	}{
		{
			desc: "err: ErrEnvironmentNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentNotFound)
			},
			req: &proto.UnarchiveEnvironmentV2Request{Id: "id01", Command: &proto.UnarchiveEnvironmentV2Command{}},
			expectedErr: createError(
				statusEnvironmentNotFound,
				localizer.MustLocalize(locale.NotFoundError),
			),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req:         &proto.UnarchiveEnvironmentV2Request{Id: "id02", Command: &proto.UnarchiveEnvironmentV2Command{}},
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
			req:         &proto.UnarchiveEnvironmentV2Request{Id: "id01", Command: &proto.UnarchiveEnvironmentV2Command{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UnarchiveEnvironmentV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
