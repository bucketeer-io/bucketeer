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
			expectedErrorMessage: "account:invalid argument[email:InvalidArgumentEmptyError]",
			expectedField:        "email",
		},
		{
			name:                 "nil field error",
			pkg:                  "feature",
			message:              "invalid input",
			errorType:            ErrorTypeInvalidArgNil,
			field:                "name",
			expectedErrorMessage: "feature:invalid input[name:InvalidArgumentNilError]",
			expectedField:        "name",
		},
		{
			name:                 "format mismatch error",
			pkg:                  "validation",
			message:              "format error",
			errorType:            ErrorTypeInvalidArgNotMatchFormat,
			field:                "date",
			expectedErrorMessage: "validation:format error[date:InvalidArgumentNotMatchFormatError]",
			expectedField:        "date",
		},
		{
			name:                 "empty message error",
			pkg:                  "test",
			message:              "",
			errorType:            ErrorTypeInvalidArgEmpty,
			field:                "field",
			expectedErrorMessage: "test:[field:InvalidArgumentEmptyError]",
			expectedField:        "field",
		},
		{
			name:                 "wrapped error",
			pkg:                  "test",
			message:              "invalid argument",
			errorType:            ErrorTypeInvalidArgEmpty,
			field:                "field",
			wrappedError:         errors.New("wrapped error"),
			expectedErrorMessage: "test:invalid argument[field:InvalidArgumentEmptyError]: wrapped error",
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
		{ErrorTypeNotFound, "NotFoundError"},
		{ErrorTypeAlreadyExists, "AlreadyExistsError"},
		{ErrorTypeUnauthenticated, "UnauthenticatedError"},
		{ErrorTypePermissionDenied, "PermissionDeniedError"},
		{ErrorTypeUnexpectedAffectedRows, "UnexpectedAffectedRowsError"},
		{ErrorTypeInternal, "InternalServerError"},
		{ErrorTypeInvalidArgUnknown, "InvalidArgumentUnknownError"},
		{ErrorTypeInvalidArgEmpty, "InvalidArgumentEmptyError"},
		{ErrorTypeInvalidArgNil, "InvalidArgumentNilError"},
		{ErrorTypeInvalidArgNotMatchFormat, "InvalidArgumentNotMatchFormatError"},
	}

	for _, tt := range tests {
		t.Run(string(tt.errorType), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, string(tt.errorType))
		})
	}
}

func TestMessageKey(t *testing.T) {
	t.Parallel()

	t.Run("default returns error type", func(t *testing.T) {
		t.Parallel()
		err := NewErrorFailedPrecondition("test", "precondition failed")
		assert.Equal(t, "FailedPreconditionError", err.MessageKey())
	})

	t.Run("override replaces default", func(t *testing.T) {
		t.Parallel()
		err := NewErrorFailedPrecondition("test", "comment required").
			WithMessageKey("CommentRequiredForUpdating")
		assert.Equal(t, "CommentRequiredForUpdating", err.MessageKey())
		assert.Equal(t, ErrorTypeFailedPrecondition, err.ErrorType())
	})

	t.Run("empty override falls back to default", func(t *testing.T) {
		t.Parallel()
		err := NewErrorFailedPrecondition("test", "generic").
			WithMessageKey("")
		assert.Equal(t, "FailedPreconditionError", err.MessageKey())
	})
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

	var targetErr *BktError
	if errors.As(fieldErr, &targetErr) {
		assert.Equal(t, "test", targetErr.PackageName())
		assert.Equal(t, "not found", targetErr.message)
		assert.Equal(t, "resource", targetErr.field)
		assert.Equal(t, "test:not found, resource: test:permission denied", targetErr.Error())
	} else {
		t.Error("Expected fieldErr to be of type *BktError")
	}

	// Check that wrapped error can be extracted
	var wrappedErr *BktError
	if errors.As(fieldErr.Unwrap(), &wrappedErr) {
		assert.Equal(t, "test", wrappedErr.PackageName())
		assert.Equal(t, "permission denied", wrappedErr.message)
		assert.Equal(t, "test:permission denied", wrappedErr.Error())
	} else {
		t.Error("Expected fieldErr to wrap originalErr")
	}
}

func TestNewErrorVariationInUse(t *testing.T) {
	t.Parallel()

	err := NewErrorVariationInUse(
		FeaturePackageName,
		ErrorTypeVariationInUseByPrerequisite,
		"feature: variation variation-1 is used as a prerequisite by feature feature-2",
		map[string]string{"featureId": "feature-2", "featureName": "Flag B"},
	)

	assert.Equal(t, FeaturePackageName, err.PackageName())
	assert.Equal(t, ErrorTypeVariationInUseByPrerequisite, err.ErrorType())
	assert.Equal(t, "VariationInUseByPrerequisiteError", err.MessageKey())
	assert.Equal(
		t,
		"feature:feature: variation variation-1 is used as a prerequisite by feature feature-2",
		err.Error(),
	)
	assert.Equal(
		t,
		map[string]string{"featureId": "feature-2", "featureName": "Flag B"},
		err.EmbeddedKeyValues(),
	)
}

func TestNewErrorVariationInUse_CopiesKeyValues(t *testing.T) {
	t.Parallel()

	keyValues := map[string]string{"featureId": "feature-2"}
	err := NewErrorVariationInUse(
		FeaturePackageName,
		ErrorTypeVariationInUseByPrerequisite,
		"feature: variation in use",
		keyValues,
	)
	keyValues["featureId"] = "mutated"

	assert.Equal(t, "feature-2", err.EmbeddedKeyValues()["featureId"])
}

func TestAsVariationInUseError(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		err      error
		expected bool
	}{
		{
			desc:     "false: nil",
			err:      nil,
			expected: false,
		},
		{
			desc:     "false: not a BktError",
			err:      errors.New("something else"),
			expected: false,
		},
		{
			desc:     "false: another BktError type",
			err:      NewErrorFailedPrecondition(FeaturePackageName, "failed precondition"),
			expected: false,
		},
		{
			desc: "true: off variation",
			err: NewErrorVariationInUse(
				FeaturePackageName, ErrorTypeVariationInUseByOffVariation, "in use", nil,
			),
			expected: true,
		},
		{
			desc: "true: individual targeting",
			err: NewErrorVariationInUse(
				FeaturePackageName, ErrorTypeVariationInUseByIndividualTarget, "in use", nil,
			),
			expected: true,
		},
		{
			desc: "true: feature flag rule",
			err: NewErrorVariationInUse(
				FeaturePackageName, ErrorTypeVariationInUseByFeatureFlagRule, "in use", nil,
			),
			expected: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			bktErr, ok := AsVariationInUseError(p.err)
			assert.Equal(t, p.expected, ok)
			if !p.expected {
				assert.Nil(t, bktErr)
			}
		})
	}
}
