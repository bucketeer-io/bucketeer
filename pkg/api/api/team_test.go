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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	teamclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/team/client/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	teamproto "github.com/bucketeer-io/bucketeer/v2/proto/team"
)

func TestGrpcGatewayService_CreateTeam(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.CreateTeamResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: create team error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.teamClient.(*teamclientmock.MockClient).EXPECT().CreateTeam(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.teamClient.(*teamclientmock.MockClient).EXPECT().CreateTeam(gomock.Any(), gomock.Any()).Return(
					&teamproto.CreateTeamResponse{
						Team: &teamproto.Team{
							Id:   "team-id",
							Name: "team-name",
						},
					}, nil)
			},
			expected: &gwproto.CreateTeamResponse{
				Team: &teamproto.Team{
					Id:   "team-id",
					Name: "team-name",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.CreateTeam(ctx, &gwproto.CreateTeamRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_DeleteTeam(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.DeleteTeamResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: create team error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.teamClient.(*teamclientmock.MockClient).EXPECT().DeleteTeam(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.teamClient.(*teamclientmock.MockClient).EXPECT().DeleteTeam(gomock.Any(), gomock.Any()).Return(
					&teamproto.DeleteTeamResponse{}, nil)
			},
			expected:    &gwproto.DeleteTeamResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.DeleteTeam(ctx, &gwproto.DeleteTeamRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_ListTeams(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.ListTeamsResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: create team error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.teamClient.(*teamclientmock.MockClient).EXPECT().ListTeams(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.teamClient.(*teamclientmock.MockClient).EXPECT().ListTeams(gomock.Any(), gomock.Any()).Return(
					&teamproto.ListTeamsResponse{
						Teams: []*teamproto.Team{
							{
								Id:          "team-id-1",
								Name:        "team-name-1",
								Description: "team-description-1",
							},
						},
						NextCursor: "1",
						TotalCount: 1,
					}, nil)
			},
			expected: &gwproto.ListTeamsResponse{
				Teams: []*teamproto.Team{
					{
						Id:          "team-id-1",
						Name:        "team-name-1",
						Description: "team-description-1",
					},
				},
				Cursor:     "1",
				TotalCount: 1,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.ListTeams(ctx, &gwproto.ListTeamsRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_TeamOrganizationScope(t *testing.T) {
	t.Parallel()

	for _, method := range []string{"create", "delete", "list"} {
		for _, tc := range []struct {
			name              string
			organizationID    string
			keyOrganizationID string
			wantErr           error
		}{
			{name: "same organization", organizationID: "org-a", keyOrganizationID: "org-a"},
			{name: "omitted organization", keyOrganizationID: "org-a"},
			{name: "other organization", organizationID: "org-b", keyOrganizationID: "org-a", wantErr: ErrNotFound},
			{name: "key without organization", wantErr: ErrInvalidAPIKey},
		} {
			t.Run(method+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				gs, ctx := newTeamWriteGateway(t, tc.keyOrganizationID)
				team := &teamproto.Team{Id: "team-a", Name: "team-a", OrganizationId: "org-a"}
				client := gs.teamClient.(*teamclientmock.MockClient)
				switch method {
				case "create":
					if tc.wantErr == nil {
						client.EXPECT().CreateTeam(teamMetadataMatcher(), &teamproto.CreateTeamRequest{
							OrganizationId: "org-a", Name: team.Name,
						}).Return(&teamproto.CreateTeamResponse{Team: team}, nil)
					}
					res, err := gs.CreateTeam(ctx, &gwproto.CreateTeamRequest{
						OrganizationId: tc.organizationID, Name: team.Name,
					})
					require.Equal(t, tc.wantErr, err)
					if tc.wantErr != nil {
						require.Nil(t, res)
						return
					}
					require.Equal(t, team, res.Team)
				case "delete":
					if tc.wantErr == nil {
						client.EXPECT().DeleteTeam(teamMetadataMatcher(), &teamproto.DeleteTeamRequest{
							OrganizationId: "org-a", Id: team.Id,
						}).Return(&teamproto.DeleteTeamResponse{}, nil)
					}
					res, err := gs.DeleteTeam(ctx, &gwproto.DeleteTeamRequest{
						OrganizationId: tc.organizationID, Id: team.Id,
					})
					require.Equal(t, tc.wantErr, err)
					if tc.wantErr != nil {
						require.Nil(t, res)
						return
					}
					require.Equal(t, &gwproto.DeleteTeamResponse{}, res)
				case "list":
					if tc.wantErr == nil {
						client.EXPECT().ListTeams(teamMetadataMatcher(), &teamproto.ListTeamsRequest{
							OrganizationId: "org-a", PageSize: 10,
						}).Return(&teamproto.ListTeamsResponse{
							Teams: []*teamproto.Team{team}, NextCursor: "1", TotalCount: 1,
						}, nil)
					}
					res, err := gs.ListTeams(ctx, &gwproto.ListTeamsRequest{
						OrganizationId: tc.organizationID, PageSize: 10,
					})
					require.Equal(t, tc.wantErr, err)
					if tc.wantErr != nil {
						require.Nil(t, res)
						return
					}
					require.Equal(t, []*teamproto.Team{team}, res.Teams)
					require.Equal(t, "1", res.Cursor)
					require.Equal(t, int64(1), res.TotalCount)
				}
			})
		}
	}
}

// newTeamWriteGateway returns a gateway whose cached key is a write key of the
// given organization, and a request context authenticated with that key.
func newTeamWriteGateway(t *testing.T, organizationID string) (*grpcGatewayService, context.Context) {
	t.Helper()
	gs := newGrpcGatewayServiceWithMock(t, gomock.NewController(t))
	key := &accountproto.EnvironmentAPIKey{
		Environment: &environmentproto.EnvironmentV2{Id: "env-a", OrganizationId: organizationID},
		ApiKey: &accountproto.APIKey{
			Id: "key-a", ApiKey: "0123456789abcdef0123456789abcdef", Name: "write-key",
			Role: accountproto.APIKey_PUBLIC_API_WRITE, Maintainer: "maintainer@example.com",
		},
	}
	gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(key.ApiKey.ApiKey).Return(key, nil)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", key.ApiKey.ApiKey))
	// Caller-supplied editor metadata must not replace the verified key's maintainer.
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(role.APIKeyMaintainerMDKey, "other-org-admin@example.com"))
	return gs, ctx
}

// teamMetadataMatcher matches a context carrying the verified key's editor metadata.
func teamMetadataMatcher() gomock.Matcher {
	return gomock.Cond(func(ctx context.Context) bool {
		md, ok := metadata.FromOutgoingContext(ctx)
		return ok &&
			len(md.Get(role.APIKeyTokenMDKey)) == 1 && md.Get(role.APIKeyTokenMDKey)[0] == "0123....cdef" &&
			len(md.Get(role.APIKeyMaintainerMDKey)) == 1 && md.Get(role.APIKeyMaintainerMDKey)[0] == "maintainer@example.com" &&
			len(md.Get(role.APIKeyNameMDKey)) == 1 && md.Get(role.APIKeyNameMDKey)[0] == "write-key"
	})
}
