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

// Permission tests for the four e2e accounts bootstrapped by
// hack/create-e2e-accounts (system admin, org owner, environment editor,
// environment viewer), each authenticating with its own access token.

package account

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestOrganizationOwnerPermissions(t *testing.T) {
	requireAccessToken(t, *orgOwnerDefaultAccessTokenPath, "org owner")

	t.Run("cannot create organization", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newEnvironmentClientWithToken(t, *orgOwnerDefaultAccessTokenPath)
		defer c.Close()
		_, err := c.CreateOrganization(ctx, newCreateOrganizationReq(ownerEmail()))
		requirePermissionDenied(t, err, "org owner CreateOrganization")
	})
}

func TestEnvironmentEditorPermissions(t *testing.T) {
	requireAccessToken(t, *envEditorAccessTokenPath, "environment editor")

	t.Run("can create feature", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newFeatureClientWithToken(t, *envEditorAccessTokenPath)
		defer c.Close()
		if _, err := c.CreateFeature(ctx, newCreateFeatureReq(featureID("editor"))); err != nil {
			t.Fatalf("environment editor should be able to create a feature, but got: %v", err)
		}
	})

	t.Run("cannot create organization", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newEnvironmentClientWithToken(t, *envEditorAccessTokenPath)
		defer c.Close()
		_, err := c.CreateOrganization(ctx, newCreateOrganizationReq(ownerEmail()))
		requirePermissionDenied(t, err, "environment editor CreateOrganization")
	})

	t.Run("cannot add new member", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newAccountClientWithToken(t, *envEditorAccessTokenPath)
		defer c.Close()
		_, err := c.CreateAccountV2(ctx, newCreateAccountReq(*organizationID, memberEmail()))
		requirePermissionDenied(t, err, "environment editor CreateAccountV2")
	})
}

func TestEnvironmentViewerPermissions(t *testing.T) {
	requireAccessToken(t, *envViewerAccessTokenPath, "environment viewer")

	t.Run("can view features", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newFeatureClientWithToken(t, *envViewerAccessTokenPath)
		defer c.Close()
		if _, err := c.ListFeatures(ctx, &featureproto.ListFeaturesRequest{
			PageSize:      1,
			EnvironmentId: *environmentID,
		}); err != nil {
			t.Fatalf("environment viewer should be able to list features, but got: %v", err)
		}
	})

	t.Run("cannot create feature", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newFeatureClientWithToken(t, *envViewerAccessTokenPath)
		defer c.Close()
		_, err := c.CreateFeature(ctx, newCreateFeatureReq(featureID("viewer")))
		requirePermissionDenied(t, err, "environment viewer CreateFeature")
	})

	t.Run("cannot create organization", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newEnvironmentClientWithToken(t, *envViewerAccessTokenPath)
		defer c.Close()
		_, err := c.CreateOrganization(ctx, newCreateOrganizationReq(ownerEmail()))
		requirePermissionDenied(t, err, "environment viewer CreateOrganization")
	})

	t.Run("cannot add new member", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newAccountClientWithToken(t, *envViewerAccessTokenPath)
		defer c.Close()
		_, err := c.CreateAccountV2(ctx, newCreateAccountReq(*organizationID, memberEmail()))
		requirePermissionDenied(t, err, "environment viewer CreateAccountV2")
	})
}

// The system admin owns the e2e (system-admin) organization but is not a member
// of the organization these tests target, so it cannot write data there.
func TestSystemAdminPermissions(t *testing.T) {
	requireAccessToken(t, *sysAdminAccessTokenPath, "system admin")

	t.Run("cannot create account", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newAccountClientWithToken(t, *sysAdminAccessTokenPath)
		defer c.Close()
		_, err := c.CreateAccountV2(ctx, newCreateAccountReq(*organizationID, memberEmail()))
		requirePermissionDenied(t, err, "system admin CreateAccountV2")
	})

	t.Run("cannot create feature", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c := newFeatureClientWithToken(t, *sysAdminAccessTokenPath)
		defer c.Close()
		_, err := c.CreateFeature(ctx, newCreateFeatureReq(featureID("sysadmin")))
		requirePermissionDenied(t, err, "system admin CreateFeature")
	})
}

func newCreateOrganizationReq(owner string) *environmentproto.CreateOrganizationRequest {
	// Name is capped at 50 chars and url code must match ^[a-z0-9-]{1,50}$.
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	return &environmentproto.CreateOrganizationRequest{
		Name:        fmt.Sprintf("e2e-perm-%s", token),
		UrlCode:     fmt.Sprintf("e2e-perm-%s", token),
		Description: "Organization created by the e2e permission tests",
		OwnerEmail:  owner,
	}
}

func newCreateAccountReq(organizationID, email string) *accountproto.CreateAccountV2Request {
	return &accountproto.CreateAccountV2Request{
		OrganizationId:   organizationID,
		Email:            email,
		Name:             fmt.Sprintf("e2e-perm-%s", uniqueSuffix()),
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{
				Role:          accountproto.AccountV2_Role_Environment_VIEWER,
				EnvironmentId: *environmentID,
			},
		},
	}
}

func newCreateFeatureReq(id string) *featureproto.CreateFeatureRequest {
	return &featureproto.CreateFeatureRequest{
		Id:            id,
		EnvironmentId: *environmentID,
		Name:          "e2e-perm-feature",
		Description:   "Feature created by the e2e permission tests",
		Variations: []*featureproto.Variation{
			{Value: "A", Name: "Variation A", Description: "Thing does A"},
			{Value: "B", Name: "Variation B", Description: "Thing does B"},
		},
		Tags:                     []string{"e2e-perm"},
		DefaultOnVariationIndex:  &wrapperspb.Int32Value{Value: 0},
		DefaultOffVariationIndex: &wrapperspb.Int32Value{Value: 1},
	}
}

func uniqueSuffix() string {
	return fmt.Sprintf("%s-%d-%s", *testID, time.Now().UnixNano(), randomString())
}

func featureID(role string) string {
	return fmt.Sprintf("e2e-perm-%s-%s", role, uniqueSuffix())
}

func ownerEmail() string {
	return fmt.Sprintf("%s-perm-owner-%s@example.com", e2eAccountAddressPrefix, uniqueSuffix())
}

func memberEmail() string {
	return fmt.Sprintf("%s-perm-member-%s@example.com", e2eAccountAddressPrefix, uniqueSuffix())
}

// requirePermissionDenied accepts PermissionDenied (member lacking the role) or
// Unauthenticated (no membership in the target organization).
func requirePermissionDenied(t *testing.T, err error, op string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected a permission error, but the call succeeded", op)
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("%s: expected a gRPC status error, but got: %v", op, err)
	}
	switch st.Code() {
	case codes.PermissionDenied, codes.Unauthenticated:
		return
	default:
		t.Fatalf("%s: expected PermissionDenied or Unauthenticated, but got %s: %v", op, st.Code(), err)
	}
}

func requireAccessToken(t *testing.T, tokenPath, role string) {
	t.Helper()
	if tokenPath == "" {
		t.Skipf("skipping: no access token provided for the %s account", role)
	}
}

func newAccountClientWithToken(t *testing.T, tokenPath string) accountclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(tokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	c, err := accountclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create account client:", err)
	}
	return c
}

func newEnvironmentClientWithToken(t *testing.T, tokenPath string) environmentclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(tokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	c, err := environmentclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create environment client:", err)
	}
	return c
}

func newFeatureClientWithToken(t *testing.T, tokenPath string) featureclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(tokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	c, err := featureclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create feature client:", err)
	}
	return c
}
