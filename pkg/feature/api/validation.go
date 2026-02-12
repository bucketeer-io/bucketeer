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
	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	maxPageSizePerRequest   = 100
	maxSegmentUsersDataSize = 2000000 // 2MB
	// Scheduled flag change limits
	maxSchedulesPerFlag    = 50
	maxChangesPerSchedule  = 50
	minScheduleTimeMinutes = 5   // 5 minutes
	maxScheduleTimeDays    = 365 // 1 year
)

var featureIDRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

func validateCreateFeatureRequest(
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

func validateDeleteFeatureRequest(req *featureproto.DeleteFeatureRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	return nil
}

func validateCloneFeatureRequest(req *featureproto.CloneFeatureRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.TargetEnvironmentId == req.EnvironmentId {
		return statusIncorrectDestinationEnvironment.Err()
	}
	return nil
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

func validateCreateFlagTriggerRequest(req *featureproto.CreateFlagTriggerRequest) error {
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
