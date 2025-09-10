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
  StrategyType,
  TargetChange,
  UserSegment
} from '@types';
import {
  DefaultRuleSchema,
  RuleSchema,
  StrategySchema,
  TargetingSchema
} from './form-schema';
import {
  DiscardChangesStateData,
  DiscardFeaturePrerequisiteChange,
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
    rolloutStrategy: {
      variations: rolloutStrategy?.variations?.length
        ? rolloutStrategy.variations?.map(item => ({
            ...item,
            weight: item.weight / 1000 || 0
          }))
        : getDefaultRolloutStrategy(feature),
      audience: rolloutStrategy?.audience
        ? {
            percentage: rolloutStrategy?.audience?.percentage || 0,
            defaultVariation: rolloutStrategy?.audience?.defaultVariation || ''
          }
        : {
            percentage: 100,
            defaultVariation: ''
          }
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

const getClauseType = (operator: FeatureRuleClauseOperator) => {
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
      clauses: (clauses || []).map(clause => ({
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

const handleGetStrategy = (
  strategy?: StrategySchema | FeatureRuleStrategy
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
        variations:
          (
            rolloutStrategy as FeatureRuleStrategy['rolloutStrategy']
          )?.variations?.map(item => ({
            ...item,
            weight: item.weight * 1000
          })) || [],
        audience: rolloutStrategy?.audience
          ? {
              percentage: rolloutStrategy?.audience?.percentage || 0,
              defaultVariation:
                rolloutStrategy?.audience?.defaultVariation || ''
            }
          : {
              percentage: 100,
              defaultVariation: ''
            }
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
        variations:
          rolloutStrategy?.variations?.map(item => ({
            ...item,
            weight: item.weight * 1000
          })) || [],
        audience: rolloutStrategy?.audience
          ? {
              percentage: rolloutStrategy?.audience?.percentage || 0,
              defaultVariation:
                rolloutStrategy?.audience?.defaultVariation || ''
            }
          : {
              percentage: 100,
              defaultVariation: ''
            }
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
        label: removedUsers.join(', ')
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
        label: addedUsers.join(', ')
      });
    }
  });

  return individualChanges.length ? individualChanges : null;
};

interface VariationFeatures {
  label: string;
  value: string;
}

const getVariationLabel = (
  variations: VariationFeatures[],
  variationId: string
) => {
  return variations.find(item => item.value === variationId)?.label || '';
};

const formatTimestamp = (ts: number): string => {
  const date = new Date(ts * 1000);
  const pad = (n: number) => String(n).padStart(2, '0');

  const year = date.getFullYear();
  const month = pad(date.getMonth() + 1);
  const day = pad(date.getDate());
  const hours = pad(date.getHours());
  const minutes = pad(date.getMinutes());

  return `${year}/${month}/${day} ${hours}:${minutes}`;
};

const getValueLabel = (
  operator: FeatureRuleClauseOperator,
  segmentUsers: UserSegment[],
  variations: { label: string; value: string }[],
  values: string[]
) => {
  if (operator === FeatureRuleClauseOperator.SEGMENT && segmentUsers) {
    return segmentUsers.find(item => values.includes(item.id))?.name || '';
  }
  if (
    operator === FeatureRuleClauseOperator.BEFORE ||
    operator === FeatureRuleClauseOperator.AFTER
  ) {
    return formatTimestamp(Number(values[0]));
  }
  if (operator === FeatureRuleClauseOperator.FEATURE_FLAG && variations) {
    return getVariationLabel(variations, values[0]);
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

const getLabelClause = (
  segmentUsers: UserSegment[],
  clause: FeatureRuleClause,
  situationOptions: VariationFeatures[],
  variations: VariationFeatures[],
  operatorOptions: VariationFeatures[],
  t: (key: string) => string
) => {
  const operator = get(clause, 'operator', '');
  const attribute = get(clause, 'attribute', '');
  const situationLabel = situationOptions.find(
    s => s.value === clause.type
  )?.label;

  const isUserSegment =
    getClauseType(clause.operator) === RuleClauseType.SEGMENT;
  const isFeatureFlag =
    getClauseType(clause.operator) === RuleClauseType.FEATURE_FLAG;
  const isCompare =
    clause.operator === FeatureRuleClauseOperator.EQUALS ||
    getClauseType(clause.operator) === RuleClauseType.COMPARE;

  const operatorLabel = () => {
    if (isUserSegment) {
      return t('is-included-in');
    }
    if (isFeatureFlag) {
      return '=';
    }
    return (
      operatorOptions?.find(operatorOption => {
        return operatorOption.value === operator;
      })?.label || operator
    );
  };

  const values = !isCompare
    ? getValueLabel(clause.operator, segmentUsers, variations, clause.values)
    : clause.values?.join(', ') || '';

  const labelRule = [situationLabel, attribute, operatorLabel(), values]
    .filter(Boolean)
    .join(' ');

  const labelField = values;
  return { labelRule, labelField };
};

export const handleCheckSegmentRulesDiscardChanges = (
  preRule: FeatureRule,
  segmentUsers: UserSegment[],
  currentRule: FeatureRule,
  situationOptions: VariationFeatures[],
  features: Feature[],
  operatorOptions: VariationFeatures[],
  variationFeatures: FeatureVariation[],
  t: (key: string) => string
): DiscardChangesStateData[] => {
  const changeClause: DiscardChangesStateData[] = [];

  const preClauses = new Map(preRule?.clauses?.map(c => [c.id, c]));
  const currentClauses = new Map(currentRule.clauses.map(c => [c.id, c]));
  const preStrategy = preRule.strategy;
  const currentStrategy = currentRule.strategy;
  const variationFeatureLabelFc = (variationId: string) => {
    return variationFeatures.find(item => item.id === variationId)?.name || '';
  };

  if (preClauses.size) {
    preRule.clauses.forEach(preClause => {
      const currentClause = currentClauses.get(preClause.id);
      const variationFeatures = getVariation(
        features,
        currentClause || preClause
      );
      const preVariationFeatures = getVariation(features, preClause);
      const isUpdate = currentClause
        ? getClauseType(preClause.operator) !== currentClause.type ||
          preClause.operator !== currentClause.operator ||
          preClause.attribute !== currentClause.attribute
        : false;

      const preValues = preClause.values || [];
      const currentValues = currentClause ? currentClause.values || [] : [];
      const isCompare = currentClause
        ? currentClause.operator === FeatureRuleClauseOperator.EQUALS ||
          currentClause.type === RuleClauseType.COMPARE
        : false;
      const removedValue = preValues.filter(
        value => !currentValues.includes(value)
      );
      const addValue = currentValues.filter(
        value => !preValues.includes(value)
      );

      if (!currentClause) {
        changeClause.push({
          label: getLabelClause(
            segmentUsers,
            preClause,
            situationOptions,
            preVariationFeatures || [],
            operatorOptions,
            t
          ).labelRule,
          changeField: 'clause',
          labelType: 'REMOVE',
          variationIndex: 0
        });
        return;
      }

      if (isUpdate) {
        changeClause.push({
          label: getLabelClause(
            segmentUsers,
            currentClause,
            situationOptions,
            variationFeatures || [],
            operatorOptions,
            t
          ).labelRule,
          labelType: 'UPDATE',
          changeField: 'clause',
          variationIndex: 0
        });
        return;
      }

      if (removedValue.length && isCompare) {
        changeClause.push({
          label: getLabelClause(
            segmentUsers,
            currentClause,
            situationOptions,
            variationFeatures || [],
            operatorOptions,
            t
          ).labelRule,
          labelType: 'REMOVE',
          changeField: 'value',
          valueLabel: getValueLabel(
            preClause.operator,
            segmentUsers,
            variationFeatures || [],
            removedValue
          ),
          variationIndex: 0
        });
      }

      if (addValue.length) {
        changeClause.push({
          valueLabel: getValueLabel(
            preClause.operator,
            segmentUsers,
            variationFeatures || [],
            addValue
          ),
          label: getLabelClause(
            segmentUsers,
            currentClause,
            situationOptions,
            variationFeatures || [],
            operatorOptions,
            t
          ).labelRule,
          labelType: isCompare ? 'ADD' : 'UPDATE',
          changeField: isCompare ? 'value' : 'clause',
          variationIndex: 0
        });
      }
    });
  }

  currentRule.clauses.forEach(currentClause => {
    if (!preClauses.has(currentClause.id)) {
      const variationFeatures = getVariation(features, currentClause);
      changeClause.push({
        label: getLabelClause(
          segmentUsers,
          currentClause,
          situationOptions,
          variationFeatures,
          operatorOptions,
          t
        ).labelRule,
        changeField: 'clause',
        labelType: 'ADD',
        variationIndex: 0
      });
    }
  });

  if (currentStrategy.type === StrategyType.FIXED && preStrategy) {
    if (
      preStrategy.fixedStrategy?.variation ===
      currentStrategy.fixedStrategy?.variation
    ) {
      return changeClause;
    }
    const variationFeatureLabel = variationFeatureLabelFc(
      currentStrategy.fixedStrategy.variation
    );
    changeClause.push({
      label: variationFeatureLabel,
      variationPercent: [{ variation: variationFeatureLabel || '' }],
      labelType: 'UPDATE',
      changeField: 'strategy',
      variationIndex: 0
    });
  }
  if (currentStrategy.type === StrategyType.ROLLOUT) {
    const preRolloutStrategy = preStrategy.rolloutStrategy?.variations.map(
      v => ({ ...v, weight: v.weight > 0 ? v.weight / 1000 : 0 })
    );
    if (
      isEqual(preRolloutStrategy, currentStrategy.rolloutStrategy.variations)
    ) {
      return changeClause;
    }
    const variationWeight = currentStrategy.rolloutStrategy.variations
      .map(variation => {
        const variationFeatureLabel = variationFeatureLabelFc(
          variation.variation
        );
        if (variationFeatureLabel) {
          return {
            variation: variationFeatureLabel,
            weight: variation.weight
          };
        }
      })
      .filter(Boolean) as { variation: string; weight: number }[];
    changeClause.push({
      label: '',
      variationPercent: variationWeight,
      labelType: 'UPDATE',
      changeField: 'strategy',
      variationIndex: 0
    });
  }

  return changeClause;
};

export const handleSwapRuleFeature = (
  feature: Feature,
  indexA: number,
  indexB: number
) => {
  let tmp = null;
  const ruleA = feature.rules[indexA];
  const ruleB = feature.rules[indexB];
  tmp = ruleA;
  feature.rules[indexA] = ruleB;
  feature.rules[indexB] = tmp;
  return feature;
};
