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

package feature

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	// The wait must cover a full redelivery cycle of the redis-stream pubsub used
	// in the dev container: a nacked message is only reclaimed by the puller's
	// recovery loop (30s ticker) after it has been idle for redisIdleTime, plus
	// the persister's flush interval. 60 retries x 2s = 120s of polling.
	segmentUserRetryTimes = 60
	// The default 60s package timeout is too short to cover the polling budget above.
	segmentUserTimeout = 3 * time.Minute
)

func TestListSegmentUsersPageSize(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), segmentUserTimeout)
	defer cancel()
	client := newFeatureClient(t)
	segmentID := createSegment(ctx, t, client).Id
	userIDs := []string{newUserID(t), newUserID(t)}
	uploadSegmentUsers(ctx, t, client, segmentID, userIDs, featureproto.SegmentUser_INCLUDED)
	waitForSegmentUsers(ctx, t, client, segmentID, len(userIDs), &wrapperspb.Int32Value{Value: int32(featureproto.SegmentUser_INCLUDED)})
	pageSize := int64(1)
	res, err := client.ListSegmentUsers(ctx, &featureproto.ListSegmentUsersRequest{
		PageSize:      pageSize,
		SegmentId:     segmentID,
		State:         &wrapperspb.Int32Value{Value: int32(featureproto.SegmentUser_INCLUDED)},
		EnvironmentId: *environmentID,
	})
	assert.NoError(t, err)
	assert.Equal(t, pageSize, int64(len(res.Users)))
}

func TestListSegmentUsersCursor(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), segmentUserTimeout)
	defer cancel()
	client := newFeatureClient(t)
	segmentID := createSegment(ctx, t, client).Id
	userIDs := []string{newUserID(t), newUserID(t), newUserID(t), newUserID(t)}
	uploadSegmentUsers(ctx, t, client, segmentID, userIDs, featureproto.SegmentUser_INCLUDED)
	state := &wrapperspb.Int32Value{Value: int32(featureproto.SegmentUser_INCLUDED)}
	waitForSegmentUsers(ctx, t, client, segmentID, len(userIDs), state)
	var lastUsers []*featureproto.SegmentUser
	pageSize := int64(2)
	cursor := ""
	for i := 0; i < 3; i++ {
		res, err := client.ListSegmentUsers(ctx, &featureproto.ListSegmentUsersRequest{
			PageSize:      pageSize,
			Cursor:        cursor,
			SegmentId:     segmentID,
			State:         state,
			EnvironmentId: *environmentID,
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
	ctx, cancel := context.WithTimeout(context.Background(), segmentUserTimeout)
	defer cancel()
	client := newFeatureClient(t)
	segmentID := createSegment(ctx, t, client).Id
	userIDs := []string{newUserID(t)}
	uploadSegmentUsers(ctx, t, client, segmentID, userIDs, featureproto.SegmentUser_INCLUDED)
	waitForSegmentUsers(ctx, t, client, segmentID, len(userIDs), nil)
	res := listSegmentUsers(ctx, t, client, segmentID, nil)
	assert.Equal(t, 1, len(res.Users))
	assert.Equal(t, segmentID, res.Users[0].SegmentId)
	assert.Equal(t, userIDs[0], res.Users[0].UserId)
	assert.Equal(t, featureproto.SegmentUser_INCLUDED, res.Users[0].State)
}

func TestBulkUploadAndDownloadSegmentUsers(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), segmentUserTimeout)
	defer cancel()
	client := newFeatureClient(t)
	segmentID := createSegment(ctx, t, client).Id
	uids := []string{newUserID(t), newUserID(t), newUserID(t)}
	sort.Strings(uids)
	userIDs := []byte(fmt.Sprintf("%s\n%s\n%s\n", uids[0], uids[1], uids[2]))
	uploadRes, err := client.BulkUploadSegmentUsers(ctx, &featureproto.BulkUploadSegmentUsersRequest{
		EnvironmentId: *environmentID,
		SegmentId:     segmentID,
		Data:          userIDs,
		State:         featureproto.SegmentUser_INCLUDED,
	})
	assert.NoError(t, err)
	assert.NotNil(t, uploadRes)
	// Wait for the background upload to complete by verifying users exist
	waitForSegmentUsers(ctx, t, client, segmentID, 3, &wrapperspb.Int32Value{Value: int32(featureproto.SegmentUser_INCLUDED)})
	// BulkDownloadSegmentUsers requires the segment status to be SUCEEDED,
	// which the persister updates after writing the users, so wait for it too.
	waitForSegmentStatus(ctx, t, client, segmentID, featureproto.Segment_SUCEEDED)
	// Now download and verify the data
	downloadRes, err := bulkDownloadSegmentUsers(t, client, segmentID)
	require.NoError(t, err)
	assert.Equal(t, string(userIDs), string(downloadRes.Data))
}

