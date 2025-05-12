import { isEqual } from 'lodash';
import { v4 as uuid } from 'uuid';
import {
  DefaultRuleStrategyType,
  Feature,
  FeaturePrerequisite,
  FeatureRule,
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
import { RuleSchema, TargetingSchema } from './form-schema';
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
) => ({
  ...defaultStrategy,
  rolloutStrategy: defaultStrategy?.rolloutStrategy?.variations?.length
    ? defaultStrategy?.rolloutStrategy?.variations
    : getDefaultRolloutStrategy(feature),
  currentOption:
    defaultStrategy?.type === StrategyType.FIXED
      ? defaultStrategy?.fixedStrategy?.variation
      : StrategyType.ROLLOUT
});

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
) =>
  targets.map(({ variation, users }) => {
    const currentVariation = variations.find(item => item.id === variation);
    return {
      variationId: variation,
      users,
      name: currentVariation?.name || currentVariation?.value || ''
    };
  });

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
      type:
        defaultStrategy?.type === StrategyType.FIXED
          ? defaultStrategy?.fixedStrategy?.variation
          : DefaultRuleStrategyType.MANUAL,
      manualStrategy: getDefaultRolloutStrategy(feature)
    },
    enabled,
    isShowRules: enabled
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

export const handleCheckSegmentRules = (
  featureRules: FeatureRule[],
  segmentRules?: RuleSchema[]
) => {
  const ruleChanges: RuleChange[] = [];
  featureRules.forEach(item => {
    const currentRule = segmentRules?.find(rule => rule.id === item.id);
    if (!currentRule) {
      ruleChanges.push({
        changeType: 'DELETE',
        rule: item
      });
    }

    if (currentRule && !isEqual(currentRule, item)) {
      ruleChanges.push({
        changeType: 'UPDATE',
        rule: {
          id: currentRule.id,
          clauses: currentRule.clauses as FeatureRuleClause[],
          strategy: currentRule.strategy as unknown as FeatureRuleStrategy
        }
      });
    }
  });

  segmentRules?.forEach(item => {
    const currentRule = featureRules?.find(rule => rule.id === item.id);
    if (!currentRule) {
      ruleChanges.push({
        changeType: 'CREATE',
        rule: {
          id: item.id,
          clauses: item.clauses as FeatureRuleClause[],
          strategy: item.strategy as unknown as FeatureRuleStrategy
        }
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

    if (currentTarget && !isEqual(currentTarget, item)) {
      targetChanges.push({
        changeType: 'UPDATE',
        target: {
          users: currentTarget.users,
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
