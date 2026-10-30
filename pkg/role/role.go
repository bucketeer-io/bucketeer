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
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	accdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
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
	// Service tokens (machine accounts used by internal batch/subscriber jobs)
	// bypass all role checks and retain full access. Both flags are required so a
	// token accidentally minted with is_service_token=true cannot escalate to
	// admin without also being a system admin.
	if token.IsSystemAdmin && token.IsServiceToken {
		// The public API gateway relays requests with a service token and
		// identifies the API key in the metadata. Those requests act as the
		// key's maintainer, whose membership and role are checked like any
		// other account, instead of inheriting the service token's access.
		if publicAPIEditor := getAPIKeyEditor(ctx); publicAPIEditor != nil && publicAPIEditor.Token != "" {
			return checkAPIKeyMaintainerEnvironmentRole(publicAPIEditor, requiredRole, environmentID, getAccountFunc)
		}
		return &eventproto.Editor{
			Email:   token.Email,
			Name:    token.Name,
			IsAdmin: true,
		}, nil
	}
	// API key metadata is only trusted on service token requests because only
	// the public API gateway sets it. On any other request it is ignored and
	// the caller is authorized as the account in its own token.

	// System admins can read from any org without membership, but must be an org
	// member with the appropriate role to write.
	if token.IsSystemAdmin && requiredRole <= accountproto.AccountV2_Role_Environment_VIEWER {
		return &eventproto.Editor{
			Email:   token.Email,
			Name:    token.Name,
			IsAdmin: true,
		}, nil
	}
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
	return checkRole(
		account.Email,
		token.Name,
		accountEnvRole,
		requiredRole,
		account.EnvironmentRoles,
		account.OrganizationRole,
		token.IsSystemAdmin,
	)
}

func CheckEnvironmentRoleWithLog(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentID string,
	getAccountFunc func(email string) (*accountproto.AccountV2, error),
	logger *zap.Logger,
	unauthenticatedErr error,
	permissionDeniedErr error,
	defaultErrFunc func(error) error,
) (*eventproto.Editor, error) {
	editor, err := CheckEnvironmentRole(ctx, requiredRole, environmentID, getAccountFunc)
	if err != nil {
		logFields := []zap.Field{
			zap.Error(err),
			zap.String("environmentId", environmentID),
			zap.String("requiredRole", requiredRole.String()),
		}
		if token, ok := rpc.GetAccessToken(ctx); ok {
			logFields = append(logFields, zap.String("email", token.Email))
		}
		fields := log.FieldsFromIncomingContext(ctx).AddFields(logFields...)
		switch status.Code(err) {
		case codes.Unauthenticated:
			logger.Error("Unauthenticated", fields...)
			return nil, unauthenticatedErr
		case codes.PermissionDenied:
			logger.Error("Permission denied", fields...)
			return nil, permissionDeniedErr
		default:
			logger.Error("Failed to check role", fields...)
			return nil, defaultErrFunc(err)
		}
	}
	return editor, nil
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
	environmentRoles []*accountproto.AccountV2_EnvironmentRole,
	organizationRole accountproto.AccountV2_Role_Organization,
	isAdmin bool,
) (*eventproto.Editor, error) {
	if organizationRole >= accountproto.AccountV2_Role_Organization_ADMIN {
		return &eventproto.Editor{
			Email:            email,
			Name:             name,
			IsAdmin:          isAdmin,
			EnvironmentRoles: environmentRoles,
			OrganizationRole: organizationRole,
		}, nil
	}
	if role < requiredRole {
		return nil, ErrPermissionDenied
	}
	return &eventproto.Editor{
		Email:            email,
		Name:             name,
		IsAdmin:          isAdmin,
		EnvironmentRoles: environmentRoles,
		OrganizationRole: organizationRole,
	}, nil
}

// checkAPIKeyMaintainerEnvironmentRole authorizes a request relayed by the
// public API gateway as the maintainer of the API key, with the same
// membership and role checks that apply to the maintainer's own requests.
func checkAPIKeyMaintainerEnvironmentRole(
	publicAPIEditor *eventproto.Editor_PublicAPIEditor,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentID string,
	getAccountFunc func(email string) (*accountproto.AccountV2, error),
) (*eventproto.Editor, error) {
	resp, err := getAccountFunc(publicAPIEditor.Maintainer)
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			return nil, ErrUnauthenticated
		}
		return nil, ErrInternal
	}
	if resp.Disabled {
		return nil, ErrUnauthenticated
	}
	account := accdomain.AccountV2{AccountV2: resp}
	editor, err := checkRole(
		publicAPIEditor.Maintainer,
		account.GetAccountFullName(),
		getRole(account.EnvironmentRoles, environmentID),
		requiredRole,
		account.EnvironmentRoles,
		account.OrganizationRole,
		false,
	)
	if err != nil {
		return nil, err
	}
	editor.PublicApiEditor = publicAPIEditor
	return editor, nil
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
	// Service tokens (machine accounts used by internal batch/subscriber jobs)
	// bypass all role checks and retain full access. Both flags are required so a
	// token accidentally minted with is_service_token=true cannot escalate to
	// admin without also being a system admin.
	if token.IsSystemAdmin && token.IsServiceToken {
		// See CheckEnvironmentRole: requests relayed by the public API gateway
		// act as the API key's maintainer and are checked against its role.
		if publicAPIEditor := getAPIKeyEditor(ctx); publicAPIEditor != nil && publicAPIEditor.Token != "" {
			return checkAPIKeyMaintainerOrganizationRole(publicAPIEditor, requiredRole, getAccountFunc)
		}
		return &eventproto.Editor{
			Email:   token.Email,
			Name:    token.Name,
			IsAdmin: true,
		}, nil
	}
	// API key metadata is only trusted on service token requests because only
	// the public API gateway sets it. On any other request it is ignored and
	// the caller is authorized as the account in its own token.

	// System admins can view any org without membership, but must be a member
	// with the appropriate role for admin/write operations (e.g. creating API keys).
	if token.IsSystemAdmin && requiredRole <= accountproto.AccountV2_Role_Organization_MEMBER {
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
		Email:            token.Email,
		Name:             token.Name,
		IsAdmin:          token.IsSystemAdmin,
		EnvironmentRoles: resp.Account.EnvironmentRoles,
		OrganizationRole: resp.Account.OrganizationRole,
	}, nil
}

