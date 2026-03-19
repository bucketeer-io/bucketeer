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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	insightsclient "github.com/bucketeer-io/bucketeer/v2/pkg/insights/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	insightsproto "github.com/bucketeer-io/bucketeer/v2/proto/insights"
)

const (
	timeout = 60 * time.Second
)

var (
	webGatewayAddr   = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort   = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert   = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath       = flag.String("api-key", "", "Client SDK API key for api-gateway")
	apiKeyServerPath = flag.String("api-key-server", "", "Server SDK API key for api-gateway")
	gatewayAddr      = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort      = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert      = flag.String("gateway-cert", "", "Gateway crt file")
	serviceTokenPath = flag.String("service-token", "", "Service token path")
	environmentID    = flag.String("environment-id", "", "Environment id")
	testID           = flag.String("test-id", "", "test ID")
	organizationID   = flag.String("organization-id", "", "Organization ID")

)

func TestGetInsightsMonthlySummary(t *testing.T) {
	c := newInsightsClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	t.Run("success", func(t *testing.T) {
		_, err := c.GetInsightsMonthlySummary(ctx, &insightsproto.GetInsightsMonthlySummaryRequest{
			EnvironmentIds: []string{*environmentID},
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error: missing environment_ids", func(t *testing.T) {
		_, err := c.GetInsightsMonthlySummary(ctx, &insightsproto.GetInsightsMonthlySummaryRequest{
			EnvironmentIds: []string{},
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected gRPC status error, got: %v", err)
		}
		if st.Code() != codes.InvalidArgument {
			t.Fatalf("expected InvalidArgument, got: %s", st.Code())
		}
	})
}

func TestGetInsightsTimeSeries(t *testing.T) {
	c := newInsightsClient(t)
	defer c.Close()

	now := time.Now()
	startAt := now.Add(-1 * time.Hour).Unix()
	endAt := now.Unix()

	t.Run("latency: success or not configured", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.GetInsightsLatency(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
			EnvironmentIds: []string{*environmentID},
			StartAt:        startAt,
			EndAt:          endAt,
		})
		if err != nil {
			assertStatusCode(t, err, codes.NotFound)
		}
	})

	t.Run("requests: success or not configured", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.GetInsightsRequests(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
			EnvironmentIds: []string{*environmentID},
			StartAt:        startAt,
			EndAt:          endAt,
		})
		if err != nil {
			assertStatusCode(t, err, codes.NotFound)
		}
	})

	t.Run("evaluations: success or not configured", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.GetInsightsEvaluations(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
			EnvironmentIds: []string{*environmentID},
			StartAt:        startAt,
			EndAt:          endAt,
		})
		if err != nil {
			assertStatusCode(t, err, codes.NotFound)
		}
	})

	t.Run("error_rates: success or not configured", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.GetInsightsErrorRates(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
			EnvironmentIds: []string{*environmentID},
			StartAt:        startAt,
			EndAt:          endAt,
		})
		if err != nil {
			assertStatusCode(t, err, codes.NotFound)
		}
	})

	t.Run("validation: missing environment_ids", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.GetInsightsLatency(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
			StartAt: startAt,
			EndAt:   endAt,
		})
		assertStatusCode(t, err, codes.InvalidArgument)
	})

	t.Run("validation: missing start_at", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.GetInsightsLatency(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
			EnvironmentIds: []string{*environmentID},
			EndAt:          endAt,
		})
		assertStatusCode(t, err, codes.InvalidArgument)
	})

	t.Run("validation: missing end_at", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.GetInsightsLatency(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
			EnvironmentIds: []string{*environmentID},
			StartAt:        startAt,
		})
		assertStatusCode(t, err, codes.InvalidArgument)
	})

	t.Run("validation: start_at after end_at", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.GetInsightsLatency(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
			EnvironmentIds: []string{*environmentID},
			StartAt:        endAt,
			EndAt:          startAt,
		})
		assertStatusCode(t, err, codes.InvalidArgument)
	})

	t.Run("validation: query range too large", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.GetInsightsLatency(ctx, &insightsproto.GetInsightsTimeSeriesRequest{
			EnvironmentIds: []string{*environmentID},
			StartAt:        now.Add(-32 * 24 * time.Hour).Unix(),
			EndAt:          endAt,
		})
		assertStatusCode(t, err, codes.InvalidArgument)
	})
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
