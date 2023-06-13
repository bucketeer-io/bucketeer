// Copyright 2022 The Bucketeer Authors.
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
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "feature flag"),
		}
	case proto.Event_FEATURE_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "feature flagの名前"),
		}
	case proto.Event_FEATURE_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.EnabledTemplate, "feature flag"),
		}
	case proto.Event_FEATURE_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DisabledTemplate, "feature flag"),
		}
	case proto.Event_FEATURE_ARCHIVED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ArchivedTemplate, "feature flag"),
		}
	case proto.Event_FEATURE_UNARCHIVED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.UnarchivedTemplate, "feature flag"),
		}
	case proto.Event_FEATURE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "feature flag"),
		}
	case proto.Event_FEATURE_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "feature flagの説明文"),
		}
	case proto.Event_FEATURE_VARIATION_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "feature flagにvariation"),
		}
	case proto.Event_FEATURE_VARIATION_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "feature flagのvariation"),
		}
	case proto.Event_FEATURE_OFF_VARIATION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "feature flagの無効時のvariation"),
		}
	case proto.Event_VARIATION_VALUE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "variationの値"),
		}
	case proto.Event_VARIATION_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "variationの名前"),
		}
	case proto.Event_VARIATION_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "variationの説明文"),
		}
	case proto.Event_VARIATION_USER_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "variationを適用するユーザー"),
		}
	case proto.Event_VARIATION_USER_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "variationを適用するユーザー"),
		}
	case proto.Event_FEATURE_RULE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "rule"),
		}
	case proto.Event_FEATURE_RULE_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "ruleの適用するvariationの選択方法"),
		}
	case proto.Event_FEATURE_RULE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "rule"),
		}
	case proto.Event_RULE_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "ruleの条件"),
		}
	case proto.Event_RULE_CLAUSE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "ruleの条件"),
		}
	case proto.Event_RULE_FIXED_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "ruleの適用するvariationの種類"),
		}
	case proto.Event_RULE_ROLLOUT_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "ruleの適用するvariationの適用割合"),
		}
	case proto.Event_CLAUSE_ATTRIBUTE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "ruleの条件のattribute"),
		}
	case proto.Event_CLAUSE_OPERATOR_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "ruleの条件のoperator"),
		}
	case proto.Event_CLAUSE_VALUE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "ruleの条件の対象の値"),
		}
	case proto.Event_CLAUSE_VALUE_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "ruleの条件の対象の値"),
		}
	case proto.Event_FEATURE_DEFAULT_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "feature flagがデフォルトで適用する条件"),
		}
	case proto.Event_FEATURE_TAG_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "タグ"),
		}
	case proto.Event_FEATURE_TAG_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "タグ"),
		}
	case proto.Event_FEATURE_VERSION_INCREMENTED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.UpdatedTemplate, "feature flagのバージョン"),
		}
	case proto.Event_FEATURE_CLONED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ClonedTemplate, "feature flag"),
		}
	case proto.Event_SAMPLING_SEED_RESET:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ResetTemplate, "ランダムサンプリング"),
		}
	case proto.Event_PREREQUISITE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "prerequisite"),
		}
	case proto.Event_PREREQUISITE_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "prerequisite"),
		}
	case proto.Event_PREREQUISITE_VARIATION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "prerequisiteのvariation"),
		}
	case proto.Event_GOAL_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "goal"),
		}
	case proto.Event_GOAL_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "goalの名前"),
		}
	case proto.Event_GOAL_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "goalの説明文"),
		}
	case proto.Event_GOAL_ARCHIVED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ArchivedTemplate, "goal"),
		}
	case proto.Event_GOAL_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "goal"),
		}
	case proto.Event_EXPERIMENT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "experiment"),
		}
	case proto.Event_EXPERIMENT_STOPPED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.StoppedTemplate, "experiment"),
		}
	case proto.Event_EXPERIMENT_START_AT_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "experimentの開始時間"),
		}
	case proto.Event_EXPERIMENT_STOP_AT_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "experimentの終了時間"),
		}
	case proto.Event_EXPERIMENT_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "experimentの名前"),
		}
	case proto.Event_EXPERIMENT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "experimentの説明文"),
		}
	case proto.Event_EXPERIMENT_ARCHIVED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ArchivedTemplate, "experiment"),
		}
	case proto.Event_EXPERIMENT_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "experiment"),
		}
	case proto.Event_EXPERIMENT_PERIOD_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "experimentの期間"),
		}
	case proto.Event_EXPERIMENT_STARTED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.StartedTemplate, "experiment"),
		}
	case proto.Event_EXPERIMENT_FINISHED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.FinishedTemplate, "experiment"),
		}
	case proto.Event_ACCOUNT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "アカウント"),
		}
	case proto.Event_ACCOUNT_ROLE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "アカウントの権限"),
		}
	case proto.Event_ACCOUNT_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.EnabledTemplate, "アカウント"),
		}
	case proto.Event_ACCOUNT_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DisabledTemplate, "アカウント"),
		}
	case proto.Event_ACCOUNT_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "アカウント"),
		}
	case proto.Event_APIKEY_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "APIキー"),
		}
	case proto.Event_APIKEY_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "APIキーの名前"),
		}
	case proto.Event_APIKEY_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.EnabledTemplate, "APIキー"),
		}
	case proto.Event_APIKEY_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DisabledTemplate, "APIキー"),
		}
	case proto.Event_SEGMENT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "segment"),
		}
	case proto.Event_SEGMENT_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "segment"),
		}
	case proto.Event_SEGMENT_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "segmentの名前"),
		}
	case proto.Event_SEGMENT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "segmentの説明文"),
		}
	case proto.Event_SEGMENT_RULE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "segmentにrule"),
		}
	case proto.Event_SEGMENT_RULE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "segmentからrule"),
		}
	case proto.Event_SEGMENT_RULE_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "segmentのruleに条件"),
		}
	case proto.Event_SEGMENT_RULE_CLAUSE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "segmentのruleから条件"),
		}
	case proto.Event_SEGMENT_CLAUSE_ATTRIBUTE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "segmentのruleの条件のattribute"),
		}
	case proto.Event_SEGMENT_CLAUSE_OPERATOR_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "segmentのruleの条件のoperator"),
		}
	case proto.Event_SEGMENT_CLAUSE_VALUE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "segmentのruleの条件の対象の値"),
		}
	case proto.Event_SEGMENT_CLAUSE_VALUE_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "segmentのruleの条件の対象の値"),
		}
	case proto.Event_SEGMENT_USER_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "segmentにユーザー"),
		}
	case proto.Event_SEGMENT_USER_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "segmentからユーザー"),
		}
	case proto.Event_SEGMENT_BULK_UPLOAD_USERS:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.UploadedTemplate, "ユーザーセグメントファイル"),
		}
	case proto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "ユーザーセグメントファイルのアップロードステータス"),
		}
	case proto.Event_ENVIRONMENT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "Environment"),
		}
	case proto.Event_ENVIRONMENT_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "Environmentの名前"),
		}
	case proto.Event_ENVIRONMENT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "Environmentの説明文"),
		}
	case proto.Event_ENVIRONMENT_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "Environment"),
		}
	case proto.Event_ADMIN_ACCOUNT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "管理者アカウント"),
		}
	case proto.Event_ADMIN_ACCOUNT_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.EnabledTemplate, "管理者アカウント"),
		}
	case proto.Event_ADMIN_ACCOUNT_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DisabledTemplate, "管理者アカウント"),
		}
	case proto.Event_AUTOOPS_RULE_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "自動オペレーションルール"),
		}
	case proto.Event_AUTOOPS_RULE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "自動オペレーションルール"),
		}
	case proto.Event_AUTOOPS_RULE_OPS_TYPE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "オペレーションタイプ"),
		}
	case proto.Event_AUTOOPS_RULE_CLAUSE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "オペレーションルール"),
		}
	case proto.Event_AUTOOPS_RULE_TRIGGERED_AT_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "自動オペレーションの実行時間"),
		}
	case proto.Event_OPS_EVENT_RATE_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "イベントレートルール"),
		}
	case proto.Event_OPS_EVENT_RATE_CLAUSE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "イベントレートルール"),
		}
	case proto.Event_DATETIME_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "日時ルール"),
		}
	case proto.Event_DATETIME_CLAUSE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "日時ルール"),
		}
	case proto.Event_PUSH_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "プッシュ設定"),
		}
	case proto.Event_PUSH_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "プッシュ設定"),
		}
	case proto.Event_PUSH_TAGS_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "プッシュ設定にタグ"),
		}
	case proto.Event_PUSH_TAGS_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "プッシュ設定からタグ"),
		}
	case proto.Event_PUSH_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "プッシュ設定の名前"),
		}
	case proto.Event_SUBSCRIPTION_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "通知設定"),
		}
	case proto.Event_SUBSCRIPTION_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "通知設定"),
		}
	case proto.Event_SUBSCRIPTION_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.EnabledTemplate, "通知設定"),
		}
	case proto.Event_SUBSCRIPTION_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DisabledTemplate, "通知設定"),
		}
	case proto.Event_SUBSCRIPTION_SOURCE_TYPE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "通知設定に通知設定ソース"),
		}
	case proto.Event_SUBSCRIPTION_SOURCE_TYPE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "通知設定から通知設定ソース"),
		}
	case proto.Event_SUBSCRIPTION_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "通知設定の名前"),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "管理者用通知設定"),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "管理者用通知設定"),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.EnabledTemplate, "管理者用通知設定"),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DisabledTemplate, "管理者用通知設定"),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_SOURCE_TYPE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "管理者用通知設定に通知設定ソース"),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_SOURCE_TYPE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "管理者用通知設定から通知設定ソース"),
		}
	case proto.Event_ADMIN_SUBSCRIPTION_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "管理者用通知設定の名前"),
		}
	case proto.Event_PROJECT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "Project"),
		}
	case proto.Event_PROJECT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "Projectの説明文"),
		}
	case proto.Event_PROJECT_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.EnabledTemplate, "Project"),
		}
	case proto.Event_PROJECT_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DisabledTemplate, "Project"),
		}
	case proto.Event_PROJECT_TRIAL_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "Trial Project"),
		}
	case proto.Event_PROJECT_TRIAL_CONVERTED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.TrialConverted),
		}
	case proto.Event_WEBHOOK_CREATED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.CreatedTemplate, "webhook"),
		}
	case proto.Event_WEBHOOK_DELETED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.DeletedTemplate, "webhook"),
		}
	case proto.Event_WEBHOOK_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "webhookの名前"),
		}
	case proto.Event_WEBHOOK_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "webhookの説明"),
		}
	case proto.Event_WEBHOOK_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.AddedTemplate, "webhookのルール"),
		}
	case proto.Event_WEBHOOK_CLAUSE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.ChangedTemplate, "webhookのルール"),
		}
	}
	return &proto.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalize(locale.UnknownOperation),
	}
}
