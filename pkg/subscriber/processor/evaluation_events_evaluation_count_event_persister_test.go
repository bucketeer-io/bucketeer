package processor

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
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
			// Note: With PFADD-first ordering, if PFADD fails, INCRBY never executes.
			// We return immediately on first failure.
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

func TestFlushAggregatedCounts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		eventCounts     map[string]int64
		userCounts      map[string]map[string]struct{}
		expectedPFAdds  int // Expected number of direct PFAdd calls
		expectedIncrBys int // Expected number of direct IncrementBy calls
		description     string
	}{
		{
			name: "single key pair - PFADD before INCRBY",
			eventCounts: map[string]int64{
				"env1:ec:hour1:feature1:varA": 100,
			},
			userCounts: map[string]map[string]struct{}{
				"env1:uc:hour1:feature1:varA": {"user1": {}, "user2": {}},
			},
			expectedPFAdds:  1,
			expectedIncrBys: 1,
			description:     "One ec/uc pair: PFADD first, then INCRBY",
		},
		{
			name: "multiple key pairs - PFADD before INCRBY for each",
			eventCounts: map[string]int64{
				"env1:ec:hour1:feature1:varA": 10,
				"env1:ec:hour1:feature2:varB": 20,
				"env1:ec:hour1:feature3:varC": 30,
			},
			userCounts: map[string]map[string]struct{}{
				"env1:uc:hour1:feature1:varA": {"user1": {}},
				"env1:uc:hour1:feature2:varB": {"user2": {}},
				"env1:uc:hour1:feature3:varC": {"user3": {}},
			},
			expectedPFAdds:  3,
			expectedIncrBys: 3,
			description:     "Three pairs: each gets PFADD then INCRBY",
		},
		{
			name:            "empty aggregation - no calls",
			eventCounts:     map[string]int64{},
			userCounts:      map[string]map[string]struct{}{},
			expectedPFAdds:  0,
			expectedIncrBys: 0,
			description:     "No data to flush means no calls",
		},
		{
			name: "real-world keys",
			eventCounts: map[string]int64{
				"env-prod:ec:1709974800:feature-login:variant-A":    500,
				"env-prod:ec:1709974800:feature-checkout:variant-B": 300,
				"env-prod:ec:1709974800:feature-sidebar:variant-A":  200,
			},
			userCounts: map[string]map[string]struct{}{
				"env-prod:uc:1709974800:feature-login:variant-A":    {"user1": {}, "user2": {}},
				"env-prod:uc:1709974800:feature-checkout:variant-B": {"user3": {}},
				"env-prod:uc:1709974800:feature-sidebar:variant-A":  {"user4": {}},
			},
			expectedPFAdds:  3,
			expectedIncrBys: 3,
			description:     "Three key pairs with realistic keys",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := &mockEvaluationCountCache{}
			persister := &evaluationCountEventPersister{
				evaluationCountCacher: mockCache,
				logger:                zap.NewNop(),
			}

			// Execute flush
			failedKeys, err := persister.flushAggregatedCounts(tt.eventCounts, tt.userCounts)

			// Verify
			assert.NoError(t, err, "flush should succeed")
			assert.Empty(t, failedKeys, "no keys should fail")

			// Verify call counts
			assert.Equal(t, tt.expectedPFAdds, mockCache.pfaddCallCount,
				"PFAdd call count should match: %s", tt.description)
			assert.Equal(t, tt.expectedIncrBys, mockCache.execCount,
				"IncrBy call count should match: %s", tt.description)

			// Verify PFADD-before-INCRBY ordering per key pair
			if len(tt.userCounts) > 0 {
				extractSuffix := func(call string) string {
					s := call
					if strings.HasPrefix(s, "PFADD:") {
						s = strings.TrimPrefix(s, "PFADD:")
						s = strings.Replace(s, ":uc:", ":", 1)
					} else if strings.HasPrefix(s, "INCRBY:") {
						s = strings.TrimPrefix(s, "INCRBY:")
						s = strings.Replace(s, ":ec:", ":", 1)
					}
					return s
				}

				pfaddIdxBySuffix := make(map[string]int)
				incrbyIdxBySuffix := make(map[string]int)
				for idx, call := range mockCache.incrbyCallOrder {
					suffix := extractSuffix(call)
					if strings.HasPrefix(call, "PFADD:") {
						if _, seen := pfaddIdxBySuffix[suffix]; !seen {
							pfaddIdxBySuffix[suffix] = idx
						}
					} else if strings.HasPrefix(call, "INCRBY:") {
						if _, seen := incrbyIdxBySuffix[suffix]; !seen {
							incrbyIdxBySuffix[suffix] = idx
						}
					}
				}

				for suffix, pfIdx := range pfaddIdxBySuffix {
					if inIdx, ok := incrbyIdxBySuffix[suffix]; ok {
						assert.Less(t, pfIdx, inIdx,
							"PFADD must come before INCRBY for key pair suffix %s", suffix)
					}
				}
			}

			// Verify data integrity
			if len(tt.eventCounts) > 0 || len(tt.userCounts) > 0 {
				assert.Equal(t, len(tt.eventCounts), len(mockCache.lastEventCounts),
					"all event count keys should be flushed")
				assert.Equal(t, len(tt.userCounts), len(mockCache.lastUserCounts),
					"all user count keys should be flushed")

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

func TestFlushAggregatedCounts_Failures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		eventCounts         map[string]int64
		userCounts          map[string]map[string]struct{}
		shouldFailPFAdd     bool
		shouldFailIncrBy    bool
		expectedSuccess     bool
		expectedPFAddCalls  int
		expectedIncrByCalls int
		description         string
	}{
		{
			name: "all operations succeed",
			eventCounts: map[string]int64{
				"env1:ec:hour1:feature1:varA": 10,
				"env1:ec:hour1:feature2:varB": 20,
			},
			userCounts: map[string]map[string]struct{}{
				"env1:uc:hour1:feature1:varA": {"user1": {}},
				"env1:uc:hour1:feature2:varB": {"user2": {}},
			},
			shouldFailPFAdd:     false,
			shouldFailIncrBy:    false,
			expectedSuccess:     true,
			expectedPFAddCalls:  2,
			expectedIncrByCalls: 2,
			description:         "No failures, all operations complete",
		},
		{
			name: "PFADD fails - INCRBY never executes (no over-count on retry)",
			eventCounts: map[string]int64{
				"env1:ec:hour1:feature1:varA": 10,
			},
			userCounts: map[string]map[string]struct{}{
				"env1:uc:hour1:feature1:varA": {"user1": {}},
			},
			shouldFailPFAdd:     true,
			shouldFailIncrBy:    false,
			expectedSuccess:     false,
			expectedPFAddCalls:  1,
			expectedIncrByCalls: 0, // INCRBY never called because PFADD failed
			description: "PFADD failure prevents INCRBY execution; " +
				"retry is safe",
		},
		{
			name: "INCRBY fails after PFADD succeeds (safe retry due to idempotency)",
			eventCounts: map[string]int64{
				"env1:ec:hour1:feature1:varA": 10,
			},
			userCounts: map[string]map[string]struct{}{
				"env1:uc:hour1:feature1:varA": {"user1": {}},
			},
			shouldFailPFAdd:     false,
			shouldFailIncrBy:    true,
			expectedSuccess:     false,
			expectedPFAddCalls:  1,
			expectedIncrByCalls: 1, // INCRBY attempted but failed
			description: "INCRBY fails but PFADD succeeded; " +
				"retry safe (PFADD idempotent)",
		},
		{
			name:                "empty data succeeds",
			eventCounts:         map[string]int64{},
			userCounts:          map[string]map[string]struct{}{},
			shouldFailPFAdd:     false,
			shouldFailIncrBy:    false,
			expectedSuccess:     true,
			expectedPFAddCalls:  0,
			expectedIncrByCalls: 0,
			description:         "Empty flush succeeds without any calls",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := &mockEvaluationCountCache{
				shouldFailPFAdd:  tt.shouldFailPFAdd,
				shouldFailIncrBy: tt.shouldFailIncrBy,
			}
			persister := &evaluationCountEventPersister{
				evaluationCountCacher: mockCache,
				logger:                zap.NewNop(),
			}

			// Execute flush
			failedKeys, err := persister.flushAggregatedCounts(tt.eventCounts, tt.userCounts)

			// Verify success/failure
			if tt.expectedSuccess {
				assert.NoError(t, err, tt.description)
				assert.Empty(t, failedKeys, tt.description)
			} else {
				assert.Error(t, err, tt.description)
				assert.NotEmpty(t, failedKeys, tt.description)
			}

			// Verify call counts match expectations
			assert.Equal(t, tt.expectedPFAddCalls, mockCache.pfaddCallCount,
				"PFAdd call count: %s", tt.description)
			assert.Equal(t, tt.expectedIncrByCalls, mockCache.execCount,
				"IncrBy call count: %s", tt.description)

			// Verify PFADD-before-INCRBY ordering per key pair when both are called
			if tt.expectedPFAddCalls > 0 && tt.expectedIncrByCalls > 0 {
				extractSuffix := func(call string) string {
					s := call
					if strings.HasPrefix(s, "PFADD:") {
						s = strings.TrimPrefix(s, "PFADD:")
						s = strings.Replace(s, ":uc:", ":", 1)
					} else if strings.HasPrefix(s, "INCRBY:") {
						s = strings.TrimPrefix(s, "INCRBY:")
						s = strings.Replace(s, ":ec:", ":", 1)
					}
					return s
				}

				pfaddIdxBySuffix := make(map[string]int)
				incrbyIdxBySuffix := make(map[string]int)
				for idx, call := range mockCache.incrbyCallOrder {
					suffix := extractSuffix(call)
					if strings.HasPrefix(call, "PFADD:") {
						if _, seen := pfaddIdxBySuffix[suffix]; !seen {
							pfaddIdxBySuffix[suffix] = idx
						}
					} else if strings.HasPrefix(call, "INCRBY:") {
						if _, seen := incrbyIdxBySuffix[suffix]; !seen {
							incrbyIdxBySuffix[suffix] = idx
						}
					}
				}

				for suffix, pfIdx := range pfaddIdxBySuffix {
					if inIdx, ok := incrbyIdxBySuffix[suffix]; ok {
						assert.Less(t, pfIdx, inIdx,
							"PFADD must come before INCRBY for key pair suffix %s", suffix)
					}
				}
			}
		})
	}
}

