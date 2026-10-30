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

	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
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
	listedPageSize int,
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
	assert.Equal(t, int(tc.pageSize), listedPageSize)
	assert.Equal(t, tc.cursor, listedCursor)
}

func systemAdminContext(t *testing.T) context.Context {
	t.Helper()
	return metadata.NewIncomingContext(createContextWithToken(t), metadata.MD{"accept-language": []string{"ja"}})
}

func TestListProjectsPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newEnvironmentService(t, mockController, nil)
			var listed *v2es.ListProjectsParams
			if tc.wantErr == nil {
				service.projectStorage.(*storagemock.MockProjectStorage).EXPECT().
					ListProjects(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2es.ListProjectsParams) ([]*proto.Project, int, int64, error) {
						listed = &p
						return []*proto.Project{}, 0, 0, nil
					})
			}

			resp, err := service.ListProjects(systemAdminContext(t), &proto.ListProjectsRequest{
				PageSize: tc.pageSize,
				Cursor:   tc.cursor,
			})
			var pageSize int
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListProjectsV2PaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newEnvironmentService(t, mockController, nil)
			var listed *v2es.ListProjectsParams
			if tc.wantErr == nil {
				service.projectStorage.(*storagemock.MockProjectStorage).EXPECT().
					ListProjects(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2es.ListProjectsParams) ([]*proto.Project, int, int64, error) {
						listed = &p
						return []*proto.Project{}, 0, 0, nil
					})
			}

			resp, err := service.ListProjectsV2(systemAdminContext(t), &proto.ListProjectsV2Request{
				OrganizationId: "org-1",
				PageSize:       tc.pageSize,
				Cursor:         tc.cursor,
			})
			var pageSize int
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListEnvironmentsV2PaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newEnvironmentService(t, mockController, nil)
			var listed *v2es.ListEnvironmentsV2Params
			if tc.wantErr == nil {
				service.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
					ListEnvironmentsV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2es.ListEnvironmentsV2Params) ([]*proto.EnvironmentV2, int, int64, error) {
						listed = &p
						return []*proto.EnvironmentV2{}, 0, 0, nil
					})
			}

			resp, err := service.ListEnvironmentsV2(systemAdminContext(t), &proto.ListEnvironmentsV2Request{
				OrganizationId: "org-1",
				PageSize:       tc.pageSize,
				Cursor:         tc.cursor,
			})
			var pageSize int
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListOrganizationsPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newEnvironmentService(t, mockController, nil)
			var listed *v2es.ListOrganizationsParams
			if tc.wantErr == nil {
				service.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().
					ListOrganizations(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2es.ListOrganizationsParams) ([]*proto.Organization, int, int64, error) {
						listed = &p
						return []*proto.Organization{}, 0, 0, nil
					})
			}

			resp, err := service.ListOrganizations(systemAdminContext(t), &proto.ListOrganizationsRequest{
				PageSize: tc.pageSize,
				Cursor:   tc.cursor,
			})
			var pageSize int
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}
