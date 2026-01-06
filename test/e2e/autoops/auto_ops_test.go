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

package autoops

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	gwapi "github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	gatewayclient "github.com/bucketeer-io/bucketeer/v2/pkg/api/client"
	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	btclient "github.com/bucketeer-io/bucketeer/v2/pkg/batch/client"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
	"github.com/bucketeer-io/bucketeer/v2/test/e2e/util"
)

const (
	goalEventType eventType = iota + 1 // eventType starts from 1 for validation.
	evaluationEventType
	metricsEventType
	prefixTestName   = "e2e-test"
	retryTimes       = 30
	timeout          = 2 * time.Minute
	prefixID         = "e2e-test"
	version          = "/v1"
	service          = "/gateway"
	eventsAPI        = "/events"
	authorizationKey = "authorization"
)

var (
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

type eventType int

type event struct {
	ID            string          `json:"id,omitempty"`
	Event         json.RawMessage `json:"event,omitempty"`
	EnvironmentId string          `json:"environment_id,omitempty"`
	Type          eventType       `json:"type,omitempty"`
}

type successResponse struct {
	Data json.RawMessage `json:"data"`
}

type registerEventsRequest struct {
	Events []event `json:"events,omitempty"`
}

type registerEventsResponse struct {
	Errors map[string]*gatewayproto.RegisterEventsResponse_Error `json:"errors,omitempty"`
}

func TestCreateAndListAutoOpsRule(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	goalID := createGoal(ctx, t, experimentClient)
	clause := createOpsEventRateClause(t, feature.Variations[0].Id, goalID)
	createAutoOpsRule(
		ctx,
		t,
		autoOpsClient,
		featureID,
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{clause},
		nil,
	)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	actual := autoOpsRules[0]
	if actual.Id == "" {
		t.Fatal("id is empty")
	}
	if actual.FeatureId != featureID {
		t.Fatalf("different feature ID, expected: %v, actual: %v", featureID, actual.FeatureId)
	}
	if actual.OpsType != autoopsproto.OpsType_EVENT_RATE {
		t.Fatalf("different ops type, expected: %v, actual: %v", autoopsproto.OpsType_EVENT_RATE, actual.OpsType)
	}
	if actual.AutoOpsStatus != autoopsproto.AutoOpsStatus_WAITING {
		t.Fatalf("different auto ops status, expected: %v, actual: %v", autoopsproto.AutoOpsStatus_WAITING, actual.AutoOpsStatus)
	}
	oerc := unmarshalOpsEventRateClause(t, actual.Clauses[0])
	if oerc.VariationId != feature.Variations[0].Id {
		t.Fatalf("different variation id, expected: %v, actual: %v", feature.Variations[0].Id, oerc.VariationId)
	}
	if oerc.GoalId != goalID {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.MinCount != int64(5) {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.ThreadsholdRate != float64(0.5) {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.Operator != autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.ActionType != autoopsproto.ActionType_DISABLE {
		t.Fatalf("different action type, expected: %v, actual: %v", "gid", oerc.ActionType)
	}
}

func TestCreateAndListAutoOpsRuleNoCommand(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	goalID := createGoal(ctx, t, experimentClient)
	clause := createOpsEventRateClause(t, feature.Variations[0].Id, goalID)
	createAutoOpsRuleNoCommand(
		ctx,
		t,
		autoOpsClient,
		featureID,
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{clause},
		nil,
	)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	actual := autoOpsRules[0]
	if actual.Id == "" {
		t.Fatal("id is empty")
	}
	if actual.FeatureId != featureID {
		t.Fatalf("different feature ID, expected: %v, actual: %v", featureID, actual.FeatureId)
	}
	if actual.OpsType != autoopsproto.OpsType_EVENT_RATE {
		t.Fatalf("different ops type, expected: %v, actual: %v", autoopsproto.OpsType_EVENT_RATE, actual.OpsType)
	}
	if actual.AutoOpsStatus != autoopsproto.AutoOpsStatus_WAITING {
		t.Fatalf("different auto ops status, expected: %v, actual: %v", autoopsproto.AutoOpsStatus_WAITING, actual.AutoOpsStatus)
	}
	oerc := unmarshalOpsEventRateClause(t, actual.Clauses[0])
	if oerc.VariationId != feature.Variations[0].Id {
		t.Fatalf("different variation id, expected: %v, actual: %v", feature.Variations[0].Id, oerc.VariationId)
	}
	if oerc.GoalId != goalID {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.MinCount != int64(5) {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.ThreadsholdRate != float64(0.5) {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.Operator != autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.ActionType != autoopsproto.ActionType_DISABLE {
		t.Fatalf("different action type, expected: %v, actual: %v", "gid", oerc.ActionType)
	}
}

func TestCreateAndListAutoOpsRuleForMultiSchedule(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	clauses := createDatetimeClausesWithActionType(t, 2)
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, autoopsproto.OpsType_SCHEDULE, nil, clauses)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	actual := autoOpsRules[0]
	if actual.Id == "" {
		t.Fatal("id is empty")
	}
	if actual.FeatureId != featureID {
		t.Fatalf("different feature ID, expected: %v, actual: %v", featureID, actual.FeatureId)
	}
	if actual.OpsType != autoopsproto.OpsType_SCHEDULE {
		t.Fatalf("different ops type, expected: %v, actual: %v", autoopsproto.OpsType_SCHEDULE, actual.OpsType)
	}
	if actual.AutoOpsStatus != autoopsproto.AutoOpsStatus_WAITING {
		t.Fatalf("different auto ops status, expected: %v, actual: %v", autoopsproto.AutoOpsStatus_WAITING, actual.AutoOpsStatus)
	}
	if len(actual.Clauses) != 2 {
		t.Fatalf("different clauses length, expected: %v, actual: %v", 2, len(actual.Clauses))
	}

	actualClause1 := actual.Clauses[0]
	oerc1 := unmarshalDatetimeClause(t, actualClause1)
	if oerc1.ActionType != autoopsproto.ActionType_DISABLE {
		t.Fatalf("different dateClause1 action type, expected: %v, actual: %v", autoopsproto.ActionType_DISABLE, actualClause1.ActionType)
	}

	actualClause2 := actual.Clauses[1]
	oerc2 := unmarshalDatetimeClause(t, actualClause2)
	if oerc2.ActionType != autoopsproto.ActionType_ENABLE {
		t.Fatalf("different dateClause2 action type, expected: %v, actual: %v", autoopsproto.ActionType_ENABLE, actualClause2.ActionType)
	}
}

func TestGetAutoOpsRule(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	goalID := createGoal(ctx, t, experimentClient)
	clause := createOpsEventRateClause(t, feature.Variations[0].Id, goalID)
	createAutoOpsRule(
		ctx,
		t,
		autoOpsClient,
		featureID,
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{clause},
		nil,
	)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	actual := getAutoOpsRules(t, autoOpsRules[0].Id)
	if actual.Id == "" {
		t.Fatal("id is empty")
	}
	if actual.FeatureId != featureID {
		t.Fatalf("different feature ID, expected: %v, actual: %v", featureID, actual.FeatureId)
	}
	if actual.OpsType != autoopsproto.OpsType_EVENT_RATE {
		t.Fatalf("different ops type, expected: %v, actual: %v", autoopsproto.OpsType_EVENT_RATE, actual.OpsType)
	}
	oerc := unmarshalOpsEventRateClause(t, actual.Clauses[0])
	if oerc.VariationId != feature.Variations[0].Id {
		t.Fatalf("different variation id, expected: %v, actual: %v", feature.Variations[0].Id, oerc.VariationId)
	}
	if oerc.GoalId != goalID {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.MinCount != int64(5) {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.ThreadsholdRate != float64(0.5) {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.Operator != autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL {
		t.Fatalf("different goal id, expected: %v, actual: %v", "gid", oerc.GoalId)
	}
	if oerc.ActionType != autoopsproto.ActionType_DISABLE {
		t.Fatalf("different action type, expected: %v, actual: %v", "gid", oerc.ActionType)
	}
}

func TestDeleteAutoOpsRule(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	goalID := createGoal(ctx, t, experimentClient)
	clause := createOpsEventRateClause(t, feature.Variations[0].Id, goalID)
	createAutoOpsRule(
		ctx,
		t,
		autoOpsClient,
		featureID,
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{clause},
		nil,
	)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	deleteAutoOpsRules(t, autoOpsClient, autoOpsRules[0].Id)
	resp, err := autoOpsClient.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentId: *environmentID,
		Id:            autoOpsRules[0].Id,
	})
	if resp != nil {
		t.Fatal("autoOpsRules is not deleted")
	}
	if err == nil {
		t.Fatal("err is empty")
	}
	if status.Code(err) != codes.NotFound {
		t.Fatalf("different error code, expected: %s, actual: %s", codes.NotFound, status.Code(err))
	}
}

func TestUpdateAutoOpsRule(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	clauses := createDatetimeClausesWithActionType(t, 1)
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, autoopsproto.OpsType_SCHEDULE, nil, clauses)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	addClause := autoopsproto.DatetimeClause{
		Time:       time.Now().Unix() + 1000,
		ActionType: autoopsproto.ActionType_DISABLE,
	}
	updateAutoOpsRules(t, autoOpsClient, autoOpsRules[0].Id, &addClause)
	resp, err := autoOpsClient.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentId: *environmentID,
		Id:            autoOpsRules[0].Id,
	})
	if resp == nil {
		t.Fatalf("failed to get AutoOpsRule, err %d", err)
	}

	odc := unmarshalDatetimeClause(t, resp.AutoOpsRule.Clauses[1])
	if odc.Time != addClause.Time {
		t.Fatalf("added DateTime is different, expected: %v, actual: %v", addClause.Time, odc.Time)
	}
	if odc.ActionType != addClause.ActionType {
		t.Fatalf("added ActionType is different, expected: %v, actual: %v", addClause.ActionType, odc.ActionType)
	}
}

func TestUpdateAutoOpsRuleNoCommand(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	clauses := createDatetimeClausesWithActionType(t, 1)
	createAutoOpsRuleNoCommand(ctx, t, autoOpsClient, featureID, autoopsproto.OpsType_SCHEDULE, nil, clauses)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	addClause := autoopsproto.DatetimeClause{
		Time:       time.Now().Unix() + 1000,
		ActionType: autoopsproto.ActionType_DISABLE,
	}
	_, err := autoOpsClient.UpdateAutoOpsRule(ctx, &autoopsproto.UpdateAutoOpsRuleRequest{
		EnvironmentId: *environmentID,
		Id:            autoOpsRules[0].Id,
		DatetimeClauseChanges: []*autoopsproto.DatetimeClauseChange{
			{
				Clause:     &addClause,
				ChangeType: autoopsproto.ChangeType_CREATE,
			},
			{
				Id:         autoOpsRules[0].Clauses[0].Id,
				ChangeType: autoopsproto.ChangeType_DELETE,
			},
		},
	})
	if err != nil {
		t.Fatal("failed to update auto ops rules", err)
	}
	resp, err := autoOpsClient.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentId: *environmentID,
		Id:            autoOpsRules[0].Id,
	})
	if resp == nil {
		t.Fatalf("failed to get AutoOpsRule, err %d", err)
	}

	odc := unmarshalDatetimeClause(t, resp.AutoOpsRule.Clauses[0])
	if odc.Time != addClause.Time {
		t.Fatalf("added DateTime is different, expected: %v, actual: %v", addClause.Time, odc.Time)
	}
	if odc.ActionType != addClause.ActionType {
		t.Fatalf("added ActionType is different, expected: %v, actual: %v", addClause.ActionType, odc.ActionType)
	}
}

func TestExecuteAutoOpsRule(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	goalID := createGoal(ctx, t, experimentClient)
	clause := createOpsEventRateClause(t, feature.Variations[0].Id, goalID)
	createAutoOpsRule(
		ctx,
		t,
		autoOpsClient,
		featureID,
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{clause},
		nil,
	)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	_, err := autoOpsClient.ExecuteAutoOps(ctx, &autoopsproto.ExecuteAutoOpsRequest{
		EnvironmentId: *environmentID,
		Id:            autoOpsRules[0].Id,
		ExecuteAutoOpsRuleCommand: &autoopsproto.ExecuteAutoOpsRuleCommand{
			ClauseId: autoOpsRules[0].Clauses[0].Id,
		},
	})
	if err != nil {
		t.Fatalf("failed to execute auto ops: %s", err.Error())
	}
	feature = getFeature(t, featureClient, featureID)
	if feature.Enabled == true {
		t.Fatalf("feature is enabled")
	}
	autoOpsRules = listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	aor := autoOpsRules[0]
	if aor.AutoOpsStatus != autoopsproto.AutoOpsStatus_RUNNING && aor.AutoOpsStatus != autoopsproto.AutoOpsStatus_FINISHED {
		t.Fatalf("The operation has been executed, but there is a problem with the status. Status: %v", aor.AutoOpsStatus)
	}
}

func TestExecuteAutoOpsRuleNoCommand(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	goalID := createGoal(ctx, t, experimentClient)
	clause := createOpsEventRateClause(t, feature.Variations[0].Id, goalID)
	createAutoOpsRule(
		ctx,
		t,
		autoOpsClient,
		featureID,
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{clause},
		nil,
	)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	_, err := autoOpsClient.ExecuteAutoOps(ctx, &autoopsproto.ExecuteAutoOpsRequest{
		EnvironmentId: *environmentID,
		Id:            autoOpsRules[0].Id,
		ClauseId:      autoOpsRules[0].Clauses[0].Id,
	})
	if err != nil {
		t.Fatalf("failed to execute auto ops: %s", err.Error())
	}
	feature = getFeature(t, featureClient, featureID)
	if feature.Enabled == true {
		t.Fatalf("feature is enabled")
	}
	autoOpsRules = listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	aor := autoOpsRules[0]
	if aor.AutoOpsStatus != autoopsproto.AutoOpsStatus_RUNNING && aor.AutoOpsStatus != autoopsproto.AutoOpsStatus_FINISHED {
		t.Fatalf("The operation has been executed, but there is a problem with the status. Status: %v", aor.AutoOpsStatus)
	}
}

func TestExecuteAutoOpsRuleForMultiSchedule(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	clauses := createDatetimeClausesWithActionType(t, 2)
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, autoopsproto.OpsType_SCHEDULE, nil, clauses)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	_, err := autoOpsClient.ExecuteAutoOps(ctx, &autoopsproto.ExecuteAutoOpsRequest{
		EnvironmentId: *environmentID,
		Id:            autoOpsRules[0].Id,
		ExecuteAutoOpsRuleCommand: &autoopsproto.ExecuteAutoOpsRuleCommand{
			ClauseId: autoOpsRules[0].Clauses[0].Id,
		},
	})
	if err != nil {
		t.Fatalf("failed to execute auto ops: %s", err.Error())
	}
	feature = getFeature(t, featureClient, featureID)
	if feature.Enabled == true {
		t.Fatalf("feature is enabled")
	}
	autoOpsRules = listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if autoOpsRules[0].AutoOpsStatus != autoopsproto.AutoOpsStatus_RUNNING {
		t.Fatalf("status is not running")
	}
}

// Test for old SDK client. Tag is not set in the EvaluationEvent and GoalEvent
// Evaluation field in the GoalEvent is deprecated.
func TestOpsEventRateBatchWithoutTag(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	goalID := createGoal(ctx, t, experimentClient)
	clause := createOpsEventRateClause(t, feature.Variations[0].Id, goalID)
	createAutoOpsRule(
		ctx,
		t,
		autoOpsClient,
		featureID,
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{clause},
		nil,
	)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}

	// Wait for the event-persister-ops subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

	userIDs := createUserIDs(t, 10)
	for _, uid := range userIDs[:6] {
		registerGoalEventWithEvaluations(t, featureID, feature.Version, goalID, uid, feature.Variations[0].Id)
	}
	for _, uid := range userIDs {
		grpcRegisterEvaluationEvent(t, featureID, feature.Version, uid, feature.Variations[0].Id, "")
	}

	checkIfAutoOpsRulesAreTriggered(t, featureID)
}

