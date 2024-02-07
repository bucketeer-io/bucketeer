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

package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
)

func (s *gatewayService) validateGoalEvent(ctx context.Context, id string, timeStamp int64) (string, error) {
	if err := uuid.ValidateUUID(id); err != nil {
		s.logger.Warn(
			"Failed to validate goal event id format",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", id),
			)...,
		)
		return codeInvalidID, errInvalidIDFormat
	}
	if !validateTimestamp(timeStamp, s.opts.oldestEventTimestamp, s.opts.furthestEventTimestamp) {
		s.logger.Debug(
			"Failed to validate goal event timestamp",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("id", id),
				zap.Int64("timestamp", timeStamp),
			)...,
		)
		return codeInvalidTimestamp, errInvalidTimestamp
	}
	return "", nil
}

func (s *gatewayService) validateEvaluationEvent(ctx context.Context, id string, timeStamp int64) (string, error) {
	if err := uuid.ValidateUUID(id); err != nil {
		s.logger.Warn(
			"Failed to validate evaluation event id format",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", id),
			)...,
		)
		return codeInvalidID, errInvalidIDFormat
	}
	if !validateTimestamp(timeStamp, s.opts.oldestEventTimestamp, s.opts.furthestEventTimestamp) {
		s.logger.Debug(
			"Failed to validate evaluation event timestamp",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("id", id),
				zap.Int64("timestamp", timeStamp),
			)...,
		)
		return codeInvalidTimestamp, errInvalidTimestamp
	}
	return "", nil
}

// For metrics events we don't need to validate the timestamp
func (s *gatewayService) validateMetricsEvent(ctx context.Context, id string) (string, error) {
	if err := uuid.ValidateUUID(id); err != nil {
		s.logger.Warn(
			"Failed to validate evaluation event id format",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", id),
			)...,
		)
		return codeInvalidID, errInvalidIDFormat
	}
	return "", nil
}
