import { useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { Feature, RuleStrategyVariation, StrategyType } from '@types';
import { RuleSchema } from '../form-schema';
import { createVariationLabel } from '../utils';
import Strategy from './strategy';

export interface VariationOption {
  label: string;
  value: string;
  type: StrategyType;
}

const SegmentVariation = ({
  feature,
  defaultRolloutStrategy,
  segmentIndex,
  targetSegmentRules,
  onChangeTargetSegmentRules
}: {
  feature: Feature;
  defaultRolloutStrategy: RuleStrategyVariation[];
  segmentIndex: number;
  targetSegmentRules: RuleSchema[];
  onChangeTargetSegmentRules: (value: RuleSchema[]) => void;
}) => {
  const { t } = useTranslation(['table', 'common', 'form']);

  const methods = useFormContext();
  const { watch, setValue, setFocus, trigger } = methods;
  const commonName = `rules.${segmentIndex}.strategy`;
  const rolloutStrategy: RuleStrategyVariation[] = watch(
    `${commonName}.rolloutStrategy`
  );

  const percentageValueCount = useMemo(
    () => rolloutStrategy?.filter(item => item.weight > 0)?.length || 0,
    [rolloutStrategy]
  );
  const variationOptions: VariationOption[] = useMemo(() => {
    const variations = feature.variations.map(item => ({
      label: createVariationLabel(item),
      value: item.id,
      type: StrategyType.FIXED
    }));
    return [
      ...variations,
      {
        label: t('form:select-rollout-percentage'),
        value: StrategyType.ROLLOUT,
        type: StrategyType.ROLLOUT
      }
    ];
  }, [feature]);

  const handleSelectStrategy = useCallback(
    (item: VariationOption) => {
      const { type, value } = item;
      const isFixed = type === StrategyType.FIXED;
      targetSegmentRules[segmentIndex] = {
        ...targetSegmentRules[segmentIndex],
        strategy: {
          ...targetSegmentRules[segmentIndex].strategy,
          type,
          currentOption: value,
          fixedStrategy: {
            variation: isFixed ? value : ''
          },
          rolloutStrategy: isFixed ? [] : defaultRolloutStrategy
        }
      };
      setValue(commonName, targetSegmentRules[segmentIndex].strategy);
      if (!isFixed) {
        let timerId: NodeJS.Timeout | null = null;
        if (timerId) clearTimeout(timerId);
        timerId = setTimeout(
          () => setFocus(`${commonName}.rolloutStrategy.0.weight`),
          100
        );
      }
      return onChangeTargetSegmentRules([...targetSegmentRules]);
    },
    [
      feature,
      variationOptions,
      defaultRolloutStrategy,
      targetSegmentRules,
      segmentIndex,
      commonName
    ]
  );

  const handleChangeRolloutWeight = useCallback(
    (weight: number, itemIndex: number) => {
      const newRollout = targetSegmentRules[
        segmentIndex
      ]?.strategy?.rolloutStrategy?.map((item, index) => {
        if (index === itemIndex) {
          return {
            ...item,
            weight
          };
        }
        return item;
      });
      targetSegmentRules[segmentIndex] = {
        ...targetSegmentRules[segmentIndex],
        strategy: {
          ...targetSegmentRules[segmentIndex].strategy,
          rolloutStrategy: newRollout
        }
      };
      setValue(commonName, targetSegmentRules[segmentIndex].strategy);
      trigger(commonName);
      return onChangeTargetSegmentRules([...targetSegmentRules]);
    },
    [targetSegmentRules, commonName, segmentIndex]
  );

  return (
    <Strategy
      feature={feature}
      strategyName={commonName}
      percentageValueCount={percentageValueCount}
      variationOptions={variationOptions}
      handleSelectStrategy={handleSelectStrategy}
      handleChangeRolloutWeight={handleChangeRolloutWeight}
    />
  );
};

export default SegmentVariation;
