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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/coderef/domain"
	pubsubmock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/coderef"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestHandleCreateCodeReference(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPublisher := pubsubmock.NewMockPublisher(ctrl)
	now := time.Now()
	codeRef := domain.NewCodeReferenceWithTime(
		"id",
		"feature-id",
		"path/to/file.go",
		10,
		"code snippet",
		"hash123",
		[]string{"alias1", "alias2"},
		"repo",
		"owner",
		proto.CodeReference_GITHUB,
		"main",
		"commit123",
		"env-id",
		now,
	)
	h, err := NewCodeReferenceCommandHandler(
		&eventproto.Editor{
			Email: "test@example.com",
		},
		codeRef,
		mockPublisher,
		"env-id",
	)
	require.NoError(t, err)
	cmd := &proto.CreateCodeReferenceCommand{
		FeatureId:        "feature-id",
		FilePath:         "path/to/file.go",
		LineNumber:       10,
		CodeSnippet:      "code snippet",
		ContentHash:      "hash123",
		Aliases:          []string{"alias1", "alias2"},
		RepositoryName:   "repo",
		RepositoryOwner:  "owner",
		RepositoryType:   proto.CodeReference_GITHUB,
		RepositoryBranch: "main",
		CommitHash:       "commit123",
		EnvironmentId:    "env-id",
	}
	mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}

func TestHandleUpdateCodeReference(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPublisher := pubsubmock.NewMockPublisher(ctrl)
	now := time.Now()
	codeRef := domain.NewCodeReferenceWithTime(
		"id",
		"feature-id",
		"path/to/file.go",
		10,
		"code snippet",
		"hash123",
		[]string{"alias1", "alias2"},
		"repo",
		"owner",
		proto.CodeReference_GITHUB,
		"main",
		"commit123",
		"env-id",
		now,
	)
	h, err := NewCodeReferenceCommandHandler(
		&eventproto.Editor{
			Email: "test@example.com",
		},
		codeRef,
		mockPublisher,
		"env-id",
	)
	require.NoError(t, err)
	cmd := &proto.UpdateCodeReferenceCommand{
		Id:               "id",
		FilePath:         "path/to/new/file.go",
		LineNumber:       20,
		CodeSnippet:      "new code snippet",
		ContentHash:      "newhash123",
		Aliases:          []string{"alias3", "alias4"},
		RepositoryBranch: "develop",
		CommitHash:       "newcommmit123",
	}
	mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}

func TestHandleDeleteCodeReference(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPublisher := pubsubmock.NewMockPublisher(ctrl)
	now := time.Now()
	codeRef := domain.NewCodeReferenceWithTime(
		"id",
		"feature-id",
		"path/to/file.go",
		10,
		"code snippet",
		"hash123",
		[]string{"alias1", "alias2"},
		"repo",
		"owner",
		proto.CodeReference_GITHUB,
		"main",
		"commit123",
		"env-id",
		now,
	)
	h, err := NewCodeReferenceCommandHandler(
		&eventproto.Editor{
			Email: "test@example.com",
		},
		codeRef,
		mockPublisher,
		"env-id",
	)
	require.NoError(t, err)
	cmd := &proto.DeleteCodeReferenceCommand{
		Id:            "id",
		EnvironmentId: "env-id",
	}
	mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}
