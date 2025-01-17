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

package coderef

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/ptypes/wrappers"

	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	coderefproto "github.com/bucketeer-io/bucketeer/proto/coderef"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	timeout = 30 * time.Second
)

var (
	webGatewayAddr       = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort       = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert       = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath           = flag.String("api-key", "", "Client SDK API key for api-gateway")
	apiKeyServerPath     = flag.String("api-key-server", "", "Server SDK API key for api-gateway")
	gatewayAddr          = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort          = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert          = flag.String("gateway-cert", "", "Gateway crt file")
	serviceTokenPath     = flag.String("service-token", "", "Service token path")
	environmentNamespace = flag.String("environment-namespace", "", "Environment namespace")
	testID               = flag.String("test-id", "", "test ID")
	organizationID       = flag.String("organization-id", "", "Organization ID")
)

func TestCreateCodeReference(t *testing.T) {
	t.Parallel()
	client := newCodeRefClient(t)
	featureClient := newFeatureClient(t)

	// First create a feature
	featureID := createFeatureID(t)
	createFeature(t, featureClient, featureID)

	// Create code reference
	createReq := newCreateCodeReferenceRequest(featureID)
	resp := createCodeReference(t, client, createReq)

	assert.Equal(t, createReq.FeatureId, resp.CodeReference.FeatureId)
	assert.Equal(t, createReq.FilePath, resp.CodeReference.FilePath)
	assert.Equal(t, createReq.LineNumber, resp.CodeReference.LineNumber)
	assert.Equal(t, createReq.CodeSnippet, resp.CodeReference.CodeSnippet)
	assert.Equal(t, createReq.RepositoryName, resp.CodeReference.RepositoryName)
	assert.Equal(t, createReq.RepositoryOwner, resp.CodeReference.RepositoryOwner)
	assert.Equal(t, createReq.RepositoryType, resp.CodeReference.RepositoryType)
	assert.Equal(t, createReq.RepositoryBranch, resp.CodeReference.RepositoryBranch)
	assert.Equal(t, createReq.CommitHash, resp.CodeReference.CommitHash)
	assert.NotEmpty(t, resp.CodeReference.Id)
}

func TestUpdateCodeReference(t *testing.T) {
	t.Parallel()
	client := newCodeRefClient(t)
	featureClient := newFeatureClient(t)

	// First create a feature
	featureID := createFeatureID(t)
	createFeature(t, featureClient, featureID)

	// Create code reference
	createReq := newCreateCodeReferenceRequest(featureID)
	createResp := createCodeReference(t, client, createReq)

	// Update code reference
	updateReq := &coderefproto.UpdateCodeReferenceRequest{
		Id:               createResp.CodeReference.Id,
		EnvironmentId:    *environmentNamespace,
		FilePath:         "updated/path/to/file.go",
		LineNumber:       200,
		CodeSnippet:      "updated code snippet",
		ContentHash:      "updated-hash-123",
		RepositoryName:   createResp.CodeReference.RepositoryName,
		RepositoryOwner:  createResp.CodeReference.RepositoryOwner,
		RepositoryType:   createResp.CodeReference.RepositoryType,
		RepositoryBranch: createResp.CodeReference.RepositoryBranch,
		CommitHash:       createResp.CodeReference.CommitHash,
	}
	_, err := client.UpdateCodeReference(context.Background(), updateReq)
	assert.NoError(t, err)

	// Get and verify update
	getResp := getCodeReference(t, client, &coderefproto.GetCodeReferenceRequest{
		Id:            createResp.CodeReference.Id,
		EnvironmentId: *environmentNamespace,
	})
	assert.Equal(t, updateReq.FilePath, getResp.CodeReference.FilePath)
	assert.Equal(t, updateReq.LineNumber, getResp.CodeReference.LineNumber)
	assert.Equal(t, updateReq.CodeSnippet, getResp.CodeReference.CodeSnippet)
	assert.Equal(t, updateReq.ContentHash, getResp.CodeReference.ContentHash)
	assert.Equal(t, updateReq.RepositoryName, getResp.CodeReference.RepositoryName)
	assert.Equal(t, updateReq.RepositoryOwner, getResp.CodeReference.RepositoryOwner)
	assert.Equal(t, updateReq.RepositoryType, getResp.CodeReference.RepositoryType)
	assert.Equal(t, updateReq.RepositoryBranch, getResp.CodeReference.RepositoryBranch)
	assert.Equal(t, updateReq.CommitHash, getResp.CodeReference.CommitHash)
}

