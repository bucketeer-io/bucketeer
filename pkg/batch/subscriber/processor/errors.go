//  Copyright 2024 The Bucketeer Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package processor

import (
	"errors"
	"fmt"
)

var (
	errSegmentInvalidConfig           = errors.New("segment: invalid config")
	errSegmentInUse                   = errors.New("segment: segment is in use")
	errSegmentExceededMaxUserIDLength = fmt.Errorf("segment: max user id length allowed is %d", maxUserIDLength)
	errUserEventInvalidConfig         = errors.New("user event: invalid config")
	errEvaluationCountInvalidConfig   = errors.New("evaluation count: invalid config")
)
