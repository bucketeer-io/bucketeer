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

	v2ss "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/storage/v2"
	v2mock "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/storage/v2/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

// listPaginationCase describes a page_size / cursor pair of a list request
// and what the handler is expected to do with it: reject it without reading
// the storage, or pass the page size and cursor through to the storage.
type listPaginationCase struct {
	name     string
	pageSize int64
	cursor   string
	wantErr  error
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
			name:     "page size and cursor are passed to the storage",
			pageSize: 10,
			cursor:   "20",
		},
		{
			name:     "zero page size and empty cursor list everything from the start",
			pageSize: 0,
			cursor:   "",
		},
	}
}

// assertListPaginationResult checks the outcome of a list call: a rejected
// request returns the expected error and never reaches the storage; an
// accepted one reaches the storage with the page size and cursor of the
// request (the storage defaults an empty cursor to the first page).
func assertListPaginationResult(
	t *testing.T,
	tc listPaginationCase,
	err error,
	listed bool,
	listedPageSize int64,
	listedCursor string,
) {
	t.Helper()
	if tc.wantErr != nil {
		require.Error(t, err)
		assert.Equal(t, tc.wantErr.Error(), err.Error())
		assert.False(t, listed)
		return
	}
	require.NoError(t, err)
	require.True(t, listed)
	assert.Equal(t, tc.pageSize, listedPageSize)
	assert.Equal(t, tc.cursor, listedCursor)
}

func listPaginationContext(t *testing.T, isSystemAdmin bool) context.Context {
	t.Helper()
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"accept-language": []string{"ja"}})
	return setToken(t, ctx, isSystemAdmin)
}

func expectListSubscriptions(service *SubscriptionService, listed **v2ss.ListSubscriptionsParams) {
	service.subscriptionStorage.(*v2mock.MockSubscriptionStorage).EXPECT().
		ListSubscriptions(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, p v2ss.ListSubscriptionsParams) ([]*proto.Subscription, int, int64, error) {
			*listed = &p
			return []*proto.Subscription{}, 0, 0, nil
		})
}

func expectListAdminSubscriptions(service *SubscriptionService, listed **v2ss.ListAdminSubscriptionsParams) {
	service.adminSubscriptionStorage.(*v2mock.MockAdminSubscriptionStorage).EXPECT().
		ListAdminSubscriptions(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, p v2ss.ListAdminSubscriptionsParams) ([]*proto.Subscription, int, int64, error) {
			*listed = &p
			return []*proto.Subscription{}, 0, 0, nil
		})
}

func TestListSubscriptionsPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newSubscriptionService(mockController, nil, nil, nil)
			var listed *v2ss.ListSubscriptionsParams
			if tc.wantErr == nil {
				expectListSubscriptions(service, &listed)
			}

			resp, err := service.ListSubscriptions(listPaginationContext(t, false), &proto.ListSubscriptionsRequest{
				EnvironmentId: "ns0",
				PageSize:      tc.pageSize,
				Cursor:        tc.cursor,
			})
			var pageSize int64
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListEnabledSubscriptionsPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newSubscriptionService(mockController, nil, nil, nil)
			var listed *v2ss.ListSubscriptionsParams
			if tc.wantErr == nil {
				expectListSubscriptions(service, &listed)
			}

			resp, err := service.ListEnabledSubscriptions(listPaginationContext(t, false), &proto.ListEnabledSubscriptionsRequest{
				EnvironmentId: "ns0",
				PageSize:      tc.pageSize,
				Cursor:        tc.cursor,
			})
			var pageSize int64
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListAdminSubscriptionsPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newSubscriptionService(mockController, nil, nil, nil)
			var listed *v2ss.ListAdminSubscriptionsParams
			if tc.wantErr == nil {
				expectListAdminSubscriptions(service, &listed)
			}

			resp, err := service.ListAdminSubscriptions(listPaginationContext(t, true), &proto.ListAdminSubscriptionsRequest{
				PageSize: tc.pageSize,
				Cursor:   tc.cursor,
			})
			var pageSize int64
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListEnabledAdminSubscriptionsPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newSubscriptionService(mockController, nil, nil, nil)
			var listed *v2ss.ListAdminSubscriptionsParams
			if tc.wantErr == nil {
				expectListAdminSubscriptions(service, &listed)
			}

			resp, err := service.ListEnabledAdminSubscriptions(
				listPaginationContext(t, true),
				&proto.ListEnabledAdminSubscriptionsRequest{
					PageSize: tc.pageSize,
					Cursor:   tc.cursor,
				},
			)
			var pageSize int64
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}
