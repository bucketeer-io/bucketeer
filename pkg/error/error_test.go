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

			bucketeerErr := &BktError{
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

	assert.ErrorIs(t, bucketeerErr, err3)
	assert.ErrorIs(t, bucketeerErr, err2)
	assert.ErrorIs(t, bucketeerErr, err1)
}

func TestNewErrorNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		pkg           string
		message       string
		field         string
		expectedMsg   string
		expectedField string
	}{
		{
			name:          "basic not found error",
			pkg:           "account",
			message:       "user not found",
			field:         "user_id",
			expectedMsg:   "account:user not found, user_id",
			expectedField: "user_id",
		},
		{
			name:          "not found error without args",
			pkg:           "test",
			message:       "resource not found",
			field:         "",
			expectedMsg:   "test:resource not found",
			expectedField: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewErrorNotFound(tt.pkg, tt.message, tt.field)

			assert.Equal(t, tt.pkg, err.PackageName())
			assert.Equal(t, ErrorTypeNotFound, err.ErrorType())
			assert.Equal(t, tt.expectedMsg, err.Message())
		})
	}
}

func TestNewErrorAlreadyExists(t *testing.T) {
	t.Parallel()

	err := NewErrorAlreadyExists("account", "user already exists")

	assert.Equal(t, "account", err.PackageName())
	assert.Equal(t, ErrorTypeAlreadyExists, err.ErrorType())
	assert.Equal(t, "account:user already exists", err.Message())
}

func TestNewErrorUnauthenticated(t *testing.T) {
	t.Parallel()

	err := NewErrorUnauthenticated("auth", "invalid token")

	assert.Equal(t, "auth", err.PackageName())
	assert.Equal(t, ErrorTypeUnauthenticated, err.ErrorType())
	assert.Equal(t, "auth:invalid token", err.Message())
}

func TestNewErrorPermissionDenied(t *testing.T) {
	t.Parallel()

	err := NewErrorPermissionDenied("feature", "insufficient permissions")

	assert.Equal(t, "feature", err.PackageName())
	assert.Equal(t, ErrorTypePermissionDenied, err.ErrorType())
	assert.Equal(t, "feature:insufficient permissions", err.Message())
}

func TestNewErrorUnexpectedAffectedRows(t *testing.T) {
	t.Parallel()

	err := NewErrorUnexpectedAffectedRows("database", "unexpected affected rows")

	assert.Equal(t, "database", err.PackageName())
	assert.Equal(t, ErrorTypeUnexpectedAffectedRows, err.ErrorType())
	assert.Equal(t, "database:unexpected affected rows", err.Message())
}

func TestNewErrorInternal(t *testing.T) {
	t.Parallel()

	err := NewErrorInternal("system", "internal server error")

	assert.Equal(t, "system", err.PackageName())
	assert.Equal(t, ErrorTypeInternal, err.ErrorType())
	assert.Equal(t, "system:internal server error", err.Message())
}

func TestNewErrorInvalidArgument(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		pkg           string
		message       string
		invalidType   InvalidType
		field         string
		expectedMsg   string
		expectedField string
	}{
		{
			name:          "empty field error",
			pkg:           "account",
			message:       "invalid argument",
			invalidType:   InvalidTypeEmpty,
			field:         "email",
			expectedMsg:   "account:invalid argument[email:empty]",
			expectedField: "email",
		},
		{
			name:          "nil field error",
			pkg:           "feature",
			message:       "invalid input",
			invalidType:   InvalidTypeNil,
			field:         "name",
			expectedMsg:   "feature:invalid input[name:nil]",
			expectedField: "name",
		},
		{
			name:          "format mismatch error",
			pkg:           "validation",
			message:       "format error",
			invalidType:   InvalidTypeNotMatchFormat,
			field:         "date",
			expectedMsg:   "validation:format error[date:not_match_format]",
			expectedField: "date",
		},
		{
			name:          "empty message error",
			pkg:           "test",
			message:       "",
			invalidType:   InvalidTypeEmpty,
			field:         "field",
			expectedMsg:   "test:invalid argument[field:empty]",
			expectedField: "field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewErrorInvalidArgument(tt.pkg, tt.message, tt.invalidType, tt.field)

			assert.Equal(t, tt.pkg, err.PackageName())
			assert.Equal(t, ErrorTypeInvalidArgument, err.ErrorType())
			assert.Equal(t, tt.expectedMsg, err.Message())
		})
	}
}

func TestNewErrorInvalidArgument_EmptyInvalidType(t *testing.T) {
	t.Parallel()

	err := NewErrorInvalidArgument("test", "invalid", "", "field")

	assert.Equal(t, "test", err.PackageName())
	assert.Equal(t, ErrorTypeInvalidArgument, err.ErrorType())
	assert.Equal(t, "test:invalid, field", err.Message())
}

func TestBucketeerError_Unwrap(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("original error")
	bucketeerErr := &BktError{
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
