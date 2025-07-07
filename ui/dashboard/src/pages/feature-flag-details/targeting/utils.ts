import { isEqual, omit } from 'lodash';
import { v4 as uuid } from 'uuid';
import {
  Feature,
  FeatureChangeType,
  FeaturePrerequisite,
  FeatureRule,
  FeatureRuleChange,
  FeatureRuleClause,
  FeatureRuleClauseOperator,
  FeatureRuleStrategy,
  FeatureTarget,
  FeatureVariation,
  PrerequisiteChange,
  RuleChange,
  StrategyType,
  TargetChange
} from '@types';
import { DiscardChangesStateData } from '.';
import {
  DefaultRuleSchema,
  RuleSchema,
  StrategySchema,
  TargetingSchema
} from './form-schema';
import {
  IndividualRuleItem,
  PrerequisiteSchema,
  RuleClauseType
} from './types';

export const getAlreadyTargetedVariation = (
  targets: IndividualRuleItem[],
  variationId: string,
  label: string
) => {
  const newTargets = targets
    .filter(target => target.variationId !== variationId)
    ?.map(item => ({
      ...item,
      users: item.users.map(user => user.toLowerCase())
    }));

  return newTargets.find(target => target.users.includes(label.toLowerCase()));
};

export const getDefaultRolloutStrategy = (feature: Feature) =>
  feature.variations?.map(val => ({
    variation: val.id,
    weight: 0
  }));

export const getDefaultStrategy = (
  feature: Feature,
  defaultStrategy: FeatureRuleStrategy
) => {
  const { type, fixedStrategy, rolloutStrategy } = defaultStrategy || {};
  return {
    ...defaultStrategy,
    fixedStrategy: {
      variation: fixedStrategy?.variation || ''
    },
    rolloutStrategy: rolloutStrategy?.variations?.length
      ? rolloutStrategy.variations?.map(item => ({
          ...item,
          weight: item.weight / 1000 || 0
        }))
      : getDefaultRolloutStrategy(feature),
    currentOption:
      type === StrategyType.FIXED
        ? fixedStrategy?.variation
        : StrategyType.ROLLOUT
  };
};

export const getDefaultRule = (feature: Feature) => ({
  id: uuid(),
  strategy: {
    currentOption: feature?.variations[0]?.id || '',
    fixedStrategy: {
      variation: feature?.variations[0]?.id || ''
    },
    rolloutStrategy: getDefaultRolloutStrategy(feature),
    type: StrategyType.FIXED
  },
  clauses: [
    {
      id: uuid(),
      type: RuleClauseType.COMPARE,
      attribute: '',
      operator: FeatureRuleClauseOperator.EQUALS,
      values: []
    }
  ]
});

const handleCreatePrerequisites = (prerequisites: FeaturePrerequisite[]) =>
  prerequisites.map(({ featureId, variationId }) => ({
    featureId,
    variationId
  }));

const handleCreateIndividualRules = (
  targets: FeatureTarget[],
  variations: FeatureVariation[]
) => {
  const hasUsers = targets.find(item => item?.users.length);
  if (hasUsers)
    return targets.map(({ variation, users }) => {
      const currentVariation = variations.find(item => item.id === variation);
      return {
        variationId: variation,
        users,
        name: currentVariation?.name || currentVariation?.value || ''
      };
    });
  return [];
};

const getClauseType = (operator: FeatureRuleClauseOperator) => {
  const { FEATURE_FLAG, BEFORE, AFTER, SEGMENT } = FeatureRuleClauseOperator;
  if (operator === FEATURE_FLAG) return RuleClauseType.FEATURE_FLAG;
  if ([BEFORE, AFTER].includes(operator)) return RuleClauseType.DATE;
  if (operator === SEGMENT) return RuleClauseType.SEGMENT;
  return RuleClauseType.COMPARE;
};

