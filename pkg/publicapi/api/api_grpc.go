// Copyright 2024 The Bucketeer Authors.
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	gmetadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	publicapiproto "github.com/bucketeer-io/bucketeer/proto/publicapi"
)

var (
	ErrContextCanceled = status.Error(codes.Canceled, "publicAPI: context canceled")
	ErrMissingAPIKey   = status.Error(codes.Unauthenticated, "publicAPI: missing APIKey")
	ErrInvalidAPIKey   = status.Error(codes.PermissionDenied, "publicAPI: invalid APIKey")
	ErrDisabledAPIKey  = status.Error(codes.PermissionDenied, "publicAPI: disabled APIKey")
	ErrBadRole         = status.Error(codes.PermissionDenied, "publicAPI: bad role")
	ErrInternal        = status.Error(codes.Internal, "publicAPI: internal")
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

type PublicAPIService struct {
	accountClient          accountclient.Client
	featureClient          featureclient.Client
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache
	flightgroup            singleflight.Group
	opts                   *options
	logger                 *zap.Logger
}

func NewPublicAPIService(
	accountClient accountclient.Client,
	featureClient featureclient.Client,
	cacher cache.MultiGetCache,
	opts ...Option,
) *PublicAPIService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &PublicAPIService{
		accountClient:          accountClient,
		featureClient:          featureClient,
		environmentAPIKeyCache: cachev3.NewEnvironmentAPIKeyCache(cacher),
		opts:                   dopts,
		logger:                 dopts.logger.Named("publicAPI"),
	}
}

func (s *PublicAPIService) Register(server *grpc.Server) {
	publicapiproto.RegisterPublicAPIServiceServer(server, s)
}

func (s *PublicAPIService) GetFeature(
	ctx context.Context,
	req *publicapiproto.GetFeatureRequest,
) (*publicapiproto.GetFeatureResponse, error) {
	if _, err := s.checkAuth(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	}); err != nil {
		s.logger.Error("Failed to check authentication",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentNamespace: req.EnvironmentId,
		Id:                   req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &publicapiproto.GetFeatureResponse{
		Feature: resp.Feature,
	}, nil
}

func (s *PublicAPIService) UpdateFeature(
	ctx context.Context,
	req *publicapiproto.UpdateFeatureRequest,
) (*publicapiproto.UpdateFeatureResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

func (s *PublicAPIService) checkAuth(
	ctx context.Context,
	roles []accountproto.APIKey_Role,
) (*accountproto.EnvironmentAPIKey, error) {
	if ctx.Err() == context.Canceled {
		s.logger.Warn(
			"Request was canceled",
			log.FieldsFromImcomingContext(ctx)...,
		)
		return nil, ErrContextCanceled
	}
	id, err := s.extractAPIKeyID(ctx)
	if err != nil {
		s.logger.Error("Failed to extract API key ID",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	envAPIKey, err := s.getEnvironmentAPIKey(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get environment API key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("apiKey", id),
			)...,
		)
		return nil, err
	}
	if err := checkEnvironmentAPIKey(envAPIKey, roles); err != nil {
		s.logger.Error("Failed to check environment API key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("envAPIKey", envAPIKey),
			)...,
		)
		return nil, err
	}
	return envAPIKey, nil
}

func (s *PublicAPIService) extractAPIKeyID(ctx context.Context) (string, error) {
	md, ok := gmetadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMissingAPIKey
	}
	keys, ok := md["authorization"]
	if !ok || len(keys) == 0 || keys[0] == "" {
		return "", ErrMissingAPIKey
	}
	return keys[0], nil
}

func environmentAPIKeyFlightID(id string) string {
	return id
}

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func (s *PublicAPIService) getEnvironmentAPIKeyFromCache(
	ctx context.Context,
	id string,
) (*accountproto.EnvironmentAPIKey, error) {
	envAPIKey, err := s.environmentAPIKeyCache.Get(id)
	if err != nil {
		return nil, err
	}
	return envAPIKey, nil
}

func (s *PublicAPIService) getEnvironmentAPIKeyFromAccountService(
	ctx context.Context,
	id string,
) (*accountproto.EnvironmentAPIKey, error) {
	resp, err := s.accountClient.GetAPIKeyBySearchingAllEnvironments(
		ctx,
		&accountproto.GetAPIKeyBySearchingAllEnvironmentsRequest{Id: id},
	)
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			return nil, ErrInvalidAPIKey
		}
		s.logger.Error(
			"Failed to get environment APIKey from account service",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, ErrInternal
	}
	return resp.EnvironmentApiKey, nil
}

func (s *PublicAPIService) getEnvironmentAPIKey(
	ctx context.Context,
	apiKey string,
) (*accountproto.EnvironmentAPIKey, error) {
	k, err, _ := s.flightgroup.Do(
		environmentAPIKeyFlightID(apiKey),
		func() (interface{}, error) {
			envAPIKey, err := s.getEnvironmentAPIKeyFromCache(ctx, apiKey)
			if err == nil {
				return envAPIKey, nil
			}
			envAPIKey, err = s.getEnvironmentAPIKeyFromAccountService(ctx, apiKey)
			if err != nil {
				return nil, err
			}
			return envAPIKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	envAPIKey := k.(*accountproto.EnvironmentAPIKey)
	return envAPIKey, nil
}

func checkEnvironmentAPIKey(
	environmentAPIKey *accountproto.EnvironmentAPIKey,
	roles []accountproto.APIKey_Role,
) error {
	if environmentAPIKey.EnvironmentDisabled {
		return ErrDisabledAPIKey
	}
	if environmentAPIKey.ApiKey.Disabled {
		return ErrDisabledAPIKey
	}
	if !contains(roles, environmentAPIKey.ApiKey.Role) {
		return ErrBadRole
	}
	return nil
}
