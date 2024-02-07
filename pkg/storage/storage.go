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
	"errors"
	"fmt"
	"io"
)

const (
	AdminEnvironmentNamespace = ""

	OrderDirectionAsc  OrderDirection = 0
	OrderDirectionDesc OrderDirection = 1

	// use 0 instead of -1 because it's used for cap to make slice in storage lister.
	QueryUnlimited = 0
)

var (
	ErrConcurrentTransaction = errors.New("storage: concurrent transaction in progress")
	ErrKeyAlreadyExists      = errors.New("storage: key already exists")
	ErrKeyNotFound           = errors.New("storage: key not found")
	ErrIteratorDone          = errors.New("storage: iterator is done")
	ErrInvalidCursor         = errors.New("storage: cursor is invalid")
	ErrBucketNotExist        = errors.New("storage: bucket doesn't exist")
	ErrObjectNotExist        = errors.New("storage: object doesn't exist")
	ErrEmptyName             = errors.New("storage: name is empty")
	ErrInvalidName           = errors.New("storage: invalid name")
)

// MultiError is returned by batch operations when there are errors with
// particular elements. Errors will be in a one-to-one correspondence with
// the input elements; successful elements will have a nil entry.
type MultiError []error

func (m MultiError) Error() string {
	s, n := "", 0
	for _, e := range m {
		if e != nil {
			if n == 0 {
				s = e.Error()
			}
			n++
		}
	}
	switch n {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, n-1)
}

type Key struct {
	ID   string
	Kind string
	// If it is empty string, query will be executed in admin namespace.
	// If not, query will be executed in namespace for target environment.
	EnvironmentNamespace string
}

func NewKey(id, kind, environmentNamespace string) *Key {
	return &Key{ID: id, Kind: kind, EnvironmentNamespace: environmentNamespace}
}

type Iterator interface {
	Next(dst interface{}) error
	Cursor() (string, error)
}

type Getter interface {
	Get(ctx context.Context, key *Key, dst interface{}) error
	GetMulti(ctx context.Context, keys []*Key, dst interface{}) error
}

type Putter interface {
	Put(ctx context.Context, key *Key, src interface{}) error
	PutMulti(ctx context.Context, keys []*Key, src interface{}) error
}

type GetPutter interface {
	Getter
	Putter
}

type Query struct {
	Kind        string
	Limit       int
	StartCursor string
	Orders      []*Order
	Filters     []*Filter
	// If it is empty string, query will be executed in admin namespace.
	// If not, query will be executed in namespace for target environment.
	EnvironmentNamespace string
}

type Filter struct {
	Property string
	Operator string
	Value    interface{}
}

func NewFilter(property, operator string, value interface{}) *Filter {
	return &Filter{
		Property: property,
		Operator: operator,
		Value:    value,
	}
}

type OrderDirection int

type Order struct {
	Property  string
	Direction OrderDirection
}

func NewOrder(property string, direction OrderDirection) *Order {
	return &Order{
		Property:  property,
		Direction: direction,
	}
}

type Querier interface {
	RunQuery(ctx context.Context, query Query) (Iterator, error)
}

type Transaction interface {
	GetPutter
	Deleter
}

type Client interface {
	GetPutter
	Deleter
	Querier
	RunInTransaction(ctx context.Context, f func(t Transaction) error) error
	Close()
}

type Writer interface {
	io.WriteCloser
}

type Reader interface {
	io.ReadCloser
}

type Deleter interface {
	Delete(ctx context.Context, key *Key) error
}

type ObjectStorageClient interface {
	Bucket(ctx context.Context, bucket string) (Bucket, error)
	Close()
}

type Bucket interface {
	Object
}

type Object interface {
	Writer(ctx context.Context, environmentNamespace, filename string, CRC32C uint32) (Writer, error)
	Reader(ctx context.Context, environmentNamespace, filename string) (Reader, error)
	Deleter
}