const handleCreateSegmentRules = (feature: Feature) => {
  return feature.rules.map(({ id, strategy, clauses }) => ({
    id,
    strategy: getDefaultStrategy(feature, strategy),
    clauses: clauses.map(clause => ({
      ...clause,
      type: getClauseType(clause.operator)
    }))
  }));
};

export const handleCreateDefaultValues = (feature: Feature) => {
  const {
    prerequisites,
    targets,
    variations,
    defaultStrategy,
    enabled,
    offVariation
  } = feature || {};
  const _prerequisites = handleCreatePrerequisites(prerequisites);
  const individualRules = handleCreateIndividualRules(targets, variations);
  const _defaultStrategy = getDefaultStrategy(feature, defaultStrategy);
  const segmentRules = handleCreateSegmentRules(feature);
  return {
    prerequisites: _prerequisites,
    individualRules,
    segmentRules,
    defaultRule: {
      ..._defaultStrategy,
      type: defaultStrategy?.type,
      currentOption:
        defaultStrategy?.type === StrategyType.FIXED
          ? defaultStrategy?.fixedStrategy?.variation
          : StrategyType.MANUAL,
      manualStrategy: _defaultStrategy?.rolloutStrategy
    },
    enabled,
    comment: '',
    requireComment: false,
    resetSampling: false,
    scheduleType: enabled ? 'DISABLE' : 'ENABLE',
    scheduleAt: String(Math.floor((new Date().getTime() + 3600000) / 1000)),
    offVariation
  } as TargetingSchema;
};

export const createVariationLabel = (variation: FeatureVariation): string => {
  if (variation == null) {
    return 'None';
  }
  const maxLength = 150;
  const ellipsis = '...';
  const label = variation.name
    ? variation.name + ' - ' + variation.value
    : variation.value;
  if (label.length > maxLength) {
    return label.slice(0, maxLength - ellipsis.length) + ellipsis;
  }
  return label;
};

const handleGetStrategy = (
  strategy?: StrategySchema | FeatureRuleStrategy
): Partial<FeatureRuleStrategy> => {
  const { type, fixedStrategy, rolloutStrategy } = strategy || {};
  if (type === StrategyType.FIXED)
    return {
      type,
      fixedStrategy: {
        variation: fixedStrategy?.variation || ''
      }
    };
  return {
    type,
    rolloutStrategy: {
      variations:
        ((rolloutStrategy as FeatureRuleStrategy['rolloutStrategy'])?.variations
          ? (rolloutStrategy as FeatureRuleStrategy['rolloutStrategy'])
              ?.variations
          : (rolloutStrategy as StrategySchema['rolloutStrategy'])
        )?.map(item => ({
          ...item,
          weight: item.weight * 1000
        })) || []
    }
  };
};

export const handleGetDefaultRuleStrategy = (
  defaultRule?: DefaultRuleSchema
): Partial<FeatureRuleStrategy> => {
  const { currentOption, manualStrategy } = defaultRule || {};
  if (currentOption === StrategyType.MANUAL) {
    return {
      type: StrategyType.ROLLOUT,
      rolloutStrategy: {
        variations:
          manualStrategy?.map(item => ({
            ...item,
            weight: item.weight * 1000
          })) || []
      }
    };
  }
  return handleGetStrategy(defaultRule as unknown as StrategySchema);
};

