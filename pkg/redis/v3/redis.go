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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v3

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/redis"
)

const (
	clientVersion = "v3"

	scanCmdName         = "SCAN"
	getCmdName          = "GET"
	getMultiCmdName     = "GET_MULTI"
	setCmdName          = "SET"
	pfAddCmdName        = "PFADD"
	pfCountCmdName      = "PFCOUNT"
	pfMergeCmdName      = "PFMERGE"
	incrByFloatCmdName  = "INCR_BY_FLOAT"
	delCmdName          = "DEL"
	incrCmdName         = "INCR"
	expireCmdName       = "EXPIRE"
	pipelineExecCmdName = "PIPELINE_EXEC"
	ttlCmdName          = "TTL"
)

var (
	ErrNil         = goredis.Nil
	ErrInvalidType = errors.New("redis: invalid type")
)

type poolStats goredis.PoolStats

func (ps *poolStats) ActiveCount() int {
	return int(ps.TotalConns)
}

func (ps *poolStats) IdleCount() int {
	return int(ps.IdleConns)
}

type Client interface {
	Close() error
	Check(context.Context) health.Status
	Stats() redis.PoolStats
	Scan(cursor uint64, key string, count int64) (uint64, []string, error)
	Get(key string) ([]byte, error)
	GetMulti(keys []string) ([]interface{}, error)
	Set(key string, val interface{}, expiration time.Duration) error
	PFAdd(key string, els ...string) (int64, error)
	PFCount(keys ...string) (int64, error)
	PFMerge(dest string, keys ...string) error
	IncrByFloat(key string, value float64) (float64, error)
	Del(key string) error
	Incr(key string) (int64, error)
	Pipeline() PipeClient
	Expire(key string, expiration time.Duration) (bool, error)
}

type client struct {
	rc     *goredis.Client
	opts   *options
	logger *zap.Logger
}

type PipeClient interface {
	PFAdd(key string, els ...string) *goredis.IntCmd
	Incr(key string) *goredis.IntCmd
	TTL(key string) *goredis.DurationCmd
	Exec() ([]goredis.Cmder, error)
	PFCount(keys ...string) *goredis.IntCmd
	Get(key string) *goredis.StringCmd
	Del(keys string) *goredis.IntCmd
}

type pipeClient struct {
	pipe   goredis.Pipeliner
	cmds   []string
	opts   *options
	logger *zap.Logger
}

type options struct {
	password     string
	maxRetries   int
	dialTimeout  time.Duration
	poolSize     int
	minIdleConns int
	poolTimeout  time.Duration
	serverName   string
	metrics      metrics.Registerer
	logger       *zap.Logger
}

func defaultOptions() *options {
	return &options{
		maxRetries:   5,
		dialTimeout:  5 * time.Second,
		poolSize:     10,
		minIdleConns: 5,
		poolTimeout:  5 * time.Second,
		logger:       zap.NewNop(),
	}
}

type Option func(*options)

func WithPassword(password string) Option {
	return func(opts *options) {
		opts.password = password
	}
}

func WithMaxRetries(maxRetries int) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(opts *options) {
		opts.dialTimeout = dialTimeout
	}
}

func WithPoolSize(poolSize int) Option {
	return func(opts *options) {
		opts.poolSize = poolSize
	}
}

func WithMinIdleConns(minIdleConns int) Option {
	return func(opts *options) {
		opts.minIdleConns = minIdleConns
	}
}

