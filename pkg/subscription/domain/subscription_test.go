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

package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

func TestValidateSlackWebhookURL(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		input    string
		expected error
	}{
		{
			desc:     "err: empty",
			input:    "",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: not a URL",
			input:    "not-a-url",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: unparsable",
			input:    "https://hooks.slack.com/services/%zz",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: http scheme",
			input:    "http://hooks.slack.com/services/T000/B000/token",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: other host",
			input:    "https://example.com/services/T000/B000/token",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: host that only ends with the Slack host",
			input:    "https://evilhooks.slack.com/services/T000/B000/token",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: host that only starts with the Slack host",
			input:    "https://hooks.slack.com.example.com/services/T000/B000/token",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: Slack host in the path only",
			input:    "https://example.com/hooks.slack.com/services/T000/B000/token",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: Slack host in the userinfo only",
			input:    "https://hooks.slack.com@example.com/services/T000/B000/token",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: userinfo",
			input:    "https://user:pass@hooks.slack.com/services/T000/B000/token",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: explicit port",
			input:    "https://hooks.slack.com:8443/services/T000/B000/token",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: loopback address",
			input:    "http://127.0.0.1:3100/internal",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: link-local address",
			input:    "http://169.254.169.254/latest/meta-data/",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: no path",
			input:    "https://hooks.slack.com",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: root path",
			input:    "https://hooks.slack.com/",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: path outside /services/",
			input:    "https://hooks.slack.com/other/T000/B000/token",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: empty /services/ path",
			input:    "https://hooks.slack.com/services/",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "err: /services without trailing slash",
			input:    "https://hooks.slack.com/services",
			expected: ErrSlackWebhookURLInvalid,
		},
		{
			desc:     "success: Slack incoming webhook",
			input:    "https://hooks.slack.com/services/T000/B000/token",
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, p.expected, ValidateSlackWebhookURL(p.input))
		})
	}
}

func TestMaskWebhookURL(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "empty: nothing to mask",
			input:    "",
			expected: "",
		},
		{
			desc:     "Slack incoming webhook: only scheme and host kept",
			input:    "https://hooks.slack.com/services/T000/B000/token",
			expected: "https://hooks.slack.com/****",
		},
		{
			desc:     "userinfo, port, query and fragment dropped",
			input:    "https://user:pass@example.com:8443/path/token?key=value#fragment",
			expected: "https://example.com:8443/****",
		},
		{
			desc:     "no path: still masked",
			input:    "https://hooks.slack.com",
			expected: "https://hooks.slack.com/****",
		},
		{
			desc:     "not a URL: replaced entirely",
			input:    "not-a-url",
			expected: "****",
		},
		{
			desc:     "unparsable: replaced entirely",
			input:    "https://hooks.slack.com/services/%zz",
			expected: "****",
		},
		{
			desc:     "already masked: unchanged",
			input:    "https://hooks.slack.com/****",
			expected: "https://hooks.slack.com/****",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, p.expected, MaskWebhookURL(p.input))
		})
	}
}

func newSlackSubscription(webhookURL string) *proto.Subscription {
	return &proto.Subscription{
		Id:          "id",
		Name:        "name",
		SourceTypes: []proto.Subscription_SourceType{proto.Subscription_DOMAIN_EVENT_FEATURE},
		Recipient: &proto.Recipient{
			Type:                  proto.Recipient_SlackChannel,
			SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: webhookURL},
			Language:              proto.Recipient_JAPANESE,
		},
		FeatureFlagTags: []string{"tag"},
	}
}

func TestMaskSubscriptionWebhookURL(t *testing.T) {
	t.Parallel()

	t.Run("nil subscription: no panic", func(t *testing.T) {
		t.Parallel()
		MaskSubscriptionWebhookURL(nil)
	})
	t.Run("no recipient: untouched", func(t *testing.T) {
		t.Parallel()
		subscription := &proto.Subscription{Id: "id", Name: "name"}
		MaskSubscriptionWebhookURL(subscription)
		assert.Equal(t, &proto.Subscription{Id: "id", Name: "name"}, subscription)
	})
	t.Run("no slack channel recipient: untouched", func(t *testing.T) {
		t.Parallel()
		subscription := &proto.Subscription{Id: "id", Recipient: &proto.Recipient{Type: proto.Recipient_SlackChannel}}
		MaskSubscriptionWebhookURL(subscription)
		assert.Equal(
			t,
			&proto.Subscription{Id: "id", Recipient: &proto.Recipient{Type: proto.Recipient_SlackChannel}},
			subscription,
		)
	})
	t.Run("webhook URL masked in place, other fields kept", func(t *testing.T) {
		t.Parallel()
		subscription := newSlackSubscription("https://hooks.slack.com/services/T000/B000/token")
		MaskSubscriptionWebhookURL(subscription)
		expected := newSlackSubscription("https://hooks.slack.com/****")
		assert.Equal(t, expected, subscription)
	})
}

