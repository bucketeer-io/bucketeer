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
	"database/sql"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
	mysqldb "github.com/go-sql-driver/mysql"
)

type mysqlOptions struct {
	logger  *zap.Logger
	metrics metrics.Registerer
}

type MySQLOption func(*mysqlOptions)

func WithMySQLLogger(l *zap.Logger) MySQLOption {
	return func(opts *mysqlOptions) {
		opts.logger = l
	}
}

func WithMySQLMetrics(r metrics.Registerer) MySQLOption {
	return func(opts *mysqlOptions) {
		opts.metrics = r
	}
}

// MySQLWriter is the interface for writing events to MySQL
type MySQLWriter interface {
	AppendRows(ctx context.Context, batches [][][]byte) ([]int, error)
	Close() error
}

// MySQLEvalEventWriter writes evaluation events to MySQL
type MySQLEvalEventWriter struct {
	client     *sql.DB
	insertStmt *sql.Stmt
	batchSize  int
	logger     *zap.Logger
}

// MySQLGoalEventWriter writes goal events to MySQL
type MySQLGoalEventWriter struct {
	client     *sql.DB
	insertStmt *sql.Stmt
	batchSize  int
	logger     *zap.Logger
}

// NewMySQLEvalEventWriter creates a new MySQL writer for evaluation events
func NewMySQLEvalEventWriter(
	ctx context.Context,
	dsn string,
	batchSize int,
	opts ...MySQLOption,
) (*MySQLEvalEventWriter, error) {
	dopts := &mysqlOptions{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}

	logger := dopts.logger.Named("mysql-eval-writer")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("Failed to connect to MySQL", zap.Error(err))
		return nil, err
	}

	// Set connection parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Create prepared statement for evaluation events
	stmt, err := db.PrepareContext(ctx, `
		INSERT INTO evaluation_event (
			id, environment_id, timestamp, feature_id, feature_version, 
			user_id, user_data, variation_id, reason, tag, source_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		logger.Error("Failed to prepare statement", zap.Error(err))
		db.Close()
		return nil, err
	}

	return &MySQLEvalEventWriter{
		client:     db,
		insertStmt: stmt,
		batchSize:  batchSize,
		logger:     logger,
	}, nil
}

// NewMySQLGoalEventWriter creates a new MySQL writer for goal events
func NewMySQLGoalEventWriter(
	ctx context.Context,
	dsn string,
	batchSize int,
	opts ...MySQLOption,
) (*MySQLGoalEventWriter, error) {
	dopts := &mysqlOptions{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}

	logger := dopts.logger.Named("mysql-goal-writer")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("Failed to connect to MySQL", zap.Error(err))
		return nil, err
	}

	// Set connection parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Create prepared statement for goal events
	stmt, err := db.PrepareContext(ctx, `
		INSERT INTO goal_event (
			id, environment_id, timestamp, goal_id, value, 
			user_id, user_data, tag, source_id, feature_id,
			feature_version, variation_id, reason
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		logger.Error("Failed to prepare statement", zap.Error(err))
		db.Close()
		return nil, err
	}

	return &MySQLGoalEventWriter{
		client:     db,
		insertStmt: stmt,
		batchSize:  batchSize,
		logger:     logger,
	}, nil
}

