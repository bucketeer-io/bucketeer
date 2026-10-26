// Copyright 2026 The Bucketeer Authors.
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

package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecimalScanAndValue(t *testing.T) {
	t.Parallel()
	var decimal Decimal
	require.NoError(t, decimal.Scan(nil))
	assert.Equal(t, Decimal(0), decimal)

	require.NoError(t, decimal.Scan([]byte("12.5")))
	assert.Equal(t, Decimal(12.5), decimal)
	value, err := decimal.Value()
	require.NoError(t, err)
	assert.Equal(t, "12.500", value)

	require.NoError(t, decimal.Scan("3.25"))
	assert.Equal(t, Decimal(3.25), decimal)
	require.NoError(t, decimal.Scan(float64(1)))
	assert.Equal(t, Decimal(1), decimal)
	require.NoError(t, decimal.Scan(float32(2)))
	assert.Equal(t, Decimal(2), decimal)
	require.NoError(t, decimal.Scan(int64(4)))
	assert.Equal(t, Decimal(4), decimal)

	require.Error(t, decimal.Scan([]byte("nope")))
	require.Error(t, decimal.Scan(true))
}

func TestNormalizeGoalIDs(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{}, NormalizeGoalIDs(nil))
	assert.Equal(t, []string{"g"}, NormalizeGoalIDs([]string{"g"}))
}