func TestCopyWithMaskedWebhookURL(t *testing.T) {
	t.Parallel()

	t.Run("nil subscription", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, CopyWithMaskedWebhookURL(nil))
	})
	t.Run("webhook URL masked in the copy, original untouched", func(t *testing.T) {
		t.Parallel()
		original := newSlackSubscription("https://hooks.slack.com/services/T000/B000/token")
		masked := CopyWithMaskedWebhookURL(original)
		assert.Equal(t, "https://hooks.slack.com/****", masked.Recipient.SlackChannelRecipient.WebhookUrl)
		assert.Equal(t, original.Id, masked.Id)
		assert.Equal(t, original.Name, masked.Name)
		assert.Equal(t, original.SourceTypes, masked.SourceTypes)
		assert.Equal(t, original.Recipient.Type, masked.Recipient.Type)
		assert.Equal(t, original.Recipient.Language, masked.Recipient.Language)
		assert.Equal(t, original.FeatureFlagTags, masked.FeatureFlagTags)
		assert.Equal(
			t,
			"https://hooks.slack.com/services/T000/B000/token",
			original.Recipient.SlackChannelRecipient.WebhookUrl,
		)
	})
}

func TestNewNotification(t *testing.T) {
	t.Parallel()
	sourceTypes := []proto.Subscription_SourceType{
		proto.Subscription_DOMAIN_EVENT_ACCOUNT,
		proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
	}
	recipient := &proto.Recipient{
		Type:                  proto.Recipient_SlackChannel,
		SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
	}
	tags := []string{"tag"}
	name := "sname"
	actual, err := NewSubscription(name, sourceTypes, recipient, tags)
	assert.NoError(t, err)
	assert.IsType(t, &Subscription{}, actual)
	assert.NotEqual(t, "", actual.Id)
	assert.Equal(t, sourceTypes, actual.SourceTypes)
	assert.Equal(t, recipient, actual.Recipient)
	assert.Equal(t, false, actual.Disabled)
	assert.NotEqual(t, 0, actual.CreatedAt)
	assert.NotEqual(t, 0, actual.UpdatedAt)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, tags, actual.FeatureFlagTags)
}

func TestUpdateNotification(t *testing.T) {
	t.Parallel()

	type input struct {
		name            *wrapperspb.StringValue
		sourceTypes     []proto.Subscription_SourceType
		disabled        *wrapperspb.BoolValue
		featureFlagTags []string
	}

	patterns := []struct {
		desc        string
		origin      *Subscription
		inputData   *input
		expected    *Subscription
		expectedErr error
	}{
		{
			desc: "Update name & sourceTypes",
			origin: &Subscription{
				&proto.Subscription{
					Id:   "id",
					Name: "origin",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{
							WebhookUrl: "https://slack-hooks.exp",
						},
					},
				},
			},
			inputData: &input{
				name: wrapperspb.String("new-name"),
				sourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_PUSH,
					proto.Subscription_DOMAIN_EVENT_SUBSCRIPTION,
				},
			},
			expected: &Subscription{
				&proto.Subscription{
					Id:   "id",
					Name: "new-name",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_PUSH,
						proto.Subscription_DOMAIN_EVENT_SUBSCRIPTION,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{
							WebhookUrl: "https://slack-hooks.exp",
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "disabled subscription",
			origin: &Subscription{
				&proto.Subscription{
					Id:   "id",
					Name: "origin",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{
							WebhookUrl: "https://slack-hooks.exp",
						},
					},
					Disabled: false,
				},
			},
			inputData: &input{
				disabled: wrapperspb.Bool(true),
			},
			expected: &Subscription{
				&proto.Subscription{
					Id:   "id",
					Name: "origin",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{
							WebhookUrl: "https://slack-hooks.exp",
						},
					},
					Disabled: true,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "err: update subscription's feature flag tags with no feature source type",
			origin: &Subscription{
				&proto.Subscription{
					Id:   "id",
					Name: "origin",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{
							WebhookUrl: "https://slack-hooks.exp",
						},
					},
					Disabled:        false,
					FeatureFlagTags: []string{"tag"},
				},
			},
			inputData: &input{
				sourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
				},
				featureFlagTags: []string{"update-tag"},
			},
			expected:    nil,
			expectedErr: ErrCannotUpdateFeatureFlagTags,
		},
		{
			desc: "success: reset feature flag tags when the feature source type is also being deleted",
			origin: &Subscription{
				&proto.Subscription{
					Id:   "id",
					Name: "origin",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{
							WebhookUrl: "https://slack-hooks.exp",
						},
					},
					Disabled:        false,
					FeatureFlagTags: []string{"tag"},
				},
			},
			inputData: &input{
				sourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
				},
			},
			expected: &Subscription{
				&proto.Subscription{
					Id:   "id",
					Name: "origin",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{
							WebhookUrl: "https://slack-hooks.exp",
						},
					},
					Disabled:        false,
					FeatureFlagTags: []string{},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: update subscription's feature flag tags",
			origin: &Subscription{
				&proto.Subscription{
					Id:   "id",
					Name: "origin",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{
							WebhookUrl: "https://slack-hooks.exp",
						},
					},
					Disabled:        false,
					FeatureFlagTags: []string{"tag"},
				},
			},
			inputData: &input{
				featureFlagTags: []string{"update-tag"},
			},
			expected: &Subscription{
				&proto.Subscription{
					Id:   "id",
					Name: "origin",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{
							WebhookUrl: "https://slack-hooks.exp",
						},
					},
					Disabled:        false,
					FeatureFlagTags: []string{"update-tag"},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := p.origin.UpdateSubscription(
				p.inputData.name,
				p.inputData.sourceTypes,
				p.inputData.disabled,
				p.inputData.featureFlagTags,
			)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.InDelta(t, time.Now().Unix(), actual.UpdatedAt, 5)
				p.expected.UpdatedAt = actual.UpdatedAt
			}
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestDisable(t *testing.T) {
	t.Parallel()
	actual := &Subscription{&proto.Subscription{Disabled: false}}
	actual.Disable()
	assert.Equal(t, true, actual.Disabled)
}

