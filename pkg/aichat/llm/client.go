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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package llm

import "context"

// Chunk represents a single piece of a streaming LLM response.
type Chunk struct {
	Content      string
	Done         bool
	FinishReason string
}

// Role represents a chat message role.
type Role = string

// Role constants for LLM messages.
const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message represents a chat message with a role.
type Message struct {
	Role    Role
	Content string
}

// StreamOptions contains configuration for a streaming chat request.
type StreamOptions struct {
	Model       string
	MaxTokens   int
	Temperature float32
}

// Client defines the interface for LLM providers.
// Implementations can wrap OpenAI, Anthropic, Google, etc.
type Client interface {
	// StreamChat sends messages to the LLM and returns channels for
	// streaming response chunks and errors.
	StreamChat(ctx context.Context, messages []Message, opts StreamOptions) (<-chan Chunk, <-chan error)

	// CreateEmbeddings generates embedding vectors for the given inputs.
	CreateEmbeddings(ctx context.Context, model string, input []string) ([][]float32, error)
}
