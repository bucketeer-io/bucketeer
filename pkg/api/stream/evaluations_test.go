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

package stream

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rest"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

func TestValidateStreamEvaluationsRequest(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		req         *gatewayproto.StreamEvaluationsRequest
		expectedErr error
	}{
		{
			desc:        "tag missing",
			req:         &gatewayproto.StreamEvaluationsRequest{User: &userproto.User{Id: "u"}},
			expectedErr: errTagRequired,
		},
		{
			desc:        "user missing",
			req:         &gatewayproto.StreamEvaluationsRequest{Tag: "t"},
			expectedErr: errUserRequired,
		},
		{
			desc:        "user id missing",
			req:         &gatewayproto.StreamEvaluationsRequest{Tag: "t", User: &userproto.User{}},
			expectedErr: errUserIDRequired,
		},
		{
			desc:        "valid",
			req:         &gatewayproto.StreamEvaluationsRequest{Tag: "t", User: &userproto.User{Id: "u"}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expectedErr, validateStreamEvaluationsRequest(p.req))
		})
	}
}

func TestCheckStreamEvaluationsRequestBodySizeLimit(t *testing.T) {
	t.Parallel()
	// bodyOfSize builds a valid request body of exactly size bytes by padding a
	// user attribute, so the limit can be exercised at its boundary.
	const bodyTemplate = `{"tag":"tag-A","user":{"id":"u-1","data":{"pad":"%s"}}}`
	bodyOfSize := func(size int) []byte {
		pad := size - len(fmt.Sprintf(bodyTemplate, ""))
		return []byte(fmt.Sprintf(bodyTemplate, strings.Repeat("x", pad)))
	}
	typicalBody, err := protojson.Marshal(&gatewayproto.StreamEvaluationsRequest{
		Tag:  "tag-A",
		User: &userproto.User{Id: "u-1", Data: map[string]string{"plan": "pro"}},
	})
	require.NoError(t, err)

	patterns := []struct {
		desc                string
		body                []byte
		expectedErr         error
		expectAPIKeyChecked bool
	}{
		{
			desc:                "typical body is accepted",
			body:                typicalBody,
			expectAPIKeyChecked: true,
		},
		{
			desc:                "body at the limit is accepted",
			body:                bodyOfSize(maxRequestBodyBytes),
			expectAPIKeyChecked: true,
		},
		{
			desc:        "body over the limit is rejected before the api key check",
			body:        bodyOfSize(maxRequestBodyBytes + 1),
			expectedErr: errRequestBodyTooLarge,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			apiKeyChecked := false
			h := &EvaluationsHandler{
				checkRequest: func(_ context.Context, _ *http.Request) (*accountproto.EnvironmentAPIKey, error) {
					apiKeyChecked = true
					return &accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-1"},
					}, nil
				},
				logger: zap.NewNop(),
			}
			httpReq := httptest.NewRequest(http.MethodPost, "/stream_evaluations", bytes.NewReader(p.body))

			envAPIKey, req, err := h.checkStreamEvaluationsRequest(httptest.NewRecorder(), httpReq)

			assert.Equal(t, p.expectAPIKeyChecked, apiKeyChecked)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
				assert.Nil(t, envAPIKey)
				assert.Nil(t, req)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "env-1", envAPIKey.Environment.Id)
			// The whole body must have been read and decoded, padding included.
			expected := &gatewayproto.StreamEvaluationsRequest{}
			require.NoError(t, protojson.Unmarshal(p.body, expected))
			assert.True(t, proto.Equal(expected, req))
		})
	}
}

// countingReader records how many bytes were read from the underlying reader.
type countingReader struct {
	io.Reader
	n int
}

func (r *countingReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.n += n
	return n, err
}

