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
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	domainproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	senderproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription/sender"
)

func TestLastDays(t *testing.T) {
	patterns := []struct {
		desc     string
		inputNow time.Time
		expected int
	}{
		{
			desc:     "now is after stopAt",
			inputNow: time.Date(2019, 12, 26, 00, 00, 00, 0, time.UTC),
			expected: 0,
		},
		{
			desc:     "now equals to stopAt",
			inputNow: time.Date(2019, 12, 25, 23, 59, 59, 0, time.UTC),
			expected: 0,
		},
		{
			desc:     "0",
			inputNow: time.Date(2019, 12, 25, 23, 00, 00, 0, time.UTC),
			expected: 0,
		},
		{
			desc:     "1",
			inputNow: time.Date(2019, 12, 24, 00, 00, 00, 0, time.UTC),
			expected: 1,
		},
	}
	stopAt := time.Date(2019, 12, 25, 23, 59, 59, 0, time.UTC)
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := lastDays(p.inputNow, stopAt)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestCreateDomainEventAttachment(t *testing.T) {
	t.Parallel()

	webURL := "https://example.com"

	patterns := []struct {
		desc        string
		entityType  domainproto.Event_EntityType
		entityID    string
		entityData  string
		expectedURL string
	}{
		{
			desc:        "feature entity",
			entityType:  domainproto.Event_FEATURE,
			entityID:    "feature-id-1",
			entityData:  `{"id": "feature-id-1", "name": "test-feature"}`,
			expectedURL: "https://example.com/test/features/feature-id-1",
		},
		{
			desc:        "autoops rule with feature_id",
			entityType:  domainproto.Event_AUTOOPS_RULE,
			entityID:    "rule-id-1",
			entityData:  `{"id": "rule-id-1", "feature_id": "feature-id-2", "ops_type": 1}`,
			expectedURL: "https://example.com/test/features/feature-id-2/autoops",
		},
		{
			desc:        "progressive rollout with feature_id",
			entityType:  domainproto.Event_PROGRESSIVE_ROLLOUT,
			entityID:    "rollout-id-1",
			entityData:  `{"id": "rollout-id-1", "feature_id": "feature-id-3", "type": 1}`,
			expectedURL: "https://example.com/test/features/feature-id-3/autoops",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			notification := &senderproto.DomainEventNotification{
				Editor: &domainproto.Editor{
					Email: "test@example.com",
				},
				EntityType:         p.entityType,
				EntityId:           p.entityID,
				EntityData:         p.entityData,
				Type:               domainproto.Event_FEATURE_CREATED,
				EnvironmentName:    "test-env",
				EnvironmentUrlCode: "test",
			}

			notifier := &slackNotifier{
				webURL: webURL,
				logger: zap.NewNop(),
			}

			md := metadata.New(map[string]string{
				"accept-language": locale.Ja,
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			localizer := locale.NewLocalizer(ctx)

			attachment, err := notifier.createDomainEventAttachment(notification, localizer)
			assert.NoError(t, err)
			assert.NotNil(t, attachment)

			assert.Contains(t, attachment.Text, p.expectedURL)
			assert.Equal(t, notification.Editor.Email, attachment.AuthorName)
			assert.Contains(t, attachment.Text, notification.EnvironmentName)
			assert.Contains(t, attachment.Text, notification.EntityId)
		})
	}
}

func newTestLocalizer(lang string) locale.Localizer {
	md := metadata.New(map[string]string{"accept-language": lang})
	return locale.NewLocalizer(metadata.NewIncomingContext(context.Background(), md))
}

func TestCreateFeatureStaleAttachment(t *testing.T) {
	t.Parallel()

	staleDays := strconv.Itoa(featuredomain.SecondsToStale / 24 / 60 / 60)

	patterns := []struct {
		desc         string
		lang         string
		expectedText string
	}{
		{
			desc:         "english",
			lang:         locale.En,
			expectedText: "There are feature flags that have not been used for more than " + staleDays + " days.",
		},
		{
			desc:         "japanese",
			lang:         locale.Ja,
			expectedText: staleDays + "日以上使用されていないフィーチャーフラグがあります。",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			notifier := &slackNotifier{webURL: "https://example.com", logger: zap.NewNop()}
			attachment, err := notifier.createFeatureStaleAttachment(
				&senderproto.FeatureStaleNotification{
					EnvironmentName:    "test-env",
					EnvironmentUrlCode: "test",
					Features: []*featureproto.Feature{
						{Id: "feature-id-1", Name: "feature-name-1"},
					},
				},
				newTestLocalizer(p.lang),
			)
			assert.NoError(t, err)
			assert.Contains(t, attachment.Text, p.expectedText)
			assert.Contains(t, attachment.Text, "https://example.com/test/features/feature-id-1")
		})
	}
}

func TestCreateExperimentRunningAttachment(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc            string
		lang            string
		expectedText    string
		expectedDaysMsg string
	}{
		{
			desc:            "english",
			lang:            locale.En,
			expectedText:    "There are running experiments.",
			expectedDaysMsg: "- `1` days left, Name: *",
		},
		{
			desc:            "japanese",
			lang:            locale.Ja,
			expectedText:    "実行中のエクスペリメントがあります。",
			expectedDaysMsg: "- 残り `1` 日, Name: *",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			notifier := &slackNotifier{webURL: "https://example.com", logger: zap.NewNop()}
			attachment, err := notifier.createExperimentRunningAttachment(
				&senderproto.ExperimentRunningNotification{
					EnvironmentName:    "test-env",
					EnvironmentUrlCode: "test",
					Experiments: []*experimentproto.Experiment{
						{
							Id:     "experiment-id-1",
							Name:   "experiment-name-1",
							StopAt: time.Now().Add(25 * time.Hour).Unix(),
						},
					},
				},
				newTestLocalizer(p.lang),
			)
			assert.NoError(t, err)
			assert.Contains(t, attachment.Text, p.expectedText)
			assert.Contains(t, attachment.Text, p.expectedDaysMsg)
		})
	}
}

func TestCreateMAUCountAttachment(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc         string
		lang         string
		expectedText string
	}{
		{
			desc:         "english",
			lang:         locale.En,
			expectedText: "This is the MAU for month 4.",
		},
		{
			desc:         "japanese",
			lang:         locale.Ja,
			expectedText: "4月のMAUです。",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			notifier := &slackNotifier{webURL: "https://example.com", logger: zap.NewNop()}
			attachment, err := notifier.createMAUCountAttachment(
				&senderproto.MauCountNotification{
					EnvironmentName: "test-env",
					EventCount:      1234,
					UserCount:       567,
					Month:           4,
				},
				newTestLocalizer(p.lang),
			)
			assert.NoError(t, err)
			assert.Contains(t, attachment.Text, p.expectedText)
			assert.Contains(t, attachment.Text, "Event count: 1,234")
			assert.Contains(t, attachment.Text, "User count: 567")
		})
	}
}
