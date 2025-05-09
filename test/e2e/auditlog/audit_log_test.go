package auditlog

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	auditlogclient "github.com/bucketeer-io/bucketeer/pkg/auditlog/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	"github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	timeout                  = 60 * time.Second
	sleepTimeBetweenRequests = 10 * time.Second
	maxRetries               = 15
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
	t.Parallel()
	featureClient := newFeatureClient(t)
	req := newCreateFeatureReq(newFeatureID(t))
	createFeatureNoCmd(t, featureClient, req)
	// wait for the audit log to be saved
	time.Sleep(20 * time.Second)

	auditlogClient := newAuditLogClient(t)

	var auditLogID string
	cursor := "0"
	maxLogs := 1000 // Maximum number of logs to fetch
	totalLogs := 0
	for {
		listResp, err := listAuditLogsWithRetry(t, auditlogClient, &auditlog.ListAuditLogsRequest{
			EnvironmentId:  *environmentID,
			EntityType:     wrapperspb.Int32(int32(eventproto.Event_FEATURE)),
			PageSize:       100,
			Cursor:         cursor,
			OrderBy:        auditlog.ListAuditLogsRequest_TIMESTAMP,
			OrderDirection: auditlog.ListAuditLogsRequest_DESC,
		})
		if err != nil {
			t.Fatal("Failed to list audit logs:", err)
		}

		totalLogs += len(listResp.AuditLogs)
		// Search for the target audit log in the current batch
		for _, log := range listResp.AuditLogs {
			if log.EntityId == req.Id {
				auditLogID = log.Id
				break
			}
		}

		// Break if we found the target log, reached max logs, or reached the end
		if auditLogID != "" || listResp.Cursor == "" || len(listResp.AuditLogs) == 0 || totalLogs >= maxLogs {
			break
		}
		cursor = listResp.Cursor
	}

	if auditLogID == "" {
		t.Fatal("Failed to find audit log for the created feature")
	}

	getResp, err := getAuditLogWithRetry(t, auditlogClient, &auditlog.GetAuditLogRequest{
		Id:            auditLogID,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal("Failed to get audit log:", err)
	}
	if getResp.AuditLog.Id != auditLogID {
		t.Fatal("GetAuditLog ID error")
	}

	listAdminResp, err := listAdminAuditLogsWithRetry(t, auditlogClient, &auditlog.ListAdminAuditLogsRequest{
		EntityType:     wrapperspb.Int32(int32(eventproto.Event_FEATURE)),
		PageSize:       10,
		Cursor:         "0",
		OrderBy:        auditlog.ListAdminAuditLogsRequest_TIMESTAMP,
		OrderDirection: auditlog.ListAdminAuditLogsRequest_DESC,
	})
	if err != nil {
		t.Fatal("Failed to list admin audit logs", err)
	}
	if len(listAdminResp.AuditLogs) > 10 {
		t.Fatal("ListAdminAuditLogs page size error")
	}
	for i := 1; i < len(listAdminResp.AuditLogs); i++ {
		if listAdminResp.AuditLogs[i].Timestamp > listAdminResp.AuditLogs[i-1].Timestamp {
			t.Fatal("ListAdminAuditLogs order error")
		}
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

func listAuditLogsWithRetry(
	t *testing.T,
	client auditlogclient.Client,
	req *auditlog.ListAuditLogsRequest,
) (*auditlog.ListAuditLogsResponse, error) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := range maxRetries {
		resp, err := client.ListAuditLogs(ctx, req)
		if err == nil {
			return resp, nil
		}
		if i == maxRetries-1 {
			return nil, fmt.Errorf("Failed to list audit logs after %d retries: %w", maxRetries, err)
		}
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable || st.Code() == codes.Internal || st.Code() == codes.DeadlineExceeded {
			fmt.Printf("Failed to list audit logs. Error code: %d. Retrying in %d seconds.\n", st.Code(), sleepTimeBetweenRequests)
			time.Sleep(sleepTimeBetweenRequests)
			continue
		}
		return nil, err
	}
	return nil, fmt.Errorf("Unexpected error: max retries reached")
}

func getAuditLogWithRetry(
	t *testing.T,
	client auditlogclient.Client,
	req *auditlog.GetAuditLogRequest,
) (*auditlog.GetAuditLogResponse, error) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := range maxRetries {
		resp, err := client.GetAuditLog(ctx, req)
		if err == nil {
			return resp, nil
		}
		if i == maxRetries-1 {
			return nil, fmt.Errorf("Failed to get audit log after %d retries: %w", maxRetries, err)
		}
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable || st.Code() == codes.Internal || st.Code() == codes.DeadlineExceeded {
			fmt.Printf("Failed to get audit log. Error code: %d. Retrying in %d seconds.\n", st.Code(), sleepTimeBetweenRequests)
			time.Sleep(sleepTimeBetweenRequests)
			continue
		}
		return nil, err
	}
	return nil, fmt.Errorf("Unexpected error: max retries reached")
}

func listAdminAuditLogsWithRetry(
	t *testing.T,
	client auditlogclient.Client,
	req *auditlog.ListAdminAuditLogsRequest,
) (*auditlog.ListAdminAuditLogsResponse, error) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := range maxRetries {
		resp, err := client.ListAdminAuditLogs(ctx, req)
		if err == nil {
			return resp, nil
		}
		if i == maxRetries-1 {
			return nil, fmt.Errorf("Failed to list admin audit logs after %d retries: %w", maxRetries, err)
		}
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable || st.Code() == codes.Internal || st.Code() == codes.DeadlineExceeded {
			fmt.Printf("Failed to list admin audit logs. Error code: %d. Retrying in %d seconds.\n", st.Code(), sleepTimeBetweenRequests)
			time.Sleep(sleepTimeBetweenRequests)
			continue
		}
		return nil, err
	}
	return nil, fmt.Errorf("Unexpected error: max retries reached")
}
