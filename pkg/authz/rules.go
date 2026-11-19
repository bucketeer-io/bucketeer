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

package authz

// The services of the web gateway, by proto service name.
const (
	AccountService       = "bucketeer.account.AccountService"
	AIChatService        = "bucketeer.aichat.AIChatService"
	AuditLogService      = "bucketeer.auditlog.AuditLogService"
	AuthService          = "bucketeer.auth.AuthService"
	AutoOpsService       = "bucketeer.autoops.AutoOpsService"
	CodeReferenceService = "bucketeer.coderef.CodeReferenceService"
	EnvironmentService   = "bucketeer.environment.EnvironmentService"
	EventCounterService  = "bucketeer.eventcounter.EventCounterService"
	ExperimentService    = "bucketeer.experiment.ExperimentService"
	FeatureService       = "bucketeer.feature.FeatureService"
	InsightsService      = "bucketeer.insights.InsightsService"
	NotificationService  = "bucketeer.notification.NotificationService"
	PushService          = "bucketeer.push.PushService"
	SubscriptionService  = "bucketeer.subscription.SubscriptionService"
	TagService           = "bucketeer.tag.TagService"
	TeamService          = "bucketeer.team.TeamService"
)

var WebServices = []string{
	AccountService,
	AIChatService,
	AuditLogService,
	AuthService,
	AutoOpsService,
	CodeReferenceService,
	EnvironmentService,
	EventCounterService,
	ExperimentService,
	FeatureService,
	InsightsService,
	NotificationService,
	PushService,
	SubscriptionService,
	TagService,
	TeamService,
}

