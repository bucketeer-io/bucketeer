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

package executor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	autoopsclientmock "github.com/bucketeer-io/bucketeer/pkg/autoops/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func TestExecuteProgressiveRollout(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*progressiveRolloutExecutor)
		expectedErr error
	}{
		{
			desc: "error: ExecuteProgressiveRollout fails",
			setup: func(e *progressiveRolloutExecutor) {
				e.autoOpsClient.(*autoopsclientmock.MockClient).EXPECT().ExecuteProgressiveRollout(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.Internal, "internal error"))
			},
			expectedErr: status.Errorf(codes.Internal, "internal error"),
		},
		{
			desc: "success",
			setup: func(e *progressiveRolloutExecutor) {
				e.autoOpsClient.(*autoopsclientmock.MockClient).EXPECT().ExecuteProgressiveRollout(gomock.Any(), gomock.Any()).Return(
					&autoopsproto.ExecuteProgressiveRolloutResponse{},
					nil,
				)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			e := newProgressiveRolloutExecutor(t, mockController)
			if p.setup != nil {
				p.setup(e)
			}
			err := e.ExecuteProgressiveRollout(context.Background(), "ns0", "rid1", "sid1")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newProgressiveRolloutExecutor(t *testing.T, mockController *gomock.Controller) *progressiveRolloutExecutor {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &progressiveRolloutExecutor{
		autoOpsClient: autoopsclientmock.NewMockClient(mockController),
		logger:        logger,
	}
}
