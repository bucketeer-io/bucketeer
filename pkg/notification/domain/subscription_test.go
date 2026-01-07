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

	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

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
					UpdatedAt: time.Now().Unix(),
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
					Disabled:  true,
					UpdatedAt: time.Now().Unix(),
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
					UpdatedAt:       time.Now().Unix(),
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
					UpdatedAt:       time.Now().Unix(),
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