func WithPoolTimeout(poolTimeout time.Duration) Option {
	return func(opts *options) {
		opts.poolTimeout = poolTimeout
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

func NewClient(addr string, opts ...Option) (Client, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	logger := options.logger.Named("redis-v3")
	rc := goredis.NewClient(&goredis.Options{
		Addr:         addr,
		Password:     options.password,
		MaxRetries:   options.maxRetries,
		DialTimeout:  options.dialTimeout,
		PoolSize:     options.poolSize,
		MinIdleConns: options.minIdleConns,
		PoolTimeout:  options.poolTimeout,
	})
	_, err := rc.Ping().Result()
	if err != nil {
		logger.Error("Failed to ping", zap.Error(err))
		return nil, err
	}
	client := &client{
		rc:     rc,
		opts:   options,
		logger: logger,
	}
	if options.metrics != nil {
		redis.RegisterMetrics(options.metrics, clientVersion, options.serverName, client)
	}
	return client, nil
}

func (c *client) Close() error {
	return c.rc.Close()
}

func (c *client) Check(ctx context.Context) health.Status {
	resultCh := make(chan health.Status, 1)
	go func() {
		_, err := c.rc.Ping().Result()
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

func (c *client) Stats() redis.PoolStats {
	return (*poolStats)(c.rc.PoolStats())
}

func (c *client) Scan(cursor uint64, key string, count int64) (uint64, []string, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, scanCmdName).Inc()
	keys, cursor, err := c.rc.Scan(cursor, key, count).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	case ErrNil:
		code = redis.CodeNotFound
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, scanCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, scanCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return cursor, keys, err
}

func (c *client) Get(key string) ([]byte, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, getCmdName).Inc()
	reply, err := c.rc.Get(key).Bytes()
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

func (c *client) GetMulti(keys []string) ([]interface{}, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, getMultiCmdName).Inc()
	reply, err := c.rc.MGet(keys...).Result()
	code := redis.CodeFail
	values := make([]interface{}, 0, len(reply))
	switch err {
	case nil:
		code = redis.CodeSuccess
		for _, r := range reply {
			s, ok := r.(string)
			if !ok {
				code = redis.CodeInvalidType
				values = nil
				err = ErrInvalidType
				break
			}
			values = append(values, []byte(s))
		}
	case ErrNil:
		code = redis.CodeNotFound
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, getMultiCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, getMultiCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return values, err
}

func (c *client) Set(key string, val interface{}, expiration time.Duration) error {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, setCmdName).Inc()
	err := c.rc.Set(key, val, expiration).Err()
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

func (c *client) PFAdd(key string, els ...string) (int64, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, pfAddCmdName).Inc()
	result, err := c.rc.PFAdd(key, els).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, pfAddCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, pfAddCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return result, err
}

func (c *client) PFMerge(dest string, keys ...string) error {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, pfMergeCmdName).Inc()
	_, err := c.rc.PFMerge(dest, keys...).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, pfMergeCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, pfMergeCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return err
}

func (c *client) PFCount(keys ...string) (int64, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, pfCountCmdName).Inc()
	count, err := c.rc.PFCount(keys...).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, pfCountCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, pfCountCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return count, err
}

func (c *client) IncrByFloat(key string, value float64) (float64, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, incrByFloatCmdName).Inc()
	v, err := c.rc.IncrByFloat(key, value).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, incrByFloatCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, incrByFloatCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return v, err
}

func (c *client) Del(key string) error {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, delCmdName).Inc()
	_, err := c.rc.Del(key).Result()
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

func (c *client) Incr(key string) (int64, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, incrCmdName).Inc()
	v, err := c.rc.Incr(key).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, incrCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, incrCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return v, err
}

func (c *client) Expire(key string, expiration time.Duration) (bool, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, expireCmdName).Inc()
	v, err := c.rc.Expire(key, expiration).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, expireCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, expireCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return v, err
}

func (c *client) Pipeline() PipeClient {
	return &pipeClient{
		pipe:   c.rc.Pipeline(),
		cmds:   []string{},
		opts:   c.opts,
		logger: c.logger,
	}
}

func (c *pipeClient) Incr(key string) *goredis.IntCmd {
	c.cmds = append(c.cmds, incrCmdName)
	return c.pipe.Incr(key)
}

func (c *pipeClient) PFAdd(key string, els ...string) *goredis.IntCmd {
	c.cmds = append(c.cmds, pfAddCmdName)
	return c.pipe.PFAdd(key, els)
}

func (c *pipeClient) TTL(key string) *goredis.DurationCmd {
	c.cmds = append(c.cmds, ttlCmdName)
	return c.pipe.TTL(key)
}

// The command name reported in the metrics handler counter
// is based on how many commands were used in the pipeline
func (c *pipeClient) Exec() ([]goredis.Cmder, error) {
	startTime := time.Now()
	cmdName := pipelineExecCmdName
	for _, cmd := range c.cmds {
		cmdName += fmt.Sprintf("_%s", cmd)
	}
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, cmdName).Inc()
	v, err := c.pipe.Exec()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	case ErrNil:
		code = redis.CodeNotFound
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, cmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, cmdName, code).Observe(
		time.Since(startTime).Seconds())
	return v, err
}

func (c *pipeClient) PFCount(keys ...string) *goredis.IntCmd {
	c.cmds = append(c.cmds, pfCountCmdName)
	return c.pipe.PFCount(keys...)
}

func (c *pipeClient) Get(key string) *goredis.StringCmd {
	c.cmds = append(c.cmds, getCmdName)
	return c.pipe.Get(key)
}

func (c *pipeClient) Del(key string) *goredis.IntCmd {
	c.cmds = append(c.cmds, pfCountCmdName)
	return c.pipe.Del(key)
}
