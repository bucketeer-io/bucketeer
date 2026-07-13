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

package gateway

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

const (
	streamEvaluationsPath = "/v1/gateway/stream_evaluations"

	sseEventPut       = "put"
	sseEventPatch     = "patch"
	sseEventError     = "error"
	sseEventHeartbeat = "heartbeat"

	// Set longer timeouts for put and patch to prevent flaky tests.
	initialPutTimeout = 30 * time.Second
	patchTimeout      = 60 * time.Second
	// The server heartbeat interval defaults to 25s (charts/api/values.yaml).
	heartbeatTimeout = 40 * time.Second
	// Includes the async bulk upload processing by the subscriber service.
	segmentPatchTimeout = 90 * time.Second
	noPatchWindow       = 10 * time.Second
)

type sseEvent struct {
	event string
	data  string
}

type sseStream struct {
	resp   *http.Response
	cancel context.CancelFunc
	events chan sseEvent
	errs   chan error
}

func TestStreamInitialPut(t *testing.T) {
	t.Parallel()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)

	userID := newUserID(t, uuid)
	// Connect and receive PUT
	stream := connectStream(t, newStreamBody(tag, userID))
	put := stream.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
	if put.UserEvaluationsId == "" {
		t.Fatal("UserEvaluationsId is empty")
	}
	if !put.Evaluations.ForceUpdate {
		t.Fatal("Wrong forceUpdate. Expected true on initial connect, actual false")
	}
	eval, err := findFeature(put.Evaluations.Evaluations, featureID)
	if err != nil {
		t.Fatalf("Failed to find evaluation in initial put. Error: %v", err)
	}
	if eval.VariationValue != "A" {
		t.Fatalf("Wrong variation value. Expected: A, actual: %s", eval.VariationValue)
	}
}

func TestStreamHeartbeat(t *testing.T) {
	t.Parallel()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)

	stream := connectStream(t, newStreamBody(tag, userID))
	stream.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
	stream.waitForEvent(t, sseEventHeartbeat, heartbeatTimeout, nil)
}

func TestStreamPatchOnFlagChange(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	featureID := newFeatureID(t, uuid)
	req := newCreateFeatureReq(featureID)
	createFeature(t, client, req)
	addTag(t, tag, featureID, client)
	updateFeatueFlagCache(t)

	userID := newUserID(t, uuid)
	// Connect and receive PUT
	stream := connectStream(t, newStreamBody(tag, userID))
	put := stream.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
	eval, err := findFeature(put.Evaluations.Evaluations, featureID)
	if err != nil {
		t.Fatalf("Failed to find evaluation in initial put. Error: %v", err)
	}
	if eval.VariationValue != "B" {
		t.Fatalf("Wrong variation value. Expected: B (off variation), actual: %s", eval.VariationValue)
	}

	// Update feature and receive PATCH
	enableFeature(t, featureID, client)
	patch := stream.waitForEvent(t, sseEventPatch, patchTimeout, func(evt *gatewayproto.StreamEvaluationsEvent) bool {
		e, err := findFeature(evt.Evaluations.Evaluations, featureID)
		return err == nil && e.VariationValue == "A"
	})
	if patch.Evaluations.ForceUpdate {
		t.Fatal("Wrong forceUpdate. Expected false on patch, actual true")
	}
	if patch.UserEvaluationsId == put.UserEvaluationsId {
		t.Fatal("UserEvaluationsId did not change after the flag update")
	}
}

