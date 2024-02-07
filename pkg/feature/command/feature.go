// Copyright 2024 The Bucketeer Authors.
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

package command

import (
	"context"
	"strings"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type FeatureCommandHandler struct {
	feature      *domain.Feature
	eventFactory *FeatureEventFactory
	Events       []*eventproto.Event
}

func NewFeatureCommandHandler(
	editor *eventproto.Editor,
	feature *domain.Feature,
	environmentNamespace string,
	comment string,
) *FeatureCommandHandler {
	return &FeatureCommandHandler{
		feature: feature,
		eventFactory: &FeatureEventFactory{
			editor:               editor,
			feature:              feature,
			environmentNamespace: environmentNamespace,
			comment:              comment,
		},
		Events: []*eventproto.Event{},
	}
}

// for unit test
func NewEmptyFeatureCommandHandler() *FeatureCommandHandler {
	return &FeatureCommandHandler{}
}

func (h *FeatureCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateFeatureCommand:
		return h.CreateFeature(ctx, c)
	case *proto.EnableFeatureCommand:
		return h.EnableFeature(ctx, c)
	case *proto.DisableFeatureCommand:
		return h.DisableFeature(ctx, c)
	case *proto.ArchiveFeatureCommand:
		return h.ArchiveFeature(ctx, c)
	case *proto.UnarchiveFeatureCommand:
		return h.UnarchiveFeature(ctx, c)
	case *proto.DeleteFeatureCommand:
		return h.DeleteFeature(ctx, c)
	case *proto.AddUserToVariationCommand:
		return h.AddUserToVariation(ctx, c)
	case *proto.RemoveUserFromVariationCommand:
		return h.RemoveUserFromVariation(ctx, c)
	case *proto.AddRuleCommand:
		return h.AddRule(ctx, c)
	case *proto.ChangeRuleStrategyCommand:
		return h.ChangeRuleStrategy(ctx, c)
	case *proto.ChangeRulesOrderCommand:
		return h.ChangeRulesOrder(ctx, c)
	case *proto.DeleteRuleCommand:
		return h.DeleteRule(ctx, c)
	case *proto.AddClauseCommand:
		return h.AddClause(ctx, c)
	case *proto.DeleteClauseCommand:
		return h.DeleteClause(ctx, c)
	case *proto.ChangeClauseAttributeCommand:
		return h.ChangeClauseAttribute(ctx, c)
	case *proto.ChangeClauseOperatorCommand:
		return h.ChangeClauseOperator(ctx, c)
	case *proto.AddClauseValueCommand:
		return h.AddClauseValue(ctx, c)
	case *proto.RemoveClauseValueCommand:
		return h.RemoveClauseValue(ctx, c)
	case *proto.ChangeDefaultStrategyCommand:
		return h.ChangeDefaultStrategy(ctx, c)
	case *proto.ChangeOffVariationCommand:
		return h.ChangeOffVariation(ctx, c)
	case *proto.ChangeFixedStrategyCommand:
		return h.ChangeFixedStrategy(ctx, c)
	case *proto.ChangeRolloutStrategyCommand:
		return h.ChangeRolloutStrategy(ctx, c)
	case *proto.AddVariationCommand:
		return h.AddVariation(ctx, c)
	case *proto.RemoveVariationCommand:
		return h.RemoveVariation(ctx, c)
	case *proto.ChangeVariationValueCommand:
		return h.ChangeVariationValue(ctx, c)
	case *proto.ChangeVariationNameCommand:
		return h.ChangeVariationName(ctx, c)
	case *proto.ChangeVariationDescriptionCommand:
		return h.ChangeVariationDescription(ctx, c)
	case *proto.RenameFeatureCommand:
		return h.RenameFeature(ctx, c)
	case *proto.ChangeDescriptionCommand:
		return h.ChangeDescription(ctx, c)
	case *proto.AddTagCommand:
		return h.AddTag(ctx, c)
	case *proto.RemoveTagCommand:
		return h.RemoveTag(ctx, c)
	case *proto.IncrementFeatureVersionCommand:
		return h.IncrementFeatureVersion(ctx, c)
	case *proto.CloneFeatureCommand:
		return h.CloneFeature(ctx, c)
	case *proto.ResetSamplingSeedCommand:
		return h.ResetSamplingSeed(ctx, c)
	case *proto.AddPrerequisiteCommand:
		return h.AddPrerequisite(ctx, c)
	case *proto.ChangePrerequisiteVariationCommand:
		return h.ChangePrerequisiteVariation(ctx, c)
	case *proto.RemovePrerequisiteCommand:
		return h.RemovePrerequisite(ctx, c)
	default:
		return errBadCommand
	}
}