func TestIncrementEnvEvents_Retry(t *testing.T) {
	t.Parallel()

	hour1 := int64(1709974800) // 2024-03-09 09:00:00 UTC

	tests := []struct {
		name              string
		envEvents         environmentEventMap
		shouldFail        bool
		expectedFailCount int
		description       string
	}{
		{
			name: "all events succeed",
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
						FeatureId:      "feature-2",
						VariationId:    "variation-B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
						UserId:         "user-2",
						User:           &userproto.User{Id: "user-2"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
				},
			},
			shouldFail:        false,
			expectedFailCount: 0,
			description:       "All events succeed, no failures",
		},
		{
			name: "flush failure causes all events to retry",
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
						FeatureId:      "feature-2",
						VariationId:    "variation-B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
						UserId:         "user-2",
						User:           &userproto.User{Id: "user-2"},
						Timestamp:      hour1,
						FeatureVersion: 1,
					},
				},
			},
			shouldFail:        true,
			expectedFailCount: 2,
			description:       "All events marked for retry on flush failure",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := &mockEvaluationCountCache{
				shouldFailFlush: tt.shouldFail,
			}
			persister := &evaluationCountEventPersister{
				evaluationCountCacher: mockCache,
				logger:                zap.NewNop(),
			}

			// Execute
			fails := persister.incrementEnvEvents(tt.envEvents)

			// Verify
			assert.Equal(t, tt.expectedFailCount, len(fails), tt.description)
			if tt.expectedFailCount > 0 {
				// All failed events should be repeatable
				for _, repeatable := range fails {
					assert.True(t, repeatable, "all failed events should be repeatable")
				}
			}
		})
	}
}

