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

	v2os "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/storage/v2"
	mockOpsCountStorage "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/storage/v2/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

func TestListOpsCountsPaginationValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		pageSize     int64
		cursor       string
		wantErr      error
		wantPageSize int
	}{
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
		},
		{
			name:         "zero page size and empty cursor list everything from the start",
			pageSize:     0,
			cursor:       "",
			wantPageSize: 0,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithTokenRoleUnassigned(t)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{"accept-language": []string{"ja"}})

			service := createAutoOpsService(mockController)
			var listed *v2os.ListOpsCountsParams
			if tc.wantErr == nil {
				service.opsCountStorage.(*mockOpsCountStorage.MockOpsCountStorage).EXPECT().
					ListOpsCounts(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2os.ListOpsCountsParams) ([]*autoopsproto.OpsCount, int, error) {
						listed = &p
						return []*autoopsproto.OpsCount{}, 0, nil
					})
			}

			resp, err := service.ListOpsCounts(ctx, &autoopsproto.ListOpsCountsRequest{
				EnvironmentId: "ns0",
				PageSize:      tc.pageSize,
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
			assert.Equal(t, tc.wantPageSize, listed.PageSize)
			// The cursor is handed to the storage as it is; the storage defaults
			// an empty cursor to the first page.
			assert.Equal(t, tc.cursor, listed.Cursor)
		})
	}
}