export const handleCheckSegmentRules = (
  featureRules: FeatureRule[],
  segmentRules?: RuleSchema[]
) => {
  const ruleChanges: RuleChange[] = [];
  const getRuleItem = (rule: RuleSchema, type: FeatureChangeType) => {
    return {
      changeType: type,
      rule: {
        id: rule.id,
        clauses: rule.clauses.map(
          clause => omit(clause, 'type') as FeatureRuleClause
        ),
        strategy: handleGetStrategy(rule.strategy)
      }
    };
  };

  featureRules.forEach(item => {
    const currentRule = segmentRules?.find(rule => rule.id === item.id);
    if (!currentRule) {
      ruleChanges.push({
        changeType: 'DELETE',
        rule: item
      });
    }
  });

  segmentRules?.forEach(item => {
    const currentRule = featureRules?.find(rule => rule.id === item.id);
    if (!currentRule) ruleChanges.push(getRuleItem(item, 'CREATE'));

    const formattedRule = {
      ...currentRule,
      clauses: currentRule?.clauses,
      strategy: handleGetStrategy(item?.strategy)
    } as FeatureRuleChange;

    const formattedItem = {
      ...item,
      clauses: item?.clauses?.map(item => omit(item, 'type')),
      strategy: handleGetStrategy(item?.strategy)
    } as FeatureRuleChange;

    if (currentRule && !isEqual(formattedRule, formattedItem)) {
      ruleChanges.push({
        changeType: 'UPDATE',
        rule: formattedItem
      });
    }
  });
  return ruleChanges;
};

export const handleCheckIndividualRules = (
  featureTargets: FeatureTarget[],
  individualRules?: TargetingSchema['individualRules']
) => {
  const targetChanges: TargetChange[] = [];

  featureTargets.forEach(item => {
    const { variation, users } = item;
    const lengthUsers = users.length;

    const currentTarget = individualRules?.find(
      rule => rule.variationId === variation
    );

    const lengthTargetUsers = currentTarget?.users?.length || 0;

    if (
      (!currentTarget && lengthUsers) ||
      (currentTarget &&
        (lengthUsers > lengthTargetUsers ||
          (lengthUsers === lengthTargetUsers &&
            !isEqual(users, currentTarget.users))))
    ) {
      targetChanges.push({
        changeType: 'DELETE',
        target: item
      });
    }

    if (currentTarget && currentTarget.users.length > lengthUsers) {
      targetChanges.push({
        changeType: 'UPDATE',
        target: {
          users: currentTarget.users.length ? currentTarget.users : item.users,
          variation: currentTarget.variationId
        }
      });
    }
  });

  individualRules?.forEach(item => {
    const { users, variationId } = item;
    const lengthUsers = users.length;

    const currentTarget = featureTargets?.find(
      target => target.variation === variationId
    );

    const lengthTargetUsers = currentTarget?.users?.length || 0;

    if (
      (!currentTarget && lengthUsers) ||
      (currentTarget &&
        lengthUsers &&
        (lengthUsers < lengthTargetUsers ||
          (lengthUsers === lengthTargetUsers &&
            !isEqual(users, currentTarget.users))))
    ) {
      targetChanges.push({
        changeType: 'CREATE',
        target: {
          users: item.users,
          variation: item.variationId
        }
      });
    }
  });
  return targetChanges;
};

export const handleCheckPrerequisites = (
  featurePrerequisites: FeaturePrerequisite[],
  prerequisites?: TargetingSchema['prerequisites']
) => {
  const prerequisiteChanges: PrerequisiteChange[] = [];
  featurePrerequisites.forEach(item => {
    const currentTarget = prerequisites?.find(
      pre => pre.variationId === item.variationId
    );
    if (!currentTarget) {
      prerequisiteChanges.push({
        changeType: 'DELETE',
        prerequisite: item
      });
    }
  });

  prerequisites?.forEach(item => {
    const currentTarget = featurePrerequisites?.find(
      pre => pre.variationId === item.variationId
    );
    if (!currentTarget) {
      prerequisiteChanges.push({
        changeType: 'CREATE',
        prerequisite: item
      });
    }
  });
  return prerequisiteChanges;
};

const getPrerequisiteDiscardChangeData = (
  item: PrerequisiteChange,
  activeFeatures: Feature[]
) => {
  const currentFeature = activeFeatures.find(
    feature => feature.id === item.prerequisite.featureId
  );
  let variationIndex = -1;
  const variation = currentFeature?.variations.find((variation, index) => {
    if (variation.id === item.prerequisite.variationId) {
      variationIndex = index;
      return variation;
    }
  });
  return {
    label: currentFeature?.name || '',
    featureId: currentFeature?.id || '',
    variationIndex,
    variation
  };
};

