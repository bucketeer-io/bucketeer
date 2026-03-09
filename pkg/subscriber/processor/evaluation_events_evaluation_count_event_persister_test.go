package processor

import (
	"context"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
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

func toSet(ids ...string) map[string]struct{} {
	s := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		s[id] = struct{}{}
	}
	return s
}

func TestBufferDAU(t *testing.T) {
	day1 := int64(1772668800) // 2026-03-05 00:00:00 UTC
	day2 := int64(1772755200) // 2026-03-06 00:00:00 UTC

	tests := []struct {
		desc        string
		existingBuf dauBuffer
		envEvents   environmentEventMap
		expectedBuf dauBuffer
	}{
		{
			desc:        "Add a new entry to an empty buffer",
			existingBuf: make(dauBuffer),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId:    "user-1",
						SourceId:  eventproto.SourceId_ANDROID,
						Timestamp: day1,
					},
				},
			},
			expectedBuf: dauBuffer{
				{dateStr: "20260305", envID: "env-1", sourceID: "ANDROID"}: toSet("user-1"),
			},
		},
		{
			desc:        "Process events from multiple environments",
			existingBuf: make(dauBuffer),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId:    "user-1",
						SourceId:  eventproto.SourceId_ANDROID,
						Timestamp: day1,
					},
				},
				"env-2": eventMap{
					"event-2": &eventproto.EvaluationEvent{
						UserId:    "user-2",
						SourceId:  eventproto.SourceId_IOS,
						Timestamp: day1,
					},
				},
			},
			expectedBuf: dauBuffer{
				{dateStr: "20260305", envID: "env-1", sourceID: "ANDROID"}: toSet("user-1"),
				{dateStr: "20260305", envID: "env-2", sourceID: "IOS"}:     toSet("user-2"),
			},
		},
		{
			desc:        "Duplicate user IDs are deduplicated in buffer",
			existingBuf: make(dauBuffer),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId:    "user-1",
						SourceId:  eventproto.SourceId_ANDROID,
						Timestamp: day1,
					},
					"event-2": &eventproto.EvaluationEvent{
						UserId:    "user-1",
						SourceId:  eventproto.SourceId_ANDROID,
						Timestamp: day1 + 3600, // same day, different hour
					},
				},
			},
			expectedBuf: dauBuffer{
				{dateStr: "20260305", envID: "env-1", sourceID: "ANDROID"}: toSet("user-1"),
			},
		},
		{
			desc:        "Different dates produce separate entries",
			existingBuf: make(dauBuffer),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId:    "user-1",
						SourceId:  eventproto.SourceId_ANDROID,
						Timestamp: day1,
					},
					"event-2": &eventproto.EvaluationEvent{
						UserId:    "user-1",
						SourceId:  eventproto.SourceId_ANDROID,
						Timestamp: day2,
					},
				},
			},
			expectedBuf: dauBuffer{
				{dateStr: "20260305", envID: "env-1", sourceID: "ANDROID"}: toSet("user-1"),
				{dateStr: "20260306", envID: "env-1", sourceID: "ANDROID"}: toSet("user-1"),
			},
		},
		{
			desc: "Append to existing buffer entries",
			existingBuf: dauBuffer{
				{dateStr: "20260305", envID: "env-1", sourceID: "ANDROID"}: toSet("user-1"),
			},
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId:    "user-2",
						SourceId:  eventproto.SourceId_ANDROID,
						Timestamp: day1,
					},
				},
			},
			expectedBuf: dauBuffer{
				{dateStr: "20260305", envID: "env-1", sourceID: "ANDROID"}: toSet("user-1", "user-2"),
			},
		},
		{
			desc:        "Skip when userID is empty and User is nil",
			existingBuf: make(dauBuffer),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId:    "",
						User:      &userproto.User{Id: ""},
						SourceId:  eventproto.SourceId_ANDROID,
						Timestamp: day1,
					},
				},
			},
			expectedBuf: dauBuffer{},
		},
		{
			desc:        "Use User.Id when UserId is empty",
			existingBuf: make(dauBuffer),
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						UserId:    "",
						User:      &userproto.User{Id: "user-from-user-id"},
						SourceId:  eventproto.SourceId_ANDROID,
						Timestamp: day1,
					},
				},
			},
			expectedBuf: dauBuffer{
				{dateStr: "20260305", envID: "env-1", sourceID: "ANDROID"}: toSet("user-from-user-id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			persister := &evaluationCountEventPersister{
				dauBuf: tt.existingBuf,
				logger: zap.NewNop(),
			}
			persister.bufferDAU(tt.envEvents)

			assert.Equal(t, len(tt.expectedBuf), len(persister.dauBuf))
			for key, expectedUsers := range tt.expectedBuf {
				actualUsers, exists := persister.dauBuf[key]
				assert.True(t, exists, "expected key %v not found in dauBuf", key)
				assert.Equal(t, expectedUsers, actualUsers)
			}
		})
	}
}

