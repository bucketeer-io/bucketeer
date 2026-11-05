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

package processor

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	fsmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	dbmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	domainproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	errSegmentStorage = errors.New("segment storage error")
	errSegmentCommit  = errors.New("segment commit error")
	errSegmentPublish = errors.New("segment publish error")
)

func newSegmentUserPersisterForTest(c *gomock.Controller) *segmentUserPersister {
	return &segmentUserPersister{
		domainPublisher: publishermock.NewMockPublisher(c),
		dbClient:        dbmock.NewMockClient(c),
		segmentStorage:  fsmock.NewMockSegmentStorage(c),
		logger:          zap.NewNop(),
	}
}

func TestUpdateSegmentStatus(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	segment := func() *domain.Segment {
		return &domain.Segment{
			Segment: &featureproto.Segment{
				Id:     "sid",
				Name:   "segment",
				Status: featureproto.Segment_UPLOADING,
			},
		}
	}
	runCallback := func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	}

	patterns := []struct {
		desc        string
		setup       func(*segmentUserPersister)
		expectedErr error
	}{
		{
			desc: "err: get segment fails, event is not published",
			setup: func(p *segmentUserPersister) {
				p.dbClient.(*dbmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(runCallback)
				p.segmentStorage.(*fsmock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), "sid", "ns0",
				).Return(nil, nil, errSegmentStorage)
				p.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Times(0)
			},
			expectedErr: errSegmentStorage,
		},
		{
			desc: "err: update segment fails, event is not published",
			setup: func(p *segmentUserPersister) {
				p.dbClient.(*dbmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(runCallback)
				p.segmentStorage.(*fsmock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), "sid", "ns0",
				).Return(segment(), nil, nil)
				p.segmentStorage.(*fsmock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), "ns0",
				).Return(errSegmentStorage)
				p.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Times(0)
			},
			expectedErr: errSegmentStorage,
		},
		{
			desc: "err: commit fails, event is not published",
			setup: func(p *segmentUserPersister) {
				p.dbClient.(*dbmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					if err := fn(ctx); err != nil {
						return err
					}
					return errSegmentCommit
				})
				p.segmentStorage.(*fsmock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), "sid", "ns0",
				).Return(segment(), nil, nil)
				p.segmentStorage.(*fsmock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), "ns0",
				).Return(nil)
				p.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Times(0)
			},
			expectedErr: errSegmentCommit,
		},
		{
			desc: "err: publish fails after commit",
			setup: func(p *segmentUserPersister) {
				p.dbClient.(*dbmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(runCallback)
				p.segmentStorage.(*fsmock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), "sid", "ns0",
				).Return(segment(), nil, nil)
				p.segmentStorage.(*fsmock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), "ns0",
				).Return(nil)
				p.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(errSegmentPublish).Times(1)
			},
			expectedErr: errSegmentPublish,
		},
		{
			desc: "success: event is published after commit",
			setup: func(p *segmentUserPersister) {
				p.dbClient.(*dbmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(runCallback)
				p.segmentStorage.(*fsmock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), "sid", "ns0",
				).Return(segment(), nil, nil)
				p.segmentStorage.(*fsmock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), "ns0",
				).DoAndReturn(func(_ context.Context, s *domain.Segment, _ string) error {
					assert.Equal(t, featureproto.Segment_SUCEEDED, s.Status)
					assert.Equal(t, int64(3), s.IncludedUserCount)
					return nil
				})
				p.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(_ context.Context, e *domainproto.Event) error {
					assert.Equal(t, domainproto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED, e.Type)
					assert.Equal(t, "sid", e.EntityId)
					assert.Equal(t, "ns0", e.EnvironmentId)
					return nil
				}).Times(1)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newSegmentUserPersisterForTest(mockController)
			p.setup(persister)
			err := persister.updateSegmentStatus(
				context.Background(),
				&domainproto.Editor{Email: "test@example.com"},
				"ns0",
				"sid",
				3,
				featureproto.SegmentUser_INCLUDED,
				featureproto.Segment_SUCEEDED,
			)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}
