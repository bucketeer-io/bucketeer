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
//

package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

const (
	rawTriggerToken = "8xK2mPqR5tV7wY9zA1bC3dE4fG6hJ8kL0mN2oP4qR6sT"
	// A service account key file as stored in the push, i.e. a JSON document.
	rawFCMServiceAccount = `{"type":"service_account","private_key":"-----BEGIN PRIVATE KEY-----\nMIIB\n-----END PRIVATE KEY-----\n"}`
	privateKeyMarker     = "PRIVATE KEY"
)

// jsonQuote returns the JSON string literal of s, as it appears in the entity data.
func jsonQuote(t *testing.T, s string) string {
	t.Helper()
	quoted, err := json.Marshal(s)
	require.NoError(t, err)
	return string(quoted)
}

func TestRedactSecretEntityData(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		entityType  eventproto.Event_EntityType
		input       string
		expected    string
		expectedErr bool
	}{
		{
			desc:       "empty",
			entityType: eventproto.Event_FLAG_TRIGGER,
			input:      "",
			expected:   "",
		},
		{
			desc:        "invalid json: dropped",
			entityType:  eventproto.Event_FLAG_TRIGGER,
			input:       "not a json",
			expected:    "",
			expectedErr: true,
		},
		{
			desc:       "entity type without secret: returned as it is even with a token field",
			entityType: eventproto.Event_FEATURE,
			input:      "{\n  \"id\": \"id-1\",\n  \"token\": \"" + rawTriggerToken + "\"\n}",
			expected:   "{\n  \"id\": \"id-1\",\n  \"token\": \"" + rawTriggerToken + "\"\n}",
		},
		{
			desc:       "flag trigger without token field: returned as it is",
			entityType: eventproto.Event_FLAG_TRIGGER,
			input:      "{\n  \"id\": \"id-1\"\n}",
			expected:   "{\n  \"id\": \"id-1\"\n}",
		},
		{
			desc:       "flag trigger with blank token: returned as it is, not re-encoded",
			entityType: eventproto.Event_FLAG_TRIGGER,
			input:      "{\n  \"token\": \"\",\n  \"id\": \"id-1\"\n}",
			expected:   "{\n  \"token\": \"\",\n  \"id\": \"id-1\"\n}",
		},
		{
			desc:       "flag trigger token blanked and numbers kept as they are",
			entityType: eventproto.Event_FLAG_TRIGGER,
			input:      "{\n  \"id\": \"id-1\",\n  \"token\": \"" + rawTriggerToken + "\",\n  \"trigger_count\": 1739746800\n}",
			expected:   "{\n  \"id\": \"id-1\",\n  \"token\": \"\",\n  \"trigger_count\": 1739746800\n}",
		},
		{
			desc:       "push fcm service account blanked",
			entityType: eventproto.Event_PUSH,
			input:      "{\n  \"fcm_service_account\": " + jsonQuote(t, rawFCMServiceAccount) + ",\n  \"id\": \"id-1\",\n  \"name\": \"push-1\"\n}",
			expected:   "{\n  \"fcm_service_account\": \"\",\n  \"id\": \"id-1\",\n  \"name\": \"push-1\"\n}",
		},
		{
			desc:       "push token field is not a secret: returned as it is",
			entityType: eventproto.Event_PUSH,
			input:      "{\n  \"id\": \"id-1\",\n  \"token\": \"not-a-trigger-token\"\n}",
			expected:   "{\n  \"id\": \"id-1\",\n  \"token\": \"not-a-trigger-token\"\n}",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := RedactSecretEntityData(p.entityType, p.input)
			assert.Equal(t, p.expected, actual)
			if p.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestRedactSecretEventData(t *testing.T) {
	t.Parallel()
	triggerCreated, err := anypb.New(&eventproto.FlagTriggerCreatedEvent{
		Id:          "trigger-1",
		FeatureId:   "feature-1",
		Description: "desc",
		Token:       rawTriggerToken,
	})
	require.NoError(t, err)
	triggerReset, err := anypb.New(&eventproto.FlagTriggerResetEvent{
		Id:        "trigger-1",
		FeatureId: "feature-1",
		Token:     rawTriggerToken,
	})
	require.NoError(t, err)
	triggerCreatedBlank, err := anypb.New(&eventproto.FlagTriggerCreatedEvent{Id: "trigger-1"})
	require.NoError(t, err)
	triggerUpdated, err := anypb.New(&eventproto.FlagTriggerUpdateEvent{Id: "trigger-1"})
	require.NoError(t, err)
	pushCreated, err := anypb.New(&eventproto.PushCreatedEvent{
		FcmServiceAccount: rawFCMServiceAccount,
		Tags:              []string{"tag-1"},
		Name:              "push-1",
	})
	require.NoError(t, err)
	pushCreatedBlank, err := anypb.New(&eventproto.PushCreatedEvent{Name: "push-1"})
	require.NoError(t, err)
	pushUpdated, err := anypb.New(&eventproto.PushUpdatedEvent{Tags: []string{"tag-1"}})
	require.NoError(t, err)

	patterns := []struct {
		desc     string
		input    *anypb.Any
		expected *anypb.Any
		// check inspects the redacted event when it is re-encoded.
		check func(*testing.T, *anypb.Any)
	}{
		{
			desc:     "nil",
			input:    nil,
			expected: nil,
		},
		{
			desc:     "flag trigger updated: same event",
			input:    triggerUpdated,
			expected: triggerUpdated,
		},
		{
			desc:     "push updated: same event",
			input:    pushUpdated,
			expected: pushUpdated,
		},
		{
			desc:     "flag trigger created with blank token: same event",
			input:    triggerCreatedBlank,
			expected: triggerCreatedBlank,
		},
		{
			desc:     "push created with blank fcm service account: same event",
			input:    pushCreatedBlank,
			expected: pushCreatedBlank,
		},
		{
			desc:  "flag trigger created: token blanked, other fields kept",
			input: triggerCreated,
			check: func(t *testing.T, actual *anypb.Any) {
				created := &eventproto.FlagTriggerCreatedEvent{}
				require.NoError(t, actual.UnmarshalTo(created))
				assert.Equal(t, "", created.Token)
				assert.Equal(t, "trigger-1", created.Id)
				assert.Equal(t, "feature-1", created.FeatureId)
				assert.Equal(t, "desc", created.Description)
			},
		},
		{
			desc:  "flag trigger reset: token blanked, other fields kept",
			input: triggerReset,
			check: func(t *testing.T, actual *anypb.Any) {
				reset := &eventproto.FlagTriggerResetEvent{}
				require.NoError(t, actual.UnmarshalTo(reset))
				assert.Equal(t, "", reset.Token)
				assert.Equal(t, "trigger-1", reset.Id)
				assert.Equal(t, "feature-1", reset.FeatureId)
			},
		},
		{
			desc:  "push created: fcm service account blanked, other fields kept",
			input: pushCreated,
			check: func(t *testing.T, actual *anypb.Any) {
				created := &eventproto.PushCreatedEvent{}
				require.NoError(t, actual.UnmarshalTo(created))
				assert.Equal(t, "", created.FcmServiceAccount)
				assert.Equal(t, []string{"tag-1"}, created.Tags)
				assert.Equal(t, "push-1", created.Name)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := RedactSecretEventData(p.input)
			require.NoError(t, err)
			if p.check != nil {
				require.NotNil(t, actual)
				assert.NotContains(t, string(actual.Value), rawTriggerToken)
				assert.NotContains(t, string(actual.Value), privateKeyMarker)
				p.check(t, actual)
				return
			}
			assert.Same(t, p.expected, actual)
		})
	}
}

func TestRedactSecretFields(t *testing.T) {
	t.Parallel()
	triggerCreated, err := anypb.New(&eventproto.FlagTriggerCreatedEvent{Id: "trigger-1", Token: rawTriggerToken})
	require.NoError(t, err)
	triggerUpdated, err := anypb.New(&eventproto.FlagTriggerUpdateEvent{Id: "trigger-1"})
	require.NoError(t, err)
	pushCreated, err := anypb.New(&eventproto.PushCreatedEvent{Name: "push-1", FcmServiceAccount: rawFCMServiceAccount})
	require.NoError(t, err)
	featureCreated, err := anypb.New(&eventproto.FeatureCreatedEvent{Id: "feature-1"})
	require.NoError(t, err)
	// The fields are sorted because the entity data is re-encoded from a map.
	triggerData := func(token string) string {
		return "{\n  \"id\": \"trigger-1\",\n  \"token\": " + jsonQuote(t, token) + "\n}"
	}
	pushData := func(fcmServiceAccount string) string {
		return "{\n  \"fcm_service_account\": " + jsonQuote(t, fcmServiceAccount) + ",\n  \"id\": \"push-1\"\n}"
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
			desc: "not a flag trigger nor a push: untouched",
			input: &proto.AuditLog{
				EntityType:         eventproto.Event_FEATURE,
				Type:               eventproto.Event_FEATURE_CREATED,
				Event:              featureCreated,
				EntityData:         triggerData(rawTriggerToken),
				PreviousEntityData: pushData(rawFCMServiceAccount),
			},
			expected: &proto.AuditLog{
				EntityType:         eventproto.Event_FEATURE,
				Type:               eventproto.Event_FEATURE_CREATED,
				Event:              featureCreated,
				EntityData:         triggerData(rawTriggerToken),
				PreviousEntityData: pushData(rawFCMServiceAccount),
			},
		},
		{
			desc: "flag trigger created: entity data and event data blanked",
			input: &proto.AuditLog{
				EntityType: eventproto.Event_FLAG_TRIGGER,
				Type:       eventproto.Event_FLAG_TRIGGER_CREATED,
				Event:      triggerCreated,
				EntityData: triggerData(rawTriggerToken),
			},
			expected: &proto.AuditLog{
				EntityType: eventproto.Event_FLAG_TRIGGER,
				Type:       eventproto.Event_FLAG_TRIGGER_CREATED,
				EntityData: triggerData(""),
			},
		},
		{
			desc: "flag trigger updated: both entity data blanked, event data kept",
			input: &proto.AuditLog{
				EntityType:         eventproto.Event_FLAG_TRIGGER,
				Type:               eventproto.Event_FLAG_TRIGGER_UPDATED,
				Event:              triggerUpdated,
				EntityData:         triggerData(rawTriggerToken),
				PreviousEntityData: triggerData(rawTriggerToken),
			},
			expected: &proto.AuditLog{
				EntityType:         eventproto.Event_FLAG_TRIGGER,
				Type:               eventproto.Event_FLAG_TRIGGER_UPDATED,
				Event:              triggerUpdated,
				EntityData:         triggerData(""),
				PreviousEntityData: triggerData(""),
			},
		},
		{
			desc: "push deleted: previous entity data and event data blanked",
			input: &proto.AuditLog{
				EntityType:         eventproto.Event_PUSH,
				Type:               eventproto.Event_PUSH_DELETED,
				Event:              pushCreated,
				PreviousEntityData: pushData(rawFCMServiceAccount),
			},
			expected: &proto.AuditLog{
				EntityType:         eventproto.Event_PUSH,
				Type:               eventproto.Event_PUSH_DELETED,
				PreviousEntityData: pushData(""),
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			require.NoError(t, RedactSecretFields(p.input))
			if p.input == nil {
				return
			}
			assert.Equal(t, p.expected.EntityData, p.input.EntityData)
			assert.Equal(t, p.expected.PreviousEntityData, p.input.PreviousEntityData)
			switch p.input.Type {
			case eventproto.Event_FLAG_TRIGGER_CREATED:
				created := &eventproto.FlagTriggerCreatedEvent{}
				require.NoError(t, p.input.Event.UnmarshalTo(created))
				assert.Equal(t, "", created.Token)
				assert.Equal(t, "trigger-1", created.Id)
			case eventproto.Event_PUSH_DELETED:
				created := &eventproto.PushCreatedEvent{}
				require.NoError(t, p.input.Event.UnmarshalTo(created))
				assert.Equal(t, "", created.FcmServiceAccount)
				assert.Equal(t, "push-1", created.Name)
			default:
				assert.Equal(t, p.expected.Event, p.input.Event)
			}
		})
	}
}

func TestNewAuditLogRedactsSecretFields(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc       string
		entityType eventproto.Event_EntityType
		eventType  eventproto.Event_Type
		event      func(*testing.T) *anypb.Any
		entityData string
		// secret must not appear anywhere in the audit log.
		secret string
		// kept is a non-secret value of the entity data that must survive.
		kept       string
		checkEvent func(*testing.T, *anypb.Any)
	}{
		{
			desc:       "flag trigger created",
			entityType: eventproto.Event_FLAG_TRIGGER,
			eventType:  eventproto.Event_FLAG_TRIGGER_CREATED,
			event: func(t *testing.T) *anypb.Any {
				event, err := anypb.New(&eventproto.FlagTriggerCreatedEvent{Id: "trigger-1", Token: rawTriggerToken})
				require.NoError(t, err)
				return event
			},
			entityData: "{\n  \"id\": \"trigger-1\",\n  \"token\": \"" + rawTriggerToken + "\"\n}",
			secret:     rawTriggerToken,
			kept:       "\"trigger-1\"",
			checkEvent: func(t *testing.T, event *anypb.Any) {
				created := &eventproto.FlagTriggerCreatedEvent{}
				require.NoError(t, event.UnmarshalTo(created))
				assert.Equal(t, "", created.Token)
				assert.Equal(t, "trigger-1", created.Id)
			},
		},
		{
			desc:       "push created",
			entityType: eventproto.Event_PUSH,
			eventType:  eventproto.Event_PUSH_CREATED,
			event: func(t *testing.T) *anypb.Any {
				event, err := anypb.New(&eventproto.PushCreatedEvent{Name: "push-1", FcmServiceAccount: rawFCMServiceAccount})
				require.NoError(t, err)
				return event
			},
			entityData: "{\n  \"fcm_service_account\": " + jsonQuote(t, rawFCMServiceAccount) + ",\n  \"id\": \"push-1\"\n}",
			secret:     privateKeyMarker,
			kept:       "\"push-1\"",
			checkEvent: func(t *testing.T, event *anypb.Any) {
				created := &eventproto.PushCreatedEvent{}
				require.NoError(t, event.UnmarshalTo(created))
				assert.Equal(t, "", created.FcmServiceAccount)
				assert.Equal(t, "push-1", created.Name)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := NewAuditLog(&eventproto.Event{
				Id:                 "auditlog-1",
				EntityType:         p.entityType,
				EntityId:           "id-1",
				Type:               p.eventType,
				Editor:             &eventproto.Editor{Email: "bucketeer@bucketeer.io"},
				Data:               p.event(t),
				EntityData:         p.entityData,
				PreviousEntityData: p.entityData,
			}, "env-1")

			assert.NotContains(t, actual.EntityData, p.secret)
			assert.Contains(t, actual.EntityData, p.kept)
			assert.NotContains(t, actual.PreviousEntityData, p.secret)
			assert.Contains(t, actual.PreviousEntityData, p.kept)
			assert.NotContains(t, string(actual.Event.Value), p.secret)
			p.checkEvent(t, actual.Event)
		})
	}
}
