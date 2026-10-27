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
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	holdoutdomain "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	v2hs "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func (s *FeatureService) CreateHoldoutGroup(
	ctx context.Context,
	req *featureproto.CreateHoldoutGroupRequest,
) (*featureproto.CreateHoldoutGroupResponse, error) {
	editor, err := s.checkHoldoutAdmin(ctx, req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := requireEnvironment(req.EnvironmentId); err != nil {
		return nil, holdoutStatus(err)
	}
	if err := holdoutdomain.RejectClientAssignment(req.BucketStart, req.SamplingSeed); err != nil {
		return nil, holdoutStatus(err)
	}
	if err := s.validateGoalIDs(ctx, req.EnvironmentId, req.GoalIds); err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	var group *holdoutdomain.HoldoutGroup
	var event *eventproto.Event
	err = s.dbClient.RunInTransactionV2(ctx, func(txCtx context.Context) error {
		if err := s.holdoutStorage.LockEnvironment(txCtx, req.EnvironmentId); err != nil {
			return err
		}
		active, err := s.holdoutStorage.ListActiveGroups(txCtx, req.EnvironmentId)
		if err != nil {
			return err
		}
		if conflict, ok := holdoutdomain.OverlappingGroup(req.StartAt, req.StopAt, "", active); ok {
			return holdoutdomain.NewPeriodOverlapError(conflict)
		}
		if req.AssignmentType == featureproto.HoldoutGroup_SEGMENT {
			if err := s.ensureSegmentExists(txCtx, req.SegmentId, req.EnvironmentId); err != nil {
				return err
			}
			if serving, ok := holdoutdomain.ServingGroupForSegment(req.SegmentId, active); ok {
				return holdoutdomain.NewSegmentLockedError(req.SegmentId, serving)
			}
		}
		var seed string
		if req.AssignmentType == featureproto.HoldoutGroup_HASH {
			seed, err = s.holdoutStorage.GetEnvironmentSamplingSeed(txCtx, req.EnvironmentId)
			if err != nil {
				return err
			}
		}
		group, err = holdoutdomain.NewHoldoutGroup(holdoutdomain.CreateParams{
			EnvironmentID:   req.EnvironmentId,
			Name:            req.Name,
			Description:     req.Description,
			AssignmentType:  req.AssignmentType,
			Percentage:      req.Percentage,
			SegmentID:       req.SegmentId,
			StartAt:         req.StartAt,
			StopAt:          req.StopAt,
			StartAnalysisAt: req.StartAnalysisAt,
			GoalIDs:         req.GoalIds,
			CreatedBy:       editor.Email,
			Now:             now,
			ExistingSeed:    seed,
		})
		if err != nil {
			return err
		}
		event, err = newHoldoutEvent(
			editor,
			req.EnvironmentId,
			group.Id,
			eventproto.Event_HOLDOUT_GROUP_CREATED,
			&eventproto.HoldoutGroupCreatedEvent{
				Id:              group.Id,
				Name:            group.Name,
				AssignmentType:  group.AssignmentType,
				Percentage:      group.Percentage,
				SegmentId:       group.SegmentId,
				StartAt:         group.StartAt,
				StopAt:          group.StopAt,
				StartAnalysisAt: group.StartAnalysisAt,
				GoalIds:         group.GoalIds,
			},
			group.HoldoutGroup,
			nil,
		)
		if err != nil {
			return err
		}
		return s.holdoutStorage.CreateHoldoutGroup(txCtx, group)
	})
	if err != nil {
		return nil, s.holdoutError(ctx, "Failed to create holdout group", req.EnvironmentId, err)
	}
	if err := s.publishHoldoutEvents(ctx, []*eventproto.Event{event}); err != nil {
		return nil, err
	}
	return &featureproto.CreateHoldoutGroupResponse{HoldoutGroup: group.HoldoutGroup}, nil
}

func (s *FeatureService) GetHoldoutGroup(
	ctx context.Context,
	req *featureproto.GetHoldoutGroupRequest,
) (*featureproto.GetHoldoutGroupResponse, error) {
	if _, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentId,
	); err != nil {
		return nil, err
	}
	if err := requireHoldoutID(req.Id, req.EnvironmentId); err != nil {
		return nil, holdoutStatus(err)
	}
	group, err := s.loadHoldoutGroup(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		return nil, s.holdoutError(ctx, "Failed to get holdout group", req.EnvironmentId, err)
	}
	return &featureproto.GetHoldoutGroupResponse{HoldoutGroup: group.HoldoutGroup}, nil
}

