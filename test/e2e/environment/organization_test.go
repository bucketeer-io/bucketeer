package environment

import (
	"context"
	"fmt"
	"testing"
	"time"

	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

const (
	defaultOrganizationID = "default"
)

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
	newDesc := fmt.Sprintf("Description %v", time.Now().Unix())
	newName := fmt.Sprintf("name-%v", time.Now().Unix())
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
