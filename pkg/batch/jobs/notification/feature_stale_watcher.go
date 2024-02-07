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

	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

type featureStaleWatcher struct {
	environmentClient environmentclient.Client
	featureClient     featureclient.Client
	sender            sender.Sender
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewFeatureStaleWatcher(
	environmentClient environmentclient.Client,
	featureClient featureclient.Client,
	sender sender.Sender,
	opts ...jobs.Option) jobs.Job {

	dopts := &jobs.Options{
		Timeout: 5 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &featureStaleWatcher{
		environmentClient: environmentClient,
		featureClient:     featureClient,
		sender:            sender,
		opts:              dopts,
		logger:            dopts.Logger.Named("feature-stale-watcher"),
	}
}

func (w *featureStaleWatcher) Run(ctx context.Context) (lastErr error) {
	ctx, cancel := context.WithTimeout(ctx, w.opts.Timeout)
	defer cancel()
	environments, err := w.listEnvironments(ctx)
	if err != nil {
		return err
	}
	for _, env := range environments {
		features, err := w.listFeatures(ctx, env.Id)
		if err != nil {
			return err
		}
		var staleFeatures []*featureproto.Feature
		for _, f := range features {
			fd := &featuredomain.Feature{Feature: f}
			now := time.Now()
			stale := fd.IsStale(now)
			if !stale {
				continue
			}
			staleFeatures = append(staleFeatures, fd.Feature)
		}
		if len(staleFeatures) == 0 {
			continue
		}
		ne, err := w.createNotificationEvent(env, staleFeatures)
		if err != nil {
			lastErr = err
		}
		if err := w.sender.Send(ctx, ne); err != nil {
			lastErr = err
		}
	}
	return
}

func (w *featureStaleWatcher) createNotificationEvent(
	environment *environmentproto.EnvironmentV2,
	features []*featureproto.Feature,
) (*senderproto.NotificationEvent, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	ne := &senderproto.NotificationEvent{
		Id:                   id.String(),
		EnvironmentNamespace: environment.Id,
		SourceType:           notificationproto.Subscription_FEATURE_STALE,
		Notification: &senderproto.Notification{
			Type: senderproto.Notification_FeatureStale,
			FeatureStaleNotification: &senderproto.FeatureStaleNotification{
				EnvironmentId: environment.Id,
				Features:      features,
			},
		},
		IsAdminEvent: false,
	}
	return ne, nil
}

func (w *featureStaleWatcher) listEnvironments(ctx context.Context) ([]*environmentproto.EnvironmentV2, error) {
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

func (w *featureStaleWatcher) listFeatures(
	ctx context.Context,
	environmentNamespace string,
) ([]*featureproto.Feature, error) {
	var features []*featureproto.Feature
	cursor := ""
	for {
		resp, err := w.featureClient.ListFeatures(ctx, &featureproto.ListFeaturesRequest{
			PageSize:             listRequestSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentNamespace,
			Archived:             &wrappers.BoolValue{Value: false},
		})
		if err != nil {
			return nil, err
		}
		for _, f := range resp.Features {
			ff := featuredomain.Feature{Feature: f}
			if ff.IsDisabledAndOffVariationEmpty() {
				continue
			}
			features = append(features, f)
		}
		featureSize := len(resp.Features)
		if featureSize == 0 || featureSize < listRequestSize {
			return features, nil
		}
		cursor = resp.Cursor
	}
}