// The REST logging middleware decodes the body of a failed request after the
// handler returns. A request rejected before its body is read, such as one
// with the wrong method, must therefore already have the capped body in place.
func TestHandleCapsRequestBodyOfRejectedMethod(t *testing.T) {
	t.Parallel()
	h := &EvaluationsHandler{logger: zap.NewNop()}
	handler := rest.LogServerMiddleware(zap.NewNop())(http.HandlerFunc(h.Handle))
	// An unterminated JSON string keeps a decoder reading until the body ends.
	body := &countingReader{Reader: io.MultiReader(
		strings.NewReader(`"`),
		bytes.NewReader(bytes.Repeat([]byte("x"), 3*maxRequestBodyBytes)),
	)}
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/stream_evaluations", body))

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	assert.LessOrEqual(t, body.n, maxRequestBodyBytes+1,
		"nothing may read more than the cap (plus the byte that detects the overflow) from the connection")
}

type stubFlusher struct{}

func (stubFlusher) Flush() {}

func TestSendSSEEvent(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	msg := &gatewayproto.StreamEvaluationsEvent{
		Evaluations: &featureproto.UserEvaluations{
			Id:          "test",
			Evaluations: []*featureproto.Evaluation{},
		},
	}
	err := sendSSEEvent(&buf, stubFlusher{}, "put", msg)
	require.NoError(t, err)
	data, err := sseMarshalOpts.Marshal(msg)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("event: put\ndata: %s\n\n", data), buf.String())
}

func TestSendHeartbeat(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	err := sendHeartbeat(&buf, stubFlusher{})
	require.NoError(t, err)
	assert.Equal(t, ":\n\n", buf.String())
}

func TestSendErrorEvent(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	sendErrorEvent(&buf, stubFlusher{}, gatewayproto.StreamErrorEvent_INTERNAL, "something broke")
	data, err := sseMarshalOpts.Marshal(&gatewayproto.StreamErrorEvent{
		Code:    gatewayproto.StreamErrorEvent_INTERNAL,
		Message: "something broke",
	})
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("event: error\ndata: %s\n\n", data), buf.String())
}

func TestInitialPut(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc      string
		evals     *featureproto.UserEvaluations
		evalErr   error
		expectErr bool
	}{
		{
			desc: "success",
			evals: &featureproto.UserEvaluations{
				Id:          "eval-1",
				Evaluations: []*featureproto.Evaluation{},
			},
		},
		{
			desc:      "evaluate error",
			evalErr:   errors.New("eval failed"),
			expectErr: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			h := &EvaluationsHandler{
				evaluate: func(
					_ context.Context, _ *userproto.User, _, _ string,
					prevUEID string, evaluatedAt int64, checkUserAttributes bool,
					features []*featureproto.Feature,
				) (string, *featureproto.UserEvaluations, error) {
					assert.Nil(t, features)
					assert.Equal(t, "", prevUEID)
					assert.Equal(t, int64(0), evaluatedAt)
					assert.True(t, checkUserAttributes)
					return "ueid-1", p.evals, p.evalErr
				},
			}
			var buf bytes.Buffer
			ueid, evalAt, err := h.sendInitialPut(context.Background(), &buf, stubFlusher{}, &userproto.User{Id: "u"}, "env1", "tag1", "source1", "", 0)
			if p.expectErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "ueid-1", ueid)
			assert.Greater(t, evalAt, int64(0))
			lines := strings.SplitN(buf.String(), "\n", 3)
			assert.Equal(t, "event: put", lines[0])
			var got gatewayproto.StreamEvaluationsEvent
			require.NoError(t, sseUnmarshalOpts.Unmarshal([]byte(strings.TrimPrefix(lines[1], "data: ")), &got))
			assert.True(t, proto.Equal(p.evals, got.Evaluations))
			assert.Equal(t, "ueid-1", got.GetUserEvaluationsId())
		})
	}
}

