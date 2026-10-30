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

	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

// listPaginationCase describes a page_size / cursor pair of a list request
// and what the handler is expected to do with it: reject it without reading
// the storage, or pass the page size and the (defaulted) cursor through.
type listPaginationCase struct {
	name         string
	pageSize     int64
	cursor       string
	wantErr      error
	wantPageSize int
	wantCursor   string
}

func listPaginationCases() []listPaginationCase {
	return []listPaginationCase{
		{
			name:     "negative page size is rejected",
			pageSize: -1,
			cursor:   "0",
			wantErr:  statusInvalidPageSize.Err(),
		},
		{
			name:     "negative cursor is rejected",
			pageSize: 10,
			cursor:   "-5",
			wantErr:  statusInvalidCursor.Err(),
		},
		{
			name:     "non-numeric cursor is rejected",
			pageSize: 10,
			cursor:   "XXX",
			wantErr:  statusInvalidCursor.Err(),
		},
		{
			name:         "page size and cursor are passed to the storage",
			pageSize:     10,
			cursor:       "20",
			wantPageSize: 10,
			wantCursor:   "20",
		},
		{
			name:         "zero page size and empty cursor list everything from the start",
			pageSize:     0,
			cursor:       "",
			wantPageSize: 0,
			wantCursor:   "0",
		},
	}
}

func TestListAPIKeysPaginationValidation(t *testing.T) {
	t.Parallel()

	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithDefaultToken(t, false)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createAccountService(t, mockController, nil)
			service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
				GetAccountV2(gomock.Any(), "bucketeer@example.com", "org0").
				Return(listScopeAccount(accountproto.AccountV2_Role_Organization_ADMIN), nil).AnyTimes()

			var listed *v2as.ListAPIKeysParams
			if tc.wantErr == nil {
				service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					ListAPIKeys(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2as.ListAPIKeysParams) ([]*accountproto.APIKey, int, int64, error) {
						listed = &p
						return []*accountproto.APIKey{}, 0, 0, nil
					})
			}

			resp, err := service.ListAPIKeys(ctx, &accountproto.ListAPIKeysRequest{
				OrganizationId: "org0",
				PageSize:       tc.pageSize,
				Cursor:         tc.cursor,
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
			assert.Equal(t, tc.wantPageSize, listed.PageSize)
			assert.Equal(t, tc.wantCursor, listed.Cursor)
		})
	}
}

func TestListAccountsV2PaginationValidation(t *testing.T) {
	t.Parallel()

	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithDefaultToken(t, false)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createAccountService(t, mockController, nil)
			service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
				GetAccountV2(gomock.Any(), "bucketeer@example.com", "org0").
				Return(listScopeAccount(accountproto.AccountV2_Role_Organization_ADMIN), nil).AnyTimes()

			var listed *v2as.ListAccountsV2Params
			if tc.wantErr == nil {
				service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					ListAccountsV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2as.ListAccountsV2Params) ([]*accountproto.AccountV2, int, int64, error) {
						listed = &p
						return []*accountproto.AccountV2{}, 0, 0, nil
					})
			}

			resp, err := service.ListAccountsV2(ctx, &accountproto.ListAccountsV2Request{
				OrganizationId: "org0",
				PageSize:       tc.pageSize,
				Cursor:         tc.cursor,
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
			assert.Equal(t, tc.wantPageSize, listed.PageSize)
			assert.Equal(t, tc.wantCursor, listed.Cursor)
		})
	}
}