func listSegmentUsers(ctx context.Context, t *testing.T, client featureclient.Client, segmentID string, state *wrapperspb.Int32Value) *featureproto.ListSegmentUsersResponse {
	t.Helper()
	req := &featureproto.ListSegmentUsersRequest{
		SegmentId:     segmentID,
		State:         state,
		EnvironmentId: *environmentID,
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
		EnvironmentId: *environmentID,
		SegmentId:     segmentID,
		State:         featureproto.SegmentUser_INCLUDED,
	})
}

func uploadSegmentUsers(
	ctx context.Context,
	t *testing.T,
	client featureclient.Client,
	segmentID string,
	userIDs []string,
	state featureproto.SegmentUser_State,
) {
	t.Helper()
	data := []byte(strings.Join(userIDs, "\n") + "\n")
	_, err := client.BulkUploadSegmentUsers(ctx, &featureproto.BulkUploadSegmentUsersRequest{
		EnvironmentId: *environmentID,
		SegmentId:     segmentID,
		Data:          data,
		State:         state,
	})
	assert.NoError(t, err)
}

func waitForSegmentUsers(
	ctx context.Context,
	t *testing.T,
	client featureclient.Client,
	segmentID string,
	expectedSize int,
	state *wrapperspb.Int32Value,
) {
	t.Helper()
	req := &featureproto.ListSegmentUsersRequest{
		SegmentId:     segmentID,
		State:         state,
		EnvironmentId: *environmentID,
	}
	var lastCount int
	var lastErr error
	for i := 0; i < segmentUserRetryTimes; i++ {
		if err := ctx.Err(); err != nil {
			t.Fatalf("waitForSegmentUsers: context done: %v (last count: %d/%d, last error: %v)",
				err, lastCount, expectedSize, lastErr)
		}
		res, err := client.ListSegmentUsers(ctx, req)
		lastErr = err
		if err == nil && res != nil {
			lastCount = len(res.Users)
			if lastCount >= expectedSize {
				return
			}
		}
		time.Sleep(2 * time.Second)
	}
	t.Fatalf("segment users not ready after %d attempts: segmentID: %s, last count: %d/%d, last error: %v",
		segmentUserRetryTimes, segmentID, lastCount, expectedSize, lastErr)
}

func waitForSegmentStatus(
	ctx context.Context,
	t *testing.T,
	client featureclient.Client,
	segmentID string,
	status featureproto.Segment_Status,
) {
	t.Helper()
	req := &featureproto.GetSegmentRequest{
		Id:            segmentID,
		EnvironmentId: *environmentID,
	}
	var lastStatus featureproto.Segment_Status
	var lastErr error
	for i := 0; i < segmentUserRetryTimes; i++ {
		if err := ctx.Err(); err != nil {
			t.Fatalf("waitForSegmentStatus: context done: %v (last status: %v, last error: %v)",
				err, lastStatus, lastErr)
		}
		res, err := client.GetSegment(ctx, req)
		lastErr = err
		if err == nil && res.Segment != nil {
			lastStatus = res.Segment.Status
			if lastStatus == status {
				return
			}
		}
		time.Sleep(2 * time.Second)
	}
	t.Fatalf("segment status did not become %v after %d attempts: segmentID: %s, last status: %v, last error: %v",
		status, segmentUserRetryTimes, segmentID, lastStatus, lastErr)
}

func newUserID(t *testing.T) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-user-id-%s", prefixID, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-user-id-%s", prefixID, newUUID(t))
}
