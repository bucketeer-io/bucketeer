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
	"time"

	"github.com/gomodule/redigo/redis"
)

type conn struct {
	redis.Conn
	clientVersion string
	serverName    string
}

func (c *conn) Do(commandName string, args ...interface{}) (interface{}, error) {
	startTime := time.Now()
	ReceivedCounter.WithLabelValues(c.clientVersion, c.serverName, commandName).Inc()
	reply, err := c.Conn.Do(commandName, args...)
	code := CodeFail
	switch err {
	case nil:
		code = CodeSuccess
	case ErrNil:
		code = CodeNotFound
	}
	HandledCounter.WithLabelValues(c.clientVersion, c.serverName, commandName, code).Inc()
	HandledHistogram.WithLabelValues(c.clientVersion, c.serverName, commandName, code).Observe(
		time.Since(startTime).Seconds())
	return reply, err
}
