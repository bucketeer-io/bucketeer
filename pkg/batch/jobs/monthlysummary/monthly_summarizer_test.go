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

package monthlysummary

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	envclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	insightsstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/insights/storage/v2/mock"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

func TestMonthlySummarizerRun(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	patterns := []struct {
		desc     string
		setup    func(*monthlySummarizer)
		expected error
	}{
		{
			desc: "fail: list environments error",
			setup: func(m *monthlySummarizer) {
				m.envClient.(*envclientmock.MockClient).EXPECT().
					ListEnvironmentsV2(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("list error"))
			},
			expected: errors.New("list error"),
		},
		{
			desc: "success: no records when all merge fail",
			setup: func(m *monthlySummarizer) {
				m.envClient.(*envclientmock.MockClient).EXPECT().
					ListEnvironmentsV2(gomock.Any(), gomock.Any()).
					Return(&envproto.ListEnvironmentsV2Response{
						Environments: []*envproto.EnvironmentV2{
							{Id: "env1"},
						},
					}, nil)
				m.mauCache.(*cachemock.MockMAUCache).EXPECT().
					MergeIntoMAUBatch(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("merge error"))
			},
			expected: nil,
		},
		{
			desc: "fail: upsert error",
			setup: func(m *monthlySummarizer) {
				m.envClient.(*envclientmock.MockClient).EXPECT().
					ListEnvironmentsV2(gomock.Any(), gomock.Any()).
					Return(&envproto.ListEnvironmentsV2Response{
						Environments: []*envproto.EnvironmentV2{
							{Id: "env1"},
						},
					}, nil)
				m.mauCache.(*cachemock.MockMAUCache).EXPECT().
					MergeIntoMAUBatch(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(map[string]int64{"ANDROID": 100}, nil)
				m.monthlySummaryStorage.(*insightsstoragemock.MockMonthlySummaryStorage).EXPECT().
					UpsertMonthlySummaryBatch(gomock.Any(), gomock.Any()).
					Return(errors.New("upsert error"))
			},
			expected: errors.New("upsert error"),
		},
		{
			desc: "success",
			setup: func(m *monthlySummarizer) {
				m.envClient.(*envclientmock.MockClient).EXPECT().
					ListEnvironmentsV2(gomock.Any(), gomock.Any()).
					Return(&envproto.ListEnvironmentsV2Response{
						Environments: []*envproto.EnvironmentV2{
							{Id: "env1"},
						},
					}, nil)
				m.mauCache.(*cachemock.MockMAUCache).EXPECT().
					MergeIntoMAUBatch(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(map[string]int64{"ANDROID": 100}, nil)
				m.monthlySummaryStorage.(*insightsstoragemock.MockMonthlySummaryStorage).EXPECT().
					UpsertMonthlySummaryBatch(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			summarizer := newMockMonthlySummarizer(t, mockController)
			p.setup(summarizer)
			err := summarizer.Run(ctx)
			if p.expected != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListSourceIDs(t *testing.T) {
	t.Parallel()
	ids := listSourceIDs()
	expected := []string{
		"ANDROID",
		"IOS",
		"WEB",
		"GO_SERVER",
		"NODE_SERVER",
		"JAVASCRIPT",
		"FLUTTER",
		"REACT",
		"REACT_NATIVE",
		"OPEN_FEATURE_KOTLIN",
		"OPEN_FEATURE_SWIFT",
		"OPEN_FEATURE_JAVASCRIPT",
		"OPEN_FEATURE_GO",
		"OPEN_FEATURE_NODE",
		"OPEN_FEATURE_REACT",
		"OPEN_FEATURE_REACT_NATIVE",
	}
	assert.ElementsMatch(t, expected, ids)
}

func newMockMonthlySummarizer(t *testing.T, c *gomock.Controller) *monthlySummarizer {
	t.Helper()
	return &monthlySummarizer{
		envClient:             envclientmock.NewMockClient(c),
		mauCache:              cachemock.NewMockMAUCache(c),
		monthlySummaryStorage: insightsstoragemock.NewMockMonthlySummaryStorage(c),
		location:              time.UTC,
		opts: &jobs.Options{
			Timeout: 30 * time.Second,
		},
		logger: zap.NewNop().Named("test-monthly-summarizer"),
	}
}
