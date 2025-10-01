package processor

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
)

func TestCacheAPIKeyLastUsedInfo(t *testing.T) {
	tests := []struct {
		name           string
		existingCache  envAPIKeyLastUsedInfoCache
		envEvents      envAPIKeyUsageEventMap
		expectedResult envAPIKeyLastUsedInfoCache
	}{
		{
			name:          "Add API key last used info to cache for new environment ID",
			existingCache: make(envAPIKeyLastUsedInfoCache),
			envEvents: envAPIKeyUsageEventMap{
				"env-1": apikeyUsageEventMap{
					"event-1": &eventproto.APIKeyUsageEvent{
						ApiKeyId:      "api-key-1",
						Timestamp:     1620000000,
						EnvironmentId: "env-1",
					},
					"event-2": &eventproto.APIKeyUsageEvent{
						ApiKeyId:      "api-key-2",
						Timestamp:     1620500000,
						EnvironmentId: "env-1",
					},
				},
			},
			expectedResult: envAPIKeyLastUsedInfoCache{
				"env-1": apikeyLastUsedInfoCache{
					"api-key-1": &domain.APIKeyLastUsedInfo{
						APIKeyLastUsedInfo: &accountproto.APIKeyLastUsedInfo{
							ApiKeyId:      "api-key-1",
							LastUsedAt:    1620000000,
							EnvironmentId: "env-1",
						},
					},
					"api-key-2": &domain.APIKeyLastUsedInfo{
						APIKeyLastUsedInfo: &accountproto.APIKeyLastUsedInfo{
							ApiKeyId:      "api-key-2",
							LastUsedAt:    1620500000,
							EnvironmentId: "env-1",
						},
					},
				},
			},
		},
		{
			name: "Update existing env API key last used info in cache",
			existingCache: envAPIKeyLastUsedInfoCache{
				"env-1": apikeyLastUsedInfoCache{
					"api-key-1": &domain.APIKeyLastUsedInfo{
						APIKeyLastUsedInfo: &accountproto.APIKeyLastUsedInfo{
							ApiKeyId:      "api-key-1",
							LastUsedAt:    1620000000,
							EnvironmentId: "env-1",
						},
					},
				},
			},
			envEvents: envAPIKeyUsageEventMap{
				"env-1": apikeyUsageEventMap{
					"event-1": &eventproto.APIKeyUsageEvent{
						ApiKeyId:      "api-key-2",
						Timestamp:     1620500000,
						EnvironmentId: "env-1",
					},
				},
				"env-2": apikeyUsageEventMap{
					"event-2": &eventproto.APIKeyUsageEvent{
						ApiKeyId:      "api-key-3",
						Timestamp:     1621000000,
						EnvironmentId: "env-2",
					},
				},
			},
			expectedResult: envAPIKeyLastUsedInfoCache{
				"env-1": apikeyLastUsedInfoCache{
					"api-key-1": &domain.APIKeyLastUsedInfo{
						APIKeyLastUsedInfo: &accountproto.APIKeyLastUsedInfo{
							ApiKeyId:      "api-key-1",
							LastUsedAt:    1620000000,
							EnvironmentId: "env-1",
						},
					},
					"api-key-2": &domain.APIKeyLastUsedInfo{
						APIKeyLastUsedInfo: &accountproto.APIKeyLastUsedInfo{
							ApiKeyId:      "api-key-2",
							LastUsedAt:    1620500000,
							EnvironmentId: "env-1",
						},
					},
				},
				"env-2": apikeyLastUsedInfoCache{
					"api-key-3": &domain.APIKeyLastUsedInfo{
						APIKeyLastUsedInfo: &accountproto.APIKeyLastUsedInfo{
							ApiKeyId:      "api-key-3",
							LastUsedAt:    1621000000,
							EnvironmentId: "env-2",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &apikeyLastUsedInfoWriter{
				apikeyLastUsedInfoCacher: tt.existingCache,
				envLastUsedCacheMutex:    sync.Mutex{},
				logger:                   zap.NewNop(),
			}

			w.cacheAPIKeyLastUsedInfoPerEnv(tt.envEvents)

			for _, cache := range tt.existingCache {
				for apiKeyID, info := range cache {
					expectedInfo, ok := tt.expectedResult[info.EnvironmentId][apiKeyID]
					if !ok {
						t.Errorf("Expected API key ID %s not found in cache", apiKeyID)
						continue
					}

					assert.Equal(t, expectedInfo.ApiKeyId, info.ApiKeyId)
					assert.Equal(t, expectedInfo.LastUsedAt, info.LastUsedAt)
					assert.Equal(t, expectedInfo.EnvironmentId, info.EnvironmentId)
				}
			}
		})
	}
}
