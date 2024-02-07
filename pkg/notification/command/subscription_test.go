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

package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

func TestCreate(t *testing.T) {
	patterns := []*struct {
		input    *proto.CreateSubscriptionCommand
		expected error
	}{
		{
			input: &proto.CreateSubscriptionCommand{
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
				},
				Recipient: &proto.Recipient{
					Type:                  proto.Recipient_SlackChannel,
					SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
				},
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, false)
		h := newSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestDelete(t *testing.T) {
	patterns := []*struct {
		input    *proto.DeleteSubscriptionCommand
		expected error
	}{
		{
			input:    &proto.DeleteSubscriptionCommand{},
			expected: nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, false)
		h := newSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestAddSourceTypes(t *testing.T) {
	patterns := []*struct {
		input    *proto.AddSourceTypesCommand
		expected error
	}{
		{
			input: &proto.AddSourceTypesCommand{
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_FEATURE,
				},
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, false)
		h := newSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestDeleteSourceTypes(t *testing.T) {
	patterns := []*struct {
		input    *proto.DeleteSourceTypesCommand
		expected error
	}{
		{
			input: &proto.DeleteSourceTypesCommand{
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
				},
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, false)
		h := newSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestEnable(t *testing.T) {
	patterns := []*struct {
		originDisabled bool
		input          *proto.EnableSubscriptionCommand
		expected       error
	}{
		{
			originDisabled: true,
			input:          &proto.EnableSubscriptionCommand{},
			expected:       nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, p.originDisabled)
		h := newSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestDisable(t *testing.T) {
	patterns := []*struct {
		originDisabled bool
		input          *proto.DisableSubscriptionCommand
		expected       error
	}{
		{
			originDisabled: false,
			input:          &proto.DisableSubscriptionCommand{},
			expected:       nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, p.originDisabled)
		h := newSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestRename(t *testing.T) {
	patterns := []*struct {
		input    *proto.RenameSubscriptionCommand
		expected error
	}{
		{
			input:    &proto.RenameSubscriptionCommand{Name: "renamed"},
			expected: nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, false)
		h := newSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func newSubscription(t *testing.T, disabled bool) *domain.Subscription {
	sourceTypes := []proto.Subscription_SourceType{
		proto.Subscription_DOMAIN_EVENT_ACCOUNT,
		proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
	}
	recipient := &proto.Recipient{
		Type:                  proto.Recipient_SlackChannel,
		SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
	}
	s, err := domain.NewSubscription("sname", sourceTypes, recipient)
	s.Disabled = disabled
	require.NoError(t, err)
	return s
}

func newSubscriptionCommandHandler(t *testing.T, subscription *domain.Subscription) Handler {
	t.Helper()
	return NewSubscriptionCommandHandler(
		&eventproto.Editor{
			Email: "email",
		},
		subscription,
		"ns0",
	)
}
