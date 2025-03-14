// Copyright 2025 The Bucketeer Authors.
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
	"errors"

	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func (s *FeatureService) ScheduleFlagChange(
	ctx context.Context,
	req *ftproto.ScheduleFlagChangeRequest,
) (*ftproto.ScheduleFlagChangeResponse, error) {
	// TODO
	return nil, errors.New("api not yet implemented")
}

func (s *FeatureService) UpdateScheduledFlagChange(
	ctx context.Context,
	req *ftproto.UpdateScheduledFlagChangeRequest,
) (*ftproto.UpdateScheduledFlagChangeResponse, error) {
	// TODO
	return nil, errors.New("api not yet implemented")
}

func (s *FeatureService) DeleteScheduledFlagChange(
	ctx context.Context,
	req *ftproto.DeleteScheduledFlagChangeRequest,
) (*ftproto.DeleteScheduledFlagChangeResponse, error) {
	// TODO
	return nil, errors.New("api not yet implemented")
}

func (s *FeatureService) ListScheduledFlagChanges(
	ctx context.Context,
	req *ftproto.ListScheduledFlagChangesRequest,
) (*ftproto.ListScheduledFlagChangesResponse, error) {
	// TODO
	return nil, errors.New("api not yet implemented")
}
