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
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

// fakeSlackPoster records each PostMessageContext call and returns optional
// per-call errors (results[i] for the i-th call, nil when out of range).
type fakeSlackPoster struct {
	mu      sync.Mutex
	calls   []string // target channel of each call
	results []error
}

func (f *fakeSlackPoster) PostMessageContext(
	_ context.Context,
	channelID string,
	_ ...slack.MsgOption,
) (string, string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	i := len(f.calls)
	f.calls = append(f.calls, channelID)
	if i < len(f.results) {
		return "", "", f.results[i]
	}
	return "", "", nil
}

func (f *fakeSlackPoster) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

func TestFailureAlerterDisabledWhenCredentialsEmpty(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc    string
		token   string
		channel string
	}{
		{desc: "empty token", token: "", channel: "#bucketeer-emergency"},
		{desc: "empty channel", token: "xoxb-token", channel: ""},
		{desc: "both empty", token: "", channel: ""},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			alerter := NewFailureAlerter(p.token, p.channel)
			_, isNoop := alerter.(*noopFailureAlerter)
			assert.True(t, isNoop)
			// Must not panic when disabled.
			alerter.NotifyBatchJobFailure(context.Background(), "job-a", errors.New("boom"))
			alerter.NotifySubscriberFailure(context.Background(), "consumer-a", errors.New("boom"))
		})
	}
}

func TestFailureAlerterPostsToConfiguredChannel(t *testing.T) {
	t.Parallel()
	alerter := NewFailureAlerter("xoxb-token", "#bucketeer-emergency").(*failureAlerter)
	fake := &fakeSlackPoster{}
	alerter.poster = fake

	alerter.NotifyBatchJobFailure(context.Background(), "ExperimentCalculator", errors.New("boom"))

	assert.Equal(t, []string{"#bucketeer-emergency"}, fake.calls)
}

func TestFailureAlerterIgnoresNilError(t *testing.T) {
	t.Parallel()
	alerter := NewFailureAlerter("xoxb-token", "#bucketeer-emergency").(*failureAlerter)
	fake := &fakeSlackPoster{}
	alerter.poster = fake

	alerter.NotifyBatchJobFailure(context.Background(), "job-a", nil)

	assert.Equal(t, 0, fake.callCount())
}

func TestFailureAlerterThrottlesWithinCooldown(t *testing.T) {
	t.Parallel()
	alerter := NewFailureAlerter(
		"xoxb-token",
		"#bucketeer-emergency",
		WithFailureAlertCooldown(30*time.Minute),
	).(*failureAlerter)
	fake := &fakeSlackPoster{}
	alerter.poster = fake
	now := time.Unix(0, 0).UTC()
	alerter.now = func() time.Time { return now }

	ctx := context.Background()
	jobErr := errors.New("boom")
	alerter.NotifyBatchJobFailure(ctx, "job-a", jobErr) // delivered (1)
	alerter.NotifyBatchJobFailure(ctx, "job-a", jobErr) // throttled
	now = now.Add(29 * time.Minute)
	alerter.NotifyBatchJobFailure(ctx, "job-a", jobErr) // still throttled
	now = now.Add(2 * time.Minute)                      // 31m total > cooldown
	alerter.NotifyBatchJobFailure(ctx, "job-a", jobErr) // delivered (2)

	assert.Equal(t, 2, fake.callCount())
}

func TestFailureAlerterThrottleIsPerKey(t *testing.T) {
	t.Parallel()
	alerter := NewFailureAlerter(
		"xoxb-token",
		"#bucketeer-emergency",
		WithFailureAlertCooldown(time.Hour),
	).(*failureAlerter)
	fake := &fakeSlackPoster{}
	alerter.poster = fake
	now := time.Unix(0, 0).UTC()
	alerter.now = func() time.Time { return now }

	ctx := context.Background()
	jobErr := errors.New("boom")
	alerter.NotifyBatchJobFailure(ctx, "job-a", jobErr)   // delivered: batch:job-a
	alerter.NotifyBatchJobFailure(ctx, "job-b", jobErr)   // delivered: batch:job-b
	alerter.NotifySubscriberFailure(ctx, "job-a", jobErr) // delivered: subscriber:job-a
	alerter.NotifyBatchJobFailure(ctx, "job-a", jobErr)   // throttled: batch:job-a

	assert.Equal(t, 3, fake.callCount())
}

func TestFailureAlerterRetriesAfterFailedDelivery(t *testing.T) {
	t.Parallel()
	alerter := NewFailureAlerter(
		"xoxb-token",
		"#bucketeer-emergency",
		WithFailureAlertCooldown(time.Hour),
	).(*failureAlerter)
	// First delivery fails; a failed delivery must not start the cooldown, so
	// the next failure for the same key is attempted again within the window.
	fake := &fakeSlackPoster{results: []error{errors.New("slack unavailable")}}
	alerter.poster = fake
	now := time.Unix(0, 0).UTC()
	alerter.now = func() time.Time { return now }

	ctx := context.Background()
	alerter.NotifyBatchJobFailure(ctx, "job-a", errors.New("boom")) // attempt 1: fails
	alerter.NotifyBatchJobFailure(ctx, "job-a", errors.New("boom")) // attempt 2: succeeds

	assert.Equal(t, 2, fake.callCount())
}
