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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachemock "github.com/bucketeer-io/bucketeer/pkg/cache/mock"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	tag                  = "bucketeer-tag"
	environmentNamespace = "bucketeer-environment"
)

func TestGetFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	features := createFeatures(t)
	dataFeatures := marshalMessage(t, features)
	key := fmt.Sprintf("%s:%s", environmentNamespace, featuresKind)

	patterns := []struct {
		desc        string
		setup       func(*featuresCache)
		expectedErr error
	}{
		{
			desc: "error_get_not_found",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().Get(key).Return(nil, cache.ErrNotFound)
			},
			expectedErr: cache.ErrNotFound,
		},
		{
			desc: "error_invalid_type",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().Get(key).Return("test", nil)
			},
			expectedErr: cache.ErrInvalidType,
		},
		{
			desc: "success",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().Get(key).Return(dataFeatures, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			tf := newFeaturesCache(t, mockController)
			p.setup(tf)
			features, err := tf.Get(environmentNamespace)
			if err == nil {
				assert.Equal(t, features.Features[0].Id, features.Features[0].Id)
				assert.Equal(t, features.Features[0].Name, features.Features[0].Name)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestPutFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	features := createFeatures(t)
	dataFeatures := marshalMessage(t, features)
	key := fmt.Sprintf("%s:%s", environmentNamespace, featuresKind)

	patterns := []struct {
		desc        string
		setup       func(*featuresCache)
		input       *featureproto.Features
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
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().Put(key, dataFeatures, featuresTTL).Return(nil)
			},
			input:       features,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			tf := newFeaturesCache(t, mockController)
			if p.setup != nil {
				p.setup(tf)
			}
			err := tf.Put(p.input, environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func createFeatures(t *testing.T) *featureproto.Features {
	t.Helper()
	f := []*featureproto.Feature{}
	for i := 0; i < 5; i++ {
		feature := &featureproto.Feature{
			Id:   fmt.Sprintf("feature-id-%d", i),
			Name: fmt.Sprintf("feature-name-%d", i),
		}
		f = append(f, feature)
	}
	return &featureproto.Features{
		Features: f,
	}
}

func marshalMessage(t *testing.T, pb proto.Message) interface{} {
	t.Helper()
	buffer, err := proto.Marshal(pb)
	require.NoError(t, err)
	return buffer
}

func newFeaturesCache(t *testing.T, mockController *gomock.Controller) *featuresCache {
	t.Helper()
	return &featuresCache{
		cache: cachemock.NewMockMultiGetCache(mockController),
	}
}