// WebRules is Table 2 for the web gateway: what each handler checks today, written as a requirement.
var WebRules = merge(
	rows(AccountService, map[string]Requirement{
		"GetMe":                     Authenticated(),
		"GetMyOrganizations":        Authenticated(),
		"GetMyOrganizationsByEmail": SystemAdminOnly(),
		"CreateAccountV2":           Need(AccountManage, InOrg("organization_id")),
		"EnableAccountV2":           Need(AccountManage, InOrg("organization_id")),
		"DisableAccountV2":          Need(AccountManage, InOrg("organization_id")),
		// A member may update its own account; the owner protection stays in the handler.
		"UpdateAccountV2":             Need(AccountManage, InOrg("organization_id")).Or(Self("email")),
		"DeleteAccountV2":             Need(AccountManage, InOrg("organization_id")),
		"GetAccountV2":                Need(AccountRead, InOrg("organization_id")),
		"GetAccountV2ByEnvironmentID": Need(AccountRead, InOrgOfEnv("environment_id")),
		"ListAccountsV2":              Need(AccountRead, InOrg("organization_id")),
		"CreateAPIKey":                Need(APIKeyManage, InOrgOfEnv("environment_id")),
		"UpdateAPIKey":                Need(APIKeyManage, InOrgOfEnv("environment_id")),
		"GetAPIKey":                   Need(APIKeyRead, InEnv("environment_id")),
		"ListAPIKeys":                 Need(APIKeyRead, InOrg("organization_id")),
		"GetEnvironmentAPIKey":        SystemAdminOnly(),
		"CreateSearchFilter":          Need(SearchFilterWrite, InEnv("environment_id")),
		"UpdateSearchFilter":          Need(SearchFilterWrite, InEnv("environment_id")),
		"DeleteSearchFilter":          Need(SearchFilterWrite, InEnv("environment_id")),
	}),
	rows(AIChatService, map[string]Requirement{
		"Chat":           Need(AIChatRead, InEnv("environment_id")),
		"GetSuggestions": Need(AIChatRead, InEnv("environment_id")),
	}),
	rows(AuditLogService, map[string]Requirement{
		"GetAuditLog":        auditLogScope(),
		"ListAuditLogs":      auditLogScope(),
		"ListFeatureHistory": Need(AuditLogRead, InEnv("environment_id")),
	}),
	rows(AuthService, map[string]Requirement{
		"ExchangeToken":        Public(),
		"GetAuthenticationURL": Public(),
		"RefreshToken":         Public(),
		"SignIn":               Public(),
		"SwitchOrganization":   Public(),
		"GetDemoSiteStatus":    Public(),
	}),
	rows(AutoOpsService, map[string]Requirement{
		"GetAutoOpsRule":            Need(AutoOpsRuleRead, InEnv("environment_id")),
		"ListAutoOpsRules":          Need(AutoOpsRuleRead, InEnv("environment_id")),
		"CreateAutoOpsRule":         Need(AutoOpsRuleWrite, InEnv("environment_id")),
		"StopAutoOpsRule":           Need(AutoOpsRuleWrite, InEnv("environment_id")),
		"DeleteAutoOpsRule":         Need(AutoOpsRuleWrite, InEnv("environment_id")),
		"UpdateAutoOpsRule":         Need(AutoOpsRuleWrite, InEnv("environment_id")),
		"ExecuteAutoOps":            Need(AutoOpsRuleWrite, InEnv("environment_id")),
		"ListOpsCounts":             Need(AutoOpsRuleRead, InEnv("environment_id")),
		"CreateProgressiveRollout":  Need(ProgressiveRolloutWrite, InEnv("environment_id")),
		"GetProgressiveRollout":     Need(ProgressiveRolloutRead, InEnv("environment_id")),
		"StopProgressiveRollout":    Need(ProgressiveRolloutWrite, InEnv("environment_id")),
		"DeleteProgressiveRollout":  Need(ProgressiveRolloutWrite, InEnv("environment_id")),
		"ListProgressiveRollouts":   Need(ProgressiveRolloutRead, InEnv("environment_id")),
		"ExecuteProgressiveRollout": Need(ProgressiveRolloutWrite, InEnv("environment_id")),
	}),
	rows(CodeReferenceService, map[string]Requirement{
		"GetCodeReference":    Need(CodeReferenceRead, InEnv("environment_id")),
		"ListCodeReferences":  Need(CodeReferenceRead, InEnv("environment_id")),
		"CreateCodeReference": Need(CodeReferenceWrite, InEnv("environment_id")),
		"UpdateCodeReference": Need(CodeReferenceWrite, InEnv("environment_id")),
		"DeleteCodeReference": Need(CodeReferenceWrite, InEnv("environment_id")),
	}),
	rows(EnvironmentService, map[string]Requirement{
		"GetEnvironmentV2":       Need(EnvironmentRead, InOrgOfEnv("id")),
		"ListEnvironmentsV2":     Need(EnvironmentRead, InOrg("organization_id")),
		"CreateEnvironmentV2":    Need(EnvironmentManage, InOrgOfProject("project_id")),
		"UpdateEnvironmentV2":    Need(EnvironmentManage, InOrgOfEnv("id")),
		"ArchiveEnvironmentV2":   Need(EnvironmentManage, InOrgOfEnv("id")),
		"UnarchiveEnvironmentV2": Need(EnvironmentManage, InOrgOfEnv("id")),
		"GetProject":             Need(ProjectRead, InOrgOfProject("id")),
		"ListProjects":           SystemAdminOnly(),
		"CreateProject":          Need(ProjectManage, InOrg("organization_id")),
		"UpdateProject":          Need(ProjectManage, InOrgOfProject("id")),
		"EnableProject":          Need(ProjectManage, InOrgOfProject("id")),
		"DisableProject":         Need(ProjectManage, InOrgOfProject("id")),
		"ConvertTrialProject":    Need(ProjectManage, InOrgOfProject("id")),
		"ListProjectsV2":         Need(ProjectRead, InOrg("organization_id")),
		"GetOrganization":        Need(OrganizationRead, InOrg("id")),
		"ListOrganizations":      SystemAdminOnly(),
		"CreateOrganization":     SystemAdminOnly(),
		// The handler also accepts the owner_email of the organization.
		"UpdateOrganization":       Need(OrganizationManage, InOrg("id")),
		"EnableOrganization":       SystemAdminOnly(),
		"DisableOrganization":      SystemAdminOnly(),
		"ArchiveOrganization":      SystemAdminOnly(),
		"UnarchiveOrganization":    SystemAdminOnly(),
		"ConvertTrialOrganization": SystemAdminOnly(),
		"ExchangeDemoToken":        Public(),
		"CreateDemoOrganization":   Public(),
	}),
	rows(EventCounterService, map[string]Requirement{
		"GetExperimentEvaluationCount": Need(EventCounterRead, InEnv("environment_id")),
		"GetEvaluationTimeseriesCount": Need(EventCounterRead, InEnv("environment_id")),
		"GetExperimentResult":          Need(EventCounterRead, InEnv("environment_id")),
		"ListExperimentResults":        Need(EventCounterRead, InEnv("environment_id")),
		"GetExperimentGoalCount":       Need(EventCounterRead, InEnv("environment_id")),
		"GetOpsEvaluationUserCount":    Need(EventCounterRead, InEnv("environment_id")),
		"GetOpsGoalUserCount":          Need(EventCounterRead, InEnv("environment_id")),
	}),
	rows(ExperimentService, map[string]Requirement{
		"GetGoal":          Need(GoalRead, InEnv("environment_id")),
		"ListGoals":        Need(GoalRead, InEnv("environment_id")),
		"CreateGoal":       Need(GoalWrite, InEnv("environment_id")),
		"UpdateGoal":       Need(GoalWrite, InEnv("environment_id")),
		"DeleteGoal":       Need(GoalWrite, InEnv("environment_id")),
		"GetExperiment":    Need(ExperimentRead, InEnv("environment_id")),
		"ListExperiments":  Need(ExperimentRead, InEnv("environment_id")),
		"CreateExperiment": Need(ExperimentWrite, InEnv("environment_id")),
		"UpdateExperiment": Need(ExperimentWrite, InEnv("environment_id")),
		"DeleteExperiment": Need(ExperimentWrite, InEnv("environment_id")),
	}),
	rows(FeatureService, map[string]Requirement{
		"GetFeature":            Need(FeatureRead, InEnv("environment_id")),
		"GetFeatures":           Need(FeatureRead, InEnv("environment_id")),
		"ListFeatures":          Need(FeatureRead, InEnv("environment_id")),
		"ListEnabledFeatures":   Need(FeatureRead, InEnv("environment_id")),
		"CreateFeature":         Need(FeatureWrite, InEnv("environment_id")),
		"UpdateFeature":         Need(FeatureWrite, InEnv("environment_id")),
		"DeleteFeature":         Need(FeatureWrite, InEnv("environment_id")),
		"EvaluateFeatures":      Need(FeatureRead, InEnv("environment_id")),
		"DebugEvaluateFeatures": Need(FeatureRead, InEnv("environment_id")),
		"CloneFeature": Need(FeatureRead, InEnv("environment_id")).
			And(Need(FeatureWrite, InEnv("target_environment_id"))),
		"BulkCloneFeature": Need(FeatureRead, InEnv("environment_id")).
			And(Need(FeatureWrite, EachEnv("target_environment_ids"))),
		"CreateScheduledFlagChange":     Need(ScheduledFlagChangeWrite, InEnv("environment_id")),
		"GetScheduledFlagChange":        Need(ScheduledFlagChangeRead, InEnv("environment_id")),
		"UpdateScheduledFlagChange":     Need(ScheduledFlagChangeWrite, InEnv("environment_id")),
		"DeleteScheduledFlagChange":     Need(ScheduledFlagChangeWrite, InEnv("environment_id")),
		"ListScheduledFlagChanges":      Need(ScheduledFlagChangeRead, InEnv("environment_id")),
		"ExecuteScheduledFlagChange":    Need(ScheduledFlagChangeWrite, InEnv("environment_id")),
		"GetScheduledFlagChangeSummary": Need(ScheduledFlagChangeRead, InEnv("environment_id")),
		"CreateSegment":                 Need(SegmentWrite, InEnv("environment_id")),
		"GetSegment":                    Need(SegmentRead, InEnv("environment_id")),
		"ListSegments":                  Need(SegmentRead, InEnv("environment_id")),
		"DeleteSegment":                 Need(SegmentWrite, InEnv("environment_id")),
		"UpdateSegment":                 Need(SegmentWrite, InEnv("environment_id")),
		"ListSegmentUsers":              Need(SegmentRead, InEnv("environment_id")),
		"BulkUploadSegmentUsers":        Need(SegmentWrite, InEnv("environment_id")),
		"BulkDownloadSegmentUsers":      Need(SegmentRead, InEnv("environment_id")),
		"ListTags":                      Need(TagRead, InEnv("environment_id")),
		"CreateFlagTrigger":             Need(FeatureWrite, InEnv("environment_id")),
		"UpdateFlagTrigger":             Need(FeatureWrite, InEnv("environment_id")),
		"DeleteFlagTrigger":             Need(FeatureWrite, InEnv("environment_id")),
		"GetFlagTrigger":                Need(FeatureRead, InEnv("environment_id")),
		"ListFlagTriggers":              Need(FeatureRead, InEnv("environment_id")),
		"GetUserAttributeKeys":          Need(UserAttributeRead, InEnv("environment_id")),
		"FlagTriggerWebhook":            Public(),
		"CreateHoldoutGroup":            Need(HoldoutGroupManage, InOrgOfEnv("environment_id")),
		"GetHoldoutGroup":               Need(HoldoutGroupRead, InEnv("environment_id")),
		"ListHoldoutGroups":             Need(HoldoutGroupRead, InEnv("environment_id")),
		"UpdateHoldoutGroup":            Need(HoldoutGroupManage, InOrgOfEnv("environment_id")),
		"StartHoldoutGroup":             Need(HoldoutGroupManage, InOrgOfEnv("environment_id")),
		"StartHoldoutGroupAnalysis":     Need(HoldoutGroupManage, InOrgOfEnv("environment_id")),
		"StopHoldoutGroup":              Need(HoldoutGroupManage, InOrgOfEnv("environment_id")),
		"ArchiveHoldoutGroup":           Need(HoldoutGroupManage, InOrgOfEnv("environment_id")),
		"DeleteHoldoutGroup":            Need(HoldoutGroupManage, InOrgOfEnv("environment_id")),
		"AddHoldoutGroupFeatures":       Need(HoldoutGroupManage, InOrgOfEnv("environment_id")),
		"RemoveHoldoutGroupFeatures":    Need(HoldoutGroupManage, InOrgOfEnv("environment_id")),
	}),
	rows(InsightsService, map[string]Requirement{
		"GetInsightsMonthlySummary": Need(InsightsRead, EachEnv("environment_ids")),
		"GetInsightsLatency":        Need(InsightsRead, EachEnv("environment_ids")),
		"GetInsightsRequests":       Need(InsightsRead, EachEnv("environment_ids")),
		"GetInsightsEvaluations":    Need(InsightsRead, EachEnv("environment_ids")),
		"GetInsightsErrorRates":     Need(InsightsRead, EachEnv("environment_ids")),
	}),
	rows(NotificationService, map[string]Requirement{
		"ListNotifications":           Authenticated(),
		"GetNotification":             Authenticated(),
		"GetNotificationUnreadCount":  Authenticated(),
		"MarkNotificationsAsRead":     Authenticated(),
		"MarkAllNotificationsAsRead":  Authenticated(),
		"ListDraftAdminNotifications": SystemAdminOnly(),
		"CreateAdminNotification":     SystemAdminOnly(),
		"UpdateAdminNotification":     SystemAdminOnly(),
		"PublishAdminNotification":    SystemAdminOnly(),
		"DeleteAdminNotification":     SystemAdminOnly(),
	}),
	rows(PushService, map[string]Requirement{
		"ListPushes": Need(PushRead, FirstSet(InOrg("organization_id"), InEnv("environment_id"))),
		"CreatePush": Need(PushWrite, InEnv("environment_id")),
		"UpdatePush": Need(PushWrite, InEnv("environment_id")),
		"GetPush":    Need(PushRead, InEnv("environment_id")),
		"DeletePush": Need(PushWrite, InEnv("environment_id")),
	}),
	rows(SubscriptionService, map[string]Requirement{
		"GetAdminSubscription":          SystemAdminOnly(),
		"ListAdminSubscriptions":        SystemAdminOnly(),
		"ListEnabledAdminSubscriptions": SystemAdminOnly(),
		"CreateAdminSubscription":       SystemAdminOnly(),
		"DeleteAdminSubscription":       SystemAdminOnly(),
		"EnableAdminSubscription":       SystemAdminOnly(),
		"DisableAdminSubscription":      SystemAdminOnly(),
		"UpdateAdminSubscription":       SystemAdminOnly(),
		"GetSubscription":               Need(SubscriptionRead, InEnv("environment_id")),
		"ListSubscriptions": Need(SubscriptionRead,
			FirstSet(InOrg("organization_id"), InEnv("environment_id"))),
		"ListEnabledSubscriptions": Need(SubscriptionRead, InEnv("environment_id")),
		"CreateSubscription":       Need(SubscriptionWrite, InEnv("environment_id")),
		"DeleteSubscription":       Need(SubscriptionWrite, InEnv("environment_id")),
		"UpdateSubscription":       Need(SubscriptionWrite, InEnv("environment_id")),
	}),
	rows(TagService, map[string]Requirement{
		"ListTags":  Need(TagRead, FirstSet(InOrg("organization_id"), InEnv("environment_id"))),
		"CreateTag": Need(TagWrite, InEnv("environment_id")),
		"DeleteTag": Need(TagWrite, InEnv("environment_id")),
	}),
	rows(TeamService, map[string]Requirement{
		"CreateTeam": Need(TeamManage, InOrg("organization_id")),
		"DeleteTeam": Need(TeamManage, InOrg("organization_id")),
		"ListTeams":  Need(TeamRead, InOrg("organization_id")),
	}),
)

// auditLogScope is checkScopeRole of the audit log service.
func auditLogScope() Requirement {
	return Need(AuditLogRead, InOrg("organization_id").IfSet()).Or(SystemAdmin()).
		And(Need(AuditLogRead, InEnv("environment_id").IfSet()))
}

// rows keys the rows of one service by full method name.
func rows(name string, methods map[string]Requirement) Rules {
	rules := make(Rules, len(methods))
	for method, r := range methods {
		rules["/"+name+"/"+method] = r
	}
	return rules
}

func merge(sets ...Rules) Rules {
	rules := make(Rules)
	for _, set := range sets {
		for method, r := range set {
			rules[method] = r
		}
	}
	return rules
}