func TestStreamReconnect(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	featureID := newFeatureID(t, uuid)
	req := newCreateFeatureReq(featureID)
	createFeature(t, client, req)
	addTag(t, tag, featureID, client)
	updateFeatueFlagCache(t)

	userID := newUserID(t, uuid)
	// Connect and receive PUT
	stream := connectStream(t, newStreamBody(tag, userID))
	put := stream.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
	ueid := put.UserEvaluationsId
	evaluatedAt := time.Now().Unix()
	// Close stream
	stream.close()

	// No flag changes since disconnect; server returns a diff put (forceUpdate=false).
	t.Run("UnchangedReconnectReturnsDiffPut", func(t *testing.T) {
		s := connectStream(t, newStreamResumeBody(tag, userID, ueid, evaluatedAt))
		resumed := s.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
		if resumed.Evaluations.ForceUpdate {
			t.Fatal("Wrong forceUpdate. Expected false on resumed connect, actual true")
		}
		s.close()
	})

	// Stale evaluatedAt (>30 days) forces a full-state put (forceUpdate=true).
	t.Run("EvaluatedAtOlderThan30DaysReturnsFullState", func(t *testing.T) {
		old := time.Now().Unix() - 31*24*60*60
		s := connectStream(t, newStreamResumeBody(tag, userID, ueid, old))
		resumed := s.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
		if !resumed.Evaluations.ForceUpdate {
			t.Fatal("Wrong forceUpdate. Expected true for evaluatedAt older than 30 days, actual false")
		}
		if !contains(resumed.Evaluations.Evaluations, featureID) {
			t.Fatalf("Full state put does not contain feature: %s", featureID)
		}
		s.close()
	})

	// Flag changed while disconnected; diff put includes the updated evaluation.
	t.Run("ReconnectAfterFlagChangeContainsUpdatedEvaluation", func(t *testing.T) {
		enableFeature(t, featureID, client)
		// Reconnect is not event-driven, so refresh the flag cache to make
		// the updated flag visible deterministically.
		updateFeatueFlagCache(t)
		s := connectStream(t, newStreamResumeBody(tag, userID, ueid, evaluatedAt))
		resumed := s.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
		if resumed.Evaluations.ForceUpdate {
			t.Fatal("Wrong forceUpdate. Expected false on resumed connect, actual true")
		}
		eval, err := findFeature(resumed.Evaluations.Evaluations, featureID)
		if err != nil {
			t.Fatalf("Failed to find evaluation in resumed put. Error: %v", err)
		}
		if eval.VariationValue != "A" {
			t.Fatalf("Wrong variation value. Expected: A, actual: %s", eval.VariationValue)
		}
		s.close()
	})
}

func TestStreamTagFilter(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	// Create features with different tags
	tagA := fmt.Sprintf("%s-tag-a-%s", prefixTestName, uuid)
	tagB := fmt.Sprintf("%s-tag-b-%s", prefixTestName, uuid)
	featureIDA := newFeatureID(t, uuid)
	featureIDB := newFeatureID(t, newUUID(t))
	reqA := newCreateFeatureReq(featureIDA)
	createFeature(t, client, reqA)
	addTag(t, tagA, featureIDA, client)
	reqB := newCreateFeatureReq(featureIDB)
	createFeature(t, client, reqB)
	addTag(t, tagB, featureIDB, client)
	updateFeatueFlagCache(t)

	userID := newUserID(t, uuid)
	// Client A - Connect and receive PUT for tag-A
	streamA := connectStream(t, newStreamBody(tagA, userID))
	putA := streamA.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
	if !contains(putA.Evaluations.Evaluations, featureIDA) {
		t.Fatalf("Initial put for tag-A does not contain feature: %s", featureIDA)
	}
	// Client B - Connect and receive PUT for tag-B
	streamB := connectStream(t, newStreamBody(tagB, userID))
	putB := streamB.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
	if contains(putB.Evaluations.Evaluations, featureIDA) {
		t.Fatalf("Initial put for tag-B contains tag-A feature: %s", featureIDA)
	}

	// Update feature of A
	enableFeature(t, featureIDA, client)
	streamA.waitForEvent(t, sseEventPatch, patchTimeout, func(evt *gatewayproto.StreamEvaluationsEvent) bool {
		e, err := findFeature(evt.Evaluations.Evaluations, featureIDA)
		return err == nil && e.VariationValue == "A"
	})
	// Env-wide dispatches from parallel tests may deliver unrelated patches to
	// tag-B, so assert on the content instead of the absence of any patch.
	streamB.assertNoPatchContaining(t, featureIDA, noPatchWindow)
}

func TestStreamPatchOnArchive(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)

	userID := newUserID(t, uuid)
	// Connect and receive PUT
	stream := connectStream(t, newStreamBody(tag, userID))
	put := stream.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
	if !contains(put.Evaluations.Evaluations, featureID) {
		t.Fatalf("Initial put does not contain feature: %s", featureID)
	}

	// Check if the archived feature is included in the archivedFeatureIds of patch response.
	archiveFeature(t, featureID, client)
	patch := stream.waitForEvent(t, sseEventPatch, patchTimeout, func(evt *gatewayproto.StreamEvaluationsEvent) bool {
		return slices.Contains(evt.Evaluations.ArchivedFeatureIds, featureID)
	})
	if patch.Evaluations.ForceUpdate {
		t.Fatal("Wrong forceUpdate. Expected false on patch, actual true")
	}
	if !slices.Contains(patch.Evaluations.ArchivedFeatureIds, featureID) {
		t.Fatalf("ArchivedFeatureIds does not contain feature: %s", featureID)
	}
}

