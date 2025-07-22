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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	acmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	accstoragemock "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/auth"
	storagemock "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

func TestNewEnvironmentService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ac := acmock.NewMockClient(mockController)
	mysqlClient := mysqlmock.NewMockClient(mockController)
	p := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewEnvironmentService(
		ac, mysqlClient, p, &auth.OAuthConfig{}, "", "", nil, nil, WithLogger(logger))
	assert.IsType(t, &EnvironmentService{}, s)
}

func createContextWithToken(t *testing.T) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: true,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, token)
}

func createDemoContextWithToken(t *testing.T) context.Context {
	t.Helper()
	token := &token.DemoCreationToken{
		Issuer:   "issuer",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "test@example.com",
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.DemoCreationTokenKey, token)
}

func createContextWithTokenRoleUnassigned(t *testing.T) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:   "issuer",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, token)
}
func newEnvironmentService(t *testing.T, mockController *gomock.Controller, s storage.Client) *EnvironmentService {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &EnvironmentService{
		accountClient:      acmock.NewMockClient(mockController),
		mysqlClient:        mysqlmock.NewMockClient(mockController),
		orgStorage:         storagemock.NewMockOrganizationStorage(mockController),
		projectStorage:     storagemock.NewMockProjectStorage(mockController),
		environmentStorage: storagemock.NewMockEnvironmentStorage(mockController),
		accountStorage:     accstoragemock.NewMockAccountStorage(mockController),
		publisher:          publishermock.NewMockPublisher(mockController),
		logger:             logger.Named("api"),
	}
}

func newDemoEnvironmentService(t *testing.T, mockController *gomock.Controller, s storage.Client) *EnvironmentService {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &EnvironmentService{
		accountClient:      acmock.NewMockClient(mockController),
		mysqlClient:        mysqlmock.NewMockClient(mockController),
		orgStorage:         storagemock.NewMockOrganizationStorage(mockController),
		projectStorage:     storagemock.NewMockProjectStorage(mockController),
		environmentStorage: storagemock.NewMockEnvironmentStorage(mockController),
		accountStorage:     accstoragemock.NewMockAccountStorage(mockController),
		publisher:          publishermock.NewMockPublisher(mockController),
		logger:             logger.Named("api"),
		opts: &options{
			logger:            logger,
			isDemoSiteEnabled: true,
		},
	}
}