func TestGrpcOpsEventRateBatch(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	goalID := createGoal(ctx, t, experimentClient)
	clause := createOpsEventRateClause(t, feature.Variations[0].Id, goalID)
	createAutoOpsRule(
		ctx,
		t,
		autoOpsClient,
		featureID,
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{clause},
		nil,
	)

	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}

	// Wait for the event-persister-ops subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

	userIDs := createUserIDs(t, 10)
	for _, uid := range userIDs[:6] {
		grpcRegisterGoalEvent(t, goalID, uid, feature.Tags[0])
	}
	for _, uid := range userIDs {
		grpcRegisterEvaluationEvent(t, featureID, feature.Version, uid, feature.Variations[0].Id, feature.Tags[0])
	}

	checkIfAutoOpsRulesAreTriggered(t, featureID)
}

func TestOpsEventRateBatch(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	now := time.Now()
	schedules := []*autoopsproto.ProgressiveRolloutSchedule{
		{
			Weight:    20000,
			ExecuteAt: now.Add(10 * time.Minute).Unix(),
		},
	}
	createProgressiveRollout(
		ctx,
		t,
		autoOpsClient,
		featureID,
		&autoopsproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: feature.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, autoOpsClient, featureID)
	if len(progressiveRollouts) != 1 {
		t.Fatal("Progressive rollout list shouldn't be empty")
	}

	goalID := createGoal(ctx, t, experimentClient)
	clause := createOpsEventRateClause(t, feature.Variations[0].Id, goalID)
	createAutoOpsRule(
		ctx,
		t,
		autoOpsClient,
		featureID,
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{clause},
		nil,
	)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}

	// Wait for the event-persister-ops subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

	userIDs := createUserIDs(t, 10)
	for _, uid := range userIDs[:6] {
		registerGoalEvent(t, goalID, uid, feature.Tags[0])
	}
	for _, uid := range userIDs {
		registerEvaluationEvent(t, featureID, feature.Version, uid, feature.Variations[0].Id, feature.Tags[0])
	}

	checkIfAutoOpsRulesAreTriggered(t, featureID)

	// As a requirement, when disabling a flag using an auto operation,
	// It must stop the progressive rollout if it is running
	pr := getProgressiveRollout(t, progressiveRollouts[0].Id)
	if pr.Status != autoopsproto.ProgressiveRollout_STOPPED {
		t.Fatalf("Progressive rollout must be stopped. Current status: %v", pr.Status)
	}
}

