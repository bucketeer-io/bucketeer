// Copyright 2026 The Bucketeer Authors.
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

package role

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

func TestCheckSystemAdminRole(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		inputCtx    context.Context
		expected    *eventproto.Editor
		expectedErr error
	}{
		{
			inputCtx:    context.Background(),
			expected:    nil,
			expectedErr: ErrUnauthenticated,
		},
		{
			inputCtx:    getContextWithToken(t, &token.AccessToken{Email: "test@example.com", IsSystemAdmin: false}),
			expected:    nil,
			expectedErr: ErrPermissionDenied,
		},
		{
			inputCtx:    getContextWithToken(t, &token.AccessToken{Email: "test@example.com", IsSystemAdmin: true}),
			expected:    &eventproto.Editor{Email: "test@example.com", IsAdmin: true},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		editor, err := CheckSystemAdminRole(p.inputCtx)
		assert.Equal(t, p.expectedErr, err)
		assert.Equal(t, p.expected, editor)
	}
}

func TestCheckRole(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	env := "ns0"

	patterns := []struct {
		desc                string
		inputCtx            context.Context
		inputRequiredRole   accountproto.AccountV2_Role_Environment
		inputGetAccountFunc func(email string) (*accountproto.AccountV2, error)
		expected            *eventproto.Editor
		expectedErr         error
	}{
		{
			desc:              "unauthenticated: no token",
			inputCtx:          context.Background(),
			inputRequiredRole: accountproto.AccountV2_Role_Environment_EDITOR,
			expected:          nil,
			expectedErr:       ErrUnauthenticated,
		},
		{
			desc:              "unauthenticated: account not found",
			inputCtx:          getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			inputRequiredRole: accountproto.AccountV2_Role_Environment_EDITOR,
			inputGetAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			expected:    nil,
			expectedErr: ErrUnauthenticated,
		},
		{
			desc:              "internalError",
			inputCtx:          getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			inputRequiredRole: accountproto.AccountV2_Role_Environment_EDITOR,
			inputGetAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return nil, status.Error(codes.Internal, "")
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc:              "permissionDenied",
			inputCtx:          getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			inputRequiredRole: accountproto.AccountV2_Role_Environment_EDITOR,
			inputGetAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				resp := &accountproto.GetAccountV2ByEnvironmentIDResponse{
					Account: &accountproto.AccountV2{
						Email:            "test@example.com",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}
				return resp.Account, nil
			},
			expected:    nil,
			expectedErr: ErrPermissionDenied,
		},
		{
			desc:              "success",
			inputCtx:          getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			inputRequiredRole: accountproto.AccountV2_Role_Environment_EDITOR,
			inputGetAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				resp := &accountproto.GetAccountV2ByEnvironmentIDResponse{
					Account: &accountproto.AccountV2{
						Email:            "test@example.com",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_EDITOR,
							},
						},
					},
				}
				return resp.Account, nil
			},
			expected: &eventproto.Editor{
				Email:   "test@example.com",
				IsAdmin: false,
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{
						EnvironmentId: "ns0",
						Role:          accountproto.AccountV2_Role_Environment_EDITOR,
					},
				},
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			editor, err := CheckEnvironmentRole(
				p.inputCtx, p.inputRequiredRole,
				env, p.inputGetAccountFunc)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, editor)
		})
	}
}