func TestIncrementEnvEvents_PartialFlushFailure(t *testing.T) {
	t.Parallel()

	hour1 := int64(1709974800) // 2024-03-09 09:00:00 UTC

	// Two events mapping to different key pairs (different features).
	// Only one key pair's INCRBY will fail. The other event should NOT be retried.
	envEvents := environmentEventMap{
		"env-1": eventMap{
			"event-success": &eventproto.EvaluationEvent{
				FeatureId:      "feature-ok",
				VariationId:    "variation-A",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-1",
				User:           &userproto.User{Id: "user-1"},
				Timestamp:      hour1,
				FeatureVersion: 1,
			},
			"event-fail": &eventproto.EvaluationEvent{
				FeatureId:      "feature-bad",
				VariationId:    "variation-B",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-2",
				User:           &userproto.User{Id: "user-2"},
				Timestamp:      hour1,
				FeatureVersion: 1,
			},
		},
	}

	// Build the ecKey that will fail so we can set failOnECKey.
	// Replicates newEvaluationCountkeyV2 logic.
	tTime := time.Unix(hour1, 0)
	date := time.Date(tTime.Year(), tTime.Month(), tTime.Day(), tTime.Hour(), 0, 0, 0, time.UTC)
	failECKey := fmt.Sprintf("env-1:%s:%d:%s:%s", eventCountKey, date.Unix(), "feature-bad", "variation-B")

	mockCache := &mockEvaluationCountCache{
		failOnECKey: failECKey,
	}
	persister := &evaluationCountEventPersister{
		evaluationCountCacher: mockCache,
		logger:                zap.NewNop(),
	}

	fails := persister.incrementEnvEvents(envEvents)

	// Only the event that maps to the failed ecKey should be retried
	assert.Contains(t, fails, "event-fail", "event mapping to failed ecKey must be retried")
	assert.True(t, fails["event-fail"], "failed event must be repeatable")
	assert.NotContains(t, fails, "event-success",
		"event mapping to succeeded ecKey must NOT be retried (would cause over-count)")
}

