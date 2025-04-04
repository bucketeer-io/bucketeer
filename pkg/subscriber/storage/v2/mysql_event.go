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

package v2

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/writer"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber/storage"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
)

// MySQLEvalEventAdapter adapts the MySQL implementation to match the storage.EvalEventWriter interface
type MySQLEvalEventAdapter struct {
	writer    writer.Writer
	batchSize int
	logger    *zap.Logger
}

// MySQLGoalEventAdapter adapts the MySQL implementation to match the storage.GoalEventWriter interface
type MySQLGoalEventAdapter struct {
	writer    writer.Writer
	batchSize int
	logger    *zap.Logger
}

// NewEvalEventAdapter creates a new event writer adapter for evaluation events using MySQL
func NewEvalEventAdapter(
	ctx context.Context,
	dsn string,
	batchSize int,
	logger *zap.Logger,
) (storage.EvalEventWriter, error) {
	mysqlWriter, err := NewMySQLEvalEventWriter(ctx, dsn, batchSize, WithMySQLLogger(logger))
	if err != nil {
		return nil, err
	}

	adapter := NewMySQLAdapter(mysqlWriter, WithAdapterLogger(logger))

	return &MySQLEvalEventAdapter{
		writer:    adapter,
		batchSize: batchSize,
		logger:    logger.Named("mysql-eval-event-adapter"),
	}, nil
}

// NewGoalEventAdapter creates a new event writer adapter for goal events using MySQL
func NewGoalEventAdapter(
	ctx context.Context,
	dsn string,
	batchSize int,
	logger *zap.Logger,
) (storage.GoalEventWriter, error) {
	mysqlWriter, err := NewMySQLGoalEventWriter(ctx, dsn, batchSize, WithMySQLLogger(logger))
	if err != nil {
		return nil, err
	}

	adapter := NewMySQLAdapter(mysqlWriter, WithAdapterLogger(logger))

	return &MySQLGoalEventAdapter{
		writer:    adapter,
		batchSize: batchSize,
		logger:    logger.Named("mysql-goal-event-adapter"),
	}, nil
}

// AppendRows implements the EvalEventWriter interface for MySQL
func (w *MySQLEvalEventAdapter) AppendRows(ctx context.Context, events []*epproto.EvaluationEvent) (map[string]bool, error) {
	fails := make(map[string]bool, len(events))

	// Encode the messages into binary format
	encoded := make([][]byte, len(events))
	for i, v := range events {
		b, err := proto.Marshal(v)
		if err != nil {
			fails[v.Id] = false
			continue
		}
		encoded[i] = b
	}

	// Create batches
	batches := getBatches(encoded, w.batchSize)

	// Use the writer to append rows
	failedBatches, err := w.writer.AppendRows(ctx, batches)

	// Map batch failures to individual event IDs
	for _, batchIdx := range failedBatches {
		start := w.batchSize * batchIdx
		end := start + w.batchSize
		if end > len(events) {
			end = len(events)
		}
		for _, evt := range events[start:end] {
			fails[evt.Id] = true
		}
	}

	return fails, err
}

// AppendRows implements the GoalEventWriter interface for MySQL
func (w *MySQLGoalEventAdapter) AppendRows(ctx context.Context, events []*epproto.GoalEvent) (map[string]bool, error) {
	fails := make(map[string]bool, len(events))

	// Encode the messages into binary format
	encoded := make([][]byte, len(events))
	for i, v := range events {
		b, err := proto.Marshal(v)
		if err != nil {
			fails[v.Id] = false
			continue
		}
		encoded[i] = b
	}

	// Create batches
	batches := getBatches(encoded, w.batchSize)

	// Use the writer to append rows
	failedBatches, err := w.writer.AppendRows(ctx, batches)

	// Map batch failures to individual event IDs
	for _, batchIdx := range failedBatches {
		start := w.batchSize * batchIdx
		end := start + w.batchSize
		if end > len(events) {
			end = len(events)
		}
		for _, evt := range events[start:end] {
			fails[evt.Id] = true
		}
	}

	return fails, err
}

// getBatches creates batches from messages
func getBatches(msgs [][]byte, batchSize int) [][][]byte {
	batches := [][][]byte{}
	for i := 0; i < len(msgs); i += batchSize {
		end := i + batchSize
		if end > len(msgs) {
			end = len(msgs)
		}
		batch := msgs[i:end]
		batches = append(batches, batch)
	}
	return batches
}
