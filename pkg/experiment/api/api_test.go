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
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	md "google.golang.org/grpc/metadata"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	autoopsclientmock "github.com/bucketeer-io/bucketeer/pkg/autoops/client/mock"
	storagemock "github.com/bucketeer-io/bucketeer/pkg/experiment/storage/v2/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
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
	autoOpsClientMock := autoopsclientmock.NewMockClient(mockController)
	mysqlClient := mysqlmock.NewMockClient(mockController)
	p := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewExperimentService(
		featureClientMock,
		accountClientMock,
		autoOpsClientMock,
		mysqlClient,
		p,
		WithLogger(logger),
	)
	assert.IsType(t, &experimentService{}, s)
}

func createExperimentService(c *gomock.Controller, specifiedEnvironmentId *string, specifiedOrgRole *accountproto.AccountV2_Role_Organization, specifiedEnvRole *accountproto.AccountV2_Role_Environment) *experimentService {
	var or accountproto.AccountV2_Role_Organization
	var er accountproto.AccountV2_Role_Environment
	var envId string
	if specifiedOrgRole != nil {
		or = *specifiedOrgRole
	} else {
		or = accountproto.AccountV2_Role_Organization_ADMIN
	}
	if specifiedEnvRole != nil {
		er = *specifiedEnvRole
	} else {
		er = accountproto.AccountV2_Role_Environment_EDITOR
	}
	if specifiedEnvironmentId != nil {
		envId = *specifiedEnvironmentId
	} else {
		envId = "ns0"
	}

	featureClientMock := featureclientmock.NewMockClient(c)
	fr := &featureproto.GetFeatureResponse{
		Feature: &featureproto.Feature{
			Id:      "fid",
			Version: 1,
			Variations: []*featureproto.Variation{
				{
					Id: "variation-a-id",
				},
				{
					Id: "variation-b-id",
				},
			},
		},
	}
	featureClientMock.EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(fr, nil).AnyTimes()
	fsr := &featureproto.GetFeaturesResponse{
		Features: []*featureproto.Feature{{
			Id:      "fid",
			Version: 1,
			Variations: []*featureproto.Variation{
				{
					Id: "variation-a-id",
				},
				{
					Id: "variation-b-id",
				},
			},
		}},
	}
	featureClientMock.EXPECT().GetFeatures(gomock.Any(), gomock.Any()).Return(fsr, nil).AnyTimes()
	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: or,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: envId,
					Role:          er,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	autoOpsClientMock := autoopsclientmock.NewMockClient(c)
	mysqlClient := mysqlmock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	p.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &experimentService{
		featureClient:     featureClientMock,
		accountClient:     accountClientMock,
		autoOpsClient:     autoOpsClientMock,
		mysqlClient:       mysqlClient,
		experimentStorage: storagemock.NewMockExperimentStorage(c),
		goalStorage:       storagemock.NewMockGoalStorage(c),
		publisher:         p,
		logger:            zap.NewNop().Named("api"),
	}
}

func createContextWithToken() context.Context {
	return createContextWithTokenAndMetadata(nil)
}

func createContextWithTokenAndMetadata(metadata map[string][]string) context.Context {
	token := &token.AccessToken{
		Issuer:   "issuer",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
	ctx := context.TODO()
	if metadata != nil {
		ctx = md.NewIncomingContext(ctx, metadata)
	}
	return context.WithValue(ctx, rpc.AccessTokenKey, token)
}

func createContextWithTokenRoleUnassigned() context.Context {
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

// convert to pointer
func toPtr[T any](value T) *T {
	return &value
}
