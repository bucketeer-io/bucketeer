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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

const (
	rawAPIKey        = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	obfuscatedAPIKey = "0123....cdef"
)

func TestObfuscateAPIKeyEntityData(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		input       string
		expected    string
		expectedErr bool
	}{
		{
			desc:     "empty",
			input:    "",
			expected: "",
		},
		{
			desc:        "invalid json: dropped",
			input:       "not a json",
			expected:    "",
			expectedErr: true,
		},
		{
			desc:     "no api key field: returned as it is",
			input:    "{\n  \"id\": \"id-1\"\n}",
			expected: "{\n  \"id\": \"id-1\"\n}",
		},
		{
			desc:     "api key obfuscated and numbers kept as they are",
			input:    "{\n  \"api_key\": \"" + rawAPIKey + "\",\n  \"created_at\": 1739746800\n}",
			expected: "{\n  \"api_key\": \"" + obfuscatedAPIKey + "\",\n  \"created_at\": 1739746800\n}",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := ObfuscateAPIKeyEntityData(p.input)
			assert.Equal(t, p.expected, actual)
			if p.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestObfuscateAPIKey(t *testing.T) {
	t.Parallel()
	createdEvent, err := anypb.New(&eventproto.APIKeyCreatedEvent{
		Id:     "id-1",
		Name:   "name-1",
		ApiKey: rawAPIKey,
	})
	require.NoError(t, err)
	changedEvent, err := anypb.New(&eventproto.APIKeyChangedEvent{Id: "id-1"})
	require.NoError(t, err)
	featureEvent, err := anypb.New(&eventproto.FeatureCreatedEvent{Id: "feature-1"})
	require.NoError(t, err)
	// The fields are sorted because the entity data is re-encoded from a map.
	entityData := func(apiKey string) string {
		return "{\n  \"api_key\": \"" + apiKey + "\",\n  \"id\": \"id-1\"\n}"
	}

	patterns := []struct {
		desc     string
		input    *proto.AuditLog
		expected *proto.AuditLog
	}{
		{
			desc:     "nil",
			input:    nil,
			expected: nil,
		},
		{
			desc: "not an api key entity: untouched",
			input: &proto.AuditLog{
				EntityType: eventproto.Event_FEATURE,
				Type:       eventproto.Event_FEATURE_CREATED,
				Event:      featureEvent,
				EntityData: entityData(rawAPIKey),
			},
			expected: &proto.AuditLog{
				EntityType: eventproto.Event_FEATURE,
				Type:       eventproto.Event_FEATURE_CREATED,
				Event:      featureEvent,
				EntityData: entityData(rawAPIKey),
			},
		},
		{
			desc: "api key created: entity data and event data obfuscated",
			input: &proto.AuditLog{
				EntityType: eventproto.Event_APIKEY,
				Type:       eventproto.Event_APIKEY_CREATED,
				Event:      createdEvent,
				EntityData: entityData(rawAPIKey),
			},
			expected: &proto.AuditLog{
				EntityType: eventproto.Event_APIKEY,
				Type:       eventproto.Event_APIKEY_CREATED,
				EntityData: entityData(obfuscatedAPIKey),
			},
		},
		{
			desc: "api key changed: both entity data obfuscated",
			input: &proto.AuditLog{
				EntityType:         eventproto.Event_APIKEY,
				Type:               eventproto.Event_APIKEY_CHANGED,
				Event:              changedEvent,
				EntityData:         entityData(rawAPIKey),
				PreviousEntityData: entityData(rawAPIKey),
			},
			expected: &proto.AuditLog{
				EntityType:         eventproto.Event_APIKEY,
				Type:               eventproto.Event_APIKEY_CHANGED,
				Event:              changedEvent,
				EntityData:         entityData(obfuscatedAPIKey),
				PreviousEntityData: entityData(obfuscatedAPIKey),
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			require.NoError(t, ObfuscateAPIKey(p.input))
			if p.input == nil {
				return
			}
			assert.Equal(t, p.expected.EntityData, p.input.EntityData)
			assert.Equal(t, p.expected.PreviousEntityData, p.input.PreviousEntityData)
			if p.input.Type == eventproto.Event_APIKEY_CREATED {
				created := &eventproto.APIKeyCreatedEvent{}
				require.NoError(t, p.input.Event.UnmarshalTo(created))
				assert.Equal(t, obfuscatedAPIKey, created.ApiKey)
				assert.Equal(t, "name-1", created.Name)
				return
			}
			assert.Equal(t, p.expected.Event, p.input.Event)
		})
	}
}

func TestObfuscateAPIKeyOnlyTouchesAPIKeyLogs(t *testing.T) {
	t.Parallel()
	auditlogs := []*proto.AuditLog{
		{
			EntityType: eventproto.Event_APIKEY,
			Type:       eventproto.Event_APIKEY_CHANGED,
			EntityData: "{\n  \"api_key\": \"" + rawAPIKey + "\"\n}",
		},
		{
			EntityType: eventproto.Event_ACCOUNT,
			Type:       eventproto.Event_ACCOUNT_V2_CREATED,
			EntityData: "{\n  \"email\": \"bucketeer@bucketeer.io\"\n}",
		},
	}
	for _, a := range auditlogs {
		require.NoError(t, ObfuscateAPIKey(a))
	}
	assert.Equal(t, "{\n  \"api_key\": \""+obfuscatedAPIKey+"\"\n}", auditlogs[0].EntityData)
	assert.Equal(t, "{\n  \"email\": \"bucketeer@bucketeer.io\"\n}", auditlogs[1].EntityData)
}

func TestNewAuditLogObfuscatesAPIKey(t *testing.T) {
	t.Parallel()
	createdEvent, err := anypb.New(&eventproto.APIKeyCreatedEvent{Id: "id-1", ApiKey: rawAPIKey})
	require.NoError(t, err)
	entityData := "{\n  \"api_key\": \"" + rawAPIKey + "\",\n  \"id\": \"id-1\"\n}"

	actual := NewAuditLog(&eventproto.Event{
		Id:                 "auditlog-1",
		EntityType:         eventproto.Event_APIKEY,
		EntityId:           "id-1",
		Type:               eventproto.Event_APIKEY_CREATED,
		Editor:             &eventproto.Editor{Email: "bucketeer@bucketeer.io"},
		Data:               createdEvent,
		EntityData:         entityData,
		PreviousEntityData: entityData,
	}, "env-1")

	assert.NotContains(t, actual.EntityData, rawAPIKey)
	assert.Contains(t, actual.EntityData, obfuscatedAPIKey)
	assert.NotContains(t, actual.PreviousEntityData, rawAPIKey)
	assert.Contains(t, actual.PreviousEntityData, obfuscatedAPIKey)
	created := &eventproto.APIKeyCreatedEvent{}
	require.NoError(t, actual.Event.UnmarshalTo(created))
	assert.Equal(t, obfuscatedAPIKey, created.ApiKey)
}
