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

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"

	gwapi "github.com/bucketeer-io/bucketeer/pkg/api/api"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const (
	defaultNumEvaluationEvents = 10000
	defaultNumGoalEvents       = 5000
	defaultConcurrentWorkers   = 10
	defaultUserDataSize        = 10 // Number of random key-value pairs in userData
	defaultJSONFieldSize       = 20 // Length of random string values

	// API paths
	version          = "/v1"
	service          = "/gateway"
	eventsAPI        = "/events"
	authorizationKey = "authorization"
)

var (
	flags = flag.NewFlagSet("event-load-test", flag.ContinueOnError)

	// Gateway connection params
	gatewayAddr = flags.String("gateway-addr", "localhost:8443", "Gateway host:port")
	apiKeyPath  = flags.String("api-key-path", "", "Path to API key file")
	insecureTLS = flags.Bool("insecure-tls", true, "Skip TLS certificate verification")

	// Test parameters
	numEvalEvents     = flags.Int("num-eval-events", defaultNumEvaluationEvents, "Number of evaluation events to generate")
	numGoalEvents     = flags.Int("num-goal-events", defaultNumGoalEvents, "Number of goal events to generate")
	concurrentWorkers = flags.Int("concurrent-workers", defaultConcurrentWorkers, "Number of concurrent workers for sending events")
	environmentID     = flags.String("environment-id", "", "Environment ID to use for events")
	featureCount      = flags.Int("feature-count", 10, "Number of unique features to use")
	userCount         = flags.Int("user-count", 100, "Number of unique users to use")
	tagCount          = flags.Int("tag-count", 5, "Number of unique tags to use")
	goalCount         = flags.Int("goal-count", 5, "Number of unique goals to use")

	// Random data generation
	userDataSize  = flags.Int("user-data-size", defaultUserDataSize, "Size of user data map")
	jsonFieldSize = flags.Int("json-field-size", defaultJSONFieldSize, "Size of random string fields")

	// JSON marshaler for protobuf
	jsonMarshaler = protojson.MarshalOptions{
		UseProtoNames: true,
	}
)

type registerEventsRequest struct {
	Events []event `json:"events,omitempty"`
}

type registerEventsResponse struct {
	Errors map[string]interface{} `json:"errors,omitempty"`
}

type successResponse struct {
	Data json.RawMessage `json:"data"`
}

type event struct {
	ID            string          `json:"id,omitempty"`
	Event         json.RawMessage `json:"event,omitempty"`
	EnvironmentID string          `json:"environment_id,omitempty"`
	Type          gwapi.EventType `json:"type,omitempty"`
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	// Validate required flags
	if *apiKeyPath == "" {
		log.Fatal("API key path is required")
	}

	// Initialize logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Generate test data
	logger.Info("Generating random test data...")
	featureIDs := generateFeatureIDs(*featureCount)
	userIDs := generateUserIDs(*userCount)
	tags := generateTags(*tagCount)
	goalIDs := generateGoalIDs(*goalCount)

	// Read API key
	apiKey, err := os.ReadFile(*apiKeyPath)
	if err != nil {
		logger.Fatal("Failed to read API key", zap.Error(err))
	}

	// Create HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: *insecureTLS,
			},
		},
		Timeout: 30 * time.Second,
	}

	// Create event URL
	eventsURL := fmt.Sprintf("https://%s%s%s%s", *gatewayAddr, version, service, eventsAPI)

	// Generate and send evaluation events
	logger.Info("Sending evaluation events...",
		zap.Int("count", *numEvalEvents),
		zap.Int("concurrentWorkers", *concurrentWorkers))

	start := time.Now()

	evalEvents := generateEvaluationEvents(
		*numEvalEvents,
		featureIDs,
		userIDs,
		tags,
	)

	sendEvents(evalEvents, eventsURL, string(apiKey), client, *concurrentWorkers, logger)

	evalDuration := time.Since(start)

	// Generate and send goal events
	logger.Info("Sending goal events...",
		zap.Int("count", *numGoalEvents),
		zap.Int("concurrentWorkers", *concurrentWorkers))

	start = time.Now()

	goalEvents := generateGoalEvents(
		*numGoalEvents,
		goalIDs,
		userIDs,
		featureIDs,
		tags,
	)

	sendEvents(goalEvents, eventsURL, string(apiKey), client, *concurrentWorkers, logger)

	goalDuration := time.Since(start)

	// Print summary
	logger.Info("Load test completed",
		zap.Int("evaluationEventsSent", *numEvalEvents),
		zap.Duration("evaluationDuration", evalDuration),
		zap.Float64("evaluationEventsPerSecond", float64(*numEvalEvents)/evalDuration.Seconds()),
		zap.Int("goalEventsSent", *numGoalEvents),
		zap.Duration("goalDuration", goalDuration),
		zap.Float64("goalEventsPerSecond", float64(*numGoalEvents)/goalDuration.Seconds()))
}

