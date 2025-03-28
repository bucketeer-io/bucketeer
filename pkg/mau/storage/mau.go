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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package storage

import (
	"context"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

type MAUStorage interface {
	DeleteRecords(ctx context.Context, partition string) error
	RebuildPartition(ctx context.Context, partition string) error
	DropPartition(ctx context.Context, partition string) error
	CreatePartition(ctx context.Context, partition, lessThan string) error
}

type mauStorage struct {
	qe mysql.QueryExecer
}

func NewMAUStorage(qe mysql.QueryExecer) MAUStorage {
	return &mauStorage{qe: qe}
}

func (s *mauStorage) DeleteRecords(ctx context.Context, partition string) error {
	query := fmt.Sprintf(`DELETE FROM mau partition(%s)`, partition)
	_, err := s.qe.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (s *mauStorage) RebuildPartition(ctx context.Context, partition string) error {
	query := fmt.Sprintf(`ALTER TABLE mau REBUILD PARTITION %s`, partition)
	_, err := s.qe.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (s *mauStorage) DropPartition(ctx context.Context, partition string) error {
	query := fmt.Sprintf(`ALTER TABLE mau DROP PARTITION %s`, partition)
	_, err := s.qe.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
func (s *mauStorage) CreatePartition(ctx context.Context, partition, lessThan string) error {
	query := fmt.Sprintf(`ALTER TABLE mau ADD PARTITION(PARTITION %s VALUES LESS THAN ('%s'))`, partition, lessThan)
	_, err := s.qe.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
