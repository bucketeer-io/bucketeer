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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	gstatus "google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const userKind = "User"

func TestValidateGetUserRequest(t *testing.T) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []struct {
		input    *userproto.GetUserRequest
		expected error
	}{
		{
			input:    &userproto.GetUserRequest{UserId: "test", EnvironmentNamespace: "ns0"},
			expected: nil,
		},
		{
			input:    &userproto.GetUserRequest{EnvironmentNamespace: "ns0"},
			expected: createError(statusMissingUserID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "user_id")),
		},
	}
	s := userService{}
	for _, p := range patterns {
		err := s.validateGetUserRequest(p.input, localizer)
		assert.Equal(t, p.expected, err)
	}
}

func TestGetUserRequest(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
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
		setup       func(s *userService)
		input       *userproto.GetUserRequest
		expected    *userproto.GetUserResponse
		expectedErr error
	}{
		{
			desc: "user not found",
			setup: func(s *userService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.storageClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:       &userproto.GetUserRequest{UserId: "user-id-0", EnvironmentNamespace: "ns0"},
			expected:    nil,
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "internal error",
			setup: func(s *userService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("internal error"))
				s.storageClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:       &userproto.GetUserRequest{UserId: "user-id-1", EnvironmentNamespace: "ns0"},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *userService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.storageClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input: &userproto.GetUserRequest{UserId: "user-id-1", EnvironmentNamespace: "ns0"},
			expected: &userproto.GetUserResponse{
				User: &userproto.User{
					Id:         "",
					Data:       nil,
					TaggedData: nil,
					LastSeen:   0,
					CreatedAt:  0,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createUserService(mockController)
		p.setup(service)
		ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
		actual, err := service.GetUser(ctx, p.input)
		assert.Equal(t, p.expected, actual)
		assert.Equal(t, p.expectedErr, err)
	}
}

func createUserService(c *gomock.Controller) *userService {
	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountResponse{
		Account: &accountproto.Account{
			Email: "email",
			Role:  accountproto.Account_VIEWER,
		},
	}
	accountClientMock.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	return &userService{
		storageClient: mysqlmock.NewMockClient(c),
		accountClient: accountClientMock,
		logger:        zap.NewNop().Named("api"),
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
