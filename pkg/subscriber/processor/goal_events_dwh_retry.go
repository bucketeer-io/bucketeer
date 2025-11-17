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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"

	ecstorage "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
)

const (
	retryGoalEventKeyKind = "goal_event_retry"
	scanBatchSize         = 100
	lockTimeout           = 30 * time.Second
)

type retryMessage struct {
	GoalEvent     *eventproto.GoalEvent `json:"goalEvent"`
	EnvironmentID string                `json:"environmentId"`
	RetryCount    int                   `json:"retryCount"`
	ID            string                `json:"id"`
	FailedEvents  []*epproto.GoalEvent  `json:"failedEvents"`
	FirstRetryAt  int64                 `json:"firstRetryAt"`
	RetryAt       int64                 `json:"retryAt"`
}

func (m *retryMessage) GetID() string { return m.ID }

// StartRetryProcessor kicks off a ticker to scan and process retry keys.
func (w *goalEvtWriter) StartRetryProcessor(ctx context.Context) {
	w.logger.Debug("Starting goal event retry processor",
		zap.Duration("interval", w.retryGoalEventInterval))
	go func() {
		ticker := time.NewTicker(w.retryGoalEventInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				w.logger.Debug("Goal event retry processor stopped")
				return
			case <-ticker.C:
				w.scanAndProcess(ctx)
			}
		}
	}()
}