func TestDatetimeBatch(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	now := time.Now()
	schedules := []*autoopsproto.ProgressiveRolloutSchedule{
		{
			Weight:    20000,
			ExecuteAt: now.Add(10 * time.Minute).Unix(),
		},
	}
	createProgressiveRollout(
		ctx,
		t,
		autoOpsClient,
		featureID,
		&autoopsproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: feature.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, autoOpsClient, featureID)
	if len(progressiveRollouts) != 1 {
		t.Fatal("Progressive rollout list shouldn't be empty")
	}

	clause := createDatetimeClause(t)
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, autoopsproto.OpsType_SCHEDULE, nil, []*autoopsproto.DatetimeClause{clause})
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}

	// Wait for the event-persister-ops subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

	checkIfAutoOpsRulesAreTriggered(t, featureID)

	// As a requirement, when disabling a flag using an auto operation,
	// It must stop the progressive rollout if it is running
	pr := getProgressiveRollout(t, progressiveRollouts[0].Id)
	if pr.Status != autoopsproto.ProgressiveRollout_STOPPED {
		t.Fatalf("Progressive rollout must be stopped. Current status: %v", pr.Status)
	}
}

func TestDatetimeBatchForMultiSchedule(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	now := time.Now()
	schedules := []*autoopsproto.ProgressiveRolloutSchedule{
		{
			Weight:    20000,
			ExecuteAt: now.Add(10 * time.Minute).Unix(),
		},
	}
	createProgressiveRollout(
		ctx,
		t,
		autoOpsClient,
		featureID,
		&autoopsproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: feature.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, autoOpsClient, featureID)
	if len(progressiveRollouts) != 1 {
		t.Fatal("Progressive rollout list shouldn't be empty")
	}

	clauses := createDatetimeClausesWithActionType(t, 2)
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, autoopsproto.OpsType_SCHEDULE, nil, clauses)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	// Wait for the event-persister-ops subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)
	checkIfAutoOpsRulesAreTriggered(t, featureID)
	// As a requirement, when disabling a flag using an auto operation,
	// It must stop the progressive rollout if it is running
	pr := getProgressiveRollout(t, progressiveRollouts[0].Id)
	if pr.Status != autoopsproto.ProgressiveRollout_STOPPED {
		t.Fatalf("Progressive rollout must be stopped. Current status: %v", pr.Status)
	}
}

