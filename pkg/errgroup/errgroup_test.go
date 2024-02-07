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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	testcases := []struct {
		err      error
		finished int32
		failed   int32
	}{
		{
			err:      nil,
			finished: 1,
			failed:   0,
		},
		{
			err:      errors.New("test"),
			finished: 1,
			failed:   1,
		},
	}
	for i, tc := range testcases {
		g := Group{}
		doneCh := g.Go(func() error {
			return tc.err
		})
		<-doneCh
		des := fmt.Sprintf("index: %d", i)
		assert.Equal(t, tc.finished, g.FinishedCount(), des)
		assert.Equal(t, tc.failed, g.FailedCount(), des)
	}
}

func TestGroupNoFinished(t *testing.T) {
	g := Group{}
	g.Go(func() error {
		time.Sleep(time.Second)
		return nil
	})
	assert.Equal(t, int32(0), g.FinishedCount())
	assert.Equal(t, int32(0), g.FailedCount())
}

func TestGroupPanic(t *testing.T) {
	g := Group{}
	doneCh := g.Go(func() error {
		panic("test")
	})
	<-doneCh
	assert.Equal(t, int32(1), g.FinishedCount())
	assert.Equal(t, int32(1), g.FailedCount())
}
