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

import "errors"

const (
	AccountPackageName = "account"
)

type BucketeerErrorInfo struct {
	PackageName string
	Message     string
	Metadatas   []map[string]string
	arguments   []string
}

func newErrorInfo(pkg, errorType, defaultMessage string, args ...string) BucketeerErrorInfo {
	msg := pkg + ":"
	if defaultMessage != "" {
		msg += defaultMessage
	}

	messageKey := pkg + "." + errorType
	metadatas := make([]map[string]string, 0, len(args))
	for _, arg := range args {
		if arg != "" {
			msg += ", " + arg
		}
		metadatas = append(metadatas, map[string]string{
			"messageKey": messageKey,
			"field":      arg,
		})
	}

	return BucketeerErrorInfo{
		PackageName: pkg,
		Message:     msg,
		Metadatas:   metadatas,
		arguments:   args,
	}
}

type ErrorNotFound struct {
	BucketeerErrorInfo
}

func NewErrorNotFound(pkg string, message string, args ...string) error {
	info := newErrorInfo(pkg, "not_found", message, args...)
	return &ErrorNotFound{BucketeerErrorInfo: info}
}

func (e *ErrorNotFound) Error() string {
	return e.Message
}

func (e *ErrorNotFound) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorAlreadyExists struct {
	BucketeerErrorInfo
}

func NewErrorAlreadyExists(pkg string, message string, args ...string) error {
	info := newErrorInfo(pkg, "already_exists", message, args...)
	return &ErrorAlreadyExists{BucketeerErrorInfo: info}
}

func (e *ErrorAlreadyExists) Error() string {
	return e.Message
}

func (e *ErrorAlreadyExists) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorUnauthenticated struct {
	BucketeerErrorInfo
}

func NewErrorUnauthenticated(pkg string, message string, args ...string) error {
	info := newErrorInfo(pkg, "unauthenticated", message, args...)
	return &ErrorUnauthenticated{BucketeerErrorInfo: info}
}

func (e *ErrorUnauthenticated) Error() string {
	return e.Message
}

func (e *ErrorUnauthenticated) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorPermissionDenied struct {
	BucketeerErrorInfo
}

func NewErrorPermissionDenied(pkg string, message string, args ...string) error {
	info := newErrorInfo(pkg, "permission_denied", message, args...)
	return &ErrorPermissionDenied{BucketeerErrorInfo: info}
}

func (e *ErrorPermissionDenied) Error() string {
	return e.Message
}

func (e *ErrorPermissionDenied) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorUnexpectedAffectedRows struct {
	BucketeerErrorInfo
}

func NewErrorUnexpectedAffectedRows(pkg string, message string, args ...string) error {
	info := newErrorInfo(pkg, "unexpected_affected_rows", message, args...)
	return &ErrorUnexpectedAffectedRows{BucketeerErrorInfo: info}
}

func (e *ErrorUnexpectedAffectedRows) Error() string {
	return e.Message
}

func (e *ErrorUnexpectedAffectedRows) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorInternal struct {
	BucketeerErrorInfo
}

func NewErrorInternal(pkg string, message string, args ...string) error {
	info := newErrorInfo(pkg, "internal", message, args...)
	return &ErrorInternal{BucketeerErrorInfo: info}
}

func (e *ErrorInternal) Error() string {
	return e.Message
}

func (e *ErrorInternal) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorInvalidAugment struct {
	BucketeerErrorInfo
	invalidType InvalidType
}

type InvalidType string

const (
	InvalidTypeEmpty          InvalidType = "empty"
	InvalidTypeNil            InvalidType = "nil"
	InvalidTypeNotMatchFormat InvalidType = "not_match_format"
)

func newInvalidAugmentErrorInfo(pkg, message string, invalidType InvalidType, args ...string) BucketeerErrorInfo {
	messageKey := pkg + ".invalid_augment"
	if invalidType != "" {
		messageKey = messageKey + "." + string(invalidType)
	}
	msg := pkg + ":"
	if message != "" {
		msg += message
	} else {
		msg += "invalid augment"
	}
	metadatas := make([]map[string]string, 0, len(args))
	for _, arg := range args {
		if arg != "" {
			msg += "[" + arg
			if invalidType != "" {
				msg += ":" + string(invalidType)
			}
			msg += "]"
		}
		metadatas = append(metadatas, map[string]string{
			"messageKey": messageKey,
			"field":      arg,
		})
	}

	return BucketeerErrorInfo{
		PackageName: pkg,
		Message:     msg,
		Metadatas:   metadatas,
		arguments:   args,
	}
}

func NewErrorInvalidAugment(pkg string, message string, invalidType InvalidType, args ...string) error {
	info := newInvalidAugmentErrorInfo(pkg, message, invalidType, args...)
	return &ErrorInvalidAugment{
		BucketeerErrorInfo: info,
		invalidType:        invalidType,
	}
}

func (e *ErrorInvalidAugment) Error() string {
	return e.Message
}

func (e *ErrorInvalidAugment) Unwrap() error {
	return errors.New(e.Message)
}
