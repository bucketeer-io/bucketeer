// Copyright 2022 The Bucketeer Authors.
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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	cachev3mock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestNewTrackHandler(t *testing.T) {
	t.Parallel()
	h := NewTrackHandler(nil, nil, nil)
	assert.IsType(t, &TrackHandler{}, h)
}

func TestServeHTTP(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()

	patterns := []struct {
		desc     string
		setup    func(*testing.T, *TrackHandler)
		input    *http.Request
		expected int
	}{
		{
			desc: "fail: bad params",
			input: httptest.NewRequest("GET",
				"/track?apikey=akey&userid=uid&goalid=gid&tag=t&timestamp=abc",
				nil),
			expected: http.StatusBadRequest,
		},
		{
			desc: "fail: publish error",
			setup: func(t *testing.T, h *TrackHandler) {
				h.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						EnvironmentNamespace: "ns0",
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				h.goalBatchPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(errors.New("internal")).MaxTimes(1)
			},
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?apikey=akey&userid=uid&goalid=gid&tag=t&timestamp=%d", now.Unix()),
				nil),
			expected: http.StatusInternalServerError,
		},
		{
			desc: "success: without value",
			setup: func(t *testing.T, h *TrackHandler) {
				h.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						EnvironmentNamespace: "ns0",
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				h.goalBatchPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?apikey=akey&userid=uid&goalid=gid&tag=t&timestamp=%d", now.Unix()),
				nil),
			expected: http.StatusOK,
		},
		{
			desc: "success: with value",
			setup: func(t *testing.T, h *TrackHandler) {
				h.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						EnvironmentNamespace: "ns0",
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				h.goalBatchPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?apikey=akey&userid=uid&goalid=gid&tag=t&timestamp=%d&value=1.234", now.Unix()),
				nil),
			expected: http.StatusOK,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			h := newTrackHandlerWithMock(t, mockController)
			if p.setup != nil {
				p.setup(t, h)
			}
			actual := httptest.NewRecorder()
			h.ServeHTTP(actual, p.input)
			assert.Equal(t, p.expected, actual.Code)
		})
	}
}

func TestValidateParams(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()

	patterns := []struct {
		desc        string
		input       *http.Request
		expected    *params
		expectedErr error
	}{
		{
			desc: "err: errAPIKeyEmpty",
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?userid=uid&goalid=gid&tag=t&timestamp=%d&value=1.234", now.Unix()),
				nil),
			expected:    nil,
			expectedErr: errAPIKeyEmpty,
		},
		{
			desc: "err: errUserIDEmpty",
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?apikey=akey&goalid=gid&tag=t&timestamp=%d&value=1.234", now.Unix()),
				nil),
			expected:    nil,
			expectedErr: errUserIDEmpty,
		},
		{
			desc: "err: errGoalIDEmpty",
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?apikey=akey&userid=uid&tag=t&timestamp=%d&value=1.234", now.Unix()),
				nil),
			expected:    nil,
			expectedErr: errGoalIDEmpty,
		},
		{
			desc: "err: errTagEmpty",
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?apikey=akey&userid=uid&goalid=gid&timestamp=%d&value=1.234", now.Unix()),
				nil),
			expected:    nil,
			expectedErr: errTagEmpty,
		},
		{
			desc: "err: errTimestampEmpty",
			input: httptest.NewRequest("GET",
				"/track?apikey=akey&userid=uid&goalid=gid&tag=t&value=1.234",
				nil),
			expected:    nil,
			expectedErr: errTimestampEmpty,
		},
		{
			desc: "err: errTimestampInvalid",
			input: httptest.NewRequest("GET",
				"/track?apikey=akey&userid=uid&goalid=gid&tag=t&timestamp=abc&value=1.234",
				nil),
			expected:    nil,
			expectedErr: errTimestampInvalid,
		},
		{
			desc: "err: errTimestampInvalid: out of window",
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?apikey=akey&userid=uid&goalid=gid&tag=t&timestamp=%d&value=1.234", now.AddDate(0, 0, 2).Unix()),
				nil),
			expected:    nil,
			expectedErr: errTimestampInvalid,
		},
		{
			desc: "err: errValueInvalid",
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?apikey=akey&userid=uid&goalid=gid&tag=t&timestamp=%d&value=abc", now.Unix()),
				nil),
			expected:    nil,
			expectedErr: errValueInvalid,
		},
		{
			desc: "success: without value",
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?apikey=akey&userid=uid&goalid=gid&tag=t&timestamp=%d", now.Unix()),
				nil),
			expected: &params{
				apiKey:    "akey",
				userID:    "uid",
				goalID:    "gid",
				tag:       "t",
				timestamp: now.Unix(),
				value:     float64(0),
			},
			expectedErr: nil,
		},
		{
			desc: "success: with value",
			input: httptest.NewRequest("GET",
				fmt.Sprintf("/track?apikey=akey&userid=uid&goalid=gid&tag=t&timestamp=%d&value=1.234", now.Unix()),
				nil),
			expected: &params{
				apiKey:    "akey",
				userID:    "uid",
				goalID:    "gid",
				tag:       "t",
				timestamp: now.Unix(),
				value:     float64(1.234),
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			h := newTrackHandlerWithMock(t, mockController)
			actual, err := h.validateParams(p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newTrackHandlerWithMock(t *testing.T, mockController *gomock.Controller) *TrackHandler {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &TrackHandler{
		accountClient:          accountclientmock.NewMockClient(mockController),
		goalBatchPublisher:     publishermock.NewMockPublisher(mockController),
		environmentAPIKeyCache: cachev3mock.NewMockEnvironmentAPIKeyCache(mockController),
		opts:                   &defaultOptions,
		logger:                 logger,
	}
}
