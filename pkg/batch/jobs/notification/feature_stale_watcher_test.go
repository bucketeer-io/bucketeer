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

package notification

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	environmentclientmock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	sendermock "github.com/bucketeer-io/bucketeer/pkg/notification/sender/mock"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestCreateNotification(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*testing.T, *featureStaleWatcher)
		expectedErr error
	}{
		{
			desc: "no featres",
			setup: func(t *testing.T, w *featureStaleWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{{Id: "ns0", Name: "ns0"}},
						Cursor:       "",
					}, nil)
				w.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any()).Return(
					&featureproto.ListFeaturesResponse{
						Features: []*featureproto.Feature{},
					}, nil)
			},
		},
		{
			desc: "no stale featres",
			setup: func(t *testing.T, w *featureStaleWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{{Id: "ns0", Name: "ns0"}},
						Cursor:       "",
					}, nil)
				w.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any()).Return(
					&featureproto.ListFeaturesResponse{
						Features: []*featureproto.Feature{{
							Id:      "fid",
							Name:    "fname",
							Enabled: true,
							LastUsedInfo: &featureproto.FeatureLastUsedInfo{
								LastUsedAt: time.Now().Unix(),
							},
						}},
					}, nil)
			},
		},
		{
			desc: "stale exists",
			setup: func(t *testing.T, w *featureStaleWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{{Id: "ns0", Name: "ns0"}},
						Cursor:       "",
					}, nil)
				w.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any()).Return(
					&featureproto.ListFeaturesResponse{
						Features: []*featureproto.Feature{{
							Id:      "fid",
							Name:    "fname",
							Enabled: true,
							LastUsedInfo: &featureproto.FeatureLastUsedInfo{
								LastUsedAt: time.Now().Unix() - 120*24*60*60,
							},
						}, {
							Id:      "fid1",
							Name:    "fname1",
							Enabled: true,
							LastUsedInfo: &featureproto.FeatureLastUsedInfo{
								LastUsedAt: time.Now().Unix() - 120*24*60*60,
							},
						}},
					}, nil)
				w.sender.(*sendermock.MockSender).EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			w := newFeatureWatcherWithMock(t, mockController)
			if p.setup != nil {
				p.setup(t, w)
			}
			err := w.Run(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newFeatureWatcherWithMock(t *testing.T, c *gomock.Controller) *featureStaleWatcher {
	t.Helper()
	return &featureStaleWatcher{
		environmentClient: environmentclientmock.NewMockClient(c),
		featureClient:     featureclientmock.NewMockClient(c),
		sender:            sendermock.NewMockSender(c),
		logger:            zap.NewNop(),
		opts: &jobs.Options{
			Timeout: 5 * time.Minute,
		},
	}
}