func TestIsErrorReason(t *testing.T) {
	tests := []struct {
		name       string
		reason     *featureproto.Reason
		expectTrue bool
	}{
		{name: "nil reason", reason: nil, expectTrue: false},
		{
			name:       "Reason_CLIENT (deprecated)",
			reason:     &featureproto.Reason{Type: featureproto.Reason_CLIENT},
			expectTrue: true,
		},
		{
			name:       "Reason_ERROR_NO_EVALUATIONS",
			reason:     &featureproto.Reason{Type: featureproto.Reason_ERROR_NO_EVALUATIONS},
			expectTrue: true,
		},
		{
			name:       "Reason_ERROR_FLAG_NOT_FOUND",
			reason:     &featureproto.Reason{Type: featureproto.Reason_ERROR_FLAG_NOT_FOUND},
			expectTrue: true,
		},
		{
			name:       "Reason_ERROR_WRONG_TYPE",
			reason:     &featureproto.Reason{Type: featureproto.Reason_ERROR_WRONG_TYPE},
			expectTrue: true,
		},
		{
			name:       "Reason_ERROR_USER_ID_NOT_SPECIFIED",
			reason:     &featureproto.Reason{Type: featureproto.Reason_ERROR_USER_ID_NOT_SPECIFIED},
			expectTrue: true,
		},
		{
			name:       "Reason_ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED",
			reason:     &featureproto.Reason{Type: featureproto.Reason_ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED},
			expectTrue: true,
		},
		{
			name:       "Reason_ERROR_EXCEPTION",
			reason:     &featureproto.Reason{Type: featureproto.Reason_ERROR_EXCEPTION},
			expectTrue: true,
		},
		{
			name:       "Reason_ERROR_CACHE_NOT_FOUND",
			reason:     &featureproto.Reason{Type: featureproto.Reason_ERROR_CACHE_NOT_FOUND},
			expectTrue: true,
		},
		{
			name:       "Reason_TARGET",
			reason:     &featureproto.Reason{Type: featureproto.Reason_TARGET},
			expectTrue: false,
		},
		{
			name:       "Reason_RULE",
			reason:     &featureproto.Reason{Type: featureproto.Reason_RULE},
			expectTrue: false,
		},
		{
			name:       "Reason_DEFAULT",
			reason:     &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
			expectTrue: false,
		},
		{
			name:       "Reason_OFF_VARIATION",
			reason:     &featureproto.Reason{Type: featureproto.Reason_OFF_VARIATION},
			expectTrue: false,
		},
		{
			name:       "Reason_PREREQUISITE",
			reason:     &featureproto.Reason{Type: featureproto.Reason_PREREQUISITE},
			expectTrue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isErrorReason(tt.reason)
			assert.Equal(t, tt.expectTrue, got, "isErrorReason must stay in sync with grpc_validation.go")
		})
	}
}

// TestIsErrorReasonCoversAllProtoErrorTypes fails when a new error reason type is added to the proto
// (CLIENT or ERROR_* naming) but isErrorReason in the persister hasn't been updated.
// This forces us to update isErrorReason and grpc_validation.go when adding new error types.
func TestIsErrorReasonCoversAllProtoErrorTypes(t *testing.T) {
	for value, name := range featureproto.Reason_Type_name {
		isErrorType := name == "CLIENT" || strings.HasPrefix(name, "ERROR_")
		if !isErrorType {
			continue
		}
		reasonType := featureproto.Reason_Type(value)
		reason := &featureproto.Reason{Type: reasonType}
		assert.True(t, isErrorReason(reason),
			"Reason %s (value=%d) is an error type per proto naming but isErrorReason returns false. "+
				"Update isErrorReason in evaluation_events_evaluation_count_event_persister.go and "+
				"isErrorReason in grpc_validation.go to include this type.",
			name, value)
	}
}

