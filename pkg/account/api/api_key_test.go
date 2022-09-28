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
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestCreateAPIKeyMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*AccountService)
		ctxRole     accountproto.Account_Role
		req         *accountproto.CreateAPIKeyRequest
		expectedErr error
	}{
		"errNoCommand": {
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAPIKeyRequest{
				Command: nil,
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"errMissingAPIKeyName": {
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAPIKeyRequest{
				Command: &accountproto.CreateAPIKeyCommand{Name: ""},
			},
			expectedErr: localizedError(statusMissingAPIKeyName, locale.JaJP),
		},
		"errInternal": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAPIKeyRequest{
				Command: &accountproto.CreateAPIKeyCommand{
					Name: "name",
					Role: accountproto.APIKey_SDK,
				},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAPIKeyRequest{
				Command: &accountproto.CreateAPIKeyCommand{
					Name: "name",
					Role: accountproto.APIKey_SDK,
				},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithDefaultToken(t, p.ctxRole)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateAPIKey(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, msg)
		})
	}
}

func TestChangeAPIKeyNameMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*AccountService)
		ctxRole     accountproto.Account_Role
		req         *accountproto.ChangeAPIKeyNameRequest
		expectedErr error
	}{
		"errMissingAPIKeyID": {
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ChangeAPIKeyNameRequest{
				Id: "",
			},
			expectedErr: localizedError(statusMissingAPIKeyID, locale.JaJP),
		},
		"errNoCommand": {
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ChangeAPIKeyNameRequest{
				Id:      "id",
				Command: nil,
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"errNotFound": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAPIKeyNotFound)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ChangeAPIKeyNameRequest{
				Id: "id",
				Command: &accountproto.ChangeAPIKeyNameCommand{
					Name: "",
				},
			},
			expectedErr: localizedError(statusNotFound, locale.JaJP),
		},
		"errInternal": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ChangeAPIKeyNameRequest{
				Id: "id",
				Command: &accountproto.ChangeAPIKeyNameCommand{
					Name: "new name",
				},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ChangeAPIKeyNameRequest{
				Id: "id",
				Command: &accountproto.ChangeAPIKeyNameCommand{
					Name: "new name",
				},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithDefaultToken(t, p.ctxRole)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ChangeAPIKeyName(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, msg)
		})
	}
}

func TestEnableAPIKeyMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*AccountService)
		ctxRole     accountproto.Account_Role
		req         *accountproto.EnableAPIKeyRequest
		expectedErr error
	}{
		"errMissingAPIKeyID": {
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAPIKeyRequest{
				Id: "",
			},
			expectedErr: localizedError(statusMissingAPIKeyID, locale.JaJP),
		},
		"errNoCommand": {
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAPIKeyRequest{
				Id:      "id",
				Command: nil,
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"errNotFound": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAPIKeyNotFound)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.EnableAPIKeyCommand{},
			},
			expectedErr: localizedError(statusNotFound, locale.JaJP),
		},
		"errInternal": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.EnableAPIKeyCommand{},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.EnableAPIKeyCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithDefaultToken(t, p.ctxRole)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.EnableAPIKey(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, msg)
		})
	}
}

func TestDisableAPIKeyMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*AccountService)
		ctxRole     accountproto.Account_Role
		req         *accountproto.DisableAPIKeyRequest
		expectedErr error
	}{
		"errMissingAPIKeyID": {
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAPIKeyRequest{
				Id: "",
			},
			expectedErr: localizedError(statusMissingAPIKeyID, locale.JaJP),
		},
		"errNoCommand": {
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAPIKeyRequest{
				Id:      "id",
				Command: nil,
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"errNotFound": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAPIKeyNotFound)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.DisableAPIKeyCommand{},
			},
			expectedErr: localizedError(statusNotFound, locale.JaJP),
		},
		"errInternal": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.DisableAPIKeyCommand{},
			},
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.DisableAPIKeyCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithDefaultToken(t, p.ctxRole)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DisableAPIKey(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, msg)
		})
	}
}

func TestGetAPIKeyMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup    func(*AccountService)
		req      *accountproto.GetAPIKeyRequest
		expected error
	}{
		"errMissingAPIKeyID": {
			req:      &accountproto.GetAPIKeyRequest{Id: ""},
			expected: localizedError(statusMissingAPIKeyID, locale.JaJP),
		},
		"errNotFound": {
			setup: func(s *AccountService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req:      &accountproto.GetAPIKeyRequest{Id: "id"},
			expected: localizedError(statusNotFound, locale.JaJP),
		},
		"success": {
			setup: func(s *AccountService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req:      &accountproto.GetAPIKeyRequest{Id: "id"},
			expected: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithDefaultToken(t, accountproto.Account_OWNER)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			res, err := service.GetAPIKey(ctx, p.req)
			assert.Equal(t, p.expected, err)
			if err == nil {
				assert.NotNil(t, res)
			}
		})
	}
}

func TestListAPIKeysMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*AccountService)
		input       *accountproto.ListAPIKeysRequest
		expected    *accountproto.ListAPIKeysResponse
		expectedErr error
	}{
		"errInvalidCursor": {
			input:       &accountproto.ListAPIKeysRequest{Cursor: "XXX"},
			expected:    nil,
			expectedErr: localizedError(statusInvalidCursor, locale.JaJP),
		},
		"errInternal": {
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			input:       &accountproto.ListAPIKeysRequest{},
			expected:    nil,
			expectedErr: localizedError(statusInternal, locale.JaJP),
		},
		"success": {
			setup: func(s *AccountService) {
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
			input:       &accountproto.ListAPIKeysRequest{PageSize: 2, Cursor: ""},
			expected:    &accountproto.ListAPIKeysResponse{ApiKeys: []*accountproto.APIKey{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithDefaultToken(t, accountproto.Account_OWNER)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.ListAPIKeys(ctx, p.input)
			assert.Equal(t, p.expectedErr, err, msg)
			assert.Equal(t, p.expected, actual, msg)
		})
	}
}