func TestListCodeReferences(t *testing.T) {
	t.Parallel()
	client := newCodeRefClient(t)
	featureClient := newFeatureClient(t)

	// First create a feature
	featureID := createFeatureID(t)
	createFeature(t, featureClient, featureID)

	// Create multiple code references
	createReq1 := newCreateCodeReferenceRequest(featureID)
	createCodeReference(t, client, createReq1)
	time.Sleep(time.Second)
	createReq2 := newCreateCodeReferenceRequest(featureID)
	createCodeReference(t, client, createReq2)

	// List code references
	listReq := &coderefproto.ListCodeReferencesRequest{
		FeatureId:      featureID,
		EnvironmentId:  *environmentNamespace,
		PageSize:       10,
		Cursor:         "0",
		OrderBy:        coderefproto.ListCodeReferencesRequest_CREATED_AT,
		OrderDirection: coderefproto.ListCodeReferencesRequest_DESC,
	}
	listResp, err := client.ListCodeReferences(context.Background(), listReq)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(listResp.CodeReferences))
}

func TestDeleteCodeReference(t *testing.T) {
	t.Parallel()
	client := newCodeRefClient(t)
	featureClient := newFeatureClient(t)

	// First create a feature
	featureID := createFeatureID(t)
	createFeature(t, featureClient, featureID)

	// Create code reference
	createReq := newCreateCodeReferenceRequest(featureID)
	createResp := createCodeReference(t, client, createReq)

	// Delete code reference
	deleteReq := &coderefproto.DeleteCodeReferenceRequest{
		Id:            createResp.CodeReference.Id,
		EnvironmentId: *environmentNamespace,
	}
	_, err := client.DeleteCodeReference(context.Background(), deleteReq)
	assert.NoError(t, err)

	// Verify deletion
	getReq := &coderefproto.GetCodeReferenceRequest{
		Id:            createResp.CodeReference.Id,
		EnvironmentId: *environmentNamespace,
	}
	_, err = client.GetCodeReference(context.Background(), getReq)
	assert.Error(t, err) // Should return error as code reference is deleted
}

func newCreateCodeReferenceRequest(featureID string) *coderefproto.CreateCodeReferenceRequest {
	return &coderefproto.CreateCodeReferenceRequest{
		EnvironmentId:    *environmentNamespace,
		FeatureId:        featureID,
		FilePath:         "path/to/file.go",
		LineNumber:       100,
		CodeSnippet:      "if (feature.enabled) { doSomething() }",
		ContentHash:      "abc123",
		Aliases:          []string{"test-alias"},
		RepositoryName:   "test-repo",
		RepositoryOwner:  "test-owner",
		RepositoryType:   coderefproto.CodeReference_GITHUB,
		RepositoryBranch: "main",
		CommitHash:       "hash123",
	}
}

func createCodeReference(
	t *testing.T,
	client coderefproto.CodeReferenceServiceClient,
	req *coderefproto.CreateCodeReferenceRequest,
) *coderefproto.CreateCodeReferenceResponse {
	resp, err := client.CreateCodeReference(context.Background(), req)
	assert.NoError(t, err)
	return resp
}

func getCodeReference(
	t *testing.T,
	client coderefproto.CodeReferenceServiceClient,
	req *coderefproto.GetCodeReferenceRequest,
) *coderefproto.GetCodeReferenceResponse {
	resp, err := client.GetCodeReference(context.Background(), req)
	assert.NoError(t, err)
	return resp
}

func newCodeRefClient(t *testing.T) coderefproto.CodeReferenceServiceClient {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	assert.NoError(t, err)
	conn, err := rpcclient.NewClientConn(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(timeout),
		rpcclient.WithBlock(),
	)
	assert.NoError(t, err)
	return coderefproto.NewCodeReferenceServiceClient(conn)
}

func newFeatureClient(t *testing.T) featureclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	assert.NoError(t, err)
	featureClient, err := featureclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(timeout),
		rpcclient.WithBlock(),
	)
	assert.NoError(t, err)
	return featureClient
}

func createFeatureID(t *testing.T) string {
	return fmt.Sprintf("feature-id-%d", time.Now().UnixNano())
}

func createFeature(t *testing.T, client featureclient.Client, featureID string) {
	cmd := &featureproto.CreateFeatureCommand{
		Id:          featureID,
		Name:        "e2e-test-feature",
		Description: "e2e test feature",
		Tags:        []string{"e2e-test"},
		Variations: []*featureproto.Variation{
			{
				Value:       "true",
				Name:        "true",
				Description: "this is a true variation",
			},
			{
				Value:       "false",
				Name:        "false",
				Description: "this is a false variation",
			},
		},
		DefaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
		DefaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
	}
	createReq := &featureproto.CreateFeatureRequest{
		Command:       cmd,
		EnvironmentId: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.CreateFeature(ctx, createReq)
	assert.NoError(t, err)
}
