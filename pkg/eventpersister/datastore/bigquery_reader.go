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
	"log"
	"sync"

	storage "cloud.google.com/go/bigquery/storage/apiv1"
	storagepb "google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1"
	"google.golang.org/grpc"

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	"github.com/googleapis/gax-go/v2"
)

// rpcOpts is used to configure the underlying gRPC client to accept large
// messages.  The BigQuery Storage API may send message blocks up to 128MB
// in size.
var rpcOpts = gax.WithGRPCOptions(
	grpc.MaxCallRecvMsgSize(1024 * 1024 * 129),
)

type EvalEventReader interface {
	Read(ctx context.Context,
		selectedFields []string,
		rowRestriction string,
	) ([]*featureproto.Evaluation, error)
}

type evalEventReader struct {
	client *storage.BigQueryReadClient
	table  string
	parent string
}

func NewEvalEventReader(
	ctx context.Context,
	project, dataset string,
) (EvalEventReader, error) {
	c, err := storage.NewBigQueryReadClient(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return &evalEventReader{
		client: c,
		table: fmt.Sprintf("projects/%s/datasets/%s/tables/%s",
			project,
			dataset,
			evaluationEventTable,
		),
		parent: fmt.Sprintf("projects/%s", project),
	}, nil
}

func (er *evalEventReader) Read(
	ctx context.Context,
	selectedFields []string,
	rowRestriction string,
) ([]*featureproto.Evaluation, error) {
	tableReadOptions := &storagepb.ReadSession_TableReadOptions{
		SelectedFields: selectedFields,
		RowRestriction: rowRestriction,
	}

	req := &storagepb.CreateReadSessionRequest{
		Parent: er.parent,
		ReadSession: &storagepb.ReadSession{
			Table:       er.table,
			ReadOptions: tableReadOptions,
		},
	}
	// Create the session from the request.
	session, err := er.client.CreateReadSession(ctx, req, rpcOpts)
	if err != nil {
		return nil, err
	}
	if len(session.GetStreams()) == 0 {
		return nil, nil
	}

	// We'll use only a single stream for reading data from the table.  Because
	// of dynamic sharding, this will yield all the rows in the table. However,
	// if you wanted to fan out multiple readers you could do so by having a
	// increasing the MaxStreamCount.
	readStream := session.GetStreams()[0].Name
	ch := make(chan *storagepb.ReadRowsResponse)

	// Use a waitgroup to coordinate the reading and decoding goroutines.
	var wg sync.WaitGroup

	// Start the reading in one goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := processStream(ctx, bqReadClient, readStream, ch); err != nil {
			log.Fatalf("processStream failure: %v", err)
		}
		close(ch)
	}()

	// Start Avro processing and decoding in another goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		switch *format {
		case ARROW_FORMAT:
			err = processArrow(ctx, session.GetArrowSchema().GetSerializedSchema(), ch)
		case AVRO_FORMAT:
			err = processAvro(ctx, session.GetAvroSchema().GetSchema(), ch)
		}
		if err != nil {
			log.Fatalf("error processing %s: %v", *format, err)
		}
	}()

	// Wait until both the reading and decoding goroutines complete.
	wg.Wait()


}
