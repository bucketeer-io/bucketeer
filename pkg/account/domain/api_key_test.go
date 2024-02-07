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
	"github.com/stretchr/testify/require"

	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestNewAPIKey(t *testing.T) {
	a, err := NewAPIKey("name", proto.APIKey_SDK)
	assert.NoError(t, err)
	assert.Equal(t, "name", a.Name)
	assert.Equal(t, proto.APIKey_SDK, a.Role)
}

func TestGenerateKey(t *testing.T) {
	key, err := generateKey()
	require.NoError(t, err)
	require.NotEmpty(t, key)
}

func TestRename(t *testing.T) {
	a, err := NewAPIKey("name", proto.APIKey_SDK)
	assert.NoError(t, err)
	a.Rename("test")
	assert.Equal(t, "test", a.Name)
}

func TestAPIKeyEnable(t *testing.T) {
	a, err := NewAPIKey("name", proto.APIKey_SDK)
	assert.NoError(t, err)
	a.Disabled = true
	a.Enable()
	assert.Equal(t, false, a.Disabled)
}

func TestAPIKeyDisable(t *testing.T) {
	a, err := NewAPIKey("name", proto.APIKey_SDK)
	assert.NoError(t, err)
	a.Disable()
	assert.Equal(t, true, a.Disabled)
}