// AppendRows appends evaluation event rows to MySQL
func (w *MySQLEvalEventWriter) AppendRows(
	ctx context.Context,
	batches [][][]byte,
) ([]int, error) {
	failedBatches := []int{}
	tx, err := w.client.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Error("Failed to start transaction", zap.Error(err))
		return []int{}, err
	}
	defer tx.Rollback()

	// Use the transaction's prepared statement
	stmt := tx.Stmt(w.insertStmt)

	for batchIdx, batch := range batches {
		batchFailed := false

		for _, eventBytes := range batch {
			event := &epproto.EvaluationEvent{}
			if err := proto.Unmarshal(eventBytes, event); err != nil {
				w.logger.Error("Failed to unmarshal evaluation event",
					zap.Error(err),
					zap.Int("batch", batchIdx),
				)
				batchFailed = true
				break
			}

			// Marshal user data to JSON
			userData, err := json.Marshal(event.UserData)
			if err != nil {
				w.logger.Error("Failed to marshal user data",
					zap.Error(err),
					zap.String("event_id", event.Id),
				)
				batchFailed = true
				break
			}

			// Convert timestamp to time.Time
			timestamp := time.Unix(event.Timestamp/1000, (event.Timestamp%1000)*1000000)

			_, err = stmt.ExecContext(ctx,
				event.Id,
				event.EnvironmentId,
				timestamp,
				event.FeatureId,
				event.FeatureVersion,
				event.UserId,
				userData,
				event.VariationId,
				event.Reason,
				event.Tag,
				event.SourceId,
			)

			if err != nil {
				// Check for duplicate key error
				if mysqlErr, ok := err.(*mysqldb.MySQLError); ok && mysqlErr.Number == 1062 {
					// Duplicate key, log but don't fail the batch
					w.logger.Warn("Duplicate evaluation event, skipping",
						zap.String("event_id", event.Id),
					)
					continue
				}

				w.logger.Error("Failed to insert evaluation event",
					zap.Error(err),
					zap.String("event_id", event.Id),
				)
				batchFailed = true
				break
			}
		}

		if batchFailed {
			failedBatches = append(failedBatches, batchIdx)
		}
	}

	if err := tx.Commit(); err != nil {
		w.logger.Error("Failed to commit transaction", zap.Error(err))
		return []int{0, len(batches) - 1}, err
	}

	return failedBatches, nil
}

// AppendRows appends goal event rows to MySQL
func (w *MySQLGoalEventWriter) AppendRows(
	ctx context.Context,
	batches [][][]byte,
) ([]int, error) {
	failedBatches := []int{}
	tx, err := w.client.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Error("Failed to start transaction", zap.Error(err))
		return []int{}, err
	}
	defer tx.Rollback()

	// Use the transaction's prepared statement
	stmt := tx.Stmt(w.insertStmt)

	for batchIdx, batch := range batches {
		batchFailed := false

		for _, eventBytes := range batch {
			event := &epproto.GoalEvent{}
			if err := proto.Unmarshal(eventBytes, event); err != nil {
				w.logger.Error("Failed to unmarshal goal event",
					zap.Error(err),
					zap.Int("batch", batchIdx),
				)
				batchFailed = true
				break
			}

			// Marshal user data to JSON
			userData, err := json.Marshal(event.UserData)
			if err != nil {
				w.logger.Error("Failed to marshal user data",
					zap.Error(err),
					zap.String("event_id", event.Id),
				)
				batchFailed = true
				break
			}

			// Convert timestamp to time.Time
			timestamp := time.Unix(event.Timestamp/1000, (event.Timestamp%1000)*1000000)

			_, err = stmt.ExecContext(ctx,
				event.Id,
				event.EnvironmentId,
				timestamp,
				event.GoalId,
				event.Value,
				event.UserId,
				userData,
				event.Tag,
				event.SourceId,
				event.FeatureId,
				event.FeatureVersion,
				event.VariationId,
				event.Reason,
			)

			if err != nil {
				// Check for duplicate key error
				if mysqlErr, ok := err.(*mysqldb.MySQLError); ok && mysqlErr.Number == 1062 {
					// Duplicate key, log but don't fail the batch
					w.logger.Warn("Duplicate goal event, skipping",
						zap.String("event_id", event.Id),
					)
					continue
				}

				w.logger.Error("Failed to insert goal event",
					zap.Error(err),
					zap.String("event_id", event.Id),
				)
				batchFailed = true
				break
			}
		}

		if batchFailed {
			failedBatches = append(failedBatches, batchIdx)
		}
	}

	if err := tx.Commit(); err != nil {
		w.logger.Error("Failed to commit transaction", zap.Error(err))
		return []int{0, len(batches) - 1}, err
	}

	return failedBatches, nil
}

// Close closes the MySQL client and statement
func (w *MySQLEvalEventWriter) Close() error {
	if err := w.insertStmt.Close(); err != nil {
		return err
	}
	return w.client.Close()
}

// Close closes the MySQL client and statement
func (w *MySQLGoalEventWriter) Close() error {
	if err := w.insertStmt.Close(); err != nil {
		return err
	}
	return w.client.Close()
}
