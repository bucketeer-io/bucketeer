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

package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestEvaluateArchivabilityWithCounts(t *testing.T) {
	t.Parallel()

	now := time.Now().Unix()
	daysAgo := func(days int) int64 {
		return now - int64(days*24*60*60)
	}

	tests := []struct {
		name            string
		features        []*ftproto.Feature
		criteria        ArchivabilityCriteria
		codeRefCounts   map[string]int64
		expectedResults map[string]struct {
			isArchivable    bool
			blockingReasons []string
		}
	}{
		{
			name: "feature is archivable when meets all criteria",
			features: []*ftproto.Feature{
				{
					Id:       "feature-1",
					Archived: false,
					LastUsedInfo: &ftproto.FeatureLastUsedInfo{
						LastUsedAt: daysAgo(100),
					},
				},
			},
			criteria: ArchivabilityCriteria{
				UnusedDaysThreshold: 90,
				CheckCodeReferences: true,
			},
			codeRefCounts: map[string]int64{
				"feature-1": 0,
			},
			expectedResults: map[string]struct {
				isArchivable    bool
				blockingReasons []string
			}{
				"feature-1": {
					isArchivable:    true,
					blockingReasons: []string{},
				},
			},
		},
		{
			name: "feature is not archivable when already archived",
			features: []*ftproto.Feature{
				{
					Id:       "feature-1",
					Archived: true,
					LastUsedInfo: &ftproto.FeatureLastUsedInfo{
						LastUsedAt: daysAgo(100),
					},
				},
			},
			criteria: ArchivabilityCriteria{
				UnusedDaysThreshold: 90,
				CheckCodeReferences: false,
			},
			codeRefCounts: map[string]int64{},
			expectedResults: map[string]struct {
				isArchivable    bool
				blockingReasons []string
			}{
				"feature-1": {
					isArchivable:    false,
					blockingReasons: []string{BlockingReasonAlreadyArchived},
				},
			},
		},
		{
			name: "feature is not archivable when not unused long enough",
			features: []*ftproto.Feature{
				{
					Id:       "feature-1",
					Archived: false,
					LastUsedInfo: &ftproto.FeatureLastUsedInfo{
						LastUsedAt: daysAgo(30), // Only 30 days ago
					},
				},
			},
			criteria: ArchivabilityCriteria{
				UnusedDaysThreshold: 90,
				CheckCodeReferences: false,
			},
			codeRefCounts: map[string]int64{},
			expectedResults: map[string]struct {
				isArchivable    bool
				blockingReasons []string
			}{
				"feature-1": {
					isArchivable:    false,
					blockingReasons: []string{BlockingReasonNotUnusedLongEnough},
				},
			},
		},
		{
			name: "feature is not archivable when has code references",
			features: []*ftproto.Feature{
				{
					Id:       "feature-1",
					Archived: false,
					LastUsedInfo: &ftproto.FeatureLastUsedInfo{
						LastUsedAt: daysAgo(100),
					},
				},
			},
			criteria: ArchivabilityCriteria{
				UnusedDaysThreshold: 90,
				CheckCodeReferences: true,
			},
			codeRefCounts: map[string]int64{
				"feature-1": 5, // Has 5 code references
			},
			expectedResults: map[string]struct {
				isArchivable    bool
				blockingReasons []string
			}{
				"feature-1": {
					isArchivable:    false,
					blockingReasons: []string{BlockingReasonHasCodeReferences},
				},
			},
		},
		{
			name: "feature is archivable when has code references but check is disabled",
			features: []*ftproto.Feature{
				{
					Id:       "feature-1",
					Archived: false,
					LastUsedInfo: &ftproto.FeatureLastUsedInfo{
						LastUsedAt: daysAgo(100),
					},
				},
			},
			criteria: ArchivabilityCriteria{
				UnusedDaysThreshold: 90,
				CheckCodeReferences: false, // Check disabled
			},
			codeRefCounts: map[string]int64{
				"feature-1": 5, // Has 5 code references but check is disabled
			},
			expectedResults: map[string]struct {
				isArchivable    bool
				blockingReasons []string
			}{
				"feature-1": {
					isArchivable:    true,
					blockingReasons: []string{},
				},
			},
		},
		{
			name: "feature is not archivable when never used (nil LastUsedInfo)",
			features: []*ftproto.Feature{
				{
					Id:           "feature-1",
					Archived:     false,
					LastUsedInfo: nil,
				},
			},
			criteria: ArchivabilityCriteria{
				UnusedDaysThreshold: 90,
				CheckCodeReferences: false,
			},
			codeRefCounts: map[string]int64{},
			expectedResults: map[string]struct {
				isArchivable    bool
				blockingReasons []string
			}{
				"feature-1": {
					isArchivable:    false,
					blockingReasons: []string{BlockingReasonNeverUsed},
				},
			},
		},
		{
			name: "feature is not archivable when has prerequisites dependencies",
			features: []*ftproto.Feature{
				{
					Id:       "feature-1",
					Archived: false,
					LastUsedInfo: &ftproto.FeatureLastUsedInfo{
						LastUsedAt: daysAgo(100),
					},
				},
				{
					Id:       "feature-2",
					Archived: false,
					LastUsedInfo: &ftproto.FeatureLastUsedInfo{
						LastUsedAt: daysAgo(100),
					},
					Prerequisites: []*ftproto.Prerequisite{
						{
							FeatureId: "feature-1", // feature-2 depends on feature-1
						},
					},
				},
			},
			criteria: ArchivabilityCriteria{
				UnusedDaysThreshold: 90,
				CheckCodeReferences: false,
			},
			codeRefCounts: map[string]int64{},
			expectedResults: map[string]struct {
				isArchivable    bool
				blockingReasons []string
			}{
				"feature-1": {
					isArchivable:    false,
					blockingReasons: []string{BlockingReasonHasDependencies},
				},
				"feature-2": {
					isArchivable:    true,
					blockingReasons: []string{},
				},
			},
		},
		{
			name: "feature is not archivable when has FEATURE_FLAG clause dependencies",
			features: []*ftproto.Feature{
				{
					Id:       "feature-1",
					Archived: false,
					LastUsedInfo: &ftproto.FeatureLastUsedInfo{
						LastUsedAt: daysAgo(100),
					},
				},
				{
					Id:       "feature-2",
					Archived: false,
					LastUsedInfo: &ftproto.FeatureLastUsedInfo{
						LastUsedAt: daysAgo(100),
					},
					Rules: []*ftproto.Rule{
						{
							Id: "rule-1",
							Clauses: []*ftproto.Clause{
								{
									Id:        "clause-1",
									Operator:  ftproto.Clause_FEATURE_FLAG,
									Attribute: "feature-1", // feature-2 depends on feature-1 via FEATURE_FLAG clause
								},
							},
						},
					},
				},
			},
			criteria: ArchivabilityCriteria{
				UnusedDaysThreshold: 90,
				CheckCodeReferences: false,
			},
			codeRefCounts: map[string]int64{},
			expectedResults: map[string]struct {
				isArchivable    bool
				blockingReasons []string
			}{
				"feature-1": {
					isArchivable:    false,
					blockingReasons: []string{BlockingReasonHasDependencies},
				},
				"feature-2": {
					isArchivable:    true,
					blockingReasons: []string{},
				},
			},
		},
		{
			name: "multiple blocking reasons",
			features: []*ftproto.Feature{
				{
					Id:           "feature-1",
					Archived:     true,
					LastUsedInfo: nil,
				},
			},
			criteria: ArchivabilityCriteria{
				UnusedDaysThreshold: 90,
				CheckCodeReferences: true,
			},
			codeRefCounts: map[string]int64{
				"feature-1": 5,
			},
			expectedResults: map[string]struct {
				isArchivable    bool
				blockingReasons []string
			}{
				"feature-1": {
					isArchivable: false,
					blockingReasons: []string{
						BlockingReasonAlreadyArchived,
						BlockingReasonNeverUsed,
						BlockingReasonHasCodeReferences,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			results := EvaluateArchivabilityWithCounts(tt.features, tt.criteria, tt.codeRefCounts)

			assert.Equal(t, len(tt.expectedResults), len(results))

			for _, result := range results {
				expected, ok := tt.expectedResults[result.Feature.Id]
				assert.True(t, ok, "unexpected feature ID: %s", result.Feature.Id)
				assert.Equal(t, expected.isArchivable, result.IsArchivable,
					"feature %s: isArchivable mismatch", result.Feature.Id)
				assert.ElementsMatch(t, expected.blockingReasons, result.BlockingReasons,
					"feature %s: blockingReasons mismatch", result.Feature.Id)
			}
		})
	}
}

type stubCodeRefStorage struct {
	counts map[string]int64
	err    error
}

func (s *stubCodeRefStorage) GetCodeReferenceCountsByFeatureIDs(
	ctx context.Context,
	environmentID string,
	featureIDs []string,
) (map[string]int64, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.counts, nil
}

func TestArchivabilityEvaluator_EvaluateArchivability(t *testing.T) {
	t.Parallel()

	now := time.Now().Unix()
	features := []*ftproto.Feature{
		{
			Id:       "feature-1",
			Archived: false,
			LastUsedInfo: &ftproto.FeatureLastUsedInfo{
				LastUsedAt: now - 100*24*60*60,
			},
		},
	}

	t.Run("returns error when CheckCodeReferences is true but storage is nil", func(t *testing.T) {
		t.Parallel()
		e := NewArchivabilityEvaluator(nil)
		_, err := e.EvaluateArchivability(context.Background(), features, ArchivabilityCriteria{
			UnusedDaysThreshold: 90,
			CheckCodeReferences: true,
		}, "env-1")
		assert.Error(t, err)
	})

	t.Run("fetches code reference counts and evaluates", func(t *testing.T) {
		t.Parallel()
		e := NewArchivabilityEvaluator(&stubCodeRefStorage{
			counts: map[string]int64{"feature-1": 0},
		})
		results, err := e.EvaluateArchivability(context.Background(), features, ArchivabilityCriteria{
			UnusedDaysThreshold: 90,
			CheckCodeReferences: true,
		}, "env-1")
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.True(t, results[0].IsArchivable)
	})
}

func TestCalculateUnusedDays(t *testing.T) {
	t.Parallel()

	now := time.Now().Unix()

	tests := []struct {
		name         string
		lastUsedAt   int64
		expectedDays int32
	}{
		{
			name:         "0 days ago",
			lastUsedAt:   now,
			expectedDays: 0,
		},
		{
			name:         "1 day ago",
			lastUsedAt:   now - 24*60*60,
			expectedDays: 1,
		},
		{
			name:         "90 days ago",
			lastUsedAt:   now - 90*24*60*60,
			expectedDays: 90,
		},
		{
			name:         "365 days ago",
			lastUsedAt:   now - 365*24*60*60,
			expectedDays: 365,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := CalculateUnusedDays(tt.lastUsedAt)
			// Allow 1 day tolerance for test execution time
			assert.InDelta(t, tt.expectedDays, result, 1)
		})
	}
}

func TestFilterArchivableFeatures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		results     []*ArchivabilityResult
		expectedIDs []string
	}{
		{
			name:        "empty results",
			results:     []*ArchivabilityResult{},
			expectedIDs: []string{},
		},
		{
			name: "all archivable",
			results: []*ArchivabilityResult{
				{Feature: &ftproto.Feature{Id: "f1"}, IsArchivable: true},
				{Feature: &ftproto.Feature{Id: "f2"}, IsArchivable: true},
			},
			expectedIDs: []string{"f1", "f2"},
		},
		{
			name: "none archivable",
			results: []*ArchivabilityResult{
				{Feature: &ftproto.Feature{Id: "f1"}, IsArchivable: false},
				{Feature: &ftproto.Feature{Id: "f2"}, IsArchivable: false},
			},
			expectedIDs: []string{},
		},
		{
			name: "mixed archivability",
			results: []*ArchivabilityResult{
				{Feature: &ftproto.Feature{Id: "f1"}, IsArchivable: true},
				{Feature: &ftproto.Feature{Id: "f2"}, IsArchivable: false},
				{Feature: &ftproto.Feature{Id: "f3"}, IsArchivable: true},
			},
			expectedIDs: []string{"f1", "f3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			features := FilterArchivableFeatures(tt.results)
			ids := make([]string, len(features))
			for i, f := range features {
				ids[i] = f.Id
			}
			assert.ElementsMatch(t, tt.expectedIDs, ids)
		})
	}
}

func TestGetArchivableFeatureIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		results     []*ArchivabilityResult
		expectedIDs []string
	}{
		{
			name:        "empty results",
			results:     []*ArchivabilityResult{},
			expectedIDs: []string{},
		},
		{
			name: "mixed archivability",
			results: []*ArchivabilityResult{
				{Feature: &ftproto.Feature{Id: "f1"}, IsArchivable: true},
				{Feature: &ftproto.Feature{Id: "f2"}, IsArchivable: false},
				{Feature: &ftproto.Feature{Id: "f3"}, IsArchivable: true},
			},
			expectedIDs: []string{"f1", "f3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ids := GetArchivableFeatureIDs(tt.results)
			assert.ElementsMatch(t, tt.expectedIDs, ids)
		})
	}
}

func TestBulkArchiveResult(t *testing.T) {
	t.Parallel()

	t.Run("new bulk archive result", func(t *testing.T) {
		t.Parallel()
		result := NewBulkArchiveResult(10)
		assert.Equal(t, 0, result.SuccessCount)
		assert.Equal(t, 0, result.FailureCount)
		assert.Empty(t, result.Results)
	})

	t.Run("add success", func(t *testing.T) {
		t.Parallel()
		result := NewBulkArchiveResult(2)
		result.AddSuccess("feature-1")
		result.AddSuccess("feature-2")

		assert.Equal(t, 2, result.SuccessCount)
		assert.Equal(t, 0, result.FailureCount)
		assert.Len(t, result.Results, 2)
		assert.True(t, result.Results[0].Success)
		assert.True(t, result.Results[1].Success)
	})

	t.Run("add failure", func(t *testing.T) {
		t.Parallel()
		result := NewBulkArchiveResult(2)
		err := assert.AnError
		result.AddFailure("feature-1", err)

		assert.Equal(t, 0, result.SuccessCount)
		assert.Equal(t, 1, result.FailureCount)
		assert.Len(t, result.Results, 1)
		assert.False(t, result.Results[0].Success)
		assert.Equal(t, err, result.Results[0].Error)
	})

	t.Run("mixed success and failure", func(t *testing.T) {
		t.Parallel()
		result := NewBulkArchiveResult(3)
		result.AddSuccess("feature-1")
		result.AddFailure("feature-2", assert.AnError)
		result.AddSuccess("feature-3")

		assert.Equal(t, 2, result.SuccessCount)
		assert.Equal(t, 1, result.FailureCount)
		assert.Len(t, result.Results, 3)
	})
}

