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
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

const (
	defaultProjectID = "default"
)

func TestGetProject(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	id := defaultProjectID
	resp, err := c.GetProject(ctx, &environmentproto.GetProjectRequest{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Project.Id != id {
		t.Fatalf("different ids, expected: %v, actual: %v", id, resp.Project.Id)
	}
}

func TestListProjects(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	pageSize := int64(1)
	resp, err := c.ListProjects(ctx, &environmentproto.ListProjectsRequest{PageSize: pageSize})
	if err != nil {
		t.Fatal(err)
	}
	responseSize := int64(len(resp.Projects))
	if responseSize != pageSize {
		t.Fatalf("different sizes, expected: %d actual: %d", pageSize, responseSize)
	}
}

func TestListProjectsRequestOrganizations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	pageSize := int64(2)

	orgResp, err := c.ListOrganizations(ctx, &environmentproto.ListOrganizationsRequest{PageSize: pageSize})
	if err != nil {
		t.Fatal(err)
	}
	length := int(math.Min(float64(len(orgResp.Organizations)), float64(pageSize)))
	orgIds := make([]string, length)
	for i := 0; i < length; i++ {
		orgIds[i] = orgResp.Organizations[rand.Intn(len(orgResp.Organizations))].Id
	}

	resp, err := c.ListProjects(ctx, &environmentproto.ListProjectsRequest{PageSize: pageSize, OrganizationIds: orgIds})
	if err != nil {
		t.Fatal(err)
	}
	for _, project := range resp.Projects {
		found := false
		for _, orgId := range orgIds {
			if orgId == project.OrganizationId {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("received a project of a not request organization expected: %v, actual: %s", orgIds, project.OrganizationId)
		}
	}
}

func TestUpdateProject(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	id := defaultProjectID
	newDesc := fmt.Sprintf("Description %v", time.Now().Unix())
	newName := fmt.Sprintf("name-%v", time.Now().Unix())
	_, err := c.UpdateProject(ctx, &environmentproto.UpdateProjectRequest{
		Id:                       id,
		ChangeDescriptionCommand: &environmentproto.ChangeDescriptionProjectCommand{Description: newDesc},
		RenameCommand:            &environmentproto.RenameProjectCommand{Name: newName},
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetProject(ctx, &environmentproto.GetProjectRequest{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if getResp.Project.Id != id {
		t.Fatalf("different ids, expected: %v, actual: %v", id, getResp.Project.Id)
	}
	if getResp.Project.Description != newDesc {
		t.Fatalf("different descriptions, expected: %v, actual: %v", newDesc, getResp.Project.Description)
	}
	if getResp.Project.Name != newName {
		t.Fatalf("different names, expected: %v, actual: %v", newName, getResp.Project.Name)
	}
}
