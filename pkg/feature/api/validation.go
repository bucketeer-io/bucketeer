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

package api

import (
	"context"
	"errors"
	"regexp"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/command"
	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	maxPageSizePerRequest   = 100
	maxSegmentUsersDataSize = 2000000 // 2MB
	totalVariationWeight    = int32(100000)
)

var featureIDRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

func validateCreateFeatureRequest(cmd *featureproto.CreateFeatureCommand) error {
	if cmd.Id == "" {
		return statusMissingID.Err()
	}
	if !featureIDRegex.MatchString(cmd.Id) {
		return statusInvalidID.Err()
	}
	if cmd.Name == "" {
		return statusMissingName.Err()
	}
	variationSize := len(cmd.Variations)
	if variationSize < 2 {
		return statusMissingFeatureVariations.Err()
	}
	if len(cmd.Tags) == 0 {
		return statusMissingFeatureTags.Err()
	}
	if cmd.DefaultOnVariationIndex == nil {
		return statusMissingDefaultOnVariation.Err()
	}
	if int(cmd.DefaultOnVariationIndex.Value) >= variationSize {
		return statusInvalidDefaultOnVariation.Err()
	}
	if cmd.DefaultOffVariationIndex == nil {
		return statusMissingDefaultOffVariation.Err()
	}
	if int(cmd.DefaultOffVariationIndex.Value) >= variationSize {
		return statusInvalidDefaultOffVariation.Err()
	}
	return nil
}

func validateCreateFeatureRequestNoCommand(
	req *featureproto.CreateFeatureRequest,
) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if !featureIDRegex.MatchString(req.Id) {
		return statusInvalidID.Err()
	}
	if req.Name == "" {
		return statusMissingName.Err()
	}
	variationSize := len(req.Variations)
	if variationSize < 2 {
		return statusMissingFeatureVariations.Err()
	}
	if len(req.Tags) == 0 {
		return statusMissingFeatureTags.Err()
	}
	if req.DefaultOnVariationIndex == nil {
		return statusMissingDefaultOnVariation.Err()
	}
	if int(req.DefaultOnVariationIndex.Value) >= variationSize {
		return statusInvalidDefaultOnVariation.Err()
	}
	if req.DefaultOffVariationIndex == nil {
		return statusMissingDefaultOffVariation.Err()
	}
	if int(req.DefaultOffVariationIndex.Value) >= variationSize {
		return statusInvalidDefaultOffVariation.Err()
	}
	return nil
}

func validateCreateSegmentRequest(
	req *featureproto.CreateSegmentRequest,
) error {
	if req.Name == "" {
		return statusMissingName.Err()
	}
	return nil
}

func validateGetSegmentRequest(req *featureproto.GetSegmentRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	return nil
}

func validateListSegmentsRequest(req *featureproto.ListSegmentsRequest) error {
	if req.PageSize > maxPageSizePerRequest {
		return statusExceededMaxPageSizePerRequest.Err()
	}
	return nil
}

func validateDeleteSegmentRequest(req *featureproto.DeleteSegmentRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	return nil
}

func validateUpdateSegmentRequest(
	req *featureproto.UpdateSegmentRequest,
) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.Name != nil && req.Name.Value == "" {
		return statusMissingName.Err()
	}
	return nil
}

func validateSegmentUserState(state featureproto.SegmentUser_State) error {
	switch state {
	case featureproto.SegmentUser_INCLUDED:
		return nil
	default:
		return statusUnknownSegmentUserState.Err()
	}
}

func validateListSegmentUsersRequest(req *featureproto.ListSegmentUsersRequest) error {
	if req.SegmentId == "" {
		return statusMissingSegmentID.Err()
	}
	if req.PageSize > maxPageSizePerRequest {
		return statusExceededMaxPageSizePerRequest.Err()
	}
	return nil
}

func validateBulkUploadSegmentUsersRequest(
	req *featureproto.BulkUploadSegmentUsersRequest,
) error {
	if req.SegmentId == "" {
		return statusMissingSegmentID.Err()
	}
	if len(req.Data) == 0 {
		return statusMissingSegmentUsersData.Err()
	}
	if len(req.Data) > maxSegmentUsersDataSize {
		return statusExceededMaxSegmentUsersDataSize.Err()
	}
	return validateSegmentUserState(req.State)
}

func validateBulkDownloadSegmentUsersRequest(
	req *featureproto.BulkDownloadSegmentUsersRequest,
) error {
	if req.SegmentId == "" {
		return statusMissingSegmentID.Err()
	}
	return validateSegmentUserState(req.State)
}