func TestStreamPatchOnSegmentBulkUpload(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	featureID := newFeatureID(t, uuid)
	req := newCreateFeatureReq(featureID)
	createFeature(t, client, req)
	addTag(t, tag, featureID, client)
	feature := getFeature(t, featureID, client)
	segmentID := createSegment(t, client)
	addSegmentRuleToFeature(t, featureID, feature.Variations[1].Id, segmentID, client) // return B for this segment
	enableFeature(t, featureID, client)
	updateFeatueFlagCache(t)

	userID := newUserID(t, uuid)
	// Connect and receive PUT
	stream := connectStream(t, newStreamBody(tag, userID))
	put := stream.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
	eval, err := findFeature(put.Evaluations.Evaluations, featureID)
	if err != nil {
		t.Fatalf("Failed to find evaluation in initial put. Error: %v", err)
	}
	if eval.VariationValue != "A" {
		t.Fatalf("Wrong variation value. Expected: A (not in segment), actual: %s", eval.VariationValue)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = client.BulkUploadSegmentUsers(ctx, &featureproto.BulkUploadSegmentUsersRequest{
		EnvironmentId: *environmentID,
		SegmentId:     segmentID,
		Data:          []byte(userID + "\n"),
		State:         featureproto.SegmentUser_INCLUDED,
	})
	if err != nil {
		t.Fatalf("Failed to bulk upload segment users. Error: %v", err)
	}
	// After updating the segment users, the patch should be sent.
	patch := stream.waitForEvent(t, sseEventPatch, segmentPatchTimeout, func(evt *gatewayproto.StreamEvaluationsEvent) bool {
		e, err := findFeature(evt.Evaluations.Evaluations, featureID)
		return err == nil && e.VariationValue == "B"
	})
	if patch.Evaluations.ForceUpdate {
		t.Fatal("Wrong forceUpdate. Expected false on patch, actual true")
	}
	patchEval, err := findFeature(patch.Evaluations.Evaluations, featureID)
	if err != nil {
		t.Fatalf("Failed to find evaluation in segment patch. Error: %v", err)
	}
	if patchEval.VariationValue != "B" {
		t.Fatalf("Wrong variation value. Expected: B (in segment), actual: %s", patchEval.VariationValue)
	}
}

func TestStreamRequestErrors(t *testing.T) {
	t.Parallel()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	clientKey := readAPIKey(t, *apiKeyPath)
	serverKey := readAPIKey(t, *apiKeyServerPath)
	invalidKey := readAPIKey(t, "testdata/invalid-apikey")

	testcases := []struct {
		desc           string
		method         string
		body           map[string]any
		apiKey         string
		expectedStatus int
	}{
		{
			desc:           "missing tag",
			method:         http.MethodPost,
			body:           map[string]any{"user": map[string]any{"id": userID}},
			apiKey:         clientKey,
			expectedStatus: http.StatusBadRequest,
		},
		{
			desc:           "missing user",
			method:         http.MethodPost,
			body:           map[string]any{"tag": tag},
			apiKey:         clientKey,
			expectedStatus: http.StatusBadRequest,
		},
		{
			desc:           "missing user id",
			method:         http.MethodPost,
			body:           map[string]any{"tag": tag, "user": map[string]any{}},
			apiKey:         clientKey,
			expectedStatus: http.StatusBadRequest,
		},
		{
			desc:           "missing api key",
			method:         http.MethodPost,
			body:           newStreamBody(tag, userID),
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			desc:           "invalid api key",
			method:         http.MethodPost,
			body:           newStreamBody(tag, userID),
			apiKey:         invalidKey,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			desc:           "server sdk api key",
			method:         http.MethodPost,
			body:           newStreamBody(tag, userID),
			apiKey:         serverKey,
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			resp, cancel := streamRequest(t, tc.method, tc.body, tc.apiKey)
			defer cancel()
			defer resp.Body.Close()
			if resp.StatusCode != tc.expectedStatus {
				body := new(bytes.Buffer)
				_, _ = body.ReadFrom(resp.Body)
				t.Fatalf("Wrong status code. Expected: %d, actual: %d, body: %s",
					tc.expectedStatus, resp.StatusCode, body.String())
			}
		})
	}
}

func newStreamBody(tag, userID string) map[string]any {
	return map[string]any{
		"tag":  tag,
		"user": map[string]any{"id": userID},
	}
}

func newStreamResumeBody(tag, userID, ueid string, evaluatedAt int64) map[string]any {
	body := newStreamBody(tag, userID)
	body["user_evaluations_id"] = ueid
	body["evaluated_at"] = evaluatedAt
	return body
}

