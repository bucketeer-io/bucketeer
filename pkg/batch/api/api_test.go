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
//

package api

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	acclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	aoclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	autoopsdomain "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	cacher "github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/cacher"
	deleter "github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/deleter"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/experiment"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/mau"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/notification"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/opsevent"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/rediscounter"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	ecclient "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/client/mock"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	notificationsender "github.com/bucketeer-io/bucketeer/v2/pkg/notification/sender/mock"
	opsexecutor "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/batch/executor/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	batchproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	ecproto "github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	jpLocation = time.FixedZone("Asia/Tokyo", 9*60*60)
)

type setupMockFunc func(
	accountMockClient *acclient.MockClient,
	environmentMockClient *environmentclient.MockClient,
	autoOpsRulesMockClient *aoclientmock.MockClient,
	experimentMockClient *experimentclient.MockClient,
	featureMockClient *featureclientmock.MockClient,
	eventCounterMockClient *ecclient.MockClient,
	notificationMockSender *notificationsender.MockSender,
	mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
	mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
	// redisCounterDeleterMock *redisCounterDeleterMock,
	mysqlMockClient *mysqlmock.MockClient,
	mysqlMockRows *mysqlmock.MockRows,
	redisMockClient *redismock.MockMultiGetCache,
	mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
)

func TestExperimentStatusUpdater(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		environmentMockClient.EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(
				&environmentproto.ListEnvironmentsV2Response{
					Environments: getEnvironments(t),
				},
				nil,
			)
		experimentMockClient.EXPECT().
			ListExperiments(gomock.Any(), gomock.Any()).
			Times(2).
			Return(
				&experimentproto.ListExperimentsResponse{
					Experiments: getExperiments(t),
				},
				nil,
			)
		experimentMockClient.EXPECT().
			StartExperiment(gomock.Any(), gomock.Any()).
			MinTimes(1).
			Return(
				&experimentproto.StartExperimentResponse{},
				nil,
			)
		experimentMockClient.EXPECT().
			FinishExperiment(gomock.Any(), gomock.Any()).
			MinTimes(1).
			Return(
				&experimentproto.FinishExperimentResponse{},
				nil,
			)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_ExperimentStatusUpdater,
	}, setupMock)
}

func TestExperimentRunningWatcher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		environmentMockClient.EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&environmentproto.ListEnvironmentsV2Response{
					Environments: getEnvironments(t),
				}, nil,
			)
		experimentMockClient.EXPECT().
			ListExperiments(gomock.Any(), gomock.Any()).
			Times(2).
			Return(
				&experimentproto.ListExperimentsResponse{
					Experiments: getExperiments(t),
				},
				nil,
			)
		notificationMockSender.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_ExperimentRunningWatcher,
	}, setupMock)
}

func TestFeatureStaleWatcher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		environmentMockClient.EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&environmentproto.ListEnvironmentsV2Response{
					Environments: getEnvironments(t),
				}, nil,
			)
		featureMockClient.EXPECT().
			ListFeatures(gomock.Any(), gomock.Any()).
			Times(2).
			Return(
				&featureproto.ListFeaturesResponse{
					Features: getFeatures(t),
				},
				nil,
			)
		notificationMockSender.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_FeatureStaleWatcher,
	}, setupMock)
}

