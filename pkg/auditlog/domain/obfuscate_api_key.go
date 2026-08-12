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
	"encoding/json"
	"errors"
	"strings"

	"google.golang.org/protobuf/types/known/anypb"

	accountdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	domainevent "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

const apiKeyJSONField = "api_key"

// ObfuscateAPIKey obfuscates the API key of an API key audit log, in place.
// The audit log is where the key would be persisted, so it is masked here instead of in the domain
// event, whose entity data is where the API key cache pipeline reads the raw key from.
func ObfuscateAPIKey(auditlog *proto.AuditLog) error {
	if auditlog == nil || auditlog.EntityType != domainevent.Event_APIKEY {
		return nil
	}
	entityData, entityErr := ObfuscateAPIKeyEntityData(auditlog.EntityData)
	auditlog.EntityData = entityData
	previousEntityData, previousErr := ObfuscateAPIKeyEntityData(auditlog.PreviousEntityData)
	auditlog.PreviousEntityData = previousEntityData
	event, eventErr := ObfuscateAPIKeyEventData(auditlog.Event)
	auditlog.Event = event
	return errors.Join(entityErr, previousErr, eventErr)
}

// ObfuscateAPIKeyEntityData obfuscates the API key of the JSON encoded API key entity data.
// It returns an empty string on failure, so a raw key is never stored nor returned.
func ObfuscateAPIKeyEntityData(data string) (string, error) {
	if data == "" {
		return data, nil
	}
	decoder := json.NewDecoder(strings.NewReader(data))
	// Avoid converting the numbers to scientific notation.
	decoder.UseNumber()
	entity := make(map[string]interface{})
	if err := decoder.Decode(&entity); err != nil {
		return "", err
	}
	apiKey, ok := entity[apiKeyJSONField].(string)
	if !ok {
		return data, nil
	}
	entity[apiKeyJSONField] = accountdomain.ObfuscateAPIKey(apiKey)
	obfuscated, err := json.MarshalIndent(entity, "", "  ")
	if err != nil {
		return "", err
	}
	return string(obfuscated), nil
}

// ObfuscateAPIKeyEventData handles the created event, the only one containing the key.
// It returns nil on failure, so a raw key is never stored nor returned.
func ObfuscateAPIKeyEventData(event *anypb.Any) (*anypb.Any, error) {
	if event == nil || !event.MessageIs(&domainevent.APIKeyCreatedEvent{}) {
		return event, nil
	}
	created := &domainevent.APIKeyCreatedEvent{}
	if err := event.UnmarshalTo(created); err != nil {
		return nil, err
	}
	created.ApiKey = accountdomain.ObfuscateAPIKey(created.ApiKey)
	obfuscated, err := anypb.New(created)
	if err != nil {
		return nil, err
	}
	return obfuscated, nil
}
