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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	acclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	accstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	aoclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	ecclient "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/client/mock"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	featurestoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	opsexecutor "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/batch/executor/mock"
	opseventstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	subscriptionsender "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/sender/mock"
	tagstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	batchproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestExecuteBatchJobRequiresSystemAdminRole(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		ctx         context.Context
		job         batchproto.BatchJob
		expectedErr error
	}{
		{
			desc:        "rejects request without access token",
			ctx:         context.Background(),
			job:         batchproto.BatchJob_FeatureFlagCacher,
			expectedErr: role.ErrUnauthenticated,
		},
		{
			desc:        "rejects request without access token before resolving the job",
			ctx:         context.Background(),
			job:         batchproto.BatchJob(-1),
			expectedErr: role.ErrUnauthenticated,
		},
		{
			desc: "rejects organization owner",
			ctx: createContextWithToken(t, &token.AccessToken{
				Email:          "owner@example.com",
				OrganizationID: "org0",
			}),
			job:         batchproto.BatchJob_FeatureFlagCacher,
			expectedErr: role.ErrPermissionDenied,
		},
		{
			desc: "rejects service token flag without system admin",
			ctx: createContextWithToken(t, &token.AccessToken{
				Email:          "service@example.com",
				OrganizationID: "org0",
				IsServiceToken: true,
			}),
			job:         batchproto.BatchJob_FeatureFlagCacher,
			expectedErr: role.ErrPermissionDenied,
		},
		{
			desc: "allows system admin",
			ctx: createContextWithToken(t, &token.AccessToken{
				Email:         "admin@example.com",
				IsSystemAdmin: true,
			}),
			job: batchproto.BatchJob_FeatureFlagCacher,
		},
		{
			desc: "allows service token",
			ctx:  createContextWithServiceToken(t),
			job:  batchproto.BatchJob_FeatureFlagCacher,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			controller := gomock.NewController(t)
			defer controller.Finish()
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
				// No expectation is registered for the rejected requests, so any
				// storage read or cache write made by the job fails the test.
				if p.expectedErr != nil {
					return
				}
				featureStorageMock.EXPECT().ListAllEnvironmentFeatures(
					gomock.Any(),
				).Return([]*featureproto.EnvironmentFeature{}, nil)
				redisMockClient.EXPECT().
					Put(gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(nil)
			}
			service := newBatchService(t, controller, setupMock)

			resp, err := service.ExecuteBatchJob(p.ctx, &batchproto.BatchJobRequest{Job: p.job})
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr, err)
				assert.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, &batchproto.BatchJobResponse{}, resp)
		})
	}
}

func createContextWithToken(t *testing.T, accessToken *token.AccessToken) context.Context {
	t.Helper()
	return context.WithValue(context.Background(), rpc.AccessTokenKey, accessToken)
}

// createContextWithServiceToken returns the context the CronJobs and the
// subscriber call the batch service with: a system admin service token.
func createContextWithServiceToken(t *testing.T) context.Context {
	t.Helper()
	return createContextWithToken(t, &token.AccessToken{
		Email:          "service@example.com",
		IsSystemAdmin:  true,
		IsServiceToken: true,
	})
}
