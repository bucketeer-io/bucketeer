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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	domain "github.com/bucketeer-io/bucketeer/pkg/push/domain"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/push"
)

func TestCreate(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		expected error
	}{
		{
			expected: nil,
		},
	}
	for _, p := range patterns {
		pm := publishermock.NewMockPublisher(mockController)
		pd := newPush(t)
		ch := newPushCommandHandler(t, pm, pd)
		if p.expected == nil {
			pm.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		}
		cmd := &proto.CreatePushCommand{Name: "name-1", FcmApiKey: "key-0", Tags: []string{"tag-0", "tag-1"}}
		err := ch.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
}

func TestDelete(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		expected error
	}{
		{
			expected: nil,
		},
	}
	for _, p := range patterns {
		pm := publishermock.NewMockPublisher(mockController)
		pd := newPush(t)
		ch := newPushCommandHandler(t, pm, pd)
		if p.expected == nil {
			pm.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		}
		cmd := &proto.DeletePushCommand{}
		err := ch.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
}

func TestAddTags(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		expected error
	}{
		{
			expected: nil,
		},
	}
	for _, p := range patterns {
		pm := publishermock.NewMockPublisher(mockController)
		pd := newPush(t)
		ch := newPushCommandHandler(t, pm, pd)
		if p.expected == nil {
			pm.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		}
		cmd := &proto.AddPushTagsCommand{Tags: []string{"tag-2"}}
		err := ch.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
}

func TestDeleteTags(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		expected error
	}{
		{
			expected: nil,
		},
	}
	for _, p := range patterns {
		pm := publishermock.NewMockPublisher(mockController)
		pd := newPush(t)
		ch := newPushCommandHandler(t, pm, pd)
		if p.expected == nil {
			pm.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		}
		cmd := &proto.DeletePushTagsCommand{Tags: []string{"tag-0"}}
		err := ch.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
}

func TestRename(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []*struct {
		expected error
	}{
		{
			expected: nil,
		},
	}
	for _, p := range patterns {
		pm := publishermock.NewMockPublisher(mockController)
		pd := newPush(t)
		ch := newPushCommandHandler(t, pm, pd)
		if p.expected == nil {
			pm.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		}
		cmd := &proto.RenamePushCommand{Name: "name-2"}
		err := ch.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
}

func newPush(t *testing.T) *domain.Push {
	d, err := domain.NewPush("name-1", "key-0", []string{"tag-0", "tag-1"})
	require.NoError(t, err)
	return d
}

func newPushCommandHandler(t *testing.T, publisher publisher.Publisher, push *domain.Push) Handler {
	t.Helper()
	return NewPushCommandHandler(
		&eventproto.Editor{
			Email: "email",
		},
		push,
		publisher,
		"ns0",
	)
}
