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
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

type mauCountWatcher struct {
	environmentClient  environmentclient.Client
	eventCounterClient ecclient.Client
	sender             sender.Sender
	location           *time.Location
	opts               *jobs.Options
	logger             *zap.Logger
}

func NewMAUCountWatcher(
	environmentClient environmentclient.Client,
	eventCounterClient ecclient.Client,
	sender sender.Sender,
	location *time.Location,
	opts ...jobs.Option) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 5 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &mauCountWatcher{
		environmentClient:  environmentClient,
		eventCounterClient: eventCounterClient,
		sender:             sender,
		location:           location,
		opts:               dopts,
		logger:             dopts.Logger.Named("mau-count-watcher"),
	}
}

func (w *mauCountWatcher) Run(ctx context.Context) (lastErr error) {
	ctx, cancel := context.WithTimeout(ctx, w.opts.Timeout)
	defer cancel()
	projects, err := w.listProjects(ctx)
	if err != nil {
		return err
	}
	year, lastMonth := w.getLastYearMonth(time.Now().In(w.location))
	for _, pj := range projects {
		environments, err := w.listEnvironments(ctx, pj.Id)
		if err != nil {
			return err
		}
		for _, env := range environments {
			eventCount, userCount, err := w.getUserCount(ctx, env.Id, w.newYearMonth(year, lastMonth))
			if err != nil {
				return err
			}
			if err := w.sendNotification(ctx, env, eventCount, userCount, lastMonth); err != nil {
				w.logger.Error("Failed to send notification",
					zap.Error(err),
					zap.String("projectId", pj.Id),
					zap.String("environmentId", env.Id),
					zap.Int64("eventCount", eventCount),
					zap.Int64("userCount", userCount),
					zap.Int32("lastMonth", lastMonth),
				)
				lastErr = err
			}
		}
	}
	return
}

func (w *mauCountWatcher) listProjects(ctx context.Context) ([]*environmentproto.Project, error) {
	var projects []*environmentproto.Project
	cursor := ""
	for {
		resp, err := w.environmentClient.ListProjects(ctx, &environmentproto.ListProjectsRequest{
			PageSize: listRequestSize,
			Cursor:   cursor,
		})
		if err != nil {
			return nil, err
		}
		projects = append(projects, resp.Projects...)
		projectSize := len(resp.Projects)
		if projectSize == 0 || projectSize < listRequestSize {
			return projects, nil
		}
		cursor = resp.Cursor
	}
}

func (w *mauCountWatcher) getLastYearMonth(now time.Time) (int32, int32) {
	targetDate := now.AddDate(0, -1, 0)
	return int32(targetDate.Year()), int32(targetDate.Month())
}

func (w *mauCountWatcher) newYearMonth(year, month int32) string {
	return fmt.Sprintf("%d%02d", year, month)
}

func (w *mauCountWatcher) listEnvironments(
	ctx context.Context,
	projectID string,
) ([]*environmentproto.EnvironmentV2, error) {
	var environments []*environmentproto.EnvironmentV2
	cursor := ""
	for {
		resp, err := w.environmentClient.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{
			PageSize:  listRequestSize,
			Cursor:    cursor,
			ProjectId: projectID,
			Archived:  wrapperspb.Bool(false),
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

func (w *mauCountWatcher) getUserCount(
	ctx context.Context,
	environmentNamespace, yearMonth string,
) (eventCount, userCount int64, err error) {
	resp, e := w.eventCounterClient.GetMAUCount(ctx, &ecproto.GetMAUCountRequest{
		EnvironmentNamespace: environmentNamespace,
		YearMonth:            yearMonth,
	})
	if e != nil {
		err = e
		return
	}
	eventCount = resp.EventCount
	userCount = resp.UserCount
	return
}

func (w *mauCountWatcher) sendNotification(
	ctx context.Context,
	environment *environmentproto.EnvironmentV2,
	eventCount, userCount int64,
	month int32,
) error {
	ne, err := w.createNotificationEvent(environment, eventCount, userCount, month)
	if err != nil {
		return err
	}
	if err := w.sender.Send(ctx, ne); err != nil {
		return err
	}
	return nil
}

func (w *mauCountWatcher) createNotificationEvent(
	environment *environmentproto.EnvironmentV2,
	eventCount, userCount int64,
	month int32,
) (*senderproto.NotificationEvent, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	ne := &senderproto.NotificationEvent{
		Id:                   id.String(),
		EnvironmentNamespace: environment.Id,
		SourceType:           notificationproto.Subscription_MAU_COUNT,
		Notification: &senderproto.Notification{
			Type: senderproto.Notification_MauCount,
			MauCountNotification: &senderproto.MauCountNotification{
				EnvironmentId: environment.Id,
				EventCount:    eventCount,
				UserCount:     userCount,
				Month:         month,
			},
		},
		IsAdminEvent: false,
	}
	return ne, nil
}
