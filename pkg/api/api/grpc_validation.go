// Copyright 2026 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	errInvalidIDFormat  = errors.New("gateway: invalid event id format")
	errInvalidTimestamp = errors.New("gateway: invalid event timestamp")
	errUnmarshalFailed  = errors.New("gateway: failed to unmarshal event")
	errEmptyFeatureID   = errors.New("gateway: feature_id is empty")
	errEmptyUserID      = errors.New("gateway: user_id is empty")
	errEmptyVariationID = errors.New("gateway: variation_id is empty")
	errEmptyGoalID      = errors.New("gateway: goal_id is empty")
	errNilReason        = errors.New("gateway: reason is nil")
)

type eventValidator interface {
	validate(ctx context.Context) (string, error)
}

type eventEvaluationValidator struct {
	event                     *eventproto.Event
	oldestTimestampDuration   time.Duration
	furthestTimestampDuration time.Duration
	logger                    *zap.Logger
	// lastUnmarshaledEvent is set on successful unmarshal for reuse (e.g., metric reporting).
	lastUnmarshaledEvent *eventproto.EvaluationEvent
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
	ev, err := v.unmarshal(ctx)
	if err != nil {
		return codeUnmarshalFailed, errUnmarshalFailed
	}

	if err := uuid.ValidateUUID(v.event.Id); err != nil {
		v.logger.Warn(
			"Failed to validate goal event id format",
			append(v.buildGoalEventLogFields(ctx, ev), zap.Error(err))...,
		)
		return codeInvalidID, errInvalidIDFormat
	}

	// validate required fields
	if ev.GoalId == "" {
		v.logger.Debug(
			"Empty goal_id",
			v.buildGoalEventLogFields(ctx, ev)...,
		)
		return codeEmptyField, errEmptyGoalID
	}
	if (ev.User == nil || (ev.User != nil && ev.User.Id == "")) && ev.UserId == "" {
		v.logger.Debug(
			"Empty user_id",
			v.buildGoalEventLogFields(ctx, ev)...,
		)
		return codeEmptyField, errEmptyUserID
	}

	if !validateTimestamp(ev.Timestamp, v.oldestTimestampDuration, v.furthestTimestampDuration) {
		v.logger.Debug(
			"Failed to validate goal event timestamp",
			v.buildGoalEventTimestampLogFields(ctx, ev)...,
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", v.event.Id),
				zap.String("environment_id", v.event.EnvironmentId),
			)...,
		)
		return nil, err
	}
	return ev, nil
}

func (v *eventEvaluationValidator) validate(ctx context.Context) (string, error) {
	ev, err := v.unmarshal(ctx)
	if err != nil {
		return codeUnmarshalFailed, errUnmarshalFailed
	}
	v.lastUnmarshaledEvent = ev

	if err := uuid.ValidateUUID(v.event.Id); err != nil {
		v.logger.Warn(
			"Failed to validate evaluation event id format",
			append(v.buildEvaluationEventLogFields(ctx, ev), zap.Error(err))...,
		)
		return codeInvalidID, errInvalidIDFormat
	}

	// validate required fields
	if ev.FeatureId == "" {
		v.logger.Debug(
			"Empty feature_id",
			v.buildEvaluationEventLogFields(ctx, ev)...,
		)
		return codeEmptyField, errEmptyFeatureID
	}
	isErrorReason := isEvaluationEventErrorReason(ev.Reason)

	if !isErrorReason && ev.VariationId == "" {
		v.logger.Debug(
			"Empty variation_id",
			v.buildEvaluationEventLogFields(ctx, ev)...,
		)
		return codeEmptyField, errEmptyVariationID
	}
	if (ev.User == nil || (ev.User != nil && ev.User.Id == "")) && ev.UserId == "" {
		v.logger.Debug(
			"Empty user_id",
			v.buildEvaluationEventLogFields(ctx, ev)...,
		)
		return codeEmptyField, errEmptyUserID
	}
	if ev.Reason == nil {
		v.logger.Debug(
			"Nil reason",
			v.buildEvaluationEventLogFields(ctx, ev)...,
		)
		return codeEmptyField, errNilReason
	}

	if !validateTimestamp(ev.Timestamp, v.oldestTimestampDuration, v.furthestTimestampDuration) {
		v.logger.Debug(
			"Failed to validate evaluation event timestamp",
			v.buildEvaluationEventTimestampLogFields(ctx, ev)...,
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", v.event.Id),
				zap.String("environment_id", v.event.EnvironmentId),
			)...,
		)
		return nil, err
	}
	return ev, nil
}

// For metrics events we don't need to validate the timestamp
func (v *eventMetricsValidator) validate(ctx context.Context) (string, error) {
	ev, err := v.unmarshal(ctx)
	if err != nil {
		return codeUnmarshalFailed, errUnmarshalFailed
	}
	if err := uuid.ValidateUUID(v.event.Id); err != nil {
		v.logger.Warn(
			"Failed to validate metrics event id format",
			append(v.buildMetricsEventLogFields(ctx, ev), zap.Error(err))...,
		)
		return codeInvalidID, errInvalidIDFormat
	}
	return "", nil
}

func (v *eventMetricsValidator) unmarshal(ctx context.Context) (*eventproto.MetricsEvent, error) {
	ev := &eventproto.MetricsEvent{}
	if err := ptypes.UnmarshalAny(v.event.Event, ev); err != nil {
		v.logger.Error(
			"Failed to extract metrics event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", v.event.Id),
				zap.String("environment_id", v.event.EnvironmentId),
			)...,
		)
		return nil, err
	}
	return ev, nil
}

