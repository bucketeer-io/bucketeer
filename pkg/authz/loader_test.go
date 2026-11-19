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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	environmentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

func TestAccountLoader(t *testing.T) {
	t.Parallel()
	account := &accountproto.AccountV2{Email: "viewer@example.com"}
	notFound := status.Error(codes.NotFound, "account not found")
	patterns := []struct {
		desc        string
		scope       ScopeID
		setup       func(*accountclientmock.MockClient)
		expected    *accountproto.AccountV2
		expectedErr error
	}{
		{
			desc:  "organization scope loads by organization",
			scope: ScopeID{Kind: OrganizationScope, ID: "org-1"},
			setup: func(c *accountclientmock.MockClient) {
				c.EXPECT().GetAccountV2(gomock.Any(), &accountproto.GetAccountV2Request{
					Email: "viewer@example.com", OrganizationId: "org-1",
				}).Return(&accountproto.GetAccountV2Response{Account: account}, nil)
			},
			expected: account,
		},
		{
			desc:  "environment scope loads by environment",
			scope: ScopeID{Kind: EnvironmentScope, ID: "env-1"},
			setup: func(c *accountclientmock.MockClient) {
				c.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), &accountproto.GetAccountV2ByEnvironmentIDRequest{
					Email: "viewer@example.com", EnvironmentId: "env-1",
				}).Return(&accountproto.GetAccountV2ByEnvironmentIDResponse{Account: account}, nil)
			},
			expected: account,
		},
		{
			desc:  "organization of environment loads by environment",
			scope: ScopeID{Kind: OrganizationOfEnvironmentScope, ID: "env-1"},
			setup: func(c *accountclientmock.MockClient) {
				c.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).
					Return(&accountproto.GetAccountV2ByEnvironmentIDResponse{Account: account}, nil)
			},
			expected: account,
		},
		{
			desc:  "organization error",
			scope: ScopeID{Kind: OrganizationScope, ID: "org-1"},
			setup: func(c *accountclientmock.MockClient) {
				c.EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).Return(nil, notFound)
			},
			expectedErr: notFound,
		},
		{
			desc:  "environment error",
			scope: ScopeID{Kind: EnvironmentScope, ID: "env-1"},
			setup: func(c *accountclientmock.MockClient) {
				c.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(nil, notFound)
			},
			expectedErr: notFound,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			client := accountclientmock.NewMockClient(gomock.NewController(t))
			p.setup(client)
			actual, err := NewAccountLoader(client).Account(context.Background(), "viewer@example.com", p.scope)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestAccountLoaderSystemAdmin(t *testing.T) {
	t.Parallel()
	account := &accountproto.AccountV2{Email: "admin@example.com", OrganizationId: "org-admin"}
	unavailable := status.Error(codes.Unavailable, "account service")
	organizations := func(c *accountclientmock.MockClient, orgs ...*environmentproto.Organization) {
		c.EXPECT().GetMyOrganizationsByEmail(gomock.Any(), &accountproto.GetMyOrganizationsByEmailRequest{
			Email: "admin@example.com",
		}).Return(&accountproto.GetMyOrganizationsResponse{Organizations: orgs}, nil)
	}
	patterns := []struct {
		desc         string
		setup        func(*accountclientmock.MockClient)
		expected     *accountproto.AccountV2
		expectedCode codes.Code
	}{
		{
			desc: "loads the account of the system admin organization",
			setup: func(c *accountclientmock.MockClient) {
				organizations(c,
					&environmentproto.Organization{Id: "org-1"},
					&environmentproto.Organization{Id: "org-admin", SystemAdmin: true},
				)
				c.EXPECT().GetAccountV2(gomock.Any(), &accountproto.GetAccountV2Request{
					Email: "admin@example.com", OrganizationId: "org-admin",
				}).Return(&accountproto.GetAccountV2Response{Account: account}, nil)
			},
			expected: account,
		},
		{
			desc: "skips a system admin organization without an account or with a disabled one",
			setup: func(c *accountclientmock.MockClient) {
				organizations(c,
					&environmentproto.Organization{Id: "org-none", SystemAdmin: true},
					&environmentproto.Organization{Id: "org-disabled", SystemAdmin: true},
					&environmentproto.Organization{Id: "org-admin", SystemAdmin: true},
				)
				c.EXPECT().GetAccountV2(gomock.Any(), &accountproto.GetAccountV2Request{
					Email: "admin@example.com", OrganizationId: "org-none",
				}).Return(nil, status.Error(codes.NotFound, "account not found"))
				c.EXPECT().GetAccountV2(gomock.Any(), &accountproto.GetAccountV2Request{
					Email: "admin@example.com", OrganizationId: "org-disabled",
				}).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{Email: "admin@example.com", Disabled: true},
				}, nil)
				c.EXPECT().GetAccountV2(gomock.Any(), &accountproto.GetAccountV2Request{
					Email: "admin@example.com", OrganizationId: "org-admin",
				}).Return(&accountproto.GetAccountV2Response{Account: account}, nil)
			},
			expected: account,
		},
		{
			desc: "no system admin organization",
			setup: func(c *accountclientmock.MockClient) {
				organizations(c, &environmentproto.Organization{Id: "org-1"})
			},
			expectedCode: codes.NotFound,
		},
		{
			desc: "account error",
			setup: func(c *accountclientmock.MockClient) {
				organizations(c, &environmentproto.Organization{Id: "org-admin", SystemAdmin: true})
				c.EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).Return(nil, unavailable)
			},
			expectedCode: codes.Unavailable,
		},
		{
			desc: "organizations error",
			setup: func(c *accountclientmock.MockClient) {
				c.EXPECT().GetMyOrganizationsByEmail(gomock.Any(), gomock.Any()).Return(nil, unavailable)
			},
			expectedCode: codes.Unavailable,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			client := accountclientmock.NewMockClient(gomock.NewController(t))
			p.setup(client)
			actual, err := NewAccountLoader(client).SystemAdmin(context.Background(), "admin@example.com")
			assert.Equal(t, p.expectedCode, status.Code(err))
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestLookups(t *testing.T) {
	t.Parallel()
	notFound := status.Error(codes.NotFound, "project not found")
	patterns := []struct {
		desc        string
		setup       func(*environmentclientmock.MockClient)
		expected    string
		expectedErr error
	}{
		{
			desc: "organization of the project",
			setup: func(c *environmentclientmock.MockClient) {
				c.EXPECT().GetProject(gomock.Any(), &environmentproto.GetProjectRequest{Id: "proj-1"}).
					Return(&environmentproto.GetProjectResponse{
						Project: &environmentproto.Project{Id: "proj-1", OrganizationId: "org-1"},
					}, nil)
			},
			expected: "org-1",
		},
		{
			desc: "unknown project",
			setup: func(c *environmentclientmock.MockClient) {
				c.EXPECT().GetProject(gomock.Any(), gomock.Any()).Return(nil, notFound)
			},
			expectedErr: notFound,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			client := environmentclientmock.NewMockClient(gomock.NewController(t))
			p.setup(client)
			actual, err := NewLookups(client).OrganizationOfProject(context.Background(), "proj-1")
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}
