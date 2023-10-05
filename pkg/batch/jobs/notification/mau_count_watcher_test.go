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

package notification

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	environmentclientmock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	ecclientmock "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client/mock"
	sendermock "github.com/bucketeer-io/bucketeer/pkg/notification/sender/mock"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

var (
	jpLocation = time.FixedZone("Asia/Tokyo", 9*60*60)
)

func TestCreateMAUNotification(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	errInternal := errors.New("internal error")
	patterns := []struct {
		desc        string
		setup       func(*testing.T, *mauCountWatcher)
		expectedErr error
	}{
		{
			desc: "err project",
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					nil, errInternal)
			},
			expectedErr: errInternal,
		},
		{
			desc: "no projects",
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{}, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "err environments",
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{
						Projects: []*environmentproto.Project{{Id: "pj0"}},
						Cursor:   "cursor",
					}, nil)
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any()).Return(
					nil, errInternal)
			},
			expectedErr: errInternal,
		},
		{
			desc: "no environments",
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{
						Projects: []*environmentproto.Project{{Id: "pj0"}},
						Cursor:   "cursor",
					}, nil)
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsV2Response{}, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "err counts",
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{
						Projects: []*environmentproto.Project{{Id: "pj0"}},
						Cursor:   "cursor",
					}, nil)
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{{Id: "eID", Name: "eID"}},
						Cursor:       "",
					}, nil)
				w.eventCounterClient.(*ecclientmock.MockClient).EXPECT().GetMAUCount(
					gomock.Any(), gomock.Any()).Return(
					nil, errInternal)
				w.sender.(*sendermock.MockSender).EXPECT().Send(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedErr: errInternal,
		},
		{
			desc: "err sender",
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{
						Projects: []*environmentproto.Project{{Id: "pj0"}},
						Cursor:   "cursor",
					}, nil)
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{{Id: "eID", Name: "eID"}},
						Cursor:       "",
					}, nil)
				w.eventCounterClient.(*ecclientmock.MockClient).EXPECT().GetMAUCount(
					gomock.Any(), gomock.Any()).Return(
					&ecproto.GetMAUCountResponse{
						EventCount: 4,
						UserCount:  2,
					}, nil)
				w.sender.(*sendermock.MockSender).EXPECT().Send(gomock.Any(), gomock.Any()).Return(errInternal)
			},
			expectedErr: errInternal,
		},
		{
			desc: "success",
			setup: func(t *testing.T, w *mauCountWatcher) {
				// list projects
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(),
					&environmentproto.ListProjectsRequest{
						PageSize: listRequestSize,
						Cursor:   "",
					},
				).Return(
					&environmentproto.ListProjectsResponse{
						Projects: []*environmentproto.Project{{Id: "pj0"}},
						Cursor:   "cursor",
					}, nil)
				// list environments
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					&environmentproto.ListEnvironmentsV2Request{
						PageSize:  listRequestSize,
						Cursor:    "",
						ProjectId: "pj0",
						Archived:  &wrappers.BoolValue{Value: false},
					},
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{{Id: "eID", Name: "eID"}},
						Cursor:       "",
					}, nil)
				// get user count
				year, month := w.getLastYearMonth(time.Now().In(w.location))
				w.eventCounterClient.(*ecclientmock.MockClient).EXPECT().GetMAUCount(
					gomock.Any(),
					&ecproto.GetMAUCountRequest{
						EnvironmentNamespace: "eID",
						YearMonth:            w.newYearMonth(year, month),
					},
				).Return(
					&ecproto.GetMAUCountResponse{
						EventCount: 4,
						UserCount:  2,
					}, nil)
				// send notification
				w.sender.(*sendermock.MockSender).EXPECT().Send(
					gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			w := newMAUCountWatcherWithMock(t, mockController)
			if p.setup != nil {
				p.setup(t, w)
			}
			err := w.Run(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetYearLastMonth(t *testing.T) {
	t.Parallel()
	watcher := &mauCountWatcher{}
	unixTime := time.Unix(1672531200, 0) // 2023/01/01 00:00:00 UTC
	unixTime.In(jpLocation)
	year, month := watcher.getLastYearMonth(unixTime)
	assert.Equal(t, int32(2022), year)
	assert.Equal(t, int32(12), month)
}

func TestNewYearMonth(t *testing.T) {
	t.Parallel()
	watcher := &mauCountWatcher{}
	uniTime := time.Unix(1675209600, 0) // 2023/02/01 00:00:00 UTC
	uniTime.In(jpLocation)
	year, month := watcher.getLastYearMonth(uniTime)
	yearMonth := watcher.newYearMonth(year, month)
	assert.Equal(t, "202301", yearMonth)
}

func newMAUCountWatcherWithMock(t *testing.T, c *gomock.Controller) *mauCountWatcher {
	t.Helper()
	return &mauCountWatcher{
		environmentClient:  environmentclientmock.NewMockClient(c),
		eventCounterClient: ecclientmock.NewMockClient(c),
		sender:             sendermock.NewMockSender(c),
		location:           jpLocation,
		logger:             zap.NewNop(),
		opts: &jobs.Options{
			Timeout: 5 * time.Minute,
		},
	}
}
