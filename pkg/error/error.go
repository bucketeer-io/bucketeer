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
)

const (
	AccountPackageName = "account"
)

type ErrorType string

const (
	ErrorTypeNotFound               ErrorType = "not_found"
	ErrorTypeAlreadyExists          ErrorType = "already_exists"
	ErrorTypeUnauthenticated        ErrorType = "unauthenticated"
	ErrorTypePermissionDenied       ErrorType = "permission_denied"
	ErrorTypeUnexpectedAffectedRows ErrorType = "unexpected_affected_rows"
	ErrorTypeInternal               ErrorType = "internal"
	ErrorTypeInvalidArgument        ErrorType = "invalid_argument"
)

// Base error - no field needed
type BktError struct {
	packageName  string
	errorType    ErrorType
	message      string
	wrappedError error
}

func (e *BktError) PackageName() string  { return e.packageName }
func (e *BktError) ErrorType() ErrorType { return e.errorType }
func (e *BktError) Error() string {
	msg := e.packageName + ":" + e.message
	if e.wrappedError != nil {
		return fmt.Sprintf("%s: %v", msg, e.wrappedError)
	}
	return msg
}

func (e *BktError) Unwrap() error { return e.wrappedError }
func (e *BktError) Wrap(err error) {
	e.wrappedError = errors.Join(e.wrappedError, err)
}

// BktFieldError represents an error with a specific field
type BktFieldError struct {
	*BktError
	field string
}

func (e *BktFieldError) Error() string {
	msg := e.packageName + ":" + e.message
	if e.field != "" {
		msg += ", " + e.field
	}
	if e.wrappedError != nil {
		return fmt.Sprintf("%s: %v", msg, e.wrappedError)
	}
	return msg
}

func (e *BktFieldError) Field() string {
	return e.field
}

type BktInvalidError struct {
	*BktFieldError
	invalidType InvalidType
}

func (e *BktInvalidError) Error() string {
	if e.message == "" {
		e.message = "invalid argument"
	}

	msg := e.packageName + ":" + e.message
	if e.field != "" {
		if e.invalidType != "" {
			msg += "[" + e.field + ":" + string(e.invalidType) + "]"
		} else {
			msg += ", " + e.field
		}
	}
	if e.wrappedError != nil {
		return fmt.Sprintf("%s: %v", msg, e.wrappedError)
	}
	return msg
}

func (e *BktInvalidError) InvalidType() InvalidType {
	return e.invalidType
}

func newBktError(pkg string, errorType ErrorType, message string) *BktError {
	return &BktError{
		packageName: pkg,
		errorType:   errorType,
		message:     message,
	}
}

func newBktFieldError(pkg string, errorType ErrorType, message string, field string) *BktFieldError {
	return &BktFieldError{
		BktError: &BktError{
			packageName: pkg,
			errorType:   errorType,
			message:     message,
		},
		field: field,
	}
}

func NewErrorNotFound(pkg string, message string, field string) *BktFieldError {
	return newBktFieldError(pkg, ErrorTypeNotFound, message, field)
}

func NewErrorAlreadyExists(pkg string, message string) *BktError {
	return newBktError(pkg, ErrorTypeAlreadyExists, message)
}

func NewErrorUnauthenticated(pkg string, message string) *BktError {
	return newBktError(pkg, ErrorTypeUnauthenticated, message)
}

func NewErrorPermissionDenied(pkg string, message string) *BktError {
	return newBktError(pkg, ErrorTypePermissionDenied, message)
}

func NewErrorUnexpectedAffectedRows(pkg string, message string) *BktError {
	return newBktError(pkg, ErrorTypeUnexpectedAffectedRows, message)
}

func NewErrorInternal(pkg string, message string) *BktError {
	return newBktError(pkg, ErrorTypeInternal, message)
}

type InvalidType string

const (
	InvalidTypeUnknown        InvalidType = "unknown"
	InvalidTypeEmpty          InvalidType = "empty"
	InvalidTypeNil            InvalidType = "nil"
	InvalidTypeNotMatchFormat InvalidType = "not_match_format"
)

func NewErrorInvalidArgument(pkg string, message string, invalidType InvalidType, field string) *BktInvalidError {
	return &BktInvalidError{
		BktFieldError: &BktFieldError{
			BktError: &BktError{
				packageName: pkg,
				errorType:   ErrorTypeInvalidArgument,
				message:     message,
			},
			field: field,
		},
		invalidType: invalidType,
	}
}
