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
//

package api

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/experiment"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/notification"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/opsevent"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	environmentdomain "github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client/mock"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	notificationsender "github.com/bucketeer-io/bucketeer/pkg/notification/sender/mock"
	opsexecutor "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor/mock"
	targetstoremock "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/targetstore/mock"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	batchproto "github.com/bucketeer-io/bucketeer/proto/batch"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	jpLocation            = time.FixedZone("Asia/Tokyo", 9*60*60)
	pullerRunningDuration = 1 * time.Second
)

type setupMockFunc func(
	environmentMockClient *environmentclient.MockClient,
	experimentMockClient *experimentclient.MockClient,
	featureMockClient *featureclientmock.MockClient,
	eventCounterMockClient *ecclient.MockClient,
	notificationMockSender *notificationsender.MockSender,
	mockTargetStore *targetstoremock.MockTargetStore,
	mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
	domainMockEventPuller *domainEventPullerMock,
	mysqlMockClient *mysqlmock.MockClient)

type domainEventPullerMock struct{}

func (d domainEventPullerMock) Pull(
	ctx context.Context,
	f func(context.Context, *puller.Message)) error {
	// using timer.Ticker to mock pubsub event
	timer := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			event := domaineventproto.Event{
				IsAdminEvent: false,
				Timestamp:    int64(time.Now().Nanosecond()),
			}
			data, _ := proto.Marshal(&event)
			f(ctx, &puller.Message{
				Attributes: map[string]string{
					"id": "event",
				},
				Data: data,
				Ack:  func() {},
			})
		}
	}
}