func TestStopAutoOpsRule(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	clauses := createDatetimeClausesWithActionType(t, 4)
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, autoopsproto.OpsType_SCHEDULE, nil, clauses)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	stopAutoOpsRule(t, autoOpsClient, autoOpsRules[0].Id)
	resp, err := autoOpsClient.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentId: *environmentID,
		Id:            autoOpsRules[0].Id,
	})
	if resp == nil {
		t.Fatalf("failed to get AutoOpsRule, err %d", err)
	}

	if resp.AutoOpsRule.AutoOpsStatus != autoopsproto.AutoOpsStatus_STOPPED {
		t.Fatalf("different auto ops status, expected: %v, actual: %v", autoopsproto.AutoOpsStatus_STOPPED, resp.AutoOpsRule.AutoOpsStatus)
	}
}

func sendHttpWebhook(t *testing.T, url, payload string) {
	t.Helper()
	// Create a custom HTTP client with insecure skip verify
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Send HTTP webhook request failed: %d, %s", resp.StatusCode, url)
	}
}

func unmarshalOpsEventRateClause(t *testing.T, clause *autoopsproto.Clause) *autoopsproto.OpsEventRateClause {
	c := &autoopsproto.OpsEventRateClause{}
	if err := ptypes.UnmarshalAny(clause.Clause, c); err != nil {
		t.Fatal(err)
	}
	return c
}

