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

	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func (s *AutoOpsService) CreateFlagTrigger(ctx context.Context,
	request *autoopsproto.CreateFlagTriggerRequest,
) (*autoopsproto.CreateFlagTriggerResponse, error) {
	//TODO implement me
	return &autoopsproto.CreateFlagTriggerResponse{}, nil
}

func (s *AutoOpsService) ResetFlagTrigger(
	ctx context.Context,
	request *autoopsproto.ResetFlagTriggerRequest,
) (*autoopsproto.ResetFlagTriggerResponse, error) {
	//TODO implement me
	return &autoopsproto.ResetFlagTriggerResponse{}, nil
}

func (s *AutoOpsService) DeleteFlagTrigger(
	ctx context.Context,
	request *autoopsproto.DeleteFlagTriggerRequest,
) (*autoopsproto.DeleteFlagTriggerResponse, error) {
	//TODO implement me
	return &autoopsproto.DeleteFlagTriggerResponse{}, nil
}

func (s *AutoOpsService) GetFlagTrigger(
	ctx context.Context,
	request *autoopsproto.GetFlagTriggerRequest,
) (*autoopsproto.GetFlagTriggerResponse, error) {
	//TODO implement me
	return &autoopsproto.GetFlagTriggerResponse{}, nil
}

func (s *AutoOpsService) ListFlagTriggers(
	ctx context.Context,
	request *autoopsproto.ListFlagTriggersRequest,
) (*autoopsproto.ListFlagTriggersResponse, error) {
	//TODO implement me
	return &autoopsproto.ListFlagTriggersResponse{}, nil
}

func (s *AutoOpsService) UpdateFlagTrigger(
	ctx context.Context,
	request *autoopsproto.UpdateFlagTriggerRequest,
) (*autoopsproto.UpdateFlagTriggerResponse, error) {
	return &autoopsproto.UpdateFlagTriggerResponse{}, nil
}

func (s *AutoOpsService) EnableFlagTrigger(
	ctx context.Context,
	request *autoopsproto.EnableFlagTriggerRequest,
) (*autoopsproto.EnableFlagTriggerResponse, error) {
	//TODO implement me
	return &autoopsproto.EnableFlagTriggerResponse{}, nil
}

func (s *AutoOpsService) DisableFlagTrigger(
	ctx context.Context,
	request *autoopsproto.DisableFlagTriggerRequest,
) (*autoopsproto.DisableFlagTriggerResponse, error) {
	//TODO implement me
	return &autoopsproto.DisableFlagTriggerResponse{}, nil
}