func TestExperimentStatusUpdater(t *testing.T) {
	setupMock := func(
		environmentMockClient *environmentclient.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockTargetStore *targetstoremock.MockTargetStore,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		domainMockEventPuller *domainEventPullerMock,
		mysqlMockClient *mysqlmock.MockClient) {
		environmentMockClient.EXPECT().
			ListEnvironments(gomock.Any(), gomock.Any()).
			Return(
				&environmentproto.ListEnvironmentsResponse{
					Environments: getEnvironments(t),
				},
				nil,
			)
		experimentMockClient.EXPECT().
			ListExperiments(gomock.Any(), gomock.Any()).
			Times(4).
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
	setupMock := func(
		environmentMockClient *environmentclient.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockTargetStore *targetstoremock.MockTargetStore,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		domainMockEventPuller *domainEventPullerMock,
		mysqlMockClient *mysqlmock.MockClient) {
		environmentMockClient.EXPECT().
			ListEnvironments(gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&environmentproto.ListEnvironmentsResponse{
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
	setupMock := func(
		environmentMockClient *environmentclient.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockTargetStore *targetstoremock.MockTargetStore,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		domainMockEventPuller *domainEventPullerMock,
		mysqlMockClient *mysqlmock.MockClient) {
		environmentMockClient.EXPECT().
			ListEnvironments(gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&environmentproto.ListEnvironmentsResponse{
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
	setupMock := func(
		environmentMockClient *environmentclient.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockTargetStore *targetstoremock.MockTargetStore,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		domainMockEventPuller *domainEventPullerMock,
		mysqlMockClient *mysqlmock.MockClient) {
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
			ListEnvironments(gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&environmentproto.ListEnvironmentsResponse{
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
	setupMock := func(
		environmentMockClient *environmentclient.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockTargetStore *targetstoremock.MockTargetStore,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		domainMockEventPuller *domainEventPullerMock,
		mysqlMockClient *mysqlmock.MockClient) {
		mockTargetStore.EXPECT().
			GetEnvironments(gomock.Any()).
			Times(1).
			Return([]*environmentdomain.Environment{
				environmentdomain.NewEnvironment("eid0", "env0", "pj0"),
				environmentdomain.NewEnvironment("eid1", "env1", "pj1"),
			})
		mockTargetStore.EXPECT().
			GetAutoOpsRules(gomock.Any(), gomock.Any()).
			Times(2).
			Return([]*autoopsdomain.AutoOpsRule{
				newAutoOpsRule(t),
			})
		mockAutoOpsExecutor.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)

	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_DatetimeWatcher,
	}, setupMock)
}

func TestEventCountWatcher(t *testing.T) {
	setupMock := func(
		environmentMockClient *environmentclient.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockTargetStore *targetstoremock.MockTargetStore,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		domainMockEventPuller *domainEventPullerMock,
		mysqlMockClient *mysqlmock.MockClient) {
		mockTargetStore.EXPECT().
			GetEnvironments(gomock.Any()).
			Times(1).
			Return([]*environmentdomain.Environment{
				environmentdomain.NewEnvironment("eid0", "env0", "pj0"),
				environmentdomain.NewEnvironment("eid1", "env1", "pj1"),
			})
		mockTargetStore.EXPECT().
			GetAutoOpsRules(gomock.Any(), gomock.Any()).
			Times(2).
			Return([]*autoopsdomain.AutoOpsRule{
				newAutoOpsRule(t),
			})
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
		mockAutoOpsExecutor.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)

	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_EventCountWatcher,
	}, setupMock)
}

func TestDomainEventInformer(t *testing.T) {
	setupMock := func(
		environmentMockClient *environmentclient.MockClient,
		experimentMockClient *experimentclient.MockClient,
		featureMockClient *featureclientmock.MockClient,
		eventCounterMockClient *ecclient.MockClient,
		notificationMockSender *notificationsender.MockSender,
		mockTargetStore *targetstoremock.MockTargetStore,
		mockAutoOpsExecutor *opsexecutor.MockAutoOpsExecutor,
		domainMockEventPuller *domainEventPullerMock,
		mysqlMockClient *mysqlmock.MockClient) {
		environmentMockClient.EXPECT().
			GetEnvironmentByNamespace(gomock.Any(), gomock.Any()).
			MinTimes(3).
			MaxTimes(5).
			Return(
				&environmentproto.GetEnvironmentByNamespaceResponse{
					Environment: &environmentproto.Environment{
						Id: "eid0",
					}},
				nil,
			)
		notificationMockSender.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			MinTimes(3).
			MaxTimes(5).
			Return(nil)
	}
	executeMockBatchJob(t, &batchproto.BatchJobRequest{
		Job: batchproto.BatchJob_DomainEventInformer,
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

	environmentMockClient := environmentclient.NewMockClient(mockController)
	experimentMockClient := experimentclient.NewMockClient(mockController)
	featureMockClient := featureclientmock.NewMockClient(mockController)
	eventCounterMockClient := ecclient.NewMockClient(mockController)
	notificationMockSender := notificationsender.NewMockSender(mockController)
	mockTargetStore := targetstoremock.NewMockTargetStore(mockController)
	mockAutoOpsExecutor := opsexecutor.NewMockAutoOpsExecutor(mockController)
	domainMockEventPuller := &domainEventPullerMock{}
	mysqlMockClient := mysqlmock.NewMockClient(mockController)

	setupMock(
		environmentMockClient,
		experimentMockClient,
		featureMockClient,
		eventCounterMockClient,
		notificationMockSender,
		mockTargetStore,
		mockAutoOpsExecutor,
		domainMockEventPuller,
		mysqlMockClient,
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
			mockTargetStore,
			mockAutoOpsExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewEventCountWatcher(
			mysqlMockClient,
			mockTargetStore,
			eventCounterMockClient,
			featureMockClient,
			mockAutoOpsExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		environmentMockClient,
		domainMockEventPuller,
		notificationMockSender,
		logger,
		notification.WithRunningDurationPerBatch(pullerRunningDuration),
	)
	return service
}

func getEnvironments(t *testing.T) []*environmentproto.Environment {
	t.Helper()
	return []*environmentproto.Environment{
		{Id: "ns0", Namespace: "ns0", ProjectId: "pj0"},
		{Id: "ns1", Namespace: "ns1", ProjectId: "pj0"},
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

func newAutoOpsRule(t *testing.T) *autoopsdomain.AutoOpsRule {
	oerc1 := &autoopsproto.OpsEventRateClause{
		GoalId:          "gid",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	oerc2 := &autoopsproto.OpsEventRateClause{
		GoalId:          "gid",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	dc1 := &autoopsproto.DatetimeClause{
		Time: 1000000001,
	}
	dc2 := &autoopsproto.DatetimeClause{
		Time: 1000000002,
	}
	aor, err := autoopsdomain.NewAutoOpsRule(
		"fid",
		autoopsproto.OpsType_ENABLE_FEATURE,
		[]*autoopsproto.OpsEventRateClause{oerc1, oerc2},
		[]*autoopsproto.DatetimeClause{dc1, dc2},
		[]*autoopsproto.WebhookClause{},
	)
	require.NoError(t, err)
	return aor
}
