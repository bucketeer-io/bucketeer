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
//

package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/anypb"

	acclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	accstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	aoclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	autoopsdomain "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/autoarchive"
	cacher "github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/cacher"
	deleter "github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/deleter"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/experiment"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/monthlysummary"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/notification"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/opsevent"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/rediscounter"
	scheduledflagchange "github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs/scheduledflagchange"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	maucachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	coderefstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/coderef/storage/mock"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	envstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2/mock"
	ecclient "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/client/mock"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	featurestoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	insightsstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/insights/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	opsexecutor "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/batch/executor/mock"
	opseventstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/storage/v2/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	subscriptionsender "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/sender/mock"
	tagstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	batchproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	ecproto "github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

type setupMockFunc func(
	accountMockClient *acclient.MockClient,
	environmentMockClient *environmentclient.MockClient,
	autoOpsRulesMockClient *aoclientmock.MockClient,
	experimentMockClient *experimentclient.MockClient,
	featureMockClient *featureclientmock.MockClient,
	eventCounterMockClient *ecclient.MockClient,
	notificationMockSender *subscriptionsender.MockSender,
	mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
	mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
	mysqlMockClient *mysqlmock.MockClient,
	mysqlMockRows *mysqlmock.MockRows,
	redisMockClient *redismock.MockMultiGetCache,
	mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
	featureStorageMock *featurestoragemock.MockFeatureStorage,
	segmentStorageMock *featurestoragemock.MockSegmentStorage,
	tagStorageMock *tagstoragemock.MockTagStorage,
	accountStorageMock *accstoragemock.MockAccountStorage,
	opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
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
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
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
			UpdateExperiment(gomock.Any(), gomock.Any()).
			MinTimes(2).
			Return(
				&experimentproto.UpdateExperimentResponse{},
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
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
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
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
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

func TestDatetimeWatcher(t *testing.T) {
	t.Parallel()
	setupMock := func(
		accountMockClient *acclient.MockClient,
		environmentMockClient *environmentclient.MockClient,
		autoOpsRulesMockClient *aoclientmock.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
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
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
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
		opsCountStorageMock.EXPECT().
			UpsertOpsCount(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)
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
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
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
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
	) {
		featureStorageMock.EXPECT().ListAllEnvironmentFeatures(
			gomock.Any(),
		).Return([]*featureproto.EnvironmentFeature{}, nil)
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
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
	) {
		segmentStorageMock.EXPECT().ListAllInUseSegments(
			gomock.Any(),
		).Return(nil, nil)
		redisMockClient.EXPECT().
			Put(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes().
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
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
	) {
		accountStorageMock.EXPECT().
			ListAllEnvironmentAPIKeys(gomock.Any()).
			Return(nil, nil)
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
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
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
		notificationMockSender *subscriptionsender.MockSender,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		mockProgressiveRolloutExecutor *opsexecutor.MockProgressiveRolloutExecutor,
		mysqlMockClient *mysqlmock.MockClient,
		mysqlMockRows *mysqlmock.MockRows,
		redisMockClient *redismock.MockMultiGetCache,
		mysqlMockQueryExecer *mysqlmock.MockQueryExecer,
		featureStorageMock *featurestoragemock.MockFeatureStorage,
		segmentStorageMock *featurestoragemock.MockSegmentStorage,
		tagStorageMock *tagstoragemock.MockTagStorage,
		accountStorageMock *accstoragemock.MockAccountStorage,
		opsCountStorageMock *opseventstoragemock.MockOpsCountStorage,
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
	notificationMockSender := subscriptionsender.NewMockSender(mockController)
	mockAutoOpsExecutor := opsexecutor.NewMockAutoOpsExecutor(mockController)
	mockProgressiveRolloutExecutor := opsexecutor.NewMockProgressiveRolloutExecutor(mockController)
	cacheMock := redismock.NewMockMultiGetDeleteCountCache(mockController)
	envStorageMock := envstoragemock.NewMockEnvironmentStorage(mockController)
	codeRefStorageMock := coderefstoragemock.NewMockCodeReferenceStorage(mockController)
	opsCountStorageMock := opseventstoragemock.NewMockOpsCountStorage(mockController)
	mysqlMockClient := mysqlmock.NewMockClient(mockController)
	mysqlMockRows := mysqlmock.NewMockRows(mockController)
	redisMockClient := redismock.NewMockMultiGetCache(mockController)
	mysqlMockQueryExecer := mysqlmock.NewMockQueryExecer(mockController)
	mauCacheMock := maucachemock.NewMockMAUCache(mockController)
	monthlySummaryStorageMock := insightsstoragemock.NewMockMonthlySummaryStorage(mockController)
	accountStorageMock := accstoragemock.NewMockAccountStorage(mockController)
	featureStorageMock := featurestoragemock.NewMockFeatureStorage(mockController)
	segmentStorageMock := featurestoragemock.NewMockSegmentStorage(mockController)
	tagStorageMock := tagstoragemock.NewMockTagStorage(mockController)
	scheduledFlagChangeStorageMock := featurestoragemock.NewMockScheduledFlagChangeStorage(mockController)

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
		featureStorageMock,
		segmentStorageMock,
		tagStorageMock,
		accountStorageMock,
		opsCountStorageMock,
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
		opsevent.NewDatetimeWatcher(
			environmentMockClient,
			autoOpsRulesMockClient,
			mockAutoOpsExecutor,
			nil, // ftCacher - not needed for this test
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewEventCountWatcher(
			opsCountStorageMock,
			environmentMockClient,
			autoOpsRulesMockClient,
			eventCounterMockClient,
			featureMockClient,
			mockAutoOpsExecutor,
			nil, // ftCacher - not needed for this test
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewProgressiveRolloutWatcher(
			environmentMockClient,
			autoOpsRulesMockClient,
			mockProgressiveRolloutExecutor,
			nil, // ftCacher - not needed for this test
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
		cacher.NewFeatureFlagCacher(
			featureStorageMock,
			[]cache.MultiGetCache{redisMockClient},
		),
		cacher.NewSegmentUserCacher(
			segmentStorageMock,
			[]cache.MultiGetCache{redisMockClient},
		),
		cacher.NewAPIKeyCacher(
			accountStorageMock,
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
		deleter.NewTagDeleter(tagStorageMock, featureStorageMock),
		autoarchive.NewFeatureAutoArchiver(
			envStorageMock,
			featureStorageMock,
			codeRefStorageMock,
			featureMockClient,
			jobs.WithTimeout(10*time.Minute),
			jobs.WithLogger(logger),
		),
		scheduledflagchange.NewScheduledFlagChangeExecutor(
			scheduledFlagChangeStorageMock,
			featureMockClient,
			jobs.WithTimeout(50*time.Second),
			jobs.WithLogger(logger),
		),
		monthlysummary.NewMonthlySummarizer(
			environmentMockClient,
			mauCacheMock,
			monthlySummaryStorageMock,
			nil,
			jobs.WithLogger(logger),
		),
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
	c, err := anypb.New(dc)
	require.NoError(t, err)
	return &autoopsproto.ProgressiveRollout{
		Id:        "prID",
		FeatureId: "fID",
		Clause:    c,
		Type:      autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
	}
}
