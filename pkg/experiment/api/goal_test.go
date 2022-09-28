// Copyright 2022 The Bucketeer Authors.
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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	v2es "github.com/bucketeer-io/bucketeer/pkg/experiment/storage/v2"
	storeclient "github.com/bucketeer-io/bucketeer/pkg/storage/testing"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

func TestGetGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup                func(*experimentService)
		id                   string
		environmentNamespace string
		expectedErr          error
	}{
		{
			setup:                nil,
			id:                   "",
			environmentNamespace: "ns0",
			expectedErr:          errGoalIDRequiredJaJP,
		},
		{
			setup: func(s *experimentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:                   "id-0",
			environmentNamespace: "ns0",
			expectedErr:          errNotFoundJaJP,
		},
		{
			setup: func(s *experimentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:                   "id-1",
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		req := &experimentproto.GetGoalRequest{Id: p.id, EnvironmentNamespace: p.environmentNamespace}
		_, err := service.GetGoal(createContextWithTokenRoleUnassigned(), req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestListGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.ListGoalsRequest
		expectedErr error
	}{
		{
			setup: func(s *experimentService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req:         &experimentproto.ListGoalsRequest{EnvironmentNamespace: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.ListGoals(createContextWithTokenRoleUnassigned(), p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestCreateGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(s *experimentService)
		req         *experimentproto.CreateGoalRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Command:              nil,
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errNoCommandJaJP,
		},
		{
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Command:              &experimentproto.CreateGoalCommand{Id: ""},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errGoalIDRequiredJaJP,
		},
		{
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Command:              &experimentproto.CreateGoalCommand{Id: "bucketeer_goal_id?"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errInvalidGoalIDJaJP,
		},
		{
			setup: nil,
			req: &experimentproto.CreateGoalRequest{
				Command:              &experimentproto.CreateGoalCommand{Id: "Bucketeer-id-2019", Name: ""},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errGoalNameRequiredJaJP,
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalAlreadyExists)
			},
			req: &experimentproto.CreateGoalRequest{
				Command:              &experimentproto.CreateGoalCommand{Id: "Bucketeer-id-2019", Name: "name-0"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errAlreadyExistsJaJP,
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.CreateGoalRequest{
				Command:              &experimentproto.CreateGoalCommand{Id: "Bucketeer-id-2020", Name: "name-1"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		ctx := createContextWithToken()
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.CreateGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestUpdateGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.UpdateGoalRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.UpdateGoalRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errGoalIDRequiredJaJP,
		},
		{
			setup: nil,
			req: &experimentproto.UpdateGoalRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errNoCommandJaJP,
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalNotFound)
			},
			req: &experimentproto.UpdateGoalRequest{
				Id:                   "id-0",
				RenameCommand:        &experimentproto.RenameGoalCommand{Name: "name-0"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errNotFoundJaJP,
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.UpdateGoalRequest{
				Id:                   "id-1",
				RenameCommand:        &experimentproto.RenameGoalCommand{Name: "name-1"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		ctx := createContextWithToken()
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.UpdateGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestArchiveGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.ArchiveGoalRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.ArchiveGoalRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errGoalIDRequiredJaJP,
		},
		{
			setup: nil,
			req: &experimentproto.ArchiveGoalRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errNoCommandJaJP,
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalNotFound)
			},
			req: &experimentproto.ArchiveGoalRequest{
				Id:                   "id-0",
				Command:              &experimentproto.ArchiveGoalCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errNotFoundJaJP,
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.ArchiveGoalRequest{
				Id:                   "id-1",
				Command:              &experimentproto.ArchiveGoalCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		ctx := createContextWithToken()
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.ArchiveGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestDeleteGoalMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.DeleteGoalRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.DeleteGoalRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errGoalIDRequiredJaJP,
		},
		{
			setup: nil,
			req: &experimentproto.DeleteGoalRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errNoCommandJaJP,
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrGoalNotFound)
			},
			req: &experimentproto.DeleteGoalRequest{
				Id:                   "id-0",
				Command:              &experimentproto.DeleteGoalCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: errNotFoundJaJP,
		},
		{
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.DeleteGoalRequest{
				Id:                   "id-1",
				Command:              &experimentproto.DeleteGoalCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		ctx := createContextWithToken()
		service := createExperimentService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.DeleteGoal(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestGoalPermissionDenied(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleUnassigned()
	s := storeclient.NewInMemoryStorage()
	service := createExperimentService(mockController, s)
	patterns := map[string]struct {
		action   func(context.Context, *experimentService) error
		expected error
	}{
		"CreateGoal": {
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.CreateGoal(ctx, &experimentproto.CreateGoalRequest{})
				return err
			},
			expected: errPermissionDeniedJaJP,
		},
		"UpdateGoal": {
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.UpdateGoal(ctx, &experimentproto.UpdateGoalRequest{})
				return err
			},
			expected: errPermissionDeniedJaJP,
		},
		"DeleteGoal": {
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.DeleteGoal(ctx, &experimentproto.DeleteGoalRequest{})
				return err
			},
			expected: errPermissionDeniedJaJP,
		},
	}
	for msg, p := range patterns {
		actual := p.action(ctx, service)
		assert.Equal(t, p.expected, actual, "%s", msg)
	}
}
