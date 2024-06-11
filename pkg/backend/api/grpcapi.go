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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	backendproto "github.com/bucketeer-io/bucketeer/proto/backend"
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

type BackendService struct {
	opts   *options
	logger *zap.Logger
}

func NewBackendService(
	opts ...Option,
) *BackendService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &BackendService{
		opts:   dopts,
		logger: dopts.logger.Named("api"),
	}
}

func (s *BackendService) GetFeature(
	ctx context.Context,
	req *backendproto.GetFeatureRequest,
) (*backendproto.GetFeatureResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

func (s *BackendService) UpdateFeature(
	ctx context.Context,
	req *backendproto.UpdateFeatureRequest,
) (*backendproto.UpdateFeatureResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}
