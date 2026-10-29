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
	"google.golang.org/grpc/metadata"
	pb "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestFlagTriggerEventsOmitToken(t *testing.T) {
	t.Parallel()

	const (
		// Environment the account of createFeatureServiceNew is an editor of.
		environmentID = "namespace"
		featureID     = "feature-1"
		triggerID     = "trigger-1"
		storedToken   = "stored-token-8xK2mPqR5tV7wY9zA1bC3dE4fG6hJ8kL0mN2oP4qR6sT"
	)

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)

	storedTrigger := func() *domain.FlagTrigger {
		return &domain.FlagTrigger{FlagTrigger: &featureproto.FlagTrigger{
			Id:            triggerID,
			FeatureId:     featureID,
			EnvironmentId: environmentID,
			Type:          featureproto.FlagTrigger_Type_WEBHOOK,
			Action:        featureproto.FlagTrigger_Action_ON,
			Description:   "description",
			Token:         storedToken,
			CreatedAt:     1,
			UpdatedAt:     2,
		}}
	}

	// capture holds what the RPC published and what it wrote to the storage.
	type capture struct {
		event  *eventproto.Event
		stored *featureproto.FlagTrigger
	}
	// snapshot copies the trigger as it is written, before the RPC blanks the
	// token of its own copy for the response.
	snapshot := func(c *capture, trigger *domain.FlagTrigger) {
		c.stored = pb.Clone(trigger.FlagTrigger).(*featureproto.FlagTrigger)
	}
	runTransaction := func(s *FeatureService) {
		s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
			gomock.Any(), gomock.Any(),
		).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	}
	capturePublish := func(s *FeatureService, c *capture) {
		s.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
			gomock.Any(), gomock.Any(),
		).DoAndReturn(func(_ context.Context, msg interface{}) error {
			event, ok := msg.(*eventproto.Event)
			require.True(t, ok)
			c.event = event
			return nil
		})
	}
	captureUpdate := func(s *FeatureService, c *capture) {
		s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().UpdateFlagTrigger(
			gomock.Any(), gomock.Any(),
		).DoAndReturn(func(_ context.Context, trigger *domain.FlagTrigger) error {
			snapshot(c, trigger)
			return nil
		})
	}

	cases := []struct {
		name  string
		setup func(*FeatureService, *capture)
		// call runs the RPC and returns the trigger URL of its response.
		call func(context.Context, *FeatureService) (string, error)
		// wantNewToken is true when the RPC generates a new token, which must be
		// stored and returned, once, in the response URL.
		wantNewToken bool
	}{
		{
			name: "create returns the token once and publishes the trigger without it",
			setup: func(s *FeatureService, c *capture) {
				runTransaction(s)
				capturePublish(s, c)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().CreateFlagTrigger(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(_ context.Context, trigger *domain.FlagTrigger) error {
					snapshot(c, trigger)
					return nil
				})
			},
			call: func(ctx context.Context, s *FeatureService) (string, error) {
				resp, err := s.CreateFlagTrigger(ctx, &featureproto.CreateFlagTriggerRequest{
					EnvironmentId: environmentID,
					FeatureId:     featureID,
					Type:          featureproto.FlagTrigger_Type_WEBHOOK,
					Action:        featureproto.FlagTrigger_Action_ON,
					Description:   "description",
				})
				if err != nil {
					return "", err
				}
				assert.Equal(t, "", resp.FlagTrigger.Token, "the response carries the token in the URL only")
				return resp.Url, nil
			},
			wantNewToken: true,
		},
		{
			name: "reset returns the new token once and publishes neither the old nor the new one",
			setup: func(s *FeatureService, c *capture) {
				runTransaction(s)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					gomock.Any(), triggerID, environmentID,
				).Return(storedTrigger(), nil)
				capturePublish(s, c)
				captureUpdate(s, c)
			},
			call: func(ctx context.Context, s *FeatureService) (string, error) {
				resp, err := s.UpdateFlagTrigger(ctx, &featureproto.UpdateFlagTriggerRequest{
					Id:            triggerID,
					EnvironmentId: environmentID,
					Reset_:        true,
				})
				if err != nil {
					return "", err
				}
				return resp.Url, nil
			},
			wantNewToken: true,
		},
		{
			name: "update keeps the stored token and publishes the trigger without it",
			setup: func(s *FeatureService, c *capture) {
				runTransaction(s)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					gomock.Any(), triggerID, environmentID,
				).Return(storedTrigger(), nil)
				capturePublish(s, c)
				captureUpdate(s, c)
			},
			call: func(ctx context.Context, s *FeatureService) (string, error) {
				resp, err := s.UpdateFlagTrigger(ctx, &featureproto.UpdateFlagTriggerRequest{
					Id:            triggerID,
					EnvironmentId: environmentID,
					Description:   wrapperspb.String("new description"),
					Disabled:      wrapperspb.Bool(true),
				})
				if err != nil {
					return "", err
				}
				return resp.Url, nil
			},
		},
		{
			name: "delete publishes the deleted trigger without its token",
			setup: func(s *FeatureService, c *capture) {
				runTransaction(s)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					gomock.Any(), triggerID, environmentID,
				).Return(storedTrigger(), nil)
				capturePublish(s, c)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().DeleteFlagTrigger(
					gomock.Any(), triggerID, environmentID,
				).Return(nil)
			},
			call: func(ctx context.Context, s *FeatureService) (string, error) {
				_, err := s.DeleteFlagTrigger(ctx, &featureproto.DeleteFlagTriggerRequest{
					Id:            triggerID,
					EnvironmentId: environmentID,
				})
				return "", err
			},
		},
		{
			name: "webhook keeps the stored token and publishes the usage without it",
			setup: func(s *FeatureService, c *capture) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTriggerByToken(
					gomock.Any(), storedToken,
				).Return(storedTrigger(), nil)
				// The flag is already on, so only the trigger usage is updated.
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().GetFeature(
					gomock.Any(), featureID, environmentID,
				).Return(&domain.Feature{Feature: &featureproto.Feature{Id: featureID, Enabled: true}}, nil)
				runTransaction(s)
				capturePublish(s, c)
				captureUpdate(s, c)
			},
			call: func(ctx context.Context, s *FeatureService) (string, error) {
				_, err := s.FlagTriggerWebhook(ctx, &featureproto.FlagTriggerWebhookRequest{Token: storedToken})
				return "", err
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			s := createFeatureServiceNew(mockController)
			c := &capture{}
			tc.setup(s, c)

			url, err := tc.call(ctx, s)
			require.NoError(t, err)
			require.NotNil(t, c.event, "the RPC publishes an event")

			// The token stays with the stored trigger and is returned, at most
			// once, to the editor.
			secrets := []string{storedToken}
			wantID := triggerID
			if tc.wantNewToken {
				require.NotNil(t, c.stored)
				require.NotEmpty(t, c.stored.Token)
				assert.NotEqual(t, storedToken, c.stored.Token)
				assert.Equal(t, s.triggerURL+"/"+c.stored.Token, url)
				secrets = append(secrets, c.stored.Token)
				wantID = c.stored.Id
			} else {
				if c.stored != nil {
					assert.Equal(t, storedToken, c.stored.Token, "the stored token is unchanged")
				}
				assert.Equal(t, "", url)
			}
			// The event, which the audit log is built from, carries none of them.
			assert.Equal(t, eventproto.Event_FLAG_TRIGGER, c.event.EntityType)
			assertFlagTriggerEventOmitsToken(t, c.event, wantID, secrets)
		})
	}
}

// assertFlagTriggerEventOmitsToken checks that the entity data, the previous
// entity data and the typed event data of a flag trigger event carry no token,
// while the other fields of the trigger are still recorded.
func assertFlagTriggerEventOmitsToken(t *testing.T, event *eventproto.Event, wantID string, secrets []string) {
	t.Helper()
	for _, data := range []string{event.EntityData, event.PreviousEntityData} {
		if data == "" {
			continue
		}
		entity := map[string]interface{}{}
		require.NoError(t, json.Unmarshal([]byte(data), &entity))
		assert.Equal(t, "", entity["token"])
		assert.Equal(t, wantID, entity["id"], "the other fields are recorded")
		for _, secret := range secrets {
			assert.NotContains(t, data, secret)
		}
	}
	require.NotNil(t, event.Data)
	message, err := event.Data.UnmarshalNew()
	require.NoError(t, err)
	if created, ok := message.(*eventproto.FlagTriggerCreatedEvent); ok {
		assert.Equal(t, "", created.Token)
	}
	for _, secret := range secrets {
		assert.NotContains(t, string(event.Data.Value), secret)
	}
}
