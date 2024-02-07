// Copyright 2024 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type FeatureLastUsedStorage interface {
	GetFeatureLastUsedInfos(
		ctx context.Context,
		ids []string,
		environmentNamespace string,
	) ([]*domain.FeatureLastUsedInfo, error)
	UpsertFeatureLastUsedInfos(
		ctx context.Context,
		featureLastUsedInfos []*domain.FeatureLastUsedInfo,
		environmentNamespace string,
	) error
}

type FeatureLastUsedLister interface {
	ListFeatureLastUsedInfo(
		ctx context.Context,
		pageSize int,
		cursor, environmentNamespace string,
		filters ...*storage.Filter,
	) ([]*proto.FeatureLastUsedInfo, string, error)
}

const featureLastUsedInfoKind = "FeatureLastUsedInfo"

type featureLastUsedInfoStorage struct {
	client storage.GetPutter
}

func NewFeatureLastUsedInfoStorage(client storage.GetPutter) FeatureLastUsedStorage {
	return &featureLastUsedInfoStorage{client: client}
}

func (s *featureLastUsedInfoStorage) GetFeatureLastUsedInfos(
	ctx context.Context,
	ids []string,
	environmentNamespace string,
) ([]*domain.FeatureLastUsedInfo, error) {
	keys := make([]*storage.Key, 0, len(ids))
	featureLastUsedInfos := make([]*proto.FeatureLastUsedInfo, 0, len(keys))
	for _, k := range ids {
		keys = append(keys, s.newKey(k, environmentNamespace))
		featureLastUsedInfos = append(featureLastUsedInfos, &proto.FeatureLastUsedInfo{})
	}
	err := s.client.GetMulti(ctx, keys, featureLastUsedInfos)
	if err != nil {
		merr, ok := err.(storage.MultiError)
		if !ok {
			return nil, err
		}
		for _, e := range merr {
			switch e {
			case nil:
			case storage.ErrKeyNotFound:
			default:
				return nil, e
			}
		}
	}
	// NOTE: If the performance matters, remove the following loop and return protos.
	domainFeatureLastUsedInfos := make([]*domain.FeatureLastUsedInfo, 0, len(keys))
	for _, f := range featureLastUsedInfos {
		if f == nil {
			continue
		}
		domainFeatureLastUsedInfos = append(
			domainFeatureLastUsedInfos,
			&domain.FeatureLastUsedInfo{FeatureLastUsedInfo: f},
		)
	}
	return domainFeatureLastUsedInfos, nil
}

func (s *featureLastUsedInfoStorage) UpsertFeatureLastUsedInfos(
	ctx context.Context,
	featureLastUsedInfos []*domain.FeatureLastUsedInfo,
	environmentNamespace string,
) error {
	keys := make([]*storage.Key, 0, len(featureLastUsedInfos))
	featureLastUsedInfoProtos := make([]*proto.FeatureLastUsedInfo, 0, len(featureLastUsedInfos))
	for _, f := range featureLastUsedInfos {
		keys = append(keys, storage.NewKey(f.ID(), featureLastUsedInfoKind, environmentNamespace))
		featureLastUsedInfoProtos = append(featureLastUsedInfoProtos, f.FeatureLastUsedInfo)
	}
	return s.client.PutMulti(ctx, keys, featureLastUsedInfoProtos)
}

func (s *featureLastUsedInfoStorage) newKey(featureLastUsedInfoKey, environmentNamespace string) *storage.Key {
	return storage.NewKey(featureLastUsedInfoKey, featureLastUsedInfoKind, environmentNamespace)
}

type featureLastUsedInfoLister struct {
	client storage.Querier
}

func NewFeatureLastUsedInfoLister(client storage.Querier) FeatureLastUsedLister {
	return &featureLastUsedInfoLister{client: client}
}

func (l *featureLastUsedInfoLister) ListFeatureLastUsedInfo(
	ctx context.Context,
	pageSize int,
	cursor, environmentNamespace string,
	filters ...*storage.Filter,
) ([]*proto.FeatureLastUsedInfo, string, error) {
	query := storage.Query{
		Kind:                 featureLastUsedInfoKind,
		Limit:                pageSize,
		StartCursor:          cursor,
		Filters:              filters,
		EnvironmentNamespace: environmentNamespace,
	}
	it, err := l.client.RunQuery(ctx, query)
	if err != nil {
		return nil, "", err
	}
	featureLastUseds := make([]*proto.FeatureLastUsedInfo, 0, pageSize)
	for {
		featureLastUsed := &proto.FeatureLastUsedInfo{}
		err := it.Next(featureLastUsed)
		if err == storage.ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, "", err
		}
		featureLastUseds = append(featureLastUseds, featureLastUsed)
	}
	nextCursor, err := it.Cursor()
	if err != nil {
		return nil, "", err
	}
	return featureLastUseds, nextCursor, nil
}