// validateTimestamp limits date range of the given timestamp to align with database
// retention policies and prevent problematic timestamps that could affect data quality.
func validateTimestamp(
	timestamp int64,
	oldestTimestampDuration, furthestTimestampDuration time.Duration,
) bool {
	given := time.Unix(timestamp, 0)

	// Reject events older than retention policy (default: 31 days)
	maxPast := time.Now().Add(-oldestTimestampDuration)
	if given.Before(maxPast) {
		return false
	}

	// Reject events too far in the future (clock skew, malicious timestamps)
	maxFuture := time.Now().Add(furthestTimestampDuration)
	return !given.After(maxFuture)
}

// buildGoalEventLogFields creates common log fields for goal event logging
func (v *eventGoalValidator) buildGoalEventLogFields(
	ctx context.Context,
	ev *eventproto.GoalEvent,
) []zap.Field {
	return log.FieldsFromIncomingContext(ctx).AddFields(
		zap.String("id", v.event.Id),
		zap.Int64("timestamp", ev.Timestamp),
		zap.String("environment_id", v.event.EnvironmentId),
		zap.String("goal_id", ev.GoalId),
		zap.Float64("value", ev.Value),
		zap.Any("user", ev.User),
		zap.String("user_id", ev.UserId),
		zap.Any("metadata", ev.Metadata),
		zap.String("tag", ev.Tag),
		zap.String("sdk_version", ev.SdkVersion),
		zap.String("source_id", ev.SourceId.String()),
	)
}

// buildGoalEventTimestampLogFields creates enhanced log fields for goal event timestamp validation
func (v *eventGoalValidator) buildGoalEventTimestampLogFields(
	ctx context.Context,
	ev *eventproto.GoalEvent,
) []zap.Field {
	now := time.Now().Unix()
	ageHours := float64(now-ev.Timestamp) / 3600.0

	return append(v.buildGoalEventLogFields(ctx, ev),
		zap.Int64("currentTime", now),
		zap.Float64("ageHours", ageHours),
		zap.String("timestampDate", time.Unix(ev.Timestamp, 0).Format(time.RFC3339)),
	)
}

// buildEvaluationEventLogFields creates common log fields for evaluation event logging
func (v *eventEvaluationValidator) buildEvaluationEventLogFields(
	ctx context.Context,
	ev *eventproto.EvaluationEvent,
) []zap.Field {
	return log.FieldsFromIncomingContext(ctx).AddFields(
		zap.String("id", v.event.Id),
		zap.Int64("timestamp", ev.Timestamp),
		zap.String("environment_id", v.event.EnvironmentId),
		zap.String("feature_id", ev.FeatureId),
		zap.Int32("feature_version", ev.FeatureVersion),
		zap.String("variation_id", ev.VariationId),
		zap.Any("reason", ev.Reason),
		zap.Any("user", ev.User),
		zap.String("user_id", ev.UserId),
		zap.Any("metadata", ev.Metadata),
		zap.String("tag", ev.Tag),
		zap.String("source_id", ev.SourceId.String()),
		zap.String("sdk_version", ev.SdkVersion),
	)
}

// buildEvaluationEventTimestampLogFields creates enhanced log fields for evaluation event timestamp validation
func (v *eventEvaluationValidator) buildEvaluationEventTimestampLogFields(
	ctx context.Context,
	ev *eventproto.EvaluationEvent,
) []zap.Field {
	now := time.Now().Unix()
	ageHours := float64(now-ev.Timestamp) / 3600.0

	return append(v.buildEvaluationEventLogFields(ctx, ev),
		zap.Int64("currentTime", now),
		zap.Float64("ageHours", ageHours),
		zap.String("timestampDate", time.Unix(ev.Timestamp, 0).Format(time.RFC3339)),
	)
}

// buildMetricsEventLogFields creates common log fields for metrics event logging
func (v *eventMetricsValidator) buildMetricsEventLogFields(
	ctx context.Context,
	ev *eventproto.MetricsEvent,
) []zap.Field {
	return log.FieldsFromIncomingContext(ctx).AddFields(
		zap.String("id", v.event.Id),
		zap.String("environment_id", v.event.EnvironmentId),
		zap.Int64("timestamp", ev.Timestamp),
		zap.String("source_id", ev.SourceId.String()),
		zap.String("sdk_version", ev.SdkVersion),
		zap.Any("metadata", ev.Metadata),
	)
}

// isEvaluationEventErrorReason returns true if the reason indicates the user
// received the default value due to an error (e.g., flag not found, cache miss).
// Returns false for nil reason (e.g., from old SDKs that omit the field).
// Must stay in sync with featureproto.Reason_Type error variants.
func isEvaluationEventErrorReason(reason *feature.Reason) bool {
	if reason == nil {
		return false
	}
	return reason.Type == feature.Reason_ERROR_NO_EVALUATIONS ||
		reason.Type == feature.Reason_ERROR_FLAG_NOT_FOUND ||
		reason.Type == feature.Reason_ERROR_WRONG_TYPE ||
		reason.Type == feature.Reason_ERROR_USER_ID_NOT_SPECIFIED ||
		reason.Type == feature.Reason_ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED ||
		reason.Type == feature.Reason_ERROR_EXCEPTION ||
		reason.Type == feature.Reason_ERROR_CACHE_NOT_FOUND ||
		reason.Type == feature.Reason_CLIENT
}
