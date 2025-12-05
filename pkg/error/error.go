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
	"strings"
)

const (
	AccountPackageName      = "account"
	FeaturePackageName      = "feature"
	NotificationPackageName = "notification"
	PushPackageName         = "push"
	TagPackageName          = "tag"
	EventCounterPackageName = "eventcounter"
	EnvironmentPackageName  = "environment"
	AuditlogPackageName     = "auditlog"
	AutoopsPackageName      = "autoops"
	CoderefPackageName      = "coderef"
	TeamPackageName         = "team"
	ExperimentPackageName   = "experiment"
	AuthPackageName         = "auth"

	invalidPrefix = "Invalid"
)

// ErrorType is also used as the message key.
type ErrorType string

const (
	ErrorTypeNotFound                 ErrorType = "NotFoundError"
	ErrorTypeAlreadyExists            ErrorType = "AlreadyExistsError"
	ErrorTypeUnauthenticated          ErrorType = "UnauthenticatedError"
	ErrorTypePermissionDenied         ErrorType = "PermissionDenied"
	ErrorTypeUnexpectedAffectedRows   ErrorType = "UnexpectedAffectedRows"
	ErrorTypeInternal                 ErrorType = "InternalServerError"
	ErrorTypeFailedPrecondition       ErrorType = "FailedPreconditionError"
	ErrorTypeUnavailable              ErrorType = "UnavailableError"
	ErrorTypeAborted                  ErrorType = "AbortedError"
	ErrorTypeInvalidArgUnknown        ErrorType = "InvalidArgumentUnknownError"
	ErrorTypeInvalidArgEmpty          ErrorType = "InvalidArgumentEmptyError"
	ErrorTypeInvalidArgNil            ErrorType = "InvalidArgumentNilError"
	ErrorTypeInvalidArgNotMatchFormat ErrorType = "InvalidArgumentNotMatchFormatError"
	ErrorTypeInvalidArgDuplicated     ErrorType = "InvalidArgumentDuplicatedError"
)

type BktError struct {
	packageName  string
	errorType    ErrorType
	message      string
	wrappedError error
	field        string // optional

	embeddedKeyValues map[string]string
}

func (e *BktError) PackageName() string  { return e.packageName }
func (e *BktError) ErrorType() ErrorType { return e.errorType }

func (e *BktError) MessageKey() string                   { return string(e.errorType) }
func (e *BktError) EmbeddedKeyValues() map[string]string { return e.embeddedKeyValues }

func (e *BktError) Error() string {
	msg := fmt.Sprintf("%s:%s", e.packageName, e.message)
	if e.field != "" {
		if strings.HasPrefix(string(e.errorType), invalidPrefix) {
			msg += fmt.Sprintf("[%s:%s]", e.field, e.errorType)
		} else {
			msg += fmt.Sprintf(", %s", e.field)
		}
	}
	if e.wrappedError != nil {
		return fmt.Sprintf("%s: %v", msg, e.wrappedError)
	}
	return msg
}

func (e *BktError) Unwrap() error { return e.wrappedError }
func (e *BktError) Wrap(err error) {
	e.wrappedError = errors.Join(e.wrappedError, err)
}

func newBktError(pkg string, errorType ErrorType, message string) *BktError {
	return &BktError{
		packageName: pkg,
		errorType:   errorType,
		message:     message,
	}
}

func newBktFieldError(pkg string, errorType ErrorType, message string, field string) *BktError {
	return &BktError{
		packageName: pkg,
		errorType:   errorType,
		message:     message,
		field:       field,
		embeddedKeyValues: map[string]string{
			"field": field,
		},
	}
}

func NewErrorNotFound(pkg string, message string, field string) *BktError {
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

func NewErrorFailedPrecondition(pkg string, message string) *BktError {
	return newBktError(pkg, ErrorTypeFailedPrecondition, message)
}

func NewErrorUnavailable(pkg string, message string) *BktError {
	return newBktError(pkg, ErrorTypeUnavailable, message)
}

func NewErrorAborted(pkg string, message string) *BktError {
	return newBktError(pkg, ErrorTypeAborted, message)
}

func NewErrorInvalidArgUnknown(pkg string, message string, field string) *BktError {
	return newBktFieldError(pkg, ErrorTypeInvalidArgUnknown, message, field)
}

func NewErrorInvalidArgEmpty(pkg string, message string, field string) *BktError {
	return newBktFieldError(pkg, ErrorTypeInvalidArgEmpty, message, field)
}

func NewErrorInvalidArgNil(pkg string, message string, field string) *BktError {
	return newBktFieldError(pkg, ErrorTypeInvalidArgNil, message, field)
}

func NewErrorInvalidArgNotMatchFormat(pkg string, message string, field string) *BktError {
	return newBktFieldError(pkg, ErrorTypeInvalidArgNotMatchFormat, message, field)
}

func NewErrorInvalidArgDuplicated(pkg string, message string, field string) *BktError {
	return newBktFieldError(pkg, ErrorTypeInvalidArgDuplicated, message, field)
}
