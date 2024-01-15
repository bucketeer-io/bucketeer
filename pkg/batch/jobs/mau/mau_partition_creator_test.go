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

func TestMAUPartitionCreatorRun(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	patterns := []struct {
		desc  string
		setup func(d *mauPartitionCreator)

		expected error
	}{
		{
			desc: "fail at CreatePartition",
			setup: func(d *mauPartitionCreator) {
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().CreatePartition(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					errors.New("test"),
				)
			},
			expected: errors.New("test"),
		},
		{
			desc: "success",
			setup: func(d *mauPartitionCreator) {
				d.mauStorage.(*storagemock.MockMAUStorage).EXPECT().CreatePartition(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil,
				)
			},
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			p := p
			creator := newMAUPartitionCreator(t, mockController)
			p.setup((*mauPartitionCreator)(creator))
			err := creator.Run(ctx)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestMAUPartitionCreatorNewPartitionName(t *testing.T) {
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
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			creator := newMAUPartitionCreator(t, nil)
			actual := creator.newPartitionName(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestMAUPartitionCreatorNewLessThan(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		input    time.Time
		expected string
	}{
		{
			desc:     "success1",
			input:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: "202310",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			p := p
			creator := newMAUPartitionCreator(t, nil)
			actual := creator.newLessThan(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newMAUPartitionCreator(t *testing.T, c *gomock.Controller) *mauPartitionCreator {
	t.Helper()
	return &mauPartitionCreator{
		mauStorage: storagemock.NewMockMAUStorage(c),
		location:   time.UTC,
		opts: &jobs.Options{
			Timeout: 60 * time.Minute,
		},
		logger: zap.NewNop().Named("test-mau-partition-deleter"),
	}
}