func unmarshalDatetimeClause(t *testing.T, clause *autoopsproto.Clause) *autoopsproto.DatetimeClause {
	c := &autoopsproto.DatetimeClause{}
	if err := ptypes.UnmarshalAny(clause.Clause, c); err != nil {
		t.Fatal(err)
	}
	return c
}

func createGoal(ctx context.Context, t *testing.T, client experimentclient.Client) string {
	t.Helper()
	goalID := createGoalID(t)
	_, err := client.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
		Id:          goalID,
		Name:        goalID,
		Description: goalID,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return goalID
}

func createOpsEventRateClause(t *testing.T, variationID, goalID string) *autoopsproto.OpsEventRateClause {
	t.Helper()
	return &autoopsproto.OpsEventRateClause{
		VariationId:     variationID,
		GoalId:          goalID,
		MinCount:        int64(5),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
		ActionType:      autoopsproto.ActionType_DISABLE,
	}
}

func createDatetimeClause(t *testing.T) *autoopsproto.DatetimeClause {
	return &autoopsproto.DatetimeClause{
		Time:       time.Now().Add(5 * time.Second).Unix(),
		ActionType: autoopsproto.ActionType_DISABLE,
	}
}

func createDatetimeClausesWithActionType(t *testing.T, createCount int) []*autoopsproto.DatetimeClause {
	t.Helper()
	var dcs []*autoopsproto.DatetimeClause
	for i := 0; i < createCount; i++ {
		at := autoopsproto.ActionType_DISABLE
		if i%2 == 1 {
			at = autoopsproto.ActionType_ENABLE
		}
		dc := &autoopsproto.DatetimeClause{
			Time:       time.Now().Add(time.Duration((i+1)*70) * time.Second).Unix(),
			ActionType: at,
		}
		dcs = append(dcs, dc)
	}
	return dcs
}

