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

func LocalizedMessage(eventType proto.Event_Type, loc string) *proto.LocalizedMessage {
	// handle loc if multi-lang is necessary
	switch eventType {
	case proto.Event_UNKNOWN:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "不明な操作を実行しました",
		}
	case proto.Event_FEATURE_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagを作成しました",
		}
	case proto.Event_FEATURE_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagの名前を変更しました",
		}
	case proto.Event_FEATURE_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagを有効化しました",
		}
	case proto.Event_FEATURE_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagを無効化しました",
		}
	case proto.Event_FEATURE_ARCHIVED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagをアーカイブしました",
		}
	case proto.Event_FEATURE_UNARCHIVED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagをアーカイブから解除しました",
		}
	case proto.Event_FEATURE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagを削除しました",
		}
	case proto.Event_FEATURE_EVALUATION_DELAYABLE_SET:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagを初回リクエスト時にキューに入れるように変更されました",
		}
	case proto.Event_FEATURE_EVALUATION_UNDELAYABLE_SET:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagを初回リクエスト時にキューに入れないように変更されました",
		}
	case proto.Event_FEATURE_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagの説明文を変更しました",
		}
	case proto.Event_FEATURE_VARIATION_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagにvariationを追加しました",
		}
	case proto.Event_FEATURE_VARIATION_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagのvariationを削除しました",
		}
	case proto.Event_FEATURE_OFF_VARIATION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagの無効時のvariationを変更しました",
		}
	case proto.Event_VARIATION_VALUE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "variationの値を変更しました",
		}
	case proto.Event_VARIATION_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "variationの名前を変更しました",
		}
	case proto.Event_VARIATION_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "variationの説明文を変更しました",
		}
	case proto.Event_VARIATION_USER_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "variationを適用するユーザーを追加しました",
		}
	case proto.Event_VARIATION_USER_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "variationを適用するユーザーを削除しました",
		}
	case proto.Event_FEATURE_RULE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleを追加しました",
		}
	case proto.Event_FEATURE_RULE_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleの適用するvariationの選択方法を変更しました",
		}
	case proto.Event_FEATURE_RULE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleを削除しました",
		}
	case proto.Event_RULE_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleの条件を追加しました",
		}
	case proto.Event_RULE_CLAUSE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleの条件を削除しました",
		}
	case proto.Event_RULE_FIXED_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleの適用するvariationの種類を変更しました",
		}
	case proto.Event_RULE_ROLLOUT_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleの適用するvariationの適用割合を変更しました",
		}
	case proto.Event_CLAUSE_ATTRIBUTE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleの条件のattributeを変更しました",
		}
	case proto.Event_CLAUSE_OPERATOR_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleの条件のoperatorを変更しました",
		}
	case proto.Event_CLAUSE_VALUE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleの条件の対象の値を追加しました",
		}
	case proto.Event_CLAUSE_VALUE_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ruleの条件の対象の値を削除しました",
		}
	case proto.Event_FEATURE_DEFAULT_STRATEGY_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagがデフォルトで適用する条件を変更しました",
		}
	case proto.Event_FEATURE_TAG_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "タグを追加しました",
		}
	case proto.Event_FEATURE_TAG_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "タグを削除しました",
		}
	case proto.Event_FEATURE_VERSION_INCREMENTED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagのバージョンを更新しました",
		}
	case proto.Event_FEATURE_CLONED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "feature flagをクローンしました",
		}
	case proto.Event_SAMPLING_SEED_RESET:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ランダムサンプリングをリセットしました",
		}
	case proto.Event_PREREQUISITE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "prerequisiteを追加しました",
		}
	case proto.Event_PREREQUISITE_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "prerequisiteを削除しました",
		}
	case proto.Event_PREREQUISITE_VARIATION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "prerequisiteのvariationを変更しました",
		}
	case proto.Event_GOAL_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "goalを作成しました",
		}
	case proto.Event_GOAL_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "goalの名前を変更しました",
		}
	case proto.Event_GOAL_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "goalの説明文を変更しました",
		}
	case proto.Event_GOAL_ARCHIVED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "goalをアーカイブしました",
		}
	case proto.Event_GOAL_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "goalを削除しました",
		}
	case proto.Event_EXPERIMENT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentを作成しました",
		}
	case proto.Event_EXPERIMENT_STOPPED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentを停止しました",
		}
	case proto.Event_EXPERIMENT_START_AT_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentの開始時間を変更しました",
		}
	case proto.Event_EXPERIMENT_STOP_AT_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentの終了時間を変更しました",
		}
	case proto.Event_EXPERIMENT_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentの名前を変更しました",
		}
	case proto.Event_EXPERIMENT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentの説明文を変更しました",
		}
	case proto.Event_EXPERIMENT_ARCHIVED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentをアーカイブしました",
		}
	case proto.Event_EXPERIMENT_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentを削除しました",
		}
	case proto.Event_EXPERIMENT_PERIOD_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentの期間を変更しました",
		}
	case proto.Event_EXPERIMENT_STARTED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentが開始しました",
		}
	case proto.Event_EXPERIMENT_FINISHED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "experimentが終了しました",
		}
	case proto.Event_ACCOUNT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "アカウントを作成しました",
		}
	case proto.Event_ACCOUNT_ROLE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "アカウントの権限を変更しました",
		}
	case proto.Event_ACCOUNT_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "アカウントを有効化しました",
		}
	case proto.Event_ACCOUNT_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "アカウントを無効化しました",
		}
	case proto.Event_ACCOUNT_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "アカウントを削除しました",
		}
	case proto.Event_APIKEY_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "APIキーを作成しました",
		}
	case proto.Event_APIKEY_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "APIキーの名前を変更しました",
		}
	case proto.Event_APIKEY_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "APIキーを有効化しました",
		}
	case proto.Event_APIKEY_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "APIキーを無効化しました",
		}
	case proto.Event_SEGMENT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentを作成しました",
		}
	case proto.Event_SEGMENT_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentを削除しました",
		}
	case proto.Event_SEGMENT_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentの名前を変更しました",
		}
	case proto.Event_SEGMENT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentの説明文を変更しました",
		}
	case proto.Event_SEGMENT_RULE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentにruleを追加しました",
		}
	case proto.Event_SEGMENT_RULE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentからruleを削除しました",
		}
	case proto.Event_SEGMENT_RULE_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentのruleに条件を追加しました",
		}
	case proto.Event_SEGMENT_RULE_CLAUSE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentのruleから条件を削除しました",
		}
	case proto.Event_SEGMENT_CLAUSE_ATTRIBUTE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentのruleの条件のattributeを変更しました",
		}
	case proto.Event_SEGMENT_CLAUSE_OPERATOR_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentのruleの条件のoperatorを変更しました",
		}
	case proto.Event_SEGMENT_CLAUSE_VALUE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentのruleの条件の対象の値を追加しました",
		}
	case proto.Event_SEGMENT_CLAUSE_VALUE_REMOVED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentのruleの条件の対象の値を削除しました",
		}
	case proto.Event_SEGMENT_USER_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentにユーザーを追加しました",
		}
	case proto.Event_SEGMENT_USER_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "segmentからユーザーを削除しました",
		}
	case proto.Event_SEGMENT_BULK_UPLOAD_USERS:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ユーザーセグメントファイルをアップロードしました",
		}
	case proto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "ユーザーセグメントファイルのアップロードステータスが変わりました",
		}
	case proto.Event_ENVIRONMENT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Environmentを作成しました",
		}
	case proto.Event_ENVIRONMENT_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Environmentの名前を変更しました",
		}
	case proto.Event_ENVIRONMENT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Environmentの説明文を変更しました",
		}
	case proto.Event_ENVIRONMENT_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Environmentを削除しました",
		}
	case proto.Event_ADMIN_ACCOUNT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "管理者アカウントを作成しました",
		}
	case proto.Event_ADMIN_ACCOUNT_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "管理者アカウントを有効化しました",
		}
	case proto.Event_ADMIN_ACCOUNT_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "管理者アカウントを無効化しました",
		}
	case proto.Event_AUTOOPS_RULE_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "自動オペレーションルールを作成しました",
		}
	case proto.Event_AUTOOPS_RULE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "自動オペレーションルールを削除しました",
		}
	case proto.Event_AUTOOPS_RULE_OPS_TYPE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "オペレーションタイプを変更しました",
		}
	case proto.Event_AUTOOPS_RULE_CLAUSE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "オペレーションルールを削除しました",
		}
	case proto.Event_AUTOOPS_RULE_TRIGGERED_AT_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "自動オペレーションの実行時間が変更されました",
		}
	case proto.Event_OPS_EVENT_RATE_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "イベントレートルールが追加されました",
		}
	case proto.Event_OPS_EVENT_RATE_CLAUSE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "イベントレートルールが変更されました",
		}
	case proto.Event_DATETIME_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "日時ルールが追加されました",
		}
	case proto.Event_DATETIME_CLAUSE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "日時ルールが変更されました",
		}
	case proto.Event_PUSH_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "プッシュ設定を作成しました",
		}
	case proto.Event_PUSH_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "プッシュ設定を削除しました",
		}
	case proto.Event_PUSH_TAGS_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "プッシュ設定にタグを追加しました",
		}
	case proto.Event_PUSH_TAGS_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "プッシュ設定からタグを削除しました",
		}
	case proto.Event_PUSH_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "プッシュ設定の名前を変更しました",
		}
	case proto.Event_SUBSCRIPTION_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "通知設定を作成しました",
		}
	case proto.Event_SUBSCRIPTION_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "通知設定を削除しました",
		}
	case proto.Event_SUBSCRIPTION_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "通知設定を有効化しました",
		}
	case proto.Event_SUBSCRIPTION_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "通知設定を無効化しました",
		}
	case proto.Event_SUBSCRIPTION_SOURCE_TYPE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "通知設定に通知設定ソースを追加しました",
		}
	case proto.Event_SUBSCRIPTION_SOURCE_TYPE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "通知設定から通知設定ソースを削除しました",
		}
	case proto.Event_SUBSCRIPTION_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "通知設定の名前を変更しました",
		}
	case proto.Event_ADMIN_SUBSCRIPTION_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "管理者用通知設定を作成しました",
		}
	case proto.Event_ADMIN_SUBSCRIPTION_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "管理者用通知設定を削除しました",
		}
	case proto.Event_ADMIN_SUBSCRIPTION_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "管理者用通知設定を有効化しました",
		}
	case proto.Event_ADMIN_SUBSCRIPTION_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "管理者用通知設定を無効化しました",
		}
	case proto.Event_ADMIN_SUBSCRIPTION_SOURCE_TYPE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "管理者用通知設定に通知設定ソースを追加しました",
		}
	case proto.Event_ADMIN_SUBSCRIPTION_SOURCE_TYPE_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "管理者用通知設定から通知設定ソースを削除しました",
		}
	case proto.Event_ADMIN_SUBSCRIPTION_RENAMED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "管理者用通知設定の名前を変更しました",
		}
	case proto.Event_PROJECT_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Projectを作成しました",
		}
	case proto.Event_PROJECT_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Projectの説明文を変更しました",
		}
	case proto.Event_PROJECT_ENABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Projectを有効化しました",
		}
	case proto.Event_PROJECT_DISABLED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Projectを無効化しました",
		}
	case proto.Event_PROJECT_TRIAL_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Trial Projectを作成しました",
		}
	case proto.Event_PROJECT_TRIAL_CONVERTED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "Trialを正式なProjectに変換しました",
		}
	case proto.Event_WEBHOOK_CREATED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "webhookを作成しました",
		}
	case proto.Event_WEBHOOK_DELETED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "webhookを削除しました",
		}
	case proto.Event_WEBHOOK_NAME_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "webhookの名前を変更しました",
		}
	case proto.Event_WEBHOOK_DESCRIPTION_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "webhookの説明を変更しました",
		}
	case proto.Event_WEBHOOK_CLAUSE_ADDED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "webhookのルールが追加されました",
		}
	case proto.Event_WEBHOOK_CLAUSE_CHANGED:
		return &proto.LocalizedMessage{
			Locale:  locale.Ja,
			Message: "webhookのルールが変更されました",
		}
	}
	return &proto.LocalizedMessage{
		Locale:  locale.Ja,
		Message: "不明な操作を実行しました",
	}
}
