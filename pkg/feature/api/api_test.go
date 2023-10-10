// Copyright 2023 The Bucketeer Authors.
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
	"fmt"
	"time"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	autoopsclientmock "github.com/bucketeer-io/bucketeer/pkg/autoops/client/mock"
	cachev3mock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	environmentNamespace = "test"
	tag                  = "tag"
	userID               = "user-id"
)

var (
	defaultOptions = options{
		logger: zap.NewNop(),
	}
	evaluation = &featureproto.Evaluation{
		FeatureId:      "feature-id",
		FeatureVersion: 1,
		UserId:         "user-id",
		VariationId:    "variation-id",
		VariationValue: "variation-value",
	}
)

func createContextWithToken() context.Context {
	token := &token.IDToken{
		Issuer:    "issuer",
		Subject:   "sub",
		Audience:  "audience",
		Expiry:    time.Now().AddDate(100, 0, 0),
		IssuedAt:  time.Now(),
		Email:     "email",
		AdminRole: accountproto.Account_OWNER,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

func createContextWithTokenRoleUnassigned() context.Context {
	token := &token.IDToken{
		Issuer:    "issuer",
		Subject:   "sub",
		Audience:  "audience",
		Expiry:    time.Now().AddDate(100, 0, 0),
		IssuedAt:  time.Now(),
		Email:     "email",
		AdminRole: accountproto.Account_UNASSIGNED,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

// FIXME: Deprecated. Do not use for a new test. Replace this with createFeatureServiceNew.
func createFeatureService(c *gomock.Controller) *FeatureService {
	p := publishermock.NewMockPublisher(c)
	p.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	p.EXPECT().PublishMulti(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	a := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountResponse{
		Account: &accountproto.Account{
			Email: "email",
			Role:  accountproto.Account_VIEWER,
		},
	}
	a.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	e := experimentclientmock.NewMockClient(c)
	e.EXPECT().ListExperiments(gomock.Any(), gomock.Any()).Return(&experimentproto.ListExperimentsResponse{}, nil).AnyTimes()
	at := autoopsclientmock.NewMockClient(c)
	at.EXPECT().ListProgressiveRollouts(gomock.Any(), gomock.Any()).Return(&autoopsproto.ListProgressiveRolloutsResponse{}, nil).AnyTimes()
	return &FeatureService{
		mysqlmock.NewMockClient(c),
		a,
		e,
		cachev3mock.NewMockFeaturesCache(c),
		at,
		cachev3mock.NewMockSegmentUsersCache(c),
		p,
		p,
		singleflight.Group{},
		&defaultOptions,
		defaultOptions.logger,
	}
}

func createFeatureServiceNew(c *gomock.Controller) *FeatureService {
	segmentUsersPublisher := publishermock.NewMockPublisher(c)
	domainPublisher := publishermock.NewMockPublisher(c)
	a := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountResponse{
		Account: &accountproto.Account{
			Email: "email",
			Role:  accountproto.Account_VIEWER,
		},
	}
	a.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	return &FeatureService{
		mysqlClient:           mysqlmock.NewMockClient(c),
		accountClient:         a,
		experimentClient:      experimentclientmock.NewMockClient(c),
		featuresCache:         cachev3mock.NewMockFeaturesCache(c),
		segmentUsersPublisher: segmentUsersPublisher,
		domainPublisher:       domainPublisher,
		opts:                  &defaultOptions,
		logger:                defaultOptions.logger,
	}
}

func createFeatureVariations() []*featureproto.Variation {
	return []*featureproto.Variation{
		{
			Value:       "variation_value_1",
			Name:        "variation_name_1",
			Description: "variation_description_1",
		},
		{
			Value:       "variation_value_2",
			Name:        "variation_name_2",
			Description: "variation_description_2",
		},
	}
}

func createFeatureTags() []string {
	return []string{"feature-tag-1", "feature-tag-2", "feature-tag-3"}
}

func contains(needle string, haystack []string) bool {
	for i := range haystack {
		if haystack[i] == needle {
			return true
		}
	}
	return false
}

type msgLengthMatcher struct{ length int }

func newMsgLengthMatcher(length int) gomock.Matcher {
	return &msgLengthMatcher{length: length}
}

func (m *msgLengthMatcher) Matches(x interface{}) bool {
	return len(x.([]publisher.Message)) == m.length
}

func (m *msgLengthMatcher) String() string {
	return fmt.Sprintf("length: %d", m.length)
}
