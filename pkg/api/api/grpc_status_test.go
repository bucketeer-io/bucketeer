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

	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
)

func TestNewGRPCStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		err                 error
		anotherDetailData   []map[string]string
		expectedCode        codes.Code
		expectedMessage     string
		expectedReason      string
		expectedMetadataLen int
	}{
		{
			name:                "ErrorInvalidArgument",
			err:                 pkgErr.NewErrorInvalidArgument("test", "invalid argument", pkgErr.InvalidTypeEmpty, "field1"),
			expectedCode:        codes.InvalidArgument,
			expectedMessage:     "test:invalid argument[field1:empty]",
			expectedReason:      "INVALID_ARGUMENT",
			expectedMetadataLen: 2,
		},
		{
			name:                "ErrorNotFound",
			err:                 pkgErr.NewErrorNotFound("test", "not found", "resource"),
			expectedCode:        codes.NotFound,
			expectedMessage:     "test:not found, resource",
			expectedReason:      "NOT_FOUND",
			expectedMetadataLen: 2,
		},
		{
			name:                "ErrorAlreadyExists",
			err:                 pkgErr.NewErrorAlreadyExists("test", "already exists", "resource"),
			expectedCode:        codes.AlreadyExists,
			expectedMessage:     "test:already exists, resource",
			expectedReason:      "ALREADY_EXISTS",
			expectedMetadataLen: 2,
		},
		{
			name:                "ErrorUnauthenticated",
			err:                 pkgErr.NewErrorUnauthenticated("test", "unauthenticated"),
			expectedCode:        codes.Unauthenticated,
			expectedMessage:     "test:unauthenticated",
			expectedReason:      "UNAUTHENTICATED",
			expectedMetadataLen: 1,
		},
		{
			name:                "ErrorPermissionDenied",
			err:                 pkgErr.NewErrorPermissionDenied("test", "permission denied"),
			expectedCode:        codes.PermissionDenied,
			expectedMessage:     "test:permission denied",
			expectedReason:      "PERMISSION_DENIED",
			expectedMetadataLen: 1,
		},
		{
			name:                "ErrorUnexpectedAffectedRows",
			err:                 pkgErr.NewErrorUnexpectedAffectedRows("test", "unexpected affected rows"),
			expectedCode:        codes.Internal,
			expectedMessage:     "test:unexpected affected rows",
			expectedReason:      "UNEXPECTED_AFFECTED_ROWS",
			expectedMetadataLen: 1,
		},
		{
			name:                "ErrorInternal",
			err:                 pkgErr.NewErrorInternal("test", "internal error"),
			expectedCode:        codes.Internal,
			expectedMessage:     "test:internal error",
			expectedReason:      "INTERNAL",
			expectedMetadataLen: 1,
		},
		{
			name:                "Non-BucketeerError",
			err:                 errors.New("standard error"),
			expectedCode:        codes.Unknown,
			expectedMessage:     "standard error",
			expectedReason:      "UNKNOWN",
			expectedMetadataLen: 0,
		},
		{
			name:                "ErrorInvalidArgument with additional metadata",
			err:                 pkgErr.NewErrorInvalidArgument("test", "invalid argument", pkgErr.InvalidTypeEmpty, "field1"),
			anotherDetailData:   []map[string]string{{"additional": "data"}},
			expectedCode:        codes.InvalidArgument,
			expectedMessage:     "test:invalid argument[field1:empty]",
			expectedReason:      "INVALID_ARGUMENT",
			expectedMetadataLen: 3,
		},
		{
			name:                "ErrorNotFound with multiple additional metadata",
			err:                 pkgErr.NewErrorNotFound("test", "not found", "resource"),
			anotherDetailData:   []map[string]string{{"key1": "value1"}, {"key2": "value2"}},
			expectedCode:        codes.NotFound,
			expectedMessage:     "test:not found, resource",
			expectedReason:      "NOT_FOUND",
			expectedMetadataLen: 4,
		},
		{
			name:                "Non-BucketeerError with additional metadata",
			err:                 errors.New("standard error"),
			anotherDetailData:   []map[string]string{{"additional": "data"}},
			expectedCode:        codes.Unknown,
			expectedMessage:     "standard error",
			expectedReason:      "UNKNOWN",
			expectedMetadataLen: 1,
		},
		{
			name:                "another metadata is nil",
			err:                 pkgErr.NewErrorNotFound("test", "not found", "resource"),
			anotherDetailData:   nil,
			expectedCode:        codes.NotFound,
			expectedMessage:     "test:not found, resource",
			expectedReason:      "NOT_FOUND",
			expectedMetadataLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			st := NewGRPCStatus(tt.err, tt.anotherDetailData...)

			assert.Equal(t, tt.expectedCode, st.Code())

			assert.Equal(t, tt.expectedMessage, st.Message())

			details := st.Details()
			assert.Len(t, details, tt.expectedMetadataLen)

			for _, detail := range details {
				if errorInfo, ok := detail.(*errdetails.ErrorInfo); ok {
					assert.Equal(t, tt.expectedReason, errorInfo.Reason)
					assert.NotEmpty(t, errorInfo.Metadata)
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

func TestNewGRPCStatus_WithDetailsError(t *testing.T) {
	t.Parallel()

	err := pkgErr.NewErrorInternal("test", "internal error")

	largeMetadata := make([]map[string]string, 1000)
	for i := 0; i < 1000; i++ {
		largeMetadata[i] = map[string]string{
			"key": "value",
		}
	}

	st := NewGRPCStatus(err, largeMetadata...)

	assert.Equal(t, codes.Internal, st.Code())
}

func TestNewGRPCStatus_ErrorInvalidArgumentTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		invalidType    pkgErr.InvalidType
		expectedReason string
	}{
		{
			name:           "InvalidTypeEmpty",
			invalidType:    pkgErr.InvalidTypeEmpty,
			expectedReason: "INVALID_ARGUMENT",
		},
		{
			name:           "InvalidTypeNil",
			invalidType:    pkgErr.InvalidTypeNil,
			expectedReason: "INVALID_ARGUMENT",
		},
		{
			name:           "InvalidTypeNotMatchFormat",
			invalidType:    pkgErr.InvalidTypeNotMatchFormat,
			expectedReason: "INVALID_ARGUMENT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := pkgErr.NewErrorInvalidArgument("test", "invalid", tt.invalidType, "field")
			st := NewGRPCStatus(err)

			assert.Equal(t, codes.InvalidArgument, st.Code())

			details := st.Details()
			assert.Len(t, details, 2)
			for _, detail := range details {
				if errorInfo, ok := detail.(*errdetails.ErrorInfo); ok {
					assert.Equal(t, tt.expectedReason, errorInfo.Reason)
				}
			}
		})
	}
}

func TestNewGRPCStatus_MetadataHandling(t *testing.T) {
	t.Parallel()

	err := pkgErr.NewErrorInvalidArgument("test", "multiple fields invalid", pkgErr.InvalidTypeEmpty, "field1", "field2", "field3")

	st := NewGRPCStatus(err)

	assert.Equal(t, codes.InvalidArgument, st.Code())

	details := st.Details()
	assert.Len(t, details, 4)

	for _, detail := range details {
		if errorInfo, ok := detail.(*errdetails.ErrorInfo); ok {
			assert.True(t, ok)
			assert.Equal(t, "INVALID_ARGUMENT", errorInfo.Reason)
			assert.Equal(t, "test.invalid_argument.empty", errorInfo.Metadata["messageKey"])
		} else if localizedMessage, ok := detail.(*errdetails.LocalizedMessage); ok {
			assert.Equal(t, "en", localizedMessage.Locale)
			assert.Equal(t, st.Message(), localizedMessage.Message)
		} else {
			assert.Fail(t, "unexpected detail type")
		}
	}
}

func TestNewGRPCStatus_ActualMessageFormat(t *testing.T) {
	t.Parallel()

	err := pkgErr.NewErrorNotFound("test", "not found", "resource")
	st := NewGRPCStatus(err)

	t.Logf("Actual message: %s", st.Message())
	t.Logf("Actual code: %v", st.Code())

	details := st.Details()
	for i, detail := range details {
		errorInfo, ok := detail.(*errdetails.ErrorInfo)
		if ok {
			t.Logf("Detail %d - Reason: %s, Domain: %s, Metadata: %v", i, errorInfo.Reason, errorInfo.Domain, errorInfo.Metadata)
		}
	}

	assert.Equal(t, codes.NotFound, st.Code())
	assert.NotEmpty(t, st.Message())
}

func TestNewGRPCStatus_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		err               error
		anotherDetailData []map[string]string
		expectedCode      codes.Code
		expectedMessage   string
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
		{
			name:              "Empty metadata",
			err:               pkgErr.NewErrorNotFound("test", "not found"),
			anotherDetailData: []map[string]string{},
			expectedCode:      codes.NotFound,
			expectedMessage:   "test:not found",
		},
		{
			name:              "Nil metadata",
			err:               pkgErr.NewErrorNotFound("test", "not found"),
			anotherDetailData: nil,
			expectedCode:      codes.NotFound,
			expectedMessage:   "test:not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			st := NewGRPCStatus(tt.err, tt.anotherDetailData...)

			assert.Equal(t, tt.expectedCode, st.Code())
			assert.Equal(t, tt.expectedMessage, st.Message())
		})
	}
}
