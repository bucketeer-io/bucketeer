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

package role

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
			expected:    &eventproto.Editor{Email: "test@example.com", Role: accountproto.Account_OWNER, IsAdmin: true},
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

	patterns := []struct {
		inputCtx            context.Context
		inputRequiredRole   accountproto.Account_Role
		inputGetAccountFunc func(email string) (*accountproto.GetAccountResponse, error)
		expected            *eventproto.Editor
		expectedErr         error
	}{
		{
			inputCtx:          context.Background(),
			inputRequiredRole: accountproto.Account_EDITOR,
			expected:          nil,
			expectedErr:       ErrUnauthenticated,
		},
		{
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
			inputRequiredRole: accountproto.Account_EDITOR,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountResponse, error) {
				return nil, status.Error(codes.NotFound, "")
			},
			expected:    nil,
			expectedErr: ErrUnauthenticated,
		},
		{
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
			inputRequiredRole: accountproto.Account_EDITOR,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountResponse, error) {
				return nil, status.Error(codes.Internal, "")
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
			inputRequiredRole: accountproto.Account_EDITOR,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountResponse, error) {
				return &accountproto.GetAccountResponse{
					Account: &accountproto.Account{Email: "test@example.com", Role: accountproto.Account_VIEWER},
				}, nil
			},
			expected:    nil,
			expectedErr: ErrPermissionDenied,
		},
		{
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_UNASSIGNED}),
			inputRequiredRole: accountproto.Account_EDITOR,
			inputGetAccountFunc: func(email string) (*accountproto.GetAccountResponse, error) {
				return &accountproto.GetAccountResponse{
					Account: &accountproto.Account{Email: "test@example.com", Role: accountproto.Account_EDITOR},
				}, nil
			},
			expected:    &eventproto.Editor{Email: "test@example.com", Role: accountproto.Account_EDITOR, IsAdmin: false},
			expectedErr: nil,
		},
		{
			inputCtx:          getContextWithToken(t, &token.IDToken{Email: "test@example.com", AdminRole: accountproto.Account_OWNER}),
			inputRequiredRole: accountproto.Account_OWNER,
			expected:          &eventproto.Editor{Email: "test@example.com", Role: accountproto.Account_OWNER, IsAdmin: true},
			expectedErr:       nil,
		},
	}
	for _, p := range patterns {
		editor, err := CheckRole(p.inputCtx, p.inputRequiredRole, p.inputGetAccountFunc)
		assert.Equal(t, p.expectedErr, err)
		assert.Equal(t, p.expected, editor)
	}
}

func getContextWithToken(t *testing.T, token *token.IDToken) context.Context {
	t.Helper()
	return context.WithValue(context.Background(), rpc.Key, token)
}