func sendEvents(events []event, url, apiKey string, client *http.Client, workers int, logger *zap.Logger) {
	// Calculate batch size based on number of events and workers
	batchSize := (len(events) + workers - 1) / workers

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup

	// Keep track of successful and failed events
	successCount := 0
	failureCount := 0
	var countMutex sync.Mutex

	// Create and start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)

		// Calculate the batch range for this worker
		start := i * batchSize
		end := start + batchSize
		if end > len(events) {
			end = len(events)
		}

		// Skip if start is beyond the range
		if start >= len(events) {
			wg.Done()
			continue
		}

		// Start the worker
		go func(workerID int, eventBatch []event) {
			defer wg.Done()

			// Process the batch in smaller chunks to avoid large requests
			const chunkSize = 50
			for j := 0; j < len(eventBatch); j += chunkSize {
				chunkEnd := j + chunkSize
				if chunkEnd > len(eventBatch) {
					chunkEnd = len(eventBatch)
				}

				chunk := eventBatch[j:chunkEnd]
				success, err := sendEventBatch(chunk, url, apiKey, client)

				countMutex.Lock()
				successCount += success
				failureCount += len(chunk) - success
				countMutex.Unlock()

				if err != nil {
					logger.Error("Error sending event batch",
						zap.Int("workerID", workerID),
						zap.Int("batchStart", j),
						zap.Int("batchEnd", chunkEnd),
						zap.Error(err))
				} else {
					logger.Info("Sent event batch",
						zap.Int("workerID", workerID),
						zap.Int("batchStart", j),
						zap.Int("batchEnd", chunkEnd),
						zap.Int("successCount", success),
						zap.Int("failureCount", len(chunk)-success))
				}

				// Add a small delay to avoid overwhelming the server
				time.Sleep(10 * time.Millisecond)
			}
		}(i, events[start:end])
	}

	// Wait for all workers to finish
	wg.Wait()

	logger.Info("Events sending completed",
		zap.Int("successCount", successCount),
		zap.Int("failureCount", failureCount))
}

