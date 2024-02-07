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

package feature

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	wrappersproto "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"

	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	segmentUserRetryTimes = 20
)

func TestAddSegmentUserCommand(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	segmentID := createSegment(ctx, t, client).Id
	userID := newUserID(t)
	testcases := []struct {
		userID string
		state  featureproto.SegmentUser_State
	}{
		{
			userID: userID,
			state:  featureproto.SegmentUser_INCLUDED,
		},
	}
	for _, tc := range testcases {
		addSegmentUser(ctx, t, client, segmentID, []string{tc.userID}, tc.state)
		user := getSegmentUser(ctx, t, client, segmentID, tc.userID, tc.state)
		id := domain.SegmentUserID(segmentID, tc.userID, tc.state)
		assert.Equal(t, id, user.Id)
		assert.Equal(t, segmentID, user.SegmentId)
		assert.Equal(t, tc.userID, user.UserId)
		assert.Equal(t, tc.state, user.State)
		assert.Equal(t, false, user.Deleted)
	}
}

func TestDeleteSegmentUserCommand(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	segmentID := createSegment(ctx, t, client).Id
	userID := newUserID(t)
	testcases := []struct {
		userID string
		state  featureproto.SegmentUser_State
	}{
		{
			userID: userID,
			state:  featureproto.SegmentUser_INCLUDED,
		},
	}
	for _, tc := range testcases {
		addSegmentUser(ctx, t, client, segmentID, []string{tc.userID}, tc.state)
		deleteSegmentUser(ctx, t, client, segmentID, []string{tc.userID}, tc.state)
		listRes := listSegmentUsers(
			ctx,
			t,
			client,
			segmentID,
			&wrappersproto.Int32Value{Value: int32(tc.state)},
		)
		assert.Empty(t, len(listRes.Users))
	}
}

func TestListSegmentUsersPageSize(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	segmentID := createSegment(ctx, t, client).Id
	userIDs := []string{newUserID(t), newUserID(t)}
	addSegmentUser(ctx, t, client, segmentID, userIDs, featureproto.SegmentUser_INCLUDED)
	pageSize := int64(1)
	res, err := client.ListSegmentUsers(ctx, &featureproto.ListSegmentUsersRequest{
		PageSize:             pageSize,
		SegmentId:            segmentID,
		State:                &wrappersproto.Int32Value{Value: int32(featureproto.SegmentUser_INCLUDED)},
		EnvironmentNamespace: *environmentNamespace,
	})
	assert.NoError(t, err)
	assert.Equal(t, pageSize, int64(len(res.Users)))
}

