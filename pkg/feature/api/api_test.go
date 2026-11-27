// Copyright 2026 The Bucketeer Authors.
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

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"github.com/bucketeer-io/bucketeer/v2/pkg/authz"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	tagstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage/mock"

	aoclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	autoopsclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	envclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	environmentId = "test"
	tag           = "tag"
	userID        = "user-id"
)

var (
	defaultOptions = options{
		logger: zap.NewNop(),
	}
)

type dummyWebhookCryptoUtil struct{}

func (u *dummyWebhookCryptoUtil) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	return []byte(data), nil
}

func (u *dummyWebhookCryptoUtil) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	return []byte(data), nil
}

// createContextWithEditor is a context after the authz interceptor allowed the request.
func createContextWithEditor() context.Context {
	return authz.WithEditor(context.TODO(), &eventproto.Editor{Email: "email"})
}

// FIXME: Deprecated. Do not use for a new test. Replace this with createFeatureServiceNew.
func createFeatureService(c *gomock.Controller) *FeatureService {
	p := publishermock.NewMockPublisher(c)
	p.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	p.EXPECT().PublishMulti(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	e := experimentclientmock.NewMockClient(c)
	e.EXPECT().ListExperiments(gomock.Any(), gomock.Any()).Return(&experimentproto.ListExperimentsResponse{}, nil).AnyTimes()
	at := autoopsclientmock.NewMockClient(c)
	env := envclientmock.NewMockClient(c)
	at.EXPECT().ListProgressiveRollouts(gomock.Any(), gomock.Any()).Return(&autoopsproto.ListProgressiveRolloutsResponse{}, nil).AnyTimes()
	return &FeatureService{
		mock.NewMockFlagTriggerStorage(c),
		mock.NewMockFeatureStorage(c),
		mock.NewMockSegmentStorage(c),
		mock.NewMockSegmentUserStorage(c),
		mock.NewMockScheduledFlagChangeStorage(c),
		nil,
		tagstoragemock.NewMockTagStorage(c),
		databasemock.NewMockClient(c),
		e,
		cachev3mock.NewMockFeaturesCache(c),
		nil,
		at,
		env,
		cachev3mock.NewMockSegmentUsersCache(c),
		cachev3mock.NewMockUserAttributesCache(c),
		p,
		p,
		singleflight.Group{},
		"http://localhost",
		&defaultOptions,
		defaultOptions.logger,
	}
}

func createFeatureServiceNew(c *gomock.Controller) *FeatureService {
	segmentUsersPublisher := publishermock.NewMockPublisher(c)
	domainPublisher := publishermock.NewMockPublisher(c)
	return &FeatureService{
		flagTriggerStorage:         mock.NewMockFlagTriggerStorage(c),
		featureStorage:             mock.NewMockFeatureStorage(c),
		segmentStorage:             mock.NewMockSegmentStorage(c),
		segmentUserStorage:         mock.NewMockSegmentUserStorage(c),
		scheduledFlagChangeStorage: mock.NewMockScheduledFlagChangeStorage(c),
		dbClient:                   databasemock.NewMockClient(c),
		tagStorage:                 tagstoragemock.NewMockTagStorage(c),
		autoOpsClient:              aoclientmock.NewMockClient(c),
		experimentClient:           experimentclientmock.NewMockClient(c),
		environmentClient:          envclientmock.NewMockClient(c),
		featuresCache:              cachev3mock.NewMockFeaturesCache(c),
		segmentUsersPublisher:      segmentUsersPublisher,
		userAttributesCache:        cachev3mock.NewMockUserAttributesCache(c),
		domainPublisher:            domainPublisher,
		triggerURL:                 "http://localhost",
		opts:                       &defaultOptions,
		logger:                     defaultOptions.logger,
	}
}

func createFeatureServiceWithSegmentUsersCache(c *gomock.Controller) *FeatureService {
	segmentUsersPublisher := publishermock.NewMockPublisher(c)
	domainPublisher := publishermock.NewMockPublisher(c)
	return &FeatureService{
		flagTriggerStorage:         mock.NewMockFlagTriggerStorage(c),
		featureStorage:             mock.NewMockFeatureStorage(c),
		scheduledFlagChangeStorage: mock.NewMockScheduledFlagChangeStorage(c),
		segmentUsersCache:          cachev3mock.NewMockSegmentUsersCache(c),
		segmentStorage:             mock.NewMockSegmentStorage(c),
		segmentUserStorage:         mock.NewMockSegmentUserStorage(c),
		dbClient:                   databasemock.NewMockClient(c),
		autoOpsClient:              aoclientmock.NewMockClient(c),
		experimentClient:           experimentclientmock.NewMockClient(c),
		environmentClient:          envclientmock.NewMockClient(c),
		featuresCache:              cachev3mock.NewMockFeaturesCache(c),
		segmentUsersPublisher:      segmentUsersPublisher,
		userAttributesCache:        cachev3mock.NewMockUserAttributesCache(c),
		domainPublisher:            domainPublisher,
		triggerURL:                 "http://localhost",
		opts:                       &defaultOptions,
		logger:                     defaultOptions.logger,
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
