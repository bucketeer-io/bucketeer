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

package storage

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
)

func TestCreateBatch(t *testing.T) {
	batchSize := 10
	t.Parallel()
	patterns := []struct {
		desc         string
		size         int
		expectedSize []int
	}{
		{
			desc:         "success: 15",
			size:         15,
			expectedSize: []int{10, 5},
		},
		{
			desc:         "success: 11",
			size:         11,
			expectedSize: []int{10, 1},
		},
		{
			desc:         "success: 20",
			size:         20,
			expectedSize: []int{10, 10},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evts := []*eventproto.GoalEvent{}
			for i := 0; i < p.size; i++ {
				evt := &eventproto.GoalEvent{
					GoalId: strconv.Itoa(i),
				}
				evts = append(evts, evt)
			}
			encoded := make([][]byte, len(evts))
			for k, v := range evts {
				b, err := proto.Marshal(v)
				require.NoError(t, err)
				encoded[k] = b
			}
			actual := getBatch(encoded, batchSize)
			for idx, exp := range p.expectedSize {
				assert.Len(t, actual[idx], exp)
			}
		})
	}
}

func TestGoalGetFailMap(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc       string
		size       int
		inputSize  int
		inputFails []int
		expected   map[string]bool
	}{
		{
			desc:       "success: size is 15",
			size:       15,
			inputFails: []int{0},
			expected: map[string]bool{
				"0": true,
				"1": true,
				"2": true,
				"3": true,
				"4": true,
				"5": true,
				"6": true,
				"7": true,
				"8": true,
				"9": true,
			},
		},
		{
			desc:       "success: size is 15, full id",
			size:       15,
			inputFails: []int{0, 1},
			expected: map[string]bool{
				"0":  true,
				"1":  true,
				"2":  true,
				"3":  true,
				"4":  true,
				"5":  true,
				"6":  true,
				"7":  true,
				"8":  true,
				"9":  true,
				"10": true,
				"11": true,
				"12": true,
				"13": true,
				"14": true,
			},
		},
		{
			desc:       "success: size is 9, full id",
			size:       9,
			inputFails: []int{0},
			expected: map[string]bool{
				"0": true,
				"1": true,
				"2": true,
				"3": true,
				"4": true,
				"5": true,
				"6": true,
				"7": true,
				"8": true,
			},
		},
		{
			desc:       "success: size is 35",
			size:       35,
			inputFails: []int{0, 2},
			expected: map[string]bool{
				"0":  true,
				"1":  true,
				"2":  true,
				"3":  true,
				"4":  true,
				"5":  true,
				"6":  true,
				"7":  true,
				"8":  true,
				"9":  true,
				"20": true,
				"21": true,
				"22": true,
				"23": true,
				"24": true,
				"25": true,
				"26": true,
				"27": true,
				"28": true,
				"29": true,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evts := []*epproto.GoalEvent{}
			for i := 0; i < p.size; i++ {
				evt := &epproto.GoalEvent{
					Id: strconv.Itoa(i),
				}
				evts = append(evts, evt)
			}
			writer := newGoalWriter(mockController)
			actual := writer.getFailMap(evts, p.inputFails)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestEvalGetFailMap(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc       string
		size       int
		inputSize  int
		inputFails []int
		expected   map[string]bool
	}{
		{
			desc:       "success: size is 15",
			size:       15,
			inputFails: []int{0},
			expected: map[string]bool{
				"0": true,
				"1": true,
				"2": true,
				"3": true,
				"4": true,
				"5": true,
				"6": true,
				"7": true,
				"8": true,
				"9": true,
			},
		},
		{
			desc:       "success: size is 15, full id",
			size:       15,
			inputFails: []int{0, 1},
			expected: map[string]bool{
				"0":  true,
				"1":  true,
				"2":  true,
				"3":  true,
				"4":  true,
				"5":  true,
				"6":  true,
				"7":  true,
				"8":  true,
				"9":  true,
				"10": true,
				"11": true,
				"12": true,
				"13": true,
				"14": true,
			},
		},
		{
			desc:       "success: size is 9, full id",
			size:       9,
			inputFails: []int{0},
			expected: map[string]bool{
				"0": true,
				"1": true,
				"2": true,
				"3": true,
				"4": true,
				"5": true,
				"6": true,
				"7": true,
				"8": true,
			},
		},
		{
			desc:       "success: size is 35",
			size:       35,
			inputFails: []int{0, 2},
			expected: map[string]bool{
				"0":  true,
				"1":  true,
				"2":  true,
				"3":  true,
				"4":  true,
				"5":  true,
				"6":  true,
				"7":  true,
				"8":  true,
				"9":  true,
				"20": true,
				"21": true,
				"22": true,
				"23": true,
				"24": true,
				"25": true,
				"26": true,
				"27": true,
				"28": true,
				"29": true,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evts := []*epproto.EvaluationEvent{}
			for i := 0; i < p.size; i++ {
				evt := &epproto.EvaluationEvent{
					Id: strconv.Itoa(i),
				}
				evts = append(evts, evt)
			}
			writer := newEvalWriter(mockController)
			actual := writer.getFailMap(evts, p.inputFails)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newGoalWriter(c *gomock.Controller) *goalEventWriter {
	return &goalEventWriter{
		queryClient: &queryClient{
			batchSize: 10,
		},
	}
}

func newEvalWriter(c *gomock.Controller) *evalEventWriter {
	return &evalEventWriter{
		queryClient: &queryClient{
			batchSize: 10,
		},
	}
}
