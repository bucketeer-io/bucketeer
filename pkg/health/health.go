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

package health

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type Status uint32

const (
	// Healthy is returned when the health check was successful
	Healthy Status = 0

	// Unhealthy is returned when the health check was unsuccessful
	Unhealthy Status = 1
)

func (s Status) String() string {
	switch s {
	case Healthy:
		return "Healthy"
	case Unhealthy:
		return "Unhealthy"
	default:
		return "Unknown"
	}
}

type check func(context.Context) Status

type checker struct {
	status uint32

	interval time.Duration
	timeout  time.Duration
	checks   map[string]check
}

type option func(*checker)

func WithCheck(name string, check check) option {
	return func(c *checker) {
		if _, ok := c.checks[name]; ok {
			panic(fmt.Sprintf("health: %s already registered", name))
		}
		c.checks[name] = check
	}
}

func WithInterval(interval time.Duration) option {
	return func(c *checker) {
		c.interval = interval
	}
}

func WithTimeout(timeout time.Duration) option {
	return func(c *checker) {
		c.timeout = timeout
	}
}

func newChecker(opts ...option) *checker {
	checker := &checker{
		status:   uint32(Unhealthy),
		interval: 10 * time.Second,
		timeout:  5 * time.Second,
		checks:   make(map[string]check),
	}
	for _, o := range opts {
		o(checker)
	}
	return checker
}

func (hc *checker) Run(ctx context.Context) {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()
	hc.check(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			hc.check(ctx)
		}
	}
}

func (hc *checker) check(ctx context.Context) {
	resultChan := make(chan Status, len(hc.checks))
	ctx, cancel := context.WithTimeout(ctx, hc.timeout)
	defer cancel()
	for _, c := range hc.checks {
		go func(c check) {
			resultChan <- c(ctx)
		}(c)
	}
	for i := 0; i < len(hc.checks); i++ {
		if res := <-resultChan; res != Healthy {
			hc.setStatus(Unhealthy)
			return
		}
	}
	hc.setStatus(Healthy)
}

func (hc *checker) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if hc.getStatus() == Unhealthy {
		resp.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (hc *checker) getStatus() Status {
	return Status(atomic.LoadUint32(&hc.status))
}

func (hc *checker) setStatus(s Status) {
	atomic.StoreUint32(&hc.status, uint32(s))
}

func (hc *checker) Stop() {
	hc.setStatus(Unhealthy)
}
