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

	v2asmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	v2als "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2"
	v2alsmock "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
)

// listCursorCase describes a cursor of a list request and what the handler is
// expected to do with it: reject it without reading the storage, or pass the
// (defaulted) cursor through.
type listCursorCase struct {
	name       string
	cursor     string
	wantErr    error
	wantCursor string
}

func listCursorCases() []listCursorCase {
	return []listCursorCase{
		{
			name:    "negative cursor is rejected",
			cursor:  "-5",
			wantErr: statusInvalidCursor.Err(),
		},
		{
			name:    "non-numeric cursor is rejected",
			cursor:  "XXX",
			wantErr: statusInvalidCursor.Err(),
		},
		{
			name:       "cursor is passed to the storage",
			cursor:     "20",
			wantCursor: "20",
		},
		{
			name:       "empty cursor lists from the start",
			cursor:     "",
			wantCursor: "0",
		},
	}
}

func TestListAuditLogsCursorValidation(t *testing.T) {
	t.Parallel()

	for _, tc := range listCursorCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newAuditLogServiceWithGetAccountByEnvironmentMock(
				t, mockController,
				accountproto.AccountV2_Role_Organization_OWNER,
				accountproto.AccountV2_Role_Environment_EDITOR,
			)
			var listed *v2als.ListAuditLogsParams
			if tc.wantErr == nil {
				service.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().
					ListAuditLogs(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2als.ListAuditLogsParams) ([]*proto.AuditLog, int, int64, error) {
						listed = &p
						return []*proto.AuditLog{}, 0, 0, nil
					})
				service.accountStorage.(*v2asmock.MockAccountStorage).EXPECT().
					GetAvatarAccountsV2(gomock.Any(), gomock.Any()).
					Return([]*accountproto.AccountV2{}, nil).AnyTimes()
			}

			resp, err := service.ListAuditLogs(createContextWithToken(t, true), &proto.ListAuditLogsRequest{
				EnvironmentId: "ns0",
				PageSize:      10,
				Cursor:        tc.cursor,
			})
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, resp)
				assert.Nil(t, listed)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, listed)
			assert.Equal(t, 10, listed.PageSize)
			assert.Equal(t, tc.wantCursor, listed.Cursor)
		})
	}
}

func TestListFeatureHistoryCursorValidation(t *testing.T) {
	t.Parallel()

	for _, tc := range listCursorCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newAuditLogServiceWithGetAccountByEnvironmentMock(
				t, mockController,
				accountproto.AccountV2_Role_Organization_OWNER,
				accountproto.AccountV2_Role_Environment_EDITOR,
			)
			var listed *v2als.ListAuditLogsParams
			if tc.wantErr == nil {
				service.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().
					ListAuditLogs(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2als.ListAuditLogsParams) ([]*proto.AuditLog, int, int64, error) {
						listed = &p
						return []*proto.AuditLog{}, 0, 0, nil
					})
				service.accountStorage.(*v2asmock.MockAccountStorage).EXPECT().
					GetAvatarAccountsV2(gomock.Any(), gomock.Any()).
					Return([]*accountproto.AccountV2{}, nil).AnyTimes()
			}

			resp, err := service.ListFeatureHistory(createContextWithToken(t, false), &proto.ListFeatureHistoryRequest{
				EnvironmentId: "ns0",
				FeatureId:     "feature-id",
				PageSize:      10,
				Cursor:        tc.cursor,
			})
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, resp)
				assert.Nil(t, listed)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, listed)
			assert.Equal(t, 10, listed.PageSize)
			assert.Equal(t, tc.wantCursor, listed.Cursor)
		})
	}
}
