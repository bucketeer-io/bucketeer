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

package api

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/v2/proto/account"
)

func TestCacheAPIKeyLastUsedAt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		apikey        *account.APIKey
		lastUsedAt    int64
		existingCache sync.Map
		expectedCache sync.Map
	}{
		{
			name:       "new entry",
			apikey:     &account.APIKey{Id: "key1"},
			lastUsedAt: 1000,
			existingCache: func() sync.Map {
				var m sync.Map
				return m
			}(),
			expectedCache: func() sync.Map {
				var m sync.Map
				m.Store("key1", int64(1000))
				return m
			}(),
		},
		{
			name:       "update existing entry with higher lastUsedAt",
			apikey:     &account.APIKey{Id: "key1"},
			lastUsedAt: 2000,
			existingCache: func() sync.Map {
				var m sync.Map
				m.Store("key1", int64(1500))
				return m
			}(),
			expectedCache: func() sync.Map {
				var m sync.Map
				m.Store("key1", int64(2000))
				return m
			}(),
		},
		{
			name:       "do not update existing entry with lower lastUsedAt",
			apikey:     &account.APIKey{Id: "key1"},
			lastUsedAt: 1000,
			existingCache: func() sync.Map {
				var m sync.Map
				m.Store("key1", int64(1500))
				return m
			}(),
			expectedCache: func() sync.Map {
				var m sync.Map
				m.Store("key1", int64(1500))
				return m
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &grpcGatewayService{
				apiKeyLastUsedInfoCacher: tt.existingCache,
			}
			service.cacheAPIKeyLastUsedAt(tt.apikey, tt.lastUsedAt)

			listExpected := make(map[string]int64)
			tt.expectedCache.Range(func(key, value interface{}) bool {
				listExpected[key.(string)] = value.(int64)
				return true
			})

			listActual := make(map[string]int64)
			service.apiKeyLastUsedInfoCacher.Range(func(key, value interface{}) bool {
				listActual[key.(string)] = value.(int64)
				return true
			})

			assert.Equal(t, listExpected, listActual)
		})
	}
}
