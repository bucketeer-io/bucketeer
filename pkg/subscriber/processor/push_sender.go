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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/protobuf/proto"

	btclient "github.com/bucketeer-io/bucketeer/pkg/batch/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	pushdomain "github.com/bucketeer-io/bucketeer/pkg/push/domain"
	pushstorage "github.com/bucketeer-io/bucketeer/pkg/push/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
	btproto "github.com/bucketeer-io/bucketeer/proto/batch"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
)

const (
	topicPrefix = "bucketeer"
	timeout     = time.Minute
)

var (
	defaultTopic = fmt.Sprintf("%s-default", topicPrefix)
)

type pushSender struct {
	featureClient featureclient.Client
	batchClient   btclient.Client
	mysqlClient   mysql.Client
	logger        *zap.Logger
}

type fcmConfig struct {
	endpointURL string
	accessToken string
}

func NewPushSender(
	featureClient featureclient.Client,
	batchClient btclient.Client,
	mysqlClient mysql.Client,
	logger *zap.Logger,
) subscriber.PubSubProcessor {
	return &pushSender{
		featureClient: featureClient,
		batchClient:   batchClient,
		mysqlClient:   mysqlClient,
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
	if err := p.send(featureID, event.EnvironmentId); err != nil {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberPushSender, codes.NonRepeatableError.String()).Inc()
		return
	}
	msg.Ack()
	subscriberHandledCounter.WithLabelValues(subscriberPushSender, codes.OK.String()).Inc()
}

func (p pushSender) send(featureID, environmentId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := p.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		Id:            featureID,
		EnvironmentId: environmentId,
	})
	if err != nil {
		p.logger.Error("Failed to get feature flag",
			zap.Error(err),
			zap.String("featureId", featureID),
			zap.String("environmentId", environmentId),
		)
		return err
	}
	pushes, err := p.listPushes(ctx, environmentId)
	if err != nil {
		p.logger.Error("Failed to list pushes",
			zap.Error(err),
			zap.String("featureId", featureID),
			zap.String("environmentId", environmentId),
		)
		return err
	}
	if len(pushes) == 0 {
		p.logger.Debug("No pushes",
			zap.String("featureId", featureID),
			zap.String("environmentId", environmentId),
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
		// Determine which topics to send notifications to
		topics := p.getTopicsForPush(resp.Feature, push)
		if len(topics) == 0 {
			continue
		}

		// Send notifications to all matching topics
		for _, topic := range topics {
			if err := p.sendPushNotification(ctx, topic, push, featureID, environmentId); err != nil {
				lastErr = err
			}
		}
	}
	return lastErr
}

// pushFCM sends a silent notification to all the devices subscrribed to the target topic
//
// IMPORTANT iOS Silent Notification Limitations:
// - Silent notifications are heavily throttled by iOS and not guaranteed to be delivered
// - They won't work when the device is in Low Power Mode
// - They won't work if the app has been force-quit by the user
// - They require Background App Refresh to be enabled for the app
// - Apple recommends sending them no more than once every 20 minutes
// - The system may delay or drop them based on battery level, memory usage, and other factors
//
// For critical updates that must reach iOS devices, consider using visible notifications
// with the "mutable-content" flag and a Notification Service Extension instead.
func (p pushSender) pushFCM(ctx context.Context, topic, fcmServiceAccount string) error {
	creds, err := p.getFCMCredentials(ctx, fcmServiceAccount)
	if err != nil {
		return err
	}
	message := map[string]interface{}{
		"message": map[string]interface{}{
			"topic": topic,
			// The values in the data payload should be converted to string type.
			"data": map[string]interface{}{
				"bucketeer_feature_flag_updated": "true",
			},
			"android": map[string]interface{}{
				"priority": "normal",
			},
			"apns": map[string]interface{}{
				"headers": map[string]string{
					"apns-priority":  "5",          // Normal priority for iOS
					"apns-push-type": "background", // Required for silent notifications
				},
				"payload": map[string]interface{}{
					"aps": map[string]interface{}{
						"content-available": 1, // Silent notification for iOS
					},
				},
			},
		},
	}
	requestBody, err := json.Marshal(message)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", creds.endpointURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", creds.accessToken))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("FCM request failed with status: %s, response body: %s", resp.Status, body)
	}
	return nil
}

func (p pushSender) getFCMCredentials(ctx context.Context, fcmServiceAccount string) (*fcmConfig, error) {
	// Create OAuth2 token source
	creds, err := google.CredentialsFromJSON(
		ctx,
		[]byte(fcmServiceAccount),
		"https://www.googleapis.com/auth/firebase.messaging",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials from JSON: %w", err)
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	return &fcmConfig{
		endpointURL: fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", creds.ProjectID),
		accessToken: token.AccessToken,
	}, nil
}

// listPushes list all the pushes
// Because the `ListPushes` API removes the FCM service account from the response
// due to security reasons, we list the pushes directly from the storage interface
func (p pushSender) listPushes(ctx context.Context, environmentId string) ([]*pushproto.Push, error) {
	options := &mysql.ListOptions{
		Limit:  mysql.QueryNoLimit,
		Offset: 0,
		Filters: []*mysql.FilterV2{
			{
				Column:   "deleted",
				Operator: mysql.OperatorEqual,
				Value:    false,
			},
			{
				Column:   "environment_id",
				Operator: mysql.OperatorEqual,
				Value:    environmentId,
			},
		},
		InFilters: nil,
		Orders:    nil,
	}

	storage := pushstorage.NewPushStorage(p.mysqlClient)
	pushes, _, _, err := storage.ListPushes(
		ctx,
		options,
	)
	if err != nil {
		return nil, err
	}
	return pushes, nil
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

// getTopicsForPush determines which topics should receive push notifications
// based on the feature tags and push configuration
func (p pushSender) getTopicsForPush(
	feature *featureproto.Feature,
	push *pushproto.Push,
) []string {
	// Handle case where both feature and push have no tags
	if len(feature.Tags) == 0 && len(push.Tags) == 0 {
		return []string{defaultTopic}
	}

	topics := []string{}
	// Handle normal case with tags
	d := pushdomain.Push{Push: push}
	for _, tag := range feature.Tags {
		if d.ExistTag(tag) {
			topics = append(topics, fmt.Sprintf("%s-%s", topicPrefix, tag))
		}
	}
	return topics
}

// sendPushNotification sends a push notification to a specific topic
func (p pushSender) sendPushNotification(
	ctx context.Context,
	topic string,
	push *pushproto.Push,
	featureID string,
	environmentID string,
) error {
	logFields := []zap.Field{
		zap.String("featureId", featureID),
		zap.String("topic", topic),
		zap.String("pushId", push.Id),
		zap.String("environmentId", environmentID),
	}
	if err := p.pushFCM(ctx, topic, push.FcmServiceAccount); err != nil {
		logFields = append(logFields, zap.Error(err))
		p.logger.Error("Failed to push notification", logFields...)
		return err
	}
	p.logger.Info("Succeeded to push notification for feature without tags", logFields...)

	return nil
}
