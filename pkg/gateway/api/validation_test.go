// Copyright 2024 The Bucketeer Authors.
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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/pkg/log"
)

func TestValidateGoalEvent(t *testing.T) {
	t.Parallel()
	logger, _ := log.NewLogger()
	now := time.Now().Unix()
	ctx := context.TODO()
	patterns := []struct {
		desc         string
		id           string
		timestamp    int64
		expectedCode string
		expectedErr  error
	}{
		{
			desc:         "err: invalid uuid",
			id:           "0efe416e 2fd2 4996 c5c3 194f05444f1f",
			timestamp:    now,
			expectedCode: codeInvalidID,
			expectedErr:  errInvalidIDFormat,
		},
		{
			desc:         "err: invalid timestamp",
			id:           newUUID(t),
			timestamp:    int64(999999999999999),
			expectedCode: codeInvalidTimestamp,
			expectedErr:  errInvalidTimestamp,
		},
		{
			desc:         "success",
			id:           newUUID(t),
			timestamp:    now,
			expectedCode: "",
			expectedErr:  nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := &gatewayService{logger: logger, opts: &defaultOptions}
			errCode, err := gs.validateGoalEvent(ctx, p.id, p.timestamp)
			assert.Equal(t, p.expectedCode, errCode)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestValidateEvaluationEvent(t *testing.T) {
	t.Parallel()
	logger, _ := log.NewLogger()
	now := time.Now().Unix()
	ctx := context.TODO()
	patterns := map[string]struct {
		id           string
		timestamp    int64
		expectedCode string
		expectedErr  error
	}{
		"err: invalid uuid": {
			id:           "0efe416e 2fd2 4996 c5c3 194f05444f1f",
			timestamp:    now,
			expectedCode: codeInvalidID,
			expectedErr:  errInvalidIDFormat,
		},
		"err: invalid timestamp": {
			id:           newUUID(t),
			timestamp:    int64(999999999999999),
			expectedCode: codeInvalidTimestamp,
			expectedErr:  errInvalidTimestamp,
		},
		"success": {
			id:           newUUID(t),
			timestamp:    now,
			expectedCode: "",
			expectedErr:  nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			gs := &gatewayService{logger: logger, opts: &defaultOptions}
			errCode, err := gs.validateEvaluationEvent(ctx, p.id, p.timestamp)
			assert.Equal(t, p.expectedCode, errCode)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestValidateMetricsEvent(t *testing.T) {
	t.Parallel()
	logger, _ := log.NewLogger()
	ctx := context.TODO()
	patterns := []struct {
		desc         string
		id           string
		expectedCode string
		expectedErr  error
	}{
		{
			desc:         "err: invalid uuid",
			id:           "0efe416e 2fd2 4996 c5c3 194f05444f1f",
			expectedCode: codeInvalidID,
			expectedErr:  errInvalidIDFormat,
		},
		{
			desc:         "success",
			id:           newUUID(t),
			expectedCode: "",
			expectedErr:  nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := &gatewayService{logger: logger}
			errCode, err := gs.validateMetricsEvent(ctx, p.id)
			assert.Equal(t, p.expectedCode, errCode)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
