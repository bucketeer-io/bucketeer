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

package apikeycacher

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	acproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	domainevent "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

const (
	listRequestPageSize = 500
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

type EnvAPIKeyCacher struct {
	puller            puller.RateLimitedPuller
	accountClient     accountclient.Client
	environmentClient environmentclient.Client
	envAPIKeyCache    cachev3.EnvironmentAPIKeyCache
	group             errgroup.Group
	opts              *options
	logger            *zap.Logger
	ctx               context.Context
	cancel            func()
	doneCh            chan struct{}
}

func NewEnvironmentAPIKeyCacher(
	p puller.Puller,
	accountClient accountclient.Client,
	environmentClient environmentclient.Client,
	v3Cache cache.Cache,
	opts ...Option,
) *EnvAPIKeyCacher {
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
	return &EnvAPIKeyCacher{
		puller:            puller.NewRateLimitedPuller(p, dopts.maxMPS),
		accountClient:     accountClient,
		environmentClient: environmentClient,
		envAPIKeyCache:    cachev3.NewEnvironmentAPIKeyCache(v3Cache),
		opts:              dopts,
		logger:            dopts.logger.Named("apikeycacher"),
		ctx:               ctx,
		cancel:            cancel,
		doneCh:            make(chan struct{}),
	}
}

func (c *EnvAPIKeyCacher) Run() error {
	defer close(c.doneCh)
	c.group.Go(func() error {
		return c.puller.Run(c.ctx)
	})
	for i := 0; i < c.opts.numWorkers; i++ {
		c.group.Go(c.batch)
	}
	return c.group.Wait()
}

func (c *EnvAPIKeyCacher) Stop() {
	c.cancel()
	<-c.doneCh
}

func (c *EnvAPIKeyCacher) Check(ctx context.Context) health.Status {
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

func (c *EnvAPIKeyCacher) batch() error {
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

func (c *EnvAPIKeyCacher) handleChunk(chunk map[string]*puller.Message) {
	for _, msg := range chunk {
		event, err := c.unmarshalMessage(msg)
		if err != nil {
			msg.Ack()
			handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
			c.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
			continue
		}
		switch event.EntityType {
		case domainevent.Event_APIKEY:
			c.handleAPIKeyEvent(msg, event)
		case domainevent.Event_PROJECT:
			c.handleProjectEvent(msg, event)
		case domainevent.Event_ENVIRONMENT:
			c.handleEnvironmentEvent(msg, event)
		default:
			msg.Ack()
			handledCounter.WithLabelValues(codes.OK.String()).Inc()
			continue
		}
	}
}

func (c *EnvAPIKeyCacher) handleAPIKeyEvent(msg *puller.Message, event *domainevent.Event) {
	apiKeyID := event.EntityId
	if apiKeyID == "" {
		msg.Ack()
		handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
		c.logger.Warn("Message contains an empty apiKeyID", zap.Any("event", event))
		return
	}
	envResp, err := c.environmentClient.GetEnvironmentV2(
		c.ctx,
		&environmentproto.GetEnvironmentV2Request{
			// EnvironmentNamespace in the domain event is same as the ID of the environment v2.
			Id: event.EnvironmentNamespace,
		},
	)
	if err != nil {
		if code := status.Code(err); code == grpccodes.NotFound {
			msg.Ack()
			handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
			c.logger.Error("The specified environmentNamespace does not exist",
				zap.Error(err),
				zap.String("apiKeyID", apiKeyID),
				zap.String("environmentNamespace", event.EnvironmentNamespace),
				zap.Any("event", event),
			)
			return
		}
		msg.Nack()
		handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
		c.logger.Error("Failed to get environment",
			zap.Error(err),
			zap.String("apiKeyID", apiKeyID),
			zap.String("environmentNamespace", event.EnvironmentNamespace),
			zap.Any("event", event),
		)
		return
	}
	if err := c.refresh(apiKeyID, false, envResp.Environment.ProjectId, envResp.Environment); err != nil {
		msg.Nack()
		handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
		c.logger.Error("Failed to refresh api key",
			zap.Error(err),
			zap.String("apiKeyID", apiKeyID),
			zap.String("environmentNamespace", event.EnvironmentNamespace),
			zap.Any("event", event),
		)
		return
	}
	msg.Ack()
	handledCounter.WithLabelValues(codes.OK.String()).Inc()
}

func (c *EnvAPIKeyCacher) handleProjectEvent(msg *puller.Message, event *domainevent.Event) {
	if !(event.Type == domainevent.Event_PROJECT_ENABLED || event.Type == domainevent.Event_PROJECT_DISABLED) {
		msg.Ack()
		handledCounter.WithLabelValues(codes.OK.String()).Inc()
		return
	}
	environmentDisabled := event.Type == domainevent.Event_PROJECT_DISABLED
	projectID := event.EntityId
	if projectID == "" {
		msg.Ack()
		handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
		c.logger.Warn("Message contains an empty projectID", zap.Any("event", event))
		return
	}
	environments, err := c.listEnvironments(projectID)
	if err != nil {
		msg.Nack()
		handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
		c.logger.Warn("Failed to list environments", zap.Any("event", event))
		return
	}
	for _, environment := range environments {
		if err := c.refreshAll(environmentDisabled, projectID, environment); err != nil {
			msg.Nack()
			handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
			c.logger.Error("Failed to refresh all api keys in the environment",
				zap.Error(err),
				zap.String("environmentNamespace", environment.Id),
				zap.Any("event", event),
			)
			return
		}
	}
	msg.Ack()
	handledCounter.WithLabelValues(codes.OK.String()).Inc()
}

func (c *EnvAPIKeyCacher) handleEnvironmentEvent(msg *puller.Message, event *domainevent.Event) {
	if event.Type != domainevent.Event_ENVIRONMENT_DELETED {
		msg.Ack()
		handledCounter.WithLabelValues(codes.OK.String()).Inc()
		return
	}
	ede := &domainevent.EnvironmentDeletedEvent{}
	if err := ptypes.UnmarshalAny(event.Data, ede); err != nil {
		msg.Ack()
		handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
		c.logger.Warn("Message doesn't contain an environment deleted event", zap.Any("event", event))
		return
	}
	environmentID := ede.Namespace
	envResp, err := c.environmentClient.GetEnvironmentV2(
		c.ctx,
		&environmentproto.GetEnvironmentV2Request{
			Id: environmentID,
		},
	)
	if err != nil {
		if code := status.Code(err); code == grpccodes.NotFound {
			msg.Ack()
			handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
			c.logger.Error("The specified environmentNamespace does not exist",
				zap.Error(err),
				zap.String("environmentNamespace", event.EnvironmentNamespace),
				zap.Any("event", event),
			)
			return
		}
		msg.Nack()
		handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
		c.logger.Error("Failed to get environment",
			zap.Error(err),
			zap.String("environmentNamespace", environmentID),
			zap.Any("event", event),
		)
		return
	}
	if err := c.refreshAll(true, envResp.Environment.ProjectId, envResp.Environment); err != nil {
		msg.Nack()
		handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
		c.logger.Error("Failed to refresh all api keys in the environment",
			zap.Error(err),
			zap.String("environmentNamespace", environmentID),
			zap.Any("event", event),
		)
		return
	}
	msg.Ack()
	handledCounter.WithLabelValues(codes.OK.String()).Inc()
}

func (c *EnvAPIKeyCacher) unmarshalMessage(msg *puller.Message) (*domainevent.Event, error) {
	event := &domainevent.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (c *EnvAPIKeyCacher) refresh(
	apiKeyID string,
	environmentDisabled bool,
	projectID string,
	environment *environmentproto.EnvironmentV2,
) error {
	req := &acproto.GetAPIKeyRequest{
		Id:                   apiKeyID,
		EnvironmentNamespace: environment.Id,
	}
	resp, err := c.accountClient.GetAPIKey(c.ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get api key: %w", err)
	}
	envAPIKey := &acproto.EnvironmentAPIKey{
		ApiKey:              resp.ApiKey,
		EnvironmentDisabled: environmentDisabled,
		ProjectId:           projectID,
		Environment:         environment,
	}
	return c.upsert(envAPIKey)
}

func (c *EnvAPIKeyCacher) refreshAll(
	environmentDisabled bool,
	projectID string,
	environment *environmentproto.EnvironmentV2,
) error {
	apiKeys, err := c.listAPIKeys(environment.Id)
	if err != nil {
		return fmt.Errorf("failed to list api keys: %w", err)
	}
	for _, key := range apiKeys {
		envAPIKey := &acproto.EnvironmentAPIKey{
			ApiKey:              key,
			EnvironmentDisabled: environmentDisabled,
			ProjectId:           projectID,
			Environment:         environment,
		}
		if err := c.upsert(envAPIKey); err != nil {
			return err
		}
	}
	return nil
}

func (c *EnvAPIKeyCacher) upsert(envAPIKey *acproto.EnvironmentAPIKey) error {
	if err := c.envAPIKeyCache.Put(envAPIKey); err != nil {
		return fmt.Errorf("failed to cache environment api key: %w", err)
	}
	c.logger.Info(
		"API key upserted successfully",
		zap.String("apiKeyID", envAPIKey.ApiKey.Id),
		zap.String("environmentID", envAPIKey.Environment.Id),
	)
	return nil
}

func (c *EnvAPIKeyCacher) listEnvironments(projectID string) ([]*environmentproto.EnvironmentV2, error) {
	var environments []*environmentproto.EnvironmentV2
	cursor := ""
	for {
		resp, err := c.environmentClient.ListEnvironmentsV2(c.ctx, &environmentproto.ListEnvironmentsV2Request{
			PageSize:  listRequestPageSize,
			Cursor:    cursor,
			ProjectId: projectID,
			Archived:  wrapperspb.Bool(false),
		})
		if err != nil {
			return nil, err
		}
		environments = append(environments, resp.Environments...)
		environmentSize := len(resp.Environments)
		if environmentSize == 0 || environmentSize < listRequestPageSize {
			return environments, nil
		}
		cursor = resp.Cursor
	}
}

func (c *EnvAPIKeyCacher) listAPIKeys(environmentNamespace string) ([]*acproto.APIKey, error) {
	apiKeys := []*acproto.APIKey{}
	cursor := ""
	for {
		resp, err := c.accountClient.ListAPIKeys(c.ctx, &acproto.ListAPIKeysRequest{
			PageSize:             listRequestPageSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentNamespace,
		})
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, resp.ApiKeys...)
		apiKeySize := len(resp.ApiKeys)
		if apiKeySize == 0 || apiKeySize < listRequestPageSize {
			return apiKeys, nil
		}
		cursor = resp.Cursor
	}
}
