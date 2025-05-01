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

package tag

import (
	"context"
	"flag"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	tagclient "github.com/bucketeer-io/bucketeer/pkg/tag/client"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	tagproto "github.com/bucketeer-io/bucketeer/proto/tag"
)

const (
	prefixID = "e2e-test"
	timeout  = 60 * time.Second
)

var (
	// FIXME: To avoid compiling the test many times, webGatewayAddr, webGatewayPort & apiKey has been also added here to prevent from getting:
	// "flag provided but not defined" error during the test. These 3 are being use in the Gateway test
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

func TestUpsertAndListTag(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newTagClient(t)
	testTags := []string{
		newTagName(t),
		newTagName(t),
		newTagName(t),
	}
	createTags(t, client, testTags, tagproto.Tag_FEATURE_FLAG)
	actual := listTags(ctx, t, client)
	// Check if the created tags are in the response
	tags := findTags(actual, testTags)
	if len(tags) != len(testTags) {
		t.Fatalf("Different sizes. Expected: %d, Actual: %d", len(testTags), len(tags))
	}
	// Wait a few seconds before upserting the same tag.
	// Otherwise, the test could fail because it could finish in less than 1 second,
	// not updating the `updateAt` correctly.
	time.Sleep(5 * time.Second)
	// Upsert tag index 1
	targetTag := tags[1]
	createTag(t, client, targetTag.Name, tagproto.Tag_FEATURE_FLAG)
	actual = listTags(ctx, t, client)
	tagUpsert := findTags(actual, []string{targetTag.Name})
	if tagUpsert == nil {
		t.Fatalf("Upserted tag wasn't found in the response. Expected: %v\n Response: %v",
			targetTag, actual)
	}
	// Check if the create time is equal
	if targetTag.CreatedAt != tagUpsert[0].CreatedAt {
		t.Fatalf("Different create time. Expected: %v\n, Actual: %v",
			targetTag, tagUpsert[0])
	}
	// Check if the update time has changed
	if targetTag.UpdatedAt == tagUpsert[0].UpdatedAt {
		t.Fatalf("The tag update time didn't change. Expected: %v\n, Actual: %v",
			targetTag, tagUpsert[0])
	}
}

func TestDeleteTag(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newTagClient(t)
	// Create tag
	createReq := &tagproto.CreateTagRequest{
		Name:          newTagName(t),
		EnvironmentId: *environmentID,
		EntityType:    tagproto.Tag_FEATURE_FLAG,
	}
	resp, err := client.CreateTag(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create tag. Error %v", err)
	}
	// Delete tag
	req := &tagproto.DeleteTagRequest{
		Id:            resp.Tag.Id,
		EnvironmentId: *environmentID,
	}
	defer cancel()
	if _, err := client.DeleteTag(ctx, req); err != nil {
		t.Fatalf("Failed to delete tag. Error: %v", err)
	}
	// List the tags
	tags := listTags(ctx, t, client)
	target := findTags(tags, []string{resp.Tag.Name})
	// Check if it has been deleted
	if len(target) != 0 {
		t.Fatalf("The tag hasn't deleted. Tag: %v", target)
	}
}

func listTags(ctx context.Context, t *testing.T, client tagclient.Client) []*tagproto.Tag {
	t.Helper()
	resp, err := client.ListTags(ctx, &tagproto.ListTagsRequest{
		PageSize:      0,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal("Failed to list tags", err)
	}
	return resp.Tags
}

func newTagClient(t *testing.T) tagclient.Client {
	t.Helper()
	creds, err := client.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	tagClient, err := tagclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(10*time.Second),
		client.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create tag client:", err)
	}
	return tagClient
}

func newTagName(t *testing.T) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-tag-%s", prefixID, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-tag-%s", prefixID, newUUID(t))
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func createTags(
	t *testing.T,
	client tagclient.Client,
	tags []string,
	entityType tagproto.Tag_EntityType,
) {
	t.Helper()
	for _, tag := range tags {
		createTag(t, client, tag, entityType)
	}
}

func createTag(
	t *testing.T,
	client tagclient.Client,
	tag string,
	entityType tagproto.Tag_EntityType,
) {
	t.Helper()
	createReq := &tagproto.CreateTagRequest{
		Name:          tag,
		EnvironmentId: *environmentID,
		EntityType:    entityType,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.CreateTag(ctx, createReq); err != nil {
		t.Fatal(err)
	}
}

func findTags(tags []*tagproto.Tag, targetNames []string) []*tagproto.Tag {
	var result []*tagproto.Tag
	for _, tag := range tags {
		if exist := slices.Contains(targetNames, tag.Name); !exist {
			continue
		}
		result = append(result, tag)
	}
	return result
}
