package processor

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

func TestCacheUserAttributes(t *testing.T) {
	tests := []struct {
		name           string
		existingCache  userAttributesCache
		envEvents      environmentEventMap
		expectedResult userAttributesCache
	}{
		{
			name:          "Add user attributes to cache for new environment ID",
			existingCache: make(userAttributesCache),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId: "user-1",
						User: &userproto.User{
							Data: map[string]string{
								"country": "Japan",
								"city":    "Tokyo",
							},
						},
					},
				},
			},
			expectedResult: userAttributesCache{
				"env-1": &userproto.UserAttributes{
					EnvironmentId: "env-1",
					UserAttributes: []*userproto.UserAttribute{
						{
							Key:    "country",
							Values: []string{"Japan"},
						},
						{
							Key:    "city",
							Values: []string{"Tokyo"},
						},
					},
				},
			},
		},
		{
			name: "Add new attributes to existing cache",
			existingCache: userAttributesCache{
				"env-1": &userproto.UserAttributes{
					EnvironmentId: "env-1",
					UserAttributes: []*userproto.UserAttribute{
						{
							Key:    "country",
							Values: []string{"Japan"},
						},
						{
							Key:    "language",
							Values: []string{"ja"},
						},
					},
				},
			},
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId: "user-2",
						User: &userproto.User{
							Data: map[string]string{
								"country": "USA",
								"city":    "New York",
							},
						},
					},
				},
			},
			expectedResult: userAttributesCache{
				"env-1": &userproto.UserAttributes{
					EnvironmentId: "env-1",
					UserAttributes: []*userproto.UserAttribute{
						{
							Key:    "country",
							Values: []string{"Japan", "USA"},
						},
						{
							Key:    "language",
							Values: []string{"ja"},
						},
						{
							Key:    "city",
							Values: []string{"New York"},
						},
					},
				},
			},
		},
		{
			name: "Handle duplicate values appropriately",
			existingCache: userAttributesCache{
				"env-1": &userproto.UserAttributes{
					EnvironmentId: "env-1",
					UserAttributes: []*userproto.UserAttribute{
						{
							Key:    "country",
							Values: []string{"Japan", "USA"},
						},
					},
				},
			},
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId: "user-3",
						User: &userproto.User{
							Data: map[string]string{
								"country": "Japan", // Duplicate with existing value
								"city":    "Osaka",
							},
						},
					},
				},
			},
			expectedResult: userAttributesCache{
				"env-1": &userproto.UserAttributes{
					EnvironmentId: "env-1",
					UserAttributes: []*userproto.UserAttribute{
						{
							Key:    "country",
							Values: []string{"Japan", "USA"}, // Duplicate is not added
						},
						{
							Key:    "city",
							Values: []string{"Osaka"},
						},
					},
				},
			},
		},
		{
			name: "Process multiple environment IDs simultaneously",
			existingCache: userAttributesCache{
				"env-1": &userproto.UserAttributes{
					EnvironmentId: "env-1",
					UserAttributes: []*userproto.UserAttribute{
						{
							Key:    "country",
							Values: []string{"Japan"},
						},
					},
				},
			},
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId: "user-4",
						User: &userproto.User{
							Data: map[string]string{
								"city": "Tokyo",
							},
						},
					},
				},
				"env-2": eventMap{
					"event-2": &eventproto.EvaluationEvent{
						UserId: "user-5",
						User: &userproto.User{
							Data: map[string]string{
								"country": "USA",
								"city":    "New York",
							},
						},
					},
				},
			},
			expectedResult: userAttributesCache{
				"env-1": &userproto.UserAttributes{
					EnvironmentId: "env-1",
					UserAttributes: []*userproto.UserAttribute{
						{
							Key:    "country",
							Values: []string{"Japan"},
						},
						{
							Key:    "city",
							Values: []string{"Tokyo"},
						},
					},
				},
				"env-2": &userproto.UserAttributes{
					EnvironmentId: "env-2",
					UserAttributes: []*userproto.UserAttribute{
						{
							Key:    "country",
							Values: []string{"USA"},
						},
						{
							Key:    "city",
							Values: []string{"New York"},
						},
					},
				},
			},
		},
		{
			name:          "Skip empty keys",
			existingCache: make(userAttributesCache),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId: "user-6",
						User: &userproto.User{
							Data: map[string]string{
								"":        "empty-key-value", // Empty key
								"country": "Japan",
								"city":    "Tokyo",
							},
						},
					},
				},
			},
			expectedResult: userAttributesCache{
				"env-1": &userproto.UserAttributes{
					EnvironmentId: "env-1",
					UserAttributes: []*userproto.UserAttribute{
						{
							Key:    "country",
							Values: []string{"Japan"},
						},
						{
							Key:    "city",
							Values: []string{"Tokyo"},
						},
					},
				},
			},
		},
		{
			name:          "Skip when User is nil",
			existingCache: make(userAttributesCache),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId: "user-7",
						User:   nil, // User is nil
					},
				},
			},
			expectedResult: userAttributesCache{},
		},
		{
			name:          "Skip when User.Data is nil",
			existingCache: make(userAttributesCache),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId: "user-8",
						User: &userproto.User{
							Data: nil, // User.Data is nil
						},
					},
				},
			},
			expectedResult: userAttributesCache{},
		},
		{
			name:          "Add different values to the same key from multiple events",
			existingCache: make(userAttributesCache),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId: "user-9",
						User: &userproto.User{
							Data: map[string]string{
								"country": "Japan",
								"city":    "Tokyo",
							},
						},
					},
					"event-2": &eventproto.EvaluationEvent{
						UserId: "user-10",
						User: &userproto.User{
							Data: map[string]string{
								"country": "USA",
								"city":    "New York",
							},
						},
					},
					"event-3": &eventproto.EvaluationEvent{
						UserId: "user-11",
						User: &userproto.User{
							Data: map[string]string{
								"country": "Germany",
								"city":    "Berlin",
							},
						},
					},
				},
			},
			expectedResult: userAttributesCache{
				"env-1": &userproto.UserAttributes{
					EnvironmentId: "env-1",
					UserAttributes: []*userproto.UserAttribute{
						{
							Key:    "country",
							Values: []string{"Japan", "USA", "Germany"},
						},
						{
							Key:    "city",
							Values: []string{"Tokyo", "New York", "Berlin"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create persister for testing
			persister := &evaluationCountEventPersister{
				userAttributesCache: tt.existingCache,
				logger:              zap.NewNop(),
			}

			// Execute cacheUserAttributes
			persister.cacheUserAttributes(tt.envEvents)

			// Verify results
			assert.Equal(t, len(tt.expectedResult), len(persister.userAttributesCache))

			for envID, expected := range tt.expectedResult {
				actual, exists := persister.userAttributesCache[envID]
				require.True(t, exists, "Environment ID %s does not exist", envID)
				require.NotNil(t, actual, "UserAttributes for environment ID %s is nil", envID)

				assert.Equal(t, expected.EnvironmentId, actual.EnvironmentId)
				assert.Equal(t, len(expected.UserAttributes), len(actual.UserAttributes))

				// UserAttributes order is not guaranteed, so compare using maps
				expectedMap := make(map[string][]string)
				for _, attr := range expected.UserAttributes {
					expectedMap[attr.Key] = attr.Values
				}

				actualMap := make(map[string][]string)
				for _, attr := range actual.UserAttributes {
					actualMap[attr.Key] = attr.Values
				}

				// Compare maps without considering order
				for key, expectedValues := range expectedMap {
					actualValues, exists := actualMap[key]
					require.True(t, exists, "Key %s does not exist in actual result", key)
					require.Len(t, actualValues, len(expectedValues), "Values count mismatch for key %s", key)

					// Sort both slices for comparison
					expectedSorted := make([]string, len(expectedValues))
					copy(expectedSorted, expectedValues)
					actualSorted := make([]string, len(actualValues))
					copy(actualSorted, actualValues)

					sort.Strings(expectedSorted)
					sort.Strings(actualSorted)

					assert.Equal(t, expectedSorted, actualSorted, "Values mismatch for key %s", key)
				}
			}
		})
	}
}
