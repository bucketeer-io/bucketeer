// Copyright 2022 The Bucketeer Authors.
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

	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	gatewayclient "github.com/bucketeer-io/bucketeer/pkg/gateway/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
	"github.com/bucketeer-io/bucketeer/test/e2e/util"
)

const (
	goalEventType eventType = iota + 1 // eventType starts from 1 for validation.
	evaluationEventType
	metricsEventType
	prefixTestName   = "e2e-test"
	retryTimes       = 60
	timeout          = 10 * time.Second
	prefixID         = "e2e-test"
	version          = "/v1"
	service          = "/gateway"
	eventsAPI        = "/events"
	authorizationKey = "authorization"
)

var (
	webGatewayAddr       = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort       = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert       = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath           = flag.String("api-key", "", "Api key path for web gateway")
	gatewayAddr          = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort          = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert          = flag.String("gateway-cert", "", "Gateway crt file")
	serviceTokenPath     = flag.String("service-token", "", "Service token path")
	environmentNamespace = flag.String("environment-namespace", "", "Environment namespace")
	testID               = flag.String("test-id", "", "test ID")
)

type eventType int

type event struct {
	ID                   string          `json:"id,omitempty"`
	Event                json.RawMessage `json:"event,omitempty"`
	EnvironmentNamespace string          `json:"environment_namespace,omitempty"`
	Type                 eventType       `json:"type,omitempty"`
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
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, []*autoopsproto.OpsEventRateClause{clause}, nil, nil)
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
	if actual.OpsType != autoopsproto.OpsType_DISABLE_FEATURE {
		t.Fatalf("different ops type, expected: %v, actual: %v", autoopsproto.OpsType_DISABLE_FEATURE, actual.OpsType)
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
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, []*autoopsproto.OpsEventRateClause{clause}, nil, nil)
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
	if actual.OpsType != autoopsproto.OpsType_DISABLE_FEATURE {
		t.Fatalf("different ops type, expected: %v, actual: %v", autoopsproto.OpsType_DISABLE_FEATURE, actual.OpsType)
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
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, []*autoopsproto.OpsEventRateClause{clause}, nil, nil)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	deleteAutoOpsRules(t, autoOpsClient, autoOpsRules[0].Id)
	resp, err := autoOpsClient.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   autoOpsRules[0].Id,
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
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, []*autoopsproto.OpsEventRateClause{clause}, nil, nil)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	_, err := autoOpsClient.ExecuteAutoOps(ctx, &autoopsproto.ExecuteAutoOpsRequest{
		EnvironmentNamespace:                *environmentNamespace,
		Id:                                  autoOpsRules[0].Id,
		ChangeAutoOpsRuleTriggeredAtCommand: &autoopsproto.ChangeAutoOpsRuleTriggeredAtCommand{},
	})
	if err != nil {
		t.Fatalf("failed to execute auto ops: %s", err.Error())
	}
	feature = getFeature(t, featureClient, featureID)
	if feature.Enabled == true {
		t.Fatalf("feature is enabled")
	}
	autoOpsRules = listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if autoOpsRules[0].TriggeredAt == 0 {
		t.Fatalf("triggered at is empty")
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
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, []*autoopsproto.OpsEventRateClause{clause}, nil, nil)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	// Wait until trasformer and watcher's targetstores are refreshed.
	time.Sleep(40 * time.Second)

	userIDs := createUserIDs(t, 10)
	for _, uid := range userIDs[:6] {
		registerGoalEventWithEvaluations(t, featureID, feature.Version, goalID, uid, feature.Variations[0].Id)
	}
	for _, uid := range userIDs {
		grpcRegisterEvaluationEvent(t, featureID, feature.Version, uid, feature.Variations[0].Id, "")
	}
	for i := 0; i < retryTimes; i++ {
		feature = getFeature(t, featureClient, featureID)
		if !feature.Enabled {
			autoOpsRules = listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
			if autoOpsRules[0].TriggeredAt == 0 {
				t.Fatalf("triggered at must not be zero")
			}
			break
		}
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)
	}
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
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, []*autoopsproto.OpsEventRateClause{clause}, nil, nil)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	// Wait until trasformer and watcher's targetstores are refreshed.
	time.Sleep(40 * time.Second)

	userIDs := createUserIDs(t, 10)
	for _, uid := range userIDs[:6] {
		grpcRegisterGoalEvent(t, goalID, uid, feature.Tags[0])
	}
	for _, uid := range userIDs {
		grpcRegisterEvaluationEvent(t, featureID, feature.Version, uid, feature.Variations[0].Id, feature.Tags[0])
	}
	for i := 0; i < retryTimes; i++ {
		feature = getFeature(t, featureClient, featureID)
		if !feature.Enabled {
			autoOpsRules = listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
			if autoOpsRules[0].TriggeredAt == 0 {
				t.Fatalf("triggered at must not be zero")
			}
			break
		}
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)
	}
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
	goalID := createGoal(ctx, t, experimentClient)
	clause := createOpsEventRateClause(t, feature.Variations[0].Id, goalID)
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, []*autoopsproto.OpsEventRateClause{clause}, nil, nil)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	// Wait until trasformer and watcher's targetstores are refreshed.
	time.Sleep(40 * time.Second)

	userIDs := createUserIDs(t, 10)
	for _, uid := range userIDs[:6] {
		registerGoalEvent(t, goalID, uid, feature.Tags[0])
	}
	for _, uid := range userIDs {
		registerEvaluationEvent(t, featureID, feature.Version, uid, feature.Variations[0].Id, feature.Tags[0])
	}
	for i := 0; i < retryTimes; i++ {
		feature = getFeature(t, featureClient, featureID)
		if !feature.Enabled {
			autoOpsRules = listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
			if autoOpsRules[0].TriggeredAt == 0 {
				t.Fatalf("triggered at must not be zero")
			}
			break
		}
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)
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
	clause := createDatetimeClause(t)
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, nil, []*autoopsproto.DatetimeClause{clause}, nil)
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if len(autoOpsRules) != 1 {
		t.Fatal("not enough rules")
	}
	// Wait until watcher's targetstore is refreshed and autoOps is executed.
	time.Sleep(50 * time.Second)

	feature := getFeature(t, featureClient, featureID)
	if feature.Enabled {
		t.Fatalf("feature must be disabled")
	}
	autoOpsRules = listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if autoOpsRules[0].TriggeredAt == 0 {
		t.Fatalf("triggered at must not be zero")
	}
}

