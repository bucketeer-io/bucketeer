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

package experiment

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	btclient "github.com/bucketeer-io/bucketeer/v2/pkg/batch/client"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	prefixTestName = "e2e-test"
	timeout        = 60 * time.Second
	retryTimes     = 250
)

var (
	// FIXME: To avoid compiling the test many times, webGatewayAddr, webGatewayPort & apiKey has been also added here to prevent from getting: "flag provided but not defined" error during the test. These 3 are being use  in the Gateway test
	webGatewayAddr   = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort   = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert   = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath       = flag.String("api-key", "", "Client SDK API key for api-gateway")
	apiKeyServerPath = flag.String("api-key-server", "", "Server SDK API key for api-gateway")
	gatewayAddr      = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort      = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert      = flag.String("gateway-cert", "", "Gateway crt file")
	serviceTokenPath = flag.String("service-token", "", "Service token path")
	environmentID    = flag.String("environment-id", "", "Environment id")
	organizationID   = flag.String("organization-id", "", "Organization ID")
	testID           = flag.String("test-id", "", "test ID")
)

func TestCreateAndGetExperiment(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureID)
	goalIDs := createGoals(ctx, t, c, 2)
	startAt := time.Now()
	stopAt := startAt.Local().Add(time.Hour * 1)
	feature := getFeature(ctx, t, featureID)
	expected := createExperimentWithMultiGoals(ctx, t, c, featureID, feature.Variations[0].Id, goalIDs, startAt, stopAt)
	getResp, err := c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	actual := getResp.Experiment
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.FeatureId, actual.FeatureId)
	assert.Equal(t, expected.FeatureVersion, actual.FeatureVersion)
	assert.Equal(t, expected.BaseVariationId, actual.BaseVariationId)
	assert.Equal(t, expected.StartAt, actual.StartAt)
	assert.Equal(t, expected.StopAt, actual.StopAt)
	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.GoalIds, actual.GoalIds)
	assert.Equal(t, expected.Deleted, actual.Deleted)
	assert.Equal(t, expected.Archived, actual.Archived)
	assert.Equal(t, expected.Maintainer, actual.Maintainer)
	stopExperiment(ctx, t, c, expected.Id)
}

func TestListExperiments(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureID)
	feature := getFeature(ctx, t, featureID)
	goalIDs := createGoals(ctx, t, c, 5)
	startAt := time.Now()
	stopAt := startAt.Local().Add(time.Hour * 1)
	expectedExps := createExperimentsWithMultiGoals(ctx, t, c, featureID, feature.Variations[0].Id, goalIDs, startAt, stopAt, 5)
	sort.Slice(expectedExps, func(i, j int) bool {
		return len(expectedExps[i].Goals) < len(expectedExps[j].Goals)
	})

	getResp, err := c.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
		EnvironmentId: *environmentID,
		OrderBy:       experimentproto.ListExperimentsRequest_GOALS_COUNT,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, getResp.Experiments)

	for i := 1; i < len(getResp.Experiments); i++ {
		if len(getResp.Experiments[i].Goals) < len(getResp.Experiments[i-1].Goals) {
			t.Fatalf("Experiments are not sorted by goals count")
		}
	}
	for _, exp := range expectedExps {
		stopExperiment(ctx, t, c, exp.Id)
	}
}

func TestStopExperiment(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureID)
	goalIDs := createGoals(ctx, t, c, 1)
	feature := getFeature(ctx, t, featureID)
	startAt := time.Now()
	stopAt := startAt.Local().Add(time.Hour * 1)
	e := createExperimentWithMultiGoals(ctx, t, c, featureID, feature.Variations[0].Id, goalIDs, startAt, stopAt)
	if _, err := c.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
		Id: e.Id,
		Status: &experimentproto.UpdateExperimentRequest_UpdatedStatus{
			Status: experimentproto.Experiment_FORCE_STOPPED,
		},
		EnvironmentId: *environmentID,
	}); err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            e.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getResp.Experiment.StoppedAt == 0 {
		t.Fatal("Experiment was not stopped")
	}
	stopExperiment(ctx, t, c, e.Id)
}

