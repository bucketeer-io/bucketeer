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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/rag"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/ratelimit"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	aichatproto "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
)

const maxRequestBodyBytes = 64 * 1024 // 64KB

const (
	httpPageTypeFeatureFlags = "feature_flags"
	httpPageTypeTargeting    = "targeting"
	httpPageTypeExperiments  = "experiments"
	httpPageTypeSegments     = "segments"
	httpPageTypeAutoops      = "autoops"
)

// sseDataEvent represents the data payload for an SSE data event.
type sseDataEvent struct {
	Content      string `json:"content"`
	Done         bool   `json:"done"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// sseErrorEvent represents the data payload for an SSE error event.
type sseErrorEvent struct {
	Error string `json:"error"`
}

// chatHTTPService handles SSE streaming for chat requests over HTTP.
// This is needed because gRPC-Gateway does not support server streaming as SSE.
type chatHTTPService struct {
	llmClient     llm.Client
	ragSearcher   rag.Searcher
	chatConfig    ChatConfig
	verifier      token.Verifier
	accountClient accountclient.Client
	featureClient featureclient.Client
	rateLimiter   *ratelimit.Limiter
	logger        *zap.Logger
}

type chatHTTPServiceOptions struct {
	rateLimiter *ratelimit.Limiter
}

// ChatHTTPServiceOption is a functional option for chatHTTPService.
type ChatHTTPServiceOption func(*chatHTTPServiceOptions)

// WithRateLimiter sets the rate limiter for the chat HTTP service.
func WithRateLimiter(l *ratelimit.Limiter) ChatHTTPServiceOption {
	return func(o *chatHTTPServiceOptions) {
		o.rateLimiter = l
	}
}

// NewChatHTTPService creates a new chat HTTP service for SSE streaming.
func NewChatHTTPService(
	llmClient llm.Client,
	ragSearcher rag.Searcher,
	chatConfig ChatConfig,
	verifier token.Verifier,
	accountClient accountclient.Client,
	featureClient featureclient.Client,
	logger *zap.Logger,
	opts ...ChatHTTPServiceOption,
) *chatHTTPService {
	dopts := &chatHTTPServiceOptions{}
	for _, opt := range opts {
		opt(dopts)
	}
	return &chatHTTPService{
		llmClient:     llmClient,
		ragSearcher:   ragSearcher,
		chatConfig:    defaultChatConfig(chatConfig),
		verifier:      verifier,
		accountClient: accountClient,
		featureClient: featureClient,
		rateLimiter:   dopts.rateLimiter,
		logger:        logger.Named("sse"),
	}
}

// ChatHTTPRequest represents the HTTP request body for chat.
type ChatHTTPRequest struct {
	Messages      []ChatMessageHTTP `json:"messages"`
	PageContext   *PageContextHTTP  `json:"pageContext,omitempty"`
	EnvironmentID string            `json:"environmentId"`
}

// ChatMessageHTTP represents a chat message in HTTP format.
type ChatMessageHTTP struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// PageContextHTTP represents page context in HTTP format.
type PageContextHTTP struct {
	PageType  string            `json:"pageType"`
	FeatureID string            `json:"featureId,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Register registers the SSE handler with the HTTP mux.
func (h *chatHTTPService) Register(mux *http.ServeMux) {
	mux.Handle("/v1/aichat/chat", h)
}

// writeJSONError writes a JSON-formatted error response.
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message}) //nolint:errcheck
}

