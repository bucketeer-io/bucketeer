import { v4 as uuid } from 'uuid';
import {
  Feature,
  FeaturePrerequisite,
  FeatureRuleClauseOperator,
  FeatureRuleStrategy,
  FeatureTarget,
  FeatureVariation,
  StrategyType
} from '@types';
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
    defaultStrategy: _defaultStrategy,
    enabled,
    isShowRules: enabled
  };
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
