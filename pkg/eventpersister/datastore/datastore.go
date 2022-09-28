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
	"sync/atomic"
)

type Writer interface {
	Write(
		ctx context.Context,
		events map[string]string,
		environmentNamespace string,
	) (map[string]bool, error)
	Close()
}

type writerPool struct {
	writes  uint64
	writers []Writer
}

func NewWriterPool(writers []Writer) Writer {
	return &writerPool{
		writers: writers,
	}
}

func (p *writerPool) Write(
	ctx context.Context,
	events map[string]string,
	environmentNamespace string,
) (map[string]bool, error) {
	writes := atomic.AddUint64(&p.writes, 1)
	index := int(writes) % len(p.writers)
	return p.writers[index].Write(ctx, events, environmentNamespace)
}

func (p *writerPool) Close() {
	for _, w := range p.writers {
		w.Close()
	}
}