func TestArchiveExperiment(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureID)
	goalIDs := createGoals(ctx, t, c, 1)
	feature := getFeature(ctx, t, featureID)
	startAt := time.Now()
	stopAt := startAt.Local().Add(time.Hour * 1)
	e := createExperimentWithMultiGoals(ctx, t, c, featureID, feature.Variations[0].Id, goalIDs, startAt, stopAt)
	if _, err := c.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
		Id:            e.Id,
		EnvironmentId: *environmentID,
		Archived:      wrapperspb.Bool(true),
	}); err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            e.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !getResp.Experiment.Archived {
		t.Fatal("Experiment was not archived")
	}
	stopExperiment(ctx, t, c, e.Id)
}

func TestDeleteExperiment(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureID)
	goalIDs := createGoals(ctx, t, c, 1)
	feature := getFeature(ctx, t, featureID)
	startAt := time.Now()
	stopAt := startAt.Local().Add(time.Hour * 1)
	e := createExperimentWithMultiGoals(ctx, t, c, featureID, feature.Variations[0].Id, goalIDs, startAt, stopAt)
	if _, err := c.DeleteExperiment(ctx, &experimentproto.DeleteExperimentRequest{
		Id:            e.Id,
		EnvironmentId: *environmentID,
	}); err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            e.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !getResp.Experiment.Deleted {
		t.Fatal("Experiment was not deleted")
	}
	stopExperiment(ctx, t, c, e.Id)
}

func TestUpdateExperiment(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureID)
	goalIDs := createGoals(ctx, t, c, 1)
	now := time.Now()
	feature := getFeature(ctx, t, featureID)
	startAt := time.Now()
	stopAt := startAt.Local().Add(time.Hour * 1)
	e := createExperimentWithMultiGoals(ctx, t, c, featureID, feature.Variations[0].Id, goalIDs, startAt, stopAt)
	startAt = now.Local().Add(time.Minute * 30)
	stopAt = now.Local().Add(time.Minute * 60)
	newName := fmt.Sprintf("%s-new-exp-name-%s", prefixTestName, newUUID(t))
	if _, err := c.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
		Id:            e.Id,
		Name:          wrapperspb.String(newName),
		StartAt:       wrapperspb.Int64(startAt.Unix()),
		StopAt:        wrapperspb.Int64(stopAt.Unix()),
		EnvironmentId: *environmentID,
	}); err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            e.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if startAt.Unix() != getResp.Experiment.StartAt {
		t.Fatalf("StartAt is not equal. Expected: %d, actual: %d", startAt.Unix(), getResp.Experiment.StartAt)
	}
	if stopAt.Unix() != getResp.Experiment.StopAt {
		t.Fatalf("StopAt is not equal. Expected: %d, actual: %d", stopAt.Unix(), getResp.Experiment.StopAt)
	}
	if newName != getResp.Experiment.Name {
		t.Fatalf("Name is not equal. Expected: %s, actual: %s", newName, getResp.Experiment.Name)
	}
	stopExperiment(ctx, t, c, e.Id)
}

