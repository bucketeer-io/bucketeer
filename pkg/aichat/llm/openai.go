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

package llm

import (
	"context"
	"errors"
	"io"

	"github.com/sashabaranov/go-openai"
)

// openaiClient wraps the go-openai client to implement the Client interface.
type openaiClient struct {
	client *openai.Client
}

// NewOpenAIClient creates a new OpenAI-compatible LLM client.
// If baseURL is empty, the default OpenAI endpoint is used.
// This supports any OpenAI-compatible API (e.g., Azure, vLLM, Ollama).
func NewOpenAIClient(apiKey, baseURL string) Client {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	return &openaiClient{
		client: openai.NewClientWithConfig(config),
	}
}

func (c *openaiClient) StreamChat(
	ctx context.Context,
	messages []Message,
	opts StreamOptions,
) (<-chan Chunk, <-chan error) {
	responseChan := make(chan Chunk, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errChan)

		oaiMessages := make([]openai.ChatCompletionMessage, len(messages))
		for i, m := range messages {
			oaiMessages[i] = openai.ChatCompletionMessage{
				Role:    m.Role,
				Content: m.Content,
			}
		}

		req := openai.ChatCompletionRequest{
			Model:       opts.Model,
			Messages:    oaiMessages,
			MaxTokens:   opts.MaxTokens,
			Temperature: opts.Temperature,
			Stream:      true,
		}

		stream, err := c.client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				select {
				case responseChan <- Chunk{Done: true, FinishReason: "stop"}:
				case <-ctx.Done():
				}
				return
			}
			if err != nil {
				errChan <- err
				return
			}

			if len(response.Choices) > 0 {
				content := response.Choices[0].Delta.Content
				if content != "" {
					select {
					case responseChan <- Chunk{Content: content, Done: false}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return responseChan, errChan
}
