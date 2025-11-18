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

	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	auditlogstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2/mock"
	ecmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	teamstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/team/storage/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
)

func TestWithLogger(t *testing.T) {
	t.Parallel()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	f := WithLogger(logger)
	opt := &options{}
	f(opt)
	assert.Equal(t, logger, opt.logger)
}

func TestNewAccountService(t *testing.T) {
	t.Parallel()
	g := NewAccountService(nil, nil, nil, nil)
	assert.IsType(t, &AccountService{}, g)
}

func createAccountService(t *testing.T, mockController *gomock.Controller, db storage.Client) *AccountService {
	t.Helper()
	logger := zap.NewNop()
	return &AccountService{
		environmentClient:    ecmock.NewMockClient(mockController),
		mysqlClient:          mysqlmock.NewMockClient(mockController),
		accountStorage:       storagemock.NewMockAccountStorage(mockController),
		teamStorage:          teamstoragemock.NewMockTeamStorage(mockController),
		adminAuditLogStorage: auditlogstoragemock.NewMockAdminAuditLogStorage(mockController),
		publisher:            publishermock.NewMockPublisher(mockController),
		logger:               logger.Named("api"),
	}
}

func createContextWithDefaultToken(t *testing.T, isSystemAdmin bool) context.Context {
	t.Helper()
	return createContextWithEmailToken(t, "bucketeer@example.com", isSystemAdmin)
}

func createContextWithEmailToken(t *testing.T, email string, isSystemAdmin bool) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         email,
		IsSystemAdmin: isSystemAdmin,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, token)
}

func createContextWithInvalidEmailToken(t *testing.T) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:   "issuer",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "bucketeer@",
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, token)
}
