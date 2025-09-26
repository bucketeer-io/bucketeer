// Copyright 2025 The Bucketeer Authors.
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

package jobs

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
)

func TestGetErrorType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error returns internal",
			err:      nil,
			expected: ErrorTypeInternal,
		},
		{
			name:     "context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: ErrorTypeTimeout,
		},
		{
			name:     "context canceled",
			err:      context.Canceled,
			expected: ErrorTypeTimeout,
		},
		{
			name:     "gRPC NotFound error",
			err:      status.Error(codes.NotFound, "environment not found"),
			expected: ErrorTypeNotFound,
		},
		{
			name:     "gRPC InvalidArgument error",
			err:      status.Error(codes.InvalidArgument, "invalid request"),
			expected: ErrorTypeValidation,
		},
		{
			name:     "gRPC FailedPrecondition error",
			err:      status.Error(codes.FailedPrecondition, "operation not allowed"),
			expected: ErrorTypeValidation,
		},
		{
			name:     "gRPC DeadlineExceeded error",
			err:      status.Error(codes.DeadlineExceeded, "request timeout"),
			expected: ErrorTypeTimeout,
		},
		{
			name:     "gRPC Internal error",
			err:      status.Error(codes.Internal, "database error"),
			expected: ErrorTypeInternal,
		},
		{
			name:     "gRPC Unauthenticated error",
			err:      status.Error(codes.Unauthenticated, "authentication failed"),
			expected: ErrorTypeInternal,
		},
		{
			name:     "cache not found error",
			err:      cache.ErrNotFound,
			expected: ErrorTypeNotFound,
		},
		{
			name:     "experiment before start error",
			err:      domain.ErrExperimentBeforeStart,
			expected: ErrorTypeValidation,
		},
		{
			name:     "experiment before stop error",
			err:      domain.ErrExperimentBeforeStop,
			expected: ErrorTypeValidation,
		},
		{
			name:     "database connection error",
			err:      errors.New("database connection failed"),
			expected: ErrorTypeInternal,
		},
		{
			name:     "network error",
			err:      errors.New("network unreachable"),
			expected: ErrorTypeInternal,
		},
		{
			name:     "unknown error",
			err:      errors.New("unknown error occurred"),
			expected: ErrorTypeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := GetErrorType(tt.err)
			if result != tt.expected {
				t.Errorf("GetErrorType(%v) = %s, want %s", tt.err, result, tt.expected)
			}
		})
	}
}
