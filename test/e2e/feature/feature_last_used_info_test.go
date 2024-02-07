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

package feature

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/encoding/protojson"

	gwapi "github.com/bucketeer-io/bucketeer/pkg/gateway/api"
	gatewayclient "github.com/bucketeer-io/bucketeer/pkg/gateway/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	"github.com/bucketeer-io/bucketeer/proto/feature"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
	"github.com/bucketeer-io/bucketeer/test/e2e/util"
)

const (
	featureRecorderRetryTimes = 60
)

func TestGprcGetFeatureLastUsedInfo(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	cmd := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	enableFeature(t, cmd.Id, client)
	f := getFeature(t, cmd.Id, client)
	lastUsedAt := time.Now()
	grpcRegisterEvaluationEvents(t, []*feature.Feature{f}, f.Tags[0], lastUsedAt)
	for i := 0; i < featureRecorderRetryTimes; i++ {
		actual := getFeature(t, cmd.Id, client)
		if actual.LastUsedInfo != nil {
			if actual.LastUsedInfo.FeatureId != f.Id {
				t.Fatalf("feature ID is not correct: expected: %s, actual: %s", f.Id, actual.LastUsedInfo.FeatureId)
			}
			if actual.LastUsedInfo.Version != f.Version {
				t.Fatalf("feature version is not correct: expected: %d, actual: %d", f.Version, actual.LastUsedInfo.Version)
			}
			if actual.LastUsedInfo.CreatedAt != lastUsedAt.Unix() {
				t.Fatalf("created at is not correct: expected: %d, actual: %d", lastUsedAt.Unix(), actual.LastUsedInfo.CreatedAt)
			}
			if actual.LastUsedInfo.LastUsedAt != lastUsedAt.Unix() {
				t.Fatalf("lastUsedAt at is not correct: expected: %d, actual: %d", lastUsedAt.Unix(), actual.LastUsedInfo.LastUsedAt)
			}
			break
		}
		if i == featureRecorderRetryTimes-1 {
			t.Fatalf("LastUsedInfo cannot be fetched.")
		}
		time.Sleep(time.Second)
	}
}

func TestGetFeatureLastUsedInfo(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	cmd := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	enableFeature(t, cmd.Id, client)
	f := getFeature(t, cmd.Id, client)
	lastUsedAt := time.Now()
	registerEvaluationEvents(t, []*feature.Feature{f}, f.Tags[0], lastUsedAt)
	for i := 0; i < featureRecorderRetryTimes; i++ {
		actual := getFeature(t, cmd.Id, client)
		if actual.LastUsedInfo != nil {
			if actual.LastUsedInfo.FeatureId != f.Id {
				t.Fatalf("feature ID is not correct: expected: %s, actual: %s", f.Id, actual.LastUsedInfo.FeatureId)
			}
			if actual.LastUsedInfo.Version != f.Version {
				t.Fatalf("feature version is not correct: expected: %d, actual: %d", f.Version, actual.LastUsedInfo.Version)
			}
			if actual.LastUsedInfo.CreatedAt != lastUsedAt.Unix() {
				t.Fatalf("created at is not correct: expected: %d, actual: %d", lastUsedAt.Unix(), actual.LastUsedInfo.CreatedAt)
			}
			if actual.LastUsedInfo.LastUsedAt != lastUsedAt.Unix() {
				t.Fatalf("lastUsedAt at is not correct: expected: %d, actual: %d", lastUsedAt.Unix(), actual.LastUsedInfo.LastUsedAt)
			}
			break
		}
		if i == featureRecorderRetryTimes-1 {
			t.Fatalf("LastUsedInfo cannot be fetched.")
		}
		time.Sleep(time.Second)
	}
}

func grpcRegisterEvaluationEvents(t *testing.T, features []*feature.Feature, tag string, now time.Time) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	events := make([]*eventproto.Event, 0)
	for _, f := range features {
		evaluation, err := ptypes.MarshalAny(&eventproto.EvaluationEvent{
			Timestamp:      now.Unix(),
			FeatureId:      f.Id,
			FeatureVersion: f.Version,
			UserId:         "user-id",
			VariationId:    "variation-id",
			User:           &userproto.User{},
			Reason:         &featureproto.Reason{},
			Tag:            tag,
		})
		if err != nil {
			t.Fatal(err)
		}
		events = append(events, &eventproto.Event{
			Id:    newUUID(t),
			Event: evaluation,
		})
	}
	req := &gatewayproto.RegisterEventsRequest{Events: events}
	_, err := c.RegisterEvents(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
}

func registerEvaluationEvents(t *testing.T, features []*feature.Feature, tag string, now time.Time) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	events := make([]util.Event, 0)
	for _, f := range features {
		evaluation, err := protojson.Marshal(&eventproto.EvaluationEvent{
			Timestamp:      now.Unix(),
			FeatureId:      f.Id,
			FeatureVersion: f.Version,
			UserId:         "user-id",
			VariationId:    "variation-id",
			User:           &userproto.User{},
			Reason:         &featureproto.Reason{},
			Tag:            tag,
		})
		if err != nil {
			t.Fatal(err)
		}
		events = append(events, util.Event{
			ID:    newUUID(t),
			Event: evaluation,
			Type:  gwapi.EvaluationEventType,
		})
	}
	response := util.RegisterEvents(t, events, *gatewayAddr, *apiKeyPath)
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func newGatewayClient(t *testing.T) gatewayclient.Client {
	t.Helper()
	creds, err := gatewayclient.NewPerRPCCredentials(*apiKeyPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := gatewayclient.NewClient(
		fmt.Sprintf("%s:%d", *gatewayAddr, *gatewayPort),
		*gatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create gateway client:", err)
	}
	return client
}
