import { FunctionComponent, useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import {
  DefaultRuleStrategyType,
  Feature,
  RuleStrategyVariation,
  StrategyType
} from '@types';
import { IconPercentage } from '@icons';
import { RuleSchema } from '../form-schema';
import Strategy from './strategy';

export interface VariationOption {
  label: string;
  value: string;
  type: StrategyType | DefaultRuleStrategyType;
  variationValue?: string;
  icon?: FunctionComponent;
}

const SegmentVariation = ({
  feature,
  defaultRolloutStrategy,
  segmentIndex,
  segmentRules
}: {
  feature: Feature;
  defaultRolloutStrategy: RuleStrategyVariation[];
  segmentIndex: number;
  segmentRules: RuleSchema[];
}) => {
  const { t } = useTranslation(['table', 'common', 'form']);

  const methods = useFormContext();
  const { watch, setValue, setFocus } = methods;
  const commonName = `segmentRules.${segmentIndex}.strategy`;
  const rolloutStrategy: RuleStrategyVariation[] = watch(
    `${commonName}.rolloutStrategy`
  );

  const percentageValueCount = useMemo(
    () => rolloutStrategy?.filter(item => item.weight > 0)?.length || 0,
    [rolloutStrategy]
  );
  const variationOptions: VariationOption[] = useMemo(() => {
    const variations = feature.variations.map(item => ({
      label: item.name || item.value,
      value: item.id,
      type: StrategyType.FIXED,
      variationValue: item.value
    }));
    return [
      ...variations,
      {
        label: t('form:manual-percentage'),
        value: StrategyType.ROLLOUT,
        type: StrategyType.ROLLOUT,
        icon: IconPercentage
      }
    ];
  }, [feature]);

  const handleSelectStrategy = useCallback(
    (item: VariationOption) => {
      const { type, value } = item;
      const isFixed = type === StrategyType.FIXED;
      segmentRules[segmentIndex] = {
        ...segmentRules[segmentIndex],
        strategy: {
          ...segmentRules[segmentIndex].strategy,
          type: type as StrategyType,
          currentOption: value,
          fixedStrategy: {
            variation: isFixed ? value : ''
          },
          rolloutStrategy: isFixed ? [] : defaultRolloutStrategy
        }
      };
      setValue(commonName, segmentRules[segmentIndex].strategy);
      if (!isFixed) {
        let timerId: NodeJS.Timeout | null = null;
        if (timerId) clearTimeout(timerId);
        timerId = setTimeout(
          () => setFocus(`${commonName}.rolloutStrategy.0.weight`),
          100
        );
      }
    },
    [
      feature,
      variationOptions,
      defaultRolloutStrategy,
      segmentRules,
      segmentIndex,
      commonName
    ]
  );

  return (
    <Strategy
      feature={feature}
      rootName={commonName}
      strategyName="rolloutStrategy"
      percentageValueCount={percentageValueCount}
      variationOptions={variationOptions}
      handleSelectStrategy={handleSelectStrategy}
    />
  );
};

export default SegmentVariation;