func TestCreateAndGetGoal(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	goalID := createGoal(ctx, t, c)
	expectedName := fmt.Sprintf("%s-goal-name", goalID)
	expectedDescription := fmt.Sprintf("%s-goal-description", goalID)
	getResp, err := c.GetGoal(ctx, &experimentproto.GetGoalRequest{
		Id:            goalID,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	actual := getResp.Goal
	if goalID != actual.Id {
		t.Fatalf("Goal id is not equal. Expected: %v, actual: %v", goalID, actual.Id)
	}
	if expectedName != actual.Name {
		t.Fatalf("Goal name is not equal. Expected: %v, actual: %v", expectedName, actual.Name)
	}
	if expectedDescription != actual.Description {
		t.Fatalf("Goal description is not equal. Expected: %v, actual: %v", expectedDescription, actual.Description)
	}
	if actual.Deleted {
		t.Fatal("Goal deleted flag is true")
	}
}

func TestListGoalsCursor(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	createGoals(ctx, t, c, 2)
	expectedSize := 1
	listResp, err := c.ListGoals(ctx, &experimentproto.ListGoalsRequest{
		PageSize:      int64(expectedSize),
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if listResp.Cursor == "" {
		t.Fatal("Cursor is empty")
	}
	actualSize := len(listResp.Goals)
	if expectedSize != actualSize {
		t.Fatalf("Different sizes. Expected: %v, actual: %v", expectedSize, actualSize)
	}
	listResp, err = c.ListGoals(ctx, &experimentproto.ListGoalsRequest{
		PageSize:      int64(expectedSize),
		Cursor:        listResp.Cursor,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	actualSize = len(listResp.Goals)
	if expectedSize != actualSize {
		t.Fatalf("Different sizes. Expected: %v, actual: %v", expectedSize, actualSize)
	}
}

func TestListGoalsPageSize(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	createGoals(ctx, t, c, 3)
	expectedSize := 3
	listResp, err := c.ListGoals(ctx, &experimentproto.ListGoalsRequest{
		PageSize:      int64(expectedSize),
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	actualSize := len(listResp.Goals)
	if expectedSize != actualSize {
		t.Fatalf("Different sizes. Expected: %v, actual: %v", expectedSize, actualSize)
	}
}

func TestUpdateGoal(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	goalID := createGoal(ctx, t, c)
	expectedName := fmt.Sprintf("%s-goal-new-name", prefixTestName)
	expectedDescription := fmt.Sprintf("%s-goal-new-description", prefixTestName)
	_, err := c.UpdateGoal(ctx, &experimentproto.UpdateGoalRequest{
		Id:            goalID,
		Name:          wrapperspb.String(expectedName),
		Description:   wrapperspb.String(expectedDescription),
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetGoal(ctx, &experimentproto.GetGoalRequest{
		Id:            goalID,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	actual := getResp.Goal
	if goalID != actual.Id {
		t.Fatalf("Goal id is not equal. Expected: %v, actual: %v", goalID, actual.Id)
	}
	if expectedName != actual.Name {
		t.Fatalf("Goal name is not equal. Expected: %v, actual: %v", expectedName, actual.Name)
	}
	if expectedDescription != actual.Description {
		t.Fatalf("Goal description is not equal. Expected: %v, actual: %v", expectedDescription, actual.Description)
	}
	if actual.Deleted {
		t.Fatal("Goal deleted flag is true")
	}
}

func TestArchiveGoal(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	goalID := createGoal(ctx, t, c)
	_, err := c.UpdateGoal(ctx, &experimentproto.UpdateGoalRequest{
		Id:            goalID,
		Archived:      wrapperspb.Bool(true),
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetGoal(ctx, &experimentproto.GetGoalRequest{
		Id:            goalID,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !getResp.Goal.Archived {
		t.Fatal("Goal archived flag is false")
	}
}

func TestDeleteGoal(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	goalID := createGoal(ctx, t, c)
	_, err := c.DeleteGoal(ctx, &experimentproto.DeleteGoalRequest{
		Id:            goalID,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetGoal(ctx, &experimentproto.GetGoalRequest{
		Id:            goalID,
		EnvironmentId: *environmentID,
	})
	if err == nil {
		t.Fatal("Expected error when getting deleted goal, got nil")
	}
	assert.Contains(t, err.Error(), "rpc error: code = NotFound")
	assert.Contains(t, err.Error(), "goal not found")
}

func TestStatusUpdateFromWaitingToRunningAndForceStopped(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureID)
	goalIDs := createGoals(ctx, t, c, 1)
	startAt := time.Now()
	stopAt := startAt.Local().Add(time.Hour)
	feature := getFeature(ctx, t, featureID)
	expected := createExperimentWithMultiGoals(ctx, t, c, featureID, feature.Variations[0].Id, goalIDs, startAt, stopAt)
	resp, err := c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Experiment.Status != experimentproto.Experiment_WAITING {
		t.Fatalf("Experiment status is not waiting. actual: %d", resp.Experiment.Status)
	}

	_, err = c.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
		Status: &experimentproto.UpdateExperimentRequest_UpdatedStatus{
			Status: experimentproto.Experiment_RUNNING,
		},
	})
	if err != nil {
		t.Fatalf("Failed to update experiment status: %v", err)
	}
	resp, err = c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Experiment.Status != experimentproto.Experiment_RUNNING {
		t.Fatalf("Experiment status is not running. actual: %d", resp.Experiment.Status)
	}
	_, err = c.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
		Status: &experimentproto.UpdateExperimentRequest_UpdatedStatus{
			Status: experimentproto.Experiment_FORCE_STOPPED,
		},
	})
	if err != nil {
		t.Fatalf("Failed to update experiment status: %v", err)
	}
	resp, err = c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Experiment.Status != experimentproto.Experiment_FORCE_STOPPED {
		t.Fatalf("Experiment status is not force stopped. actual: %d", resp.Experiment.Status)
	}
	stopExperiment(ctx, t, c, expected.Id)
}

func TestStatusUpdateFromWaitingToRunning(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureID)
	goalIDs := createGoals(ctx, t, c, 1)
	startAt := time.Now()
	stopAt := startAt.Local().Add(time.Hour)
	feature := getFeature(ctx, t, featureID)
	expected := createExperimentWithMultiGoals(ctx, t, c, featureID, feature.Variations[0].Id, goalIDs, startAt, stopAt)
	resp, err := c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Experiment.Status != experimentproto.Experiment_WAITING {
		t.Fatalf("Experiment status is not waiting. actual: %d", resp.Experiment.Status)
	}
	for i := 0; i < retryTimes; i++ {
		resp, err = c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
			Id:            expected.Id,
			EnvironmentId: *environmentID,
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Experiment.Status == experimentproto.Experiment_RUNNING {
			break
		}
		if i == retryTimes-1 {
			t.Fatalf("retry timeout: %s", resp.Experiment.Name)
		}
		time.Sleep(time.Second)
	}
	stopExperiment(ctx, t, c, expected.Id)
}

func TestStatusUpdateFromRunningToStopped(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureID)
	goalIDs := createGoals(ctx, t, c, 1)
	now := time.Now()
	startAt := now.Local().Add(-4 * 24 * time.Hour)
	stopAt := now.Local().Add(-3 * 24 * time.Hour)
	feature := getFeature(ctx, t, featureID)
	expected := createExperimentWithMultiGoals(ctx, t, c, featureID, feature.Variations[0].Id, goalIDs, startAt, stopAt)
	resp, err := c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Experiment.Status != experimentproto.Experiment_WAITING {
		t.Fatalf("Experiment status is not waiting. actual: %d", resp.Experiment.Status)
	}
	if _, err = c.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
		Status: &experimentproto.UpdateExperimentRequest_UpdatedStatus{
			Status: experimentproto.Experiment_RUNNING,
		},
	}); err != nil {
		t.Fatal(err)
	}
	resp, err = c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
	})
	if resp.Experiment.Status != experimentproto.Experiment_RUNNING {
		t.Fatalf("Experiment status is not running. actual: %d", resp.Experiment.Status)
	}
	for i := 0; i < retryTimes; i++ {
		resp, err = c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
			Id:            expected.Id,
			EnvironmentId: *environmentID,
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Experiment.Status == experimentproto.Experiment_STOPPED {
			break
		}
		if i == retryTimes-1 {
			t.Fatalf("retry timeout: %s", resp.Experiment.Name)
		}
		time.Sleep(time.Second)
	}
	stopExperiment(ctx, t, c, expected.Id)
}

func TestStatusUpdateFromWaitingToStopped(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureID)
	goalIDs := createGoals(ctx, t, c, 1)
	now := time.Now()
	startAt := now.Local().Add(-4 * 24 * time.Hour)
	stopAt := now.Local().Add(-3 * 24 * time.Hour)
	feature := getFeature(ctx, t, featureID)
	expected := createExperimentWithMultiGoals(ctx, t, c, featureID, feature.Variations[0].Id, goalIDs, startAt, stopAt)
	resp, err := c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id:            expected.Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Experiment.Status != experimentproto.Experiment_WAITING {
		t.Fatalf("Experiment status is not waiting. actual: %d", resp.Experiment.Status)
	}
	for i := 0; i < retryTimes; i++ {
		resp, err = c.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
			Id:            expected.Id,
			EnvironmentId: *environmentID,
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Experiment.Status == experimentproto.Experiment_STOPPED {
			break
		}
		if i == retryTimes-1 {
			t.Fatalf("retry timeout: %s", resp.Experiment.Name)
		}
		time.Sleep(time.Second)
	}
	stopExperiment(ctx, t, c, expected.Id)
}

func TestCreateListGoals(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	goalID := createGoalID(t)
	createGoalResp, err := c.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
		EnvironmentId:  *environmentID,
		Id:             goalID,
		Name:           fmt.Sprintf("%s-goal-name", goalID),
		Description:    fmt.Sprintf("%s-goal-description", goalID),
		ConnectionType: experimentproto.Goal_EXPERIMENT,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, createGoalResp)

	listGoalsResp, err := c.ListGoals(ctx, &experimentproto.ListGoalsRequest{
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if listGoalsResp == nil || len(listGoalsResp.Goals) == 0 {
		t.Fatal("No goals")
	}

	var pbGoal *experimentproto.Goal
	for _, goal := range listGoalsResp.Goals {
		if goal.Id == createGoalResp.Goal.Id {
			pbGoal = goal
			break
		}
	}
	if pbGoal == nil {
		t.Fatalf("Goal not found: %s", createGoalResp.Goal.Id)
	}
	assert.Equal(t, createGoalResp.Goal.Id, pbGoal.Id)
	assert.Equal(t, createGoalResp.Goal.Name, pbGoal.Name)
	assert.Equal(t, createGoalResp.Goal.Description, pbGoal.Description)

	_, err = c.DeleteGoal(ctx, &experimentproto.DeleteGoalRequest{
		Id:            goalID,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateUpdateGoal(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newExperimentClient(t)
	defer c.Close()
	goalID := createGoalID(t)
	createGoalResp, err := c.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
		EnvironmentId:  *environmentID,
		Id:             goalID,
		Name:           fmt.Sprintf("%s-goal-name", goalID),
		Description:    fmt.Sprintf("%s-goal-description", goalID),
		ConnectionType: experimentproto.Goal_OPERATION,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, createGoalResp)

	expectedName := fmt.Sprintf("%s-goal-new-name-%s", prefixTestName, newUUID(t))
	expectedDescription := fmt.Sprintf("%s-goal-new-description-%s", prefixTestName, newUUID(t))
	updateGoalResp, err := c.UpdateGoal(ctx, &experimentproto.UpdateGoalRequest{
		EnvironmentId: *environmentID,
		Id:            goalID,
		Name:          wrapperspb.String(expectedName),
		Description:   wrapperspb.String(expectedDescription),
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, createGoalResp)

	assert.Equal(t, expectedName, updateGoalResp.Goal.Name)
	assert.Equal(t, expectedDescription, updateGoalResp.Goal.Description)

	getGoalResp, err := c.GetGoal(ctx, &experimentproto.GetGoalRequest{
		Id:            goalID,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, getGoalResp)

	assert.Equal(t, expectedName, getGoalResp.Goal.Name)
	assert.Equal(t, expectedDescription, getGoalResp.Goal.Description)

	_, err = c.DeleteGoal(ctx, &experimentproto.DeleteGoalRequest{
		Id:            goalID,
		EnvironmentId: *environmentID,
	})
}

func createGoal(ctx context.Context, t *testing.T, client experimentclient.Client) string {
	t.Helper()
	goalID := createGoalID(t)
	_, err := client.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
		Id:             goalID,
		Name:           fmt.Sprintf("%s-goal-name", goalID),
		Description:    fmt.Sprintf("%s-goal-description", goalID),
		ConnectionType: experimentproto.Goal_EXPERIMENT,
		EnvironmentId:  *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return goalID
}

func createGoals(ctx context.Context, t *testing.T, client experimentclient.Client, total int) []string {
	t.Helper()
	goalIDs := make([]string, 0, total)
	for i := 0; i < total; i++ {
		goalIDs = append(goalIDs, createGoal(ctx, t, client))
	}
	return goalIDs
}

func createExperimentsWithMultiGoals(
	ctx context.Context,
	t *testing.T,
	client experimentclient.Client,
	featureID, baseVariationID string,
	goalIDs []string,
	startAt, stopAt time.Time,
	total int,
) []*experimentproto.Experiment {
	e := []*experimentproto.Experiment{}
	for i := 0; i < total; i++ {
		e = append(e, createExperimentWithMultiGoals(ctx, t, client, featureID, baseVariationID, goalIDs, startAt, stopAt))
	}
	return e
}

func createExperimentWithMultiGoals(
	ctx context.Context,
	t *testing.T,
	client experimentclient.Client,
	featureID, baseVariationID string,
	goalIDs []string,
	startAt, stopAt time.Time,
) *experimentproto.Experiment {
	resp, err := client.CreateExperiment(ctx, &experimentproto.CreateExperimentRequest{
		FeatureId:       featureID,
		StartAt:         startAt.Unix(),
		StopAt:          stopAt.Unix(),
		GoalIds:         goalIDs,
		Name:            strings.Join(goalIDs, ","),
		BaseVariationId: baseVariationID,
		EnvironmentId:   *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.Experiment
}

func newExperimentClient(t *testing.T) experimentclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := experimentclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create experiment client:", err)
	}
	return client
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func getFeature(ctx context.Context, t *testing.T, featureID string) *featureproto.Feature {
	t.Helper()
	client := newFeatureClient(t)
	defer client.Close()
	req := &featureproto.GetFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
	}
	resp, err := client.GetFeature(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	return resp.Feature
}

func createFeature(ctx context.Context, t *testing.T, featureID string) {
	t.Helper()
	client := newFeatureClient(t)
	defer client.Close()
	createReq := &featureproto.CreateFeatureRequest{
		EnvironmentId: *environmentID,
		Id:            featureID,
		Name:          featureID,
		Description:   "e2e-test-gateway-feature-description",
		Variations: []*featureproto.Variation{
			{
				Value:       "A",
				Name:        "Variation A",
				Description: "Thing does A",
			},
			{
				Value:       "B",
				Name:        "Variation B",
				Description: "Thing does B",
			},
		},
		Tags: []string{
			"e2e-test-tag-1",
			"e2e-test-tag-2",
			"e2e-test-tag-3",
		},
		DefaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
		DefaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
	}
	if _, err := client.CreateFeature(ctx, createReq); err != nil {
		t.Fatal(err)
	}
	enableFeature(t, featureID, client)
}

func newFeatureClient(t *testing.T) featureclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	featureClient, err := featureclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create feature client:", err)
	}
	return featureClient
}

func enableFeature(t *testing.T, featureID string, client featureclient.Client) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
		Id:            featureID,
		Enabled:       wrapperspb.Bool(true),
		EnvironmentId: *environmentID,
	}); err != nil {
		t.Fatalf("Failed to enable feature id: %s. Error: %v", featureID, err)
	}
}

func createFeatureID(t *testing.T) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-feature-id-%s", prefixTestName, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-feature-id-%s", prefixTestName, newUUID(t))
}

func createGoalID(t *testing.T) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-goal-id-%s", prefixTestName, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-goal-id-%s", prefixTestName, newUUID(t))
}

// This helper tries to stop the running experiments
// that are finished testing and waiting for deletion.
// This will improve the load on the http-stan while analysing the other experiments
// speeding up and improve timeout flaky tests.
// Since this is optional, it will ignore any errors.
func stopExperiment(
	ctx context.Context,
	t *testing.T,
	client experimentclient.Client,
	id string,
) {
	t.Helper()
	_, err := client.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
		EnvironmentId: *environmentID,
		Id:            id,
		Status: &experimentproto.UpdateExperimentRequest_UpdatedStatus{
			Status: experimentproto.Experiment_FORCE_STOPPED,
		},
	})
	if err != nil {
		// Ignore
		return
	}
	// Update experiment cache
	batchClient := newBatchClient(t)
	defer batchClient.Close()
	numRetries := 3
	for i := 0; i < numRetries; i++ {
		_, err = batchClient.ExecuteBatchJob(
			ctx,
			&btproto.BatchJobRequest{Job: btproto.BatchJob_ExperimentCacher})
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		if st.Code() != codes.Unavailable {
			return
		}
		fmt.Printf("Failed to execute experiment cacher batch (Called by stopExperiment). Error code: %d. Retrying in 5 seconds.\n", st.Code())
		time.Sleep(5 * time.Second)
	}
}

func newBatchClient(t *testing.T) btclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := btclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create batch client:", err)
	}
	return client
}
