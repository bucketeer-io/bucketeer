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
	"crypto/rand"
	"encoding/base32"
	"flag"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

const (
	defaultOrganizationID   = "e2e"
	e2eAccountAddressPrefix = "e2e-test"
	timeout                 = 60 * time.Second
	firstName               = "first-name"
	lastName                = "last-name"
	language                = "language"
)

var (
	webGatewayAddr                 = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort                 = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert                 = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath                     = flag.String("api-key", "", "Client SDK API key for api-gateway")
	apiKeyServerPath               = flag.String("api-key-server", "", "Server SDK API key for api-gateway")
	gatewayAddr                    = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort                    = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert                    = flag.String("gateway-cert", "", "Gateway crt file")
	sysAdminAccessTokenPath        = flag.String("sys-admin-access-token", "", "System admin access token path")
	orgOwnerDefaultAccessTokenPath = flag.String("org-owner-default-access-token", "", "Organization admin access token path")
	orgOwnerE2EAccessTokenPath     = flag.String("org-owner-e2e-access-token", "", "Organization admin (e2e org) access token path")
	envEditorAccessTokenPath       = flag.String("env-editor-access-token", "", "Environment editor access token path")
	envViewerAccessTokenPath       = flag.String("env-viewer-access-token", "", "Environment viewer access token path")
	environmentID                  = flag.String("environment-id", "", "Environment id")
	organizationID                 = flag.String("organization-id", "", "Organization ID")
	testID                         = flag.String("test-id", "", "test ID")
)

func TestGetAccount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerE2EAccessTokenPath)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId:   defaultOrganizationID,
		Name:             name,
		Email:            email,
		FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
		LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
		Language:         language,
		Tags:             []string{"tag"},
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{
				Role:          accountproto.AccountV2_Role_Environment_VIEWER,
				EnvironmentId: "test",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Account.Email != email {
		t.Fatalf("different email, expected: %v, actual: %v", email, resp.Account.Email)
	}
	if resp.Account.OrganizationId != defaultOrganizationID {
		t.Fatalf("different organization id, expected: %v, actual: %v", defaultOrganizationID, resp.Account.OrganizationId)
	}
}

func TestCreateAccountV2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerE2EAccessTokenPath)
	defer c.Close()
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	firstName := fmt.Sprintf("first-name-%v", time.Now().Unix())
	lastName := fmt.Sprintf("last-name-%v", time.Now().Unix())
	avatarURL := fmt.Sprintf("https://example.com/avatar-%v.png", time.Now().Unix())
	tags := []string{fmt.Sprintf("tag-%d", time.Now().Unix())}
	envRole := &accountproto.AccountV2_EnvironmentRole{
		Role:          accountproto.AccountV2_Role_Environment_VIEWER,
		EnvironmentId: "test",
	}
	resp, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId:   defaultOrganizationID,
		Name:             name,
		Email:            email,
		FirstName:        firstName,
		LastName:         lastName,
		Language:         language,
		AvatarImageUrl:   avatarURL,
		Tags:             tags,
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{envRole},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Account.Email != email {
		t.Fatalf("different email, expected: %v, actual: %v", email, resp.Account.Email)
	}
	if resp.Account.OrganizationId != defaultOrganizationID {
		t.Fatalf("different organization id, expected: %v, actual: %v", defaultOrganizationID, resp.Account.OrganizationId)
	}
	if resp.Account.Name != name {
		t.Fatalf("different name, expected: %v, actual: %v", name, resp.Account.Name)
	}
	if resp.Account.FirstName != firstName {
		t.Fatalf("different first name, expected: %v, actual: %v", firstName, resp.Account.FirstName)
	}
	if resp.Account.LastName != lastName {
		t.Fatalf("different last name, expected: %v, actual: %v", lastName, resp.Account.LastName)
	}
	if resp.Account.Language != language {
		t.Fatalf("different language, expected: %v, actual: %v", language, resp.Account.Language)
	}
	if resp.Account.AvatarImageUrl != avatarURL {
		t.Fatalf("different avatar url, expected: %v, actual: %v", avatarURL, resp.Account.AvatarImageUrl)
	}
	if !reflect.DeepEqual(resp.Account.Tags, tags) {
		t.Fatalf("different tags, expected: %v, actual: %v", tags, resp.Account.Tags)
	}
	if resp.Account.OrganizationRole != accountproto.AccountV2_Role_Organization_MEMBER {
		t.Fatalf("different organization role, expected: %v, actual: %v",
			accountproto.AccountV2_Role_Organization_MEMBER,
			resp.Account.OrganizationRole,
		)
	}
	if resp.Account.EnvironmentRoles[0].Role != envRole.Role {
		t.Fatalf("different organization role, expected: %v, actual: %v",
			envRole.Role,
			resp.Account.EnvironmentRoles[0].Role,
		)
	}
	if resp.Account.EnvironmentRoles[0].EnvironmentId != envRole.EnvironmentId {
		t.Fatalf("different environment id, expected: %s, actual: %s",
			envRole.EnvironmentId,
			resp.Account.EnvironmentRoles[0].EnvironmentId,
		)
	}
}