func TestGetVariationID(t *testing.T) {
	tests := []struct {
		name        string
		reason      *featureproto.Reason
		vID         string
		expectedVID string
		expectErr   bool
	}{
		{
			name:      "nil reason returns error",
			reason:    nil,
			vID:       "variation-1",
			expectErr: true,
		},
		{
			name:        "Reason_CLIENT returns default (deprecated)",
			reason:      &featureproto.Reason{Type: featureproto.Reason_CLIENT},
			vID:         "variation-1",
			expectedVID: defaultVariationID,
		},
		{
			name:        "Reason_ERROR_NO_EVALUATIONS returns default",
			reason:      &featureproto.Reason{Type: featureproto.Reason_ERROR_NO_EVALUATIONS},
			vID:         "",
			expectedVID: defaultVariationID,
		},
		{
			name:        "Reason_ERROR_FLAG_NOT_FOUND returns default",
			reason:      &featureproto.Reason{Type: featureproto.Reason_ERROR_FLAG_NOT_FOUND},
			vID:         "",
			expectedVID: defaultVariationID,
		},
		{
			name:        "Reason_ERROR_WRONG_TYPE returns default",
			reason:      &featureproto.Reason{Type: featureproto.Reason_ERROR_WRONG_TYPE},
			vID:         "",
			expectedVID: defaultVariationID,
		},
		{
			name:        "Reason_ERROR_USER_ID_NOT_SPECIFIED returns default",
			reason:      &featureproto.Reason{Type: featureproto.Reason_ERROR_USER_ID_NOT_SPECIFIED},
			vID:         "",
			expectedVID: defaultVariationID,
		},
		{
			name:        "Reason_ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED returns default",
			reason:      &featureproto.Reason{Type: featureproto.Reason_ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED},
			vID:         "",
			expectedVID: defaultVariationID,
		},
		{
			name:        "Reason_ERROR_EXCEPTION returns default",
			reason:      &featureproto.Reason{Type: featureproto.Reason_ERROR_EXCEPTION},
			vID:         "",
			expectedVID: defaultVariationID,
		},
		{
			name:        "Reason_ERROR_CACHE_NOT_FOUND returns default",
			reason:      &featureproto.Reason{Type: featureproto.Reason_ERROR_CACHE_NOT_FOUND},
			vID:         "",
			expectedVID: defaultVariationID,
		},
		{
			name:        "error reason with non-empty vID still returns default (overrides client value)",
			reason:      &featureproto.Reason{Type: featureproto.Reason_ERROR_FLAG_NOT_FOUND},
			vID:         "variation-A",
			expectedVID: defaultVariationID,
		},
		{
			name:        "Reason_TARGET returns actual variation id",
			reason:      &featureproto.Reason{Type: featureproto.Reason_TARGET},
			vID:         "variation-1",
			expectedVID: "variation-1",
		},
		{
			name:        "Reason_RULE returns actual variation id",
			reason:      &featureproto.Reason{Type: featureproto.Reason_RULE},
			vID:         "variation-2",
			expectedVID: "variation-2",
		},
		{
			name:        "Reason_DEFAULT returns actual variation id",
			reason:      &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
			vID:         "variation-3",
			expectedVID: "variation-3",
		},
		{
			name:        "Reason_OFF_VARIATION returns actual variation id",
			reason:      &featureproto.Reason{Type: featureproto.Reason_OFF_VARIATION},
			vID:         "variation-4",
			expectedVID: "variation-4",
		},
		{
			name:        "Reason_PREREQUISITE returns actual variation id",
			reason:      &featureproto.Reason{Type: featureproto.Reason_PREREQUISITE},
			vID:         "variation-5",
			expectedVID: "variation-5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getVariationID(tt.reason, tt.vID)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, ErrReasonNil, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedVID, got, "getVariationID should return expected variation ID for default value counter")
		})
	}
}

