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

package bigquery

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/bigquery/writer"
	dwhstorage "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
)

const (
	evaluationEventTable = "evaluation_event"
	goalEventTable       = "goal_event"
)

type evalEventWriter struct {
	*queryClient
}

type goalEventWriter struct {
	*queryClient
}

type queryClient struct {
	writer    writer.Writer
	batchSize int
}

// Close closes the underlying BigQuery managed writer. It is promoted to the
// evalEventWriter/goalEventWriter types so the server can close them on shutdown.
func (q *queryClient) Close() error {
	return q.writer.Close()
}

// NewEvaluationEventWriter creates a BigQuery-backed evaluation event writer.
func NewEvaluationEventWriter(
	ctx context.Context,
	logger *zap.Logger,
	project, dataset string,
	batchSize int,
	registerer metrics.Registerer,
) (dwhstorage.EvalEventWriter, error) {
	evt := epproto.EvaluationEvent{}
	w, err := writer.NewWriter(
		ctx,
		project,
		dataset,
		evaluationEventTable,
		evt.ProtoReflect().Descriptor(),
		writer.WithLogger(logger),
		writer.WithMetrics(registerer),
	)
	if err != nil {
		return nil, err
	}
	return &evalEventWriter{
		queryClient: &queryClient{
			writer:    w,
			batchSize: batchSize,
		},
	}, nil
}

// NewGoalEventWriter creates a BigQuery-backed goal event writer.
func NewGoalEventWriter(
	ctx context.Context,
	logger *zap.Logger,
	project, dataset string,
	batchSize int,
	registerer metrics.Registerer,
) (dwhstorage.GoalEventWriter, error) {
	evt := epproto.GoalEvent{}
	w, err := writer.NewWriter(
		ctx,
		project,
		dataset,
		goalEventTable,
		evt.ProtoReflect().Descriptor(),
		writer.WithLogger(logger),
		writer.WithMetrics(registerer),
	)
	if err != nil {
		return nil, err
	}
	return &goalEventWriter{
		queryClient: &queryClient{
			writer:    w,
			batchSize: batchSize,
		},
	}, nil
}

func (ew *evalEventWriter) AppendRows(
	ctx context.Context,
	events []*epproto.EvaluationEvent,
) (map[string]bool, error) {
	fails := make(map[string]bool, len(events))
	var err error
	// Encode the messages into binary format.
	encoded := make([][]byte, len(events))
	for k, v := range events {
		b, err := proto.Marshal(v)
		if err != nil {
			fails[v.Id] = false
			continue
		}
		encoded[k] = b
	}
	batches := getBatch(encoded, ew.batchSize)
	fs, err := ew.writer.AppendRows(ctx, batches)
	failMap := ew.getFailMap(events, fs)
	for id, f := range failMap {
		fails[id] = f
	}
	return fails, err
}

func (ew *evalEventWriter) getFailMap(
	es []*epproto.EvaluationEvent,
	fails []int,
) map[string]bool {
	failMap := map[string]bool{}
	for _, f := range fails {
		start := ew.batchSize * f
		end := start + ew.batchSize
		if end > len(es) {
			end = len(es)
		}
		evts := es[start:end]
		for _, evt := range evts {
			failMap[evt.Id] = true
		}
	}
	return failMap
}

func (gw *goalEventWriter) AppendRows(
	ctx context.Context,
	events []*epproto.GoalEvent,
) (map[string]bool, error) {
	fails := make(map[string]bool, len(events))
	var err error
	// Encode the messages into binary format.
	encoded := make([][]byte, len(events))
	for k, v := range events {
		b, err := proto.Marshal(v)
		if err != nil {
			fails[v.Id] = false
			continue
		}
		encoded[k] = b
	}
	batches := getBatch(encoded, gw.batchSize)
	fs, err := gw.writer.AppendRows(ctx, batches)
	failMap := gw.getFailMap(events, fs)
	for id, f := range failMap {
		fails[id] = f
	}
	return fails, err
}

func (gw *goalEventWriter) getFailMap(
	es []*epproto.GoalEvent,
	fails []int,
) map[string]bool {
	failMap := map[string]bool{}
	for _, f := range fails {
		start := gw.batchSize * f
		end := start + gw.batchSize
		if end > len(es) {
			end = len(es)
		}
		evts := es[start:end]
		for _, evt := range evts {
			failMap[evt.Id] = true
		}
	}
	return failMap
}

func getBatch(msgs [][]byte, batchSize int) [][][]byte {
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
