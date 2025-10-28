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
	"github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

func TestCacheAPIKeyLastUsedAt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		envAPIKey     *account.EnvironmentAPIKey
		lastUsedAt    int64
		existingCache envAPIKeyLastUsedInfoCache
		expectedCache envAPIKeyLastUsedInfoCache
	}{
		{
			name: "Add API key last used at to cache for new environment ID",
			envAPIKey: &account.EnvironmentAPIKey{
				Environment: &environment.EnvironmentV2{
					Id:   "env-1",
					Name: "Environment 1",
				},
				ApiKey: &account.APIKey{
					Id:   "api-key-1",
					Name: "API Key 1",
				},
			},
			lastUsedAt:    1625155200,
			existingCache: envAPIKeyLastUsedInfoCache{},
			expectedCache: envAPIKeyLastUsedInfoCache{
				"env-1": {
					"api-key-1": lastUsedInfo{
						envAPIKey: &account.EnvironmentAPIKey{
							Environment: &environment.EnvironmentV2{
								Id:   "env-1",
								Name: "Environment 1",
							},
							ApiKey: &account.APIKey{
								Id:   "api-key-1",
								Name: "API Key 1",
							},
						},
						lastUsedAt: 1625155200,
					},
				},
			},
		},
		{
			name: "Update existing API key last used at in cache",
			envAPIKey: &account.EnvironmentAPIKey{
				Environment: &environment.EnvironmentV2{
					Id:   "env-1",
					Name: "Environment 1",
				},
				ApiKey: &account.APIKey{
					Id:   "api-key-1",
					Name: "API Key 1",
				},
			},
			lastUsedAt: 1625241600,
			existingCache: envAPIKeyLastUsedInfoCache{
				"env-1": {
					"api-key-1": lastUsedInfo{
						envAPIKey: &account.EnvironmentAPIKey{
							Environment: &environment.EnvironmentV2{
								Id:   "env-1",
								Name: "Environment 1",
							},
							ApiKey: &account.APIKey{
								Id:   "api-key-1",
								Name: "API Key 1",
							},
						},
						lastUsedAt: 1620055200,
					},
				},
				"env-2": {
					"api-key-2": lastUsedInfo{
						envAPIKey: &account.EnvironmentAPIKey{
							Environment: &environment.EnvironmentV2{
								Id:   "env-2",
								Name: "Environment 2",
							},
							ApiKey: &account.APIKey{
								Id:         "api-key-2",
								LastUsedAt: 1620055200,
							},
						},
						lastUsedAt: 1621055200,
					},
				},
			},
			expectedCache: envAPIKeyLastUsedInfoCache{
				"env-1": {
					"api-key-1": lastUsedInfo{
						envAPIKey: &account.EnvironmentAPIKey{
							Environment: &environment.EnvironmentV2{
								Id:   "env-1",
								Name: "Environment 1",
							},
							ApiKey: &account.APIKey{
								Id:   "api-key-1",
								Name: "API Key 1",
							},
						},
						lastUsedAt: 1625241600,
					},
				},
				"env-2": {
					"api-key-2": lastUsedInfo{
						envAPIKey: &account.EnvironmentAPIKey{
							Environment: &environment.EnvironmentV2{
								Id:   "env-2",
								Name: "Environment 2",
							},
							ApiKey: &account.APIKey{
								Id:         "api-key-2",
								LastUsedAt: 1620055200,
							},
						},
						lastUsedAt: 1621055200,
					},
				},
			},
		},
		{
			name: "Do not update API key last used at if existing is newer",
			envAPIKey: &account.EnvironmentAPIKey{
				Environment: &environment.EnvironmentV2{
					Id:   "env-1",
					Name: "Environment 1",
				},
				ApiKey: &account.APIKey{
					Id:   "api-key-1",
					Name: "API Key 1",
				},
			},
			lastUsedAt: 1620055200,
			existingCache: envAPIKeyLastUsedInfoCache{
				"env-1": {
					"api-key-1": lastUsedInfo{
						envAPIKey: &account.EnvironmentAPIKey{
							Environment: &environment.EnvironmentV2{
								Id:   "env-1",
								Name: "Environment 1",
							},
							ApiKey: &account.APIKey{
								Id:   "api-key-1",
								Name: "API Key 1",
							},
						},
						lastUsedAt: 1625241600,
					},
				},
			},
			expectedCache: envAPIKeyLastUsedInfoCache{
				"env-1": {
					"api-key-1": lastUsedInfo{
						envAPIKey: &account.EnvironmentAPIKey{
							Environment: &environment.EnvironmentV2{
								Id:   "env-1",
								Name: "Environment 1",
							},
							ApiKey: &account.APIKey{
								Id:   "api-key-1",
								Name: "API Key 1",
							},
						},
						lastUsedAt: 1625241600,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &grpcGatewayService{
				envAPIKeyLastUsedInfoCacher: tt.existingCache,
				envAPIKeyLastUsedInfoMutex:  sync.Mutex{},
			}
			service.cacheAPIKeyLastUsedAt(tt.envAPIKey, tt.lastUsedAt)
			assert.Equal(t, tt.expectedCache, service.envAPIKeyLastUsedInfoCacher)
		})
	}
}
