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

package api

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	accountdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

const apiKeyJSONField = "api_key"

// obfuscateAPIKeys obfuscates the API keys saved before they were obfuscated at the source.
func (s *auditlogService) obfuscateAPIKeys(auditlogs []*proto.AuditLog) {
	for i := range auditlogs {
		s.obfuscateAPIKey(auditlogs[i])
	}
}

func (s *auditlogService) obfuscateAPIKey(auditlog *proto.AuditLog) {
	if auditlog == nil || auditlog.EntityType != eventproto.Event_APIKEY {
		return
	}
	auditlog.EntityData = s.obfuscateAPIKeyEntityData(auditlog.EntityData)
	auditlog.PreviousEntityData = s.obfuscateAPIKeyEntityData(auditlog.PreviousEntityData)
	auditlog.Event = s.obfuscateAPIKeyEventData(auditlog.Event)
}

// obfuscateAPIKeyEntityData drops the data when it fails, so a raw key is never returned.
func (s *auditlogService) obfuscateAPIKeyEntityData(data string) string {
	if data == "" {
		return data
	}
	decoder := json.NewDecoder(strings.NewReader(data))
	// Avoid converting the numbers to scientific notation.
	decoder.UseNumber()
	entity := make(map[string]interface{})
	if err := decoder.Decode(&entity); err != nil {
		s.logger.Error("Failed to decode the api key entity data", zap.Error(err))
		return ""
	}
	apiKey, ok := entity[apiKeyJSONField].(string)
	if !ok {
		return data
	}
	entity[apiKeyJSONField] = accountdomain.ObfuscateAPIKey(apiKey)
	obfuscated, err := json.MarshalIndent(entity, "", "  ")
	if err != nil {
		s.logger.Error("Failed to encode the obfuscated api key entity data", zap.Error(err))
		return ""
	}
	return string(obfuscated)
}

// obfuscateAPIKeyEventData handles the created event, the only one containing the key.
// It drops the event when it fails, so a raw key is never returned.
func (s *auditlogService) obfuscateAPIKeyEventData(event *anypb.Any) *anypb.Any {
	if event == nil || !event.MessageIs(&eventproto.APIKeyCreatedEvent{}) {
		return event
	}
	created := &eventproto.APIKeyCreatedEvent{}
	if err := event.UnmarshalTo(created); err != nil {
		s.logger.Error("Failed to unmarshal the api key created event", zap.Error(err))
		return nil
	}
	created.ApiKey = accountdomain.ObfuscateAPIKey(created.ApiKey)
	obfuscated, err := anypb.New(created)
	if err != nil {
		s.logger.Error("Failed to marshal the obfuscated api key created event", zap.Error(err))
		return nil
	}
	return obfuscated
}
