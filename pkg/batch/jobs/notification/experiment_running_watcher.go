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

package notification

import (
	"context"
	"time"

	wrappersproto "github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

const (
	listRequestSize = 500
)

type experimentRunningWatcher struct {
	environmentClient environmentclient.Client
	experimentClient  experimentclient.Client
	sender            sender.Sender
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewExperimentRunningWatcher(
	environmentClient environmentclient.Client,
	experimentClient experimentclient.Client,
	sender sender.Sender,
	opts ...jobs.Option) jobs.Job {

	dopts := &jobs.Options{
		Timeout: 5 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &experimentRunningWatcher{
		environmentClient: environmentClient,
		experimentClient:  experimentClient,
		sender:            sender,
		opts:              dopts,
		logger:            dopts.Logger.Named("experiment-running-watcher"),
	}
}

func (w *experimentRunningWatcher) Run(ctx context.Context) (lastErr error) {
	ctx, cancel := context.WithTimeout(ctx, w.opts.Timeout)
	defer cancel()
	environments, err := w.listEnvironments(ctx)
	if err != nil {
		return err
	}
	for _, env := range environments {
		experiments, err := w.listExperiments(ctx, env.Id)
		if err != nil {
			return err
		}
		if len(experiments) == 0 {
			continue
		}
		ne, err := w.createNotificationEvent(env, experiments)
		if err != nil {
			lastErr = err
		}
		if err := w.sender.Send(ctx, ne); err != nil {
			lastErr = err
		}
	}
	return
}

func (w *experimentRunningWatcher) createNotificationEvent(
	environment *environmentproto.EnvironmentV2,
	experiments []*experimentproto.Experiment,
) (*senderproto.NotificationEvent, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	ne := &senderproto.NotificationEvent{
		Id:                   id.String(),
		EnvironmentNamespace: environment.Id,
		SourceType:           notificationproto.Subscription_EXPERIMENT_RUNNING,
		Notification: &senderproto.Notification{
			Type: senderproto.Notification_ExperimentRunning,
			ExperimentRunningNotification: &senderproto.ExperimentRunningNotification{
				EnvironmentId: environment.Id,
				Experiments:   experiments,
			},
		},
		IsAdminEvent: false,
	}
	return ne, nil
}

func (w *experimentRunningWatcher) listEnvironments(ctx context.Context) ([]*environmentproto.EnvironmentV2, error) {
	var environments []*environmentproto.EnvironmentV2
	cursor := ""
	for {
		resp, err := w.environmentClient.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{
			PageSize: listRequestSize,
			Cursor:   cursor,
			Archived: wrapperspb.Bool(false),
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

func (w *experimentRunningWatcher) listExperiments(
	ctx context.Context,
	environmentNamespace string,
) ([]*experimentproto.Experiment, error) {
	var experiments []*experimentproto.Experiment
	cursor := ""
	for {
		resp, err := w.experimentClient.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
			PageSize:             listRequestSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentNamespace,
			Status:               &wrappersproto.Int32Value{Value: int32(experimentproto.Experiment_RUNNING)},
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
