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

import "context"

// DocChunk represents a document chunk returned by a search.
type DocChunk struct {
	ID       string  `json:"id"`
	Content  string  `json:"content"`
	Metadata DocMeta `json:"metadata"`
}

// DocMeta contains metadata about a document chunk.
type DocMeta struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Category string `json:"category"`
}

// Searcher defines the interface for document search.
type Searcher interface {
	Search(ctx context.Context, query string, topK int) ([]DocChunk, error)
}
