package processor

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

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
