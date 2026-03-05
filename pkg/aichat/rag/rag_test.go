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

package rag

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm/mock"
)

func TestCosineSimilarity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "opposite vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{-1, 0, 0},
			expected: -1.0,
		},
		{
			name:     "similar vectors",
			a:        []float32{1, 1, 0},
			b:        []float32{1, 0, 0},
			expected: float32(1.0 / math.Sqrt(2)),
		},
		{
			name:     "empty vectors",
			a:        []float32{},
			b:        []float32{},
			expected: 0.0,
		},
		{
			name:     "different length vectors",
			a:        []float32{1, 0},
			b:        []float32{1, 0, 0},
			expected: 0.0,
		},
		{
			name:     "zero vector",
			a:        []float32{0, 0, 0},
			b:        []float32{1, 1, 1},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := CosineSimilarity(tt.a, tt.b)
			assert.InDelta(t, float64(tt.expected), float64(result), 0.0001)
		})
	}
}

func TestServiceSearch(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)

	// Create a service with test documents
	svc := &Service{
		llmClient:      mockClient,
		embeddingModel: "test-model",
		docs: []DocChunk{
			{
				ID:        "doc-1",
				Content:   "Feature flags help with gradual rollouts",
				Embedding: []float32{1, 0, 0},
				Metadata:  DocMeta{Title: "Feature Flags", URL: "https://example.com/flags", Category: "feature-flags"},
			},
			{
				ID:        "doc-2",
				Content:   "Segments group users for targeting",
				Embedding: []float32{0, 1, 0},
				Metadata:  DocMeta{Title: "Segments", URL: "https://example.com/segments", Category: "segments"},
			},
			{
				ID:        "doc-3",
				Content:   "A/B testing with experiments",
				Embedding: []float32{0, 0, 1},
				Metadata:  DocMeta{Title: "Experiments", URL: "https://example.com/experiments", Category: "experiments"},
			},
		},
		logger: zap.NewNop(),
	}

	t.Run("returns top K results sorted by similarity", func(t *testing.T) {
		// Query embedding is closest to doc-1
		mockClient.EXPECT().
			CreateEmbeddings(gomock.Any(), "test-model", []string{"feature flags"}).
			Return([][]float32{{0.9, 0.1, 0}}, nil)

		results, err := svc.Search(context.Background(), "feature flags", 2)
		require.NoError(t, err)
		require.Len(t, results, 2)
		assert.Equal(t, "doc-1", results[0].ID)
	})

	t.Run("returns empty for empty query", func(t *testing.T) {
		results, err := svc.Search(context.Background(), "", 3)
		require.NoError(t, err)
		assert.Nil(t, results)
	})

	t.Run("returns empty for empty docs", func(t *testing.T) {
		emptySvc := &Service{
			llmClient:      mockClient,
			embeddingModel: "test-model",
			docs:           []DocChunk{},
			logger:         zap.NewNop(),
		}
		results, err := emptySvc.Search(context.Background(), "test", 3)
		require.NoError(t, err)
		assert.Nil(t, results)
	})

	t.Run("handles topK larger than docs", func(t *testing.T) {
		mockClient.EXPECT().
			CreateEmbeddings(gomock.Any(), "test-model", []string{"test"}).
			Return([][]float32{{1, 0, 0}}, nil)

		results, err := svc.Search(context.Background(), "test", 10)
		require.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("returns error on embedding failure", func(t *testing.T) {
		mockClient.EXPECT().
			CreateEmbeddings(gomock.Any(), "test-model", []string{"error"}).
			Return(nil, assert.AnError)

		results, err := svc.Search(context.Background(), "error", 3)
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestLoadEmbeddedDocs(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)

	svc, err := NewService(mockClient, "", zap.NewNop())
	require.NoError(t, err)
	assert.NotEmpty(t, svc.docs)
}