func (w *goalEvtWriter) scanAndProcess(ctx context.Context) {
	var cursor uint64
	for {
		nextCursor, keys, err := w.redisClient.Scan(cursor, fmt.Sprintf("*:%s:*", retryGoalEventKeyKind), scanBatchSize)
		if err != nil {
			w.logger.Error("Scan failed", zap.Error(err), zap.Uint64("cursor", cursor))
			break
		}
		for _, key := range keys {
			w.processRetryKey(ctx, key)
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
}

func (w *goalEvtWriter) processRetryKey(ctx context.Context, key string) {
	acquired, lockVal, err := w.lockGoalEventRetryLock(ctx, key)
	if err != nil {
		w.logger.Error("Lock acquisition failed", zap.Error(err), zap.String("key", key))
		return
	}
	if !acquired {
		return
	}
	defer w.unlockGoalEventRetryLock(ctx, key, lockVal)

	data, err := w.redisClient.Get(key)
	if err != nil {
		w.logger.Error("Redis GET failed", zap.Error(err), zap.String("key", key))
		return
	}

	var msg retryMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		w.logger.Error("JSON unmarshal failed", zap.Error(err), zap.String("key", key))
		w.deleteKey(ctx, key)
		return
	}

	now := time.Now().Unix()

	// Check if retry period has been exceeded - delete the key if so
	if msg.FirstRetryAt != 0 {
		elapsed := time.Since(time.Unix(msg.FirstRetryAt, 0))
		if elapsed > w.maxRetryGoalEventPeriod {
			w.logger.Info("Retry period exceeded, deleting retry key",
				zap.String("retryId", msg.ID),
				zap.String("goalId", msg.GoalEvent.GoalId),
				zap.String("userId", msg.GoalEvent.UserId),
				zap.Duration("elapsed", elapsed),
				zap.Duration("maxPeriod", w.maxRetryGoalEventPeriod),
				zap.Int("retryCount", msg.RetryCount),
				zap.String("key", key),
			)
			// Delete the expired retry key
			if err := w.redisClient.Del(key); err != nil {
				w.logger.Error("Failed to delete expired retry key",
					zap.Error(err),
					zap.String("key", key),
				)
			}
			return
		}
	}

	if now < msg.RetryAt {
		w.logger.Debug("Not time for retry",
			zap.String("retryId", msg.ID),
			zap.String("goalId", msg.GoalEvent.GoalId),
			zap.String("userId", msg.GoalEvent.UserId),
			zap.String("retryAt", time.Unix(msg.RetryAt, 0).Format(time.RFC3339)),
			zap.Int64("now", now),
			zap.Int64("retryAtUnix", msg.RetryAt),
			zap.Int64("secondsUntilRetry", msg.RetryAt-now),
			zap.String("key", key),
		)
		return
	}

	w.logger.Info("Processing retry message",
		zap.String("retryId", msg.ID),
		zap.String("goalId", msg.GoalEvent.GoalId),
		zap.String("userId", msg.GoalEvent.UserId),
		zap.Int("retryCount", msg.RetryCount),
		zap.String("key", key),
	)
	w.handleMessage(ctx, &msg, key)
}

func (w *goalEvtWriter) handleMessage(ctx context.Context, msg *retryMessage, key string) {
	lg := w.logger.With(
		zap.String("environmentId", msg.EnvironmentID),
		zap.String("goalId", msg.GoalEvent.GoalId),
		zap.Int("retryCount", msg.RetryCount),
	)
	if len(msg.FailedEvents) > 0 {
		w.handleFailedBatch(ctx, msg, key, lg)
	} else {
		w.handleNewRetry(ctx, msg, key, lg)
	}
}

func (w *goalEvtWriter) handleFailedBatch(ctx context.Context, msg *retryMessage, key string, lg *zap.Logger) {
	failures, err := w.writer.AppendRows(ctx, msg.FailedEvents)
	if err != nil {
		lg.Error("Append failed batch", zap.Error(err))
		msg.RetryCount++
		if err := w.storeRetryMessage(msg); err != nil {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
			lg.Error("Failed to store retry message", zap.Error(err))
		}
		return
	}

	// filter out successful ones
	var still []*epproto.GoalEvent
	for id, failed := range failures {
		if failed {
			for _, ev := range msg.FailedEvents {
				if ev.Id == id {
					still = append(still, ev)
				}
			}
		}
	}
	if len(still) > 0 {
		msg.RetryCount++
		msg.FailedEvents = still
		if err := w.storeRetryMessage(msg); err != nil {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
			lg.Error("Failed to store retry message", zap.Error(err))
		}
	} else {
		w.deleteKey(ctx, key)
	}
}

func (w *goalEvtWriter) handleNewRetry(ctx context.Context, msg *retryMessage, key string, lg *zap.Logger) {
	experiments, err := w.listExperiments(ctx, msg.EnvironmentID)
	if err != nil {
		lg.Error("List experiments failed", zap.Error(err))
		msg.RetryCount++
		if err := w.storeRetryMessage(msg); err != nil {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
			lg.Error("Failed to store retry message", zap.Error(err))
		}
		return
	}
	if len(experiments) == 0 {
		// Cache might be stale - do DB fallback before marking as missing
		lg.Debug("RETRY: No experiments found in cache, checking DB",
			zap.String("retryId", msg.ID),
			zap.String("goalId", msg.GoalEvent.GoalId),
			zap.String("environmentId", msg.EnvironmentID),
		)
		dbExperiments, dbErr := w.listExperimentsFromDB(ctx, msg.EnvironmentID)
		if dbErr != nil {
			lg.Error("RETRY: Failed to list experiments from DB in fallback path",
				zap.Error(dbErr),
				zap.String("retryId", msg.ID),
				zap.String("goalId", msg.GoalEvent.GoalId),
				zap.String("environmentId", msg.EnvironmentID),
			)
			// Treat as infra error: retry later
			msg.RetryCount++
			if err := w.storeRetryMessage(msg); err != nil {
				subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
				lg.Error("Failed to store retry message", zap.Error(err))
			}
			return
		}
		if len(dbExperiments) == 0 {
			// Confirmed no experiments exist - mark as missing
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeRetryMessageNoExperiments).Inc()
			lg.Warn("RETRY: No experiments found in DB, marking as missing and deleting retry message",
				zap.String("retryId", msg.ID),
				zap.String("goalId", msg.GoalEvent.GoalId),
				zap.String("environmentId", msg.EnvironmentID),
			)
			// Mark as missing in sync.Map so future events are rejected immediately
			w.markExperimentMissing(msg.EnvironmentID, msg.GoalEvent.GoalId)
			w.deleteKey(ctx, key)
			return
		}
		// Found experiments in DB - use them
		lg.Info("RETRY: Found experiments in DB fallback",
			zap.String("retryId", msg.ID),
			zap.String("goalId", msg.GoalEvent.GoalId),
			zap.String("environmentId", msg.EnvironmentID),
			zap.Int("experimentCount", len(dbExperiments)),
		)
		experiments = dbExperiments
	}

	lg.Info("🔄 RETRY: Looking for evaluation",
		zap.String("retryId", msg.ID),
		zap.String("goalId", msg.GoalEvent.GoalId),
		zap.String("userId", msg.GoalEvent.UserId),
		zap.Int64("goalTimestamp", msg.GoalEvent.Timestamp),
		zap.String("goalTimestampISO", time.Unix(msg.GoalEvent.Timestamp, 0).Format(time.RFC3339)),
		zap.Int("retryCount", msg.RetryCount),
		zap.Int64("firstRetryAt", msg.FirstRetryAt),
		zap.String("firstRetryAtISO", func() string {
			if msg.FirstRetryAt > 0 {
				return time.Unix(msg.FirstRetryAt, 0).Format(time.RFC3339)
			}
			return "N/A"
		}()),
	)
	evals, _, err := w.linkGoalEventByExperiment(ctx, msg.GoalEvent, msg.ID, msg.EnvironmentID, msg.GoalEvent.Tag, experiments)
	if err != nil {
		if errors.Is(err, ErrExperimentNotFound) {
			// Experiment not found in the list - could be:
			// 1. Experiment doesn't match goal ID
			// 2. Goal event timestamp outside experiment window
			// 3. Experiment was deleted/stopped
			// Do DB fallback to confirm before marking as missing
			lg.Debug("RETRY: Experiment not found in list, checking DB",
				zap.String("retryId", msg.ID),
				zap.String("goalId", msg.GoalEvent.GoalId),
				zap.String("environmentId", msg.EnvironmentID),
				zap.Int("experimentCount", len(experiments)),
			)
			dbExperiments, dbErr := w.listExperimentsFromDB(ctx, msg.EnvironmentID)
			if dbErr != nil {
				lg.Error("RETRY: Failed to list experiments from DB when experiment not found",
					zap.Error(dbErr),
					zap.String("retryId", msg.ID),
					zap.String("goalId", msg.GoalEvent.GoalId),
				)
				// Treat as infra error: retry later
				msg.RetryCount++
				if err := w.storeRetryMessage(msg); err != nil {
					subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
					lg.Error("Failed to store retry message", zap.Error(err))
				}
				return
			}
			// Check if any experiment matches the goal ID (regardless of timestamp window)
			// This helps distinguish between "experiment doesn't exist" vs "experiment exists but timestamp outside window"
			hasMatchingGoalID := false
			for _, exp := range dbExperiments {
				if w.findGoalID(msg.GoalEvent.GoalId, exp.GoalIds) {
					hasMatchingGoalID = true
					break
				}
			}

			// Try linking again with DB experiments
			dbEvals, _, dbLinkErr := w.linkGoalEventByExperiment(ctx, msg.GoalEvent, msg.ID, msg.EnvironmentID, msg.GoalEvent.Tag, dbExperiments)
			if dbLinkErr != nil {
				if errors.Is(dbLinkErr, ErrExperimentNotFound) {
					if !hasMatchingGoalID {
						// No experiment exists with this goal ID - mark as missing
						subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeExperimentNotFound).Inc()
						lg.Warn("RETRY: No experiment found with matching goal ID, marking as missing and deleting retry",
							zap.String("retryId", msg.ID),
							zap.String("goalId", msg.GoalEvent.GoalId),
							zap.String("environmentId", msg.EnvironmentID),
						)
						w.markExperimentMissing(msg.EnvironmentID, msg.GoalEvent.GoalId)
					} else {
						// Experiment exists with matching goal ID but doesn't match (e.g., timestamp outside window)
						// Don't mark as missing - just delete the retry as the event is invalid
						subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeExperimentNotFound).Inc()
						lg.Warn("RETRY: Experiment exists but doesn't match goal event (e.g., timestamp outside window), deleting retry",
							zap.String("retryId", msg.ID),
							zap.String("goalId", msg.GoalEvent.GoalId),
							zap.String("environmentId", msg.EnvironmentID),
							zap.Int64("goalTimestamp", msg.GoalEvent.Timestamp),
						)
					}
					w.deleteKey(ctx, key)
				} else {
					// Other error (e.g., evaluation not found) - continue retrying
					lg.Error("RETRY: DB fallback linking failed with non-NotFound error",
						zap.Error(dbLinkErr),
						zap.String("retryId", msg.ID),
						zap.String("goalId", msg.GoalEvent.GoalId),
					)
					msg.RetryCount++
					if err := w.storeRetryMessage(msg); err != nil {
						subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
						lg.Error("Failed to store retry message", zap.Error(err))
					}
				}
				return
			}
			// Found experiment via DB fallback - use the evaluations
			lg.Info("RETRY: Found experiment via DB fallback after initial NotFound",
				zap.String("retryId", msg.ID),
				zap.String("goalId", msg.GoalEvent.GoalId),
				zap.Int("evaluationCount", len(dbEvals)),
			)
			evals = dbEvals
			// Clear negative cache since we found the experiment
			w.clearExperimentMissing(msg.EnvironmentID, msg.GoalEvent.GoalId)
			// Clear error to continue processing evaluations below
			err = nil
		} else if errors.Is(err, ecstorage.ErrNoResultsFound) {
			// Evaluation not found - this is expected, continue retrying
			lg.Warn("🔄 RETRY: Still no evaluations found, will retry again",
				zap.String("retryId", msg.ID),
				zap.String("goalId", msg.GoalEvent.GoalId),
				zap.String("userId", msg.GoalEvent.UserId),
				zap.Int64("goalTimestamp", msg.GoalEvent.Timestamp),
				zap.String("goalTimestampISO", time.Unix(msg.GoalEvent.Timestamp, 0).Format(time.RFC3339)),
				zap.Int("retryCount", msg.RetryCount),
				zap.String("tag", msg.GoalEvent.Tag),
			)
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeRetryMessageNoEvaluations).Inc()
			msg.RetryCount++
			if err := w.storeRetryMessage(msg); err != nil {
				subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
				lg.Error("Failed to store retry message", zap.Error(err))
			}
		} else {
			lg.Error("❌ RETRY: Linking failed",
				zap.Error(err),
				zap.String("retryId", msg.ID),
				zap.String("goalId", msg.GoalEvent.GoalId),
				zap.String("userId", msg.GoalEvent.UserId),
				zap.Int64("goalTimestamp", msg.GoalEvent.Timestamp),
				zap.Int("retryCount", msg.RetryCount),
			)
			msg.RetryCount++
			if err := w.storeRetryMessage(msg); err != nil {
				subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
				lg.Error("Failed to store retry message", zap.Error(err))
			}
			return
		}
		// If we found experiments via DB fallback (evals was set above), continue processing below
		// Otherwise, return here
		if err != nil {
			return
		}
	}
	// Successfully found experiment (even if no evaluations yet) - clear any negative cache entry
	// This ensures that if experiment was previously marked as missing, we clear it now
	// Note: Cache may have already been cleared above if found via DB fallback
	if err == nil {
		w.clearExperimentMissing(msg.EnvironmentID, msg.GoalEvent.GoalId)
	}
	// If err == nil but evals is empty, this means evaluations weren't found yet
	// Treat this the same as ErrNoResultsFound and retry
	if len(evals) == 0 {
		lg.Warn("🔄 RETRY: Still no evaluations found (err=nil but evals empty), will retry again",
			zap.String("retryId", msg.ID),
			zap.String("goalId", msg.GoalEvent.GoalId),
			zap.String("userId", msg.GoalEvent.UserId),
			zap.Int64("goalTimestamp", msg.GoalEvent.Timestamp),
			zap.String("goalTimestampISO", time.Unix(msg.GoalEvent.Timestamp, 0).Format(time.RFC3339)),
			zap.Int("retryCount", msg.RetryCount),
			zap.String("tag", msg.GoalEvent.Tag),
		)
		subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeRetryMessageNoEvaluations).Inc()
		msg.RetryCount++
		if err := w.storeRetryMessage(msg); err != nil {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
			lg.Error("Failed to store retry message", zap.Error(err))
		}
		return
	}
	evalDetails := make([]map[string]interface{}, len(evals))
	for i, eval := range evals {
		evalDetails[i] = map[string]interface{}{
			"featureId":      eval.FeatureID,
			"featureVersion": eval.FeatureVersion,
			"variationId":    eval.VariationID,
			"timestamp":      eval.Timestamp,
		}
	}
	lg.Info("✅ RETRY: Evaluation found!",
		zap.String("retryId", msg.ID),
		zap.String("goalId", msg.GoalEvent.GoalId),
		zap.String("userId", msg.GoalEvent.UserId),
		zap.Int64("goalTimestamp", msg.GoalEvent.Timestamp),
		zap.String("goalTimestampISO", time.Unix(msg.GoalEvent.Timestamp, 0).Format(time.RFC3339)),
		zap.Int("evaluationCount", len(evals)),
		zap.Int("retryCount", msg.RetryCount),
		zap.Any("evaluations", evalDetails),
	)

	// Note: sync.Map already cleared above when experiment was found

	var events []*epproto.GoalEvent
	for _, ev := range evals {
		ge, _, err := w.convToGoalEvent(msg.GoalEvent, ev, msg.ID, msg.GoalEvent.Tag, msg.EnvironmentID)
		if err == nil {
			events = append(events, ge)
		}
	}
	if len(events) == 0 {
		subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeRetryMessageNoGoalEvents).Inc()
		w.deleteKey(ctx, key)
		return
	}

	if len(events) > 0 {
		eventIds := make([]string, len(events))
		for i, e := range events {
			eventIds[i] = e.Id
		}
		lg.Info("💾 RETRY: Writing goal events to data warehouse",
			zap.String("retryId", msg.ID),
			zap.String("goalId", msg.GoalEvent.GoalId),
			zap.String("userId", msg.GoalEvent.UserId),
			zap.Int("retryCount", msg.RetryCount),
			zap.Int("eventCount", len(events)),
			zap.Strings("eventIds", eventIds),
		)
	}
	fs, err := w.writer.AppendRows(ctx, events)
	if err != nil {
		lg.Error("❌ RETRY: AppendRows failed",
			zap.Error(err),
			zap.String("retryId", msg.ID),
			zap.String("goalId", msg.GoalEvent.GoalId),
			zap.String("userId", msg.GoalEvent.UserId),
			zap.Int("eventCount", len(events)),
		)
		msg.RetryCount++
		msg.FailedEvents = events
		if err := w.storeRetryMessage(msg); err != nil {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
			lg.Error("Failed to store retry message", zap.Error(err))
		}
		return
	}
	var failed []*epproto.GoalEvent
	successfulIds := make([]string, 0)
	for id, f := range fs {
		if f {
			for _, ev := range events {
				if ev.Id == id {
					failed = append(failed, ev)
				}
			}
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failed) > 0 {
		failedIds := make([]string, len(failed))
		for i, e := range failed {
			failedIds[i] = e.Id
		}
		lg.Warn("🔄 Retry: Some goal events failed to write, will retry again",
			zap.String("retryId", msg.ID),
			zap.Int("failedCount", len(failed)),
			zap.Int("successfulCount", len(successfulIds)),
			zap.Strings("failedEventIds", failedIds),
		)
		msg.RetryCount++
		msg.FailedEvents = failed
		if err := w.storeRetryMessage(msg); err != nil {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
			lg.Error("Failed to store retry message", zap.Error(err))
		}
	} else {
		lg.Info("✅ Retry: Successfully wrote all goal events, deleting retry key",
			zap.String("retryId", msg.ID),
			zap.String("goalId", msg.GoalEvent.GoalId),
			zap.String("userId", msg.GoalEvent.UserId),
			zap.Int("retryCount", msg.RetryCount),
			zap.Int("eventCount", len(events)),
			zap.Strings("eventIds", successfulIds),
		)
		subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeRetryMessageAppendSuccess).Inc()
		w.deleteKey(ctx, key)
	}
}

