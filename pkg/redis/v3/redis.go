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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v3

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/health"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/redis"
)

const (
	clientVersion = "v3"

	scanCmdName          = "SCAN"
	getCmdName           = "GET"
	getMultiCmdName      = "GET_MULTI"
	setCmdName           = "SET"
	pfAddCmdName         = "PFADD"
	pfCountCmdName       = "PFCOUNT"
	pfMergeCmdName       = "PFMERGE"
	incrByFloatCmdName   = "INCR_BY_FLOAT"
	delCmdName           = "DEL"
	incrCmdName          = "INCR"
	incrByCmdName        = "INCRBY"
	expireCmdName        = "EXPIRE"
	pipelineExecCmdName  = "PIPELINE_EXEC"
	ttlCmdName           = "TTL"
	SetNXCmdName         = "SETNX"
	saddCmdName          = "SADD"
	smembersCmdName      = "SMEMBERS"
	xAddCmdName          = "XADD"
	xGroupCreateCmdName  = "XGROUP_CREATE"
	xReadGroupCmdName    = "XREADGROUP"
	xAckCmdName          = "XACK"
	xPendingCmdName      = "XPENDING"
	xClaimCmdName        = "XCLAIM"
	xInfoGroupsCmdName   = "XINFO_GROUPS"
	xGroupDestroyCmdName = "XGROUP_DESTROY"
)

// RedisMode specifies how the Redis client should be created.
// The default is RedisModeAuto, which detects the Redis deployment mode
// with mismatch detection. For stricter production setups, you can
// explicitly set RedisModeCluster or RedisModeStandalone.
type RedisMode string

