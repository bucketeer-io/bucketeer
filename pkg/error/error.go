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

const (
	AccountPackageName = "account"
)

type BucketeerError interface {
	error
	PackageName() string
	Message() string
	Metadatas() []map[string]string
	AddMetadata(metadatas ...map[string]string)
}

type bucketeerError struct {
	packageName string
	message     string
	metadatas   []map[string]string
	err         error
}

func (e *bucketeerError) PackageName() string            { return e.packageName }
func (e *bucketeerError) Message() string                { return e.message }
func (e *bucketeerError) Metadatas() []map[string]string { return e.metadatas }
func (e *bucketeerError) Error() string                  { return e.message }
func (e *bucketeerError) Unwrap() error                  { return e.err }
func (e *bucketeerError) AddMetadata(metadatas ...map[string]string) {
	e.metadatas = append(e.metadatas, metadatas...)
}

func newError(pkg, errorType, defaultMessage string, err error, args ...string) *bucketeerError {
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
	return &bucketeerError{
		packageName: pkg,
		message:     msg,
		metadatas:   metadatas,
		err:         err,
	}
}

type ErrorNotFound struct {
	*bucketeerError
}

func (e *ErrorNotFound) Is(target error) bool {
	_, ok := target.(*ErrorNotFound)
	return ok
}

func NewErrorNotFound(pkg string, message string, err error, args ...string) error {
	return &ErrorNotFound{newError(pkg, "not_found", message, err, args...)}
}

type ErrorAlreadyExists struct {
	*bucketeerError
}

func (e *ErrorAlreadyExists) Is(target error) bool {
	_, ok := target.(*ErrorAlreadyExists)
	return ok
}

func NewErrorAlreadyExists(pkg string, message string, err error, args ...string) error {
	return &ErrorAlreadyExists{newError(pkg, "already_exists", message, err, args...)}
}

type ErrorUnauthenticated struct {
	*bucketeerError
}

func (e *ErrorUnauthenticated) Is(target error) bool {
	_, ok := target.(*ErrorUnauthenticated)
	return ok
}

func NewErrorUnauthenticated(pkg string, message string, err error, args ...string) error {
	return &ErrorUnauthenticated{newError(pkg, "unauthenticated", message, err, args...)}
}

type ErrorPermissionDenied struct {
	*bucketeerError
}

func (e *ErrorPermissionDenied) Is(target error) bool {
	_, ok := target.(*ErrorPermissionDenied)
	return ok
}

func NewErrorPermissionDenied(pkg string, message string, err error, args ...string) error {
	return &ErrorPermissionDenied{newError(pkg, "permission_denied", message, err, args...)}
}

type ErrorUnexpectedAffectedRows struct {
	*bucketeerError
}

func (e *ErrorUnexpectedAffectedRows) Is(target error) bool {
	_, ok := target.(*ErrorUnexpectedAffectedRows)
	return ok
}

func NewErrorUnexpectedAffectedRows(pkg string, message string, err error, args ...string) error {
	return &ErrorUnexpectedAffectedRows{newError(pkg, "unexpected_affected_rows", message, err, args...)}
}

type ErrorInternal struct {
	*bucketeerError
}

func (e *ErrorInternal) Is(target error) bool {
	_, ok := target.(*ErrorInternal)
	return ok
}

func NewErrorInternal(pkg string, message string, err error, args ...string) error {
	return &ErrorInternal{newError(pkg, "internal", message, err, args...)}
}

type ErrorInvalidAugment struct {
	*bucketeerError
	invalidType InvalidType
}

func (e *ErrorInvalidAugment) Is(target error) bool {
	_, ok := target.(*ErrorInvalidAugment)
	return ok
}

type InvalidType string

const (
	InvalidTypeEmpty          InvalidType = "empty"
	InvalidTypeNil            InvalidType = "nil"
	InvalidTypeNotMatchFormat InvalidType = "not_match_format"
)

func NewErrorInvalidAugment(pkg string, message string, invalidType InvalidType, err error, args ...string) error {
	return &ErrorInvalidAugment{
		bucketeerError: newInvalidAugmentErrorBase(pkg, message, invalidType, err, args...),
		invalidType:    invalidType,
	}
}

func newInvalidAugmentErrorBase(pkg, message string, invalidType InvalidType, err error, args ...string) *bucketeerError {
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
	return &bucketeerError{
		packageName: pkg,
		message:     msg,
		metadatas:   metadatas,
		err:         err,
	}
}
