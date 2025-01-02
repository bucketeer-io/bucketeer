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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	ftdomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	ftmock "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestEnableFeature(t *testing.T) {
	t.Parallel()

	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storageMock := ftmock.NewMockFeatureStorage(mockController)

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	logger := zap.NewNop()
	localizer := locale.NewLocalizer(ctx)

	patterns := []struct {
		desc            string
		feature         *featureproto.Feature
		autoOpsRuleType autoopsproto.OpsType
		updateCallTimes int
		expected        bool
		expectedErr     error
	}{
		{
			desc: "err: internal",
			feature: &featureproto.Feature{
				Enabled: false,
			},
			updateCallTimes: 1,
			expected:        true,
			expectedErr:     errors.New("err: internal"),
		},
		{
			desc: "success: is already enabled",
			feature: &featureproto.Feature{
				Enabled: true,
			},
			updateCallTimes: 0,
			expected:        true,
			expectedErr:     nil,
		},
		{
			desc: "success",
			feature: &featureproto.Feature{
				Enabled: false,
			},
			updateCallTimes: 1,
			expected:        true,
			expectedErr:     nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			feature := &ftdomain.Feature{Feature: p.feature}

			if p.expectedErr != nil {
				storageMock.EXPECT().UpdateFeature(ctx, feature, "env").
					Return(p.expectedErr).Times(p.updateCallTimes)
			} else {
				storageMock.EXPECT().UpdateFeature(ctx, feature, "env").
					Return(nil).Times(p.updateCallTimes)
			}

			err := executeAutoOpsRuleOperation(
				ctx,
				storageMock,
				"env",
				autoopsproto.ActionType_ENABLE,
				feature,
				logger,
				localizer,
			)
			assert.Equal(t, p.expected, feature.Enabled)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDisableFeature(t *testing.T) {
	t.Parallel()

	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storageMock := ftmock.NewMockFeatureStorage(mockController)

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	logger := zap.NewNop()
	localizer := locale.NewLocalizer(ctx)

	patterns := []struct {
		desc            string
		feature         *featureproto.Feature
		autoOpsRuleType autoopsproto.OpsType
		updateCallTimes int
		expected        bool
		expectedErr     error
	}{
		{
			desc: "err: internal",
			feature: &featureproto.Feature{
				Enabled: true,
			},
			updateCallTimes: 1,
			expected:        false,
			expectedErr:     errors.New("err: internal"),
		},
		{
			desc: "success: is already disabled",
			feature: &featureproto.Feature{
				Enabled: false,
			},
			updateCallTimes: 0,
			expected:        false,
			expectedErr:     nil,
		},
		{
			desc: "success",
			feature: &featureproto.Feature{
				Enabled: true,
			},
			updateCallTimes: 1,
			expected:        false,
			expectedErr:     nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			feature := &ftdomain.Feature{Feature: p.feature}

			if p.expectedErr != nil {
				storageMock.EXPECT().UpdateFeature(ctx, feature, "env").
					Return(p.expectedErr).Times(p.updateCallTimes)
			} else {
				storageMock.EXPECT().UpdateFeature(ctx, feature, "env").
					Return(nil).Times(p.updateCallTimes)
			}

			err := executeAutoOpsRuleOperation(
				ctx,
				storageMock,
				"env",
				autoopsproto.ActionType_DISABLE,
				feature,
				logger,
				localizer,
			)
			assert.Equal(t, p.expected, feature.Enabled)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
