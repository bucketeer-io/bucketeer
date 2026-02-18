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
//

package cacher

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	ftcachermock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/cacher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
)

func TestSegmentUserCacherJobRun(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	internalErr := errors.New("internal error")

	patterns := []struct {
		desc        string
		setup       func(*segmentUserCacherJob)
		expectedErr error
	}{
		{
			desc: "err: cacher fails",
			setup: func(job *segmentUserCacherJob) {
				job.cacher.(*ftcachermock.MockSegmentUserCacher).EXPECT().
					RefreshAllEnvironmentCaches(gomock.Any()).
					Return(internalErr)
			},
			expectedErr: internalErr,
		},
		{
			desc: "success",
			setup: func(job *segmentUserCacherJob) {
				job.cacher.(*ftcachermock.MockSegmentUserCacher).EXPECT().
					RefreshAllEnvironmentCaches(gomock.Any()).
					Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			job := newSegmentUserCacherJobWithMock(t, controller)
			p.setup(job)
			err := job.Run(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newSegmentUserCacherJobWithMock(t *testing.T, controller *gomock.Controller) *segmentUserCacherJob {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &segmentUserCacherJob{
		cacher: ftcachermock.NewMockSegmentUserCacher(controller),
		logger: logger,
	}
}
