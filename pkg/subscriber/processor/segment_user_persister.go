// Copyright 2025 The Bucketeer Authors.
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

package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	btclient "github.com/bucketeer-io/bucketeer/pkg/batch/client"
	"github.com/bucketeer-io/bucketeer/pkg/feature/command"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/factory"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	v3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
	btproto "github.com/bucketeer-io/bucketeer/proto/batch"
	domainproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	serviceevent "github.com/bucketeer-io/bucketeer/proto/event/service"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	maxUserIDLength = 100
)

type segmentUserPersisterConfig struct {
	DomainEventProject string `json:"domainEventProject"`
	DomainEventTopic   string `json:"domainEventTopic"`
	FlushSize          int    `json:"flushSize"`
	FlushInterval      int    `json:"flushInterval"`
	// PubSub configuration
	PubSubType          string `json:"pubSubType"`          // google or redis-stream
	RedisServerName     string `json:"redisServerName"`     // Redis server name
	RedisAddr           string `json:"redisAddr"`           // Redis address
	RedisPoolSize       int    `json:"redisPoolSize"`       // Redis pool size
	RedisMinIdle        int    `json:"redisMinIdle"`        // Redis min idle connections
	RedisPartitionCount int    `json:"redisPartitionCount"` // Redis partition count
	Project             string `json:"project"`             // Google Cloud project ID
}

type segmentUserPersister struct {
	segmentUserPersisterConfig segmentUserPersisterConfig
	domainPublisher            publisher.Publisher
	batchClient                btclient.Client
	mysqlClient                mysql.Client
	segmentStorage             v2fs.SegmentStorage
	segmentUserStorage         v2fs.SegmentUserStorage
	logger                     *zap.Logger
}

