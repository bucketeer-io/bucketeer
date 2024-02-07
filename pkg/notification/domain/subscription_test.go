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

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/bucketeer-io/bucketeer/proto/notification"
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
	name := "sname"
	actual, err := NewSubscription(name, sourceTypes, recipient)
	assert.NoError(t, err)
	assert.IsType(t, &Subscription{}, actual)
	assert.NotEqual(t, "", actual.Id)
	assert.Equal(t, sourceTypes, actual.SourceTypes)
	assert.Equal(t, recipient, actual.Recipient)
	assert.Equal(t, false, actual.Disabled)
	assert.NotEqual(t, 0, actual.CreatedAt)
	assert.NotEqual(t, 0, actual.UpdatedAt)
	assert.Equal(t, name, actual.Name)
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