func (s *FeatureService) ListHoldoutGroups(
	ctx context.Context,
	req *featureproto.ListHoldoutGroupsRequest,
) (*featureproto.ListHoldoutGroupsResponse, error) {
	if _, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentId,
	); err != nil {
		return nil, err
	}
	if err := requireEnvironment(req.EnvironmentId); err != nil {
		return nil, holdoutStatus(err)
	}
	groups, next, total, err := s.holdoutStorage.ListHoldoutGroups(ctx, v2hs.ListHoldoutGroupsParams{
		EnvironmentID:  req.EnvironmentId,
		Archived:       optionalBool(req.Archived),
		Statuses:       req.Statuses,
		SearchKeyword:  req.SearchKeyword,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		PageSize:       int(req.PageSize),
		Cursor:         req.Cursor,
	})
	if err != nil {
		return nil, s.holdoutError(ctx, "Failed to list holdout groups", req.EnvironmentId, err)
	}
	if err := s.attachCurrentFeatures(ctx, req.EnvironmentId, groups); err != nil {
		return nil, s.holdoutError(ctx, "Failed to list holdout group features", req.EnvironmentId, err)
	}
	out := make([]*featureproto.HoldoutGroup, 0, len(groups))
	for _, group := range groups {
		out = append(out, group.HoldoutGroup)
	}
	return &featureproto.ListHoldoutGroupsResponse{
		HoldoutGroups: out,
		Cursor:        strconv.Itoa(next),
		TotalCount:    total,
	}, nil
}

