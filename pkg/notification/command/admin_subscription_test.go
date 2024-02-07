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

	"github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

func TestAdminCreate(t *testing.T) {
	patterns := []*struct {
		input    *proto.CreateAdminSubscriptionCommand
		expected error
	}{
		{
			input: &proto.CreateAdminSubscriptionCommand{
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
		h := newAdminSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestAdminDelete(t *testing.T) {
	patterns := []*struct {
		input    *proto.DeleteAdminSubscriptionCommand
		expected error
	}{
		{
			input:    &proto.DeleteAdminSubscriptionCommand{},
			expected: nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, false)
		h := newAdminSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestAdminAddSourceTypes(t *testing.T) {
	patterns := []*struct {
		input    *proto.AddAdminSubscriptionSourceTypesCommand
		expected error
	}{
		{
			input: &proto.AddAdminSubscriptionSourceTypesCommand{
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_FEATURE,
				},
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, false)
		h := newAdminSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestAdminDeleteSourceTypes(t *testing.T) {
	patterns := []*struct {
		input    *proto.DeleteAdminSubscriptionSourceTypesCommand
		expected error
	}{
		{
			input: &proto.DeleteAdminSubscriptionSourceTypesCommand{
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
				},
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, false)
		h := newAdminSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestAdminEnable(t *testing.T) {
	patterns := []*struct {
		originDisabled bool
		input          *proto.EnableAdminSubscriptionCommand
		expected       error
	}{
		{
			originDisabled: true,
			input:          &proto.EnableAdminSubscriptionCommand{},
			expected:       nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, p.originDisabled)
		h := newAdminSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestAdminDisable(t *testing.T) {
	patterns := []*struct {
		originDisabled bool
		input          *proto.DisableAdminSubscriptionCommand
		expected       error
	}{
		{
			originDisabled: false,
			input:          &proto.DisableAdminSubscriptionCommand{},
			expected:       nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, p.originDisabled)
		h := newAdminSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func TestAdminRename(t *testing.T) {
	patterns := []*struct {
		input    *proto.RenameAdminSubscriptionCommand
		expected error
	}{
		{
			input:    &proto.RenameAdminSubscriptionCommand{Name: "renamed"},
			expected: nil,
		},
	}
	for _, p := range patterns {
		s := newSubscription(t, false)
		h := newAdminSubscriptionCommandHandler(t, s)
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expected, err)
		assert.Equal(t, 1, len(h.Events()))
	}
}

func newAdminSubscriptionCommandHandler(t *testing.T, subscription *domain.Subscription) Handler {
	t.Helper()
	return NewAdminSubscriptionCommandHandler(
		&eventproto.Editor{
			Email:   "email",
			IsAdmin: true,
		},
		subscription,
	)
}