func TestIncrementEnvEvents_SharedKeyPairFailure(t *testing.T) {
	t.Parallel()

	hour1 := int64(1709974800)

	// Two events that map to the SAME key pair (same feature/variation/env/hour).
	// When that shared key pair fails, both events should be retried.
	envEvents := environmentEventMap{
		"env-1": eventMap{
			"event-A": &eventproto.EvaluationEvent{
				FeatureId:      "feature-shared",
				VariationId:    "variation-A",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-1",
				User:           &userproto.User{Id: "user-1"},
				Timestamp:      hour1,
				FeatureVersion: 1,
			},
			"event-B": &eventproto.EvaluationEvent{
				FeatureId:      "feature-shared",
				VariationId:    "variation-A",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-2",
				User:           &userproto.User{Id: "user-2"},
				Timestamp:      hour1,
				FeatureVersion: 1,
			},
		},
	}

	tTime := time.Unix(hour1, 0)
	date := time.Date(tTime.Year(), tTime.Month(), tTime.Day(), tTime.Hour(), 0, 0, 0, time.UTC)
	failECKey := fmt.Sprintf("env-1:%s:%d:%s:%s", eventCountKey, date.Unix(), "feature-shared", "variation-A")

	mockCache := &mockEvaluationCountCache{
		failOnECKey: failECKey,
	}
	persister := &evaluationCountEventPersister{
		evaluationCountCacher: mockCache,
		logger:                zap.NewNop(),
	}

	fails := persister.incrementEnvEvents(envEvents)

	assert.Equal(t, 2, len(fails), "both events sharing the failed key pair must be retried")
	assert.Contains(t, fails, "event-A")
	assert.Contains(t, fails, "event-B")
	assert.True(t, fails["event-A"], "should be repeatable")
	assert.True(t, fails["event-B"], "should be repeatable")
}

