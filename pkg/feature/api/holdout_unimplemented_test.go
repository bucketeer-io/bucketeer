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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestHoldoutRPCsUnimplemented(t *testing.T) {
	t.Parallel()
	service := &FeatureService{}
	ctx := context.Background()
	calls := []func() error{
		func() error {
			_, err := service.CreateHoldoutGroup(ctx, &featureproto.CreateHoldoutGroupRequest{})
			return err
		},
		func() error {
			_, err := service.GetHoldoutGroup(ctx, &featureproto.GetHoldoutGroupRequest{})
			return err
		},
		func() error {
			_, err := service.ListHoldoutGroups(ctx, &featureproto.ListHoldoutGroupsRequest{})
			return err
		},
		func() error {
			_, err := service.UpdateHoldoutGroup(ctx, &featureproto.UpdateHoldoutGroupRequest{})
			return err
		},
		func() error {
			_, err := service.StartHoldoutGroup(ctx, &featureproto.StartHoldoutGroupRequest{})
			return err
		},
		func() error {
			_, err := service.StartHoldoutGroupAnalysis(ctx, &featureproto.StartHoldoutGroupAnalysisRequest{})
			return err
		},
		func() error {
			_, err := service.StopHoldoutGroup(ctx, &featureproto.StopHoldoutGroupRequest{})
			return err
		},
		func() error {
			_, err := service.ArchiveHoldoutGroup(ctx, &featureproto.ArchiveHoldoutGroupRequest{})
			return err
		},
		func() error {
			_, err := service.DeleteHoldoutGroup(ctx, &featureproto.DeleteHoldoutGroupRequest{})
			return err
		},
		func() error {
			_, err := service.AddHoldoutGroupFeatures(ctx, &featureproto.AddHoldoutGroupFeaturesRequest{})
			return err
		},
		func() error {
			_, err := service.RemoveHoldoutGroupFeatures(ctx, &featureproto.RemoveHoldoutGroupFeaturesRequest{})
			return err
		},
	}
	for _, call := range calls {
		err := call()
		require.Error(t, err)
		assert.Equal(t, codes.Unimplemented, status.Code(err))
	}
}
