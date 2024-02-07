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
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	md "google.golang.org/grpc/metadata"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestNewExperimentService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	featureClientMock := featureclientmock.NewMockClient(mockController)
	accountClientMock := accountclientmock.NewMockClient(mockController)
	mysqlClient := mysqlmock.NewMockClient(mockController)
	p := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewExperimentService(
		featureClientMock,
		accountClientMock,
		mysqlClient,
		p,
		WithLogger(logger),
	)
	assert.IsType(t, &experimentService{}, s)
}

func createExperimentService(c *gomock.Controller, s storage.Client) *experimentService {
	featureClientMock := featureclientmock.NewMockClient(c)
	fr := &featureproto.GetFeatureResponse{
		Feature: &featureproto.Feature{
			Id:         "fid",
			Version:    1,
			Variations: []*featureproto.Variation{},
		},
	}
	featureClientMock.EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(fr, nil).AnyTimes()
	fsr := &featureproto.GetFeaturesResponse{
		Features: []*featureproto.Feature{{
			Id:         "fid",
			Version:    1,
			Variations: []*featureproto.Variation{},
		}},
	}
	featureClientMock.EXPECT().GetFeatures(gomock.Any(), gomock.Any()).Return(fsr, nil).AnyTimes()
	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "ns0",
					Role:          accountproto.AccountV2_Role_Environment_EDITOR,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	mysqlClient := mysqlmock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	p.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	es := NewExperimentService(featureClientMock, accountClientMock, mysqlClient, p)
	return es.(*experimentService)
}

func createContextWithToken() context.Context {
	return createContextWithTokenAndMetadata(nil)
}

func createContextWithTokenAndMetadata(metadata map[string][]string) context.Context {
	token := &token.IDToken{
		Issuer:   "issuer",
		Subject:  "sub",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
	ctx := context.TODO()
	if metadata != nil {
		ctx = md.NewIncomingContext(ctx, metadata)
	}
	return context.WithValue(ctx, rpc.Key, token)
}

func createContextWithTokenRoleUnassigned() context.Context {
	token := &token.IDToken{
		Issuer:   "issuer",
		Subject:  "sub",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}
