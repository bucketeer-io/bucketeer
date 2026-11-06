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

package account

import (
	"context"
	"fmt"
	"sync"
	"testing"

	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

const (
	// The project and environment of the e2e organization that the account
	// tests grant environment roles in. They are created on first use and
	// kept, so every later run finds them by URL code.
	e2eOrgProjectURLCode     = "e2e-account-tests"
	e2eOrgEnvironmentURLCode = "e2e-account-tests"
	listPageSize             = 100
)

var (
	e2eOrgEnvironmentOnce sync.Once
	e2eOrgEnvironment     string
	e2eOrgEnvironmentErr  error
)

// e2eOrgEnvironmentID returns the ID of an environment of the e2e
// organization. An account's environment roles must refer to environments of
// the account's organization, and the e2e organization has none seeded, so
// the tests create a project and an environment there once and reuse them.
func e2eOrgEnvironmentID(t *testing.T) string {
	t.Helper()
	e2eOrgEnvironmentOnce.Do(func() {
		e2eOrgEnvironment, e2eOrgEnvironmentErr = ensureE2EOrgEnvironment(t)
	})
	if e2eOrgEnvironmentErr != nil {
		t.Fatalf("Failed to prepare an environment in the %s organization: %v",
			defaultOrganizationID, e2eOrgEnvironmentErr)
	}
	return e2eOrgEnvironment
}

func ensureE2EOrgEnvironment(t *testing.T) (string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// The owner of the e2e organization may list and create its projects and
	// environments, so the setup needs no other credential than the tests.
	c := newEnvironmentClientWithToken(t, *orgOwnerE2EAccessTokenPath)
	defer c.Close()

	projectID, err := findProjectID(ctx, c)
	if err != nil {
		return "", err
	}
	if projectID == "" {
		created, err := c.CreateProject(ctx, &environmentproto.CreateProjectRequest{
			OrganizationId: defaultOrganizationID,
			Name:           "E2E account tests",
			UrlCode:        e2eOrgProjectURLCode,
			Description:    "Project of the account e2e tests",
		})
		if err != nil {
			return "", fmt.Errorf("create project: %w", err)
		}
		projectID = created.Project.Id
	}

	environmentID, err := findEnvironmentID(ctx, c, projectID)
	if err != nil {
		return "", err
	}
	if environmentID == "" {
		created, err := c.CreateEnvironmentV2(ctx, &environmentproto.CreateEnvironmentV2Request{
			Name:        "E2E account tests",
			UrlCode:     e2eOrgEnvironmentURLCode,
			Description: "Environment of the account e2e tests",
			ProjectId:   projectID,
		})
		if err != nil {
			return "", fmt.Errorf("create environment: %w", err)
		}
		environmentID = created.Environment.Id
	}
	return environmentID, nil
}

// findProjectID returns the ID of the test project of the e2e organization,
// or an empty string when it does not exist yet.
func findProjectID(ctx context.Context, c environmentclient.Client) (string, error) {
	resp, err := c.ListProjectsV2(ctx, &environmentproto.ListProjectsV2Request{
		OrganizationId: defaultOrganizationID,
		SearchKeyword:  e2eOrgProjectURLCode,
		PageSize:       listPageSize,
	})
	if err != nil {
		return "", fmt.Errorf("list projects: %w", err)
	}
	for _, p := range resp.Projects {
		if p.UrlCode == e2eOrgProjectURLCode {
			return p.Id, nil
		}
	}
	return "", nil
}

// findEnvironmentID returns the ID of the test environment of the project, or
// an empty string when it does not exist yet.
func findEnvironmentID(ctx context.Context, c environmentclient.Client, projectID string) (string, error) {
	resp, err := c.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{
		OrganizationId: defaultOrganizationID,
		ProjectId:      projectID,
		PageSize:       listPageSize,
	})
	if err != nil {
		return "", fmt.Errorf("list environments: %w", err)
	}
	for _, e := range resp.Environments {
		if e.UrlCode == e2eOrgEnvironmentURLCode {
			return e.Id, nil
		}
	}
	return "", nil
}
