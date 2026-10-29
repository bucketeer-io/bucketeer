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
	"errors"
	"strings"

	"google.golang.org/protobuf/types/known/anypb"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	domainevent "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

// secretJSONFields maps the entity types whose JSON encoded entity data holds
// a secret to the field holding it: the webhook token of a flag trigger and
// the FCM service account (which contains a private key) of a push.
var secretJSONFields = map[domainevent.Event_EntityType]string{
	domainevent.Event_FLAG_TRIGGER: "token",
	domainevent.Event_PUSH:         "fcm_service_account",
}

// RedactSecretFields blanks, in place, the secret fields of a flag trigger or
// push audit log: in the entity data, the previous entity data and the event
// data. The events no longer carry these fields; this covers the audit logs
// persisted before that, and keeps the fields out of the ones persisted from
// now on. The audit logs of the other entity types are left untouched.
func RedactSecretFields(auditlog *proto.AuditLog) error {
	if auditlog == nil {
		return nil
	}
	if _, ok := secretJSONFields[auditlog.EntityType]; !ok {
		return nil
	}
	entityData, entityErr := RedactSecretEntityData(auditlog.EntityType, auditlog.EntityData)
	auditlog.EntityData = entityData
	previousEntityData, previousErr := RedactSecretEntityData(auditlog.EntityType, auditlog.PreviousEntityData)
	auditlog.PreviousEntityData = previousEntityData
	event, eventErr := RedactSecretEventData(auditlog.Event)
	auditlog.Event = event
	return errors.Join(entityErr, previousErr, eventErr)
}

// RedactSecretEntityData blanks the secret field of the JSON encoded entity
// data of the given entity type. The data is returned as it is when the entity
// type holds no secret, or when the field is missing or already blank.
// It returns an empty string on failure, so a secret is never stored nor returned.
func RedactSecretEntityData(entityType domainevent.Event_EntityType, data string) (string, error) {
	field, ok := secretJSONFields[entityType]
	if !ok || data == "" {
		return data, nil
	}
	decoder := json.NewDecoder(strings.NewReader(data))
	// Avoid converting the numbers to scientific notation.
	decoder.UseNumber()
	entity := make(map[string]interface{})
	if err := decoder.Decode(&entity); err != nil {
		return "", err
	}
	if value, ok := entity[field].(string); !ok || value == "" {
		return data, nil
	}
	entity[field] = ""
	redacted, err := json.MarshalIndent(entity, "", "  ")
	if err != nil {
		return "", err
	}
	return string(redacted), nil
}

// RedactSecretEventData blanks the secret field of the events that carry one:
// the token of the flag trigger created and reset events, and the FCM service
// account of the push created event, which is also recorded when a push is
// deleted. The other events are returned as they are.
// It returns nil on failure, so a secret is never stored nor returned.
func RedactSecretEventData(event *anypb.Any) (*anypb.Any, error) {
	if event == nil {
		return nil, nil
	}
	if !event.MessageIs(&domainevent.FlagTriggerCreatedEvent{}) &&
		!event.MessageIs(&domainevent.FlagTriggerResetEvent{}) &&
		!event.MessageIs(&domainevent.PushCreatedEvent{}) {
		return event, nil
	}
	message, err := event.UnmarshalNew()
	if err != nil {
		return nil, err
	}
	switch m := message.(type) {
	case *domainevent.FlagTriggerCreatedEvent:
		if m.Token == "" {
			return event, nil
		}
		m.Token = ""
	case *domainevent.FlagTriggerResetEvent:
		if m.Token == "" {
			return event, nil
		}
		m.Token = ""
	case *domainevent.PushCreatedEvent:
		if m.FcmServiceAccount == "" {
			return event, nil
		}
		m.FcmServiceAccount = ""
	default:
		return event, nil
	}
	return anypb.New(message)
}
