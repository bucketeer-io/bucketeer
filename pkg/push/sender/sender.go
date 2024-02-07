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

package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	pushclient "github.com/bucketeer-io/bucketeer/pkg/push/client"
	pushdomain "github.com/bucketeer-io/bucketeer/pkg/push/domain"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
)

type options struct {
	maxMPS     int
	numWorkers int
	metrics    metrics.Registerer
	logger     *zap.Logger
}

const (
	listRequestSize = 500
	fcmSendURL      = "https://fcm.googleapis.com/fcm/send"
	topicPrefix     = "bucketeer-"
)

var defaultOptions = options{
	maxMPS:     1000,
	numWorkers: 1,
	logger:     zap.NewNop(),
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

type Sender interface {
	Check(context.Context) health.Status
	Run() error
	Stop()
}

type sender struct {
	puller        puller.RateLimitedPuller
	pushClient    pushclient.Client
	featureClient featureclient.Client
	featuresCache cachev3.FeaturesCache
	group         errgroup.Group
	opts          *options
	logger        *zap.Logger
	ctx           context.Context
	cancel        func()
	doneCh        chan struct{}
}

func NewSender(
	p puller.Puller,
	pushClient pushclient.Client,
	featureClient featureclient.Client,
	v3Cache cache.MultiGetCache,
	opts ...Option) Sender {

	ctx, cancel := context.WithCancel(context.Background())
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics != nil {
		registerMetrics(options.metrics)
	}
	return &sender{
		puller:        puller.NewRateLimitedPuller(p, options.maxMPS),
		pushClient:    pushClient,
		featureClient: featureClient,
		featuresCache: cachev3.NewFeaturesCache(v3Cache),
		opts:          &options,
		logger:        options.logger.Named("sender"),
		ctx:           ctx,
		cancel:        cancel,
		doneCh:        make(chan struct{}),
	}
}

func (s *sender) Run() error {
	defer close(s.doneCh)
	s.group.Go(func() error {
		return s.puller.Run(s.ctx)
	})
	for i := 0; i < s.opts.numWorkers; i++ {
		s.group.Go(s.runWorker)
	}
	return s.group.Wait()
}

func (s *sender) Stop() {
	s.cancel()
	<-s.doneCh
}

func (s *sender) Check(ctx context.Context) health.Status {
	select {
	case <-s.ctx.Done():
		s.logger.Error("Unhealthy due to context Done is closed", zap.Error(s.ctx.Err()))
		return health.Unhealthy
	default:
		if s.group.FinishedCount() > 0 {
			s.logger.Error("Unhealthy", zap.Int32("FinishedCount", s.group.FinishedCount()))
			return health.Unhealthy
		}
		return health.Healthy
	}
}

func (s *sender) runWorker() error {
	record := func(code codes.Code, startTime time.Time) {
		handledCounter.WithLabelValues(code.String()).Inc()
		handledHistogram.WithLabelValues(code.String()).Observe(time.Since(startTime).Seconds())
	}
	for {
		select {
		case msg, ok := <-s.puller.MessageCh():
			if !ok {
				return nil
			}
			receivedCounter.Inc()
			startTime := time.Now()
			if id := msg.Attributes["id"]; id == "" {
				msg.Ack()
				record(codes.MissingID, startTime)
				continue
			}
			s.handle(msg)
			msg.Ack()
			record(codes.OK, startTime)
		case <-s.ctx.Done():
			return nil
		}
	}
}

func (s *sender) handle(msg *puller.Message) {
	event, err := s.unmarshalMessage(msg)
	if err != nil {
		msg.Ack()
		handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
		return
	}
	featureID, isTarget := s.extractFeatureID(event)
	if !isTarget {
		msg.Ack()
		handledCounter.WithLabelValues(codes.OK.String()).Inc()
		return
	}
	if featureID == "" {
		msg.Ack()
		handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
		s.logger.Warn("Message contains an empty FeatureID", zap.Any("event", event))
		return
	}
	if err := s.send(featureID, event.EnvironmentNamespace); err != nil {
		msg.Ack()
		handledCounter.WithLabelValues(codes.NonRepeatableError.String()).Inc()
		return
	}
	msg.Ack()
	handledCounter.WithLabelValues(codes.OK.String()).Inc()
}

func (s *sender) send(featureID, environmentNamespace string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		Id:                   featureID,
		EnvironmentNamespace: environmentNamespace,
	})
	if err != nil {
		return err
	}
	pushes, err := s.listPushes(ctx, environmentNamespace)
	if err != nil {
		return err
	}
	if len(pushes) == 0 {
		s.logger.Info("No pushes",
			zap.String("featureId", featureID),
			zap.String("environmentNamespace", environmentNamespace),
		)
		return nil
	}
	var lastErr error
	for _, p := range pushes {
		d := pushdomain.Push{Push: p}
		for _, t := range resp.Feature.Tags {
			if !d.ExistTag(t) {
				continue
			}
			if !s.isFeaturesCacheLatest(ctx, environmentNamespace, t, resp.Feature.Id, resp.Feature.Version) {
				if err = s.updateFeatures(ctx, environmentNamespace, t); err != nil {
					s.logger.Error("Failed to update features", zap.Error(err),
						zap.String("featureId", featureID),
						zap.String("tag", t),
						zap.String("pushId", d.Push.Id),
						zap.String("environmentNamespace", environmentNamespace),
					)
				}
			}
			topic := topicPrefix + t
			if err = s.pushFCM(ctx, d.FcmApiKey, topic); err != nil {
				s.logger.Error("Failed to push notification", zap.Error(err),
					zap.String("featureId", featureID),
					zap.String("tag", t),
					zap.String("topic", topic),
					zap.String("pushId", d.Push.Id),
					zap.String("environmentNamespace", environmentNamespace),
				)
				lastErr = err
				continue
			}
			s.logger.Info("Succeeded to push notification",
				zap.String("featureId", featureID),
				zap.String("tag", t),
				zap.String("topic", topic),
				zap.String("pushId", d.Push.Id),
				zap.String("environmentNamespace", environmentNamespace),
			)
		}
	}
	return lastErr
}