func TestIncrementEnvEvents_Aggregation(t *testing.T) {
	t.Parallel()

	hour1 := int64(1709974800) // 2024-03-09 09:00:00 UTC

	tests := []struct {
		name                   string
		envEvents              environmentEventMap
		expectedEventCountKeys int // number of unique event count keys
		expectedUserCountKeys  int // number of unique user count keys
		expectedFailCount      int // number of expected failures
		expectFlushCalled      bool
		simulateFlushError     bool
	}{
		{
			name: "single event creates one aggregated entry",
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						FeatureId:      "feature-1",
						VariationId:    "variation-A",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
						UserId:         "user-1",
						User:           &userproto.User{Id: "user-1"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
				},
			},
			expectedEventCountKeys: 1,
			expectedUserCountKeys:  1,
			expectedFailCount:      0,
			expectFlushCalled:      true,
		},
		{
			name: "multiple events same key aggregate",
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						FeatureId:      "feature-1",
						VariationId:    "variation-A",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
						UserId:         "user-1",
						User:           &userproto.User{Id: "user-1"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
					"event-2": &eventproto.EvaluationEvent{
						FeatureId:      "feature-1",
						VariationId:    "variation-A",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
						UserId:         "user-2",
						User:           &userproto.User{Id: "user-2"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
					"event-3": &eventproto.EvaluationEvent{
						FeatureId:      "feature-1",
						VariationId:    "variation-A",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
						UserId:         "user-3",
						User:           &userproto.User{Id: "user-3"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
				},
			},
			expectedEventCountKeys: 1, // all 3 events share same key
			expectedUserCountKeys:  1, // all 3 events share same key
			expectedFailCount:      0,
			expectFlushCalled:      true,
		},
		{
			name: "different variations create separate keys",
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						FeatureId:      "feature-1",
						VariationId:    "variation-A",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
						UserId:         "user-1",
						User:           &userproto.User{Id: "user-1"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
					"event-2": &eventproto.EvaluationEvent{
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
						UserId:         "user-2",
						User:           &userproto.User{Id: "user-2"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
				},
			},
			expectedEventCountKeys: 2, // different variations
			expectedUserCountKeys:  2,
			expectedFailCount:      0,
			expectFlushCalled:      true,
		},
		{
			name: "error reasons map to default variation",
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						FeatureId:      "feature-1",
						VariationId:    "variation-A",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_ERROR_FLAG_NOT_FOUND},
						UserId:         "user-1",
						User:           &userproto.User{Id: "user-1"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
					"event-2": &eventproto.EvaluationEvent{
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_ERROR_CACHE_NOT_FOUND},
						UserId:         "user-2",
						User:           &userproto.User{Id: "user-2"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
				},
			},
			expectedEventCountKeys: 1, // both map to "default" variation
			expectedUserCountKeys:  1,
			expectedFailCount:      0,
			expectFlushCalled:      true,
		},
		{
			name: "nil reason events are marked as failed",
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						FeatureId:      "feature-1",
						VariationId:    "variation-A",
						Reason:         nil, // nil reason
						UserId:         "user-1",
						User:           &userproto.User{Id: "user-1"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
				},
			},
			expectedEventCountKeys: 0, // not added to aggregator
			expectedUserCountKeys:  0,
			expectedFailCount:      1,     // marked as non-repeatable error
			expectFlushCalled:      false, // flush not called when no events to aggregate
		},
		{
			name: "flush error marks all events as failed",
			envEvents: environmentEventMap{
				"env-1": eventMap{
					"event-1": &eventproto.EvaluationEvent{
						FeatureId:      "feature-1",
						VariationId:    "variation-A",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
						UserId:         "user-1",
						User:           &userproto.User{Id: "user-1"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
				},
			},
			// Note: ec and uc keys may hash to different slots, so on error we might
			// only attempt one slot before returning. Just verify flush was called.
			expectedEventCountKeys: 0, // Don't verify exact keys when error occurs
			expectedUserCountKeys:  0, // Don't verify exact keys when error occurs
			expectedFailCount:      1,
			expectFlushCalled:      true,
			simulateFlushError:     true, // simulate Redis failure
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock cache
			mockCache := &mockEvaluationCountCache{
				shouldFailFlush: tt.simulateFlushError,
			}

			persister := &evaluationCountEventPersister{
				evaluationCountCacher: mockCache,
				logger:                zap.NewNop(),
			}

			// Execute
			fails := persister.incrementEnvEvents(tt.envEvents)

			// Verify
			assert.Equal(t, tt.expectedFailCount, len(fails), "fail count mismatch")

			if tt.expectFlushCalled {
				assert.True(t, mockCache.flushCalled, "flush should be called")
				// Only verify key counts when not simulating errors
				// (errors may cause early return before all slots are attempted)
				if !tt.simulateFlushError {
					assert.Equal(t, tt.expectedEventCountKeys, len(mockCache.lastEventCounts), "event count keys mismatch")
					assert.Equal(t, tt.expectedUserCountKeys, len(mockCache.lastUserCounts), "user count keys mismatch")
				}
			}
		})
	}
}

func TestFlushAggregatedCounts_SlotGrouping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		eventCounts          map[string]int64
		userCounts           map[string]map[string]struct{}
		expectedPipelineExec int // Expected number of pipeline Exec() calls (one per unique slot); -1 = calculate dynamically
		description          string
	}{
		{
			name: "keys with same hash tag in both ec and uc",
			eventCounts: map[string]int64{
				"ec:{shared}:hour1:feature1:varA:env1": 100,
			},
			userCounts: map[string]map[string]struct{}{
				"uc:{shared}:hour1:feature1:varA:env1": {"user1": {}, "user2": {}},
			},
			expectedPipelineExec: 1,
			description:          "Hash tag forces both ec and uc keys to same slot",
		},
		{
			name: "ec and uc keys without hash tags may hash to different slots",
			eventCounts: map[string]int64{
				"ec:hour1:feature1:varA:env1": 100,
			},
			userCounts: map[string]map[string]struct{}{
				"uc:hour1:feature1:varA:env1": {"user1": {}, "user2": {}},
			},
			expectedPipelineExec: -1, // Calculate dynamically based on actual slot distribution
			description:          "Without hash tags, ec and uc prefixes may cause different slots",
		},
		{
			name: "keys with same hash tag execute in one pipeline",
			eventCounts: map[string]int64{
				"ec:{slot1}:key1": 10,
				"ec:{slot1}:key2": 20,
				"ec:{slot1}:key3": 30,
			},
			userCounts: map[string]map[string]struct{}{
				"uc:{slot1}:key1": {"user1": {}},
				"uc:{slot1}:key2": {"user2": {}},
				"uc:{slot1}:key3": {"user3": {}},
			},
			expectedPipelineExec: 1,
			description:          "All keys use {slot1} hash tag, forcing same slot",
		},
		{
			name: "keys with different hash tags execute in multiple pipelines",
			eventCounts: map[string]int64{
				"ec:{slotA}:key1": 10,
				"ec:{slotB}:key2": 20,
				"ec:{slotC}:key3": 30,
			},
			userCounts: map[string]map[string]struct{}{
				"uc:{slotA}:key1": {"user1": {}},
				"uc:{slotB}:key2": {"user2": {}},
				"uc:{slotC}:key3": {"user3": {}},
			},
			expectedPipelineExec: 3,
			description:          "Different hash tags create different slots",
		},
		{
			name: "mixed hash tags create appropriate number of pipelines",
			eventCounts: map[string]int64{
				"ec:{groupA}:key1": 10,
				"ec:{groupA}:key2": 20,
				"ec:{groupB}:key3": 30,
				"ec:{groupB}:key4": 40,
			},
			userCounts: map[string]map[string]struct{}{
				"uc:{groupA}:key1": {"user1": {}},
				"uc:{groupA}:key2": {"user2": {}},
				"uc:{groupB}:key3": {"user3": {}},
				"uc:{groupB}:key4": {"user4": {}},
			},
			expectedPipelineExec: 2,
			description:          "Two groups with same hash tags create 2 pipelines",
		},
		{
			name:                 "empty aggregation creates no pipelines",
			eventCounts:          map[string]int64{},
			userCounts:           map[string]map[string]struct{}{},
			expectedPipelineExec: 0,
			description:          "No data to flush means no pipelines created",
		},
		{
			name: "real-world keys may span multiple slots",
			eventCounts: map[string]int64{
				// These keys don't have hash tags, so they will hash based on full key
				// They will likely hash to different slots
				"ec:1709974800:feature-login:variant-A:env-prod":    500,
				"ec:1709974800:feature-checkout:variant-B:env-prod": 300,
				"ec:1709974800:feature-sidebar:variant-A:env-prod":  200,
			},
			userCounts: map[string]map[string]struct{}{
				"uc:1709974800:feature-login:variant-A:env-prod":    {"user1": {}, "user2": {}},
				"uc:1709974800:feature-checkout:variant-B:env-prod": {"user3": {}},
				"uc:1709974800:feature-sidebar:variant-A:env-prod":  {"user4": {}},
			},
			// Don't hardcode expected count - calculate it dynamically
			expectedPipelineExec: -1, // Special value: calculate based on actual slot distribution
			description:          "Real keys distribute across slots based on CRC16",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Calculate expected pipeline count for real-world scenario
			expectedExecs := tt.expectedPipelineExec
			if expectedExecs == -1 {
				// Calculate unique slots for this test case
				slotSet := make(map[int]bool)
				for key := range tt.eventCounts {
					slotSet[redisv3.KeyHashSlot(key)] = true
				}
				for key := range tt.userCounts {
					slotSet[redisv3.KeyHashSlot(key)] = true
				}
				expectedExecs = len(slotSet)
			}

			mockCache := &mockEvaluationCountCache{}
			persister := &evaluationCountEventPersister{
				evaluationCountCacher: mockCache,
				logger:                zap.NewNop(),
			}

			// Execute flush
			err := persister.flushAggregatedCounts(tt.eventCounts, tt.userCounts)

			// Verify
			if len(tt.eventCounts) == 0 && len(tt.userCounts) == 0 {
				assert.NoError(t, err, "empty flush should succeed")
				assert.Equal(t, 0, mockCache.execCount, "no pipelines should be executed for empty data")
			} else {
				assert.NoError(t, err, "flush should succeed")
				assert.Equal(t, expectedExecs, mockCache.execCount,
					"pipeline Exec() count should match number of unique slots: %s", tt.description)

				// Verify all data was flushed
				assert.Equal(t, len(tt.eventCounts), len(mockCache.lastEventCounts),
					"all event count keys should be flushed")
				assert.Equal(t, len(tt.userCounts), len(mockCache.lastUserCounts),
					"all user count keys should be flushed")

				// Verify data integrity
				for key, expectedCount := range tt.eventCounts {
					assert.Equal(t, expectedCount, mockCache.lastEventCounts[key],
						"event count for key %s should match", key)
				}
				for key, expectedUsers := range tt.userCounts {
					assert.Equal(t, len(expectedUsers), len(mockCache.lastUserCounts[key]),
						"user count for key %s should match", key)
				}
			}
		})
	}
}

