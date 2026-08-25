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

package api

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
		setupCache    func(m *sync.Map)
		expectedCache map[string]apikeyLastUsedAt
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
			setupCache: func(m *sync.Map) {},
			expectedCache: map[string]apikeyLastUsedAt{
				"key1": {
					apiKeyID:      "key1",
					lastUsedAt:    1000,
					environmentID: "env1",
				},
			},
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
			setupCache: func(m *sync.Map) {
				m.Store("key1", apikeyLastUsedAt{
					apiKeyID:      "key1",
					lastUsedAt:    1500,
					environmentID: "env1",
				})
			},
			expectedCache: map[string]apikeyLastUsedAt{
				"key1": {
					apiKeyID:      "key1",
					lastUsedAt:    2000,
					environmentID: "env1",
				},
			},
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
			setupCache: func(m *sync.Map) {
				m.Store("key1", apikeyLastUsedAt{
					apiKeyID:      "key1",
					lastUsedAt:    1500,
					environmentID: "env1",
				})
			},
			expectedCache: map[string]apikeyLastUsedAt{
				"key1": {
					apiKeyID:      "key1",
					lastUsedAt:    1500,
					environmentID: "env1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &grpcGatewayService{}
			tt.setupCache(&service.apiKeyLastUsedInfoCacher)
			service.cacheAPIKeyLastUsedAt(tt.apikey, tt.lastUsedAt)

			listActual := make(map[string]apikeyLastUsedAt)
			service.apiKeyLastUsedInfoCacher.Range(func(key, value interface{}) bool {
				listActual[key.(string)] = value.(apikeyLastUsedAt)
				return true
			})

			assert.Equal(t, tt.expectedCache, listActual)
		})
	}
}