func TestMAUCountWatcher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		environmentMockClient.EXPECT().
			ListProjects(gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&environmentproto.ListProjectsResponse{
					Projects: getProjects(t),
				},
				nil,
			)
		environmentMockClient.EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&environmentproto.ListEnvironmentsV2Response{
					Environments: getEnvironments(t),
				},
				nil,
			)
		eventCounterMockClient.EXPECT().
			GetMAUCount(gomock.Any(), gomock.Any()).
			Times(2).
			Return(
				&ecproto.GetMAUCountResponse{
					EventCount: 1,
					UserCount:  1,
				},
				nil,
			)
		notificationMockSender.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_MauCountWatcher,
	}, setupMock)
}
func TestDatetimeWatcher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		environmentMockClient.EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&environmentproto.ListEnvironmentsV2Response{
					Environments: []*environmentproto.EnvironmentV2{
						{Id: "env0", ProjectId: "pj0"},
					},
				},
				nil,
			)
		autoOpsRulesMockClient.EXPECT().ListAutoOpsRules(
			gomock.Any(),
			&autoopsproto.ListAutoOpsRulesRequest{
				PageSize:      0,
				EnvironmentId: "env0",
			},
		).Return(
			&autoopsproto.ListAutoOpsRulesResponse{
				AutoOpsRules: []*autoopsproto.AutoOpsRule{
					newAutoOpsRuleForScheduleType(t),
				},
			},
			nil,
		)
		mockAutoOpsExecutor.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_DatetimeWatcher,
	}, setupMock)
}

func TestEventCountWatcher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		environmentMockClient.EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&environmentproto.ListEnvironmentsV2Response{
					Environments: []*environmentproto.EnvironmentV2{
						{Id: "env0", ProjectId: "pj0"},
						{Id: "env1", ProjectId: "pj1"},
					},
				},
				nil,
			)
		autoOpsRulesMockClient.EXPECT().ListAutoOpsRules(
			gomock.Any(),
			&autoopsproto.ListAutoOpsRulesRequest{
				PageSize:      0,
				EnvironmentId: "env0",
			},
		).Return(
			&autoopsproto.ListAutoOpsRulesResponse{
				AutoOpsRules: []*autoopsproto.AutoOpsRule{
					newAutoOpsRuleForEventRate(t),
				},
			},
			nil,
		)
		autoOpsRulesMockClient.EXPECT().ListAutoOpsRules(
			gomock.Any(),
			&autoopsproto.ListAutoOpsRulesRequest{
				PageSize:      0,
				EnvironmentId: "env1",
			},
		).Return(
			&autoopsproto.ListAutoOpsRulesResponse{
				AutoOpsRules: []*autoopsproto.AutoOpsRule{
					newAutoOpsRuleForEventRate(t),
				},
			},
			nil,
		)
		featureMockClient.EXPECT().
			GetFeature(gomock.Any(), gomock.Any()).
			Times(2).
			Return(
				&featureproto.GetFeatureResponse{
					Feature: &featureproto.Feature{
						Version: 1,
					},
				},
				nil,
			)
		eventCounterMockClient.EXPECT().
			GetOpsGoalUserCount(gomock.Any(), gomock.Any()).
			Times(2).
			Return(
				&ecproto.GetOpsGoalUserCountResponse{
					Count: 12,
				},
				nil,
			)
		eventCounterMockClient.EXPECT().
			GetOpsEvaluationUserCount(gomock.Any(), gomock.Any()).
			Times(2).
			Return(
				&ecproto.GetOpsEvaluationUserCountResponse{
					Count: 10,
				},
				nil,
			)
		mysqlMockClient.EXPECT().
			ExecContext(
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil, nil)
		mysqlMockQueryExecer.EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).AnyTimes().Return(nil, nil)
		mockAutoOpsExecutor.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)

	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_EventCountWatcher,
	}, setupMock)
}

