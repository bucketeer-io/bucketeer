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
	"flag"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aoclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	aoproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	"github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
	"github.com/bucketeer-io/bucketeer/test/util"
)

const (
	prefixID = "e2e-test"
	timeout  = 10 * time.Second
)

var (
	// FIXME: To avoid compiling the test many times, webGatewayAddr, webGatewayPort & apiKey has been also added here to prevent from getting: "flag provided but not defined" error during the test. These 3 are being use  in the Gateway test
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

	tags = []string{"e2e-test-tag-1", "e2e-test-tag-2", "e2e-test-tag-3"}
)

func TestCreateFeature(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	cmd := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	f := getFeature(t, cmd.Id, client)
	if cmd.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if cmd.Name != f.Name {
		t.Fatalf("Different names. Expected: %s actual: %s", cmd.Name, f.Name)
	}
	if cmd.Description != f.Description {
		t.Fatalf("Different descriptions. Expected: %s actual: %s", cmd.Description, f.Description)
	}
	if f.Enabled {
		t.Fatalf("Enabled flag is true")
	}
	for i := range f.Variations {
		compareVariation(t, cmd.Variations[i], f.Variations[i])
	}
	if !reflect.DeepEqual(cmd.Tags, f.Tags) {
		t.Fatalf("Different tags. Expected: %v actual: %v: ", cmd.Tags, f.Tags)
	}
	defaultOnVariation := findVariation(f.DefaultStrategy.FixedStrategy.Variation, f.Variations)
	cmdDefaultOnVariation := cmd.Variations[int(cmd.DefaultOnVariationIndex.Value)]
	if cmdDefaultOnVariation.Value != defaultOnVariation.Value {
		t.Fatalf("Different default on variation value. Expected: %s actual: %s", cmdDefaultOnVariation.Value, defaultOnVariation.Value)
	}
	defaultOffVariation := findVariation(f.OffVariation, f.Variations)
	cmdDefaultOffVariation := cmd.Variations[int(cmd.DefaultOffVariationIndex.Value)]
	if cmdDefaultOffVariation.Value != defaultOffVariation.Value {
		t.Fatalf("Different default off variation value. Expected: %s actual: %s", cmdDefaultOffVariation.Value, defaultOffVariation.Value)
	}
}