func TestEnable(t *testing.T) {
	t.Parallel()
	actual := &Subscription{&proto.Subscription{Disabled: true}}
	actual.Enable()
	assert.Equal(t, false, actual.Disabled)
}

func TestAddSourceTypes(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		origin      *Subscription
		input       []proto.Subscription_SourceType
		expectedErr error
		expected    []proto.Subscription_SourceType
	}{
		{
			desc: "success: one",
			origin: &Subscription{&proto.Subscription{SourceTypes: []proto.Subscription_SourceType{
				proto.Subscription_DOMAIN_EVENT_ACCOUNT,
				proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
			}}},
			input:       []proto.Subscription_SourceType{proto.Subscription_DOMAIN_EVENT_FEATURE},
			expectedErr: nil,
			expected: []proto.Subscription_SourceType{
				proto.Subscription_DOMAIN_EVENT_FEATURE,
				proto.Subscription_DOMAIN_EVENT_ACCOUNT,
				proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := p.origin.AddSourceTypes(p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, p.origin.SourceTypes)
		})
	}
}

func TestUpdateFeatureFlagTags(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc        string
		tags        []string
		sourceTypes []proto.Subscription_SourceType
		expected    []string
		expectedErr error
	}{
		{
			desc: "err: update subscription's feature flag tags with no feature source type",
			tags: []string{"android", "ios"},
			sourceTypes: []proto.Subscription_SourceType{
				proto.Subscription_DOMAIN_EVENT_ACCOUNT,
			},
			expected:    []string{"android", "ios", "web"},
			expectedErr: ErrCannotUpdateFeatureFlagTags,
		},
		{
			desc: "success",
			tags: []string{"android", "ios"},
			sourceTypes: []proto.Subscription_SourceType{
				proto.Subscription_DOMAIN_EVENT_ACCOUNT,
				proto.Subscription_DOMAIN_EVENT_FEATURE,
			},
			expected:    []string{"android", "ios", "web"},
			expectedErr: ErrCannotUpdateFeatureFlagTags,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := &Subscription{&proto.Subscription{
				SourceTypes:     p.sourceTypes,
				FeatureFlagTags: p.tags,
			}}
			err := actual.UpdateFeatureFlagTags(p.expected)
			if err != nil {
				assert.Equal(t, p.expectedErr, err)
				assert.Equal(t, p.tags, actual.FeatureFlagTags)
				assert.Zero(t, actual.UpdatedAt)
			} else {
				assert.Equal(t, p.expected, actual.FeatureFlagTags)
				assert.NotZero(t, actual.UpdatedAt)
			}
		})
	}
}
