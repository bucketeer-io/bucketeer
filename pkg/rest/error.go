// Copyright 2024 The Bucketeer Authors.
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

package rest

import (
	"errors"
)

type status struct {
	code int
	err  error
}

type errStatus interface {
	GetErrMessage() string
	GetStatusCode() int
}

func NewErrStatus(code int, msg string) error {
	s := &status{
		code: code,
		err:  errors.New(msg),
	}
	return s
}

func (s *status) Error() string {
	return s.err.Error()
}

func (s *status) GetErrMessage() string {
	return s.err.Error()
}

func (s *status) GetStatusCode() int {
	return s.code
}

func convertToErrStatus(err error) (errStatus, bool) {
	s, ok := err.(errStatus)
	if !ok {
		return nil, false
	}
	return s, true
}
