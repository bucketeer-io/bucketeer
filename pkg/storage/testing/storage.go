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

package testing

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto" // nolint:staticcheck

	"github.com/bucketeer-io/bucketeer/pkg/storage"
)

var (
	errSourceMustBeProtoMessage = errors.New("storage: source is not a proto message")
	errMultiArgTypeInvalid      = errors.New("storage: src has invalid type")
	errDifferentLength          = errors.New("storage: keys and src slices have different length")
)

type iterator struct {
}

func (i *iterator) Next(dst interface{}) error {
	return storage.ErrIteratorDone
}

func (i *iterator) Cursor() (string, error) {
	return "", nil
}

type inMemoryStorage struct {
	data  map[string]interface{}
	mutex sync.Mutex
}

func NewInMemoryStorage() storage.Client {
	return &inMemoryStorage{
		data: make(map[string]interface{}),
	}
}

func (s *inMemoryStorage) Get(ctx context.Context, key *storage.Key, dst interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if val, ok := s.data[s.key(key)]; ok {
		err := proto.Unmarshal(val.([]byte), dst.(proto.Message))
		if err != nil {
			return err
		}
		return nil
	}
	return storage.ErrKeyNotFound
}

func (s *inMemoryStorage) GetMulti(ctx context.Context, keys []*storage.Key, dst interface{}) error {
	values := reflect.ValueOf(dst)
	if values.Kind() != reflect.Slice {
		return errMultiArgTypeInvalid
	}
	if len(keys) != values.Len() {
		return errDifferentLength
	}
	if len(keys) == 0 {
		return nil
	}
	multiErr, any := make(storage.MultiError, len(keys)), false
	for i, key := range keys {
		if val, ok := s.data[s.key(key)]; ok {
			if !values.Index(i).CanInterface() {
				return errMultiArgTypeInvalid
			}
			msg, ok := values.Index(i).Interface().(proto.Message)
			if !ok {
				return errSourceMustBeProtoMessage
			}
			err := proto.Unmarshal(val.([]byte), msg)
			if err != nil {
				return err
			}
		} else {
			multiErr[i], any = storage.ErrKeyNotFound, true
		}
	}
	if any {
		return multiErr
	}
	return nil
}

func (s *inMemoryStorage) Put(ctx context.Context, key *storage.Key, src interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	msg, ok := src.(proto.Message)
	if !ok {
		return errSourceMustBeProtoMessage
	}
	buffer, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	s.data[s.key(key)] = buffer
	return nil
}

func (s *inMemoryStorage) PutMulti(ctx context.Context, keys []*storage.Key, src interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	values := reflect.ValueOf(src)
	if values.Kind() != reflect.Slice {
		return errMultiArgTypeInvalid
	}
	if len(keys) != values.Len() {
		return errDifferentLength
	}
	if len(keys) == 0 {
		return nil
	}
	for i, key := range keys {
		if !values.Index(i).CanInterface() {
			return errMultiArgTypeInvalid
		}
		msg, ok := values.Index(i).Interface().(proto.Message)
		if !ok {
			return errSourceMustBeProtoMessage
		}
		buffer, err := proto.Marshal(msg)
		if err != nil {
			return err
		}
		s.data[s.key(key)] = buffer
	}
	return nil
}

func (s *inMemoryStorage) RunQuery(ctx context.Context, query storage.Query) (storage.Iterator, error) {
	return &iterator{}, nil
}

func (s *inMemoryStorage) RunInTransaction(ctx context.Context, f func(t storage.Transaction) error) error {
	return f(s)
}

func (s *inMemoryStorage) Delete(ctx context.Context, key *storage.Key) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, s.key(key))
	return nil
}

func (s *inMemoryStorage) Close() {
}

func (s *inMemoryStorage) key(key *storage.Key) string {
	if key.EnvironmentNamespace == storage.AdminEnvironmentNamespace {
		return fmt.Sprintf("%s:%s", key.Kind, key.ID)
	}
	return fmt.Sprintf("%s:%s:%s", key.EnvironmentNamespace, key.Kind, key.ID)
}
