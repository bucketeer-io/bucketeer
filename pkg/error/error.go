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
	Pkg         string
	InvalidArgs []string
}

func NewErrorInvalidAugment(pkg string, invalidArgs []string) *ErrorInvalidAugment {
	return &ErrorInvalidAugment{
		Pkg:         pkg,
		InvalidArgs: invalidArgs,
	}
}

func (e *ErrorInvalidAugment) Error() string {
	msg := e.Pkg + ":invalid augment arguments"
	for _, arg := range e.InvalidArgs {
		msg += " " + arg
	}
	return msg
}

func (e *ErrorInvalidAugment) Unwrap() error {
	errMsg := e.Error()
	return errors.New(errMsg)
}
