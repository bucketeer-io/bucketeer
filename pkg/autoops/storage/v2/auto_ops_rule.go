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

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

var (
	ErrAutoOpsRuleAlreadyExists          = err.NewErrorAlreadyExists(err.AutoopsPackageName, "already exists")
	ErrAutoOpsRuleNotFound               = err.NewErrorNotFound(err.AutoopsPackageName, "not found", "autoOpsRule")
	ErrAutoOpsRuleUnexpectedAffectedRows = err.NewErrorUnexpectedAffectedRows(
		err.AutoopsPackageName,
		"unexpected affected rows",
	)
)

// Shared list-query errors returned by AutoOpsRuleStorage implementations.
var ErrInvalidCursor = errors.New("autoops/storage/v2: invalid cursor")

type AutoOpsRuleStorage interface {
	CreateAutoOpsRule(ctx context.Context, e *domain.AutoOpsRule, environmentId string) error
	UpdateAutoOpsRule(ctx context.Context, e *domain.AutoOpsRule, environmentId string) error
	GetAutoOpsRule(ctx context.Context, id, environmentId string) (*domain.AutoOpsRule, error)
	ListAutoOpsRules(
		ctx context.Context,
		params ListAutoOpsRulesParams,
	) ([]*proto.AutoOpsRule, int, error)
}

type ListAutoOpsRulesParams struct {
	EnvironmentID string
	FeatureIDs    []string
	PageSize      int
	Cursor        string
}
