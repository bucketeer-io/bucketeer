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
	"testing"
	"time"

	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestCreateAndListTag(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	cmd := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	actual := listTags(ctx, t, client)
	tags := findTags(actual, cmd.Tags)
	if len(tags) != len(cmd.Tags) {
		t.Fatalf("Different sizes. Expected: %d, Actual: %d", len(cmd.Tags), len(tags))
	}
}

func TestUpdateTag(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	actual := listTags(ctx, t, client)
	tags := findTags(actual, cmd.Tags)
	if len(tags) != len(cmd.Tags) {
		t.Fatalf("Different sizes. Expected: %d, Actual: %d", len(cmd.Tags), len(tags))
	}

	newTag := "tag-1"
	addTag(t, newTag, featureID, client)
	expected := append(cmd.Tags, newTag)
	time.Sleep(time.Second * 3)
	actual = listTags(ctx, t, client)
	tags = findTags(actual, expected)
	if len(tags) != len(expected) {
		t.Fatalf("Different sizes. Expected: %d, Actual: %d", len(expected), len(tags))
	}
}

func findTags(tags []*feature.Tag, targetIDs []string) []*feature.Tag {
	var result []*feature.Tag
	for _, tag := range tags {
		if exist := existTag(targetIDs, tag.Id); !exist {
			continue
		}
		result = append(result, tag)
	}
	return result
}

func existTag(tags []string, target string) bool {
	for _, tag := range tags {
		if tag == target {
			return true
		}
	}
	return false
}

func listTags(ctx context.Context, t *testing.T, client featureclient.Client) []*feature.Tag {
	t.Helper()
	resp, err := client.ListTags(ctx, &feature.ListTagsRequest{
		PageSize:             int64(500),
		EnvironmentNamespace: *environmentNamespace,
	})
	if err != nil {
		t.Fatal("failed to list tags", err)
	}
	return resp.Tags
}

func addTag(t *testing.T, tag string, featureID string, client featureclient.Client) {
	t.Helper()
	addReq := &feature.UpdateFeatureDetailsRequest{
		Id: featureID,
		AddTagCommands: []*feature.AddTagCommand{
			{Tag: tag},
		},
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureDetails(ctx, addReq); err != nil {
		t.Fatal(err)
	}
}
