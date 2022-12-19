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
	"fmt"

	"cloud.google.com/go/bigquery/storage/managedwriter"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	"google.golang.org/protobuf/proto"

	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

const (
	evaluationEventTable = "evaluation_event"
	goalEventTable       = "goal_event"
)

type EvalEventWriter interface {
	Write(ctx context.Context, events []*ecproto.EvaluationEvent) error
}

type GoalEventWriter interface {
	Write(ctx context.Context, events []*ecproto.GoalEvent) error
}

type evalEventWriter struct {
	stream *managedwriter.ManagedStream
}

type goalEventWriter struct {
	stream *managedwriter.ManagedStream
}

func NewEvalEventWriter(
	ctx context.Context,
	project, dataset string,
) (EvalEventWriter, error) {
	c, err := managedwriter.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	evt := ecproto.EvaluationEvent{}
	descriptorProto, err := adapt.NormalizeDescriptor(evt.ProtoReflect().Descriptor())
	if err != nil {
		return nil, err
	}

	tableName := fmt.Sprintf(
		"projects/%s/datasets/%s/tables/%s",
		project,
		dataset,
		evaluationEventTable,
	)

	managedStream, err := c.NewManagedStream(
		ctx,
		managedwriter.WithSchemaDescriptor(descriptorProto),
		managedwriter.WithDestinationTable(tableName),
		managedwriter.WithType(managedwriter.CommittedStream),
	)
	if err != nil {
		return nil, err
	}
	return &evalEventWriter{
		stream: managedStream,
	}, nil
}

func NewGoalEventWriter(
	ctx context.Context,
	project, dataset string,
) (GoalEventWriter, error) {
	c, err := managedwriter.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	evt := ecproto.GoalEvent{}
	descriptorProto, err := adapt.NormalizeDescriptor(evt.ProtoReflect().Descriptor())
	if err != nil {
		return nil, err
	}

	tableName := fmt.Sprintf(
		"projects/%s/datasets/%s/tables/%s",
		project,
		dataset,
		evaluationEventTable,
	)

	managedStream, err := c.NewManagedStream(
		ctx,
		managedwriter.WithSchemaDescriptor(descriptorProto),
		managedwriter.WithDestinationTable(tableName),
		managedwriter.WithType(managedwriter.CommittedStream),
	)
	if err != nil {
		return nil, err
	}
	return &goalEventWriter{
		stream: managedStream,
	}, nil
}

func (ew *evalEventWriter) Write(
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
	result, err := ew.stream.AppendRows(ctx, encoded)
	returnedOffset, err := result.GetResult(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ew *goalEventWriter) Write(
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
	result, err := ew.stream.AppendRows(ctx, encoded)
	returnedOffset, err := result.GetResult(ctx)
	if err != nil {
		return err
	}
	return nil
}
