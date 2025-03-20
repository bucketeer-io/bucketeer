import { useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { Feature, RuleStrategyVariation, StrategyType } from '@types';
import { IconInfo } from '@icons';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import { RuleSchema } from '../form-schema';
import { createVariationLabel } from '../utils';

interface VariationOption {
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
  const { control, watch, setValue } = methods;
  const commonName = `rules.${segmentIndex}.strategy`;
  const currentOption = watch(`${commonName}.currentOption`);
  const type = watch(`${commonName}.type`);

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

  return (
    <div>
      <Form.Label required className="relative w-fit mb-5">
        {t('feature-flags.variation')}
        <Icon
          icon={IconInfo}
          size="xs"
          color="gray-500"
          className="absolute -right-6"
        />
      </Form.Label>
      <div className="flex w-full gap-x-4">
        <p className="typo-para-small text-gray-600 mt-3">
          {t('feature-flags.serve')}
        </p>
        <div className="flex flex-col w-full gap-x-2">
          <Form.Field
            control={control}
            name={`${commonName}.currentOption`}
            render={({ field }) => (
              <Form.Item className="flex flex-col flex-1 py-0 w-full">
                <Form.Control>
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      label={
                        variationOptions.find(
                          item => item.value === currentOption
                        )?.label || ''
                      }
                      isExpand
                      className="w-full"
                    />
                    <DropdownMenuContent align="start">
                      {variationOptions.map((item, index) => (
                        <DropdownMenuItem
                          {...field}
                          key={index}
                          label={item.label}
                          value={item.value}
                          onSelectOption={() => handleSelectStrategy(item)}
                        />
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />

          {type === StrategyType.ROLLOUT && (
            <Form.Field
              control={control}
              name={`${commonName}.currentOption`}
              render={() => (
                <Form.Item className="flex flex-col w-full gap-y-2">
                  <Form.Control>
                    <div className="flex items-center min-w-full"></div>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          )}
        </div>
      </div>
    </div>
  );
};

export default SegmentVariation;
