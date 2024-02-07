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

package persister

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/service"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

func TestValidateEvent(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		input    *eventproto.UserEvent
		expected bool
	}{
		{
			input: &eventproto.UserEvent{
				UserId:   "hoge",
				LastSeen: 3456789,
			},
			expected: true,
		},
		{
			input:    &eventproto.UserEvent{},
			expected: false,
		},
		{
			input: &eventproto.UserEvent{
				UserId:   "",
				LastSeen: 3456789,
			},
			expected: false,
		},
		{
			input: &eventproto.UserEvent{
				UserId:   "hoge",
				LastSeen: 0,
			},
			expected: false,
		},
	}
	logger, _ := log.NewLogger()
	pst := persister{logger: logger}
	for _, p := range patterns {
		actual := pst.validateEvent(p.input)
		assert.Equal(t, p.expected, actual)
	}
}

func TestUpsert(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()
	uuid, err := uuid.NewUUID()
	require.NoError(t, err)

	patterns := []struct {
		desc               string
		setup              func(*persister)
		input              *eventproto.UserEvent
		expectedOK         bool
		expectedRepeatable bool
		expected           *userproto.User
	}{
		{
			desc: "upsert mau error",
			setup: func(p *persister) {
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("internal"))
			},
			input: &eventproto.UserEvent{
				EnvironmentNamespace: "ns0",
				UserId:               "id-1",
				LastSeen:             3,
			},
			expectedOK:         false,
			expectedRepeatable: true,
			expected: &userproto.User{
				Id:       "id-1",
				LastSeen: 3,
			},
		},
		{
			desc: "upsert success",
			setup: func(p *persister) {
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &eventproto.UserEvent{
				EnvironmentNamespace: "ns0",
				UserId:               "id-1",
				LastSeen:             3,
			},
			expectedOK:         true,
			expectedRepeatable: false,
			expected: &userproto.User{
				Id:       "id-1",
				LastSeen: 3,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pst := newPersisterWithMock(t, mockController, now, uuid)
			if p.setup != nil {
				p.setup(pst)
			}
			ok, repeatable := pst.upsert(p.input)
			assert.Equal(t, p.expectedOK, ok)
			assert.Equal(t, p.expectedRepeatable, repeatable)
		})
	}
}

func newPersisterWithMock(
	t *testing.T,
	mockController *gomock.Controller,
	now time.Time,
	id *uuid.UUID,
) *persister {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &persister{
		mysqlClient: mysqlmock.NewMockClient(mockController),
		timeNow:     func() time.Time { return now },
		newUUID:     func() (*uuid.UUID, error) { return id, nil },
		opts:        defaultOptions,
		logger:      logger,
	}
}
