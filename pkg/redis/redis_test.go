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

package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/health"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyCluster struct {
	healthy bool
}

func (c *dummyCluster) Get() redis.Conn {
	return &dummyConn{healthy: c.healthy}
}

func (c *dummyCluster) Stats() map[string]redis.PoolStats {
	return nil
}

func (c *dummyCluster) Close() error {
	return nil
}

type dummyConn struct {
	healthy bool
}

func (c *dummyConn) Close() error {
	return nil
}

func (c *dummyConn) Err() error {
	return nil
}

func (c *dummyConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if c.healthy {
		return nil, nil
	}
	return nil, errors.New("error")
}

func (c *dummyConn) Send(commandName string, args ...interface{}) error {
	return nil
}

func (c *dummyConn) Flush() error {
	return nil
}

func (c *dummyConn) Receive() (reply interface{}, err error) {
	return nil, nil
}

func TestRedisCheckHealthy(t *testing.T) {
	cluster := &cluster{
		redisCluster: &dummyCluster{healthy: true},
		opts:         &options{},
		logger:       zap.NewNop(),
	}
	status := cluster.Check(context.TODO())
	if status != health.Healthy {
		t.Fail()
	}
}

func TestRedisCheckUnealthy(t *testing.T) {
	cluster := &cluster{
		redisCluster: &dummyCluster{healthy: false},
		opts:         &options{},
		logger:       zap.NewNop(),
	}
	status := cluster.Check(context.TODO())
	if status != health.Unhealthy {
		t.Fail()
	}
}

func TestRedisCheckTimeout(t *testing.T) {
	cluster := &cluster{
		redisCluster: &dummyCluster{healthy: true},
		opts:         &options{},
		logger:       zap.NewNop(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	status := cluster.Check(ctx)
	if status != health.Unhealthy {
		t.Fail()
	}
}

func TestWithDialPassword(t *testing.T) {
	t.Parallel()
	opts := &options{}
	require.Equal(t, "", opts.dialPassword)
	WithDialPassword("test-password")(opts)
	assert.Equal(t, "test-password", opts.dialPassword)
}

func TestWithDialConnectTimeout(t *testing.T) {
	t.Parallel()
	opts := &options{}
	require.Equal(t, time.Duration(0), opts.dialConnectTimeout)
	WithDialConnectTimeout(time.Minute)(opts)
	assert.Equal(t, time.Minute, opts.dialConnectTimeout)
}

func TestWithPoolMaxIdle(t *testing.T) {
	t.Parallel()
	opts := &options{}
	require.Equal(t, 0, opts.poolMaxIdle)
	WithPoolMaxIdle(1)(opts)
	assert.Equal(t, 1, opts.poolMaxIdle)
}

func TestWithPoolMaxActive(t *testing.T) {
	t.Parallel()
	opts := &options{}
	require.Equal(t, 0, opts.poolMaxActive)
	WithPoolMaxActive(1)(opts)
	assert.Equal(t, 1, opts.poolMaxActive)
}

func TestWithPoolIdleTimeout(t *testing.T) {
	t.Parallel()
	opts := &options{}
	require.Equal(t, time.Duration(0), opts.poolIdleTimeout)
	WithPoolIdleTimeout(time.Second)(opts)
	assert.Equal(t, time.Second, opts.poolIdleTimeout)
}

func TestWithServerName(t *testing.T) {
	t.Parallel()
	opts := &options{}
	require.Equal(t, "", opts.serverName)
	WithServerName("non-persistent-redis")(opts)
	assert.Equal(t, "non-persistent-redis", opts.serverName)
}
