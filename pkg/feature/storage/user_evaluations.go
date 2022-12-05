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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package storage

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto" // nolint:staticcheck

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	storage "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	tableName    = "user_evaluations"
	columnFamily = "ue"
	columnName   = "evaluations"
)

type UserEvaluationsStorage interface {
	UpsertUserEvaluation(ctx context.Context, evaluation *featureproto.Evaluation, environmentNamespace, tag string) error
	GetUserEvaluations(ctx context.Context, userID, environmentNamespace, tag string) ([]*featureproto.Evaluation, error)
	GetUserEvaluation(
		ctx context.Context,
		userID, environmentNamespace, tag, featureID string,
		featureVersion int32,
	) (*featureproto.Evaluation, error)
}

type userEvaluationsStorage struct {
	client storage.Client
}

func NewUserEvaluationsStorage(client storage.Client) UserEvaluationsStorage {
	return &userEvaluationsStorage{client: client}
}

func (s *userEvaluationsStorage) UpsertUserEvaluation(
	ctx context.Context,
	evaluation *featureproto.Evaluation,
	environmentNamespace, tag string,
) error {
	key := newKey(
		environmentNamespace,
		tag,
		evaluation.UserId,
		evaluation.FeatureId,
		evaluation.FeatureVersion,
	)
	value, err := proto.Marshal(evaluation)
	if err != nil {
		return err
	}
	req := &storage.WriteRequest{
		TableName:    tableName,
		ColumnFamily: columnFamily,
		ColumnName:   columnName,
		Items: []*storage.WriteItem{
			{
				Key:   key,
				Value: value,
			},
		},
	}
	if err := s.client.WriteRow(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *userEvaluationsStorage) GetUserEvaluation(
	ctx context.Context,
	userID, environmentNamespace, tag, featureID string,
	featureVersion int32,
) (*featureproto.Evaluation, error) {
	req := &storage.ReadRowRequest{
		TableName:    tableName,
		ColumnFamily: columnFamily,
		RowKey: newKey(
			environmentNamespace,
			tag,
			userID,
			featureID,
			featureVersion,
		),
	}
	it, err := s.client.ReadRow(ctx, req)
	if err != nil {
		return nil, err
	}
	item, err := it.ReadItem(columnFamily, columnName)
	if err != nil {
		return nil, err
	}
	evaluation := &featureproto.Evaluation{}
	if err := proto.Unmarshal(item.Value, evaluation); err != nil {
		return nil, err
	}
	return evaluation, nil
}

func (s *userEvaluationsStorage) GetUserEvaluations(
	ctx context.Context,
	userID, environmentNamespace, tag string,
) ([]*featureproto.Evaluation, error) {
	prefix := newPrefix(environmentNamespace, tag, userID)
	req := &storage.ReadRequest{
		TableName:    tableName,
		ColumnFamily: columnFamily,
		RowSet:       storage.RowPrefix(prefix),
		RowFilters: []storage.RowFilter{
			storage.LatestNFilter(1),
		},
	}
	it, err := s.client.ReadRows(ctx, req)
	if err != nil {
		return nil, err
	}
	items, err := it.ReadItems(columnFamily, columnName)
	if err != nil {
		return nil, err
	}
	evaluations := make([]*featureproto.Evaluation, 0, len(items))
	for _, item := range items {
		evaluation := &featureproto.Evaluation{}
		if err := proto.Unmarshal(item.Value, evaluation); err != nil {
			return nil, err
		}
		evaluations = append(evaluations, evaluation)
	}
	return evaluations, nil
}

func newKey(
	environmentNamespace, tag, userID, featureID string,
	featureVersion int32,
) string {
	evaluationID := domain.EvaluationID(featureID, featureVersion, userID)
	key := fmt.Sprintf("%s#%s#%s", userID, tag, evaluationID)
	return storage.NewKey(environmentNamespace, key)
}

func newPrefix(environmentNamespace, tag, userID string) string {
	prefix := fmt.Sprintf("%s#%s", userID, tag)
	return storage.NewKey(environmentNamespace, prefix)
}
