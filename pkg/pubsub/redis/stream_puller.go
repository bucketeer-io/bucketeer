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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/puller.go
package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	v3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

const (
	// Default values
	defaultBatchSize = 10
	defaultBlockTime = 5 * time.Second
	defaultIdleTime  = 60 * time.Second
	defaultBackoff   = 1 * time.Second
	maxBackoff       = 30 * time.Second
)

// StreamPuller is a Redis Stream implementation of the Puller interface
type StreamPuller struct {
	redisClient    v3.Client
	subscription   string // Consumer group name
	topicBase      string // Base stream name
	partitionCount int    // Number of partitions
	consumer       string // Consumer name
	batchSize      int64
	blockTime      time.Duration
	idleTime       time.Duration
	closed         bool
	mutex          sync.Mutex
	done           chan struct{}
	logger         *zap.Logger
	messages       chan *puller.Message
	handler        func(context.Context, *puller.Message) // Store the handler function
}

type StreamPullerOption func(*StreamPuller)

// WithStreamPullerMetrics sets the metrics registerer for the puller
func WithStreamPullerMetrics(registerer metrics.Registerer) StreamPullerOption {
	return func(p *StreamPuller) {
		// Redis client metrics are already registered by the Redis client
	}
}

// WithStreamPullerLogger sets the logger for the puller
func WithStreamPullerLogger(logger *zap.Logger) StreamPullerOption {
	return func(p *StreamPuller) {
		p.logger = logger
	}
}

// WithStreamPullerBatchSize sets the batch size for reading from the stream
func WithStreamPullerBatchSize(size int64) StreamPullerOption {
	return func(p *StreamPuller) {
		p.batchSize = size
	}
}

// WithStreamPullerBlockTime sets the block time for reading from the stream
func WithStreamPullerBlockTime(d time.Duration) StreamPullerOption {
	return func(p *StreamPuller) {
		p.blockTime = d
	}
}

// WithStreamPullerIdleTime sets the idle time for detecting stale messages
func WithStreamPullerIdleTime(d time.Duration) StreamPullerOption {
	return func(p *StreamPuller) {
		p.idleTime = d
	}
}

// WithStreamPullerPartitionCount sets the number of partitions for the puller
func WithStreamPullerPartitionCount(count int) StreamPullerOption {
	return func(p *StreamPuller) {
		p.partitionCount = count
	}
}

// getStreamKey returns the partitioned stream name
func (p *StreamPuller) getStreamKey(partition int) string {
	return fmt.Sprintf("%s-%d", p.topicBase, partition)
}

// NewStreamPuller creates a new Redis Stream puller
func NewStreamPuller(
	client v3.Client,
	subscription string, // Consumer group name
	topic string, // Stream name
	opts ...StreamPullerOption,
) puller.Puller {
	p := &StreamPuller{
		redisClient:    client,
		subscription:   subscription,
		topicBase:      topic,
		partitionCount: defaultStreamPartitionCount,
		consumer:       fmt.Sprintf("%s-%s-%d", subscription, topic, time.Now().UnixNano()),
		batchSize:      defaultBatchSize,
		blockTime:      defaultBlockTime,
		idleTime:       defaultIdleTime,
		logger:         zap.NewNop(),
		done:           make(chan struct{}),
		messages:       make(chan *puller.Message),
	}

	for _, opt := range opts {
		opt(p)
	}

	p.logger = p.logger.Named("redis-stream-puller")

	return p
}

