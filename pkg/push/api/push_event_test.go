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
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	pb "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/push/domain"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2/mock"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	pushproto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

func TestPushEventsOmitFCMServiceAccount(t *testing.T) {
	t.Parallel()

	const (
		environmentID = "ns0"
		pushID        = "push-1"
		// Part of fcmServiceAccountDummy that must not appear in any event.
		privateKeyMarker = "PRIVATE KEY"
	)
	ctx := createContextWithToken(t, true)

	storedPush := func() *domain.Push {
		return &domain.Push{Push: &pushproto.Push{
			Id:                pushID,
			Name:              "name-1",
			Tags:              []string{"tag-1"},
			FcmServiceAccount: string(fcmServiceAccountDummy),
			CreatedAt:         1,
			UpdatedAt:         2,
		}}
	}

	// capture holds what the RPC published and what it wrote to the storage.
	type capture struct {
		event  *eventproto.Event
		stored *pushproto.Push
	}
	// snapshot copies the push as it is written, before the RPC blanks the
	// service account of its own copy for the response.
	snapshot := func(c *capture, push *domain.Push) {
		c.stored = pb.Clone(push.Push).(*pushproto.Push)
	}
	runTransaction := func(s *PushService) {
		s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
			gomock.Any(), gomock.Any(),
		).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	}
	capturePublish := func(s *PushService, c *capture) {
		s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
			gomock.Any(), gomock.Any(),
		).DoAndReturn(func(_ context.Context, msg interface{}) error {
			event, ok := msg.(*eventproto.Event)
			require.True(t, ok)
			c.event = event
			return nil
		})
	}

	cases := []struct {
		name  string
		setup func(*PushService, *capture)
		call  func(context.Context, *PushService) error
		// wantStored is true when the RPC writes the push, with its service
		// account, to the storage.
		wantStored bool
	}{
		{
			name: "create stores the service account and publishes the push without it",
			setup: func(s *PushService, c *capture) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return([]*pushproto.Push{}, 0, int64(0), nil)
				runTransaction(s)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().CreatePush(
					gomock.Any(), gomock.Any(), environmentID,
				).DoAndReturn(func(_ context.Context, push *domain.Push, _ string) error {
					snapshot(c, push)
					return nil
				})
				capturePublish(s, c)
			},
			call: func(ctx context.Context, s *PushService) error {
				resp, err := s.CreatePush(ctx, &pushproto.CreatePushRequest{
					EnvironmentId:     environmentID,
					FcmServiceAccount: fcmServiceAccountDummy,
					Tags:              []string{"tag-1"},
					Name:              "name-1",
				})
				if err != nil {
					return err
				}
				assert.Equal(t, "", resp.Push.FcmServiceAccount, "the response carries no service account")
				return nil
			},
			wantStored: true,
		},
		{
			name: "update keeps the stored service account and publishes the push without it",
			setup: func(s *PushService, c *capture) {
				runTransaction(s)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), pushID, environmentID,
				).Return(storedPush(), nil)
				capturePublish(s, c)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().UpdatePush(
					gomock.Any(), gomock.Any(), environmentID,
				).DoAndReturn(func(_ context.Context, push *domain.Push, _ string) error {
					snapshot(c, push)
					return nil
				})
			},
			call: func(ctx context.Context, s *PushService) error {
				resp, err := s.UpdatePush(ctx, &pushproto.UpdatePushRequest{
					Id:            pushID,
					EnvironmentId: environmentID,
					Name:          wrapperspb.String("name-2"),
					TagChanges: []*pushproto.TagChange{
						{Tag: "tag-2", ChangeType: pushproto.ChangeType_CREATE},
					},
				})
				if err != nil {
					return err
				}
				assert.Equal(t, "", resp.Push.FcmServiceAccount, "the response carries no service account")
				return nil
			},
			wantStored: true,
		},
		{
			name: "delete publishes the deleted push without its service account",
			setup: func(s *PushService, c *capture) {
				runTransaction(s)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), pushID, environmentID,
				).Return(storedPush(), nil)
				capturePublish(s, c)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().DeletePush(
					gomock.Any(), pushID, environmentID,
				).Return(nil)
			},
			call: func(ctx context.Context, s *PushService) error {
				_, err := s.DeletePush(ctx, &pushproto.DeletePushRequest{
					Id:            pushID,
					EnvironmentId: environmentID,
				})
				return err
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			s := newPushServiceWithMock(t, mockController)
			c := &capture{}
			tc.setup(s, c)

			require.NoError(t, tc.call(ctx, s))
			require.NotNil(t, c.event, "the RPC publishes an event")

			// The service account stays with the stored push, which the sender reads.
			if tc.wantStored {
				require.NotNil(t, c.stored)
				assert.Equal(t, string(fcmServiceAccountDummy), c.stored.FcmServiceAccount)
			}
			// The event, which the audit log is built from, does not carry it.
			assert.Equal(t, eventproto.Event_PUSH, c.event.EntityType)
			for _, data := range []string{c.event.EntityData, c.event.PreviousEntityData} {
				if data == "" {
					continue
				}
				entity := map[string]interface{}{}
				require.NoError(t, json.Unmarshal([]byte(data), &entity))
				assert.Equal(t, "", entity["fcm_service_account"])
				assert.NotEmpty(t, entity["id"], "the other fields are recorded")
				assert.NotContains(t, data, privateKeyMarker)
			}
			require.NotNil(t, c.event.Data)
			message, err := c.event.Data.UnmarshalNew()
			require.NoError(t, err)
			if created, ok := message.(*eventproto.PushCreatedEvent); ok {
				assert.Equal(t, "", created.FcmServiceAccount)
				assert.Equal(t, "name-1", created.Name, "the other fields are recorded")
			}
			assert.NotContains(t, string(c.event.Data.Value), privateKeyMarker)
		})
	}
}
