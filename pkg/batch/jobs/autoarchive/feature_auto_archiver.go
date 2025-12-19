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

package autoarchive

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	coderefstorage "github.com/bucketeer-io/bucketeer/v2/pkg/coderef/storage"
	environmentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	environmentstorage "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	featurestorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// featureAutoArchiver is a batch job that automatically archives
// feature flags based on environment-level auto-archive settings.
type featureAutoArchiver struct {
	envStorage     environmentstorage.EnvironmentStorage
	ftStorage      featurestorage.FeatureStorage
	codeRefStorage coderefstorage.CodeReferenceStorage
	evaluator      *domain.ArchivabilityEvaluator
	featureClient  featureclient.Client
	opts           *jobs.Options
	logger         *zap.Logger
}

// NewFeatureAutoArchiver creates a new feature auto-archiver batch job.
func NewFeatureAutoArchiver(
	mysqlClient mysql.Client,
	featureClient featureclient.Client,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 10 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	codeRefStorage := coderefstorage.NewCodeReferenceStorage(mysqlClient)
	return &featureAutoArchiver{
		envStorage:     environmentstorage.NewEnvironmentStorage(mysqlClient),
		ftStorage:      featurestorage.NewFeatureStorage(mysqlClient),
		codeRefStorage: codeRefStorage,
		evaluator:      domain.NewArchivabilityEvaluator(codeRefStorage),
		featureClient:  featureClient,
		opts:           dopts,
		logger:         dopts.Logger.Named("feature-auto-archiver"),
	}
}

// Run executes the feature auto-archive batch job.
// It iterates through all environments with auto-archive enabled,
// evaluates feature archivability, and archives eligible features.
func (a *featureAutoArchiver) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobFeatureAutoArchiver, lastErr, time.Since(startTime))
	}()

	ctx, cancel := context.WithTimeout(ctx, a.opts.Timeout)
	defer cancel()

	a.logger.Info("FeatureAutoArchiver start running")

	// Get environments with auto-archive enabled
	environments, err := a.envStorage.ListAutoArchiveEnabledEnvironments(ctx)
	if err != nil {
		a.logger.Error("Failed to list auto-archive enabled environments", zap.Error(err))
		return err
	}

	if len(environments) == 0 {
		a.logger.Info("No environments with auto-archive enabled")
		return nil
	}

	a.logger.Info("Found auto-archive enabled environments", zap.Int("count", len(environments)))

	// Process each environment
	for _, env := range environments {
		if err := a.processEnvironment(ctx, env); err != nil {
			a.logger.Error("Failed to process environment",
				zap.String("environmentId", env.Id),
				zap.Error(err),
			)
			lastErr = err
			// Continue processing other environments even if one fails
			continue
		}
	}

	a.logger.Info("FeatureAutoArchiver finished",
		zap.Duration("elapsedTime", time.Since(startTime)),
	)
	return lastErr
}

// processEnvironment processes a single environment for auto-archiving.
func (a *featureAutoArchiver) processEnvironment(
	ctx context.Context,
	env *environmentdomain.EnvironmentV2,
) error {
	a.logger.Info("Processing environment",
		zap.String("environmentId", env.Id),
		zap.Int32("unusedDaysThreshold", env.AutoArchiveUnusedDays),
		zap.Bool("checkCodeRefs", env.AutoArchiveCheckCodeRefs),
	)

	// Get non-archived features in the environment
	features, err := a.listNonArchivedFeatures(ctx, env.Id)
	if err != nil {
		return err
	}

	if len(features) == 0 {
		a.logger.Info("No non-archived features found", zap.String("environmentId", env.Id))
		return nil
	}

	// Build archivability criteria from environment settings
	criteria := domain.ArchivabilityCriteria{
		UnusedDaysThreshold: env.AutoArchiveUnusedDays,
		CheckCodeReferences: env.AutoArchiveCheckCodeRefs,
	}

	// Evaluate archivability using pre-created evaluator
	results, err := a.evaluator.EvaluateArchivability(ctx, features, criteria, env.Id)
	if err != nil {
		a.logger.Error("Failed to evaluate archivability",
			zap.String("environmentId", env.Id),
			zap.Error(err),
		)
		return err
	}

	// Get archivable feature IDs
	archivableIDs := domain.GetArchivableFeatureIDs(results)
	if len(archivableIDs) == 0 {
		a.logger.Info("No archivable features found",
			zap.String("environmentId", env.Id),
			zap.Int("evaluatedCount", len(results)),
		)
		return nil
	}

	a.logger.Info("Found archivable features",
		zap.String("environmentId", env.Id),
		zap.Int("archivableCount", len(archivableIDs)),
		zap.Int("evaluatedCount", len(results)),
	)

	// Archive features in bulk
	bulkResult := domain.ArchiveFeaturesInBulk(ctx, archivableIDs, env.Id, a)

	a.logger.Info("Archive results",
		zap.String("environmentId", env.Id),
		zap.Int("successCount", bulkResult.SuccessCount),
		zap.Int("failureCount", bulkResult.FailureCount),
	)

	// Return error if any archive failed
	if bulkResult.FailureCount > 0 {
		for _, result := range bulkResult.Results {
			if !result.Success {
				return result.Error
			}
		}
	}

	return nil
}

// listNonArchivedFeatures retrieves all non-archived features for the given environment.
func (a *featureAutoArchiver) listNonArchivedFeatures(
	ctx context.Context,
	environmentID string,
) ([]*featureproto.Feature, error) {
	archived := false
	options := &mysql.ListOptions{
		Filters: []*mysql.FilterV2{
			{
				Column:   "feature.archived",
				Operator: mysql.OperatorEqual,
				Value:    archived,
			},
			{
				Column:   "feature.environment_id",
				Operator: mysql.OperatorEqual,
				Value:    environmentID,
			},
		},
		Limit:  mysql.QueryNoLimit,
		Offset: mysql.QueryNoOffset,
	}

	features, _, _, err := a.ftStorage.ListFeatures(ctx, options)
	if err != nil {
		a.logger.Error("Failed to list features",
			zap.String("environmentId", environmentID),
			zap.Error(err),
		)
		return nil, err
	}

	return features, nil
}

// ArchiveFeature implements domain.FeatureArchiver interface.
// It archives a single feature by calling the feature service API.
func (a *featureAutoArchiver) ArchiveFeature(
	ctx context.Context,
	featureID, environmentID string,
) error {
	_, err := a.featureClient.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
		Id:            featureID,
		EnvironmentId: environmentID,
		Archived:      wrapperspb.Bool(true),
		Comment:       "Automatically archived by auto-archive batch job",
	})
	if err != nil {
		a.logger.Error("Failed to archive feature",
			zap.String("featureId", featureID),
			zap.String("environmentId", environmentID),
			zap.Error(err),
		)
		return err
	}

	a.logger.Info("Successfully archived feature",
		zap.String("featureId", featureID),
		zap.String("environmentId", environmentID),
	)
	return nil
}
