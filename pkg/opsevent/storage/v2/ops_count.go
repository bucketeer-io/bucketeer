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

	"github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

// Shared list-query errors returned by OpsCountStorage implementations.
var (
	ErrInvalidCursor = errors.New("opsevent/storage/v2: invalid cursor")
)

type OpsCountStorage interface {
	UpsertOpsCount(ctx context.Context, environmentId string, oc *domain.OpsCount) error
	ListOpsCounts(
		ctx context.Context,
		params ListOpsCountsParams,
	) ([]*proto.OpsCount, int, error)
}

type ListOpsCountsParams struct {
	EnvironmentID  string
	FeatureIDs     []string
	AutoOpsRuleIDs []string
	PageSize       int
	Cursor         string
}