func TestArchiveFeature(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	req := &feature.ArchiveFeatureRequest{
		Id:                   featureID,
		Command:              &feature.ArchiveFeatureCommand{},
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.ArchiveFeature(ctx, req); err != nil {
		t.Fatal(err)
	}
	f := getFeature(t, featureID, client)
	if cmd.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if !f.Archived {
		t.Fatal("Delete flag is false")
	}
}

func TestUnarchiveFeature(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &feature.ArchiveFeatureRequest{
		Id:                   featureID,
		Command:              &feature.ArchiveFeatureCommand{},
		EnvironmentNamespace: *environmentNamespace,
	}
	if _, err := client.ArchiveFeature(ctx, req); err != nil {
		t.Fatal(err)
	}
	reqUnarchive := &feature.UnarchiveFeatureRequest{
		Id:                   featureID,
		Command:              &feature.UnarchiveFeatureCommand{},
		EnvironmentNamespace: *environmentNamespace,
	}
	if _, err := client.UnarchiveFeature(ctx, reqUnarchive); err != nil {
		t.Fatal(err)
	}
	f := getFeature(t, featureID, client)
	if cmd.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if f.Archived {
		t.Fatal("Delete flag is true")
	}
}

func TestDeleteFeature(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	deleteReq := &feature.DeleteFeatureRequest{
		Id:                   featureID,
		Command:              &feature.DeleteFeatureCommand{},
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.DeleteFeature(ctx, deleteReq); err != nil {
		t.Fatal(err)
	}
	f := getFeature(t, featureID, client)
	if cmd.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if !f.Deleted {
		t.Fatal("Delete flag is false")
	}
}

func TestEnableFeature(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	cmd := newCreateFeatureWithTwoVariationsCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	f := getFeature(t, cmd.Id, client)
	aoCLient := newAutoOpsClient(t)
	schedules := []*aoproto.ProgressiveRolloutSchedule{
		{
			Weight:    20000,
			ExecuteAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}
	createProgressiveRollout(
		ctx,
		t,
		aoCLient,
		cmd.Id,
		&aoproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: f.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, aoCLient, cmd.Id)
	if len(progressiveRollouts) != 1 {
		t.Fatal("Progressive rollout list shouldn't be empty")
	}
	enableFeature(t, cmd.Id, client)
	f = getFeature(t, cmd.Id, client)
	if cmd.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if !f.Enabled {
		t.Fatal("Enabled flag is false")
	}
	// As a requirement, when disabling a flag,
	// It must stop the progressive rollout if it is running.
	// In this case, we ensure that the status didn't change after enabling the flag.
	pr := getProgressiveRollout(t, aoCLient, progressiveRollouts[0].Id)
	if pr.Status != aoproto.ProgressiveRollout_WAITING {
		t.Fatalf("Progressive rollout must be in waiting status. Current status: %v", pr.Status)
	}
}

func TestUpdateTargetingEnableFeature(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureWithTwoVariationsCommand(featureID)
	createFeature(t, client, cmd)
	f := getFeature(t, cmd.Id, client)
	aoCLient := newAutoOpsClient(t)
	schedules := []*aoproto.ProgressiveRolloutSchedule{
		{
			Weight:    20000,
			ExecuteAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}
	createProgressiveRollout(
		ctx,
		t,
		aoCLient,
		cmd.Id,
		&aoproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: f.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, aoCLient, cmd.Id)
	if len(progressiveRollouts) != 1 {
		t.Fatal("Progressive rollout list shouldn't be empty")
	}
	f = getFeature(t, featureID, client)
	disableCmd, _ := util.MarshalCommand(&feature.EnableFeatureCommand{})
	updateFeatureTargeting(t, client, disableCmd, featureID)
	f = getFeature(t, cmd.Id, client)
	if !f.Enabled {
		t.Fatal("Flag must be enabled")
	}
	// As a requirement, when disabling a flag,
	// It must stop the progressive rollout if it is running.
	// In this case, we ensure that the status didn't change after enabling the flag.
	pr := getProgressiveRollout(t, aoCLient, progressiveRollouts[0].Id)
	if pr.Status != aoproto.ProgressiveRollout_WAITING {
		t.Fatalf("Progressive rollout must be in waiting status. Current status: %v", pr.Status)
	}
}

func TestDisableFeature(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	cmd := newCreateFeatureWithTwoVariationsCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	enableFeature(t, cmd.Id, client)
	f := getFeature(t, cmd.Id, client)
	aoCLient := newAutoOpsClient(t)
	schedules := []*aoproto.ProgressiveRolloutSchedule{
		{
			Weight:    20000,
			ExecuteAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}
	createProgressiveRollout(
		ctx,
		t,
		aoCLient,
		cmd.Id,
		&aoproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: f.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, aoCLient, cmd.Id)
	if len(progressiveRollouts) != 1 {
		t.Fatal("Progressive rollout list shouldn't be empty")
	}
	disableReq := &feature.DisableFeatureRequest{
		Id:                   cmd.Id,
		Command:              &feature.DisableFeatureCommand{},
		EnvironmentNamespace: *environmentNamespace,
	}
	if _, err := client.DisableFeature(ctx, disableReq); err != nil {
		t.Fatal(err)
	}
	f = getFeature(t, cmd.Id, client)
	if cmd.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if f.Enabled {
		t.Fatal("Enabled flag is true")
	}
	// As a requirement, when disabling a flag using an auto operation,
	// It must stop the progressive rollout if it is running
	pr := getProgressiveRollout(t, aoCLient, progressiveRollouts[0].Id)
	if pr.Status != aoproto.ProgressiveRollout_STOPPED {
		t.Fatalf("Progressive rollout must be stopped. Current status: %v", pr.Status)
	}
}

func TestUpdateTargetingDisableFeature(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureWithTwoVariationsCommand(featureID)
	createFeature(t, client, cmd)
	enableFeature(t, cmd.Id, client)
	f := getFeature(t, cmd.Id, client)
	aoCLient := newAutoOpsClient(t)
	schedules := []*aoproto.ProgressiveRolloutSchedule{
		{
			Weight:    20000,
			ExecuteAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}
	createProgressiveRollout(
		ctx,
		t,
		aoCLient,
		cmd.Id,
		&aoproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: f.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, aoCLient, cmd.Id)
	if len(progressiveRollouts) != 1 {
		t.Fatal("Progressive rollout list shouldn't be empty")
	}
	f = getFeature(t, featureID, client)
	disableCmd, _ := util.MarshalCommand(&feature.DisableFeatureCommand{})
	updateFeatureTargeting(t, client, disableCmd, featureID)
	f = getFeature(t, cmd.Id, client)
	if f.Enabled {
		t.Fatal("Flag must be disabled")
	}
	// As a requirement, when disabling a flag using an auto operation,
	// It must stop the progressive rollout if it is running
	pr := getProgressiveRollout(t, aoCLient, progressiveRollouts[0].Id)
	if pr.Status != aoproto.ProgressiveRollout_STOPPED {
		t.Fatalf("Progressive rollout must be stopped. Current status: %v", pr.Status)
	}
}

func TestListArchivedFeatures(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	size := int64(1)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	req := &feature.ArchiveFeatureRequest{
		Id:                   featureID,
		Command:              &feature.ArchiveFeatureCommand{},
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.ArchiveFeature(ctx, req); err != nil {
		t.Fatal(err)
	}
	listReq := &feature.ListFeaturesRequest{
		PageSize:             size,
		Archived:             &wrappers.BoolValue{Value: true},
		EnvironmentNamespace: *environmentNamespace,
	}
	response, err := client.ListFeatures(ctx, listReq)
	if err != nil {
		t.Fatal(err)
	}
	responseSize := int64(len(response.Features))
	if responseSize != size {
		t.Fatalf("Different sizes. Expected: %d actual: %d", size, responseSize)
	}
}

func TestListFeaturesPageSize(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	size := int64(1)
	createRandomIDFeatures(t, 2, client)
	listReq := &feature.ListFeaturesRequest{
		PageSize:             size,
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := client.ListFeatures(ctx, listReq)
	if err != nil {
		t.Fatal(err)
	}
	responseSize := int64(len(response.Features))
	if responseSize != size {
		t.Fatalf("Different sizes. Expected: %d actual: %d", size, responseSize)
	}
}

func TestListFeaturesCursor(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	createRandomIDFeatures(t, 3, client)
	size := int64(1)
	listReq := &feature.ListFeaturesRequest{
		PageSize:             size,
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := client.ListFeatures(ctx, listReq)
	if err != nil {
		t.Fatal(err)
	}
	if response.Cursor == "" {
		t.Fatal("Cursor is empty")
	}
	features := response.Features
	listReq = &feature.ListFeaturesRequest{
		PageSize:             size,
		Cursor:               response.Cursor,
		EnvironmentNamespace: *environmentNamespace,
	}
	response, err = client.ListFeatures(ctx, listReq)
	if err != nil {
		t.Fatal(err)
	}
	require.EqualValues(t, size, len(features))
	require.EqualValues(t, size, len(response.Features))
	require.NotEqual(t, features[0].Id, response.Features[0].Id)
}

func TestListFeaturesOrderByName(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	size := int64(3)
	createRandomIDFeatures(t, 3, client)

	testcases := []struct {
		orderDirection  feature.ListFeaturesRequest_OrderDirection
		checkSortedFunc func(a []string) bool
	}{
		{
			orderDirection:  feature.ListFeaturesRequest_ASC,
			checkSortedFunc: sort.StringsAreSorted,
		},
		{
			orderDirection:  feature.ListFeaturesRequest_DESC,
			checkSortedFunc: util.StringsAreReverseSorted,
		},
	}

	for _, tc := range testcases {
		listReq := &feature.ListFeaturesRequest{
			PageSize:             size,
			OrderBy:              feature.ListFeaturesRequest_NAME,
			OrderDirection:       tc.orderDirection,
			EnvironmentNamespace: *environmentNamespace,
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		response, err := client.ListFeatures(ctx, listReq)
		if err != nil {
			t.Fatal(err)
		}
		names := make([]string, 0, len(response.Features))
		for _, f := range response.Features {
			names = append(names, f.Name)
		}
		if !tc.checkSortedFunc(names) {
			t.Fatalf("Features aren't sorted by Name %s. Features: %v", tc.orderDirection, response.Features)
		}
	}
}

func TestListFeaturesOrderByCreatedAt(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	size := int64(3)
	createRandomIDFeatures(t, 3, client)

	testcases := []struct {
		orderDirection  feature.ListFeaturesRequest_OrderDirection
		checkSortedFunc func(a []int64) bool
	}{
		{
			orderDirection:  feature.ListFeaturesRequest_ASC,
			checkSortedFunc: util.Int64sAreSorted,
		},
		{
			orderDirection:  feature.ListFeaturesRequest_DESC,
			checkSortedFunc: util.Int64sAreReverseSorted,
		},
	}

	for _, tc := range testcases {
		listReq := &feature.ListFeaturesRequest{
			PageSize:             size,
			OrderBy:              feature.ListFeaturesRequest_CREATED_AT,
			OrderDirection:       tc.orderDirection,
			EnvironmentNamespace: *environmentNamespace,
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		response, err := client.ListFeatures(ctx, listReq)
		if err != nil {
			t.Fatal(err)
		}
		createdAts := make([]int64, 0, len(response.Features))
		for _, f := range response.Features {
			createdAts = append(createdAts, f.CreatedAt)
		}
		if !tc.checkSortedFunc(createdAts) {
			t.Fatalf("Features aren't sorted by CreatedAt %s. Features: %v", tc.orderDirection, response.Features)
		}
	}
}

func TestListFeaturesOrderByUpdatedAt(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	size := int64(3)
	createRandomIDFeatures(t, 3, client)

	testcases := []struct {
		orderDirection  feature.ListFeaturesRequest_OrderDirection
		checkSortedFunc func(a []int64) bool
	}{
		{
			orderDirection:  feature.ListFeaturesRequest_ASC,
			checkSortedFunc: util.Int64sAreSorted,
		},
		{
			orderDirection:  feature.ListFeaturesRequest_DESC,
			checkSortedFunc: util.Int64sAreReverseSorted,
		},
	}

	for _, tc := range testcases {
		listReq := &feature.ListFeaturesRequest{
			PageSize:             size,
			OrderBy:              feature.ListFeaturesRequest_UPDATED_AT,
			OrderDirection:       tc.orderDirection,
			EnvironmentNamespace: *environmentNamespace,
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		response, err := client.ListFeatures(ctx, listReq)
		if err != nil {
			t.Fatal(err)
		}
		updatedAts := make([]int64, 0, len(response.Features))
		for _, f := range response.Features {
			updatedAts = append(updatedAts, f.UpdatedAt)
		}
		if !tc.checkSortedFunc(updatedAts) {
			t.Fatalf("Features aren't sorted by UpdatedAt %s. Features: %v", tc.orderDirection, response.Features)
		}
	}
}

func TestListEnabledFeaturesPageSize(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	ids := []string{newFeatureID(t), newFeatureID(t), newFeatureID(t)}
	createFeatures(t, ids, client)
	enableFeatures(t, ids, client)
	size := int64(2)
	listReq := &feature.ListEnabledFeaturesRequest{
		PageSize:             size,
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := client.ListEnabledFeatures(ctx, listReq)
	if err != nil {
		t.Fatal(err)
	}
	responseSize := int64(len(response.Features))
	if responseSize != size {
		t.Fatalf("Different sizes. Expected: %d actual: %d", size, responseSize)
	}
	checkEnabledFlag(t, response.Features)
}

func TestListEnabledFeaturesCursor(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	ids := []string{newFeatureID(t), newFeatureID(t), newFeatureID(t), newFeatureID(t)}
	createFeatures(t, ids, client)
	enableFeatures(t, ids, client)
	size := int64(2)
	listReq := &feature.ListEnabledFeaturesRequest{
		PageSize:             size,
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := client.ListEnabledFeatures(ctx, listReq)
	if err != nil {
		t.Fatal(err)
	}
	if response.Cursor == "" {
		t.Fatal("Cursor is empty")
	}
	features := response.Features
	firstPageIds := make([]string, 0, len(features))
	for _, feature := range features {
		firstPageIds = append(firstPageIds, feature.Id)
	}
	checkEnabledFlag(t, features)
	listReq = &feature.ListEnabledFeaturesRequest{
		PageSize:             size,
		Cursor:               response.Cursor,
		EnvironmentNamespace: *environmentNamespace,
	}
	response, err = client.ListEnabledFeatures(ctx, listReq)
	if err != nil {
		t.Fatal(err)
	}
	checkEnabledFlag(t, features)
	for _, feature := range response.Features {
		// TODO: Features should be tagged while creating and then check the returned features are in that created list.
		// assert.Contains(t, ids, feature.Id)
		assert.NotContains(t, firstPageIds, feature.Id)
	}
}

func TestRename(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	cmd := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	expected := "new-feature-name"
	updateReq := &feature.UpdateFeatureDetailsRequest{
		Id:                   cmd.Id,
		RenameFeatureCommand: &feature.RenameFeatureCommand{Name: expected},
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureDetails(ctx, updateReq); err != nil {
		t.Fatal(err)
	}
	f := getFeature(t, cmd.Id, client)
	if cmd.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if expected != f.Name {
		t.Fatalf("Different names. Expected: %s actual: %s", expected, f.Name)
	}
}

func TestChangeDescription(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	cmd := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	expected := "new-feature-description"
	updateReq := &feature.UpdateFeatureDetailsRequest{
		Id:                       cmd.Id,
		ChangeDescriptionCommand: &feature.ChangeDescriptionCommand{Description: expected},
		EnvironmentNamespace:     *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureDetails(ctx, updateReq); err != nil {
		t.Fatal(err)
	}
	f := getFeature(t, cmd.Id, client)
	if cmd.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if expected != f.Description {
		t.Fatalf("Different names. Expected: %s actual: %s", expected, f.Description)
	}
}

func TestAddTags(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	cmd := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	newTags := []string{"e2e-test-tag-4", "e2e-test-tag-5", "e2e-test-tag-6"}
	addReq := &feature.UpdateFeatureDetailsRequest{
		Id: cmd.Id,
		AddTagCommands: []*feature.AddTagCommand{
			{Tag: newTags[0]},
			{Tag: newTags[1]},
			{Tag: newTags[2]},
		},
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureDetails(ctx, addReq); err != nil {
		t.Fatal(err)
	}
	f := getFeature(t, cmd.Id, client)
	if cmd.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	cmd.Tags = append(cmd.Tags, newTags...)
	if !reflect.DeepEqual(cmd.Tags, f.Tags) {
		t.Fatalf("Different tags. Expected: %v actual: %v: ", cmd.Tags, f.Tags)
	}
}

func TestRemoveTags(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	cmd := newCreateFeatureCommand(newFeatureID(t))
	createFeature(t, client, cmd)
	removeTargetTags := []*feature.RemoveTagCommand{
		{Tag: cmd.Tags[0]},
		{Tag: cmd.Tags[2]},
	}
	removeReq := &feature.UpdateFeatureDetailsRequest{
		Id:                   cmd.Id,
		RemoveTagCommands:    removeTargetTags,
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureDetails(ctx, removeReq); err != nil {
		t.Fatal(err)
	}
	f := getFeature(t, cmd.Id, client)
	if cmd.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if len(f.Tags) != 1 {
		t.Fatalf("Tags should have only 1 element. Expected: %s actual: %v", cmd.Tags[1], f.Tags)
	}
	if f.Tags[0] != cmd.Tags[1] {
		t.Fatalf("The wrong tag might has been deleted. Expected to be deleted: %v actual: %v", removeTargetTags, f.Tags)
	}
}

func TestAddVariation(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	targetVariationValues := []string{newUUID(t), newUUID(t)}
	targetVariations := newVariations(targetVariationValues)
	addCmd := newAddVariationsCommand(t, targetVariations)
	updateVariations(t, featureID, addCmd, client)
	feature := getFeature(t, featureID, client)
	if feature.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, feature.Id)
	}
	var matched int
	for _, e := range targetVariations {
		for _, g := range feature.Variations {
			if e.Value == g.Value {
				compareVariation(t, e, g)
				matched++
			}
		}
	}
	size := len(targetVariations)
	if matched != size {
		t.Fatalf("The number of variations added does not match. Expected: %d actual: %d", size, matched)
	}
}

func TestRemoveVariation(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	targetVariationID := getFeature(t, featureID, client).Variations[2].Id
	removeCmd := newRemoveVariationsCommand(t, []string{targetVariationID})
	updateVariations(t, featureID, removeCmd, client)
	feature := getFeature(t, featureID, client)
	if feature.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, feature.Id)
	}
	if len(feature.Variations) != 3 {
		t.Fatal("Variations should have 3 elements. Actual:", feature.Variations)
	}
	if findVariation(targetVariationID, feature.Variations) != nil {
		t.Fatalf("The wrong variation might has been deleted. Expected: %s actual: %v", targetVariationID, feature.Variations)
	}
}

func TestRemoveVariations(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	feature := getFeature(t, featureID, client)
	targetVariationIDS := []string{feature.Variations[2].Id, feature.Variations[3].Id}
	cmds := newRemoveVariationsCommand(t, targetVariationIDS)
	updateVariations(t, featureID, cmds, client)
	feature = getFeature(t, featureID, client)
	if feature.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, feature.Id)
	}
	if len(feature.Variations) != 2 {
		t.Fatal("Variations should have only 2 elements. Actual:", feature.Variations)
	}
	if variation := findOneOfVariations(targetVariationIDS, feature.Variations); variation != nil {
		t.Fatalf("The wrong variation might has been deleted. Expected: %v actual: %v", targetVariationIDS, feature.Variations)
	}
}

func TestChangeVariationValue(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	targetVariationID := getFeature(t, featureID, client).Variations[1].Id
	targetVariationValue := "new-variation-value"
	changeCmd, err := util.MarshalCommand(&feature.ChangeVariationValueCommand{
		Id:    targetVariationID,
		Value: targetVariationValue,
	})
	if err != nil {
		t.Fatal(err)
	}
	updateVariations(t, featureID, []*feature.Command{{Command: changeCmd}}, client)
	f := getFeature(t, featureID, client)
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	expected := &feature.Variation{
		Value:       targetVariationValue,
		Name:        cmd.Variations[1].Name,
		Description: cmd.Variations[1].Description,
	}
	compareVariation(t, expected, f.Variations[1])
}

func TestChangeVariationName(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	targetVariationID := getFeature(t, featureID, client).Variations[1].Id
	targetVariationName := "new-variation-name"
	changeCmd, err := util.MarshalCommand(&feature.ChangeVariationNameCommand{
		Id:   targetVariationID,
		Name: targetVariationName,
	})
	if err != nil {
		t.Fatal(err)
	}
	updateVariations(t, featureID, []*feature.Command{{Command: changeCmd}}, client)
	f := getFeature(t, featureID, client)
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	expected := &feature.Variation{
		Value:       cmd.Variations[1].Value,
		Name:        targetVariationName,
		Description: cmd.Variations[1].Description,
	}
	compareVariation(t, expected, f.Variations[1])
}

func TestChangeVariationDescription(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	targetVariationID := getFeature(t, featureID, client).Variations[1].Id
	targetVariationDescription := "new-variation-description"
	changeCmd, err := util.MarshalCommand(&feature.ChangeVariationDescriptionCommand{
		Id:          targetVariationID,
		Description: targetVariationDescription,
	})
	if err != nil {
		t.Fatal(err)
	}
	updateVariations(t, featureID, []*feature.Command{{Command: changeCmd}}, client)
	f := getFeature(t, featureID, client)
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	expected := &feature.Variation{
		Value:       cmd.Variations[1].Value,
		Name:        cmd.Variations[1].Name,
		Description: targetVariationDescription,
	}
	compareVariation(t, expected, f.Variations[1])
}

func TestChangeFixedStrategy(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	f := getFeature(t, featureID, client)
	rule := newFixedStrategyRule(f.Variations[0].Id)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	expected := f.Variations[1].Id
	changeCmd, err := util.MarshalCommand(&feature.ChangeFixedStrategyCommand{
		Id:       featureID,
		RuleId:   rule.Id,
		Strategy: &feature.FixedStrategy{Variation: expected},
	})
	require.NoError(t, err)
	updateFeatureTargeting(t, client, changeCmd, featureID)
	f = getFeature(t, featureID, client)
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	actual := f.Rules[0].Strategy.FixedStrategy.Variation
	if expected != actual {
		t.Fatalf("Variation id is not equal. Expected: %s actual: %s", expected, actual)
	}
}

func TestChangeRolloutStrategy(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	f := getFeature(t, featureID, client)
	rule := newRolloutStrategyRule(f.Variations)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	expected := &feature.RolloutStrategy{
		Variations: []*feature.RolloutStrategy_Variation{
			{
				Variation: f.Variations[0].Id,
				Weight:    12000,
			},
			{
				Variation: f.Variations[1].Id,
				Weight:    30000,
			},
			{
				Variation: f.Variations[2].Id,
				Weight:    50000,
			},
			{
				Variation: f.Variations[3].Id,
				Weight:    8000,
			},
		},
	}
	changeCmd, err := util.MarshalCommand(&feature.ChangeRolloutStrategyCommand{
		Id:       featureID,
		RuleId:   rule.Id,
		Strategy: expected,
	})
	require.NoError(t, err)
	updateFeatureTargeting(t, client, changeCmd, featureID)
	f = getFeature(t, featureID, client)
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	actual := f.Rules[0].Strategy.RolloutStrategy
	if !proto.Equal(expected, actual) {
		t.Fatalf("Strategy is not equal. Expected: %s actual: %s", expected, actual)
	}
}

func TestChangeOffVariation(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	expected := getFeature(t, featureID, client).Variations[1].Id
	changeCmd, err := util.MarshalCommand(&feature.ChangeOffVariationCommand{Id: expected})
	if err != nil {
		t.Fatal(err)
	}
	updateFeatureTargeting(t, client, changeCmd, featureID)
	f := getFeature(t, featureID, client)
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if expected != f.OffVariation {
		t.Fatalf("Off variation does not match. Expected: %s actual: %s", expected, f.OffVariation)
	}
}

func TestAddUserToVariation(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	expected := "new-user"
	addCmd, err := util.MarshalCommand(&feature.AddUserToVariationCommand{
		Id:   getFeature(t, featureID, client).Variations[1].Id,
		User: expected,
	})
	if err != nil {
		t.Fatal(err)
	}
	updateFeatureTargeting(t, client, addCmd, featureID)
	f := getFeature(t, featureID, client)
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if expected != f.Targets[1].Users[0] {
		t.Fatalf("User does not match. Expected to be deleted: %s actual: %s", expected, f.Targets[1].Users[0])
	}
}

func TestRemoveUserToVariation(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	variationID := getFeature(t, featureID, client).Variations[1].Id
	expected := "new-user"
	addCmd, err := util.MarshalCommand(&feature.AddUserToVariationCommand{
		Id:   variationID,
		User: expected,
	})
	if err != nil {
		t.Fatal(err)
	}
	updateFeatureTargeting(t, client, addCmd, featureID)
	removeCmd, err := util.MarshalCommand(&feature.RemoveUserFromVariationCommand{
		Id:   variationID,
		User: expected,
	})
	if err != nil {
		t.Fatal(err)
	}
	updateFeatureTargeting(t, client, removeCmd, featureID)
	f := getFeature(t, featureID, client)
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if len(f.Targets[1].Users) > 0 {
		t.Fatalf("User was not deleted. Expected: %s actual: %s", expected, f.Targets[0].Users[0])
	}
}

func TestAddRule(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	rule := newFixedStrategyRule(getFeature(t, featureID, client).Variations[1].Id)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	f := getFeature(t, featureID, client)
	r := f.Rules[0]
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if !proto.Equal(rule.Strategy, r.Strategy) {
		t.Fatalf("Strategy is not equal. Expected: %v actual: %v", rule.Strategy, r.Strategy)
	}
	expectedSize := len(rule.Clauses)
	actualSize := len(r.Clauses)
	if expectedSize != actualSize {
		t.Fatalf("Clauses have different sizes. Expected: %d actual: %d", expectedSize, actualSize)
	}
	for i := range rule.Clauses {
		if r.Clauses[i].Id == "" {
			t.Fatalf("ID is empty")
		}
		compareClause(t, rule.Clauses[i], r.Clauses[i])
	}
}

func TestChangeRuleToFixedStrategy(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	f := getFeature(t, featureID, client)
	rule := newRolloutStrategyRule(f.Variations)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	expected := &feature.Strategy{
		Type:          feature.Strategy_FIXED,
		FixedStrategy: &feature.FixedStrategy{Variation: f.Variations[1].Id},
	}
	changeCmd, err := util.MarshalCommand(&feature.ChangeRuleStrategyCommand{
		Id:       featureID,
		RuleId:   rule.Id,
		Strategy: expected,
	})
	require.NoError(t, err)
	updateFeatureTargeting(t, client, changeCmd, featureID)
	f = getFeature(t, featureID, client)
	assert.Equal(t, cmd.Id, f.Id, fmt.Sprintf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id))
	actual := f.Rules[0].Strategy
	if !proto.Equal(expected, actual) {
		t.Fatalf("Strategy is not equal. Expected: %s actual: %s", expected, actual)
	}
}

func TestChangeRuleToRolloutStrategy(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	f := getFeature(t, featureID, client)
	rule := newFixedStrategyRule(f.Variations[0].Id)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	expected := &feature.Strategy{
		Type: feature.Strategy_ROLLOUT,
		RolloutStrategy: &feature.RolloutStrategy{
			Variations: []*feature.RolloutStrategy_Variation{
				{
					Variation: f.Variations[0].Id,
					Weight:    12000,
				},
				{
					Variation: f.Variations[1].Id,
					Weight:    30000,
				},
				{
					Variation: f.Variations[2].Id,
					Weight:    50000,
				},
				{
					Variation: f.Variations[3].Id,
					Weight:    8000,
				},
			},
		},
	}
	changeCmd, err := util.MarshalCommand(&feature.ChangeRuleStrategyCommand{
		Id:       featureID,
		RuleId:   rule.Id,
		Strategy: expected,
	})
	require.NoError(t, err)
	updateFeatureTargeting(t, client, changeCmd, featureID)
	f = getFeature(t, featureID, client)
	assert.Equal(t, cmd.Id, f.Id, fmt.Sprintf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id))
	actual := f.Rules[0].Strategy
	if !proto.Equal(expected, actual) {
		t.Fatalf("Strategy is not equal. Expected: %s actual: %s", expected, actual)
	}
}

func TestDeleteRule(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	rule := newFixedStrategyRule(getFeature(t, featureID, client).Variations[1].Id)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	f := getFeature(t, featureID, client)
	rule = f.Rules[0]
	removeCmd, _ := util.MarshalCommand(&feature.DeleteRuleCommand{Id: rule.Id})
	updateFeatureTargeting(t, client, removeCmd, featureID)
	r := getFeature(t, featureID, client).Rules
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	if len(r) > 0 {
		t.Fatalf("The wrong rule might has been delete. Expected: %v actual: %v", rule, r)
	}
}

func TestAddClause(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	rule := newFixedStrategyRule(getFeature(t, featureID, client).Variations[1].Id)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	r := getFeature(t, featureID, client).Rules[0]
	clause := newClause()
	addCmd, _ = util.MarshalCommand(&feature.AddClauseCommand{
		RuleId: r.Id,
		Clause: clause,
	})
	updateFeatureTargeting(t, client, addCmd, featureID)
	f := getFeature(t, featureID, client)
	r = f.Rules[0]
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	expectedSize := 3
	actualSize := len(r.Clauses)
	if expectedSize != actualSize {
		t.Fatalf("Clauses have different sizes. Expected: %d actual: %d", expectedSize, actualSize)
	}
	compareClause(t, clause, r.Clauses[2])
}

func TestDeleteClause(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	rule := newFixedStrategyRule(getFeature(t, featureID, client).Variations[1].Id)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	r := getFeature(t, featureID, client).Rules[0]
	expected := r.Clauses[1]
	removeCmd, _ := util.MarshalCommand(&feature.DeleteClauseCommand{
		Id:     expected.Id,
		RuleId: r.Id,
	})
	updateFeatureTargeting(t, client, removeCmd, featureID)
	f := getFeature(t, featureID, client)
	r = f.Rules[0]
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	expectedSize := 1
	actualSize := len(r.Clauses)
	if expectedSize != actualSize {
		t.Fatalf("Clauses have different sizes. Expected: %d actual: %d", expectedSize, actualSize)
	}
	if proto.Equal(expected, r.Clauses[0]) {
		t.Fatalf("The wrong clause might has been delete. Expected: %v actual: %v", expected, r.Clauses[0])
	}
}

func TestChangeClauseAttribute(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	rule := newFixedStrategyRule(getFeature(t, featureID, client).Variations[1].Id)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	r := getFeature(t, featureID, client).Rules[0]
	c := r.Clauses[1]
	expected := &feature.Clause{
		Attribute: "change-clause-attribute",
		Operator:  c.Operator,
		Values:    c.Values,
	}
	changeCmd, _ := util.MarshalCommand(&feature.ChangeClauseAttributeCommand{
		Id:        c.Id,
		RuleId:    r.Id,
		Attribute: expected.Attribute,
	})
	updateFeatureTargeting(t, client, changeCmd, featureID)
	f := getFeature(t, featureID, client)
	r = f.Rules[0]
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	compareClause(t, expected, r.Clauses[1])
}

func TestChangeClauseOperator(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	rule := newFixedStrategyRule(getFeature(t, featureID, client).Variations[1].Id)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	r := getFeature(t, featureID, client).Rules[0]
	c := r.Clauses[1]
	expected := &feature.Clause{
		Attribute: c.Attribute,
		Operator:  feature.Clause_EQUALS,
		Values:    c.Values,
	}
	changeCmd, _ := util.MarshalCommand(&feature.ChangeClauseOperatorCommand{
		Id:       c.Id,
		RuleId:   r.Id,
		Operator: expected.Operator,
	})
	updateFeatureTargeting(t, client, changeCmd, featureID)
	f := getFeature(t, featureID, client)
	r = f.Rules[0]
	if f.Id != cmd.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", cmd.Id, f.Id)
	}
	compareClause(t, expected, r.Clauses[1])
}

func TestAddClauseValue(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	rule := newFixedStrategyRule(getFeature(t, featureID, client).Variations[1].Id)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	r := getFeature(t, featureID, client).Rules[0]
	c := r.Clauses[1]
	values := append(c.Values, "new-value")
	expected := &feature.Clause{
		Attribute: c.Attribute,
		Operator:  c.Operator,
		Values:    values,
	}
	changeCmd, _ := util.MarshalCommand(&feature.AddClauseValueCommand{
		Id:     c.Id,
		RuleId: r.Id,
		Value:  expected.Values[2],
	})
	updateFeatureTargeting(t, client, changeCmd, featureID)
	f := getFeature(t, featureID, client)
	r = f.Rules[0]
	if f.Id != cmd.Id {
		t.Fatalf("different ids. expected: %s actual: %s", cmd.Id, f.Id)
	}
	compareClause(t, expected, r.Clauses[1])
}

func TestRemoveClauseValue(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	rule := newFixedStrategyRule(getFeature(t, featureID, client).Variations[1].Id)
	addCmd, _ := util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	r := getFeature(t, featureID, client).Rules[0]
	c := r.Clauses[0]
	expected := &feature.Clause{
		Attribute: c.Attribute,
		Operator:  c.Operator,
		Values:    []string{c.Values[0]},
	}
	removeCmd, _ := util.MarshalCommand(&feature.RemoveClauseValueCommand{
		Id:     c.Id,
		RuleId: r.Id,
		Value:  c.Values[1],
	})
	updateFeatureTargeting(t, client, removeCmd, featureID)
	f := getFeature(t, featureID, client)
	r = f.Rules[0]
	if f.Id != cmd.Id {
		t.Fatalf("different ids. expected: %s actual: %s", cmd.Id, f.Id)
	}
	expectedSize := 1
	actualSize := len(r.Clauses[0].Values)
	if expectedSize != actualSize {
		t.Fatalf("Values have different sizes. Expected: %d actual: %d", expectedSize, actualSize)
	}
	compareClause(t, expected, r.Clauses[0])
}

func TestEvaluateFeatures(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID1 := newFeatureID(t)
	cmd1 := newCreateFeatureCommand(featureID1)
	createFeature(t, client, cmd1)
	featureID2 := newFeatureID(t)
	cmd2 := newCreateFeatureCommand(featureID2)
	createFeature(t, client, cmd2)
	enableFeature(t, cmd2.Id, client)
	userID := "user-id-01"
	tag := tags[0]
	res := evaluateFeatures(t, client, userID, tag)
	if len(res.UserEvaluations.Evaluations) < 2 {
		t.Fatalf("length of user evaluations is not enough. Expected: >=%d, Actual: %d", 2, len(res.UserEvaluations.Evaluations))
	}
}

func TestEvaluateFeaturesWithEmptyTag(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID1 := newFeatureID(t)
	cmd1 := newCreateFeatureCommand(featureID1)
	createFeature(t, client, cmd1)
	featureID2 := newFeatureID(t)
	cmd2 := newCreateFeatureCommand(featureID2)
	createFeature(t, client, cmd2)
	enableFeature(t, cmd2.Id, client)
	userID := "user-id-01"
	res := evaluateFeatures(t, client, userID, "")
	if len(res.UserEvaluations.Evaluations) < 2 {
		t.Fatalf("length of user evaluations is not enough. Expected: >=%d, Actual: %d", 2, len(res.UserEvaluations.Evaluations))
	}
}

// TODO: implement the process to delete new environments so that we can run "TestCloneFeature"
/*
func TestCloneFeature(t *testing.T) {
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	enableFeature(t, featureID, client)
	f := getFeature(t, featureID, client)
	expected := "new-user"
	addCmd, err := util.MarshalCommand(&feature.AddUserToVariationCommand{
		Id:   f.Variations[1].Id,
		User: expected,
	})
	require.NoError(t, err)
	updateFeatureTargeting(t, client, addCmd, featureID)
	targetEnvironmentNamespace := newUUID(t)
	c := newEnvironmentClient(t)
	defer c.Close()
	envCmd := newEnvironmentCommand(targetEnvironmentNamespace)
	createEnvironment(t, c, envCmd)
	req := &feature.CloneFeatureRequest{
		Id: featureID,
		Command: &feature.CloneFeatureCommand{
			EnvironmentNamespace: targetEnvironmentNamespace,
		},
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.CloneFeature(ctx, req); err != nil {
		t.Fatal(err)
	}
	f = getFeature(t, featureID, client)
	cf := getClonedFeature(t, featureID, targetEnvironmentNamespace, client)
	if cf.Id != f.Id {
		t.Fatalf("Different ids. Expected: %s actual: %s", f.Id, cf.Id)
	}
	if cf.Name != f.Name {
		t.Fatalf("Different names. Expected: %s actual: %s", f.Name, cf.Name)
	}
	if cf.Description != f.Description {
		t.Fatalf("Different descriptions. Expected: %s actual: %s", f.Description, cf.Description)
	}
	if cf.Enabled {
		t.Fatalf("Enabled flag is true")
	}
	if cf.Targets[1].Users[0] != f.Targets[1].Users[0] {
		t.Fatalf("User does not match. Expected to be deleted: %s actual: %s", f.Targets[1].Users[0], cf.Targets[1].Users[0])
	}
	expectedVersion := int32(1)
	if cf.Version != expectedVersion {
		t.Fatalf("Different version. Expected: %d actual %d", expectedVersion, f.Version)
	}
	for i := range cf.Variations {
		compareVariation(t, f.Variations[i], cf.Variations[i])
	}
	if !reflect.DeepEqual(f.Tags, cf.Tags) {
		t.Fatalf("Different tags. Expected: %v actual: %v: ", f.Tags, cf.Tags)
	}
	featureDefaultOnVariation := findVariation(f.DefaultStrategy.FixedStrategy.Variation, f.Variations)
	clonedFeatureDefaultOnVariation := findVariation(cf.DefaultStrategy.FixedStrategy.Variation, cf.Variations)
	if clonedFeatureDefaultOnVariation.Value != featureDefaultOnVariation.Value {
		t.Fatalf("Different default on variation value. Expected: %s actual: %s", featureDefaultOnVariation.Value, clonedFeatureDefaultOnVariation.Value)
	}
	featureDefaultOffVariation := findVariation(f.OffVariation, f.Variations)
	clonedFeatureDefaultOffVariation := findVariation(cf.OffVariation, cf.Variations)
	if clonedFeatureDefaultOffVariation.Value != featureDefaultOffVariation.Value {
		t.Fatalf("Different default off variation value. Expected: %s actual: %s", featureDefaultOffVariation.Value, clonedFeatureDefaultOffVariation.Value)
	}
	for i := range cf.Rules {
		if cf.Rules[i].Strategy.FixedStrategy.Variation != f.Rules[i].Strategy.FixedStrategy.Variation {
			t.Fatalf("Different variation in rules. Expected: %s actual %s", f.Rules[i].Strategy.FixedStrategy.Variation, cf.Rules[i].Strategy.FixedStrategy.Variation)
		}
	}
	rule := newRolloutStrategyRule(f.Variations)
	addCmd, err = util.MarshalCommand(&feature.AddRuleCommand{Rule: rule})
	require.NoError(t, err)
	updateFeatureTargeting(t, client, addCmd, featureID)
	strategy := &feature.RolloutStrategy{
		Variations: []*feature.RolloutStrategy_Variation{
			{
				Variation: f.Variations[0].Id,
				Weight:    12000,
			},
			{
				Variation: f.Variations[1].Id,
				Weight:    30000,
			},
			{
				Variation: f.Variations[2].Id,
				Weight:    50000,
			},
			{
				Variation: f.Variations[3].Id,
				Weight:    8000,
			},
		},
	}
	changeCmd, err := util.MarshalCommand(&feature.ChangeRolloutStrategyCommand{
		Id:       featureID,
		RuleId:   rule.Id,
		Strategy: strategy,
	})
	require.NoError(t, err)
	updateFeatureTargeting(t, client, changeCmd, featureID)
	changeCmd, err = util.MarshalCommand(&feature.ChangeDefaultStrategyCommand{
		Strategy: &feature.Strategy{
			Type:            feature.Strategy_ROLLOUT,
			RolloutStrategy: strategy,
		},
	})
	require.NoError(t, err)
	updateFeatureTargeting(t, client, changeCmd, featureID)
	anotherTargetEnvironmentNamespace := newUUID(t)
	c = newEnvironmentClient(t)
	defer c.Close()
	envCmd = newEnvironmentCommand(anotherTargetEnvironmentNamespace)
	createEnvironment(t, c, envCmd)
	req = &feature.CloneFeatureRequest{
		Id: featureID,
		Command: &feature.CloneFeatureCommand{
			EnvironmentNamespace: anotherTargetEnvironmentNamespace,
		},
		EnvironmentNamespace: *environmentNamespace,
	}
	if _, err := client.CloneFeature(ctx, req); err != nil {
		t.Fatal(err)
	}
	f = getFeature(t, featureID, client)
	cf = getClonedFeature(t, featureID, anotherTargetEnvironmentNamespace, client)
	for i := range cf.Rules {
		for idx := range cf.Rules[i].Strategy.RolloutStrategy.Variations {
			if cf.Rules[i].Strategy.RolloutStrategy.Variations[idx].Weight != f.Rules[i].Strategy.RolloutStrategy.Variations[idx].Weight {
				t.Fatalf("Diffrent strategy on variation weight. Expected: %d actual: %d", f.Rules[i].Strategy.RolloutStrategy.Variations[idx].Weight, cf.Rules[i].Strategy.RolloutStrategy.Variations[idx].Weight)
			}
		}
	}
	for i := range cf.DefaultStrategy.RolloutStrategy.Variations {
		if cf.DefaultStrategy.RolloutStrategy.Variations[i].Weight != f.DefaultStrategy.RolloutStrategy.Variations[i].Weight {
			t.Fatalf("Different default on variation weight. Expected: %d actual %d", f.DefaultStrategy.RolloutStrategy.Variations[i].Weight, cf.DefaultStrategy.RolloutStrategy.Variations[i].Weight)
		}
	}
}
*/

func newFeatureID(t *testing.T) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-feature-id-%s", prefixID, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-feature-id-%s", prefixID, newUUID(t))
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func newFeatureClient(t *testing.T) featureclient.Client {
	t.Helper()
	creds, err := client.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	featureClient, err := featureclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create feature client:", err)
	}
	return featureClient
}

func newCreateFeatureCommand(featureID string) *feature.CreateFeatureCommand {
	return &feature.CreateFeatureCommand{
		Id:          featureID,
		Name:        "e2e-test-feature-name",
		Description: "e2e-test-feature-description",
		Variations: []*feature.Variation{
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
			{
				Value:       "C",
				Name:        "Variation C",
				Description: "Thing does C",
			},
			{
				Value:       "D",
				Name:        "Variation D",
				Description: "Thing does D",
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

func newCreateFeatureWithTwoVariationsCommand(featureID string) *feature.CreateFeatureCommand {
	return &feature.CreateFeatureCommand{
		Id:          featureID,
		Name:        "e2e-test-feature-name",
		Description: "e2e-test-feature-description",
		Variations: []*feature.Variation{
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

func newAddVariationsCommand(t *testing.T, vs []*feature.Variation) []*feature.Command {
	var cmds []*feature.Command
	for _, v := range vs {
		cmd, err := util.MarshalCommand(&feature.AddVariationCommand{
			Value:       v.Value,
			Name:        v.Name,
			Description: v.Description,
		})
		if err != nil {
			t.Fatal(err)
		}
		cmds = append(cmds, &feature.Command{Command: cmd})
	}
	return cmds
}

func newRemoveVariationsCommand(t *testing.T, featureIDS []string) []*feature.Command {
	var cmds []*feature.Command
	for _, id := range featureIDS {
		cmd, err := util.MarshalCommand(&feature.RemoveVariationCommand{Id: id})
		if err != nil {
			t.Fatal(err)
		}
		cmds = append(cmds, &feature.Command{Command: cmd})
	}
	return cmds
}

func newVariations(randomValues []string) []*feature.Variation {
	var vs []*feature.Variation
	for _, value := range randomValues {
		v := &feature.Variation{
			Value:       fmt.Sprintf("%s", value),
			Name:        fmt.Sprintf("Variation %s", value),
			Description: fmt.Sprintf("Thing does %s", value),
		}
		vs = append(vs, v)
	}
	return vs
}

func createRandomIDFeatures(t *testing.T, size int, client featureclient.Client) {
	t.Helper()
	for i := 0; i < size; i++ {
		createFeature(t, client, newCreateFeatureCommand(newFeatureID(t)))
	}
}

func createFeatures(t *testing.T, featureIDS []string, client featureclient.Client) {
	t.Helper()
	for _, id := range featureIDS {
		createFeature(t, client, newCreateFeatureCommand(id))
	}
}

func createFeature(t *testing.T, client featureclient.Client, cmd *feature.CreateFeatureCommand) {
	t.Helper()
	createReq := &feature.CreateFeatureRequest{
		Command:              cmd,
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.CreateFeature(ctx, createReq); err != nil {
		t.Fatal(err)
	}
}

func getFeature(t *testing.T, featureID string, client featureclient.Client) *feature.Feature {
	t.Helper()
	getReq := &feature.GetFeatureRequest{
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

/*
func getClonedFeature(t *testing.T, featureID, en string, client featureclient.Client) *feature.Feature {
	t.Helper()
	getReq := &feature.GetFeatureRequest{
		Id:                   featureID,
		EnvironmentNamespace: en,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := client.GetFeature(ctx, getReq)
	if err != nil {
		t.Fatal("Failed to get feature:", err)
	}
	return response.Feature
}
*/

func enableFeatures(t *testing.T, featureIDS []string, client featureclient.Client) {
	t.Helper()
	for _, featureID := range featureIDS {
		enableFeature(t, featureID, client)
	}
}

func enableFeature(t *testing.T, featureID string, client featureclient.Client) {
	t.Helper()
	enableReq := &feature.EnableFeatureRequest{
		Id:                   featureID,
		Command:              &feature.EnableFeatureCommand{},
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.EnableFeature(ctx, enableReq); err != nil {
		t.Fatalf("Failed to enable feature id: %s. Error: %v", featureID, err)
	}
}

func checkEnabledFlag(t *testing.T, features []*feature.Feature) {
	t.Helper()
	for _, feature := range features {
		if !feature.Enabled {
			t.Fatal("Feature enabled flag is false. ID:", feature.Id)
		}
	}
}

func compareVariation(t *testing.T, expected *feature.Variation, actual *feature.Variation) {
	t.Helper()
	if expected.Value != actual.Value {
		t.Fatalf("Different values. Expected: %s actual: %s", expected.Value, actual.Value)
	}
	if expected.Name != actual.Name {
		t.Fatalf("Different names. Expected: %s actual: %s", expected.Name, actual.Name)
	}
	if expected.Description != actual.Description {
		t.Fatalf("Different descriptions. Expected: %s actual: %s", expected.Description, actual.Description)
	}
}

func updateVariations(t *testing.T, featureID string, commands []*feature.Command, client featureclient.Client) {
	t.Helper()
	updateReq := &feature.UpdateFeatureVariationsRequest{
		Id:                   featureID,
		Commands:             commands,
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureVariations(ctx, updateReq); err != nil {
		t.Fatal(err)
	}
}

func findOneOfVariations(ids []string, vs []*feature.Variation) *feature.Variation {
	for _, id := range ids {
		if variation := findVariation(id, vs); variation != nil {
			return variation
		}
	}
	return nil
}

func findVariation(id string, vs []*feature.Variation) *feature.Variation {
	for i := range vs {
		if vs[i].Id == id {
			return vs[i]
		}
	}
	return nil
}

func newFixedStrategyRule(variationID string) *feature.Rule {
	uuid, _ := uuid.NewUUID()
	return &feature.Rule{
		Id: uuid.String(),
		Strategy: &feature.Strategy{
			Type: feature.Strategy_FIXED,
			FixedStrategy: &feature.FixedStrategy{
				Variation: variationID,
			},
		},
		Clauses: []*feature.Clause{
			{
				Attribute: "attribute-1",
				Operator:  feature.Clause_EQUALS,
				Values:    []string{"value-1", "value-2"},
			},
			{
				Attribute: "attribute-2",
				Operator:  feature.Clause_IN,
				Values:    []string{"value-1", "value-2"},
			},
		},
	}
}

func newFixedStrategyRuleWithSegment(variationID, segmentID string) *feature.Rule {
	uuid, _ := uuid.NewUUID()
	return &feature.Rule{
		Id: uuid.String(),
		Strategy: &feature.Strategy{
			Type: feature.Strategy_FIXED,
			FixedStrategy: &feature.FixedStrategy{
				Variation: variationID,
			},
		},
		Clauses: []*feature.Clause{
			{
				Attribute: "attribute-1",
				Operator:  feature.Clause_SEGMENT,
				Values:    []string{segmentID},
			},
		},
	}
}

func newRolloutStrategyRule(variations []*feature.Variation) *feature.Rule {
	uuid, _ := uuid.NewUUID()
	return &feature.Rule{
		Id: uuid.String(),
		Strategy: &feature.Strategy{
			Type: feature.Strategy_ROLLOUT,
			RolloutStrategy: &feature.RolloutStrategy{
				Variations: []*feature.RolloutStrategy_Variation{
					{
						Variation: variations[0].Id,
						Weight:    70000,
					},
					{
						Variation: variations[1].Id,
						Weight:    12000,
					},
					{
						Variation: variations[2].Id,
						Weight:    10000,
					},
					{
						Variation: variations[3].Id,
						Weight:    8000,
					},
				},
			},
		},
		Clauses: []*feature.Clause{
			{
				Attribute: "attribute-1",
				Operator:  feature.Clause_EQUALS,
				Values:    []string{"value-1", "value-2"},
			},
			{
				Attribute: "attribute-2",
				Operator:  feature.Clause_IN,
				Values:    []string{"value-1", "value-2"},
			},
		},
	}
}

func newClause() *feature.Clause {
	return &feature.Clause{
		Attribute: "attribute-3",
		Operator:  feature.Clause_EQUALS,
		Values:    []string{"value-3-a", "value-3-b"},
	}
}

func compareClause(t *testing.T, expected *feature.Clause, actual *feature.Clause) {
	t.Helper()
	if expected.Attribute != actual.Attribute {
		t.Fatalf("Attribute does not match. Expected: %s actual %s", expected.Attribute, actual.Attribute)
	}
	if expected.Operator != actual.Operator {
		t.Fatalf("Operator does not match. Expected: %v actual %v", expected.Operator, actual.Operator)
	}
	if !reflect.DeepEqual(expected.Values, actual.Values) {
		t.Fatalf("Values does not match. Expected: %v actual %v", expected.Values, actual.Values)
	}
}

func updateFeatureTargeting(t *testing.T, client featureclient.Client, cmd *any.Any, featureID string) {
	t.Helper()
	updateReq := &feature.UpdateFeatureTargetingRequest{
		Id: featureID,
		Commands: []*feature.Command{
			{Command: cmd},
		},
		EnvironmentNamespace: *environmentNamespace,
		From:                 feature.UpdateFeatureTargetingRequest_USER,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureTargeting(ctx, updateReq); err != nil {
		t.Fatal(err)
	}
}

func compareFeatures(t *testing.T, expected []*feature.Feature, actual []*feature.Feature) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("Different sizes. Expected: %d, actual: %d", len(expected), len(actual))
	}
	for i := 0; i < len(expected); i++ {
		if !proto.Equal(actual[i], expected[i]) {
			t.Fatalf("Features do not match. Expected: %v, actual: %v", expected[i], actual[i])
		}
	}
}

func evaluateFeatures(t *testing.T, client featureclient.Client, userID, tag string) *feature.EvaluateFeaturesResponse {
	t.Helper()
	req := &feature.EvaluateFeaturesRequest{
		User:                 &userproto.User{Id: userID},
		EnvironmentNamespace: *environmentNamespace,
		Tag:                  tag,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := client.EvaluateFeatures(ctx, req)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return res
}

func createProgressiveRollout(
	ctx context.Context,
	t *testing.T,
	client aoclient.Client,
	featureID string,
	manual *aoproto.ProgressiveRolloutManualScheduleClause,
	template *aoproto.ProgressiveRolloutTemplateScheduleClause,
) {
	t.Helper()
	cmd := &aoproto.CreateProgressiveRolloutCommand{
		FeatureId:                                featureID,
		ProgressiveRolloutManualScheduleClause:   manual,
		ProgressiveRolloutTemplateScheduleClause: template,
	}
	_, err := client.CreateProgressiveRollout(ctx, &aoproto.CreateProgressiveRolloutRequest{
		EnvironmentNamespace: *environmentNamespace,
		Command:              cmd,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func listProgressiveRollouts(t *testing.T, client aoclient.Client, featureID string) []*aoproto.ProgressiveRollout {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.ListProgressiveRollouts(ctx, &aoproto.ListProgressiveRolloutsRequest{
		EnvironmentNamespace: *environmentNamespace,
		PageSize:             0,
		FeatureIds:           []string{featureID},
	})
	if err != nil {
		t.Fatal("Failed to list progressive rollout", err)
	}
	return resp.ProgressiveRollouts
}

func getProgressiveRollout(t *testing.T, client aoclient.Client, id string) *aoproto.ProgressiveRollout {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.GetProgressiveRollout(ctx, &aoproto.GetProgressiveRolloutRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
	})
	if err != nil {
		t.Fatal("Failed to get progressive rollout", err)
	}
	return resp.ProgressiveRollout
}

func newAutoOpsClient(t *testing.T) aoclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := aoclient.NewClient(
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

/*
func newEnvironmentClient(t *testing.T) environmentclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := environmentclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create environment client:", err)
	}
	return client
}

func newEnvironmentCommand(id string) *environment.CreateEnvironmentCommand {
	namespace := strings.Replace(id, "-", "", -1)
	return &environment.CreateEnvironmentCommand{
		Namespace:   namespace,
		Name:        "e2e-test-environment-namespace",
		Description: "e2e-test-environment-namespace-description",
		Id:          id,
		ProjectId:   defaultProjectID,
	}
}

func createEnvironment(t *testing.T, client environmentclient.Client, cmd *environment.CreateEnvironmentCommand) {
	t.Helper()
	createReq := &environment.CreateEnvironmentRequest{
		Command: cmd,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.CreateEnvironment(ctx, createReq); err != nil {
		t.Fatal(err)
	}
}
*/