func (s *FeatureService) UpdateHoldoutGroup(
	ctx context.Context,
	req *featureproto.UpdateHoldoutGroupRequest,
) (*featureproto.UpdateHoldoutGroupResponse, error) {
	editor, err := s.checkHoldoutAdmin(ctx, req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := requireHoldoutID(req.Id, req.EnvironmentId); err != nil {
		return nil, holdoutStatus(err)
	}
	goalIDs := []string(nil)
	goalIDsSet := req.GoalIds != nil
	if goalIDsSet {
		goalIDs = req.GoalIds.Values
		if err := s.validateGoalIDs(ctx, req.EnvironmentId, goalIDs); err != nil {
			return nil, err
		}
	}
	now := time.Now().Unix()
	var updated *holdoutdomain.HoldoutGroup
	var event *eventproto.Event
	err = s.dbClient.RunInTransactionV2(ctx, func(txCtx context.Context) error {
		if err := s.holdoutStorage.LockEnvironment(txCtx, req.EnvironmentId); err != nil {
			return err
		}
		current, err := s.holdoutStorage.GetHoldoutGroupForUpdate(txCtx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		previous := current.Clone()
		scheduleChanged := req.StartAt != nil || req.StopAt != nil || req.StartAnalysisAt != nil
		if err := current.Update(holdoutdomain.UpdateParams{
			Name:            req.Name,
			Description:     req.Description,
			StartAt:         req.StartAt,
			StopAt:          req.StopAt,
			StartAnalysisAt: req.StartAnalysisAt,
			GoalIDsSet:      goalIDsSet,
			GoalIDs:         goalIDs,
			UpdatedBy:       editor.Email,
			Now:             now,
		}); err != nil {
			return err
		}
		if scheduleChanged {
			active, err := s.holdoutStorage.ListActiveGroups(txCtx, req.EnvironmentId)
			if err != nil {
				return err
			}
			if conflict, ok := holdoutdomain.OverlappingGroup(
				current.StartAt, current.StopAt, current.Id, active,
			); ok {
				return holdoutdomain.NewPeriodOverlapError(conflict)
			}
		}
		event, err = newHoldoutEvent(
			editor,
			req.EnvironmentId,
			current.Id,
			eventproto.Event_HOLDOUT_GROUP_UPDATED,
			&eventproto.HoldoutGroupUpdatedEvent{Id: current.Id},
			current.HoldoutGroup,
			previous.HoldoutGroup,
		)
		if err != nil {
			return err
		}
		updated = current
		return s.holdoutStorage.UpdateHoldoutGroup(txCtx, current)
	})
	if err != nil {
		return nil, s.holdoutError(ctx, "Failed to update holdout group", req.EnvironmentId, err)
	}
	if err := s.publishHoldoutEvents(ctx, []*eventproto.Event{event}); err != nil {
		return nil, err
	}
	if err := s.attachCurrentFeatures(ctx, req.EnvironmentId, []*holdoutdomain.HoldoutGroup{updated}); err != nil {
		return nil, s.holdoutError(ctx, "Failed to load holdout group features", req.EnvironmentId, err)
	}
	return &featureproto.UpdateHoldoutGroupResponse{HoldoutGroup: updated.HoldoutGroup}, nil
}

func (s *FeatureService) StartHoldoutGroup(
	ctx context.Context,
	req *featureproto.StartHoldoutGroupRequest,
) (*featureproto.StartHoldoutGroupResponse, error) {
	group, err := s.transitionHoldout(ctx, req.Id, req.EnvironmentId, func(
		group *holdoutdomain.HoldoutGroup, now int64, active []holdoutdomain.GroupSnapshot,
	) error {
		if err := group.Start(now); err != nil {
			return err
		}
		if serving, ok := holdoutdomain.OtherServingGroup(group.Id, active); ok {
			return holdoutdomain.NewServingGroupExistsError(serving)
		}
		if conflict, ok := holdoutdomain.OverlappingGroup(group.StartAt, group.StopAt, group.Id, active); ok {
			return holdoutdomain.NewPeriodOverlapError(conflict)
		}
		return nil
	}, eventproto.Event_HOLDOUT_GROUP_STARTED, func(group *holdoutdomain.HoldoutGroup) proto.Message {
		return &eventproto.HoldoutGroupStartedEvent{Id: group.Id}
	})
	if err != nil {
		return nil, err
	}
	return &featureproto.StartHoldoutGroupResponse{HoldoutGroup: group}, nil
}

func (s *FeatureService) StartHoldoutGroupAnalysis(
	ctx context.Context,
	req *featureproto.StartHoldoutGroupAnalysisRequest,
) (*featureproto.StartHoldoutGroupAnalysisResponse, error) {
	group, err := s.transitionHoldout(ctx, req.Id, req.EnvironmentId, func(
		group *holdoutdomain.HoldoutGroup, now int64, _ []holdoutdomain.GroupSnapshot,
	) error {
		return group.StartAnalysis(now)
	}, eventproto.Event_HOLDOUT_GROUP_ANALYSIS_STARTED, func(group *holdoutdomain.HoldoutGroup) proto.Message {
		return &eventproto.HoldoutGroupAnalysisStartedEvent{Id: group.Id}
	})
	if err != nil {
		return nil, err
	}
	return &featureproto.StartHoldoutGroupAnalysisResponse{HoldoutGroup: group}, nil
}

func (s *FeatureService) StopHoldoutGroup(
	ctx context.Context,
	req *featureproto.StopHoldoutGroupRequest,
) (*featureproto.StopHoldoutGroupResponse, error) {
	group, err := s.transitionHoldout(ctx, req.Id, req.EnvironmentId, func(
		group *holdoutdomain.HoldoutGroup, now int64, _ []holdoutdomain.GroupSnapshot,
	) error {
		if req.Mode == featureproto.StopHoldoutGroupRequest_SCHEDULED {
			return group.FinishScheduled(now)
		}
		return group.StopManual(now)
	}, eventproto.Event_HOLDOUT_GROUP_STOPPED, func(group *holdoutdomain.HoldoutGroup) proto.Message {
		return &eventproto.HoldoutGroupStoppedEvent{Id: group.Id, Status: group.Status}
	})
	if err != nil {
		return nil, err
	}
	return &featureproto.StopHoldoutGroupResponse{HoldoutGroup: group}, nil
}

func (s *FeatureService) ArchiveHoldoutGroup(
	ctx context.Context,
	req *featureproto.ArchiveHoldoutGroupRequest,
) (*featureproto.ArchiveHoldoutGroupResponse, error) {
	group, err := s.transitionHoldout(ctx, req.Id, req.EnvironmentId, func(
		group *holdoutdomain.HoldoutGroup, now int64, _ []holdoutdomain.GroupSnapshot,
	) error {
		return group.Archive(now)
	}, eventproto.Event_HOLDOUT_GROUP_ARCHIVED, func(group *holdoutdomain.HoldoutGroup) proto.Message {
		return &eventproto.HoldoutGroupArchivedEvent{Id: group.Id}
	})
	if err != nil {
		return nil, err
	}
	return &featureproto.ArchiveHoldoutGroupResponse{HoldoutGroup: group}, nil
}

func (s *FeatureService) DeleteHoldoutGroup(
	ctx context.Context,
	req *featureproto.DeleteHoldoutGroupRequest,
) (*featureproto.DeleteHoldoutGroupResponse, error) {
	_, err := s.transitionHoldout(ctx, req.Id, req.EnvironmentId, func(
		group *holdoutdomain.HoldoutGroup, now int64, _ []holdoutdomain.GroupSnapshot,
	) error {
		return group.Delete(now)
	}, eventproto.Event_HOLDOUT_GROUP_DELETED, func(group *holdoutdomain.HoldoutGroup) proto.Message {
		return &eventproto.HoldoutGroupDeletedEvent{Id: group.Id}
	})
	if err != nil {
		return nil, err
	}
	return &featureproto.DeleteHoldoutGroupResponse{}, nil
}

func (s *FeatureService) AddHoldoutGroupFeatures(
	ctx context.Context,
	req *featureproto.AddHoldoutGroupFeaturesRequest,
) (*featureproto.AddHoldoutGroupFeaturesResponse, error) {
	editor, err := s.checkHoldoutAdmin(ctx, req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := requireHoldoutID(req.Id, req.EnvironmentId); err != nil {
		return nil, holdoutStatus(err)
	}
	if err := holdoutdomain.ValidateFeatureInputs(req.Features); err != nil {
		return nil, holdoutStatus(err)
	}
	now := time.Now().Unix()
	var updated *holdoutdomain.HoldoutGroup
	var events []*eventproto.Event
	err = s.dbClient.RunInTransactionV2(ctx, func(txCtx context.Context) error {
		if err := s.holdoutStorage.LockEnvironment(txCtx, req.EnvironmentId); err != nil {
			return err
		}
		group, err := s.holdoutStorage.GetHoldoutGroupForUpdate(txCtx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		previous := group.Clone()
		if err := group.CanAddFeatures(); err != nil {
			return err
		}
		created := make([]*holdoutdomain.HoldoutGroupFeature, 0, len(req.Features))
		for _, input := range req.Features {
			row, feature, err := s.prepareHoldoutFeature(txCtx, group, input, now)
			if err != nil {
				return err
			}
			created = append(created, row)
			groupEvent, err := newHoldoutEvent(
				editor, req.EnvironmentId, group.Id, eventproto.Event_HOLDOUT_GROUP_FEATURE_ADDED,
				&eventproto.HoldoutGroupFeatureAddedEvent{
					Id: group.Id, FeatureId: input.FeatureId, HoldoutVariationId: input.HoldoutVariationId,
				},
				group.HoldoutGroup, previous.HoldoutGroup,
			)
			if err != nil {
				return err
			}
			featureEvent, err := newFeatureHoldoutEvent(
				editor, req.EnvironmentId, input.FeatureId, eventproto.Event_HOLDOUT_GROUP_FEATURE_ADDED,
				&eventproto.HoldoutGroupFeatureAddedEvent{
					Id: group.Id, FeatureId: input.FeatureId, HoldoutVariationId: input.HoldoutVariationId,
				},
				feature,
			)
			if err != nil {
				return err
			}
			events = append(events, groupEvent, featureEvent)
		}
		if err := s.holdoutStorage.CreateHoldoutGroupFeatures(txCtx, created); err != nil {
			return err
		}
		group.Touch(editor.Email, now)
		updated = group
		return s.holdoutStorage.UpdateHoldoutGroup(txCtx, group)
	})
	if err != nil {
		return nil, s.holdoutError(ctx, "Failed to add holdout group features", req.EnvironmentId, err)
	}
	if err := s.publishHoldoutEvents(ctx, events); err != nil {
		return nil, err
	}
	if err := s.attachCurrentFeatures(ctx, req.EnvironmentId, []*holdoutdomain.HoldoutGroup{updated}); err != nil {
		return nil, s.holdoutError(ctx, "Failed to load holdout group features", req.EnvironmentId, err)
	}
	return &featureproto.AddHoldoutGroupFeaturesResponse{HoldoutGroup: updated.HoldoutGroup}, nil
}

func (s *FeatureService) RemoveHoldoutGroupFeatures(
	ctx context.Context,
	req *featureproto.RemoveHoldoutGroupFeaturesRequest,
) (*featureproto.RemoveHoldoutGroupFeaturesResponse, error) {
	editor, err := s.checkHoldoutAdmin(ctx, req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := requireHoldoutID(req.Id, req.EnvironmentId); err != nil {
		return nil, holdoutStatus(err)
	}
	if err := holdoutdomain.ValidateFeatureInputs(featureInputs(req.FeatureIds)); err != nil {
		return nil, holdoutStatus(err)
	}
	now := time.Now().Unix()
	var updated *holdoutdomain.HoldoutGroup
	var events []*eventproto.Event
	err = s.dbClient.RunInTransactionV2(ctx, func(txCtx context.Context) error {
		if err := s.holdoutStorage.LockEnvironment(txCtx, req.EnvironmentId); err != nil {
			return err
		}
		group, err := s.holdoutStorage.GetHoldoutGroupForUpdate(txCtx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		previous := group.Clone()
		if err := group.CanRemoveFeatures(); err != nil {
			return err
		}
		for _, featureID := range req.FeatureIds {
			if err := s.holdoutStorage.RemoveHoldoutGroupFeature(
				txCtx, req.EnvironmentId, group.Id, featureID, now,
			); err != nil {
				return err
			}
			feature, err := s.featureStorage.GetFeature(txCtx, featureID, req.EnvironmentId)
			if err != nil {
				if errors.Is(err, v2fs.ErrFeatureNotFound) {
					return pkgErr.NewErrorNotFound(pkgErr.HoldoutPackageName, "feature not found", "feature_id")
				}
				return err
			}
			groupEvent, err := newHoldoutEvent(
				editor, req.EnvironmentId, group.Id, eventproto.Event_HOLDOUT_GROUP_FEATURE_REMOVED,
				&eventproto.HoldoutGroupFeatureRemovedEvent{Id: group.Id, FeatureId: featureID},
				group.HoldoutGroup, previous.HoldoutGroup,
			)
			if err != nil {
				return err
			}
			featureEvent, err := newFeatureHoldoutEvent(
				editor, req.EnvironmentId, featureID, eventproto.Event_HOLDOUT_GROUP_FEATURE_REMOVED,
				&eventproto.HoldoutGroupFeatureRemovedEvent{Id: group.Id, FeatureId: featureID},
				feature.Feature,
			)
			if err != nil {
				return err
			}
			events = append(events, groupEvent, featureEvent)
		}
		group.Touch(editor.Email, now)
		updated = group
		return s.holdoutStorage.UpdateHoldoutGroup(txCtx, group)
	})
	if err != nil {
		return nil, s.holdoutError(ctx, "Failed to remove holdout group features", req.EnvironmentId, err)
	}
	if err := s.publishHoldoutEvents(ctx, events); err != nil {
		return nil, err
	}
	if err := s.attachCurrentFeatures(ctx, req.EnvironmentId, []*holdoutdomain.HoldoutGroup{updated}); err != nil {
		return nil, s.holdoutError(ctx, "Failed to load holdout group features", req.EnvironmentId, err)
	}
	return &featureproto.RemoveHoldoutGroupFeaturesResponse{HoldoutGroup: updated.HoldoutGroup}, nil
}

func (s *FeatureService) transitionHoldout(
	ctx context.Context,
	id, environmentID string,
	apply func(*holdoutdomain.HoldoutGroup, int64, []holdoutdomain.GroupSnapshot) error,
	eventType eventproto.Event_Type,
	eventMessage func(*holdoutdomain.HoldoutGroup) proto.Message,
) (*featureproto.HoldoutGroup, error) {
	editor, err := s.checkHoldoutAdmin(ctx, environmentID)
	if err != nil {
		return nil, err
	}
	if err := requireHoldoutID(id, environmentID); err != nil {
		return nil, holdoutStatus(err)
	}
	now := time.Now().Unix()
	var updated *holdoutdomain.HoldoutGroup
	var event *eventproto.Event
	err = s.dbClient.RunInTransactionV2(ctx, func(txCtx context.Context) error {
		if err := s.holdoutStorage.LockEnvironment(txCtx, environmentID); err != nil {
			return err
		}
		group, err := s.holdoutStorage.GetHoldoutGroupForUpdate(txCtx, id, environmentID)
		if err != nil {
			return err
		}
		previous := group.Clone()
		active, err := s.holdoutStorage.ListActiveGroups(txCtx, environmentID)
		if err != nil {
			return err
		}
		if err := apply(group, now, active); err != nil {
			return err
		}
		if err := s.revalidateSegmentForActivation(txCtx, previous, group); err != nil {
			return err
		}
		group.UpdatedBy = editor.Email
		event, err = newHoldoutEvent(
			editor, environmentID, group.Id, eventType, eventMessage(group), group.HoldoutGroup, previous.HoldoutGroup,
		)
		if err != nil {
			return err
		}
		updated = group
		return s.holdoutStorage.UpdateHoldoutGroup(txCtx, group)
	})
	if err != nil {
		return nil, s.holdoutError(ctx, "Failed to update holdout group status", environmentID, err)
	}
	if err := s.publishHoldoutEvents(ctx, []*eventproto.Event{event}); err != nil {
		return nil, err
	}
	if err := s.attachCurrentFeatures(ctx, environmentID, []*holdoutdomain.HoldoutGroup{updated}); err != nil {
		return nil, s.holdoutError(ctx, "Failed to load holdout group features", environmentID, err)
	}
	return updated.HoldoutGroup, nil
}

func (s *FeatureService) prepareHoldoutFeature(
	ctx context.Context,
	group *holdoutdomain.HoldoutGroup,
	input *featureproto.HoldoutGroupFeatureInput,
	now int64,
) (*holdoutdomain.HoldoutGroupFeature, *featureproto.Feature, error) {
	feature, err := s.featureStorage.GetFeature(ctx, input.FeatureId, group.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2fs.ErrFeatureNotFound) {
			return nil, nil, pkgErr.NewErrorNotFound(pkgErr.HoldoutPackageName, "feature not found", "feature_id")
		}
		return nil, nil, err
	}
	if feature.Archived || feature.Deleted {
		return nil, nil, holdoutdomain.NewFeatureUnavailableError(input.FeatureId)
	}
	if !holdoutdomain.VariationBelongsToFeature(input.HoldoutVariationId, feature.Variations) {
		return nil, nil, holdoutdomain.NewVariationNotOnFeatureError(input.FeatureId, input.HoldoutVariationId)
	}
	targets, err := s.holdoutStorage.ListLiveFeatureTargets(ctx, group.EnvironmentId, input.FeatureId)
	if err != nil {
		return nil, nil, err
	}
	for _, target := range targets {
		if target.GroupID == group.Id {
			return nil, nil, holdoutdomain.NewFeatureAlreadyAddedError(input.FeatureId)
		}
	}
	if conflict, ok := holdoutdomain.ConflictingFeatureTarget(input.FeatureId, group.Id, targets); ok {
		return nil, nil, holdoutdomain.NewFeatureInUseError(input.FeatureId, conflict)
	}
	row, err := holdoutdomain.NewHoldoutGroupFeature(
		group.Id, group.EnvironmentId, input.FeatureId, input.HoldoutVariationId, now,
	)
	if err != nil {
		return nil, nil, err
	}
	return row, feature.Feature, nil
}

func (s *FeatureService) loadHoldoutGroup(
	ctx context.Context,
	id, environmentID string,
) (*holdoutdomain.HoldoutGroup, error) {
	group, err := s.holdoutStorage.GetHoldoutGroup(ctx, id, environmentID)
	if err != nil {
		return nil, err
	}
	if err := s.attachCurrentFeatures(ctx, environmentID, []*holdoutdomain.HoldoutGroup{group}); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *FeatureService) attachCurrentFeatures(
	ctx context.Context,
	environmentID string,
	groups []*holdoutdomain.HoldoutGroup,
) error {
	if len(groups) == 0 {
		return nil
	}
	ids := make([]string, 0, len(groups))
	byID := make(map[string]*holdoutdomain.HoldoutGroup, len(groups))
	for _, group := range groups {
		ids = append(ids, group.Id)
		byID[group.Id] = group
		group.Features = nil
	}
	features, err := s.holdoutStorage.ListCurrentFeatures(ctx, environmentID, ids)
	if err != nil {
		return err
	}
	for _, feature := range features {
		group := byID[feature.HoldoutGroupId]
		if group == nil {
			continue
		}
		group.Features = append(group.Features, feature.HoldoutGroupFeature)
	}
	return nil
}

// revalidateSegmentForActivation reloads a SEGMENT cohort inside the transition
// transaction. Create checked the segment, but a later delete or upload can
// finish while the group is still WAITING.
func (s *FeatureService) revalidateSegmentForActivation(
	ctx context.Context,
	previous, group *holdoutdomain.HoldoutGroup,
) error {
	if previous.Status != featureproto.HoldoutGroup_WAITING ||
		group.Status != featureproto.HoldoutGroup_RUNNING ||
		group.AssignmentType != featureproto.HoldoutGroup_SEGMENT {
		return nil
	}
	segment, _, err := s.segmentStorage.GetSegment(ctx, group.SegmentId, group.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2fs.ErrSegmentNotFound) {
			return pkgErr.NewErrorNotFound(pkgErr.HoldoutPackageName, "segment not found", "segment_id")
		}
		return err
	}
	if segment.Deleted {
		return pkgErr.NewErrorNotFound(pkgErr.HoldoutPackageName, "segment not found", "segment_id")
	}
	if segment.Status == featureproto.Segment_UPLOADING {
		return pkgErr.NewErrorFailedPrecondition(
			pkgErr.FeaturePackageName, "segment users already uploading",
		).WithMessageKey("SegmentUsersAlreadyUploading")
	}
	return nil
}

func (s *FeatureService) ensureSegmentExists(ctx context.Context, segmentID, environmentID string) error {
	if segmentID == "" {
		return holdoutdomain.ErrSegmentRequired
	}
	segment, _, err := s.segmentStorage.GetSegment(ctx, segmentID, environmentID)
	if err != nil {
		if errors.Is(err, v2fs.ErrSegmentNotFound) {
			return pkgErr.NewErrorNotFound(pkgErr.HoldoutPackageName, "segment not found", "segment_id")
		}
		return err
	}
	if segment.Deleted {
		return pkgErr.NewErrorNotFound(pkgErr.HoldoutPackageName, "segment not found", "segment_id")
	}
	return nil
}

func (s *FeatureService) validateGoalIDs(ctx context.Context, environmentID string, goalIDs []string) error {
	seen := make(map[string]struct{}, len(goalIDs))
	for _, goalID := range goalIDs {
		if goalID == "" {
			return holdoutStatus(holdoutdomain.ErrGoalIDRequired)
		}
		if _, ok := seen[goalID]; ok {
			continue
		}
		seen[goalID] = struct{}{}
		resp, err := s.experimentClient.GetGoal(ctx, &experimentproto.GetGoalRequest{
			Id:            goalID,
			EnvironmentId: environmentID,
		})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return holdoutStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(
					pkgErr.HoldoutPackageName, "goal not found", "goal_id",
				).WithMessageKey("HoldoutGoalNotFound"))
			}
			return s.holdoutError(ctx, "Failed to get goal", environmentID, err)
		}
		if resp.GetGoal().GetArchived() || resp.GetGoal().GetDeleted() {
			return holdoutStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(
				pkgErr.HoldoutPackageName, "goal is archived or deleted", "goal_id",
			).WithMessageKey("HoldoutGoalNotFound"))
		}
	}
	return nil
}

