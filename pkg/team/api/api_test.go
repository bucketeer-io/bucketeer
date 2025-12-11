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

package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/team/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/team/storage"
	teamstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/team/storage/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/team"
)

func TestNewTeamService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	mysqlClientMock := mysqlmock.NewMockClient(mockController)
	accountClientMock := accountclientmock.NewMockClient(mockController)
	p := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewTeamService(
		mysqlClientMock,
		accountClientMock,
		p,
		WithLogger(logger),
	)
	assert.IsType(t, &TeamService{}, s)
}

func TestTeamService_CreateTeam(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"en"},
	})

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*TeamService)
		req         *proto.CreateTeamRequest
		expectedErr error
	}{
		{
			desc: "err: unauthenticated",
			ctx:  context.TODO(),
			req: &proto.CreateTeamRequest{
				OrganizationId: "ns0",
				Name:           "test-team",
			},
			expectedErr: statusUnauthenticated.Err(),
		},
		{
			desc: "err: ErrNameRequired",
			ctx:  ctx,
			req: &proto.CreateTeamRequest{
				OrganizationId: "ns0",
			},
			expectedErr: statusNameRequired.Err(),
		},
		{
			desc: "success: insert team",
			ctx:  ctx,
			setup: func(s *TeamService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.teamStorage.(*teamstoragemock.MockTeamStorage).EXPECT().GetTeamByName(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, storage.ErrTeamNotFound)
				s.teamStorage.(*teamstoragemock.MockTeamStorage).EXPECT().UpsertTeam(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.teamStorage.(*teamstoragemock.MockTeamStorage).EXPECT().GetTeamByName(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Team{
					Team: &proto.Team{
						Id:   "team1",
						Name: "test-team",
					},
				}, nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateTeamRequest{
				OrganizationId: "ns0",
				Name:           "test-team",
			},
			expectedErr: nil,
		},
		{
			desc: "success: team already exists",
			ctx:  ctx,
			setup: func(s *TeamService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.teamStorage.(*teamstoragemock.MockTeamStorage).EXPECT().UpsertTeam(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.teamStorage.(*teamstoragemock.MockTeamStorage).EXPECT().GetTeamByName(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Team{
					Team: &proto.Team{
						Id:   "team1",
						Name: "test-team",
					},
				}, nil).Times(2)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateTeamRequest{
				OrganizationId: "ns0",
				Name:           "test-team",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createTeamService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.CreateTeam(p.ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestTeamService_ListTeams(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"en"},
	})

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*TeamService)
		req         *proto.ListTeamsRequest
		expectedRes *proto.ListTeamsResponse
		expectedErr error
	}{
		{
			desc: "err: unauthenticated",
			ctx:  context.TODO(),
			req: &proto.ListTeamsRequest{
				OrganizationId: "ns0",
			},
			expectedRes: nil,
			expectedErr: statusUnauthenticated.Err(),
		},
		{
			desc: "err: invalid cursor",
			ctx:  ctx,
			req: &proto.ListTeamsRequest{
				OrganizationId: "ns0",
				Cursor:         "invalid",
			},
			expectedRes: nil,
			expectedErr: statusInvalidCursor.Err(),
		},
		{
			desc: "err: invalid order_by",
			ctx:  ctx,
			req: &proto.ListTeamsRequest{
				OrganizationId: "ns0",
				OrderBy:        proto.ListTeamsRequest_OrderBy(999),
			},
			expectedRes: nil,
			expectedErr: statusInvalidOrderBy.Err(),
		},
		{
			desc: "success: no teams",
			ctx:  ctx,
			setup: func(s *TeamService) {
				s.teamStorage.(*teamstoragemock.MockTeamStorage).EXPECT().ListTeams(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Team{}, 0, int64(0), nil)
			},
			req: &proto.ListTeamsRequest{
				OrganizationId: "ns0",
				Cursor:         "0",
				PageSize:       10,
			},
			expectedRes: &proto.ListTeamsResponse{
				Teams:      []*proto.Team{},
				TotalCount: 0,
				NextCursor: "0",
			},
			expectedErr: nil,
		},
		{
			desc: "success: with teams",
			ctx:  ctx,
			setup: func(s *TeamService) {
				s.teamStorage.(*teamstoragemock.MockTeamStorage).EXPECT().ListTeams(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Team{
					{Id: "team1", Name: "Team 1"},
					{Id: "team2", Name: "Team 2"},
				}, 2, int64(2), nil)
			},
			req: &proto.ListTeamsRequest{
				OrganizationId: "ns0",
				Cursor:         "0",
				PageSize:       10,
			},
			expectedRes: &proto.ListTeamsResponse{
				Teams: []*proto.Team{
					{Id: "team1", Name: "Team 1"},
					{Id: "team2", Name: "Team 2"},
				},
				TotalCount: 2,
				NextCursor: "2",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createTeamService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			res, err := s.ListTeams(p.ctx, p.req)
			assert.Equal(t, p.expectedRes, res)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestTeamService_DeleteTeam(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"en"},
	})

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*TeamService)
		req         *proto.DeleteTeamRequest
		expectedRes *proto.DeleteTeamResponse
		expectedErr error
	}{
		{
			desc: "err: unauthenticated",
			ctx:  context.TODO(),
			req: &proto.DeleteTeamRequest{
				OrganizationId: "ns0",
				Id:             "team1",
			},
			expectedRes: nil,
			expectedErr: statusUnauthenticated.Err(),
		},
		{
			desc: "err: team id required",
			ctx:  ctx,
			req: &proto.DeleteTeamRequest{
				OrganizationId: "ns0",
				Id:             "",
			},
			expectedRes: nil,
			expectedErr: statusTeamIDRequired.Err(),
		},
		{
			desc: "err: team not found",
			ctx:  ctx,
			setup: func(s *TeamService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context, mysql.Transaction) error) error {
					return storage.ErrTeamNotFound
				})
			},
			req: &proto.DeleteTeamRequest{
				OrganizationId: "ns0",
				Id:             "team1",
			},
			expectedRes: nil,
			expectedErr: statusTeamNotFound.Err(),
		},
		{
			desc: "err: team is in use",
			ctx:  ctx,
			setup: func(s *TeamService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) error {
					return fn(ctx, nil)
				})
				s.teamStorage.(*teamstoragemock.MockTeamStorage).EXPECT().GetTeam(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Team{
					Team: &proto.Team{
						Id:             "team1",
						OrganizationId: "ns0",
						Name:           "team1",
					},
				}, nil)
				s.accountClient.(*accountclientmock.MockClient).EXPECT().ListAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.ListAccountsV2Response{
					Accounts: []*accountproto.AccountV2{
						{
							Email: "bucketeer@gmail.com",
							Teams: []string{"team1"},
						},
					},
				}, nil)
			},
			req: &proto.DeleteTeamRequest{
				OrganizationId: "ns0",
				Id:             "team1",
			},
			expectedRes: nil,
			expectedErr: statusTeamInUsed.Err(),
		},
		{
			desc: "success",
			ctx:  ctx,
			setup: func(s *TeamService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.teamStorage.(*teamstoragemock.MockTeamStorage).EXPECT().GetTeam(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Team{
					Team: &proto.Team{
						Id:             "team1",
						OrganizationId: "ns0",
					},
				}, nil)
				s.accountClient.(*accountclientmock.MockClient).EXPECT().ListAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.ListAccountsV2Response{}, nil)
				s.teamStorage.(*teamstoragemock.MockTeamStorage).EXPECT().DeleteTeam(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.DeleteTeamRequest{
				OrganizationId: "ns0",
				Id:             "team1",
			},
			expectedRes: &proto.DeleteTeamResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createTeamService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			res, err := s.DeleteTeam(p.ctx, p.req)
			assert.Equal(t, p.expectedRes, res)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func createTeamService(c *gomock.Controller) *TeamService {
	mysqlClientMock := mysqlmock.NewMockClient(c)
	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2Response{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
		},
	}
	accountClientMock.EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	p := publishermock.NewMockPublisher(c)
	logger := zap.NewNop()
	return &TeamService{
		mysqlClient:   mysqlClientMock,
		teamStorage:   teamstoragemock.NewMockTeamStorage(c),
		accountClient: accountClientMock,
		publisher:     p,
		opts: &options{
			logger: zap.NewNop(),
		},
		logger: logger,
	}
}

func createContextWithToken(t *testing.T) context.Context {
	t.Helper()
	accessToken := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: false,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, accessToken)
}
