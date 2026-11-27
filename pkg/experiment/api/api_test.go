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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	md "google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/v2/pkg/authz"
	autoopsclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	bkterr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	dbmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// errPublish is returned by the publisher mock to simulate a publish failure
// after the transaction has committed.
var errPublish = bkterr.NewErrorInternal(bkterr.ExperimentPackageName, "publish error")

// errStorage is returned by storage mocks to simulate a failed write inside
// the transaction.
var errStorage = bkterr.NewErrorInternal(bkterr.ExperimentPackageName, "storage error")

// errCommit simulates a transaction that fails to commit after the callback succeeded.
var errCommit = bkterr.NewErrorInternal(bkterr.ExperimentPackageName, "commit error")

func TestNewExperimentService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	featureClientMock := featureclientmock.NewMockClient(mockController)
	autoOpsClientMock := autoopsclientmock.NewMockClient(mockController)
	dbClientMock := dbmock.NewMockClient(mockController)
	experimentStorageMock := storagemock.NewMockExperimentStorage(mockController)
	goalStorageMock := storagemock.NewMockGoalStorage(mockController)
	p := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewExperimentService(
		featureClientMock,
		autoOpsClientMock,
		dbClientMock,
		experimentStorageMock,
		goalStorageMock,
		p,
		WithLogger(logger),
	)
	assert.IsType(t, &experimentService{}, s)
}

func createExperimentService(c *gomock.Controller) *experimentService {
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
	autoOpsClientMock := autoopsclientmock.NewMockClient(c)
	dbClientMock := dbmock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	return &experimentService{
		featureClient:     featureClientMock,
		autoOpsClient:     autoOpsClientMock,
		dbClient:          dbClientMock,
		experimentStorage: storagemock.NewMockExperimentStorage(c),
		goalStorage:       storagemock.NewMockGoalStorage(c),
		publisher:         p,
		logger:            zap.NewNop().Named("api"),
	}
}

// createContextWithEditor is a context after the authz interceptor allowed the request.
func createContextWithEditor() context.Context {
	return createContextWithEditorAndMetadata(nil)
}

func createContextWithEditorAndMetadata(metadata map[string][]string) context.Context {
	ctx := context.TODO()
	if metadata != nil {
		ctx = md.NewIncomingContext(ctx, metadata)
	}
	return authz.WithEditor(ctx, &eventproto.Editor{Email: "email"})
}
