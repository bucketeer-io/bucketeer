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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aoproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	exproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

// inMemoryMultiGetCache adapts InMemoryCache, which already supports conditional writes, to MultiGetCache
type inMemoryMultiGetCache struct {
	*InMemoryCache
}

func (inMemoryMultiGetCache) GetMulti(interface{}, bool) ([]interface{}, error) { return nil, nil }
func (inMemoryMultiGetCache) ScanKeys(string, int64) ([]string, error)          { return nil, nil }
func (inMemoryMultiGetCache) SMembers(string) ([]string, error)                 { return nil, nil }

// The batch cacher reads an old snapshot (generation 1), the refresher then writes a newer one (generation 2),
// and the batch cacher's late write must not overwrite it.
func TestSnapshotPutIfNewerWriteOrdering(t *testing.T) {
	t.Parallel()
	const (
		batchGeneration     = int64(1)
		refresherGeneration = int64(2)
	)
	patterns := []struct {
		desc string
		run  func(t *testing.T, c inMemoryMultiGetCache) string
	}{
		{
			desc: "auto ops rules",
			run: func(t *testing.T, c inMemoryMultiGetCache) string {
				rc := NewAutoOpsRulesCache(c)
				newer := &aoproto.AutoOpsRules{AutoOpsRules: []*aoproto.AutoOpsRule{{Id: "new-rule"}}}
				accepted, err := rc.PutIfNewer(newer, environmentId, refresherGeneration)
				require.NoError(t, err)
				require.True(t, accepted)
				accepted, err = rc.PutIfNewer(&aoproto.AutoOpsRules{}, environmentId, batchGeneration)
				require.NoError(t, err)
				assert.False(t, accepted)
				actual, err := rc.Get(environmentId)
				require.NoError(t, err)
				require.Len(t, actual.AutoOpsRules, 1)
				return actual.AutoOpsRules[0].Id
			},
		},
		{
			desc: "experiments",
			run: func(t *testing.T, c inMemoryMultiGetCache) string {
				ec := NewExperimentsCache(c)
				newer := &exproto.Experiments{Experiments: []*exproto.Experiment{{Id: "new-rule"}}}
				accepted, err := ec.PutIfNewer(newer, environmentId, refresherGeneration)
				require.NoError(t, err)
				require.True(t, accepted)
				accepted, err = ec.PutIfNewer(&exproto.Experiments{}, environmentId, batchGeneration)
				require.NoError(t, err)
				assert.False(t, accepted)
				actual, err := ec.Get(environmentId)
				require.NoError(t, err)
				require.Len(t, actual.Experiments, 1)
				return actual.Experiments[0].Id
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			c := inMemoryMultiGetCache{NewInMemoryCache()}
			assert.Equal(t, "new-rule", p.run(t, c))
		})
	}
}
