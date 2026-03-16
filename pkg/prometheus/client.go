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
package prometheus

import (
	"context"
	"errors"
	"time"

	"github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"
)

// Client defines the interface for querying Prometheus.
type Client interface {
	QueryInstant(ctx context.Context, query string, ts time.Time) (model.Vector, error)
}

// Option configures a Client.
type Option func(*options)

// WithLogger sets the logger for the client.
func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type options struct {
	logger *zap.Logger
}

type client struct {
	api    promv1.API
	logger *zap.Logger
}

// NewClient creates a new Prometheus client.
func NewClient(baseURL string, opts ...Option) (Client, error) {
	o := &options{logger: zap.NewNop()}
	for _, opt := range opts {
		opt(o)
	}

	apiClient, err := api.NewClient(api.Config{Address: baseURL})
	if err != nil {
		return nil, err
	}

	return &client{
		api:    promv1.NewAPI(apiClient),
		logger: o.logger.Named("prometheus"),
	}, nil
}

func (c *client) QueryInstant(ctx context.Context, query string, ts time.Time) (model.Vector, error) {
	result, warnings, err := c.api.Query(ctx, query, ts)
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		c.logger.Warn("Prometheus warnings", zap.Strings("warnings", warnings))
	}
	vector, ok := result.(model.Vector)
	if !ok {
		return nil, errors.New("prometheus: unexpected response type, expected vector")
	}
	return vector, nil
}