const (
	RedisModeCluster    RedisMode = "cluster"
	RedisModeStandalone RedisMode = "standalone"
	RedisModeAuto       RedisMode = "auto"

	mismatchCheckInterval = 5 * time.Minute
	detectionTimeout      = 3 * time.Second
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
	IncrBy(key string, value int64) (int64, error)
	SAdd(key string, members ...interface{}) (int64, error)
	SMembers(key string) ([]string, error)
	Pipeline(tx bool) PipeClient
	Expire(key string, expiration time.Duration) (bool, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *goredis.Cmd
	Dump(key string) (string, error)
	Restore(key string, ttl int64, value string) error
	Exists(key string) (int64, error)

	// Redis Stream methods
	XAdd(ctx context.Context, stream string, values map[string]interface{}) (string, error)
	XGroupCreateMkStream(stream, group, start string) error
	XReadGroup(ctx context.Context,
		group, consumer string,
		streams []string,
		count int64,
		block time.Duration,
	) ([]goredis.XStream, error)
	XAck(stream, group, id string) error
	XPendingExt(ctx context.Context,
		stream, group, start, end string,
		count int64,
		idle time.Duration,
	) ([]goredis.XPendingExt, error)
	XClaim(ctx context.Context,
		stream, group, consumer string,
		minIdle time.Duration,
		ids []string) ([]goredis.XMessage, error)
	XInfoGroups(ctx context.Context, stream string) ([]goredis.XInfoGroup, error)
	XGroupDestroy(ctx context.Context, stream, group string) error
}

type client struct {
	rc         goredis.UniversalClient
	opts       *options
	logger     *zap.Logger
	clientType ClientType
	done       chan struct{}
}

type PipeClient interface {
	PFAdd(key string, els ...string) *goredis.IntCmd
	PFMerge(dest string, keys ...string) *goredis.StatusCmd
	Incr(key string) *goredis.IntCmd
	IncrBy(key string, value int64) *goredis.IntCmd
	TTL(key string) *goredis.DurationCmd
	SAdd(key string, members ...interface{}) *goredis.IntCmd
	Expire(key string, expiration time.Duration) *goredis.BoolCmd
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
	redisMode    RedisMode
	metrics      metrics.Registerer
	logger       *zap.Logger
}

func defaultOptions() *options {
	return &options{
		maxRetries:   5,
		dialTimeout:  5 * time.Second,
		poolSize:     10,
		minIdleConns: 5,
		poolTimeout:  30 * time.Second,
		redisMode:    RedisModeAuto,
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

// WithRedisMode sets the Redis client creation mode.
// "cluster" always creates a ClusterClient, "standalone" always creates a standard Client,
// "auto" (default) tries detection and runs background mismatch checking.
func WithRedisMode(mode RedisMode) Option {
	return func(opts *options) {
		m := RedisMode(strings.ToLower(string(mode)))
		switch m {
		case RedisModeCluster, RedisModeStandalone, RedisModeAuto:
			opts.redisMode = m
		default:
			opts.redisMode = RedisModeAuto
		}
	}
}

func NewClient(addr string, opts ...Option) (Client, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	logger := options.logger.Named("redis-v3")

	clusterOpts := &goredis.ClusterOptions{
		Addrs:        []string{addr},
		Password:     options.password,
		MaxRetries:   options.maxRetries,
		DialTimeout:  options.dialTimeout,
		PoolSize:     options.poolSize,
		MinIdleConns: options.minIdleConns,
		PoolTimeout:  options.poolTimeout,
	}
	standardOpts := &goredis.Options{
		Addr:         addr,
		Password:     options.password,
		MaxRetries:   options.maxRetries,
		DialTimeout:  options.dialTimeout,
		PoolSize:     options.poolSize,
		MinIdleConns: options.minIdleConns,
		PoolTimeout:  options.poolTimeout,
	}

	var rc goredis.UniversalClient
	var clientType ClientType

	switch options.redisMode {
	case RedisModeCluster:
		rc = goredis.NewClusterClient(clusterOpts)
		clientType = ClientTypeCluster
		logger.Info("Creating Redis cluster client (explicit mode)",
			zap.String("addr", addr),
			zap.String("mode", string(options.redisMode)),
			zap.String("clientType", clientTypeString(clientType)),
		)

	case RedisModeStandalone:
		rc = goredis.NewClient(standardOpts)
		clientType = ClientTypeStandard
		logger.Info("Creating Redis standalone client (explicit mode)",
			zap.String("addr", addr),
			zap.String("mode", string(options.redisMode)),
			zap.String("clientType", clientTypeString(clientType)),
		)

	default: // RedisModeAuto
		clientType, rc = detectRedisMode(addr, clusterOpts, standardOpts, logger)
	}

	// Non-blocking startup: try to ping but don't fail if Redis is unavailable.
	// This allows the service to start in degraded mode using database fallback,
	// while automatically reconnecting when Redis becomes available.
	ctx, cancel := context.WithTimeout(context.Background(), detectionTimeout)
	defer cancel()

	if _, err := rc.Ping(ctx).Result(); err != nil {
		logger.Warn("Redis not available at startup, service will retry on first use",
			zap.Error(err),
			zap.String("addr", addr),
			zap.String("mode", string(options.redisMode)),
		)
	}

	c := &client{
		rc:         rc,
		opts:       options,
		logger:     logger,
		clientType: clientType,
		done:       make(chan struct{}),
	}
	if options.metrics != nil {
		redis.RegisterMetrics(options.metrics, clientVersion, options.serverName, c)
	}

	// In auto mode, run a background goroutine that periodically checks
	// whether the detected mode matches the actual Redis topology.
	if options.redisMode == RedisModeAuto {
		go c.runMismatchDetector(addr)
	}

	return c, nil
}

// detectRedisMode tries to determine whether the Redis server is a cluster or standalone
// by issuing CLUSTER INFO with a short timeout. Falls back to standalone if detection fails.
// Note: if Redis is unreachable at startup and the actual topology is a cluster,
// the standalone fallback will receive MOVED/ASK errors once the cluster recovers.
// For cluster deployments, prefer setting RedisMode to "cluster" explicitly.
func detectRedisMode(
	addr string,
	clusterOpts *goredis.ClusterOptions,
	standardOpts *goredis.Options,
	logger *zap.Logger,
) (ClientType, goredis.UniversalClient) {
	probe := goredis.NewClient(standardOpts)
	defer probe.Close()

	ctx, cancel := context.WithTimeout(context.Background(), detectionTimeout)
	defer cancel()

	if isCluster, _ := probeClusterMode(ctx, probe); isCluster {
		logger.Info("Redis cluster detected (auto mode)",
			zap.String("addr", addr),
		)
		return ClientTypeCluster, goredis.NewClusterClient(clusterOpts)
	}

	logger.Info("Redis standalone detected or unavailable, defaulting to standalone (auto mode)",
		zap.String("addr", addr),
	)
	return ClientTypeStandard, goredis.NewClient(standardOpts)
}

// probeClusterMode checks CLUSTER INFO output to determine if the server is in cluster mode.
// Returns (true, nil) for cluster, (false, nil) for standalone, (false, err) on failure.
func probeClusterMode(ctx context.Context, c *goredis.Client) (bool, error) {
	info, err := c.ClusterInfo(ctx).Result()
	if err != nil {
		return false, err
	}
	return strings.Contains(info, "cluster_enabled:1"), nil
}

// runMismatchDetector periodically checks if the configured client type matches
// the actual Redis topology. Logs a warning on mismatch.
func (c *client) runMismatchDetector(addr string) {
	ticker := time.NewTicker(mismatchCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
		}

		ctx, cancel := context.WithTimeout(context.Background(), detectionTimeout)
		probe := goredis.NewClient(&goredis.Options{
			Addr:        addr,
			Password:    c.opts.password,
			DialTimeout: c.opts.dialTimeout,
		})

		actualCluster := false
		if isCluster, err := probeClusterMode(ctx, probe); err == nil {
			actualCluster = isCluster
		} else {
			probe.Close()
			cancel()
			continue
		}
		probe.Close()
		cancel()

		configuredCluster := c.clientType == ClientTypeCluster
		if actualCluster != configuredCluster {
			expected := clientTypeString(c.clientType)
			actual := "standalone"
			if actualCluster {
				actual = "cluster"
			}
			c.logger.Warn("Redis mode mismatch detected",
				zap.String("addr", addr),
				zap.String("configured", expected),
				zap.String("actual", actual),
			)
		}
	}
}

func clientTypeString(ct ClientType) string {
	if ct == ClientTypeCluster {
		return "cluster"
	}
	return "standalone"
}

func (c *client) Close() error {
	close(c.done)
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

func (c *client) SAdd(key string, members ...interface{}) (int64, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, saddCmdName).Inc()
	result, err := c.rc.SAdd(context.TODO(), key, members...).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, saddCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, saddCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return result, err
}

func (c *client) SMembers(key string) ([]string, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, smembersCmdName).Inc()
	result, err := c.rc.SMembers(context.TODO(), key).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, smembersCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, smembersCmdName, code).Observe(
		time.Since(startTime).Seconds())
	return result, err
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

func (c *client) IncrBy(key string, value int64) (int64, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, incrByCmdName).Inc()
	v, err := c.rc.IncrBy(context.TODO(), key, value).Result()
	code := redis.CodeFail
	switch err {
	case nil:
		code = redis.CodeSuccess
	}
	redis.HandledCounter.WithLabelValues(clientVersion, c.opts.serverName, incrByCmdName, code).Inc()
	redis.HandledHistogram.WithLabelValues(clientVersion, c.opts.serverName, incrByCmdName, code).Observe(
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

func (c *pipeClient) IncrBy(key string, value int64) *goredis.IntCmd {
	c.cmds = append(c.cmds, incrByCmdName)
	return c.pipe.IncrBy(c.ctx, key, value)
}

func (c *pipeClient) PFAdd(key string, els ...string) *goredis.IntCmd {
	c.cmds = append(c.cmds, pfAddCmdName)
	return c.pipe.PFAdd(c.ctx, key, els)
}

func (c *pipeClient) PFMerge(dest string, keys ...string) *goredis.StatusCmd {
	c.cmds = append(c.cmds, pfMergeCmdName)
	return c.pipe.PFMerge(c.ctx, dest, keys...)
}

func (c *pipeClient) TTL(key string) *goredis.DurationCmd {
	c.cmds = append(c.cmds, ttlCmdName)
	return c.pipe.TTL(c.ctx, key)
}

func (c *pipeClient) SAdd(key string, members ...interface{}) *goredis.IntCmd {
	c.cmds = append(c.cmds, saddCmdName)
	return c.pipe.SAdd(c.ctx, key, members...)
}

func (c *pipeClient) Expire(key string, expiration time.Duration) *goredis.BoolCmd {
	c.cmds = append(c.cmds, expireCmdName)
	return c.pipe.Expire(c.ctx, key, expiration)
}

func (c *pipeClient) Exec() ([]goredis.Cmder, error) {
	startTime := time.Now()
	cmdName := pipelineExecCmdName
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
	c.cmds = append(c.cmds, delCmdName)
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

// Helper function to convert error to code string for metrics
func convertErrorToMetricsCode(err error) string {
	if err == nil {
		return redis.CodeSuccess
	}
	if errors.Is(err, ErrNil) {
		return redis.CodeNotFound
	}
	if errors.Is(err, ErrInvalidType) {
		return redis.CodeInvalidType
	}
	return redis.CodeFail
}

// XAdd adds a message to a stream
func (c *client) XAdd(ctx context.Context, stream string, values map[string]interface{}) (string, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, xAddCmdName).Inc()

	// Create XAddArgs with auto-generated ID and values
	args := &goredis.XAddArgs{
		Stream: stream,
		Values: values,
		ID:     "*", // Auto-generate ID
	}

	// Execute XAdd command
	cmd := c.rc.XAdd(ctx, args)
	id, err := cmd.Result()

	code := convertErrorToMetricsCode(err)
	redis.HandledCounter.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xAddCmdName,
		code,
	).Inc()
	redis.HandledHistogram.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xAddCmdName,
		code,
	).Observe(time.Since(startTime).Seconds())

	if err != nil {
		c.logger.Error("Failed to add message to stream",
			zap.String("stream", stream),
			zap.Error(err),
		)
		return "", err
	}

	return id, nil
}

