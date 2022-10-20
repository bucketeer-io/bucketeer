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

package api

import (
	"regexp"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/feature/command"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	maxPageSizePerRequest   = 100
	maxUserIDsLength        = 100000
	maxSegmentUsersDataSize = 2000000 // 2MB
	totalVariationWeight    = int32(100000)
)

var featureIDRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

func validateCreateFeatureRequest(cmd *featureproto.CreateFeatureCommand) error {
	if cmd == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	if cmd.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	if !featureIDRegex.MatchString(cmd.Id) {
		return localizedError(statusInvalidID, locale.JaJP)
	}
	if cmd.Name == "" {
		return localizedError(statusMissingName, locale.JaJP)
	}
	variationSize := len(cmd.Variations)
	if variationSize < 2 {
		return localizedError(statusMissingFeatureVariations, locale.JaJP)
	}
	if len(cmd.Tags) == 0 {
		return localizedError(statusMissingFeatureTags, locale.JaJP)
	}
	if cmd.DefaultOnVariationIndex == nil {
		return localizedError(statusMissingDefaultOnVariation, locale.JaJP)
	}
	if int(cmd.DefaultOnVariationIndex.Value) >= variationSize {
		return localizedError(statusInvalidDefaultOnVariation, locale.JaJP)
	}
	if cmd.DefaultOffVariationIndex == nil {
		return localizedError(statusMissingDefaultOffVariation, locale.JaJP)
	}
	if int(cmd.DefaultOffVariationIndex.Value) >= variationSize {
		return localizedError(statusInvalidDefaultOffVariation, locale.JaJP)
	}
	return nil
}

func validateCreateSegmentRequest(cmd *featureproto.CreateSegmentCommand) error {
	if cmd == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	if cmd.Name == "" {
		return localizedError(statusMissingName, locale.JaJP)
	}
	return nil
}

func validateGetSegmentRequest(req *featureproto.GetSegmentRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	return nil
}

func validateListSegmentsRequest(req *featureproto.ListSegmentsRequest) error {
	if req.PageSize > maxPageSizePerRequest {
		return localizedError(statusExceededMaxPageSizePerRequest, locale.JaJP)
	}
	return nil
}

func validateDeleteSegmentRequest(req *featureproto.DeleteSegmentRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	return nil
}

func validateUpdateSegment(segmentID string, commands []command.Command) error {
	if segmentID == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	return validateUpdateSegmentCommands(commands)
}

func validateUpdateSegmentCommands(commands []command.Command) error {
	for _, cmd := range commands {
		switch c := cmd.(type) {
		case *featureproto.ChangeSegmentNameCommand:
			return validateChangeSegmentName(c)
		case *featureproto.ChangeSegmentDescriptionCommand:
			return nil
		case *featureproto.AddRuleCommand:
			return validateAddSegmentRule(c)
		case *featureproto.DeleteRuleCommand:
			return validateDeleteSegmentRule(c)
		case *featureproto.AddClauseCommand:
			return validateAddSegmentClauseCommand(c)
		case *featureproto.DeleteClauseCommand:
			return validateDeleteSegmentClauseCommand(c)
		case *featureproto.ChangeClauseAttributeCommand:
			return validateChangeClauseAttributeCommand(c)
		case *featureproto.ChangeClauseOperatorCommand:
			return validateChangeClauseOperatorCommand(c)
		case *featureproto.AddClauseValueCommand:
			return validateAddClauseValueCommand(c)
		case *featureproto.RemoveClauseValueCommand:
			return validateRemoveClauseValueCommand(c)
		default:
			return localizedError(statusUnknownCommand, locale.JaJP)
		}
	}
	return localizedError(statusMissingCommand, locale.JaJP)
}

func validateChangeSegmentName(cmd *featureproto.ChangeSegmentNameCommand) error {
	if cmd.Name == "" {
		return localizedError(statusMissingName, locale.JaJP)
	}
	return nil
}

func validateAddSegmentRule(cmd *featureproto.AddRuleCommand) error {
	if cmd.Rule == nil {
		return localizedError(statusMissingRule, locale.JaJP)
	}
	if cmd.Rule.Id == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	if err := uuid.ValidateUUID(cmd.Rule.Id); err != nil {
		return localizedError(statusIncorrectUUIDFormat, locale.JaJP)
	}
	if len(cmd.Rule.Clauses) == 0 {
		return localizedError(statusMissingRuleClause, locale.JaJP)
	}
	return validateClauses(cmd.Rule.Clauses)
}

func validateClauses(clauses []*featureproto.Clause) error {
	for _, clause := range clauses {
		if clause.Attribute == "" {
			return localizedError(statusMissingClauseAttribute, locale.JaJP)
		}
		if len(clause.Values) == 0 {
			return localizedError(statusMissingClauseValues, locale.JaJP)
		}
	}
	return nil
}

