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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

type AccountStorage interface {
	CreateAccountV2(ctx context.Context, a *domain.AccountV2) error
	UpdateAccountV2(ctx context.Context, a *domain.AccountV2) error
	DeleteAccountV2(ctx context.Context, a *domain.AccountV2) error
	GetAccountV2(ctx context.Context, email, organizationID string) (*domain.AccountV2, error)
	GetAccountV2ByEnvironmentID(ctx context.Context, email, environmentID string) (*domain.AccountV2, error)
	GetAvatarAccountsV2(ctx context.Context, options *mysql.ListOptions) ([]*proto.AccountV2, error)
	GetAccountsWithOrganization(ctx context.Context, email string) ([]*domain.AccountWithOrganization, error)
	ListAccountsV2(
		ctx context.Context,
		options *mysql.ListOptions,
	) ([]*proto.AccountV2, int, int64, error)
	GetSystemAdminAccountV2(ctx context.Context, email string) (*domain.AccountV2, error)
	CreateAPIKey(ctx context.Context, k *domain.APIKey, environmentID string) error
	UpdateAPIKey(ctx context.Context, k *domain.APIKey, environmentID string) error
	UpdateAPIKeyLastUsedAt(ctx context.Context, id, environmentID string, lastUsedAt int64) (bool, error)
	GetAPIKey(ctx context.Context, id, environmentID string) (*domain.APIKey, error)
	GetAPIKeyByAPIKey(ctx context.Context, apiKey string, environmentID string) (*domain.APIKey, error)
	GetEnvironmentAPIKey(ctx context.Context, apiKey string) (*domain.EnvironmentAPIKey, error)
	ListAllEnvironmentAPIKeys(ctx context.Context) ([]*domain.EnvironmentAPIKey, error)
	ListAPIKeys(ctx context.Context, options *mysql.ListOptions) ([]*proto.APIKey, int, int64, error)
}

type accountStorage struct {
	qe mysql.QueryExecer
}

func NewAccountStorage(qe mysql.QueryExecer) AccountStorage {
	return &accountStorage{qe}
}
