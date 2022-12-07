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
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/user"
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

type userService struct {
	storageClient mysql.Client
	accountClient accountclient.Client
	opts          *options
	logger        *zap.Logger
}

func NewUserService(
	client mysql.Client,
	accountClient accountclient.Client,
	opts ...Option,
) rpc.Service {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &userService{
		storageClient: client,
		accountClient: accountClient,
		opts:          &options,
		logger:        options.logger.Named("api"),
	}
}

func (s *userService) Register(server *grpc.Server) {
	proto.RegisterUserServiceServer(server, s)
}

func (s *userService) checkRole(
	ctx context.Context,
	requiredRole accountproto.Account_Role,
	environmentNamespace string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckRole(ctx, requiredRole, func(email string) (*accountproto.GetAccountResponse, error) {
		return s.accountClient.GetAccount(ctx, &accountproto.GetAccountRequest{
			Email:                email,
			EnvironmentNamespace: environmentNamespace,
		})
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
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Info(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}
	return editor, nil
}