func (w *goalEvtWriter) deleteKey(ctx context.Context, key string) {
	if err := w.redisClient.Del(key); err != nil {
		w.logger.Error("Failed to delete key", zap.Error(err), zap.String("key", key))
	}
}

func (w *goalEvtWriter) lockGoalEventRetryLock(ctx context.Context, key string) (bool, string, error) {
	lockCtx, cancel := context.WithTimeout(ctx, lockTimeout)
	defer cancel()
	parts := strings.SplitN(key, ":", 3)
	if len(parts) < 3 {
		return false, "", fmt.Errorf("invalid retry key format: %s", key)
	}
	environmentID, eventID := parts[0], parts[2]
	acquired, value, err := w.locker.Lock(lockCtx, environmentID, eventID)
	if err != nil {
		w.logger.Error("Failed to acquire lock",
			zap.Error(err),
			zap.String("environmentId", environmentID),
			zap.String("eventId", eventID),
		)
		return false, "", err
	}
	return acquired, value, nil
}

func (w *goalEvtWriter) unlockGoalEventRetryLock(ctx context.Context, key, value string) {
	parts := strings.SplitN(key, ":", 3)
	if len(parts) < 3 {
		w.logger.Error("Invalid retry key format for unlock", zap.String("key", key))
		return
	}
	environmentID, eventID := parts[0], parts[2]
	unlocked, err := w.locker.Unlock(ctx, environmentID, eventID, value)
	if err != nil {
		w.logger.Error("Failed to release lock",
			zap.Error(err),
			zap.String("environmentId", environmentID),
			zap.String("eventId", eventID),
			zap.String("value", value),
		)
		return
	}
	if !unlocked {
		w.logger.Warn("Lock was not released, possibly expired",
			zap.String("environmentId", environmentID),
			zap.String("eventId", eventID),
			zap.String("value", value),
		)
	}
}

