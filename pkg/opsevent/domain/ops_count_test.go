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

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOpsCount(t *testing.T) {
	t.Parallel()
	oc := NewOpsCount("fid", "aid", "cid", int64(1), int64(2))
	assert.Equal(t, "fid", oc.FeatureId)
	assert.Equal(t, "cid", oc.Id)
	assert.Equal(t, "aid", oc.AutoOpsRuleId)
	assert.Equal(t, "cid", oc.ClauseId)
	assert.Equal(t, int64(1), oc.OpsEventCount)
	assert.Equal(t, int64(2), oc.EvaluationCount)
	assert.NotEqual(t, int64(0), oc.UpdatedAt)
}
