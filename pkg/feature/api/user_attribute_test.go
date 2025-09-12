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

	"github.com/bucketeer-io/bucketeer/pkg/api/api"
	cachev3mock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestGetUserAttributeKeys(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(*FeatureService)
		input          *featureproto.GetUserAttributeKeysRequest
		expected       *featureproto.GetUserAttributeKeysResponse
		getExpectedErr func(localizer locale.Localizer) error
	}{
		{
			desc:    "error: permission denied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: createContextWithTokenRoleUnassigned(),
			setup:   func(s *FeatureService) {},
			input: &featureproto.GetUserAttributeKeysRequest{
				EnvironmentId: "ns0",
			},
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied), localizer)
			},
		},
		{
			desc:    "error: cache error",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			setup: func(s *FeatureService) {
				s.userAttributesCache.(*cachev3mock.MockUserAttributesCache).EXPECT().
					GetUserAttributeKeyAll("ns0").
					Return(nil, pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "cache error"))
			},
			input: &featureproto.GetUserAttributeKeysRequest{
				EnvironmentId: "ns0",
			},
			getExpectedErr: func(localizer locale.Localizer) error {
				return api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "cache error")).Err()
			},
		},
		{
			desc:    "success: empty user attribute keys",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			setup: func(s *FeatureService) {
				s.userAttributesCache.(*cachev3mock.MockUserAttributesCache).EXPECT().
					GetUserAttributeKeyAll("ns0").
					Return([]string{
						"key1",
						"key2",
					}, nil)
			},
			input: &featureproto.GetUserAttributeKeysRequest{
				EnvironmentId: "ns0",
			},
			expected: &featureproto.GetUserAttributeKeysResponse{
				UserAttributeKeys: []string{
					"key1",
					"key2",
				},
			},
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "success: with user attribute keys",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			setup: func(s *FeatureService) {
				expectedKeys := []string{"app_version", "appVersion", "platform", "country", "deviceType"}
				s.userAttributesCache.(*cachev3mock.MockUserAttributesCache).EXPECT().
					GetUserAttributeKeyAll("ns0").
					Return(expectedKeys, nil)
			},
			input: &featureproto.GetUserAttributeKeysRequest{
				EnvironmentId: "ns0",
			},
			expected: &featureproto.GetUserAttributeKeysResponse{
				UserAttributeKeys: []string{"appVersion", "app_version", "country", "deviceType", "platform"},
			},
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "success: with Viewer Account",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: createContextWithTokenRoleUnassigned(),
			setup: func(s *FeatureService) {
				expectedKeys := []string{"appVersion", "platform"}
				s.userAttributesCache.(*cachev3mock.MockUserAttributesCache).EXPECT().
					GetUserAttributeKeyAll("ns0").
					Return(expectedKeys, nil)
			},
			input: &featureproto.GetUserAttributeKeysRequest{
				EnvironmentId: "ns0",
			},
			expected: &featureproto.GetUserAttributeKeysResponse{
				UserAttributeKeys: []string{"appVersion", "platform"},
			},
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			fs := p.service
			if p.setup != nil {
				p.setup(fs)
			}
			ctx := p.context
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"ja"},
			})
			localizer := locale.NewLocalizer(ctx)

			resp, err := fs.GetUserAttributeKeys(ctx, p.input)
			assert.Equal(t, p.getExpectedErr(localizer), err)
			if err == nil {
				assert.Equal(t, p.expected, resp)
			}
		})
	}
}
