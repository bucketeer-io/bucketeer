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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

// overDefaultMessageSizeUsers is a number of segment users whose list is over
// 4MiB, the size of the largest message a gRPC client accepts by default.
const overDefaultMessageSizeUsers = 50000

// A SEGMENT group on a segment whose users do not fit in a gRPC message of
// the default size, and that no flag rule references. The group is started
// right after it is created. Its members are held out, and a server SDK gets
// the users of the segment from the REST endpoint it calls.
func TestServingSegmentOverDefaultMessageSize(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	tag := newTag(t)
	flag := createEnabledFeature(t, client, tag)
	refreshFeatureCache(t, tag, flag.Id)
	segment := createSegment(t, client)
	userIDs := make([]string, overDefaultMessageSizeUsers)
	for i := range userIDs {
		userIDs[i] = fmt.Sprintf("u%d", i)
	}
	uploadSegmentUsersAndWait(t, client, segment.Id, userIDs)
	member, nonMember := userIDs[len(userIDs)-1], newUserID(t)

	start, stop := servingWindowNow()
	group := createHoldoutGroup(t, client, groupSpec{
		assignment: featureproto.HoldoutGroup_SEGMENT, segmentID: segment.Id, startAt: start, stopAt: stop,
	}.request(t))
	addHoldoutFeatures(t, client, group.Id, target(flag, "C"))
	group = startHoldoutGroup(t, client, group.Id)

	running := waitForHoldoutServing(t, tag, member, flag.Id, variationID(flag, "C"))
	requireHoldoutEvaluation(t, running.Evaluations.Evaluations, flag.Id, variationID(flag, "C"), group.Id)
	other := grpcGetEvaluations(t, tag, nonMember)
	requireNormalEvaluation(t, other.Evaluations.Evaluations, flag.Id, variationID(flag, "A"), featureproto.Reason_DEFAULT)

	var segmentUsers *featureproto.SegmentUsers
	deadline := time.Now().Add(cachePropagationTimeout)
	for {
		status, body, err := restGetSegmentUsers()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status, "%.300s", body)
		resp := &gatewayproto.GetSegmentUsersResponse{}
		require.NoError(t, protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(body, resp))
		segmentUsers = findSegmentUsers(resp.SegmentUsers, segment.Id)
		if segmentUsers != nil {
			break
		}
		require.True(t, time.Now().Before(deadline), "the holdout segment never reached the server SDK over REST")
		time.Sleep(pollInterval)
	}
	assert.Len(t, segmentUsers.Users, len(userIDs))

	stopHoldoutGroup(t, client, group.Id)
	waitForHoldoutCleared(t, tag, member)
}

// restGetSegmentUsers is the first /get_segment_users request of a server
// SDK: it holds no segment yet. It returns the status and the body as they
// are, so that the caller can report a rejected response.
func restGetSegmentUsers() (int, []byte, error) {
	apiKey, err := os.ReadFile(holdoutEnv.serverKeyPath)
	if err != nil {
		return 0, nil, err
	}
	body, err := json.Marshal(map[string]interface{}{
		"segment_ids":  []string{},
		"requested_at": "0",
		"source_id":    "GO_SERVER",
		"sdk_version":  "v0.0.1-e2e-test",
	})
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf("https://%s/get_segment_users", *gatewayAddr), bytes.NewReader(body),
	)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Add("authorization", string(apiKey))
	req.Header.Add("Content-Type", "application/json")
	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			// The certificate of the e2e gateway may be expired or self-signed.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, data, nil
}
