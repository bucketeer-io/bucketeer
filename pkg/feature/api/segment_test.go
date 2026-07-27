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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"

	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestCreateSegmentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	testcases := []struct {
		desc     string
		setup    func(*FeatureService)
		req      *featureproto.CreateSegmentRequest
		expected error
	}{
		{
			desc:  "error: missing name",
			setup: nil,
			req: &featureproto.CreateSegmentRequest{
				Name:          "",
				Description:   "description",
				EnvironmentId: "ns0",
			},
			expected: statusMissingName.Err(),
		},
		{
			desc:  "error: rule with SEGMENT operator",
			setup: nil,
			req: &featureproto.CreateSegmentRequest{
				Name:          "name",
				EnvironmentId: "ns0",
				Rules: []*featureproto.Rule{
					{
						Clauses: []*featureproto.Clause{
							{
								Attribute: "",
								Operator:  featureproto.Clause_SEGMENT,
								Values:    []string{"segment-id"},
							},
						},
					},
				},
			},
			expected: statusSegmentRuleOperatorNotAllowed.Err(),
		},
		{
			desc:  "error: rule with strategy",
			setup: nil,
			req: &featureproto.CreateSegmentRequest{
				Name:          "name",
				EnvironmentId: "ns0",
				Rules: []*featureproto.Rule{
					{
						Strategy: &featureproto.Strategy{
							Type: featureproto.Strategy_FIXED,
							FixedStrategy: &featureproto.FixedStrategy{
								Variation: "variation-a",
							},
						},
						Clauses: []*featureproto.Clause{
							{
								Attribute: "plan",
								Operator:  featureproto.Clause_EQUALS,
								Values:    []string{"premium"},
							},
						},
					},
				},
			},
			expected: statusSegmentRuleStrategyNotAllowed.Err(),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					err := fn(ctx)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().CreateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.CreateSegmentRequest{
				Name:          "name",
				Description:   "description",
				EnvironmentId: "ns0",
			},
			expected: nil,
		},
		{
			desc: "success with rules: ids are generated server-side",
			setup: func(s *FeatureService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					err := fn(ctx)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().CreateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.CreateSegmentRequest{
				Name:          "name",
				Description:   "description",
				EnvironmentId: "ns0",
				Rules: []*featureproto.Rule{
					{
						// No ids: they must be generated server-side.
						Clauses: []*featureproto.Clause{
							{
								Attribute: "plan",
								Operator:  featureproto.Clause_EQUALS,
								Values:    []string{"premium"},
							},
						},
					},
				},
			},
			expected: nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		resp, err := service.CreateSegment(ctx, tc.req)
		assert.Equal(t, tc.expected, err)
		if err == nil && len(tc.req.Rules) > 0 {
			require.Equal(t, len(tc.req.Rules), len(resp.Segment.Rules))
			for _, rule := range resp.Segment.Rules {
				assert.NoError(t, uuid.ValidateUUID(rule.Id))
				for _, clause := range rule.Clauses {
					assert.NoError(t, uuid.ValidateUUID(clause.Id))
				}
			}
		}
	}
}

func TestDeleteSegmentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	testcases := []struct {
		desc     string
		setup    func(*FeatureService)
		req      *featureproto.DeleteSegmentRequest
		expected error
	}{
		{
			desc:  "error: missing id",
			setup: nil,
			req: &featureproto.DeleteSegmentRequest{
				Id:            "",
				EnvironmentId: "ns0",
			},
			expected: statusMissingID.Err(),
		},
		{
			desc: "error: segment not found",
			setup: func(s *FeatureService) {
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{}, 0, int64(0), nil)
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrSegmentNotFound)
			},
			req: &featureproto.DeleteSegmentRequest{
				Id:            "id",
				EnvironmentId: "ns0",
			},
			expected: statusSegmentNotFound.Err(),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{}, 0, int64(0), nil)
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					err := fn(ctx)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{
						Id: "id",
					},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().DeleteSegment(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.DeleteSegmentRequest{
				Id:            "id",
				EnvironmentId: "ns0",
			},
			expected: nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		_, err := service.DeleteSegment(ctx, tc.req)
		assert.Equal(t, tc.expected, err)
	}
}

func TestUpdateSegmentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	testcases := []struct {
		desc     string
		setup    func(*FeatureService)
		req      *featureproto.UpdateSegmentRequest
		expected error
	}{
		{
			desc:  "error: missing id",
			setup: nil,
			req: &featureproto.UpdateSegmentRequest{
				EnvironmentId: "ns0",
			},
			expected: statusMissingID.Err(),
		},
		{
			desc:  "error: update empty name",
			setup: nil,
			req: &featureproto.UpdateSegmentRequest{
				Id:            "id0",
				EnvironmentId: "ns0",
				Name:          wrapperspb.String(""),
			},
			expected: statusMissingName.Err(),
		},
		{
			desc: "success update name",
			setup: func(s *FeatureService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					err := fn(ctx)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.All(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{
						Id: "id0",
					},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.UpdateSegmentRequest{
				Id:            "id0",
				EnvironmentId: "ns0",
				Name:          wrapperspb.String("new-name"),
			},
			expected: nil,
		},
		{
			desc: "success update description",
			setup: func(s *FeatureService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					err := fn(ctx)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.All(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{
						Id: "id0",
					},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.UpdateSegmentRequest{
				Id:            "id0",
				EnvironmentId: "ns0",
				Description:   wrapperspb.String("new-description"),
			},
		},
		{
			desc:  "error: invalid rules",
			setup: nil,
			req: &featureproto.UpdateSegmentRequest{
				Id:            "id0",
				EnvironmentId: "ns0",
				Rules: &featureproto.RuleListValue{
					Values: []*featureproto.Rule{
						{
							Clauses: []*featureproto.Clause{
								{
									Attribute: "plan",
									Operator:  featureproto.Clause_EQUALS,
									// Missing values.
								},
							},
						},
					},
				},
			},
			expected: statusSegmentRuleClauseValuesRequired.Err(),
		},
		{
			desc: "success replace rules: cache is refreshed reusing the cached user list",
			setup: func(s *FeatureService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					err := fn(ctx)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.All(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{
						Id: "id0",
					},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				// Cache hit: the user list is reused, storage is not queried.
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(
					"id0", "ns0",
				).Return(&featureproto.SegmentUsers{
					SegmentId: "id0",
					Users:     []*featureproto.SegmentUser{{UserId: "user-0"}},
				}, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(nil)
			},
			req: &featureproto.UpdateSegmentRequest{
				Id:            "id0",
				EnvironmentId: "ns0",
				Rules: &featureproto.RuleListValue{
					Values: []*featureproto.Rule{
						{
							Clauses: []*featureproto.Clause{
								{
									Attribute: "plan",
									Operator:  featureproto.Clause_EQUALS,
									Values:    []string{"premium"},
								},
							},
						},
					},
				},
			},
			expected: nil,
		},
		{
			desc: "success clear rules: present with no values replaces with empty list",
			setup: func(s *FeatureService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					err := fn(ctx)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.All(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{
						Id: "id0",
						Rules: []*featureproto.Rule{
							{
								Id: "b52d3181-e6f0-4d4c-b40f-9891d56a708e",
								Clauses: []*featureproto.Clause{
									{
										Id:        "3ecb45b5-90e4-4d0c-9d4c-468ba9ee2b0c",
										Attribute: "plan",
										Operator:  featureproto.Clause_EQUALS,
										Values:    []string{"premium"},
									},
								},
							},
						},
					},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				// Cache miss: the user list is loaded from storage.
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(
					"id0", "ns0",
				).Return(nil, errors.New("cache miss"))
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().ListSegmentUsersBySegment(
					gomock.Any(), "id0", "ns0",
				).Return([]*featureproto.SegmentUser{}, nil)
				s.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Put(
					gomock.Any(), "ns0",
				).Return(nil)
			},
			req: &featureproto.UpdateSegmentRequest{
				Id:            "id0",
				EnvironmentId: "ns0",
				Rules:         &featureproto.RuleListValue{},
			},
			expected: nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		resp, err := service.UpdateSegment(ctx, tc.req)
		assert.Equal(t, tc.expected, err)
		if err == nil && tc.req.Rules != nil {
			// Present = full replacement, including replacement with an empty list.
			assert.Equal(t, len(tc.req.Rules.Values), len(resp.Segment.Rules))
		}
	}
}

func TestGetSegmentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testcases := []struct {
		desc           string
		setup          func(*FeatureService)
		service        *FeatureService
		context        context.Context
		id             string
		environmentId  string
		getExpectedErr func() error
	}{
		{
			desc:    "error: missing id",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:         nil,
			id:            "",
			environmentId: "ns0",
			getExpectedErr: func() error {
				return statusMissingID.Err()
			},
		},
		{
			desc:    "error: segment not found",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil, v2fs.ErrSegmentNotFound)
			},
			id:            "id",
			environmentId: "ns0",
			getExpectedErr: func() error {
				return statusSegmentNotFound.Err()
			},
		},
		{
			desc:    "success",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{}}, nil, nil)
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					{
						Id: "id",
					},
				}, 0, int64(0), nil)
			},
			id:            "id",
			environmentId: "ns0",
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success with Viewer account",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{}}, nil, nil)
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					{
						Id: "id",
					},
				}, 0, int64(0), nil)
			},
			id:            "id",
			environmentId: "ns0",
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:         func(s *FeatureService) {},
			id:            "id",
			environmentId: "ns0",
			getExpectedErr: func() error {
				return statusPermissionDenied.Err()
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			service := tc.service
			if tc.setup != nil {
				tc.setup(service)
			}
			ctx := tc.context

			req := &featureproto.GetSegmentRequest{Id: tc.id, EnvironmentId: tc.environmentId}
			_, err := service.GetSegment(ctx, req)
			assert.Equal(t, tc.getExpectedErr(), err)
		})
	}
}

func TestListSegmentsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testcases := []struct {
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(*FeatureService)
		pageSize       int64
		environmentId  string
		getExpectedErr func() error
	}{
		{
			desc:    "error: exceeded max page size per request",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:         nil,
			pageSize:      int64(maxPageSizePerRequest + 1),
			environmentId: "ns0",
			getExpectedErr: func() error {
				return statusExceededMaxPageSizePerRequest.Err()
			},
		},
		{
			desc:    "success",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().ListSegments(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Segment{
					{
						Id: "id",
					},
				}, 0, int64(0), map[string][]string{
					"id": {"id"},
				}, nil)
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					{
						Id: "id",
					},
				}, 0, int64(0), nil)
			},
			pageSize:      int64(maxPageSizePerRequest),
			environmentId: "ns0",
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "success with Viewer account",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().ListSegments(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Segment{
					{
						Id: "id",
					},
				}, 0, int64(0), map[string][]string{
					"id": {"id"},
				}, nil)
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					{
						Id: "id",
					},
				}, 0, int64(0), nil)
			},
			pageSize:      int64(maxPageSizePerRequest),
			environmentId: "ns0",
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:         func(s *FeatureService) {},
			pageSize:      int64(maxPageSizePerRequest),
			environmentId: "ns0",
			getExpectedErr: func() error {
				return statusPermissionDenied.Err()
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			service := tc.service
			if tc.setup != nil {
				tc.setup(service)
			}
			ctx := tc.context

			req := &featureproto.ListSegmentsRequest{PageSize: tc.pageSize, EnvironmentId: tc.environmentId}
			_, err := service.ListSegments(ctx, req)
			assert.Equal(t, tc.getExpectedErr(), err)
		})
	}
}

func setToken(ctx context.Context) context.Context {
	t := &token.AccessToken{
		Issuer:   "issuer",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
	return context.WithValue(ctx, rpc.AccessTokenKey, t)
}
