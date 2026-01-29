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
//

package feature

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestCreateFeatureFlagTrigger(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	cmd := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	// Create flag trigger
	createFlagTriggerReq := newCreateFlagTriggerReq(
		cmd.Id,
		newTriggerDescription(t),
		featureproto.FlagTrigger_Action_ON,
	)
	resp := createFeatureFlagTrigger(t, client, createFlagTriggerReq)
	if resp.FlagTrigger.FeatureId != cmd.Id {
		t.Fatalf("unexpected flag feature id: %s, feature id: %s", resp.FlagTrigger.FeatureId, cmd.Id)
	}
	if resp.FlagTrigger.Type != createFlagTriggerReq.Type {
		t.Fatalf("unexpected trigger type: %s, type: %s", resp.FlagTrigger.Type, createFlagTriggerReq.Type)
	}
	if resp.FlagTrigger.Action != createFlagTriggerReq.Action {
		t.Fatalf("unexpected trigger action: %s, action: %s",
			resp.FlagTrigger.Action, createFlagTriggerReq.Action)
	}
	if resp.FlagTrigger.Description != createFlagTriggerReq.Description {
		t.Fatalf("unexpected trigger description: %s, description: %s",
			resp.FlagTrigger.Description, createFlagTriggerReq.Description)
	}
	if resp.GetUrl() == "" {
		t.Fatal("unexpected empty url")
	}
}