// stubFeatureArchiver is a test implementation of FeatureArchiver.
type stubFeatureArchiver struct {
	archiveFunc func(ctx context.Context, featureID, environmentID string) error
	calls       []struct {
		featureID     string
		environmentID string
	}
}

func (s *stubFeatureArchiver) ArchiveFeature(ctx context.Context, featureID, environmentID string) error {
	s.calls = append(s.calls, struct {
		featureID     string
		environmentID string
	}{featureID, environmentID})
	if s.archiveFunc != nil {
		return s.archiveFunc(ctx, featureID, environmentID)
	}
	return nil
}

func TestArchiveFeaturesInBulk(t *testing.T) {
	t.Parallel()

	t.Run("empty feature list", func(t *testing.T) {
		t.Parallel()
		archiver := &stubFeatureArchiver{}
		result := ArchiveFeaturesInBulk(context.Background(), []string{}, "env-1", archiver)

		assert.Equal(t, 0, result.SuccessCount)
		assert.Equal(t, 0, result.FailureCount)
		assert.Empty(t, result.Results)
		assert.Empty(t, archiver.calls)
	})

	t.Run("all features archived successfully", func(t *testing.T) {
		t.Parallel()
		archiver := &stubFeatureArchiver{}
		featureIDs := []string{"feature-1", "feature-2", "feature-3"}

		result := ArchiveFeaturesInBulk(context.Background(), featureIDs, "env-1", archiver)

		assert.Equal(t, 3, result.SuccessCount)
		assert.Equal(t, 0, result.FailureCount)
		assert.Len(t, result.Results, 3)
		for _, r := range result.Results {
			assert.True(t, r.Success)
			assert.Nil(t, r.Error)
		}
		// Verify all features were processed
		assert.Len(t, archiver.calls, 3)
		for _, call := range archiver.calls {
			assert.Equal(t, "env-1", call.environmentID)
		}
	})

	t.Run("all features fail to archive", func(t *testing.T) {
		t.Parallel()
		expectedErr := assert.AnError
		archiver := &stubFeatureArchiver{
			archiveFunc: func(ctx context.Context, featureID, environmentID string) error {
				return expectedErr
			},
		}
		featureIDs := []string{"feature-1", "feature-2"}

		result := ArchiveFeaturesInBulk(context.Background(), featureIDs, "env-1", archiver)

		assert.Equal(t, 0, result.SuccessCount)
		assert.Equal(t, 2, result.FailureCount)
		assert.Len(t, result.Results, 2)
		for _, r := range result.Results {
			assert.False(t, r.Success)
			assert.Equal(t, expectedErr, r.Error)
		}
	})

	t.Run("partial success", func(t *testing.T) {
		t.Parallel()
		archiver := &stubFeatureArchiver{
			archiveFunc: func(ctx context.Context, featureID, environmentID string) error {
				if featureID == "feature-2" {
					return assert.AnError
				}
				return nil
			},
		}
		featureIDs := []string{"feature-1", "feature-2", "feature-3"}

		result := ArchiveFeaturesInBulk(context.Background(), featureIDs, "env-1", archiver)

		assert.Equal(t, 2, result.SuccessCount)
		assert.Equal(t, 1, result.FailureCount)
		assert.Len(t, result.Results, 3)

		// Verify order is preserved
		assert.Equal(t, "feature-1", result.Results[0].FeatureID)
		assert.True(t, result.Results[0].Success)
		assert.Equal(t, "feature-2", result.Results[1].FeatureID)
		assert.False(t, result.Results[1].Success)
		assert.Equal(t, "feature-3", result.Results[2].FeatureID)
		assert.True(t, result.Results[2].Success)
	})

	t.Run("continues after failure", func(t *testing.T) {
		t.Parallel()
		callCount := 0
		archiver := &stubFeatureArchiver{
			archiveFunc: func(ctx context.Context, featureID, environmentID string) error {
				callCount++
				if featureID == "feature-1" {
					return assert.AnError
				}
				return nil
			},
		}
		featureIDs := []string{"feature-1", "feature-2", "feature-3"}

		result := ArchiveFeaturesInBulk(context.Background(), featureIDs, "env-1", archiver)

		// Verify all features were attempted despite first failure
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, result.SuccessCount)
		assert.Equal(t, 1, result.FailureCount)
	})
}