// XGroupCreateMkStream creates a consumer group, creating the stream if it doesn't exist
func (c *client) XGroupCreateMkStream(stream, group, start string) error {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, xGroupCreateCmdName).Inc()

	// Execute XGroupCreateMkStream command
	cmd := c.rc.XGroupCreateMkStream(context.Background(), stream, group, start)
	err := cmd.Err()

	code := convertErrorToMetricsCode(err)
	redis.HandledCounter.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xGroupCreateCmdName,
		code,
	).Inc()
	redis.HandledHistogram.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xGroupCreateCmdName,
		code,
	).Observe(time.Since(startTime).Seconds())

	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		c.logger.Error("Failed to create consumer group",
			zap.String("stream", stream),
			zap.String("group", group),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// XReadGroup reads messages from a stream using a consumer group
func (c *client) XReadGroup(ctx context.Context,
	group, consumer string,
	streams []string,
	count int64,
	block time.Duration,
) ([]goredis.XStream, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, xReadGroupCmdName).Inc()

	// Create XReadGroupArgs
	args := &goredis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  streams,
		Count:    count,
		Block:    block,
	}

	// Execute XReadGroup command
	cmd := c.rc.XReadGroup(ctx, args)
	result, err := cmd.Result()

	code := convertErrorToMetricsCode(err)
	redis.HandledCounter.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xReadGroupCmdName,
		code,
	).Inc()
	redis.HandledHistogram.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xReadGroupCmdName,
		code,
	).Observe(time.Since(startTime).Seconds())

	if err != nil && err != goredis.Nil {
		c.logger.Error("Failed to read from stream",
			zap.String("group", group),
			zap.String("consumer", consumer),
			zap.Strings("streams", streams),
			zap.Error(err),
		)
		return nil, err
	}

	return result, err
}

