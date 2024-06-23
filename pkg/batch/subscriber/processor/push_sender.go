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

package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	btclient "github.com/bucketeer-io/bucketeer/pkg/batch/client"
	"github.com/bucketeer-io/bucketeer/pkg/batch/subscriber"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	pushclient "github.com/bucketeer-io/bucketeer/pkg/push/client"
	pushdomain "github.com/bucketeer-io/bucketeer/pkg/push/domain"
	btproto "github.com/bucketeer-io/bucketeer/proto/batch"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
)

const (
	listRequestSize = 500
	fcmSendURL      = "https://fcm.googleapis.com/fcm/send"
	topicPrefix     = "bucketeer-"
)

type pushSender struct {
	pushClient    pushclient.Client
	featureClient featureclient.Client
	batchClient   btclient.Client
	logger        *zap.Logger
}

func NewPushSender(
	pushClient pushclient.Client,
	featureClient featureclient.Client,
	batchClient btclient.Client,
	logger *zap.Logger,
) subscriber.Processor {
	return &pushSender{
		pushClient:    pushClient,
		featureClient: featureClient,
		batchClient:   batchClient,
		logger:        logger,
	}
}

func (p pushSender) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	record := func(code codes.Code, startTime time.Time) {
		subscriberHandledCounter.WithLabelValues(subscriberPushSender, code.String()).Inc()
		subscriberHandledHistogram.WithLabelValues(
			subscriberPushSender,
			code.String(),
		).Observe(time.Since(startTime).Seconds())
	}
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberPushSender).Inc()
			startTime := time.Now()
			if id := msg.Attributes["id"]; id == "" {
				msg.Ack()
				record(codes.MissingID, startTime)
				continue
			}
			p.handle(msg)
			msg.Ack()
			record(codes.OK, startTime)
		case <-ctx.Done():
			return nil
		}
	}
}

func (p pushSender) handle(msg *puller.Message) {
	event, err := p.unmarshalMessage(msg)
	if err != nil {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberPushSender, codes.BadMessage.String()).Inc()
		return
	}
	featureID, isTarget := p.extractFeatureID(event)
	if !isTarget {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberPushSender, codes.OK.String()).Inc()
		return
	}
	if featureID == "" {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberPushSender, codes.BadMessage.String()).Inc()
		p.logger.Warn("Message contains an empty FeatureID", zap.Any("event", event))
		return
	}
	if err := p.send(featureID, event.EnvironmentNamespace); err != nil {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberPushSender, codes.NonRepeatableError.String()).Inc()
		return
	}
	msg.Ack()
	subscriberHandledCounter.WithLabelValues(subscriberPushSender, codes.OK.String()).Inc()
}

func (p pushSender) send(featureID, environmentNamespace string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := p.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		Id:                   featureID,
		EnvironmentNamespace: environmentNamespace,
	})
	if err != nil {
		return err
	}
	pushes, err := p.listPushes(ctx, environmentNamespace)
	if err != nil {
		return err
	}
	if len(pushes) == 0 {
		p.logger.Info("No pushes",
			zap.String("featureId", featureID),
			zap.String("environmentNamespace", environmentNamespace),
		)
		return nil
	}
	// Before sending the push notification we must update the cache
	// so the api-gateway can evaluate the user correctly.
	if err := p.updateFeatureFlagCache(ctx); err != nil {
		p.logger.Error("Failed to update feature flag cache", zap.Error(err))
		return err
	}
	if err := p.updateSegmentUserCache(ctx); err != nil {
		p.logger.Error("Failed to update segment user cache", zap.Error(err))
		return err
	}
	var lastErr error
	for _, push := range pushes {
		d := pushdomain.Push{Push: push}
		for _, tag := range resp.Feature.Tags {
			if !d.ExistTag(tag) {
				continue
			}
			topic := topicPrefix + tag
			if err = p.pushFCM(ctx, d.FcmApiKey, topic); err != nil {
				p.logger.Error("Failed to push notification", zap.Error(err),
					zap.String("featureId", featureID),
					zap.String("tag", tag),
					zap.String("topic", topic),
					zap.String("pushId", d.Push.Id),
					zap.String("environmentNamespace", environmentNamespace),
				)
				lastErr = err
				continue
			}
			p.logger.Info("Succeeded to push notification",
				zap.String("featureId", featureID),
				zap.String("tag", tag),
				zap.String("topic", topic),
				zap.String("pushId", d.Push.Id),
				zap.String("environmentNamespace", environmentNamespace),
			)
		}
	}
	return lastErr
}

func (p pushSender) pushFCM(ctx context.Context, fcmAPIKey, topic string) error {
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

func (p pushSender) listPushes(ctx context.Context, environmentNamespace string) ([]*pushproto.Push, error) {
	var pushes []*pushproto.Push
	cursor := ""
	for {
		resp, err := p.pushClient.ListPushes(ctx, &pushproto.ListPushesRequest{
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

func (p pushSender) unmarshalMessage(msg *puller.Message) (*domaineventproto.Event, error) {
	event := &domaineventproto.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		p.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
		return nil, err
	}
	return event, nil
}

func (p pushSender) extractFeatureID(event *domaineventproto.Event) (string, bool) {
	if event.EntityType != domaineventproto.Event_FEATURE {
		return "", false
	}
	if event.Type != domaineventproto.Event_FEATURE_VERSION_INCREMENTED &&
		event.Type != domaineventproto.Event_FEATURE_UPDATED {
		return "", false
	}
	return event.EntityId, true
}

// The batch API updates the feature flag cache in all the environments
func (p pushSender) updateFeatureFlagCache(ctx context.Context) error {
	req := &btproto.BatchJobRequest{
		Job: btproto.BatchJob_FeatureFlagCacher,
	}
	_, err := p.batchClient.ExecuteBatchJob(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// The batch API updates the segment user cache in all the environments
func (p pushSender) updateSegmentUserCache(ctx context.Context) error {
	req := &btproto.BatchJobRequest{
		Job: btproto.BatchJob_SegmentUserCacher,
	}
	_, err := p.batchClient.ExecuteBatchJob(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
