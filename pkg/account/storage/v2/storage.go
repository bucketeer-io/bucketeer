// Copyright 2024 The Bucketeer Authors.
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
	GetSystemAdminAccountV2(ctx context.Context, email string) (*domain.AccountV2, error)
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

const transactionKey = "transaction"

type accountStorage struct {
	client mysql.Client
}

func NewAccountStorage(client mysql.Client) AccountStorage {
	return &accountStorage{client}
}

func (s *accountStorage) RunInTransaction(ctx context.Context, f func() error) error {
	tx, err := s.client.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("account: begin tx: %w", err)
	}
	ctx = context.WithValue(ctx, transactionKey, tx)
	return s.client.RunInTransaction(ctx, tx, f)
}

func (s *accountStorage) qe(ctx context.Context) mysql.QueryExecer {
	tx, ok := ctx.Value(transactionKey).(mysql.Transaction)
	if ok {
		return tx
	}
	return s.client
}
