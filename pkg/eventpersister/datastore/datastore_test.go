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

package datastore

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWriter struct {
	writes int
}

func (w *testWriter) Write(ctx context.Context, events map[string]string, environmentNamespace string) (map[string]bool, error) {
	w.writes++
	return nil, nil
}

func (w *testWriter) Close() {}

func newTestWriters(num int) []Writer {
	writers := make([]Writer, 0, num)
	for i := 0; i < num; i++ {
		writers = append(writers, &testWriter{})
	}
	return writers
}

func TestWriterPool(t *testing.T) {
	testcases := []struct {
		writes               uint64
		writers              []Writer
		environmentNamespace string
	}{
		{
			writes:               0,
			writers:              newTestWriters(5),
			environmentNamespace: "ns0",
		},
		{
			writes:               1<<64 - 1,
			writers:              newTestWriters(2),
			environmentNamespace: "ns0",
		},
	}
	for i, tc := range testcases {
		des := fmt.Sprintf("Index %d", i)
		pool := writerPool{
			writes:  tc.writes,
			writers: tc.writers,
		}
		for j := 0; j < len(tc.writers); j++ {
			pool.Write(context.Background(), nil, tc.environmentNamespace)
		}
		for j := 0; j < len(tc.writers); j++ {
			writer, ok := tc.writers[j].(*testWriter)
			require.Equal(t, true, ok, des)
			assert.Equal(t, 1, writer.writes, des)
		}
	}
}
