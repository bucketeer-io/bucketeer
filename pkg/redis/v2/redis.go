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

package v2

import (
	"context"
	"time"

	goredis "github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/redis"
)

const (
	clientVersion = "v2"

	getCmdName           = "GET"
	setCmdName           = "SET"
	delCmdName           = "DEL"
	forEachMasterCmdName = "FOR_EACH_MASTER"
)

var ErrNil = goredis.Nil

type poolStats goredis.PoolStats

func (ps *poolStats) ActiveCount() int {
	return int(ps.TotalConns)
}

func (ps *poolStats) IdleCount() int {
	return int(ps.IdleConns)
}

type Cluster interface {
	Get(key string) ([]byte, error)
	Set(key string, val interface{}, expiration time.Duration) error
	Del(key string) error
	ForEachMaster(fn func(client *goredis.Client) error) error
	Check(context.Context) health.Status
	Stats() redis.PoolStats
	Close() error
}

type cluster struct {
	cc     *goredis.ClusterClient
	opts   *options
	logger *zap.Logger
}

type options struct {
	dialPassword       string
	dialConnectTimeout time.Duration
	poolMaxIdle        int
	poolIdleTimeout    time.Duration
	routeByLatency     bool
	serverName         string
	metrics            metrics.Registerer
	logger             *zap.Logger
}

func defaultOptions() *options {
	return &options{
		dialConnectTimeout: 5 * time.Second,
		poolMaxIdle:        5,
		poolIdleTimeout:    time.Minute,
		routeByLatency:     true,
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

func NewCluster(nodes []string, opts ...Option) (Cluster, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	logger := options.logger.Named("redis-v2")
	cc := goredis.NewClusterClient(&goredis.ClusterOptions{
		Password:       options.dialPassword,
		Addrs:          nodes,
		IdleTimeout:    options.poolIdleTimeout,
		DialTimeout:    options.dialConnectTimeout,
		RouteByLatency: options.routeByLatency,
	})
	if err := cc.ReloadState(); err != nil {
		logger.Error("Failed to refresh", zap.Error(err))
		return nil, err
	}
	cluster := &cluster{
		cc:     cc,
		opts:   options,
		logger: logger,
	}
	if options.metrics != nil {
		redis.RegisterMetrics(options.metrics, clientVersion, options.serverName, cluster)
	}
	cluster.logger.Debug("redis/v2 client was initialized")
	return cluster, nil
}

func (c *cluster) Get(key string) ([]byte, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, getCmdName).Inc()
	reply, err := c.cc.Get(key).Bytes()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	case ErrNil:
		code = redis.CodeNotFound
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, getCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, getCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return reply, err
}

func (c *cluster) Set(key string, val interface{}, expiration time.Duration) error {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, setCmdName).Inc()
	_, err := c.cc.Set(key, val, expiration).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	case ErrNil:
		code = redis.CodeNotFound
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, setCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, setCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return err
}

func (c *cluster) Del(key string) error {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, delCmdName).Inc()
	_, err := c.cc.Del(key).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	case ErrNil:
		code = redis.CodeNotFound
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, delCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, delCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return err
}

func (c *cluster) ForEachMaster(fn func(client *goredis.Client) error) error {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, forEachMasterCmdName).Inc()
	err := c.cc.ForEachMaster(fn)
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	case ErrNil:
		code = redis.CodeNotFound
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, forEachMasterCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, forEachMasterCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return err
}

func (c *cluster) Check(ctx context.Context) health.Status {
	resultCh := make(chan health.Status, 1)
	go func() {
		_, err := c.cc.Ping().Result()
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

func (c *cluster) Stats() redis.PoolStats {
	return (*poolStats)(c.cc.PoolStats())
}

func (c *cluster) Close() error {
	return c.cc.Close()
}
