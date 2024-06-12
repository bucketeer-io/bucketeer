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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3mock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestCheckAuth(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*PublicAPIService) context.Context
		input       []accountproto.APIKey_Role
		expected    *accountproto.EnvironmentAPIKey
		expectedErr error
	}{
		{
			desc: "error: context canceled",
			setup: func(s *PublicAPIService) context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			input: []accountproto.APIKey_Role{
				accountproto.APIKey_PUBLIC_API_READ_ONLY,
			},
			expected:    nil,
			expectedErr: ErrContextCanceled,
		},
		{
			desc: "error: invalid api key",
			setup: func(s *PublicAPIService) context.Context {
				s.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				s.accountClient.(*accountclientmock.MockClient).EXPECT().GetAPIKeyBySearchingAllEnvironments(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.NotFound, "error: apy key not found"))
				return metadata.NewIncomingContext(context.TODO(), metadata.MD{
					"authorization": []string{"test-key"},
				})
			},
			input: []accountproto.APIKey_Role{
				accountproto.APIKey_PUBLIC_API_READ_ONLY,
			},
			expected:    nil,
			expectedErr: ErrInvalidAPIKey,
		},
		{
			desc: "error: internal",
			setup: func(s *PublicAPIService) context.Context {
				s.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				s.accountClient.(*accountclientmock.MockClient).EXPECT().GetAPIKeyBySearchingAllEnvironments(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.Internal, "error: internal"))
				return metadata.NewIncomingContext(context.TODO(), metadata.MD{
					"authorization": []string{"test-key"},
				})
			},
			input: []accountproto.APIKey_Role{
				accountproto.APIKey_PUBLIC_API_READ_ONLY,
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "error: api key missing",
			setup: func(s *PublicAPIService) context.Context {
				return context.Background()
			},
			input: []accountproto.APIKey_Role{
				accountproto.APIKey_PUBLIC_API_READ_ONLY,
			},
			expected:    nil,
			expectedErr: ErrMissingAPIKey,
		},
		{
			desc: "error: api key disabled",
			setup: func(s *PublicAPIService) context.Context {
				s.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment:         &environmentproto.EnvironmentV2{Id: "ns0"},
						EnvironmentDisabled: true,
						ApiKey: &accountproto.APIKey{
							Id:   "test-key",
							Role: accountproto.APIKey_PUBLIC_API_READ_ONLY,
						},
					}, nil)
				return metadata.NewIncomingContext(context.TODO(), metadata.MD{
					"authorization": []string{"test-key"},
				})
			},
			input: []accountproto.APIKey_Role{
				accountproto.APIKey_PUBLIC_API_READ_ONLY,
			},
			expected:    nil,
			expectedErr: ErrDisabledAPIKey,
		},
		{
			desc: "success",
			setup: func(s *PublicAPIService) context.Context {
				s.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:   "test-key",
							Role: accountproto.APIKey_PUBLIC_API_READ_ONLY,
						},
					}, nil)
				return metadata.NewIncomingContext(context.TODO(), metadata.MD{
					"authorization": []string{"test-key"},
				})
			},
			input: []accountproto.APIKey_Role{
				accountproto.APIKey_PUBLIC_API_READ_ONLY,
			},
			expected: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey: &accountproto.APIKey{
					Id:   "test-key",
					Role: accountproto.APIKey_PUBLIC_API_READ_ONLY,
				},
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			p := p
			s := newPublicAPIServiceWithMock(t, mockController)
			ctx := p.setup(s)
			actual, err := s.checkAuth(ctx, p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newPublicAPIServiceWithMock(t *testing.T, mockController *gomock.Controller) *PublicAPIService {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &PublicAPIService{
		featureClient:          featureclientmock.NewMockClient(mockController),
		accountClient:          accountclientmock.NewMockClient(mockController),
		environmentAPIKeyCache: cachev3mock.NewMockEnvironmentAPIKeyCache(mockController),
		logger:                 logger,
	}
}
