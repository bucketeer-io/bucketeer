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

package error

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketeerError_BasicProperties(t *testing.T) {
	t.Parallel()

	err := NewErrorNotFound("test", "not found", "resource")

	assert.Equal(t, "test", err.PackageName())
	assert.Equal(t, ErrorTypeNotFound, err.ErrorType())
	assert.Equal(t, "test:not found, resource", err.Message())
	assert.Equal(t, "test:not found, resource", err.Error())
}

func TestBucketeerError_Metadatas(t *testing.T) {
	t.Parallel()

	err := NewErrorNotFound("test", "not found", "resource1", "resource2")

	metadatas := err.Metadatas()
	assert.Len(t, metadatas, 2)

	expectedMessageKey := "test.not_found"
	for _, metadata := range metadatas {
		assert.Equal(t, expectedMessageKey, metadata["messageKey"])
		assert.Contains(t, []string{"resource1", "resource2"}, metadata["field"])
	}
}

func TestBucketeerError_AddMetadata(t *testing.T) {
	t.Parallel()

	err := NewErrorNotFound("test", "not found", "resource")
	initialLen := len(err.Metadatas())

	newMetadata := map[string]string{"key1": "value1", "key2": "value2"}
	err.AddMetadata(newMetadata)

	metadatas := err.Metadatas()
	assert.Len(t, metadatas, initialLen+1)
	assert.Equal(t, newMetadata, metadatas[initialLen])
}

func TestBucketeerError_JoinError(t *testing.T) {
	t.Parallel()

	err := NewErrorNotFound("test", "not found", "resource")
	originalErr := err.err

	joinedErr := errors.New("additional error")
	err.JoinError(joinedErr)

	assert.NotEqual(t, originalErr, err.err)
	assert.ErrorIs(t, err.err, joinedErr)
}

func TestBucketeerError_JoinError_ErrorsIs(t *testing.T) {
	t.Parallel()
	joinedError := errors.New("database error")
	otherError := errors.New("non-existent error")

	tests := []struct {
		name          string
		joinedError   error
		checkError    error
		shouldBeFound bool
	}{
		{
			name:          "joined error should be found with errors.Is",
			joinedError:   joinedError,
			checkError:    joinedError,
			shouldBeFound: true,
		},
		{
			name:          "non-existent error should not be found",
			joinedError:   joinedError,
			checkError:    otherError,
			shouldBeFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bucketeerErr := &BucketeerError{
				packageName: "test",
				errorType:   ErrorTypeNotFound,
				message:     "test error",
			}

			bucketeerErr.JoinError(tt.joinedError)

			if tt.shouldBeFound {
				assert.True(t, errors.Is(bucketeerErr, tt.checkError), "Expected error to be found")
			} else {
				assert.False(t, errors.Is(bucketeerErr, tt.checkError), "Expected error to not be found")
			}
		})
	}
}

func TestBucketeerError_JoinError_MultipleErrors(t *testing.T) {
	t.Parallel()

	bucketeerErr := NewErrorNotFound("test", "not found", "resource")

	err1 := errors.New("database error")
	err2 := errors.New("network error")
	err3 := errors.New("validation error")

	bucketeerErr.JoinError(err1)
	bucketeerErr.JoinError(err2)
	bucketeerErr.JoinError(err3)

	assert.ErrorIs(t, bucketeerErr, err1)
	assert.ErrorIs(t, bucketeerErr, err2)
	assert.ErrorIs(t, bucketeerErr, err3)

	nonExistentErr := errors.New("non-existent error")
	assert.False(t, errors.Is(bucketeerErr, nonExistentErr))
}

func TestBucketeerError_JoinError_WithWrappedErrors(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("original error")
	wrappedErr := fmt.Errorf("wrapped: %w", originalErr)

	bucketeerErr := NewErrorNotFound("test", "not found", "resource")
	bucketeerErr.JoinError(wrappedErr)

	assert.ErrorIs(t, bucketeerErr, wrappedErr)

	assert.ErrorIs(t, bucketeerErr, originalErr)
}

func TestBucketeerError_JoinError_Chain(t *testing.T) {
	t.Parallel()

	bucketeerErr := NewErrorNotFound("test", "not found", "resource")

	err1 := errors.New("level 1 error")
	err2 := fmt.Errorf("level 2: %w", err1)
	err3 := fmt.Errorf("level 3: %w", err2)

	bucketeerErr.JoinError(err3)

	// チェーンの各レベルが検知できることを確認
	assert.ErrorIs(t, bucketeerErr, err3)
	assert.ErrorIs(t, bucketeerErr, err2)
	assert.ErrorIs(t, bucketeerErr, err1)
}

func TestNewErrorNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		pkg            string
		message        string
		args           []string
		expectedMsg    string
		expectedFields []string
	}{
		{
			name:           "basic not found error",
			pkg:            "account",
			message:        "user not found",
			args:           []string{"user_id"},
			expectedMsg:    "account:user not found, user_id",
			expectedFields: []string{"user_id"},
		},
		{
			name:           "not found error with multiple fields",
			pkg:            "feature",
			message:        "feature not found",
			args:           []string{"feature_id", "environment_id"},
			expectedMsg:    "feature:feature not found, feature_id, environment_id",
			expectedFields: []string{"feature_id", "environment_id"},
		},
		{
			name:           "not found error without args",
			pkg:            "test",
			message:        "resource not found",
			args:           []string{},
			expectedMsg:    "test:resource not found",
			expectedFields: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewErrorNotFound(tt.pkg, tt.message, tt.args...)

			assert.Equal(t, tt.pkg, err.PackageName())
			assert.Equal(t, ErrorTypeNotFound, err.ErrorType())
			assert.Equal(t, tt.expectedMsg, err.Message())

			metadatas := err.Metadatas()
			assert.Len(t, metadatas, len(tt.expectedFields))

			for i, field := range tt.expectedFields {
				assert.Equal(t, field, metadatas[i]["field"])
				assert.Equal(t, tt.pkg+".not_found", metadatas[i]["messageKey"])
			}
		})
	}
}

func TestNewErrorAlreadyExists(t *testing.T) {
	t.Parallel()

	err := NewErrorAlreadyExists("account", "user already exists", "email")

	assert.Equal(t, "account", err.PackageName())
	assert.Equal(t, ErrorTypeAlreadyExists, err.ErrorType())
	assert.Equal(t, "account:user already exists, email", err.Message())

	metadatas := err.Metadatas()
	assert.Len(t, metadatas, 1)
	assert.Equal(t, "email", metadatas[0]["field"])
	assert.Equal(t, "account.already_exists", metadatas[0]["messageKey"])
}

func TestNewErrorUnauthenticated(t *testing.T) {
	t.Parallel()

	err := NewErrorUnauthenticated("auth", "invalid token")

	assert.Equal(t, "auth", err.PackageName())
	assert.Equal(t, ErrorTypeUnauthenticated, err.ErrorType())
	assert.Equal(t, "auth:invalid token", err.Message())

	metadatas := err.Metadatas()
	assert.Len(t, metadatas, 0)
}

func TestNewErrorPermissionDenied(t *testing.T) {
	t.Parallel()

	err := NewErrorPermissionDenied("feature", "insufficient permissions", "feature_id")

	assert.Equal(t, "feature", err.PackageName())
	assert.Equal(t, ErrorTypePermissionDenied, err.ErrorType())
	assert.Equal(t, "feature:insufficient permissions, feature_id", err.Message())

	metadatas := err.Metadatas()
	assert.Len(t, metadatas, 1)
	assert.Equal(t, "feature_id", metadatas[0]["field"])
	assert.Equal(t, "feature.permission_denied", metadatas[0]["messageKey"])
}

func TestNewErrorUnexpectedAffectedRows(t *testing.T) {
	t.Parallel()

	err := NewErrorUnexpectedAffectedRows("database", "unexpected affected rows")

	assert.Equal(t, "database", err.PackageName())
	assert.Equal(t, ErrorTypeUnexpectedAffectedRows, err.ErrorType())
	assert.Equal(t, "database:unexpected affected rows", err.Message())

	metadatas := err.Metadatas()
	assert.Len(t, metadatas, 0)
}

func TestNewErrorInternal(t *testing.T) {
	t.Parallel()

	err := NewErrorInternal("system", "internal server error")

	assert.Equal(t, "system", err.PackageName())
	assert.Equal(t, ErrorTypeInternal, err.ErrorType())
	assert.Equal(t, "system:internal server error", err.Message())

	metadatas := err.Metadatas()
	assert.Len(t, metadatas, 0)
}

