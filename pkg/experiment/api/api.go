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

	"go.uber.org/zap"
	"google.golang.org/grpc"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	storage "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type experimentService struct {
	featureClient     featureclient.Client
	accountClient     accountclient.Client
	autoOpsClient     autoopsclient.Client
	mysqlClient       mysql.Client
	experimentStorage storage.ExperimentStorage
	goalStorage       storage.GoalStorage
	publisher         publisher.Publisher
	opts              *options
	logger            *zap.Logger
}

func NewExperimentService(
	featureClient featureclient.Client,
	accountClient accountclient.Client,
	autoOpsClient autoopsclient.Client,
	mysqlClient mysql.Client,
	publisher publisher.Publisher,
	opts ...Option,
) rpc.Service {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &experimentService{
		featureClient:     featureClient,
		accountClient:     accountClient,
		autoOpsClient:     autoOpsClient,
		mysqlClient:       mysqlClient,
		experimentStorage: storage.NewExperimentStorage(mysqlClient),
		goalStorage:       storage.NewGoalStorage(mysqlClient),
		publisher:         publisher,
		opts:              dopts,
		logger:            dopts.logger.Named("api"),
	}
}

func (s *experimentService) Register(server *grpc.Server) {
	proto.RegisterExperimentServiceServer(server, s)
}

func (s *experimentService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentId string,
) (*eventproto.Editor, error) {
	return role.CheckEnvironmentRoleWithLog(
		ctx,
		requiredRole,
		environmentId,
		func(email string) (*accountproto.AccountV2, error) {
			resp, err := s.accountClient.GetAccountV2ByEnvironmentID(ctx, &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         email,
				EnvironmentId: environmentId,
			})
			if err != nil {
				return nil, err
			}
			return resp.Account, nil
		},
		s.logger,
		statusUnauthenticated.Err(),
		statusPermissionDenied.Err(),
		func(err error) error { return api.NewGRPCStatus(err).Err() },
	)
}
