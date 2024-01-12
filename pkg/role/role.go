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
	return &eventproto.Editor{
		Email:   token.Email,
		IsAdmin: true,
	}, nil
}

func CheckRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentID string,
	getAccountFunc func(email string) (*accountproto.AccountV2, error),
) (*eventproto.Editor, error) {
	token, ok := rpc.GetIDToken(ctx)
	if !ok {
		return nil, ErrUnauthenticated
	}
	// TODO remove this condition after migration to AccountV2
	if !token.IsAdmin() {
		// get account for the environment namespace
		account, err := getAccountFunc(token.Email)
		if err != nil {
			if code := status.Code(err); code == codes.NotFound {
				return nil, ErrUnauthenticated
			}
			return nil, ErrInternal
		}
		accountEnvRole := getRole(account.EnvironmentRoles, environmentID)
		return checkRole(account.Email, accountEnvRole, requiredRole, false)
	}
	return checkRole(token.Email, accountproto.AccountV2_Role_Environment_EDITOR, requiredRole, true)
}

func getRole(roles []*accountproto.AccountV2_EnvironmentRole, envID string) accountproto.AccountV2_Role_Environment {
	for _, role := range roles {
		if role.EnvironmentId == envID {
			return role.Role
		}
	}
	return accountproto.AccountV2_Role_Environment_UNASSIGNED
}

func checkRole(
	email string,
	role, requiredRole accountproto.AccountV2_Role_Environment,
	isAdmin bool,
) (*eventproto.Editor, error) {
	if role < requiredRole {
		return nil, ErrPermissionDenied
	}
	return &eventproto.Editor{
		Email:   email,
		IsAdmin: isAdmin,
	}, nil
}

func CheckOrganizationRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Organization,
	getAccountFunc func(email string) (*accountproto.GetAccountV2Response, error),
) (*eventproto.Editor, error) {
	token, ok := rpc.GetIDToken(ctx)
	if !ok {
		return nil, ErrUnauthenticated
	}
	// TODO remove this condition after migration to AccountV2
	if token.IsAdmin() {
		return &eventproto.Editor{
			Email:   token.Email,
			IsAdmin: true,
		}, nil
	}
	resp, err := getAccountFunc(token.Email)
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			return nil, ErrUnauthenticated
		}
		return nil, ErrInternal
	}
	if resp.Account.OrganizationRole < requiredRole {
		return nil, ErrPermissionDenied
	}
	return &eventproto.Editor{
		Email:   token.Email,
		IsAdmin: token.IsAdmin(),
	}, nil
}
