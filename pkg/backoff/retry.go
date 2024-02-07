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

package backoff

import (
	"context"
	"time"
)

type Retry interface {
	WaitNext() bool
	Calls() int
}

func NewRetry(ctx context.Context, max int, backoff Backoff) Retry {
	return &retry{
		max:     max,
		backoff: backoff,
		ctx:     ctx,
	}
}

type retry struct {
	max     int
	calls   int
	ctx     context.Context
	backoff Backoff
}

func (r *retry) WaitNext() bool {
	defer func() {
		r.calls++
	}()
	if r.calls >= r.max {
		return false
	}
	d := r.backoff.Next()
	if d == 0 {
		return true
	}
	t := time.NewTimer(d)
	select {
	case <-r.ctx.Done():
		return false
	case <-t.C:
	}
	return true
}

func (r *retry) Calls() int {
	return r.calls
}
