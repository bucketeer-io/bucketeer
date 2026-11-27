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

package authz

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

type accountLoader struct {
	client accountproto.AccountServiceClient
}

// NewAccountLoader loads accounts through the account service.
func NewAccountLoader(client accountproto.AccountServiceClient) AccountLoader {
	return &accountLoader{client: client}
}

func (l *accountLoader) Account(ctx context.Context, email string, scope ScopeID) (*accountproto.AccountV2, error) {
	if scope.Kind == OrganizationScope {
		resp, err := l.client.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
			Email:          email,
			OrganizationId: scope.ID,
		})
		if err != nil {
			return nil, err
		}
		return resp.Account, nil
	}
	resp, err := l.client.GetAccountV2ByEnvironmentID(ctx, &accountproto.GetAccountV2ByEnvironmentIDRequest{
		Email:         email,
		EnvironmentId: scope.ID,
	})
	if err != nil {
		return nil, err
	}
	return resp.Account, nil
}

// SystemAdmin loads the caller's enabled account in a system admin organization.
func (l *accountLoader) SystemAdmin(ctx context.Context, email string) (*accountproto.AccountV2, error) {
	resp, err := l.client.GetMyOrganizationsByEmail(ctx, &accountproto.GetMyOrganizationsByEmailRequest{
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	for _, org := range resp.Organizations {
		if !org.SystemAdmin {
			continue
		}
		account, err := l.Account(ctx, email, ScopeID{Kind: OrganizationScope, ID: org.Id})
		if status.Code(err) == codes.NotFound || (err == nil && account.Disabled) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, status.Error(codes.NotFound, "authz: no system admin account")
}

type lookups struct {
	client environmentproto.EnvironmentServiceClient
}

// NewLookups finds scopes through the environment service.
func NewLookups(client environmentproto.EnvironmentServiceClient) Lookups {
	return &lookups{client: client}
}

func (l *lookups) OrganizationOfProject(ctx context.Context, projectID string) (string, error) {
	resp, err := l.client.GetProject(ctx, &environmentproto.GetProjectRequest{Id: projectID})
	if err != nil {
		return "", err
	}
	return resp.Project.OrganizationId, nil
}
