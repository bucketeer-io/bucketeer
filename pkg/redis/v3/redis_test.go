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

package v3

import (
	"testing"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewClientIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name          string
		addr          string
		expectError   bool
		expectCluster bool
	}{
		{
			name:          "standalone redis on default port",
			addr:          "localhost:6379",
			expectError:   false,
			expectCluster: false,
		},
		{
			name:          "unreachable redis",
			addr:          "localhost:9999",
			expectError:   false,
			expectCluster: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			logger := zap.NewNop()
			client, err := NewClient(tt.addr, WithLogger(logger))

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)

				if client != nil {
					client.Close()
				}
			}
		})
	}
}

func TestNewClientBehavior(t *testing.T) {
	t.Parallel()

	t.Run("unreachable redis returns client with auto mode", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient("localhost:9999", WithLogger(logger))

		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeStandard, rc.clientType)
			c.Close()
		}
	})

	t.Run("options are applied", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithPoolSize(20),
			WithMinIdleConns(5),
			WithPassword("test-password"),
		)

		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			c.Close()
		}
	})
}

func TestNewClientWithRedisMode(t *testing.T) {
	t.Parallel()

	t.Run("cluster mode creates ClusterClient", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode(RedisModeCluster),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeCluster, rc.clientType)
			_, ok := rc.rc.(*goredis.ClusterClient)
			assert.True(t, ok)
			c.Close()
		}
	})

	t.Run("standalone mode creates standard Client", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode(RedisModeStandalone),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeStandard, rc.clientType)
			_, ok := rc.rc.(*goredis.Client)
			assert.True(t, ok)
			c.Close()
		}
	})

	t.Run("auto mode defaults to standalone when unreachable", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode(RedisModeAuto),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeStandard, rc.clientType)
			c.Close()
		}
	})

	t.Run("invalid mode falls back to auto", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode("invalid"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			// auto mode defaults to standalone when Redis is unreachable
			assert.Equal(t, ClientTypeStandard, rc.clientType)
			c.Close()
		}
	})

	t.Run("case-insensitive mode parsing", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode("CLUSTER"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeCluster, rc.clientType)
			c.Close()
		}
	})
}

func TestWithRedisMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    RedisMode
		expected RedisMode
	}{
		{"cluster", RedisModeCluster, RedisModeCluster},
		{"standalone", RedisModeStandalone, RedisModeStandalone},
		{"auto", RedisModeAuto, RedisModeAuto},
		{"uppercase CLUSTER", "CLUSTER", RedisModeCluster},
		{"mixed case Standalone", "Standalone", RedisModeStandalone},
		{"invalid falls back to auto", "invalid", RedisModeAuto},
		{"empty falls back to auto", "", RedisModeAuto},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithRedisMode(tt.input)(opts)
			assert.Equal(t, tt.expected, opts.redisMode)
		})
	}
}

func TestClientTypeString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "cluster", clientTypeString(ClientTypeCluster))
	assert.Equal(t, "standalone", clientTypeString(ClientTypeStandard))
}
