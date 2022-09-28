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

package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	storage "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	btmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable/mock"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	environmentNamespace = "environmentNamespace"
	tag                  = "tag"
	userID               = "user-id"
)

var (
	evaluation = &featureproto.Evaluation{
		FeatureId:      "feature-id",
		FeatureVersion: 1,
		UserId:         "user-id",
		VariationId:    "variation-id",
		VariationValue: "variation-value",
	}
)

func TestNewUserEvaluationsStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := NewUserEvaluationsStorage(btmock.NewMockClient(mockController))
	assert.IsType(t, &userEvaluationsStorage{}, db)
}

type rows struct {
	columnFamily string
	value        []byte
}

func (r *rows) ReadItems(column string) ([]*storage.ReadItem, error) {
	items := []*storage.ReadItem{
		{
			RowKey:    "Row-1",
			Column:    fmt.Sprintf("%s:%s", r.columnFamily, column),
			Timestamp: 0,
			Value:     r.value,
		},
	}
	return items, nil
}

func TestGetUserEvaluations(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	prefix := newPrefix(environmentNamespace, tag, userID)
	req := &storage.ReadRequest{
		TableName:    tableName,
		ColumnFamily: columnFamily,
		RowSet:       storage.RowPrefix(prefix),
		RowFilters: []storage.RowFilter{
			storage.LatestNFilter(1),
		},
	}
	value, err := proto.Marshal(evaluation)
	assert.NoError(t, err)
	patterns := []struct {
		desc        string
		setup       func(context.Context, *userEvaluationsStorage)
		expected    []*featureproto.Evaluation
		expectedErr error
	}{
		{
			desc: "ErrInternal",
			setup: func(ctx context.Context, s *userEvaluationsStorage) {
				s.client.(*btmock.MockClient).EXPECT().ReadRows(
					ctx,
					req,
				).Return(nil, storage.ErrInternal)
			},
			expected:    nil,
			expectedErr: storage.ErrInternal,
		},
		{
			desc: "ErrKeyNotFound",
			setup: func(ctx context.Context, s *userEvaluationsStorage) {
				s.client.(*btmock.MockClient).EXPECT().ReadRows(
					ctx,
					req,
				).Return(nil, storage.ErrKeyNotFound)
			},
			expected:    nil,
			expectedErr: storage.ErrKeyNotFound,
		},
		{
			desc: "Success",
			setup: func(ctx context.Context, s *userEvaluationsStorage) {
				s.client.(*btmock.MockClient).EXPECT().ReadRows(
					ctx,
					req,
				).Return(
					&rows{
						columnFamily: columnFamily,
						value:        value,
					},
					nil,
				)
			},
			expected:    []*featureproto.Evaluation{evaluation},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		s := createNewUserEvaluationsStorage(mockController)
		p.setup(ctx, s)
		actual, err := s.GetUserEvaluations(
			ctx,
			userID,
			environmentNamespace,
			tag,
		)
		if p.expected != nil {
			assert.True(t, proto.Equal(p.expected[0], actual[0]), p.desc)
		}
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestUpsertUserEvaluation(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	value, err := proto.Marshal(evaluation)
	assert.NoError(t, err)
	key := newKey(
		environmentNamespace,
		tag,
		evaluation.UserId,
		evaluation.FeatureId,
		evaluation.FeatureVersion,
	)
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
	patterns := []struct {
		desc     string
		setup    func(context.Context, *userEvaluationsStorage)
		expected error
	}{
		{
			desc: "ErrInternal",
			setup: func(ctx context.Context, s *userEvaluationsStorage) {
				s.client.(*btmock.MockClient).EXPECT().WriteRow(
					ctx,
					req,
				).Return(storage.ErrInternal)
			},
			expected: storage.ErrInternal,
		},
		{
			desc: "Success",
			setup: func(ctx context.Context, s *userEvaluationsStorage) {
				s.client.(*btmock.MockClient).EXPECT().WriteRow(
					ctx,
					req,
				).Return(nil)
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		s := createNewUserEvaluationsStorage(mockController)
		p.setup(ctx, s)
		actual := s.UpsertUserEvaluation(
			ctx,
			evaluation,
			environmentNamespace,
			tag,
		)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
	}
}

func createNewUserEvaluationsStorage(c *gomock.Controller) *userEvaluationsStorage {
	return &userEvaluationsStorage{
		client: btmock.NewMockClient(c),
	}
}
