package environment

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"

	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	defaultOrganizationID = "e2e"
)

func TestCreateDeleteOrganization(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	envc := newEnvironmentClient(t)
	ftc := newFeatureClient(t)
	defer envc.Close()

	// 1. Create Organization
	createOrgResp, err := envc.CreateOrganization(ctx, &environmentproto.CreateOrganizationRequest{
		Name:          fmt.Sprintf("org-e2e-%d", time.Now().UnixNano()),
		UrlCode:       fmt.Sprintf("org-url-%d", time.Now().UnixNano()),
		IsSystemAdmin: false,
		OwnerEmail:    "demo@bucketeer.io",
	})
	if err != nil {
		t.Fatal(err)
	}
	if createOrgResp == nil || createOrgResp.Organization == nil {
		t.Fatal("create organization response or organization is nil")
	}
	orgID := createOrgResp.Organization.Id

	// 2. get environment
	getEnvResp, err := envc.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{
		OrganizationId: orgID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getEnvResp == nil || len(getEnvResp.Environments) == 0 {
		t.Fatalf("no environments found for organization %s", orgID)
	}
	envID := getEnvResp.Environments[0].Id

	// 3. create data in organization
	createFfResp, err := ftc.CreateFeature(ctx, newCreateFeatureReq(
		fmt.Sprintf("feature-e2e-%d", time.Now().UnixNano()),
		envID,
	))
	if err != nil {
		t.Fatal(err)
	}
	if createFfResp == nil || createFfResp.Feature == nil {
		t.Fatal("create feature response or feature is nil")
	}

	// 4. delete organization data
	_, err = envc.DeleteOrganizationData(ctx, &environmentproto.DeleteOrganizationDataRequest{
		OrganizationIds: []string{orgID},
	})
	if err != nil {
		t.Fatal(err)
	}

	// 5. verify data is deleted
	_, err = ftc.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentId: envID,
		Id:            createFfResp.Feature.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting feature from deleted organization, but got nil")
	}
}

func TestGetOrganization(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	id := defaultOrganizationID
	resp, err := c.GetOrganization(ctx, &environmentproto.GetOrganizationRequest{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Organization.Id != id {
		t.Fatalf("different ids, expected: %v, actual: %v", id, resp.Organization.Id)
	}
}

func TestListOrganizations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	pageSize := int64(1)
	resp, err := c.ListOrganizations(ctx, &environmentproto.ListOrganizationsRequest{PageSize: pageSize})
	if err != nil {
		t.Fatal(err)
	}
	responseSize := int64(len(resp.Organizations))
	if responseSize != pageSize {
		t.Fatalf("different sizes, expected: %d actual: %d", pageSize, responseSize)
	}
}

func TestUpdateOrganization(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	id := defaultOrganizationID
	newDesc := fmt.Sprintf("This organization is for organization e2e tests (Updated at %d)", time.Now().Unix())
	newName := fmt.Sprintf("E2E organization (Updated at %d)", time.Now().Unix())
	_, err := c.UpdateOrganization(ctx, &environmentproto.UpdateOrganizationRequest{
		Id:                       id,
		ChangeDescriptionCommand: &environmentproto.ChangeDescriptionOrganizationCommand{Description: newDesc},
		RenameCommand:            &environmentproto.ChangeNameOrganizationCommand{Name: newName},
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetOrganization(ctx, &environmentproto.GetOrganizationRequest{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if getResp.Organization.Id != id {
		t.Fatalf("different ids, expected: %v, actual: %v", id, getResp.Organization.Id)
	}
	if getResp.Organization.Description != newDesc {
		t.Fatalf("different descriptions, expected: %v, actual: %v", newDesc, getResp.Organization.Description)
	}
	if getResp.Organization.Name != newName {
		t.Fatalf("different names, expected: %v, actual: %v", newName, getResp.Organization.Name)
	}
}

func TestEnableAndDisableOrganization(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	id := defaultOrganizationID
	_, err := c.DisableOrganization(ctx, &environmentproto.DisableOrganizationRequest{
		Id:      id,
		Command: &environmentproto.DisableOrganizationCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp1, err := c.GetOrganization(ctx, &environmentproto.GetOrganizationRequest{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if getResp1.Organization.Disabled != true {
		t.Fatalf("different ids, expected: %v, actual: %v", id, getResp1.Organization.Id)
	}

	_, err = c.EnableOrganization(ctx, &environmentproto.EnableOrganizationRequest{
		Id:      id,
		Command: &environmentproto.EnableOrganizationCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp2, err := c.GetOrganization(ctx, &environmentproto.GetOrganizationRequest{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if getResp2.Organization.Disabled != false {
		t.Fatalf("different ids, expected: %v, actual: %v", id, getResp2.Organization.Id)
	}
}

func TestArchiveAndUnarchiveOrganization(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	id := defaultOrganizationID
	_, err := c.ArchiveOrganization(ctx, &environmentproto.ArchiveOrganizationRequest{
		Id:      id,
		Command: &environmentproto.ArchiveOrganizationCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp1, err := c.GetOrganization(ctx, &environmentproto.GetOrganizationRequest{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if getResp1.Organization.Archived != true {
		t.Fatalf("different ids, expected: %v, actual: %v", id, getResp1.Organization.Id)
	}

	_, err = c.UnarchiveOrganization(ctx, &environmentproto.UnarchiveOrganizationRequest{
		Id:      id,
		Command: &environmentproto.UnarchiveOrganizationCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp2, err := c.GetOrganization(ctx, &environmentproto.GetOrganizationRequest{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if getResp2.Organization.Archived != false {
		t.Fatalf("different ids, expected: %v, actual: %v", id, getResp2.Organization.Id)
	}
}

func newCreateFeatureReq(featureID, envID string) *featureproto.CreateFeatureRequest {
	return &featureproto.CreateFeatureRequest{
		Id:            featureID,
		EnvironmentId: envID,
		Name:          "e2e-test-feature-name",
		Description:   "e2e-test-feature-description",
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
