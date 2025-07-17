// Copyright 2025 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package error

import "errors"

type ErrorInvalidAugment struct {
	PackageName string
	Message     string
	Metadata    map[string]string
	argument    string
	invalidType InvalidType
}

type InvalidType string

const (
	InvalidTypeEmpty          InvalidType = "empty"
	InvalidTypeNil            InvalidType = "nil"
	InvalidTypeNotMatchFormat InvalidType = "not_match_format"
)

func NewErrorInvalidAugment(pkg string, message string, invalidType InvalidType, arg string) *ErrorInvalidAugment {
	//example: account:invalid augment[account_id:empty]
	msg := pkg + ":"
	if message != "" {
		msg += message
	} else {
		msg += "invalid augment"
		if arg != "" {
			msg += "[" + arg
			if invalidType != "" {
				msg += ":" + string(invalidType)
			}
			msg += "]"
		}
	}
	messageKey := pkg + ".invalid_augment"
	if invalidType != "" {
		messageKey = messageKey + "." + string(invalidType)
	}
	//example: {"messageKey": "account.invalid_augment.empty", "field": "account_id"}
	metadata := map[string]string{
		"messageKey": messageKey,
		"field":      arg,
	}

	return &ErrorInvalidAugment{
		PackageName: pkg,
		Message:     msg,
		Metadata:    metadata,
		argument:    arg,
		invalidType: invalidType,
	}
}

func (e *ErrorInvalidAugment) Error() string {
	return e.Message
}

func (e *ErrorInvalidAugment) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorNotFound struct {
	PackageName string
	Message     string
	Metadata    map[string]string
	argument    string
}

func NewErrorNotFound(pkg string, message string, arg string) *ErrorNotFound {
	//example: account:account_id not found
	msg := pkg + ":"
	if message != "" {
		msg += message
	} else {
		msg := pkg + ":"
		if arg != "" {
			msg += arg + " "
		}
		msg += "not found"
	}

	messageKey := pkg + ".not_found"
	//example: {"messageKey": "account.not_found", "field": "account_id"}
	metadata := map[string]string{
		"messageKey": messageKey,
		"field":      arg,
	}
	return &ErrorNotFound{
		PackageName: pkg,
		argument:    arg,
		Message:     msg,
		Metadata:    metadata,
	}
}

func (e *ErrorNotFound) Error() string {
	return e.Message
}

func (e *ErrorNotFound) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorAlreadyExists struct {
	PackageName string
	Message     string
	Metadata    map[string]string
	argument    string
}

func NewErrorAlreadyExists(pkg string, message string, arg string) *ErrorAlreadyExists {
	msg := pkg + ":"
	if message != "" {
		msg += message
	} else {
		msg += "already exists"
	}

	messageKey := pkg + ".already_exists"
	metadata := map[string]string{
		"messageKey": messageKey,
		"field":      arg,
	}
	return &ErrorAlreadyExists{
		PackageName: pkg,
		Message:     msg,
		Metadata:    metadata,
		argument:    arg,
	}
}

func (e *ErrorAlreadyExists) Error() string {
	return e.Message
}

func (e *ErrorAlreadyExists) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorUnauthenticated struct {
	PackageName string
	Message     string
	Metadata    map[string]string
}

func NewErrorUnauthenticated(pkg string, message string) *ErrorUnauthenticated {
	msg := pkg + ":"
	if message != "" {
		msg += message
	} else {
		msg += "unauthenticated"
	}

	messageKey := pkg + ".unauthenticated"
	metadata := map[string]string{
		"messageKey": messageKey,
	}
	return &ErrorUnauthenticated{
		PackageName: pkg,
		Message:     msg,
		Metadata:    metadata,
	}
}

func (e *ErrorUnauthenticated) Error() string {
	return e.Message
}

func (e *ErrorUnauthenticated) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorPermissionDenied struct {
	PackageName string
	Message     string
	Metadata    map[string]string
}

func NewErrorPermissionDenied(pkg string, message string) *ErrorPermissionDenied {
	msg := pkg + ":"
	if message != "" {
		msg += message
	} else {
		msg += "permission denied"
	}

	messageKey := pkg + ".permission_denied"
	metadata := map[string]string{
		"messageKey": messageKey,
	}
	return &ErrorPermissionDenied{
		PackageName: pkg,
		Message:     msg,
		Metadata:    metadata,
	}
}

func (e *ErrorPermissionDenied) Error() string {
	return e.Message
}

func (e *ErrorPermissionDenied) Unwrap() error {
	return errors.New(e.Message)
}

type ErrorInternal struct {
	PackageName string
	Message     string
	Metadata    map[string]string
}

func NewErrorInternal(pkg string, message string) *ErrorInternal {
	msg := pkg + ":"
	if message != "" {
		msg += message
	} else {
		msg += "internal"
	}

	messageKey := pkg + ".internal"
	metadata := map[string]string{
		"messageKey": messageKey,
	}
	return &ErrorInternal{
		PackageName: pkg,
		Message:     msg,
		Metadata:    metadata,
	}
}

func (e *ErrorInternal) Error() string {
	return e.Message
}

func (e *ErrorInternal) Unwrap() error {
	return errors.New(e.Message)
}
