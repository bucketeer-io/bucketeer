// Copyright 2022 The Bucketeer Authors.
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

package writer

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

var defaultOptions = options{
	logger:    zap.NewNop(),
	batchSize: 10,
}

func TestCreateBatch(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
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
			w := newWriter(mockController)
			actual := w.getBatch(encoded)
			for idx, exp := range p.expectedSize {
				assert.Len(t, actual[idx], exp)
			}
		})
	}
}

func newWriter(c *gomock.Controller) *writer {
	return &writer{
		opts: &defaultOptions,
	}
}
