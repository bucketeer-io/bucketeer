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

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
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
	SetNXCmdName        = "SETNX"
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

type ClientType int

const (
	ClientTypeStandard ClientType = iota
	ClientTypeCluster
)

type Client interface {
	Close() error
	Check(context.Context) health.Status
	Stats() redis.PoolStats
	Scan(cursor uint64, key string, count int64) (uint64, []string, error)
	Get(key string) ([]byte, error)
	GetMulti(keys []string, ignoreNotFound bool) ([]interface{}, error)
	Set(key string, val interface{}, expiration time.Duration) error
	PFAdd(key string, els ...string) (int64, error)
	PFCount(keys ...string) (int64, error)
	PFMerge(dest string, expiration time.Duration, keys ...string) error
	IncrByFloat(key string, value float64) (float64, error)
	Del(key string) error
	Incr(key string) (int64, error)
	Pipeline(tx bool) PipeClient
	Expire(key string, expiration time.Duration) (bool, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *goredis.Cmd
	Dump(key string) (string, error)
	Restore(key string, ttl int64, value string) error
	Exists(key string) (int64, error)
}

type client struct {
	rc         goredis.UniversalClient
	opts       *options
	logger     *zap.Logger
	clientType ClientType
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

func NewClient(addr string, opts ...Option) (Client, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	logger := options.logger.Named("redis-v3")

	standardClientOpts := &goredis.Options{
		Addr:         addr,
		Password:     options.password,
		MaxRetries:   options.maxRetries,
		DialTimeout:  options.dialTimeout,
		PoolSize:     options.poolSize,
		MinIdleConns: options.minIdleConns,
		PoolTimeout:  options.poolTimeout,
	}

	tmpClient := goredis.NewClient(standardClientOpts)
	defer tmpClient.Close()

	var rc goredis.UniversalClient
	var clientType ClientType
	if _, err := tmpClient.ClusterInfo(context.TODO()).Result(); err == nil {
		logger.Debug("Redis cluster detected, creating cluster client",
			zap.String("addr", addr),
			zap.Int("maxRetries", options.maxRetries),
			zap.Duration("dialTimeout", options.dialTimeout),
			zap.Int("poolSize", options.poolSize),
			zap.Int("minIdleConns", options.minIdleConns),
			zap.Duration("poolTimeout", options.poolTimeout),
		)
		rc = goredis.NewClusterClient(&goredis.ClusterOptions{
			Addrs:        []string{addr},
			Password:     options.password,
			MaxRetries:   options.maxRetries,
			DialTimeout:  options.dialTimeout,
			PoolSize:     options.poolSize,
			MinIdleConns: options.minIdleConns,
			PoolTimeout:  options.poolTimeout,
		})
		clientType = ClientTypeCluster
	} else {
		logger.Debug("Redis standalone detected, creating standard client",
			zap.String("addr", addr),
			zap.Int("maxRetries", options.maxRetries),
			zap.Duration("dialTimeout", options.dialTimeout),
			zap.Int("poolSize", options.poolSize),
			zap.Int("minIdleConns", options.minIdleConns),
			zap.Duration("poolTimeout", options.poolTimeout),
		)
		rc = goredis.NewClient(standardClientOpts)
		clientType = ClientTypeStandard
	}
	_, err := rc.Ping(context.TODO()).Result()
	if err != nil {
		logger.Error("Failed to ping", zap.Error(err))
		return nil, err
	}
	client := &client{
		rc:         rc,
		opts:       options,
		logger:     logger,
		clientType: clientType,
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
		_, err := c.rc.Ping(context.TODO()).Result()
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
	keys, cursor, err := c.rc.Scan(context.TODO(), cursor, key, count).Result()
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
	reply, err := c.rc.Get(context.TODO(), key).Bytes()
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

func (c *client) GetMulti(keys []string, ignoreNotFound bool) ([]interface{}, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, getMultiCmdName).Inc()

	var reply []interface{}
	var err error

	if c.clientType == ClientTypeCluster {
		// Use cluster-aware approach
		slotMap := make(map[int][]string)
		for _, key := range keys {
			slot := keyHashSlot(key)
			slotMap[slot] = append(slotMap[slot], key)
		}

		// Create a map to store results temporarily
		resultMap := make(map[string]interface{}, len(keys))

		for _, slotKeys := range slotMap {
			slotReply, slotErr := c.rc.MGet(context.TODO(), slotKeys...).Result()
			if slotErr != nil {
				err = slotErr
				break
			}
			// Store results in the map
			for i, key := range slotKeys {
				resultMap[key] = slotReply[i]
			}
		}

		if err == nil {
			// Populate the reply slice in the original order
			reply = make([]interface{}, len(keys))
			for i, key := range keys {
				reply[i] = resultMap[key]
			}
		}
	} else {
		// Use standard approach for non-cluster client
		reply, err = c.rc.MGet(context.TODO(), keys...).Result()
	}

	code := redis.CodeFail
	values := make([]interface{}, 0, len(reply))
	if err == nil {
		code = redis.CodeSuccess
		for _, r := range reply {
			if r == nil {
				if ignoreNotFound {
					values = append(values, "")
					continue
				}
				code = redis.CodeNotFound
				err = ErrNil
				values = nil
				break
			}
			s, ok := r.(string)
			if !ok {
				code = redis.CodeInvalidType
				err = ErrInvalidType
				values = nil
				break
			}
			values = append(values, []byte(s))
		}
	} else if err == ErrNil {
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
	err := c.rc.Set(context.TODO(), key, val, expiration).Err()
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
	result, err := c.rc.PFAdd(context.TODO(), key, els).Result()
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

func (c *client) PFMerge(dest string, expiration time.Duration, keys ...string) error {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, pfMergeCmdName).Inc()

	var err error
	ctx := context.TODO()

	// Metrics reporting deferred at the end of the function
	defer func() {
		code := redis.CodeFail
		if err == nil {
			code = redis.CodeSuccess
		}
		// Reporting metrics
		redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, pfMergeCmdName, code).Inc()
		redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, pfMergeCmdName, code).Observe(
			time.Since(startTime).Seconds())
	}()

	if c.clientType == ClientTypeCluster {
		// Cluster-aware approach
		var allK []string

		// Fetch all HLL objects via GetMulti and store them client side as strings
		hllObjects, err := c.GetMulti(keys, true)
		if err != nil {
			return err
		}
		allHLLObjects := make([]string, 0, len(hllObjects))
		for _, hllObj := range hllObjects {
			if obj, ok := hllObj.([]byte); ok {
				allHLLObjects = append(allHLLObjects, string(obj))
			}
		}

		// Randomize a keyslot hash
		randomHashSlot := c.randomID()

		// Special handling of dest variable if it already exists
		destData, err := c.rc.Get(ctx, dest).Result()
		if err != nil && !errors.Is(err, goredis.Nil) {
			return err
		}
		if !errors.Is(err, goredis.Nil) {
			allHLLObjects = append(allHLLObjects, destData)
		}

		// MSet all stored HLL objects with {RandomHash}RandomKey hll_obj
		pairs := make([]interface{}, 0, len(allHLLObjects)*2)
		for _, hllObject := range allHLLObjects {
			k := c.randomHashSlotKey(randomHashSlot)
			allK = append(allK, k)
			pairs = append(pairs, k, hllObject)
		}
		if len(pairs) > 0 {
			err = c.rc.MSet(ctx, pairs...).Err()
			if err != nil {
				return err
			}
		}

		// Do regular PFMERGE operation and store value in random key in {RandomHash}
		tmpDest := c.randomHashSlotKey(randomHashSlot)
		_, err = c.rc.PFMerge(ctx, tmpDest, allK...).Result()
		if err != nil {
			return err
		}

		// Do GET and SET so that result will be stored in the destination object anywhere in the cluster
		parsedDest, err := c.rc.Get(ctx, tmpDest).Result()
		if err != nil {
			return err
		}
		err = c.rc.Set(ctx, dest, parsedDest, expiration).Err()
		if err != nil {
			return err
		}

		// Cleanup temporary keys
		pipe := c.rc.Pipeline()
		for _, k := range allK {
			pipe.Del(ctx, k)
		}
		pipe.Del(ctx, tmpDest)

		// Execute the pipeline
		_, err = pipe.Exec(ctx)
		if err != nil {
			return err
		}
	} else {
		// Standard non-cluster approach
		// We use transaction here to ensure the expiration will be set when merging the key
		pipe := c.rc.TxPipeline()
		pipe.PFMerge(ctx, dest, keys...)
		pipe.Expire(ctx, dest, expiration)

		// Execute the pipeline
		_, err = pipe.Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *client) randomID() string {
	return uuid.New().String()
}

func (c *client) randomHashSlotKey(randomHashSlot string) string {
	return fmt.Sprintf("{%s}%s", randomHashSlot, uuid.New().String())
}

func (c *client) PFCount(keys ...string) (int64, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, pfCountCmdName).Inc()

	var count int64
	var err error

	if c.clientType == ClientTypeCluster {
		// Use cluster-aware approach
		slotMap := make(map[int][]string)
		for _, key := range keys {
			slot := keyHashSlot(key)
			slotMap[slot] = append(slotMap[slot], key)
		}

		for _, slotKeys := range slotMap {
			slotCount, slotErr := c.rc.PFCount(context.TODO(), slotKeys...).Result()
			if slotErr != nil {
				err = slotErr
				break
			}
			count += slotCount
		}
	} else {
		// Use standard approach for non-cluster client
		count, err = c.rc.PFCount(context.TODO(), keys...).Result()
	}

	code := redis.CodeFail
	if err == nil {
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
	v, err := c.rc.IncrByFloat(context.TODO(), key, value).Result()
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
	_, err := c.rc.Del(context.TODO(), key).Result()
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
	v, err := c.rc.Incr(context.TODO(), key).Result()
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
	v, err := c.rc.Expire(context.TODO(), key, expiration).Result()
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

func (c *client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, SetNXCmdName).Inc()

	result, err := c.rc.SetNX(context.TODO(), key, value, expiration).Result()

	code := redis.CodeFail
	if err == nil {
		code = redis.CodeSuccess
	}

	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, SetNXCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, SetNXCmdName, code).Observe(
		time.Since(startTime).Seconds())

	return result, err
}

func (c *client) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *goredis.Cmd {
	return c.rc.Eval(context.TODO(), script, keys, args...)
}

func (c *client) Pipeline(tx bool) PipeClient {
	var pipeliner goredis.Pipeliner
	if tx {
		// All commands in the transaction will either all succeed or none will be applied if an error occurs.
		pipeliner = c.rc.TxPipeline()
	} else {
		pipeliner = c.rc.Pipeline()
	}
	return &pipeClient{
		ctx:    context.TODO(),
		pipe:   pipeliner,
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
	return c.pipe.PFCount(c.ctx, keys...)
}

func (c *pipeClient) Get(key string) *goredis.StringCmd {
	c.cmds = append(c.cmds, getCmdName)
	return c.pipe.Get(c.ctx, key)
}

func (c *pipeClient) Del(key string) *goredis.IntCmd {
	c.cmds = append(c.cmds, pfCountCmdName)
	return c.pipe.Del(c.ctx, key)
}

func (c *client) Dump(key string) (string, error) {
	result, err := c.rc.Dump(context.TODO(), key).Result()
	return result, err
}

func (c *client) Restore(key string, ttl int64, value string) error {
	_, err := c.rc.Restore(context.TODO(), key, time.Duration(ttl)*time.Millisecond, value).Result()
	return err
}

func (c *client) Exists(key string) (int64, error) {
	return c.rc.Exists(context.TODO(), key).Result()
}