// XAck acknowledges a message in a consumer group
func (c *client) XAck(stream, group, id string) error {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, xAckCmdName).Inc()

	// Execute XAck command
	cmd := c.rc.XAck(context.Background(), stream, group, id)
	_, err := cmd.Result()

	code := convertErrorToMetricsCode(err)
	redis.HandledCounter.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xAckCmdName,
		code,
	).Inc()
	redis.HandledHistogram.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xAckCmdName,
		code,
	).Observe(time.Since(startTime).Seconds())

	if err != nil {
		c.logger.Error("Failed to acknowledge message",
			zap.String("stream", stream),
			zap.String("group", group),
			zap.String("id", id),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// XPendingExt gets extended information about pending messages in a consumer group
func (c *client) XPendingExt(ctx context.Context,
	stream, group, start, end string,
	count int64, idle time.Duration) ([]goredis.XPendingExt, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, xPendingCmdName).Inc()

	// Create XPendingExtArgs
	args := &goredis.XPendingExtArgs{
		Stream: stream,
		Group:  group,
		Start:  start,
		End:    end,
		Count:  count,
		Idle:   idle,
	}

	// Execute XPendingExt command
	cmd := c.rc.XPendingExt(ctx, args)
	pending, err := cmd.Result()

	code := convertErrorToMetricsCode(err)
	redis.HandledCounter.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xPendingCmdName,
		code,
	).Inc()
	redis.HandledHistogram.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xPendingCmdName,
		code,
	).Observe(time.Since(startTime).Seconds())

	if err != nil && err != goredis.Nil {
		c.logger.Error("Failed to get pending messages",
			zap.String("stream", stream),
			zap.String("group", group),
			zap.Error(err),
		)
		return nil, err
	}

	return pending, err
}

