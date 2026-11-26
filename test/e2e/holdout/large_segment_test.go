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

package holdout

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// largeSegmentUsers is a number of users whose list does not fit in one gRPC
// response of the default size (4 MiB), which is what the cache refresher
// reads to cache the users of a segment. Such a segment is cached by the
// SegmentUserCacher batch job only.
const largeSegmentUsers = 50000

// A SEGMENT group holds out the members of a large segment that no flag rule
// references. The batch job caches the users of the segment while the group
// is still WAITING, so they are held out as soon as the group starts.
func TestServingLargeSegment(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	tag := newTag(t)
	flag := createEnabledFeature(t, client, tag)
	refreshFeatureCache(t, tag, flag.Id)
	segment := createSegment(t, client)
	userIDs := make([]string, largeSegmentUsers)
	for i := range userIDs {
		userIDs[i] = fmt.Sprintf("u%d", i)
	}
	uploadManySegmentUsers(t, client, segment.Id, userIDs)
	member, nonMember := userIDs[len(userIDs)-1], newUserID(t)

	start, stop := servingWindowNow()
	group := createHoldoutGroup(t, client, groupSpec{
		assignment: featureproto.HoldoutGroup_SEGMENT, segmentID: segment.Id, startAt: start, stopAt: stop,
	}.request(t))
	addHoldoutFeatures(t, client, group.Id, target(flag, "C"))
	runBatchJob(t, btproto.BatchJob_SegmentUserCacher)
	group = startHoldoutGroup(t, client, group.Id)

	running := waitForHoldoutServing(t, tag, member, flag.Id, variationID(flag, "C"))
	requireHoldoutEvaluation(t, running.Evaluations.Evaluations, flag.Id, variationID(flag, "C"), group.Id)
	other := grpcGetEvaluations(t, tag, nonMember)
	requireNormalEvaluation(t, other.Evaluations.Evaluations, flag.Id, variationID(flag, "A"), featureproto.Reason_DEFAULT)

	stopHoldoutGroup(t, client, group.Id)
	waitForHoldoutCleared(t, tag, member)
}

// uploadManySegmentUsers uploads included users and waits until the upload
// has been stored. It does not list the users as uploadSegmentUsersAndWait
// does, because a list of this many users does not fit in one response.
func uploadManySegmentUsers(t *testing.T, client featureclient.Client, segmentID string, userIDs []string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.BulkUploadSegmentUsers(ctx, &featureproto.BulkUploadSegmentUsersRequest{
		EnvironmentId: holdoutEnv.id,
		SegmentId:     segmentID,
		Data:          []byte(strings.Join(userIDs, "\n") + "\n"),
		State:         featureproto.SegmentUser_INCLUDED,
	})
	require.NoError(t, err)
	// The upload sets the segment to UPLOADING before it returns, and the
	// persister stores the number of users of this upload with SUCEEDED.
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		res, err := client.GetSegment(ctx, &featureproto.GetSegmentRequest{
			Id: segmentID, EnvironmentId: holdoutEnv.id,
		})
		return err == nil &&
			res.Segment.Status == featureproto.Segment_SUCEEDED &&
			res.Segment.IncludedUserCount == int64(len(userIDs))
	}, segmentUploadTimeout, 2*time.Second, "segment upload of %d users not stored", len(userIDs))
}
