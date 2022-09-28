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

package cacher

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	featureservice "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	domainevent "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	listRequestSize = 500
)

type options struct {
	maxMPS        int
	numWorkers    int
	flushSize     int
	flushInterval time.Duration
	metrics       metrics.Registerer
	logger        *zap.Logger
}

type Option func(*options)

func WithMaxMPS(mps int) Option {
	return func(opts *options) {
		opts.maxMPS = mps
	}
}

func WithNumWorkers(n int) Option {
	return func(opts *options) {
		opts.numWorkers = n
	}
}

func WithFlushSize(size int) Option {
	return func(opts *options) {
		opts.flushSize = size
	}
}

func WithFlushInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.flushInterval = interval
	}
}

func WithMetrics(r metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = r
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type FeatureCacher struct {
	puller        puller.RateLimitedPuller
	featuresCache cachev3.FeaturesCache
	featureClient featureservice.Client
	group         errgroup.Group
	opts          *options
	logger        *zap.Logger
	ctx           context.Context
	cancel        func()
	doneCh        chan struct{}
}

func NewFeatureCacher(
	p puller.Puller,
	client featureservice.Client,
	v3Cache cache.MultiGetCache,
	opts ...Option,
) *FeatureCacher {
	ctx, cancel := context.WithCancel(context.Background())
	dopts := &options{
		maxMPS:        1000,
		numWorkers:    1,
		flushSize:     100,
		flushInterval: time.Minute,
		logger:        zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	return &FeatureCacher{
		puller:        puller.NewRateLimitedPuller(p, dopts.maxMPS),
		featuresCache: cachev3.NewFeaturesCache(v3Cache),
		featureClient: client,
		opts:          dopts,
		logger:        dopts.logger.Named("cacher"),
		ctx:           ctx,
		cancel:        cancel,
		doneCh:        make(chan struct{}),
	}
}

func (c *FeatureCacher) Run() error {
	defer close(c.doneCh)
	c.group.Go(func() error {
		return c.puller.Run(c.ctx)
	})
	for i := 0; i < c.opts.numWorkers; i++ {
		c.group.Go(c.batch)
	}
	return c.group.Wait()
}

func (c *FeatureCacher) Stop() {
	c.cancel()
	<-c.doneCh
}

func (c *FeatureCacher) Check(ctx context.Context) health.Status {
	select {
	case <-c.ctx.Done():
		c.logger.Error("Unhealthy due to context Done is closed", zap.Error(c.ctx.Err()))
		return health.Unhealthy
	default:
		if c.group.FinishedCount() > 0 {
			c.logger.Error("Unhealthy", zap.Int32("FinishedCount", c.group.FinishedCount()))
			return health.Unhealthy
		}
		return health.Healthy
	}
}

func (c *FeatureCacher) batch() error {
	chunk := make(map[string]*puller.Message, c.opts.flushSize)
	timer := time.NewTimer(c.opts.flushInterval)
	defer timer.Stop()
	for {
		select {
		case msg, ok := <-c.puller.MessageCh():
			if !ok {
				return nil
			}
			receivedCounter.Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				handledCounter.WithLabelValues(codes.MissingID.String()).Inc()
				continue
			}
			if _, ok := chunk[id]; ok {
				c.logger.Warn("Message with duplicate id", zap.String("id", id))
				handledCounter.WithLabelValues(codes.DuplicateID.String()).Inc()
			}
			chunk[id] = msg
			if len(chunk) >= c.opts.flushSize {
				c.handleChunk(chunk)
				chunk = make(map[string]*puller.Message, c.opts.flushSize)
				timer.Reset(c.opts.flushInterval)
			}
		case <-timer.C:
			if len(chunk) > 0 {
				c.handleChunk(chunk)
				chunk = make(map[string]*puller.Message, c.opts.flushSize)
			}
			timer.Reset(c.opts.flushInterval)
		case <-c.ctx.Done():
			return nil
		}
	}
}

func (c *FeatureCacher) handleChunk(chunk map[string]*puller.Message) {
	handledFeatures := make(map[string]struct{}, len(chunk))
	for _, msg := range chunk {
		event, err := c.unmarshalMessage(msg)
		if err != nil {
			msg.Ack()
			handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
			continue
		}
		featureID, isTarget := c.extractFeatureID(event)
		if !isTarget {
			msg.Ack()
			handledCounter.WithLabelValues(codes.OK.String()).Inc()
			continue
		}
		if featureID == "" {
			msg.Ack()
			handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
			c.logger.Warn("Message contains an empty FeatureID", zap.Any("event", event))
			continue
		}
		if _, ok := handledFeatures[c.handledFeatureKey(featureID, event.EnvironmentNamespace)]; ok {
			msg.Ack()
			handledCounter.WithLabelValues(codes.OK.String()).Inc()
			continue
		}
		if ok := c.refresh(event.EnvironmentNamespace); ok {
			msg.Ack()
			handledFeatures[c.handledFeatureKey(featureID, event.EnvironmentNamespace)] = struct{}{}
			handledCounter.WithLabelValues(codes.OK.String()).Inc()
		} else {
			msg.Nack()
			handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
		}
	}
}

func (c *FeatureCacher) handledFeatureKey(featureID, environmentNamespace string) string {
	if environmentNamespace == "" {
		return featureID
	}
	return fmt.Sprintf("%s:%s", environmentNamespace, featureID)
}

func (c *FeatureCacher) refresh(environmentNamespace string) bool {
	features, err := c.listFeatures(environmentNamespace)
	if err != nil {
		c.logger.Error("Failed to retrieve features", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace))
		return false
	}
	err = c.featuresCache.Put(&featureproto.Features{
		Features: features,
	}, environmentNamespace)
	if err != nil {
		c.logger.Error(
			"Failed to cache Features",
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
		)
		return false
	}
	return true
}

func (c *FeatureCacher) listFeatures(environmentNamespace string) ([]*featureproto.Feature, error) {
	features := []*featureproto.Feature{}
	cursor := ""
	for {
		resp, err := c.featureClient.ListFeatures(c.ctx, &featureproto.ListFeaturesRequest{
			PageSize:             listRequestSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentNamespace,
			Archived:             &wrappers.BoolValue{Value: false},
		})
		if err != nil {
			return nil, err
		}
		for _, f := range resp.Features {
			if !f.Enabled && f.OffVariation == "" {
				continue
			}
			features = append(features, f)
		}
		featureSize := len(resp.Features)
		if featureSize == 0 || featureSize < listRequestSize {
			return features, nil
		}
		cursor = resp.Cursor
	}
}

func (c *FeatureCacher) unmarshalMessage(msg *puller.Message) (*domainevent.Event, error) {
	event := &domainevent.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		c.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
		return nil, err
	}
	return event, nil
}

func (c *FeatureCacher) extractFeatureID(event *domainevent.Event) (string, bool) {
	if event.EntityType != domainevent.Event_FEATURE {
		return "", false
	}
	return event.EntityId, true
}
