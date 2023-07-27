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

package environment

import (
	"context"
	"fmt"
	"testing"
	"time"

	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestGetEnvironmentV2(t *testing.T) {
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
	// TODO: replace namespace to name after migration to environment-v2 API
	if resp.Environment.Name != *environmentNamespace {
		t.Fatalf("different name, expected: %v, actual: %v", *environmentNamespace, resp.Environment.Name)
	}
}

func TestListEnvironmentsV2ByProject(t *testing.T) {
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
}

func TestListEnvironmentsV2(t *testing.T) {
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
}

func TestUpdateEnvironmentV2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	id := getEnvironmentID(t)
	newDesc := fmt.Sprintf("Description %v", time.Now().Unix())
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
	// TODO: replace namespace to name after migration to environment-v2 API
	if getResp.Environment.Name != *environmentNamespace {
		t.Fatalf("different name, expected: %v, actual: %v", *environmentNamespace, getResp.Environment.Name)
	}
	if getResp.Environment.Description != newDesc {
		t.Fatalf("different descriptions, expected: %v, actual: %v", newDesc, getResp.Environment.Description)
	}
}
