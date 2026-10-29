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

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"

	v2asmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	v2alsmock "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

const (
	redactTestEnvironmentID = "ns0"
	redactTestTriggerToken  = "8xK2mPqR5tV7wY9zA1bC3dE4fG6hJ8kL0mN2oP4qR6sT"
	redactTestPrivateKey    = "PRIVATE KEY"
)

// storedAuditLogs returns audit logs as persisted before the flag trigger
// token and the push FCM service account were kept out of them, plus a
// feature audit log that must be returned as it is.
func storedAuditLogs(t *testing.T) []*proto.AuditLog {
	t.Helper()
	triggerCreated, err := anypb.New(&domaineventproto.FlagTriggerCreatedEvent{
		Id:    "trigger-1",
		Token: redactTestTriggerToken,
	})
	require.NoError(t, err)
	pushCreated, err := anypb.New(&domaineventproto.PushCreatedEvent{
		Name:              "push-1",
		FcmServiceAccount: `{"type":"service_account","private_key":"-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----\n"}`,
	})
	require.NoError(t, err)
	featureCreated, err := anypb.New(&domaineventproto.FeatureCreatedEvent{Id: "feature-1"})
	require.NoError(t, err)
	editor := func() *domaineventproto.Editor {
		return &domaineventproto.Editor{Email: "editor@example.com"}
	}
	return []*proto.AuditLog{
		{
			Id:         "trigger-created",
			EntityType: domaineventproto.Event_FLAG_TRIGGER,
			Type:       domaineventproto.Event_FLAG_TRIGGER_CREATED,
			Event:      triggerCreated,
			Editor:     editor(),
			EntityData: "{\n  \"id\": \"trigger-1\",\n  \"token\": \"" + redactTestTriggerToken + "\"\n}",
		},
		{
			Id:                 "trigger-updated",
			EntityType:         domaineventproto.Event_FLAG_TRIGGER,
			Type:               domaineventproto.Event_FLAG_TRIGGER_UPDATED,
			Editor:             editor(),
			EntityData:         "{\n  \"id\": \"trigger-1\",\n  \"token\": \"" + redactTestTriggerToken + "\"\n}",
			PreviousEntityData: "{\n  \"id\": \"trigger-1\",\n  \"token\": \"" + redactTestTriggerToken + "\"\n}",
		},
		{
			Id:                 "push-deleted",
			EntityType:         domaineventproto.Event_PUSH,
			Type:               domaineventproto.Event_PUSH_DELETED,
			Event:              pushCreated,
			Editor:             editor(),
			PreviousEntityData: "{\n  \"fcm_service_account\": \"{\\\"private_key\\\":\\\"-----BEGIN PRIVATE KEY-----\\\"}\",\n  \"id\": \"push-1\"\n}",
		},
		{
			Id:         "feature-created",
			EntityType: domaineventproto.Event_FEATURE,
			Type:       domaineventproto.Event_FEATURE_CREATED,
			Event:      featureCreated,
			Editor:     editor(),
			EntityData: "{\n  \"id\": \"feature-1\",\n  \"name\": \"token " + redactTestTriggerToken + "\"\n}",
		},
	}
}

// assertAuditLogRedacted checks the audit log returned for the stored one.
func assertAuditLogRedacted(t *testing.T, actual *proto.AuditLog) {
	t.Helper()
	switch actual.Id {
	case "feature-created":
		assert.Contains(t, actual.EntityData, redactTestTriggerToken, "other entity types are returned as they are")
		return
	case "trigger-created":
		created := &domaineventproto.FlagTriggerCreatedEvent{}
		require.NoError(t, actual.Event.UnmarshalTo(created))
		assert.Equal(t, "", created.Token)
		assert.Equal(t, "trigger-1", created.Id)
	case "push-deleted":
		created := &domaineventproto.PushCreatedEvent{}
		require.NoError(t, actual.Event.UnmarshalTo(created))
		assert.Equal(t, "", created.FcmServiceAccount)
		assert.Equal(t, "push-1", created.Name)
	}
	assert.NotContains(t, actual.EntityData, redactTestTriggerToken)
	assert.NotContains(t, actual.EntityData, redactTestPrivateKey)
	assert.NotContains(t, actual.PreviousEntityData, redactTestTriggerToken)
	assert.NotContains(t, actual.PreviousEntityData, redactTestPrivateKey)
	if actual.EntityData != "" {
		assert.Contains(t, actual.EntityData, "\"id\": \"trigger-1\"", "the other fields are kept")
	}
	if actual.PreviousEntityData != "" {
		assert.Contains(t, actual.PreviousEntityData, "\"id\": \"", "the other fields are kept")
	}
	assert.NotEmpty(t, actual.LocalizedMessage)
}

func TestListAuditLogsRedactsSecretFields(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	s := newAuditLogServiceWithGetAccountByEnvironmentMock(
		t, mockController,
		accountproto.AccountV2_Role_Organization_MEMBER,
		accountproto.AccountV2_Role_Environment_VIEWER,
	)
	stored := storedAuditLogs(t)
	s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
		gomock.Any(), gomock.Any(),
	).Return(stored, len(stored), int64(len(stored)), nil)
	s.accountStorage.(*v2asmock.MockAccountStorage).EXPECT().GetAvatarAccountsV2(
		gomock.Any(), gomock.Any(),
	).Return([]*accountproto.AccountV2{}, nil)
	ctx := metadata.NewIncomingContext(createContextWithToken(t, false), metadata.MD{
		"accept-language": []string{"ja"},
	})

	actual, err := s.ListAuditLogs(ctx, &proto.ListAuditLogsRequest{
		EnvironmentId: redactTestEnvironmentID,
	})
	require.NoError(t, err)
	require.Len(t, actual.AuditLogs, len(stored))
	for _, auditlog := range actual.AuditLogs {
		assertAuditLogRedacted(t, auditlog)
	}
}

func TestGetAuditLogRedactsSecretFields(t *testing.T) {
	t.Parallel()
	for _, stored := range storedAuditLogs(t) {
		t.Run(stored.Id, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			s := newAuditLogServiceWithGetAccountByEnvironmentMock(
				t, mockController,
				accountproto.AccountV2_Role_Organization_MEMBER,
				accountproto.AccountV2_Role_Environment_VIEWER,
			)
			s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().GetAuditLog(
				gomock.Any(), stored.Id, redactTestEnvironmentID, "",
			).Return(stored, nil)
			s.accountStorage.(*v2asmock.MockAccountStorage).EXPECT().GetAvatarAccountsV2(
				gomock.Any(), gomock.Any(),
			).Return([]*accountproto.AccountV2{}, nil)
			ctx := metadata.NewIncomingContext(createContextWithToken(t, false), metadata.MD{
				"accept-language": []string{"ja"},
			})

			actual, err := s.GetAuditLog(ctx, &proto.GetAuditLogRequest{
				Id:            stored.Id,
				EnvironmentId: redactTestEnvironmentID,
			})
			require.NoError(t, err)
			assertAuditLogRedacted(t, actual.AuditLog)
		})
	}
}
