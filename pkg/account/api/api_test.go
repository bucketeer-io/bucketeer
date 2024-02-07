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
	"encoding/base64"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	storagemock "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2/mock"
	ecmock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
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
	g := NewAccountService(nil, nil, nil)
	assert.IsType(t, &AccountService{}, g)
}

func createAccountService(t *testing.T, mockController *gomock.Controller, db storage.Client) *AccountService {
	t.Helper()
	logger := zap.NewNop()
	return &AccountService{
		environmentClient: ecmock.NewMockClient(mockController),
		accountStorage:    storagemock.NewMockAccountStorage(mockController),
		publisher:         publishermock.NewMockPublisher(mockController),
		logger:            logger.Named("api"),
	}
}

func createContextWithDefaultToken(t *testing.T, role accountproto.Account_Role, isSystemAdmin bool) context.Context {
	t.Helper()
	return createContextWithEmailToken(t, "bucketeer@example.com", role, isSystemAdmin)
}

func createContextWithEmailToken(t *testing.T, email string, role accountproto.Account_Role, isSystemAdmin bool) context.Context {
	t.Helper()
	sub := &authproto.IDTokenSubject{
		UserId: email,
		ConnId: "test-connector-id",
	}
	data, err := proto.Marshal(sub)
	require.NoError(t, err)
	token := &token.IDToken{
		Issuer:        "issuer",
		Subject:       base64.RawURLEncoding.EncodeToString([]byte(data)),
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         email,
		IsSystemAdmin: isSystemAdmin,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

func createContextWithInvalidSubjectToken(t *testing.T, role accountproto.Account_Role) context.Context {
	t.Helper()
	token := &token.IDToken{
		Issuer:   "issuer",
		Subject:  base64.RawURLEncoding.EncodeToString([]byte("bucketeer@example.com")),
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "bucketeer@example.com",
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

func createContextWithInvalidEmailToken(t *testing.T, role accountproto.Account_Role) context.Context {
	t.Helper()
	sub := &authproto.IDTokenSubject{
		UserId: "bucketeer@example.com",
		ConnId: "test-connector-id",
	}
	data, err := proto.Marshal(sub)
	require.NoError(t, err)
	token := &token.IDToken{
		Issuer:   "issuer",
		Subject:  base64.RawURLEncoding.EncodeToString([]byte(data)),
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "bucketeer@",
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}