func validateEvaluateFeatures(req *featureproto.EvaluateFeaturesRequest) error {
	if req.User == nil {
		return statusMissingUser.Err()
	}
	if req.User.Id == "" {
		return statusMissingUserID.Err()
	}
	return nil
}

func validateDebugEvaluateFeatures(req *featureproto.DebugEvaluateFeaturesRequest) error {
	if len(req.Users) == 0 {
		return statusMissingUser.Err()
	}

	for _, user := range req.Users {
		if user.Id == "" {
			return statusMissingUserID.Err()
		}
	}

	if len(req.FeatureIds) == 0 {
		return statusMissingFeatureIDs.Err()
	}

	return nil
}

func validateGetFeatureRequest(req *featureproto.GetFeatureRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	return nil
}

func validateGetFeaturesRequest(req *featureproto.GetFeaturesRequest) error {
	if len(req.Ids) == 0 {
		return statusMissingIDs.Err()
	}
	for _, id := range req.Ids {
		if id == "" {
			return statusMissingIDs.Err()
		}
	}
	return nil
}

func validateEnableFeatureRequest(req *featureproto.EnableFeatureRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.Command == nil {
		return statusMissingCommand.Err()
	}
	return nil
}

func validateDisableFeatureRequest(req *featureproto.DisableFeatureRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.Command == nil {
		return statusMissingCommand.Err()
	}
	return nil
}

func validateDeleteFeatureRequest(req *featureproto.DeleteFeatureRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.Command == nil {
		return statusMissingCommand.Err()
	}
	return nil
}

// We can't add or remove a variation when a progressive rollout is in progress, but changes can be made.
func (s *FeatureService) validateFeatureVariationsCommand(
	ctx context.Context,
	fs []*featureproto.Feature,
	environmentID string,
	f *featureproto.Feature,
	cmd command.Command,
) error {
	switch c := cmd.(type) {
	case *featureproto.AddVariationCommand:
		if err := s.checkProgressiveRolloutInProgress(ctx, environmentID, f.Id); err != nil {
			return err
		}
		return nil
	case *featureproto.RemoveVariationCommand:
		if err := s.checkProgressiveRolloutInProgress(ctx, environmentID, f.Id); err != nil {
			return err
		}
		return validateRemoveVariationCommand(c, fs, f)
	case *featureproto.ChangeVariationValueCommand:
		return validateVariationCommand(fs, f)
	default:
		return nil
	}
}

func (s *FeatureService) checkProgressiveRolloutInProgress(
	ctx context.Context,
	environmentID, featureID string,
) error {
	exists, err := s.existsRunningProgressiveRollout(ctx, featureID, environmentID)
	if err != nil {
		return api.NewGRPCStatus(err).Err()
	}
	if exists {
		return statusProgressiveRolloutWaitingOrRunningState.Err()
	}
	return nil
}

func validateArchiveFeatureRequest(
	req *featureproto.ArchiveFeatureRequest,
) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.Command == nil {
		return statusMissingCommand.Err()
	}
	return nil
}

func validateVariationCommand(fs []*featureproto.Feature, tgt *featureproto.Feature) error {
	if featuredomain.HasFeaturesDependsOnTargets([]*featureproto.Feature{tgt}, fs) {
		return statusInvalidChangingVariation.Err()
	}
	return nil
}

// validateRemoveVariationCommand validates that a specific variation can be safely removed
func validateRemoveVariationCommand(
	cmd *featureproto.RemoveVariationCommand,
	fs []*featureproto.Feature,
	tgt *featureproto.Feature,
) error {
	// Find the variation being removed to get its value
	var deletedVariationValue string
	for _, variation := range tgt.Variations {
		if variation.Id == cmd.Id {
			deletedVariationValue = variation.Value
			break
		}
	}

	// Even if we can't find the variation value, we should still check for prerequisites
	// since they reference variation IDs, not values

	// Optimization: First check if ANY features depend on our target
	// This reuses existing logic to quickly filter relevant features
	allFeaturesMap := make(map[string]*featureproto.Feature, len(fs))
	for _, f := range fs {
		allFeaturesMap[f.Id] = f
	}

	dependentFeatures := featuredomain.GetFeaturesDependsOnTargets([]*featureproto.Feature{tgt}, allFeaturesMap)
	delete(dependentFeatures, tgt.Id) // Remove the target itself

	if len(dependentFeatures) == 0 {
		// No features depend on our target, so variation deletion is safe
		return nil
	}

	// Convert dependent features back to slice for ValidateVariationUsage
	dependentFeaturesSlice := make([]*featureproto.Feature, 0, len(dependentFeatures))
	for _, f := range dependentFeatures {
		dependentFeaturesSlice = append(dependentFeaturesSlice, f)
	}

	// Use our precise cross-feature validation only on dependent features
	deletedVariations := map[string]string{}
	if deletedVariationValue != "" {
		deletedVariations[cmd.Id] = deletedVariationValue
	} else {
		// We don't have the variation value, but we still need to check prerequisites
		// For prerequisites, we only need the variation ID (key), not the value
		deletedVariations[cmd.Id] = "" // Empty value, but we'll check keys for prerequisites
	}

	if len(deletedVariations) == 0 {
		return nil // No variations being deleted
	}

	if err := featuredomain.ValidateVariationUsage(dependentFeaturesSlice, tgt.Id, deletedVariations); err != nil {
		if errors.Is(err, featuredomain.ErrVariationInUse) {
			// Use the legacy error status for RemoveVariationCommand for backward compatibility
			return statusInvalidChangingVariation.Err()
		}
		return err
	}

	return nil
}

