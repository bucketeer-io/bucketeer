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
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	ftmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestEnableFeature(t *testing.T) {
	t.Parallel()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	logger := zap.NewNop()

	patterns := []struct {
		desc         string
		setupFunc    func() *ftdomain.Feature
		setupMocks   func(storageMock *ftmock.MockFeatureStorage, publisherMock *publishermock.MockPublisher)
		expectedFunc func() *ftdomain.Feature
		expectedErr  error
	}{
		{
			desc: "err: internal",
			setupFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: false,
					Version: 1,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			setupMocks: func(storageMock *ftmock.MockFeatureStorage, publisherMock *publishermock.MockPublisher) {
				storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
					Return(errors.New("internal error")).Times(1)
				// No publisher calls expected when storage update fails
			},
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: true,
					Version: 2,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			expectedErr: errors.New("internal error"),
		},
		{
			desc: "success",
			setupFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: false,
					Version: 1,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			setupMocks: func(storageMock *ftmock.MockFeatureStorage, publisherMock *publishermock.MockPublisher) {
				storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
					DoAndReturn(func(ctx context.Context, f *ftdomain.Feature, env string) error {
						// Verify the updated feature has the correct state
						assert.Equal(t, true, f.Enabled)
						assert.Equal(t, int32(2), f.Version)
						return nil
					}).Times(1)
				// Expect feature domain event to be published
				publisherMock.EXPECT().Publish(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, event interface{}) error {
						// Verify it's a feature domain event
						domainEvent, ok := event.(*eventproto.Event)
						assert.True(t, ok)
						assert.Equal(t, eventproto.Event_FEATURE, domainEvent.EntityType)
						assert.Equal(t, eventproto.Event_FEATURE_ENABLED, domainEvent.Type)
						return nil
					}).Times(1)
			},
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: true,
					Version: 2,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			expectedErr: nil,
		},
		{
			desc: "no change: already enabled",
			setupFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: true,
					Version: 1,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			setupMocks: func(storageMock *ftmock.MockFeatureStorage, publisherMock *publishermock.MockPublisher) {
				// No storage or publisher calls expected when there's no change
			},
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: true,
					Version: 1,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storageMock := ftmock.NewMockFeatureStorage(mockController)
			publisherMock := publishermock.NewMockPublisher(mockController)

			// Setup mocks for this specific test case
			p.setupMocks(storageMock, publisherMock)

			err := enableFeature(
				ctx,
				storageMock,
				"env",
				p.setupFunc(),
				logger,
				publisherMock,
				&eventproto.Editor{
					Email: "email",
				},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDisableFeature(t *testing.T) {
	t.Parallel()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	logger := zap.NewNop()

	patterns := []struct {
		desc         string
		setupFunc    func() *ftdomain.Feature
		setupMocks   func(storageMock *ftmock.MockFeatureStorage, publisherMock *publishermock.MockPublisher)
		expectedFunc func() *ftdomain.Feature
		expectedErr  error
	}{
		{
			desc: "err: internal",
			setupFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: true,
					Version: 1,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			setupMocks: func(storageMock *ftmock.MockFeatureStorage, publisherMock *publishermock.MockPublisher) {
				storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
					Return(errors.New("internal error")).Times(1)
				// No publisher calls expected when storage update fails
			},
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: false,
					Version: 2,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			expectedErr: errors.New("internal error"),
		},
		{
			desc: "success",
			setupFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: true,
					Version: 1,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			setupMocks: func(storageMock *ftmock.MockFeatureStorage, publisherMock *publishermock.MockPublisher) {
				storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
					DoAndReturn(func(ctx context.Context, f *ftdomain.Feature, env string) error {
						// Verify the updated feature has the correct state
						assert.Equal(t, false, f.Enabled)
						assert.Equal(t, int32(2), f.Version)
						return nil
					}).Times(1)
				// Expect feature domain event to be published
				publisherMock.EXPECT().Publish(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, event interface{}) error {
						// Verify it's a feature domain event
						domainEvent, ok := event.(*eventproto.Event)
						assert.True(t, ok)
						assert.Equal(t, eventproto.Event_FEATURE, domainEvent.EntityType)
						assert.Equal(t, eventproto.Event_FEATURE_DISABLED, domainEvent.Type)
						return nil
					}).Times(1)
			},
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: false,
					Version: 2,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			expectedErr: nil,
		},
		{
			desc: "no change: already disabled",
			setupFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: false,
					Version: 1,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			setupMocks: func(storageMock *ftmock.MockFeatureStorage, publisherMock *publishermock.MockPublisher) {
				// No storage or publisher calls expected when there's no change
			},
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: false,
					Version: 1,
					Variations: []*featureproto.Variation{
						{Id: "vid-1", Value: "true", Name: "variation-1"},
						{Id: "vid-2", Value: "false", Name: "variation-2"},
					},
					OffVariation: "vid-2",
				}}
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storageMock := ftmock.NewMockFeatureStorage(mockController)
			publisherMock := publishermock.NewMockPublisher(mockController)

			// Setup mocks for this specific test case
			p.setupMocks(storageMock, publisherMock)

			err := disableFeature(
				ctx,
				storageMock,
				"env",
				p.setupFunc(),
				logger,
				publisherMock,
				&eventproto.Editor{
					Email: "email",
				},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
