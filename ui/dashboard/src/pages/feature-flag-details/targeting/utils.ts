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
import {
  DefaultRuleSchema,
  RuleSchema,
  StrategySchema,
  TargetingSchema
} from './form-schema';
import { IndividualRuleItem, RuleClauseType } from './types';

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
  const { prerequisites, targets, variations, defaultStrategy, enabled } =
    feature || {};
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
    isShowRules: enabled,
    comment: '',
    requireComment: false,
    resetSampling: false,
    scheduleType: enabled ? 'DISABLE' : 'ENABLE',
    scheduleAt: String(Math.floor((new Date().getTime() + 3600000) / 1000))
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
    const currentTarget = individualRules?.find(
      rule => rule.variationId === item.variation
    );
    if (!currentTarget) {
      targetChanges.push({
        changeType: 'DELETE',
        target: item
      });
    }

    const targetObj = {
      users: currentTarget?.users,
      variation: currentTarget?.variationId
    };

    if (currentTarget && !isEqual(targetObj, item)) {
      targetChanges.push({
        changeType: currentTarget.users.length ? 'UPDATE' : 'DELETE',
        target: {
          users: currentTarget.users.length ? currentTarget.users : item.users,
          variation: currentTarget.variationId
        }
      });
    }
  });

  individualRules?.forEach(item => {
    const currentTarget = featureTargets?.find(
      target => target.variation === item.variationId
    );
    if (!currentTarget) {
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