func TestCheckOrganizationRole(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc                string
		inputCtx            context.Context
		inputRequiredRole   accountproto.AccountV2_Role_Organization
		inputGetAccountFunc func(email string) (*accountproto.GetAccountV2Response, error)
		expected            *eventproto.Editor
		expectedErr         error
	}{
		{
			desc:              "unauthenticated: no token",
			inputCtx:          context.Background(),
			inputRequiredRole: accountproto.AccountV2_Role_Organization_MEMBER,
			expected:          nil,
			expectedErr:       ErrUnauthenticated,
		},
		{
			desc:              "unauthenticated: account not found",
			inputCtx:          getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			inputRequiredRole: accountproto.AccountV2_Role_Organization_MEMBER,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			expected:    nil,
			expectedErr: ErrUnauthenticated,
		},
		{
			desc:              "internalError",
			inputCtx:          getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			inputRequiredRole: accountproto.AccountV2_Role_Organization_MEMBER,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return nil, status.Error(codes.Internal, "")
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc:              "unauthenticated: account disabled",
			inputCtx:          getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			inputRequiredRole: accountproto.AccountV2_Role_Organization_ADMIN,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return &accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						Email:            "test@example.com",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						Disabled:         true,
					},
				}, nil
			},
			expected:    nil,
			expectedErr: ErrUnauthenticated,
		},
		{
			desc:              "permissionDenied",
			inputCtx:          getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			inputRequiredRole: accountproto.AccountV2_Role_Organization_ADMIN,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return &accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{Email: "test@example.com", OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER},
				}, nil
			},
			expected:    nil,
			expectedErr: ErrPermissionDenied,
		},
		{
			desc:              "success",
			inputCtx:          getContextWithToken(t, &token.AccessToken{Email: "test@example.com", Name: "test"}),
			inputRequiredRole: accountproto.AccountV2_Role_Organization_ADMIN,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return &accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						Email:            "test@example.com",
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil
			},
			expected: &eventproto.Editor{
				Email:   "test@example.com",
				IsAdmin: false, Name: "test",
				OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			editor, err := CheckOrganizationRole(p.inputCtx, p.inputRequiredRole, p.inputGetAccountFunc)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, editor)
		})
	}
}

func getContextWithToken(t *testing.T, token *token.AccessToken) context.Context {
	t.Helper()
	return context.WithValue(context.Background(), rpc.AccessTokenKey, token)
}

func getContextWithTokenAndAPIKey(
	t *testing.T,
	token *token.AccessToken,
	apiKeyToken string,
	apiKeyMaintainer string,
	apiKeyName string,
) context.Context {
	t.Helper()
	ctx := context.WithValue(context.Background(), rpc.AccessTokenKey, token)
	headerMetaData := metadata.New(map[string]string{
		APIKeyTokenMDKey:      apiKeyToken,
		APIKeyMaintainerMDKey: apiKeyMaintainer,
		APIKeyNameMDKey:       apiKeyName,
	})

	return metadata.NewIncomingContext(ctx, headerMetaData)
}

func TestCheckEnvironmentRole(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		ctx            context.Context
		requiredRole   accountproto.AccountV2_Role_Environment
		environmentID  string
		getAccountFunc func(email string) (*accountproto.AccountV2, error)
		expected       *eventproto.Editor
		expectedErr    error
	}{
		{
			desc:          "err: account not found",
			ctx:           getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			expected:    nil,
			expectedErr: ErrUnauthenticated,
		},
		{
			desc:          "err: account disabled",
			ctx:           getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return &accountproto.AccountV2{
					Disabled: true,
				}, nil
			},
			expected:    nil,
			expectedErr: ErrUnauthenticated,
		},
		{
			desc:          "success: environment role satisfied",
			ctx:           getContextWithToken(t, &token.AccessToken{Email: "test@example.com", Name: "test"}),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return &accountproto.AccountV2{
					Email: "test@example.com",
					EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
						{EnvironmentId: "ns0", Role: accountproto.AccountV2_Role_Environment_EDITOR},
					},
					Name:     "test",
					Disabled: false,
				}, nil
			},
			expected: &eventproto.Editor{
				Email: "test@example.com",
				Name:  "test",
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "ns0", Role: accountproto.AccountV2_Role_Environment_EDITOR},
				},
			},
			expectedErr: nil,
		},
		{
			desc:          "success: organization role satisfied",
			ctx:           getContextWithToken(t, &token.AccessToken{Email: "test@example.com", Name: "test"}),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return &accountproto.AccountV2{
					Email:            "test@example.com",
					OrganizationRole: accountproto.AccountV2_Role_Organization_OWNER,
					Name:             "test",
					Disabled:         false,
				}, nil
			},
			expected: &eventproto.Editor{
				Email:            "test@example.com",
				Name:             "test",
				OrganizationRole: accountproto.AccountV2_Role_Organization_OWNER,
			},
			expectedErr: nil,
		},
		{
			desc: "success get API key editor",
			ctx: getContextWithTokenAndAPIKey(
				t,
				&token.AccessToken{Email: "localenv@bucketeer.io", IsSystemAdmin: true},
				"apikey_token",
				"apikey_maintainer@example.com",
				"apikey_name",
			),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return &accountproto.AccountV2{
					Email: "apikey_maintainer@example.com",
					EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
						{EnvironmentId: "ns0", Role: accountproto.AccountV2_Role_Environment_EDITOR},
					},
					FirstName: "apikey",
					LastName:  "maintainer",
					Disabled:  false,
				}, nil
			},
			expected: &eventproto.Editor{
				Email: "apikey_maintainer@example.com",
				Name:  "apikey maintainer",
				PublicApiEditor: &eventproto.Editor_PublicAPIEditor{
					Token:      "apikey_token",
					Maintainer: "apikey_maintainer@example.com",
					Name:       "apikey_name",
				},
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "ns0", Role: accountproto.AccountV2_Role_Environment_EDITOR},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			editor, err := CheckEnvironmentRole(
				p.ctx,
				p.requiredRole,
				p.environmentID,
				p.getAccountFunc)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, editor)
		})
	}
}