func validateDeleteSegmentRule(cmd *featureproto.DeleteRuleCommand) error {
	if cmd == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	if cmd.Id == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	return nil
}

func validateAddSegmentClauseCommand(cmd *featureproto.AddClauseCommand) error {
	if cmd.RuleId == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	if cmd.Clause == nil {
		return localizedError(statusMissingRuleClause, locale.JaJP)
	}
	return validateClauses([]*featureproto.Clause{cmd.Clause})
}

func validateDeleteSegmentClauseCommand(cmd *featureproto.DeleteClauseCommand) error {
	if cmd.Id == "" {
		return localizedError(statusMissingClauseID, locale.JaJP)
	}
	if cmd.RuleId == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	return nil
}

func validateChangeClauseAttributeCommand(cmd *featureproto.ChangeClauseAttributeCommand) error {
	if cmd.Id == "" {
		return localizedError(statusMissingClauseID, locale.JaJP)
	}
	if cmd.RuleId == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	if cmd.Attribute == "" {
		return localizedError(statusMissingClauseAttribute, locale.JaJP)
	}
	return nil
}

func validateChangeClauseOperatorCommand(cmd *featureproto.ChangeClauseOperatorCommand) error {
	if cmd.Id == "" {
		return localizedError(statusMissingClauseID, locale.JaJP)
	}
	if cmd.RuleId == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	return nil
}

func validateAddClauseValueCommand(cmd *featureproto.AddClauseValueCommand) error {
	return validateClauseValueCommand(cmd.Id, cmd.RuleId, cmd.Value)
}

func validateRemoveClauseValueCommand(cmd *featureproto.RemoveClauseValueCommand) error {
	return validateClauseValueCommand(cmd.Id, cmd.RuleId, cmd.Value)
}

func validateClauseValueCommand(clauseID string, ruleID string, value string) error {
	if clauseID == "" {
		return localizedError(statusMissingClauseID, locale.JaJP)
	}
	if ruleID == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	if value == "" {
		return localizedError(statusMissingClauseValue, locale.JaJP)
	}
	return nil
}

func validateAddSegmentUserRequest(req *featureproto.AddSegmentUserRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	return validateSegmentUserState(req.Command.State)
}

func validateAddSegmentUserCommand(cmd *featureproto.AddSegmentUserCommand) error {
	return validateUserIDs(cmd.UserIds)
}

func validateDeleteSegmentUserRequest(req *featureproto.DeleteSegmentUserRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	return validateSegmentUserState(req.Command.State)
}

func validateSegmentUserState(state featureproto.SegmentUser_State) error {
	switch state {
	case featureproto.SegmentUser_INCLUDED:
		return nil
	default:
		return localizedError(statusUnknownSegmentUserState, locale.JaJP)
	}
}

func validateDeleteSegmentUserCommand(cmd *featureproto.DeleteSegmentUserCommand) error {
	return validateUserIDs(cmd.UserIds)
}

func validateUserIDs(userIDs []string) error {
	size := len(userIDs)
	if size == 0 {
		return localizedError(statusMissingUserIDs, locale.JaJP)
	}
	if size > maxUserIDsLength {
		return localizedError(statusExceededMaxUserIDsLength, locale.JaJP)
	}
	for _, id := range userIDs {
		if id == "" {
			return localizedError(statusMissingUserID, locale.JaJP)
		}
	}
	return nil
}

func validateGetSegmentUserRequest(req *featureproto.GetSegmentUserRequest) error {
	if req.SegmentId == "" {
		return localizedError(statusMissingSegmentID, locale.JaJP)
	}
	if req.UserId == "" {
		return localizedError(statusMissingUserID, locale.JaJP)
	}
	return nil
}

func validateListSegmentUsersRequest(req *featureproto.ListSegmentUsersRequest) error {
	if req.SegmentId == "" {
		return localizedError(statusMissingSegmentID, locale.JaJP)
	}
	if req.PageSize > maxPageSizePerRequest {
		return localizedError(statusExceededMaxPageSizePerRequest, locale.JaJP)
	}
	return nil
}

