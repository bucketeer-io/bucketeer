// Copyright 2024 The Bucketeer Authors.
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

package experiment

import (
	"context"
	"time"

	wrappersproto "github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

const (
	listRequestSize = 500
)

type experimentStatusUpdater struct {
	environmentClient environmentclient.Client
	experimentClient  experimentclient.Client
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewExperimentStatusUpdater(
	environmentClient environmentclient.Client,
	experimentClient experimentclient.Client,
	opts ...jobs.Option) jobs.Job {

	dopts := &jobs.Options{
		Timeout: 1 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &experimentStatusUpdater{
		environmentClient: environmentClient,
		experimentClient:  experimentClient,
		opts:              dopts,
		logger:            dopts.Logger.Named("status-updater"),
	}
}

func (u *experimentStatusUpdater) Run(ctx context.Context) (lastErr error) {
	ctx, cancel := context.WithTimeout(ctx, u.opts.Timeout)
	defer cancel()
	environments, err := u.listEnvironments(ctx)
	if err != nil {
		u.logger.Error("Failed to list environments", zap.Error(err))
		lastErr = err
		return
	}
	for _, env := range environments {
		var experiments []*experimentproto.Experiment
		statuses := []experimentproto.Experiment_Status{
			experimentproto.Experiment_WAITING,
			experimentproto.Experiment_RUNNING,
		}
		for _, status := range statuses {
			exps, err := u.listExperiments(ctx, env.Id, status)
			if err != nil {
				u.logger.Error("Failed to list experiments", zap.Error(err),
					zap.String("environmentNamespace", env.Id),
					zap.Int32("status", int32(status)),
				)
				lastErr = err
				continue
			}
			experiments = append(experiments, exps...)
		}
		for _, e := range experiments {
			if err = u.updateStatus(ctx, env.Id, e); err != nil {
				lastErr = err
			}
		}
	}
	return
}

func (u *experimentStatusUpdater) updateStatus(
	ctx context.Context,
	environmentNamespace string,
	experiment *experimentproto.Experiment,
) error {
	if experiment.Status == experimentproto.Experiment_WAITING {
		if err := u.updateToRunning(ctx, environmentNamespace, experiment); err != nil {
			return err
		}
		return nil
	}
	if experiment.Status == experimentproto.Experiment_RUNNING {
		if err := u.updateToStopped(ctx, environmentNamespace, experiment); err != nil {
			return err
		}
	}
	return nil
}

func (u *experimentStatusUpdater) updateToRunning(
	ctx context.Context,
	environmentNamespace string,
	experiment *experimentproto.Experiment,
) error {
	de := domain.Experiment{Experiment: experiment}
	if err := de.Start(); err != nil {
		if err != domain.ErrExperimentBeforeStart {
			u.logger.Error("Failed to start check if experiment running", zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("id", experiment.Id))
			return err
		}
		return nil
	}
	_, err := u.experimentClient.StartExperiment(ctx, &experimentproto.StartExperimentRequest{
		EnvironmentNamespace: environmentNamespace,
		Id:                   experiment.Id,
		Command:              &experimentproto.StartExperimentCommand{},
	})
	if err != nil {
		u.logger.Error("Failed to update status to running", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("id", experiment.Id))
		return err
	}
	return nil
}

func (u *experimentStatusUpdater) updateToStopped(
	ctx context.Context,
	environmentNamespace string,
	experiment *experimentproto.Experiment,
) error {
	de := domain.Experiment{Experiment: experiment}
	if err := de.Finish(); err != nil {
		if err != domain.ErrExperimentBeforeStop {
			u.logger.Error("Failed to end check if experiment running", zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("id", experiment.Id))
			return err
		}
		return nil
	}
	_, err := u.experimentClient.FinishExperiment(ctx, &experimentproto.FinishExperimentRequest{
		EnvironmentNamespace: environmentNamespace,
		Id:                   experiment.Id,
		Command:              &experimentproto.FinishExperimentCommand{},
	})
	if err != nil {
		u.logger.Error("Failed to update status to stopped", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("id", experiment.Id))
		return err
	}
	return nil
}

func (u *experimentStatusUpdater) listExperiments(
	ctx context.Context,
	environmentNamespace string,
	status experimentproto.Experiment_Status,
) ([]*experimentproto.Experiment, error) {
	var experiments []*experimentproto.Experiment
	cursor := ""
	for {
		resp, err := u.experimentClient.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
			PageSize:             listRequestSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentNamespace,
			Status:               &wrappersproto.Int32Value{Value: int32(status)},
		})
		if err != nil {
			return nil, err
		}
		experiments = append(experiments, resp.Experiments...)
		size := len(resp.Experiments)
		if size == 0 || size < listRequestSize {
			return experiments, nil
		}
		cursor = resp.Cursor
	}
}

func (u *experimentStatusUpdater) listEnvironments(ctx context.Context) ([]*environmentproto.EnvironmentV2, error) {
	var environments []*environmentproto.EnvironmentV2
	cursor := ""
	for {
		resp, err := u.environmentClient.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{
			PageSize: listRequestSize,
			Cursor:   cursor,
			Archived: &wrappersproto.BoolValue{Value: false},
		})
		if err != nil {
			return nil, err
		}
		environments = append(environments, resp.Environments...)
		environmentSize := len(resp.Environments)
		if environmentSize == 0 || environmentSize < listRequestSize {
			return environments, nil
		}
		cursor = resp.Cursor
	}
}
