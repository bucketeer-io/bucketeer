// Copyright 2026 The Bucketeer Authors.
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

package notifier

import (
	"context"
	"sync"
	"time"

	"github.com/slack-go/slack"
	"go.uber.org/zap"

	senderproto "github.com/bucketeer-io/bucketeer/v2/proto/notification/sender"
)

// defaultFailureAlertCooldown is the minimum interval between two alerts for
// the same job/consumer. Batch cron jobs run every minute, so without a
// cooldown a persistently failing job would post a Slack message every minute.
const defaultFailureAlertCooldown = 30 * time.Minute

// slackMessagePoster is the subset of the Slack Web API client used to post
// failure alerts. *slack.Client satisfies it; tests inject a fake.
type slackMessagePoster interface {
	PostMessageContext(
		ctx context.Context,
		channelID string,
		options ...slack.MsgOption,
	) (string, string, error)
}

type failureAlerterOptions struct {
	cooldown time.Duration
	logger   *zap.Logger
}

type FailureAlerterOption func(*failureAlerterOptions)

func WithFailureAlertCooldown(d time.Duration) FailureAlerterOption {
	return func(o *failureAlerterOptions) {
		o.cooldown = d
	}
}

func WithFailureAlertLogger(l *zap.Logger) FailureAlerterOption {
	return func(o *failureAlerterOptions) {
		o.logger = l
	}
}

// NewFailureAlerter returns a FailureAlerter that posts to the given Slack
// channel through the Slack Web API (chat.postMessage), authenticated with the
// given bot token. When token or channel is empty the alerter is disabled (a
// no-op is returned).
func NewFailureAlerter(
	token, channel string,
	opts ...FailureAlerterOption,
) FailureAlerter {
	if token == "" || channel == "" {
		return &noopFailureAlerter{}
	}
	options := &failureAlerterOptions{
		cooldown: defaultFailureAlertCooldown,
		logger:   zap.NewNop(),
	}
	for _, opt := range opts {
		opt(options)
	}
	return &failureAlerter{
		poster:   slack.New(token),
		channel:  channel,
		cooldown: options.cooldown,
		now:      time.Now,
		lastSent: make(map[string]time.Time),
		logger:   options.logger.Named("failure-alerter"),
	}
}

type failureAlerter struct {
	poster   slackMessagePoster
	channel  string
	cooldown time.Duration
	now      func() time.Time
	mu       sync.Mutex
	lastSent map[string]time.Time
	logger   *zap.Logger
}

func (a *failureAlerter) NotifyBatchJobFailure(ctx context.Context, jobName string, jobErr error) {
	a.notify(ctx, senderproto.JobFailureNotification_BATCH, jobName, jobErr)
}

func (a *failureAlerter) NotifySubscriberFailure(
	ctx context.Context,
	consumerName string,
	consumerErr error,
) {
	a.notify(ctx, senderproto.JobFailureNotification_SUBSCRIBER, consumerName, consumerErr)
}

func (a *failureAlerter) notify(
	ctx context.Context,
	serviceType senderproto.JobFailureNotification_ServiceType,
	name string,
	failureErr error,
) {
	if failureErr == nil {
		return
	}
	key := serviceType.String() + ":" + name
	if a.throttled(key) {
		a.logger.Debug("Failure alert throttled within cooldown",
			zap.String("serviceType", serviceType.String()),
			zap.String("name", name),
			zap.Duration("cooldown", a.cooldown),
		)
		return
	}
	attachment := jobFailureAttachment(&senderproto.JobFailureNotification{
		ServiceType:  serviceType,
		JobName:      name,
		ErrorMessage: failureErr.Error(),
	})
	if _, _, err := a.poster.PostMessageContext(
		ctx,
		a.channel,
		slack.MsgOptionAttachments(*attachment),
	); err != nil {
		a.logger.Error("Failed to send failure alert",
			zap.Error(err),
			zap.String("serviceType", serviceType.String()),
			zap.String("name", name),
			zap.String("channel", a.channel),
		)
		// Do not start the cooldown when the alert was not delivered, so the
		// next failure can try again.
		return
	}
	a.markSent(key)
}

// throttled reports whether an alert for key is still within the cooldown
// window of the previous successful alert.
func (a *failureAlerter) throttled(key string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	last, ok := a.lastSent[key]
	return ok && a.now().Sub(last) < a.cooldown
}

func (a *failureAlerter) markSent(key string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lastSent[key] = a.now()
}

type noopFailureAlerter struct{}

func (n *noopFailureAlerter) NotifyBatchJobFailure(context.Context, string, error)   {}
func (n *noopFailureAlerter) NotifySubscriberFailure(context.Context, string, error) {}
