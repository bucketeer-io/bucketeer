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
//

package api

import (
	"context"

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func (s *FeatureService) CreateFlagTrigger(
	ctx context.Context,
	request *featureproto.CreateFlagTriggerRequest,
) (*featureproto.CreateFlagTriggerResponse, error) {
	return &featureproto.CreateFlagTriggerResponse{}, nil
}

func (s *FeatureService) UpdateFlagTrigger(
	ctx context.Context,
	request *featureproto.UpdateFlagTriggerRequest,
) (*featureproto.UpdateFlagTriggerResponse, error) {
	return &featureproto.UpdateFlagTriggerResponse{}, nil
}

func (s *FeatureService) EnableFlagTrigger(
	ctx context.Context,
	request *featureproto.EnableFlagTriggerRequest,
) (*featureproto.EnableFlagTriggerResponse, error) {
	return &featureproto.EnableFlagTriggerResponse{}, nil
}

func (s *FeatureService) DisableFlagTrigger(
	ctx context.Context,
	request *featureproto.DisableFlagTriggerRequest,
) (*featureproto.DisableFlagTriggerResponse, error) {
	return &featureproto.DisableFlagTriggerResponse{}, nil
}

func (s *FeatureService) ResetFlagTrigger(
	ctx context.Context,
	request *featureproto.ResetFlagTriggerRequest,
) (*featureproto.ResetFlagTriggerResponse, error) {
	return &featureproto.ResetFlagTriggerResponse{}, nil
}

func (s *FeatureService) DeleteFlagTrigger(
	ctx context.Context,
	request *featureproto.DeleteFlagTriggerRequest,
) (*featureproto.DeleteFlagTriggerResponse, error) {
	return &featureproto.DeleteFlagTriggerResponse{}, nil
}

func (s *FeatureService) GetFlagTrigger(
	ctx context.Context,
	request *featureproto.GetFlagTriggerRequest,
) (*featureproto.GetFlagTriggerResponse, error) {
	return &featureproto.GetFlagTriggerResponse{}, nil
}

func (s *FeatureService) ListFlagTriggers(
	ctx context.Context,
	request *featureproto.ListFlagTriggersRequest,
) (*featureproto.ListFlagTriggersResponse, error) {
	return &featureproto.ListFlagTriggersResponse{}, nil
}
