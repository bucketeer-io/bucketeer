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

package demo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	orgmock "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func newMockDemoDataDeleter(t *testing.T, mockController *gomock.Controller) *demoDataDeleter {
	logger, err := log.NewLogger()
	require.NoError(t, err)

	return &demoDataDeleter{
		mysqlClient:         mysqlmock.NewMockClient(mockController),
		organizationStorage: orgmock.NewMockOrganizationStorage(mockController),
		environmentStorage:  orgmock.NewMockEnvironmentStorage(mockController),
		logger:              logger,
		opts: &jobs.Options{
			Timeout: time.Minute,
		},
	}
}

func TestDemoDataDeleter_Run(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	deleter := newMockDemoDataDeleter(t, mockController)
	err := os.Setenv("DEMO_TRIAL_PERIOD_DAY", "7")
	require.NoError(t, err)

	patterns := []struct {
		desc     string
		setup    func(t *testing.T, d *demoDataDeleter)
		expected error
	}{
		{
			desc: "error: internal",
			setup: func(t *testing.T, d *demoDataDeleter) {
				err := os.Setenv("BUCKETEER_BATCH_DEMO_ENABLED", "true")
				if err != nil {
					return
				}
				d.organizationStorage.(*orgmock.MockOrganizationStorage).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), status.Error(codes.Internal, "internal error"))
			},
			expected: status.Error(codes.Internal, "internal error"),
		},
		{
			desc: "success: demo mode disabled",
			setup: func(t *testing.T, d *demoDataDeleter) {
				err := os.Setenv("BUCKETEER_BATCH_DEMO_ENABLED", "false")
				if err != nil {
					return
				}
			},
		},
		{
			desc: "success: no outdated organizations",
			setup: func(t *testing.T, d *demoDataDeleter) {
				err := os.Setenv("BUCKETEER_BATCH_DEMO_ENABLED", "true")
				if err != nil {
					return
				}
				d.organizationStorage.(*orgmock.MockOrganizationStorage).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Organization{}, 0, int64(0), nil)
			},
		},
		{
			desc: "success: no outdated environments",
			setup: func(t *testing.T, d *demoDataDeleter) {
				err := os.Setenv("BUCKETEER_BATCH_DEMO_ENABLED", "true")
				if err != nil {
					return
				}
				d.organizationStorage.(*orgmock.MockOrganizationStorage).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Organization{
					{Id: "org1"},
					{Id: "org2"},
				}, 0, int64(0), nil)
				d.environmentStorage.(*orgmock.MockEnvironmentStorage).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.EnvironmentV2{}, 0, int64(0), nil)
				d.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				d.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				d.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				for i := 0; i < len(targetEntitiesInOrganization); i++ {
					d.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
						gomock.Any(), gomock.Any(), gomock.Any(),
					).Return(nil, nil)
				}
				d.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
		},
		{
			desc: "success",
			setup: func(t *testing.T, d *demoDataDeleter) {
				err := os.Setenv("BUCKETEER_BATCH_DEMO_ENABLED", "true")
				if err != nil {
					return
				}
				d.organizationStorage.(*orgmock.MockOrganizationStorage).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Organization{
					{Id: "org1"},
					{Id: "org2"},
				}, 0, int64(0), nil)
				d.environmentStorage.(*orgmock.MockEnvironmentStorage).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.EnvironmentV2{
					{Id: "env1", OrganizationId: "org1"},
					{Id: "env2", OrganizationId: "org2"},
				}, 0, int64(0), nil)
				// equal to number of environments
				for i := 0; i < 2; i++ {
					d.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
						gomock.Any(), gomock.Any(),
					).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
						_ = fn(ctx, nil)
					}).Return(nil)
					for j := 0; j < len(targetEntities); j++ {
						d.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
							gomock.Any(), gomock.Any(), gomock.Any(),
						).Return(nil, nil)
					}
				}
				d.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				d.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				d.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				for i := 0; i < len(targetEntitiesInOrganization); i++ {
					d.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
						gomock.Any(), gomock.Any(), gomock.Any(),
					).Return(nil, nil)
				}
				d.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
		},
	}
	for _, pattern := range patterns {
		t.Run(pattern.desc, func(t *testing.T) {
			if pattern.setup != nil {
				pattern.setup(t, deleter)
			}
			ctx := context.Background()
			err := deleter.Run(ctx)
			if pattern.expected != nil {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if err.Error() != pattern.expected.Error() {
					t.Errorf("expected error %v but got %v", pattern.expected, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