func TestCreateAccountV2Admin(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerE2EAccessTokenPath)
	defer c.Close()
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	firstName := fmt.Sprintf("first-name-%v", time.Now().Unix())
	lastName := fmt.Sprintf("last-name-%v", time.Now().Unix())
	avatarURL := fmt.Sprintf("https://example.com/avatar-%v.png", time.Now().Unix())
	tags := []string{fmt.Sprintf("tag-%d", time.Now().Unix())}
	envRole := &accountproto.AccountV2_EnvironmentRole{
		Role:          accountproto.AccountV2_Role_Environment_VIEWER,
		EnvironmentId: "test",
	}
	resp, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId:   defaultOrganizationID,
		Name:             name,
		Email:            email,
		FirstName:        firstName,
		LastName:         lastName,
		Language:         language,
		AvatarImageUrl:   avatarURL,
		Tags:             tags,
		OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{envRole},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Account.Email != email {
		t.Fatalf("different email, expected: %v, actual: %v", email, resp.Account.Email)
	}
	if resp.Account.OrganizationRole != accountproto.AccountV2_Role_Organization_ADMIN {
		t.Fatalf("different organization role, expected: %v, actual: %v",
			accountproto.AccountV2_Role_Organization_ADMIN,
			resp.Account.OrganizationRole,
		)
	}
	if len(resp.Account.EnvironmentRoles) != 0 {
		t.Fatalf("env roles length should be zero, actual: %d", len(resp.Account.EnvironmentRoles))
	}
}

func TestListAccounts(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerE2EAccessTokenPath)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId:   defaultOrganizationID,
		Name:             name,
		Email:            email,
		FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
		LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
		Language:         language,
		Tags:             []string{"tag"},
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{
				Role:          accountproto.AccountV2_Role_Environment_VIEWER,
				EnvironmentId: "test",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	pageSize := int64(1)
	resp, err := c.ListAccountsV2(ctx, &accountproto.ListAccountsV2Request{
		OrganizationId: defaultOrganizationID,
		PageSize:       pageSize,
	})
	if err != nil {
		t.Fatal(err)
	}
	responseSize := int64(len(resp.Accounts))
	if responseSize != pageSize {
		t.Fatalf("different sizes, expected: %d actual: %d", pageSize, responseSize)
	}
}

