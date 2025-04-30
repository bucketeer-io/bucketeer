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
	"strings"
	"time"

	"go.uber.org/zap"

	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
)

const (
	retryGoalEventKeyKind = "goal_event_retry"
	scanBatchSize         = 100
	lockTimeout           = 30 * time.Second
	maxBackoffExponent    = 10 // cap exponential backoff at 2^10
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

// StartRetryProcessor starts a goroutine to process retry messages
func (w *goalEvtWriter) StartRetryProcessor(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.retryGoalEventInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				w.logger.Debug("Goal event retry processor stopped")
				return
			case <-ticker.C:
				w.logger.Debug("Starting goal event retry processor",
					zap.String("now", time.Now().Format(time.RFC3339)),
					zap.Duration("interval", w.retryGoalEventInterval))
				w.scanAndProcess(ctx)
			}
		}
	}()
}

func (w *goalEvtWriter) scanAndProcess(ctx context.Context) {
	var cursor uint64
	total := 0
	for {
		// Scan for keys in format: environmentID:goal_event_retry:eventID
		nextCursor, keys, err := w.redisClient.Scan(cursor, fmt.Sprintf("*:%s:*", retryGoalEventKeyKind), scanBatchSize)
		if err != nil {
			w.logger.Error("Scan failed", zap.Error(err), zap.Uint64("cursor", cursor))
			break
		}
		for _, key := range keys {
			w.processRetryKey(ctx, key)
		}
		total += len(keys)
		cursor = nextCursor
		if cursor == 0 {
			w.logger.Debug("Scan complete", zap.Int("totalKeys", total))
			break
		}
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
	if now < msg.RetryAt {
		w.logger.Debug("Not time for retry",
			zap.String("retryAt", time.Unix(msg.RetryAt, 0).Format(time.RFC3339)),
			zap.String("key", key),
		)
		return
	}

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
	if err != nil && !errors.Is(err, ErrExperimentNotFound) {
		lg.Error("Append failed batch", zap.Error(err))
		msg.RetryCount++
		if err := w.storeRetryMessage(msg); err != nil {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
			lg.Error("Failed to store retry message", zap.Error(err))
		}
		return
	}
	if errors.Is(err, ErrExperimentNotFound) {
		lg.Warn("Experiment not found, deleting retry message")
		w.deleteKey(ctx, key)
		return
	}
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
		if err := w.storeRetryMessage(msg); err != nil {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
			lg.Error("Failed to store retry message", zap.Error(err))
		}
		return
	}

	if len(experiments) == 0 {
		lg.Info("No experiments found, deleting retry message")
		subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeRetryMessageNoExperiments).Inc()
		w.deleteKey(ctx, key)
		return
	}

	evals, _, err := w.linkGoalEventByExperiment(ctx, msg.GoalEvent, msg.EnvironmentID, msg.GoalEvent.Tag, experiments)
	if err != nil {
		if errors.Is(err, ErrExperimentNotFound) {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeExperimentNotFound).Inc()
			lg.Error("Experiment not found, deleting retry message")
			w.deleteKey(ctx, key)
		} else {
			lg.Error("Linking failed", zap.Error(err))
			msg.RetryCount++
			if err := w.storeRetryMessage(msg); err != nil {
				subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
				lg.Error("Failed to store retry message", zap.Error(err))
			}
		}
		return
	}
	if len(evals) == 0 {
		subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeRetryMessageNoEvaluations).Inc()
		w.deleteKey(ctx, key)
		return
	}

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

	fs, err := w.writer.AppendRows(ctx, events)
	if err != nil {
		lg.Error("AppendRows failed", zap.Error(err))
		msg.RetryCount++
		msg.FailedEvents = events
		if err := w.storeRetryMessage(msg); err != nil {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
			lg.Error("Failed to store retry message", zap.Error(err))
		}
		return
	}
	var failed []*epproto.GoalEvent
	for id, f := range fs {
		if f {
			for _, ev := range events {
				if ev.Id == id {
					failed = append(failed, ev)
				}
			}
		}
	}
	if len(failed) > 0 {
		msg.RetryCount++
		msg.FailedEvents = failed
		if err := w.storeRetryMessage(msg); err != nil {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
			lg.Error("Failed to store retry message", zap.Error(err))
		}
	} else {
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

func (w *goalEvtWriter) storeRetryMessage(msg *retryMessage) error {
	now := time.Now().Unix()
	key := fmt.Sprintf("%s:%s:%s", msg.EnvironmentID, retryGoalEventKeyKind, msg.ID)

	// First retry
	if msg.FirstRetryAt == 0 {
		msg.FirstRetryAt = now
		msg.RetryAt = now + int64(w.retryGoalEventInterval.Seconds())
		ttl := w.maxRetryGoalEventPeriod
		data, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		w.logger.Debug("Storing first retry message",
			zap.String("key", key),
			zap.Duration("ttl", ttl),
			zap.Int64("retryAt", msg.RetryAt),
		)
		return w.redisClient.Set(key, data, ttl)
	}

	// Subsequent retries
	w.logger.Debug("Storing subsequent retry message",
		zap.String("key", key),
		zap.Int64("retryAt", msg.RetryAt),
	)
	elapsed := time.Since(time.Unix(msg.FirstRetryAt, 0))
	remaining := w.maxRetryGoalEventPeriod - elapsed
	if remaining <= 0 {
		return fmt.Errorf("retry period exceeded %v since first retry", w.maxRetryGoalEventPeriod)
	}
	if remaining < w.retryGoalEventInterval {
		remaining = w.retryGoalEventInterval
	}
	backoff := w.calculateExponentialBackoff(msg)
	if backoff > remaining {
		backoff = remaining
	}
	msg.RetryAt = now + int64(backoff.Seconds())
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return w.redisClient.Set(key, data, remaining)
}

func (w *goalEvtWriter) calculateExponentialBackoff(msg *retryMessage) time.Duration {
	exp := msg.RetryCount
	if exp > maxBackoffExponent {
		exp = maxBackoffExponent
	}
	return time.Duration(1<<uint(exp)) * w.retryGoalEventInterval
}
