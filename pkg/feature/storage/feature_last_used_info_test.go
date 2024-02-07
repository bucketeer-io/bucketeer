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

package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	storagetesting "github.com/bucketeer-io/bucketeer/pkg/storage/testing"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestNewFeatureLastUsedStorage(t *testing.T) {
	db := NewFeatureLastUsedInfoStorage(storagetesting.NewInMemoryStorage())
	assert.IsType(t, &featureLastUsedInfoStorage{}, db)
}

func TestGetFeatureLastUsedInfos(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		input                []string
		environmentNamespace string
		expected             []*domain.FeatureLastUsedInfo
		expectedErr          error
	}{
		{
			input:                []string{},
			environmentNamespace: "ns0",
			expected:             []*domain.FeatureLastUsedInfo{},
			expectedErr:          nil,
		},
		{
			input:                []string{"feature-id-1:1"},
			environmentNamespace: "ns0",
			expected: []*domain.FeatureLastUsedInfo{
				{
					FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
						FeatureId:  "feature-id-1",
						Version:    1,
						LastUsedAt: 2,
						CreatedAt:  1,
					},
				},
			},
			expectedErr: nil,
		},
		{
			input: []string{
				"feature-id-1:1",
				"feature-id-2:1",
			},
			environmentNamespace: "ns0",
			expected: []*domain.FeatureLastUsedInfo{
				{
					FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
						FeatureId:  "feature-id-1",
						Version:    1,
						LastUsedAt: 2,
						CreatedAt:  1,
					},
				},
				{
					FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
						FeatureId:  "feature-id-2",
						Version:    1,
						LastUsedAt: 2,
						CreatedAt:  1,
					},
				},
			},
			expectedErr: nil,
		},
	}
	client := storagetesting.NewInMemoryStorage()
	keys := []*storage.Key{
		storage.NewKey("feature-id-1:1", featureLastUsedInfoKind, "ns0"),
		storage.NewKey("feature-id-2:1", featureLastUsedInfoKind, "ns0"),
	}
	existedEls := []*proto.FeatureLastUsedInfo{
		{
			FeatureId:  "feature-id-1",
			Version:    1,
			LastUsedAt: 2,
			CreatedAt:  1,
		},
		{
			FeatureId:  "feature-id-2",
			Version:    1,
			LastUsedAt: 2,
			CreatedAt:  1,
		},
	}
	err := client.PutMulti(context.Background(), keys, existedEls)
	assert.NoError(t, err)
	db := NewFeatureLastUsedInfoStorage(client)
	for _, p := range patterns {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		actual, err := db.GetFeatureLastUsedInfos(ctx, p.input, p.environmentNamespace)
		assert.Equal(t, p.expectedErr, err)
		if err == nil && len(p.input) > 0 {
			for i, e := range p.expected {
				assert.NoError(t, err)
				assert.Equal(t, e.FeatureId, actual[i].FeatureId)
				assert.Equal(t, e.Version, actual[i].Version)
				assert.Equal(t, e.LastUsedAt, actual[i].LastUsedAt)
				assert.Equal(t, e.CreatedAt, actual[i].CreatedAt)
			}
		}
	}
}

func TestUpsertFeatureLastUsedInfo(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		data                 []*domain.FeatureLastUsedInfo
		environmentNamespace string
		expectedErr          error
	}{
		// insert
		{
			data: []*domain.FeatureLastUsedInfo{
				{
					FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
						FeatureId:  "feature-id-10",
						Version:    1,
						LastUsedAt: 2,
						CreatedAt:  1,
					},
				},
			},
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
		// multi insert
		{
			data: []*domain.FeatureLastUsedInfo{
				{
					FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
						FeatureId:  "feature-id-11",
						Version:    1,
						LastUsedAt: 2,
						CreatedAt:  1,
					},
				},
				{
					FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
						FeatureId:  "feature-id-12",
						Version:    1,
						LastUsedAt: 2,
						CreatedAt:  1,
					},
				},
			},
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
		// update
		{
			data: []*domain.FeatureLastUsedInfo{
				{
					FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
						FeatureId:  "feature-id-1",
						Version:    1,
						LastUsedAt: 3,
						CreatedAt:  1,
					},
				},
			},
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
		// insert & update
		{
			data: []*domain.FeatureLastUsedInfo{
				{
					FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
						FeatureId:  "feature-id-2",
						Version:    1,
						LastUsedAt: 3,
						CreatedAt:  1,
					},
				},
				{
					FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
						FeatureId:  "feature-id-13",
						Version:    1,
						LastUsedAt: 3,
						CreatedAt:  1,
					},
				},
			},
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
	}
	client := storagetesting.NewInMemoryStorage()
	keys := []*storage.Key{
		storage.NewKey("feature-id-1:1", featureLastUsedInfoKind, "ns0"),
		storage.NewKey("feature-id-2:1", featureLastUsedInfoKind, "ns0"),
	}
	existedEls := []*proto.FeatureLastUsedInfo{
		{
			FeatureId:  "feature-id-1",
			Version:    1,
			LastUsedAt: 2,
			CreatedAt:  1,
		},
		{
			FeatureId:  "feature-id-2",
			Version:    1,
			LastUsedAt: 2,
			CreatedAt:  1,
		},
	}
	err := client.PutMulti(context.Background(), keys, existedEls)
	assert.NoError(t, err)
	db := NewFeatureLastUsedInfoStorage(client)
	for _, p := range patterns {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := db.UpsertFeatureLastUsedInfos(ctx, p.data, p.environmentNamespace)
		assert.NoError(t, err)
		actual := make([]*proto.FeatureLastUsedInfo, len(p.data))
		for i := range actual {
			actual[i] = &proto.FeatureLastUsedInfo{}
		}
		keys := make([]*storage.Key, 0, len(p.data))
		for _, d := range p.data {
			keys = append(keys, storage.NewKey(d.ID(), featureLastUsedInfoKind, p.environmentNamespace))
		}
		err = client.GetMulti(ctx, keys, actual)
		assert.NoError(t, err)
		for i, e := range p.data {
			assert.NoError(t, err)
			assert.Equal(t, e.FeatureLastUsedInfo.FeatureId, actual[i].FeatureId)
			assert.Equal(t, e.FeatureLastUsedInfo.Version, actual[i].Version)
			assert.Equal(t, e.FeatureLastUsedInfo.LastUsedAt, actual[i].LastUsedAt)
			assert.Equal(t, e.FeatureLastUsedInfo.CreatedAt, actual[i].CreatedAt)
		}
	}
}

func TestNewFeatureLastUsedLister(t *testing.T) {
	db := NewFeatureLastUsedInfoLister(storagetesting.NewInMemoryStorage())
	assert.IsType(t, &featureLastUsedInfoLister{}, db)
}
