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

package api

import (
	"github.com/bucketeer-io/bucketeer/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
)

var (
	statusInternal                                = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "internal"))
	statusMissingFrom                             = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "missing from", "from"))
	statusMissingID                               = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing id", "id"))
	statusMissingIDs                              = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing ids", "ids"))
	statusInvalidID                               = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "invalid id", "id"))
	statusMissingUser                             = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing user", "user"))
	statusMissingUserID                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing user id", "user_id"))
	statusMissingUserIDs                          = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing user ids", "user_ids"))
	statusMissingFeatureIDs                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing feature ids", "feature_ids"))
	statusMissingCommand                          = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "missing command", "command"))
	statusMissingDefaultOnVariation               = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing default on variation", "default_on_variation"))
	statusMissingDefaultOffVariation              = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing default off variation", "default_off_variation"))
	statusInvalidDefaultOnVariation               = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "invalid default on variation", "default_on_variation"))
	statusInvalidDefaultOffVariation              = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "invalid default off variation", "default_off_variation"))
	statusMissingVariationID                      = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing variation id", "variation_id"))
	statusInvalidVariationID                      = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "invalid variation id", "variation_id"))
	statusDifferentVariationsSize                 = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "feature variations and rollout variations must have the same size", "feature_variations_and_rollout_variations"))
	statusExceededMaxVariationWeight              = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "the sum of all weights value is %d", "feature_variations_and_rollout_variations"))
	statusIncorrectVariationWeight                = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "command: variation weight must be between 0 and %d", "feature_variations_and_rollout_variations"))
	statusInvalidCursor                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "cursor is invalid", "cursor"))
	statusInvalidOrderBy                          = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "order_by is invalid", "order_by"))
	statusMissingName                             = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing name", "name"))
	statusMissingFeatureVariations                = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "feature must contain more than one variation", "feature_variations"))
	statusMissingFeatureTags                      = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "feature must contain one or more tags", "feature_tags"))
	statusUnknownCommand                          = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "unknown command", "command"))
	statusCommentRequiredForUpdating              = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "a comment is required for updating"))
	statusMissingRule                             = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing rule", "rule"))
	statusMissingRuleID                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing rule id", "rule_id"))
	statusMissingRuleClause                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing rule clause", "rule_clause"))
	statusMissingClauseID                         = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing clause id", "clause_id"))
	statusMissingClauseAttribute                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing clause attribute", "clause_attribute"))
	statusMissingClauseValues                     = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing clause values", "clause_values"))
	statusMissingClauseValue                      = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing clause value", "clause_value"))
	statusMissingSegmentID                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing segment id", "segment_id"))
	statusMissingSegmentUsersData                 = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing segment users data", "segment_users_data"))
	statusMissingRuleStrategy                     = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing rule strategy", "rule_strategy"))
	statusUnknownStrategy                         = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "unknown strategy", "strategy"))
	statusMissingFixedStrategy                    = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing fixed strategy", "fixed_strategy"))
	statusMissingRolloutStrategy                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing rollout strategy", "rollout_strategy"))
	statusExceededMaxSegmentUsersDataSize         = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "max segment users data size allowed is %d bytes", "segment_users_data"))
	statusUnknownSegmentUserState                 = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "unknown segment user state", "segment_user_state"))
	statusIncorrectUUIDFormat                     = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "uuid format must be an uuid version 4", "uuid"))
	statusExceededMaxUserIDsLength                = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "max user ids length allowed is %d", "user_ids"))
	statusIncorrectDestinationEnvironment         = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "destination environment is the same as origin one", "destination_environment"))
	statusExceededMaxPageSizePerRequest           = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "max page size allowed is %d", "page_size"))
	statusNotFound                                = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "not found", "feature"))
	statusSegmentNotFound                         = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "segment not found", "segment"))
	statusAlreadyExists                           = api.NewGRPCStatus(pkgErr.NewErrorAlreadyExists(pkgErr.FeaturePackageName, "already exists"))
	statusSegmentUsersAlreadyUploading            = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "segment users already uploading"))
	statusSegmentStatusNotSuceeded                = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "segment status is not suceeded"))
	statusSegmentInUse                            = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "segment is in use"))
	statusUnauthenticated                         = api.NewGRPCStatus(pkgErr.NewErrorUnauthenticated(pkgErr.FeaturePackageName, "unauthenticated"))
	statusPermissionDenied                        = api.NewGRPCStatus(pkgErr.NewErrorPermissionDenied(pkgErr.FeaturePackageName, "permission denied"))
	statusWaitingOrRunningExperimentExists        = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "experiment in waiting or running status exists"))
	statusCycleExists                             = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "circular dependency detected"))
	statusInvalidArchive                          = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "cant't archive because this feature is used as a prerequsite"))
	statusInvalidChangingVariation                = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "can't change or remove this variation because it is used as a prerequsite"))
	statusVariationInUseByOtherFeatures           = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "can't remove this variation because it is used as a prerequisite or rule in other features"))
	statusInvalidPrerequisite                     = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "invalid prerequisite"))
	statusProgressiveRolloutWaitingOrRunningState = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "there is a progressive rollout in the waiting or running state"))
	// flag trigger
	statusMissingTriggerFeatureID  = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing trigger feature id", "trigger_feature_id"))
	statusMissingTriggerType       = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing trigger type", "trigger_type"))
	statusMissingTriggerAction     = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing trigger action", "trigger_action"))
	statusMissingTriggerID         = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing trigger id", "trigger_id"))
	statusSecretRequired           = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "trigger secret is required", "trigger_secret"))
	statusTriggerAlreadyDisabled   = api.NewGRPCStatus(pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "trigger already disabled"))
	statusTriggerNotFound          = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "trigger not found", "trigger"))
	statusTriggerDisableFailed     = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "trigger disable failed"))
	statusTriggerEnableFailed      = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "trigger enable failed"))
	statusTriggerActionInvalid     = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "trigger action is invalid", "trigger_action"))
	statusTriggerUsageUpdateFailed = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "trigger usage update failed"))
)
