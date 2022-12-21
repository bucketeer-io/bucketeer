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

package v3

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/go-redis/redis/v8"
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
	incrByFloatCmdName  = "INCR_BY_FLOAT"
	delCmdName          = "DEL"
	incrCmdName         = "INCR"
	expireCmdName       = "EXPIRE"
	pipelineExecCmdName = "PIPELINE_EXEC"
	ttlCmdName          = "TTL"

	incrementWithExpireScript string = `
    local key =  KEYS[1]
    local result = redis.pcall("INCR", key)
    if type(result) == 'table' and result.err then
      redis.log(redis.LOG_WARNING, "failed to increment", key, result.err)
      return -1
    end
    if result == 1 then
      redis.log(redis.LOG_NOTICE, "setting expiration in the key", key)
      local expireInSeconds = ARGV[1]
      local res = redis.pcall("EXPIRE", key, expireInSeconds)
      if type(res) == 'table' and res.err then
        redis.log(redis.LOG_WARNING, "failed to set expiration", key, expireInSeconds, res.err)
        return -1
      end
    end
    return result
  `
)

var (
	ErrNil               = goredis.Nil
	ErrInvalidType       = errors.New("redis: invalid type")
	ErrFailedToIncrement = errors.New("redis: failed to increment and set expiration using the lua script")
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
	IncrByFloat(key string, value float64) (float64, error)
	Del(key string) error
	Incr(key string, expire time.Duration) (int64, error)
	Pipeline() PipeClient
	Expire(key string, seconds time.Duration) (bool, error)
}

type client struct {
	ctx    context.Context
	rc     *goredis.Client
	opts   *options
	logger *zap.Logger
}

type PipeClient interface {
	PFAdd(key string, els ...string) *goredis.IntCmd
	Incr(key string) *goredis.IntCmd
	TTL(key string) *goredis.DurationCmd
	Exec() ([]goredis.Cmder, error)
}

type pipeClient struct {
	ctx    context.Context
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

func NewClient(ctx context.Context, addr string, opts ...Option) (Client, error) {
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
	_, err := rc.Ping(ctx).Result()
	if err != nil {
		logger.Error("Failed to ping", zap.Error(err))
		return nil, err
	}
	client := &client{
		ctx:    ctx,
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
		_, err := c.rc.Ping(ctx).Result()
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
	keys, cursor, err := c.rc.Scan(c.ctx, cursor, key, count).Result()
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
	reply, err := c.rc.Get(c.ctx, key).Bytes()
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
	reply, err := c.rc.MGet(c.ctx, keys...).Result()
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
	err := c.rc.Set(c.ctx, key, val, expiration).Err()
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
	result, err := c.rc.PFAdd(c.ctx, key, els).Result()
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

func (c *client) PFCount(keys ...string) (int64, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, pfCountCmdName).Inc()
	count, err := c.rc.PFCount(c.ctx, keys...).Result()
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
	v, err := c.rc.IncrByFloat(c.ctx, key, value).Result()
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
	_, err := c.rc.Del(c.ctx, key).Result()
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

// Incr interface only set the expiration when the key is created.
// Once is created, the expiration setting will be ignored in the Lua script
func (c *client) Incr(key string, expiration time.Duration) (int64, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, incrCmdName).Inc()
	var val int64
	var err error
	if expiration < 1*time.Second {
		val, err = c.rc.Incr(c.ctx, key).Result()
	} else {
		incr := goredis.NewScript(incrementWithExpireScript)
		keys := []string{key}
		values := []interface{}{expiration.Seconds()}
		val, err = incr.Run(c.ctx, c.rc, keys, values...).Int64()
		if val == -1 {
			err = ErrFailedToIncrement
		}
	}
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, incrCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, incrCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return val, err
}

func (c *client) Expire(key string, seconds time.Duration) (bool, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, expireCmdName).Inc()
	v, err := c.rc.Expire(c.ctx, key, seconds).Result()
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
		ctx:    c.ctx,
		pipe:   c.rc.Pipeline(),
		cmds:   []string{},
		opts:   c.opts,
		logger: c.logger,
	}
}

func (c *pipeClient) Incr(key string) *goredis.IntCmd {
	c.cmds = append(c.cmds, incrCmdName)
	return c.pipe.Incr(c.ctx, key)
}

func (c *pipeClient) PFAdd(key string, els ...string) *goredis.IntCmd {
	c.cmds = append(c.cmds, pfAddCmdName)
	return c.pipe.PFAdd(c.ctx, key, els)
}

func (c *pipeClient) TTL(key string) *goredis.DurationCmd {
	c.cmds = append(c.cmds, ttlCmdName)
	return c.pipe.TTL(c.ctx, key)
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
	v, err := c.pipe.Exec(c.ctx)
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, cmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, cmdName, code).Observe(
		time.Since(startTime).Seconds())
	return v, err
}
