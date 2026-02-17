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
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	btclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/batch/client/mock"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	envclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	exprclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	tagstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	exprproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

func TestGetFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(*FeatureService)
		input          *featureproto.GetFeatureRequest
		getExpectedErr func() error
	}{
		{
			desc:    "error: id is empty",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			input: &featureproto.GetFeatureRequest{
				Id:            "",
				EnvironmentId: "ns0",
			},
			getExpectedErr: func() error {
				return statusMissingID.Err()
			},
		},
		{
			desc:    "success",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.fluiStorage.(*mock.MockFeatureLastUsedInfoStorage).EXPECT().GetFeatureLastUsedInfos(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &featureproto.GetFeatureRequest{
				Id:            "fid",
				EnvironmentId: "ns0",
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success with Viewer Account",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: createContextWithTokenRoleUnassigned(),
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.fluiStorage.(*mock.MockFeatureLastUsedInfoStorage).EXPECT().GetFeatureLastUsedInfos(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &featureproto.GetFeatureRequest{
				Id:            "fid",
				EnvironmentId: "ns0",
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success get feature by version",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.fluiStorage.(*mock.MockFeatureLastUsedInfoStorage).EXPECT().GetFeatureLastUsedInfos(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &featureproto.GetFeatureRequest{
				Id:             "fid",
				EnvironmentId:  "ns0",
				FeatureVersion: wrapperspb.Int32(1),
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: createContextWithTokenRoleUnassigned(),
			setup:   func(s *FeatureService) {},
			input: &featureproto.GetFeatureRequest{
				Id:            "fid",
				EnvironmentId: "ns0",
			},
			getExpectedErr: func() error {
				return statusPermissionDenied.Err()
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

			_, err := fs.GetFeature(ctx, p.input)
			assert.Equal(t, p.getExpectedErr(), err)
		})
	}
}

func TestGetFeaturesMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(*FeatureService)
		input          []string
		getExpectedErr func() error
	}{
		{
			desc:    "error: id is nil",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			input:   nil,
			getExpectedErr: func() error {
				return statusMissingIDs.Err()
			},
		},
		{
			desc:    "error: contains empty id",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			input:   []string{"id", ""},
			getExpectedErr: func() error {
				return statusMissingIDs.Err()
			},
		},
		{
			desc:    "success",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input: []string{"fid"},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success with Viewer Account",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: createContextWithTokenRoleUnassigned(),
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input: []string{"fid"},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: createContextWithTokenRoleUnassigned(),
			setup:   func(s *FeatureService) {},
			input:   []string{"fid"},
			getExpectedErr: func() error {
				return statusPermissionDenied.Err()
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			fs := p.service
			if p.setup != nil {
				p.setup(fs)
			}
			req := &featureproto.GetFeaturesRequest{
				EnvironmentId: "ns0",
				Ids:           p.input,
			}
			ctx := p.context
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"ja"},
			})

			_, err := fs.GetFeatures(ctx, req)
			assert.Equal(t, p.getExpectedErr(), err)
		})
	}
}

func TestListFeaturesMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(*FeatureService)
		orderBy        featureproto.ListFeaturesRequest_OrderBy
		hasExperiment  bool
		environmentId  string
		getExpectedErr func() error
	}{
		{
			desc:          "error: invalid order by",
			service:       createFeatureService(mockController),
			context:       createContextWithToken(),
			setup:         nil,
			orderBy:       featureproto.ListFeaturesRequest_OrderBy(999),
			hasExperiment: false,
			environmentId: "ns0",
			getExpectedErr: func() error {
				return statusInvalidOrderBy.Err()
			},
		},
		{
			desc:    "success do not has Experiment",
			service: createFeatureService(mockController),
			context: createContextWithToken(),
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().GetFeatureSummary(
					gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			orderBy:        featureproto.ListFeaturesRequest_DEFAULT,
			hasExperiment:  false,
			environmentId:  "ns0",
			getExpectedErr: func() error { return nil },
		},
		{
			desc:    "success has Experiment",
			service: createFeatureService(mockController),
			context: createContextWithToken(),
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().GetFeatureSummary(
					gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			orderBy:       featureproto.ListFeaturesRequest_DEFAULT,
			hasExperiment: true,
			environmentId: "ns0",
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success with Viewer Account",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: createContextWithTokenRoleUnassigned(),
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().GetFeatureSummary(
					gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			orderBy:        featureproto.ListFeaturesRequest_DEFAULT,
			hasExperiment:  false,
			environmentId:  "ns0",
			getExpectedErr: func() error { return nil },
		},
		{
			desc:          "errPermissionDenied",
			service:       createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context:       createContextWithTokenRoleUnassigned(),
			setup:         func(s *FeatureService) {},
			orderBy:       featureproto.ListFeaturesRequest_DEFAULT,
			hasExperiment: false,
			environmentId: "ns0",
			getExpectedErr: func() error {
				return statusPermissionDenied.Err()
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := p.service
			if p.setup != nil {
				p.setup(service)
			}
			req := &featureproto.ListFeaturesRequest{
				OrderBy:       p.orderBy,
				EnvironmentId: "ns0",
				HasExperiment: &wrappers.BoolValue{Value: p.hasExperiment},
			}
			ctx := p.context
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"ja"},
			})

			_, err := service.ListFeatures(ctx, req)
			assert.Equal(t, p.getExpectedErr(), err)
		})
	}
}

func TestCreateFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	variations := createFeatureVariations()
	tags := createFeatureTags()
	patterns := []struct {
		setup                                             func(*FeatureService)
		id, name, description                             string
		variations                                        []*featureproto.Variation
		tags                                              []string
		defaultOnVariationIndex, defaultOffVariationIndex *wrappers.Int32Value
		environmentId                                     string
		expected                                          error
	}{
		{
			setup:                    nil,
			id:                       "",
			name:                     "name",
			description:              "error: statusMissingID",
			variations:               nil,
			tags:                     nil,
			defaultOnVariationIndex:  nil,
			defaultOffVariationIndex: nil,
			environmentId:            "ns0",
			expected:                 statusMissingID.Err(),
		},
		{
			setup:                    nil,
			id:                       "bucketeer_id",
			name:                     "name",
			description:              "error: statusInvalidID",
			variations:               nil,
			tags:                     nil,
			defaultOnVariationIndex:  nil,
			defaultOffVariationIndex: nil,
			environmentId:            "ns0",
			expected:                 statusInvalidID.Err(),
		},
		{
			setup:                    nil,
			id:                       "Bucketeer-id-2019",
			name:                     "",
			description:              "error: statusMissingName",
			variations:               nil,
			tags:                     nil,
			defaultOnVariationIndex:  nil,
			defaultOffVariationIndex: nil,
			environmentId:            "ns0",
			expected:                 statusMissingName.Err(),
		},
		{
			setup:                    nil,
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "error: statusMissingFeatureVariations",
			variations:               nil,
			tags:                     nil,
			defaultOnVariationIndex:  nil,
			defaultOffVariationIndex: nil,
			environmentId:            "ns0",
			expected:                 statusMissingFeatureVariations.Err(),
		},
		{
			setup:         nil,
			id:            "Bucketeer-id-2019",
			name:          "name",
			description:   "error: statusMissingFeatureTags",
			variations:    variations,
			tags:          nil,
			environmentId: "ns0",
			expected:      statusMissingFeatureTags.Err(),
		},
		{
			setup:                    nil,
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "error: statusMissingDefaultOnVariation",
			variations:               variations,
			tags:                     tags,
			defaultOnVariationIndex:  nil,
			defaultOffVariationIndex: nil,
			environmentId:            "ns0",
			expected:                 statusMissingDefaultOnVariation.Err(),
		},
		{
			setup:                    nil,
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "error: statusMissingDefaultOffVariation",
			variations:               variations,
			tags:                     tags,
			defaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
			defaultOffVariationIndex: nil,
			environmentId:            "ns0",
			expected:                 statusMissingDefaultOffVariation.Err(),
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrFeatureAlreadyExists)
			},
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "error: statusAlreadyExists",
			variations:               variations,
			tags:                     tags,
			defaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
			defaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
			environmentId:            "ns0",
			expected:                 statusAlreadyExists.Err(),
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.batchClient.(*btclientmock.MockClient).EXPECT().ExecuteBatchJob(gomock.Any(), gomock.Any())
			},
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "success",
			variations:               variations,
			tags:                     tags,
			defaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
			defaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
			environmentId:            "ns0",
			expected:                 nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.description, func(t *testing.T) {
			service := createFeatureService(mockController)
			if p.setup != nil {
				p.setup(service)
			}
			req := &featureproto.CreateFeatureRequest{
				Id:                       p.id,
				Name:                     p.name,
				Description:              p.description,
				Variations:               p.variations,
				Tags:                     p.tags,
				DefaultOnVariationIndex:  p.defaultOnVariationIndex,
				DefaultOffVariationIndex: p.defaultOffVariationIndex,
				EnvironmentId:            p.environmentId,
			}
			actual, err := service.CreateFeature(ctx, req)
			if p.expected == nil {
				assert.NotNil(t, actual.Feature)
			}
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestSetFeatureToLastUsedInfosByChunk(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		setup         func(*FeatureService)
		input         []*featureproto.Feature
		environmentId string
		expected      error
	}{
		{
			setup: func(s *FeatureService) {
				s.fluiStorage.(*mock.MockFeatureLastUsedInfoStorage).EXPECT().GetFeatureLastUsedInfos(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*domain.FeatureLastUsedInfo{
					{
						FeatureLastUsedInfo: &featureproto.FeatureLastUsedInfo{
							FeatureId:  "feature-id-0",
							LastUsedAt: time.Now().Unix(),
						},
					},
				}, nil)
			},
			input: []*featureproto.Feature{
				{
					Id:      "feature-id-0",
					Version: 1,
				},
			},
			environmentId: "ns0",
			expected:      nil,
		},
	}
	for _, p := range patterns {
		fs := createFeatureServiceNew(mockController)
		p.setup(fs)
		err := fs.setLastUsedInfosToFeatureByChunk(context.Background(), p.input, p.environmentId)
		assert.Equal(t, p.expected, err)
	}
}

