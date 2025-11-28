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

package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

func TestNewGRPCStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		err              error
		expectedCode     codes.Code
		expectedMessage  string
		expectedReason   string
		expectedMetadata map[string]string
	}{
		{
			name:            "ErrorInvalidEmpty",
			err:             pkgErr.NewErrorInvalidArgEmpty("test", "invalid argument", "field1"),
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "test:invalid argument[field1:InvalidArgumentEmptyError]",
			expectedReason:  "INVALID_ARGUMENT_EMPTY",
			expectedMetadata: map[string]string{
				"messagekey": "InvalidArgumentEmptyError",
				"field":      "field1",
			},
		},
		{
			name:            "ErrorInvalidNil",
			err:             pkgErr.NewErrorInvalidArgNil("test", "invalid argument", "field1"),
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "test:invalid argument[field1:InvalidArgumentNilError]",
			expectedReason:  "INVALID_ARGUMENT_NIL",
			expectedMetadata: map[string]string{
				"messagekey": "InvalidArgumentNilError",
				"field":      "field1",
			},
		},
		{
			name:            "ErrorInvalidNotMatchFormat",
			err:             pkgErr.NewErrorInvalidArgNotMatchFormat("test", "invalid argument", "field1"),
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "test:invalid argument[field1:InvalidArgumentNotMatchFormatError]",
			expectedReason:  "INVALID_ARGUMENT_NOT_MATCH_FORMAT",
			expectedMetadata: map[string]string{
				"messagekey": "InvalidArgumentNotMatchFormatError",
				"field":      "field1",
			},
		},
		{
			name:            "ErrorInvalidUnknown",
			err:             pkgErr.NewErrorInvalidArgUnknown("test", "invalid argument", "field1"),
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "test:invalid argument[field1:InvalidArgumentUnknownError]",
			expectedReason:  "INVALID_ARGUMENT",
			expectedMetadata: map[string]string{
				"messagekey": "InvalidArgumentUnknownError",
				"field":      "field1",
			},
		},
		{
			name:            "ErrorNotFound",
			err:             pkgErr.NewErrorNotFound("test", "not found", "resource"),
			expectedCode:    codes.NotFound,
			expectedMessage: "test:not found, resource",
			expectedReason:  "NOT_FOUND",
			expectedMetadata: map[string]string{
				"messagekey": "NotFoundError",
				"field":      "resource",
			},
		},
		{
			name:            "ErrorAlreadyExists",
			err:             pkgErr.NewErrorAlreadyExists("test", "already exists"),
			expectedCode:    codes.AlreadyExists,
			expectedMessage: "test:already exists",
			expectedReason:  "ALREADY_EXISTS",
			expectedMetadata: map[string]string{
				"messagekey": "AlreadyExistsError",
			},
		},
		{
			name:            "ErrorUnauthenticated",
			err:             pkgErr.NewErrorUnauthenticated("test", "unauthenticated"),
			expectedCode:    codes.Unauthenticated,
			expectedMessage: "test:unauthenticated",
			expectedReason:  "UNAUTHENTICATED",
			expectedMetadata: map[string]string{
				"messagekey": "UnauthenticatedError",
			},
		},
		{
			name:            "ErrorPermissionDenied",
			err:             pkgErr.NewErrorPermissionDenied("test", "permission denied"),
			expectedCode:    codes.PermissionDenied,
			expectedMessage: "test:permission denied",
			expectedReason:  "PERMISSION_DENIED",
			expectedMetadata: map[string]string{
				"messagekey": "PermissionDenied",
			},
		},
		{
			name:            "ErrorUnexpectedAffectedRows",
			err:             pkgErr.NewErrorUnexpectedAffectedRows("test", "unexpected affected rows"),
			expectedCode:    codes.Internal,
			expectedMessage: "test:unexpected affected rows",
			expectedReason:  "UNEXPECTED_AFFECTED_ROWS",
			expectedMetadata: map[string]string{
				"messagekey": "UnexpectedAffectedRows",
			},
		},
		{
			name:            "ErrorInternal",
			err:             pkgErr.NewErrorInternal("test", "internal error"),
			expectedCode:    codes.Internal,
			expectedMessage: "test:internal error",
			expectedReason:  "INTERNAL",
			expectedMetadata: map[string]string{
				"messagekey": "InternalServerError",
			},
		},
		{
			name:            "Non-BucketeerError",
			err:             errors.New("standard error"),
			expectedCode:    codes.Unknown,
			expectedMessage: "standard error",
			expectedReason:  "UNKNOWN",
			expectedMetadata: map[string]string{
				"messagekey": "unknown",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			st := NewGRPCStatus(tt.err)

			assert.Equal(t, tt.expectedCode, st.Code())

			assert.Equal(t, tt.expectedMessage, st.Message())

			details := st.Details()

			for _, detail := range details {
				if errorInfo, ok := detail.(*errdetails.ErrorInfo); ok {
					assert.Equal(t, tt.expectedReason, errorInfo.Reason)
					assert.NotEmpty(t, errorInfo.Metadata)
					assert.Equal(t, tt.expectedMetadata["messagekey"], errorInfo.Metadata["messageKey"])
					assert.Equal(t, tt.expectedMetadata["field"], errorInfo.Metadata["field"])
				} else if localizedMessage, ok := detail.(*errdetails.LocalizedMessage); ok {
					assert.Equal(t, "en", localizedMessage.Locale)
					assert.Equal(t, st.Message(), localizedMessage.Message)
				}
			}
		})
	}
}

func TestNewGRPCStatus_NilError(t *testing.T) {
	t.Parallel()

	st := NewGRPCStatus(nil)

	assert.Equal(t, codes.Unknown, st.Code())
	assert.Equal(t, "", st.Message())
	assert.Len(t, st.Details(), 0)
}

func TestNewGRPCStatus_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		err             error
		expectedCode    codes.Code
		expectedMessage string
	}{
		{
			name:            "Empty string error",
			err:             errors.New(""),
			expectedCode:    codes.Unknown,
			expectedMessage: "",
		},
		{
			name:            "Special characters in error message",
			err:             errors.New("error with special chars: !@#$%^&*()"),
			expectedCode:    codes.Unknown,
			expectedMessage: "error with special chars: !@#$%^&*()",
		},
		{
			name:            "Unicode characters in error message",
			err:             errors.New("error message with unicode: 🚀"),
			expectedCode:    codes.Unknown,
			expectedMessage: "error message with unicode: 🚀",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			st := NewGRPCStatus(tt.err)

			assert.Equal(t, tt.expectedCode, st.Code())
			assert.Equal(t, tt.expectedMessage, st.Message())
		})
	}
}
