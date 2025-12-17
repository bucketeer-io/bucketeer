import { TFunction } from 'i18next';
import { get } from 'lodash';
import isEqual from 'lodash/isEqual';
import omit from 'lodash/omit';
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
  RuleStrategyVariation,
  StrategyType,
  TargetChange,
  UserSegment
} from '@types';
import { formatLongDateTime } from 'utils/date-time';
import {
  DefaultRuleSchema,
  RuleSchema,
  StrategySchema,
  TargetingSchema
} from './form-schema';
import {
  ClauseLabel,
  DiscardChangesStateData,
  DiscardFeaturePrerequisiteChange,
  IndividualRuleItem,
  PrerequisiteSchema,
  RuleClauseType,
  VariationFeatures,
  VariationPercent
} from './types';

const createAudienceConfig = (audience?: {
  percentage?: number;
  defaultVariation?: string;
}) => ({
  percentage: audience?.percentage || 0,
  defaultVariation: audience?.defaultVariation || ''
});

const getDefaultAudienceConfig = () => ({
  percentage: 100,
  defaultVariation: ''
});

const convertVariationWeights = (
  variations: FeatureRuleStrategy['rolloutStrategy']['variations'],
  factor: number
) =>
  variations?.map(item => ({
    ...item,
    weight: item.weight * factor
  })) || [];

const findVariationWithIndex = <T extends { id: string }>(
  variations: T[],
  targetId: string
): { variation?: T; index: number } => {
  const index = variations.findIndex(v => v?.id === targetId);
  const variation = index !== -1 ? variations[index] : undefined;

  return { variation, index };
};

const createClauseLabelParams = (
  segmentUsers: UserSegment[],
  clause: FeatureRuleClause,
  situationOptions: VariationFeatures[],
  operatorOptions: VariationFeatures[],
  features: Feature[],
  t: TFunction
) => {
  const variationFeatures = getVariation(features, clause);
  return {
    segmentUsers,
    clause,
    situationOptions,
    variationFeatures,
    operatorOptions,
    t
  };
};

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
    rolloutStrategy: {
      variations: rolloutStrategy?.variations?.length
        ? convertVariationWeights(rolloutStrategy.variations, 1 / 1000)
        : getDefaultRolloutStrategy(feature),
      audience: rolloutStrategy?.audience
        ? createAudienceConfig(rolloutStrategy.audience)
        : getDefaultAudienceConfig()
    },
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
    rolloutStrategy: {
      variations: getDefaultRolloutStrategy(feature),
      audience: {
        percentage: 100,
        defaultVariation: ''
      }
    },
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

export const handleCreatePrerequisites = (
  prerequisites: FeaturePrerequisite[]
) =>
  prerequisites.map(({ featureId, variationId }) => ({
    featureId,
    variationId
  }));

