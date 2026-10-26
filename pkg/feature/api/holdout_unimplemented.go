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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// Holdout RPCs are part of the FeatureService contract. Handlers land in a
// later change; until then each method reports that it is not implemented.
func holdoutUnimplemented() error {
	return status.Error(codes.Unimplemented, "holdout group is not implemented")
}

func (*FeatureService) CreateHoldoutGroup(
	context.Context,
	*featureproto.CreateHoldoutGroupRequest,
) (*featureproto.CreateHoldoutGroupResponse, error) {
	return nil, holdoutUnimplemented()
}

func (*FeatureService) GetHoldoutGroup(
	context.Context,
	*featureproto.GetHoldoutGroupRequest,
) (*featureproto.GetHoldoutGroupResponse, error) {
	return nil, holdoutUnimplemented()
}

func (*FeatureService) ListHoldoutGroups(
	context.Context,
	*featureproto.ListHoldoutGroupsRequest,
) (*featureproto.ListHoldoutGroupsResponse, error) {
	return nil, holdoutUnimplemented()
}

func (*FeatureService) UpdateHoldoutGroup(
	context.Context,
	*featureproto.UpdateHoldoutGroupRequest,
) (*featureproto.UpdateHoldoutGroupResponse, error) {
	return nil, holdoutUnimplemented()
}

func (*FeatureService) StartHoldoutGroup(
	context.Context,
	*featureproto.StartHoldoutGroupRequest,
) (*featureproto.StartHoldoutGroupResponse, error) {
	return nil, holdoutUnimplemented()
}

func (*FeatureService) StartHoldoutGroupAnalysis(
	context.Context,
	*featureproto.StartHoldoutGroupAnalysisRequest,
) (*featureproto.StartHoldoutGroupAnalysisResponse, error) {
	return nil, holdoutUnimplemented()
}

func (*FeatureService) StopHoldoutGroup(
	context.Context,
	*featureproto.StopHoldoutGroupRequest,
) (*featureproto.StopHoldoutGroupResponse, error) {
	return nil, holdoutUnimplemented()
}

func (*FeatureService) ArchiveHoldoutGroup(
	context.Context,
	*featureproto.ArchiveHoldoutGroupRequest,
) (*featureproto.ArchiveHoldoutGroupResponse, error) {
	return nil, holdoutUnimplemented()
}

func (*FeatureService) DeleteHoldoutGroup(
	context.Context,
	*featureproto.DeleteHoldoutGroupRequest,
) (*featureproto.DeleteHoldoutGroupResponse, error) {
	return nil, holdoutUnimplemented()
}

func (*FeatureService) AddHoldoutGroupFeatures(
	context.Context,
	*featureproto.AddHoldoutGroupFeaturesRequest,
) (*featureproto.AddHoldoutGroupFeaturesResponse, error) {
	return nil, holdoutUnimplemented()
}

func (*FeatureService) RemoveHoldoutGroupFeatures(
	context.Context,
	*featureproto.RemoveHoldoutGroupFeaturesRequest,
) (*featureproto.RemoveHoldoutGroupFeaturesResponse, error) {
	return nil, holdoutUnimplemented()
}
