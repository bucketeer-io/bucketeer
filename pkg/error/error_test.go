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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketeerError_BasicProperties(t *testing.T) {
	t.Parallel()

	err := NewErrorNotFound("test", "not found", "resource")

	assert.Equal(t, "test", err.PackageName())
	assert.Equal(t, ErrorTypeNotFound, err.ErrorType())
	assert.Equal(t, "test:not found, resource", err.Error())
}

func TestBucketeerError_Wrap(t *testing.T) {
	t.Parallel()

	wrappedErr := errors.New("additional error")
	err := NewErrorNotFound("test", "not found", "resource")
	err.Wrap(wrappedErr)

	assert.ErrorIs(t, err, err)
	assert.ErrorIs(t, err, wrappedErr)
}

func TestNewErrorNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		pkg                  string
		message              string
		field                string
		wrappedError         error
		expectedErrorMessage string
		expectedField        string
	}{
		{
			name:                 "basic not found error",
			pkg:                  "account",
			message:              "user not found",
			field:                "user_id",
			expectedErrorMessage: "account:user not found, user_id",
			expectedField:        "user_id",
		},
		{
			name:                 "not found error without args",
			pkg:                  "test",
			message:              "resource not found",
			field:                "",
			expectedErrorMessage: "test:resource not found",
			expectedField:        "",
		},
		{
			name:                 "wrapped error",
			pkg:                  "test",
			message:              "resource not found",
			field:                "",
			wrappedError:         errors.New("wrapped error"),
			expectedErrorMessage: "test:resource not found: wrapped error",
			expectedField:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewErrorNotFound(tt.pkg, tt.message, tt.field)
			err.Wrap(tt.wrappedError)

			assert.Equal(t, tt.pkg, err.PackageName())
			assert.Equal(t, ErrorTypeNotFound, err.ErrorType())
			assert.Equal(t, tt.expectedErrorMessage, err.Error())
		})
	}
}

func TestNewErrorAlreadyExists(t *testing.T) {
	t.Parallel()

	err := NewErrorAlreadyExists("account", "user already exists")

	assert.Equal(t, "account", err.PackageName())
	assert.Equal(t, ErrorTypeAlreadyExists, err.ErrorType())
	assert.Equal(t, "user already exists", err.message)
	assert.Equal(t, "account:user already exists", err.Error())
}

func TestNewErrorUnauthenticated(t *testing.T) {
	t.Parallel()

	err := NewErrorUnauthenticated("auth", "invalid token")

	assert.Equal(t, "auth", err.PackageName())
	assert.Equal(t, ErrorTypeUnauthenticated, err.ErrorType())
	assert.Equal(t, "invalid token", err.message)
	assert.Equal(t, "auth:invalid token", err.Error())
}

func TestNewErrorPermissionDenied(t *testing.T) {
	t.Parallel()

	err := NewErrorPermissionDenied("feature", "insufficient permissions")

	assert.Equal(t, "feature", err.PackageName())
	assert.Equal(t, ErrorTypePermissionDenied, err.ErrorType())
	assert.Equal(t, "feature:insufficient permissions", err.Error())
}

func TestNewErrorUnexpectedAffectedRows(t *testing.T) {
	t.Parallel()

	err := NewErrorUnexpectedAffectedRows("database", "unexpected affected rows")

	assert.Equal(t, "database", err.PackageName())
	assert.Equal(t, ErrorTypeUnexpectedAffectedRows, err.ErrorType())
	assert.Equal(t, "unexpected affected rows", err.message)
	assert.Equal(t, "database:unexpected affected rows", err.Error())
}

func TestNewErrorInternal(t *testing.T) {
	t.Parallel()

	err := NewErrorInternal("system", "internal server error")

	assert.Equal(t, "system", err.PackageName())
	assert.Equal(t, ErrorTypeInternal, err.ErrorType())
	assert.Equal(t, "internal server error", err.message)
	assert.Equal(t, "system:internal server error", err.Error())
}

