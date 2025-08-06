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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	ftdomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	ftmock "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestEnableFeature(t *testing.T) {
	t.Parallel()

	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storageMock := ftmock.NewMockFeatureStorage(mockController)
	publisherMock := publishermock.NewMockPublisher(mockController)

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	logger := zap.NewNop()
	localizer := locale.NewLocalizer(ctx)
	editor := &eventproto.Editor{
		Email: "test@example.com",
	}

	patterns := []struct {
		desc            string
		setupFunc       func() *ftdomain.Feature
		updateCallTimes int
		expectedFunc    func() *ftdomain.Feature
		expectedErr     error
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
			updateCallTimes: 1,
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: true,
					Version: 2,
				}}
			},
			expectedErr: errors.New("err: internal"),
		},
		{
			desc: "success: is already enabled - no storage call",
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
			updateCallTimes: 0, // No storage call when already enabled
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: true,
					Version: 1, // Version unchanged
				}}
			},
			expectedErr: nil,
		},
		{
			desc: "success: disabled to enabled",
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
			updateCallTimes: 1,
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: true,
					Version: 2, // Version incremented
				}}
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			feature := p.setupFunc()

			if p.updateCallTimes > 0 {
				if p.expectedErr != nil {
					storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
						Return(p.expectedErr).Times(p.updateCallTimes)
					// No publisher calls expected when storage update fails
				} else {
					storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
						DoAndReturn(func(ctx context.Context, f *ftdomain.Feature, env string) error {
							// Verify the updated feature has the correct state
							expected := p.expectedFunc()
							assert.Equal(t, expected.Enabled, f.Enabled)
							assert.Equal(t, expected.Version, f.Version)
							return nil
						}).Times(p.updateCallTimes)
					// Expect feature domain event to be published
					publisherMock.EXPECT().Publish(ctx, gomock.Any()).
						DoAndReturn(func(ctx context.Context, event interface{}) error {
							// Verify it's a feature domain event
							domainEvent, ok := event.(*eventproto.Event)
							assert.True(t, ok)
							assert.Equal(t, eventproto.Event_FEATURE, domainEvent.EntityType)
							assert.Equal(t, eventproto.Event_FEATURE_ENABLED, domainEvent.Type)
							return nil
						}).Times(p.updateCallTimes)
				}
			}
			// No storage or publisher calls expected when no changes occur

			err := executeAutoOpsRuleOperation(
				ctx,
				storageMock,
				"env",
				autoopsproto.ActionType_ENABLE,
				feature,
				logger,
				localizer,
				publisherMock,
				editor,
			)

			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDisableFeature(t *testing.T) {
	t.Parallel()

	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storageMock := ftmock.NewMockFeatureStorage(mockController)
	publisherMock := publishermock.NewMockPublisher(mockController)

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	logger := zap.NewNop()
	localizer := locale.NewLocalizer(ctx)
	editor := &eventproto.Editor{
		Email: "test@example.com",
	}

	patterns := []struct {
		desc            string
		setupFunc       func() *ftdomain.Feature
		updateCallTimes int
		expectedFunc    func() *ftdomain.Feature
		expectedErr     error
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
			updateCallTimes: 1,
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: false,
					Version: 2,
				}}
			},
			expectedErr: errors.New("err: internal"),
		},
		{
			desc: "success: is already disabled - no storage call",
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
			updateCallTimes: 0, // No storage call when already disabled
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: false,
					Version: 1, // Version unchanged
				}}
			},
			expectedErr: nil,
		},
		{
			desc: "success: enabled to disabled",
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
			updateCallTimes: 1,
			expectedFunc: func() *ftdomain.Feature {
				return &ftdomain.Feature{Feature: &featureproto.Feature{
					Enabled: false,
					Version: 2, // Version incremented
				}}
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			feature := p.setupFunc()

			if p.updateCallTimes > 0 {
				if p.expectedErr != nil {
					storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
						Return(p.expectedErr).Times(p.updateCallTimes)
					// No publisher calls expected when storage update fails
				} else {
					storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
						DoAndReturn(func(ctx context.Context, f *ftdomain.Feature, env string) error {
							// Verify the updated feature has the correct state
							expected := p.expectedFunc()
							assert.Equal(t, expected.Enabled, f.Enabled)
							assert.Equal(t, expected.Version, f.Version)
							return nil
						}).Times(p.updateCallTimes)
					// Expect feature domain event to be published
					publisherMock.EXPECT().Publish(ctx, gomock.Any()).
						DoAndReturn(func(ctx context.Context, event interface{}) error {
							// Verify it's a feature domain event
							domainEvent, ok := event.(*eventproto.Event)
							assert.True(t, ok)
							assert.Equal(t, eventproto.Event_FEATURE, domainEvent.EntityType)
							assert.Equal(t, eventproto.Event_FEATURE_DISABLED, domainEvent.Type)
							return nil
						}).Times(p.updateCallTimes)
				}
			}
			// No storage or publisher calls expected when no changes occur

			err := executeAutoOpsRuleOperation(
				ctx,
				storageMock,
				"env",
				autoopsproto.ActionType_DISABLE,
				feature,
				logger,
				localizer,
				publisherMock,
				editor,
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