// XClaim claims pending messages from a consumer group
func (c *client) XClaim(ctx context.Context,
	stream, group, consumer string,
	minIdle time.Duration,
	ids []string,
) ([]goredis.XMessage, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, xClaimCmdName).Inc()

	// Create XClaimArgs
	args := &goredis.XClaimArgs{
		Stream:   stream,
		Group:    group,
		Consumer: consumer,
		MinIdle:  minIdle,
		Messages: ids,
	}

	// Execute XClaim command
	cmd := c.rc.XClaim(ctx, args)
	messages, err := cmd.Result()

	code := convertErrorToMetricsCode(err)
	redis.HandledCounter.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xClaimCmdName,
		code,
	).Inc()
	redis.HandledHistogram.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xClaimCmdName,
		code,
	).Observe(time.Since(startTime).Seconds())

	if err != nil && err != goredis.Nil {
		c.logger.Error("Failed to claim messages",
			zap.String("stream", stream),
			zap.String("group", group),
			zap.String("consumer", consumer),
			zap.Strings("ids", ids),
			zap.Error(err),
		)
		return nil, err
	}

	return messages, err
}

// XInfoGroups gets information about consumer groups for a stream
func (c *client) XInfoGroups(ctx context.Context, stream string) ([]goredis.XInfoGroup, error) {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, xInfoGroupsCmdName).Inc()

	// Execute XInfoGroups command
	cmd := c.rc.XInfoGroups(ctx, stream)
	groups, err := cmd.Result()

	code := convertErrorToMetricsCode(err)
	redis.HandledCounter.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xInfoGroupsCmdName,
		code,
	).Inc()
	redis.HandledHistogram.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xInfoGroupsCmdName,
		code,
	).Observe(time.Since(startTime).Seconds())

	if err != nil && err != goredis.Nil {
		c.logger.Error("Failed to get stream groups info",
			zap.String("stream", stream),
			zap.Error(err),
		)
		return nil, err
	}

	return groups, err
}

// XGroupDestroy destroys a consumer group from a stream
func (c *client) XGroupDestroy(ctx context.Context, stream, group string) error {
	startTime := time.Now()
	redis.ReceivedCounter.WithLabelValues(clientVersion, c.opts.serverName, xGroupDestroyCmdName).Inc()

	err := c.rc.XGroupDestroy(ctx, stream, group).Err()

	code := convertErrorToMetricsCode(err)
	redis.HandledCounter.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xGroupDestroyCmdName,
		code,
	).Inc()
	redis.HandledHistogram.WithLabelValues(
		clientVersion,
		c.opts.serverName,
		xGroupDestroyCmdName,
		code,
	).Observe(time.Since(startTime).Seconds())

	if err != nil && err != goredis.Nil {
		// Missing stream/group is expected during startup cleanup or when the
		// topic partition never received messages; log at debug to avoid noisy
		// error logs during normal rollouts.
		if isBenignXGroupDestroyErr(err) {
			c.logger.Debug("Consumer group or stream not found during destroy (expected)",
				zap.String("stream", stream),
				zap.String("group", group),
				zap.Error(err),
			)
			return err
		}
		c.logger.Error("Failed to destroy consumer group",
			zap.String("stream", stream),
			zap.String("group", group),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// isBenignXGroupDestroyErr returns true for errors that indicate the stream or
// consumer group simply does not exist, which is expected during cleanup.
func isBenignXGroupDestroyErr(err error) bool {
	if err == nil {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "nogroup") ||
		strings.Contains(msg, "no such key") ||
		strings.Contains(msg, "stream key not found") ||
		strings.Contains(msg, "the xgroup subcommand requires the key to exist") ||
		strings.Contains(msg, "requires the key to exist")
}
