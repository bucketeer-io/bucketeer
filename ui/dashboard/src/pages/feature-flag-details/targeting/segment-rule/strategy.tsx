import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import {
  DefaultRuleStrategyType,
  Feature,
  RuleStrategyVariation,
  StrategyType
} from '@types';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import { StrategySchema } from '../form-schema';
import PercentageBar from './percentage-bar';
import PercentageInput from './percentage-input';
import { VariationOption } from './variation';

interface Props {
  feature: Feature;
  rootName: string;
  strategyName: 'rolloutStrategy' | 'manualStrategy';
  variationOptions: VariationOption[];
  percentageValueCount: number;
  label?: string;
  isRequired?: boolean;
  handleSelectStrategy: (item: VariationOption) => void;
}

const Strategy = ({
  feature,
  rootName,
  strategyName,
  variationOptions,
  percentageValueCount,
  label,
  isRequired = true,
  handleSelectStrategy
}: Props) => {
  const { t } = useTranslation(['table', 'common', 'form', 'message']);
  const { control, watch, setError, clearErrors } = useFormContext();

  const type = watch(`${rootName}.type`);
  const currentOption = watch(`${rootName}.currentOption`);

  const handleCheckError = (values: StrategySchema['rolloutStrategy']) => {
    const total = values
      ?.map(v => Number(v.weight))
      .reduce((acc: number, current: number) => {
        return acc + (current || 0);
      }, 0);

    if (total !== 100) {
      return setError(`${rootName}.${strategyName}`, {
        message: t('message:validation.should-be-percent')
      });
    }
    clearErrors(`${rootName}.${strategyName}`);
  };

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
            name={`${rootName}.currentOption`}
            render={({ field }) => {
              const option = variationOptions.find(
                item => item.value === currentOption
              );

              return (
                <Form.Item className="flex flex-col flex-1 py-0 w-full">
                  <Form.Control>
                    <DropdownMenu>
                      <div className="flex flex-col gap-y-2 w-full">
                        <DropdownMenuTrigger
                          trigger={
                            <div className="flex items-center gap-x-2 typo-para-medium text-gray-700">
                              {feature.variationType === 'BOOLEAN' &&
                                option?.variationValue && (
                                  <FlagVariationPolygon
                                    index={
                                      option?.variationValue === 'true' ? 0 : 1
                                    }
                                    className="!z-0"
                                  />
                                )}
                              {option?.icon && (
                                <Icon icon={option.icon} size={'sm'} />
                              )}
                              {option?.label || ''}
                            </div>
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
              );
            }}
          />

          {[StrategyType.ROLLOUT, DefaultRuleStrategyType.MANUAL].includes(
            type
          ) && (
            <>
              <Form.Field
                control={control}
                name={`${rootName}.${strategyName}`}
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
                                  variationOptions={variationOptions}
                                  feature={feature}
                                  name={`${rootName}.${strategyName}.${index}.weight`}
                                  variationId={rollout.variation}
                                  handleChangeRolloutWeight={value => {
                                    field.value[index] = {
                                      ...field.value[index],
                                      weight: +value
                                    };
                                    field.onChange(field.value);
                                    handleCheckError(field.value);
                                  }}
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
