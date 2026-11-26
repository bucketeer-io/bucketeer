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

	tagstorage "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage"
	tagstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
)

// listScopeEditor returns the editor of the caller for the environment
// scoping tests: a member of the organization with roles in env0 and env1 only.
func listScopeEditor(orgRole accountproto.AccountV2_Role_Organization) *eventproto.Editor {
	return &eventproto.Editor{
		Email:            "email",
		OrganizationRole: orgRole,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_VIEWER},
			{EnvironmentId: "env1", Role: accountproto.AccountV2_Role_Environment_EDITOR},
		},
	}
}

func TestListTagsMemberEnvironmentScope(t *testing.T) {
	t.Parallel()

	storedTags := func() []*proto.Tag {
		return []*proto.Tag{{Id: "tag0", Name: "tag0", EnvironmentId: "env1", EntityType: proto.Tag_FEATURE_FLAG}}
	}

	cases := []struct {
		name            string
		editor          *eventproto.Editor
		req             *proto.ListTagsRequest
		wantEnvFilter   []string
		wantStorageCall bool
		wantErr         error
	}{
		{
			// An UNASSIGNED entry grants no access to its environment.
			name: "member's unassigned environment entry is not an accessible environment",
			editor: &eventproto.Editor{
				Email:            "email",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_VIEWER},
					{EnvironmentId: "env2", Role: accountproto.AccountV2_Role_Environment_UNASSIGNED},
				},
			},
			req: &proto.ListTagsRequest{
				OrganizationId: "org0",
				EnvironmentId:  "env2",
				PageSize:       10,
			},
		},
		{
			name:   "member requesting an environment outside their roles gets an empty page",
			editor: listScopeEditor(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &proto.ListTagsRequest{
				OrganizationId: "org0",
				EnvironmentId:  "env2",
				PageSize:       10,
				Cursor:         "5",
			},
		},
		{
			name: "member without any environment role gets an empty page",
			editor: &eventproto.Editor{
				Email:            "email",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			},
			req: &proto.ListTagsRequest{OrganizationId: "org0", EnvironmentId: "env0", PageSize: 10},
		},
		{
			name:   "member with no accessible environment still gets an invalid cursor rejected",
			editor: listScopeEditor(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &proto.ListTagsRequest{
				OrganizationId: "org0",
				EnvironmentId:  "env2",
				Cursor:         "XXX",
			},
			wantErr: statusInvalidCursor.Err(),
		},
		{
			name:   "member with no accessible environment still gets an invalid order rejected",
			editor: listScopeEditor(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &proto.ListTagsRequest{
				OrganizationId: "org0",
				EnvironmentId:  "env2",
				OrderBy:        proto.ListTagsRequest_OrderBy(999),
			},
			wantErr: statusInvalidOrderBy.Err(),
		},
		{
			name:   "member requesting an environment they have a role in is scoped to it",
			editor: listScopeEditor(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &proto.ListTagsRequest{
				OrganizationId: "org0",
				EnvironmentId:  "env1",
				PageSize:       10,
			},
			wantEnvFilter:   []string{"env1"},
			wantStorageCall: true,
		},
		{
			name:   "admin requesting an environment outside their roles is not restricted",
			editor: listScopeEditor(accountproto.AccountV2_Role_Organization_ADMIN),
			req: &proto.ListTagsRequest{
				OrganizationId: "org0",
				EnvironmentId:  "env2",
				PageSize:       10,
			},
			wantEnvFilter:   []string{"env2"},
			wantStorageCall: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := createTagService(mockController)

			var listed *tagstorage.ListTagsParams
			if tc.wantStorageCall {
				service.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().
					ListTags(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p tagstorage.ListTagsParams) ([]*proto.Tag, int, int64, error) {
						listed = &p
						return storedTags(), 1, 1, nil
					})
			}

			resp, err := service.ListTags(createContextWithEditor(t, tc.editor), tc.req)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			if !tc.wantStorageCall {
				// Nothing is read from the storage and the page is empty.
				assert.Empty(t, resp.Tags)
				assert.Equal(t, int64(0), resp.TotalCount)
				expectedCursor := tc.req.Cursor
				if expectedCursor == "" {
					expectedCursor = "0"
				}
				assert.Equal(t, expectedCursor, resp.Cursor)
				return
			}
			require.NotNil(t, listed)
			assert.Equal(t, "org0", listed.OrganizationID)
			assert.Equal(t, tc.wantEnvFilter, listed.EnvironmentIDs)
			require.Len(t, resp.Tags, 1)
			assert.Equal(t, "tag0", resp.Tags[0].Id)
			assert.Equal(t, int64(1), resp.TotalCount)
		})
	}
}
