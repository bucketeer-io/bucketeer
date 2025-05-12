import { useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { DefaultRuleStrategyType, Feature } from '@types';
import { IconCircleDashed, IconInfo, IconPercentage } from '@icons';
import Card from 'pages/feature-flag-details/elements/card';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import { TargetingSchema } from '../form-schema';
import Strategy from '../segment-rule/strategy';
import { VariationOption } from '../segment-rule/variation';
import { getDefaultRolloutStrategy } from '../utils';

const DefaultRule = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['form']);

  const { watch, setFocus, setValue } = useFormContext<TargetingSchema>();

  const commonName = 'defaultRule';
  const defaultRule = watch(commonName);
  const manualStrategy = watch('defaultRule.manualStrategy');
  const defaultRolloutStrategy = getDefaultRolloutStrategy(feature);

  const variationOptions: VariationOption[] = useMemo(() => {
    const variations = feature.variations.map(item => ({
      label: item.name || item.value,
      value: item.id,
      type: DefaultRuleStrategyType.FIXED,
      variationValue: item.value
    }));
    return [
      ...variations,
      {
        label: t('form:manual-percentage'),
        value: DefaultRuleStrategyType.MANUAL,
        type: DefaultRuleStrategyType.MANUAL,
        icon: IconPercentage
      },
      {
        label: t('common:source-type.progressive-rollout'),
        value: DefaultRuleStrategyType.ROLLOUT,
        type: DefaultRuleStrategyType.ROLLOUT,
        icon: IconCircleDashed
      }
    ];
  }, [feature]);

  const percentageValueCount = useMemo(
    () => manualStrategy?.filter(item => item.weight > 0)?.length || 0,
    [manualStrategy]
  );

  const handleSelectStrategy = useCallback(
    (item: VariationOption) => {
      const { type, value } = item;
      const isFixed = type === DefaultRuleStrategyType.FIXED;
      const isRollout = type === DefaultRuleStrategyType.ROLLOUT;

      setValue(commonName, {
        ...defaultRule,
        type: type as DefaultRuleStrategyType,
        currentOption: value,
        fixedStrategy: {
          variation: isFixed ? value : ''
        },
        rolloutStrategy: isFixed ? [] : defaultRolloutStrategy
      });
      if (!isFixed) {
        let timerId: NodeJS.Timeout | null = null;
        if (timerId) clearTimeout(timerId);
        timerId = setTimeout(
          () =>
            setFocus(
              `${commonName}.${isRollout ? 'rolloutStrategy' : 'manualStrategy'}.0.weight`
            ),
          100
        );
      }
    },
    [feature, variationOptions, commonName]
  );

  return (
    <Card>
      <div className="flex flex-col w-full gap-y-3">
        <div className="flex items-center gap-x-2">
          <p className="typo-para-medium text-gray-700">
            {t('targeting.default-rule')}
          </p>
          <Tooltip
            trigger={
              <div className="flex-center size-fit">
                <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
              </div>
            }
          />
        </div>
        <p className="typo-para-small text-gray-500">
          {t('targeting.default-rule-desc')}
        </p>
      </div>
      <Strategy
        feature={feature}
        rootName={commonName}
        strategyName={'manualStrategy'}
        variationOptions={variationOptions}
        percentageValueCount={percentageValueCount}
        handleSelectStrategy={handleSelectStrategy}
      />
    </Card>
  );
};

export default DefaultRule;