func TestUpdateAccountThenDeleteAccount(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerE2EAccessTokenPath)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId:   defaultOrganizationID,
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
	})
	if err != nil {
		t.Fatal(err)
	}
	newName := fmt.Sprintf("name-%v", time.Now().Unix())
	newFirstName := fmt.Sprintf("first-name-%v", time.Now().Unix())
	newLastName := fmt.Sprintf("last-name-%v", time.Now().Unix())
	newAvatarURL := fmt.Sprintf("https://example.com/avatar-%v.png", time.Now().Unix())
	_, err = c.UpdateAccountV2(ctx, &accountproto.UpdateAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		Name:           wrapperspb.String(newName),
		FirstName:      wrapperspb.String(newFirstName),
		LastName:       wrapperspb.String(newLastName),
		AvatarImageUrl: wrapperspb.String(newAvatarURL),
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{
				EnvironmentId: "test",
				Role:          accountproto.AccountV2_Role_Environment_EDITOR,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getResp.Account.Email != email {
		t.Fatalf("different email, expected: %v, actual: %v", email, getResp.Account.Email)
	}
	if getResp.Account.OrganizationId != defaultOrganizationID {
		t.Fatalf("different organization id, expected: %v, actual: %v", defaultOrganizationID, getResp.Account.OrganizationId)
	}
	if getResp.Account.Name != newName {
		t.Fatalf("different name, expected: %v, actual: %v", newName, getResp.Account.Name)
	}
	if getResp.Account.FirstName != newFirstName {
		t.Fatalf("different first name, expected: %v, actual: %v", newFirstName, getResp.Account.FirstName)
	}
	if getResp.Account.LastName != newLastName {
		t.Fatalf("different last name, expected: %v, actual: %v", newLastName, getResp.Account.LastName)
	}
	if getResp.Account.AvatarImageUrl != newAvatarURL {
		t.Fatalf("different avatar url, expected: %v, actual: %v", newAvatarURL, getResp.Account.AvatarImageUrl)
	}

	// disable then enable account
	_, err = c.DisableAccountV2(ctx, &accountproto.DisableAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp1, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !getResp1.Account.Disabled {
		t.Fatalf("different enabled, expected: %v, actual: %v", true, getResp1.Account.Disabled)
	}

	_, err = c.EnableAccountV2(ctx, &accountproto.EnableAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp2, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getResp2.Account.Disabled {
		t.Fatalf("different enabled, expected: %v, actual: %v", false, getResp2.Account.Disabled)
	}

	// delete account
	_, err = c.DeleteAccountV2(ctx, &accountproto.DeleteAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err == nil {
		t.Fatal("account is not deleted")
	}
}

func TestCreateUpdateAccountTeams(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerE2EAccessTokenPath)
	defer c.Close()
	// create account
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId:   defaultOrganizationID,
		Name:             name,
		Email:            email,
		FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
		LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
		Language:         language,
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		Teams:            []string{"team-id-1", "team-id-2"},
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{
				EnvironmentId: "test",
				Role:          accountproto.AccountV2_Role_Environment_VIEWER,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// update account team
	_, err = c.UpdateAccountV2(ctx, &accountproto.UpdateAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		TeamChanges: []*accountproto.TeamChange{
			{
				ChangeType: accountproto.ChangeType_CREATE,
				Team:       "team-id-3",
			},
			{
				ChangeType: accountproto.ChangeType_DELETE,
				Team:       "team-id-2",
			},
		},
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{
				EnvironmentId: "test",
				Role:          accountproto.AccountV2_Role_Environment_VIEWER,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedTeams := []string{"team-id-1", "team-id-3"}
	if !reflect.DeepEqual(getResp.Account.Teams, expectedTeams) {
		t.Fatalf("different teams, expected: %v, actual: %v", expectedTeams, getResp.Account.Teams)
	}

	// delete account
	_, err = c.DeleteAccountV2(ctx, &accountproto.DeleteAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnableAndDisableAccount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerE2EAccessTokenPath)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId:   defaultOrganizationID,
		Name:             name,
		Email:            email,
		FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
		LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
		Language:         language,
		Tags:             []string{"tag"},
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{
				Role:          accountproto.AccountV2_Role_Environment_VIEWER,
				EnvironmentId: "test",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.UpdateAccountV2(ctx, &accountproto.UpdateAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		Disabled:       wrapperspb.Bool(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp1, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !getResp1.Account.Disabled {
		t.Fatalf("different enabled, expected: %v, actual: %v", true, getResp1.Account.Disabled)
	}

	_, err = c.UpdateAccountV2(ctx, &accountproto.UpdateAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		Disabled:       wrapperspb.Bool(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp2, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getResp2.Account.Disabled {
		t.Fatalf("different enabled, expected: %v, actual: %v", false, getResp2.Account.Disabled)
	}
}

func TestDeleteAccount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerE2EAccessTokenPath)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId:   defaultOrganizationID,
		Name:             name,
		Email:            email,
		FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
		LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
		Language:         language,
		Tags:             []string{"tag"},
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{
				Role:          accountproto.AccountV2_Role_Environment_VIEWER,
				EnvironmentId: "test",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.DeleteAccountV2(ctx, &accountproto.DeleteAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err == nil {
		t.Fatal("account is not deleted")
	}
}

// Search filters are personal settings: the account service only lets callers
// manage the filters of their own account, in the organization that owns the
// environment they are authorized for. The tests below therefore operate on the
// account behind the org owner token and clean up the filters they create.
func TestCreateSearchFilter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerDefaultAccessTokenPath)
	defer c.Close()
	email := callerEmail(t, ctx, c)

	requestSearchFilter := &accountproto.SearchFilter{
		Name:             fmt.Sprintf("name-%s", uniqueSuffix()),
		Query:            "query",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		DefaultFilter:    false,
		EnvironmentId:    *environmentID,
	}
	created, cleanup := createOwnSearchFilter(t, ctx, c, email, requestSearchFilter)
	defer cleanup()

	if created.Query != requestSearchFilter.Query {
		t.Fatalf("different query of filters, expected: %v, actual: %v", requestSearchFilter.Query, created.Query)
	}
	if created.FilterTargetType != requestSearchFilter.FilterTargetType {
		t.Fatalf("different filter target type of filters, expected: %v, actual: %v", requestSearchFilter.FilterTargetType, created.FilterTargetType)
	}
	if created.DefaultFilter != requestSearchFilter.DefaultFilter {
		t.Fatalf("different default filter of filters, expected: %v, actual: %v", requestSearchFilter.DefaultFilter, created.DefaultFilter)
	}
	if created.EnvironmentId != requestSearchFilter.EnvironmentId {
		t.Fatalf("different environment id of filters, expected: %v, actual: %v", requestSearchFilter.EnvironmentId, created.EnvironmentId)
	}
}

func TestUpdateSearchFilter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerDefaultAccessTokenPath)
	defer c.Close()
	email := callerEmail(t, ctx, c)

	initialSearchFilter := &accountproto.SearchFilter{
		Name:             fmt.Sprintf("name-%s", uniqueSuffix()),
		Query:            "query",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		DefaultFilter:    false,
		EnvironmentId:    *environmentID,
	}
	created, cleanup := createOwnSearchFilter(t, ctx, c, email, initialSearchFilter)
	defer cleanup()

	requestSearchFilter := &accountproto.SearchFilter{
		Name:             fmt.Sprintf("new-name-%s", uniqueSuffix()),
		Query:            "new-query",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		DefaultFilter:    true,
		EnvironmentId:    *environmentID,
	}
	_, err := c.UpdateSearchFilter(ctx, &accountproto.UpdateSearchFilterRequest{
		Email:          email,
		OrganizationId: *organizationID,
		EnvironmentId:  *environmentID,
		ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
			Id:   created.Id,
			Name: requestSearchFilter.Name,
		},
		ChangeQueryCommand: &accountproto.ChangeSearchFilterQueryCommand{
			Id:    created.Id,
			Query: requestSearchFilter.Query,
		},
		ChangeDefaultFilterCommand: &accountproto.ChangeDefaultSearchFilterCommand{
			Id:            created.Id,
			DefaultFilter: requestSearchFilter.DefaultFilter,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	updateSearchFilter := findSearchFilterByID(ownSearchFilters(t, ctx, c, email), created.Id)
	if updateSearchFilter == nil {
		t.Fatalf("search filter %s not found after update", created.Id)
	}
	if updateSearchFilter.Name != requestSearchFilter.Name {
		t.Fatalf("different name of filters, expected: %v, actual: %v", requestSearchFilter.Name, updateSearchFilter.Name)
	}
	if updateSearchFilter.Query != requestSearchFilter.Query {
		t.Fatalf("different query of filters, expected: %v, actual: %v", requestSearchFilter.Query, updateSearchFilter.Query)
	}
	if updateSearchFilter.FilterTargetType != requestSearchFilter.FilterTargetType {
		t.Fatalf("different filter target type of filters, expected: %v, actual: %v", requestSearchFilter.FilterTargetType, updateSearchFilter.FilterTargetType)
	}
	if updateSearchFilter.DefaultFilter != requestSearchFilter.DefaultFilter {
		t.Fatalf("different default filter of filters, expected: %v, actual: %v", requestSearchFilter.DefaultFilter, updateSearchFilter.DefaultFilter)
	}
}

func TestDeleteSearchFilter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t, orgOwnerDefaultAccessTokenPath)
	defer c.Close()
	email := callerEmail(t, ctx, c)

	first, cleanupFirst := createOwnSearchFilter(t, ctx, c, email, &accountproto.SearchFilter{
		Name:             fmt.Sprintf("name1-%s", uniqueSuffix()),
		Query:            "query1",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		DefaultFilter:    false,
		EnvironmentId:    *environmentID,
	})
	defer cleanupFirst()
	second, cleanupSecond := createOwnSearchFilter(t, ctx, c, email, &accountproto.SearchFilter{
		Name:             fmt.Sprintf("name2-%s", uniqueSuffix()),
		Query:            "query2",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		DefaultFilter:    false,
		EnvironmentId:    *environmentID,
	})
	defer cleanupSecond()
	before := ownSearchFilters(t, ctx, c, email)

	_, err := c.DeleteSearchFilter(ctx, &accountproto.DeleteSearchFilterRequest{
		Email:          email,
		OrganizationId: *organizationID,
		EnvironmentId:  *environmentID,
		Command: &accountproto.DeleteSearchFilterCommand{
			Id: first.Id,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	after := ownSearchFilters(t, ctx, c, email)
	if len(after) != len(before)-1 {
		t.Fatalf("different count of filters, expected: %d, actual: %d", len(before)-1, len(after))
	}
	if findSearchFilterByID(after, first.Id) != nil {
		t.Fatalf("search filter is not deleted")
	}
	if findSearchFilterByID(after, second.Id) == nil {
		t.Fatalf("search filter %s was removed unexpectedly", second.Id)
	}
}

// callerEmail returns the email of the account behind the client's access token.
func callerEmail(t *testing.T, ctx context.Context, c accountclient.Client) string {
	t.Helper()
	me, err := c.GetMe(ctx, &accountproto.GetMeRequest{OrganizationId: *organizationID})
	if err != nil {
		t.Fatal(err)
	}
	return me.Account.Email
}

// ownSearchFilters returns the search filters stored on the caller's account in
// the organization that owns the e2e environment.
func ownSearchFilters(t *testing.T, ctx context.Context, c accountclient.Client, email string) []*accountproto.SearchFilter {
	t.Helper()
	resp, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: *organizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.Account.SearchFilters
}

func findSearchFilterByName(filters []*accountproto.SearchFilter, name string) *accountproto.SearchFilter {
	for _, f := range filters {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func findSearchFilterByID(filters []*accountproto.SearchFilter, id string) *accountproto.SearchFilter {
	for _, f := range filters {
		if f.Id == id {
			return f
		}
	}
	return nil
}

// createOwnSearchFilter creates a search filter on the caller's own account and
// returns the stored filter together with a function that deletes it again, so
// the shared e2e account does not accumulate filters across runs.
func createOwnSearchFilter(
	t *testing.T,
	ctx context.Context,
	c accountclient.Client,
	email string,
	filter *accountproto.SearchFilter,
) (*accountproto.SearchFilter, func()) {
	t.Helper()
	before := ownSearchFilters(t, ctx, c, email)
	if findSearchFilterByName(before, filter.Name) != nil {
		t.Fatalf("search filter %q already exists", filter.Name)
	}
	_, err := c.CreateSearchFilter(ctx, &accountproto.CreateSearchFilterRequest{
		Email:          email,
		OrganizationId: *organizationID,
		EnvironmentId:  *environmentID,
		Command: &accountproto.CreateSearchFilterCommand{
			Name:             filter.Name,
			Query:            filter.Query,
			FilterTargetType: filter.FilterTargetType,
			EnvironmentId:    filter.EnvironmentId,
			DefaultFilter:    filter.DefaultFilter,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	after := ownSearchFilters(t, ctx, c, email)
	if len(after) != len(before)+1 {
		t.Fatalf("different count of filters, expected: %d, actual: %d", len(before)+1, len(after))
	}
	created := findSearchFilterByName(after, filter.Name)
	if created == nil {
		t.Fatalf("search filter %q was not created", filter.Name)
	}
	cleanup := func() {
		_, _ = c.DeleteSearchFilter(ctx, &accountproto.DeleteSearchFilterRequest{
			Email:          email,
			OrganizationId: *organizationID,
			EnvironmentId:  *environmentID,
			Command:        &accountproto.DeleteSearchFilterCommand{Id: created.Id},
		})
	}
	return created, cleanup
}

func newAccountClient(t *testing.T, token *string) accountclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*token)
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
