import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { RuleStrategyVariation, StrategyType } from '@types';
import { cn } from 'utils/style';
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
  rootName: string;
  strategyName: 'rolloutStrategy' | 'manualStrategy';
  variationOptions: VariationOption[];
  percentageValueCount: number;
  label?: string;
  isRequired?: boolean;
  isDisabled?: boolean;
  handleSelectStrategy: (item: VariationOption) => void;
}

const Strategy = ({
  rootName,
  strategyName,
  variationOptions,
  percentageValueCount,
  label,
  isRequired = true,
  isDisabled,
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

  const isShowPercentage =
    (type === StrategyType.ROLLOUT && strategyName === 'rolloutStrategy') ||
    currentOption === StrategyType.MANUAL;
  return (
    <div>
      {label && (
        <Form.Label
          required={isRequired}
          className={cn('relative w-fit mb-2 ml-14 text-gray-700')}
        >
          {label}
        </Form.Label>
      )}
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
                              {option?.icon && (
                                <Icon icon={option.icon} size={'sm'} />
                              )}
                              {option?.label || ''}
                            </div>
                          }
                          isExpand
                          disabled={isDisabled}
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
                            icon={item?.icon}
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

          <div className={isShowPercentage ? 'flex' : 'hidden'}>
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
                              (item: RuleStrategyVariation, index: number) => (
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
                            (rollout: RuleStrategyVariation, index: number) => (
                              <PercentageInput
                                key={index}
                                isDisabled={isDisabled}
                                variationOptions={variationOptions}
                                name={`${rootName}.${strategyName}.${index}.weight`}
                                variationId={rollout.variation}
                                handleChangeRolloutWeight={value => {
                                  field.value[index] = {
                                    ...field.value[index],
                                    weight: +value
                                  };
                                  field.onChange(field.value, {
                                    shouldValidate: true
                                  });
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
          </div>
        </div>
      </div>
    </div>
  );
};

export default Strategy;
