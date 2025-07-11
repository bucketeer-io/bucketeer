// Copyright 2025 The Bucketeer Authors.
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
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	cachemock "github.com/bucketeer-io/bucketeer/pkg/cache/mock"
	redismock "github.com/bucketeer-io/bucketeer/pkg/redis/v3/mock"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const (
	testEnvironmentId = "env-id"
)

func TestGetUserAttributeKeyAllCache(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	keyPrefix := fmt.Sprintf("%s:%s", testEnvironmentId, userAttributeKind)
	key := keyPrefix + ":*"
	var cursor uint64
	// Redis keys in the format: environmentId:user_attr:attributeKey
	redisKeys := []string{
		fmt.Sprintf("%s:%s:key1", testEnvironmentId, userAttributeKind),
		fmt.Sprintf("%s:%s:key2", testEnvironmentId, userAttributeKind),
	}
	expectedAttributeKeys := []string{"key1", "key2"}

	patterns := []struct {
		desc        string
		setup       func(*userAttributesCache)
		expectedErr error
	}{
		{
			desc: "error_scan",
			setup: func(uac *userAttributesCache) {
				uac.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(cursor, key, userAttributesMaxSize).Return(cursor, nil, errors.New("scan error"))
			},
			expectedErr: errors.New("scan error"),
		},
		{
			desc: "success",
			setup: func(uac *userAttributesCache) {
				uac.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(cursor, key, userAttributesMaxSize).Return(uint64(0), redisKeys, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: scan twice",
			setup: func(uac *userAttributesCache) {
				uac.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(cursor, key, userAttributesMaxSize).Return(uint64(1), redisKeys[:1], nil)
				uac.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(1), key, userAttributesMaxSize).Return(uint64(0), redisKeys[1:], nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			uac := newUserAttributesCache(t, mockController)
			p.setup(uac)
			attributeKeys, err := uac.GetUserAttributeKeyAll(testEnvironmentId)
			if err == nil {
				assert.Len(t, attributeKeys, len(expectedAttributeKeys))
				for i, expectedKey := range expectedAttributeKeys {
					assert.Equal(t, expectedKey, attributeKeys[i])
				}
			}
			if p.expectedErr != nil {
				assert.EqualError(t, err, p.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPutUserAttributesCache(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userAttrs := &userproto.UserAttributes{
		EnvironmentId: testEnvironmentId,
		UserAttributes: []*userproto.UserAttribute{
			{Key: "key1", Values: []string{"v1", "v2"}},
			{Key: "key2", Values: []string{"v3"}},
		},
	}

	patterns := []struct {
		desc        string
		setup       func(*userAttributesCache, *redismock.MockPipeClient)
		input       *userproto.UserAttributes
		expectedErr error
	}{
		{
			desc:        "nil_input",
			setup:       nil,
			input:       nil,
			expectedErr: errors.New("user attributes is nil"),
		},
		{
			desc: "success",
			setup: func(uac *userAttributesCache, pipe *redismock.MockPipeClient) {
				for _, attr := range userAttrs.UserAttributes {
					key := fmt.Sprintf("%s:%s:%s", testEnvironmentId, userAttributeKind, attr.Key)
					for _, v := range attr.Values {
						pipe.EXPECT().SAdd(key, v)
					}
					pipe.EXPECT().Expire(key, userAttributeTTL)
				}
				pipe.EXPECT().Exec().Return(nil, nil)
			},
			input:       userAttrs,
			expectedErr: nil,
		},
		{
			desc: "error_exec",
			setup: func(uac *userAttributesCache, pipe *redismock.MockPipeClient) {
				for _, attr := range userAttrs.UserAttributes {
					key := fmt.Sprintf("%s:%s:%s", testEnvironmentId, userAttributeKind, attr.Key)
					for _, v := range attr.Values {
						pipe.EXPECT().SAdd(key, v)
					}
					pipe.EXPECT().Expire(key, userAttributeTTL)
				}
				pipe.EXPECT().Exec().Return(nil, errors.New("exec error"))
			},
			input:       userAttrs,
			expectedErr: errors.New("exec error"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			uac := newUserAttributesCache(t, mockController)
			if p.input == nil {
				err := uac.Put(p.input)
				assert.Equal(t, p.expectedErr, err)
				return
			}
			pipe := redismock.NewMockPipeClient(mockController)
			uac.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Pipeline(true).Return(pipe)
			if p.setup != nil {
				p.setup(uac, pipe)
			}
			err := uac.Put(p.input)
			if p.expectedErr != nil {
				assert.EqualError(t, err, p.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newUserAttributesCache(t *testing.T, mockController *gomock.Controller) *userAttributesCache {
	t.Helper()
	return &userAttributesCache{
		cache: cachemock.NewMockMultiGetDeleteCountCache(mockController),
	}
}
