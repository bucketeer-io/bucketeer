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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	v2ps "github.com/bucketeer-io/bucketeer/pkg/push/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	storagetesting "github.com/bucketeer-io/bucketeer/pkg/storage/testing"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
)

func TestNewPushService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mysqlClient := mysqlmock.NewMockClient(mockController)
	featureClientMock := featureclientmock.NewMockClient(mockController)
	experimentClientMock := experimentclientmock.NewMockClient(mockController)
	accountClientMock := accountclientmock.NewMockClient(mockController)
	pm := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewPushService(
		mysqlClient,
		featureClientMock,
		experimentClientMock,
		accountClientMock,
		pm,
		WithLogger(logger),
	)
	assert.IsType(t, &PushService{}, s)
}

func TestCreatePushMySQL(t *testing.T) {
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
		setup       func(*PushService)
		req         *pushproto.CreatePushRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrFCMAPIKeyRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				Command: &pushproto.CreatePushCommand{},
			},
			expectedErr: createError(statusFCMAPIKeyRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "fcm_api_key")),
		},
		{
			desc:  "err: ErrTagsRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				Command: &pushproto.CreatePushCommand{
					FcmApiKey: "key-0",
				},
			},
			expectedErr: createError(statusTagsRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag")),
		},
		{
			desc:  "err: ErrNameRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				Command: &pushproto.CreatePushCommand{
					FcmApiKey: "key-1",
					Tags:      []string{"tag-0"},
				},
			},
			expectedErr: createError(statusNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc: "err: ErrAlreadyExists",
			setup: func(s *PushService) {
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2ps.ErrPushAlreadyExists)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentNamespace: "ns0",
				Command: &pushproto.CreatePushCommand{
					FcmApiKey: "key-0",
					Tags:      []string{"tag-0"},
					Name:      "name-1",
				},
			},
			expectedErr: createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "success",
			setup: func(s *PushService) {
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentNamespace: "ns0",
				Command: &pushproto.CreatePushCommand{
					FcmApiKey: "key-1",
					Tags:      []string{"tag-0"},
					Name:      "name-1",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreatePush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdatePushMySQL(t *testing.T) {
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
		setup       func(*PushService)
		req         *pushproto.UpdatePushRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &pushproto.UpdatePushRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			req: &pushproto.UpdatePushRequest{
				Id: "key-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: ErrTagsRequired: delete",
			req: &pushproto.UpdatePushRequest{
				Id:                    "key-0",
				DeletePushTagsCommand: &pushproto.DeletePushTagsCommand{Tags: []string{}},
			},
			expectedErr: createError(statusTagsRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag")),
		},
		{
			desc: "err: ErrTagsRequired: add",
			req: &pushproto.UpdatePushRequest{
				Id:                 "key-0",
				AddPushTagsCommand: &pushproto.AddPushTagsCommand{Tags: []string{}},
			},
			expectedErr: createError(statusTagsRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag")),
		},
		{
			desc: "err: ErrNameRequired: add",
			req: &pushproto.UpdatePushRequest{
				Id:                "key-0",
				RenamePushCommand: &pushproto.RenamePushCommand{Name: ""},
			},
			expectedErr: createError(statusNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *PushService) {
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2ps.ErrPushNotFound)
			},
			req: &pushproto.UpdatePushRequest{
				Id:                 "key-1",
				AddPushTagsCommand: &pushproto.AddPushTagsCommand{Tags: []string{"tag-1"}},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success: rename",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentNamespace: "ns0",
				Id:                   "key-0",
				RenamePushCommand:    &pushproto.RenamePushCommand{Name: "name-1"},
			},
			expectedErr: nil,
		},
		{
			desc: "success: deletePushTags",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentNamespace:  "ns0",
				Id:                    "key-0",
				DeletePushTagsCommand: &pushproto.DeletePushTagsCommand{Tags: []string{"tag-0"}},
			},
			expectedErr: nil,
		},
		{
			desc: "success: addPushTags",
			setup: func(s *PushService) {
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentNamespace: "ns0",
				Id:                   "key-0",
				AddPushTagsCommand:   &pushproto.AddPushTagsCommand{Tags: []string{"tag-2"}},
			},
			expectedErr: nil,
		},
		{
			desc: "success",
			setup: func(s *PushService) {
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentNamespace:  "ns0",
				Id:                    "key-0",
				AddPushTagsCommand:    &pushproto.AddPushTagsCommand{Tags: []string{"tag-2"}},
				DeletePushTagsCommand: &pushproto.DeletePushTagsCommand{Tags: []string{"tag-0"}},
				RenamePushCommand:     &pushproto.RenamePushCommand{Name: "name-1"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdatePush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeletePushMySQL(t *testing.T) {
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
		setup       func(*PushService)
		req         *pushproto.DeletePushRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &pushproto.DeletePushRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			req: &pushproto.DeletePushRequest{
				Id: "key-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2ps.ErrPushNotFound)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentNamespace: "ns0",
				Id:                   "key-1",
				Command:              &pushproto.DeletePushCommand{},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentNamespace: "ns0",
				Id:                   "key-0",
				Command:              &pushproto.DeletePushCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DeletePush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListPushesMySQL(t *testing.T) {
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
		setup       func(*PushService)
		input       *pushproto.ListPushesRequest
		expected    *pushproto.ListPushesResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			setup:       nil,
			input:       &pushproto.ListPushesRequest{Cursor: "XXX"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       &pushproto.ListPushesRequest{},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *PushService) {
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
			input:       &pushproto.ListPushesRequest{PageSize: 2, Cursor: ""},
			expected:    &pushproto.ListPushesResponse{Pushes: []*pushproto.Push{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newPushServiceWithMock(t, mockController, storagetesting.NewInMemoryStorage())
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListPushes(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newPushServiceWithMock(t *testing.T, c *gomock.Controller, s storage.Client) *PushService {
	t.Helper()
	return &PushService{
		mysqlClient:      mysqlmock.NewMockClient(c),
		featureClient:    featureclientmock.NewMockClient(c),
		experimentClient: experimentclientmock.NewMockClient(c),
		accountClient:    accountclientmock.NewMockClient(c),
		publisher:        publishermock.NewMockPublisher(c),
		logger:           zap.NewNop(),
	}
}

func createContextWithToken(t *testing.T) context.Context {
	t.Helper()
	token := &token.IDToken{
		Issuer:    "issuer",
		Subject:   "sub",
		Audience:  "audience",
		Expiry:    time.Now().AddDate(100, 0, 0),
		IssuedAt:  time.Now(),
		Email:     "email",
		AdminRole: accountproto.Account_OWNER,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}
