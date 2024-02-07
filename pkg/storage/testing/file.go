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

	"github.com/bucketeer-io/bucketeer/pkg/storage"
)

type inMemoryStorageBucket struct {
}

func NewInMemoryStorageBucket() storage.Bucket {
	return &inMemoryStorageBucket{}
}

func (f *inMemoryStorageBucket) Delete(ctx context.Context, key *storage.Key) error {
	// TODO
	return nil
}

func (f *inMemoryStorageBucket) Writer(
	ctx context.Context,
	environmentNamespace,
	filename string,
	CRC32C uint32,
) (storage.Writer, error) {
	// TODO
	return nil, nil
}

func (f *inMemoryStorageBucket) Reader(
	ctx context.Context,
	environmentNamespace,
	filename string,
) (storage.Reader, error) {
	// TODO
	return nil, nil
}
