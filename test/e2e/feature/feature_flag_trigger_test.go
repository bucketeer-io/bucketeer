// Copyright 2023 The Bucketeer Authors.
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
	"testing"

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestCreateFeatureFlagTrigger(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	cmd := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	// Create flag trigger
	createFlagTriggerCommand := newCreateFlagTriggerCmd(cmd, "create flag trigger test")
	resp := createFeatureFlagTrigger(t, client, createFlagTriggerCommand)
	if resp.FlagTrigger.FeatureId == cmd.Id {
		t.Fatal("unexpected feature id")
	}
	if resp.FlagTrigger.Type != createFlagTriggerCommand.Type {
		t.Fatal("unexpected type")
	}
	if resp.FlagTrigger.Action != createFlagTriggerCommand.Action {
		t.Fatal("unexpected action")
	}
	if resp.FlagTrigger.Description != createFlagTriggerCommand.Description {
		t.Fatal("unexpected description")
	}
}

func TestUpdateFlagTrigger(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	createFeature(t, client, newCreateFeatureCommand(newFeatureID(t)))
	// Create flag trigger
	createFlagTriggerCommand := newCreateFlagTriggerCmd(
		newCreateFeatureCommand(newFeatureID(t)),
		"create flag trigger test",
	)
	createResp := createFeatureFlagTrigger(t, client, createFlagTriggerCommand)
	// Update flag trigger
	updateFlagTriggerReq := &featureproto.UpdateFlagTriggerRequest{
		EnvironmentNamespace: *environmentNamespace,
		ChangeFlagTriggerDescriptionCommand: &featureproto.ChangeFlagTriggerDescriptionCommand{
			Id:          createResp.FlagTrigger.Id,
			Description: "change flag trigger description test",
		},
	}
	_, err := client.UpdateFlagTrigger(context.Background(), updateFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	// Get flag trigger
	getFlagTriggerReq := &featureproto.GetFlagTriggerRequest{
		Id:                   createResp.FlagTrigger.Id,
		EnvironmentNamespace: *environmentNamespace,
	}
	resp := getFeatureFlagTrigger(t, client, getFlagTriggerReq)
	if resp.FlagTrigger.Description != updateFlagTriggerReq.ChangeFlagTriggerDescriptionCommand.Description {
		t.Fatal("unexpected description")
	}
}

func TestDisableEnableFlagTrigger(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	createFeature(t, client, newCreateFeatureCommand(newFeatureID(t)))
	// Create flag trigger
	createFlagTriggerCommand := newCreateFlagTriggerCmd(
		newCreateFeatureCommand(newFeatureID(t)),
		"create flag trigger test",
	)
	createResp := createFeatureFlagTrigger(t, client, createFlagTriggerCommand)
	// Disable flag trigger
	disableFlagTriggerReq := &featureproto.DisableFlagTriggerRequest{
		EnvironmentNamespace: *environmentNamespace,
		DisableFlagTriggerCommand: &featureproto.DisableFlagTriggerCommand{
			Id: createResp.FlagTrigger.Id,
		},
	}
	_, err := client.DisableFlagTrigger(context.Background(), disableFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	getFlagTriggerReq := &featureproto.GetFlagTriggerRequest{
		Id:                   createResp.FlagTrigger.Id,
		EnvironmentNamespace: *environmentNamespace,
	}
	// Get flag trigger
	resp := getFeatureFlagTrigger(t, client, getFlagTriggerReq)
	if resp.FlagTrigger.Disabled != true {
		t.Fatal("unexpected disabled")
	}
	// Enable flag trigger
	enableFlagTriggerReq := &featureproto.EnableFlagTriggerRequest{
		EnvironmentNamespace: *environmentNamespace,
		EnableFlagTriggerCommand: &featureproto.EnableFlagTriggerCommand{
			Id: createResp.FlagTrigger.Id,
		},
	}
	_, err = client.EnableFlagTrigger(context.Background(), enableFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	// Get flag trigger
	resp = getFeatureFlagTrigger(t, client, getFlagTriggerReq)
	if resp.FlagTrigger.Disabled != false {
		t.Fatal("unexpected disabled")
	}
}

func TestResetFlagTrigger(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	createFeature(t, client, newCreateFeatureCommand(newFeatureID(t)))
	// Create flag trigger
	createFlagTriggerCommand := newCreateFlagTriggerCmd(
		newCreateFeatureCommand(newFeatureID(t)),
		"create flag trigger test",
	)
	createResp := createFeatureFlagTrigger(t, client, createFlagTriggerCommand)
	// Reset flag trigger
	resetFlagTriggerReq := &featureproto.ResetFlagTriggerRequest{
		EnvironmentNamespace: *environmentNamespace,
		ResetFlagTriggerCommand: &featureproto.ResetFlagTriggerCommand{
			Id: createResp.FlagTrigger.Id,
		},
	}
	_, err := client.ResetFlagTrigger(context.Background(), resetFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	// Get flag trigger
	getFlagTriggerReq := &featureproto.GetFlagTriggerRequest{
		Id:                   createResp.FlagTrigger.Id,
		EnvironmentNamespace: *environmentNamespace,
	}
	resp := getFeatureFlagTrigger(t, client, getFlagTriggerReq)
	if resp.FlagTrigger.Uuid == createResp.FlagTrigger.Uuid {
		t.Fatal("unexpected uuid")
	}
	if resp.Url == createResp.Url {
		t.Fatal("unexpected url")
	}
}

func TestDeleteFlagTrigger(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	createFeature(t, client, newCreateFeatureCommand(newFeatureID(t)))
	// Create flag trigger
	createFlagTriggerCommand := newCreateFlagTriggerCmd(
		newCreateFeatureCommand(newFeatureID(t)),
		"create flag trigger test",
	)
	createResp := createFeatureFlagTrigger(t, client, createFlagTriggerCommand)
	// Delete flag trigger
	deleteFlagTriggerReq := &featureproto.DeleteFlagTriggerRequest{
		EnvironmentNamespace: *environmentNamespace,
		DeleteFlagTriggerCommand: &featureproto.DeleteFlagTriggerCommand{
			Id: createResp.FlagTrigger.Id,
		},
	}
	_, err := client.DeleteFlagTrigger(context.Background(), deleteFlagTriggerReq)
	if err != nil {
		t.Fatal(err)
	}
	// Get flag trigger
	getFlagTriggerReq := &featureproto.GetFlagTriggerRequest{
		Id:                   createResp.FlagTrigger.Id,
		EnvironmentNamespace: *environmentNamespace,
	}
	resp := getFeatureFlagTrigger(t, client, getFlagTriggerReq)
	if resp.FlagTrigger.Deleted != true {
		t.Fatal("unexpected deleted")
	}
}

func TestListFlagTriggers(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Create feature
	command := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, command)
	// Create flag triggers
	trigger1, err := client.CreateFlagTrigger(context.Background(), &featureproto.CreateFlagTriggerRequest{
		EnvironmentNamespace:     *environmentNamespace,
		CreateFlagTriggerCommand: newCreateFlagTriggerCmd(command, "create flag trigger test 1"),
	})
	if err != nil {
		t.Fatal(err)
	}
	trigger2, err := client.CreateFlagTrigger(context.Background(), &featureproto.CreateFlagTriggerRequest{
		EnvironmentNamespace:     *environmentNamespace,
		CreateFlagTriggerCommand: newCreateFlagTriggerCmd(command, "create flag trigger test 2"),
	})
	if err != nil {
		t.Fatal(err)
	}
	// List flag triggers
	listFlagTriggersReq := &featureproto.ListFlagTriggersRequest{
		FeatureId:            command.Id,
		EnvironmentNamespace: *environmentNamespace,
		Cursor:               0,
		PageSize:             10,
		OrderBy:              featureproto.ListFlagTriggersRequest_CREATED_AT,
		OrderDirection:       featureproto.ListFlagTriggersRequest_ASC,
	}
	triggers, err := client.ListFlagTriggers(context.Background(), listFlagTriggersReq)
	if err != nil {
		t.Fatal(err)
	}
	if len(triggers.FlagTriggers) != 2 {
		t.Fatal("unexpected length")
	}
	if triggers.FlagTriggers[0].FlagTrigger.GetId() != trigger1.FlagTrigger.Id {
		t.Fatal("unexpected id")
	}
	if triggers.FlagTriggers[1].FlagTrigger.GetId() != trigger2.FlagTrigger.Id {
		t.Fatal("unexpected id")
	}
	if triggers.TotalCount != 2 {
		t.Fatal("unexpected total count")
	}
	if triggers.Cursor != "2" {
		t.Fatal("unexpected cursor")
	}
}

func newCreateFlagTriggerCmd(
	cmd *featureproto.CreateFeatureCommand,
	description string,
) *featureproto.CreateFlagTriggerCommand {
	createFlagTriggerCommand := &featureproto.CreateFlagTriggerCommand{
		FeatureId:   cmd.Id,
		Type:        featureproto.FlagTrigger_Type_WEBHOOK,
		Action:      featureproto.FlagTrigger_Action_ON,
		Description: description,
	}
	return createFlagTriggerCommand
}

func getFeatureFlagTrigger(
	t *testing.T,
	client featureproto.FeatureServiceClient,
	cmd *featureproto.GetFlagTriggerRequest,
) *featureproto.GetFlagTriggerResponse {
	t.Helper()
	resp, err := client.GetFlagTrigger(context.Background(), cmd)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func createFeatureFlagTrigger(
	t *testing.T,
	client featureproto.FeatureServiceClient,
	cmd *featureproto.CreateFlagTriggerCommand,
) *featureproto.CreateFlagTriggerResponse {
	t.Helper()
	resp, err := client.CreateFlagTrigger(context.Background(), &featureproto.CreateFlagTriggerRequest{
		EnvironmentNamespace:     *environmentNamespace,
		CreateFlagTriggerCommand: cmd,
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp
}
