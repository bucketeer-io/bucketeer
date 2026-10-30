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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/encoding/protojson"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	streamLogUserID = "user-1"
	// streamLogUserAttribute is an end-user attribute value that must never
	// be written to the log.
	streamLogUserAttribute = "user-1@example.com"
)

func TestHandleLogsUserIDOnly(t *testing.T) {
	t.Parallel()
	checkErr := errInternal
	patterns := []struct {
		desc           string
		checkRequest   CheckRequestFunc
		req            *gatewayproto.StreamEvaluationsRequest
		expectedStatus int
		expectedLogMsg string
	}{
		{
			desc: "err: request check failure is logged with the user id only",
			checkRequest: func(context.Context, *http.Request) (*accountproto.EnvironmentAPIKey, error) {
				return nil, checkErr
			},
			req: &gatewayproto.StreamEvaluationsRequest{
				Tag: "tag-A",
				User: &userproto.User{
					Id:   streamLogUserID,
					Data: map[string]string{"email": streamLogUserAttribute},
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedLogMsg: "Failed to check StreamEvaluations request",
		},
		{
			desc: "err: validation failure is logged with the user id only",
			checkRequest: func(context.Context, *http.Request) (*accountproto.EnvironmentAPIKey, error) {
				return &accountproto.EnvironmentAPIKey{
					Environment: &environmentproto.EnvironmentV2{Id: "env-1"},
				}, nil
			},
			req: &gatewayproto.StreamEvaluationsRequest{
				User: &userproto.User{
					Id:   streamLogUserID,
					Data: map[string]string{"email": streamLogUserAttribute},
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedLogMsg: "Failed to validate StreamEvaluations request",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			core, logs := observer.New(zapcore.DebugLevel)
			h := NewEvaluationsHandler(
				NewDispatcher(10000, nil, zap.NewNop()),
				time.Hour,
				fixedEvaluate("ueid-1"),
				p.checkRequest,
				prometheus.NewCounterVec(prometheus.CounterOpts{Name: "test_log_request_total"},
					[]string{"organization_id", "project_id", "project_url_code",
						"environment_id", "environment_url_code", "method", "source_id"}),
				zap.New(core),
			)
			body, err := protojson.Marshal(p.req)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			h.Handle(rec, httptest.NewRequest(http.MethodPost, "/stream_evaluations", bytes.NewReader(body)))

			assert.Equal(t, p.expectedStatus, rec.Code)
			for _, entry := range logs.All() {
				assert.NotContains(t, entry.Message, streamLogUserAttribute)
				for key, field := range entry.ContextMap() {
					assert.NotContains(t, fmt.Sprint(field), streamLogUserAttribute,
						"field %q of %q", key, entry.Message)
				}
			}
			entries := logs.FilterMessage(p.expectedLogMsg).All()
			require.Len(t, entries, 1, "expected exactly one %q entry", p.expectedLogMsg)
			assert.Equal(t, streamLogUserID, entries[0].ContextMap()["userID"])
		})
	}
}