func NewSegmentUserPersister(
	config interface{},
	batchClient btclient.Client,
	mysqlClient mysql.Client,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (subscriber.PubSubProcessor, error) {
	segmentPersisterJsonConfig, ok := config.(map[string]interface{})
	if !ok {
		logger.Error("SegmentUserPersister: invalid config")
		return nil, ErrSegmentInvalidConfig
	}
	configBytes, err := json.Marshal(segmentPersisterJsonConfig)
	if err != nil {
		logger.Error("SegmentUserPersister: failed to marshal config", zap.Error(err))
		return nil, err
	}
	var segmentPersisterConfig segmentUserPersisterConfig
	err = json.Unmarshal(configBytes, &segmentPersisterConfig)
	if err != nil {
		logger.Error("SegmentUserPersister: failed to unmarshal config", zap.Error(err))
		return nil, err
	}

	// Create domain publisher using factory pattern
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Determine PubSub type and create appropriate client
	var pubSubType factory.PubSubType
	switch segmentPersisterConfig.PubSubType {
	case "redis-stream":
		pubSubType = factory.RedisStream
	case "google":
		pubSubType = factory.Google
	default:
		// Default to Google for backward compatibility
		pubSubType = factory.Google
	}

	// Create factory client options
	var opts []factory.Option
	opts = append(opts, factory.WithPubSubType(pubSubType))
	opts = append(opts, factory.WithLogger(logger))

	if pubSubType == factory.Google {
		opts = append(opts, factory.WithProjectID(segmentPersisterConfig.Project))
	} else if pubSubType == factory.RedisStream {
		// Create Redis client internally
		redisClient, err := createRedisClientForSegmentPersister(ctx, segmentPersisterConfig, logger, registerer)
		if err != nil {
			logger.Error("SegmentUserPersister: failed to create Redis client", zap.Error(err))
			return nil, err
		}
		opts = append(opts, factory.WithRedisClient(redisClient))
		if segmentPersisterConfig.RedisPartitionCount > 0 {
			opts = append(opts, factory.WithPartitionCount(segmentPersisterConfig.RedisPartitionCount))
		}
	}

	client, err := factory.NewClient(ctx, opts...)
	if err != nil {
		logger.Error("SegmentUserPersister: failed to create pubsub client", zap.Error(err))
		return nil, err
	}

	domainPublisher, err := client.CreatePublisher(segmentPersisterConfig.DomainEventTopic)
	if err != nil {
		logger.Error("SegmentUserPersister: failed to create domain publisher", zap.Error(err))
		return nil, err
	}
	return &segmentUserPersister{
		segmentUserPersisterConfig: segmentPersisterConfig,
		domainPublisher:            domainPublisher,
		batchClient:                batchClient,
		mysqlClient:                mysqlClient,
		segmentStorage:             v2fs.NewSegmentStorage(mysqlClient),
		segmentUserStorage:         v2fs.NewSegmentUserStorage(mysqlClient),
		logger:                     logger,
	}, nil
}

// createRedisClientForSegmentPersister creates a Redis client from the configuration
func createRedisClientForSegmentPersister(
	ctx context.Context,
	conf segmentUserPersisterConfig,
	logger *zap.Logger,
	registerer metrics.Registerer,
) (v3.Client, error) {
	redisAddr := conf.RedisAddr
	if redisAddr == "" {
		return nil, fmt.Errorf("redis address is required for Redis PubSub")
	}

	redisPoolSize := 10
	if conf.RedisPoolSize > 0 {
		redisPoolSize = conf.RedisPoolSize
	}

	redisMinIdle := 3
	if conf.RedisMinIdle > 0 {
		redisMinIdle = conf.RedisMinIdle
	}

	logger.Debug("Creating Redis client for segment user persister",
		zap.String("address", redisAddr),
		zap.Int("poolSize", redisPoolSize),
		zap.Int("minIdle", redisMinIdle),
		zap.String("serverName", conf.RedisServerName),
		zap.String("pubSubType", conf.PubSubType),
	)

	// Create Redis client
	return v3.NewClient(
		redisAddr,
		v3.WithPoolSize(redisPoolSize),
		v3.WithMinIdleConns(redisMinIdle),
		v3.WithServerName(conf.RedisServerName),
		v3.WithMetrics(registerer),
		v3.WithLogger(logger),
	)
}

func (p *segmentUserPersister) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	chunk := make(map[string]*puller.Message, p.segmentUserPersisterConfig.FlushSize)
	ticker := time.NewTicker(time.Duration(p.segmentUserPersisterConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberSegmentUser).Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				subscriberHandledCounter.WithLabelValues(subscriberSegmentUser, codes.MissingID.String()).Inc()
				continue
			}
			if _, ok := chunk[id]; ok {
				subscriberHandledCounter.WithLabelValues(subscriberSegmentUser, codes.DuplicateID.String()).Inc()
			}
			chunk[id] = msg
			if len(chunk) >= p.segmentUserPersisterConfig.FlushSize {
				p.handleChunk(ctx, chunk)
				chunk = make(map[string]*puller.Message, p.segmentUserPersisterConfig.FlushSize)
			}
		case <-ticker.C:
			if len(chunk) > 0 {
				p.handleChunk(ctx, chunk)
				chunk = make(map[string]*puller.Message, p.segmentUserPersisterConfig.FlushSize)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (p *segmentUserPersister) handleChunk(ctx context.Context, chunk map[string]*puller.Message) {
	for _, msg := range chunk {
		p.logger.Debug("handling a message", zap.String("msgID", msg.ID))
		event, err := p.unmarshalMessage(msg)
		if err != nil {
			msg.Ack()
			p.logger.Error("failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
			subscriberHandledCounter.WithLabelValues(subscriberSegmentUser, codes.BadMessage.String()).Inc()
			continue
		}
		if !validateSegmentUserState(event.State) {
			msg.Ack()
			p.logger.Error(
				"invalid state",
				zap.String("environmentId", event.EnvironmentId),
				zap.Int32("state", int32(event.State)),
			)
			subscriberHandledCounter.WithLabelValues(subscriberSegmentUser, codes.BadMessage.String()).Inc()
			if err := p.updateSegmentStatus(
				ctx,
				event.Editor,
				event.EnvironmentId,
				event.SegmentId,
				0,
				event.State,
				featureproto.Segment_FAILED,
			); err != nil {
				p.logger.Error(
					"failed to update segment status",
					zap.Error(err),
					zap.String("environmentId", event.EnvironmentId),
				)
			}
			continue
		}
		if err := p.handleEvent(ctx, event); err != nil {
			switch {
			case errors.Is(err, storage.ErrKeyNotFound), errors.Is(err, v2fs.ErrSegmentNotFound):
				msg.Ack()
				p.logger.Warn(
					"segment not found",
					zap.Error(err),
					zap.String("environmentId", event.EnvironmentId),
				)
				subscriberHandledCounter.WithLabelValues(subscriberSegmentUser, codes.NonRepeatableError.String()).Inc()
			case errors.Is(err, ErrSegmentExceededMaxUserIDLength):
				msg.Ack()
				p.logger.Warn(
					"exceeded max user id length",
					zap.Error(err),
					zap.String("environmentId", event.EnvironmentId),
				)
				subscriberHandledCounter.WithLabelValues(subscriberSegmentUser, codes.NonRepeatableError.String()).Inc()
				if err := p.updateSegmentStatus(
					ctx,
					event.Editor,
					event.EnvironmentId,
					event.SegmentId,
					0,
					event.State,
					featureproto.Segment_FAILED,
				); err != nil {
					p.logger.Error(
						"failed to update segment status",
						zap.Error(err),
						zap.String("environmentId", event.EnvironmentId),
					)
				}
			default:
				// retryable
				msg.Nack()
				p.logger.Error(
					"failed to handle event",
					zap.Error(err),
					zap.String("environmentId", event.EnvironmentId),
				)
				subscriberHandledCounter.WithLabelValues(subscriberSegmentUser, codes.RepeatableError.String()).Inc()
			}
			continue
		}
		msg.Ack()
		p.logger.Debug(
			"suceeded to persist segment users",
			zap.String("msgID", msg.ID),
			zap.String("environmentId", event.EnvironmentId),
			zap.String("segmentId", event.SegmentId),
		)
		subscriberHandledCounter.WithLabelValues(subscriberSegmentUser, codes.OK.String()).Inc()
	}
}

func (p *segmentUserPersister) unmarshalMessage(msg *puller.Message,
) (*serviceevent.BulkSegmentUsersReceivedEvent, error) {
	event := &serviceevent.BulkSegmentUsersReceivedEvent{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func validateSegmentUserState(state featureproto.SegmentUser_State) bool {
	switch state {
	case featureproto.SegmentUser_INCLUDED:
		return true
	default:
		return false
	}
}

func (p *segmentUserPersister) handleEvent(
	ctx context.Context, event *serviceevent.BulkSegmentUsersReceivedEvent) error {
	cnt, err := p.persistSegmentUsers(ctx, event.EnvironmentId, event.SegmentId, event.Data, event.State)
	if err != nil {
		if err := p.updateSegmentStatus(
			ctx,
			event.Editor,
			event.EnvironmentId,
			event.SegmentId,
			cnt,
			event.State,
			featureproto.Segment_FAILED,
		); err != nil {
			p.logger.Error(
				"failed to update to segment status to failed",
				zap.Error(err),
				zap.String("segmentId", event.SegmentId),
				zap.Int64("userCount", cnt),
				zap.String("environmentId", event.EnvironmentId),
			)
			return err
		}
		return err
	}
	return p.updateSegmentStatus(
		ctx,
		event.Editor,
		event.EnvironmentId,
		event.SegmentId,
		cnt,
		event.State,
		featureproto.Segment_SUCEEDED,
	)
}

func (p *segmentUserPersister) persistSegmentUsers(
	ctx context.Context,
	environmentId string,
	segmentID string,
	data []byte,
	state featureproto.SegmentUser_State,
) (int64, error) {
	segmentUserIDs := strings.Split(
		strings.NewReplacer(
			",", "\n",
			"\r\n", "\n",
		).Replace(string(data)),
		"\n",
	)
	uniqueSegmentUserIDs := make(map[string]struct{}, len(segmentUserIDs))
	for _, id := range segmentUserIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if len(id) > maxUserIDLength {
			return 0, ErrSegmentExceededMaxUserIDLength
		}
		uniqueSegmentUserIDs[id] = struct{}{}
	}
	allSegmentUsers := make([]*featureproto.SegmentUser, 0, len(uniqueSegmentUserIDs))
	var cnt int64
	for id := range uniqueSegmentUserIDs {
		cnt++
		user := domain.NewSegmentUser(segmentID, id, state, false)
		allSegmentUsers = append(allSegmentUsers, user.SegmentUser)
	}
	err := p.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		if err := p.segmentUserStorage.UpsertSegmentUsers(contextWithTx, allSegmentUsers, environmentId); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	p.updateSegmentUserCache(ctx)
	return cnt, nil
}

func (p *segmentUserPersister) updateSegmentStatus(
	ctx context.Context,
	editor *domainproto.Editor,
	environmentId string,
	segmentID string,
	cnt int64,
	state featureproto.SegmentUser_State,
	status featureproto.Segment_Status,
) error {
	return p.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		segment, _, err := p.segmentStorage.GetSegment(contextWithTx, segmentID, environmentId)
		if err != nil {
			return err
		}
		changeCmd := &featureproto.ChangeBulkUploadSegmentUsersStatusCommand{
			Status: status,
			State:  state,
			Count:  cnt,
		}
		handler, err := command.NewSegmentCommandHandler(editor, segment, p.domainPublisher, environmentId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, changeCmd); err != nil {
			return err
		}
		return p.segmentStorage.UpdateSegment(contextWithTx, segment, environmentId)
	})
}

// Even if the update request fails, the cronjob will keep trying
// to update the cache every minute, so we don't need to retry.
func (p *segmentUserPersister) updateSegmentUserCache(ctx context.Context) {
	req := &btproto.BatchJobRequest{
		Job: btproto.BatchJob_SegmentUserCacher,
	}
	_, err := p.batchClient.ExecuteBatchJob(ctx, req)
	if err != nil {
		p.logger.Error("Failed to update segment user cache", zap.Error(err))
	}
}