// mockEvaluationCountCache mocks the cache interface for testing
type mockEvaluationCountCache struct {
	mu               sync.Mutex
	flushCalled      bool
	lastEventCounts  map[string]int64
	lastUserCounts   map[string]map[string]struct{}
	shouldFailFlush  bool
	pipelineExecuted bool
	pipelineCount    int // Number of pipelines created
	execCount        int // Number of times Exec() was called
}

func (m *mockEvaluationCountCache) Pipeline(tx bool) redisv3.PipeClient {
	m.mu.Lock()
	m.pipelineCount++
	m.mu.Unlock()

	return &mockPipeClient{
		cache:          m,
		transactional:  tx,
		commands:       []string{},
		eventCounts:    make(map[string]int64),
		userCounts:     make(map[string][]string),
		shouldFailExec: m.shouldFailFlush,
	}
}

// Cache interface methods
func (m *mockEvaluationCountCache) Get(key interface{}) (interface{}, error) { return nil, nil }
func (m *mockEvaluationCountCache) Put(key interface{}, value interface{}, expiration time.Duration) error {
	return nil
}

// MultiGetter interface methods
func (m *mockEvaluationCountCache) GetMulti(keys interface{}, ignoreNotFound bool) ([]interface{}, error) {
	return nil, nil
}
func (m *mockEvaluationCountCache) Scan(cursor, key, count interface{}) (uint64, []string, error) {
	return 0, nil, nil
}
func (m *mockEvaluationCountCache) SMembers(key string) ([]string, error) { return nil, nil }

