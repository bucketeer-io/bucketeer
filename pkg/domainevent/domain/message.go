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

package domain

import (
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	proto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func LocalizedMessage(eventType proto.Event_Type, localizer locale.Localizer) *proto.LocalizedMessage {
	// handle loc if multi-lang is necessary
	switch eventType {
	case proto.Event_UNKNOWN:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.UnknownOperation),
		}
	case proto.Event_FEATURE_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlag),
			),
		}
	case proto.Event_FEATURE_RENAMED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlag),
			),
		}
	case proto.Event_FEATURE_ENABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.EnabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlag),
			),
		}
	case proto.Event_FEATURE_DISABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DisabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlag),
			),
		}
	case proto.Event_FEATURE_ARCHIVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ArchivedTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlag),
			),
		}
	case proto.Event_FEATURE_UNARCHIVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.UnarchivedTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlag),
			),
		}
	case proto.Event_FEATURE_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlag),
			),
		}
	case proto.Event_FEATURE_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DescriptionUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlag),
			),
		}
	case proto.Event_FEATURE_VARIATION_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlagVariation),
			),
		}
	case proto.Event_FEATURE_VARIATION_REMOVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlagVariation),
			),
		}
	case proto.Event_FEATURE_OFF_VARIATION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.OffVariation),
			),
		}
	case proto.Event_VARIATION_VALUE_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ValueUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Variation),
			),
		}
	case proto.Event_VARIATION_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Variation),
			),
		}
	case proto.Event_VARIATION_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DescriptionUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Variation),
			),
		}
	case proto.Event_VARIATION_USER_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.IndividualUser),
			),
		}
	case proto.Event_VARIATION_USER_REMOVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.IndividualUser),
			),
		}
	case proto.Event_FEATURE_RULE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Rule),
			),
		}
	case proto.Event_FEATURE_RULE_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.RuleStrategy),
			),
		}
	case proto.Event_FEATURE_RULE_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Rule),
			),
		}
	case proto.Event_RULE_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ConditionAddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Rule)),
		}
	case proto.Event_RULE_CLAUSE_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ConditionDeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Rule),
			),
		}
	case proto.Event_RULE_FIXED_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.RuleFixedStrategyVariation),
			),
		}
	case proto.Event_RULE_ROLLOUT_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.RuleRolloutStrategyVariation),
			),
		}
	case proto.Event_CLAUSE_ATTRIBUTE_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.RuleAttribute),
			),
		}
	case proto.Event_CLAUSE_OPERATOR_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.RuleOperator),
			),
		}
	case proto.Event_CLAUSE_VALUE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ValueAddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.RuleCondition),
			),
		}
	case proto.Event_CLAUSE_VALUE_REMOVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ValueAddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.RuleCondition),
			),
		}
	case proto.Event_FEATURE_DEFAULT_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.DefaultStrategy),
			),
		}
	case proto.Event_FEATURE_TAG_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Tag),
			),
		}
	case proto.Event_FEATURE_TAG_REMOVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Tag),
			),
		}
	case proto.Event_FEATURE_VERSION_INCREMENTED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.IncrementedTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlagVersion),
			),
		}
	case proto.Event_FEATURE_CLONED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ClonedTemplate,
				localizer.MustLocalizeWithTemplate(locale.FeatureFlag),
			),
		}
	case proto.Event_SAMPLING_SEED_RESET:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ResetTemplate,
				localizer.MustLocalizeWithTemplate(locale.RandomSampling),
			),
		}
	case proto.Event_PREREQUISITE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Prerequisite),
			),
		}
	case proto.Event_PREREQUISITE_REMOVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Prerequisite),
			),
		}
	case proto.Event_PREREQUISITE_VARIATION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.PrerequisiteVariation),
			),
		}
	case proto.Event_GOAL_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Goal),
			),
		}
	case proto.Event_GOAL_RENAMED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Goal),
			),
		}
	case proto.Event_GOAL_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DescriptionUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Goal),
			),
		}
	case proto.Event_GOAL_ARCHIVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ArchivedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Goal),
			),
		}
	case proto.Event_GOAL_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Goal),
			),
		}
	case proto.Event_EXPERIMENT_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Experiment),
			),
		}
	case proto.Event_EXPERIMENT_STOPPED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.StoppedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Experiment),
			),
		}
	case proto.Event_EXPERIMENT_START_AT_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.ExperimentStartDate),
			),
		}
	case proto.Event_EXPERIMENT_STOP_AT_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.ExperimentEndDate)),
		}
	case proto.Event_EXPERIMENT_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Experiment)),
		}
	case proto.Event_EXPERIMENT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DescriptionUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Experiment)),
		}
	case proto.Event_EXPERIMENT_ARCHIVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ArchivedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Experiment),
			),
		}
	case proto.Event_EXPERIMENT_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Experiment)),
		}
	case proto.Event_EXPERIMENT_PERIOD_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.ExperimentPeriod),
			),
		}
	case proto.Event_EXPERIMENT_STARTED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.StartedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Experiment),
			),
		}
	case proto.Event_EXPERIMENT_FINISHED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.FinishedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Experiment)),
		}
	case proto.Event_ACCOUNT_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Account),
			),
		}
	case proto.Event_ACCOUNT_ROLE_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AccountRole),
			),
		}
	case proto.Event_ACCOUNT_ENABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.EnabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.Account),
			),
		}
	case proto.Event_ACCOUNT_DISABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DisabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.Account),
			),
		}
	case proto.Event_ACCOUNT_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Account),
			),
		}
	case proto.Event_APIKEY_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.APIKey),
			),
		}
	case proto.Event_APIKEY_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.APIKey),
			),
		}
	case proto.Event_APIKEY_ENABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.EnabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.APIKey),
			),
		}
	case proto.Event_APIKEY_DISABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DisabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.APIKey),
			),
		}
	case proto.Event_SEGMENT_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Segment),
			),
		}
	case proto.Event_SEGMENT_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Segment),
			),
		}
	case proto.Event_SEGMENT_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Segment),
			),
		}
	case proto.Event_SEGMENT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DescriptionUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Segment),
			),
		}
	case proto.Event_SEGMENT_RULE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentRule),
			),
		}
	case proto.Event_SEGMENT_RULE_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentRule),
			),
		}
	case proto.Event_SEGMENT_RULE_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ConditionAddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentRule)),
		}
	case proto.Event_SEGMENT_RULE_CLAUSE_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ConditionDeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentRule)),
		}
	case proto.Event_SEGMENT_CLAUSE_ATTRIBUTE_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentRuleAttribute),
			),
		}
	case proto.Event_SEGMENT_CLAUSE_OPERATOR_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentRuleOperator),
			),
		}
	case proto.Event_SEGMENT_CLAUSE_VALUE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ValueAddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentRule),
			),
		}
	case proto.Event_SEGMENT_CLAUSE_VALUE_REMOVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ValueDeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentRule),
			),
		}
	case proto.Event_SEGMENT_USER_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentUser),
			),
		}
	case proto.Event_SEGMENT_USER_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentUser),
			),
		}
	case proto.Event_SEGMENT_BULK_UPLOAD_USERS:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.UploadedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentUser),
			),
		}
	case proto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.UpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.SegmentUserUploadStatus),
			),
		}
	case proto.Event_ENVIRONMENT_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Environment),
			),
		}
	case proto.Event_ENVIRONMENT_RENAMED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Environment),
			),
		}
	case proto.Event_ENVIRONMENT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DescriptionUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Environment),
			),
		}
	case proto.Event_ENVIRONMENT_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Environment),
			),
		}
	case proto.Event_ENVIRONMENT_V2_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Environment),
			),
		}
	case proto.Event_ENVIRONMENT_V2_RENAMED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Environment),
			),
		}
	case proto.Event_ENVIRONMENT_V2_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DescriptionUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Environment),
			),
		}
	case proto.Event_ENVIRONMENT_V2_ARCHIVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ArchivedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Environment),
			),
		}
	case proto.Event_ENVIRONMENT_V2_UNARCHIVED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.UnarchivedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Environment)),
		}
	case proto.Event_ADMIN_ACCOUNT_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AdminAccount),
			),
		}
	case proto.Event_ADMIN_ACCOUNT_ENABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.EnabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.AdminAccount),
			),
		}
	case proto.Event_ADMIN_ACCOUNT_DISABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DisabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.AdminAccount),
			),
		}
	case proto.Event_AUTOOPS_RULE_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AutoOperation),
			),
		}
	case proto.Event_AUTOOPS_RULE_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AutoOperation),
			),
		}
	case proto.Event_AUTOOPS_RULE_OPS_TYPE_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AutoOperationType),
			),
		}
	case proto.Event_AUTOOPS_RULE_CLAUSE_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ConditionDeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AutoOperation),
			),
		}
	case proto.Event_AUTOOPS_RULE_TRIGGERED_AT_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AutoOperationExecutionTime),
			),
		}
	case proto.Event_OPS_EVENT_RATE_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.EventRate),
			),
		}
	case proto.Event_OPS_EVENT_RATE_CLAUSE_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.EventRate),
			),
		}
	case proto.Event_DATETIME_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Datetime),
			),
		}
	case proto.Event_DATETIME_CLAUSE_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Datetime),
			),
		}
	case proto.Event_PUSH_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Push),
			),
		}
	case proto.Event_PUSH_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Push),
			),
		}
	case proto.Event_PUSH_TAGS_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Push),
			),
		}
	case proto.Event_PUSH_TAGS_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.PushTag),
			),
		}
	case proto.Event_PUSH_RENAMED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Push),
			),
		}
	case proto.Event_SUBSCRIPTION_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Notification),
			),
		}
	case proto.Event_SUBSCRIPTION_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Notification),
			),
		}
	case proto.Event_SUBSCRIPTION_ENABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.EnabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.Notification),
			),
		}
	case proto.Event_SUBSCRIPTION_DISABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DisabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.Notification),
			),
		}
	case proto.Event_SUBSCRIPTION_SOURCE_TYPE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.NotificationType),
			),
		}
	case proto.Event_SUBSCRIPTION_SOURCE_TYPE_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.NotificationType),
			),
		}
	case proto.Event_SUBSCRIPTION_RENAMED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Notification),
			),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AdminNotification),
			),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AdminNotification),
			),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_ENABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.EnabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.AdminNotification),
			),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_DISABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DisabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.AdminNotification),
			),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_SOURCE_TYPE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AdminNotificationType),
			),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_SOURCE_TYPE_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AdminNotificationType),
			),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_RENAMED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.AdminNotification),
			),
		}
	case proto.Event_PROJECT_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Project),
			),
		}
	case proto.Event_PROJECT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DescriptionUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Project),
			),
		}
	case proto.Event_PROJECT_ENABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.EnabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.Project),
			),
		}
	case proto.Event_PROJECT_DISABLED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DisabledTemplate,
				localizer.MustLocalizeWithTemplate(locale.Project),
			),
		}
	case proto.Event_PROJECT_TRIAL_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.TrialProject)),
		}
	case proto.Event_PROJECT_TRIAL_CONVERTED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.TrialConverted),
		}
	case proto.Event_PROJECT_RENAMED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Project),
			),
		}
	case proto.Event_WEBHOOK_CREATED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.CreatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Webhook),
			),
		}
	case proto.Event_WEBHOOK_DELETED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DeletedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Webhook),
			),
		}
	case proto.Event_WEBHOOK_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.NameUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Webhook),
			),
		}
	case proto.Event_WEBHOOK_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.DescriptionUpdatedTemplate,
				localizer.MustLocalizeWithTemplate(locale.Webhook),
			),
		}
	case proto.Event_WEBHOOK_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.AddedTemplate,
				localizer.MustLocalizeWithTemplate(locale.WebhookRule),
			),
		}
	case proto.Event_WEBHOOK_CLAUSE_CHANGED:
		return &proto.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.ChangedTemplate,
				localizer.MustLocalizeWithTemplate(locale.WebhookRule),
			),
		}
	case proto.Event_PROGRESSIVE_ROLLOUT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Progressive Rolloutを作成しました",
		}
	case proto.Event_PROGRESSIVE_ROLLOUT_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Progressive Rolloutを削除しました",
		}
	case proto.Event_PROGRESSIVE_ROLLOUT_SCHEDULE_TRIGGERED_AT_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Progressive Rolloutの実行時間が変更されました",
		}
	}
	return &proto.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalize(locale.UnknownOperation),
	}
}
