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

package apikey_last_used_at_writer

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/v2/proto/account"
	"github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

func TestCacheAPIKeyLastUsedAt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		apikey        *account.EnvironmentAPIKey
		lastUsedAt    int64
		existingCache sync.Map
		expectedCache sync.Map
	}{
		{
			name: "new entry",
			apikey: &account.EnvironmentAPIKey{
				ApiKey: &account.APIKey{
					Id: "key1",
				},
				Environment: &environment.EnvironmentV2{
					Id: "env1",
				},
			},
			lastUsedAt: 1000,
			existingCache: func() sync.Map {
				var m sync.Map
				return m
			}(),
			expectedCache: func() sync.Map {
				var m sync.Map
				m.Store("key1", apikeyLastUsedAt{
					apiKeyID:      "key1",
					lastUsedAt:    1000,
					environmentID: "env1",
				})
				return m
			}(),
		},
		{
			name: "update existing entry with higher lastUsedAt",
			apikey: &account.EnvironmentAPIKey{
				ApiKey: &account.APIKey{
					Id: "key1",
				},
				Environment: &environment.EnvironmentV2{
					Id: "env1",
				},
			},
			lastUsedAt: 2000,
			existingCache: func() sync.Map {
				var m sync.Map
				m.Store("key1", apikeyLastUsedAt{
					apiKeyID:      "key1",
					lastUsedAt:    1500,
					environmentID: "env1",
				})
				return m
			}(),
			expectedCache: func() sync.Map {
				var m sync.Map
				m.Store("key1", apikeyLastUsedAt{
					apiKeyID:      "key1",
					lastUsedAt:    2000,
					environmentID: "env1",
				})
				return m
			}(),
		},
		{
			name: "do not update existing entry with lower lastUsedAt",
			apikey: &account.EnvironmentAPIKey{
				ApiKey: &account.APIKey{
					Id: "key1",
				},
				Environment: &environment.EnvironmentV2{
					Id: "env1",
				},
			},
			lastUsedAt: 1000,
			existingCache: func() sync.Map {
				var m sync.Map
				m.Store("key1", apikeyLastUsedAt{
					apiKeyID:      "key1",
					lastUsedAt:    1500,
					environmentID: "env1",
				})
				return m
			}(),
			expectedCache: func() sync.Map {
				var m sync.Map
				m.Store("key1", apikeyLastUsedAt{
					apiKeyID:      "key1",
					lastUsedAt:    1500,
					environmentID: "env1",
				})
				return m
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &APIKeyLastUsedWriter{
				APIKeyLastUsedInfoCacher: tt.existingCache,
			}
			service.CacheAPIKeyLastUsedAt(tt.apikey, tt.lastUsedAt)

			listExpected := make(map[string]apikeyLastUsedAt)
			tt.expectedCache.Range(func(key, value interface{}) bool {
				listExpected[key.(string)] = value.(apikeyLastUsedAt)
				return true
			})

			listActual := make(map[string]apikeyLastUsedAt)
			service.APIKeyLastUsedInfoCacher.Range(func(key, value interface{}) bool {
				listActual[key.(string)] = value.(apikeyLastUsedAt)
				return true
			})

			assert.Equal(t, listExpected, listActual)
		})
	}
}