func TestCreateAndListWebhook(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()

	name := newWebhookName(t)
	description := newUUID(t)
	resp := createWebhook(ctx, t, autoOpsClient, name, description)
	webhooks := listWebhooks(ctx, t, autoOpsClient)
	var actual *autoopsproto.Webhook
	for _, w := range webhooks {
		if w.Id == resp.Webhook.Id {
			actual = w
			break
		}
	}
	if actual == nil {
		t.Fatal("webhook is nil")
	}
	if actual.Id == "" {
		t.Fatal("id is empty")
	}
	if actual.Name != name {
		t.Fatalf("diffrent name, expected: %v, actual: %v", name, actual.Name)
	}
	if actual.Description != description {
		t.Fatalf("diffrent description, expected: %v, actual: %v", description, actual.Description)
	}
}

func TestGetWebhook(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()

	name := newWebhookName(t)
	description := newUUID(t)
	resp := createWebhook(ctx, t, autoOpsClient, name, description)
	actual := getWebhook(ctx, t, autoOpsClient, resp.Webhook.Id)
	if actual == nil {
		t.Fatal("webhook is nil")
	}
	if actual.Id == "" {
		t.Fatal("id is empty")
	}
	if actual.Name != name {
		t.Fatalf("diffrent name, expected: %v, actual: %v", name, actual.Name)
	}
	if actual.Description != description {
		t.Fatalf("diffrent description, expected: %v, actual: %v", description, actual.Description)
	}
}

func TestUpdateWebhook(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()

	name := newWebhookName(t)
	description := newUUID(t)
	resp := createWebhook(ctx, t, autoOpsClient, name, description)
	webhook := getWebhook(ctx, t, autoOpsClient, resp.Webhook.Id)
	if webhook == nil {
		t.Fatal("webhook is nil")
	}
	newDesc := newUUID(t)
	newName := newWebhookName(t)
	updateWebhookDescription(ctx, t, autoOpsClient, resp.Webhook.Id, newDesc)
	updateWebhookName(ctx, t, autoOpsClient, resp.Webhook.Id, newName)
	actual := getWebhook(ctx, t, autoOpsClient, resp.Webhook.Id)
	if actual.Id == "" {
		t.Fatal("id is empty")
	}
	if actual.Name != newName {
		t.Fatalf("diffrent name, expected: %v, actual: %v", name, actual.Name)
	}
	if actual.Description != newDesc {
		t.Fatalf("diffrent description, expected: %v, actual: %v", description, actual.Description)
	}
}