func validateUnarchiveFeatureRequest(req *featureproto.UnarchiveFeatureRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.Command == nil {
		return statusMissingCommand.Err()
	}
	return nil
}

func validateCloneFeatureRequest(req *featureproto.CloneFeatureRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.Command.EnvironmentId == req.EnvironmentId {
		return statusIncorrectDestinationEnvironment.Err()
	}
	return nil
}

func validateCloneFeatureRequestNoCommand(req *featureproto.CloneFeatureRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.TargetEnvironmentId == req.EnvironmentId {
		return statusIncorrectDestinationEnvironment.Err()
	}
	return nil
}

func validateUpdateFeatureTargetingRequest(
	req *featureproto.UpdateFeatureTargetingRequest,
) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.From == featureproto.UpdateFeatureTargetingRequest_UNKNOWN {
		return statusMissingFrom.Err()
	}
	return nil
}

func (s *FeatureService) validateFeatureTargetingCommand(
	ctx context.Context,
	from featureproto.UpdateFeatureTargetingRequest_From,
	environmentID string,
	fs []*featureproto.Feature,
	tarF *featureproto.Feature,
	cmd command.Command,
) error {
	switch c := cmd.(type) {
	case *featureproto.AddRuleCommand:
		return validateRule(fs, tarF, c.Rule)
	case *featureproto.ChangeRuleStrategyCommand:
		return validateChangeRuleStrategy(tarF.Variations, c)
	case *featureproto.ChangeDefaultStrategyCommand:
		return s.validateChangeDefaultStrategy(ctx, from, environmentID, tarF.Id, tarF.Variations, c)
	case *featureproto.ChangeFixedStrategyCommand:
		return validateChangeFixedStrategy(c)
	case *featureproto.ChangeRolloutStrategyCommand:
		return validateChangeRolloutStrategy(tarF.Variations, c)
	case *featureproto.AddPrerequisiteCommand:
		return validateAddPrerequisite(fs, tarF, c.Prerequisite)
	case *featureproto.ChangePrerequisiteVariationCommand:
		return validateChangePrerequisiteVariation(fs, c.Prerequisite)
	default:
		return nil
	}
}

func validateRule(
	fs []*featureproto.Feature,
	tarF *featureproto.Feature,
	rule *featureproto.Rule) error {
	if rule.Id == "" {
		return statusMissingRuleID.Err()
	}
	if err := uuid.ValidateUUID(rule.Id); err != nil {
		return statusIncorrectUUIDFormat.Err()
	}
	// Check dependency.
	tarF.Rules = append(tarF.Rules, rule)
	defer func() { tarF.Rules = tarF.Rules[:len(tarF.Rules)-1] }()
	if err := featuredomain.ValidateFeatureDependencies(fs); err != nil {
		if errors.Is(err, featuredomain.ErrCycleExists) {
			return statusCycleExists.Err()
		}
		return api.NewGRPCStatus(err).Err()
	}
	return validateStrategy(tarF.Variations, rule.Strategy)
}

func validateChangeRuleStrategy(
	variations []*featureproto.Variation,
	cmd *featureproto.ChangeRuleStrategyCommand,
) error {
	if cmd.RuleId == "" {
		return statusMissingRuleID.Err()
	}
	return validateStrategy(variations, cmd.Strategy)
}