func TestCheckEnvironmentRoleWithLog(t *testing.T) {
	t.Parallel()

	var (
		errCustomUnauthenticated  = errors.New("custom unauthenticated")
		errCustomPermissionDenied = errors.New("custom permission denied")
		defaultErrFunc            = func(err error) error {
			return fmt.Errorf("wrapped: %w", err)
		}
	)

	patterns := []struct {
		desc             string
		ctx              context.Context
		requiredRole     accountproto.AccountV2_Role_Environment
		environmentID    string
		getAccountFunc   func(email string) (*accountproto.AccountV2, error)
		expected         *eventproto.Editor
		expectedErr      error
		expectedLogCount int
		expectedLogMsg   string
		expectedEmail    string
	}{
		{
			desc:          "unauthenticated: no token returns custom unauthenticated error",
			ctx:           context.Background(),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			expected:         nil,
			expectedErr:      errCustomUnauthenticated,
			expectedLogCount: 1,
			expectedLogMsg:   "Unauthenticated",
			expectedEmail:    "",
		},
		{
			desc:          "unauthenticated: account not found returns custom unauthenticated error",
			ctx:           getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			expected:         nil,
			expectedErr:      errCustomUnauthenticated,
			expectedLogCount: 1,
			expectedLogMsg:   "Unauthenticated",
			expectedEmail:    "test@example.com",
		},
		{
			desc:          "unauthenticated: account disabled returns custom unauthenticated error",
			ctx:           getContextWithToken(t, &token.AccessToken{Email: "disabled@example.com"}),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return &accountproto.AccountV2{Disabled: true}, nil
			},
			expected:         nil,
			expectedErr:      errCustomUnauthenticated,
			expectedLogCount: 1,
			expectedLogMsg:   "Unauthenticated",
			expectedEmail:    "disabled@example.com",
		},
		{
			desc:          "permission denied: insufficient role returns custom permission denied error",
			ctx:           getContextWithToken(t, &token.AccessToken{Email: "viewer@example.com"}),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return &accountproto.AccountV2{
					Email:            "viewer@example.com",
					OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
					EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
						{EnvironmentId: "ns0", Role: accountproto.AccountV2_Role_Environment_VIEWER},
					},
				}, nil
			},
			expected:         nil,
			expectedErr:      errCustomPermissionDenied,
			expectedLogCount: 1,
			expectedLogMsg:   "Permission denied",
			expectedEmail:    "viewer@example.com",
		},
		{
			desc:          "default error: internal error uses defaultErrFunc",
			ctx:           getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return nil, status.Error(codes.Internal, "something broke")
			},
			expected:         nil,
			expectedErr:      fmt.Errorf("wrapped: %w", ErrInternal),
			expectedLogCount: 1,
			expectedLogMsg:   "Failed to check role",
			expectedEmail:    "test@example.com",
		},
		{
			desc:          "success: returns editor from underlying CheckEnvironmentRole",
			ctx:           getContextWithToken(t, &token.AccessToken{Email: "test@example.com", Name: "test"}),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return &accountproto.AccountV2{
					Email: "test@example.com",
					EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
						{EnvironmentId: "ns0", Role: accountproto.AccountV2_Role_Environment_EDITOR},
					},
				}, nil
			},
			expected: &eventproto.Editor{
				Email: "test@example.com",
				Name:  "test",
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "ns0", Role: accountproto.AccountV2_Role_Environment_EDITOR},
				},
			},
			expectedErr:      nil,
			expectedLogCount: 0,
		},
		{
			desc:          "success: system admin bypasses role check",
			ctx:           getContextWithToken(t, &token.AccessToken{Email: "admin@example.com", Name: "admin", IsSystemAdmin: true}),
			requiredRole:  accountproto.AccountV2_Role_Environment_EDITOR,
			environmentID: "ns0",
			getAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				t.Fatal("getAccountFunc should not be called for system admin")
				return nil, nil
			},
			expected: &eventproto.Editor{
				Email:   "admin@example.com",
				Name:    "admin",
				IsAdmin: true,
			},
			expectedErr:      nil,
			expectedLogCount: 0,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			var buf bytes.Buffer
			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(&buf),
				zapcore.ErrorLevel,
			)
			logger := zap.New(core)

			editor, err := CheckEnvironmentRoleWithLog(
				p.ctx,
				p.requiredRole,
				p.environmentID,
				p.getAccountFunc,
				logger,
				errCustomUnauthenticated,
				errCustomPermissionDenied,
				defaultErrFunc,
			)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, editor)

			if p.expectedLogCount == 0 {
				assert.Empty(t, buf.String(), "expected no log output")
				return
			}

			require.NotEmpty(t, buf.String(), "expected log output")
			var logEntry map[string]interface{}
			require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry),
				"log output should be valid JSON")

			assert.Equal(t, p.expectedLogMsg, logEntry["msg"])
			assert.Equal(t, p.environmentID, logEntry["environmentId"])
			assert.Equal(t, p.requiredRole.String(), logEntry["requiredRole"])
			if p.expectedEmail != "" {
				assert.Equal(t, p.expectedEmail, logEntry["email"])
			} else {
				_, hasEmail := logEntry["email"]
				assert.False(t, hasEmail,
					"email field should not be present when there is no token")
			}
		})
	}
}