func TestProgressiveRolloutWatcher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		environmentMockClient.EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&environmentproto.ListEnvironmentsV2Response{
					Environments: []*environmentproto.EnvironmentV2{
						{Id: "env0", ProjectId: "pj0"},
						{Id: "env1", ProjectId: "pj1"},
					},
				},
				nil,
			)
		autoOpsRulesMockClient.EXPECT().ListProgressiveRollouts(
			gomock.Any(),
			&autoopsproto.ListProgressiveRolloutsRequest{
				PageSize:      0,
				EnvironmentId: "env0",
			},
		).Return(
			&autoopsproto.ListProgressiveRolloutsResponse{
				ProgressiveRollouts: []*autoopsproto.ProgressiveRollout{
					newProgressiveRollout(t),
				},
			},
			nil,
		)
		autoOpsRulesMockClient.EXPECT().ListProgressiveRollouts(
			gomock.Any(),
			&autoopsproto.ListProgressiveRolloutsRequest{
				PageSize:      0,
				EnvironmentId: "env1",
			},
		).Return(
			&autoopsproto.ListProgressiveRolloutsResponse{
				ProgressiveRollouts: []*autoopsproto.ProgressiveRollout{
					newProgressiveRollout(t),
				},
			},
			nil,
		)
		mockProgressiveRolloutExecutor.EXPECT().
			ExecuteProgressiveRollout(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_ProgressiveRolloutWatcher,
	}, setupMock)
}

func TestFeatureFlagCacher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		mysqlMockRows.EXPECT().Close().Return(nil)
		mysqlMockRows.EXPECT().Next().Return(false)
		mysqlMockRows.EXPECT().Err().Return(nil)

		mysqlMockClient.EXPECT().QueryContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(mysqlMockRows, nil)
		redisMockClient.EXPECT().
			Put(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes().
			Return(nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_FeatureFlagCacher,
	}, setupMock)
}

func TestSegmentUserCacher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		environmentMockClient.EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(
				&environmentproto.ListEnvironmentsV2Response{
					Environments: []*environmentproto.EnvironmentV2{
						{Id: "env0", ProjectId: "pj0"},
						{Id: "env1", ProjectId: "pj1"},
					},
				},
				nil,
			)
		featureMockClient.EXPECT().
			ListSegments(gomock.Any(), gomock.Any()).
			Times(2).
			Return(
				&featureproto.ListSegmentsResponse{
					Segments: []*featureproto.Segment{
						{
							Id: "segment-id",
						},
					},
				},
				nil,
			)
		featureMockClient.EXPECT().
			ListSegmentUsers(gomock.Any(), gomock.Any()).
			Times(2).
			Return(
				&featureproto.ListSegmentUsersResponse{
					Users: []*featureproto.SegmentUser{},
				},
				nil,
			)
		redisMockClient.EXPECT().
			Put(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_SegmentUserCacher,
	}, setupMock)
}

func TestAPIKeyCacher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		mysqlMockRows.EXPECT().Close().Return(nil)
		mysqlMockRows.EXPECT().Next().Return(false)
		mysqlMockRows.EXPECT().Err().Return(nil)

		mysqlMockClient.EXPECT().QueryContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(mysqlMockRows, nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_ApiKeyCacher,
	}, setupMock)
}

func TestExperimentCacher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		environmentMockClient.EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(
				&environmentproto.ListEnvironmentsV2Response{
					Environments: []*environmentproto.EnvironmentV2{
						{Id: "env0", ProjectId: "pj0"},
						{Id: "env1", ProjectId: "pj1"},
					},
				},
				nil,
			)
		experimentMockClient.EXPECT().
			ListExperiments(gomock.Any(), gomock.Any()).
			Times(2).
			Return(
				&experimentproto.ListExperimentsResponse{
					Experiments: getExperiments(t),
				},
				nil,
			)
		redisMockClient.EXPECT().
			Put(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_ExperimentCacher,
	}, setupMock)
}

func TestAutoOpsRulesCacher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	) {
		environmentMockClient.EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(
				&environmentproto.ListEnvironmentsV2Response{
					Environments: []*environmentproto.EnvironmentV2{
						{Id: "env0", ProjectId: "pj0"},
						{Id: "env1", ProjectId: "pj1"},
					},
				},
				nil,
			)
		autoOpsRulesMockClient.EXPECT().ListAutoOpsRules(
			gomock.Any(),
			gomock.Any(),
		).
			Times(2).
			Return(
				&autoopsproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*autoopsproto.AutoOpsRule{
						newAutoOpsRuleForScheduleType(t),
					},
				},
				nil,
			)
		redisMockClient.EXPECT().
			Put(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_AutoOpsRulesCacher,
	}, setupMock)
}