func TestUpdateFeatureFlagDescriptionTriggerNoCommand(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	req := newCreateFeatureReq(newFeatureID(t))
	createFeatureNoCmd(t, client, req)
	// Create flag trigger
	createFlagTriggerReq := newCreateFlagTriggerReq(
		req.Id,
		newTriggerDescription(t),
		featureproto.FlagTrigger_Action_ON,
	)
	resp, err := client.CreateFlagTrigger(context.Background(), createFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp.FlagTrigger.FeatureId != req.Id {
		t.Fatalf("unexpected flag feature id: %s, feature id: %s", resp.FlagTrigger.FeatureId, req.Id)
	}
	if resp.FlagTrigger.Type != createFlagTriggerReq.Type {
		t.Fatalf("unexpected trigger type: %s, type: %s", resp.FlagTrigger.Type, createFlagTriggerReq.Type)
	}
	if resp.FlagTrigger.Action != createFlagTriggerReq.Action {
		t.Fatalf("unexpected trigger action: %s, action: %s",
			resp.FlagTrigger.Action, createFlagTriggerReq.Action)
	}
	if resp.FlagTrigger.Description != createFlagTriggerReq.Description {
		t.Fatalf("unexpected trigger description: %s, description: %s",
			resp.FlagTrigger.Description, createFlagTriggerReq.Description)
	}
	if resp.GetUrl() == "" {
		t.Fatal("unexpected empty url")
	}
	updateFlagTriggerReq := &featureproto.UpdateFlagTriggerRequest{
		Id:            resp.FlagTrigger.Id,
		EnvironmentId: *environmentID,
		Description:   wrapperspb.String(newTriggerDescription(t)),
	}
	_, err = client.UpdateFlagTrigger(context.Background(), updateFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	// Get flag trigger
	getFlagTriggerReq := &featureproto.GetFlagTriggerRequest{
		Id:            resp.FlagTrigger.Id,
		EnvironmentId: *environmentID,
	}
	getResp := getFeatureFlagTrigger(t, client, getFlagTriggerReq)
	if getResp.FlagTrigger.Description != updateFlagTriggerReq.Description.Value {
		t.Fatal("unexpected description")
	}
}

func TestUpdateFlagTrigger(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	command := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, command)
	// Create flag trigger
	createFlagTriggerReq := newCreateFlagTriggerReq(
		command.Id,
		newTriggerDescription(t),
		featureproto.FlagTrigger_Action_ON,
	)
	createResp := createFeatureFlagTrigger(t, client, createFlagTriggerReq)
	// Update flag trigger
	updateFlagTriggerReq := &featureproto.UpdateFlagTriggerRequest{
		Id:            createResp.FlagTrigger.Id,
		EnvironmentId: *environmentID,
		Description:   wrapperspb.String(newTriggerDescription(t)),
	}
	_, err := client.UpdateFlagTrigger(context.Background(), updateFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	// Get flag trigger
	getFlagTriggerReq := &featureproto.GetFlagTriggerRequest{
		Id:            createResp.FlagTrigger.Id,
		EnvironmentId: *environmentID,
	}
	resp := getFeatureFlagTrigger(t, client, getFlagTriggerReq)
	if resp.FlagTrigger.Description != updateFlagTriggerReq.Description.Value {
		t.Fatal("unexpected description")
	}
}

func TestResetFlagUpdateTrigger(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	req := newCreateFeatureReq(newFeatureID(t))
	createFeatureNoCmd(t, client, req)
	// Create flag trigger
	createFlagTriggerReq := newCreateFlagTriggerReq(
		req.Id,
		newTriggerDescription(t),
		featureproto.FlagTrigger_Action_ON,
	)
	createResp, err := client.CreateFlagTrigger(context.Background(), createFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1)
	updateFlagTriggerReq := &featureproto.UpdateFlagTriggerRequest{
		Id:            createResp.FlagTrigger.Id,
		EnvironmentId: *environmentID,
		Reset_:        true,
	}
	updateResp, err := client.UpdateFlagTrigger(context.Background(), updateFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, createResp.Url, updateResp.Url)
}

func TestDeleteFlagTrigger(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	command := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, command)
	// Create flag trigger
	createFlagTriggerReq := newCreateFlagTriggerReq(
		command.Id,
		newTriggerDescription(t),
		featureproto.FlagTrigger_Action_ON,
	)
	createResp := createFeatureFlagTrigger(t, client, createFlagTriggerReq)
	// Delete flag trigger
	deleteFlagTriggerReq := &featureproto.DeleteFlagTriggerRequest{
		Id:            createResp.FlagTrigger.Id,
		EnvironmentId: *environmentID,
	}
	_, err := client.DeleteFlagTrigger(context.Background(), deleteFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	// Get flag trigger
	getFlagTriggerReq := &featureproto.GetFlagTriggerRequest{
		Id:            createResp.FlagTrigger.Id,
		EnvironmentId: *environmentID,
	}
	_, err = client.GetFlagTrigger(context.Background(), getFlagTriggerReq)
	if err == nil {
		t.Fatal("flag trigger delete error")
	}
}

func TestListFlagTriggers(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	command := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, command)
	// Create flag triggers
	trigger1, err := client.CreateFlagTrigger(context.Background(), newCreateFlagTriggerReq(
		command.Id,
		newTriggerDescription(t),
		featureproto.FlagTrigger_Action_ON,
	))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	trigger2, err := client.CreateFlagTrigger(context.Background(), newCreateFlagTriggerReq(
		command.Id,
		newTriggerDescription(t),
		featureproto.FlagTrigger_Action_ON,
	))
	if err != nil {
		t.Fatal(err)
	}
	// List flag triggers
	listFlagTriggersReq := &featureproto.ListFlagTriggersRequest{
		FeatureId:      command.Id,
		EnvironmentId:  *environmentID,
		Cursor:         "0",
		PageSize:       10,
		OrderBy:        featureproto.ListFlagTriggersRequest_CREATED_AT,
		OrderDirection: featureproto.ListFlagTriggersRequest_ASC,
	}
	triggers, err := client.ListFlagTriggers(context.Background(), listFlagTriggersReq)
	if err != nil {
		t.Fatal(err)
	}
	if len(triggers.FlagTriggers) != 2 {
		t.Fatalf("unexpected length: %d", len(triggers.FlagTriggers))
	}
	if triggers.FlagTriggers[0].FlagTrigger.GetId() != trigger1.FlagTrigger.Id {
		t.Fatalf(
			"unexpected id: %s , id: %s",
			triggers.FlagTriggers[0].FlagTrigger.GetId(),
			trigger1.FlagTrigger.Id,
		)
	}
	if triggers.FlagTriggers[1].FlagTrigger.GetId() != trigger2.FlagTrigger.Id {
		t.Fatalf(
			"unexpected id: %s , id: %s",
			triggers.FlagTriggers[1].FlagTrigger.GetId(),
			trigger2.FlagTrigger.Id,
		)
	}
	if triggers.TotalCount != 2 {
		t.Fatalf("unexpected total count: %d", triggers.TotalCount)
	}
	if triggers.Cursor != "2" {
		t.Fatalf("unexpected cursor: %s", triggers.Cursor)
	}
}

func TestFeatureFlagWebhook(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	command := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, command)
	// Create Enable flag triggers
	enableTrigger, err := client.CreateFlagTrigger(context.Background(), newCreateFlagTriggerReq(
		command.Id,
		newTriggerDescription(t),
		featureproto.FlagTrigger_Action_ON,
	))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("enable trigger URL: %s", enableTrigger.GetUrl())
	// Send post request
	resp, err := sendPostRequestIgnoreSSL(enableTrigger.GetUrl())
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %s", resp.Status)
	}
	enabledFeature := getFeature(t, command.Id, client)
	if enabledFeature.Enabled != true {
		t.Fatalf("unexpected enabled: %v", enabledFeature.Enabled)
	}
	enabledTrigger := getFeatureFlagTrigger(t, client, &featureproto.GetFlagTriggerRequest{
		Id:            enableTrigger.FlagTrigger.Id,
		EnvironmentId: *environmentID,
	})
	if enabledTrigger.FlagTrigger.TriggerCount != 1 {
		t.Fatalf("unexpected trigger times: %d", enabledTrigger.FlagTrigger.TriggerCount)
	}
	if enabledTrigger.FlagTrigger.LastTriggeredAt == 0 {
		t.Fatal("unexpected last triggered at")
	}
	// Create Disable flag triggers
	disableTrigger, err := client.CreateFlagTrigger(context.Background(), newCreateFlagTriggerReq(
		command.Id,
		newTriggerDescription(t),
		featureproto.FlagTrigger_Action_OFF,
	))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("disable trigger URL: %s", disableTrigger.GetUrl())
	// Send post request
	resp, err = sendPostRequestIgnoreSSL(disableTrigger.GetUrl())
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %s", resp.Status)
	}
	disabledFeature := getFeature(t, command.Id, client)
	if disabledFeature.Enabled != false {
		t.Fatalf("unexpected enabled: %v", disabledFeature.Enabled)
	}
	disabledTrigger := getFeatureFlagTrigger(t, client, &featureproto.GetFlagTriggerRequest{
		Id:            disableTrigger.FlagTrigger.Id,
		EnvironmentId: *environmentID,
	})
	if disabledTrigger.FlagTrigger.TriggerCount != 1 {
		t.Fatalf("unexpected trigger times: %d", disabledTrigger.FlagTrigger.TriggerCount)
	}
	if disabledTrigger.FlagTrigger.LastTriggeredAt == 0 {
		t.Fatal("unexpected last triggered at")
	}
}

