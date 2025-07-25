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

package api

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	v2ts "github.com/bucketeer-io/bucketeer/pkg/tag/storage"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	btclient "github.com/bucketeer-io/bucketeer/pkg/batch/client"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
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
	fluiStorage           v2fs.FeatureLastUsedInfoStorage
	flagTriggerStorage    v2fs.FlagTriggerStorage
	featureStorage        v2fs.FeatureStorage
	segmentStorage        v2fs.SegmentStorage
	segmentUserStorage    v2fs.SegmentUserStorage
	tagStorage            v2ts.TagStorage
	mysqlClient           mysql.Client
	accountClient         accountclient.Client
	experimentClient      experimentclient.Client
	featuresCache         cachev3.FeaturesCache
	autoOpsClient         autoopsclient.Client
	batchClient           btclient.Client
	environmentClient     envclient.Client
	segmentUsersCache     cachev3.SegmentUsersCache
	userAttributesCache   cachev3.UserAttributesCache
	segmentUsersPublisher publisher.Publisher
	domainPublisher       publisher.Publisher
	flightgroup           singleflight.Group
	triggerURL            string
	opts                  *options
	logger                *zap.Logger
}

func NewFeatureService(
	mysqlClient mysql.Client,
	accountClient accountclient.Client,
	experimentClient experimentclient.Client,
	autoOpsClient autoopsclient.Client,
	batchClient btclient.Client,
	environmentClient envclient.Client,
	v3Cache cache.MultiGetDeleteCountCache,
	persistentRedisV3Cache cache.MultiGetDeleteCountCache,
	segmentUsersPublisher publisher.Publisher,
	domainPublisher publisher.Publisher,
	triggerURL string,
	opts ...Option,
) *FeatureService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &FeatureService{
		fluiStorage:           v2fs.NewFeatureLastUsedInfoStorage(mysqlClient),
		flagTriggerStorage:    v2fs.NewFlagTriggerStorage(mysqlClient),
		featureStorage:        v2fs.NewFeatureStorage(mysqlClient),
		segmentStorage:        v2fs.NewSegmentStorage(mysqlClient),
		segmentUserStorage:    v2fs.NewSegmentUserStorage(mysqlClient),
		tagStorage:            v2ts.NewTagStorage(mysqlClient),
		mysqlClient:           mysqlClient,
		accountClient:         accountClient,
		experimentClient:      experimentClient,
		autoOpsClient:         autoOpsClient,
		batchClient:           batchClient,
		environmentClient:     environmentClient,
		featuresCache:         cachev3.NewFeaturesCache(v3Cache),
		segmentUsersCache:     cachev3.NewSegmentUsersCache(v3Cache),
		userAttributesCache:   cachev3.NewUserAttributesCache(persistentRedisV3Cache),
		segmentUsersPublisher: segmentUsersPublisher,
		domainPublisher:       domainPublisher,
		triggerURL:            triggerURL,
		opts:                  dopts,
		logger:                dopts.logger.Named("api"),
	}
}

func (s *FeatureService) Register(server *grpc.Server) {
	featureproto.RegisterFeatureServiceServer(server, s)
}

func (s *FeatureService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentId string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckEnvironmentRole(
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
		})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
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
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
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
					zap.String("environmentId", environmentId),
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

func (s *FeatureService) reportInternalServerError(
	ctx context.Context,
	err error,
	environmentId string,
	localizer locale.Localizer,
) error {
	s.logger.Error(
		"Internal server error",
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentId", environmentId),
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
