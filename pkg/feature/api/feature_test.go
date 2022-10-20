// Copyright 2022 The Bucketeer Authors.
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

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"

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

	patterns := []struct {
		desc     string
		setup    func(*FeatureService)
		input    string
		expected error
	}{
		{
			desc:     "error: id is empty",
			input:    "",
			expected: errMissingIDJaJP,
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
			ctx := createContextWithToken()
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

	patterns := []struct {
		desc     string
		setup    func(*FeatureService)
		input    []string
		expected error
	}{
		{
			desc:     "error: id is nil",
			input:    nil,
			expected: errMissingIDsJaJP,
		},
		{
			desc:     "error: contains empty id",
			input:    []string{"id", ""},
			expected: errMissingIDsJaJP,
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
			ctx := createContextWithToken()
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
			expected:             errInvalidOrderByJaJP,
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
		ctx := createContextWithToken()
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
			description:              "description",
			variations:               nil,
			tags:                     nil,
			defaultOnVariationIndex:  nil,
			defaultOffVariationIndex: nil,
			environmentNamespace:     "ns0",
			expected:                 errMissingIDJaJP,
		},
		{
			setup:                    nil,
			id:                       "bucketeer_id",
			name:                     "name",
			description:              "description",
			variations:               nil,
			tags:                     nil,
			defaultOnVariationIndex:  nil,
			defaultOffVariationIndex: nil,
			environmentNamespace:     "ns0",
			expected:                 errInvalidIDJaJP,
		},
		{
			setup:                    nil,
			id:                       "Bucketeer-id-2019",
			name:                     "",
			description:              "description",
			variations:               nil,
			tags:                     nil,
			defaultOnVariationIndex:  nil,
			defaultOffVariationIndex: nil,
			environmentNamespace:     "ns0",
			expected:                 errMissingNameJaJP,
		},
		{
			setup:                    nil,
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "description",
			variations:               nil,
			tags:                     nil,
			defaultOnVariationIndex:  nil,
			defaultOffVariationIndex: nil,
			environmentNamespace:     "ns0",
			expected:                 errMissingFeatureVariationsJaJP,
		},
		{
			setup:                nil,
			id:                   "Bucketeer-id-2019",
			name:                 "name",
			description:          "description",
			variations:           variations,
			tags:                 nil,
			environmentNamespace: "ns0",
			expected:             errMissingFeatureTagsJaJP,
		},
		{
			setup:                    nil,
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "description",
			variations:               variations,
			tags:                     tags,
			defaultOnVariationIndex:  nil,
			defaultOffVariationIndex: nil,
			environmentNamespace:     "ns0",
			expected:                 errMissingDefaultOnVariationJaJP,
		},
		{
			setup:                    nil,
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "description",
			variations:               variations,
			tags:                     tags,
			defaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
			defaultOffVariationIndex: nil,
			environmentNamespace:     "ns0",
			expected:                 errMissingDefaultOffVariationJaJP,
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
			description:              "description",
			variations:               variations,
			tags:                     tags,
			defaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
			defaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
			environmentNamespace:     "ns0",
			expected:                 errAlreadyExistsJaJP,
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			id:                       "Bucketeer-id-2019",
			name:                     "name",
			description:              "description",
			variations:               variations,
			tags:                     tags,
			defaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
			defaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
			environmentNamespace:     "ns0",
			expected:                 nil,
		},
	}
	for _, p := range patterns {
		ctx := createContextWithToken()
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
	}
}

func TestSetFeatureToLastUsedInfosByChunk(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
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
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	patterns := []struct {
		input       error
		expectedErr error
	}{
		{
			input:       v2fs.ErrFeatureNotFound,
			expectedErr: errNotFoundJaJP,
		},
		{
			input:       v2fs.ErrFeatureUnexpectedAffectedRows,
			expectedErr: errNotFoundJaJP,
		},
		{
			input:       storage.ErrKeyNotFound,
			expectedErr: errNotFoundJaJP,
		},
		{
			input:       domain.ErrAlreadyDisabled,
			expectedErr: errNothingChangeJaJP,
		},
		{
			input:       domain.ErrAlreadyEnabled,
			expectedErr: errNothingChangeJaJP,
		},
		{
			input:       errors.New("test"),
			expectedErr: errInternalJaJP,
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
			expectedErr: errMissingUserJaJP,
		},
		{
			desc:        "fail: ErrMissingUserID",
			setup:       nil,
			input:       &featureproto.EvaluateFeaturesRequest{User: &userproto.User{}},
			expected:    nil,
			expectedErr: errMissingUserIDJaJP,
		},
		{
			desc:        "fail: ErrMissingFeatureTag",
			setup:       nil,
			input:       &featureproto.EvaluateFeaturesRequest{User: &userproto.User{Id: "test-id"}, EnvironmentNamespace: "ns0"},
			expected:    nil,
			expectedErr: errMissingFeatureTagJaJP,
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
			expectedErr: localizedError(statusInternal, locale.JaJP),
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
			expectedErr: localizedError(statusInternal, locale.JaJP),
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
		ctx := createContextWithToken()
		service := createFeatureService(mockController)
		if p.setup != nil {
			p.setup(service)
		}
		resp, err := service.EvaluateFeatures(ctx, p.input)
		if err == nil {
			if len(resp.UserEvaluations.Evaluations) > 0 {
				assert.Equal(t, p.expected.UserEvaluations.Evaluations[0].VariationId, resp.UserEvaluations.Evaluations[0].VariationId, p.desc)
				assert.Equal(t, p.expected.UserEvaluations.Evaluations[0].Reason, resp.UserEvaluations.Evaluations[0].Reason)
			} else {
				assert.Equal(t, p.expected.UserEvaluations.Evaluations, resp.UserEvaluations.Evaluations, p.desc)
			}
		} else {
			assert.Equal(t, p.expected, resp, p.desc)
		}
		assert.Equal(t, p.expectedErr, err, p.desc)
	}
}

func TestUnauthenticated(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
			expected: errUnauthenticatedJaJP,
		},
		{
			desc: "GetFeatures",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.GetFeatures(ctx, &featureproto.GetFeaturesRequest{})
				return err
			},
			expected: errUnauthenticatedJaJP,
		},
		{
			desc: "ListFeatures",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.ListFeatures(ctx, &featureproto.ListFeaturesRequest{})
				return err
			},
			expected: errUnauthenticatedJaJP,
		},
		{
			desc: "ListFeaturesEnabled",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.ListEnabledFeatures(ctx, &featureproto.ListEnabledFeaturesRequest{})
				return err
			},
			expected: errUnauthenticatedJaJP,
		},
		{
			desc: "EvaluateFeatures",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.EvaluateFeatures(ctx, &featureproto.EvaluateFeaturesRequest{})
				return err
			},
			expected: errUnauthenticatedJaJP,
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
			expected: errPermissionDeniedJaJP,
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
			expected: errPermissionDeniedJaJP,
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
			expected: errPermissionDeniedJaJP,
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
			expected: errPermissionDeniedJaJP,
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
			expected: errPermissionDeniedJaJP,
		},
		{
			desc: "UpdateFeatureVariations",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.UpdateFeatureVariations(ctx, &featureproto.UpdateFeatureVariationsRequest{})
				return err
			},
			expected: errPermissionDeniedJaJP,
		},
		{
			desc: "UpdateFeatureTargeting",
			action: func(ctx context.Context, fs *FeatureService) error {
				_, err := fs.UpdateFeatureTargeting(ctx, &featureproto.UpdateFeatureTargetingRequest{})
				return err
			},
			expected: errPermissionDeniedJaJP,
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
			expected: errPermissionDeniedJaJP,
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

	patterns := []struct {
		setup       func(*FeatureService)
		req         *featureproto.EnableFeatureRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &featureproto.EnableFeatureRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errMissingIDJaJP,
		},
		{
			setup: nil,
			req: &featureproto.EnableFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errMissingCommandJaJP,
		},
		{
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
			expectedErr: errNotFoundJaJP,
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.EnableFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.EnableFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		ctx := createContextWithToken()
		service := createFeatureService(mockController)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.EnableFeature(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestDisableFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(*FeatureService)
		req         *featureproto.DisableFeatureRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &featureproto.DisableFeatureRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errMissingIDJaJP,
		},
		{
			setup: nil,
			req: &featureproto.DisableFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errMissingCommandJaJP,
		},
		{
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
			expectedErr: errNotFoundJaJP,
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.DisableFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.DisableFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		ctx := createContextWithToken()
		service := createFeatureService(mockController)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.DisableFeature(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
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
			expectedErr: errMissingIDJaJP,
		},
		{
			req: &featureproto.ArchiveFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			fs:          nil,
			expectedErr: errMissingCommandJaJP,
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
			expectedErr: localizedError(statusInvalidArchive, locale.JaJP),
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
		err := validateArchiveFeatureRequest(p.req, p.fs)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestUnarchiveFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(*FeatureService)
		req         *featureproto.UnarchiveFeatureRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &featureproto.UnarchiveFeatureRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errMissingIDJaJP,
		},
		{
			setup: nil,
			req: &featureproto.UnarchiveFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errMissingCommandJaJP,
		},
		{
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
			expectedErr: errNotFoundJaJP,
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.UnarchiveFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.UnarchiveFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		ctx := createContextWithToken()
		service := createFeatureService(mockController)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.UnarchiveFeature(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestDeleteFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(*FeatureService)
		req         *featureproto.DeleteFeatureRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &featureproto.DeleteFeatureRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errMissingIDJaJP,
		},
		{
			setup: nil,
			req: &featureproto.DeleteFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errMissingCommandJaJP,
		},
		{
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
			expectedErr: errNotFoundJaJP,
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.DeleteFeatureRequest{
				Id:                   "id-1",
				Command:              &featureproto.DeleteFeatureCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		ctx := createContextWithToken()
		service := createFeatureService(mockController)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.DeleteFeature(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestCloneFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(*FeatureService)
		req         *featureproto.CloneFeatureRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &featureproto.CloneFeatureRequest{
				Id: "",
			},
			expectedErr: errMissingIDJaJP,
		},
		{
			setup: nil,
			req: &featureproto.CloneFeatureRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errMissingCommandJaJP,
		},
		{
			setup: nil,
			req: &featureproto.CloneFeatureRequest{
				Id: "id-0",
				Command: &featureproto.CloneFeatureCommand{
					EnvironmentNamespace: "ns0",
				},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errIncorrectDestinationEnvironmentJaJP,
		},
		{
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
			expectedErr: errAlreadyExistsJaJP,
		},
		{
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
	}

	for _, p := range patterns {
		ctx := createContextWithToken()
		service := createFeatureService(mockController)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.CloneFeature(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
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
			expected: localizedError(statusMissingRuleID, locale.JaJP),
		},
		{
			rule: &featureproto.Rule{
				Id: "rule-id",
				Strategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{Variation: vID},
				},
			},
			expected: localizedError(statusIncorrectUUIDFormat, locale.JaJP),
		},
		{
			rule: &featureproto.Rule{
				Id:       rID,
				Strategy: nil,
			},
			expected: localizedError(statusMissingRuleStrategy, locale.JaJP),
		},
		{
			rule: &featureproto.Rule{
				Id: rID,
				Strategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{},
				},
			},
			expected: localizedError(statusMissingVariationID, locale.JaJP),
		},
		{
			rule:     expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := validateRule(f.Variations, p.rule)
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
			expected: localizedError(statusMissingRuleID, locale.JaJP),
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
			expected: localizedError(statusIncorrectUUIDFormat, locale.JaJP),
		},
		{
			rule: &featureproto.Rule{
				Id:       rID,
				Strategy: nil,
			},
			expected: localizedError(statusMissingRuleStrategy, locale.JaJP),
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
			expected: localizedError(statusDifferentVariationsSize, locale.JaJP),
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
			expected: localizedError(statusMissingVariationID, locale.JaJP),
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
			expected: localizedError(statusMissingVariationID, locale.JaJP),
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
			expected: localizedError(statusIncorrectVariationWeight, locale.JaJP),
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
			expected: localizedError(statusIncorrectVariationWeight, locale.JaJP),
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
			expected: localizedError(statusExceededMaxVariationWeight, locale.JaJP),
		},
		{
			rule:     expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := validateRule(f.Variations, p.rule)
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
	patterns := []*struct {
		ruleID   string
		strategy *featureproto.Strategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: expected,
			expected: localizedError(statusMissingRuleID, locale.JaJP),
		},
		{
			ruleID:   rID,
			strategy: nil,
			expected: localizedError(statusMissingRuleStrategy, locale.JaJP),
		},
		{
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_FIXED,
			},
			expected: localizedError(statusMissingFixedStrategy, locale.JaJP),
		},
		{
			ruleID: rID,
			strategy: &featureproto.Strategy{
				Type:          featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{},
			},
			expected: localizedError(statusMissingVariationID, locale.JaJP),
		},
		{
			ruleID:   "",
			strategy: nil,
			expected: localizedError(statusMissingRuleID, locale.JaJP),
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
		err := validateChangeRuleStrategy(f.Variations, cmd)
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
			expected: localizedError(statusMissingRuleID, locale.JaJP),
		},
		{
			desc:     "fail: errMissingRuleStrategy",
			ruleID:   rID,
			strategy: nil,
			expected: localizedError(statusMissingRuleStrategy, locale.JaJP),
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
			expected: localizedError(statusDifferentVariationsSize, locale.JaJP),
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
			expected: localizedError(statusMissingVariationID, locale.JaJP),
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
			expected: localizedError(statusMissingVariationID, locale.JaJP),
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
			expected: localizedError(statusIncorrectVariationWeight, locale.JaJP),
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
			expected: localizedError(statusIncorrectVariationWeight, locale.JaJP),
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
			expected: localizedError(statusExceededMaxVariationWeight, locale.JaJP),
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
			expected: localizedError(statusExceededMaxVariationWeight, locale.JaJP),
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
			err := validateChangeRuleStrategy(f.Variations, cmd)
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
	patterns := []*struct {
		ruleID   string
		strategy *featureproto.FixedStrategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: &featureproto.FixedStrategy{Variation: vID},
			expected: localizedError(statusMissingRuleID, locale.JaJP),
		},
		{
			ruleID:   rID,
			strategy: nil,
			expected: localizedError(statusMissingFixedStrategy, locale.JaJP),
		},
		{
			ruleID:   rID,
			strategy: &featureproto.FixedStrategy{},
			expected: localizedError(statusMissingVariationID, locale.JaJP),
		},
		{
			ruleID:   "",
			strategy: nil,
			expected: localizedError(statusMissingRuleID, locale.JaJP),
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
		err := validateChangeFixedStrategy(cmd)
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
	patterns := []*struct {
		ruleID   string
		strategy *featureproto.RolloutStrategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: &featureproto.RolloutStrategy{},
			expected: localizedError(statusMissingRuleID, locale.JaJP),
		},
		{
			ruleID:   rID,
			strategy: nil,
			expected: localizedError(statusMissingRolloutStrategy, locale.JaJP),
		},
		{
			ruleID: rID,
			strategy: &featureproto.RolloutStrategy{Variations: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    70000,
				},
			}},
			expected: localizedError(statusDifferentVariationsSize, locale.JaJP),
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
			expected: localizedError(statusMissingVariationID, locale.JaJP),
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
			expected: localizedError(statusMissingVariationID, locale.JaJP),
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
			expected: localizedError(statusIncorrectVariationWeight, locale.JaJP),
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
			expected: localizedError(statusIncorrectVariationWeight, locale.JaJP),
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
			expected: localizedError(statusExceededMaxVariationWeight, locale.JaJP),
		},
		{
			ruleID:   "",
			strategy: nil,
			expected: localizedError(statusMissingRuleID, locale.JaJP),
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
		err := validateChangeRolloutStrategy(f.Variations, cmd)
		assert.Equal(t, p.expected, err)
	}
}

func TestChangeDefaultStrategy(t *testing.T) {
	t.Parallel()
	f := makeFeature("feature-id")
	patterns := []struct {
		desc        string
		strategy    *featureproto.Strategy
		expectedErr error
	}{
		{
			desc:        "fail: errMissingRuleStrategy",
			strategy:    nil,
			expectedErr: localizedError(statusMissingRuleStrategy, locale.JaJP),
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
			expectedErr: localizedError(statusExceededMaxVariationWeight, locale.JaJP),
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
			expectedErr: localizedError(statusExceededMaxVariationWeight, locale.JaJP),
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
			err := validateChangeDefaultStrategy(f.Variations, cmd)
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
			expectedErr: localizedError(statusInvalidChangingVariation, locale.JaJP),
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
		err := validateFeatureVariationsCommand(p.fs, p.cmd)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestValidateAddPrerequisite(t *testing.T) {
	t.Parallel()
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	fID0 := "fID-0"
	fID1 := "fID-1"
	fID2 := "fID-2"
	fID3 := "fID-3"
	fID4 := "fID-4"
	fID5 := "fID-5"
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
			expectedErr: localizedError(statusCycleExists, locale.JaJP),
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
			expectedErr: localizedError(statusInvalidPrerequisite, locale.JaJP),
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
			expectedErr: localizedError(statusInvalidPrerequisite, locale.JaJP),
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
			expectedErr: localizedError(statusInvalidVariationID, locale.JaJP),
		},
	}
	for _, p := range pattens {
		err := validateAddPrerequisite(p.fs, p.fs[0], p.prerequisite, localizer)
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
			expectedErr: localizedError(statusInvalidVariationID, locale.JaJP),
		},
	}
	for _, p := range pattens {
		err := validateChangePrerequisiteVariation(p.fs, p.prerequisite)
		assert.Equal(t, p.expectedErr, err)
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
