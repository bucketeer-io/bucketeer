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

package insights

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	gatewayclient "github.com/bucketeer-io/bucketeer/v2/pkg/api/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	insightsclient "github.com/bucketeer-io/bucketeer/v2/pkg/insights/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	insightsproto "github.com/bucketeer-io/bucketeer/v2/proto/insights"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	prefixTestName = "e2e-insights"
	timeout        = 3 * time.Minute
	retryInterval  = 5 * time.Second
)

var (
	webGatewayAddr   = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort   = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert   = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath       = flag.String("api-key", "", "Client SDK API key for api-gateway")
	_                = flag.String("api-key-server", "", "Server SDK API key for api-gateway")
	gatewayAddr      = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort      = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert      = flag.String("gateway-cert", "", "Gateway crt file")
	serviceTokenPath = flag.String("service-token", "", "Service token path")
	environmentID    = flag.String("environment-id", "", "Environment id")
	testID           = flag.String("test-id", "", "test ID")
	_                = flag.String("organization-id", "", "Organization ID")
)

func TestGetInsightsMonthlySummary(t *testing.T) {
	t.Parallel()

	// --- Setup: create feature and send evaluation events ---
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := newFeatureID(t)
	tag := newTagName(t)
	createFeature(t, featureClient, featureID, tag)
	t.Cleanup(func() { archiveFeature(t, featureID) })
	f := getFeature(t, featureClient, featureID)

	// MonthlySummarizer processes yesterday's data, so we use yesterday's timestamp.
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	userIDs := make([]string, 5)
	for i := range userIDs {
		userIDs[i] = newUserID(t)
	}
	registerEvaluationEvents(t, f.Id, f.Version, userIDs, f.Variations[0].Id, tag, yesterday)

	// --- Verify ---
	insightsClient := newInsightsClient(t)
	defer insightsClient.Close()

	expectedYearmonth := yesterday.Format("200601")
	var found bool
	deadline := time.Now().Add(timeout)
	for !found && time.Now().Before(deadline) {
		time.Sleep(retryInterval)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := insightsClient.GetInsightsMonthlySummary(ctx, &insightsproto.GetInsightsMonthlySummaryRequest{
			EnvironmentIds: []string{*environmentID},
			SourceIds:      []eventproto.SourceId{eventproto.SourceId_ANDROID},
		})
		cancel()
		if err != nil {
			t.Logf("polling: %v", err)
			continue
		}
		for _, series := range resp.Series {
			if series.EnvironmentId != *environmentID || series.SourceId != eventproto.SourceId_ANDROID {
				continue
			}
			for _, dp := range series.Data {
				if dp.Yearmonth == expectedYearmonth && dp.Mau > 0 {
					found = true
					break
				}
			}
		}
	}
	if !found {
		t.Fatalf("timeout: monthly summary MAU not found for yearmonth=%s", expectedYearmonth)
	}
}

func TestGetInsightsMonthlySummaryValidation(t *testing.T) {
	t.Parallel()
	c := newInsightsClient(t)
	t.Cleanup(func() { c.Close() })

	t.Run("missing environment_ids", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.GetInsightsMonthlySummary(ctx, &insightsproto.GetInsightsMonthlySummaryRequest{
			EnvironmentIds: []string{},
		})
		assertStatusCode(t, err, codes.InvalidArgument)
	})
}

func TestGetInsightsTimeSeries(t *testing.T) {
	t.Parallel()

	// --- Setup: create feature and generate gateway traffic ---
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := newFeatureID(t)
	tag := newTagName(t)
	createFeature(t, featureClient, featureID, tag)
	t.Cleanup(func() { archiveFeature(t, featureID) })
	f := getFeature(t, featureClient, featureID)
	generateGatewayTraffic(t, f, tag)

	// --- Tests ---
	insightsClient := newInsightsClient(t)
	t.Cleanup(func() { insightsClient.Close() })

	tests := []struct {
		name string
		call func(
			ctx context.Context,
			in *insightsproto.GetInsightsTimeSeriesRequest,
			opts ...grpc.CallOption,
		) (*insightsproto.GetInsightsTimeSeriesResponse, error)
		// requireData indicates the endpoint must return at least one data point.
		requireData bool
	}{
		{"latency", insightsClient.GetInsightsLatency, true},
		{"requests", insightsClient.GetInsightsRequests, true},
		{"evaluations", insightsClient.GetInsightsEvaluations, true},
		// error_rates may be empty when no errors have occurred.
		{"error_rates", insightsClient.GetInsightsErrorRates, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if !tt.requireData {
				now := time.Now()
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				_, err := tt.call(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
					EnvironmentIds: []string{*environmentID},
					StartAt:        now.Add(-10 * time.Minute).Unix(),
					EndAt:          now.Unix(),
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			var found bool
			deadline := time.Now().Add(timeout)
			for !found && time.Now().Before(deadline) {
				time.Sleep(retryInterval)
				now := time.Now()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				resp, err := tt.call(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
					EnvironmentIds: []string{*environmentID},
					StartAt:        now.Add(-10 * time.Minute).Unix(),
					EndAt:          now.Unix(),
				})
				cancel()
				if err != nil {
					t.Logf("polling %s: %v", tt.name, err)
					continue
				}
				for _, ts := range resp.Timeseries {
					if len(ts.Data) > 0 {
						found = true
						break
					}
				}
			}
			if !found {
				t.Fatalf("timeout: no %s time series data found for environment=%s", tt.name, *environmentID)
			}
		})
	}
}

