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

	storage "cloud.google.com/go/bigquery/storage/apiv1"
	storagepb "google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1"

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type EvalEventReader interface {
	Read(ctx context.Context,
		selectedFields []string,
		rowRestriction string,
	) ([]*featureproto.Evaluation, error)
}

type evalEventReader struct {
	client *storage.BigQueryReadClient
}

func NewEvalEventReader(
	ctx context.Context,
) (EvalEventReader, error) {
	c, err := storage.NewBigQueryReadClient(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return &evalEventReader{
		client: c,
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
		// TODO: Fill request struct fields.
		// See https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1#CreateReadSessionRequest.
	}
}