func validateBulkUploadSegmentUsersRequest(req *featureproto.BulkUploadSegmentUsersRequest) error {
	if req.SegmentId == "" {
		return localizedError(statusMissingSegmentID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	return nil
}

func validateBulkUploadSegmentUsersCommand(cmd *featureproto.BulkUploadSegmentUsersCommand) error {
	if len(cmd.Data) == 0 {
		return localizedError(statusMissingSegmentUsersData, locale.JaJP)
	}
	if len(cmd.Data) > maxSegmentUsersDataSize {
		return localizedError(statusExceededMaxSegmentUsersDataSize, locale.JaJP)
	}
	return validateSegmentUserState(cmd.State)
}

func validateBulkDownloadSegmentUsersRequest(req *featureproto.BulkDownloadSegmentUsersRequest) error {
	if req.SegmentId == "" {
		return localizedError(statusMissingSegmentID, locale.JaJP)
	}
	return validateSegmentUserState(req.State)
}

func validateEvaluateFeatures(req *featureproto.EvaluateFeaturesRequest) error {
	if req.User == nil {
		return localizedError(statusMissingUser, locale.JaJP)
	}
	if req.User.Id == "" {
		return localizedError(statusMissingUserID, locale.JaJP)
	}
	if req.Tag == "" {
		return localizedError(statusMissingFeatureTag, locale.JaJP)
	}
	return nil
}

func validateUpsertUserEvaluationRequest(req *featureproto.UpsertUserEvaluationRequest) error {
	if req.Tag == "" {
		return localizedError(statusMissingFeatureTag, locale.JaJP)
	}
	if req.Evaluation == nil {
		return localizedError(statusMissingEvaluation, locale.JaJP)
	}
	return nil
}

func validateGetUserEvaluationsRequest(req *featureproto.GetUserEvaluationsRequest) error {
	if req.Tag == "" {
		return localizedError(statusMissingFeatureTag, locale.JaJP)
	}
	if req.UserId == "" {
		return localizedError(statusMissingUserID, locale.JaJP)
	}
	return nil
}

func validateGetFeatureRequest(req *featureproto.GetFeatureRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	return nil
}

func validateGetFeaturesRequest(req *featureproto.GetFeaturesRequest) error {
	if len(req.Ids) == 0 {
		return localizedError(statusMissingIDs, locale.JaJP)
	}
	for _, id := range req.Ids {
		if id == "" {
			return localizedError(statusMissingIDs, locale.JaJP)
		}
	}
	return nil
}

func validateEnableFeatureRequest(req *featureproto.EnableFeatureRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	return nil
}

func validateDisableFeatureRequest(req *featureproto.DisableFeatureRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	return nil
}

func validateDeleteFeatureRequest(req *featureproto.DeleteFeatureRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	return nil
}

func validateFeatureVariationsCommand(
	fs []*featureproto.Feature,
	cmd command.Command,
) error {
	switch c := cmd.(type) {
	case *featureproto.RemoveVariationCommand:
		return validateVariationCommand(fs, c.Id)
	case *featureproto.ChangeVariationValueCommand:
		return validateVariationCommand(fs, c.Id)
	default:
		return nil
	}
}

func validateArchiveFeatureRequest(req *featureproto.ArchiveFeatureRequest, fs []*featureproto.Feature) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	for _, f := range fs {
		for _, p := range f.Prerequisites {
			if p.FeatureId == req.Id {
				return localizedError(statusInvalidArchive, locale.JaJP)
			}
		}
	}
	return nil
}

func validateVariationCommand(fs []*featureproto.Feature, vID string) error {
	for _, f := range fs {
		for _, p := range f.Prerequisites {
			if p.VariationId == vID {
				return localizedError(statusInvalidChangingVariation, locale.JaJP)
			}
		}
	}
	return nil
}

func validateUnarchiveFeatureRequest(req *featureproto.UnarchiveFeatureRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	return nil
}

func validateCloneFeatureRequest(req *featureproto.CloneFeatureRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusMissingCommand, locale.JaJP)
	}
	if req.Command.EnvironmentNamespace == req.EnvironmentNamespace {
		return localizedError(statusIncorrectDestinationEnvironment, locale.JaJP)
	}
	return nil
}

func validateFeatureTargetingCommand(
	fs []*featureproto.Feature,
	tarF *featureproto.Feature,
	cmd command.Command,
	localizer locale.Localizer,
) error {
	switch c := cmd.(type) {
	case *featureproto.AddRuleCommand:
		return validateRule(tarF.Variations, c.Rule)
	case *featureproto.ChangeRuleStrategyCommand:
		return validateChangeRuleStrategy(tarF.Variations, c)
	case *featureproto.ChangeDefaultStrategyCommand:
		return validateChangeDefaultStrategy(tarF.Variations, c)
	case *featureproto.ChangeFixedStrategyCommand:
		return validateChangeFixedStrategy(c)
	case *featureproto.ChangeRolloutStrategyCommand:
		return validateChangeRolloutStrategy(tarF.Variations, c)
	case *featureproto.AddPrerequisiteCommand:
		return validateAddPrerequisite(fs, tarF, c.Prerequisite, localizer)
	case *featureproto.ChangePrerequisiteVariationCommand:
		return validateChangePrerequisiteVariation(fs, c.Prerequisite, localizer)
	default:
		return nil
	}
}