func TestIncrementEnvEvents_MultiSourceMetricsAttribution(t *testing.T) {
	hour1 := int64(1709974800) // 2024-03-09 09:00:00 UTC

	origCounter := evaluationEventCounter
	evaluationEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "subscriber",
			Name:      "evaluation_event_total_multisource_test",
		}, []string{"environment_id", "source_id", "feature_id", "variation_id"})
	defer func() { evaluationEventCounter = origCounter }()

	envEvents := environmentEventMap{
		"env-1": eventMap{
			"event-android": &eventproto.EvaluationEvent{
				FeatureId:      "feature-X",
				VariationId:    "variation-A",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-1",
				User:           &userproto.User{Id: "user-1"},
				SourceId:       eventproto.SourceId_ANDROID,
				Timestamp:      hour1,
				FeatureVersion: 1,
			},
			"event-ios": &eventproto.EvaluationEvent{
				FeatureId:      "feature-X",
				VariationId:    "variation-A",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-2",
				User:           &userproto.User{Id: "user-2"},
				SourceId:       eventproto.SourceId_IOS,
				Timestamp:      hour1,
				FeatureVersion: 1,
			},
		},
	}

	tTime := time.Unix(hour1, 0)
	date := time.Date(tTime.Year(), tTime.Month(), tTime.Day(), tTime.Hour(), 0, 0, 0, time.UTC)
	failECKey := fmt.Sprintf("env-1:%s:%d:%s:%s", eventCountKey, date.Unix(), "feature-X", "variation-A")

	mockCache := &mockEvaluationCountCache{
		failOnECKey: failECKey,
	}
	persister := &evaluationCountEventPersister{
		evaluationCountCacher: mockCache,
		logger:                zap.NewNop(),
	}

	fails := persister.incrementEnvEvents(envEvents)

	assert.Contains(t, fails, "event-android")
	assert.Contains(t, fails, "event-ios")
	assert.True(t, fails["event-android"])
	assert.True(t, fails["event-ios"])

	// Both sources must show zero: succeeded = total(1) - failed(1) = 0 for each.
	// Before the fix, ecKeyToMetricsKey only kept the last-written sourceId,
	// so one source would incorrectly report succeeded > 0.
	androidCounter, err := evaluationEventCounter.GetMetricWithLabelValues(
		"env-1", eventproto.SourceId_ANDROID.String(), "feature-X", "variation-A")
	require.NoError(t, err)
	iosCounter, err := evaluationEventCounter.GetMetricWithLabelValues(
		"env-1", eventproto.SourceId_IOS.String(), "feature-X", "variation-A")
	require.NoError(t, err)

	var androidMetric dto.Metric
	require.NoError(t, androidCounter.Write(&androidMetric))
	var iosMetric dto.Metric
	require.NoError(t, iosCounter.Write(&iosMetric))

	assert.Equal(t, float64(0), androidMetric.GetCounter().GetValue(),
		"ANDROID counter must be 0: all events for this source failed")
	assert.Equal(t, float64(0), iosMetric.GetCounter().GetValue(),
		"IOS counter must be 0: all events for this source failed")
}

func TestFlushAggregatedCounts_AllPairsFail(t *testing.T) {
	t.Parallel()

	eventCounts := map[string]int64{
		"env1:ec:hour1:feature1:varA": 10,
		"env1:ec:hour1:feature2:varB": 20,
	}
	userCounts := map[string]map[string]struct{}{
		"env1:uc:hour1:feature1:varA": {"user1": {}},
		"env1:uc:hour1:feature2:varB": {"user2": {}},
	}

	mockCache := &mockEvaluationCountCache{
		shouldFailIncrBy: true,
	}
	persister := &evaluationCountEventPersister{
		evaluationCountCacher: mockCache,
		logger:                zap.NewNop(),
	}

	failedKeys, err := persister.flushAggregatedCounts(eventCounts, userCounts)

	assert.Error(t, err)
	assert.Equal(t, 2, len(failedKeys), "all ecKeys must be in failedKeys when all pairs fail")
	assert.Contains(t, failedKeys, "env1:ec:hour1:feature1:varA")
	assert.Contains(t, failedKeys, "env1:ec:hour1:feature2:varB")

	// Both PFAdds should still have been attempted (no early return)
	assert.Equal(t, 2, mockCache.pfaddCallCount, "all PFAdds should be attempted even when INCRBY fails")
	// Both INCRBYs were attempted (both failed)
	assert.Equal(t, 2, mockCache.execCount, "all INCRBYs should be attempted")
}

func TestFlushAggregatedCounts_PartialFailure(t *testing.T) {
	t.Parallel()

	// Two key pairs: one will succeed, one will fail via failOnECKey
	eventCounts := map[string]int64{
		"env1:ec:hour1:feature-ok:varA":  50,
		"env1:ec:hour1:feature-bad:varB": 30,
	}
	userCounts := map[string]map[string]struct{}{
		"env1:uc:hour1:feature-ok:varA":  {"user1": {}},
		"env1:uc:hour1:feature-bad:varB": {"user2": {}},
	}

	mockCache := &mockEvaluationCountCache{
		failOnECKey: "env1:ec:hour1:feature-bad:varB",
	}
	persister := &evaluationCountEventPersister{
		evaluationCountCacher: mockCache,
		logger:                zap.NewNop(),
	}

	failedKeys, err := persister.flushAggregatedCounts(eventCounts, userCounts)

	// Should return error for partial failure
	assert.Error(t, err)
	assert.Contains(t, failedKeys, "env1:ec:hour1:feature-bad:varB",
		"failed ecKey must be in failedKeys")
	assert.NotContains(t, failedKeys, "env1:ec:hour1:feature-ok:varA",
		"succeeded ecKey must NOT be in failedKeys")

	// Both PFAdds should have been attempted (no early return)
	assert.Equal(t, 2, mockCache.pfaddCallCount, "all PFAdds should be attempted")

	// The succeeded key pair's data should be persisted
	assert.Equal(t, int64(50), mockCache.lastEventCounts["env1:ec:hour1:feature-ok:varA"],
		"succeeded key pair's count should be persisted")
	assert.NotContains(t, mockCache.lastEventCounts, "env1:ec:hour1:feature-bad:varB",
		"failed key pair's count should NOT be persisted")
}