func TestDeleteWebhook(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()

	name := newWebhookName(t)
	description := newUUID(t)
	resp := createWebhook(ctx, t, autoOpsClient, name, description)
	webhook := getWebhook(ctx, t, autoOpsClient, resp.Webhook.Id)
	if webhook == nil {
		t.Fatal("webhook is nil")
	}
	deleteWebhook(ctx, t, autoOpsClient, resp.Webhook.Id)
	getResp, err := autoOpsClient.GetWebhook(ctx, &autoopsproto.GetWebhookRequest{
		Id:                   resp.Webhook.Id,
		EnvironmentNamespace: *environmentNamespace,
	})
	if getResp != nil {
		t.Fatal("webhook is not deleted")
	}
	if status.Code(err) != codes.NotFound {
		t.Fatalf("different error code, expected: %s, actual: %s", codes.NotFound, status.Code(err))
	}
}

func TestHttpWebhook(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()

	name := newWebhookName(t)
	description := newUUID(t)
	resp := createWebhook(ctx, t, autoOpsClient, name, description)
	webhook := getWebhook(ctx, t, autoOpsClient, resp.Webhook.Id)
	if webhook == nil {
		t.Fatal("webhook is nil")
	}
	featureID := createFeatureID(t)
	createFeature(ctx, t, featureClient, featureID)
	condition := createWebhookClause_Condition(autoopsproto.WebhookClause_Condition_EQUAL, `.body."Alert id"`, `123`)
	clause := createWebhookClause(resp.Webhook.Id, []*autoopsproto.WebhookClause_Condition{condition})
	createAutoOpsRule(ctx, t, autoOpsClient, featureID, nil, nil, []*autoopsproto.WebhookClause{clause})
	autoOpsRules := listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	expectedNum := 1
	if len(autoOpsRules) != expectedNum {
		t.Fatal("not enough rules")
	}

	sendHttpWebhook(t, resp.Url, `{"body":{"Alert id": 123}}`)
	feature := getFeature(t, featureClient, featureID)
	if feature.Enabled == true {
		t.Fatalf("feature is enabled")
	}
	autoOpsRules = listAutoOpsRulesByFeatureID(t, autoOpsClient, featureID)
	if autoOpsRules[0].TriggeredAt == 0 {
		t.Fatalf("triggered at is empty")
	}
}

func sendHttpWebhook(t *testing.T, url, payload string) {
	t.Helper()
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))
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

func createGoal(ctx context.Context, t *testing.T, client experimentclient.Client) string {
	t.Helper()
	goalID := createGoalID(t)
	cmd := &experimentproto.CreateGoalCommand{
		Id:          goalID,
		Name:        goalID,
		Description: goalID,
	}
	_, err := client.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
		Command:              cmd,
		EnvironmentNamespace: *environmentNamespace,
	})
	if err != nil {
		t.Fatal(err)
	}
	return cmd.Id
}

