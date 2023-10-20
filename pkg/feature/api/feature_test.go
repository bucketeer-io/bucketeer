// Copyright 2023 The Bucketeer Authors.
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
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/command"
	cachev3mock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

func TestGetFeatureMySQL(t *testing.T) {
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
		input    string
		expected error
	}{
		{
			desc:     "error: id is empty",
			input:    "",
			expected: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "success",
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
			input:    "fid",
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			fs := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(fs)
			}
			req := &featureproto.GetFeatureRequest{
				EnvironmentNamespace: "ns0",
				Id:                   p.input,
			}
			_, err := fs.GetFeature(ctx, req)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestGetFeaturesMySQL(t *testing.T) {
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
		input    []string
		expected error
	}{
		{
			desc:     "error: id is nil",
			input:    nil,
			expected: createError(statusMissingIDs, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ids")),
		},
		{
			desc:     "error: contains empty id",
			input:    []string{"id", ""},
			expected: createError(statusMissingIDs, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ids")),
		},
		{
			desc: "success",
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
			input:    []string{"fid"},
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			fs := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(fs)
			}
			req := &featureproto.GetFeaturesRequest{
				EnvironmentNamespace: "ns0",
				Ids:                  p.input,
			}
			_, err := fs.GetFeatures(ctx, req)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestListFeaturesMySQL(t *testing.T) {
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
		setup                func(*FeatureService)
		orderBy              featureproto.ListFeaturesRequest_OrderBy
		hasExperiment        bool
		environmentNamespace string
		expected             error
	}{
		{
			setup:                nil,
			orderBy:              featureproto.ListFeaturesRequest_OrderBy(999),
			hasExperiment:        false,
			environmentNamespace: "ns0",
			expected:             createError(statusInvalidOrderBy, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "order_by")),
		},
		{
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
			orderBy:              featureproto.ListFeaturesRequest_DEFAULT,
			hasExperiment:        false,
			environmentNamespace: "ns0",
			expected:             nil,
		},
		{
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
			orderBy:              featureproto.ListFeaturesRequest_DEFAULT,
			hasExperiment:        true,
			environmentNamespace: "ns0",
			expected:             nil,
		},
	}
	for _, p := range patterns {
		service := createFeatureService(mockController)
		if p.setup != nil {
			p.setup(service)
		}
		req := &featureproto.ListFeaturesRequest{
			OrderBy:              p.orderBy,
			EnvironmentNamespace: "ns0",
		}
		_, err := service.ListFeatures(ctx, req)
		assert.Equal(t, p.expected, err)
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
		environmentNamespace                              string
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
			environmentNamespace:     "ns0",
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
			environmentNamespace:     "ns0",
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
			environmentNamespace:     "ns0",
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
			environmentNamespace:     "ns0",
			expected:                 createError(statusMissingFeatureVariations, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variations")),
		},
		{
			setup:                nil,
			id:                   "Bucketeer-id-2019",
			name:                 "name",
			description:          "error: statusMissingFeatureTags",
			variations:           variations,
			tags:                 nil,
			environmentNamespace: "ns0",
			expected:             createError(statusMissingFeatureTags, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tags")),
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
			environmentNamespace:     "ns0",
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
			environmentNamespace:     "ns0",
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
			environmentNamespace:     "ns0",
			expected:                 createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
					ctx, gomock.Any(), gomock.Any(),
				).Return(row)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(nil)
			},
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "success",
			variations:               variations,
			tags:                     tags,
			defaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
			defaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
			environmentNamespace:     "ns0",
			expected:                 nil,
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
					ctx, gomock.Any(), gomock.Any(),
				).Return(row)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(errors.New("error"))
			},
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "success to create, but fail to refresh cache",
			variations:               variations,
			tags:                     tags,
			defaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
			defaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
			environmentNamespace:     "ns0",
			expected:                 createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
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
				EnvironmentNamespace: p.environmentNamespace,
			}
			_, err := service.CreateFeature(ctx, req)
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
		setup                func(*FeatureService)
		input                []*featureproto.Feature
		environmentNamespace string
		expected             error
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
			environmentNamespace: "ns0",
			expected:             nil,
		},
	}
	for _, p := range patterns {
		fs := createFeatureServiceNew(mockController)
		p.setup(fs)
		err := fs.setLastUsedInfosToFeatureByChunk(context.Background(), p.input, p.environmentNamespace, localizer)
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
		input       *featureproto.EvaluateFeaturesRequest
		expected    *featureproto.EvaluateFeaturesResponse
		expectedErr error
	}{
		{
			desc:        "fail: ErrMissingUser",
			setup:       nil,
			input:       &featureproto.EvaluateFeaturesRequest{},
			expected:    nil,
			expectedErr: createError(statusMissingUser, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "user")),
		},
		{
			desc:        "fail: ErrMissingUserID",
			setup:       nil,
			input:       &featureproto.EvaluateFeaturesRequest{User: &userproto.User{}},
			expected:    nil,
			expectedErr: createError(statusMissingUserID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "user_id")),
		},
		{
			desc: "fail: return errInternal when getting features",
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentNamespace: "ns0", Tag: "android"},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success: get from cache",
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
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "user-id-1"}, EnvironmentNamespace: "ns0", Tag: "ios"},
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
			expectedErr: nil,
		},
		{
			desc: "success: get from cache and filter by tag: return empty",
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
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "user-id-1"}, EnvironmentNamespace: "ns0", Tag: "web"},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: get features from storage",
			setup: func(s *FeatureService) {
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, errors.New("error"))
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
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentNamespace: "ns0", Tag: "android"},
			expected: &featureproto.EvaluateFeaturesResponse{
				UserEvaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "fail: return errInternal when getting segment users",
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentNamespace: "ns0", Tag: "android"},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success: get users from storage",
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
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentNamespace: "ns0", Tag: "android"},
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
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createFeatureService(mockController)
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
		assert.Equal(t, p.expectedErr, err, p.desc)
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
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &featureproto.EvaluateFeaturesRequest{
				User:                 &userproto.User{Id: "user-id"},
				EnvironmentNamespace: "ns0",
				Tag:                  "android",
				FeatureId:            "fid-2",
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
				User:                 &userproto.User{Id: "user-id"},
				EnvironmentNamespace: "ns0",
				Tag:                  "android",
				FeatureId:            "fid-1",
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
						EnvironmentNamespace: "ns1",
					},
					EnvironmentNamespace: "ns0",
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
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error: statusMissingCommand",
			setup: nil,
			req: &featureproto.EnableFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc: "error: statusNotFound",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrFeatureNotFound)
			},
			req: &featureproto.EnableFeatureRequest{
				Id:                   "id-0",
				Command:              &featureproto.EnableFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
					ctx, gomock.Any(), gomock.Any(),
				).Return(row)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(nil)
			},
			req: &featureproto.EnableFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.EnableFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
		{
			desc: "success to enable, but fail to refresh cache",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
					ctx, gomock.Any(), gomock.Any(),
				).Return(row)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(errors.New("error"))
			},
			req: &featureproto.EnableFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.EnableFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
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
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error: statusMissingCommand",
			setup: nil,
			req: &featureproto.DisableFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc: "error: statusNotFound",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrFeatureNotFound)
			},
			req: &featureproto.DisableFeatureRequest{
				Id:                   "id-0",
				Command:              &featureproto.DisableFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
					ctx, gomock.Any(), gomock.Any(),
				).Return(row)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(nil)
			},
			req: &featureproto.DisableFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.DisableFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
		{
			desc: "success to disable, but fail to refresh cache",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
					ctx, gomock.Any(), gomock.Any(),
				).Return(row)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(errors.New("error"))
			},
			req: &featureproto.DisableFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.DisableFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
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
	f0 := makeFeature("fID-0")
	f1 := makeFeature("fID-1")
	f2 := makeFeature("fID-2")
	f3 := makeFeature("fID-3")
	f4 := makeFeature("fID-4")
	f5 := makeFeature("fID-5")
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
		fs          []*featureproto.Feature
		expectedErr error
	}{
		{
			req: &featureproto.ArchiveFeatureRequest{
				EnvironmentNamespace: "ns0",
			},
			fs:          nil,
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			req: &featureproto.ArchiveFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			fs:          nil,
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			req: &featureproto.ArchiveFeatureRequest{
				Id:                   "fID-0",
				EnvironmentNamespace: "ns0",
				Command:              &featureproto.ArchiveFeatureCommand{},
			},
			fs: []*featureproto.Feature{
				{
					Id: f0.Id,
				},
				{
					Id: f1.Id,
				},
				{
					Id: f2.Id,
				},
				{
					Id: f3.Id,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: f4.Id,
						},
						{
							FeatureId: f5.Id,
						},
					},
				},
				{
					Id: f4.Id,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: f0.Id,
						},
						{
							FeatureId: f2.Id,
						},
					},
				},
				{
					Id: f5.Id,
				},
			},
			expectedErr: createError(statusInvalidArchive, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "archive")),
		},
		{
			req: &featureproto.ArchiveFeatureRequest{
				Id:                   "fID-0",
				EnvironmentNamespace: "ns0",
				Command:              &featureproto.ArchiveFeatureCommand{},
			},
			fs: []*featureproto.Feature{
				{
					Id: f0.Id,
				},
				{
					Id: f1.Id,
				},
				{
					Id: f2.Id,
				},
				{
					Id: f3.Id,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: f2.Id,
						},
						{
							FeatureId: f1.Id,
						},
					},
				},
				{
					Id: f4.Id,
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId: f5.Id,
						},
					},
				},
				{
					Id: f5.Id,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		err := validateArchiveFeatureRequest(p.req, p.fs, localizer)
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
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error: statusMissingCommand",
			setup: nil,
			req: &featureproto.UnarchiveFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc: "error: statusNotFound",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrFeatureNotFound)
			},
			req: &featureproto.UnarchiveFeatureRequest{
				Id:                   "id-0",
				Command:              &featureproto.UnarchiveFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
					ctx, gomock.Any(), gomock.Any(),
				).Return(row)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(nil)
			},
			req: &featureproto.UnarchiveFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.UnarchiveFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
		{
			desc: "success to unarchive, but fail to refresh cache",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
					ctx, gomock.Any(), gomock.Any(),
				).Return(row)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(errors.New("error"))
			},
			req: &featureproto.UnarchiveFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.UnarchiveFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
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
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error: statusMissingCommand",
			setup: nil,
			req: &featureproto.DeleteFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc: "error: statusNotFound",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrFeatureNotFound)
			},
			req: &featureproto.DeleteFeatureRequest{
				Id:                   "id-0",
				Command:              &featureproto.DeleteFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
					ctx, gomock.Any(), gomock.Any(),
				).Return(row)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(nil)
			},
			req: &featureproto.DeleteFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.DeleteFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
		{
			desc: "success to delete, but fail to refresh cache",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
					ctx, gomock.Any(), gomock.Any(),
				).Return(row)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(errors.New("error"))
			},
			req: &featureproto.DeleteFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.DeleteFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
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
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc:  "error: statusIncorrectDestinationEnvironment",
			setup: nil,
			req: &featureproto.CloneFeatureRequest{
				Id: "id-0",
				Command: &featureproto.CloneFeatureCommand{
					EnvironmentNamespace: "ns0",
				},
				EnvironmentNamespace: "ns0",
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
					EnvironmentNamespace: "ns1",
				},
				EnvironmentNamespace: "ns0",
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
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				r := mysqlmock.NewMockRow(mockController)
				r.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					ctx, gomock.Any(), gomock.Any(),
				).Return(r)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(nil)
			},
			req: &featureproto.CloneFeatureRequest{
				Id: "id-0",
				Command: &featureproto.CloneFeatureCommand{
					EnvironmentNamespace: "ns1",
				},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
		{
			desc: "success to clone, but fail to refresh cache",
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
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				r := mysqlmock.NewMockRow(mockController)
				r.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					ctx, gomock.Any(), gomock.Any(),
				).Return(r)
				s.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(errors.New("error"))
			},
			req: &featureproto.CloneFeatureRequest{
				Id: "id-0",
				Command: &featureproto.CloneFeatureCommand{
					EnvironmentNamespace: "ns1",
				},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
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
		rule     *featureproto.Rule
		expected error
	}{
		{
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
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{},
				},
			},
			expected: createError(statusMissingVariationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			rule:     expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := validateRule(f.Variations, p.rule, localizer)
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
		err := validateRule(f.Variations, p.rule, localizer)
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
	f := makeFeature("feature-id")
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
		strategy    *featureproto.Strategy
		expectedErr error
	}{
		{
			desc:        "fail: errMissingRuleStrategy",
			strategy:    nil,
			expectedErr: createError(statusMissingRuleStrategy, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "rule_strategy")),
		},
		{
			desc: "fail: errIncorrectVariationWeightJaJP: more than total weight",
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
			desc: "success",
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
			cmd := &featureproto.ChangeDefaultStrategyCommand{
				Strategy: p.strategy,
			}
			err := validateChangeDefaultStrategy(f.Variations, cmd, localizer)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestValidateFeatureVariationsCommand(t *testing.T) {
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
		cmd         command.Command
		fs          []*featureproto.Feature
		expectedErr error
	}{
		{
			cmd: &featureproto.CreateFeatureCommand{},
			fs: []*featureproto.Feature{
				{
					Id: fID0,
				},
			},
			expectedErr: nil,
		},
		{
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
	}
	for _, p := range pattens {
		err := validateFeatureVariationsCommand(p.fs, p.cmd, localizer)
		assert.Equal(t, p.expectedErr, err)
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

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}
