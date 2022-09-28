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

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

const (
	listRequestPageSize = 500
)

type options struct {
	logger *zap.Logger
}

var defaultOptions = options{
	logger: zap.NewNop(),
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type AccountService struct {
	environmentClient environmentclient.Client
	mysqlClient       mysql.Client
	publisher         publisher.Publisher
	opts              *options
	logger            *zap.Logger
}

func NewAccountService(
	e environmentclient.Client,
	mysqlClient mysql.Client,
	publisher publisher.Publisher,
	opts ...Option,
) *AccountService {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &AccountService{
		environmentClient: e,
		mysqlClient:       mysqlClient,
		publisher:         publisher,
		opts:              &options,
		logger:            options.logger.Named("api"),
	}
}

func (s *AccountService) Register(server *grpc.Server) {
	proto.RegisterAccountServiceServer(server, s)
}

func (s *AccountService) makeProjectSet(projects []*environmentproto.Project) map[string]*environmentproto.Project {
	projectSet := make(map[string]*environmentproto.Project)
	for _, p := range projects {
		projectSet[p.Id] = p
	}
	return projectSet
}

func (s *AccountService) listProjects(ctx context.Context) ([]*environmentproto.Project, error) {
	projects := []*environmentproto.Project{}
	cursor := ""
	for {
		resp, err := s.environmentClient.ListProjects(ctx, &environmentproto.ListProjectsRequest{
			PageSize: listRequestPageSize,
			Cursor:   cursor,
		})
		if err != nil {
			return nil, err
		}
		projects = append(projects, resp.Projects...)
		projectSize := len(resp.Projects)
		if projectSize == 0 || projectSize < listRequestPageSize {
			return projects, nil
		}
		cursor = resp.Cursor
	}
}

func (s *AccountService) listEnvironments(ctx context.Context) ([]*environmentproto.Environment, error) {
	environments := []*environmentproto.Environment{}
	cursor := ""
	for {
		resp, err := s.environmentClient.ListEnvironments(ctx, &environmentproto.ListEnvironmentsRequest{
			PageSize: listRequestPageSize,
			Cursor:   cursor,
		})
		if err != nil {
			return nil, err
		}
		environments = append(environments, resp.Environments...)
		environmentSize := len(resp.Environments)
		if environmentSize == 0 || environmentSize < listRequestPageSize {
			return environments, nil
		}
		cursor = resp.Cursor
	}
}

func (s *AccountService) checkAdminRole(ctx context.Context) (*eventproto.Editor, error) {
	editor, err := role.CheckAdminRole(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Info(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, localizedError(statusUnauthenticated, locale.JaJP)
		case codes.PermissionDenied:
			s.logger.Info(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, localizedError(statusPermissionDenied, locale.JaJP)
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, localizedError(statusInternal, locale.JaJP)
		}
	}
	return editor, nil
}

func (s *AccountService) checkRole(
	ctx context.Context,
	requiredRole proto.Account_Role,
	environmentNamespace string,
) (*eventproto.Editor, error) {
	editor, err := role.CheckRole(ctx, requiredRole, func(email string) (*proto.GetAccountResponse, error) {
		account, err := s.getAccount(ctx, email, environmentNamespace)
		if err != nil {
			return nil, err
		}
		return &proto.GetAccountResponse{Account: account.Account}, nil
	})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Info(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return nil, localizedError(statusUnauthenticated, locale.JaJP)
		case codes.PermissionDenied:
			s.logger.Info(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return nil, localizedError(statusPermissionDenied, locale.JaJP)
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return nil, localizedError(statusInternal, locale.JaJP)
		}
	}
	return editor, nil
}