func readAPIKey(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(data))
}

func streamRequest(
	t *testing.T,
	method string,
	body map[string]any,
	apiKey string,
) (*http.Response, context.CancelFunc) {
	t.Helper()
	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	url := fmt.Sprintf("https://%s%s", *gatewayAddr, streamEvaluationsPath)
	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(encoded))
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	if apiKey != "" {
		req.Header.Add("authorization", apiKey)
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
			ResponseHeaderTimeout: initialPutTimeout,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	return resp, cancel
}

func connectStream(t *testing.T, body map[string]any) *sseStream {
	t.Helper()
	resp, cancel := streamRequest(t, http.MethodPost, body, readAPIKey(t, *apiKeyPath))
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
	s := &sseStream{
		resp:   resp,
		cancel: cancel,
		events: make(chan sseEvent, 64),
		errs:   make(chan error, 1),
	}
	go s.readLoop()
	t.Cleanup(s.close)
	return s
}

func (s *sseStream) readLoop() {
	defer close(s.events)
	scanner := bufio.NewScanner(s.resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var eventType, data string
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == "":
			if eventType != "" || data != "" {
				s.events <- sseEvent{event: eventType, data: data}
				eventType, data = "", ""
			}
		case strings.HasPrefix(line, ":"):
			s.events <- sseEvent{event: sseEventHeartbeat}
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

func (s *sseStream) close() {
	s.cancel()
	s.resp.Body.Close()
}

// waitForEvent returns the next event of the given type accepted by match,
// skipping heartbeats and non-matching events. Fails on error events or timeout.
func (s *sseStream) waitForEvent(
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
			if ev.event == sseEventError {
				t.Fatalf("Received error event: %s", ev.data)
			}
			if ev.event != eventType {
				continue
			}
			if eventType == sseEventHeartbeat {
				return nil
			}
			evt := &gatewayproto.StreamEvaluationsEvent{}
			if err := protojson.Unmarshal([]byte(ev.data), evt); err != nil {
				t.Fatalf("Failed to unmarshal %s event: %v, data: %s", eventType, err, ev.data)
			}
			if match == nil || match(evt) {
				return evt
			}
			t.Logf("Skipping non-matching %s event: %s", eventType, ev.data)
		case err := <-s.errs:
			t.Fatalf("Stream read error while waiting for %s event: %v", eventType, err)
		case <-deadline:
			t.Fatalf("Timed out waiting for %s event", eventType)
		}
	}
}

// assertNoPatchContaining fails if a patch containing featureID arrives within wait.
func (s *sseStream) assertNoPatchContaining(t *testing.T, featureID string, wait time.Duration) {
	t.Helper()
	deadline := time.After(wait)
	for {
		select {
		case ev, ok := <-s.events:
			if !ok {
				t.Fatal("Stream closed while asserting patch absence")
			}
			if ev.event != sseEventPatch {
				continue
			}
			evt := &gatewayproto.StreamEvaluationsEvent{}
			if err := protojson.Unmarshal([]byte(ev.data), evt); err != nil {
				t.Fatalf("Failed to unmarshal patch event: %v, data: %s", err, ev.data)
			}
			if contains(evt.Evaluations.Evaluations, featureID) {
				t.Fatalf("Received a patch containing feature: %s", featureID)
			}
		case err := <-s.errs:
			t.Fatalf("Stream read error while asserting patch absence: %v", err)
		case <-deadline:
			return
		}
	}
}

func createSegment(t *testing.T, client featureclient.Client) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.CreateSegment(ctx, &featureproto.CreateSegmentRequest{
		EnvironmentId: *environmentID,
		Name:          fmt.Sprintf("%s-segment-%s", prefixTestName, newUUID(t)),
		Description:   "e2e-test-stream-segment",
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.Segment.Id
}

func addSegmentRuleToFeature(t *testing.T, featureID, variationID, segmentID string, client featureclient.Client) {
	t.Helper()
	rule := &featureproto.Rule{
		Id: newUUID(t),
		Strategy: &featureproto.Strategy{
			Type: featureproto.Strategy_FIXED,
			FixedStrategy: &featureproto.FixedStrategy{
				Variation: variationID,
			},
		},
		Clauses: []*featureproto.Clause{
			{
				Id:       newUUID(t),
				Operator: featureproto.Clause_SEGMENT,
				Values:   []string{segmentID},
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
		RuleChanges: []*featureproto.RuleChange{
			{
				ChangeType: featureproto.ChangeType_CREATE,
				Rule:       rule,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
