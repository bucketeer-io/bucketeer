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

	"github.com/bucketeer-io/bucketeer/pkg/coderef/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

type CodeReferenceStorage interface {
	RunInTransaction(ctx context.Context, f func() error) error
	CreateCodeReference(ctx context.Context, codeRef *domain.CodeReference) error
	UpdateCodeReference(ctx context.Context, codeRef *domain.CodeReference) error
	GetCodeReference(ctx context.Context, id, environmentID string) (*domain.CodeReference, error)
	ListCodeReferences(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*domain.CodeReference, int, int64, error)
	DeleteCodeReference(ctx context.Context, id, environmentID string) error
}

const transactionKey = "transaction"