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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

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
				evaluate: func(_ context.Context, _ *userproto.User, _, _ string, evaluatedAt int64) (*featureproto.UserEvaluations, error) {
					assert.Equal(t, int64(0), evaluatedAt)
					return p.evals, p.evalErr
				},
			}
			var buf bytes.Buffer
			evalAt, err := h.sendInitialPut(context.Background(), &buf, stubFlusher{}, &userproto.User{Id: "u"}, "env1", "tag1", "source1")
			if p.expectErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Greater(t, evalAt, int64(0))
			lines := strings.SplitN(buf.String(), "\n", 3)
			assert.Equal(t, "event: put", lines[0])
			var got gatewayproto.StreamEvaluationsEvent
			require.NoError(t, sseUnmarshalOpts.Unmarshal([]byte(strings.TrimPrefix(lines[1], "data: ")), &got))
			assert.True(t, proto.Equal(p.evals, got.Evaluations))
		})
	}
}

func TestPatch(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		evaluatedAt int64
		evals       *featureproto.UserEvaluations
		evalErr     error
		expectErr   bool
	}{
		{
			desc:        "success",
			evaluatedAt: 1000,
			evals: &featureproto.UserEvaluations{
				Id:          "eval-2",
				Evaluations: []*featureproto.Evaluation{},
			},
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
				evaluate: func(_ context.Context, _ *userproto.User, _, _ string, evaluatedAt int64) (*featureproto.UserEvaluations, error) {
					assert.Equal(t, p.evaluatedAt, evaluatedAt)
					return p.evals, p.evalErr
				},
			}
			var buf bytes.Buffer
			newEvalAt, err := h.sendPatch(context.Background(), &buf, stubFlusher{}, &userproto.User{Id: "u"}, "env1", "tag1", "source1", p.evaluatedAt)
			if p.expectErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Greater(t, newEvalAt, p.evaluatedAt)
			lines := strings.SplitN(buf.String(), "\n", 3)
			assert.Equal(t, "event: patch", lines[0])
			var got gatewayproto.StreamEvaluationsEvent
			require.NoError(t, sseUnmarshalOpts.Unmarshal([]byte(strings.TrimPrefix(lines[1], "data: ")), &got))
			assert.True(t, proto.Equal(p.evals, got.Evaluations))
		})
	}
}
