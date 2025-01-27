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

package deleter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	ftstoragemock "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	tagstoragemock "github.com/bucketeer-io/bucketeer/pkg/tag/storage/mock"
	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
	tagproto "github.com/bucketeer-io/bucketeer/proto/tag"
)

func TestRun(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	patterns := []struct {
		desc     string
		setup    func(td *tagDeleter)
		expected error
	}{
		{
			desc: "err: internal error when listing all environment tags",
			setup: func(td *tagDeleter) {
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().ListAllEnvironmentTags(gomock.Any()).Return(
					nil,
					errors.New("internal error"),
				)
			},
			expected: errInternal,
		},
		{
			desc: "err: internal error when listing all environment features",
			setup: func(td *tagDeleter) {
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().ListAllEnvironmentTags(gomock.Any()).Return(
					[]*tagproto.EnvironmentTag{},
					nil,
				)
				td.ftStorage.(*ftstoragemock.MockFeatureStorage).EXPECT().ListAllEnvironmentFeatures(gomock.Any()).Return(
					nil,
					errors.New("internal error"),
				)
			},
			expected: errInternal,
		},
		{
			desc: "err: internal error when deleting tag",
			setup: func(td *tagDeleter) {
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().ListAllEnvironmentTags(gomock.Any()).Return(
					[]*tagproto.EnvironmentTag{
						{
							EnvironmentId: "env-id-1",
							Tags: []*tagproto.Tag{
								{
									Id:            "tag-id-1",
									Name:          "android",
									EnvironmentId: "env-id-1",
								},
							},
						},
					},
					nil,
				)
				td.ftStorage.(*ftstoragemock.MockFeatureStorage).EXPECT().ListAllEnvironmentFeatures(gomock.Any()).Return(
					[]*ftproto.EnvironmentFeature{},
					nil,
				)
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(
					gomock.Any(), "tag-id-1").Return(errors.New("internal error"))
			},
			expected: errInternal,
		},
		{
			desc: "success: delete tags by checking the flags",
			setup: func(td *tagDeleter) {
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().ListAllEnvironmentTags(gomock.Any()).Return(
					[]*tagproto.EnvironmentTag{
						{
							EnvironmentId: "env-id-1",
							Tags: []*tagproto.Tag{
								{
									Id:            "tag-id-1",
									Name:          "android",
									EnvironmentId: "env-id-1",
								},
								{
									Id:            "tag-id-2",
									Name:          "web",
									EnvironmentId: "env-id-1",
								},
								{
									Id:            "tag-id-3",
									Name:          "server",
									EnvironmentId: "env-id-1",
								},
							},
						},
						{
							EnvironmentId: "env-id-2",
							Tags: []*tagproto.Tag{
								{
									Id:            "tag-id-4",
									Name:          "android",
									EnvironmentId: "env-id-2",
								},
								{
									Id:            "tag-id-5",
									Name:          "ios",
									EnvironmentId: "env-id-2",
								},
							},
						},
					},
					nil,
				)
				td.ftStorage.(*ftstoragemock.MockFeatureStorage).EXPECT().ListAllEnvironmentFeatures(gomock.Any()).Return(
					[]*ftproto.EnvironmentFeature{
						{
							EnvironmentId: "env-id-1",
							Features: []*ftproto.Feature{
								{
									Id:   "feature-id-1",
									Tags: []string{"android"},
								},
								{
									Id:   "feature-id-2",
									Tags: []string{"server"},
								},
							},
						},
						{
							EnvironmentId: "env-id-2",
							Features: []*ftproto.Feature{
								{
									Id:   "feature-id-1",
									Tags: []string{"ios"},
								},
							},
						},
					},
					nil,
				)
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(
					gomock.Any(), "tag-id-2").Return(nil)
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(
					gomock.Any(), "tag-id-4").Return(nil)
			},
			expected: nil,
		},
		{
			desc: "success: delete all the tags when there are no flags",
			setup: func(td *tagDeleter) {
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().ListAllEnvironmentTags(gomock.Any()).Return(
					[]*tagproto.EnvironmentTag{
						{
							EnvironmentId: "env-id-1",
							Tags: []*tagproto.Tag{
								{
									Id:            "tag-id-1",
									Name:          "android",
									EnvironmentId: "env-id-1",
								},
								{
									Id:            "tag-id-2",
									Name:          "web",
									EnvironmentId: "env-id-1",
								},
								{
									Id:            "tag-id-3",
									Name:          "server",
									EnvironmentId: "env-id-1",
								},
							},
						},
						{
							EnvironmentId: "env-id-2",
							Tags: []*tagproto.Tag{
								{
									Id:            "tag-id-4",
									Name:          "android",
									EnvironmentId: "env-id-2",
								},
								{
									Id:            "tag-id-5",
									Name:          "ios",
									EnvironmentId: "env-id-2",
								},
							},
						},
					},
					nil,
				)
				td.ftStorage.(*ftstoragemock.MockFeatureStorage).EXPECT().ListAllEnvironmentFeatures(gomock.Any()).Return(
					[]*ftproto.EnvironmentFeature{},
					nil,
				)
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(
					gomock.Any(), "tag-id-1").Return(nil)
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(
					gomock.Any(), "tag-id-2").Return(nil)
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(
					gomock.Any(), "tag-id-3").Return(nil)
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(
					gomock.Any(), "tag-id-4").Return(nil)
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(
					gomock.Any(), "tag-id-5").Return(nil)
			},
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMockTagDeleter(t, mockController)
			p.setup(deleter)
			err := deleter.Run(ctx)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestDeleteUnusedTags(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	patterns := []struct {
		desc          string
		setup         func(td *tagDeleter)
		inputTags     []*tagproto.Tag
		inputFts      []*ftproto.Feature
		expectedCount int
		expectedError error
	}{
		{
			desc: "err: internal error while deleting 2 entries, but only 1 one is deleted",
			setup: func(td *tagDeleter) {
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(ctx, "tag-id-2").Return(nil)
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(ctx, "tag-id-3").
					Return(errors.New("internal error"))
			},
			inputTags: []*tagproto.Tag{
				{
					Id:            "tag-id-1",
					Name:          "android",
					EnvironmentId: "env-id",
				},
				{
					Id:            "tag-id-2",
					Name:          "server",
					EnvironmentId: "env-id",
				},
				{
					Id:            "tag-id-3",
					Name:          "ios",
					EnvironmentId: "env-id",
				},
			},
			inputFts: []*ftproto.Feature{
				{
					Id:   "feature-id-1",
					Tags: []string{"android"},
				},
			},
			expectedCount: 1,
			expectedError: errInternal,
		},
		{
			desc: "success",
			setup: func(td *tagDeleter) {
				td.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(ctx, "tag-id-2").Return(nil)
			},
			inputTags: []*tagproto.Tag{
				{
					Id:            "tag-id-1",
					Name:          "android",
					EnvironmentId: "env-id",
				},
				{
					Id:            "tag-id-2",
					Name:          "server",
					EnvironmentId: "env-id",
				},
				{
					Id:            "tag-id-3",
					Name:          "ios",
					EnvironmentId: "env-id",
				},
			},
			inputFts: []*ftproto.Feature{
				{
					Id:   "feature-id-1",
					Tags: []string{"android"},
				},
				{
					Id:   "feature-id-2",
					Tags: []string{"ios"},
				},
			},
			expectedCount: 1,
			expectedError: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			deleter := newMockTagDeleter(t, mockController)
			p.setup(deleter)
			deletedCount, err := deleter.deleteUnusedTags(ctx, p.inputTags, p.inputFts)
			assert.Equal(t, p.expectedError, err)
			assert.Equal(t, p.expectedCount, deletedCount)
		})
	}
}

func newMockTagDeleter(t *testing.T, c *gomock.Controller) *tagDeleter {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &tagDeleter{
		tagStorage: tagstoragemock.NewMockTagStorage(c),
		ftStorage:  ftstoragemock.NewMockFeatureStorage(c),
		opts: &jobs.Options{
			Timeout: 5 * time.Second,
		},
		logger: logger,
	}
}