func sendPostRequestIgnoreSSL(targetURL string) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func newCreateFlagTriggerReq(
	featureID string,
	description string,
	action featureproto.FlagTrigger_Action,
) *featureproto.CreateFlagTriggerRequest {
	return &featureproto.CreateFlagTriggerRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		Type:          featureproto.FlagTrigger_Type_WEBHOOK,
		Action:        action,
		Description:   description,
	}
}

func getFeatureFlagTrigger(
	t *testing.T,
	client featureproto.FeatureServiceClient,
	req *featureproto.GetFlagTriggerRequest,
) *featureproto.GetFlagTriggerResponse {
	t.Helper()
	resp, err := client.GetFlagTrigger(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func createFeatureFlagTrigger(
	t *testing.T,
	client featureproto.FeatureServiceClient,
	req *featureproto.CreateFlagTriggerRequest,
) *featureproto.CreateFlagTriggerResponse {
	t.Helper()
	resp, err := client.CreateFlagTrigger(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func newTriggerDescription(t *testing.T) string {
	t.Helper()
	now := time.Now()
	salt, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	salted := fmt.Sprintf("%s-%v", salt.String(), now.Unix())
	if *testID != "" {
		return fmt.Sprintf("%s-%s-%v-trigger-description", prefixID, *testID, salted)
	}
	return fmt.Sprintf("%s-%v-trigger-description", prefixID, salted)
}