// Pull reads messages from the stream and calls the handler for each message
func (p *StreamPuller) Pull(ctx context.Context, handler func(context.Context, *puller.Message)) error {
	p.mutex.Lock()
	if p.closed {
		p.mutex.Unlock()
		return fmt.Errorf("redis stream puller is closed")
	}

	// Store the handler function for use in recoveryLoop
	p.handler = handler
	p.mutex.Unlock()

	// Create consumer groups for all partitions
	for partition := 0; partition < p.partitionCount; partition++ {
		streamKey := p.getStreamKey(partition)
		err := p.redisClient.XGroupCreateMkStream(streamKey, p.subscription, "0")
		if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
			p.logger.Error("Failed to create consumer group",
				zap.Error(err),
				zap.String("subscription", p.subscription),
				zap.String("stream", streamKey),
			)
			// Continue with other partitions even if this one fails
		}
	}

	// Start a background goroutine to reclaim stale messages
	go p.recoveryLoop(ctx)

	// Setup mechanism to handle context cancellation
	ctxDone := ctx.Done()
	backoff := defaultBackoff

	// Process messages in the current goroutine to block until context is canceled
	for {
		select {
		case <-ctxDone:
			p.logger.Debug("Context canceled, stopping pull",
				zap.String("subscription", p.subscription),
				zap.String("topicBase", p.topicBase),
			)
			p.Close()
			return ctx.Err()

		case <-p.done:
			p.logger.Debug("Puller closed, stopping pull",
				zap.String("subscription", p.subscription),
				zap.String("topicBase", p.topicBase),
			)
			return nil

		default:
			// Build streams argument for XREADGROUP
			// It should alternate between stream keys and ">" for each partition
			streams := make([]string, 0, p.partitionCount*2)
			for partition := 0; partition < p.partitionCount; partition++ {
				streamKey := p.getStreamKey(partition)
				streams = append(streams, streamKey, ">") // ">" means only new messages
			}

			// Read from all partitions
			streamResults, err := p.redisClient.XReadGroup(
				ctx,
				p.subscription,
				p.consumer,
				streams,
				p.batchSize,
				p.blockTime,
			)

			if err != nil {
				if err == context.Canceled || err == context.DeadlineExceeded {
					// Context canceled or deadline exceeded
					return err
				}

				if err == goredis.Nil {
					// No messages available, try again
					continue
				}

				// Exponential backoff on error
				p.logger.Error("Failed to read from streams",
					zap.Error(err),
					zap.String("subscription", p.subscription),
					zap.String("topicBase", p.topicBase),
					zap.Duration("backoff", backoff),
				)

				select {
				case <-time.After(backoff):
					// Exponential backoff with maximum
					backoff = min(backoff*2, maxBackoff)
				case <-ctxDone:
					return ctx.Err()
				case <-p.done:
					return nil
				}
				continue
			}

			// Reset backoff on success
			backoff = defaultBackoff

			// Process messages from all streams
			for _, streamResult := range streamResults {
				for _, msg := range streamResult.Messages {
					// Create ack and nack functions that will be called by the handler
					streamKey := streamResult.Stream
					ackFunc := func() {
						if err := p.redisClient.XAck(streamKey, p.subscription, msg.ID); err != nil {
							p.logger.Error("Failed to acknowledge message",
								zap.Error(err),
								zap.String("subscription", p.subscription),
								zap.String("stream", streamKey),
								zap.String("id", msg.ID),
							)
						}
					}

					nackFunc := func() {
						// For Redis Streams, not acknowledging a message is equivalent to a NACK
						// The message will remain in the pending entries list and will be redelivered
						p.logger.Debug("Message not acknowledged (NACK)",
							zap.String("subscription", p.subscription),
							zap.String("stream", streamKey),
							zap.String("id", msg.ID),
						)
					}

					// Extract data from the message
					// In Redis Streams, values are a map where the key is the field name and the value is the field value
					// We expect one field with the message ID as the key and the serialized message as the value
					var data []byte
					for _, value := range msg.Values {
						if s, ok := value.(string); ok {
							data = []byte(s)
							break
						} else if b, ok := value.([]byte); ok {
							data = b
							break
						}
					}

					// Create a message with Ack/Nack functions
					message := &puller.Message{
						ID:   msg.ID,
						Data: data,
						Attributes: map[string]string{
							"id":     msg.ID,
							"stream": streamKey,
						},
						Ack:  ackFunc,
						Nack: nackFunc,
					}

					// Handle message
					p.handler(ctx, message)
				}
			}
		}
	}
}