// ServeHTTP handles POST /v1/aichat/chat with SSE streaming.
func (h *chatHTTPService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authenticate
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		writeJSONError(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.verifier.VerifyAccessToken(parts[1])
	if err != nil {
		h.logger.Warn("authentication failed", zap.Error(err))
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Rate limit check (per email)
	if h.rateLimiter != nil && !h.rateLimiter.Allow(accessToken.Email) {
		writeJSONError(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// Parse request body with size limit
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req ChatHTTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Input validation (before authorization to avoid unnecessary RPC)
	if len(req.Messages) == 0 {
		writeJSONError(w, "messages required", http.StatusBadRequest)
		return
	}
	if len(req.Messages) > maxMessages {
		writeJSONError(w, "too many messages", http.StatusBadRequest)
		return
	}
	if req.EnvironmentID == "" {
		writeJSONError(w, "environmentId required", http.StatusBadRequest)
		return
	}

	// Authorization check (Viewer role minimum)
	if h.accountClient == nil {
		h.logger.Error("accountClient is not configured")
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	resp, err := h.accountClient.GetAccountV2ByEnvironmentID(
		r.Context(),
		&accountproto.GetAccountV2ByEnvironmentIDRequest{
			Email:         accessToken.Email,
			EnvironmentId: req.EnvironmentID,
		},
	)
	if err != nil {
		h.logger.Error("authorization failed",
			zap.Error(err),
			zap.String("email", accessToken.Email),
			zap.String("environmentId", req.EnvironmentID),
		)
		writeJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}
	if resp.Account.Disabled {
		h.logger.Warn("disabled account attempted access",
			zap.String("email", accessToken.Email),
		)
		writeJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}
	envRole := getEnvironmentRole(resp.Account.EnvironmentRoles, req.EnvironmentID)
	if resp.Account.OrganizationRole < accountproto.AccountV2_Role_Organization_ADMIN &&
		envRole < accountproto.AccountV2_Role_Environment_VIEWER {
		h.logger.Warn("insufficient environment role",
			zap.String("email", accessToken.Email),
			zap.String("environmentId", req.EnvironmentID),
		)
		writeJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Convert to proto request
	protoReq := h.toProtoRequest(&req)

	// Check for SSE support before setting headers
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSONError(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	// Stream response using nil-channel pattern to avoid race conditions.
	// Both channels are closed when the ChatService goroutine finishes,
	// so we must drain responseChan fully before exiting.
	responseChan, errChan := streamChat(
		r.Context(), h.llmClient, h.ragSearcher, h.featureClient,
		h.chatConfig, protoReq, h.logger,
	)

	for responseChan != nil || errChan != nil {
		select {
		case <-r.Context().Done():
			return
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			if err != nil {
				h.writeSSEError(w, flusher, err)
				h.writeSSEDone(w, flusher)
				return
			}
		case resp, ok := <-responseChan:
			if !ok {
				responseChan = nil
				continue
			}
			h.writeSSEData(w, flusher, resp)
			if resp.Done {
				h.writeSSEDone(w, flusher)
				return
			}
		}
	}
	h.writeSSEDone(w, flusher)
}

func (h *chatHTTPService) writeSSEData(w http.ResponseWriter, f http.Flusher, resp *aichatproto.ChatStreamResponse) {
	data, err := json.Marshal(sseDataEvent{
		Content:      resp.Content,
		Done:         resp.Done,
		FinishReason: resp.FinishReason,
	})
	if err != nil {
		h.logger.Error("failed to marshal SSE data", zap.Error(err))
		return
	}
	fmt.Fprintf(w, "data: %s\n\n", data)
	f.Flush()
}

func (h *chatHTTPService) writeSSEDone(w http.ResponseWriter, f http.Flusher) {
	fmt.Fprintf(w, "data: [DONE]\n\n")
	f.Flush()
}

func (h *chatHTTPService) writeSSEError(w http.ResponseWriter, f http.Flusher, streamErr error) {
	h.logger.Error("SSE stream error", zap.Error(streamErr))
	data, err := json.Marshal(sseErrorEvent{
		Error: "An error occurred while processing your request",
	})
	if err != nil {
		h.logger.Error("failed to marshal SSE error", zap.Error(err))
		return
	}
	fmt.Fprintf(w, "event: error\ndata: %s\n\n", data)
	f.Flush()
}

// getEnvironmentRole returns the role for a given environment ID from the account's roles.
func getEnvironmentRole(
	roles []*accountproto.AccountV2_EnvironmentRole,
	envID string,
) accountproto.AccountV2_Role_Environment {
	for _, r := range roles {
		if r.EnvironmentId == envID {
			return r.Role
		}
	}
	return accountproto.AccountV2_Role_Environment_UNASSIGNED
}

func (h *chatHTTPService) toProtoRequest(req *ChatHTTPRequest) *aichatproto.ChatRequest {
	messages := make([]*aichatproto.ChatMessage, len(req.Messages))
	for i, m := range req.Messages {
		// Default to ROLE_USER for any unrecognized role string.
		r := aichatproto.ChatMessage_ROLE_USER
		if m.Role == llm.RoleAssistant {
			r = aichatproto.ChatMessage_ROLE_ASSISTANT
		}
		messages[i] = &aichatproto.ChatMessage{
			Role:    r,
			Content: m.Content,
		}
	}

	var pageContext *aichatproto.PageContext
	if req.PageContext != nil {
		pageType := aichatproto.PageContext_PAGE_TYPE_UNSPECIFIED
		switch req.PageContext.PageType {
		case httpPageTypeFeatureFlags:
			pageType = aichatproto.PageContext_PAGE_TYPE_FEATURE_FLAGS
		case httpPageTypeTargeting:
			pageType = aichatproto.PageContext_PAGE_TYPE_TARGETING
		case httpPageTypeExperiments:
			pageType = aichatproto.PageContext_PAGE_TYPE_EXPERIMENTS
		case httpPageTypeSegments:
			pageType = aichatproto.PageContext_PAGE_TYPE_SEGMENTS
		case httpPageTypeAutoops:
			pageType = aichatproto.PageContext_PAGE_TYPE_AUTOOPS
		}
		pageContext = &aichatproto.PageContext{
			PageType:  pageType,
			FeatureId: req.PageContext.FeatureID,
			Metadata:  req.PageContext.Metadata,
		}
	}

	return &aichatproto.ChatRequest{
		Messages:      messages,
		PageContext:   pageContext,
		EnvironmentId: req.EnvironmentID,
	}
}
