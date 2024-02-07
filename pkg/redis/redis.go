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
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	clientVersion = "v1"
)

var ErrNil = redis.ErrNil

type poolStats map[string]redis.PoolStats

func (ps poolStats) ActiveCount() int {
	conns := 0
	for _, stat := range ps {
		conns += stat.ActiveCount
	}
	return conns
}

func (ps poolStats) IdleCount() int {
	conns := 0
	for _, stat := range ps {
		conns += stat.IdleCount
	}
	return conns
}

type Cluster interface {
	Get(opts ...ConnectionOptions) redis.Conn
	Check(context.Context) health.Status
	Stats() PoolStats
	Close() error
}

type cluster struct {
	redisCluster
	opts   *options
	logger *zap.Logger
}

type redisCluster interface {
	Get() redis.Conn
	Stats() map[string]redis.PoolStats
	Close() error
}

type options struct {
	dialPassword       string
	dialConnectTimeout time.Duration
	poolMaxIdle        int
	poolMaxActive      int
	poolIdleTimeout    time.Duration
	serverName         string
	metrics            metrics.Registerer
	logger             *zap.Logger
}

func defaultOptions() *options {
	return &options{
		dialConnectTimeout: 5 * time.Second,
		poolMaxIdle:        5,
		poolMaxActive:      10,
		poolIdleTimeout:    time.Minute,
		logger:             zap.NewNop(),
	}
}

type Option func(*options)

func WithDialPassword(password string) Option {
	return func(opts *options) {
		opts.dialPassword = password
	}
}

func WithDialConnectTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.dialConnectTimeout = timeout
	}
}

func WithPoolMaxIdle(num int) Option {
	return func(opts *options) {
		opts.poolMaxIdle = num
	}
}

func WithPoolMaxActive(num int) Option {
	return func(opts *options) {
		opts.poolMaxActive = num
	}
}

func WithPoolIdleTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.poolIdleTimeout = timeout
	}
}

func WithServerName(serverName string) Option {
	return func(opts *options) {
		opts.serverName = serverName
	}
}

func WithMetrics(r metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = r
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type connectionOptions struct {
	readOnly      bool
	retries       int
	tryAgainDelay time.Duration
}

type ConnectionOptions func(*connectionOptions)

var defaultConnectionOptions = connectionOptions{
	retries:       3,
	tryAgainDelay: 250 * time.Millisecond,
}

func WithReadOnly() ConnectionOptions {
	return func(opts *connectionOptions) {
		opts.readOnly = true
	}
}

func WithRetry(retries int, tryAgainDelay time.Duration) ConnectionOptions {
	return func(opts *connectionOptions) {
		opts.retries = retries
		opts.tryAgainDelay = tryAgainDelay
	}
}

func WithoutRetry() ConnectionOptions {
	return func(opts *connectionOptions) {
		opts.retries = 0
	}
}

func NewCluster(nodes []string, opts ...Option) (Cluster, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	logger := options.logger.Named("redis-v1")
	c := &redisc.Cluster{
		StartupNodes: nodes,
		DialOptions: []redis.DialOption{
			redis.DialConnectTimeout(options.dialConnectTimeout),
			redis.DialPassword(options.dialPassword),
		},
		CreatePool: createPool(options),
	}
	if err := c.Refresh(); err != nil {
		logger.Error("Failed to refresh", zap.Error(err))
		return nil, err
	}
	cluster := &cluster{
		redisCluster: c,
		opts:         options,
		logger:       logger,
	}
	if options.metrics != nil {
		RegisterMetrics(options.metrics, clientVersion, options.serverName, cluster)
	}
	return cluster, nil
}

func (c *cluster) Check(ctx context.Context) health.Status {
	resultCh := make(chan health.Status, 1)
	go func() {
		conn := c.Get(WithoutRetry())
		defer conn.Close()
		_, err := conn.Do("PING")
		if err != nil {
			c.logger.Error("Unhealthy", zap.Error(err))
			resultCh <- health.Unhealthy
			return
		}
		resultCh <- health.Healthy
	}()
	select {
	case <-ctx.Done():
		c.logger.Error("Unhealthy due to context Done is closed", zap.Error(ctx.Err()))
		return health.Unhealthy
	case status := <-resultCh:
		return status
	}
}

// This function does not return the error, that occurs while altering the original connection.
func (c *cluster) Get(opts ...ConnectionOptions) redis.Conn {
	options := defaultConnectionOptions
	for _, opt := range opts {
		opt(&options)
	}
	connection := c.redisCluster.Get()
	if options.readOnly {
		if err := redisc.ReadOnlyConn(connection); err != nil {
			c.logger.Error("Failed to create read-only connection",
				zap.Error(err),
				zap.Any("options", options))
		}
	}
	if options.retries > 0 {
		retryConnection, err := redisc.RetryConn(connection, options.retries, options.tryAgainDelay)
		if err != nil {
			c.logger.Error("Failed to create retry connection",
				zap.Error(err),
				zap.Any("options", options))
		} else {
			connection = retryConnection
		}
	}
	return &conn{Conn: connection, clientVersion: clientVersion, serverName: c.opts.serverName}
}

func (c *cluster) Stats() PoolStats {
	return poolStats(c.redisCluster.Stats())
}

func createPool(opts *options) func(string, ...redis.DialOption) (*redis.Pool, error) {
	return func(address string, options ...redis.DialOption) (*redis.Pool, error) {
		p := &redis.Pool{
			MaxIdle:     opts.poolMaxIdle,
			MaxActive:   opts.poolMaxActive,
			IdleTimeout: opts.poolIdleTimeout,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", address, options...)
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
		return p, nil
	}
}
