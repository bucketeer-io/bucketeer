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
	_ "embed"
	"encoding/json"
	"math"
	"sort"
	"sync"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm"
)

//go:embed data/bucketeer-docs.json
var embeddedDocsJSON []byte

// DocChunk represents a document chunk with its embedding vector.
type DocChunk struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Embedding []float32 `json:"embedding"`
	Metadata  DocMeta   `json:"metadata"`
}

// DocMeta contains metadata about a document chunk.
type DocMeta struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Category string `json:"category"`
}

// Service provides RAG (Retrieval-Augmented Generation) capabilities.
type Service struct {
	llmClient      llm.Client
	embeddingModel string
	docs           []DocChunk
	mu             sync.RWMutex
	logger         *zap.Logger
}

// NewService creates a new RAG service with pre-loaded embedded documents.
func NewService(
	llmClient llm.Client,
	embeddingModel string,
	logger *zap.Logger,
) (*Service, error) {
	if embeddingModel == "" {
		embeddingModel = "text-embedding-3-small"
	}

	svc := &Service{
		llmClient:      llmClient,
		embeddingModel: embeddingModel,
		logger:         logger.Named("rag-service"),
	}

	if err := svc.loadEmbeddedDocs(); err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *Service) loadEmbeddedDocs() error {
	var docs []DocChunk
	if err := json.Unmarshal(embeddedDocsJSON, &docs); err != nil {
		return err
	}
	s.docs = docs
	s.logger.Info("Loaded embedded documents", zap.Int("count", len(docs)))
	return nil
}

// Search finds the most relevant documents for a given query.
func (s *Service) Search(ctx context.Context, query string, topK int) ([]DocChunk, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.docs) == 0 || query == "" {
		return nil, nil
	}

	embeddings, err := s.llmClient.CreateEmbeddings(ctx, s.embeddingModel, []string{query})
	if err != nil {
		return nil, err
	}

	if len(embeddings) == 0 || len(embeddings[0]) == 0 {
		return nil, nil
	}

	queryEmbedding := embeddings[0]

	type scored struct {
		doc   DocChunk
		score float32
	}

	scoredDocs := make([]scored, len(s.docs))
	for i, doc := range s.docs {
		scoredDocs[i] = scored{
			doc:   doc,
			score: CosineSimilarity(queryEmbedding, doc.Embedding),
		}
	}

	sort.Slice(scoredDocs, func(i, j int) bool {
		return scoredDocs[i].score > scoredDocs[j].score
	})

	result := make([]DocChunk, 0, topK)
	for i := 0; i < topK && i < len(scoredDocs); i++ {
		result = append(result, scoredDocs[i].doc)
	}

	return result, nil
}

// CosineSimilarity computes the cosine similarity between two vectors.
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / float32(math.Sqrt(float64(normA)*float64(normB)))
}