func generateGatewayTraffic(t *testing.T, f *featureproto.Feature, tag string) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	userID := newUserID(t)

	// Call RegisterEvents to generate request count and latency metrics.
	for i := range 3 {
		ev, err := anypb.New(&eventproto.EvaluationEvent{
			Timestamp:      time.Now().Unix(),
			FeatureId:      f.Id,
			FeatureVersion: f.Version,
			UserId:         userID,
			VariationId:    f.Variations[0].Id,
			User:           &userproto.User{Id: userID},
			Reason:         &featureproto.Reason{},
			Tag:            tag,
			SourceId:       eventproto.SourceId_ANDROID,
		})
		if err != nil {
			t.Fatal(err)
		}
		resp, err := c.RegisterEvents(ctx, &gatewayproto.RegisterEventsRequest{
			Events: []*eventproto.Event{{Id: newUUID(t), Event: ev}},
		})
		if err != nil {
			t.Fatalf("RegisterEvents call %d failed: %v", i, err)
		}
		if len(resp.Errors) > 0 {
			t.Logf("RegisterEvents call %d had errors: %v", i, resp.Errors)
		}
	}

	// Call GetEvaluations to generate evaluation counter metrics.
	for i := range 3 {
		_, err := c.GetEvaluations(ctx, &gatewayproto.GetEvaluationsRequest{
			Tag:      tag,
			User:     &userproto.User{Id: userID},
			SourceId: eventproto.SourceId_ANDROID,
		})
		if err != nil {
			t.Logf("GetEvaluations call %d: %v", i, err)
		}
	}
}

func registerEvaluationEvents(
	t *testing.T,
	featureID string,
	featureVersion int32,
	userIDs []string,
	variationID, tag string,
	timestamp time.Time,
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	events := make([]*eventproto.Event, 0, len(userIDs))
	for _, userID := range userIDs {
		evaluation, err := anypb.New(&eventproto.EvaluationEvent{
			Timestamp:      timestamp.Unix(),
			FeatureId:      featureID,
			FeatureVersion: featureVersion,
			UserId:         userID,
			VariationId:    variationID,
			User:           &userproto.User{Id: userID},
			Reason:         &featureproto.Reason{},
			Tag:            tag,
			SourceId:       eventproto.SourceId_ANDROID,
		})
		if err != nil {
			t.Fatal(err)
		}
		events = append(events, &eventproto.Event{
			Id:    newUUID(t),
			Event: evaluation,
		})
	}

	response, err := c.RegisterEvents(ctx, &gatewayproto.RegisterEventsRequest{
		Events: events,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register evaluation events: %v", response.Errors)
	}
}

func createFeature(t *testing.T, client featureclient.Client, featureID, tag string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.CreateFeature(ctx, &featureproto.CreateFeatureRequest{
		Id:          featureID,
		Name:        featureID,
		Description: "e2e-test-insights-feature",
		Variations: []*featureproto.Variation{
			{Value: "true", Name: "True", Description: "Enabled"},
			{Value: "false", Name: "False", Description: "Disabled"},
		},
		Tags:                     []string{tag},
		DefaultOnVariationIndex:  &wrapperspb.Int32Value{Value: 0},
		DefaultOffVariationIndex: &wrapperspb.Int32Value{Value: 1},
		EnvironmentId:            *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
		Id:            featureID,
		Enabled:       wrapperspb.Bool(true),
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatalf("Failed to enable feature %s: %v", featureID, err)
	}
}

func getFeature(t *testing.T, client featureclient.Client, featureID string) *featureproto.Feature {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.GetFeature(ctx, &featureproto.GetFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatalf("Failed to get feature %s: %v", featureID, err)
	}
	return resp.Feature
}

// archiveFeature is a best-effort cleanup helper.
// Since this is optional, errors are only logged.
func archiveFeature(t *testing.T, featureID string) {
	t.Helper()
	c := newFeatureClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := c.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
		Id:            featureID,
		Archived:      wrapperspb.Bool(true),
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Logf("failed to archive feature %s: %v", featureID, err)
	}
}

func assertStatusCode(t *testing.T, err error, expected codes.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %s, got nil", expected)
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got: %v", err)
	}
	if st.Code() != expected {
		t.Fatalf("expected %s, got: %s (message: %s)", expected, st.Code(), st.Message())
	}
}

func newInsightsClient(t *testing.T) insightsclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	c, err := insightsclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create insights client:", err)
	}
	return c
}

func newGatewayClient(t *testing.T) gatewayclient.Client {
	t.Helper()
	creds, err := gatewayclient.NewPerRPCCredentials(*apiKeyPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	c, err := gatewayclient.NewClient(
		fmt.Sprintf("%s:%d", *gatewayAddr, *gatewayPort),
		*gatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create gateway client:", err)
	}
	return c
}

func newFeatureClient(t *testing.T) featureclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	c, err := featureclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create feature client:", err)
	}
	return c
}

func newFeatureID(t *testing.T) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-feature-id-%s", prefixTestName, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-feature-id-%s", prefixTestName, newUUID(t))
}

func newTagName(t *testing.T) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-tag-%s", prefixTestName, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-tag-%s", prefixTestName, newUUID(t))
}

func newUserID(t *testing.T) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-user-%s", prefixTestName, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-user-%s", prefixTestName, newUUID(t))
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}
