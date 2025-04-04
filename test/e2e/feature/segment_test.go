// Copyright 2025 The Bucketeer Authors.
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
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	"github.com/bucketeer-io/bucketeer/test/util"
)

const (
	prefixSegment = "e2e-test"
)

func TestCreateListSegmentNoCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	createResp, err := client.CreateSegment(ctx, &featureproto.CreateSegmentRequest{
		EnvironmentId: *environmentID,
		Name:          newSegmentName(t),
		Description:   fmt.Sprintf("%s-description", prefixSegment),
	})
	assert.NoError(t, err)

	listResp, err := client.ListSegments(ctx, &featureproto.ListSegmentsRequest{
		EnvironmentId: *environmentID,
	})
	assert.NoError(t, err)
	assert.NotNil(t, listResp)

	var segment *featureproto.Segment
	for _, s := range listResp.Segments {
		if s.Id == createResp.Segment.Id {
			assert.Equal(t, createResp.Segment.Name, s.Name)
			segment = s
			break
		}
	}
	if segment == nil {
		t.Fatalf("segment not found in list response")
	}

	// delete segment
	_, err = client.DeleteSegment(ctx, &featureproto.DeleteSegmentRequest{
		Id:            segment.Id,
		EnvironmentId: *environmentID,
	})
	assert.NoError(t, err)
}

func TestCreateUpdateNoCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	createResp, err := client.CreateSegment(ctx, &featureproto.CreateSegmentRequest{
		EnvironmentId: *environmentID,
		Name:          newSegmentName(t),
		Description:   fmt.Sprintf("%s-description", prefixSegment),
	})
	assert.NoError(t, err)
	assert.NotNil(t, createResp)
	segment := createResp.Segment

	updateResp, err := client.UpdateSegment(ctx, &featureproto.UpdateSegmentRequest{
		Id:            createResp.Segment.Id,
		EnvironmentId: *environmentID,
		Name:          wrapperspb.String(fmt.Sprintf("%s-update", segment.Name)),
	})
	assert.NoError(t, err)
	if updateResp == nil {
		t.Fatalf("update response is nil")
	}
	assert.Equal(t, fmt.Sprintf("%s-update", segment.Name), updateResp.Segment.Name)

	// delete segment
	_, err = client.DeleteSegment(ctx, &featureproto.DeleteSegmentRequest{
		Id:            segment.Id,
		EnvironmentId: *environmentID,
	})
	assert.NoError(t, err)
}

func TestCreateSegment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	cmd := &featureproto.CreateSegmentCommand{
		Name:        newSegmentName(t),
		Description: fmt.Sprintf("%s-description", prefixSegment),
	}
	req := &featureproto.CreateSegmentRequest{
		Command:       cmd,
		EnvironmentId: *environmentID,
	}
	res, err := client.CreateSegment(ctx, req)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Segment.Id)
	assert.Equal(t, cmd.Name, res.Segment.Name)
	assert.Equal(t, cmd.Description, res.Segment.Description)
	assert.Zero(t, res.Segment.Rules)
	assert.NotZero(t, res.Segment.CreatedAt)
	assert.Zero(t, res.Segment.UpdatedAt)
	assert.Equal(t, int64(1), res.Segment.Version)
	assert.Equal(t, false, res.Segment.Deleted)
}

func TestGetSegment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	expected := createSegment(ctx, t, client)
	actual := getSegment(ctx, t, client, expected.Id)
	if !proto.Equal(expected, actual) {
		t.Fatalf("Different segments. Expected: %v, actual: %v", expected, actual)
	}
}

func TestGetUsedSegment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	segment := createSegment(ctx, t, client)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	feature := getFeature(t, featureID, client)
	rule := newFixedStrategyRuleWithSegment(feature.Variations[0].Id, segment.Id)
	addCmd, err := util.MarshalCommand(&featureproto.AddRuleCommand{Rule: rule})
	require.NoError(t, err)
	updateFeatureTargeting(t, client, addCmd, featureID)
	feature = getFeature(t, featureID, client)
	actual := getSegment(ctx, t, client, segment.Id)
	if !proto.Equal(feature, actual.Features[0]) {
		t.Fatalf("Different feature. Expected: %v actual: %v", feature, actual.Features[0])
	}
}