func (s *FeatureService) checkHoldoutAdmin(
	ctx context.Context,
	environmentID string,
) (*eventproto.Editor, error) {
	return s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		environmentID,
	)
}

func (s *FeatureService) publishHoldoutEvents(ctx context.Context, events []*eventproto.Event) error {
	if errs := s.publishDomainEvents(ctx, events); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish holdout events",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Any("errors", errs))...,
		)
		return statusInternal.Err()
	}
	return nil
}

func (s *FeatureService) holdoutError(ctx context.Context, message, environmentID string, err error) error {
	var bktErr *pkgErr.BktError
	if !errors.As(err, &bktErr) {
		s.logger.Error(
			message,
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentID),
			)...,
		)
	}
	return holdoutStatus(err)
}

// rejectHoldoutVariationRemoval blocks deletion of a variation that a live holdout serves.
func (s *FeatureService) rejectHoldoutVariationRemoval(
	ctx context.Context,
	environmentID, featureID string,
	changes []*featureproto.VariationChange,
) error {
	if s.holdoutStorage == nil {
		return nil
	}
	var variationIDs []string
	for _, change := range changes {
		if change.GetChangeType() != featureproto.ChangeType_DELETE {
			continue
		}
		if id := change.GetVariation().GetId(); id != "" {
			variationIDs = append(variationIDs, id)
		}
	}
	if len(variationIDs) == 0 {
		return nil
	}
	targets, err := s.holdoutStorage.ListLiveFeatureTargets(ctx, environmentID, featureID)
	if err != nil {
		return err
	}
	for _, variationID := range variationIDs {
		for _, target := range targets {
			if target.VariationID == variationID {
				return holdoutdomain.NewVariationInUseError(variationID, target)
			}
		}
	}
	return nil
}

