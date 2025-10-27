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

package rediscounter

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"

	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	evmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	evproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

var (
	inputRequest = &evproto.ListEnvironmentsV2Request{
		PageSize: 0,
		Archived: &wrapperspb.BoolValue{Value: false},
	}

	inputEnvironments = []*evproto.EnvironmentV2{
		{
			Id: "dev",
		},
		{
			Id: "prd",
		},
	}
)

func TestRun(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	patterns := []struct {
		desc     string
		setup    func(r *redisCounterDeleter)
		expected error
	}{
		{
			desc: "list environments internal error",
			setup: func(r *redisCounterDeleter) {
				r.envClient.(*evmock.MockClient).EXPECT().ListEnvironmentsV2(gomock.Any(), inputRequest).Return(
					nil, errors.New("internal error"))
			},
			expected: errors.New("internal error"),
		},
		{
			desc: "redis scan internal error",
			setup: func(r *redisCounterDeleter) {
				r.envClient.(*evmock.MockClient).EXPECT().ListEnvironmentsV2(gomock.Any(), inputRequest).Return(
					&evproto.ListEnvironmentsV2Response{
						Environments: inputEnvironments,
					}, nil)
				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "dev:uc:*", redisScanMaxSize).Return(
					uint64(0), nil, errors.New("internal error"))
			},
			expected: errors.New("internal error"),
		},
		{
			desc: "success: no keys found",
			setup: func(r *redisCounterDeleter) {
				r.envClient.(*evmock.MockClient).EXPECT().ListEnvironmentsV2(gomock.Any(), inputRequest).Return(
					&evproto.ListEnvironmentsV2Response{
						Environments: inputEnvironments,
					}, nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "dev:uc:*", redisScanMaxSize).Return(
					uint64(0), nil, nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "dev:ec:*", redisScanMaxSize).Return(
					uint64(0), nil, nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "prd:uc:*", redisScanMaxSize).Return(
					uint64(0), nil, nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "prd:ec:*", redisScanMaxSize).Return(
					uint64(0), nil, nil)
			},
		},
		{
			desc: "error while deleting keys",
			setup: func(r *redisCounterDeleter) {
				r.envClient.(*evmock.MockClient).EXPECT().ListEnvironmentsV2(gomock.Any(), inputRequest).Return(
					&evproto.ListEnvironmentsV2Response{
						Environments: inputEnvironments,
					}, nil)
				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "dev:uc:*", redisScanMaxSize).Return(
					uint64(0), makeDummyKeys(t, "dev", "uc", 31, 1), nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Delete(gomock.Any()).Return(errors.New("delete error")).AnyTimes()
			},
			expected: fmt.Errorf("failed to delete key %s: %w",
				fmt.Sprintf("dev:uc:%d:feature_id_0:variation_id_0", time.Now().Unix()-(31*24*60*60)),
				errors.New("delete error")),
		},
		{
			desc: "success: no keys older than 31 days",
			setup: func(r *redisCounterDeleter) {
				r.envClient.(*evmock.MockClient).EXPECT().ListEnvironmentsV2(gomock.Any(), inputRequest).Return(
					&evproto.ListEnvironmentsV2Response{
						Environments: inputEnvironments,
					}, nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "dev:uc:*", redisScanMaxSize).Return(
					uint64(0), makeDummyKeys(t, "dev", "uc", 7, 100), nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "dev:ec:*", redisScanMaxSize).Return(
					uint64(0), makeDummyKeys(t, "dev", "ec", 7, 100), nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "prd:uc:*", redisScanMaxSize).Return(
					uint64(0), makeDummyKeys(t, "prd", "uc", 7, 100), nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "prd:ec:*", redisScanMaxSize).Return(
					uint64(0), makeDummyKeys(t, "prd", "ec", 7, 100), nil)
			},
			expected: nil,
		},
		{
			desc: "success",
			setup: func(r *redisCounterDeleter) {
				r.envClient.(*evmock.MockClient).EXPECT().ListEnvironmentsV2(gomock.Any(), inputRequest).Return(
					&evproto.ListEnvironmentsV2Response{
						Environments: inputEnvironments,
					}, nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "dev:uc:*", redisScanMaxSize).Return(
					uint64(0), makeDummyKeys(t, "dev", "uc", 31, 100), nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "dev:ec:*", redisScanMaxSize).Return(
					uint64(0), makeDummyKeys(t, "dev", "ec", 31, 100), nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "prd:uc:*", redisScanMaxSize).Return(
					uint64(0), makeDummyKeys(t, "prd", "uc", 31, 100), nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "prd:ec:*", redisScanMaxSize).Return(
					uint64(0), makeDummyKeys(t, "prd", "ec", 31, 150), nil)

				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Delete(gomock.Any()).Return(nil).Times(450)
			},
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMockRedisCounterDeleter(t, mockController)
			p.setup(deleter)
			err := deleter.Run(ctx)
			if p.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, p.expected.Error())
			}
		})
	}
}

