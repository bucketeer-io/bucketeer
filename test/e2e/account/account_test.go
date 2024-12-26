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

package account

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"flag"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
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
	webGatewayAddr       = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort       = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert       = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath           = flag.String("api-key", "", "Client SDK API key for api-gateway")
	apiKeyServerPath     = flag.String("api-key-server", "", "Server SDK API key for api-gateway")
	gatewayAddr          = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort          = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert          = flag.String("gateway-cert", "", "Gateway crt file")
	serviceTokenPath     = flag.String("service-token", "", "Service token path")
	environmentNamespace = flag.String("environment-namespace", "", "Environment namespace")
	organizationID       = flag.String("organization-id", "", "Organization ID")
	testID               = flag.String("test-id", "", "test ID")
)

func TestGetAccount(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateAccountV2Command{
			Name:             name,
			Email:            email,
			FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
			LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
			Language:         language,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					Role:          accountproto.AccountV2_Role_Environment_VIEWER,
					EnvironmentId: "test",
				},
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

func TestListAccounts(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateAccountV2Command{
			Name:             name,
			Email:            email,
			FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
			LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
			Language:         language,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					Role:          accountproto.AccountV2_Role_Environment_VIEWER,
					EnvironmentId: "test",
				},
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

func TestUpdateAccount(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateAccountV2Command{
			Name:             name,
			Email:            email,
			FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
			LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
			Language:         language,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					Role:          accountproto.AccountV2_Role_Environment_VIEWER,
					EnvironmentId: "test",
				},
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
		ChangeNameCommand: &accountproto.ChangeAccountV2NameCommand{
			Name: newName,
		},
		ChangeFirstNameCommand: &accountproto.ChangeAccountV2FirstNameCommand{
			FirstName: newFirstName,
		},
		ChangeLastNameCommand: &accountproto.ChangeAccountV2LastNameCommand{
			LastName: newLastName,
		},
		ChangeAvatarUrlCommand: &accountproto.ChangeAccountV2AvatarImageUrlCommand{
			AvatarImageUrl: newAvatarURL,
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
}

func TestUpdateAccountThenDeleteAccountNoCommand(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t)
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
		Command:        &accountproto.DeleteAccountV2Command{},
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

func TestEnableAndDisableAccount(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateAccountV2Command{
			Name:             name,
			Email:            email,
			FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
			LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
			Language:         language,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					Role:          accountproto.AccountV2_Role_Environment_VIEWER,
					EnvironmentId: "test",
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.DisableAccountV2(ctx, &accountproto.DisableAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		Command:        &accountproto.DisableAccountV2Command{},
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
		Command:        &accountproto.EnableAccountV2Command{},
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
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateAccountV2Command{
			Name:             name,
			Email:            email,
			FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
			LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
			Language:         language,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					Role:          accountproto.AccountV2_Role_Environment_VIEWER,
					EnvironmentId: "test",
				},
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
		Command:        &accountproto.DeleteAccountV2Command{},
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

func TestCreateSearchFilter(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateAccountV2Command{
			Name:             name,
			Email:            email,
			FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
			LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
			Language:         language,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					Role:          accountproto.AccountV2_Role_Environment_VIEWER,
					EnvironmentId: "test",
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	baseAccount, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if baseAccount.Account.SearchFilters != nil {
		t.Fatal("search filters are not nil")
	}

	requestSearchFilter := &accountproto.SearchFilter{
		Name:             "name",
		Query:            "query",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		DefaultFilter:    false,
		EnvironmentId:    "environment-id",
	}
	_, err = c.CreateSearchFilter(ctx, &accountproto.CreateSearchFilterRequest{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateSearchFilterCommand{
			Name:             requestSearchFilter.Name,
			Query:            requestSearchFilter.Query,
			FilterTargetType: requestSearchFilter.FilterTargetType,
			EnvironmentId:    requestSearchFilter.EnvironmentId,
			DefaultFilter:    requestSearchFilter.DefaultFilter,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	account, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(account.Account.SearchFilters) != 1 {
		t.Fatalf("different count of filters, expected: 1, actual: %v", len(account.Account.SearchFilters))
	}
	if account.Account.SearchFilters[0].Name != requestSearchFilter.Name {
		t.Fatalf("different name of filters, expected: %v, actual: %v", requestSearchFilter.Name, account.Account.SearchFilters[0].Name)
	}
	if account.Account.SearchFilters[0].Query != requestSearchFilter.Query {
		t.Fatalf("different query of filters, expected: %v, actual: %v", requestSearchFilter.Query, account.Account.SearchFilters[0].Query)
	}
	if account.Account.SearchFilters[0].FilterTargetType != requestSearchFilter.FilterTargetType {
		t.Fatalf("different filter target type of filters, expected: %v, actual: %v", requestSearchFilter.FilterTargetType, account.Account.SearchFilters[0].FilterTargetType)
	}
	if account.Account.SearchFilters[0].DefaultFilter != requestSearchFilter.DefaultFilter {
		t.Fatalf("different default filter of filters, expected: %v, actual: %v", requestSearchFilter.DefaultFilter, account.Account.SearchFilters[0].DefaultFilter)
	}
}

func TestUpdateSearchFilter(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateAccountV2Command{
			Name:             name,
			Email:            email,
			FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
			LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
			Language:         language,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					Role:          accountproto.AccountV2_Role_Environment_VIEWER,
					EnvironmentId: "test",
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	baseAccount, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if baseAccount.Account.SearchFilters != nil {
		t.Fatal("search filters are not nil")
	}
	initialSearchFilter := &accountproto.SearchFilter{
		Name:             "name",
		Query:            "query",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		DefaultFilter:    false,
		EnvironmentId:    "environment-id",
	}
	_, err = c.CreateSearchFilter(ctx, &accountproto.CreateSearchFilterRequest{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateSearchFilterCommand{
			Name:             initialSearchFilter.Name,
			Query:            initialSearchFilter.Query,
			FilterTargetType: initialSearchFilter.FilterTargetType,
			EnvironmentId:    initialSearchFilter.EnvironmentId,
			DefaultFilter:    initialSearchFilter.DefaultFilter,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	account, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(account.Account.SearchFilters) != 1 {
		t.Fatalf("different count of filters, expected: 1, actual: %v", len(account.Account.SearchFilters))
	}
	if account.Account.SearchFilters[0].Name != initialSearchFilter.Name {
		t.Fatalf("different name of filters, expected: %v, actual: %v", initialSearchFilter.Name, account.Account.SearchFilters[0].Name)
	}
	if account.Account.SearchFilters[0].Query != initialSearchFilter.Query {
		t.Fatalf("different query of filters, expected: %v, actual: %v", initialSearchFilter.Query, account.Account.SearchFilters[0].Query)
	}
	if account.Account.SearchFilters[0].FilterTargetType != initialSearchFilter.FilterTargetType {
		t.Fatalf("different filter target type of filters, expected: %v, actual: %v", initialSearchFilter.FilterTargetType, account.Account.SearchFilters[0].FilterTargetType)
	}
	if account.Account.SearchFilters[0].DefaultFilter != initialSearchFilter.DefaultFilter {
		t.Fatalf("different default filter of filters, expected: %v, actual: %v", initialSearchFilter.DefaultFilter, account.Account.SearchFilters[0].DefaultFilter)
	}

	requestSearchFilter := &accountproto.SearchFilter{
		Name:             "new-name",
		Query:            "new-query",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		DefaultFilter:    true,
		EnvironmentId:    "environment-id",
	}

	_, err = c.UpdateSearchFilter(ctx, &accountproto.UpdateSearchFilterRequest{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		EnvironmentId:  requestSearchFilter.EnvironmentId,
		ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
			Id:   account.Account.SearchFilters[0].Id,
			Name: requestSearchFilter.Name,
		},
		ChangeQueryCommand: &accountproto.ChangeSearchFilterQueryCommand{
			Id:    account.Account.SearchFilters[0].Id,
			Query: requestSearchFilter.Query,
		},
		ChangeDefaultFilterCommand: &accountproto.ChangeDefaultSearchFilterCommand{
			Id:            account.Account.SearchFilters[0].Id,
			DefaultFilter: requestSearchFilter.DefaultFilter,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	updateAccount, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(updateAccount.Account.SearchFilters) != 1 {
		t.Fatalf("different count of filters, expected: 1, actual: %v", len(account.Account.SearchFilters))
	}
	updateSearchFilter := updateAccount.Account.SearchFilters[0]
	if updateSearchFilter.Name != requestSearchFilter.Name {
		t.Fatalf("different name of filters, expected: %v, actual: %v", initialSearchFilter.Name, requestSearchFilter.Name)
	}
	if updateSearchFilter.Query != requestSearchFilter.Query {
		t.Fatalf("different query of filters, expected: %v, actual: %v", initialSearchFilter.Query, requestSearchFilter.Query)
	}
	if updateSearchFilter.FilterTargetType != requestSearchFilter.FilterTargetType {
		t.Fatalf("different filter target type of filters, expected: %v, actual: %v", initialSearchFilter.FilterTargetType, requestSearchFilter.FilterTargetType)
	}
	if updateSearchFilter.DefaultFilter != requestSearchFilter.DefaultFilter {
		t.Fatalf("different default filter of filters, expected: %v, actual: %v", initialSearchFilter.DefaultFilter, requestSearchFilter.DefaultFilter)
	}
}

func TestDeleteSearchFilter(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAccountClient(t)
	defer c.Close()
	email := fmt.Sprintf("%s-%s-%v-%s@example.com", e2eAccountAddressPrefix, *testID, time.Now().Unix(), randomString())
	name := fmt.Sprintf("name-%v-%v", time.Now().Unix(), randomString())
	_, err := c.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateAccountV2Command{
			Name:             name,
			Email:            email,
			FirstName:        fmt.Sprintf("%s-%v", firstName, time.Now().Unix()),
			LastName:         fmt.Sprintf("%s-%v", lastName, time.Now().Unix()),
			Language:         language,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					Role:          accountproto.AccountV2_Role_Environment_VIEWER,
					EnvironmentId: "test",
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	baseAccount, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if baseAccount.Account.SearchFilters != nil {
		t.Fatal("search filters are not nil")
	}

	initial1SearchFilter := &accountproto.SearchFilter{
		Name:             "name1",
		Query:            "query1",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		DefaultFilter:    false,
		EnvironmentId:    "environment-id",
	}
	_, err = c.CreateSearchFilter(ctx, &accountproto.CreateSearchFilterRequest{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateSearchFilterCommand{
			Name:             initial1SearchFilter.Name,
			Query:            initial1SearchFilter.Query,
			FilterTargetType: initial1SearchFilter.FilterTargetType,
			EnvironmentId:    initial1SearchFilter.EnvironmentId,
			DefaultFilter:    initial1SearchFilter.DefaultFilter,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	initial2SearchFilter := &accountproto.SearchFilter{
		Name:             "name2",
		Query:            "query2",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		DefaultFilter:    false,
		EnvironmentId:    "environment-id",
	}
	_, err = c.CreateSearchFilter(ctx, &accountproto.CreateSearchFilterRequest{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.CreateSearchFilterCommand{
			Name:             initial2SearchFilter.Name,
			Query:            initial2SearchFilter.Query,
			FilterTargetType: initial2SearchFilter.FilterTargetType,
			EnvironmentId:    initial2SearchFilter.EnvironmentId,
			DefaultFilter:    initial2SearchFilter.DefaultFilter,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedAccount, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(updatedAccount.Account.SearchFilters) != 2 {
		t.Fatalf("different count of filters, expected: 2, actual: %v", len(updatedAccount.Account.SearchFilters))
	}

	deleteFilterID := updatedAccount.Account.SearchFilters[0].Id
	_, err = c.DeleteSearchFilter(ctx, &accountproto.DeleteSearchFilterRequest{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		Command: &accountproto.DeleteSearchFilterCommand{
			Id: deleteFilterID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	filterRemovalAccount, err := c.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(filterRemovalAccount.Account.SearchFilters) != 1 {
		t.Fatalf("different count of filters, expected: 1, actual: %v", len(filterRemovalAccount.Account.SearchFilters))
	}
	for _, f := range filterRemovalAccount.Account.SearchFilters {
		if f.Id == deleteFilterID {
			t.Fatalf("search filter is not deleted")
		}
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