// Deleter interface methods
func (m *mockEvaluationCountCache) Delete(key string) error { return nil }

// Counter interface methods
func (m *mockEvaluationCountCache) Increment(key string) (int64, error) { return 0, nil }
func (m *mockEvaluationCountCache) PFAdd(key string, els ...string) (int64, error) {
	return 0, nil
}

// PFGetter interface methods
func (m *mockEvaluationCountCache) PFCount(keys ...string) (int64, error) { return 0, nil }

// PFMerger interface methods
func (m *mockEvaluationCountCache) PFMerge(dest string, expiration time.Duration, keys ...string) error {
	return nil
}

// Expirer interface methods
func (m *mockEvaluationCountCache) Expire(key string, expiration time.Duration) (bool, error) {
	return false, nil
}

type mockPipeClient struct {
	cache          *mockEvaluationCountCache
	transactional  bool
	commands       []string
	eventCounts    map[string]int64
	userCounts     map[string][]string
	shouldFailExec bool
}

func (m *mockPipeClient) IncrBy(key string, value int64) *goredis.IntCmd {
	m.commands = append(m.commands, "INCRBY")
	m.eventCounts[key] = value
	return goredis.NewIntCmd(context.Background())
}

func (m *mockPipeClient) PFAdd(key string, els ...string) *goredis.IntCmd {
	m.commands = append(m.commands, "PFADD")
	m.userCounts[key] = els
	return goredis.NewIntCmd(context.Background())
}

