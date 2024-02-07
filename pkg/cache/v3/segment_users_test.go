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

package v3

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachemock "github.com/bucketeer-io/bucketeer/pkg/cache/mock"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	segmentID = "segment-id"
)

func TestGetSegmentUser(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	segmentUsers := createSegmentUsersCache(t)
	dataSegmentUsers := marshalMessage(t, segmentUsers)
	key := cache.MakeKey(segmentUsersKind, segmentID, environmentNamespace)

	patterns := []struct {
		desc        string
		setup       func(*segmentUsersCache)
		expectedErr error
	}{
		{
			desc: "error_get_not_found",
			setup: func(sc *segmentUsersCache) {
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().Get(key).Return(nil, cache.ErrNotFound)
			},
			expectedErr: cache.ErrNotFound,
		},
		{
			desc: "error_invalid_type",
			setup: func(sc *segmentUsersCache) {
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().Get(key).Return("test", nil)
			},
			expectedErr: cache.ErrInvalidType,
		},
		{
			desc: "success",
			setup: func(sc *segmentUsersCache) {
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().Get(key).Return(dataSegmentUsers, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sc := newSegmentUsersCache(t, mockController)
			p.setup(sc)
			cache, err := sc.Get(segmentID, environmentNamespace)
			if err == nil {
				assert.Equal(t, segmentUsers.SegmentId, cache.SegmentId)
				assert.Equal(t, segmentUsers.Users[0].Id, cache.Users[0].Id)
				assert.Equal(t, segmentUsers.Users[0].SegmentId, cache.Users[0].SegmentId)
				assert.Equal(t, segmentUsers.Users[0].UserId, cache.Users[0].UserId)
				assert.Equal(t, segmentUsers.Users[0].State, cache.Users[0].State)
				assert.Equal(t, segmentUsers.Users[0].Deleted, cache.Users[0].Deleted)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAllSegmentUser(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	segmentUsers := createSegmentUsersCache(t)
	dataSegmentUsers := marshalMessage(t, segmentUsers)
	keys := []string{
		fmt.Sprintf("%s:%s:segment-id-1", environmentNamespace, segmentUsersKind),
		fmt.Sprintf("%s:%s:segment-id-2", environmentNamespace, segmentUsersKind),
	}

	keyPrefix := cache.MakeKeyPrefix(segmentUsersKind, environmentNamespace)
	key := keyPrefix + "*"
	var cursor uint64

	patterns := []struct {
		desc        string
		setup       func(*segmentUsersCache)
		expectedErr error
	}{
		{
			desc: "error_scan_not_found",
			setup: func(sc *segmentUsersCache) {
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().Scan(cursor, key, segmentUsersMaxSize).Return(
					cursor, nil, cache.ErrNotFound)
			},
			expectedErr: cache.ErrNotFound,
		},
		{
			desc: "error_get_multi_not_found",
			setup: func(sc *segmentUsersCache) {
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().Scan(cursor, key, segmentUsersMaxSize).Return(
					cursor, keys, nil)
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().GetMulti(keys).Return(nil, cache.ErrNotFound)
			},
			expectedErr: cache.ErrNotFound,
		},
		{
			desc: "error_invalid_type",
			setup: func(sc *segmentUsersCache) {
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().Scan(cursor, key, segmentUsersMaxSize).Return(
					cursor, keys, nil)
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().GetMulti(keys).Return([]interface{}{"test"}, nil)
			},
			expectedErr: cache.ErrInvalidType,
		},
		{
			desc: "success",
			setup: func(sc *segmentUsersCache) {
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().Scan(cursor, key, segmentUsersMaxSize).Return(
					cursor, keys, nil)
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().GetMulti(keys).Return([]interface{}{dataSegmentUsers}, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sc := newSegmentUsersCache(t, mockController)
			p.setup(sc)
			allUsers, err := sc.GetAll(environmentNamespace)
			if err == nil {
				users := allUsers[0]
				assert.Equal(t, segmentUsers.SegmentId, users.SegmentId)
				for i := 0; i < len(segmentUsers.Users); i++ {
					assert.Equal(t, segmentUsers.Users[i].Id, users.Users[i].Id)
					assert.Equal(t, segmentUsers.Users[i].SegmentId, users.Users[i].SegmentId)
					assert.Equal(t, segmentUsers.Users[i].UserId, users.Users[i].UserId)
					assert.Equal(t, segmentUsers.Users[i].State, users.Users[i].State)
					assert.Equal(t, segmentUsers.Users[i].Deleted, users.Users[i].Deleted)
				}
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestPutSegmentUser(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	segmentUsers := createSegmentUsersCache(t)
	dataSegmentUsers := marshalMessage(t, segmentUsers)
	key := cache.MakeKey(segmentUsersKind, segmentID, environmentNamespace)

	patterns := []struct {
		desc        string
		setup       func(*segmentUsersCache)
		input       *featureproto.SegmentUsers
		expectedErr error
	}{
		{
			desc:        "error_proto_message_nil",
			setup:       nil,
			input:       nil,
			expectedErr: proto.ErrNil,
		},
		{
			desc: "success",
			setup: func(sc *segmentUsersCache) {
				sc.cache.(*cachemock.MockMultiGetCache).EXPECT().Put(key, dataSegmentUsers, segmentUsersTTL).Return(nil)
			},
			input:       segmentUsers,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sc := newSegmentUsersCache(t, mockController)
			if p.setup != nil {
				p.setup(sc)
			}
			err := sc.Put(p.input, environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func createSegmentUsersCache(t *testing.T) *featureproto.SegmentUsers {
	t.Helper()
	u := []*featureproto.SegmentUser{}
	for i := 0; i < 5; i++ {
		user := &featureproto.SegmentUser{
			Id:        fmt.Sprintf("segment-user-id-%d", i),
			SegmentId: segmentID,
			UserId:    fmt.Sprintf("user-id-%d", i),
			State:     featureproto.SegmentUser_INCLUDED,
			Deleted:   false,
		}
		u = append(u, user)
	}
	return &featureproto.SegmentUsers{
		SegmentId: segmentID,
		Users:     u,
	}
}

func newSegmentUsersCache(t *testing.T, mockController *gomock.Controller) *segmentUsersCache {
	t.Helper()
	return &segmentUsersCache{
		cache: cachemock.NewMockMultiGetCache(mockController),
	}
}
