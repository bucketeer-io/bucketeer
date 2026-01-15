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

	domainevent "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

func TestNewAuditLog(t *testing.T) {
	t.Parallel()
	actual := NewAuditLog(
		&domainevent.Event{
			Id:                 "id",
			Timestamp:          1,
			EntityType:         domainevent.Event_ACCOUNT,
			EntityId:           "entityId",
			Type:               domainevent.Event_FEATURE_CREATED,
			Editor:             &domainevent.Editor{Email: "email"},
			EntityData:         "entityData",
			PreviousEntityData: "previousEntityData",
			Options:            &domainevent.Options{Comment: "comment"},
		},
		"en",
	)
	assert.IsType(t, &AuditLog{}, actual)
	assert.Equal(t, "id", actual.Id)
	assert.Equal(t, int64(1), actual.Timestamp)
	assert.Equal(t, domainevent.Event_ACCOUNT, actual.EntityType)
	assert.Equal(t, "entityId", actual.EntityId)
	assert.Equal(t, domainevent.Event_FEATURE_CREATED, actual.Type)
	assert.Equal(t, &domainevent.Editor{Email: "email"}, actual.Editor)
	assert.Equal(t, "entityData", actual.EntityData)
	assert.Equal(t, "previousEntityData", actual.PreviousEntityData)
	assert.Equal(t, &domainevent.Options{Comment: "comment"}, actual.Options)
	assert.Equal(t, "en", actual.EnvironmentId)
}
