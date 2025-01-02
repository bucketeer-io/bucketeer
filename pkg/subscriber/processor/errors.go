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

package processor

import (
	"errors"
	"fmt"
)

var (
	errAuditLogInvalidConfig                     = errors.New("auditLog: invalid config")
	ErrSegmentInvalidConfig                      = errors.New("segment: invalid config")
	ErrSegmentInUse                              = errors.New("segment: segment is in use")
	ErrSegmentExceededMaxUserIDLength            = fmt.Errorf("segment: max user id length allowed is %d", maxUserIDLength) //nolint:lll
	ErrUserEventInvalidConfig                    = errors.New("user event: invalid config")
	ErrEvaluationCountInvalidConfig              = errors.New("evaluation count: invalid config")
	ErrEventsDWHPersisterInvalidConfig           = errors.New("eventpersister: invalid config")
	ErrEventsOPSPersisterInvalidConfig           = errors.New("eventpersister: invalid config")
	ErrExperimentNotFound                        = errors.New("eventpersister: experiment not found")
	ErrReasonNil                                 = errors.New("eventpersister: reason is nil")
	ErrEvaluationsAreEmpty                       = errors.New("eventpersister: evaluations are empty")
	ErrEvaluationEventIssuedAfterExperimentEnded = errors.New("eventpersister: evaluation event issued after experiment ended") //nolint:lll
	ErrFailedToEvaluateUser                      = errors.New("eventpersister: failed to evaluate user")
	ErrAutoOpsRuleNotFound                       = errors.New("eventpersister: auto ops rule not found")
	ErrFeatureEmptyList                          = errors.New("eventpersister: list feature returned empty")
	ErrFeatureVersionNotFound                    = errors.New("eventpersister: feature version not found")
	ErrUnknownEvent                              = errors.New("metricsevent persister: unknown metrics event")
	ErrInvalidDuration                           = errors.New("metricsevent persister: invalid duration")
	ErrUnknownApiId                              = errors.New("metricsevent persister: unknown api id")
)