func TestCheckOrganizationRoleWithLog(t *testing.T) {
	t.Parallel()

	var (
		errCustomUnauthenticated  = errors.New("custom unauthenticated")
		errCustomPermissionDenied = errors.New("custom permission denied")
		defaultErrFunc            = func(err error) error {
			return fmt.Errorf("wrapped: %w", err)
		}
	)

	patterns := []struct {
		desc             string
		ctx              context.Context
		requiredRole     accountproto.AccountV2_Role_Organization
		organizationID   string
		getAccountFunc   func(email string) (*accountproto.GetAccountV2Response, error)
		expected         *eventproto.Editor
		expectedErr      error
		expectedLogCount int
		expectedLogMsg   string
		expectedEmail    string
	}{
		{
			desc:           "unauthenticated: no token returns custom unauthenticated error",
			ctx:            context.Background(),
			requiredRole:   accountproto.AccountV2_Role_Organization_ADMIN,
			organizationID: "org0",
			getAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			expected:         nil,
			expectedErr:      errCustomUnauthenticated,
			expectedLogCount: 1,
			expectedLogMsg:   "Unauthenticated",
			expectedEmail:    "",
		},
		{
			desc:           "unauthenticated: account not found returns custom unauthenticated error",
			ctx:            getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			requiredRole:   accountproto.AccountV2_Role_Organization_MEMBER,
			organizationID: "org0",
			getAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			expected:         nil,
			expectedErr:      errCustomUnauthenticated,
			expectedLogCount: 1,
			expectedLogMsg:   "Unauthenticated",
			expectedEmail:    "test@example.com",
		},
		{
			desc:           "unauthenticated: account disabled returns custom unauthenticated error",
			ctx:            getContextWithToken(t, &token.AccessToken{Email: "disabled@example.com"}),
			requiredRole:   accountproto.AccountV2_Role_Organization_MEMBER,
			organizationID: "org0",
			getAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return &accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{Disabled: true},
				}, nil
			},
			expected:         nil,
			expectedErr:      errCustomUnauthenticated,
			expectedLogCount: 1,
			expectedLogMsg:   "Unauthenticated",
			expectedEmail:    "disabled@example.com",
		},
		{
			desc:           "permission denied: insufficient role returns custom permission denied error",
			ctx:            getContextWithToken(t, &token.AccessToken{Email: "member@example.com"}),
			requiredRole:   accountproto.AccountV2_Role_Organization_ADMIN,
			organizationID: "org0",
			getAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return &accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						Email:            "member@example.com",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
					},
				}, nil
			},
			expected:         nil,
			expectedErr:      errCustomPermissionDenied,
			expectedLogCount: 1,
			expectedLogMsg:   "Permission denied",
			expectedEmail:    "member@example.com",
		},
		{
			desc:           "default error: internal error uses defaultErrFunc",
			ctx:            getContextWithToken(t, &token.AccessToken{Email: "test@example.com"}),
			requiredRole:   accountproto.AccountV2_Role_Organization_ADMIN,
			organizationID: "org0",
			getAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return nil, status.Error(codes.Internal, "something broke")
			},
			expected:         nil,
			expectedErr:      fmt.Errorf("wrapped: %w", ErrInternal),
			expectedLogCount: 1,
			expectedLogMsg:   "Failed to check role",
			expectedEmail:    "test@example.com",
		},
		{
			desc:           "success: returns editor from underlying CheckOrganizationRole",
			ctx:            getContextWithToken(t, &token.AccessToken{Email: "test@example.com", Name: "test"}),
			requiredRole:   accountproto.AccountV2_Role_Organization_ADMIN,
			organizationID: "org0",
			getAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return &accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						Email:            "test@example.com",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil
			},
			expected: &eventproto.Editor{
				Email:            "test@example.com",
				Name:             "test",
				OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			},
			expectedErr:      nil,
			expectedLogCount: 0,
		},
		{
			desc:           "success: system admin bypasses role check",
			ctx:            getContextWithToken(t, &token.AccessToken{Email: "admin@example.com", Name: "admin", IsSystemAdmin: true}),
			requiredRole:   accountproto.AccountV2_Role_Organization_OWNER,
			organizationID: "org0",
			getAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				t.Fatal("getAccountFunc should not be called for system admin")
				return nil, nil
			},
			expected: &eventproto.Editor{
				Email:   "admin@example.com",
				Name:    "admin",
				IsAdmin: true,
			},
			expectedErr:      nil,
			expectedLogCount: 0,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			var buf bytes.Buffer
			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(&buf),
				zapcore.ErrorLevel,
			)
			logger := zap.New(core)

			editor, err := CheckOrganizationRoleWithLog(
				p.ctx,
				p.requiredRole,
				p.organizationID,
				p.getAccountFunc,
				logger,
				errCustomUnauthenticated,
				errCustomPermissionDenied,
				defaultErrFunc,
			)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, editor)

			if p.expectedLogCount == 0 {
				assert.Empty(t, buf.String(), "expected no log output")
				return
			}

			require.NotEmpty(t, buf.String(), "expected log output")
			var logEntry map[string]interface{}
			require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry),
				"log output should be valid JSON")

			assert.Equal(t, p.expectedLogMsg, logEntry["msg"])
			assert.Equal(t, p.organizationID, logEntry["organizationId"])
			assert.Equal(t, p.requiredRole.String(), logEntry["requiredRole"])
			if p.expectedEmail != "" {
				assert.Equal(t, p.expectedEmail, logEntry["email"])
			} else {
				_, hasEmail := logEntry["email"]
				assert.False(t, hasEmail,
					"email field should not be present when there is no token")
			}
		})
	}
}

func TestCheckOrganizationRoleByEnvironmentIDWithLog(t *testing.T) {
	t.Parallel()

	var (
		errCustomUnauthenticated = errors.New("custom unauthenticated")
		defaultErrFunc           = func(err error) error { return err }
	)

	t.Run("unauthenticated logs environmentID field", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(&buf),
			zapcore.ErrorLevel,
		)
		logger := zap.New(core)

		editor, err := CheckOrganizationRoleByEnvironmentIDWithLog(
			context.Background(),
			accountproto.AccountV2_Role_Organization_MEMBER,
			"env0",
			func(email string) (*accountproto.GetAccountV2Response, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			logger,
			errCustomUnauthenticated,
			ErrPermissionDenied,
			defaultErrFunc,
		)
		assert.Nil(t, editor)
		assert.Equal(t, errCustomUnauthenticated, err)

		require.NotEmpty(t, buf.String(), "expected log output")
		var logEntry map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry),
			"log output should be valid JSON")
		assert.Equal(t, "Unauthenticated", logEntry["msg"])
		assert.Equal(t, "env0", logEntry["environmentId"])
		_, hasOrganizationID := logEntry["organizationId"]
		assert.False(t, hasOrganizationID)
	})
}