func TestListEnvironments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	patterns := []struct {
		desc          string
		setup         func(r *redisCounterDeleter)
		expected      []*evproto.EnvironmentV2
		expectedError error
	}{
		{
			desc: "err: internal",
			setup: func(r *redisCounterDeleter) {
				r.envClient.(*evmock.MockClient).EXPECT().ListEnvironmentsV2(ctx, inputRequest).Return(
					nil, errors.New("internal error"))
			},
			expected:      nil,
			expectedError: errors.New("internal error"),
		},
		{
			desc: "success",
			setup: func(r *redisCounterDeleter) {
				r.envClient.(*evmock.MockClient).EXPECT().ListEnvironmentsV2(gomock.Any(), inputRequest).Return(
					&evproto.ListEnvironmentsV2Response{
						Environments: inputEnvironments,
					}, nil)
			},
			expected:      inputEnvironments,
			expectedError: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMockRedisCounterDeleter(t, mockController)
			p.setup(deleter)
			envs, err := deleter.listEnvironments(ctx)
			for i := 0; i < len(envs); i++ {
				assert.True(t, proto.Equal(p.expected[i], envs[i]))
			}
			assert.Equal(t, p.expectedError, err)
		})
	}
}

func TestNewKeyPrefix(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	t.Run("success", func(t *testing.T) {
		deleter := newMockRedisCounterDeleter(t, mockController)
		keyPrefix := deleter.newKeyPrefix("dev", "uc")
		assert.Equal(t, "dev:uc:*", keyPrefix)
	})
}

func TestScan(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	now := time.Now()
	inputKeys := []string{
		fmt.Sprintf("dev:uc:%d:feature_id_1:variation_id", now.Unix()-3*day),
		fmt.Sprintf("dev:uc:%d:feature_id_2:variation_id", now.Unix()-31*day),
		fmt.Sprintf("dev:uc:%d:feature_id_3:variation_id", now.Unix()-7*day),
	}

	patterns := []struct {
		desc          string
		setup         func(r *redisCounterDeleter)
		expected      []string
		expectedError error
	}{
		{
			desc: "err: internal",
			setup: func(r *redisCounterDeleter) {
				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "key", redisScanMaxSize).Return(
					uint64(0), nil, errors.New("internal error"))
			},
			expected:      nil,
			expectedError: errors.New("internal error"),
		},
		{
			desc: "success",
			setup: func(r *redisCounterDeleter) {
				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Scan(uint64(0), "key", redisScanMaxSize).Return(
					uint64(0), inputKeys, nil)
			},
			expected:      inputKeys,
			expectedError: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMockRedisCounterDeleter(t, mockController)
			p.setup(deleter)
			keys, err := deleter.scan("dev", "uc", "key")
			assert.Equal(t, p.expected, keys)
			assert.Equal(t, p.expectedError, err)
		})
	}
}

