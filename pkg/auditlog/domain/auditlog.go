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
	proto "github.com/bucketeer-io/bucketeer/proto/auditlog"
	domainevent "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type AuditLog struct {
	*proto.AuditLog
	EnvironmentNamespace string
}

func NewAuditLog(event *domainevent.Event, envirronmentNamespace string) *AuditLog {
	return &AuditLog{
		AuditLog: &proto.AuditLog{
			Id:         event.Id,
			Timestamp:  event.Timestamp,
			EntityType: event.EntityType,
			EntityId:   event.EntityId,
			Type:       event.Type,
			Event:      event.Data,
			Editor:     event.Editor,
			Options:    event.Options,
		},
		EnvironmentNamespace: envirronmentNamespace,
	}
}
