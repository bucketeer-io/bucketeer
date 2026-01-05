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

package v3

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3/mock"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	testEnvironmentId = "env-id"
)

func TestGetUserAttributeKeyAllCache(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	indexKey := fmt.Sprintf("%s:%s:keys", testEnvironmentId, userAttributeKind)
	expectedAttributeKeys := []string{"key1", "key2"}

	patterns := []struct {
		desc        string
		setup       func(*userAttributesCache)
		expectedErr error
	}{
		{
			desc: "error_smembers",
			setup: func(uac *userAttributesCache) {
				uac.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().SMembers(indexKey).Return(nil, errors.New("smembers error"))
			},
			expectedErr: errors.New("smembers error"),
		},
		{
			desc: "success",
			setup: func(uac *userAttributesCache) {
				uac.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().SMembers(indexKey).Return(expectedAttributeKeys, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success_empty",
			setup: func(uac *userAttributesCache) {
				uac.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().SMembers(indexKey).Return([]string{}, nil)
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
				if p.desc == "success" {
					assert.Len(t, attributeKeys, len(expectedAttributeKeys))
					for i, expectedKey := range expectedAttributeKeys {
						assert.Equal(t, expectedKey, attributeKeys[i])
					}
				} else if p.desc == "success_empty" {
					assert.Len(t, attributeKeys, 0)
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
		ttl         time.Duration
		expectedErr error
	}{
		{
			desc:        "nil_input",
			setup:       nil,
			input:       nil,
			ttl:         30 * time.Second,
			expectedErr: errors.New("userAttributes cannot be nil"),
		},
		{
			desc: "success_with_custom_ttl",
			setup: func(uac *userAttributesCache, pipe *redismock.MockPipeClient) {
				indexKey := fmt.Sprintf("%s:%s:keys", testEnvironmentId, userAttributeKind)

				for _, attr := range userAttrs.UserAttributes {
					// 1) per-attribute value set
					attrKey := fmt.Sprintf("%s:%s:%s", testEnvironmentId, userAttributeKind, attr.Key)
					members := make([]interface{}, len(attr.Values))
					for i, v := range attr.Values {
						members[i] = v
					}
					pipe.EXPECT().SAdd(attrKey, members...)
					pipe.EXPECT().Expire(attrKey, 7*time.Second)

					// 2) record the attribute.Key in the index set
					pipe.EXPECT().SAdd(indexKey, attr.Key)
					pipe.EXPECT().Expire(indexKey, 7*time.Second)
				}
				pipe.EXPECT().Exec().Return(nil, nil)
			},
			input:       userAttrs,
			ttl:         7 * time.Second,
			expectedErr: nil,
		},
		{
			desc: "success_with_zero_ttl",
			setup: func(uac *userAttributesCache, pipe *redismock.MockPipeClient) {
				indexKey := fmt.Sprintf("%s:%s:keys", testEnvironmentId, userAttributeKind)

				for _, attr := range userAttrs.UserAttributes {
					// 1) per-attribute value set
					attrKey := fmt.Sprintf("%s:%s:%s", testEnvironmentId, userAttributeKind, attr.Key)
					members := make([]interface{}, len(attr.Values))
					for i, v := range attr.Values {
						members[i] = v
					}
					pipe.EXPECT().SAdd(attrKey, members...)
					pipe.EXPECT().Expire(attrKey, 0*time.Second)

					// 2) record the attribute.Key in the index set
					pipe.EXPECT().SAdd(indexKey, attr.Key)
					pipe.EXPECT().Expire(indexKey, 0*time.Second)
				}
				pipe.EXPECT().Exec().Return(nil, nil)
			},
			input:       userAttrs,
			ttl:         0 * time.Second,
			expectedErr: nil,
		},
		{
			desc: "error_exec",
			setup: func(uac *userAttributesCache, pipe *redismock.MockPipeClient) {
				indexKey := fmt.Sprintf("%s:%s:keys", testEnvironmentId, userAttributeKind)

				for _, attr := range userAttrs.UserAttributes {
					// 1) per-attribute value set
					attrKey := fmt.Sprintf("%s:%s:%s", testEnvironmentId, userAttributeKind, attr.Key)
					members := make([]interface{}, len(attr.Values))
					for i, v := range attr.Values {
						members[i] = v
					}
					pipe.EXPECT().SAdd(attrKey, members...)
					pipe.EXPECT().Expire(attrKey, 1*time.Hour)

					// 2) record the attribute.Key in the index set
					pipe.EXPECT().SAdd(indexKey, attr.Key)
					pipe.EXPECT().Expire(indexKey, 1*time.Hour)
				}
				pipe.EXPECT().Exec().Return(nil, errors.New("exec error"))
			},
			input:       userAttrs,
			ttl:         1 * time.Hour,
			expectedErr: errors.New("exec error"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			uac := newUserAttributesCache(t, mockController)
			if p.input == nil {
				err := uac.Put(p.input, p.ttl)
				assert.Equal(t, p.expectedErr, err)
				return
			}
			pipe := redismock.NewMockPipeClient(mockController)
			uac.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Pipeline(false).Return(pipe)
			if p.setup != nil {
				p.setup(uac, pipe)
			}
			err := uac.Put(p.input, p.ttl)
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
