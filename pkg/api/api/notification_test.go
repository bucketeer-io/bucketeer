// Copyright 2025 The Bucketeer Authors.
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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	cachev3mock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	mocknotificationclient "github.com/bucketeer-io/bucketeer/pkg/notification/client/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	"github.com/bucketeer-io/bucketeer/proto/notification"
)

func TestGrpcGatewayService_GetSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetSubscriptionResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: GetSubscription error",
			ctx:  context.Background(),
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.notificationClient.(*mocknotificationclient.MockClient).EXPECT().
					GetSubscription(gomock.Any(), gomock.Any()).
					Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			ctx:  context.Background(),
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.notificationClient.(*mocknotificationclient.MockClient).EXPECT().
					GetSubscription(gomock.Any(), gomock.Any()).
					Return(&notification.GetSubscriptionResponse{
						Subscription: &notification.Subscription{
							Id: "sub-1",
							SourceTypes: []notification.Subscription_SourceType{
								notification.Subscription_DOMAIN_EVENT_ACCOUNT,
								notification.Subscription_DOMAIN_EVENT_FEATURE,
							},
							Recipient: &notification.Recipient{
								Type: notification.Recipient_SlackChannel,
								SlackChannelRecipient: &notification.SlackChannelRecipient{
									WebhookUrl: "webhook-url",
								},
								Language: notification.Recipient_ENGLISH,
							},
							Name:            "sub-1",
							EnvironmentId:   "env-id-1",
							EnvironmentName: "env-name-1",
							FeatureFlagTags: []string{"log"},
						},
					}, nil)
			},
			expected: &gwproto.GetSubscriptionResponse{
				Subscription: &notification.Subscription{
					Id: "sub-1",
					SourceTypes: []notification.Subscription_SourceType{
						notification.Subscription_DOMAIN_EVENT_ACCOUNT,
						notification.Subscription_DOMAIN_EVENT_FEATURE,
					},
					Recipient: &notification.Recipient{
						Type: notification.Recipient_SlackChannel,
						SlackChannelRecipient: &notification.SlackChannelRecipient{
							WebhookUrl: "webhook-url",
						},
						Language: notification.Recipient_ENGLISH,
					},
					Name:            "sub-1",
					EnvironmentId:   "env-id-1",
					EnvironmentName: "env-name-1",
					FeatureFlagTags: []string{"log"},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetSubscription(ctx, &gwproto.GetSubscriptionRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_ListSubscriptions(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.ListSubscriptionsResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: ListSubscriptions error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.notificationClient.(*mocknotificationclient.MockClient).EXPECT().ListSubscriptions(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.notificationClient.(*mocknotificationclient.MockClient).EXPECT().ListSubscriptions(
					gomock.Any(), gomock.Any(),
				).Return(&notification.ListSubscriptionsResponse{
					Subscriptions: []*notification.Subscription{
						{
							Id: "sub-1",
							SourceTypes: []notification.Subscription_SourceType{
								notification.Subscription_DOMAIN_EVENT_ACCOUNT,
								notification.Subscription_DOMAIN_EVENT_FEATURE,
							},
							Recipient: &notification.Recipient{
								Type: notification.Recipient_SlackChannel,
								SlackChannelRecipient: &notification.SlackChannelRecipient{
									WebhookUrl: "webhook-url",
								},
								Language: notification.Recipient_ENGLISH,
							},
							Name:            "sub-1",
							EnvironmentId:   "env-id-1",
							EnvironmentName: "env-name-1",
							FeatureFlagTags: []string{"log"},
						},
					},
					Cursor:     "0",
					TotalCount: 1,
				}, nil)
			},
			expected: &gwproto.ListSubscriptionsResponse{
				Subscriptions: []*notification.Subscription{
					{
						Id: "sub-1",
						SourceTypes: []notification.Subscription_SourceType{
							notification.Subscription_DOMAIN_EVENT_ACCOUNT,
							notification.Subscription_DOMAIN_EVENT_FEATURE,
						},
						Recipient: &notification.Recipient{
							Type: notification.Recipient_SlackChannel,
							SlackChannelRecipient: &notification.SlackChannelRecipient{
								WebhookUrl: "webhook-url",
							},
							Language: notification.Recipient_ENGLISH,
						},
						Name:            "sub-1",
						EnvironmentId:   "env-id-1",
						EnvironmentName: "env-name-1",
						FeatureFlagTags: []string{"log"},
					},
				},
				Cursor:     "0",
				TotalCount: 1,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.ListSubscriptions(ctx, &gwproto.ListSubscriptionsRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_CreateSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.CreateSubscriptionResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: CreateSubscriptions error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.notificationClient.(*mocknotificationclient.MockClient).EXPECT().CreateSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			ctx:  context.Background(),
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.notificationClient.(*mocknotificationclient.MockClient).EXPECT().CreateSubscription(
					gomock.Any(), gomock.Any(),
				).Return(&notification.CreateSubscriptionResponse{
					Subscription: &notification.Subscription{
						Id: "sub-1",
						SourceTypes: []notification.Subscription_SourceType{
							notification.Subscription_DOMAIN_EVENT_ACCOUNT,
							notification.Subscription_DOMAIN_EVENT_FEATURE,
						},
						Recipient: &notification.Recipient{
							Type: notification.Recipient_SlackChannel,
							SlackChannelRecipient: &notification.SlackChannelRecipient{
								WebhookUrl: "webhook-url",
							},
							Language: notification.Recipient_ENGLISH,
						},
						Name:            "sub-1",
						EnvironmentId:   "env-id-1",
						EnvironmentName: "env-name-1",
						FeatureFlagTags: []string{"log"},
					},
				}, nil)
			},
			expected: &gwproto.CreateSubscriptionResponse{
				Subscription: &notification.Subscription{
					Id: "sub-1",
					SourceTypes: []notification.Subscription_SourceType{
						notification.Subscription_DOMAIN_EVENT_ACCOUNT,
						notification.Subscription_DOMAIN_EVENT_FEATURE,
					},
					Recipient: &notification.Recipient{
						Type: notification.Recipient_SlackChannel,
						SlackChannelRecipient: &notification.SlackChannelRecipient{
							WebhookUrl: "webhook-url",
						},
						Language: notification.Recipient_ENGLISH,
					},
					Name:            "sub-1",
					EnvironmentId:   "env-id-1",
					EnvironmentName: "env-name-1",
					FeatureFlagTags: []string{"log"},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.CreateSubscription(ctx, &gwproto.CreateSubscriptionRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_DeleteSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.DeleteSubscriptionResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: DeleteSubscription error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.notificationClient.(*mocknotificationclient.MockClient).EXPECT().DeleteSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.notificationClient.(*mocknotificationclient.MockClient).EXPECT().DeleteSubscription(
					gomock.Any(), gomock.Any(),
				).Return(&notification.DeleteSubscriptionResponse{}, nil)
			},
			expected:    &gwproto.DeleteSubscriptionResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.DeleteSubscription(ctx, &gwproto.DeleteSubscriptionRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_UpdateSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.UpdateSubscriptionResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: UpdateSubscription error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.notificationClient.(*mocknotificationclient.MockClient).EXPECT().UpdateSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.notificationClient.(*mocknotificationclient.MockClient).EXPECT().UpdateSubscription(
					gomock.Any(), gomock.Any(),
				).Return(&notification.UpdateSubscriptionResponse{}, nil)
			},
			expected:    &gwproto.UpdateSubscriptionResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.UpdateSubscription(ctx, &gwproto.UpdateSubscriptionRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}
