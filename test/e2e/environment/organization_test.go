package environment

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	notificationproto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
	pushproto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

const (
	defaultOrganizationID = "e2e"
)

var (
	fcmServiceAccountDummy = `{
		"type": "service_account",
		"project_id": "e2e-%d",
		"private_key_id": "private-key-id",
		"private_key": "-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----\n",
		"client_email": "fcm-service-account@test.iam.gserviceaccount.com",
		"client_id": "client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/fcm-service-account@test.iam.gserviceaccount.com",
		"universe_domain": "googleapis.com"
	}`
)

func TestCreateDeleteOrganization(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	envc := newEnvironmentClient(t)
	ftc := newFeatureClient(t)
	expc := newExperimentClient(t)
	pushc := newPushClient(t)
	notic := newNotificationClient(t)
	opsc := newAutoOpsClient(t)
	accountc := newAccountClient(t)
	defer envc.Close()
	defer ftc.Close()
	defer expc.Close()
	defer pushc.Close()
	defer notic.Close()
	defer opsc.Close()
	defer accountc.Close()

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
	// 3.1 create feature
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

	// 3.2 create goals
	createGoalResp, err := expc.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
		EnvironmentId:  envID,
		Id:             fmt.Sprintf("goal-id-e2e-%d", time.Now().UnixNano()),
		Name:           fmt.Sprintf("goal-name-e2e-%d", time.Now().UnixNano()),
		ConnectionType: experimentproto.Goal_EXPERIMENT,
	})
	if err != nil {
		t.Fatal(err)
	}
	if createGoalResp == nil || createGoalResp.Goal == nil {
		t.Fatal("create goal response or goal is nil")
	}

	// 3.3 create push
	createPush, err := pushc.CreatePush(ctx, &pushproto.CreatePushRequest{
		EnvironmentId:     envID,
		Name:              fmt.Sprintf("push-name-e2e-%d", time.Now().UnixNano()),
		FcmServiceAccount: []byte(fmt.Sprintf(fcmServiceAccountDummy, time.Now().UnixNano())),
		Tags:              []string{"e2e-test-tag-1", "e2e-test-tag-2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if createPush == nil || createPush.Push == nil {
		t.Fatal("create push response or push is nil")
	}

	// 3.4 create subscription
	createSubscription, err := notic.CreateSubscription(ctx, &notificationproto.CreateSubscriptionRequest{
		EnvironmentId: envID,
		Name:          fmt.Sprintf("subscription-name-e2e-%d", time.Now().UnixNano()),
		SourceTypes: []notificationproto.Subscription_SourceType{
			notificationproto.Subscription_DOMAIN_EVENT_FEATURE,
			notificationproto.Subscription_DOMAIN_EVENT_APIKEY,
		},
		Recipient: &notificationproto.Recipient{
			Type: notificationproto.Recipient_SlackChannel,
			SlackChannelRecipient: &notificationproto.SlackChannelRecipient{
				WebhookUrl: "https://hooks.slack",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if createSubscription == nil || createSubscription.Subscription == nil {
		t.Fatal("create subscription response or subscription is nil")
	}

	// 3.5 create operation
	createOpsResp, err := opsc.CreateAutoOpsRule(ctx, &autoopsproto.CreateAutoOpsRuleRequest{
		EnvironmentId: envID,
		FeatureId:     createFfResp.Feature.Id,
		OpsType:       autoopsproto.OpsType_SCHEDULE,
		DatetimeClauses: []*autoopsproto.DatetimeClause{
			{
				Time:       time.Now().Add(2 * time.Hour).Unix(),
				ActionType: autoopsproto.ActionType_ENABLE,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if createOpsResp == nil || createOpsResp.AutoOpsRule == nil {
		t.Fatal("create auto ops rule response or rule is nil")
	}

	// 3.6 create experiment
	createExpResp, err := expc.CreateExperiment(ctx, &experimentproto.CreateExperimentRequest{
		EnvironmentId:   envID,
		FeatureId:       createFfResp.Feature.Id,
		Name:            fmt.Sprintf("experiment-name-e2e-%d", time.Now().UnixNano()),
		Description:     "e2e test experiment description",
		GoalIds:         []string{createGoalResp.Goal.Id},
		StartAt:         time.Now().Unix(),
		StopAt:          time.Now().Add(24 * time.Hour).Unix(),
		BaseVariationId: createFfResp.Feature.Variations[0].Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if createExpResp == nil || createExpResp.Experiment == nil {
		t.Fatal("create experiment response or experiment is nil")
	}

	// 3.7 create account
	createAccResp, err := accountc.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
		OrganizationId:   orgID,
		Email:            fmt.Sprintf("e2e%d@bucketeer.io", time.Now().UnixNano()),
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{
				EnvironmentId: envID,
				Role:          accountproto.AccountV2_Role_Environment_EDITOR,
			},
		},
		FirstName: fmt.Sprintf("first-name-e2e-%d", time.Now().UnixNano()),
		LastName:  fmt.Sprintf("last-name-e2e-%d", time.Now().UnixNano()),
	})
	if err != nil {
		t.Fatal(err)
	}
	if createAccResp == nil || createAccResp.Account == nil {
		t.Fatal("create account response or account is nil")
	}

	// 3.8 create API Key
	createAPIKeyResp, err := accountc.CreateAPIKey(ctx, &accountproto.CreateAPIKeyRequest{
		EnvironmentId: envID,
		Name:          fmt.Sprintf("api-key-name-e2e-%d", time.Now().UnixNano()),
		Role:          accountproto.APIKey_PUBLIC_API_WRITE,
	})
	if err != nil {
		t.Fatal(err)
	}
	if createAPIKeyResp == nil || createAPIKeyResp.ApiKey == nil {
		t.Fatal("create API key response or API key is nil")
	}

	// 4.0 dry-run delete organization data
	_, err = envc.DeleteOrganizationData(ctx, &environmentproto.DeleteOrganizationDataRequest{
		OrganizationIds: []string{orgID},
		DryRun:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	// 4.1 verify feature still exists
	getFeatureResp, err := ftc.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentId: envID,
		Id:            createFfResp.Feature.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getFeatureResp == nil || getFeatureResp.Feature == nil {
		t.Fatal("get feature response or feature is nil after dry-run delete")
	}

	// 4.2 verify goal still exists
	getGoalResp, err := expc.GetGoal(ctx, &experimentproto.GetGoalRequest{
		EnvironmentId: envID,
		Id:            createGoalResp.Goal.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getGoalResp == nil || getGoalResp.Goal == nil {
		t.Fatal("get goal response or goal is nil after dry-run delete")
	}

	// 4.3 verify push still exists
	getPushResp, err := pushc.GetPush(ctx, &pushproto.GetPushRequest{
		EnvironmentId: envID,
		Id:            createPush.Push.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getPushResp == nil || getPushResp.Push == nil {
		t.Fatal("get push response or push is nil after dry-run delete")
	}

	// 4.4 verify subscription still exists
	getSubscriptionResp, err := notic.GetSubscription(ctx, &notificationproto.GetSubscriptionRequest{
		EnvironmentId: envID,
		Id:            createSubscription.Subscription.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getSubscriptionResp == nil || getSubscriptionResp.Subscription == nil {
		t.Fatal("get subscription response or subscription is nil after dry-run delete")
	}

	// 4.5 verify auto ops rule still exists
	getAutoOpsRuleResp, err := opsc.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentId: envID,
		Id:            createOpsResp.AutoOpsRule.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getAutoOpsRuleResp == nil || getAutoOpsRuleResp.AutoOpsRule == nil {
		t.Fatal("get auto ops rule response or rule is nil after dry-run delete")
	}

	// 4.6 verify experiment still exists
	getExperimentResp, err := expc.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		EnvironmentId: envID,
		Id:            createExpResp.Experiment.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getExperimentResp == nil || getExperimentResp.Experiment == nil {
		t.Fatal("get experiment response or experiment is nil after dry-run delete")
	}

	// 4.7 verify account still exists
	getAccountResp, err := accountc.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		OrganizationId: orgID,
		Email:          createAccResp.Account.Email,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getAccountResp == nil || getAccountResp.Account == nil {
		t.Fatal("get account response or account is nil after dry-run delete")
	}

	// 4.8 verify API key still exists
	getAPIKeyResp, err := accountc.GetAPIKey(ctx, &accountproto.GetAPIKeyRequest{
		EnvironmentId: envID,
		Id:            createAPIKeyResp.ApiKey.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getAPIKeyResp == nil || getAPIKeyResp.ApiKey == nil {
		t.Fatal("get API key response or API key is nil after dry-run delete")
	}

	// 5. delete organization data
	_, err = envc.DeleteOrganizationData(ctx, &environmentproto.DeleteOrganizationDataRequest{
		OrganizationIds: []string{orgID},
	})
	if err != nil {
		t.Fatal(err)
	}

	// 6.0 verify data is deleted
	// 6.1 verify feature is deleted
	_, err = ftc.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentId: envID,
		Id:            createFfResp.Feature.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting feature from deleted organization, but got nil")
	}

	// 6.2 verify goal is deleted
	_, err = expc.GetGoal(ctx, &experimentproto.GetGoalRequest{
		EnvironmentId: envID,
		Id:            createGoalResp.Goal.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting goal from deleted organization, but got nil")
	}

	// 6.3 verify push is deleted
	_, err = pushc.GetPush(ctx, &pushproto.GetPushRequest{
		EnvironmentId: envID,
		Id:            createPush.Push.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting push from deleted organization, but got nil")
	}

	// 6.4 verify subscription is deleted
	_, err = notic.GetSubscription(ctx, &notificationproto.GetSubscriptionRequest{
		EnvironmentId: envID,
		Id:            createSubscription.Subscription.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting subscription from deleted organization, but got nil")
	}

	// 6.5 verify auto ops rule is deleted
	_, err = opsc.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentId: envID,
		Id:            createOpsResp.AutoOpsRule.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting auto ops rule from deleted organization, but got nil")
	}

	// 6.6 verify experiment is deleted
	_, err = expc.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		EnvironmentId: envID,
		Id:            createExpResp.Experiment.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting experiment from deleted organization, but got nil")
	}

	// 6.7 verify account is deleted
	_, err = accountc.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		OrganizationId: orgID,
		Email:          createAccResp.Account.Email,
	})
	if err == nil {
		t.Fatal("expected error when getting account from deleted organization, but got nil")
	}

	// 6.8 verify API key is deleted
	_, err = accountc.GetAPIKey(ctx, &accountproto.GetAPIKeyRequest{
		EnvironmentId: envID,
		Id:            createAPIKeyResp.ApiKey.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting API key from deleted organization, but got nil")
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
