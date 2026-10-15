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

package v3

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type scanPage struct {
	keys []string
	next uint64
	err  error
}

func TestScanNode(t *testing.T) {
	t.Parallel()
	errScan := errors.New("scan failed")
	errCallback := errors.New("callback failed")
	patterns := []struct {
		desc        string
		pages       []scanPage
		callbackErr error
		wantCursors []uint64
		wantKeys    [][]string
		wantErr     error
	}{
		{
			desc: "follows the cursor until it returns to zero",
			pages: []scanPage{
				{keys: []string{"a", "b"}, next: 7},
				{keys: []string{}, next: 9},
				{keys: []string{"c"}, next: 0},
			},
			wantCursors: []uint64{0, 7, 9},
			wantKeys:    [][]string{{"a", "b"}, {"c"}},
		},
		{
			desc:        "scan error stops the iteration",
			pages:       []scanPage{{keys: []string{"a"}, next: 3}, {err: errScan}},
			wantCursors: []uint64{0, 3},
			wantKeys:    [][]string{{"a"}},
			wantErr:     errScan,
		},
		{
			desc:        "callback error stops the iteration",
			pages:       []scanPage{{keys: []string{"a"}, next: 3}},
			callbackErr: errCallback,
			wantCursors: []uint64{0},
			wantKeys:    [][]string{{"a"}},
			wantErr:     errCallback,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			c := &client{opts: defaultOptions(), logger: zap.NewNop()}
			var cursors []uint64
			scan := func(_ context.Context, cursor uint64, match string, count int64) *goredis.ScanCmd {
				assert.Equal(t, "*:retry:*", match)
				assert.Equal(t, int64(100), count)
				page := p.pages[len(cursors)]
				cursors = append(cursors, cursor)
				return goredis.NewScanCmdResult(page.keys, page.next, page.err)
			}
			var keys [][]string
			err := c.scanNode(context.Background(), scan, "*:retry:*", 100, func(k []string) error {
				keys = append(keys, k)
				return p.callbackErr
			})
			assert.ErrorIs(t, err, p.wantErr)
			assert.Equal(t, p.wantCursors, cursors)
			assert.Equal(t, p.wantKeys, keys)
		})
	}
}

func TestSerialize(t *testing.T) {
	t.Parallel()
	var inFlight, maxInFlight, calls atomic.Int32
	fn := serialize(func(_ []string) error {
		n := inFlight.Add(1)
		for {
			m := maxInFlight.Load()
			if n <= m || maxInFlight.CompareAndSwap(m, n) {
				break
			}
		}
		calls.Add(1)
		inFlight.Add(-1)
		return nil
	})
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(t, fn([]string{"k"}))
		}()
	}
	wg.Wait()
	assert.Equal(t, int32(50), calls.Load())
	assert.Equal(t, int32(1), maxInFlight.Load())
}

func TestSerializeStopsAfterError(t *testing.T) {
	t.Parallel()
	fnErr := errors.New("fn failed")
	var calls atomic.Int32
	fn := serialize(func(_ []string) error {
		calls.Add(1)
		return fnErr
	})
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.ErrorIs(t, fn([]string{"k"}), fnErr)
		}()
	}
	wg.Wait()
	assert.Equal(t, int32(1), calls.Load())
}