func TestPatch(t *testing.T) {
	t.Parallel()
	snapshot := []*featureproto.Feature{{Id: "f1", Version: 3}}
	patterns := []struct {
		desc        string
		evaluatedAt int64
		evals       *featureproto.UserEvaluations
		evalErr     error
		expectErr   bool
		expectSend  bool
	}{
		{
			desc:        "diff: sends patch event",
			evaluatedAt: 1000,
			evals: &featureproto.UserEvaluations{
				Id: "eval-2",
				Evaluations: []*featureproto.Evaluation{
					{FeatureId: "f1"},
				},
			},
			expectSend: true,
		},
		{
			desc:        "none: skips sending",
			evaluatedAt: 1000,
			evals: &featureproto.UserEvaluations{
				Id:          "eval-2",
				Evaluations: []*featureproto.Evaluation{},
			},
			expectSend: false,
		},
		{
			desc:        "empty force update is delivered",
			evaluatedAt: 1000,
			evals: &featureproto.UserEvaluations{
				Id:          "eval-2",
				Evaluations: []*featureproto.Evaluation{},
				ForceUpdate: true,
			},
			expectSend: true,
		},
		{
			desc:        "evaluate error",
			evaluatedAt: 1000,
			evalErr:     errors.New("eval failed"),
			expectErr:   true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			h := &EvaluationsHandler{
				evaluate: func(
					_ context.Context, _ *userproto.User, _, _ string,
					prevUEID string, evaluatedAt int64, checkUserAttributes bool,
					features []*featureproto.Feature,
				) (string, *featureproto.UserEvaluations, error) {
					assert.Equal(t, snapshot, features)
					assert.Equal(t, "prev-ueid", prevUEID)
					assert.Equal(t, p.evaluatedAt, evaluatedAt)
					assert.False(t, checkUserAttributes)
					return "new-ueid", p.evals, p.evalErr
				},
			}
			var buf bytes.Buffer
			ueid, newEvalAt, err := h.sendPatch(context.Background(), &buf, stubFlusher{}, &userproto.User{Id: "u"}, "env1", "tag1", "source1", "prev-ueid", p.evaluatedAt, snapshot)
			if p.expectErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "new-ueid", ueid)
			assert.Greater(t, newEvalAt, p.evaluatedAt)
			if !p.expectSend {
				assert.Empty(t, buf.String())
				return
			}
			lines := strings.SplitN(buf.String(), "\n", 3)
			assert.Equal(t, "event: patch", lines[0])
			var got gatewayproto.StreamEvaluationsEvent
			require.NoError(t, sseUnmarshalOpts.Unmarshal([]byte(strings.TrimPrefix(lines[1], "data: ")), &got))
			assert.True(t, proto.Equal(p.evals, got.Evaluations))
			assert.Equal(t, "new-ueid", got.GetUserEvaluationsId())
		})
	}
}

// Handle blocks in its event loop while the SSE connection is active.
// Shutdown must unblock it so httpServer.Shutdown is not delayed.
func TestHandleExitsOnDispatcherShutdown(t *testing.T) {
	t.Parallel()
	d := NewDispatcher(10000, nil, zap.NewNop())
	handleReturned := startBlockedStreamHandle(t, d)

	d.Shutdown()

	select {
	case <-handleReturned:
		// Success: Shutdown unblocked the handler.
	case <-time.After(time.Second):
		t.Fatal("Handle is still blocked in its event loop after Dispatcher.Shutdown")
	}
}

func TestHandleClosesStreamWhenHoldoutSnapshotUnavailable(t *testing.T) {
	t.Parallel()
	d := NewDispatcher(10000, nil, zap.NewNop())
	// The gateway marks IDs issued without a holdout snapshot. A stream must not
	// keep such assignments alive until an unrelated domain event: it has to
	// close so the reconnect forces a full evaluation that restores HOLDOUT.
	markedUEID := "unavailable-holdout:ueid-1"
	require.True(t, evaluation.IssuedWithoutHoldout(markedUEID))
	rec := httptest.NewRecorder()
	handleReturned := runStreamHandle(t, d, rec, fixedEvaluate(markedUEID))

	select {
	case <-handleReturned:
	case <-time.After(time.Second):
		t.Fatal("Handle kept the stream open after an evaluation without a holdout snapshot")
	}
	body := rec.Body.String()
	assert.Contains(t, body, "event: put\n", "the fresh assignments must still be delivered")
	assert.Contains(t, body, markedUEID)
	assert.Contains(t, body, "event: error\n", "the stream must close with a retryable error")
	assert.Contains(t, body, "holdout snapshot unavailable")
	assert.Empty(t, snapshotConnCounts(d), "the connection must be deregistered")
}