export const handleCreateIndividualRules = (
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

export const getClauseType = (operator: FeatureRuleClauseOperator) => {
  const { FEATURE_FLAG, BEFORE, AFTER, SEGMENT } = FeatureRuleClauseOperator;
  if (operator === FEATURE_FLAG) return RuleClauseType.FEATURE_FLAG;
  if ([BEFORE, AFTER].includes(operator)) return RuleClauseType.DATE;
  if (operator === SEGMENT) return RuleClauseType.SEGMENT;
  return RuleClauseType.COMPARE;
};

export const handleCreateSegmentRules = (feature: Feature) => {
  return (feature.rules || [])
    .filter(rule => rule !== undefined && rule !== null)
    .map(({ id, strategy, clauses }) => ({
      id,
      strategy: getDefaultStrategy(feature, strategy),
      clauses: (clauses || [])
        .filter(clause => clause !== undefined && clause !== null)
        .map(clause => ({
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
      type:
        defaultStrategy?.type === StrategyType.FIXED
          ? StrategyType.FIXED
          : StrategyType.MANUAL,
      currentOption:
        defaultStrategy?.type === StrategyType.FIXED
          ? defaultStrategy?.fixedStrategy?.variation
          : StrategyType.MANUAL,
      rolloutStrategy: _defaultStrategy?.rolloutStrategy
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

export const handleGetStrategy = (
  strategy?: StrategySchema | FeatureRuleStrategy,
  factor?: number
): Partial<FeatureRuleStrategy> => {
  const { type, fixedStrategy, rolloutStrategy } = strategy || {};
  if (type === StrategyType.FIXED) {
    return {
      type,
      fixedStrategy: {
        variation: fixedStrategy?.variation || ''
      }
    };
  } else {
    return {
      type,
      rolloutStrategy: {
        variations: convertVariationWeights(
          (rolloutStrategy as FeatureRuleStrategy['rolloutStrategy'])
            ?.variations || [],
          factor ?? 1000
        ),
        audience: rolloutStrategy?.audience
          ? createAudienceConfig(rolloutStrategy.audience)
          : getDefaultAudienceConfig()
      }
    };
  }
};

export const handleGetDefaultRuleStrategy = (
  defaultRule?: DefaultRuleSchema
): Partial<FeatureRuleStrategy> => {
  const { currentOption, rolloutStrategy } = defaultRule || {};
  if (currentOption === StrategyType.MANUAL) {
    return {
      type: StrategyType.ROLLOUT,
      rolloutStrategy: {
        variations: convertVariationWeights(
          rolloutStrategy?.variations || [],
          1000
        ),
        audience: rolloutStrategy?.audience
          ? createAudienceConfig(rolloutStrategy.audience)
          : getDefaultAudienceConfig()
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
      strategy: handleGetStrategy(currentRule?.strategy)
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

export const isEquallyVariations = (
  variations: StrategySchema['rolloutStrategy']['variations'] = []
): boolean => {
  if (variations.length === 0) return false;
  const expectedWeight = 100 / variations.length;

  return variations.every(
    item => Math.abs(item.weight - expectedWeight) < 0.0001
  );
};

const getPrerequisiteDiscardChangeData = (
  item: PrerequisiteChange,
  activeFeatures: Feature[]
) => {
  const currentFeature = activeFeatures.find(
    feature => feature.id === item.prerequisite.featureId
  );
  const { variation, index: variationIndex } = findVariationWithIndex(
    currentFeature?.variations || [],
    item.prerequisite.variationId
  );
  return {
    label: currentFeature?.name || '',
    featureId: currentFeature?.id || '',
    variationIndex,
    variation
  };
};

export const handleCheckPrerequisiteDiscardChanges = (
  prerequisites: PrerequisiteSchema[],
  feature: Feature,
  activeFeatures: Feature[]
) => {
  const diffPrerequisites = handleCheckPrerequisites(
    feature.prerequisites,
    prerequisites
  );
  if (!diffPrerequisites.length) return null;

  const changesByFeature = new Map<string, DiscardFeaturePrerequisiteChange>();

  diffPrerequisites.forEach(change => {
    const featureId = change.prerequisite.featureId;
    if (!featureId) return;

    if (!changesByFeature.has(featureId)) {
      changesByFeature.set(featureId, {});
    }

    const changeSet = changesByFeature.get(featureId)!;
    if (change.changeType === 'DELETE') changeSet.deleted = change;
    if (change.changeType === 'CREATE') changeSet.created = change;
  });

  const normalizedChanges: DiscardChangesStateData[] = [];
  changesByFeature.forEach(({ deleted, created }) => {
    if (deleted && created) {
      normalizedChanges.push({
        ...getPrerequisiteDiscardChangeData(created, activeFeatures),
        labelType: 'UPDATE'
      });
    } else if (deleted) {
      normalizedChanges.push({
        ...getPrerequisiteDiscardChangeData(deleted, activeFeatures),
        labelType: 'REMOVE'
      });
    } else if (created) {
      normalizedChanges.push({
        ...getPrerequisiteDiscardChangeData(created, activeFeatures),
        labelType: 'ADD'
      });
    }
  });

  return normalizedChanges;
};

const getIndividualDiscardChangeData = (
  feature: Feature,
  individual: FeatureTarget
) => {
  const { variation, index: variationIndex } = findVariationWithIndex(
    feature?.variations || [],
    individual.variation
  );

  return {
    label: individual.users?.join(', ') || '',
    variationIndex,
    variation
  };
};

export const handleCheckIndividualDiscardChanges = (
  feature: Feature,
  rules: IndividualRuleItem[]
): DiscardChangesStateData[] | null => {
  const individualChanges: DiscardChangesStateData[] = [];

  feature.targets.forEach((preIndividual, index) => {
    const currentIndividual = rules.find(
      r => r.variationId === preIndividual.variation
    );

    const preUsers = preIndividual.users || [];
    const currentUsers = currentIndividual?.users || [];

    const removedUsers = preUsers.filter(u => !currentUsers.includes(u));
    if (removedUsers.length) {
      individualChanges.push({
        featureId: feature.id,
        variationIndex: index,
        variation: getIndividualDiscardChangeData(feature, preIndividual)
          .variation,
        labelType: 'REMOVE',
        label: '',
        groupLabel: removedUsers
      });
    }

    const addedUsers = currentUsers.filter(u => !preUsers.includes(u));
    if (addedUsers.length) {
      individualChanges.push({
        featureId: feature.id,
        variationIndex: index,
        variation: getIndividualDiscardChangeData(feature, preIndividual)
          .variation,
        labelType: 'ADD',
        label: '',
        groupLabel: addedUsers
      });
    }
  });

  return individualChanges.length ? individualChanges : null;
};

const getVariationLabel = (
  variations: VariationFeatures[],
  variationId: string
) => {
  const { variation, index } = findVariationWithIndex(
    variations.map(v => ({ id: v.value, ...v })),
    variationId
  );

  if (index === -1) {
    return { label: '', index: -1 };
  }

  return { label: variation?.label || '', index };
};

const getValueLabel = (
  operator: FeatureRuleClauseOperator,
  segmentUsers: UserSegment[],
  variations: { label: string; value: string }[],
  values: string[]
) => {
  const { SEGMENT, FEATURE_FLAG, BEFORE, AFTER } = FeatureRuleClauseOperator;
  if (operator === SEGMENT && segmentUsers) {
    return segmentUsers.find(item => values.includes(item.id))?.name || '';
  }
  if ([BEFORE, AFTER].includes(operator as FeatureRuleClauseOperator)) {
    return values[0]
      ? formatLongDateTime({
          value: values[0],
          overrideOptions: {
            month: 'long',
            hour: '2-digit',
            minute: '2-digit',
            hourCycle: 'h23'
          }
        })
      : '';
  }
  if (operator === FEATURE_FLAG && variations) {
    return getVariationLabel(variations, values[0]).label;
  }
  return values.join(', ');
};

const getVariation = (features: Feature[], clause: FeatureRuleClause) => {
  const featureId =
    clause.operator === FeatureRuleClauseOperator.FEATURE_FLAG
      ? clause.attribute
      : '';
  const variationFeatures = features
    ?.find(item => item.id === featureId)
    ?.variations?.map(v => ({
      label: v.name || v.value,
      value: v.id
    }));
  return variationFeatures || [];
};

const normalizedWeightVariations = (variations: RuleStrategyVariation[]) => {
  return variations.map(v => ({
    ...v,
    weight: v.weight >= 0 && v.weight <= 100 ? v.weight : v.weight / 1000
  }));
};

export const getClauseLabel = (
  segmentUsers: UserSegment[],
  clause: FeatureRuleClause,
  situationOptions: VariationFeatures[],
  variations: VariationFeatures[],
  operatorOptions: VariationFeatures[],
  t: TFunction<['common', 'form', 'message'], undefined>
): ClauseLabel => {
  const operator = get(clause, 'operator', '');

  const attribute = get(clause, 'attribute', '');
  const situationText = situationOptions.find(
    s => s.value === clause.type
  )?.label;

  const isUserSegment =
    getClauseType(clause.operator) === RuleClauseType.SEGMENT;
  const isFeatureFlag =
    getClauseType(clause.operator) === RuleClauseType.FEATURE_FLAG;
  const isCompare =
    clause.operator === FeatureRuleClauseOperator.EQUALS ||
    getClauseType(clause.operator) === RuleClauseType.COMPARE;

  const operatorText = isUserSegment
    ? t('is-included-in')
    : isFeatureFlag
      ? '='
      : operatorOptions.find(opt => opt.value === operator)?.label || operator;

  const values = !isCompare
    ? getValueLabel(clause.operator, segmentUsers, variations, clause.values)
    : clause.values?.join(', ') || '';

  const fullLabel = [situationText, attribute, operatorText, values]
    .filter(Boolean)
    .join(' ');

  const valueLabel = values;

  return { fullLabel, valueLabel };
};

const getVariationInfo = (
  variationFeatures: FeatureVariation[],
  variationId: string
) => {
  const { variation, index: indexVariation } = findVariationWithIndex(
    variationFeatures,
    variationId
  );

  return indexVariation === -1
    ? { variationLabel: '', indexVariation: -1, variationId: '' }
    : {
        variationLabel: variation?.name || '',
        indexVariation,
        variationId: variation?.id || ''
      };
};

const getStrategyVariationWeight = (
  strategy: FeatureRuleStrategy,
  variationFeatures: FeatureVariation[]
): VariationPercent[] => {
  const { ROLLOUT, FIXED, MANUAL } = StrategyType;
  if (strategy.type === FIXED) {
    const { variationLabel, indexVariation, variationId } = getVariationInfo(
      variationFeatures,
      strategy.fixedStrategy.variation
    );

    return [
      {
        variationId: variationId,
        variation: variationLabel,
        weight: null,
        variationIndex: indexVariation
      }
    ];
  }
  if ([ROLLOUT, MANUAL].includes(strategy.type)) {
    return strategy.rolloutStrategy.variations
      .map(variation => {
        const { variationLabel, indexVariation, variationId } =
          getVariationInfo(variationFeatures, variation.variation);

        if (variationLabel) {
          return {
            variationId: variationId,
            variation: variationLabel,
            weight: variation.weight,
            variationIndex: indexVariation
          };
        }
      })
      .filter(Boolean) as VariationPercent[];
  }

  return [];
};

export const getClauseLabelsFromRule = (
  rule: FeatureRule,
  features: Feature[],
  segmentUsers: UserSegment[],
  situationOptions: VariationFeatures[],
  operatorOptions: VariationFeatures[],
  t: TFunction<['common', 'form', 'message'], undefined>
): string[] =>
  rule.clauses.map(clause => {
    const variationFeaturesLabel = getVariation(features, clause);
    const { fullLabel } = getClauseLabel(
      segmentUsers,
      clause,
      situationOptions,
      variationFeaturesLabel,
      operatorOptions,
      t
    );

    return `${fullLabel}`;
  });

export const getFeatureRuleLabels = (
  originFeature: Feature,
  currentFeature: Feature,
  features: Feature[],
  segmentUsers: UserSegment[],
  situationOptions: VariationFeatures[],
  operatorOptions: VariationFeatures[],
  t: TFunction<['common', 'form', 'message'], undefined>,
  variationFeatures: FeatureVariation[]
) => {
  const originFeatueRuleIds = originFeature.rules.map(r => r.id);
  const filteredRules = currentFeature.rules;

  const variations: VariationPercent[][] = [];
  const labels: string[][] = [];
  const isNewRules: boolean[] = [];

  filteredRules.forEach(rule => {
    isNewRules.push(!originFeatueRuleIds.includes(rule.id));
    variations.push(
      getStrategyVariationWeight(
        rule.strategy,
        variationFeatures
      ) as VariationPercent[]
    );
    labels.push(
      getClauseLabelsFromRule(
        rule,
        features,
        segmentUsers,
        situationOptions,
        operatorOptions,
        t
      )
    );
  });

  return { variations, labels, isNewRules };
};

export const handleSwapRuleFeature = (
  feature: Feature,
  indexA: number,
  indexB: number
): Feature => {
  const newRules = [...feature.rules];
  [newRules[indexA], newRules[indexB]] = [newRules[indexB], newRules[indexA]];

  return { ...feature, rules: newRules };
};

const getAudienceChangeData = (
  strategy: FeatureRuleStrategy,
  variationFeatures: FeatureVariation[]
) => {
  const variationPercents = getStrategyVariationWeight(
    strategy,
    variationFeatures
  );

  const defaultVariation = strategy.rolloutStrategy?.audience?.defaultVariation;
  const includedPercent = strategy.rolloutStrategy?.audience?.percentage;

  const audienceExcluded = variationPercents.find(
    v => v.variationId === defaultVariation
  );

  const audienceIncluded = variationPercents.filter(
    v => v.variationId !== defaultVariation
  );

  return {
    variationPercents,
    audienceExcluded: {
      ...audienceExcluded,
      weight: Number(100 - (includedPercent ?? 0))
    } as VariationPercent,
    audienceIncluded: audienceIncluded
  };
};

const diffClauses = (
  preRule: FeatureRule,
  currentRule: FeatureRule,
  features: Feature[],
  segmentUsers: UserSegment[],
  situationOptions: VariationFeatures[],
  operatorOptions: VariationFeatures[],
  t: TFunction
): {
  changes: DiscardChangesStateData[];
  action: 'new-rule' | 'edit-rule' | undefined;
} => {
  const changes: DiscardChangesStateData[] = [];
  let actionChange = undefined;

  const preClauses = new Map(preRule.clauses.map(c => [c.id, c]));
  const currentClauses = new Map(currentRule.clauses.map(c => [c.id, c]));

  preRule.clauses.forEach(preClause => {
    const currentClause = currentClauses.get(preClause.id);

    if (!currentClause) {
      const preParams = createClauseLabelParams(
        segmentUsers,
        preClause,
        situationOptions,
        operatorOptions,
        features,
        t
      );
      changes.push({
        label: getClauseLabel(
          preParams.segmentUsers,
          preParams.clause,
          preParams.situationOptions,
          preParams.variationFeatures,
          preParams.operatorOptions,
          preParams.t
        ).fullLabel,
        changeType: 'clause',
        labelType: 'REMOVE',
        variationIndex: 0
      });
      return;
    }

    const isUpdate =
      getClauseType(preClause.operator) !==
        getClauseType(currentClause.operator) ||
      preClause.operator !== currentClause.operator ||
      preClause.attribute !== currentClause.attribute;

    if (isUpdate) {
      const currentParams = createClauseLabelParams(
        segmentUsers,
        currentClause,
        situationOptions,
        operatorOptions,
        features,
        t
      );
      changes.push({
        label: getClauseLabel(
          currentParams.segmentUsers,
          currentParams.clause,
          currentParams.situationOptions,
          currentParams.variationFeatures,
          currentParams.operatorOptions,
          currentParams.t
        ).fullLabel,
        labelType: 'UPDATE',
        changeType: 'clause',
        variationIndex: 0
      });
      return;
    }

    const preValues = preClause.values || [];
    const currentValues = currentClause.values || [];

    const isCompare =
      currentClause.operator === FeatureRuleClauseOperator.EQUALS ||
      currentClause.type === RuleClauseType.COMPARE;

    const removedValue = preValues.filter(v => !currentValues.includes(v));
    const addValue = currentValues.filter(v => !preValues.includes(v));

    if (removedValue.length && isCompare) {
      actionChange = 'edit-rule';
      const currentParams = createClauseLabelParams(
        segmentUsers,
        currentClause,
        situationOptions,
        operatorOptions,
        features,
        t
      );
      changes.push({
        label: getClauseLabel(
          currentParams.segmentUsers,
          currentParams.clause,
          currentParams.situationOptions,
          currentParams.variationFeatures,
          currentParams.operatorOptions,
          currentParams.t
        ).fullLabel,
        labelType: 'REMOVE',
        changeType: 'value',
        valueLabel: getValueLabel(
          preClause.operator,
          segmentUsers,
          currentParams.variationFeatures,
          removedValue
        ),
        variationIndex: 0
      });
    }

    if (addValue.length) {
      const currentParams = createClauseLabelParams(
        segmentUsers,
        currentClause,
        situationOptions,
        operatorOptions,
        features,
        t
      );
      changes.push({
        valueLabel: getValueLabel(
          preClause.operator,
          segmentUsers,
          currentParams.variationFeatures,
          addValue
        ),
        label: getClauseLabel(
          currentParams.segmentUsers,
          currentParams.clause,
          currentParams.situationOptions,
          currentParams.variationFeatures,
          currentParams.operatorOptions,
          currentParams.t
        ).fullLabel,
        labelType: isCompare ? 'ADD' : 'UPDATE',
        changeType: isCompare ? 'value' : 'clause',
        variationIndex: 0
      });
    }
  });

  currentRule.clauses.forEach(currentClause => {
    if (!preClauses.has(currentClause.id)) {
      const currentParams = createClauseLabelParams(
        segmentUsers,
        currentClause,
        situationOptions,
        operatorOptions,
        features,
        t
      );
      changes.push({
        label: getClauseLabel(
          currentParams.segmentUsers,
          currentParams.clause,
          currentParams.situationOptions,
          currentParams.variationFeatures,
          currentParams.operatorOptions,
          currentParams.t
        ).fullLabel,
        changeType: 'clause',
        labelType: 'ADD',
        variationIndex: 0
      });
    }
  });

  return { changes, action: actionChange };
};

const diffStrategy = (
  preStrategy: FeatureRuleStrategy,
  currentStrategy: FeatureRuleStrategy,
  variationFeatures: FeatureVariation[],
  isDefaultRule?: boolean
): DiscardChangesStateData[] => {
  const changes: DiscardChangesStateData[] = [];

  if (currentStrategy.type === StrategyType.FIXED) {
    if (
      preStrategy?.fixedStrategy?.variation !==
      currentStrategy.fixedStrategy?.variation
    ) {
      changes.push({
        label: '',
        variationPercent: getStrategyVariationWeight(
          currentStrategy,
          variationFeatures
        ),
        labelType: 'UPDATE',
        changeType: isDefaultRule ? 'default-strategy' : 'strategy',
        variationIndex: 0
      });
    }
  }

  if (
    currentStrategy.type === StrategyType.ROLLOUT ||
    currentStrategy.type === StrategyType.MANUAL
  ) {
    const { variationPercents, audienceExcluded, audienceIncluded } =
      getAudienceChangeData(currentStrategy, variationFeatures);
    const preRollout = normalizedWeightVariations(
      preStrategy.rolloutStrategy?.variations || []
    );

    if (
      !isEqual(
        preStrategy.rolloutStrategy?.audience,
        currentStrategy.rolloutStrategy?.audience
      )
    ) {
      changes.push({
        label: '',
        audienceExcluded,
        audienceIncluded,
        labelType: 'UPDATE',
        changeType: isDefaultRule ? 'default-audience' : 'audience',
        variationIndex: 0
      });
    }

    if (!isEqual(preRollout, currentStrategy.rolloutStrategy.variations)) {
      changes.push({
        label: '',
        variationPercent: variationPercents,
        labelType: 'UPDATE',
        changeType: isDefaultRule ? 'default-strategy' : 'strategy',
        variationIndex: 0
      });
    }
  }

  return changes;
};

export const checkDefaultRuleDiscardChanges = (
  preStrategy: FeatureRuleStrategy,
  currentStrategy: FeatureRuleStrategy,
  variationFeatures: FeatureVariation[]
) => {
  return diffStrategy(preStrategy, currentStrategy, variationFeatures, true);
};

export const handleCheckSegmentRulesDiscardChanges = (
  preRule: FeatureRule | null,
  segmentUsers: UserSegment[],
  currentRule: FeatureRule,
  situationOptions: VariationFeatures[],
  features: Feature[],
  operatorOptions: VariationFeatures[],
  variationFeatures: FeatureVariation[],
  t: TFunction
): {
  changes: DiscardChangesStateData[];
  action: 'new-rule' | 'edit-rule' | undefined;
} => {
  let actionChange = undefined;
  const { variationPercents, audienceExcluded, audienceIncluded } =
    getAudienceChangeData(currentRule.strategy, variationFeatures);

  if (!preRule) {
    const isAddience =
      currentRule.strategy.type !== StrategyType.FIXED
        ? {
            label: '',
            audienceExcluded: audienceExcluded,
            audienceIncluded: audienceIncluded,
            labelType: 'ADD',
            changeType: 'audience',
            variationIndex: 0
          }
        : null;
    return {
      changes: [
        {
          label: getClauseLabelsFromRule(
            currentRule,
            features,
            segmentUsers,
            situationOptions,
            operatorOptions,
            t
          ).join(` ${t('common:and')} `),
          variationPercent: variationPercents,
          changeType: 'new-rule',
          labelType: 'ADD',
          variationIndex: 0
        },

        isAddience as DiscardChangesStateData
      ],
      action: 'new-rule'
    };
  }

  const { changes: clauseChanges, action } = diffClauses(
    preRule,
    currentRule,
    features,
    segmentUsers,
    situationOptions,
    operatorOptions,
    t
  );
  const strategyChanges = diffStrategy(
    preRule.strategy,
    currentRule.strategy,
    variationFeatures
  );

  if (clauseChanges.length || strategyChanges.length) {
    actionChange = action === 'new-rule' ? 'new-rule' : 'edit-rule';
  }

  return {
    changes: [...clauseChanges, ...strategyChanges],
    action: actionChange as 'new-rule' | 'edit-rule' | undefined
  };
};

export const handleCheckRuleDeleted = (
  currentRules: FeatureRule[],
  oldRules: FeatureRule[],
  features: Feature[],
  segmentUsers: UserSegment[],
  situationOptions: VariationFeatures[],
  operatorOptions: VariationFeatures[],
  variationFeatures: FeatureVariation[],
  t: TFunction
) => {
  let deletes: DiscardChangesStateData[] = [];
  const deletedRules: FeatureRule[] = [];
  oldRules.forEach(oldRule => {
    const isExist = currentRules.find(
      currentRule => currentRule.id === oldRule.id
    );
    if (!isExist) {
      deletedRules.push(oldRule);
    }
  });
  deletes = deletedRules.map(deletedRule => {
    const clauseLabels = getClauseLabelsFromRule(
      deletedRule,
      features,
      segmentUsers,
      situationOptions,
      operatorOptions,
      t
    ).join(` ${t('common:and')} `);
    const normalizedWeight = {
      ...deletedRule.strategy,
      rolloutStrategy: {
        ...deletedRule.strategy?.rolloutStrategy,
        variations: normalizedWeightVariations(
          deletedRule.strategy?.rolloutStrategy?.variations ?? []
        )
      }
    };
    const variationPercents = getStrategyVariationWeight(
      normalizedWeight,
      variationFeatures
    );
    return {
      label: clauseLabels,
      variationPercent: variationPercents,
      changeType: 'deleted-rule' as const,
      labelType: 'REMOVE' as const,
      variationIndex: 0
    };
  });

  return deletes;
};
