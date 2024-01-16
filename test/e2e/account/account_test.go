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

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

const (
	defaultOrganizationID   = "e2e"
	e2eAccountAddressPrefix = "e2e-test"
	timeout                 = 10 * time.Second
)

var (
	webGatewayAddr       = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort       = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert       = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath           = flag.String("api-key", "", "Api key path for web gateway")
	gatewayAddr          = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort          = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert          = flag.String("gateway-cert", "", "Gateway crt file")
	serviceTokenPath     = flag.String("service-token", "", "Service token path")
	environmentNamespace = flag.String("environment-namespace", "", "Environment namespace")
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
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
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
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
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
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	newName := fmt.Sprintf("name-%v", time.Now().Unix())
	newAvatarURL := fmt.Sprintf("https://example.com/avatar-%v.png", time.Now().Unix())
	_, err = c.UpdateAccountV2(ctx, &accountproto.UpdateAccountV2Request{
		Email:          email,
		OrganizationId: defaultOrganizationID,
		ChangeNameCommand: &accountproto.ChangeAccountV2NameCommand{
			Name: newName,
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
	if getResp.Account.AvatarImageUrl != newAvatarURL {
		t.Fatalf("different avatar url, expected: %v, actual: %v", newAvatarURL, getResp.Account.AvatarImageUrl)
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
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
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
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
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
