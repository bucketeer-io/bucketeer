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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/metadata"

	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	subscriptiondomain "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/domain"
	domainproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	subscriptionproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
	senderproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription/sender"
)

const testSlackWebhookURL = "https://hooks.slack.com/services/T000/B000/token"

// roundTripperFunc lets a test intercept the outgoing webhook request
// without opening a network connection.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newMAUCountNotification() *senderproto.Notification {
	return &senderproto.Notification{
		Type: senderproto.Notification_MauCount,
		MauCountNotification: &senderproto.MauCountNotification{
			EnvironmentName: "test-env",
			EventCount:      1234,
			UserCount:       567,
			Month:           4,
		},
	}
}

func newSlackRecipient(webhookURL string) *subscriptionproto.Recipient {
	return &subscriptionproto.Recipient{
		Type:     subscriptionproto.Recipient_SlackChannel,
		Language: subscriptionproto.Recipient_ENGLISH,
		SlackChannelRecipient: &subscriptionproto.SlackChannelRecipient{
			WebhookUrl: webhookURL,
		},
	}
}

func TestNewWebhookHTTPClient(t *testing.T) {
	t.Parallel()

	client := newWebhookHTTPClient()
	assert.Equal(t, webhookTimeout, client.Timeout)
	assert.NotZero(t, client.Timeout)
	require.NotNil(t, client.CheckRedirect)
	assert.Equal(t, http.ErrUseLastResponse, client.CheckRedirect(nil, nil))
}

func TestNotifyDeliversOnlyToSlackIncomingWebhooks(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc             string
		webhookURL       string
		expectedErr      error
		expectedRequests int32
	}{
		{
			desc:             "err: loopback address",
			webhookURL:       "http://127.0.0.1:3100/internal",
			expectedErr:      subscriptiondomain.ErrSlackWebhookURLInvalid,
			expectedRequests: 0,
		},
		{
			desc:             "err: link-local address",
			webhookURL:       "http://169.254.169.254/latest/meta-data/",
			expectedErr:      subscriptiondomain.ErrSlackWebhookURLInvalid,
			expectedRequests: 0,
		},
		{
			desc:             "err: other https host",
			webhookURL:       "https://example.com/services/T000/B000/token",
			expectedErr:      subscriptiondomain.ErrSlackWebhookURLInvalid,
			expectedRequests: 0,
		},
		{
			desc:             "success: Slack incoming webhook",
			webhookURL:       testSlackWebhookURL,
			expectedErr:      nil,
			expectedRequests: 1,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			var requests int32
			client := newWebhookHTTPClient()
			client.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
				atomic.AddInt32(&requests, 1)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, testSlackWebhookURL, req.URL.String())
				assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				var msg slack.WebhookMessage
				require.NoError(t, json.Unmarshal(body, &msg))
				require.Len(t, msg.Attachments, 1)
				assert.Contains(t, msg.Attachments[0].Text, "test-env")
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       io.NopCloser(strings.NewReader("ok")),
					Request:    req,
				}, nil
			})
			notifier := &slackNotifier{
				webURL:     "https://example.com",
				httpClient: client,
				logger:     zap.NewNop(),
			}
			recipient := newSlackRecipient(p.webhookURL)
			err := notifier.Notify(context.Background(), newMAUCountNotification(), recipient, recipient.Language)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedRequests, atomic.LoadInt32(&requests))
		})
	}
}

func TestPostWebhook(t *testing.T) {
	t.Parallel()

	var redirectTargetHits int32
	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/redirect-target", http.StatusFound)
	})
	mux.HandleFunc("/redirect-target", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&redirectTargetHits, 1)
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	patterns := []struct {
		desc        string
		ctx         context.Context
		path        string
		expectedErr bool
	}{
		{
			desc:        "success: 200 OK",
			ctx:         context.Background(),
			path:        "/ok",
			expectedErr: false,
		},
		{
			desc:        "err: redirect is returned to the caller instead of being followed",
			ctx:         context.Background(),
			path:        "/redirect",
			expectedErr: true,
		},
		{
			desc:        "err: canceled context is honored",
			ctx:         canceledCtx,
			path:        "/ok",
			expectedErr: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			notifier := &slackNotifier{
				webURL:     "https://example.com",
				httpClient: newWebhookHTTPClient(),
				logger:     zap.NewNop(),
			}
			msg := &slack.WebhookMessage{Text: "hello"}
			err := notifier.postWebhook(p.ctx, msg, server.URL+p.path)
			if p.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	assert.Equal(t, int32(0), atomic.LoadInt32(&redirectTargetHits))
}

func TestPostWebhookMasksWebhookURLInErrors(t *testing.T) {
	t.Parallel()

	maskedURL := subscriptiondomain.MaskWebhookURL(testSlackWebhookURL)
	patterns := []struct {
		desc            string
		roundTrip       func(*http.Request) (*http.Response, error)
		expectedErr     bool
		expectedErrType error
		expectedErrText string
	}{
		{
			desc: "err: transport failure embeds the url and is masked",
			roundTrip: func(*http.Request) (*http.Response, error) {
				return nil, errors.New("dial tcp: connection refused")
			},
			expectedErr:     true,
			expectedErrText: "connection refused",
		},
		{
			desc: "err: status code failure is returned unchanged",
			roundTrip: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body:       io.NopCloser(strings.NewReader("")),
					Request:    req,
				}, nil
			},
			expectedErr:     true,
			expectedErrType: slack.StatusCodeError{},
			expectedErrText: "500",
		},
		{
			desc: "success: 200 OK",
			roundTrip: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       io.NopCloser(strings.NewReader("ok")),
					Request:    req,
				}, nil
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			core, logs := observer.New(zapcore.DebugLevel)
			client := newWebhookHTTPClient()
			client.Transport = roundTripperFunc(p.roundTrip)
			notifier := &slackNotifier{
				webURL:     "https://example.com",
				httpClient: client,
				logger:     zap.New(core),
			}
			err := notifier.postWebhook(context.Background(), &slack.WebhookMessage{Text: "hello"}, testSlackWebhookURL)
			for _, entry := range logs.All() {
				assert.NotContains(t, entry.Message, testSlackWebhookURL)
				for key, field := range entry.ContextMap() {
					assert.NotContains(t, fmt.Sprint(field), testSlackWebhookURL,
						"field %q of %q", key, entry.Message)
				}
			}
			if !p.expectedErr {
				require.NoError(t, err)
				assert.Empty(t, logs.All())
				return
			}
			require.Error(t, err)
			assert.NotContains(t, err.Error(), testSlackWebhookURL)
			assert.Contains(t, err.Error(), p.expectedErrText)
			if p.expectedErrType != nil {
				assert.IsType(t, p.expectedErrType, err)
			} else {
				assert.Contains(t, err.Error(), maskedURL)
			}
			entries := logs.FilterMessage("Failed to post a message").All()
			require.Len(t, entries, 1)
			assert.Equal(t, err.Error(), entries[0].ContextMap()["error"])
			assert.Equal(t,
				subscriptiondomain.SlackChannelRecipientID(testSlackWebhookURL),
				entries[0].ContextMap()["slackRecipientId"],
			)
		})
	}
}

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
			expectedDaysMsg: "- Days left: `1`, Name: *",
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
