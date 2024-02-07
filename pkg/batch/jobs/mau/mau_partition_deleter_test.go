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
	storagemock "github.com/bucketeer-io/bucketeer/pkg/mau/storage/mock"
)

func TestMAUPartitionDeleterRun(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	patterns := []struct {
		desc     string
		setup    func(d *mauPartitionDeleter)
		expected error
	}{
		{
			desc: "fail at DeleteRecords",
			setup: func(d *mauPartitionDeleter) {
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().DeleteRecords(gomock.Any(), gomock.Any()).Return(
					errors.New("test"),
				)
			},
			expected: errors.New("test"),
		},
		{
			desc: "fail at RebuildPartition",
			setup: func(d *mauPartitionDeleter) {
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().DeleteRecords(gomock.Any(), gomock.Any()).Return(
					nil,
				)
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().RebuildPartition(gomock.Any(), gomock.Any()).Return(
					errors.New("test"),
				)
			},
			expected: errors.New("test"),
		},
		{
			desc: "fail at DropPartition",
			setup: func(d *mauPartitionDeleter) {
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().DeleteRecords(gomock.Any(), gomock.Any()).Return(
					nil,
				)
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().RebuildPartition(gomock.Any(), gomock.Any()).Return(
					nil,
				)
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().DropPartition(gomock.Any(), gomock.Any()).Return(
					errors.New("test"),
				)
			},
			expected: errors.New("test"),
		},
		{
			desc: "success",
			setup: func(d *mauPartitionDeleter) {
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().DeleteRecords(gomock.Any(), gomock.Any()).Return(
					nil,
				)
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().RebuildPartition(gomock.Any(), gomock.Any()).Return(
					nil,
				)
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().DropPartition(gomock.Any(), gomock.Any()).Return(
					nil,
				)
			},
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMAUPartitionDeleter(t, mockController)
			p.setup(deleter)
			err := deleter.Run(ctx)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestMAUPartitionDeleterGetDeleteTargetPartition(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		input    time.Time
		expected time.Time
	}{
		{
			desc:     "success1",
			input:    time.Date(2023, 10, 1, 1, 0, 0, 0, time.UTC),
			expected: time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:     "success2",
			input:    time.Date(2023, 10, 31, 1, 0, 0, 0, time.UTC),
			expected: time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:     "success3",
			input:    time.Date(2023, 2, 15, 23, 00, 00, 0, time.UTC),
			expected: time.Date(2022, 11, 1, 0, 00, 00, 0, time.UTC),
		},
		{
			desc:     "success4",
			input:    time.Date(2023, 5, 31, 1, 00, 00, 0, time.UTC),
			expected: time.Date(2023, 2, 1, 0, 00, 00, 0, time.UTC),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMAUPartitionDeleter(t, nil)
			actual := deleter.getDeleteTargetPartition(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestMAUPartitionDeleterNewPartitionName(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		input    time.Time
		expected string
	}{
		{
			desc:     "success1",
			input:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: "p202310",
		},
		{
			desc:     "success2",
			input:    time.Date(2023, 10, 31, 1, 0, 0, 0, time.UTC),
			expected: "p202310",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMAUPartitionDeleter(t, nil)
			actual := deleter.newPartitionName(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newMAUPartitionDeleter(t *testing.T, c *gomock.Controller) *mauPartitionDeleter {
	t.Helper()
	return &mauPartitionDeleter{
		mauStorage: storagemock.NewMockMAUStorage(c),
		location:   time.UTC,
		opts: &jobs.Options{
			Timeout: 60 * time.Minute,
		},
		logger: zap.NewNop().Named("test-mau-partition-deleter"),
	}
}