// lockHoldoutEnvironment serializes flag and segment writes with holdout
// mutations. Both sides lock environment_v2 first, so neither waits on the
// other's feature or segment row. A nil storage keeps existing tests working.
func (s *FeatureService) lockHoldoutEnvironment(ctx context.Context, environmentID string) error {
	if s.holdoutStorage == nil || environmentID == "" {
		return nil
	}
	return s.holdoutStorage.LockEnvironment(ctx, environmentID)
}

// rejectHoldoutFeatureArchive blocks archiving or deleting a flag that is a current holdout target.
func (s *FeatureService) rejectHoldoutFeatureArchive(
	ctx context.Context,
	environmentID, featureID string,
) error {
	if s.holdoutStorage == nil {
		return nil
	}
	targets, err := s.holdoutStorage.ListLiveFeatureTargets(ctx, environmentID, featureID)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return nil
	}
	return holdoutdomain.NewFeatureTargetedError(featureID, targets[0])
}

// ensureSegmentNotLockedByHoldout rejects edits to a segment referenced by a serving group.
func (s *FeatureService) ensureSegmentNotLockedByHoldout(
	ctx context.Context,
	segmentID, environmentID string,
) error {
	if s.holdoutStorage == nil || segmentID == "" {
		return nil
	}
	// Same environment lock as StartHoldoutGroup, taken before the segment row
	// is written, so a start cannot commit RUNNING between this read and the write.
	if err := s.holdoutStorage.LockEnvironment(ctx, environmentID); err != nil {
		return err
	}
	groups, err := s.holdoutStorage.ListActiveGroups(ctx, environmentID)
	if err != nil {
		return err
	}
	if group, ok := holdoutdomain.ServingGroupForSegment(segmentID, groups); ok {
		return holdoutdomain.NewSegmentLockedError(segmentID, group)
	}
	return nil
}

