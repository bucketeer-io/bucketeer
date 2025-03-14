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
//

package calculator

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/domain"
	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/experimentcalc"
	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/stan"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/proto/environment"
	"github.com/bucketeer-io/bucketeer/proto/experiment"
)

const (
	day = 24 * time.Hour
)

type experimentCalculate struct {
	environmentClient environmentclient.Client
	experimentClient  experimentclient.Client
	calculator        *experimentcalc.ExperimentCalculator
	experimentLock    *ExperimentLock
	location          *time.Location
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewExperimentCalculate(
	httpStan *stan.Stan,
	stanModelID string,
	environmentClient environmentclient.Client,
	experimentClient experimentclient.Client,
	ecClient ecclient.Client,
	mysqlClient mysql.Client,
	experimentLock *ExperimentLock,
	location *time.Location,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 1 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	calculator := experimentcalc.NewExperimentCalculator(
		httpStan,
		stanModelID,
		environmentClient,
		ecClient,
		experimentClient,
		mysqlClient,
		dopts.Metrics,
		location,
		dopts.Logger,
	)
	return &experimentCalculate{
		environmentClient: environmentClient,
		experimentClient:  experimentClient,
		calculator:        calculator,
		experimentLock:    experimentLock,
		location:          location,
		opts:              dopts,
		logger:            dopts.Logger.Named("experiment-calculate"),
	}
}

func (e *experimentCalculate) Run(ctx context.Context) error {
	// Create a context with timeout from the options
	ctxWithTimeout, cancel := context.WithTimeout(ctx, e.opts.Timeout)

	// Run the calculation logic in a goroutine and return immediately
	go func() {
		defer cancel() // Ensure context is canceled when goroutine completes

		e.logger.Info("Started experiment calculation job")
		startTime := time.Now().In(e.location)
		environments, environmentErr := e.listEnvironments(ctxWithTimeout)
		if environmentErr != nil {
			e.logger.Error("Failed to list environments when calculating experiments",
				log.FieldsFromImcomingContext(ctxWithTimeout).AddFields(
					zap.Error(environmentErr),
				)...,
			)
			return
		}
		var calculatedCount int
		for _, env := range environments {
			experiments, experimentErr := e.listExperiments(ctxWithTimeout, env.Id)
			if experimentErr != nil {
				e.logger.Error("Failed to list experiments when running experiment calculation",
					log.FieldsFromImcomingContext(ctxWithTimeout).AddFields(
						zap.Error(experimentErr),
					)...,
				)
				return
			}
			if experiments == nil {
				e.logger.Info("There are no experiments for calculation in the specified environment",
					log.FieldsFromImcomingContext(ctxWithTimeout).AddFields(
						zap.String("environmentId", env.Id),
					)...,
				)
				continue
			}
			for _, ex := range experiments {
				calculateErr := e.calculateExperimentWithLock(ctxWithTimeout, env, ex)
				if calculateErr != nil {
					e.logger.Error("Failed to calculate experiment",
						log.FieldsFromImcomingContext(ctxWithTimeout).AddFields(
							zap.Error(calculateErr),
							zap.String("environmentId", env.Id),
							zap.Any("experiment", ex),
						)...,
					)
					continue
				}
				calculatedCount++
				e.logger.Info("Experiment calculated successfully in the specified environment",
					log.FieldsFromImcomingContext(ctxWithTimeout).AddFields(
						zap.String("environmentId", env.Id),
						zap.Any("experiment", ex),
					)...,
				)
			}
		}
		e.logger.Info("Finished experiment calculation job",
			zap.Duration("elapsedTime", time.Since(startTime)),
			zap.Int("totalCalculatedExperiments", calculatedCount),
		)
	}()

	// Return immediately, not waiting for the goroutine to complete
	return nil
}

func (e *experimentCalculate) calculateExperimentWithLock(ctx context.Context,
	env *environment.EnvironmentV2,
	experiment *experiment.Experiment,
) error {
	locked, lockValue, err := e.experimentLock.Lock(ctx, env.Id, experiment.Id)
	if err != nil {
		return fmt.Errorf("Failed to acquire lock when calculating experiment. Error: %w", err)
	}
	if !locked {
		e.logger.Info("Experiment is being calculated by another instance",
			zap.String("environmentId", env.Id),
			zap.String("experimentId", experiment.Id),
			zap.String("experimentName", experiment.Name),
		)
		return nil
	}
	if calcErr := e.calculateExperiment(ctx, env, experiment); calcErr != nil {
		// To prevent calculating the same experiment multiple times in a short time,
		// we set the TTL for the lock key and only unlock it when an error occurs so that it can retry.
		unlocked, unlockErr := e.experimentLock.Unlock(ctx, env.Id, experiment.Id, lockValue)
		if unlockErr != nil {
			e.logger.Error("Failed to release lock when calculating experiment",
				zap.Error(unlockErr),
				zap.String("environmentId", env.Id),
				zap.Any("experiment", experiment),
			)
		}
		if !unlocked {
			e.logger.Warn("Lock was not released when calculating experiment, possibly expired",
				zap.String("environmentId", env.Id),
				zap.Any("experiment", experiment),
			)
		}
		return calcErr
	}
	return nil
}

func (e *experimentCalculate) calculateExperiment(ctx context.Context,
	env *environment.EnvironmentV2,
	experiment *experiment.Experiment,
) error {
	err := e.calculator.Run(ctx, &domain.ExperimentCalculatorReq{
		EnvironmentId: env.Id,
		Experiment:    experiment,
	})
	if err != nil {
		e.logger.Error("Failed experiment calculation",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", env.Id),
				zap.Any("experiment", experiment),
			)...,
		)
		return err
	}
	return nil
}

func (e *experimentCalculate) listEnvironments(
	ctx context.Context,
) ([]*environment.EnvironmentV2, error) {
	listEnvironmentsRequest := environment.ListEnvironmentsV2Request{
		PageSize: 0,
		Cursor:   "",
		Archived: wrapperspb.Bool(false),
	}
	resp, err := e.environmentClient.ListEnvironmentsV2(ctx, &listEnvironmentsRequest)
	if err != nil {
		return nil, err
	}
	return resp.Environments, err
}

func (e *experimentCalculate) listExperiments(
	ctx context.Context,
	environmentId string,
) ([]*experiment.Experiment, error) {
	req := &experiment.ListExperimentsRequest{
		// Because the evaluation and goal events may be sent with a delay
		// for many reasons from the client side, we still calculate
		// the results for two days after it stopped.
		StopAt:        time.Now().In(e.location).Add(-2 * day).Unix(),
		PageSize:      0,
		Cursor:        "",
		EnvironmentId: environmentId,
		Statuses: []experiment.Experiment_Status{
			experiment.Experiment_RUNNING,
			experiment.Experiment_STOPPED,
		},
	}
	resp, err := e.experimentClient.ListExperiments(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Experiments, nil
}