func TestChangeSegmentName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	seg := createSegment(ctx, t, client)
	id := seg.Id
	name := seg.Name
	cmd := &featureproto.ChangeSegmentNameCommand{
		Name: fmt.Sprintf("%s-change-name", prefixSegment),
	}
	cmdChange, err := ptypes.MarshalAny(cmd)
	assert.NoError(t, err)
	res, err := client.UpdateSegment(
		ctx,
		&featureproto.UpdateSegmentRequest{
			Id: id,
			Commands: []*featureproto.Command{
				{Command: cmdChange},
			},
			EnvironmentId: *environmentID,
		},
	)
	assert.NotNil(t, res)
	assert.NoError(t, err)
	segment := getSegment(ctx, t, client, id)
	assert.Equal(t, cmd.Name, segment.Name)

	// After confirming that the name has changed correctly,
	// We must change back the original name, so this e2e test
	// can be deleted correctly when running the delete e2e data workflow
	cmd = &featureproto.ChangeSegmentNameCommand{
		Name: name,
	}
	cmdChange, err = ptypes.MarshalAny(cmd)
	assert.NoError(t, err)
	res, err = client.UpdateSegment(
		ctx,
		&featureproto.UpdateSegmentRequest{
			Id: id,
			Commands: []*featureproto.Command{
				{Command: cmdChange},
			},
			EnvironmentId: *environmentID,
		},
	)
	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestChangeSegmentDescription(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	id := createSegment(ctx, t, client).Id
	cmd := &featureproto.ChangeSegmentDescriptionCommand{
		Description: fmt.Sprintf("%s-change-description", prefixSegment),
	}
	cmdChange, err := ptypes.MarshalAny(cmd)
	assert.NoError(t, err)
	res, err := client.UpdateSegment(
		ctx,
		&featureproto.UpdateSegmentRequest{
			Id: id,
			Commands: []*featureproto.Command{
				{Command: cmdChange},
			},
			EnvironmentId: *environmentID,
		},
	)
	assert.NotNil(t, res)
	assert.NoError(t, err)
	segment := getSegment(ctx, t, client, id)
	assert.Equal(t, cmd.Description, segment.Description)
}

func TestDeleteSegment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	id := createSegment(ctx, t, client).Id
	res, err := client.DeleteSegment(
		ctx,
		&featureproto.DeleteSegmentRequest{
			Id:            id,
			Command:       &featureproto.DeleteSegmentCommand{},
			EnvironmentId: *environmentID,
		},
	)
	assert.NotNil(t, res)
	assert.NoError(t, err)
	segment := getSegment(ctx, t, client, id)
	assert.Equal(t, true, segment.Deleted)
}

func TestListSegmentsPageSize(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	for i := 0; i < 2; i++ {
		createSegment(ctx, t, client)
	}
	pageSize := int64(1)
	res, err := client.ListSegments(ctx, &featureproto.ListSegmentsRequest{
		PageSize:      pageSize,
		EnvironmentId: *environmentID,
	})
	assert.NoError(t, err)
	assert.Equal(t, pageSize, int64(len(res.Segments)))
}

func TestListSegmentsCursor(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	for i := 0; i < 4; i++ {
		createSegment(ctx, t, client)
	}
	pageSize := int64(2)
	res, err := client.ListSegments(ctx, &featureproto.ListSegmentsRequest{
		PageSize:      pageSize,
		EnvironmentId: *environmentID,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Cursor)
	resCursor, err := client.ListSegments(
		ctx,
		&featureproto.ListSegmentsRequest{
			PageSize:      pageSize,
			Cursor:        res.Cursor,
			EnvironmentId: *environmentID,
		},
	)
	assert.NoError(t, err)
	segmentsSize := len(res.Segments)
	assert.Equal(t, segmentsSize, len(resCursor.Segments))
	for i := 0; i < segmentsSize; i++ {
		if proto.Equal(res.Segments[i], resCursor.Segments[i]) {
			t.Fatalf("Equal segments. Expected: %v, actual: %v", res.Segments, resCursor.Segments)
		}
	}
}

func getSegment(ctx context.Context, t *testing.T, client featureclient.Client, id string) *featureproto.Segment {
	t.Helper()
	req := &featureproto.GetSegmentRequest{
		Id:            id,
		EnvironmentId: *environmentID,
	}
	res, err := client.GetSegment(ctx, req)
	assert.NoError(t, err)
	return res.Segment
}

func createSegment(ctx context.Context, t *testing.T, client featureclient.Client) *featureproto.Segment {
	t.Helper()
	cmd := &featureproto.CreateSegmentCommand{
		Name:        newSegmentName(t),
		Description: fmt.Sprintf("%s-%s", "description", prefixSegment),
	}
	req := &featureproto.CreateSegmentRequest{
		Command:       cmd,
		EnvironmentId: *environmentID,
	}
	res, err := client.CreateSegment(ctx, req)
	assert.NoError(t, err)
	return res.Segment
}

func newSegmentName(t *testing.T) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-name-%s", prefixSegment, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-name-%s", prefixSegment, newUUID(t))
}
