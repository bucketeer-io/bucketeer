package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	eventcounterclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/client/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventcounterproto "github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

func TestGrpcGatewayService_GetExperimentEvaluationCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetExperimentEvaluationCountResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: get auto ops error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetExperimentEvaluationCount(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetExperimentEvaluationCount(
					gomock.Any(), gomock.Any(),
				).Return(&eventcounterproto.GetExperimentEvaluationCountResponse{
					FeatureId: "feature-1",
					VariationCounts: []*eventcounterproto.VariationCount{
						{
							VariationId: "variation-1",
							UserCount:   100,
							EventCount:  100,
						},
						{
							VariationId: "variation-2",
							UserCount:   200,
							EventCount:  200,
						},
					},
				}, nil)
			},
			expected: &gwproto.GetExperimentEvaluationCountResponse{
				FeatureId: "feature-1",
				VariationCounts: []*eventcounterproto.VariationCount{
					{
						VariationId: "variation-1",
						UserCount:   100,
						EventCount:  100,
					},
					{
						VariationId: "variation-2",
						UserCount:   200,
						EventCount:  200,
					},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetExperimentEvaluationCount(ctx, &gwproto.GetExperimentEvaluationCountRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_GetEvaluationTimeseriesCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetEvaluationTimeseriesCountResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: get timeseries error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetEvaluationTimeseriesCount(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetEvaluationTimeseriesCount(
					gomock.Any(), gomock.Any(),
				).Return(&eventcounterproto.GetEvaluationTimeseriesCountResponse{
					UserCounts: []*eventcounterproto.VariationTimeseries{
						{
							VariationId: "variation-1",
							Timeseries: &eventcounterproto.Timeseries{
								Timestamps:  []int64{1622548800, 1622635200},
								Values:      []float64{100, 150},
								TotalCounts: 100,
							},
						},
					},
					EventCounts: []*eventcounterproto.VariationTimeseries{
						{
							VariationId: "variation-2",
							Timeseries: &eventcounterproto.Timeseries{
								Timestamps:  []int64{1622548800, 1622635200},
								Values:      []float64{100, 150},
								TotalCounts: 100,
							},
						},
					},
				}, nil)
			},
			expected: &gwproto.GetEvaluationTimeseriesCountResponse{
				UserCounts: []*eventcounterproto.VariationTimeseries{
					{
						VariationId: "variation-1",
						Timeseries: &eventcounterproto.Timeseries{
							Timestamps:  []int64{1622548800, 1622635200},
							Values:      []float64{100, 150},
							TotalCounts: 100,
						},
					},
				},
				EventCounts: []*eventcounterproto.VariationTimeseries{
					{
						VariationId: "variation-2",
						Timeseries: &eventcounterproto.Timeseries{
							Timestamps:  []int64{1622548800, 1622635200},
							Values:      []float64{100, 150},
							TotalCounts: 100,
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetEvaluationTimeseriesCount(ctx, &gwproto.GetEvaluationTimeseriesCountRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_GetExperimentResult(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetExperimentResultResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: get timeseries error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetExperimentResult(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetExperimentResult(
					gomock.Any(), gomock.Any(),
				).Return(&eventcounterproto.GetExperimentResultResponse{
					ExperimentResult: &eventcounterproto.ExperimentResult{
						Id:           "result-1",
						ExperimentId: "exp-1",
						UpdatedAt:    1622548800,
						GoalResults: []*eventcounterproto.GoalResult{
							{
								GoalId: "goal-1",
								VariationResults: []*eventcounterproto.VariationResult{
									{
										VariationId:  "variation-1",
										ExpectedLoss: 0.1,
									},
								},
							},
						},
					},
				}, nil)
			},
			expected: &gwproto.GetExperimentResultResponse{
				ExperimentResult: &eventcounterproto.ExperimentResult{
					Id:           "result-1",
					ExperimentId: "exp-1",
					UpdatedAt:    1622548800,
					GoalResults: []*eventcounterproto.GoalResult{
						{
							GoalId: "goal-1",
							VariationResults: []*eventcounterproto.VariationResult{
								{
									VariationId:  "variation-1",
									ExpectedLoss: 0.1,
								},
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetExperimentResult(ctx, &gwproto.GetExperimentResultRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_ListExperimentResults(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		expected    *gwproto.ListExperimentResultsResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: get timeseries error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().ListExperimentResults(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().ListExperimentResults(
					gomock.Any(), gomock.Any(),
				).Return(&eventcounterproto.ListExperimentResultsResponse{
					Results: map[string]*eventcounterproto.ExperimentResult{
						"result-1": {
							Id:           "result-1",
							ExperimentId: "exp-1",
							UpdatedAt:    1622548800,
							GoalResults: []*eventcounterproto.GoalResult{
								{
									GoalId: "goal-1",
									VariationResults: []*eventcounterproto.VariationResult{
										{
											VariationId:  "variation-1",
											ExpectedLoss: 0.1,
										},
									},
								},
							},
						},
					},
				}, nil)
			},
			expected: &gwproto.ListExperimentResultsResponse{
				Results: map[string]*eventcounterproto.ExperimentResult{
					"result-1": {
						Id:           "result-1",
						ExperimentId: "exp-1",
						UpdatedAt:    1622548800,
						GoalResults: []*eventcounterproto.GoalResult{
							{
								GoalId: "goal-1",
								VariationResults: []*eventcounterproto.VariationResult{
									{
										VariationId:  "variation-1",
										ExpectedLoss: 0.1,
									},
								},
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.ListExperimentResults(ctx, &gwproto.ListExperimentResultsRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_GetExperimentGoalCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetExperimentGoalCountResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: get timeseries error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetExperimentGoalCount(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetExperimentGoalCount(
					gomock.Any(), gomock.Any(),
				).Return(&eventcounterproto.GetExperimentGoalCountResponse{
					GoalId: "goal-1",
					VariationCounts: []*eventcounterproto.VariationCount{
						{
							VariationId: "variation-1",
							UserCount:   100,
							EventCount:  100,
						},
						{
							VariationId: "variation-2",
							UserCount:   200,
							EventCount:  200,
						},
					},
				}, nil)
			},
			expected: &gwproto.GetExperimentGoalCountResponse{
				GoalId: "goal-1",
				VariationCounts: []*eventcounterproto.VariationCount{
					{
						VariationId: "variation-1",
						UserCount:   100,
						EventCount:  100,
					},
					{
						VariationId: "variation-2",
						UserCount:   200,
						EventCount:  200,
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetExperimentGoalCount(ctx, &gwproto.GetExperimentGoalCountRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_GetOpsEvaluationUserCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetOpsEvaluationUserCountResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: get timeseries error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetOpsEvaluationUserCount(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetOpsEvaluationUserCount(
					gomock.Any(), gomock.Any(),
				).Return(&eventcounterproto.GetOpsEvaluationUserCountResponse{
					OpsRuleId: "ops-rule-1",
					ClauseId:  "clause-1",
					Count:     100,
				}, nil)
			},
			expected: &gwproto.GetOpsEvaluationUserCountResponse{
				OpsRuleId: "ops-rule-1",
				ClauseId:  "clause-1",
				Count:     100,
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetOpsEvaluationUserCount(ctx, &gwproto.GetOpsEvaluationUserCountRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_GetOpsGoalUserCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetOpsGoalUserCountResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: get timeseries error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetOpsGoalUserCount(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.eventCounterClient.(*eventcounterclientmock.MockClient).EXPECT().GetOpsGoalUserCount(
					gomock.Any(), gomock.Any(),
				).Return(&eventcounterproto.GetOpsGoalUserCountResponse{
					OpsRuleId: "ops-rule-1",
					ClauseId:  "clause-1",
					Count:     100,
				}, nil)
			},
			expected: &gwproto.GetOpsGoalUserCountResponse{
				OpsRuleId: "ops-rule-1",
				ClauseId:  "clause-1",
				Count:     100,
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetOpsGoalUserCount(ctx, &gwproto.GetOpsGoalUserCountRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}