func (h *FeatureCommandHandler) CreateFeature(ctx context.Context, cmd *proto.CreateFeatureCommand) error {
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_CREATED, &eventproto.FeatureCreatedEvent{
		Id:                       h.feature.Id,
		Name:                     h.feature.Name,
		Description:              h.feature.Description,
		User:                     "default",
		Variations:               h.feature.Variations,
		DefaultOnVariationIndex:  cmd.DefaultOnVariationIndex,
		DefaultOffVariationIndex: cmd.DefaultOffVariationIndex,
		VariationType:            cmd.VariationType,
		Tags:                     h.feature.Tags,
		Prerequisites:            h.feature.Prerequisites,
		Targets:                  h.feature.Targets,
		Rules:                    h.feature.Rules,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) EnableFeature(ctx context.Context, cmd *proto.EnableFeatureCommand) error {
	if err := h.feature.Enable(); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_ENABLED, &eventproto.FeatureEnabledEvent{
		Id: h.feature.Id,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) DisableFeature(ctx context.Context, cmd *proto.DisableFeatureCommand) error {
	if err := h.feature.Disable(); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_DISABLED, &eventproto.FeatureDisabledEvent{
		Id: h.feature.Id,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ArchiveFeature(ctx context.Context, cmd *proto.ArchiveFeatureCommand) error {
	if err := h.feature.Archive(); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_ARCHIVED, &eventproto.FeatureArchivedEvent{
		Id: h.feature.Id,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) UnarchiveFeature(ctx context.Context, cmd *proto.UnarchiveFeatureCommand) error {
	if err := h.feature.Unarchive(); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_UNARCHIVED, &eventproto.FeatureUnarchivedEvent{
		Id: h.feature.Id,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) DeleteFeature(ctx context.Context, cmd *proto.DeleteFeatureCommand) error {
	if err := h.feature.Delete(); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_DELETED, &eventproto.FeatureDeletedEvent{
		Id: h.feature.Id,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) IncrementFeatureVersion(
	ctx context.Context,
	cmd *proto.IncrementFeatureVersionCommand,
) error {
	err := h.feature.IncrementVersion()
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_FEATURE_VERSION_INCREMENTED,
		&eventproto.FeatureVersionIncrementedEvent{
			Id:      h.feature.Id,
			Version: h.feature.Version,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) CloneFeature(ctx context.Context, cmd *proto.CloneFeatureCommand) error {
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_CLONED, &eventproto.FeatureClonedEvent{
		Id:              h.feature.Id,
		Name:            h.feature.Name,
		Description:     h.feature.Description,
		Variations:      h.feature.Variations,
		Targets:         h.feature.Targets,
		Rules:           h.feature.Rules,
		DefaultStrategy: h.feature.DefaultStrategy,
		OffVariation:    h.feature.OffVariation,
		Tags:            h.feature.Tags,
		Maintainer:      h.feature.Maintainer,
		VariationType:   h.feature.VariationType,
		Prerequisites:   h.feature.Prerequisites,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ResetSamplingSeed(ctx context.Context, cmd *proto.ResetSamplingSeedCommand) error {
	if err := h.feature.ResetSamplingSeed(); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_SAMPLING_SEED_RESET,
		&eventproto.FeatureSamplingSeedResetEvent{
			SamplingSeed: h.feature.SamplingSeed,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) AddPrerequisite(ctx context.Context, cmd *proto.AddPrerequisiteCommand) error {
	if err := h.feature.AddPrerequisite(cmd.Prerequisite.FeatureId, cmd.Prerequisite.VariationId); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_PREREQUISITE_ADDED,
		&eventproto.PrerequisiteAddedEvent{
			Prerequisite: cmd.Prerequisite,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangePrerequisiteVariation(
	ctx context.Context,
	cmd *proto.ChangePrerequisiteVariationCommand,
) error {
	if err := h.feature.ChangePrerequisiteVariation(cmd.Prerequisite.FeatureId, cmd.Prerequisite.VariationId); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_PREREQUISITE_VARIATION_CHANGED,
		&eventproto.PrerequisiteVariationChangedEvent{
			Prerequisite: cmd.Prerequisite,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) RemovePrerequisite(ctx context.Context, cmd *proto.RemovePrerequisiteCommand) error {
	if err := h.feature.RemovePrerequisite(cmd.FeatureId); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_PREREQUISITE_REMOVED,
		&eventproto.PrerequisiteRemovedEvent{
			FeatureId: cmd.FeatureId,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) AddUserToVariation(ctx context.Context, cmd *proto.AddUserToVariationCommand) error {
	userID := strings.TrimSpace(cmd.User)
	err := h.feature.AddUserToVariation(cmd.Id, userID)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_VARIATION_USER_ADDED, &eventproto.VariationUserAddedEvent{
		FeatureId: h.feature.Id,
		Id:        cmd.Id,
		User:      userID,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) RemoveUserFromVariation(
	ctx context.Context,
	cmd *proto.RemoveUserFromVariationCommand,
) error {
	userID := strings.TrimSpace(cmd.User)
	err := h.feature.RemoveUserFromVariation(cmd.Id, userID)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_VARIATION_USER_REMOVED,
		&eventproto.VariationUserRemovedEvent{
			FeatureId: h.feature.Id,
			Id:        cmd.Id,
			User:      userID,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) AddRule(ctx context.Context, cmd *proto.AddRuleCommand) error {
	for _, clause := range cmd.Rule.Clauses {
		id, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		clause.Id = id.String()
	}
	err := h.feature.AddRule(cmd.Rule)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_RULE_ADDED, &eventproto.FeatureRuleAddedEvent{
		Id:   h.feature.Id,
		Rule: cmd.Rule,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeRuleStrategy(ctx context.Context, cmd *proto.ChangeRuleStrategyCommand) error {
	if err := h.feature.ChangeRuleStrategy(cmd.RuleId, cmd.Strategy); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_FEATURE_RULE_STRATEGY_CHANGED,
		&eventproto.FeatureChangeRuleStrategyEvent{
			FeatureId: h.feature.Id,
			RuleId:    cmd.RuleId,
			Strategy:  cmd.Strategy,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeRulesOrder(ctx context.Context, cmd *proto.ChangeRulesOrderCommand) error {
	if err := h.feature.ChangeRulesOrder(cmd.RuleIds); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_FEATURE_RULES_ORDER_CHANGED,
		&eventproto.FeatureRulesOrderChangedEvent{
			FeatureId: h.feature.Id,
			RuleIds:   cmd.RuleIds,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) DeleteRule(ctx context.Context, cmd *proto.DeleteRuleCommand) error {
	err := h.feature.DeleteRule(cmd.Id)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_FEATURE_RULE_DELETED, &eventproto.FeatureRuleDeletedEvent{
		Id:     h.feature.Id,
		RuleId: cmd.Id,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) AddClause(ctx context.Context, cmd *proto.AddClauseCommand) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	cmd.Clause.Id = id.String()
	err = h.feature.AddClause(cmd.RuleId, cmd.Clause)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_RULE_CLAUSE_ADDED, &eventproto.RuleClauseAddedEvent{
		FeatureId: h.feature.Id,
		RuleId:    cmd.RuleId,
		Clause:    cmd.Clause,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) DeleteClause(ctx context.Context, cmd *proto.DeleteClauseCommand) error {
	err := h.feature.DeleteClause(cmd.RuleId, cmd.Id)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_RULE_CLAUSE_DELETED, &eventproto.RuleClauseDeletedEvent{
		FeatureId: h.feature.Id,
		RuleId:    cmd.RuleId,
		Id:        cmd.Id,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeClauseAttribute(
	ctx context.Context,
	cmd *proto.ChangeClauseAttributeCommand,
) error {
	err := h.feature.ChangeClauseAttribute(cmd.RuleId, cmd.Id, cmd.Attribute)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_CLAUSE_ATTRIBUTE_CHANGED,
		&eventproto.ClauseAttributeChangedEvent{
			FeatureId: h.feature.Id,
			RuleId:    cmd.RuleId,
			Id:        cmd.Id,
			Attribute: cmd.Attribute,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeClauseOperator(
	ctx context.Context,
	cmd *proto.ChangeClauseOperatorCommand,
) error {
	err := h.feature.ChangeClauseOperator(cmd.RuleId, cmd.Id, cmd.Operator)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_CLAUSE_OPERATOR_CHANGED,
		&eventproto.ClauseOperatorChangedEvent{
			FeatureId: h.feature.Id,
			RuleId:    cmd.RuleId,
			Id:        cmd.Id,
			Operator:  cmd.Operator,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) AddClauseValue(ctx context.Context, cmd *proto.AddClauseValueCommand) error {
	err := h.feature.AddClauseValue(cmd.RuleId, cmd.Id, cmd.Value)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_CLAUSE_VALUE_ADDED, &eventproto.ClauseValueAddedEvent{
		FeatureId: h.feature.Id,
		RuleId:    cmd.RuleId,
		Id:        cmd.Id,
		Value:     cmd.Value,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) RemoveClauseValue(ctx context.Context, cmd *proto.RemoveClauseValueCommand) error {
	err := h.feature.RemoveClauseValue(cmd.RuleId, cmd.Id, cmd.Value)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(eventproto.Event_CLAUSE_VALUE_REMOVED, &eventproto.ClauseValueRemovedEvent{
		FeatureId: h.feature.Id,
		RuleId:    cmd.RuleId,
		Id:        cmd.Id,
		Value:     cmd.Value,
	})
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeDefaultStrategy(
	ctx context.Context,
	cmd *proto.ChangeDefaultStrategyCommand,
) error {
	err := h.feature.ChangeDefaultStrategy(cmd.Strategy)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_FEATURE_DEFAULT_STRATEGY_CHANGED,
		&eventproto.FeatureDefaultStrategyChangedEvent{
			Id:       h.feature.Id,
			Strategy: cmd.Strategy,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeOffVariation(ctx context.Context, cmd *proto.ChangeOffVariationCommand) error {
	err := h.feature.ChangeOffVariation(cmd.Id)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_FEATURE_OFF_VARIATION_CHANGED,
		&eventproto.FeatureOffVariationChangedEvent{
			Id:           h.feature.Id,
			OffVariation: cmd.Id,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeFixedStrategy(
	ctx context.Context,
	cmd *proto.ChangeFixedStrategyCommand,
) error {
	if err := h.feature.ChangeFixedStrategy(cmd.RuleId, cmd.Strategy); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_RULE_FIXED_STRATEGY_CHANGED,
		&eventproto.FeatureFixedStrategyChangedEvent{
			FeatureId: h.feature.Id,
			RuleId:    cmd.RuleId,
			Strategy:  cmd.Strategy,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeRolloutStrategy(
	ctx context.Context,
	cmd *proto.ChangeRolloutStrategyCommand,
) error {
	if err := h.feature.ChangeRolloutStrategy(cmd.RuleId, cmd.Strategy); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_RULE_ROLLOUT_STRATEGY_CHANGED,
		&eventproto.FeatureRolloutStrategyChangedEvent{
			FeatureId: h.feature.Id,
			RuleId:    cmd.RuleId,
			Strategy:  cmd.Strategy,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) AddVariation(ctx context.Context, cmd *proto.AddVariationCommand) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	if err = h.feature.AddVariation(id.String(), cmd.Value, cmd.Name, cmd.Description); err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_FEATURE_VARIATION_ADDED,
		&eventproto.FeatureVariationAddedEvent{
			Id: h.feature.Id,
			Variation: &proto.Variation{
				Id:          id.String(),
				Value:       cmd.Value,
				Name:        cmd.Name,
				Description: cmd.Description,
			},
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) RemoveVariation(ctx context.Context, cmd *proto.RemoveVariationCommand) error {
	err := h.feature.RemoveVariation(cmd.Id)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_FEATURE_VARIATION_REMOVED,
		&eventproto.FeatureVariationRemovedEvent{
			Id:          h.feature.Id,
			VariationId: cmd.Id,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeVariationValue(
	ctx context.Context,
	cmd *proto.ChangeVariationValueCommand,
) error {
	err := h.feature.ChangeVariationValue(cmd.Id, cmd.Value)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_VARIATION_VALUE_CHANGED,
		&eventproto.VariationValueChangedEvent{
			FeatureId: h.feature.Id,
			Id:        cmd.Id,
			Value:     cmd.Value,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeVariationName(
	ctx context.Context,
	cmd *proto.ChangeVariationNameCommand,
) error {
	err := h.feature.ChangeVariationName(cmd.Id, cmd.Name)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_VARIATION_NAME_CHANGED,
		&eventproto.VariationNameChangedEvent{
			FeatureId: h.feature.Id,
			Id:        cmd.Id,
			Name:      cmd.Name,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}

func (h *FeatureCommandHandler) ChangeVariationDescription(
	ctx context.Context,
	cmd *proto.ChangeVariationDescriptionCommand,
) error {
	err := h.feature.ChangeVariationDescription(cmd.Id, cmd.Description)
	if err != nil {
		return err
	}
	event, err := h.eventFactory.CreateEvent(
		eventproto.Event_VARIATION_DESCRIPTION_CHANGED,
		&eventproto.VariationDescriptionChangedEvent{
			FeatureId:   h.feature.Id,
			Id:          cmd.Id,
			Description: cmd.Description,
		},
	)
	if err != nil {
		return err
	}
	h.Events = append(h.Events, event)
	return nil
}
