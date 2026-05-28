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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	"errors"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

// nolint:lll
var (
	ErrAccountAlreadyExists          = pkgErr.NewErrorAlreadyExists(pkgErr.AccountPackageName, "account already exists")
	ErrAccountNotFound               = pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "account not found", "account")
	ErrAccountUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(pkgErr.AccountPackageName, " unexpected affected rows")
	ErrSystemAdminAccountNotFound    = pkgErr.NewErrorNotFound(
		pkgErr.AccountPackageName,
		"admin account not found",
		"admin_account",
	)
	ErrAPIKeyAlreadyExists          = pkgErr.NewErrorAlreadyExists(pkgErr.AccountPackageName, "api key already exists")
	ErrAPIKeyNotFound               = pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "api key not found", "api_key")
	ErrAPIKeyUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.AccountPackageName,
		"api key unexpected affected rows",
	)
)

var (
	ErrInvalidOrderBy = errors.New("account storage: invalid order by")
	ErrInvalidCursor  = errors.New("account storage: invalid cursor")
)

type AccountStorage interface {
	CreateAccountV2(ctx context.Context, a *domain.AccountV2) error
	UpdateAccountV2(ctx context.Context, a *domain.AccountV2) error
	DeleteAccountV2(ctx context.Context, a *domain.AccountV2) error
	GetAccountV2(ctx context.Context, email, organizationID string) (*domain.AccountV2, error)
	GetAccountV2ByEnvironmentID(ctx context.Context, email, environmentID string) (*domain.AccountV2, error)
	GetSystemAdminAccountV2(ctx context.Context, email string) (*domain.AccountV2, error)
	GetAccountsWithOrganization(ctx context.Context, email string) ([]*domain.AccountWithOrganization, error)
	GetAvatarAccountsV2(ctx context.Context, params GetAvatarAccountsV2Params) ([]*proto.AccountV2, error)
	ListAccountsV2(ctx context.Context, params ListAccountsV2Params) ([]*proto.AccountV2, int, int64, error)
	CreateAPIKey(ctx context.Context, k *domain.APIKey, environmentID string) error
	UpdateAPIKey(ctx context.Context, k *domain.APIKey, environmentID string) error
	UpdateAPIKeyLastUsedAt(ctx context.Context, id, environmentID string, lastUsedAt int64) (bool, error)
	GetAPIKey(ctx context.Context, id, environmentID string) (*domain.APIKey, error)
	GetAPIKeyByAPIKey(ctx context.Context, apiKey string, environmentID string) (*domain.APIKey, error)
	GetEnvironmentAPIKey(ctx context.Context, apiKey string) (*domain.EnvironmentAPIKey, error)
	ListAllEnvironmentAPIKeys(ctx context.Context) ([]*domain.EnvironmentAPIKey, error)
	ListAPIKeys(ctx context.Context, params ListAPIKeysParams) ([]*proto.APIKey, int, int64, error)
}

type GetAvatarAccountsV2Params struct {
	Emails        []string
	EnvironmentID string
}

type ListAccountsV2Params struct {
	OrganizationID   string
	Disabled         *bool
	Tags             []string
	Teams            []string
	OrganizationRole *int32
	EnvironmentID    *string
	EnvironmentRole  *int32
	// EnvironmentRoles is used for members who can only see accounts in their environments.
	// When set, generates OR conditions: (env_roles contains role1) OR (env_roles contains role2) OR (org_role >= admin)
	EnvironmentRoles []*proto.AccountV2_EnvironmentRole
	SearchKeyword    string
	OrderBy          proto.ListAccountsV2Request_OrderBy
	OrderDirection   proto.ListAccountsV2Request_OrderDirection
	PageSize         int
	Cursor           string
}

type ListAPIKeysParams struct {
	OrganizationID  string
	EnvironmentIDs  []string
	Disabled        *bool
	MaintainerEmail string
	SearchKeyword   string
	OrderBy         proto.ListAPIKeysRequest_OrderBy
	OrderDirection  proto.ListAPIKeysRequest_OrderDirection
	PageSize        int
	Cursor          string
}
