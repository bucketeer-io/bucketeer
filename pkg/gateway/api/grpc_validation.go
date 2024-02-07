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
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

var (
	errInvalidIDFormat  = errors.New("gateway: invalid event id format")
	errInvalidTimestamp = errors.New("gateway: invalid event timestamp")
	errUnmarshalFailed  = errors.New("gateway: failed to unmarshal event")
)

type eventValidator interface {
	validate(ctx context.Context) (string, error)
}

type eventEvaluationValidator struct {
	event                     *eventproto.Event
	oldestTimestampDuration   time.Duration
	furthestTimestampDuration time.Duration
	logger                    *zap.Logger
}

type eventGoalValidator struct {
	event                     *eventproto.Event
	oldestTimestampDuration   time.Duration
	furthestTimestampDuration time.Duration
	logger                    *zap.Logger
}

type eventMetricsValidator struct {
	event                     *eventproto.Event
	oldestTimestampDuration   time.Duration
	furthestTimestampDuration time.Duration
	logger                    *zap.Logger
}

func newEventValidator(
	event *eventproto.Event,
	oldestTimestampDuration, furthestTimestampDuration time.Duration,
	logger *zap.Logger,
) eventValidator {
	if ptypes.Is(event.Event, grpcGoalEvent) {
		return &eventGoalValidator{
			event:                     event,
			oldestTimestampDuration:   oldestTimestampDuration,
			furthestTimestampDuration: furthestTimestampDuration,
			logger:                    logger,
		}
	}
	if ptypes.Is(event.Event, grpcEvaluationEvent) {
		return &eventEvaluationValidator{
			event:                     event,
			oldestTimestampDuration:   oldestTimestampDuration,
			furthestTimestampDuration: furthestTimestampDuration,
			logger:                    logger,
		}
	}
	if ptypes.Is(event.Event, grpcMetricsEvent) {
		return &eventMetricsValidator{
			event:                     event,
			oldestTimestampDuration:   oldestTimestampDuration,
			furthestTimestampDuration: furthestTimestampDuration,
			logger:                    logger,
		}
	}
	return nil
}

func (v *eventGoalValidator) validate(ctx context.Context) (string, error) {
	if err := uuid.ValidateUUID(v.event.Id); err != nil {
		v.logger.Warn(
			"Failed to validate goal event id format",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", v.event.Id),
			)...,
		)
		return codeInvalidID, errInvalidIDFormat
	}
	ev, err := v.unmarshal(ctx)
	if err != nil {
		return codeUnmarshalFailed, errUnmarshalFailed
	}
	if !validateTimestamp(ev.Timestamp, v.oldestTimestampDuration, v.furthestTimestampDuration) {
		v.logger.Debug(
			"Failed to validate goal event timestamp",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("id", v.event.Id),
				zap.Int64("timestamp", ev.Timestamp),
			)...,
		)
		return codeInvalidTimestamp, errInvalidTimestamp
	}
	return "", nil
}

func (v *eventGoalValidator) unmarshal(ctx context.Context) (*eventproto.GoalEvent, error) {
	ev := &eventproto.GoalEvent{}
	if err := ptypes.UnmarshalAny(v.event.Event, ev); err != nil {
		v.logger.Error(
			"Failed to extract goal event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", v.event.Id),
			)...,
		)
		return nil, err
	}
	return ev, nil
}

func (v *eventEvaluationValidator) validate(ctx context.Context) (string, error) {
	if err := uuid.ValidateUUID(v.event.Id); err != nil {
		v.logger.Warn(
			"Failed to validate evaluation event id format",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", v.event.Id),
			)...,
		)
		return codeInvalidID, errInvalidIDFormat
	}
	ev, err := v.unmarshal(ctx)
	if err != nil {
		return codeUnmarshalFailed, errUnmarshalFailed
	}
	if !validateTimestamp(ev.Timestamp, v.oldestTimestampDuration, v.furthestTimestampDuration) {
		v.logger.Debug(
			"Failed to validate evaluation event timestamp",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("id", v.event.Id),
				zap.Int64("timestamp", ev.Timestamp),
			)...,
		)
		return codeInvalidTimestamp, errInvalidTimestamp
	}
	return "", nil
}

func (v *eventEvaluationValidator) unmarshal(ctx context.Context) (*eventproto.EvaluationEvent, error) {
	ev := &eventproto.EvaluationEvent{}
	if err := ptypes.UnmarshalAny(v.event.Event, ev); err != nil {
		v.logger.Error(
			"Failed to extract evaluation event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", v.event.Id),
			)...,
		)
		return nil, err
	}
	return ev, nil
}

// For metrics events we don't need to validate the timestamp
func (v *eventMetricsValidator) validate(ctx context.Context) (string, error) {
	if err := uuid.ValidateUUID(v.event.Id); err != nil {
		v.logger.Warn(
			"Failed to validate metrics event id format",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", v.event.Id),
			)...,
		)
		return codeInvalidID, errInvalidIDFormat
	}
	_, err := v.unmarshal(ctx)
	if err != nil {
		return codeUnmarshalFailed, errUnmarshalFailed
	}
	return "", nil
}

func (v *eventMetricsValidator) unmarshal(ctx context.Context) (*eventproto.MetricsEvent, error) {
	ev := &eventproto.MetricsEvent{}
	if err := ptypes.UnmarshalAny(v.event.Event, ev); err != nil {
		v.logger.Error(
			"Failed to extract metrics event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", v.event.Id),
			)...,
		)
		return nil, err
	}
	return ev, nil
}

// validateTimestamp limits date range of the given timestamp
// because we can't stream data outside the allowed bounds into a persistent datastore.
func validateTimestamp(
	timestamp int64,
	oldestTimestampDuration, furthestTimestampDuration time.Duration,
) bool {
	given := time.Unix(timestamp, 0)
	maxPast := time.Now().Add(-oldestTimestampDuration)
	if given.Before(maxPast) {
		return false
	}
	maxFuture := time.Now().Add(furthestTimestampDuration)
	return !given.After(maxFuture)
}
