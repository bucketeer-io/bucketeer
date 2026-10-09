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
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	redis "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
)

// putIfNewerScript is the Redis-side mirror of cache.IsStaleGeneration. It
// must stay in sync with pkg/cache/generation.go.
//
// KEYS[1] = value key
// KEYS[2] = generation marker key (must hash to the same slot as KEYS[1])
// ARGV[1] = value
// ARGV[2] = new generation, fixed-width decimal (cache.FormatGeneration)
// ARGV[3] = TTL in milliseconds; 0 or less means no expiry
// ARGV[4] = marker to store on accept (cache.EncodeGenerationMarker), computed
//
//	in Go so Redis never hashes the new value
//
// Returns 1 when the value was written, 0 when it was rejected as stale.
//
// The write is rejected only when the value key exists, the marker is
// well-formed, the marker's generation is strictly greater than ARGV[2], and
// the marker describes the stored value (recorded length equals STRLEN and
// recorded digest equals sha1hex of the value). The digest is only computed
// on that last step, so the common accept path never hashes the stored value.
// Generations are compared as equal-length strings, which is exact;
// converting them with tonumber would lose precision above 2^53.
var putIfNewerScript = strings.ReplaceAll(
	putIfNewerScriptTemplate,
	"__GENERATION_DIGITS__",
	strconv.Itoa(cache.GenerationDigits),
)

const putIfNewerScriptTemplate = `
local new_gen = ARGV[2]
if redis.call('EXISTS', KEYS[1]) == 1 then
  local marker = redis.call('GET', KEYS[2])
  if marker then
    local pattern = '^(' .. string.rep('%d', __GENERATION_DIGITS__) .. '):(%d+):(%x+)$'
    local cur_gen, cur_len, cur_digest = string.match(marker, pattern)
    if cur_gen and new_gen < cur_gen and tonumber(cur_len) == redis.call('STRLEN', KEYS[1]) then
      if cur_digest == redis.sha1hex(redis.call('GET', KEYS[1])) then
        return 0
      end
    end
  end
end
local marker = ARGV[4]
local ttl = tonumber(ARGV[3])
if ttl and ttl > 0 then
  redis.call('SET', KEYS[1], ARGV[1], 'PX', ttl)
  redis.call('SET', KEYS[2], marker, 'PX', ttl)
else
  redis.call('SET', KEYS[1], ARGV[1])
  redis.call('SET', KEYS[2], marker)
end
return 1
`

type redisCache struct {
	client redis.Client
}

func NewRedisCache(client redis.Client) cache.MultiGetDeleteCountCache {
	return &redisCache{
		client: client,
	}
}

func (r *redisCache) Get(key interface{}) (interface{}, error) {
	value, err := r.client.Get(key.(string))
	if err != nil {
		if err == redis.ErrNil {
			return nil, cache.ErrNotFound
		}
		return nil, err
	}
	return value, nil
}

func (r *redisCache) Put(key interface{}, value interface{}, expiration time.Duration) error {
	return r.client.Set(key.(string), value, expiration)
}

// PutIfNewer implements cache.ConditionalPutter atomically via a Lua script.
func (r *redisCache) PutIfNewer(
	key, genKey string,
	value []byte,
	generation int64,
	expiration time.Duration,
) (bool, error) {
	formatted, err := cache.FormatGeneration(generation)
	if err != nil {
		return false, err
	}
	cmd := r.client.Eval(
		context.Background(),
		putIfNewerScript,
		[]string{key, genKey},
		value,
		formatted,
		expiration.Milliseconds(),
		cache.EncodeGenerationMarker(formatted, value),
	)
	res, err := cmd.Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func (r *redisCache) SAdd(key string, members ...interface{}) (int64, error) {
	return r.client.SAdd(key, members...)
}

func (r *redisCache) SMembers(key string) ([]string, error) {
	return r.client.SMembers(key)
}

func (r *redisCache) GetMulti(keys interface{}, ignoreNotFound bool) ([]interface{}, error) {
	value, err := r.client.GetMulti(keys.([]string), ignoreNotFound)
	switch err {
	case nil:
		return value, nil
	case redis.ErrNil:
		return nil, cache.ErrNotFound
	case redis.ErrInvalidType:
		return nil, cache.ErrInvalidType
	default:
		return nil, err
	}
}

func (r *redisCache) Scan(cursor, key, count interface{}) (uint64, []string, error) {
	c, keys, err := r.client.Scan(cursor.(uint64), key.(string), count.(int64))
	switch err {
	case nil:
		return c, keys, nil
	case redis.ErrNil:
		return 0, nil, cache.ErrNotFound
	default:
		return 0, nil, err
	}
}

func (r *redisCache) Delete(key string) error {
	return r.client.Del(key)
}

func (r *redisCache) Increment(key string) (int64, error) {
	return r.client.Incr(key)
}

func (r *redisCache) IncrementBy(key string, value int64) (int64, error) {
	return r.client.IncrBy(key, value)
}

func (r *redisCache) PFCount(keys ...string) (int64, error) {
	return r.client.PFCount(keys...)
}

func (r *redisCache) PFMerge(dest string, expiration time.Duration, keys ...string) error {
	return r.client.PFMerge(dest, expiration, keys...)
}

func (r *redisCache) PFAdd(key string, els ...string) (int64, error) {
	return r.client.PFAdd(key, els...)
}

func (r *redisCache) Pipeline(tx bool) redis.PipeClient {
	return r.client.Pipeline(tx)
}

func (r *redisCache) Expire(key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(key, expiration)
}
