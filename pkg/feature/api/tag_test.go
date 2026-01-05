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

package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	tagstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage/mock"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	tagproto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
)

func TestListTagsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithToken()
	service := createFeatureServiceNew(mockController)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*FeatureService)
		input       *featureproto.ListTagsRequest
		expected    *featureproto.ListTagsResponse
		expectedErr error
	}{
		{
			desc:        "errInvalidCursor",
			setup:       nil,
			input:       &featureproto.ListTagsRequest{EnvironmentId: environmentId, Cursor: "foo"},
			expected:    nil,
			expectedErr: statusInvalidCursor.Err(),
		},
		{
			desc: "errInternal",
			setup: func(fs *FeatureService) {
				fs.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().ListTags(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "test"))
			},
			input:       &featureproto.ListTagsRequest{EnvironmentId: environmentId},
			expected:    nil,
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "test")).Err(),
		},
		{
			desc: "success",
			setup: func(fs *FeatureService) {
				fs.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().ListTags(
					gomock.Any(), gomock.Any(),
				).Return([]*tagproto.Tag{}, 0, int64(0), nil)
			},
			input: &featureproto.ListTagsRequest{
				PageSize:      2,
				Cursor:        "",
				EnvironmentId: environmentId,
			},
			expected:    &featureproto.ListTagsResponse{Tags: []*featureproto.Tag{}, Cursor: "0"},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.ListTags(ctx, p.input)
			assert.Equal(t, p.expected, actual, p.desc)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestUpsertTags(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	service := createFeatureServiceNew(mockController)
	ctx := createContextWithToken()
	internalErr := errors.New("test")
	patterns := []struct {
		desc        string
		tags        []string
		setup       func(fs *FeatureService)
		expectedErr error
	}{
		{
			desc: "error: internal error when creating tag",
			tags: []string{"tag"},
			setup: func(fs *FeatureService) {
				fs.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().UpsertTag(
					gomock.Any(), gomock.Any(),
				).Return(internalErr)
			},
			expectedErr: internalErr,
		},
		{
			desc:        "success: tag with whitespaces",
			tags:        []string{" "},
			expectedErr: nil,
		},
		{
			desc: "success: create new tag",
			tags: []string{"tag"},
			setup: func(fs *FeatureService) {
				fs.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().UpsertTag(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			if p.setup != nil {
				p.setup(service)
			}
			err := service.upsertTags(ctx, p.tags, environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
