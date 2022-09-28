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
	"golang.org/x/sync/singleflight"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featurestorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	bigtable "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
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

type FeatureService struct {
	mysqlClient           mysql.Client
	userEvaluationStorage featurestorage.UserEvaluationsStorage
	accountClient         accountclient.Client
	experimentClient      experimentclient.Client
	featuresCache         cachev3.FeaturesCache
	segmentUsersCache     cachev3.SegmentUsersCache
	segmentUsersPublisher publisher.Publisher
	domainPublisher       publisher.Publisher
	flightgroup           singleflight.Group
	opts                  *options
	logger                *zap.Logger
}

func NewFeatureService(
	mysqlClient mysql.Client,
	btClient bigtable.Client,
	accountClient accountclient.Client,
	experimentClient experimentclient.Client,
	v3Cache cache.MultiGetCache,
	segmentUsersPublisher publisher.Publisher,
	domainPublisher publisher.Publisher,
	opts ...Option,
) *FeatureService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &FeatureService{
		mysqlClient:           mysqlClient,
		userEvaluationStorage: featurestorage.NewUserEvaluationsStorage(btClient),
		accountClient:         accountClient,
		experimentClient:      experimentClient,
		featuresCache:         cachev3.NewFeaturesCache(v3Cache),
		segmentUsersCache:     cachev3.NewSegmentUsersCache(v3Cache),
		segmentUsersPublisher: segmentUsersPublisher,
		domainPublisher:       domainPublisher,
		opts:                  dopts,
		logger:                dopts.logger.Named("api"),
	}
}

func (s *FeatureService) Register(server *grpc.Server) {
	featureproto.RegisterFeatureServiceServer(server, s)
}

func (s *FeatureService) checkRole(
	ctx context.Context,
	requiredRole accountproto.Account_Role,
	environmentNamespace string,
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

func (s *FeatureService) reportInternalServerError(
	ctx context.Context,
	err error,
	environmentNamespace string,
	localizer locale.Localizer,
) error {
	s.logger.Error(
		"Internal server error",
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
		return statusInternal.Err()
	}
	return dt.Err()
}