func TestConvUpdateFeatureError(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		input       error
		expectedErr error
	}{
		{
			input:       v2fs.ErrFeatureNotFound,
			expectedErr: statusFeatureNotFound.Err(),
		},
		{
			input:       v2fs.ErrFeatureUnexpectedAffectedRows,
			expectedErr: statusFeatureNotFound.Err(),
		},
		{
			input:       storage.ErrKeyNotFound,
			expectedErr: statusFeatureNotFound.Err(),
		},
		{
			input:       pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "test"),
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "test")).Err(),
		},
	}
	for _, p := range patterns {
		fs := &FeatureService{}
		err := fs.convUpdateFeatureError(p.input)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestEvaluateFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	vID1 := newUUID(t)
	vID2 := newUUID(t)
	vID3 := newUUID(t)
	vID4 := newUUID(t)

	patterns := []struct {
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(*FeatureService)
		input          *featureproto.EvaluateFeaturesRequest
		expected       *featureproto.EvaluateFeaturesResponse
		getExpectedErr func() error
	}{
		{
			desc:     "fail: ErrMissingUser",
			context:  createContextWithToken(),
			service:  createFeatureService(mockController),
			setup:    nil,
			input:    &featureproto.EvaluateFeaturesRequest{},
			expected: nil,
			getExpectedErr: func() error {
				return statusMissingUser.Err()
			},
		},
		{
			desc:     "fail: ErrMissingUserID",
			context:  createContextWithToken(),
			service:  createFeatureService(mockController),
			setup:    nil,
			input:    &featureproto.EvaluateFeaturesRequest{User: &userproto.User{}},
			expected: nil,
			getExpectedErr: func() error {
				return statusMissingUserID.Err()
			},
		},
		{
			desc:    "fail: return errInternal when getting features",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error"))
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error"))
			},
			input:    &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentId: "ns0", Tag: "android"},
			expected: nil,
			getExpectedErr: func() error {
				return api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error")).Err()
			},
		},
		{
			desc:    "success: get from cache",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: newUUID(t),
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID2,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"segment-id",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID1,
									},
								},
								Tags: []string{"android"},
							},
							{
								Id: newUUID(t),
								Variations: []*featureproto.Variation{
									{
										Id:    vID3,
										Value: "true",
									},
									{
										Id:    vID4,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID4,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"segment-id",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID3,
									},
								},
								Tags: []string{"ios"},
							},
						},
					}, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&featureproto.SegmentUsers{
						SegmentId: "segment-id",
						Users: []*featureproto.SegmentUser{
							{
								SegmentId: "segment-id",
								UserId:    "user-id-1",
								State:     featureproto.SegmentUser_INCLUDED,
								Deleted:   false,
							},
							{
								SegmentId: "segment-id",
								UserId:    "user-id-2",
								State:     featureproto.SegmentUser_INCLUDED,
								Deleted:   false,
							},
						},
					}, nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "user-id-1"}, EnvironmentId: "ns0", Tag: "ios"},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							VariationId: vID4,
							Reason: &featureproto.Reason{
								Type:   featureproto.Reason_RULE,
								RuleId: "rule-1",
							},
						},
					},
				},
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success: get from cache and filter by tag: return empty",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: newUUID(t),
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID2,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"segment-id",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID1,
									},
								},
								Tags: []string{"android"},
							},
							{
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID2,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"segment-id",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID1,
									},
								},
								Tags: []string{"ios"},
							},
						},
					}, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&featureproto.SegmentUsers{
						SegmentId: "segment-id",
						Users: []*featureproto.SegmentUser{
							{
								SegmentId: "segment-id",
								UserId:    "user-id-1",
								State:     featureproto.SegmentUser_INCLUDED,
								Deleted:   false,
							},
							{
								SegmentId: "segment-id",
								UserId:    "user-id-2",
								State:     featureproto.SegmentUser_INCLUDED,
								Deleted:   false,
							},
						},
					}, nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "user-id-1"}, EnvironmentId: "ns0", Tag: "web"},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{},
				},
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success: get features from storage",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, errors.New("error"))
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{}, 0, int64(0), nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentId: "ns0", Tag: "android"},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{},
				},
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "fail: return errInternal when getting segment users",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID2,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"id-0",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID1,
									},
								},
								Tags: []string{"android"},
							},
						}}, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "random error"))
				s.segmentUserStorage.(*mock.MockSegmentUserStorage).EXPECT().ListSegmentUsers(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, 0, pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error"))
			},
			input:    &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentId: "ns0", Tag: "android"},
			expected: nil,
			getExpectedErr: func() error {
				return api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "random error")).Err()
			},
		},
		{
			desc:    "success: get users from storage",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID2,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"id-0",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID2,
									},
								},
								Tags: []string{"android"},
							},
						}}, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("random error"))
				s.segmentUserStorage.(*mock.MockSegmentUserStorage).EXPECT().ListSegmentUsers(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.SegmentUser{}, 0, nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentId: "ns0", Tag: "android"},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							VariationId: vID2,
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
				},
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success with viewer account",
			context: createContextWithTokenRoleUnassigned(),
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, errors.New("error"))
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{}, 0, int64(0), nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentId: "ns0", Tag: "android"},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{},
				},
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:     "errPermissionDenied",
			context:  createContextWithTokenRoleUnassigned(),
			service:  createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:    func(s *FeatureService) {},
			input:    &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentId: "ns0", Tag: "android"},
			expected: nil,
			getExpectedErr: func() error {
				return statusPermissionDenied.Err()
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := p.context
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"ja"},
			})

			service := p.service
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.EvaluateFeatures(ctx, p.input)
			if err == nil {
				if len(resp.UserEvaluations.Evaluations) == 1 {
					assert.Equal(t, p.expected.UserEvaluations.Evaluations[0].VariationId, resp.UserEvaluations.Evaluations[0].VariationId, p.desc)
					assert.Equal(t, p.expected.UserEvaluations.Evaluations[0].Reason, resp.UserEvaluations.Evaluations[0].Reason, p.desc)
				} else {
					assert.Equal(t, p.expected.UserEvaluations.Evaluations, resp.UserEvaluations.Evaluations, p.desc)
				}
			} else {
				assert.Equal(t, p.expected, resp, p.desc)
			}
			assert.Equal(t, p.getExpectedErr(), err, p.desc)
		})
	}
}

func TestDebugEvaluateFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	vID1 := newUUID(t)
	vID2 := newUUID(t)
	vID3 := newUUID(t)
	vID4 := newUUID(t)

	patterns := []struct {
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(*FeatureService)
		input          *featureproto.DebugEvaluateFeaturesRequest
		expected       *featureproto.DebugEvaluateFeaturesResponse
		getExpectedErr func() error
	}{
		{
			desc:     "fail: ErrMissingUser",
			context:  createContextWithToken(),
			service:  createFeatureService(mockController),
			setup:    nil,
			input:    &featureproto.DebugEvaluateFeaturesRequest{},
			expected: nil,
			getExpectedErr: func() error {
				return statusMissingUser.Err()
			},
		},
		{
			desc:    "fail: ErrMissingUserID",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup:   nil,
			input: &featureproto.DebugEvaluateFeaturesRequest{
				Users: []*userproto.User{
					{},
				},
			},
			expected: nil,
			getExpectedErr: func() error {
				return statusMissingUserID.Err()
			},
		},
		{
			desc:    "fail: return errInternal when getting features",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error"))
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error"))
			},
			input: &featureproto.DebugEvaluateFeaturesRequest{
				FeatureIds:    []string{"feature-id-1", "feature-id-2"},
				Users:         []*userproto.User{{Id: "test-id"}},
				EnvironmentId: "ns0",
			},
			expected: nil,
			getExpectedErr: func() error {
				return api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error")).Err()
			},
		},
		{
			desc:    "success: get from cache",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: "feature-id-1",
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID2,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"segment-id",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID1,
									},
								},
								Tags: []string{"android"},
							},
							{
								Id: "feature-id-2",
								Variations: []*featureproto.Variation{
									{
										Id:    vID3,
										Value: "true",
									},
									{
										Id:    vID4,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID4,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"segment-id",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID3,
									},
								},
								Tags: []string{"ios"},
							},
						},
					}, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&featureproto.SegmentUsers{
						SegmentId: "segment-id",
						Users: []*featureproto.SegmentUser{
							{
								SegmentId: "segment-id",
								UserId:    "user-id-1",
								State:     featureproto.SegmentUser_INCLUDED,
								Deleted:   false,
							},
							{
								SegmentId: "segment-id",
								UserId:    "user-id-2",
								State:     featureproto.SegmentUser_INCLUDED,
								Deleted:   false,
							},
						},
					}, nil)
			},
			input: &featureproto.DebugEvaluateFeaturesRequest{
				FeatureIds:    []string{"feature-id-1", "feature-id-2"},
				Users:         []*userproto.User{{Id: "user-id-1"}},
				EnvironmentId: "ns0",
			},
			expected: &featureproto.DebugEvaluateFeaturesResponse{
				Evaluations: []*featureproto.Evaluation{
					{
						VariationId: vID2,
						FeatureId:   "feature-id-1",
						UserId:      "user-id-1",
						Reason: &featureproto.Reason{
							Type:   featureproto.Reason_RULE,
							RuleId: "rule-1",
						},
					},
					{
						VariationId: vID4,
						FeatureId:   "feature-id-2",
						UserId:      "user-id-1",
						Reason: &featureproto.Reason{
							Type:   featureproto.Reason_RULE,
							RuleId: "rule-1",
						},
					},
				},
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "fail: return errInternal when getting segment users",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: "feature-id-1",
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID2,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"id-0",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID1,
									},
								},
								Tags: []string{"android"},
							},
						}}, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "random error"))
				s.segmentUserStorage.(*mock.MockSegmentUserStorage).EXPECT().ListSegmentUsers(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, 0, pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error"))
			},
			input: &featureproto.DebugEvaluateFeaturesRequest{
				FeatureIds:    []string{"feature-id-1"},
				Users:         []*userproto.User{{Id: "test-id"}},
				EnvironmentId: "ns0",
			},
			expected: nil,
			getExpectedErr: func() error {
				return api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "random error")).Err()
			},
		},
		{
			desc:    "success: get users from storage",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: "feature-id-1",
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID2,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"id-0",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID2,
									},
								},
								Tags: []string{"android"},
							},
						}}, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("random error"))
				s.segmentUserStorage.(*mock.MockSegmentUserStorage).EXPECT().ListSegmentUsers(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.SegmentUser{}, 0, nil)
			},
			input: &featureproto.DebugEvaluateFeaturesRequest{
				FeatureIds:    []string{"feature-id-1"},
				Users:         []*userproto.User{{Id: "test-id"}},
				EnvironmentId: "ns0",
			},
			expected: &featureproto.DebugEvaluateFeaturesResponse{
				Evaluations: []*featureproto.Evaluation{
					{
						FeatureId:   "feature-id-1",
						UserId:      "test-id",
						VariationId: vID2,
						Reason: &featureproto.Reason{
							Type: featureproto.Reason_DEFAULT,
						},
					},
				},
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success with viewer account",
			context: createContextWithTokenRoleUnassigned(),
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, errors.New("error"))
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					{
						Id: "feature-id-1",
						Variations: []*featureproto.Variation{
							{
								Id:    vID1,
								Value: "true",
							},
							{
								Id:    vID2,
								Value: "false",
							},
						},
						DefaultStrategy: &featureproto.Strategy{
							Type: featureproto.Strategy_FIXED,
							FixedStrategy: &featureproto.FixedStrategy{
								Variation: vID2,
							},
						},
						Tags: []string{"android"},
					},
				}, 0, int64(0), nil)
			},
			input: &featureproto.DebugEvaluateFeaturesRequest{
				FeatureIds:    []string{"feature-id-1"},
				Users:         []*userproto.User{{Id: "test-id"}},
				EnvironmentId: "ns0",
			},
			expected: &featureproto.DebugEvaluateFeaturesResponse{
				Evaluations: []*featureproto.Evaluation{
					{
						FeatureId:   "feature-id-1",
						UserId:      "test-id",
						VariationId: vID2,
						Reason: &featureproto.Reason{
							Type: featureproto.Reason_DEFAULT,
						},
					},
				},
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			context: createContextWithTokenRoleUnassigned(),
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:   func(s *FeatureService) {},
			input: &featureproto.DebugEvaluateFeaturesRequest{
				Users:         []*userproto.User{{Id: "test-id"}},
				EnvironmentId: "ns0",
			},
			expected: nil,
			getExpectedErr: func() error {
				return statusPermissionDenied.Err()
			},
		},
		{
			desc:    "success evaluate multiple users",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: "feature-id-1",
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Targets: []*featureproto.Target{
									{
										Variation: vID1,
										Users:     []string{"test-id-1"},
									},
									{
										Variation: vID2,
										Users:     []string{"test-id-2"},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID2,
									},
								},
								Tags: []string{"android"},
							},
							{
								Id: "feature-id-2",
								Variations: []*featureproto.Variation{
									{
										Id:    vID3,
										Value: "true",
									},
									{
										Id:    vID4,
										Value: "false",
									},
								},
								Targets: []*featureproto.Target{
									{
										Variation: vID3,
										Users:     []string{"test-id-1"},
									},
									{
										Variation: vID4,
										Users:     []string{"test-id-2"},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID3,
									},
								},
								Tags: []string{"ios"},
							},
						}}, nil)
			},
			input: &featureproto.DebugEvaluateFeaturesRequest{
				FeatureIds:    []string{"feature-id-1", "feature-id-2"},
				Users:         []*userproto.User{{Id: "test-id-1"}, {Id: "test-id-2"}},
				EnvironmentId: "ns0",
			},
			expected: &featureproto.DebugEvaluateFeaturesResponse{
				Evaluations: []*featureproto.Evaluation{
					{
						FeatureId:   "feature-id-1",
						UserId:      "test-id-1",
						VariationId: vID1,
						Reason: &featureproto.Reason{
							Type: featureproto.Reason_TARGET,
						},
					},
					{
						FeatureId:   "feature-id-2",
						UserId:      "test-id-1",
						VariationId: vID3,
						Reason: &featureproto.Reason{
							Type: featureproto.Reason_TARGET,
						},
					},
					{
						FeatureId:   "feature-id-1",
						UserId:      "test-id-2",
						VariationId: vID2,
						Reason: &featureproto.Reason{
							Type: featureproto.Reason_TARGET,
						},
					},
					{
						FeatureId:   "feature-id-2",
						UserId:      "test-id-2",
						VariationId: vID4,
						Reason: &featureproto.Reason{
							Type: featureproto.Reason_TARGET,
						},
					},
				},
			},
			getExpectedErr: func() error {
				return nil
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := p.context
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"ja"},
			})

			service := p.service
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.DebugEvaluateFeatures(ctx, p.input)
			if err != nil {
				assert.Equal(t, p.getExpectedErr(), err, p.desc)
				return
			}

			assert.Equal(t, len(p.expected.Evaluations), len(resp.Evaluations))
			for i := 0; i < len(resp.Evaluations); i++ {
				assert.Equal(t, p.expected.Evaluations[i].VariationId, resp.Evaluations[i].VariationId)
				assert.Equal(t, p.expected.Evaluations[i].Reason, resp.Evaluations[i].Reason)
				assert.Equal(t, p.expected.Evaluations[i].FeatureId, resp.Evaluations[i].FeatureId)
				assert.Equal(t, p.expected.Evaluations[i].UserId, resp.Evaluations[i].UserId)
			}
		})
	}
}

