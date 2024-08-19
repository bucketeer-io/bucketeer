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

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestNewSearchFilter(t *testing.T) {
	s, err := NewSearchFilter("name", "query", true, proto.FilterTargetType_GOAL, "environmentID")
	assert.Nil(t, err)
	assert.Equal(t, "name", s.Name)
	assert.Equal(t, "query", s.Query)
	assert.Equal(t, true, s.DefaultFilter)
	assert.Equal(t, proto.FilterTargetType_GOAL, s.FilterTargetType)
}

func TestSetDefaultFilter(t *testing.T) {
	s, _ := NewSearchFilter("name", "query", true, proto.FilterTargetType_GOAL, "environmentID")
	assert.Equal(t, true, s.DefaultFilter)

	s.SetDefaultFilter(false)
	assert.Equal(t, false, s.DefaultFilter)
}

func TestUpdateSearchFilter(t *testing.T) {
	s, _ := NewSearchFilter("name", "query", true, proto.FilterTargetType_GOAL, "environmentID")
	assert.Equal(t, "name", s.Name)
	assert.Equal(t, "query", s.Query)
	assert.Equal(t, true, s.DefaultFilter)
	assert.Equal(t, proto.FilterTargetType_GOAL, s.FilterTargetType)

	s.UpdateSearchFilter("newName", "newQuery", false)
	assert.Equal(t, "newName", s.Name)
	assert.Equal(t, "newQuery", s.Query)
	assert.Equal(t, false, s.DefaultFilter)
}