func createOpsEventRateClause(t *testing.T, variationID string, goalID string) *autoopsproto.OpsEventRateClause {
	return &autoopsproto.OpsEventRateClause{
		VariationId:     variationID,
		GoalId:          goalID,
		MinCount:        int64(5),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
}

func createDatetimeClause(t *testing.T) *autoopsproto.DatetimeClause {
	return &autoopsproto.DatetimeClause{
		Time: time.Now().Add(5 * time.Second).Unix(),
	}
}

func createWebhookClause(webhookID string, condition []*autoopsproto.WebhookClause_Condition) *autoopsproto.WebhookClause {
	return &autoopsproto.WebhookClause{
		WebhookId:  webhookID,
		Conditions: condition,
	}
}

func createWebhookClause_Condition(operator autoopsproto.WebhookClause_Condition_Operator, filter, value string) *autoopsproto.WebhookClause_Condition {
	return &autoopsproto.WebhookClause_Condition{
		Filter:   filter,
		Value:    value,
		Operator: operator,
	}
}

func createAutoOpsRule(ctx context.Context, t *testing.T, client autoopsclient.Client, featureID string, oercs []*autoopsproto.OpsEventRateClause, dcs []*autoopsproto.DatetimeClause, wc []*autoopsproto.WebhookClause) {
	cmd := &autoopsproto.CreateAutoOpsRuleCommand{
		FeatureId:           featureID,
		OpsType:             autoopsproto.OpsType_DISABLE_FEATURE,
		OpsEventRateClauses: oercs,
		DatetimeClauses:     dcs,
		WebhookClauses:      wc,
	}
	_, err := client.CreateAutoOpsRule(ctx, &autoopsproto.CreateAutoOpsRuleRequest{
		EnvironmentNamespace: *environmentNamespace,
		Command:              cmd,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createWebhook(ctx context.Context, t *testing.T, client autoopsclient.Client, name, description string) *autoopsproto.CreateWebhookResponse {
	t.Helper()
	cmd := &autoopsproto.CreateWebhookCommand{
		Name:        name,
		Description: description,
	}
	resp, err := client.CreateWebhook(ctx, &autoopsproto.CreateWebhookRequest{
		EnvironmentNamespace: *environmentNamespace,
		Command:              cmd,
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp
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
	if *testID != "" {
		return fmt.Sprintf("%s-%s-webhook-name-%s", prefixID, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-webhook-name-%s", prefixID, newUUID(t))
}

func createFeature(ctx context.Context, t *testing.T, client featureclient.Client, featureID string) {
	t.Helper()
	cmd := newCreateFeatureCommand(featureID)
	createReq := &featureproto.CreateFeatureRequest{
		Command:              cmd,
		EnvironmentNamespace: *environmentNamespace,
	}
	if _, err := client.CreateFeature(ctx, createReq); err != nil {
		t.Fatal(err)
	}
	enableFeature(t, featureID, client)
}

func getFeature(t *testing.T, client featureclient.Client, featureID string) *featureproto.Feature {
	t.Helper()
	getReq := &featureproto.GetFeatureRequest{
		Id:                   featureID,
		EnvironmentNamespace: *environmentNamespace,
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
		Id:                   featureID,
		Command:              &featureproto.EnableFeatureCommand{},
		EnvironmentNamespace: *environmentNamespace,
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
		EnvironmentNamespace: *environmentNamespace,
		PageSize:             int64(500),
		FeatureIds:           []string{featureID},
	})
	if err != nil {
		t.Fatal("failed to list auto ops rules", err)
	}
	return resp.AutoOpsRules
}

func listWebhooks(ctx context.Context, t *testing.T, client autoopsclient.Client) []*autoopsproto.Webhook {
	t.Helper()
	resp, err := client.ListWebhooks(ctx, &autoopsproto.ListWebhooksRequest{
		PageSize:             int64(500),
		EnvironmentNamespace: *environmentNamespace,
	})
	if err != nil {
		t.Fatal("failed to list webhooks", err)
	}
	return resp.Webhooks
}

func getWebhook(ctx context.Context, t *testing.T, client autoopsclient.Client, id string) *autoopsproto.Webhook {
	t.Helper()
	resp, err := client.GetWebhook(ctx, &autoopsproto.GetWebhookRequest{
		Id:                   id,
		EnvironmentNamespace: *environmentNamespace,
	})
	if err != nil {
		t.Fatal("failed to get webhook", err)
	}
	return resp.Webhook
}

func getAutoOpsRules(t *testing.T, id string) *autoopsproto.AutoOpsRule {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAutoOpsClient(t)
	defer c.Close()
	resp, err := c.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
	})
	if err != nil {
		t.Fatal("failed to list auto ops rules", err)
	}
	return resp.AutoOpsRule
}

func updateWebhookDescription(ctx context.Context, t *testing.T, client autoopsclient.Client, id, desc string) {
	t.Helper()
	_, err := client.UpdateWebhook(ctx, &autoopsproto.UpdateWebhookRequest{
		Id:                   id,
		EnvironmentNamespace: *environmentNamespace,
		ChangeWebhookDescriptionCommand: &autoopsproto.ChangeWebhookDescriptionCommand{
			Description: desc,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func updateWebhookName(ctx context.Context, t *testing.T, client autoopsclient.Client, id, name string) {
	t.Helper()
	_, err := client.UpdateWebhook(ctx, &autoopsproto.UpdateWebhookRequest{
		Id:                   id,
		EnvironmentNamespace: *environmentNamespace,
		ChangeWebhookNameCommand: &autoopsproto.ChangeWebhookNameCommand{
			Name: name,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func deleteAutoOpsRules(t *testing.T, client autoopsclient.Client, id string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.DeleteAutoOpsRule(ctx, &autoopsproto.DeleteAutoOpsRuleRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
		Command:              &autoopsproto.DeleteAutoOpsRuleCommand{},
	})
	if err != nil {
		t.Fatal("failed to list auto ops rules", err)
	}
}

func deleteWebhook(ctx context.Context, t *testing.T, client autoopsclient.Client, id string) {
	t.Helper()
	_, err := client.DeleteWebhook(ctx, &autoopsproto.DeleteWebhookRequest{
		Id:                   id,
		EnvironmentNamespace: *environmentNamespace,
		Command:              &autoopsproto.DeleteWebhookCommand{},
	})
	if err != nil {
		t.Fatal(err)
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
			Type:  util.GoalEventType,
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
			Type:  util.EvaluationEventType,
		},
	}
	response := util.RegisterEvents(t, events, *gatewayAddr, *apiKeyPath)
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func createUserIDs(t *testing.T, total int) []string {
	userIDs := make([]string, 0)
	for i := 0; i < total; i++ {
		id := newUUID(t)
		userID := fmt.Sprintf("%s-user-%s", prefixTestName, id)
		userIDs = append(userIDs, userID)
	}
	return userIDs
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
