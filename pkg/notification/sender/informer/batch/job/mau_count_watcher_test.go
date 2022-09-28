// Copyright 2022 The Bucketeer Authors.
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

package job

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	environmentclientmock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	ecclientmock "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client/mock"
	sendermock "github.com/bucketeer-io/bucketeer/pkg/notification/sender/mock"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

func TestCreateMAUNotification(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	errInternal := errors.New("internal error")
	patterns := map[string]struct {
		setup       func(*testing.T, *mauCountWatcher)
		expectedErr error
	}{
		"err project": {
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					nil, errInternal).Times(1)
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironments(
					gomock.Any(), gomock.Any()).Times(0)
				w.eventCounterClient.(*ecclientmock.MockClient).EXPECT().GetUserCountV2(
					gomock.Any(), gomock.Any()).Times(0)
				w.sender.(*sendermock.MockSender).EXPECT().Send(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedErr: errInternal,
		},
		"no projects": {
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{}, nil).Times(1)
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironments(
					gomock.Any(), gomock.Any()).Times(0)
				w.eventCounterClient.(*ecclientmock.MockClient).EXPECT().GetUserCountV2(
					gomock.Any(), gomock.Any()).Times(0)
				w.sender.(*sendermock.MockSender).EXPECT().Send(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedErr: nil,
		},
		"err environments": {
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{
						Projects: []*environmentproto.Project{{Id: "pj0"}},
						Cursor:   "cursor",
					}, nil).Times(1)
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironments(
					gomock.Any(), gomock.Any()).Return(
					nil, errInternal).Times(1)
				w.eventCounterClient.(*ecclientmock.MockClient).EXPECT().GetUserCountV2(
					gomock.Any(), gomock.Any()).Times(0)
				w.sender.(*sendermock.MockSender).EXPECT().Send(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedErr: errInternal,
		},
		"no environments": {
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{
						Projects: []*environmentproto.Project{{Id: "pj0"}},
						Cursor:   "cursor",
					}, nil).Times(1)
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironments(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsResponse{}, nil).Times(1)
				w.eventCounterClient.(*ecclientmock.MockClient).EXPECT().GetUserCountV2(
					gomock.Any(), gomock.Any()).Times(0)
				w.sender.(*sendermock.MockSender).EXPECT().Send(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedErr: nil,
		},
		"err counts": {
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{
						Projects: []*environmentproto.Project{{Id: "pj0"}},
						Cursor:   "cursor",
					}, nil).Times(1)
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironments(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsResponse{
						Environments: []*environmentproto.Environment{{Id: "eID", Namespace: "eNamespace"}},
						Cursor:       "",
					}, nil).Times(1)
				w.eventCounterClient.(*ecclientmock.MockClient).EXPECT().GetUserCountV2(
					gomock.Any(), gomock.Any()).Return(
					nil, errInternal).Times(1)
				w.sender.(*sendermock.MockSender).EXPECT().Send(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedErr: errInternal,
		},
		"err sender": {
			setup: func(t *testing.T, w *mauCountWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{
						Projects: []*environmentproto.Project{{Id: "pj0"}},
						Cursor:   "cursor",
					}, nil).Times(1)
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironments(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsResponse{
						Environments: []*environmentproto.Environment{{Id: "eID", Namespace: "eNamespace"}},
						Cursor:       "",
					}, nil).Times(1)
				w.eventCounterClient.(*ecclientmock.MockClient).EXPECT().GetUserCountV2(
					gomock.Any(), gomock.Any()).Return(
					&ecproto.GetUserCountV2Response{
						EventCount: 4,
						UserCount:  2,
					}, nil).Times(1)
				w.sender.(*sendermock.MockSender).EXPECT().Send(gomock.Any(), gomock.Any()).Return(errInternal).Times(1)
			},
			expectedErr: errInternal,
		},
		"success": {
			setup: func(t *testing.T, w *mauCountWatcher) {
				// list projects
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListProjects(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListProjectsResponse{
						Projects: []*environmentproto.Project{{Id: "pj0"}},
						Cursor:   "cursor",
					}, nil).Times(1)
				// list environments
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironments(
					gomock.Any(),
					&environmentproto.ListEnvironmentsRequest{
						PageSize:  listRequestSize,
						Cursor:    "",
						ProjectId: "pj0",
					},
				).Return(
					&environmentproto.ListEnvironmentsResponse{
						Environments: []*environmentproto.Environment{{Id: "eID", Namespace: "eNamespace"}},
						Cursor:       "",
					}, nil).Times(1)
				// get user count
				startAt, endAt := w.getMAUInterval(time.Now())
				w.eventCounterClient.(*ecclientmock.MockClient).EXPECT().GetUserCountV2(
					gomock.Any(),
					&ecproto.GetUserCountV2Request{
						EnvironmentNamespace: "eNamespace",
						StartAt:              startAt,
						EndAt:                endAt,
					},
				).Return(
					&ecproto.GetUserCountV2Response{
						EventCount: 4,
						UserCount:  2,
					}, nil).Times(1)
				// send notification
				w.sender.(*sendermock.MockSender).EXPECT().Send(
					gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			w := newMAUCountWatcherWithMock(t, mockController)
			if p.setup != nil {
				p.setup(t, w)
			}
			err := w.Run(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newMAUCountWatcherWithMock(t *testing.T, c *gomock.Controller) *mauCountWatcher {
	t.Helper()
	return &mauCountWatcher{
		environmentClient:  environmentclientmock.NewMockClient(c),
		eventCounterClient: ecclientmock.NewMockClient(c),
		sender:             sendermock.NewMockSender(c),
		logger:             zap.NewNop(),
		opts: &options{
			timeout: 5 * time.Minute,
		},
	}
}
