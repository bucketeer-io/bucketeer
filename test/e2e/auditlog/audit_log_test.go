package auditlog

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"

	auditlogclient "github.com/bucketeer-io/bucketeer/pkg/auditlog/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	"github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	timeout = 60 * time.Second
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

func TestListAndGetAuditLog(t *testing.T) {
	featureClient := newFeatureClient(t)
	req := newCreateFeatureReq(newFeatureID(t))
	createFeatureNoCmd(t, featureClient, req)
	// wait for the audit log to save
	time.Sleep(5 * time.Second)

	auditlogClient := newAuditLogClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	listResp, err := auditlogClient.ListAuditLogs(ctx, &auditlog.ListAuditLogsRequest{
		EnvironmentId:  *environmentID,
		EntityType:     wrapperspb.Int32(int32(eventproto.Event_FEATURE)),
		PageSize:       10,
		Cursor:         "0",
		OrderBy:        auditlog.ListAuditLogsRequest_TIMESTAMP,
		OrderDirection: auditlog.ListAuditLogsRequest_DESC,
	})
	if err != nil {
		t.Fatal("Failed to list audit logs:", err)
	}

	_, err = auditlogClient.GetAuditLog(ctx, &auditlog.GetAuditLogRequest{
		Id:            listResp.AuditLogs[0].Id,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal("Failed to get audit log:", err)
	}
}

func newFeatureID(t *testing.T) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-feature-id-%s", "audit-log", *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-feature-id-%s", "audit-log", newUUID(t))
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func createFeatureNoCmd(t *testing.T, client featureclient.Client, req *feature.CreateFeatureRequest) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.CreateFeature(ctx, req); err != nil {
		t.Fatal(err)
	}
}

func newCreateFeatureReq(featureID string) *feature.CreateFeatureRequest {
	return &feature.CreateFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
		Name:          "e2e-test-feature-name",
		Description:   "e2e-test-feature-description",
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

func newAuditLogClient(t *testing.T) auditlogclient.Client {
	t.Helper()
	creds, err := client.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	auditlogClient, err := auditlogclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create auditlog client:", err)
	}
	return auditlogClient
}