// We can't change the default strategy when there is a progressive rollout in progress.
// Otherwise, it could conflict with the rollout rules.
func (s *FeatureService) validateChangeDefaultStrategy(
	ctx context.Context,
	from featureproto.UpdateFeatureTargetingRequest_From,
	environmentID, featureID string,
	variations []*featureproto.Variation,
	cmd *featureproto.ChangeDefaultStrategyCommand,
) error {
	// Because the progressive rollout changes the default strategy,
	// We must check from where the request comes
	if from == featureproto.UpdateFeatureTargetingRequest_USER {
		if err := s.checkProgressiveRolloutInProgress(ctx, environmentID, featureID); err != nil {
			return err
		}
	}
	if cmd.Strategy == nil {
		return statusMissingRuleStrategy.Err()
	}
	return validateStrategy(variations, cmd.Strategy)
}

func validateStrategy(
	variations []*featureproto.Variation,
	strategy *featureproto.Strategy,
) error {
	if strategy == nil {
		return statusMissingRuleStrategy.Err()
	}
	if strategy.Type == featureproto.Strategy_FIXED {
		return validateFixedStrategy(strategy.FixedStrategy)
	}
	if strategy.Type == featureproto.Strategy_ROLLOUT {
		return validateRolloutStrategy(variations, strategy.RolloutStrategy)
	}
	return statusUnknownStrategy.Err()
}

func validateChangeFixedStrategy(cmd *featureproto.ChangeFixedStrategyCommand) error {
	if cmd.RuleId == "" {
		return statusMissingRuleID.Err()
	}
	return validateFixedStrategy(cmd.Strategy)
}

func validateChangeRolloutStrategy(
	variations []*featureproto.Variation,
	cmd *featureproto.ChangeRolloutStrategyCommand,
) error {
	if cmd.RuleId == "" {
		return statusMissingRuleID.Err()
	}
	return validateRolloutStrategy(variations, cmd.Strategy)
}

func validateFixedStrategy(strategy *featureproto.FixedStrategy) error {
	if strategy == nil {
		return statusMissingFixedStrategy.Err()
	}
	if strategy.Variation == "" {
		return statusMissingVariationID.Err()
	}
	return nil
}

func validateRolloutStrategy(
	variations []*featureproto.Variation,
	strategy *featureproto.RolloutStrategy,
) error {
	if strategy == nil {
		return statusMissingRolloutStrategy.Err()
	}
	if len(variations) != len(strategy.Variations) {
		return statusDifferentVariationsSize.Err()
	}
	sum := int32(0)
	for _, v := range strategy.Variations {
		if v.Variation == "" {
			return statusMissingVariationID.Err()
		}
		if v.Weight < 0 {
			return statusIncorrectVariationWeight.Err()
		}
		sum += v.Weight
	}
	if sum != totalVariationWeight {
		return statusExceededMaxVariationWeight.Err()
	}
	return nil
}

func validateAddPrerequisite(
	fs []*featureproto.Feature,
	tarF *featureproto.Feature,
	p *featureproto.Prerequisite,
) error {
	if tarF.Id == p.FeatureId {
		return statusInvalidPrerequisite.Err()
	}
	for _, pf := range tarF.Prerequisites {
		if pf.FeatureId == p.FeatureId {
			return statusInvalidPrerequisite.Err()
		}
	}
	if err := validateVariationID(fs, p); err != nil {
		return err
	}
	prevPrerequisites := tarF.Prerequisites
	tarF.Prerequisites = append(tarF.Prerequisites, p)
	defer func() { tarF.Prerequisites = prevPrerequisites }()
	if err := featuredomain.ValidateFeatureDependencies(fs); err != nil {
		if err == featuredomain.ErrCycleExists {
			return statusCycleExists.Err()
		}
		return api.NewGRPCStatus(err).Err()
	}
	return nil
}

func validateChangePrerequisiteVariation(
	fs []*featureproto.Feature,
	p *featureproto.Prerequisite,
) error {
	if err := validateVariationID(fs, p); err != nil {
		return err
	}
	return nil
}

func validateVariationID(fs []*featureproto.Feature, p *featureproto.Prerequisite) error {
	f, err := findFeature(fs, p.FeatureId)
	if err != nil {
		return err
	}
	for _, v := range f.Variations {
		if v.Id == p.VariationId {
			return nil
		}
	}
	return statusInvalidVariationID.Err()
}

func (s *FeatureService) validateFeatureStatus(
	ctx context.Context,
	id, environmentId string,
) error {
	runningExperimentExists, err := s.existsRunningExperiment(ctx, id, environmentId)
	if err != nil {
		return api.NewGRPCStatus(err).Err()
	}
	if runningExperimentExists {
		return statusWaitingOrRunningExperimentExists.Err()
	}
	return nil
}