func TestFilterKeysOlderThanThirtyOneDays(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	now := time.Now()

	patterns := []struct {
		desc              string
		inputEnvNamespace string
		inputKind         string
		inputKeys         []string
		expected          []string
		expectedError     error
	}{
		{
			desc:              "success: skip pfmerge keys",
			inputEnvNamespace: "",
			inputKind:         "uc",
			inputKeys: []string{
				"uc:pfmerge-key:ios-premium-tab-enabled",
				"uc:pfmerge-key:android-feature",
				fmt.Sprintf("uc:%d:feature_id:variation_id", now.Unix()-31*day),
			},
			expected: []string{
				fmt.Sprintf("uc:%d:feature_id:variation_id", now.Unix()-31*day),
			},
			expectedError: nil,
		},
		{
			desc:              "success: skip pfmerge keys with environment id",
			inputEnvNamespace: "dev",
			inputKind:         "uc",
			inputKeys: []string{
				"dev:pfmerge:uc:pfmerge-key:feature_id:variation_id",
				fmt.Sprintf("dev:uc:%d:feature_id_1:variation_id", now.Unix()-3*day),
				fmt.Sprintf("dev:uc:%d:feature_id_2:variation_id", now.Unix()-31*day),
			},
			expected: []string{
				fmt.Sprintf("dev:uc:%d:feature_id_2:variation_id", now.Unix()-31*day),
			},
			expectedError: nil,
		},
		{
			desc:              "errSubmatchStringNotFound",
			inputEnvNamespace: "dev",
			inputKind:         "uc",
			inputKeys:         []string{fmt.Sprintf("dev:uc:%d", now.Unix())},
			expected:          nil,
			expectedError:     errSubmatchStringNotFound,
		},
		{
			desc:              "success: using empty environment id",
			inputEnvNamespace: "",
			inputKind:         "uc",
			inputKeys: []string{
				fmt.Sprintf("uc:%d:feature_id_1:variation_id", now.Unix()-3*day),
				fmt.Sprintf("uc:%d:feature_id_2:variation_id", now.Unix()-31*day),
				fmt.Sprintf("uc:%d:feature_id_3:variation_id", now.Unix()-7*day),
			},
			expected: []string{
				fmt.Sprintf("uc:%d:feature_id_2:variation_id", now.Unix()-31*day),
			},
			expectedError: nil,
		},
		{
			desc:              "success",
			inputEnvNamespace: "dev",
			inputKind:         "uc",
			inputKeys: []string{
				fmt.Sprintf("dev:uc:%d:feature_id_1:variation_id", now.Unix()-3*day),
				fmt.Sprintf("dev:uc:%d:feature_id_2:variation_id", now.Unix()-31*day),
				fmt.Sprintf("dev:uc:%d:feature_id_3:variation_id", now.Unix()-7*day),
			},
			expected: []string{
				fmt.Sprintf("dev:uc:%d:feature_id_2:variation_id", now.Unix()-31*day),
			},
			expectedError: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMockRedisCounterDeleter(t, mockController)
			keys, err := deleter.filterKeysOlderThanThirtyOneDays(
				p.inputEnvNamespace,
				p.inputKind,
				p.inputKeys,
			)
			assert.Equal(t, p.expected, keys)
			assert.Equal(t, p.expectedError, err)
		})
	}
}

func TestChunkSlice(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		inputKeys      []string
		inputChunkSize int
		expected       int
	}{
		{
			desc:           "success: chunks size 11",
			inputKeys:      makeDummyKeys(t, "dev", "uc", 7, 315),
			inputChunkSize: 30,
			expected:       11,
		},
		{
			desc:           "success: chunks size 6",
			inputKeys:      makeDummyKeys(t, "dev", "uc", 7, 550),
			inputChunkSize: redisChunkMaxSize,
			expected:       6,
		},
		{
			desc:           "success: chunks size 1",
			inputKeys:      makeDummyKeys(t, "dev", "uc", 7, 99),
			inputChunkSize: redisChunkMaxSize,
			expected:       1,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMockRedisCounterDeleter(t, mockController)
			chunks := deleter.chunkSlice(
				p.inputKeys,
				p.inputChunkSize,
			)
			assert.Equal(t, p.expected, len(chunks))
		})
	}
}

func TestDeleteKeys(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc     string
		setup    func(r *redisCounterDeleter)
		input    []string
		expected error
	}{
		{
			desc: "error while deleting keys",
			setup: func(r *redisCounterDeleter) {
				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Delete(gomock.Any()).Return(errors.New("delete error")).AnyTimes()
			},
			input: makeDummyKeys(t, "prd", "ec", 31, 100),
			expected: fmt.Errorf("failed to delete key %s: %w",
				fmt.Sprintf("prd:ec:%d:feature_id_0:variation_id_0", time.Now().Unix()-(31*24*60*60)),
				errors.New("delete error")),
		},
		{
			desc: "success",
			setup: func(r *redisCounterDeleter) {
				r.cache.(*cachemock.MockMultiGetDeleteCountCache).EXPECT().Delete(gomock.Any()).Return(nil).Times(150)
			},
			input:    makeDummyKeys(t, "prd", "ec", 31, 150),
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMockRedisCounterDeleter(t, mockController)
			p.setup(deleter)
			err := deleter.deleteKeys(p.input)
			if p.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, p.expected.Error())
			}
		})
	}
}

func makeDummyKeys(t *testing.T, environmentId, kind string, days, size int) []string {
	t.Helper()
	now := time.Now()
	keys := make([]string, 0, size)
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("%s:%s:%d:feature_id_%d:variation_id_%d", environmentId, kind, now.Unix()-(int64(days)*day), i, i)
		keys = append(keys, key)
	}
	return keys
}

func newMockRedisCounterDeleter(t *testing.T, c *gomock.Controller) *redisCounterDeleter {
	t.Helper()
	return &redisCounterDeleter{
		envClient: evmock.NewMockClient(c),
		cache:     cachemock.NewMockMultiGetDeleteCountCache(c),
		opts: &jobs.Options{
			Timeout: 5 * time.Second,
		},
		logger: zap.NewNop().Named("test-redis-counter-deleter"),
	}
}
