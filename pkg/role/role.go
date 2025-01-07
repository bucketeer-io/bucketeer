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

package role

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	accdomain "github.com/bucketeer-io/bucketeer/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

const (
	APIKeyTokenMDKey      = "apikey-token"
	APIKeyMaintainerMDKey = "apikey-maintainer"
	APIKeyNameMDKey       = "apikey-name"
)

var (
	ErrUnauthenticated  = status.Error(codes.Unauthenticated, "Unauthenticated user")
	ErrPermissionDenied = status.Error(codes.PermissionDenied, "Permission denied")
	ErrInternal         = status.Error(codes.Internal, "Internal")
)

func CheckSystemAdminRole(ctx context.Context) (*eventproto.Editor, error) {
	token, ok := rpc.GetAccessToken(ctx)
	if !ok {
		return nil, ErrUnauthenticated
	}
	if !token.IsSystemAdmin {
		return nil, ErrPermissionDenied
	}
	return &eventproto.Editor{
		Email:   token.Email,
		IsAdmin: true,
	}, nil
}

func CheckEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentID string,
	getAccountFunc func(email string) (*accountproto.AccountV2, error),
) (*eventproto.Editor, error) {
	token, ok := rpc.GetAccessToken(ctx)
	if !ok {
		return nil, ErrUnauthenticated
	}
	publicAPIEditor := getAPIKeyEditor(ctx)
	if publicAPIEditor != nil && publicAPIEditor.Token != "" && token.IsSystemAdmin {
		var accountName string
		resp, err := getAccountFunc(publicAPIEditor.Maintainer)
		if err == nil && resp != nil {
			account := accdomain.AccountV2{AccountV2: resp}
			accountName = account.GetAccountFullName()
		}
		return &eventproto.Editor{
			Email:           publicAPIEditor.Maintainer,
			Name:            accountName,
			PublicApiEditor: publicAPIEditor,
		}, nil
	}

	if token.IsSystemAdmin {
		return checkRole(token.Email, token.Name, accountproto.AccountV2_Role_Environment_EDITOR, requiredRole, true)
	}
	// get account for the environment namespace
	account, err := getAccountFunc(token.Email)
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			return nil, ErrUnauthenticated
		}
		return nil, ErrInternal
	}
	if account.Disabled {
		return nil, ErrUnauthenticated
	}
	accountEnvRole := getRole(account.EnvironmentRoles, environmentID)
	return checkRole(account.Email, token.Name, accountEnvRole, requiredRole, false)
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
	name string,
	role, requiredRole accountproto.AccountV2_Role_Environment,
	isAdmin bool,
) (*eventproto.Editor, error) {
	if role < requiredRole {
		return nil, ErrPermissionDenied
	}
	return &eventproto.Editor{
		Email:   email,
		Name:    name,
		IsAdmin: isAdmin,
	}, nil
}

func CheckOrganizationRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Organization,
	getAccountFunc func(email string) (*accountproto.GetAccountV2Response, error),
) (*eventproto.Editor, error) {
	token, ok := rpc.GetAccessToken(ctx)
	if !ok {
		return nil, ErrUnauthenticated
	}
	publicAPIEditor := getAPIKeyEditor(ctx)
	if publicAPIEditor != nil && publicAPIEditor.Token != "" && token.IsSystemAdmin {
		resp, err := getAccountFunc(publicAPIEditor.Maintainer)
		if err != nil {
			return nil, err
		}
		account := accdomain.AccountV2{AccountV2: resp.Account}
		return &eventproto.Editor{
			Email:           publicAPIEditor.Maintainer,
			Name:            account.GetAccountFullName(),
			PublicApiEditor: publicAPIEditor,
		}, nil
	}

	if token.IsSystemAdmin {
		return &eventproto.Editor{
			Email:   token.Email,
			Name:    token.Name,
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
	if resp.Account.Disabled {
		return nil, ErrUnauthenticated
	}
	if resp.Account.OrganizationRole < requiredRole {
		return nil, ErrPermissionDenied
	}
	return &eventproto.Editor{
		Email:   token.Email,
		Name:    token.Name,
		IsAdmin: token.IsSystemAdmin,
	}, nil
}

func getAPIKeyEditor(ctx context.Context) *eventproto.Editor_PublicAPIEditor {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	apikeyToken := md.Get(APIKeyTokenMDKey)
	if len(apikeyToken) == 0 {
		return nil
	}

	publicAPIEditor := &eventproto.Editor_PublicAPIEditor{}
	publicAPIEditor.Token = apikeyToken[0]

	if len(md.Get(APIKeyMaintainerMDKey)) > 0 {
		publicAPIEditor.Maintainer = md.Get(APIKeyMaintainerMDKey)[0]
	}

	if len(md.Get(APIKeyNameMDKey)) > 0 {
		publicAPIEditor.Name = md.Get(APIKeyNameMDKey)[0]
	}
	return publicAPIEditor
}
