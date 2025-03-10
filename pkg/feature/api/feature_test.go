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
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2/mock"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	acmock "github.com/bucketeer-io/bucketeer/pkg/autoops/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/command"
	btclientmock "github.com/bucketeer-io/bucketeer/pkg/batch/client/mock"
	cachev3mock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	envclientmock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	exprclientmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	aoproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
	exprproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
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
		input          string
		getExpectedErr func(localizer locale.Localizer) error
	}{
		{
			desc:    "error: id is empty",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			input:   "",
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"), localizer)
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
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			input: "fid",
			getExpectedErr: func(localizer locale.Localizer) error {
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
			input: "fid",
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: createContextWithTokenRoleUnassigned(),
			setup:   func(s *FeatureService) {},
			input:   "fid",
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied), localizer)
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
			req := &featureproto.GetFeatureRequest{
				EnvironmentId: "ns0",
				Id:            p.input,
			}
			localizer := locale.NewLocalizer(ctx)

			_, err := fs.GetFeature(ctx, req)
			assert.Equal(t, p.getExpectedErr(localizer), err)
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
		getExpectedErr func(localizer locale.Localizer) error
	}{
		{
			desc:    "error: id is nil",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			input:   nil,
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusMissingIDs, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ids"), localizer)
			},
		},
		{
			desc:    "error: contains empty id",
			service: createFeatureServiceNew(mockController),
			context: createContextWithToken(),
			input:   []string{"id", ""},
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusMissingIDs, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ids"), localizer)
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
			getExpectedErr: func(localizer locale.Localizer) error {
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
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: createContextWithTokenRoleUnassigned(),
			setup:   func(s *FeatureService) {},
			input:   []string{"fid"},
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied), localizer)
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
			localizer := locale.NewLocalizer(ctx)

			_, err := fs.GetFeatures(ctx, req)
			assert.Equal(t, p.getExpectedErr(localizer), err)
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
		getExpectedErr func(localizer locale.Localizer) error
	}{
		{
			desc:          "error: invalid order by",
			service:       createFeatureService(mockController),
			context:       createContextWithToken(),
			setup:         nil,
			orderBy:       featureproto.ListFeaturesRequest_OrderBy(999),
			hasExperiment: false,
			environmentId: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusInvalidOrderBy, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "order_by"), localizer)
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
			getExpectedErr: func(localizer locale.Localizer) error { return nil },
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
			getExpectedErr: func(localizer locale.Localizer) error {
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
			getExpectedErr: func(localizer locale.Localizer) error { return nil },
		},
		{
			desc:          "errPermissionDenied",
			service:       createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context:       createContextWithTokenRoleUnassigned(),
			setup:         func(s *FeatureService) {},
			orderBy:       featureproto.ListFeaturesRequest_DEFAULT,
			hasExperiment: false,
			environmentId: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied), localizer)
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
			localizer := locale.NewLocalizer(ctx)

			_, err := service.ListFeatures(ctx, req)
			assert.Equal(t, p.getExpectedErr(localizer), err)
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
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

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
			expected:                 createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
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
			expected:                 createError(statusInvalidID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "id")),
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
			expected:                 createError(statusMissingName, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
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
			expected:                 createError(statusMissingFeatureVariations, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variations")),
		},
		{
			setup:         nil,
			id:            "Bucketeer-id-2019",
			name:          "name",
			description:   "error: statusMissingFeatureTags",
			variations:    variations,
			tags:          nil,
			environmentId: "ns0",
			expected:      createError(statusMissingFeatureTags, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tags")),
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
			expected:                 createError(statusMissingDefaultOnVariation, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "default_on_variation")),
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
			expected:                 createError(statusMissingDefaultOffVariation, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "default_off_variation")),
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
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
			expected:                 createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
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
				Command: &featureproto.CreateFeatureCommand{
					Id:                       p.id,
					Name:                     p.name,
					Description:              p.description,
					Variations:               p.variations,
					Tags:                     p.tags,
					DefaultOnVariationIndex:  p.defaultOnVariationIndex,
					DefaultOffVariationIndex: p.defaultOffVariationIndex,
				},
				EnvironmentId: p.environmentId,
			}
			actual, err := service.CreateFeature(ctx, req)
			if p.expected == nil {
				assert.NotNil(t, actual.Feature)
			}
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestCreateFeatureNoCommandMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

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
			expected:                 createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
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
			expected:                 createError(statusInvalidID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "id")),
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
			expected:                 createError(statusMissingName, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
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
			expected:                 createError(statusMissingFeatureVariations, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variations")),
		},
		{
			setup:         nil,
			id:            "Bucketeer-id-2019",
			name:          "name",
			description:   "error: statusMissingFeatureTags",
			variations:    variations,
			tags:          nil,
			environmentId: "ns0",
			expected:      createError(statusMissingFeatureTags, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tags")),
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
			expected:                 createError(statusMissingDefaultOnVariation, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "default_on_variation")),
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
			expected:                 createError(statusMissingDefaultOffVariation, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "default_off_variation")),
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
			expected:                 createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
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
	localizer := locale.NewLocalizer(ctx)
	patterns := []struct {
		setup         func(*FeatureService)
		input         []*featureproto.Feature
		environmentId string
		expected      error
	}{
		{
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
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
		err := fs.setLastUsedInfosToFeatureByChunk(context.Background(), p.input, p.environmentId, localizer)
		assert.Equal(t, p.expected, err)
	}
}

func TestConvUpdateFeatureError(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []struct {
		input       error
		expectedErr error
	}{
		{
			input:       v2fs.ErrFeatureNotFound,
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			input:       v2fs.ErrFeatureUnexpectedAffectedRows,
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			input:       storage.ErrKeyNotFound,
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			input:       domain.ErrAlreadyDisabled,
			expectedErr: createError(statusNothingChange, localizer.MustLocalize(locale.NothingToChange)),
		},
		{
			input:       domain.ErrAlreadyEnabled,
			expectedErr: createError(statusNothingChange, localizer.MustLocalize(locale.NothingToChange)),
		},
		{
			input:       errors.New("test"),
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
	}
	for _, p := range patterns {
		fs := &FeatureService{}
		err := fs.convUpdateFeatureError(p.input, localizer)
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
		getExpectedErr func(localizer locale.Localizer) error
	}{
		{
			desc:     "fail: ErrMissingUser",
			context:  createContextWithToken(),
			service:  createFeatureService(mockController),
			setup:    nil,
			input:    &featureproto.EvaluateFeaturesRequest{},
			expected: nil,
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusMissingUser, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "user"), localizer)
			},
		},
		{
			desc:     "fail: ErrMissingUserID",
			context:  createContextWithToken(),
			service:  createFeatureService(mockController),
			setup:    nil,
			input:    &featureproto.EvaluateFeaturesRequest{User: &userproto.User{}},
			expected: nil,
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusMissingUserID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "user_id"), localizer)
			},
		},
		{
			desc:    "fail: return errInternal when getting features",
			context: createContextWithToken(),
			service: createFeatureService(mockController),
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, errors.New("error"))
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("error"))
			},
			input:    &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentId: "ns0", Tag: "android"},
			expected: nil,
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusInternal, localizer.MustLocalize(locale.InternalServerError), localizer)
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
			getExpectedErr: func(localizer locale.Localizer) error {
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
			getExpectedErr: func(localizer locale.Localizer) error {
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
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{}, 0, int64(0), nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentId: "ns0", Tag: "android"},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{},
				},
			},
			getExpectedErr: func(localizer locale.Localizer) error {
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
					nil, errors.New("random error"))
				s.segmentUserStorage.(*storagemock.MockSegmentUserStorage).EXPECT().ListSegmentUsers(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, 0, errors.New("error"))
			},
			input:    &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentId: "ns0", Tag: "android"},
			expected: nil,
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusInternal, localizer.MustLocalize(locale.InternalServerError), localizer)
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
				s.segmentUserStorage.(*storagemock.MockSegmentUserStorage).EXPECT().ListSegmentUsers(
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
			getExpectedErr: func(localizer locale.Localizer) error {
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
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{}, 0, int64(0), nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentId: "ns0", Tag: "android"},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{},
				},
			},
			getExpectedErr: func(localizer locale.Localizer) error {
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
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied), localizer)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := p.context
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"ja"},
			})
			localizer := locale.NewLocalizer(ctx)

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
			assert.Equal(t, p.getExpectedErr(localizer), err, p.desc)
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
				s.segmentUserStorage.(*storagemock.MockSegmentUserStorage).EXPECT().ListSegmentUsers(
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
		getExpectedErr func(localizer locale.Localizer) error
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
			getExpectedErr: func(localizer locale.Localizer) error {
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
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: createContextWithTokenRoleUnassigned(),
			setup:   func(s *FeatureService) {},
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied), localizer)
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
			localizer := locale.NewLocalizer(ctx)

			_, err := fs.ListEnabledFeatures(ctx, req)
			assert.Equal(t, p.getExpectedErr(localizer), err)
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
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

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
			expected: createError(statusUnauthenticated, localizer.MustLocalize(locale.UnauthenticatedError)),
		},
		{
			desc: "GetFeatures",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.GetFeatures(ctx, &featureproto.GetFeaturesRequest{})
				return err
			},
			expected: createError(statusUnauthenticated, localizer.MustLocalize(locale.UnauthenticatedError)),
		},
		{
			desc: "ListFeatures",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.ListFeatures(ctx, &featureproto.ListFeaturesRequest{})
				return err
			},
			expected: createError(statusUnauthenticated, localizer.MustLocalize(locale.UnauthenticatedError)),
		},
		{
			desc: "ListFeaturesEnabled",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.ListEnabledFeatures(ctx, &featureproto.ListEnabledFeaturesRequest{})
				return err
			},
			expected: createError(statusUnauthenticated, localizer.MustLocalize(locale.UnauthenticatedError)),
		},
		{
			desc: "EvaluateFeatures",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.EvaluateFeatures(ctx, &featureproto.EvaluateFeaturesRequest{})
				return err
			},
			expected: createError(statusUnauthenticated, localizer.MustLocalize(locale.UnauthenticatedError)),
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
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	service := createFeatureService(mockController)
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
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "EnableFeature",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.EnableFeature(ctx, &featureproto.EnableFeatureRequest{
					Id:      "id",
					Command: &featureproto.EnableFeatureCommand{},
				})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "DisableFeature",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.DisableFeature(ctx, &featureproto.DisableFeatureRequest{
					Id:      "id",
					Command: &featureproto.DisableFeatureCommand{},
				})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "UnarchiveFeature",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.UnarchiveFeature(ctx, &featureproto.UnarchiveFeatureRequest{
					Id:      "id",
					Command: &featureproto.UnarchiveFeatureCommand{},
				})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "DeleteFeature",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.DeleteFeature(ctx, &featureproto.DeleteFeatureRequest{
					Id:      "id",
					Command: &featureproto.DeleteFeatureCommand{},
				})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "UpdateFeatureVariations",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.UpdateFeatureVariations(ctx, &featureproto.UpdateFeatureVariationsRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "UpdateFeatureTargeting",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.UpdateFeatureTargeting(ctx, &featureproto.UpdateFeatureTargetingRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "CloneFeature",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.CloneFeature(ctx, &featureproto.CloneFeatureRequest{
					Id: "id",
					Command: &featureproto.CloneFeatureCommand{
						EnvironmentId: "ns1",
					},
					EnvironmentId: "ns0",
				})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
	}
	for _, p := range patterns {
		actual := p.action(ctx, service)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
	}
}

func TestEnableFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*FeatureService)
		req         *featureproto.EnableFeatureRequest
		expectedErr error
	}{
		{
			desc:  "error: statusMissingID",
			setup: nil,
			req: &featureproto.EnableFeatureRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error: statusMissingCommand",
			setup: nil,
			req: &featureproto.EnableFeatureRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc: "error: statusNotFound",
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
			req: &featureproto.EnableFeatureRequest{
				Id:            "id-0",
				Command:       &featureproto.EnableFeatureCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
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
			req: &featureproto.EnableFeatureRequest{
				Id:            "id-1",
				Command:       &featureproto.EnableFeatureCommand{},
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
			_, err := service.EnableFeature(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDisableFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*FeatureService)
		req         *featureproto.DisableFeatureRequest
		expectedErr error
	}{
		{
			desc:  "error: statusMissingID",
			setup: nil,
			req: &featureproto.DisableFeatureRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error: statusMissingCommand",
			setup: nil,
			req: &featureproto.DisableFeatureRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc: "error: statusNotFound",
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
			req: &featureproto.DisableFeatureRequest{
				Id:            "id-0",
				Command:       &featureproto.DisableFeatureCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
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
			req: &featureproto.DisableFeatureRequest{
				Id:            "id-1",
				Command:       &featureproto.DisableFeatureCommand{},
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
			_, err := service.DisableFeature(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})

	}
}

func TestValidateArchiveFeature(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		req         *featureproto.ArchiveFeatureRequest
		expectedErr error
	}{
		{
			req: &featureproto.ArchiveFeatureRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			req: &featureproto.ArchiveFeatureRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			req: &featureproto.ArchiveFeatureRequest{
				Id:            "fID-0",
				EnvironmentId: "ns0",
				Command:       &featureproto.ArchiveFeatureCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		err := validateArchiveFeatureRequest(p.req, localizer)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestUnarchiveFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*FeatureService)
		req         *featureproto.UnarchiveFeatureRequest
		expectedErr error
	}{
		{
			desc:  "error: statusMissingID",
			setup: nil,
			req: &featureproto.UnarchiveFeatureRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error: statusMissingCommand",
			setup: nil,
			req: &featureproto.UnarchiveFeatureRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc: "error: statusNotFound",
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
			req: &featureproto.UnarchiveFeatureRequest{
				Id:            "id-0",
				Command:       &featureproto.UnarchiveFeatureCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
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
			req: &featureproto.UnarchiveFeatureRequest{
				Id:            "id-1",
				Command:       &featureproto.UnarchiveFeatureCommand{},
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
			_, err := service.UnarchiveFeature(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
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
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

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
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error: statusMissingCommand",
			setup: nil,
			req: &featureproto.DeleteFeatureRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc: "error: statusNotFound",
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
				Command:       &featureproto.DeleteFeatureCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
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
				Command:       &featureproto.DeleteFeatureCommand{},
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
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

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
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error: statusMissingCommand",
			setup: nil,
			req: &featureproto.CloneFeatureRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc:  "error: statusIncorrectDestinationEnvironment",
			setup: nil,
			req: &featureproto.CloneFeatureRequest{
				Id: "id-0",
				Command: &featureproto.CloneFeatureCommand{
					EnvironmentId: "ns0",
				},
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusIncorrectDestinationEnvironment, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "environment")),
		},
		{
			desc: "error: statusAlreadyExists",
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrFeatureAlreadyExists)
			},
			req: &featureproto.CloneFeatureRequest{
				Id: "id-0",
				Command: &featureproto.CloneFeatureCommand{
					EnvironmentId: "ns1",
				},
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.batchClient.(*btclientmock.MockClient).EXPECT().ExecuteBatchJob(gomock.Any(), gomock.Any())
			},
			req: &featureproto.CloneFeatureRequest{
				Id: "id-0",
				Command: &featureproto.CloneFeatureCommand{
					EnvironmentId: "ns1",
				},
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
			_, err := service.CloneFeature(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestAddFixedStrategyRule(t *testing.T) {
	t.Parallel()
	f := makeFeature("feature-id")
	rID := newUUID(t)
	vID := f.Variations[0].Id
	expected := &featureproto.Rule{
		Id: rID,
		Strategy: &featureproto.Strategy{
			Type:          featureproto.Strategy_FIXED,
			FixedStrategy: &featureproto.FixedStrategy{Variation: vID},
		},
	}
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []struct {
		fs       []*featureproto.Feature
		rule     *featureproto.Rule
		expected error
	}{
		{
			fs: []*featureproto.Feature{},
			rule: &featureproto.Rule{
				Id: "",
				Strategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{Variation: vID},
				},
			},
			expected: createError(statusMissingRuleID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_id")),
		},
		{
			fs: []*featureproto.Feature{},
			rule: &featureproto.Rule{
				Id: "rule-id",
				Strategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{Variation: vID},
				},
			},
			expected: createError(statusIncorrectUUIDFormat, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "id")),
		},
		{
			fs: []*featureproto.Feature{},
			rule: &featureproto.Rule{
				Id:       rID,
				Strategy: nil,
			},
			expected: createError(statusMissingRuleStrategy, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_strategy")),
		},
		{
			fs: []*featureproto.Feature{},
			rule: &featureproto.Rule{
				Id: rID,
				Strategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{},
				},
			},
			expected: createError(statusMissingVariationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			fs: []*featureproto.Feature{
				f.Feature,
				{Id: "feature-1",
					Prerequisites: []*featureproto.Prerequisite{
						{FeatureId: "feature-id"},
					}},
			},
			rule: &featureproto.Rule{
				Id: rID,
				Clauses: []*featureproto.Clause{
					{Operator: featureproto.Clause_FEATURE_FLAG, Attribute: "feature-1"},
				},
				Strategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{},
				},
			},
			expected: createError(statusCycleExists, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "rule")),
		},
		{
			fs:       []*featureproto.Feature{},
			rule:     expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := validateRule(p.fs, f.Feature, p.rule, localizer)
		assert.Equal(t, p.expected, err)
	}
}

func TestAddRolloutStrategyRule(t *testing.T) {
	t.Parallel()
	f := makeFeature("feature-id")
	rID := newUUID(t)
	vID1 := f.Variations[0].Id
	vID2 := f.Variations[1].Id
	expected := &featureproto.Rule{
		Id: rID,
		Strategy: &featureproto.Strategy{
			Type: featureproto.Strategy_ROLLOUT,
			RolloutStrategy: &featureproto.RolloutStrategy{
				Variations: []*featureproto.RolloutStrategy_Variation{
					{
						Variation: vID1,
						Weight:    30000,
					},
					{
						Variation: vID2,
						Weight:    70000,
					},
				},
			},
		},
	}
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []*struct {
		rule     *featureproto.Rule
		expected error
	}{
		{
			rule: &featureproto.Rule{
				Id: "",
				Strategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{
						Variations: []*featureproto.RolloutStrategy_Variation{
							{
								Variation: vID1,
								Weight:    30000,
							},
							{
								Variation: vID2,
								Weight:    70000,
							},
						},
					},
				},
			},
			expected: createError(statusMissingRuleID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_id")),
		},
		{
			rule: &featureproto.Rule{
				Id: "rule-id",
				Strategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{
						Variations: []*featureproto.RolloutStrategy_Variation{
							{
								Variation: vID1,
								Weight:    30000,
							},
							{
								Variation: vID2,
								Weight:    70000,
							},
						},
					},
				},
			},
			expected: createError(statusIncorrectUUIDFormat, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "id")),
		},
		{
			rule: &featureproto.Rule{
				Id:       rID,
				Strategy: nil,
			},
			expected: createError(statusMissingRuleStrategy, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_strategy")),
		},
		{
			rule: &featureproto.Rule{
				Id: rID,
				Strategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{
						Variations: []*featureproto.RolloutStrategy_Variation{
							{
								Variation: vID1,
								Weight:    30000,
							},
						},
					},
				},
			},
			expected: createError(statusDifferentVariationsSize, localizer.MustLocalize(locale.DifferentVariationsSize)),
		},
		{
			rule: &featureproto.Rule{
				Id: rID,
				Strategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{
						Variations: []*featureproto.RolloutStrategy_Variation{
							{
								Variation: "",
								Weight:    30000,
							},
							{
								Variation: vID2,
								Weight:    70000,
							},
						},
					},
				},
			},
			expected: createError(statusMissingVariationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			rule: &featureproto.Rule{
				Id: rID,
				Strategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{
						Variations: []*featureproto.RolloutStrategy_Variation{
							{
								Variation: vID1,
								Weight:    30000,
							},
							{
								Variation: "",
								Weight:    70000,
							},
						},
					},
				},
			},
			expected: createError(statusMissingVariationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			rule: &featureproto.Rule{
				Id: rID,
				Strategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{
						Variations: []*featureproto.RolloutStrategy_Variation{
							{
								Variation: vID1,
								Weight:    -1,
							},
							{
								Variation: vID2,
								Weight:    70000,
							},
						},
					},
				},
			},
			expected: createError(statusIncorrectVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			rule: &featureproto.Rule{
				Id: rID,
				Strategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{
						Variations: []*featureproto.RolloutStrategy_Variation{
							{
								Variation: vID1,
								Weight:    30000,
							},
							{
								Variation: vID2,
								Weight:    -1,
							},
						},
					},
				},
			},
			expected: createError(statusIncorrectVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			rule: &featureproto.Rule{
				Id: rID,
				Strategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{
						Variations: []*featureproto.RolloutStrategy_Variation{
							{
								Variation: vID1,
								Weight:    30000,
							},
							{
								Variation: vID2,
								Weight:    71000,
							},
						},
					},
				},
			},
			expected: createError(statusExceededMaxVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			rule:     expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := validateRule([]*featureproto.Feature{}, f.Feature, p.rule, localizer)
		assert.Equal(t, p.expected, err)
	}
}

func TestChangeRuleToFixedStrategy(t *testing.T) {
	t.Parallel()
	f := makeFeature("feature-id")
	r := f.Rules[0]
	rID := r.Id
	vID := f.Variations[0].Id
	expected := &featureproto.Strategy{
		Type:          featureproto.Strategy_FIXED,
		FixedStrategy: &featureproto.FixedStrategy{Variation: vID},
	}
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []*struct {
		ruleID   string
		strategy *featureproto.Strategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: expected,
			expected: createError(statusMissingRuleID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_id")),
		},
		{
			ruleID:   rID,
			strategy: nil,
			expected: createError(statusMissingRuleStrategy, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_strategy")),
		},
		{
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_FIXED,
			},
			expected: createError(statusMissingFixedStrategy, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "fixed_strategy")),
		},
		{
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type:          featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{},
			},
			expected: createError(statusMissingVariationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			ruleID:   "",
			strategy: nil,
			expected: createError(statusMissingRuleID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_id")),
		},
		{
			ruleID:   rID,
			strategy: expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		cmd := &featureproto.ChangeRuleStrategyCommand{
			RuleId:   p.ruleID,
			Strategy: p.strategy,
		}
		err := validateChangeRuleStrategy(f.Variations, cmd, localizer)
		assert.Equal(t, p.expected, err)
	}
}

func TestChangeRuleToRolloutStrategy(t *testing.T) {
	t.Parallel()
	f := makeFeature("feature-id")
	r := f.Rules[0]
	rID := r.Id
	vID1 := f.Variations[0].Id
	vID2 := f.Variations[1].Id
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	expected := &featureproto.Strategy{
		Type: featureproto.Strategy_ROLLOUT,
		RolloutStrategy: &featureproto.RolloutStrategy{
			Variations: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    30000,
				},
				{
					Variation: vID2,
					Weight:    70000,
				},
			},
		},
	}
	patterns := []struct {
		desc     string
		ruleID   string
		strategy *featureproto.Strategy
		expected error
	}{
		{
			desc:     "fail: errMissingRuleID",
			ruleID:   "",
			strategy: expected,
			expected: createError(statusMissingRuleID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_id")),
		},
		{
			desc:     "fail: errMissingRuleStrategy",
			ruleID:   rID,
			strategy: nil,
			expected: createError(statusMissingRuleStrategy, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_strategy")),
		},
		{
			desc:   "fail: errDifferentVariationsSizeJaJP",
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: vID1,
							Weight:    30000,
						},
					},
				},
			},
			expected: createError(statusDifferentVariationsSize, localizer.MustLocalize(locale.DifferentVariationsSize)),
		},
		{
			desc:   "fail: errMissingVariationIDJaJP: idx-0",
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: "",
							Weight:    30000,
						},
						{
							Variation: vID2,
							Weight:    70000,
						},
					},
				},
			},
			expected: createError(statusMissingVariationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			desc:   "fail: errMissingVariationIDJaJP: idx-1",
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: vID1,
							Weight:    30000,
						},
						{
							Variation: "",
							Weight:    70000,
						},
					},
				},
			},
			expected: createError(statusMissingVariationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			desc:   "fail: errIncorrectVariationWeightJaJP: idx-0",
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: vID1,
							Weight:    -1,
						},
						{
							Variation: vID2,
							Weight:    70000,
						},
					},
				},
			},
			expected: createError(statusIncorrectVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			desc:   "fail: errIncorrectVariationWeightJaJP: idx-1",
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: vID1,
							Weight:    30000,
						},
						{
							Variation: vID2,
							Weight:    -1,
						},
					},
				},
			},
			expected: createError(statusIncorrectVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			desc:   "fail: errIncorrectVariationWeightJaJP: more than total weight",
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: vID1,
							Weight:    30000,
						},
						{
							Variation: vID2,
							Weight:    70001,
						},
					},
				},
			},
			expected: createError(statusExceededMaxVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			desc:   "fail: errIncorrectVariationWeightJaJP: less than total weight",
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: vID1,
							Weight:    29999,
						},
						{
							Variation: vID2,
							Weight:    70000,
						},
					},
				},
			},
			expected: createError(statusExceededMaxVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			desc:     "success",
			ruleID:   rID,
			strategy: expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cmd := &featureproto.ChangeRuleStrategyCommand{
				RuleId:   p.ruleID,
				Strategy: p.strategy,
			}
			err := validateChangeRuleStrategy(f.Variations, cmd, localizer)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestChangeFixedStrategy(t *testing.T) {
	t.Parallel()
	f := makeFeature("feature-id")
	r := f.Rules[0]
	rID := r.Id
	vID := f.Variations[0].Id
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []*struct {
		ruleID   string
		strategy *featureproto.FixedStrategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: &featureproto.FixedStrategy{Variation: vID},
			expected: createError(statusMissingRuleID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_id")),
		},
		{
			ruleID:   rID,
			strategy: nil,
			expected: createError(statusMissingFixedStrategy, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "fixed_strategy")),
		},
		{
			ruleID:   rID,
			strategy: &featureproto.FixedStrategy{},
			expected: createError(statusMissingVariationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			ruleID:   "",
			strategy: nil,
			expected: createError(statusMissingRuleID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_id")),
		},
		{
			ruleID:   rID,
			strategy: &featureproto.FixedStrategy{Variation: vID},
			expected: nil,
		},
	}
	for _, p := range patterns {
		cmd := &featureproto.ChangeFixedStrategyCommand{
			RuleId:   p.ruleID,
			Strategy: p.strategy,
		}
		err := validateChangeFixedStrategy(cmd, localizer)
		assert.Equal(t, p.expected, err)
	}
}

func TestChangeRolloutStrategy(t *testing.T) {
	t.Parallel()
	f := makeFeature("feature-id")
	r := f.Rules[0]
	rID := r.Id
	vID1 := f.Variations[0].Id
	vID2 := f.Variations[1].Id
	expected := &featureproto.RolloutStrategy{Variations: []*featureproto.RolloutStrategy_Variation{
		{
			Variation: vID1,
			Weight:    70000,
		},
		{
			Variation: vID2,
			Weight:    30000,
		},
	}}
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []*struct {
		ruleID   string
		strategy *featureproto.RolloutStrategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: &featureproto.RolloutStrategy{},
			expected: createError(statusMissingRuleID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_id")),
		},
		{
			ruleID:   rID,
			strategy: nil,
			expected: createError(statusMissingRolloutStrategy, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rollout_strategy")),
		},
		{
			ruleID: rID,
			strategy: &featureproto.RolloutStrategy{Variations: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    70000,
				},
			}},
			expected: createError(statusDifferentVariationsSize, localizer.MustLocalize(locale.DifferentVariationsSize)),
		},
		{
			ruleID: rID,
			strategy: &featureproto.RolloutStrategy{Variations: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: "",
					Weight:    70000,
				},
				{
					Variation: vID2,
					Weight:    30000,
				},
			}},
			expected: createError(statusMissingVariationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			ruleID: rID,
			strategy: &featureproto.RolloutStrategy{Variations: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    70000,
				},
				{
					Variation: "",
					Weight:    30000,
				},
			}},
			expected: createError(statusMissingVariationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			ruleID: rID,
			strategy: &featureproto.RolloutStrategy{Variations: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    -1,
				},
				{
					Variation: vID2,
					Weight:    30000,
				},
			}},
			expected: createError(statusIncorrectVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			ruleID: rID,
			strategy: &featureproto.RolloutStrategy{Variations: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    70000,
				},
				{
					Variation: vID2,
					Weight:    -1,
				},
			}},
			expected: createError(statusIncorrectVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			ruleID: rID,
			strategy: &featureproto.RolloutStrategy{Variations: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    62000,
				},
				{
					Variation: vID2,
					Weight:    59000,
				},
			}},
			expected: createError(statusExceededMaxVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			ruleID:   "",
			strategy: nil,
			expected: createError(statusMissingRuleID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_id")),
		},
		{
			ruleID:   rID,
			strategy: expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		cmd := &featureproto.ChangeRolloutStrategyCommand{
			RuleId:   p.ruleID,
			Strategy: p.strategy,
		}
		err := validateChangeRolloutStrategy(f.Variations, cmd, localizer)
		assert.Equal(t, p.expected, err)
	}
}

