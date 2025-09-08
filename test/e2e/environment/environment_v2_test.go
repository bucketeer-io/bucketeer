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

package environment

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

const (
	timeout = 60 * time.Second
)

var (
	// FIXME: To avoid compiling the test many times, webGatewayAddr, webGatewayPort & apiKey has been also added here to prevent from getting: "flag provided but not defined" error during the test. These 3 are being use in the Gateway test
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
	organizationID   = flag.String("organization-id", "", "Organization ID")
	testID           = flag.String("test-id", "", "test ID")
)

const (
	environmentName = "E2E environment"
)

func TestGetEnvironmentV2(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newEnvironmentClient(t)
		defer c.Close()
		id := getEnvironmentID(t)
		resp, err := c.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{Id: id})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Environment.Id != id {
			t.Fatalf("different ids, expected: %v, actual: %v", id, resp.Environment.Id)
		}
		if resp.Environment.Name != environmentName {
			t.Fatalf("different name, expected: %v, actual: %v", environmentName, resp.Environment.Name)
		}
	})
}

func TestListEnvironmentsV2ByProject(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newEnvironmentClient(t)
		defer c.Close()
		resp, err := c.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{ProjectId: defaultProjectID})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Environments) == 0 {
			t.Fatal("environments is empty, expected at least 1")
		}
		for _, env := range resp.Environments {
			if env.ProjectId != defaultProjectID {
				t.Fatalf("different project id, expected: %s, actual: %s", defaultProjectID, env.ProjectId)
			}
		}
	})
}

func TestListEnvironmentsV2(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newEnvironmentClient(t)
		defer c.Close()
		pageSize := int64(1)
		resp, err := c.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{PageSize: pageSize})
		if err != nil {
			t.Fatal(err)
		}
		responseSize := int64(len(resp.Environments))
		if responseSize != pageSize {
			t.Fatalf("different sizes, expected: %d actual: %d", pageSize, responseSize)
		}
	})
}

func TestUpdateEnvironmentV2(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newEnvironmentClient(t)
		defer c.Close()
		id := getEnvironmentID(t)
		newDesc := fmt.Sprintf("This environment is for local development (Updated at %d)", time.Now().Unix())
		_, err := c.UpdateEnvironmentV2(ctx, &environmentproto.UpdateEnvironmentV2Request{
			Id:                       id,
			ChangeDescriptionCommand: &environmentproto.ChangeDescriptionEnvironmentV2Command{Description: newDesc},
		})
		if err != nil {
			t.Fatal(err)
		}
		getResp, err := c.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{Id: id})
		if err != nil {
			t.Fatal(err)
		}
		if getResp.Environment.Id != id {
			t.Fatalf("different ids, expected: %v, actual: %v", id, getResp.Environment.Id)
		}
		if getResp.Environment.Name != environmentName {
			t.Fatalf("different name, expected: %v, actual: %v", environmentName, getResp.Environment.Name)
		}
		if getResp.Environment.Description != newDesc {
			t.Fatalf("different descriptions, expected: %v, actual: %v", newDesc, getResp.Environment.Description)
		}

		newDesc = fmt.Sprintf("This environment is for local development (Updated at %d with no command)", time.Now().Unix())
		_, err = c.UpdateEnvironmentV2(ctx, &environmentproto.UpdateEnvironmentV2Request{
			Id:          id,
			Description: wrapperspb.String(newDesc),
		})
		if err != nil {
			t.Fatal(err)
		}
		getResp, err = c.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{Id: id})
		if err != nil {
			t.Fatal(err)
		}
		if getResp.Environment.Id != id {
			t.Fatalf("different ids, expected: %v, actual: %v", id, getResp.Environment.Id)
		}
		if getResp.Environment.Name != environmentName {
			t.Fatalf("different name, expected: %v, actual: %v", environmentName, getResp.Environment.Name)
		}
		if getResp.Environment.Description != newDesc {
			t.Fatalf("different descriptions, expected: %v, actual: %v", newDesc, getResp.Environment.Description)
		}
	})
}

func getEnvironmentID(t *testing.T) string {
	t.Helper()
	if *environmentID == "" {
		return "production"
	}
	return *environmentID
}

func newEnvironmentClient(t *testing.T) environmentclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := environmentclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create environment client:", err)
	}
	return client
}
