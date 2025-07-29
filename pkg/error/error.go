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

type BucketeerError struct {
	packageName string
	errorType   ErrorType
	message     string
	metadatas   []map[string]string
	err         error
}

type ErrorType string

const (
	ErrorTypeNotFound               ErrorType = "not_found"
	ErrorTypeAlreadyExists          ErrorType = "already_exists"
	ErrorTypeUnauthenticated        ErrorType = "unauthenticated"
	ErrorTypePermissionDenied       ErrorType = "permission_denied"
	ErrorTypeUnexpectedAffectedRows ErrorType = "unexpected_affected_rows"
	ErrorTypeInternal               ErrorType = "internal"
	ErrorTypeInvalidAugment         ErrorType = "invalid_augment"
)

func (e *BucketeerError) PackageName() string            { return e.packageName }
func (e *BucketeerError) ErrorType() ErrorType           { return e.errorType }
func (e *BucketeerError) Message() string                { return e.message }
func (e *BucketeerError) Metadatas() []map[string]string { return e.metadatas }
func (e *BucketeerError) Error() string                  { return e.message }
func (e *BucketeerError) Unwrap() error                  { return e.err }
func (e *BucketeerError) AddMetadata(metadatas ...map[string]string) {
	e.metadatas = append(e.metadatas, metadatas...)
}
func (e *BucketeerError) JoinError(err error) {
	e.err = errors.Join(e.err, err)
}

func newError(pkg string, errorType ErrorType, defaultMessage string, args ...string) *BucketeerError {
	msg := pkg + ":"
	if defaultMessage != "" {
		msg += defaultMessage
	}

	messageKey := pkg + "." + string(errorType)
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

	// example: NotFound {
	// 	packageName: "account",
	// 	message:     "account not found, user_id",
	//  metadatas:  []map[string]string{
	// 		{
	// 			"messageKey": "account.not_found",
	// 			"field":      "user_id",
	// 		},
	// 	},
	//}
	return &BucketeerError{
		packageName: pkg,
		errorType:   errorType,
		message:     msg,
		metadatas:   metadatas,
	}
}

func NewErrorNotFound(pkg string, message string, args ...string) *BucketeerError {
	return newError(pkg, ErrorTypeNotFound, message, args...)
}

func NewErrorAlreadyExists(pkg string, message string, args ...string) *BucketeerError {
	return newError(pkg, ErrorTypeAlreadyExists, message, args...)
}

func NewErrorUnauthenticated(pkg string, message string, args ...string) *BucketeerError {
	return newError(pkg, ErrorTypeUnauthenticated, message, args...)
}

func NewErrorPermissionDenied(pkg string, message string, args ...string) *BucketeerError {
	return newError(pkg, ErrorTypePermissionDenied, message, args...)
}

func NewErrorUnexpectedAffectedRows(pkg string, message string, args ...string) *BucketeerError {
	return newError(pkg, ErrorTypeUnexpectedAffectedRows, message, args...)
}

func NewErrorInternal(pkg string, message string, args ...string) *BucketeerError {
	return newError(pkg, ErrorTypeInternal, message, args...)
}

type InvalidType string

const (
	invalidTypeUnknown        InvalidType = "unknown"
	InvalidTypeEmpty          InvalidType = "empty"
	InvalidTypeNil            InvalidType = "nil"
	InvalidTypeNotMatchFormat InvalidType = "not_match_format"
)

func NewErrorInvalidAugment(pkg string, message string, invalidType InvalidType, args ...string) *BucketeerError {
	return newInvalidAugmentErrorBase(pkg, message, invalidType, args...)
}

func newInvalidAugmentErrorBase(pkg, message string, invalidType InvalidType, args ...string) *BucketeerError {
	errorType := ErrorTypeInvalidAugment
	messageKey := pkg + "." + string(errorType)
	if invalidType == "" {
		invalidType = invalidTypeUnknown
	}
	messageKey = messageKey + "." + string(invalidType)
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

	// example: two invalid extensions found {
	// 	packageName: "account",
	// 	message:     "account invalid augment, user_id[empty], name[empty]",
	//  metadatas:  []map[string]string{
	// 		{
	// 			"messageKey": "account.invalid_augment.empty",
	// 			"field":      "user_id",
	// 		},
	// 		{
	// 			"messageKey": "account.invalid_augment.empty",
	// 			"field":      "name",
	// 		},
	// 	},
	//}
	return &BucketeerError{
		packageName: pkg,
		errorType:   errorType,
		message:     msg,
		metadatas:   metadatas,
	}
}
