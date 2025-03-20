import { useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { Feature, StrategyType } from '@types';
import { IconInfo } from '@icons';
import Icon from 'components/icon';
import Card from '../elements/card';
import { StrategySchema } from './form-schema';
import Strategy from './target-segment-rule/strategy';
import { createVariationLabel } from './utils';

export interface VariationOption {
  label: string;
  value: string;
  type: StrategyType;
}

const DefaultStrategy = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['table', 'common', 'form']);

  const methods = useFormContext();
  const { watch, setValue, setFocus, trigger } = methods;
  const commonName = `defaultStrategy`;
  const defaultStrategy: StrategySchema = watch(commonName);
  const rolloutStrategy = useMemo(
    () => defaultStrategy?.rolloutStrategy || [],
    [defaultStrategy]
  );
  const isFixed = useMemo(
    () => defaultStrategy?.type === StrategyType.FIXED,
    [defaultStrategy]
  );

  const percentageValueCount = useMemo(
    () => rolloutStrategy?.filter(item => item?.weight > 0)?.length || 0,
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
      const _defaultStrategy = {
        ...defaultStrategy,
        type,
        currentOption: value,
        fixedStrategy: {
          variation: isFixed ? value : defaultStrategy.fixedStrategy?.variation
        },
        rolloutStrategy: defaultStrategy.rolloutStrategy || []
      };
      setValue(commonName, _defaultStrategy);
      if (!isFixed) {
        let timerId: NodeJS.Timeout | null = null;
        if (timerId) clearTimeout(timerId);
        timerId = setTimeout(
          () => setFocus(`${commonName}.rolloutStrategy.0.weight`),
          100
        );
      }
      // return onChangeTargetSegmentRules([...targetSegmentRules]);
    },
    [feature, variationOptions, defaultStrategy, commonName]
  );

  const handleChangeRolloutWeight = useCallback(
    (weight: number, itemIndex: number) => {
      const newRollout = defaultStrategy.rolloutStrategy?.map((item, index) => {
        if (index === itemIndex) {
          return {
            ...item,
            weight
          };
        }
        return item;
      });
      const _defaultStrategy = {
        ...defaultStrategy,

        rolloutStrategy: newRollout
      };
      setValue(commonName, _defaultStrategy);
      trigger(commonName);
    },
    [defaultStrategy, commonName]
  );

  return (
    <Card>
      {isFixed && (
        <div>
          <div className="flex items-center gap-x-2">
            <p className="typo-para-medium leading-4 text-gray-700">
              {t('feature-flags.default-rule')}
              <span className="text-accent-red-400">*</span>
            </p>
            <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
          </div>
          <p className="typo-para-small text-gray-600 mt-4">
            {t('feature-flags.default-rule-desc')}
          </p>
        </div>
      )}
      <Strategy
        label={isFixed ? t('feature-flags.variation') : t('form:rollout')}
        isRequired={!isFixed}
        feature={feature}
        strategyName={commonName}
        percentageValueCount={percentageValueCount}
        variationOptions={variationOptions}
        handleSelectStrategy={handleSelectStrategy}
        handleChangeRolloutWeight={handleChangeRolloutWeight}
      />
    </Card>
  );
};

export default DefaultStrategy;
