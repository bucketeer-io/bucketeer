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
//

package cacher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	mockcachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	expproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

func TestExperimentsPutCache(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	envID := "env-id"
	experiments := &expproto.Experiments{
		Experiments: []*expproto.Experiment{
			{
				Id: "exp-id-1",
			},
			{
				Id: "exp-id-2",
			},
		},
	}

	patterns := []struct {
		desc     string
		setup    func(*experimentCacher)
		expected int
	}{
		{
			desc: "err: error at index 0",
			setup: func(fc *experimentCacher) {
				fc.caches[0].(*mockcachev3.MockExperimentsCache).EXPECT().Put(experiments, envID).
					Return(errors.New("internal error"))
				fc.caches[1].(*mockcachev3.MockExperimentsCache).EXPECT().Put(experiments, envID).
					Return(nil)
			},
			expected: 1,
		},
		{
			desc: "err: error at index 1",
			setup: func(fc *experimentCacher) {
				fc.caches[0].(*mockcachev3.MockExperimentsCache).EXPECT().Put(experiments, envID).
					Return(nil)
				fc.caches[1].(*mockcachev3.MockExperimentsCache).EXPECT().Put(experiments, envID).
					Return(errors.New("internal error"))
			},
			expected: 1,
		},
		{
			desc: "success",
			setup: func(fc *experimentCacher) {
				fc.caches[0].(*mockcachev3.MockExperimentsCache).EXPECT().Put(experiments, envID).
					Return(nil)
				fc.caches[1].(*mockcachev3.MockExperimentsCache).EXPECT().Put(experiments, envID).
					Return(nil)
			},
			expected: 2,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newExperimentsFlagCacher(t, controller)
			p.setup(cacher)
			updatedInstances := cacher.putCache(experiments, envID)
			assert.Equal(t, p.expected, updatedInstances)
		})
	}
}

func newExperimentsFlagCacher(t *testing.T, controller *gomock.Controller) *experimentCacher {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &experimentCacher{
		caches: []cachev3.ExperimentsCache{
			mockcachev3.NewMockExperimentsCache(controller),
			mockcachev3.NewMockExperimentsCache(controller),
		},
		logger: logger,
	}
}