func newHoldoutEvent(
	editor *eventproto.Editor,
	environmentID, entityID string,
	eventType eventproto.Event_Type,
	event proto.Message,
	current, previous interface{},
) (*eventproto.Event, error) {
	return domainevent.NewEvent(
		editor,
		eventproto.Event_HOLDOUT_GROUP,
		entityID,
		eventType,
		event,
		environmentID,
		current,
		previous,
	)
}

func newFeatureHoldoutEvent(
	editor *eventproto.Editor,
	environmentID, featureID string,
	eventType eventproto.Event_Type,
	event proto.Message,
	feature *featureproto.Feature,
) (*eventproto.Event, error) {
	return domainevent.NewEvent(
		editor,
		eventproto.Event_FEATURE,
		featureID,
		eventType,
		event,
		environmentID,
		feature,
		nil,
	)
}

func requireHoldoutID(id, environmentID string) error {
	if id == "" {
		return pkgErr.NewErrorInvalidArgEmpty(pkgErr.HoldoutPackageName, "id is required", "id").
			WithMessageKey("HoldoutIDRequired")
	}
	return requireEnvironment(environmentID)
}

func requireEnvironment(environmentID string) error {
	if environmentID == "" {
		return pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.HoldoutPackageName, "environment_id is required", "environment_id",
		).WithMessageKey("HoldoutEnvironmentIDRequired")
	}
	return nil
}

func holdoutStatus(err error) error {
	if err == nil {
		return nil
	}
	return api.NewGRPCStatus(err).Err()
}

func optionalBool(value *wrapperspb.BoolValue) *bool {
	if value == nil {
		return nil
	}
	v := value.GetValue()
	return &v
}

func featureInputs(featureIDs []string) []*featureproto.HoldoutGroupFeatureInput {
	inputs := make([]*featureproto.HoldoutGroupFeatureInput, 0, len(featureIDs))
	for _, featureID := range featureIDs {
		inputs = append(inputs, &featureproto.HoldoutGroupFeatureInput{
			FeatureId:          featureID,
			HoldoutVariationId: "placeholder",
		})
	}
	return inputs
}