// recoveryLoop periodically checks for stale messages and reclaims them
func (p *StreamPuller) recoveryLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Debug("Context canceled, stopping recovery loop",
				zap.String("subscription", p.subscription),
				zap.String("topicBase", p.topicBase),
			)
			return
		case <-p.done:
			p.logger.Debug("Puller closed, stopping recovery loop",
				zap.String("subscription", p.subscription),
				zap.String("topicBase", p.topicBase),
			)
			return
		case <-ticker.C:
			// Skip processing if we don't have a handler
			if p.handler == nil {
				p.logger.Debug("No handler registered, skipping recovery",
					zap.String("subscription", p.subscription),
					zap.String("topicBase", p.topicBase),
				)
				continue
			}

			// Check each partition for stale messages
			for partition := 0; partition < p.partitionCount; partition++ {
				streamKey := p.getStreamKey(partition)

				// Retrieve pending messages that have been idle longer than idleTime
				pendingMessages, err := p.redisClient.XPendingExt(
					ctx,
					streamKey,
					p.subscription,
					"-", // Start
					"+", // End
					10,  // Count
					p.idleTime,
				)
				if err != nil {
					p.logger.Error("Failed to get pending messages",
						zap.Error(err),
						zap.String("subscription", p.subscription),
						zap.String("stream", streamKey),
					)
					continue // Continue with next partition
				}

				if len(pendingMessages) == 0 {
					// No stale messages to reclaim for this partition
					continue
				}

				p.logger.Info("Found stale messages to reclaim",
					zap.Int("count", len(pendingMessages)),
					zap.String("subscription", p.subscription),
					zap.String("stream", streamKey),
				)

				// Collect message IDs
				messageIDs := make([]string, len(pendingMessages))
				for i, pm := range pendingMessages {
					messageIDs[i] = pm.ID
					p.logger.Debug("Claiming stale message",
						zap.String("id", pm.ID),
						zap.String("previous_consumer", pm.Consumer),
						zap.String("subscription", p.subscription),
						zap.String("stream", streamKey),
					)
				}

				// Claim the messages for the current consumer
				claimed, err := p.redisClient.XClaim(
					ctx,
					streamKey,
					p.subscription,
					p.consumer,
					p.idleTime,
					messageIDs,
				)
				if err != nil {
					p.logger.Error("Failed to claim messages",
						zap.Error(err),
						zap.String("subscription", p.subscription),
						zap.String("stream", streamKey),
					)
					continue
				}

				p.logger.Info("Successfully claimed stale messages",
					zap.Int("claimed_count", len(claimed)),
					zap.Int("requested_count", len(messageIDs)),
					zap.String("subscription", p.subscription),
					zap.String("stream", streamKey),
				)

				// Reprocess the claimed messages
				p.reprocessClaimedMessages(ctx, claimed, streamKey)
			}
		}
	}
}

// reprocessClaimedMessages handles claimed messages and reprocesses them
func (p *StreamPuller) reprocessClaimedMessages(ctx context.Context, claimed []goredis.XMessage, streamKey string) {
	reprocessedCount := 0
	for _, msg := range claimed {
		p.logger.Debug("Reprocessing claimed message",
			zap.String("subscription", p.subscription),
			zap.String("stream", streamKey),
			zap.String("id", msg.ID),
		)

		// Create ack and nack functions just like in the Pull method
		ackFunc := func() {
			if err := p.redisClient.XAck(streamKey, p.subscription, msg.ID); err != nil {
				p.logger.Error("Failed to acknowledge claimed message",
					zap.Error(err),
					zap.String("subscription", p.subscription),
					zap.String("stream", streamKey),
					zap.String("id", msg.ID),
				)
			}
		}

		nackFunc := func() {
			// For Redis Streams, not acknowledging a message is equivalent to a NACK
			// The message will remain in the pending entries list and will be redelivered
			p.logger.Debug("Claimed message not acknowledged (NACK)",
				zap.String("subscription", p.subscription),
				zap.String("stream", streamKey),
				zap.String("id", msg.ID),
			)
		}

		// Extract data from the message
		var data []byte
		for _, value := range msg.Values {
			if s, ok := value.(string); ok {
				data = []byte(s)
				break
			} else if b, ok := value.([]byte); ok {
				data = b
				break
			}
		}

		// Create a message with Ack/Nack functions
		message := &puller.Message{
			ID:   msg.ID,
			Data: data,
			Attributes: map[string]string{
				"stream":  streamKey,
				"claimed": "true",
			},
			Ack:  ackFunc,
			Nack: nackFunc,
		}

		// Process the message in a goroutine to avoid blocking the recovery loop
		// Create a new context with timeout for the handler
		handlerCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		go func(ctx context.Context, msg *puller.Message) {
			defer cancel() // Ensure the context is cancelled when the goroutine exits
			// Use the same handler that was passed to Pull
			p.handler(ctx, msg)
		}(handlerCtx, message)

		reprocessedCount++
	}

	p.logger.Info("Reprocessing claimed messages",
		zap.Int("reprocessed_count", reprocessedCount),
		zap.String("subscription", p.subscription),
		zap.String("stream", streamKey),
	)
}

// Close closes the puller
func (p *StreamPuller) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	close(p.done)

	return nil
}

// SubscriptionName returns the name of the subscription
func (p *StreamPuller) SubscriptionName() string {
	return fmt.Sprintf("%s:%s", p.subscription, p.topicBase)
}

// Helper function to get the minimum of two durations
func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
