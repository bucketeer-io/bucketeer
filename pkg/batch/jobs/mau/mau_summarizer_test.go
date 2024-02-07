// Copyright 2024 The Bucketeer Authors.
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

package mau

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	ecclientmock "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

func TestRun(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	patterns := []struct {
		desc     string
		setup    func(r *mauSummarizer)
		expected error
	}{
		{
			desc: "fail",
			setup: func(r *mauSummarizer) {
				r.eventCounterClient.(*ecclientmock.MockClient).EXPECT().SummarizeMAUCounts(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("test"))
			},
			expected: errors.New("test"),
		},
		{
			desc: "success",
			setup: func(r *mauSummarizer) {
				r.eventCounterClient.(*ecclientmock.MockClient).EXPECT().SummarizeMAUCounts(gomock.Any(), gomock.Any()).Return(
					&ecproto.SummarizeMAUCountsResponse{}, nil)
			},
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			summarizer := newMockMAUSummarizer(t, mockController)
			p.setup(summarizer)
			err := summarizer.Run(ctx)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestGetYesterday(t *testing.T) {
	t.Parallel()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	patterns := []struct {
		desc     string
		input    time.Time
		expected time.Time
	}{
		{
			desc:     "success",
			input:    time.Date(2023, 1, 31, 1, 00, 00, 0, time.UTC),
			expected: time.Date(2023, 1, 30, 1, 00, 00, 0, time.UTC),
		},
		{
			desc:     "over month",
			input:    time.Date(2023, 2, 1, 1, 00, 00, 0, time.UTC),
			expected: time.Date(2023, 1, 31, 1, 00, 00, 0, time.UTC),
		},
		{
			desc:     "over year",
			input:    time.Date(2023, 1, 1, 1, 00, 00, 0, time.UTC),
			expected: time.Date(2022, 12, 31, 1, 00, 00, 0, time.UTC),
		},
		{
			desc:     "success JST",
			input:    time.Date(2023, 1, 31, 1, 00, 00, 0, jst),
			expected: time.Date(2023, 1, 30, 1, 00, 00, 0, jst),
		},
		{
			desc:     "over month JST",
			input:    time.Date(2023, 2, 1, 1, 00, 00, 0, jst),
			expected: time.Date(2023, 1, 31, 1, 00, 00, 0, jst),
		},
		{
			desc:     "over year JST",
			input:    time.Date(2023, 1, 1, 1, 00, 00, 0, jst),
			expected: time.Date(2022, 12, 31, 1, 00, 00, 0, jst),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			summarizer := &mauSummarizer{}
			actual := summarizer.getYesterday(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestIsEndDateOfMonth(t *testing.T) {
	t.Parallel()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	patterns := []struct {
		desc     string
		input    time.Time
		expected bool
	}{
		{
			desc:     "not end of month",
			input:    time.Date(2023, 1, 30, 1, 00, 00, 0, time.UTC),
			expected: false,
		},
		{
			desc:     "success1",
			input:    time.Date(2023, 1, 31, 1, 00, 00, 0, time.UTC),
			expected: true,
		},
		{
			desc:     "success2",
			input:    time.Date(2023, 2, 1, 1, 00, 00, 0, jst),
			expected: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			summarizer := &mauSummarizer{}
			actual := summarizer.isEndDateOfMonth(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newMockMAUSummarizer(t *testing.T, c *gomock.Controller) *mauSummarizer {
	t.Helper()
	return &mauSummarizer{
		mysqlClient:        mysqlmock.NewMockClient(c),
		eventCounterClient: ecclientmock.NewMockClient(c),
		location:           time.UTC,
		opts: &jobs.Options{
			Timeout: 30 * time.Second,
		},
		logger: zap.NewNop().Named("test-mau-summarizer"),
	}
}
