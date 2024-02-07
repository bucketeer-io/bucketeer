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

package errgroup

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"

	eg "golang.org/x/sync/errgroup"
)

type Group struct {
	eg.Group
	finishedCount int32
	failedCount   int32
}

func (g *Group) Go(f func() error) <-chan struct{} {
	doneCh := make(chan struct{})
	g.Group.Go(func() (err error) {
		defer func() {
			atomic.AddInt32(&g.finishedCount, 1)
			if r := recover(); r != nil {
				fmt.Printf("errgroup: recovered, err: %v, stacktrace: %s\n", r, takeStacktrace())
				err = errors.New("errgroup: panic")
			}
			if err != nil {
				atomic.AddInt32(&g.failedCount, 1)
			}
			close(doneCh)
		}()
		err = f()
		return
	})
	return doneCh
}

func (g *Group) FinishedCount() int32 {
	return atomic.LoadInt32(&g.finishedCount)
}

func (g *Group) FailedCount() int32 {
	return atomic.LoadInt32(&g.failedCount)
}

func takeStacktrace() string {
	callers := make([]string, 0, 16)
	var pc [16]uintptr
	n := runtime.Callers(2, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc)
		name := fn.Name()
		callers = append(callers, fmt.Sprintf("%s\n\t%s:%d", file, name, line))
	}
	return strings.Join(callers, "\n")
}