func TestEvaluateSingleFeature(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	vID1 := newUUID(t)
	vID2 := newUUID(t)
	vID3 := newUUID(t)
	vID4 := newUUID(t)

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		desc     string
		setup    func(*FeatureService)
		input    *featureproto.EvaluateFeaturesRequest
		expected *featureproto.EvaluateFeaturesResponse
	}{
		{
			desc: "success: evaluate single feature",
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: "fid-1",
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID2,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"id-0",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID2,
									},
								},
								Tags: []string{"android"},
							},
							{
								Id: "fid-2",
								Variations: []*featureproto.Variation{
									{
										Id:    vID3,
										Value: "true",
									},
									{
										Id:    vID4,
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: vID4,
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"id-0",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID4,
									},
								},
								Tags: []string{"android"},
							},
						}}, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("random error"))
				s.segmentUserStorage.(*mock.MockSegmentUserStorage).EXPECT().ListSegmentUsers(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.SegmentUser{}, 0, nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{
				User:          &userproto.User{Id: "user-id"},
				EnvironmentId: "ns0",
				Tag:           "android",
				FeatureId:     "fid-2",
			},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							FeatureId:   "fid-2",
							VariationId: vID4,
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
				},
			},
		},
		{
			desc: "success: evaluate single feature with prerequisite",
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id:      "fid-1",
								Enabled: true,
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Value: "true",
									},
									{
										Id:    vID2,
										Value: "false",
									},
								},
								Prerequisites: []*featureproto.Prerequisite{
									{
										FeatureId:   "fid-2",
										VariationId: vID4,
									},
								},
								Targets: []*featureproto.Target{
									{
										Variation: vID1,
										Users:     []string{"user-id"},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID2,
									},
								},
								Tags:         []string{"android"},
								OffVariation: vID2,
							},
							{
								Id:      "fid-2",
								Enabled: true,
								Variations: []*featureproto.Variation{
									{
										Id:    vID3,
										Value: "true",
									},
									{
										Id:    vID4,
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID4,
									},
								},
								OffVariation: vID3,
								Tags:         []string{"android"},
							},
						}}, nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{
				User:          &userproto.User{Id: "user-id"},
				EnvironmentId: "ns0",
				Tag:           "android",
				FeatureId:     "fid-1",
			},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							FeatureId:   "fid-1",
							VariationId: vID1,
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_TARGET,
							},
						},
					},
				},
			},
		},
	}
	for _, p := range patterns {
		service := createFeatureService(mockController)
		if p.setup != nil {
			p.setup(service)
		}
		resp, _ := service.EvaluateFeatures(ctx, p.input)
		assert.True(t, len(resp.UserEvaluations.Evaluations) == 1)
		assert.Equal(t, p.input.FeatureId, p.expected.UserEvaluations.Evaluations[0].FeatureId, p.desc)
		assert.Equal(t, p.expected.UserEvaluations.Evaluations[0].VariationId, resp.UserEvaluations.Evaluations[0].VariationId, p.desc)
		assert.Equal(t, p.expected.UserEvaluations.Evaluations[0].Reason, resp.UserEvaluations.Evaluations[0].Reason, p.desc)
	}
}

func TestListEnabledFeaturesMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(*FeatureService)
		getExpectedErr func() error
	}{
		{
			desc:    "success",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success with Viewer Account",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: createContextWithTokenRoleUnassigned(),
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: createContextWithTokenRoleUnassigned(),
			setup:   func(s *FeatureService) {},
			getExpectedErr: func() error {
				return statusPermissionDenied.Err()
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
			req := &featureproto.ListEnabledFeaturesRequest{
				EnvironmentId: "ns0",
			}

			_, err := fs.ListEnabledFeatures(ctx, req)
			assert.Equal(t, p.getExpectedErr(), err)
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	service := createFeatureService(mockController)
	patterns := []struct {
		desc     string
		action   func(context.Context, *FeatureService) error
		expected error
	}{
		{
			desc: "GetFeature",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.GetFeature(ctx, &featureproto.GetFeatureRequest{})
				return err
			},
			expected: statusUnauthenticated.Err(),
		},
		{
			desc: "GetFeatures",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.GetFeatures(ctx, &featureproto.GetFeaturesRequest{})
				return err
			},
			expected: statusUnauthenticated.Err(),
		},
		{
			desc: "ListFeatures",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.ListFeatures(ctx, &featureproto.ListFeaturesRequest{})
				return err
			},
			expected: statusUnauthenticated.Err(),
		},
		{
			desc: "ListFeaturesEnabled",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.ListEnabledFeatures(ctx, &featureproto.ListEnabledFeaturesRequest{})
				return err
			},
			expected: statusUnauthenticated.Err(),
		},
		{
			desc: "EvaluateFeatures",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.EvaluateFeatures(ctx, &featureproto.EvaluateFeaturesRequest{})
				return err
			},
			expected: statusUnauthenticated.Err(),
		},
	}
	for _, p := range patterns {
		actual := p.action(ctx, service)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
	}
}

func TestPermissionDenied(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleUnassigned()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	// Use a service with unassigned roles instead of admin
	service := createFeatureServiceWithGetAccountByEnvironmentMock(
		mockController,
		accountproto.AccountV2_Role_Organization_UNASSIGNED,
		accountproto.AccountV2_Role_Environment_UNASSIGNED,
	)
	patterns := []struct {
		desc     string
		action   func(context.Context, *FeatureService) error
		expected error
	}{
		{
			desc: "CreateFeature",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.CreateFeature(ctx, &featureproto.CreateFeatureRequest{})
				return err
			},
			expected: statusPermissionDenied.Err(),
		},
		{
			desc: "DeleteFeature",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.DeleteFeature(ctx, &featureproto.DeleteFeatureRequest{
					Id: "id",
				})
				return err
			},
			expected: statusPermissionDenied.Err(),
		},
		{
			desc: "CloneFeature",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.CloneFeature(ctx, &featureproto.CloneFeatureRequest{
					Id:            "id",
					EnvironmentId: "ns0",
				})
				return err
			},
			expected: statusPermissionDenied.Err(),
		},
	}
	for _, p := range patterns {
		actual := p.action(ctx, service)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
	}
}

func TestDeleteFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*FeatureService)
		req         *featureproto.DeleteFeatureRequest
		expectedErr error
	}{
		{
			desc:  "error: statusMissingID",
			setup: nil,
			req: &featureproto.DeleteFeatureRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: statusMissingID.Err(),
		},
		{
			desc: "error: statusFeatureNotFound",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrFeatureNotFound)
				s.environmentClient.(*envclientmock.MockClient).EXPECT().GetEnvironmentV2(gomock.Any(), gomock.Any()).Return(
					&envproto.GetEnvironmentV2Response{
						Environment: &envproto.EnvironmentV2{},
					},
					nil,
				)
			},
			req: &featureproto.DeleteFeatureRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
			},
			expectedErr: statusFeatureNotFound.Err(),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.batchClient.(*btclientmock.MockClient).EXPECT().ExecuteBatchJob(gomock.Any(), gomock.Any())
				s.environmentClient.(*envclientmock.MockClient).EXPECT().GetEnvironmentV2(
					gomock.Any(),
					&envproto.GetEnvironmentV2Request{
						Id: "ns0",
					},
				).Return(
					&envproto.GetEnvironmentV2Response{
						Environment: &envproto.EnvironmentV2{
							RequireComment: true,
						},
					},
					nil,
				)
			},
			req: &featureproto.DeleteFeatureRequest{
				Id:            "id-1",
				Comment:       "test comment",
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createFeatureService(mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DeleteFeature(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCloneFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*FeatureService)
		req         *featureproto.CloneFeatureRequest
		expectedErr error
	}{
		{
			desc:  "error: statusMissingID",
			setup: nil,
			req: &featureproto.CloneFeatureRequest{
				Id: "",
			},
			expectedErr: statusMissingID.Err(),
		},
		{
			desc:  "error: statusIncorrectDestinationEnvironment",
			setup: nil,
			req: &featureproto.CloneFeatureRequest{
				Id:                  "id-0",
				TargetEnvironmentId: "ns0",
				EnvironmentId:       "ns0",
			},
			expectedErr: statusIncorrectDestinationEnvironment.Err(),
		},
		{
			desc: "error: statusAlreadyExists",
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2fs.ErrFeatureAlreadyExists)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2fs.ErrFeatureAlreadyExists)
			},
			req: &featureproto.CloneFeatureRequest{
				Id:                  "id-0",
				TargetEnvironmentId: "ns1",
				EnvironmentId:       "ns0",
			},
			expectedErr: statusAlreadyExists.Err(),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				s.batchClient.(*btclientmock.MockClient).EXPECT().ExecuteBatchJob(gomock.Any(), gomock.Any())
			},
			req: &featureproto.CloneFeatureRequest{
				Id:                  "id-0",
				TargetEnvironmentId: "ns1",
				EnvironmentId:       "ns0",
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createFeatureService(mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CloneFeature(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetTargetFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	multipleFs := []*featureproto.Feature{
		{
			Id: "fid3",
		},
		{
			Id: "fid2",
		},
		{
			Id: "fid10",
		},
		{
			Id: "fid",
		},
	}
	multiplePreFs := []*featureproto.Feature{
		{
			Id: "fid3",
			Prerequisites: []*featureproto.Prerequisite{
				{
					FeatureId: "fid10",
				},
			},
		},
		{
			Id: "fid2",
		},
		{
			Id: "fid10",
		},
		{
			Id: "fid",
			Prerequisites: []*featureproto.Prerequisite{
				{
					FeatureId: "fid3",
				},
			},
		},
	}
	patterns := []struct {
		desc        string
		fs          []*featureproto.Feature
		id          string
		expected    []*featureproto.Feature
		expectedErr error
	}{
		{
			desc:        "err: feature not found",
			id:          "not_found",
			fs:          multipleFs,
			expected:    nil,
			expectedErr: statusInternal.Err(),
		},
		{
			desc:        "success: feature id is empty",
			id:          "",
			fs:          multipleFs,
			expected:    multipleFs,
			expectedErr: nil,
		},
		{
			desc: "success: prerequisite not configured",
			id:   "fid",
			fs:   multipleFs,
			expected: []*featureproto.Feature{
				multipleFs[3],
			},
			expectedErr: nil,
		},
		{
			desc:        "success: prerequisite configured",
			id:          "fid",
			fs:          multiplePreFs,
			expected:    multiplePreFs,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createFeatureService(mockController)
			actual, err := service.getTargetFeatures(p.fs, p.id)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func makeFeature(id string) *domain.Feature {
	return &domain.Feature{
		Feature: &featureproto.Feature{
			Id:        id,
			Name:      "test feature",
			Version:   1,
			CreatedAt: time.Now().Unix(),
			Variations: []*featureproto.Variation{
				{
					Id:          "variation-A",
					Value:       "A",
					Name:        "Variation A",
					Description: "Thing does A",
				},
				{
					Id:          "variation-B",
					Value:       "B",
					Name:        "Variation B",
					Description: "Thing does B",
				},
			},
			Targets: []*featureproto.Target{
				{
					Variation: "variation-B",
					Users: []string{
						"user1",
					},
				},
			},
			Rules: []*featureproto.Rule{
				{
					Id: "rule-1",
					Strategy: &featureproto.Strategy{
						Type: featureproto.Strategy_FIXED,
						FixedStrategy: &featureproto.FixedStrategy{
							Variation: "variation-A",
						},
					},
					Clauses: []*featureproto.Clause{
						{
							Id:        "clause-1",
							Attribute: "name",
							Operator:  featureproto.Clause_EQUALS,
							Values: []string{
								"user1",
								"user2",
							},
						},
					},
				},
			},
			DefaultStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{
					Variation: "variation-B",
				},
			},
		},
	}
}

func TestValidateEnvironmentSettings(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc     string
		setup    func(*FeatureService)
		env      string
		comment  string
		expected error
	}{
		{
			desc: "error: comment is required",
			setup: func(s *FeatureService) {
				s.environmentClient.(*envclientmock.MockClient).EXPECT().GetEnvironmentV2(
					gomock.Any(),
					&envproto.GetEnvironmentV2Request{
						Id: "env-id",
					},
				).Return(
					&envproto.GetEnvironmentV2Response{
						Environment: &envproto.EnvironmentV2{
							RequireComment: true,
						},
					},
					nil,
				)
			},
			env:      "env-id",
			comment:  "",
			expected: statusCommentRequiredForUpdating.Err(),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				s.environmentClient.(*envclientmock.MockClient).EXPECT().GetEnvironmentV2(
					gomock.Any(),
					&envproto.GetEnvironmentV2Request{
						Id: "env-id",
					},
				).Return(
					&envproto.GetEnvironmentV2Response{
						Environment: &envproto.EnvironmentV2{
							RequireComment: true,
						},
					},
					nil,
				)
			},
			env:      "env-id",
			comment:  "test comment",
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createFeatureService(mockController)
			if p.setup != nil {
				p.setup(service)
			}
			err := service.validateEnvironmentSettings(ctx, p.env, p.comment)
			assert.Equal(t, p.expected, err)
		})
	}
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func TestUpdateFeature(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []*struct {
		desc        string
		ctx         context.Context
		setup       func(*FeatureService)
		input       *featureproto.UpdateFeatureRequest
		expectedErr error
	}{
		{
			desc:        "fail: checkEnvironmentRole",
			ctx:         context.Background(),
			input:       &featureproto.UpdateFeatureRequest{},
			expectedErr: statusUnauthenticated.Err(),
		},
		{
			desc: "fail: id is empty",
			ctx:  createContextWithToken(),
			input: &featureproto.UpdateFeatureRequest{
				EnvironmentId: "eid",
				Comment:       "comment",
				Name:          wrapperspb.String("name"),
				Description:   wrapperspb.String("desc"),
			},
			expectedErr: statusMissingID.Err(),
		},
		{
			desc: "fail: validateFeatureStatus",
			setup: func(s *FeatureService) {
				s.experimentClient.(*exprclientmock.MockClient).EXPECT().ListExperiments(gomock.Any(), gomock.Any()).Return(
					nil, pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "internal"),
				)
			},
			ctx: createContextWithToken(),
			input: &featureproto.UpdateFeatureRequest{
				EnvironmentId: "eid",
				Comment:       "comment",
				Id:            "fid",
				Name:          wrapperspb.String("name"),
				Description:   wrapperspb.String("desc"),
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "internal")).Err(),
		},
		{
			desc: "fail: validateEnvironmentSettings",
			setup: func(s *FeatureService) {
				s.experimentClient.(*exprclientmock.MockClient).EXPECT().ListExperiments(gomock.Any(), gomock.Any()).Return(
					&exprproto.ListExperimentsResponse{}, nil,
				)
				s.environmentClient.(*envclientmock.MockClient).EXPECT().GetEnvironmentV2(
					gomock.Any(),
					&envproto.GetEnvironmentV2Request{Id: "eid"},
				).Return(nil, pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "internal"))
			},
			ctx: createContextWithToken(),
			input: &featureproto.UpdateFeatureRequest{
				EnvironmentId: "eid",
				Comment:       "comment",
				Id:            "fid",
				Name:          wrapperspb.String("name"),
				Description:   wrapperspb.String("desc"),
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "internal")).Err(),
		},
		{
			desc: "fail: publish domain event",
			setup: func(s *FeatureService) {
				s.experimentClient.(*exprclientmock.MockClient).EXPECT().ListExperiments(gomock.Any(), gomock.Any()).Return(
					&exprproto.ListExperimentsResponse{},
					nil,
				)
				s.environmentClient.(*envclientmock.MockClient).EXPECT().GetEnvironmentV2(
					gomock.Any(),
					&envproto.GetEnvironmentV2Request{Id: "eid"},
				).Return(
					&envproto.GetEnvironmentV2Response{Environment: &envproto.EnvironmentV2{RequireComment: true}},
					nil,
				)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.domainPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(map[string]error{"key": errors.New("internal")})
			},
			ctx: createContextWithToken(),
			input: &featureproto.UpdateFeatureRequest{
				EnvironmentId: "eid",
				Comment:       "comment",
				Id:            "fid",
				Name:          wrapperspb.String("name"),
				Description:   wrapperspb.String("desc"),
			},
			expectedErr: statusInternal.Err(),
		},
		{
			desc: "fail: archive feature with dependencies",
			setup: func(s *FeatureService) {
				targetVID := newUUID(t)
				targetVID2 := newUUID(t)
				dependentVID1 := newUUID(t)
				dependentVID2 := newUUID(t)
				s.experimentClient.(*exprclientmock.MockClient).EXPECT().ListExperiments(gomock.Any(), gomock.Any()).Return(
					&exprproto.ListExperimentsResponse{},
					nil,
				)
				s.environmentClient.(*envclientmock.MockClient).EXPECT().GetEnvironmentV2(
					gomock.Any(),
					&envproto.GetEnvironmentV2Request{Id: "eid"},
				).Return(
					&envproto.GetEnvironmentV2Response{Environment: &envproto.EnvironmentV2{RequireComment: true}},
					nil,
				)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					// The error is expected because another feature depends on the target
					assert.Error(t, err)
				}).Return(statusInvalidArchive.Err())
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					// Target feature to be archived
					{
						Id: "target-feature",
						Variations: []*featureproto.Variation{
							{
								Id:    targetVID,
								Value: "true",
								Name:  "true",
							},
							{
								Id:    targetVID2,
								Value: "false",
								Name:  "false",
							},
						},
						OffVariation: targetVID2,
						DefaultStrategy: &featureproto.Strategy{
							Type: featureproto.Strategy_FIXED,
							FixedStrategy: &featureproto.FixedStrategy{
								Variation: targetVID,
							},
						},
						Tags: []string{"test"},
					},
					// Dependent feature that uses target-feature as a prerequisite
					{
						Id: "dependent-feature",
						Variations: []*featureproto.Variation{
							{
								Id:    dependentVID1,
								Value: "true",
								Name:  "true",
							},
							{
								Id:    dependentVID2,
								Value: "false",
								Name:  "false",
							},
						},
						OffVariation: dependentVID2,
						DefaultStrategy: &featureproto.Strategy{
							Type: featureproto.Strategy_FIXED,
							FixedStrategy: &featureproto.FixedStrategy{
								Variation: dependentVID1,
							},
						},
						// This feature depends on target-feature
						Prerequisites: []*featureproto.Prerequisite{
							{
								FeatureId:   "target-feature",
								VariationId: targetVID,
							},
						},
						Tags: []string{"test"},
					},
				}, 0, int64(0), nil)
			},
			ctx: createContextWithToken(),
			input: &featureproto.UpdateFeatureRequest{
				EnvironmentId: "eid",
				Comment:       "comment",
				Id:            "target-feature",
				Archived:      wrapperspb.Bool(true),
			},
			expectedErr: statusInvalidArchive.Err(),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				vID1 := newUUID(t)
				vID2 := newUUID(t)
				rID := newUUID(t)
				cID := newUUID(t)
				s.experimentClient.(*exprclientmock.MockClient).EXPECT().ListExperiments(gomock.Any(), gomock.Any()).Return(
					&exprproto.ListExperimentsResponse{},
					nil,
				)
				s.environmentClient.(*envclientmock.MockClient).EXPECT().GetEnvironmentV2(
					gomock.Any(),
					&envproto.GetEnvironmentV2Request{Id: "eid"},
				).Return(
					&envproto.GetEnvironmentV2Response{Environment: &envproto.EnvironmentV2{RequireComment: true}},
					nil,
				)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					{
						Id: "fid",
						Variations: []*featureproto.Variation{
							{
								Id:    vID1,
								Value: "true",
								Name:  "true",
							},
							{
								Id:    vID2,
								Value: "false",
								Name:  "false",
							},
						},
						OffVariation: vID2,
						Rules: []*featureproto.Rule{
							{
								Id: rID,
								Strategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID2,
									},
								},
								Clauses: []*featureproto.Clause{
									{
										Id:       cID,
										Operator: featureproto.Clause_SEGMENT,
										Values: []string{
											"segment-id",
										},
									},
								},
							},
						},
						DefaultStrategy: &featureproto.Strategy{
							Type: featureproto.Strategy_FIXED,
							FixedStrategy: &featureproto.FixedStrategy{
								Variation: vID1,
							},
						},
						Tags: []string{"android"},
					},
				}, 0, int64(0), nil)
				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().UpdateFeature(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().UpsertTag(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				s.domainPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				// Same-flag conflict detection runs after update
				s.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage).EXPECT().ListScheduledFlagChanges(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.ScheduledFlagChange{}, 0, int64(0), nil)
				// Cross-flag conflict detection
				s.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage).EXPECT().ListScheduledFlagChanges(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.ScheduledFlagChange{}, 0, int64(0), nil)
				s.batchClient.(*btclientmock.MockClient).EXPECT().ExecuteBatchJob(gomock.Any(), gomock.Any())
			},
			ctx: createContextWithToken(),
			input: &featureproto.UpdateFeatureRequest{
				EnvironmentId: "eid",
				Comment:       "comment",
				Id:            "fid",
				Name:          wrapperspb.String("name"),
				Description:   wrapperspb.String("desc"),
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdateFeature(p.ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