func (s *FeatureService) validateEnvironmentSettings(
	ctx context.Context,
	environmentId, updateComment string,
) error {
	req := &envproto.GetEnvironmentV2Request{
		Id: environmentId,
	}
	resp, err := s.environmentClient.GetEnvironmentV2(ctx, req)
	if err != nil {
		return api.NewGRPCStatus(err).Err()
	}
	if resp.Environment.RequireComment {
		if updateComment == "" {
			return statusCommentRequiredForUpdating.Err()
		}
	}
	return nil
}

func validateCreateFlagTriggerCommand(cmd *featureproto.CreateFlagTriggerCommand) error {
	if cmd.FeatureId == "" {
		return statusMissingTriggerFeatureID.Err()
	}
	if cmd.Type == featureproto.FlagTrigger_Type_UNKNOWN {
		return statusMissingTriggerType.Err()
	}
	if cmd.Action == featureproto.FlagTrigger_Action_UNKNOWN {
		return statusMissingTriggerAction.Err()
	}
	return nil
}

func validateCreateFlagTriggerNoCommand(req *featureproto.CreateFlagTriggerRequest) error {
	if req.FeatureId == "" {
		return statusMissingTriggerFeatureID.Err()
	}
	if req.Type == featureproto.FlagTrigger_Type_UNKNOWN {
		return statusMissingTriggerType.Err()
	}
	if req.Action == featureproto.FlagTrigger_Action_UNKNOWN {
		return statusMissingTriggerAction.Err()
	}
	return nil
}

func validateEnableFlagTriggerCommand(cmd *featureproto.EnableFlagTriggerCommand) error {
	if cmd == nil {
		return statusMissingCommand.Err()
	}
	return nil
}

func validateDisableFlagTriggerCommand(cmd *featureproto.DisableFlagTriggerCommand) error {
	if cmd == nil {
		return statusMissingCommand.Err()
	}
	return nil
}

func validateResetFlagTriggerCommand(cmd *featureproto.ResetFlagTriggerCommand) error {
	if cmd == nil {
		return statusMissingCommand.Err()
	}
	return nil
}

func validateGetFlagTriggerRequest(req *featureproto.GetFlagTriggerRequest) error {
	if req.Id == "" {
		return statusMissingTriggerID.Err()
	}
	return nil
}

func validateListFlagTriggersRequest(req *featureproto.ListFlagTriggersRequest) error {
	if req.FeatureId == "" {
		return statusMissingTriggerFeatureID.Err()
	}
	return nil
}

// validateVariationDeletion validates that variations being deleted are not used as cross-feature dependencies.
func validateVariationDeletion(
	variationChanges []*featureproto.VariationChange,
	features []*featureproto.Feature,
	targetFeatureID string,
) error {
	// Extract variations being deleted (variationID -> variationValue)
	deletedVariations := make(map[string]string)
	for _, change := range variationChanges {
		if change.ChangeType == featureproto.ChangeType_DELETE {
			deletedVariations[change.Variation.Id] = change.Variation.Value
		}
	}

	if len(deletedVariations) == 0 {
		return nil // No variations being deleted
	}

	// Find the target feature
	var targetFeature *featureproto.Feature
	for _, f := range features {
		if f.Id == targetFeatureID {
			targetFeature = f
			break
		}
	}

	if targetFeature == nil {
		// Target feature not found in the list, cannot proceed with validation
		return nil
	}

	// Optimization: First check if ANY features depend on our target
	// This reuses existing logic to quickly filter relevant features
	allFeaturesMap := make(map[string]*featureproto.Feature, len(features))
	for _, f := range features {
		allFeaturesMap[f.Id] = f
	}

	dependentFeatures := featuredomain.GetFeaturesDependsOnTargets(
		[]*featureproto.Feature{targetFeature},
		allFeaturesMap,
	)
	delete(dependentFeatures, targetFeatureID) // Remove the target itself

	if len(dependentFeatures) == 0 {
		// No features depend on our target, so variation deletion is safe
		return nil
	}

	// Convert dependent features back to slice for ValidateVariationUsage
	dependentFeaturesSlice := make([]*featureproto.Feature, 0, len(dependentFeatures))
	for _, f := range dependentFeatures {
		dependentFeaturesSlice = append(dependentFeaturesSlice, f)
	}

	// Check if the deleted variation is used as a prerequisite or rule in other features
	if err := featuredomain.ValidateVariationUsage(
		dependentFeaturesSlice,
		targetFeatureID,
		deletedVariations,
	); err != nil {
		if errors.Is(err, featuredomain.ErrVariationInUse) {
			return statusVariationInUseByOtherFeatures.Err()
		}
		return err
	}

	return nil
}
