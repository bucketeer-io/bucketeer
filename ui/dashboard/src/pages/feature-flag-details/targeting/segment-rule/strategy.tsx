import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { Feature, RuleStrategyVariation, StrategyType } from '@types';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import PercentageBar from './percentage-bar';
import PercentageInput from './percentage-input';
import { VariationOption } from './variation';

interface Props {
  feature: Feature;
  strategyName: string;
  variationOptions: VariationOption[];
  percentageValueCount: number;
  label?: string;
  isRequired?: boolean;
  handleSelectStrategy: (item: VariationOption) => void;
  handleChangeRolloutWeight: (value: number, index: number) => void;
}

const Strategy = ({
  feature,
  strategyName,
  variationOptions,
  percentageValueCount,
  label,
  isRequired = true,
  handleSelectStrategy,
  handleChangeRolloutWeight
}: Props) => {
  const { t } = useTranslation(['table', 'common', 'form']);
  const { control, watch } = useFormContext();

  const type = watch(`${strategyName}.type`);
  const currentOption = watch(`${strategyName}.currentOption`);
  return (
    <div>
      <Form.Label
        required={isRequired}
        className={cn('relative w-fit mb-5 text-gray-700', {
          'mb-2 ml-14': label
        })}
      >
        {label ? label : t('feature-flags.variation')}
        {!label && (
          <Icon
            icon={IconInfo}
            size="xs"
            color="gray-500"
            className="absolute -right-6"
          />
        )}
      </Form.Label>
      <div className="flex w-full gap-x-4">
        <p className="typo-para-small text-gray-600 mt-3 uppercase">
          {t('feature-flags.serve')}
        </p>
        <div className="flex flex-col w-full gap-x-2">
          <Form.Field
            control={control}
            name={`${strategyName}.currentOption`}
            render={({ field }) => (
              <Form.Item className="flex flex-col flex-1 py-0 w-full">
                <Form.Control>
                  <DropdownMenu>
                    <div className="flex flex-col gap-y-2 w-full">
                      <DropdownMenuTrigger
                        label={
                          variationOptions.find(
                            item => item.value === currentOption
                          )?.label || ''
                        }
                        isExpand
                        className="w-full"
                      />
                    </div>
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
            <>
              <Form.Field
                control={control}
                name={`${strategyName}.rolloutStrategy`}
                render={({ field }) => {
                  return (
                    <Form.Item className="flex flex-col w-full gap-y-2">
                      <Form.Control>
                        <>
                          {percentageValueCount > 0 && (
                            <div className="flex items-center w-full p-0.5 border border-gray-400 rounded-full">
                              {field.value?.map(
                                (
                                  item: RuleStrategyVariation,
                                  index: number
                                ) => (
                                  <PercentageBar
                                    key={index}
                                    weight={item.weight}
                                    currentIndex={index}
                                    isRoundedFull={percentageValueCount === 1}
                                  />
                                )
                              )}
                            </div>
                          )}
                          <div className="flex flex-col gap-y-4">
                            {field.value?.map(
                              (
                                rollout: RuleStrategyVariation,
                                index: number
                              ) => (
                                <PercentageInput
                                  key={index}
                                  feature={feature}
                                  name={`${strategyName}.rolloutStrategy.${index}.weight`}
                                  variationId={rollout.variation}
                                  handleChangeRolloutWeight={value =>
                                    handleChangeRolloutWeight(value, index)
                                  }
                                />
                              )
                            )}
                          </div>
                        </>
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  );
                }}
              />
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default Strategy;
