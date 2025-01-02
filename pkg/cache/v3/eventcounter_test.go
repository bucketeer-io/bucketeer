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

package v3

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/cache/mock"
)

func TestGetEventValues(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		input    []interface{}
		expected []float64
		isValid  bool
	}{
		{
			desc:     "success: empty input",
			input:    []interface{}{},
			expected: []float64{},
			isValid:  true,
		},
		{
			desc:     "success: valid input",
			input:    []interface{}{[]byte("1.5"), []byte("2.7"), []byte("3.0")},
			expected: []float64{1.5, 2.7, 3.0},
			isValid:  true,
		},
		{
			desc:     "failure: invalid input",
			input:    []interface{}{[]byte("not a number")},
			expected: nil,
			isValid:  false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			c := &eventCounterCache{}
			actual, err := c.getEventValues(p.input)
			if p.isValid {
				assert.NoError(t, err)
				assert.Equal(t, p.expected, actual)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestGetEventCounts(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	c := &eventCounterCache{cache: mockCache}

	patterns := []struct {
		desc     string
		keys     []string
		mockResp []interface{}
		expected []float64
		isValid  bool
	}{
		{
			desc:     "success: valid input",
			keys:     []string{"key1", "key2", "key3"},
			mockResp: []interface{}{[]byte("1.5"), []byte("2.7"), []byte("3.0")},
			expected: []float64{1.5, 2.7, 3.0},
			isValid:  true,
		},
		{
			desc:     "failure: invalid input",
			keys:     []string{"key1"},
			mockResp: []interface{}{[]byte("not a number")},
			expected: nil,
			isValid:  false,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockCache.EXPECT().GetMulti(p.keys, true).Return(p.mockResp, nil)
			actual, err := c.GetEventCounts(p.keys)
			if p.isValid {
				assert.NoError(t, err)
				assert.Equal(t, p.expected, actual)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestGetUserCounts(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	c := &eventCounterCache{cache: mockCache}

	patterns := []struct {
		desc     string
		keys     []string
		mockResp []int64
		expected []float64
		isValid  bool
	}{
		{
			desc:     "success: valid input",
			keys:     []string{"key1", "key2", "key3"},
			mockResp: []int64{10, 20, 30},
			expected: []float64{10, 20, 30},
			isValid:  true,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			for i, key := range p.keys {
				mockCache.EXPECT().PFCount(key).Return(p.mockResp[i], nil)
			}
			actual, err := c.GetUserCounts(p.keys)
			if p.isValid {
				assert.NoError(t, err)
				assert.Equal(t, p.expected, actual)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