func TestListSegmentUsersCursor(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	segmentID := createSegment(ctx, t, client).Id
	userIDs := []string{newUserID(t), newUserID(t), newUserID(t), newUserID(t)}
	addSegmentUser(ctx, t, client, segmentID, userIDs, featureproto.SegmentUser_INCLUDED)
	var lastUsers []*featureproto.SegmentUser
	pageSize := int64(2)
	state := &wrappersproto.Int32Value{Value: int32(featureproto.SegmentUser_INCLUDED)}
	cursor := ""
	for i := 0; i < 3; i++ {
		res, err := client.ListSegmentUsers(ctx, &featureproto.ListSegmentUsersRequest{
			PageSize:             pageSize,
			Cursor:               cursor,
			SegmentId:            segmentID,
			State:                state,
			EnvironmentNamespace: *environmentNamespace,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Cursor)
		cursor = res.Cursor
		switch i {
		case 0:
			assert.Equal(t, int(pageSize), len(res.Users))
			copySegmentUsers(lastUsers, res.Users)
			break
		case 1:
			assert.Equal(t, int(pageSize), len(res.Users))
			if containsSegmentUser(lastUsers, res.Users) {
				t.Fatalf("Segment user from the last response was found in the actual response. Last response: %v, actual response: %v", lastUsers, res.Users)
			}
			break
		case 2:
			assert.Zero(t, len(res.Users))
			break
		}
	}
}

func TestListSegmentUsersWithoutState(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	segmentID := createSegment(ctx, t, client).Id
	userIDs := []string{newUserID(t)}
	addSegmentUser(ctx, t, client, segmentID, userIDs, featureproto.SegmentUser_INCLUDED)
	res := listSegmentUsers(ctx, t, client, segmentID, nil)
	assert.Equal(t, 1, len(res.Users))
	assert.Equal(t, segmentID, res.Users[0].SegmentId)
	assert.Equal(t, userIDs[0], res.Users[0].UserId)
	assert.Equal(t, featureproto.SegmentUser_INCLUDED, res.Users[0].State)
}

func TestBulkUploadAndDownloadSegmentUsers(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	segmentID := createSegment(ctx, t, client).Id
	uids := []string{newUserID(t), newUserID(t), newUserID(t)}
	sort.Strings(uids)
	userIDs := []byte(fmt.Sprintf("%s\n%s\n%s\n", uids[0], uids[1], uids[2]))
	uploadRes, err := client.BulkUploadSegmentUsers(ctx, &featureproto.BulkUploadSegmentUsersRequest{
		EnvironmentNamespace: *environmentNamespace,
		SegmentId:            segmentID,
		Command: &featureproto.BulkUploadSegmentUsersCommand{
			Data:  userIDs,
			State: featureproto.SegmentUser_INCLUDED,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, uploadRes)
	for i := 0; i < segmentUserRetryTimes; i++ {
		downloadRes, err := bulkDownloadSegmentUsers(t, client, segmentID)
		if err == nil {
			assert.Equal(t, string(userIDs), string(downloadRes.Data))
			break
		}
		if i == segmentUserRetryTimes-1 {
			t.Fatalf("SegmentUsers cannot be downloaded.")
		}
		time.Sleep(time.Second)
	}
}

func addSegmentUser(ctx context.Context, t *testing.T, client featureclient.Client, segmentID string, userIDs []string, state featureproto.SegmentUser_State) {
	t.Helper()
	req := &featureproto.AddSegmentUserRequest{
		Id: segmentID,
		Command: &featureproto.AddSegmentUserCommand{
			UserIds: userIDs,
			State:   state,
		},
		EnvironmentNamespace: *environmentNamespace,
	}
	res, err := client.AddSegmentUser(ctx, req)
	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func deleteSegmentUser(ctx context.Context, t *testing.T, client featureclient.Client, segmentID string, userIDs []string, state featureproto.SegmentUser_State) {
	req := &featureproto.DeleteSegmentUserRequest{
		Id: segmentID,
		Command: &featureproto.DeleteSegmentUserCommand{
			UserIds: userIDs,
			State:   state,
		},
		EnvironmentNamespace: *environmentNamespace,
	}
	res, err := client.DeleteSegmentUser(ctx, req)
	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func getSegmentUser(ctx context.Context, t *testing.T, client featureclient.Client, segmentID string, userID string, state featureproto.SegmentUser_State) *featureproto.SegmentUser {
	t.Helper()
	req := &featureproto.GetSegmentUserRequest{
		SegmentId:            segmentID,
		UserId:               userID,
		State:                state,
		EnvironmentNamespace: *environmentNamespace,
	}
	res, err := client.GetSegmentUser(ctx, req)
	assert.NoError(t, err)
	return res.User
}

func listSegmentUsers(ctx context.Context, t *testing.T, client featureclient.Client, segmentID string, state *wrappersproto.Int32Value) *featureproto.ListSegmentUsersResponse {
	t.Helper()
	req := &featureproto.ListSegmentUsersRequest{
		SegmentId:            segmentID,
		State:                state,
		EnvironmentNamespace: *environmentNamespace,
	}
	res, err := client.ListSegmentUsers(ctx, req)
	assert.NoError(t, err)
	return res
}

func copySegmentUsers(dst []*featureproto.SegmentUser, src []*featureproto.SegmentUser) {
	dst = make([]*featureproto.SegmentUser, 0, len(src))
	for _, s := range src {
		dst = append(dst, &featureproto.SegmentUser{
			Id:        s.Id,
			SegmentId: s.SegmentId,
			UserId:    s.UserId,
			State:     s.State,
			Deleted:   s.Deleted,
		})
	}
}

func containsSegmentUser(lastUsers []*featureproto.SegmentUser, actualUsers []*featureproto.SegmentUser) bool {
	for _, user := range lastUsers {
		for _, u := range actualUsers {
			if proto.Equal(user, u) {
				return true
			}
		}
	}
	return false
}

func bulkDownloadSegmentUsers(t *testing.T, client featureclient.Client, segmentID string) (*featureproto.BulkDownloadSegmentUsersResponse, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return client.BulkDownloadSegmentUsers(ctx, &featureproto.BulkDownloadSegmentUsersRequest{
		EnvironmentNamespace: *environmentNamespace,
		SegmentId:            segmentID,
		State:                featureproto.SegmentUser_INCLUDED,
	})
}

func newUserID(t *testing.T) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-user-id-%s", prefixID, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-user-id-%s", prefixID, newUUID(t))
}
