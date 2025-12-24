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
	"errors"
	"time"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// Blocking reason constants for archivability evaluation.
// These reasons indicate why a feature flag cannot be archived.
const (
	// BlockingReasonAlreadyArchived indicates the feature is already archived.
	BlockingReasonAlreadyArchived = "already_archived"
	// BlockingReasonHasDependencies indicates other features depend on this feature.
	BlockingReasonHasDependencies = "has_dependencies"
	// BlockingReasonNotUnusedLongEnough indicates the feature was used recently.
	BlockingReasonNotUnusedLongEnough = "not_unused_long_enough"
	// BlockingReasonHasCodeReferences indicates the feature has code references.
	BlockingReasonHasCodeReferences = "has_code_references"
	// BlockingReasonNeverUsed indicates the feature has never been used (no LastUsedInfo).
	BlockingReasonNeverUsed = "never_used"
)

// ArchivabilityCriteria defines the criteria for evaluating feature archivability.
type ArchivabilityCriteria struct {
	// UnusedDaysThreshold is the minimum number of days a feature must be unused
	// before it can be considered for archiving.
	UnusedDaysThreshold int32
	// CheckCodeReferences determines whether to check for code references
	// when evaluating archivability.
	CheckCodeReferences bool
}

// ArchivabilityResult contains the result of archivability evaluation for a single feature.
type ArchivabilityResult struct {
	// Feature is the evaluated feature.
	Feature *ftproto.Feature
	// IsArchivable indicates whether the feature can be archived.
	IsArchivable bool
	// UnusedDays is the number of days since the feature was last used.
	UnusedDays int32
	// CodeRefCount is the number of code references for this feature.
	CodeRefCount int64
	// BlockingReasons contains the reasons why the feature cannot be archived.
	// Empty if IsArchivable is true.
	BlockingReasons []string
}

// BulkArchiveResult contains the result of a bulk archive operation.
type BulkArchiveResult struct {
	// SuccessCount is the number of features successfully archived.
	SuccessCount int
	// FailureCount is the number of features that failed to archive.
	FailureCount int
	// Results contains individual results for each feature.
	Results []*SingleArchiveResult
}

// SingleArchiveResult contains the result of archiving a single feature.
type SingleArchiveResult struct {
	// FeatureID is the ID of the feature.
	FeatureID string
	// Success indicates whether the archive operation succeeded.
	Success bool
	// Error contains the error if the operation failed.
	Error error
}

// CodeRefStorage defines the interface for accessing code reference counts.
// This interface is implemented by the coderef storage layer.
type CodeRefStorage interface {
	// GetCodeReferenceCountsByFeatureIDs returns a map of feature ID to code reference count
	// for the given environment.
	GetCodeReferenceCountsByFeatureIDs(
		ctx context.Context,
		environmentID string,
		featureIDs []string,
	) (map[string]int64, error)
}

// ArchivabilityEvaluator provides centralized logic for evaluating feature archivability.
// It can be used by both batch jobs and API endpoints.
type ArchivabilityEvaluator struct {
	codeRefStorage CodeRefStorage
}

var errCodeRefStorageRequired = errors.New("feature: code reference storage required when CheckCodeReferences is true")

// NewArchivabilityEvaluator creates a new ArchivabilityEvaluator instance.
func NewArchivabilityEvaluator(codeRefStorage CodeRefStorage) *ArchivabilityEvaluator {
	return &ArchivabilityEvaluator{
		codeRefStorage: codeRefStorage,
	}
}

// EvaluateArchivability evaluates the archivability of multiple features
// based on the given criteria. It fetches code reference counts internally
// to avoid N+1 queries.
func (e *ArchivabilityEvaluator) EvaluateArchivability(
	ctx context.Context,
	features []*ftproto.Feature,
	criteria ArchivabilityCriteria,
	environmentID string,
) ([]*ArchivabilityResult, error) {
	if criteria.CheckCodeReferences && e.codeRefStorage == nil {
		return nil, errCodeRefStorageRequired
	}

	// Build feature IDs list for bulk code reference fetch
	featureIDs := make([]string, len(features))
	for i, f := range features {
		featureIDs[i] = f.Id
	}

	// Bulk fetch code reference counts to avoid N+1 queries
	var codeRefCounts map[string]int64
	var err error
	if criteria.CheckCodeReferences {
		codeRefCounts, err = e.codeRefStorage.GetCodeReferenceCountsByFeatureIDs(ctx, environmentID, featureIDs)
		if err != nil {
			return nil, err
		}
	} else {
		codeRefCounts = make(map[string]int64)
	}

	// Use the standalone function for actual evaluation
	return EvaluateArchivabilityWithCounts(features, criteria, codeRefCounts), nil
}

// EvaluateArchivabilityWithCounts evaluates the archivability of multiple features
// using pre-fetched code reference counts. This is useful for testing or when
// code reference counts are already available.
func EvaluateArchivabilityWithCounts(
	features []*ftproto.Feature,
	criteria ArchivabilityCriteria,
	codeRefCounts map[string]int64,
) []*ArchivabilityResult {
	results := make([]*ArchivabilityResult, 0, len(features))

	// Build dependency map: which features depend on which other features
	dependencyTargets := buildDependencyTargets(features)

	for _, f := range features {
		result := evaluateSingleFeature(f, criteria, codeRefCounts, dependencyTargets)
		results = append(results, result)
	}

	return results
}

