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

package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestChangeBulkUploadSegmentUsersStatus(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	patterns := []struct {
		desc        string
		setup       func(*segmentCommandHandler)
		cmd         *featureproto.ChangeBulkUploadSegmentUsersStatusCommand
		expectedErr error
	}{
		{
			desc: "succeeded included",
			setup: func(s *segmentCommandHandler) {
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			cmd: &featureproto.ChangeBulkUploadSegmentUsersStatusCommand{
				Status: featureproto.Segment_SUCEEDED,
				State:  featureproto.SegmentUser_INCLUDED,
				Count:  1,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			segment, err := domain.NewSegment("test-name", "test-description")
			assert.NoError(t, err)
			handler := newMockSegmentCommandHandler(t, mockController, segment)
			p.setup(handler)
			err = handler.Handle(ctx, p.cmd)
			assert.Equal(t, p.expectedErr, err, p.desc)
			assert.Equal(t, segment.Status, p.cmd.Status)
			switch p.cmd.State {
			case featureproto.SegmentUser_INCLUDED:
				assert.Equal(t, segment.IncludedUserCount, p.cmd.Count)
			default:
				t.Fatal("unknown segment user state")
			}
		})
	}
}

func newMockSegmentCommandHandler(t *testing.T, mockController *gomock.Controller, segment *domain.Segment) *segmentCommandHandler {
	t.Helper()
	return &segmentCommandHandler{
		&eventproto.Editor{
			Email: "email",
		},
		segment,
		publishermock.NewMockPublisher(mockController),
		"bucketeer-environment-space",
	}
}