func TestNewErrorInvalidArgument(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		pkg            string
		message        string
		invalidType    InvalidType
		args           []string
		expectedMsg    string
		expectedFields []string
	}{
		{
			name:           "empty field error",
			pkg:            "account",
			message:        "invalid argument",
			invalidType:    InvalidTypeEmpty,
			args:           []string{"email"},
			expectedMsg:    "account:invalid argument[email:empty]",
			expectedFields: []string{"email"},
		},
		{
			name:           "nil field error",
			pkg:            "feature",
			message:        "invalid input",
			invalidType:    InvalidTypeNil,
			args:           []string{"name"},
			expectedMsg:    "feature:invalid input[name:nil]",
			expectedFields: []string{"name"},
		},
		{
			name:           "format mismatch error",
			pkg:            "validation",
			message:        "format error",
			invalidType:    InvalidTypeNotMatchFormat,
			args:           []string{"date"},
			expectedMsg:    "validation:format error[date:not_match_format]",
			expectedFields: []string{"date"},
		},
		{
			name:           "multiple fields error",
			pkg:            "form",
			message:        "validation failed",
			invalidType:    InvalidTypeEmpty,
			args:           []string{"username", "password"},
			expectedMsg:    "form:validation failed[username:empty][password:empty]",
			expectedFields: []string{"username", "password"},
		},
		{
			name:           "empty message error",
			pkg:            "test",
			message:        "",
			invalidType:    InvalidTypeEmpty,
			args:           []string{"field"},
			expectedMsg:    "test:invalid argument[field:empty]",
			expectedFields: []string{"field"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewErrorInvalidArgument(tt.pkg, tt.message, tt.invalidType, tt.args...)

			assert.Equal(t, tt.pkg, err.PackageName())
			assert.Equal(t, ErrorTypeInvalidArgument, err.ErrorType())
			assert.Equal(t, tt.expectedMsg, err.Message())

			metadatas := err.Metadatas()
			assert.Len(t, metadatas, len(tt.expectedFields))

			expectedMessageKey := tt.pkg + ".invalid_argument." + string(tt.invalidType)
			for i, field := range tt.expectedFields {
				assert.Equal(t, field, metadatas[i]["field"])
				assert.Equal(t, expectedMessageKey, metadatas[i]["messageKey"])
			}
		})
	}
}

func TestNewErrorInvalidArgument_EmptyInvalidType(t *testing.T) {
	t.Parallel()

	err := NewErrorInvalidArgument("test", "invalid", "", "field")

	assert.Equal(t, "test", err.PackageName())
	assert.Equal(t, ErrorTypeInvalidArgument, err.ErrorType())
	assert.Equal(t, "test:invalid[field:unknown]", err.Message())

	metadatas := err.Metadatas()
	assert.Len(t, metadatas, 1)
	assert.Equal(t, "field", metadatas[0]["field"])
	assert.Equal(t, "test.invalid_argument.unknown", metadatas[0]["messageKey"])
}

func TestBucketeerError_Unwrap(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("original error")
	bucketeerErr := &BucketeerError{
		packageName: "test",
		errorType:   ErrorTypeInternal,
		message:     "test error",
		err:         originalErr,
	}

	unwrapped := bucketeerErr.Unwrap()
	assert.Equal(t, originalErr, unwrapped)
}

func TestBucketeerError_EmptyArgs(t *testing.T) {
	t.Parallel()

	err := NewErrorNotFound("test", "not found", "")

	assert.Equal(t, "test:not found", err.Message())

	metadatas := err.Metadatas()
	assert.Len(t, metadatas, 1)
	assert.Equal(t, "", metadatas[0]["field"])
}

func TestBucketeerError_MultipleEmptyArgs(t *testing.T) {
	t.Parallel()

	err := NewErrorNotFound("test", "not found", "", "", "valid_field")

	assert.Equal(t, "test:not found, valid_field", err.Message())

	metadatas := err.Metadatas()
	assert.Len(t, metadatas, 3)
	assert.Equal(t, "", metadatas[0]["field"])
	assert.Equal(t, "", metadatas[1]["field"])
	assert.Equal(t, "valid_field", metadatas[2]["field"])
}

func TestErrorType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		errorType ErrorType
		expected  string
	}{
		{ErrorTypeNotFound, "not_found"},
		{ErrorTypeAlreadyExists, "already_exists"},
		{ErrorTypeUnauthenticated, "unauthenticated"},
		{ErrorTypePermissionDenied, "permission_denied"},
		{ErrorTypeUnexpectedAffectedRows, "unexpected_affected_rows"},
		{ErrorTypeInternal, "internal"},
		{ErrorTypeInvalidArgument, "invalid_argument"},
	}

	for _, tt := range tests {
		t.Run(string(tt.errorType), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, string(tt.errorType))
		})
	}
}

func TestInvalidType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		invalidType InvalidType
		expected    string
	}{
		{InvalidTypeEmpty, "empty"},
		{InvalidTypeNil, "nil"},
		{InvalidTypeNotMatchFormat, "not_match_format"},
	}

	for _, tt := range tests {
		t.Run(string(tt.invalidType), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, string(tt.invalidType))
		})
	}
}