func validateRule(variations []*featureproto.Variation, rule *featureproto.Rule) error {
	if rule.Id == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	if err := uuid.ValidateUUID(rule.Id); err != nil {
		return localizedError(statusIncorrectUUIDFormat, locale.JaJP)
	}
	return validateStrategy(variations, rule.Strategy)
}

func validateChangeRuleStrategy(variations []*featureproto.Variation, cmd *featureproto.ChangeRuleStrategyCommand,
) error {
	if cmd.RuleId == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	return validateStrategy(variations, cmd.Strategy)
}

func validateChangeDefaultStrategy(
	variations []*featureproto.Variation,
	cmd *featureproto.ChangeDefaultStrategyCommand,
) error {
	if cmd.Strategy == nil {
		return localizedError(statusMissingRuleStrategy, locale.JaJP)
	}
	return validateStrategy(variations, cmd.Strategy)
}

func validateStrategy(variations []*featureproto.Variation, strategy *featureproto.Strategy) error {
	if strategy == nil {
		return localizedError(statusMissingRuleStrategy, locale.JaJP)
	}
	if strategy.Type == featureproto.Strategy_FIXED {
		return validateFixedStrategy(strategy.FixedStrategy)
	}
	if strategy.Type == featureproto.Strategy_ROLLOUT {
		return validateRolloutStrategy(variations, strategy.RolloutStrategy)
	}
	return localizedError(statusUnknownStrategy, locale.JaJP)
}

func validateChangeFixedStrategy(cmd *featureproto.ChangeFixedStrategyCommand) error {
	if cmd.RuleId == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	return validateFixedStrategy(cmd.Strategy)
}

func validateChangeRolloutStrategy(
	variations []*featureproto.Variation,
	cmd *featureproto.ChangeRolloutStrategyCommand,
) error {
	if cmd.RuleId == "" {
		return localizedError(statusMissingRuleID, locale.JaJP)
	}
	return validateRolloutStrategy(variations, cmd.Strategy)
}

func validateFixedStrategy(strategy *featureproto.FixedStrategy) error {
	if strategy == nil {
		return localizedError(statusMissingFixedStrategy, locale.JaJP)
	}
	if strategy.Variation == "" {
		return localizedError(statusMissingVariationID, locale.JaJP)
	}
	return nil
}

func validateRolloutStrategy(variations []*featureproto.Variation, strategy *featureproto.RolloutStrategy) error {
	if strategy == nil {
		return localizedError(statusMissingRolloutStrategy, locale.JaJP)
	}
	if len(variations) != len(strategy.Variations) {
		return localizedError(statusDifferentVariationsSize, locale.JaJP)
	}
	sum := int32(0)
	for _, v := range strategy.Variations {
		if v.Variation == "" {
			return localizedError(statusMissingVariationID, locale.JaJP)
		}
		if v.Weight < 0 {
			return localizedError(statusIncorrectVariationWeight, locale.JaJP)
		}
		sum += v.Weight
	}
	if sum != totalVariationWeight {
		return localizedError(statusExceededMaxVariationWeight, locale.JaJP)
	}
	return nil
}

func validateAddPrerequisite(
	fs []*featureproto.Feature,
	tarF *featureproto.Feature,
	p *featureproto.Prerequisite,
	localizer locale.Localizer,
) error {
	if tarF.Id == p.FeatureId {
		return localizedError(statusInvalidPrerequisite, locale.JaJP)
	}
	for _, pf := range tarF.Prerequisites {
		if pf.FeatureId == p.FeatureId {
			return localizedError(statusInvalidPrerequisite, locale.JaJP)
		}
	}
	if err := validateVariationID(fs, p, localizer); err != nil {
		return err
	}
	tarF.Prerequisites = append(tarF.Prerequisites, p)
	_, err := domain.TopologicalSort(fs)
	if err != nil {
		if err == domain.ErrCycleExists {
			return localizedError(statusCycleExists, locale.JaJP)
		}
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateChangePrerequisiteVariation(
	fs []*featureproto.Feature,
	p *featureproto.Prerequisite,
	localizer locale.Localizer,
) error {
	if err := validateVariationID(fs, p, localizer); err != nil {
		return err
	}
	return nil
}

func validateVariationID(fs []*featureproto.Feature, p *featureproto.Prerequisite, localizer locale.Localizer) error {
	f, err := findFeature(fs, p.FeatureId, localizer)
	if err != nil {
		return err
	}
	for _, v := range f.Variations {
		if v.Id == p.VariationId {
			return nil
		}
	}
	return localizedError(statusInvalidVariationID, locale.JaJP)
}
