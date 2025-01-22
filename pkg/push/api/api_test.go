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
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountproto "github.com/bucketeer-io/bucketeer/proto/account"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	v2ps "github.com/bucketeer-io/bucketeer/pkg/push/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
)

var fcmServiceAccountDummy = []byte(`
	{
		"type": "service_account",
		"project_id": "test",
		"private_key_id": "private-key-id",
		"private_key": "-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----\n",
		"client_email": "fcm-service-account@test.iam.gserviceaccount.com",
		"client_id": "client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/fcm-service-account@test.iam.gserviceaccount.com",
		"universe_domain": "googleapis.com"
	}
`)

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

	ctx := createContextWithToken(t, true)
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
		// command is deprecating
		//{
		//	desc:  "err: ErrNoCommand",
		//	setup: nil,
		//	req: &pushproto.CreatePushRequest{
		//		Command: nil,
		//	},
		//	expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		//},
		{
			desc:  "err: ErrFCMServiceAccountRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				Command: &pushproto.CreatePushCommand{},
			},
			expectedErr: createError(statusFCMServiceAccountRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "fcm_service_account")),
		},
		{
			desc:  "err: ErrTagsRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				Command: &pushproto.CreatePushCommand{
					FcmServiceAccount: fcmServiceAccountDummy,
				},
			},
			expectedErr: createError(statusTagsRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag")),
		},
		{
			desc:  "err: ErrNameRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				Command: &pushproto.CreatePushCommand{
					FcmServiceAccount: fcmServiceAccountDummy,
					Tags:              []string{"tag-0"},
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
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.All(),
				).Return(qe).AnyTimes()
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2ps.ErrPushAlreadyExists)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentId: "ns0",
				Command: &pushproto.CreatePushCommand{
					FcmServiceAccount: fcmServiceAccountDummy,
					Tags:              []string{"tag-0"},
					Name:              "name-1",
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
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.All(),
				).Return(qe).AnyTimes()
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentId: "ns0",
				Command: &pushproto.CreatePushCommand{
					FcmServiceAccount: fcmServiceAccountDummy,
					Tags:              []string{"tag-0"},
					Name:              "name-1",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreatePush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreatePushNoCommandMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)
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
			desc:  "err: ErrFCMServiceAccountRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				FcmServiceAccount: nil,
			},
			expectedErr: createError(statusFCMServiceAccountRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "fcm_service_account")),
		},
		{
			desc:  "err: ErrTagsRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				FcmServiceAccount: fcmServiceAccountDummy,
			},
			expectedErr: createError(statusTagsRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag")),
		},
		{
			desc:  "err: ErrNameRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				FcmServiceAccount: fcmServiceAccountDummy,
				Tags:              []string{"tag-0"},
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
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.All(),
				).Return(qe).AnyTimes()
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2ps.ErrPushAlreadyExists)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentId:     "ns0",
				FcmServiceAccount: fcmServiceAccountDummy,
				Tags:              []string{"tag-0"},
				Name:              "name-1",
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
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.All(),
				).Return(qe).AnyTimes()
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentId:     "ns0",
				FcmServiceAccount: fcmServiceAccountDummy,
				Tags:              []string{"tag-0"},
				Name:              "name-1",
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
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

	ctx := createContextWithToken(t, true)
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
			desc: "err: ErrIDRequired",
			req: &pushproto.UpdatePushRequest{
				RenamePushCommand: &pushproto.RenamePushCommand{},
			},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		// command is deprecating
		//{
		//	desc: "err: ErrNoCommand",
		//	req: &pushproto.UpdatePushRequest{
		//		Id: "key-0",
		//	},
		//	expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		//},
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
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.All(),
				).Return(qe).AnyTimes()
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentId:     "ns0",
				Id:                "key-0",
				RenamePushCommand: &pushproto.RenamePushCommand{Name: "name-1"},
			},
			expectedErr: nil,
		},
		{
			desc: "success: deletePushTags",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentId:         "ns0",
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
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.All(),
				).Return(qe).AnyTimes()
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentId:      "ns0",
				Id:                 "key-0",
				AddPushTagsCommand: &pushproto.AddPushTagsCommand{Tags: []string{"tag-2"}},
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
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.All(),
				).Return(qe).AnyTimes()
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentId:         "ns0",
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
			service := newPushServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdatePush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdatePushNoCommandMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)
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
			desc: "err: ErrNotFound",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2ps.ErrPushNotFound)
			},
			req: &pushproto.UpdatePushRequest{
				Id:   "key-1",
				Name: wrapperspb.String("push-0"),
				Tags: []string{"tag-0"},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success update name",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				Id:   "key-0",
				Name: wrapperspb.String("push-0"),
			},
			expectedErr: nil,
		},
		{
			desc: "success update tags",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				Id:   "key-0",
				Tags: []string{"tag-0"},
			},
			expectedErr: nil,
		},
		{
			desc: "success",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
				Name:          wrapperspb.String("name-1"),
				Tags:          []string{"tag-0"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdatePush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCheckFCMServiceAccount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)
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
		desc              string
		fcmServiceAccount []byte
		pushes            []*pushproto.Push
		expected          error
	}{
		{
			desc:              "err: invalid service account",
			fcmServiceAccount: []byte(`"key":"value"`),
			pushes:            nil,
			expected: createError(
				statusFCMServiceAccountInvalid,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "FCM service account"),
			),
		},
		{
			desc:              "err: internal error",
			fcmServiceAccount: fcmServiceAccountDummy,
			pushes: []*pushproto.Push{
				{
					FcmServiceAccount: "`{\"key\":\"value\"}`",
				},
			},
			expected: createError(
				statusInternal,
				localizer.MustLocalize(locale.InternalServerError),
			),
		},
		{
			desc:              "err: service account already exists",
			fcmServiceAccount: fcmServiceAccountDummy,
			pushes: []*pushproto.Push{
				{
					FcmServiceAccount: string(fcmServiceAccountDummy),
				},
			},
			expected: createError(
				statusFCMServiceAccountAlreadyExists,
				localizer.MustLocalizeWithTemplate(locale.AlreadyExistsError, "FCM service account"),
			),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
			err := service.checkFCMServiceAccount(
				ctx,
				p.pushes,
				p.fcmServiceAccount,
				localizer,
			)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestDeletePushMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)
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
		// command is deprecating
		//{
		//	desc: "err: ErrNoCommand",
		//	req: &pushproto.DeletePushRequest{
		//		Id: "key-0",
		//	},
		//	expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		//},
		{
			desc: "err: ErrNotFound",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2ps.ErrPushNotFound)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
				Command:       &pushproto.DeletePushCommand{},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *PushService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
				Command:       &pushproto.DeletePushCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
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

	ctx := createContextWithToken(t, false)
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
		setup       func(*PushService)
		input       *pushproto.ListPushesRequest
		expected    *pushproto.ListPushesResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			orgRole:     toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:     toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup:       nil,
			input:       &pushproto.ListPushesRequest{Cursor: "XXX", EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *PushService) {
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe).AnyTimes()
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       &pushproto.ListPushesRequest{EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc:        "err: ErrPermissionDenied",
			orgRole:     toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:     toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			input:       &pushproto.ListPushesRequest{EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:    "success",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *PushService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe).AnyTimes()
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:       &pushproto.ListPushesRequest{PageSize: 2, Cursor: "", EnvironmentId: "ns0"},
			expected:    &pushproto.ListPushesResponse{Pushes: []*pushproto.Push{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newPushService(mockController, nil, p.orgRole, p.envRole)
			if p.setup != nil {
				p.setup(s)
			}

			actual, err := s.ListPushes(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestGetPushMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)
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
		req         *pushproto.GetPushRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &pushproto.GetPushRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *PushService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(v2ps.ErrPushNotFound)
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &pushproto.GetPushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *PushService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				qe := mock.NewMockQueryExecer(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &pushproto.GetPushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.GetPush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newPushServiceWithMock(t *testing.T, c *gomock.Controller) *PushService {
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

func newPushService(c *gomock.Controller, specifiedEnvironmentId *string, specifiedOrgRole *accountproto.AccountV2_Role_Organization, specifiedEnvRole *accountproto.AccountV2_Role_Environment) *PushService {
	var or accountproto.AccountV2_Role_Organization
	var er accountproto.AccountV2_Role_Environment
	var envId string
	if specifiedOrgRole != nil {
		or = *specifiedOrgRole
	} else {
		or = accountproto.AccountV2_Role_Organization_ADMIN
	}
	if specifiedEnvRole != nil {
		er = *specifiedEnvRole
	} else {
		er = accountproto.AccountV2_Role_Environment_EDITOR
	}
	if specifiedEnvironmentId != nil {
		envId = *specifiedEnvironmentId
	} else {
		envId = "ns0"
	}

	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: or,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: envId,
					Role:          er,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	mysqlClient := mysqlmock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	p.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &PushService{
		mysqlClient:      mysqlClient,
		featureClient:    featureclientmock.NewMockClient(c),
		experimentClient: experimentclientmock.NewMockClient(c),
		accountClient:    accountClientMock,
		publisher:        publishermock.NewMockPublisher(c),
		logger:           zap.NewNop(),
	}
}

func createContextWithToken(t *testing.T, isSystemAdmin bool) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: isSystemAdmin,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

// convert to pointer
func toPtr[T any](value T) *T {
	return &value
}
