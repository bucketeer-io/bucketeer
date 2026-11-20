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

package api

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

type methodStream struct {
	grpc.ServerTransportStream
	method string
}

func (s *methodStream) Method() string { return s.method }

func contextWithMethod(method string) context.Context {
	return grpc.NewContextWithServerTransportStream(
		context.Background(),
		&methodStream{method: "/bucketeer.gateway.Gateway/" + method},
	)
}

func TestMaintainerReadsAreGatewayMethods(t *testing.T) {
	t.Parallel()
	service := gwproto.File_proto_gateway_service_proto.Services().ByName("Gateway")
	for name := range maintainerReads {
		assert.NotNil(t, service.Methods().ByName(protoreflect.Name(name)), name)
	}
}

func TestMaintainerReadDenial(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		account  *accountproto.AccountV2
		scope    maintainerReadScope
		expected string
	}{
		{
			desc:     "no account",
			account:  nil,
			scope:    maintainerReadEnvironment,
			expected: maintainerReadReasonNoAccount,
		},
		{
			desc:     "disabled",
			account:  &accountproto.AccountV2{Disabled: true, OrganizationRole: accountproto.AccountV2_Role_Organization_OWNER},
			scope:    maintainerReadOrganization,
			expected: maintainerReadReasonDisabled,
		},
		{
			desc: "no role in the environment",
			account: &accountproto.AccountV2{
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "other", Role: accountproto.AccountV2_Role_Environment_EDITOR},
				},
			},
			scope:    maintainerReadEnvironment,
			expected: maintainerReadReasonNoEnvironmentRole,
		},
		{
			desc: "viewer in the environment",
			account: &accountproto.AccountV2{
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env", Role: accountproto.AccountV2_Role_Environment_VIEWER},
				},
			},
			scope:    maintainerReadEnvironment,
			expected: "",
		},
		{
			desc:     "organization admin without environment roles",
			account:  &accountproto.AccountV2{OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN},
			scope:    maintainerReadEnvironment,
			expected: "",
		},
		{
			desc:     "member reading the organization",
			account:  &accountproto.AccountV2{OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER},
			scope:    maintainerReadOrganization,
			expected: "",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, maintainerReadDenial(p.account, p.scope, "env"))
		})
	}
}

func TestReportMaintainerRead(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	envAPIKey := func(id string) *accountproto.EnvironmentAPIKey {
		return &accountproto.EnvironmentAPIKey{
			ApiKey:      &accountproto.APIKey{Id: id, Maintainer: "maintainer@example.com"},
			Environment: &environmentproto.EnvironmentV2{Id: "env-" + id, OrganizationId: "org-" + id},
		}
	}
	counted := func(id, method, reason string) float64 {
		return testutil.ToFloat64(maintainerReadDeniedCounter.WithLabelValues("org-"+id, "env-"+id, method, reason))
	}
	patterns := []struct {
		desc     string
		keyID    string
		method   string
		calls    int
		setup    func(*grpcGatewayService)
		reason   string
		expected float64
	}{
		{
			desc:   "a write is not checked",
			keyID:  "write",
			method: "CreateFeature",
			calls:  1,
		},
		{
			desc:   "a read with maintainer metadata is not checked",
			keyID:  "metadata",
			method: "ListTags",
			calls:  1,
		},
		{
			desc:   "no account: counted on every read, looked up once",
			keyID:  "none",
			method: "GetFeature",
			calls:  2,
			setup: func(s *grpcGatewayService) {
				s.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), &accountproto.GetAccountV2Request{
						Email:          "maintainer@example.com",
						OrganizationId: "org-none",
					},
				).Return(nil, status.Error(codes.NotFound, "not found")).Times(1)
			},
			reason:   maintainerReadReasonNoAccount,
			expected: 2,
		},
		{
			desc:   "no environment role: counted for an environment read",
			keyID:  "norole",
			method: "ListExperiments",
			calls:  1,
			setup: func(s *grpcGatewayService) {
				s.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.GetAccountV2Response{Account: &accountproto.AccountV2{
					OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				}}, nil).Times(1)
			},
			reason:   maintainerReadReasonNoEnvironmentRole,
			expected: 1,
		},
		{
			desc:   "no environment role: not counted for an organization read",
			keyID:  "orgread",
			method: "ListProjects",
			calls:  1,
			setup: func(s *grpcGatewayService) {
				s.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.GetAccountV2Response{Account: &accountproto.AccountV2{
					OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				}}, nil).Times(1)
			},
			reason:   maintainerReadReasonNoEnvironmentRole,
			expected: 0,
		},
		{
			desc:   "lookup error: not counted",
			keyID:  "error",
			method: "GetFeature",
			calls:  1,
			setup: func(s *grpcGatewayService) {
				s.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, status.Error(codes.Internal, "internal")).Times(1)
			},
			reason:   maintainerReadReasonNoAccount,
			expected: 0,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(gs)
			}
			ctx := metadata.NewIncomingContext(contextWithMethod(p.method), metadata.MD{})
			for i := 0; i < p.calls; i++ {
				gs.reportMaintainerRead(ctx, envAPIKey(p.keyID))
			}
			assert.Equal(t, p.expected, counted(p.keyID, p.method, p.reason))
		})
	}
}
