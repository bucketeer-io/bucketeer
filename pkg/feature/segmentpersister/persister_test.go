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

package segmentpersister

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	cachemock "github.com/bucketeer-io/bucketeer/pkg/cache/mock"
	cachev3mock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	metricsmock "github.com/bucketeer-io/bucketeer/pkg/metrics/mock"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	pullermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/mock"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	serviceevent "github.com/bucketeer-io/bucketeer/proto/event/service"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestNewPersister(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	puller := pullermock.NewMockPuller(mockController)
	publisher := publishermock.NewMockPublisher(mockController)
	mysqlClient := mysqlmock.NewMockClient(mockController)
	redis := cachemock.NewMockMultiGetCache(mockController)
	registerer := metricsmock.NewMockRegisterer(mockController)
	registerer.EXPECT().MustRegister(gomock.Any()).Return()
	p := NewPersister(
		puller,
		publisher,
		mysqlClient,
		redis,
		WithMaxMPS(100),
		WithNumWorkers(1),
		WithFlushSize(1),
		WithFlushInterval(time.Second),
		WithMetrics(registerer),
		WithLogger(zap.NewNop()),
	)
	assert.IsType(t, &Persister{}, p)
}

func TestHandleEventMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc          string
		setup         func(*Persister)
		event         *serviceevent.BulkSegmentUsersReceivedEvent
		segment       *domain.Segment
		expectedCount int64
		expectedErr   error
	}{
		{
			desc: "err: ErrSegmentNotFound",
			setup: func(p *Persister) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			event: &serviceevent.BulkSegmentUsersReceivedEvent{
				SegmentId:            "sid",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: v2fs.ErrSegmentNotFound,
		},
		{
			desc: "err: errExceededMaxUserIDLength",
			setup: func(p *Persister) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			event: &serviceevent.BulkSegmentUsersReceivedEvent{
				SegmentId:            "sid",
				EnvironmentNamespace: "ns0",
				Data:                 []byte(strings.Repeat("a", maxUserIDLength+1)),
				State:                featureproto.SegmentUser_INCLUDED,
			},
			expectedErr: errExceededMaxUserIDLength,
		},
		{
			desc: "success",
			setup: func(p *Persister) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil).Times(2)
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil).Times(2)
			},
			event: &serviceevent.BulkSegmentUsersReceivedEvent{
				SegmentId:            "sid",
				EnvironmentNamespace: "ns0",
				Data:                 []byte("user1\nuser2\r\nuser2\n"),
				State:                featureproto.SegmentUser_INCLUDED,
				Editor: &eventproto.Editor{
					Email: "email",
				},
			},
			expectedErr: nil,
		},
	}
	for _, pat := range patterns {
		t.Run(pat.desc, func(t *testing.T) {
			persister := newPersister(t, mockController)
			if pat.setup != nil {
				pat.setup(persister)
			}
			err := persister.handleEvent(context.Background(), pat.event)
			assert.Equal(t, pat.expectedErr, err)
		})
	}
}

func newPersister(t *testing.T, mockController *gomock.Controller) *Persister {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	logger := zap.NewNop()
	return &Persister{
		puller:            pullermock.NewMockRateLimitedPuller(mockController),
		domainPublisher:   publishermock.NewMockPublisher(mockController),
		mysqlClient:       mysqlmock.NewMockClient(mockController),
		segmentUsersCache: cachev3mock.NewMockSegmentUsersCache(mockController),
		logger:            logger.Named("persister"),
		ctx:               ctx,
		cancel:            cancel,
		doneCh:            make(chan struct{}),
	}
}