func executeMockBatchJob(t *testing.T,
	request *batchproto.BatchJobRequest, setupMock setupMockFunc) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := newBatchService(t, controller, setupMock)
	_, err := service.ExecuteBatchJob(context.Background(), request)
	assert.NoError(t, err)
}

func newBatchService(t *testing.T,
	mockController *gomock.Controller, setupMock setupMockFunc) *batchService {
	logger, err := log.NewLogger()
	require.NoError(t, err)

	accountMockClient := acclient.NewMockClient(mockController)
	environmentMockClient := environmentclient.NewMockClient(mockController)
	autoOpsRulesMockClient := aoclientmock.NewMockClient(mockController)
	experimentMockClient := experimentclient.NewMockClient(mockController)
	featureMockClient := featureclientmock.NewMockClient(mockController)
	eventCounterMockClient := ecclient.NewMockClient(mockController)
	notificationMockSender := notificationsender.NewMockSender(mockController)
	mockAutoOpsExecutor := opsexecutor.NewMockAutoOpsExecutor(mockController)
	mockProgressiveRolloutExecutor := opsexecutor.NewMockProgressiveRolloutExecutor(mockController)
	cacheMock := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	mysqlMockClient := mysqlmock.NewMockClient(mockController)
	mysqlMockRows := mysqlmock.NewMockRows(mockController)
	redisMockClient := redismock.NewMockMultiGetCache(mockController)
	mysqlMockQueryExecer := mysqlmock.NewMockQueryExecer(mockController)

	setupMock(
		accountMockClient,
		environmentMockClient,
		autoOpsRulesMockClient,
		experimentMockClient,
		featureMockClient,
		eventCounterMockClient,
		notificationMockSender,
		mockAutoOpsExecutor,
		mockProgressiveRolloutExecutor,
		mysqlMockClient,
		mysqlMockRows,
		redisMockClient,
		mysqlMockQueryExecer,
	)

	service := NewBatchService(
		experiment.NewExperimentStatusUpdater(
			environmentMockClient,
			experimentMockClient,
			jobs.WithLogger(logger),
		),
		notification.NewExperimentRunningWatcher(
			environmentMockClient,
			experimentMockClient,
			notificationMockSender,
			jobs.WithTimeout(1*time.Minute),
			jobs.WithLogger(logger),
		),
		notification.NewFeatureStaleWatcher(
			environmentMockClient,
			featureMockClient,
			notificationMockSender,
			jobs.WithTimeout(1*time.Minute),
			jobs.WithLogger(logger),
		),
		notification.NewMAUCountWatcher(
			environmentMockClient,
			eventCounterMockClient,
			notificationMockSender,
			jpLocation,
			jobs.WithTimeout(60*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewDatetimeWatcher(
			environmentMockClient,
			autoOpsRulesMockClient,
			mockAutoOpsExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewEventCountWatcher(
			mysqlMockClient,
			environmentMockClient,
			autoOpsRulesMockClient,
			eventCounterMockClient,
			featureMockClient,
			mockAutoOpsExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewProgressiveRolloutWacher(
			environmentMockClient,
			autoOpsRulesMockClient,
			mockProgressiveRolloutExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		rediscounter.NewRedisCounterDeleter(
			cacheMock,
			environmentMockClient,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		nil,
		mau.NewMAUSummarizer(
			mysqlMockClient,
			eventCounterMockClient,
			jpLocation,
			jobs.WithTimeout(30*time.Minute),
			jobs.WithLogger(logger),
		),
		mau.NewMAUPartitionDeleter(
			mysqlMockClient,
			jpLocation,
			jobs.WithTimeout(60*time.Minute),
			jobs.WithLogger(logger),
		),
		mau.NewMAUPartitionCreator(
			mysqlMockClient,
			jpLocation,
			jobs.WithTimeout(60*time.Minute),
			jobs.WithLogger(logger),
		),
		cacher.NewFeatureFlagCacher(
			mysqlMockClient,
			[]cache.MultiGetCache{redisMockClient},
		),
		cacher.NewSegmentUserCacher(
			environmentMockClient,
			featureMockClient,
			[]cache.MultiGetCache{redisMockClient},
		),
		cacher.NewAPIKeyCacher(
			mysqlMockClient,
			[]cache.MultiGetCache{redisMockClient},
		),
		cacher.NewExperimentCacher(
			environmentMockClient,
			experimentMockClient,
			[]cache.MultiGetCache{redisMockClient},
		),
		cacher.NewAutoOpsRulesCacher(
			environmentMockClient,
			autoOpsRulesMockClient,
			redisMockClient,
		),
		deleter.NewTagDeleter(mysqlMockClient),
		logger,
	)
	return service
}

func getEnvironments(t *testing.T) []*environmentproto.EnvironmentV2 {
	t.Helper()
	return []*environmentproto.EnvironmentV2{
		{Id: "ns0", Name: "ns0", ProjectId: "pj0"},
		{Id: "ns1", Name: "ns1", ProjectId: "pj0"},
	}
}

func getExperiments(t *testing.T) []*experimentproto.Experiment {
	t.Helper()
	return []*experimentproto.Experiment{
		{
			Id:     "eid0",
			Status: experimentproto.Experiment_WAITING,
		},
		{
			Id:     "eid1",
			Status: experimentproto.Experiment_RUNNING,
		},
	}
}

func getFeatures(t *testing.T) []*featureproto.Feature {
	t.Helper()
	return []*featureproto.Feature{
		{
			Id:           "fid0",
			Enabled:      true,
			OffVariation: "variation",
			LastUsedInfo: &featureproto.FeatureLastUsedInfo{
				// Stale feature
				LastUsedAt: time.Now().Unix() - featuredomain.SecondsToStale - 20,
			},
		},
	}
}

func getProjects(t *testing.T) []*environmentproto.Project {
	t.Helper()
	return []*environmentproto.Project{
		{
			Id:          "pj0",
			Description: "pj0",
		},
	}
}

func newAutoOpsRuleForScheduleType(t *testing.T) *autoopsproto.AutoOpsRule {
	dc1 := &autoopsproto.DatetimeClause{
		Time: 1000000001,
	}
	dc2 := &autoopsproto.DatetimeClause{
		Time: 1000000002,
	}
	aor, err := autoopsdomain.NewAutoOpsRule(
		"fid",
		autoopsproto.OpsType_SCHEDULE,
		[]*autoopsproto.OpsEventRateClause{},
		[]*autoopsproto.DatetimeClause{dc1, dc2},
	)
	require.NoError(t, err)
	return aor.AutoOpsRule
}

func newAutoOpsRuleForEventRate(t *testing.T) *autoopsproto.AutoOpsRule {
	oerc := &autoopsproto.OpsEventRateClause{
		GoalId:          "gid",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	aor, err := autoopsdomain.NewAutoOpsRule(
		"fid",
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{oerc},
		[]*autoopsproto.DatetimeClause{},
	)
	require.NoError(t, err)
	return aor.AutoOpsRule
}

func newProgressiveRollout(t *testing.T) *autoopsproto.ProgressiveRollout {
	dc := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
		Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
			{
				ScheduleId: "sID",
				ExecuteAt:  time.Now().Unix(),
			},
		},
	}
	c, err := ptypes.MarshalAny(dc)
	require.NoError(t, err)
	return &autoopsproto.ProgressiveRollout{
		Id:        "prID",
		FeatureId: "fID",
		Clause:    c,
		Type:      autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
	}
}
