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

package api

import (
	"context"
	"sort"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func (s *FeatureService) GetUserAttributeKeys(
	ctx context.Context,
	req *featureproto.GetUserAttributeKeysRequest,
) (*featureproto.GetUserAttributeKeysResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		s.logger.Error("Failed to get user attribute keys", zap.Error(err))
		return nil, err
	}

	userAttributeKeys, err := s.userAttributesCache.GetUserAttributeKeyAll(req.EnvironmentId)
	sort.Strings(userAttributeKeys)
	if err != nil {
		s.logger.Error(
			"Failed to get user attribute keys",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	return &featureproto.GetUserAttributeKeysResponse{
		UserAttributeKeys: userAttributeKeys,
	}, nil
}
