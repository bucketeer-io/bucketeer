// Copyright 2023 The Bucketeer Authors.
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
//

package batch

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	batchclient "github.com/bucketeer-io/bucketeer/pkg/batch/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	proto "github.com/bucketeer-io/bucketeer/proto/batch"
)

var (
	// FIXME: To avoid compiling the test many times, webGatewayAddr, webGatewayPort & apiKey has been also added here to prevent from getting: "flag provided but not defined" error during the test. These 3 are being use  in the Gateway test
	webGatewayAddr   = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort   = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert   = flag.String("web-gateway-cert", "", "Web gateway crt file")
	serviceTokenPath = flag.String("service-token", "", "Service token path")
)

const timeout = 10 * time.Second

func TestExperimentStatusUpdater(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newBatchServiceClient(t)
	defer client.Close()
	_, err := client.ExecuteBatchJob(ctx, &proto.BatchJobRequest{
		Job: proto.BatchJob_ExperimentStatusUpdater,
	})
	if err != nil {
		t.Fatal("failed to execute batch job:", err)
	}
}

func TestExperimentRunningWatcher(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newBatchServiceClient(t)
	defer client.Close()
	_, err := client.ExecuteBatchJob(ctx, &proto.BatchJobRequest{
		Job: proto.BatchJob_ExperimentRunningWatcher,
	})
	if err != nil {
		t.Fatal("failed to execute batch job:", err)
	}
}

func TestFeatureStaleWatcher(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newBatchServiceClient(t)
	defer client.Close()
	_, err := client.ExecuteBatchJob(ctx, &proto.BatchJobRequest{
		Job: proto.BatchJob_FeatureStaleWatcher,
	})
	if err != nil {
		t.Fatal("failed to execute batch job:", err)
	}
}

func TestMAUCountWatcher(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newBatchServiceClient(t)
	defer client.Close()
	_, err := client.ExecuteBatchJob(ctx, &proto.BatchJobRequest{
		Job: proto.BatchJob_MauCountWatcher,
	})
	if err != nil {
		t.Fatal("failed to execute batch job:", err)
	}
}

func TestDatetimeWatcher(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newBatchServiceClient(t)
	defer client.Close()
	_, err := client.ExecuteBatchJob(ctx, &proto.BatchJobRequest{
		Job: proto.BatchJob_DatetimeWatcher,
	})
	if err != nil {
		t.Fatal("failed to execute batch job:", err)
	}
}

func TestCountWatcher(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newBatchServiceClient(t)
	defer client.Close()
	_, err := client.ExecuteBatchJob(ctx, &proto.BatchJobRequest{
		Job: proto.BatchJob_EventCountWatcher,
	})
	if err != nil {
		t.Fatal("failed to execute batch job:", err)
	}
}

func newBatchServiceClient(t *testing.T) batchclient.Client {
	t.Helper()
	credentials, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("failed to create credentials:", err)
	}
	client, err := batchclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(credentials),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("failed to create batch client:", err)
	}
	return client
}