const getPrerequisiteDiscardChangeItem = (
  checkPrerequisites: PrerequisiteChange[],
  currentPrerequisite: PrerequisiteChange,
  changeType: FeatureChangeType
) => {
  return checkPrerequisites.find(
    pre =>
      pre.changeType === changeType &&
      pre.prerequisite.featureId === currentPrerequisite.prerequisite.featureId
  );
};

export const handleCheckPrerequisiteDiscardChanges = (
  prerequisites: PrerequisiteSchema[],
  feature: Feature,
  activeFeatures: Feature[]
) => {
  const checkPrerequisites = handleCheckPrerequisites(
    feature.prerequisites,
    prerequisites
  );
  if (!checkPrerequisites.length) return null;
  const prerequisiteChanges: DiscardChangesStateData[] = [];

  checkPrerequisites.forEach(item => {
    const isDeleteItem = item.changeType === 'DELETE';
    const isCreateItem = item.changeType === 'CREATE';
    const prerequisiteData = getPrerequisiteDiscardChangeData(
      item,
      activeFeatures
    );
    const isExistedItem = prerequisiteChanges.find(
      pre => pre.featureId === item.prerequisite.featureId
    );
    if (isExistedItem) return;
    if (isDeleteItem) {
      const updateItem = getPrerequisiteDiscardChangeItem(
        checkPrerequisites,
        item,
        'CREATE'
      );
      return prerequisiteChanges.push({
        ...prerequisiteData,
        labelType: updateItem ? 'UPDATE' : 'REMOVE'
      });
    }
    if (isCreateItem) {
      const deleteItem = getPrerequisiteDiscardChangeItem(
        checkPrerequisites,
        item,
        'DELETE'
      );
      return prerequisiteChanges.push({
        ...prerequisiteData,
        labelType: deleteItem ? 'UPDATE' : 'ADD'
      });
    }
  });
  return prerequisiteChanges;
};

const getIndividualDiscardChangeData = (
  feature: Feature,
  individual: FeatureTarget
) => {
  let variationIndex = -1;
  const variation = feature?.variations.find((variation, index) => {
    if (variation.id === individual.variation) {
      variationIndex = index;
      return variation;
    }
  });

  return {
    label: individual.users?.join(', ') || '',
    variationIndex,
    variation
  };
};

const getIndividualDiscardChangeItem = (
  checkIndividuals: TargetChange[],
  currentIndividual: TargetChange,
  changeType: FeatureChangeType
) => {
  return checkIndividuals.find(
    pre =>
      pre.changeType === changeType &&
      pre.target.variation === currentIndividual.target.variation
  );
};

export const handleCheckIndividualDiscardChanges = (
  feature: Feature,
  rules: IndividualRuleItem[]
) => {
  const checkIndividualRules = handleCheckIndividualRules(
    feature.targets,
    rules
  );
  if (!checkIndividualRules.length) return null;
  const individualChanges: DiscardChangesStateData[] = [];
  checkIndividualRules.forEach(item => {
    const isDeleteItem = item.changeType === 'DELETE';
    const isCreateItem = item.changeType === 'CREATE';
    const individualData = getIndividualDiscardChangeData(feature, item.target);
    const isExistedItem = individualChanges.find(
      pre => pre.variation?.id === item.target.variation
    );
    if (isExistedItem) return;
    if (isDeleteItem) {
      const updateItem = getIndividualDiscardChangeItem(
        checkIndividualRules,
        item,
        'CREATE'
      );
      return individualChanges.push({
        ...individualData,
        labelType: updateItem ? 'UPDATE' : 'REMOVE'
      });
    }
    if (isCreateItem) {
      const deleteItem = getIndividualDiscardChangeItem(
        checkIndividualRules,
        item,
        'DELETE'
      );
      return individualChanges.push({
        ...individualData,
        labelType: deleteItem ? 'UPDATE' : 'ADD'
      });
    }
  });
  return individualChanges;
};
