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

package team

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"flag"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"
	"time"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	teamclient "github.com/bucketeer-io/bucketeer/v2/pkg/team/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	teamproto "github.com/bucketeer-io/bucketeer/v2/proto/team"
)

const (
	prefixID                = "e2e-test"
	timeout                 = 60 * time.Second
	firstName               = "first-name"
	lastName                = "last-name"
	language                = "language"
	e2eAccountAddressPrefix = "e2e-test"
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

func TestUpsertAndListTeam(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := newTeamClient(t)
	testTeams := []string{
		newTeamName(t),
		newTeamName(t),
		newTeamName(t),
	}
	createTeams(t, client, testTeams)
	actual := listTeams(ctx, t, client)
	// Check if the created teams are in the response
	teams := findTeams(actual, testTeams)
	if len(teams) != len(testTeams) {
		t.Fatalf("Different sizes. Expected: %d, Actual: %d", len(testTeams), len(teams))
	}
	// Wait a few seconds before upserting the same team.
	// Otherwise, the test could fail because it could finish in less than 1 second,
	// not updating the `updateAt` correctly.
	time.Sleep(5 * time.Second)
	// Upsert team index 1
	targetTeam := teams[1]
	createTeam(t, client, targetTeam.Name)
	actual = listTeams(ctx, t, client)
	teamUpsert := findTeams(actual, []string{targetTeam.Name})
	if teamUpsert == nil {
		t.Fatalf("Upserted team wasn't found in the response. Expected: %v\n Response: %v",
			targetTeam, actual)
	}
	// Check if the create time is equal
	if targetTeam.CreatedAt != teamUpsert[0].CreatedAt {
		t.Fatalf("Different create time. Expected: %v\n, Actual: %v",
			targetTeam, teamUpsert[0])
	}
	// Check if the update time has changed
	if targetTeam.UpdatedAt == teamUpsert[0].UpdatedAt {
		t.Fatalf("The team update time didn't change. Expected: %v\n, Actual: %v",
			targetTeam, teamUpsert[0])
	}
}

func TestDeleteTeam(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newTeamClient(t)
	// Create team
	createReq := &teamproto.CreateTeamRequest{
		Name:           newTeamName(t),
		OrganizationId: *organizationID,
	}
	resp, err := client.CreateTeam(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create team. Error %v", err)
	}
	// Delete team
	req := &teamproto.DeleteTeamRequest{
		Id:             resp.Team.Id,
		OrganizationId: *organizationID,
	}
	defer cancel()
	if _, err := client.DeleteTeam(ctx, req); err != nil {
		t.Fatalf("Failed to delete team. Error: %v", err)
	}
	// List the teams
	teams := listTeams(ctx, t, client)
	target := findTeams(teams, []string{resp.Team.Name})
	// Check if it has been deleted
	if len(target) != 0 {
		t.Fatalf("The team hasn't deleted. Team: %v", target)
	}
}

func TestFailedDeleteTeam(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	teamClient := newTeamClient(t)
	defer teamClient.Close()
	accountClient := newAccountClient(t)
	defer accountClient.Close()

	// create account with team
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := accountClient.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId:   *organizationID,
		Name:             name,
		Email:            email,
		FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
		LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
		Language:         language,
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{
				EnvironmentId: "test",
				Role:          accountproto.AccountV2_Role_Environment_VIEWER,
			},
		},
		Teams: []string{"team1"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// get team id
	teams := listTeams(ctx, t, teamClient)
	var teamID string
	for _, team := range teams {
		if team.Name == "team1" {
			teamID = team.Id
			break
		}
	}
	req := &teamproto.DeleteTeamRequest{
		Id:             teamID,
		OrganizationId: *organizationID,
	}
	defer cancel()
	if _, err := teamClient.DeleteTeam(ctx, req); err == nil {
		t.Fatal("Expected error when deleting team with existing account, but got nil")
	}
}

func newAccountClient(t *testing.T) accountclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := accountclient.NewClient(
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

func randomString() string {
	b := make([]byte, 8)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	return strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")
}

func listTeams(ctx context.Context, t *testing.T, client teamclient.Client) []*teamproto.Team {
	t.Helper()
	resp, err := client.ListTeams(ctx, &teamproto.ListTeamsRequest{
		PageSize:       0,
		OrganizationId: *organizationID,
	})
	if err != nil {
		t.Fatal("Failed to list teams", err)
	}
	return resp.Teams
}

func newTeamName(t *testing.T) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-team-%s", prefixID, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-team-%s", prefixID, newUUID(t))
}

func createTeams(
	t *testing.T,
	client teamclient.Client,
	teams []string,
) {
	t.Helper()
	for _, team := range teams {
		createTeam(t, client, team)
	}
}

func createTeam(
	t *testing.T,
	client teamclient.Client,
	team string,
) {
	t.Helper()
	createReq := &teamproto.CreateTeamRequest{
		Name:           team,
		OrganizationId: *organizationID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.CreateTeam(ctx, createReq); err != nil {
		t.Fatal(err)
	}
}

func findTeams(teams []*teamproto.Team, targetNames []string) []*teamproto.Team {
	var result []*teamproto.Team
	for _, team := range teams {
		if exist := slices.Contains(targetNames, team.Name); !exist {
			continue
		}
		result = append(result, team)
	}
	return result
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func newTeamClient(t *testing.T) teamclient.Client {
	t.Helper()
	creds, err := client.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	teamClient, err := teamclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(10*time.Second),
		client.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create team client:", err)
	}
	return teamClient
}