// computeBackoffAndTTL calculates the next retry interval and TTL for a retry message.
// It handles exponential backoff with dynamic caps based on the max retry period.
func (w *goalEvtWriter) computeBackoffAndTTL(
	retryCount int,
	firstRetryAt int64,
	initialInterval time.Duration,
	maxRetryPeriod time.Duration,
) (nextInterval time.Duration, ttl time.Duration, err error) {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = initialInterval
	bo.Multiplier = 2.0
	// dynamic cap so we never back off past the max period
	ratio := float64(maxRetryPeriod) / float64(initialInterval)
	maxExp := int(math.Floor(math.Log2(ratio)))
	bo.MaxInterval = time.Duration(1<<uint(maxExp)) * initialInterval
	bo.MaxElapsedTime = maxRetryPeriod
	bo.Reset()

	// Advance the backoff state retryCount times to get the correct interval
	for i := 0; i <= retryCount; i++ {
		nextInterval = bo.NextBackOff()
	}
	if nextInterval == backoff.Stop {
		return 0, 0, fmt.Errorf("retry period exceeded %v", maxRetryPeriod)
	}

	// Cap the retry interval to a reasonable maximum (1 minute) to prevent
	// extremely long waits between retries, especially for Redis Stream where
	// events can arrive out of order and need faster retries.
	// This balances between retrying quickly enough for tests and not overwhelming
	// the system with too frequent retries in production.
	const maxRetryInterval = 1 * time.Minute
	if nextInterval > maxRetryInterval {
		nextInterval = maxRetryInterval
	}

	// Compute TTL: ensure it's at least 2x nextInterval to account for:
	// - Clock skew
	// - Processing delays
	// - Retry processor not running exactly on time
	// But don't exceed the maxRetryPeriod
	minTTL := nextInterval * 2
	if firstRetryAt != 0 {
		elapsed := time.Since(time.Unix(firstRetryAt, 0))
		remaining := maxRetryPeriod - elapsed
		if remaining <= 0 {
			return 0, 0, fmt.Errorf("retry period exceeded %v since first retry", maxRetryPeriod)
		}
		// Use minTTL (2x nextInterval) if we have enough remaining time,
		// otherwise use remaining (we're close to max period)
		if remaining >= minTTL {
			ttl = minTTL
		} else {
			ttl = remaining
		}
	} else {
		// First retry: use minTTL or maxRetryPeriod, whichever is smaller
		if maxRetryPeriod < minTTL {
			ttl = maxRetryPeriod
		} else {
			ttl = minTTL
		}
	}

	return nextInterval, ttl, nil
}

func (w *goalEvtWriter) storeRetryMessage(msg *retryMessage) error {
	now := time.Now().Unix()
	key := fmt.Sprintf("%s:%s:%s", msg.EnvironmentID, retryGoalEventKeyKind, msg.ID)

	nextInterval, ttl, err := w.computeBackoffAndTTL(
		msg.RetryCount,
		msg.FirstRetryAt,
		w.retryGoalEventInterval,
		w.maxRetryGoalEventPeriod,
	)
	if err != nil {
		return err
	}

	msg.RetryAt = now + int64(nextInterval.Seconds())
	if msg.FirstRetryAt == 0 {
		msg.FirstRetryAt = now
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	w.logger.Info("Storing retry message",
		zap.String("key", key),
		zap.String("retryId", msg.ID),
		zap.String("goalId", msg.GoalEvent.GoalId),
		zap.String("userId", msg.GoalEvent.UserId),
		zap.Int("retryCount", msg.RetryCount),
		zap.Duration("ttl", ttl),
		zap.Int64("retryAt", msg.RetryAt),
		zap.String("retryAtTime", time.Unix(msg.RetryAt, 0).Format(time.RFC3339)),
	)
	return w.redisClient.Set(key, data, ttl)
}
