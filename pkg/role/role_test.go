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

package role

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestCheckAdminRole(t *testing.T) {
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
			inputCtx:    getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
			expected:    nil,
			expectedErr: ErrPermissionDenied,
		},
		{
			inputCtx:    getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_OWNER}),
			expected:    &eventproto.Editor{Email: "test@example.com", IsAdmin: true},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		editor, err := CheckAdminRole(p.inputCtx)
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
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
			inputRequiredRole: accountproto.AccountV2_Role_Environment_EDITOR,
			inputGetAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			expected:    nil,
			expectedErr: ErrUnauthenticated,
		},
		{
			desc:              "internalError",
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
			inputRequiredRole: accountproto.AccountV2_Role_Environment_EDITOR,
			inputGetAccountFunc: func(email string) (*accountproto.AccountV2, error) {
				return nil, status.Error(codes.Internal, "")
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc:              "permissionDenied",
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
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
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
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
			expected:    &eventproto.Editor{Email: "test@example.com", IsAdmin: false},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		editor, err := CheckRole(p.inputCtx, p.inputRequiredRole, env, p.inputGetAccountFunc)
		assert.Equal(t, p.expectedErr, err)
		assert.Equal(t, p.expected, editor)
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
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
			inputRequiredRole: accountproto.AccountV2_Role_Organization_MEMBER,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			expected:    nil,
			expectedErr: ErrUnauthenticated,
		},
		{
			desc:              "internalError",
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
			inputRequiredRole: accountproto.AccountV2_Role_Organization_MEMBER,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return nil, status.Error(codes.Internal, "")
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc:              "permissionDenied",
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
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
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
			inputRequiredRole: accountproto.AccountV2_Role_Organization_ADMIN,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountV2Response, error) {
				return &accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{Email: "test@example.com", OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN},
				}, nil
			},
			expected:    &eventproto.Editor{Email: "test@example.com", IsAdmin: false},
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

func getContextWithToken(t *testing.T, token *token.IDToken) context.Context {
	t.Helper()
	return context.WithValue(context.Background(), rpc.Key, token)
}