func TestNewErrorInvalidArgument(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		pkg                  string
		message              string
		errorType            ErrorType
		field                string
		wrappedError         error
		expectedErrorMessage string
		expectedField        string
	}{
		{
			name:                 "empty field error",
			pkg:                  "account",
			message:              "invalid argument",
			errorType:            ErrorTypeInvalidArgEmpty,
			field:                "email",
			expectedErrorMessage: "account:invalid argument[email:invalid_empty]",
			expectedField:        "email",
		},
		{
			name:                 "nil field error",
			pkg:                  "feature",
			message:              "invalid input",
			errorType:            ErrorTypeInvalidArgNil,
			field:                "name",
			expectedErrorMessage: "feature:invalid input[name:invalid_nil]",
			expectedField:        "name",
		},
		{
			name:                 "format mismatch error",
			pkg:                  "validation",
			message:              "format error",
			errorType:            ErrorTypeInvalidArgNotMatchFormat,
			field:                "date",
			expectedErrorMessage: "validation:format error[date:invalid_not_match_format]",
			expectedField:        "date",
		},
		{
			name:                 "empty message error",
			pkg:                  "test",
			message:              "",
			errorType:            ErrorTypeInvalidArgEmpty,
			field:                "field",
			expectedErrorMessage: "test:[field:invalid_empty]",
			expectedField:        "field",
		},
		{
			name:                 "wrapped error",
			pkg:                  "test",
			message:              "invalid argument",
			errorType:            ErrorTypeInvalidArgEmpty,
			field:                "field",
			wrappedError:         errors.New("wrapped error"),
			expectedErrorMessage: "test:invalid argument[field:invalid_empty]: wrapped error",
			expectedField:        "field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := newBktFieldError(tt.pkg, tt.errorType, tt.message, tt.field)
			err.Wrap(tt.wrappedError)
			assert.Equal(t, tt.pkg, err.PackageName())
			assert.Equal(t, tt.errorType, err.ErrorType())
			assert.Equal(t, tt.expectedErrorMessage, err.Error())
		})
	}
}

func TestBucketeerError_Unwrap(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("original error")
	bucketeerErr := &BktError{
		packageName:  "test",
		errorType:    ErrorTypeInternal,
		message:      "test error",
		wrappedError: originalErr,
	}

	unwrapped := bucketeerErr.Unwrap()
	assert.Equal(t, originalErr, unwrapped)
	assert.ErrorIs(t, bucketeerErr, originalErr)
	assert.True(t, errors.Is(bucketeerErr, originalErr))
}

func TestBucketeerError_EmptyArgs(t *testing.T) {
	t.Parallel()

	err := NewErrorNotFound("test", "not found", "")
	assert.Equal(t, "not found", err.message)
	assert.Equal(t, "test:not found", err.Error())
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
		{ErrorTypeInvalidArgUnknown, "invalid_unknown"},
		{ErrorTypeInvalidArgEmpty, "invalid_empty"},
		{ErrorTypeInvalidArgNil, "invalid_nil"},
		{ErrorTypeInvalidArgNotMatchFormat, "invalid_not_match_format"},
	}

	for _, tt := range tests {
		t.Run(string(tt.errorType), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, string(tt.errorType))
		})
	}
}

func TestErrorWrapComplex(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("original error")
	fieldErr := NewErrorNotFound("test", "not found", "resource")
	fieldErr.Wrap(originalErr)

	invalidErr := NewErrorInvalidArgEmpty("test", "invalid argument", "field")
	invalidErr.Wrap(fieldErr)

	assert.ErrorIs(t, invalidErr, originalErr)
	assert.ErrorIs(t, invalidErr, fieldErr)

	additionalErr := errors.New("additional error")
	invalidErr.Wrap(additionalErr)

	assert.ErrorIs(t, invalidErr, additionalErr)
	assert.ErrorIs(t, invalidErr, originalErr)
	assert.ErrorIs(t, invalidErr, fieldErr)
}

func TestErrorAs(t *testing.T) {
	t.Parallel()

	originalErr := NewErrorPermissionDenied("test", "permission denied")
	fieldErr := NewErrorNotFound("test", "not found", "resource")
	fieldErr.Wrap(originalErr)

	var targetErr *BktFieldError
	if errors.As(fieldErr, &targetErr) {
		assert.Equal(t, "test", targetErr.PackageName())
		assert.Equal(t, "not found", targetErr.message)
		assert.Equal(t, "resource", targetErr.Field())
	} else {
		t.Error("Expected fieldErr to be of type *BktFieldError")
	}

	var wrappedErr *BktError
	if errors.As(fieldErr, &wrappedErr) {
		assert.Equal(t, "test", wrappedErr.PackageName())
		assert.Equal(t, "permission denied", wrappedErr.message)
		assert.Equal(t, "test:permission denied", wrappedErr.Error())
	} else {
		t.Error("Expected fieldErr to wrap originalErr")
	}
}
