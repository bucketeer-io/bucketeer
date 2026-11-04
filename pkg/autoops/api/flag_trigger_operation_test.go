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
	"google.golang.org/grpc/metadata"

	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	ftmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
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

	patterns := []struct {
		desc              string
		setupFunc         func() *ftdomain.Feature
		setupMocks        func(storageMock *ftmock.MockFeatureStorage)
		expectEvent       bool
		expectedEventType eventproto.Event_Type
		expectedErr       error
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
			setupMocks: func(storageMock *ftmock.MockFeatureStorage) {
				storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
					Return(errors.New("internal error")).Times(1)
			},
			// No event is returned when the storage update fails
			expectEvent: false,
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
			setupMocks: func(storageMock *ftmock.MockFeatureStorage) {
				storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
					DoAndReturn(func(ctx context.Context, f *ftdomain.Feature, env string) error {
						// Verify the updated feature has the correct state
						assert.Equal(t, true, f.Enabled)
						assert.Equal(t, int32(2), f.Version)
						return nil
					}).Times(1)
			},
			// The feature domain event is returned to the caller, not published here
			expectEvent:       true,
			expectedEventType: eventproto.Event_FEATURE_ENABLED,
			expectedErr:       nil,
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
			setupMocks: func(storageMock *ftmock.MockFeatureStorage) {
				// No storage calls expected when there's no change
			},
			expectEvent: false,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storageMock := ftmock.NewMockFeatureStorage(mockController)

			// Setup mocks for this specific test case
			p.setupMocks(storageMock)

			event, err := enableFeature(
				ctx,
				storageMock,
				"env",
				p.setupFunc(),
				&eventproto.Editor{
					Email: "email",
				},
			)
			assert.Equal(t, p.expectedErr, err)
			if !p.expectEvent {
				assert.Nil(t, event)
				return
			}
			assert.NotNil(t, event)
			assert.Equal(t, eventproto.Event_FEATURE, event.EntityType)
			assert.Equal(t, p.expectedEventType, event.Type)
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

	patterns := []struct {
		desc              string
		setupFunc         func() *ftdomain.Feature
		setupMocks        func(storageMock *ftmock.MockFeatureStorage)
		expectEvent       bool
		expectedEventType eventproto.Event_Type
		expectedErr       error
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
			setupMocks: func(storageMock *ftmock.MockFeatureStorage) {
				storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
					Return(errors.New("internal error")).Times(1)
			},
			// No event is returned when the storage update fails
			expectEvent: false,
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
			setupMocks: func(storageMock *ftmock.MockFeatureStorage) {
				storageMock.EXPECT().UpdateFeature(ctx, gomock.Any(), "env").
					DoAndReturn(func(ctx context.Context, f *ftdomain.Feature, env string) error {
						// Verify the updated feature has the correct state
						assert.Equal(t, false, f.Enabled)
						assert.Equal(t, int32(2), f.Version)
						return nil
					}).Times(1)
			},
			// The feature domain event is returned to the caller, not published here
			expectEvent:       true,
			expectedEventType: eventproto.Event_FEATURE_DISABLED,
			expectedErr:       nil,
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
			setupMocks: func(storageMock *ftmock.MockFeatureStorage) {
				// No storage calls expected when there's no change
			},
			expectEvent: false,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storageMock := ftmock.NewMockFeatureStorage(mockController)

			// Setup mocks for this specific test case
			p.setupMocks(storageMock)

			event, err := disableFeature(
				ctx,
				storageMock,
				"env",
				p.setupFunc(),
				&eventproto.Editor{
					Email: "email",
				},
			)
			assert.Equal(t, p.expectedErr, err)
			if !p.expectEvent {
				assert.Nil(t, event)
				return
			}
			assert.NotNil(t, event)
			assert.Equal(t, eventproto.Event_FEATURE, event.EntityType)
			assert.Equal(t, p.expectedEventType, event.Type)
		})
	}
}