// checkAPIKeyMaintainerOrganizationRole authorizes a request relayed by the
// public API gateway as the maintainer of the API key, with the same
// membership and role checks that apply to the maintainer's own requests.
func checkAPIKeyMaintainerOrganizationRole(
	publicAPIEditor *eventproto.Editor_PublicAPIEditor,
	requiredRole accountproto.AccountV2_Role_Organization,
	getAccountFunc func(email string) (*accountproto.GetAccountV2Response, error),
) (*eventproto.Editor, error) {
	resp, err := getAccountFunc(publicAPIEditor.Maintainer)
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
	account := accdomain.AccountV2{AccountV2: resp.Account}
	return &eventproto.Editor{
		Email:            publicAPIEditor.Maintainer,
		Name:             account.GetAccountFullName(),
		PublicApiEditor:  publicAPIEditor,
		EnvironmentRoles: account.EnvironmentRoles,
		OrganizationRole: account.OrganizationRole,
	}, nil
}

func CheckOrganizationRoleWithLog(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Organization,
	organizationID string,
	getAccountFunc func(email string) (*accountproto.GetAccountV2Response, error),
	logger *zap.Logger,
	unauthenticatedErr error,
	permissionDeniedErr error,
	defaultErrFunc func(error) error,
) (*eventproto.Editor, error) {
	return checkOrganizationRoleWithLog(
		ctx,
		requiredRole,
		zap.String("organizationId", organizationID),
		getAccountFunc,
		logger,
		unauthenticatedErr,
		permissionDeniedErr,
		defaultErrFunc,
	)
}

func CheckOrganizationRoleByEnvironmentIDWithLog(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Organization,
	environmentID string,
	getAccountFunc func(email string) (*accountproto.GetAccountV2Response, error),
	logger *zap.Logger,
	unauthenticatedErr error,
	permissionDeniedErr error,
	defaultErrFunc func(error) error,
) (*eventproto.Editor, error) {
	return checkOrganizationRoleWithLog(
		ctx,
		requiredRole,
		zap.String("environmentId", environmentID),
		getAccountFunc,
		logger,
		unauthenticatedErr,
		permissionDeniedErr,
		defaultErrFunc,
	)
}

func checkOrganizationRoleWithLog(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Organization,
	resourceField zap.Field,
	getAccountFunc func(email string) (*accountproto.GetAccountV2Response, error),
	logger *zap.Logger,
	unauthenticatedErr error,
	permissionDeniedErr error,
	defaultErrFunc func(error) error,
) (*eventproto.Editor, error) {
	editor, err := CheckOrganizationRole(ctx, requiredRole, getAccountFunc)
	if err != nil {
		logFields := []zap.Field{
			zap.Error(err),
			resourceField,
			zap.String("requiredRole", requiredRole.String()),
		}
		if token, ok := rpc.GetAccessToken(ctx); ok {
			logFields = append(logFields, zap.String("email", token.Email))
		}
		fields := log.FieldsFromIncomingContext(ctx).AddFields(logFields...)
		switch status.Code(err) {
		case codes.Unauthenticated:
			logger.Error("Unauthenticated", fields...)
			return nil, unauthenticatedErr
		case codes.PermissionDenied:
			logger.Error("Permission denied", fields...)
			return nil, permissionDeniedErr
		default:
			logger.Error("Failed to check role", fields...)
			return nil, defaultErrFunc(err)
		}
	}
	return editor, nil
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
	// Only the obfuscated key may reach the domain events and the audit log.
	publicAPIEditor.Token = accdomain.ObfuscateAPIKey(apikeyToken[0])

	if len(md.Get(APIKeyMaintainerMDKey)) > 0 {
		publicAPIEditor.Maintainer = md.Get(APIKeyMaintainerMDKey)[0]
	}

	if len(md.Get(APIKeyNameMDKey)) > 0 {
		publicAPIEditor.Name = md.Get(APIKeyNameMDKey)[0]
	}
	return publicAPIEditor
}