// buildDependencyTargets builds a map of feature IDs that are depended upon by other features.
// Returns a map where keys are feature IDs that have dependents.
func buildDependencyTargets(features []*ftproto.Feature) map[string]bool {
	targets := make(map[string]bool)

	for _, f := range features {
		// Check prerequisites
		for _, prereq := range f.Prerequisites {
			targets[prereq.FeatureId] = true
		}

		// Check rules for FEATURE_FLAG clause type dependencies
		for _, rule := range f.Rules {
			for _, clause := range rule.Clauses {
				if clause.Operator == ftproto.Clause_FEATURE_FLAG {
					// The attribute contains the feature ID being referenced
					targets[clause.Attribute] = true
				}
			}
		}
	}

	return targets
}

// evaluateSingleFeature evaluates the archivability of a single feature.
func evaluateSingleFeature(
	f *ftproto.Feature,
	criteria ArchivabilityCriteria,
	codeRefCounts map[string]int64,
	dependencyTargets map[string]bool,
) *ArchivabilityResult {
	result := &ArchivabilityResult{
		Feature:         f,
		IsArchivable:    true,
		BlockingReasons: []string{},
	}

	// Check if already archived
	if f.Archived {
		result.IsArchivable = false
		result.BlockingReasons = append(result.BlockingReasons, BlockingReasonAlreadyArchived)
	}

	// Check if feature has never been used
	if f.LastUsedInfo == nil || f.LastUsedInfo.LastUsedAt == 0 {
		result.IsArchivable = false
		result.BlockingReasons = append(result.BlockingReasons, BlockingReasonNeverUsed)
	} else {
		// Calculate unused days
		result.UnusedDays = CalculateUnusedDays(f.LastUsedInfo.LastUsedAt)

		// Check if not unused long enough
		if result.UnusedDays < criteria.UnusedDaysThreshold {
			result.IsArchivable = false
			result.BlockingReasons = append(result.BlockingReasons, BlockingReasonNotUnusedLongEnough)
		}
	}

	// Check code references if enabled
	if criteria.CheckCodeReferences {
		codeRefCount, exists := codeRefCounts[f.Id]
		if exists {
			result.CodeRefCount = codeRefCount
		}
		if codeRefCount > 0 {
			result.IsArchivable = false
			result.BlockingReasons = append(result.BlockingReasons, BlockingReasonHasCodeReferences)
		}
	}

	// Check if other features depend on this feature
	if dependencyTargets[f.Id] {
		result.IsArchivable = false
		result.BlockingReasons = append(result.BlockingReasons, BlockingReasonHasDependencies)
	}

	return result
}

// CalculateUnusedDays calculates the number of days since the given timestamp.
func CalculateUnusedDays(lastUsedAt int64) int32 {
	now := time.Now().Unix()
	diff := now - lastUsedAt
	if diff < 0 {
		return 0
	}
	days := diff / (24 * 60 * 60)
	return int32(days)
}

// FilterArchivableFeatures returns only the features that are archivable from the evaluation results.
// This is a convenience function for batch jobs and APIs that need to process archivable features.
func FilterArchivableFeatures(results []*ArchivabilityResult) []*ftproto.Feature {
	archivable := make([]*ftproto.Feature, 0)
	for _, r := range results {
		if r.IsArchivable {
			archivable = append(archivable, r.Feature)
		}
	}
	return archivable
}

// GetArchivableFeatureIDs returns the IDs of features that are archivable from the evaluation results.
func GetArchivableFeatureIDs(results []*ArchivabilityResult) []string {
	ids := make([]string, 0)
	for _, r := range results {
		if r.IsArchivable {
			ids = append(ids, r.Feature.Id)
		}
	}
	return ids
}

// NewBulkArchiveResult creates a new BulkArchiveResult with the given capacity.
func NewBulkArchiveResult(capacity int) *BulkArchiveResult {
	return &BulkArchiveResult{
		Results: make([]*SingleArchiveResult, 0, capacity),
	}
}

// AddSuccess adds a successful archive result.
func (r *BulkArchiveResult) AddSuccess(featureID string) {
	r.Results = append(r.Results, &SingleArchiveResult{
		FeatureID: featureID,
		Success:   true,
	})
	r.SuccessCount++
}

// AddFailure adds a failed archive result.
func (r *BulkArchiveResult) AddFailure(featureID string, err error) {
	r.Results = append(r.Results, &SingleArchiveResult{
		FeatureID: featureID,
		Success:   false,
		Error:     err,
	})
	r.FailureCount++
}

// FeatureArchiver defines the interface for archiving a single feature.
// This interface is implemented by higher-level layers (API, batch) that have
// access to infrastructure dependencies (mysql, publisher, command handler).
// The implementation is responsible for:
// - Running the archive operation in a transaction
// - Creating and publishing domain events
// - Incrementing the feature version
type FeatureArchiver interface {
	// ArchiveFeature archives a single feature by ID.
	// Returns an error if the archive operation fails.
	ArchiveFeature(ctx context.Context, featureID, environmentID string) error
}

// ArchiveFeaturesInBulk archives multiple features using the provided archiver.
// Each feature is archived independently to support partial success.
// Returns BulkArchiveResult containing success/failure counts and individual results.
func ArchiveFeaturesInBulk(
	ctx context.Context,
	featureIDs []string,
	environmentID string,
	archiver FeatureArchiver,
) *BulkArchiveResult {
	result := NewBulkArchiveResult(len(featureIDs))

	for _, featureID := range featureIDs {
		if err := archiver.ArchiveFeature(ctx, featureID, environmentID); err != nil {
			result.AddFailure(featureID, err)
		} else {
			result.AddSuccess(featureID)
		}
	}

	return result
}