// pfaddCall records a PFAdd call for verification
type pfaddCall struct {
	key string
	els []string
	err error
}

// mockEvaluationCountCache mocks the cache interface for testing
type mockEvaluationCountCache struct {
	mu               sync.Mutex
	flushCalled      bool
	lastEventCounts  map[string]int64
	lastUserCounts   map[string]map[string]struct{}
	shouldFailFlush  bool
	shouldFailPFAdd  bool   // Simulates PFADD failure for all keys
	shouldFailIncrBy bool   // Simulates INCRBY failure for all keys
	failOnECKey      string // If set, IncrementBy fails only for this specific key
	pipelineExecuted bool
	pipelineCount    int // Number of pipelines created
	execCount        int // Number of IncrementBy calls (direct INCRBY operations)
	pfaddCallCount   int // Number of direct PFAdd calls
	pfaddCalls       []pfaddCall
	incrbyCallOrder  []string // Track order of operations
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

func (m *mockEvaluationCountCache) IncrementBy(key string, value int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.flushCalled = true
	m.execCount++
	m.incrbyCallOrder = append(m.incrbyCallOrder, "INCRBY:"+key)

	if m.shouldFailIncrBy || m.shouldFailFlush {
		return 0, assert.AnError
	}
	if m.failOnECKey != "" && m.failOnECKey == key {
		return 0, assert.AnError
	}

	if m.lastEventCounts == nil {
		m.lastEventCounts = make(map[string]int64)
	}

	m.lastEventCounts[key] += value

	return m.lastEventCounts[key], nil
}

func (m *mockEvaluationCountCache) PFAdd(key string, els ...string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.flushCalled = true // Mark that flush operations are happening
	m.pfaddCallCount++
	m.incrbyCallOrder = append(m.incrbyCallOrder, "PFADD:"+key)

	if m.shouldFailPFAdd || m.shouldFailFlush {
		err := assert.AnError
		m.pfaddCalls = append(m.pfaddCalls, pfaddCall{key: key, els: els, err: err})
		return 0, err
	}

	// Record successful call
	m.pfaddCalls = append(m.pfaddCalls, pfaddCall{key: key, els: els, err: nil})

	// Initialize user counts map
	if m.lastUserCounts == nil {
		m.lastUserCounts = make(map[string]map[string]struct{})
	}
	if m.lastUserCounts[key] == nil {
		m.lastUserCounts[key] = make(map[string]struct{})
	}

	// Add users (simulating HyperLogLog idempotency)
	for _, user := range els {
		m.lastUserCounts[key][user] = struct{}{}
	}

	return int64(len(els)), nil
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

	// Record order in parent cache
	m.cache.mu.Lock()
	m.cache.incrbyCallOrder = append(m.cache.incrbyCallOrder, "INCRBY:"+key)
	m.cache.mu.Unlock()

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

	// Check if this pipeline should fail
	if m.shouldFailExec {
		return nil, assert.AnError
	}

	// Check if IncrBy should fail (only if this pipeline contains INCRBY)
	hasIncrBy := false
	for _, cmd := range m.commands {
		if cmd == "INCRBY" {
			hasIncrBy = true
			break
		}
	}
	if hasIncrBy && m.cache.shouldFailIncrBy {
		return nil, assert.AnError
	}

	// For successful pipelines, accumulate the data
	if m.cache.lastEventCounts == nil {
		m.cache.lastEventCounts = make(map[string]int64)
	}

	// Merge event counts (only for INCRBY commands)
	for key, count := range m.eventCounts {
		m.cache.lastEventCounts[key] += count
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
