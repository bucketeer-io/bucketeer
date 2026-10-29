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

package util

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

// SSE event names emitted by /v1/gateway/stream_evaluations. Heartbeats are
// SSE comment lines and surface as SSEEventHeartbeat.
const (
	SSEEventPut       = "put"
	SSEEventPatch     = "patch"
	SSEEventError     = "error"
	SSEEventHeartbeat = "heartbeat"

	streamEvaluationsPath = "/v1/gateway/stream_evaluations"
	streamConnectTimeout  = 30 * time.Second
)

// SSEEvent is one raw server-sent event.
type SSEEvent struct {
	Event string
	Data  string
}

// SSEStream is an open /stream_evaluations connection. It mirrors the
// unexported helper in test/e2e/gateway/stream_test.go so other suites can
// drive the SSE path without copying it.
type SSEStream struct {
	resp   *http.Response
	cancel context.CancelFunc
	events chan SSEEvent
	errs   chan error
	// Errors keeps every error event received so a test can assert that the
	// server never closed the stream with one.
	Errors []*gatewayproto.StreamErrorEvent
}

// StreamBody builds the request body for a fresh stream.
func StreamBody(tag, userID string) map[string]interface{} {
	return map[string]interface{}{
		"tag":  tag,
		"user": map[string]interface{}{"id": userID},
	}
}

// StreamResumeBody builds the request body for a reconnect carrying the
// previous UserEvaluationsID and evaluatedAt.
func StreamResumeBody(tag, userID, ueid string, evaluatedAt int64) map[string]interface{} {
	body := StreamBody(tag, userID)
	body["user_evaluations_id"] = ueid
	body["evaluated_at"] = evaluatedAt
	return body
}

// ConnectStream opens an SSE stream on gatewayAddr with the API key stored at
// apiKeyPath and starts reading it. The stream is closed in t.Cleanup.
func ConnectStream(t *testing.T, gatewayAddr, apiKeyPath string, body map[string]interface{}) *SSEStream {
	t.Helper()
	key, err := os.ReadFile(apiKeyPath)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	url := fmt.Sprintf("https://%s%s", gatewayAddr, streamEvaluationsPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(encoded))
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	req.Header.Add(authorizationKey, strings.TrimSpace(string(key)))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
			ResponseHeaderTimeout: streamConnectTimeout,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		resp.Body.Close()
		cancel()
		t.Fatalf("Failed to connect stream. Status: %d, body: %s", resp.StatusCode, buf.String())
	}
	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		resp.Body.Close()
		cancel()
		t.Fatalf("Wrong content type. Expected: text/event-stream, actual: %s", ct)
	}
	s := &SSEStream{
		resp:   resp,
		cancel: cancel,
		events: make(chan SSEEvent, 64),
		errs:   make(chan error, 1),
	}
	go s.readLoop()
	t.Cleanup(s.Close)
	return s
}

func (s *SSEStream) readLoop() {
	defer close(s.events)
	scanner := bufio.NewScanner(s.resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var eventType, data string
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == "":
			if eventType != "" || data != "" {
				s.events <- SSEEvent{Event: eventType, Data: data}
				eventType, data = "", ""
			}
		case strings.HasPrefix(line, ":"):
			s.events <- SSEEvent{Event: SSEEventHeartbeat}
		case strings.HasPrefix(line, "event: "):
			eventType = strings.TrimPrefix(line, "event: ")
		case strings.HasPrefix(line, "data: "):
			data = strings.TrimPrefix(line, "data: ")
		}
	}
	if err := scanner.Err(); err != nil {
		select {
		case s.errs <- err:
		default:
		}
	}
}

// Close ends the stream.
func (s *SSEStream) Close() {
	s.cancel()
	s.resp.Body.Close()
}

// WaitForEvent returns the next event of eventType accepted by match (nil
// accepts any), skipping heartbeats and non-matching events. It fails on an
// error event, a closed stream, or the timeout.
func (s *SSEStream) WaitForEvent(
	t *testing.T,
	eventType string,
	timeout time.Duration,
	match func(*gatewayproto.StreamEvaluationsEvent) bool,
) *gatewayproto.StreamEvaluationsEvent {
	t.Helper()
	deadline := time.After(timeout)
	for {
		select {
		case ev, ok := <-s.events:
			if !ok {
				t.Fatalf("Stream closed while waiting for %s event", eventType)
			}
			if ev.Event == SSEEventError {
				s.recordError(t, ev)
				t.Fatalf("Received error event: %s", ev.Data)
			}
			if ev.Event != eventType {
				continue
			}
			if eventType == SSEEventHeartbeat {
				return nil
			}
			evt := &gatewayproto.StreamEvaluationsEvent{}
			if err := protojson.Unmarshal([]byte(ev.Data), evt); err != nil {
				t.Fatalf("Failed to unmarshal %s event: %v, data: %s", eventType, err, ev.Data)
			}
			if match == nil || match(evt) {
				return evt
			}
			t.Logf("Skipping non-matching %s event: %s", eventType, ev.Data)
		case err := <-s.errs:
			t.Fatalf("Stream read error while waiting for %s event: %v", eventType, err)
		case <-deadline:
			t.Fatalf("Timed out waiting for %s event", eventType)
		}
	}
}

// DrainPatches reads events for wait and returns every patch received. It
// fails on an error event or a closed stream, so a caller can assert both
// "no patch matched X" and "the stream stayed open".
func (s *SSEStream) DrainPatches(t *testing.T, wait time.Duration) []*gatewayproto.StreamEvaluationsEvent {
	t.Helper()
	deadline := time.After(wait)
	var patches []*gatewayproto.StreamEvaluationsEvent
	for {
		select {
		case ev, ok := <-s.events:
			if !ok {
				t.Fatal("Stream closed while draining patches")
			}
			if ev.Event == SSEEventError {
				s.recordError(t, ev)
				t.Fatalf("Received error event while draining patches: %s", ev.Data)
			}
			if ev.Event != SSEEventPatch {
				continue
			}
			evt := &gatewayproto.StreamEvaluationsEvent{}
			if err := protojson.Unmarshal([]byte(ev.Data), evt); err != nil {
				t.Fatalf("Failed to unmarshal patch event: %v, data: %s", err, ev.Data)
			}
			patches = append(patches, evt)
		case err := <-s.errs:
			t.Fatalf("Stream read error while draining patches: %v", err)
		case <-deadline:
			return patches
		}
	}
}

func (s *SSEStream) recordError(t *testing.T, ev SSEEvent) {
	t.Helper()
	evt := &gatewayproto.StreamErrorEvent{}
	if err := protojson.Unmarshal([]byte(ev.Data), evt); err == nil {
		s.Errors = append(s.Errors, evt)
	}
}