func TestHandleKeepsStreamOpenWithHoldoutSnapshot(t *testing.T) {
	t.Parallel()
	d := NewDispatcher(10000, nil, zap.NewNop())
	rec := httptest.NewRecorder()
	handleReturned := runStreamHandle(t, d, rec, fixedEvaluate("ueid-1"))

	select {
	case <-handleReturned:
		t.Fatal("Handle closed a stream whose evaluation had the holdout snapshot")
	case <-time.After(100 * time.Millisecond):
	}
	d.Shutdown()
	<-handleReturned
	assert.NotContains(t, rec.Body.String(), "event: error\n")
}

// writeDeadlineRecorder stands in for the server connection: it records the
// write deadlines set through http.ResponseController.
type writeDeadlineRecorder struct {
	*httptest.ResponseRecorder
	writeDeadlines []time.Time
}

func (r *writeDeadlineRecorder) SetWriteDeadline(deadline time.Time) error {
	r.writeDeadlines = append(r.writeDeadlines, deadline)
	return nil
}

// The REST server applies a write timeout to every response. A stream stays
// open far longer than that, so Handle must clear the deadline of its response.
func TestHandleClearsWriteDeadlineForStream(t *testing.T) {
	t.Parallel()
	d := NewDispatcher(10000, nil, zap.NewNop())
	rec := &writeDeadlineRecorder{ResponseRecorder: httptest.NewRecorder()}
	handleReturned := runStreamHandle(t, d, rec, fixedEvaluate("ueid-1"))

	d.Shutdown()
	<-handleReturned

	assert.Equal(t, []time.Time{{}}, rec.writeDeadlines, "the write deadline must be cleared exactly once")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "event: put\n", "the stream must still be served")
}

func fixedEvaluate(ueid string) EvaluateFunc {
	return func(
		_ context.Context, _ *userproto.User, _, _ string, _ string, _ int64, _ bool,
		_ []*featureproto.Feature,
	) (string, *featureproto.UserEvaluations, error) {
		return ueid, &featureproto.UserEvaluations{
			Id:          ueid,
			Evaluations: []*featureproto.Evaluation{},
		}, nil
	}
}

// startBlockedStreamHandle runs Handle for one SSE connection in a goroutine and waits
// until the connection is registered. The returned channel is closed when Handle returns.
func startBlockedStreamHandle(t *testing.T, d *Dispatcher) <-chan struct{} {
	t.Helper()
	return runStreamHandle(t, d, httptest.NewRecorder(), fixedEvaluate("ueid-1"))
}

// runStreamHandle runs Handle for one SSE connection in a goroutine and waits
// until the connection is registered or Handle returns. The returned channel
// is closed when Handle returns.
func runStreamHandle(
	t *testing.T, d *Dispatcher, rec http.ResponseWriter, evaluate EvaluateFunc,
) <-chan struct{} {
	t.Helper()
	h := NewEvaluationsHandler(
		d,
		time.Hour, // ensure the heartbeat ticker never fires during the test
		evaluate,
		func(_ context.Context, _ *http.Request) (*accountproto.EnvironmentAPIKey, error) {
			return &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "env-1"},
			}, nil
		},
		prometheus.NewCounterVec(prometheus.CounterOpts{Name: "test_request_total"},
			[]string{"organization_id", "project_id", "project_url_code",
				"environment_id", "environment_url_code", "method", "source_id"}),
		zap.NewNop(),
	)

	body, err := protojson.Marshal(&gatewayproto.StreamEvaluationsRequest{
		Tag:  "tag-A",
		User: &userproto.User{Id: "u-1"},
	})
	require.NoError(t, err)
	httpReq := httptest.NewRequest(http.MethodPost, "/stream_evaluations", bytes.NewReader(body))

	handleReturned := make(chan struct{})
	go func() {
		defer close(handleReturned)
		h.Handle(rec, httpReq)
	}()

	require.Eventually(t, func() bool {
		select {
		case <-handleReturned:
			return true
		default:
			return len(snapshotConnCounts(d)) > 0
		}
	}, time.Second, 10*time.Millisecond)

	return handleReturned
}
