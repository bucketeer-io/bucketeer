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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

var (
	ErrUnauthenticated  = status.Error(codes.Unauthenticated, "Unauthenticated user")
	ErrPermissionDenied = status.Error(codes.PermissionDenied, "Permission denied")
	ErrInternal         = status.Error(codes.Internal, "Internal")
)

func CheckAdminRole(ctx context.Context) (*eventproto.Editor, error) {
	token, ok := rpc.GetIDToken(ctx)
	if !ok {
		return nil, ErrUnauthenticated
	}
	if !token.IsAdmin() {
		return nil, ErrPermissionDenied
	}
	return checkRole(token.Email, accountproto.Account_OWNER, accountproto.Account_OWNER, true)
}

func CheckRole(
	ctx context.Context,
	requiredRole accountproto.Account_Role,
	getAccountFunc func(email string) (*accountproto.GetAccountResponse, error),
) (*eventproto.Editor, error) {
	token, ok := rpc.GetIDToken(ctx)
	if !ok {
		return nil, ErrUnauthenticated
	}
	if !token.IsAdmin() {
		// get account for the environment namespace
		resp, err := getAccountFunc(token.Email)
		if err != nil {
			if code := status.Code(err); code == codes.NotFound {
				return nil, ErrUnauthenticated
			}
			return nil, ErrInternal
		}
		return checkRole(resp.Account.Email, resp.Account.Role, requiredRole, false)
	}
	return checkRole(token.Email, accountproto.Account_OWNER, requiredRole, true)
}

func checkRole(email string, role, requiredRole accountproto.Account_Role, isAdmin bool) (*eventproto.Editor, error) {
	if role == accountproto.Account_UNASSIGNED || role < requiredRole {
		return nil, ErrPermissionDenied
	}
	return &eventproto.Editor{
		Email:   email,
		Role:    role,
		IsAdmin: isAdmin,
	}, nil
}
