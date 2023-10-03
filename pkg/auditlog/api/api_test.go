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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	v2alsmock "github.com/bucketeer-io/bucketeer/pkg/auditlog/storage/v2/mock"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	proto "github.com/bucketeer-io/bucketeer/proto/auditlog"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestNewAuditLogService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	accountClientMock := accountclientmock.NewMockClient(mockController)
	mysqlClient := mysqlmock.NewMockClient(mockController)
	logger := zap.NewNop()
	s := NewAuditLogService(accountClientMock, mysqlClient, WithLogger(logger))
	assert.IsType(t, &auditlogService{}, s)
}

func TestListAuditLogsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
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
		setup       func(*auditlogService)
		input       *proto.ListAuditLogsRequest
		expected    *proto.ListAuditLogsResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			setup:       nil,
			input:       &proto.ListAuditLogsRequest{Cursor: "XXX"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *auditlogService) {
				s.mysqlStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("test"))
			},
			input:       &proto.ListAuditLogsRequest{},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *auditlogService) {
				s.mysqlStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(createAuditLogs(t), 2, int64(10), nil)
			},
			input:       &proto.ListAuditLogsRequest{PageSize: 2, Cursor: "", EnvironmentNamespace: "ns0"},
			expected:    &proto.ListAuditLogsResponse{AuditLogs: createAuditLogs(t), Cursor: "2", TotalCount: 10},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newAuditLogService(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListAuditLogs(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestListAdminAuditLogsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, accountproto.Account_OWNER)
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
		setup       func(*auditlogService)
		input       *proto.ListAdminAuditLogsRequest
		expected    *proto.ListAdminAuditLogsResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			setup:       nil,
			input:       &proto.ListAdminAuditLogsRequest{Cursor: "invalid"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *auditlogService) {
				s.mysqlAdminStorage.(*v2alsmock.MockAdminAuditLogStorage).EXPECT().ListAdminAuditLogs(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("test"))
			},
			input:       &proto.ListAdminAuditLogsRequest{},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *auditlogService) {
				s.mysqlAdminStorage.(*v2alsmock.MockAdminAuditLogStorage).EXPECT().ListAdminAuditLogs(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(createAuditLogs(t), 2, int64(10), nil)
			},
			input:       &proto.ListAdminAuditLogsRequest{PageSize: 2, Cursor: ""},
			expected:    &proto.ListAdminAuditLogsResponse{AuditLogs: createAuditLogs(t), Cursor: "2", TotalCount: 10},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newAuditLogService(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListAdminAuditLogs(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestListFeatureHistoryMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
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
		setup       func(*auditlogService)
		input       *proto.ListFeatureHistoryRequest
		expected    *proto.ListFeatureHistoryResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			setup:       nil,
			input:       &proto.ListFeatureHistoryRequest{Cursor: "XXX"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *auditlogService) {
				s.mysqlStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("test"))
			},
			input:       &proto.ListFeatureHistoryRequest{},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *auditlogService) {
				s.mysqlStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(createAuditLogs(t), 2, int64(10), nil)
			},
			input: &proto.ListFeatureHistoryRequest{
				FeatureId: "fid-1", PageSize: 2, Cursor: "", EnvironmentNamespace: "ns0",
			},
			expected:    &proto.ListFeatureHistoryResponse{AuditLogs: createAuditLogs(t), Cursor: "2", TotalCount: int64(10)},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newAuditLogService(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListFeatureHistory(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newAuditLogService(t *testing.T, mockController *gomock.Controller) *auditlogService {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	accountClientMock := accountclientmock.NewMockClient(mockController)
	ar := &accountproto.GetAccountResponse{
		Account: &accountproto.Account{
			Email: "email",
			Role:  accountproto.Account_VIEWER,
		},
	}
	accountClientMock.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	return &auditlogService{
		accountClient:     accountClientMock,
		mysqlStorage:      v2alsmock.NewMockAuditLogStorage(mockController),
		mysqlAdminStorage: v2alsmock.NewMockAdminAuditLogStorage(mockController),
		logger:            logger.Named("api"),
	}
}

func createAuditLogs(t *testing.T) []*proto.AuditLog {
	t.Helper()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	msgUnknown := domainevent.LocalizedMessage(domaineventproto.Event_UNKNOWN, localizer)
	return []*proto.AuditLog{
		{Id: "id-0", LocalizedMessage: msgUnknown},
		{Id: "id-1", LocalizedMessage: msgUnknown},
	}
}

func createContextWithToken(t *testing.T, role accountproto.Account_Role) context.Context {
	t.Helper()
	token := &token.IDToken{
		Email:     "test@example.com",
		AdminRole: role,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}