func sendEventBatch(events []event, url, apiKey string, client *http.Client) (int, error) {
	req := registerEventsRequest{
		Events: events,
	}

	// Marshal the request body
	reqBody, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Add(authorizationKey, apiKey)
	httpReq.Header.Add("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	// Parse the response
	var sr successResponse
	if err = json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for errors in the response
	var rer registerEventsResponse
	if err = json.Unmarshal(sr.Data, &rer); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	// Count successful events
	return len(events) - len(rer.Errors), nil
}

func generateFeatureIDs(count int) []string {
	ids := make([]string, count)
	for i := 0; i < count; i++ {
		ids[i] = fmt.Sprintf("feature-load-test-%d", i+1)
	}
	return ids
}

func generateUserIDs(count int) []string {
	ids := make([]string, count)
	for i := 0; i < count; i++ {
		ids[i] = fmt.Sprintf("user-load-test-%d", i+1)
	}
	return ids
}

func generateTags(count int) []string {
	tags := make([]string, count)
	for i := 0; i < count; i++ {
		tags[i] = fmt.Sprintf("tag-load-test-%d", i+1)
	}
	return tags
}

func generateGoalIDs(count int) []string {
	ids := make([]string, count)
	for i := 0; i < count; i++ {
		ids[i] = fmt.Sprintf("goal-load-test-%d", i+1)
	}
	return ids
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateUserData() map[string]string {
	userData := make(map[string]string, *userDataSize)
	for i := 0; i < *userDataSize; i++ {
		key := fmt.Sprintf("attribute-%d", i+1)
		userData[key] = generateRandomString(*jsonFieldSize)
	}
	return userData
}

func generateEvaluationEvents(count int, featureIDs, userIDs, tags []string) []event {
	events := make([]event, count)

	for i := 0; i < count; i++ {
		featureID := featureIDs[rand.Intn(len(featureIDs))]
		userID := userIDs[rand.Intn(len(userIDs))]
		tag := tags[rand.Intn(len(tags))]

		evalEvent := &eventproto.EvaluationEvent{
			Timestamp:      time.Now().Unix(),
			FeatureId:      featureID,
			FeatureVersion: int32(rand.Intn(5) + 1),
			UserId:         userID,
			VariationId:    fmt.Sprintf("variation-%d", rand.Intn(3)+1),
			User: &userproto.User{
				Id:   userID,
				Data: generateUserData(),
			},
			Reason: &featureproto.Reason{
				Type: featureproto.Reason_Type(rand.Intn(5) + 1),
			},
			Tag:      tag,
			SourceId: eventproto.SourceId(rand.Intn(3) + 1),
		}

		// Use protojson to marshal the protobuf message
		eventJSON, _ := jsonMarshaler.Marshal(evalEvent)

		events[i] = event{
			ID:            uuid.New().String(),
			Event:         eventJSON,
			EnvironmentID: *environmentID,
			Type:          gwapi.EvaluationEventType,
		}
	}

	return events
}

func generateGoalEvents(count int, goalIDs, userIDs, featureIDs, tags []string) []event {
	events := make([]event, count)

	for i := 0; i < count; i++ {
		goalID := goalIDs[rand.Intn(len(goalIDs))]
		userID := userIDs[rand.Intn(len(userIDs))]
		tag := tags[rand.Intn(len(tags))]

		// Convert float32 to float64 for the Value field
		value := float64(rand.Float32() * 100)

		goalEvent := &eventproto.GoalEvent{
			Timestamp: time.Now().Unix(),
			GoalId:    goalID,
			UserId:    userID,
			Value:     value,
			User: &userproto.User{
				Id:   userID,
				Data: generateUserData(),
			},
			Tag:      tag,
			SourceId: eventproto.SourceId(rand.Intn(3) + 1),
		}

		// Randomly decide to link with a feature evaluation (50% chance)
		if rand.Intn(2) == 1 {
			featureID := featureIDs[rand.Intn(len(featureIDs))]
			variationID := fmt.Sprintf("variation-%d", rand.Intn(3)+1)
			featureVersion := int32(rand.Intn(5) + 1)

			// Create evaluation for the goal event
			goalEvent.Evaluations = []*featureproto.Evaluation{
				{
					Id:             uuid.New().String(),
					FeatureId:      featureID,
					FeatureVersion: featureVersion,
					UserId:         userID,
					VariationId:    variationID,
				},
			}
		}

		// Use protojson to marshal the protobuf message
		eventJSON, _ := jsonMarshaler.Marshal(goalEvent)

		events[i] = event{
			ID:            uuid.New().String(),
			Event:         eventJSON,
			EnvironmentID: *environmentID,
			Type:          gwapi.GoalEventType,
		}
	}

	return events
}
