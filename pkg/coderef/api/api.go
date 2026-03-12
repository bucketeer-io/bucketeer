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
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	mysql "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/coderef"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type CodeReferenceService struct {
	accountClient accountclient.Client
	mysqlClient   mysql.Client
	publisher     publisher.Publisher
	opts          *options
	logger        *zap.Logger
}

func NewCodeReferenceService(
	ac accountclient.Client,
	mysqlClient mysql.Client,
	p publisher.Publisher,
	opts ...Option,
) *CodeReferenceService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &CodeReferenceService{
		accountClient: ac,
		mysqlClient:   mysqlClient,
		publisher:     p,
		opts:          dopts,
		logger:        dopts.logger.Named("api"),
	}
}

func (s *CodeReferenceService) Register(server *grpc.Server) {
	proto.RegisterCodeReferenceServiceServer(server, s)
}

func (s *CodeReferenceService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentID string,
) (*eventproto.Editor, error) {
	return role.CheckEnvironmentRoleWithLog(
		ctx,
		requiredRole,
		environmentID,
		func(email string) (*accountproto.AccountV2, error) {
			resp, err := s.accountClient.GetAccountV2ByEnvironmentID(ctx, &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         email,
				EnvironmentId: environmentID,
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
