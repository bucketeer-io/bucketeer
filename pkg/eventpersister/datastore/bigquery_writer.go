// Copyright 2022 The Bucketeer Authors.
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

package datastore

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/query"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

type EvalEventWriter interface {
	AppendRows(ctx context.Context, events []*ecproto.EvaluationEvent) error
}

type GoalEventWriter interface {
	AppendRows(ctx context.Context, events []*ecproto.GoalEvent) error
}

type evalEventWriter struct {
	*queryClient
}

type goalEventWriter struct {
	*queryClient
}

type queryClient struct {
	query query.Query
}

func NewEvalEventWriter(q query.Query) EvalEventWriter {
	return &evalEventWriter{
		queryClient: &queryClient{
			query: q,
		},
	}
}

func NewGoalEventWriter(q query.Query) GoalEventWriter {
	return &goalEventWriter{
		queryClient: &queryClient{
			query: q,
		},
	}
}

func (ew *evalEventWriter) AppendRows(
	ctx context.Context,
	events []*ecproto.EvaluationEvent,
) error {
	// Encode the messages into binary format.
	encoded := make([][]byte, len(events))
	for k, v := range events {
		b, err := proto.Marshal(v)
		if err != nil {
			return err
		}
		encoded[k] = b
	}
	if err := ew.query.AppendRows(ctx, encoded); err != nil {
		return err
	}
	return nil
}

func (gw *goalEventWriter) AppendRows(
	ctx context.Context,
	events []*ecproto.GoalEvent,
) error {
	// Encode the messages into binary format.
	encoded := make([][]byte, len(events))
	for k, v := range events {
		b, err := proto.Marshal(v)
		if err != nil {
			return err
		}
		encoded[k] = b
	}
	if err := gw.query.AppendRows(ctx, encoded); err != nil {
		return err
	}
	return nil
}