func TestChangeDefaultStrategy(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	f := makeFeature("feature-id")
	environmentID := "envID"
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []struct {
		desc        string
		setup       func(*FeatureService)
		from        featureproto.UpdateFeatureTargetingRequest_From
		strategy    *featureproto.Strategy
		expectedErr error
	}{
		{
			desc: "err: internal error while getting the progressive rollout list",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{f.Id},
					},
				).Return(nil, errors.New("internal"))
			},
			from:     featureproto.UpdateFeatureTargetingRequest_USER,
			strategy: nil,
			expectedErr: createError(
				statusInternal,
				localizer.MustLocalizeWithTemplate(locale.InternalServerError),
			),
		},
		{
			desc: "err: there is a progressive in progressive",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{f.Id},
					},
				).Return(&aoproto.ListProgressiveRolloutsResponse{
					ProgressiveRollouts: []*aoproto.ProgressiveRollout{
						{
							Status: aoproto.ProgressiveRollout_WAITING,
						},
					},
				}, nil)
			},
			from:     featureproto.UpdateFeatureTargetingRequest_USER,
			strategy: nil,
			expectedErr: createError(
				statusProgressiveRolloutWaitingOrRunningState,
				localizer.MustLocalizeWithTemplate(locale.AutoOpsProgressiveRolloutInProgress),
			),
		},
		{
			desc: "fail: errMissingRuleStrategy",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{f.Id},
					},
				).Return(&aoproto.ListProgressiveRolloutsResponse{
					ProgressiveRollouts: []*aoproto.ProgressiveRollout{
						{
							Status: aoproto.ProgressiveRollout_FINISHED,
						},
					},
				}, nil)
			},
			from:        featureproto.UpdateFeatureTargetingRequest_USER,
			strategy:    nil,
			expectedErr: createError(statusMissingRuleStrategy, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_strategy")),
		},
		{
			desc: "fail: errIncorrectVariationWeightJaJP: more than total weight",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{f.Id},
					},
				).Return(&aoproto.ListProgressiveRolloutsResponse{
					ProgressiveRollouts: []*aoproto.ProgressiveRollout{
						{
							Status: aoproto.ProgressiveRollout_FINISHED,
						},
					},
				}, nil)
			},
			from: featureproto.UpdateFeatureTargetingRequest_USER,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: "variation-A",
							Weight:    30000,
						},
						{
							Variation: "variation-B",
							Weight:    70001,
						},
					},
				},
			},
			expectedErr: createError(statusExceededMaxVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			desc: "fail: errIncorrectVariationWeightJaJP: less than total weight",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{f.Id},
					},
				).Return(&aoproto.ListProgressiveRolloutsResponse{
					ProgressiveRollouts: []*aoproto.ProgressiveRollout{
						{
							Status: aoproto.ProgressiveRollout_FINISHED,
						},
					},
				}, nil)
			},
			from: featureproto.UpdateFeatureTargetingRequest_USER,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: "variation-A",
							Weight:    29999,
						},
						{
							Variation: "variation-B",
							Weight:    70000,
						},
					},
				},
			},
			expectedErr: createError(statusExceededMaxVariationWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_weight")),
		},
		{
			desc: "success: request from user",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{f.Id},
					},
				).Return(&aoproto.ListProgressiveRolloutsResponse{
					ProgressiveRollouts: []*aoproto.ProgressiveRollout{
						{
							Status: aoproto.ProgressiveRollout_FINISHED,
						},
					},
				}, nil)
			},
			from: featureproto.UpdateFeatureTargetingRequest_USER,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: "variation-A",
							Weight:    30000,
						},
						{
							Variation: "variation-B",
							Weight:    70000,
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc:  "success: request from ops",
			setup: func(fs *FeatureService) {},
			from:  featureproto.UpdateFeatureTargetingRequest_OPS,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{
							Variation: "variation-A",
							Weight:    30000,
						},
						{
							Variation: "variation-B",
							Weight:    70000,
						},
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createFeatureServiceNew(mockController)
			p.setup(service)
			cmd := &featureproto.ChangeDefaultStrategyCommand{
				Strategy: p.strategy,
			}
			err := service.validateChangeDefaultStrategy(ctx, p.from, "envID", f.Id, f.Variations, cmd, localizer)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestValidateFeatureVariationsCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	fID0 := "fID-0"
	fID1 := "fID-1"
	fID2 := "fID-2"
	fID3 := "fID-3"
	fID4 := "fID-4"
	fID5 := "fID-5"
	environmentID := "envID"
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	pattens := []*struct {
		desc        string
		setup       func(*FeatureService)
		cmd         command.Command
		fs          []*featureproto.Feature
		expectedErr error
	}{
		{
			desc: "err RemoveVariationCommand: internal error while getting the progressive rollout list",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{fID0},
					},
				).Return(nil, errors.New("internal"))
			},
			cmd: &featureproto.RemoveVariationCommand{
				Id: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
			},
			expectedErr: createError(
				statusInternal,
				localizer.MustLocalizeWithTemplate(locale.InternalServerError),
			),
		},
		{
			desc: "err AddVariationCommand: internal error while getting the progressive rollout list",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{fID0},
					},
				).Return(nil, errors.New("internal"))
			},
			cmd: &featureproto.AddVariationCommand{
				Name: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
			},
			expectedErr: createError(
				statusInternal,
				localizer.MustLocalizeWithTemplate(locale.InternalServerError),
			),
		},
		{
			desc: "err RemoveVariationCommand: there is a progressive in progressive",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{fID0},
					},
				).Return(&aoproto.ListProgressiveRolloutsResponse{
					ProgressiveRollouts: []*aoproto.ProgressiveRollout{
						{
							Status: aoproto.ProgressiveRollout_WAITING,
						},
					},
				}, nil)
			},
			cmd: &featureproto.RemoveVariationCommand{
				Id: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
			},
			expectedErr: createError(
				statusProgressiveRolloutWaitingOrRunningState,
				localizer.MustLocalizeWithTemplate(locale.AutoOpsProgressiveRolloutInProgress),
			),
		},
		{
			desc: "err AddVariationCommand: there is a progressive in progressive",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{fID0},
					},
				).Return(&aoproto.ListProgressiveRolloutsResponse{
					ProgressiveRollouts: []*aoproto.ProgressiveRollout{
						{
							Status: aoproto.ProgressiveRollout_RUNNING,
						},
					},
				}, nil)
			},
			cmd: &featureproto.AddVariationCommand{
				Name: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
			},
			expectedErr: createError(
				statusProgressiveRolloutWaitingOrRunningState,
				localizer.MustLocalizeWithTemplate(locale.AutoOpsProgressiveRolloutInProgress),
			),
		},
		{
			desc: "success: do nothing",
			cmd:  &featureproto.CreateFeatureCommand{},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "err: statusInvalidChangingVariation",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{fID0},
					},
				).Return(&aoproto.ListProgressiveRolloutsResponse{
					ProgressiveRollouts: []*aoproto.ProgressiveRollout{
						{
							Status: aoproto.ProgressiveRollout_FINISHED,
						},
					},
				}, nil)
			},
			cmd: &featureproto.RemoveVariationCommand{
				Id: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
				{
					Id: fID1,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId:   fID0,
							VariationId: "variation-A",
						},
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID2,
				},
				{
					Id: fID3,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID4,
						},
						{
							FeatureId: fID5,
						},
					},
				},
				{
					Id: fID4,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID0,
						},
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID5,
				},
			},
			expectedErr: createError(statusInvalidChangingVariation, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation")),
		},
		{
			desc: "success: RemoveVariationCommand",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{fID0},
					},
				).Return(&aoproto.ListProgressiveRolloutsResponse{
					ProgressiveRollouts: []*aoproto.ProgressiveRollout{
						{
							Status: aoproto.ProgressiveRollout_FINISHED,
						},
					},
				}, nil)
			},
			cmd: &featureproto.RemoveVariationCommand{
				Id: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
				{
					Id: fID1,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID2,
				},
				{
					Id: fID3,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID4,
						},
						{
							FeatureId: fID5,
						},
					},
				},
				{
					Id: fID4,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID5,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: AddVariationCommand",
			setup: func(fs *FeatureService) {
				fs.autoOpsClient.(*acmock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&aoproto.ListProgressiveRolloutsRequest{
						EnvironmentId: environmentID,
						PageSize:      listRequestSize,
						Cursor:        "",
						FeatureIds:    []string{fID0},
					},
				).Return(&aoproto.ListProgressiveRolloutsResponse{
					ProgressiveRollouts: []*aoproto.ProgressiveRollout{
						{
							Status: aoproto.ProgressiveRollout_FINISHED,
						},
					},
				}, nil)
			},
			cmd: &featureproto.AddVariationCommand{
				Name: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
				{
					Id: fID1,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID2,
				},
				{
					Id: fID3,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID4,
						},
						{
							FeatureId: fID5,
						},
					},
				},
				{
					Id: fID4,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID5,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range pattens {
		service := createFeatureServiceNew(mockController)
		if p.setup != nil {
			p.setup(service)
		}
		err := service.validateFeatureVariationsCommand(ctx, p.fs, "envID", &featureproto.Feature{Id: fID0}, p.cmd, localizer)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestValidateAddPrerequisite(t *testing.T) {
	t.Parallel()
	fID0 := "fID-0"
	fID1 := "fID-1"
	fID2 := "fID-2"
	fID3 := "fID-3"
	fID4 := "fID-4"
	fID5 := "fID-5"
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	pattens := []*struct {
		prerequisite *featureproto.Prerequisite
		fs           []*featureproto.Feature
		expectedErr  error
	}{
		{
			prerequisite: &featureproto.Prerequisite{
				FeatureId:   fID1,
				VariationId: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id:            fID0,
					Prerequisites: []*featureproto.Prerequisite{},
				},
				{
					Id: fID1,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID0,
						},
						{
							FeatureId: fID2,
						},
					},
					Variations: []*featureproto.Variation{
						{
							Id: "variation-A",
						},
					},
				},
				{
					Id:            fID2,
					Prerequisites: []*featureproto.Prerequisite{},
				},
				{
					Id: fID3,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID4,
						},
						{
							FeatureId: fID5,
						},
					},
				},
				{
					Id: fID4,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID0,
						},
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id:            fID5,
					Prerequisites: []*featureproto.Prerequisite{},
				},
			},
			expectedErr: createError(statusCycleExists, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "prerequisite")),
		},
		{
			prerequisite: &featureproto.Prerequisite{
				FeatureId:   fID1,
				VariationId: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
				{
					Id: fID1,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
					Variations: []*featureproto.Variation{
						{
							Id: "variation-A",
						},
					},
				},
				{
					Id: fID2,
				},
				{
					Id: fID3,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID4,
						},
						{
							FeatureId: fID5,
						},
					},
				},
				{
					Id: fID4,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID5,
				},
			},
			expectedErr: nil,
		},
		{
			prerequisite: &featureproto.Prerequisite{
				FeatureId:   fID0,
				VariationId: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
					Variations: []*featureproto.Variation{
						{
							Id: "variation-A",
						},
					},
				},
				{
					Id: fID1,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID2,
				},
				{
					Id: fID3,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID4,
						},
						{
							FeatureId: fID5,
						},
					},
				},
				{
					Id: fID4,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID5,
				},
			},
			expectedErr: createError(statusInvalidPrerequisite, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "prerequisite")),
		},
		{
			prerequisite: &featureproto.Prerequisite{
				FeatureId:   fID1,
				VariationId: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId:   fID1,
							VariationId: "variation-B",
						},
					},
				},
				{
					Id: fID1,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
					Variations: []*featureproto.Variation{
						{
							Id: "variation-A",
						},
					},
				},
				{
					Id: fID2,
				},
				{
					Id: fID3,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID4,
						},
						{
							FeatureId: fID5,
						},
					},
				},
				{
					Id: fID4,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID5,
				},
			},
			expectedErr: createError(statusInvalidPrerequisite, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "prerequisite")),
		},
		{
			prerequisite: &featureproto.Prerequisite{
				FeatureId:   fID1,
				VariationId: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
				{
					Id: fID1,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
					Variations: []*featureproto.Variation{
						{
							Id: "variation-B",
						},
					},
				},
				{
					Id: fID2,
				},
				{
					Id: fID3,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID4,
						},
						{
							FeatureId: fID5,
						},
					},
				},
				{
					Id: fID4,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID5,
				},
			},
			expectedErr: createError(statusInvalidVariationID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_id")),
		},
	}
	for _, p := range pattens {
		prevPre := p.fs[0].Prerequisites
		err := validateAddPrerequisite(p.fs, p.fs[0], p.prerequisite, localizer)
		if err == nil {
			assert.Equal(t, p.fs[0].Prerequisites, prevPre)
		}
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestValidateChangePrerequisiteVariation(t *testing.T) {
	t.Parallel()
	fID0 := "fID-0"
	fID1 := "fID-1"
	fID2 := "fID-2"
	fID3 := "fID-3"
	fID4 := "fID-4"
	fID5 := "fID-5"
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	pattens := []*struct {
		prerequisite *featureproto.Prerequisite
		fs           []*featureproto.Feature
		expectedErr  error
	}{
		{
			prerequisite: &featureproto.Prerequisite{
				FeatureId:   fID1,
				VariationId: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
				{
					Id: fID1,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
					Variations: []*featureproto.Variation{
						{
							Id: "variation-A",
						},
					},
				},
				{
					Id: fID2,
				},
				{
					Id: fID3,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID4,
						},
						{
							FeatureId: fID5,
						},
					},
				},
				{
					Id: fID4,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID5,
				},
			},
			expectedErr: nil,
		},
		{
			prerequisite: &featureproto.Prerequisite{
				FeatureId:   fID1,
				VariationId: "variation-A",
			},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
				{
					Id: fID1,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
					Variations: []*featureproto.Variation{
						{
							Id: "variation-B",
						},
					},
				},
				{
					Id: fID2,
				},
				{
					Id: fID3,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID4,
						},
						{
							FeatureId: fID5,
						},
					},
				},
				{
					Id: fID4,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: fID2,
						},
					},
				},
				{
					Id: fID5,
				},
			},
			expectedErr: createError(statusInvalidVariationID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_id")),
		},
	}
	for _, p := range pattens {
		err := validateChangePrerequisiteVariation(p.fs, p.prerequisite, localizer)
		assert.Equal(t, p.expectedErr, err)
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
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
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
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
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
			actual, err := service.getTargetFeatures(p.fs, p.id, localizer)
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
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

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
			expected: createError(statusCommentRequiredForUpdating, localizer.MustLocalizeWithTemplate(locale.CommentRequiredForUpdating, "command")),
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
			err := service.validateEnvironmentSettings(ctx, p.env, p.comment, localizer)
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

	ctx := createContextWithToken()
	localizer := locale.NewLocalizer(ctx)
	unauthenticatedErr, _ := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalize(locale.UnauthenticatedError),
	})
	missingIDErr, _ := statusMissingID.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
	})
	internalErr, _ := statusInternal.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalizeWithTemplate(locale.InternalServerError, "id"),
	})

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
			expectedErr: unauthenticatedErr.Err(),
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
			expectedErr: missingIDErr.Err(),
		},
		{
			desc: "fail: validateFeatureStatus",
			setup: func(s *FeatureService) {
				s.experimentClient.(*exprclientmock.MockClient).EXPECT().ListExperiments(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("internal"),
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
			expectedErr: internalErr.Err(),
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
				).Return(nil, errors.New("internal"))
			},
			ctx: createContextWithToken(),
			input: &featureproto.UpdateFeatureRequest{
				EnvironmentId: "eid",
				Comment:       "comment",
				Id:            "fid",
				Name:          wrapperspb.String("name"),
				Description:   wrapperspb.String("desc"),
			},
			expectedErr: internalErr.Err(),
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
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
			expectedErr: internalErr.Err(),
		},
		{
			desc: "success",
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.domainPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
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
		service := createFeatureServiceNew(mockController)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.UpdateFeature(p.ctx, p.input)
		assert.Equal(t, p.expectedErr, err)
	}
}