func (s *sender) pushFCM(ctx context.Context, fcmAPIKey, topic string) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"to": "/topics/" + topic,
		// The values in the data payload should be converted to string type.
		// https://firebase.google.com/docs/cloud-messaging/http-server-ref
		"data": map[string]string{
			"bucketeer_feature_flag_updated": "true",
		},
		"content_available": true,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fcmSendURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", fcmAPIKey))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *sender) listPushes(ctx context.Context, environmentNamespace string) ([]*pushproto.Push, error) {
	pushes := []*pushproto.Push{}
	cursor := ""
	for {
		resp, err := s.pushClient.ListPushes(ctx, &pushproto.ListPushesRequest{
			PageSize:             listRequestSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentNamespace,
		})
		if err != nil {
			return nil, err
		}
		pushes = append(pushes, resp.Pushes...)
		pushSize := len(resp.Pushes)
		if pushSize == 0 || pushSize < listRequestSize {
			return pushes, nil
		}
		cursor = resp.Cursor
	}
}

func (s *sender) unmarshalMessage(msg *puller.Message) (*domaineventproto.Event, error) {
	event := &domaineventproto.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		s.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
		return nil, err
	}
	return event, nil
}

func (s *sender) extractFeatureID(event *domaineventproto.Event) (string, bool) {
	if event.EntityType != domaineventproto.Event_FEATURE {
		return "", false
	}
	if event.Type != domaineventproto.Event_FEATURE_VERSION_INCREMENTED {
		return "", false
	}
	return event.EntityId, true
}

func (s *sender) isFeaturesCacheLatest(
	ctx context.Context,
	environmentNamespace,
	tag, featureID string,
	featureVersion int32,
) bool {
	features, err := s.featuresCache.Get(environmentNamespace)
	if err != nil {
		s.logger.Info(
			"Failed to get Features",
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("tag", tag),
			zap.String("featureId", featureID),
			zap.Int32("featureVersion", featureVersion),
		)
		return false
	}
	return s.isFeaturesLatest(features, featureID, featureVersion)
}

func (s *sender) updateFeatures(ctx context.Context, environmentNamespace, tag string) error {
	fs, err := s.listFeatures(ctx, environmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve features from storage",
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("tag", tag),
		)
		return err
	}
	features := &featureproto.Features{
		Features: fs,
	}
	if err := s.featuresCache.Put(features, environmentNamespace); err != nil {
		s.logger.Error(
			"Failed to cache features",
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
		)
		return err
	}
	return nil
}

func (s *sender) isFeaturesLatest(
	features *featureproto.Features,
	featureID string,
	featureVersion int32,
) bool {
	for _, f := range features.Features {
		if f.Id == featureID {
			return f.Version >= featureVersion
		}
	}
	return false
}

func (s *sender) listFeatures(ctx context.Context, environmentNamespace string) ([]*featureproto.Feature, error) {
	features := []*featureproto.Feature{}
	cursor := ""
	for {
		resp, err := s.featureClient.ListFeatures(ctx, &featureproto.ListFeaturesRequest{
			PageSize:             listRequestSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentNamespace,
		})
		if err != nil {
			return nil, err
		}
		for _, f := range resp.Features {
			ff := featuredomain.Feature{Feature: f}
			if ff.IsDisabledAndOffVariationEmpty() {
				continue
			}
			// To keep the cache size small, we exclude feature flags archived more than thirty days ago.
			if ff.IsArchivedBeforeLastThirtyDays() {
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