func createAutoOpsRule(
	ctx context.Context,
	t *testing.T,
	client autoopsclient.Client,
	featureID string,
	opsType autoopsproto.OpsType,
	oercs []*autoopsproto.OpsEventRateClause,
	dcs []*autoopsproto.DatetimeClause,
) {
	t.Helper()
	cmd := &autoopsproto.CreateAutoOpsRuleCommand{
		FeatureId:           featureID,
		OpsType:             opsType,
		OpsEventRateClauses: oercs,
		DatetimeClauses:     dcs,
	}
	_, err := client.CreateAutoOpsRule(ctx, &autoopsproto.CreateAutoOpsRuleRequest{
		EnvironmentId: *environmentID,
		Command:       cmd,
	})
	if err != nil {
		t.Fatal(err)
	}
	// Update auto ops rules cache
	batchClient := newBatchClient(t)
	defer batchClient.Close()
	numRetries := 5
	for i := 0; i < numRetries; i++ {
		_, err = batchClient.ExecuteBatchJob(
			ctx,
			&btproto.BatchJobRequest{Job: btproto.BatchJob_AutoOpsRulesCacher})
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		if st.Code() == codes.Internal {
			t.Fatal(err)
		}
		fmt.Printf("Failed to execute auto ops rules cacher batch. Error code: %d\n. Retrying in 5 seconds.", st.Code())
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func createAutoOpsRuleNoCommand(
	ctx context.Context,
	t *testing.T,
	client autoopsclient.Client,
	featureID string,
	opsType autoopsproto.OpsType,
	oercs []*autoopsproto.OpsEventRateClause,
	dcs []*autoopsproto.DatetimeClause,
) {
	t.Helper()
	_, err := client.CreateAutoOpsRule(ctx, &autoopsproto.CreateAutoOpsRuleRequest{
		EnvironmentId:       *environmentID,
		FeatureId:           featureID,
		OpsType:             opsType,
		OpsEventRateClauses: oercs,
		DatetimeClauses:     dcs,
	})
	if err != nil {
		t.Fatal(err)
	}
	// Update auto ops rules cache
	batchClient := newBatchClient(t)
	defer batchClient.Close()
	numRetries := 5
	for i := 0; i < numRetries; i++ {
		_, err = batchClient.ExecuteBatchJob(
			ctx,
			&btproto.BatchJobRequest{Job: btproto.BatchJob_AutoOpsRulesCacher})
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		if st.Code() == codes.Internal {
			t.Fatal(err)
		}
		fmt.Printf("Failed to execute auto ops rules cacher batch. Error code: %d\n. Retrying in 5 seconds.", st.Code())
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func newAutoOpsClient(t *testing.T) autoopsclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := autoopsclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create auto ops client:", err)
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

func newWebhookName(t *testing.T) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-webhook-name-%s", prefixID, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-webhook-name-%s", prefixID, newUUID(t))
}

func createFeature(ctx context.Context, t *testing.T, client featureclient.Client, featureID string) {
	t.Helper()
	cmd := newCreateFeatureCommand(featureID)
	createReq := &featureproto.CreateFeatureRequest{
		Command:       cmd,
		EnvironmentId: *environmentID,
	}
	if _, err := client.CreateFeature(ctx, createReq); err != nil {
		t.Fatal(err)
	}
	enableFeature(t, featureID, client)
}

func createDisabledFeature(ctx context.Context, t *testing.T, client featureclient.Client, featureID string) {
	t.Helper()
	cmd := newCreateFeatureCommand(featureID)
	createReq := &featureproto.CreateFeatureRequest{
		Command:       cmd,
		EnvironmentId: *environmentID,
	}
	if _, err := client.CreateFeature(ctx, createReq); err != nil {
		t.Fatal(err)
	}
}

func getFeature(t *testing.T, client featureclient.Client, featureID string) *featureproto.Feature {
	t.Helper()
	getReq := &featureproto.GetFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := client.GetFeature(ctx, getReq)
	if err != nil {
		t.Fatal("Failed to get feature:", err)
	}
	return response.Feature
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

func newCreateFeatureCommand(featureID string) *featureproto.CreateFeatureCommand {
	return &featureproto.CreateFeatureCommand{
		Id:          featureID,
		Name:        "e2e-test-gateway-feature-name",
		Description: "e2e-test-gateway-feature-description",
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
}

func enableFeature(t *testing.T, featureID string, client featureclient.Client) {
	t.Helper()
	enableReq := &featureproto.EnableFeatureRequest{
		Id:            featureID,
		Command:       &featureproto.EnableFeatureCommand{},
		EnvironmentId: *environmentID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.EnableFeature(ctx, enableReq); err != nil {
		t.Fatalf("Failed to enable feature id: %s. Error: %v", featureID, err)
	}
}

func listAutoOpsRulesByFeatureID(t *testing.T, client autoopsclient.Client, featureID string) []*autoopsproto.AutoOpsRule {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.ListAutoOpsRules(ctx, &autoopsproto.ListAutoOpsRulesRequest{
		EnvironmentId: *environmentID,
		PageSize:      int64(500),
		FeatureIds:    []string{featureID},
	})
	if err != nil {
		t.Fatal("failed to list auto ops rules", err)
	}
	return resp.AutoOpsRules
}

func getAutoOpsRules(t *testing.T, id string) *autoopsproto.AutoOpsRule {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAutoOpsClient(t)
	defer c.Close()
	resp, err := c.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentId: *environmentID,
		Id:            id,
	})
	if err != nil {
		t.Fatal("failed to list auto ops rules", err)
	}
	return resp.AutoOpsRule
}

func deleteAutoOpsRules(t *testing.T, client autoopsclient.Client, id string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.DeleteAutoOpsRule(ctx, &autoopsproto.DeleteAutoOpsRuleRequest{
		EnvironmentId: *environmentID,
		Id:            id,
	})
	if err != nil {
		t.Fatal("failed to list auto ops rules", err)
	}
}

func stopAutoOpsRule(t *testing.T, client autoopsclient.Client, id string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.StopAutoOpsRule(ctx, &autoopsproto.StopAutoOpsRuleRequest{
		EnvironmentId: *environmentID,
		Id:            id,
	})
	if err != nil {
		t.Fatal("failed to stop auto ops rules", err)
	}
}

func updateAutoOpsRules(t *testing.T, client autoopsclient.Client, id string, dateClause *autoopsproto.DatetimeClause) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.UpdateAutoOpsRule(ctx, &autoopsproto.UpdateAutoOpsRuleRequest{
		EnvironmentId: *environmentID,
		Id:            id,
		AddDatetimeClauseCommands: []*autoopsproto.AddDatetimeClauseCommand{
			{DatetimeClause: dateClause},
		},
	})
	if err != nil {
		t.Fatal("failed to update auto ops rules", err)
	}
}

// Test for old SDK client
// Evaluation field in the GoalEvent is deprecated.
func registerGoalEventWithEvaluations(
	t *testing.T,
	featureID string,
	featureVersion int32,
	goalID, userID, variationID string,
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	goal, err := ptypes.MarshalAny(&eventproto.GoalEvent{
		Timestamp: time.Now().Unix(),
		GoalId:    goalID,
		UserId:    userID,
		Value:     0.3,
		User:      &userproto.User{},
		Evaluations: []*featureproto.Evaluation{
			{
				Id:             fmt.Sprintf("%s-evaluation-id-%s", prefixTestName, newUUID(t)),
				FeatureId:      featureID,
				FeatureVersion: featureVersion,
				UserId:         userID,
				VariationId:    variationID,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	req := &gatewayproto.RegisterEventsRequest{
		Events: []*eventproto.Event{
			{
				Id:    newUUID(t),
				Event: goal,
			},
		},
	}
	response, err := c.RegisterEvents(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func grpcRegisterGoalEvent(
	t *testing.T,
	goalID, userID, tag string,
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	goal, err := ptypes.MarshalAny(&eventproto.GoalEvent{
		Timestamp: time.Now().Unix(),
		GoalId:    goalID,
		UserId:    userID,
		Value:     0.3,
		User:      &userproto.User{},
		Tag:       tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := &gatewayproto.RegisterEventsRequest{
		Events: []*eventproto.Event{
			{
				Id:    newUUID(t),
				Event: goal,
			},
		},
	}
	response, err := c.RegisterEvents(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func registerGoalEvent(
	t *testing.T,
	goalID, userID, tag string,
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	goal, err := protojson.Marshal(&eventproto.GoalEvent{
		Timestamp: time.Now().Unix(),
		GoalId:    goalID,
		UserId:    userID,
		Value:     0.3,
		User:      &userproto.User{},
		Tag:       tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	events := []util.Event{
		{
			ID:    newUUID(t),
			Event: goal,
			Type:  gwapi.GoalEventType,
		},
	}
	response := util.RegisterEvents(t, events, *gatewayAddr, *apiKeyPath)
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func newGatewayClient(t *testing.T) gatewayclient.Client {
	t.Helper()
	creds, err := gatewayclient.NewPerRPCCredentials(*apiKeyPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := gatewayclient.NewClient(
		fmt.Sprintf("%s:%d", *gatewayAddr, *gatewayPort),
		*gatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create gateway client:", err)
	}
	return client
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

func grpcRegisterEvaluationEvent(
	t *testing.T,
	featureID string,
	featureVersion int32,
	userID, variationID, tag string,
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	evaluation, err := ptypes.MarshalAny(&eventproto.EvaluationEvent{
		Timestamp:      time.Now().Unix(),
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		UserId:         userID,
		VariationId:    variationID,
		User:           &userproto.User{},
		Reason:         &featureproto.Reason{},
		Tag:            tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := &gatewayproto.RegisterEventsRequest{
		Events: []*eventproto.Event{
			{
				Id:    newUUID(t),
				Event: evaluation,
			},
		},
	}
	response, err := c.RegisterEvents(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func registerEvaluationEvent(
	t *testing.T,
	featureID string,
	featureVersion int32,
	userID, variationID, tag string,
) {
	t.Helper()
	evaluation, err := protojson.Marshal(&eventproto.EvaluationEvent{
		Timestamp:      time.Now().Unix(),
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		UserId:         userID,
		VariationId:    variationID,
		User:           &userproto.User{},
		Reason:         &featureproto.Reason{},
		Tag:            tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	events := []util.Event{
		{
			ID:    newUUID(t),
			Event: evaluation,
			Type:  gwapi.EvaluationEventType,
		},
	}
	response := util.RegisterEvents(t, events, *gatewayAddr, *apiKeyPath)
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func createUserIDs(t *testing.T, total int) []string {
	t.Helper()
	userIDs := make([]string, 0)
	for i := 0; i < total; i++ {
		id := newUUID(t)
		userID := fmt.Sprintf("%s-user-%s", prefixTestName, id)
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

func createFeatureID(t *testing.T) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-feature-id-%s", prefixTestName, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-feature-id-%s", prefixTestName, newUUID(t))
}

func createGoalID(t *testing.T) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-goal-id-%s", prefixTestName, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-goal-id-%s", prefixTestName, newUUID(t))
}

func checkIfAutoOpsRulesAreTriggered(t *testing.T, featureID string) {
	t.Helper()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)
		feature := getFeature(t, featureClient, featureID)
		if feature.Enabled {
			continue
		}
		autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
		aor := autoOpsRules[0]
		if aor.AutoOpsStatus != autoopsproto.AutoOpsStatus_RUNNING && aor.AutoOpsStatus != autoopsproto.AutoOpsStatus_FINISHED {
			t.Fatalf("The operation has been executed, but there is a problem with the status. Status: %v", aor.AutoOpsStatus)
		}
		break
	}
}
