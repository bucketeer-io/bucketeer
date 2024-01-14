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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

type AccountStorage interface {
	RunInTransaction(ctx context.Context, f func() error) error
	CreateAccount(ctx context.Context, a *domain.Account, environmentNamespace string) error
	UpdateAccount(ctx context.Context, a *domain.Account, environmentNamespace string) error
	GetAccount(ctx context.Context, id, environmentNamespace string) (*domain.Account, error)
	ListAccounts(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Account, int, int64, error)
	CreateAccountV2(ctx context.Context, a *domain.AccountV2) error
	UpdateAccountV2(ctx context.Context, a *domain.AccountV2) error
	DeleteAccountV2(ctx context.Context, a *domain.AccountV2) error
	GetAccountV2(ctx context.Context, email, organizationID string) (*domain.AccountV2, error)
	GetAccountV2ByEnvironmentID(ctx context.Context, email, environmentID string) (*domain.AccountV2, error)
	GetAccountsWithOrganization(ctx context.Context, email string) ([]*domain.AccountWithOrganization, error)
	ListAccountsV2(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.AccountV2, int, int64, error)
	CreateAdminAccount(ctx context.Context, a *domain.Account) error
	UpdateAdminAccount(ctx context.Context, a *domain.Account) error
	GetAdminAccount(ctx context.Context, id string) (*domain.Account, error)
	ListAdminAccounts(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Account, int, int64, error)
	CreateAPIKey(ctx context.Context, k *domain.APIKey, environmentNamespace string) error
	UpdateAPIKey(ctx context.Context, k *domain.APIKey, environmentNamespace string) error
	GetAPIKey(ctx context.Context, id, environmentNamespace string) (*domain.APIKey, error)
	ListAPIKeys(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.APIKey, int, int64, error)
}

type accountStorage struct {
	client mysql.Client
	tx     mysql.Transaction
}

func NewAccountStorage(client mysql.Client) AccountStorage {
	return &accountStorage{client, nil}
}

func (s *accountStorage) RunInTransaction(ctx context.Context, f func() error) error {
	tx, err := s.client.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("account: begin tx: %w", err)
	}
	s.tx = tx
	defer func() {
		s.tx = nil
	}()
	return s.client.RunInTransaction(ctx, tx, f)
}

func (s *accountStorage) qe() mysql.QueryExecer {
	if s.tx != nil {
		return s.tx
	}
	return s.client
}