func (m *mockPipeClient) Exec() ([]goredis.Cmder, error) {
	m.cache.mu.Lock()
	defer m.cache.mu.Unlock()

	m.cache.flushCalled = true
	m.cache.pipelineExecuted = true
	m.cache.execCount++

	// Accumulate across multiple pipeline calls (one per slot)
	// Always track what was attempted, even if exec fails
	if m.cache.lastEventCounts == nil {
		m.cache.lastEventCounts = make(map[string]int64)
	}
	if m.cache.lastUserCounts == nil {
		m.cache.lastUserCounts = make(map[string]map[string]struct{})
	}

	// Merge event counts
	for key, count := range m.eventCounts {
		m.cache.lastEventCounts[key] += count
	}

	// Merge user counts
	for key, users := range m.userCounts {
		if m.cache.lastUserCounts[key] == nil {
			m.cache.lastUserCounts[key] = make(map[string]struct{})
		}
		for _, user := range users {
			m.cache.lastUserCounts[key][user] = struct{}{}
		}
	}

	// Check for error after tracking (to verify what was attempted)
	if m.shouldFailExec {
		return nil, assert.AnError
	}
	return nil, nil
}

// Unused pipeline methods
func (m *mockPipeClient) Incr(key string) *goredis.IntCmd {
	return goredis.NewIntCmd(context.Background())
}
func (m *mockPipeClient) TTL(key string) *goredis.DurationCmd {
	return goredis.NewDurationCmd(context.Background(), 0)
}
func (m *mockPipeClient) SAdd(key string, members ...interface{}) *goredis.IntCmd {
	return goredis.NewIntCmd(context.Background())
}
func (m *mockPipeClient) Expire(key string, expiration time.Duration) *goredis.BoolCmd {
	return goredis.NewBoolCmd(context.Background())
}
func (m *mockPipeClient) PFCount(keys ...string) *goredis.IntCmd {
	return goredis.NewIntCmd(context.Background())
}
func (m *mockPipeClient) Get(key string) *goredis.StringCmd {
	return goredis.NewStringCmd(context.Background())
}
func (m *mockPipeClient) Del(keys string) *goredis.IntCmd {
	return goredis.NewIntCmd(context.Background())
}
