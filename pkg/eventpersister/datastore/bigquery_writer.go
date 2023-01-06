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

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

const (
	evaluationEventTable = "evaluation_event"
	goalEventTable       = "goal_event"
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
	query bigquery.Query
}

func NewEvalEventWriter(
	ctx context.Context,
	project, dataset string,
	opts ...bigquery.QueryOption,
) (EvalEventWriter, error) {
	evt := ecproto.EvaluationEvent{}
	q, err := bigquery.NewQuery(
		ctx,
		project,
		dataset,
		evaluationEventTable,
		evt.ProtoReflect().Descriptor(),
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return &evalEventWriter{
		queryClient: &queryClient{
			query: q,
		},
	}, nil
}

func NewGoalEventWriter(
	ctx context.Context,
	project, dataset string,
	opts ...bigquery.QueryOption,
) (GoalEventWriter, error) {
	evt := ecproto.GoalEvent{}
	q, err := bigquery.NewQuery(
		ctx,
		project,
		dataset,
		goalEventTable,
		evt.ProtoReflect().Descriptor(),
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return &goalEventWriter{
		queryClient: &queryClient{
			query: q,
		},
	}, nil
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
