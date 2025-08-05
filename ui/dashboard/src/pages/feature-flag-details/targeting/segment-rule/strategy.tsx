import { useEffect, useRef, useState } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { RuleStrategyVariation, StrategyType } from '@types';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import ExperimentSelect from 'pages/experiments/experiments-modal/experiment-create-update/define-audience/experiment-select';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import { Tooltip } from 'components/tooltip';
import { StrategySchema } from '../form-schema';
import { isEquallyVariations } from '../utils';
import PercentageBar from './percentage-bar';
import PercentageInput from './percentage-input';
import { VariationOption } from './variation';

const experimentOptions = [
  {
    label: '5%',
    value: 5
  },
  {
    label: '10%',
    value: 10
  },
  {
    label: '50%',
    value: 50
  },
  {
    label: '90%',
    value: 90
  },
  {
    label: 'Custom',
    value: 'custom'
  }
];

const startOptions = [
  {
    label: '= Split equally',
    value: 'equally'
  },
  {
    label: '% Split by percentage',
    value: 'percentage'
  }
];

interface Props {
  rootName: string;
  strategyName: 'rolloutStrategy';
  variationOptions: VariationOption[];
  percentageValueCount: number;
  label?: string;
  isRequired?: boolean;
  isDisabled?: boolean;
  handleSelectStrategy: (item: VariationOption) => void;
}

type SplitOptionType = 'equally' | 'percentage';
type StrategyVariation = StrategySchema['rolloutStrategy']['variations'];

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
  const { control, watch, setError, setValue, clearErrors, setFocus } =
    useFormContext();

  const type = watch(`${rootName}.type`);
  const currentOption = watch(`${rootName}.currentOption`);
  const rolloutStrategy = watch(`${rootName}.${strategyName}`);
  const experimentPercentage = rolloutStrategy?.audience?.percentage || 0;
  const variations = (rolloutStrategy?.variations as StrategyVariation) || [];

  const [isCustomExperiment, setIsCustomExperiment] = useState(false);
  const [splitOptionType, setSplitOptionType] = useState<SplitOptionType>();
  const inputRef = useRef<HTMLInputElement>(null);

  const handleCheckError = (values: StrategyVariation) => {
    const total = values
      ?.map(v => Number(v.weight))
      .reduce((acc: number, current: number) => {
        return acc + (current || 0);
      }, 0);

    if (total !== 100) {
      return setError(`${rootName}.${strategyName}.variations`, {
        message: t('message:validation.should-be-percent')
      });
    } else {
      return clearErrors(`${rootName}.${strategyName}.variations`);
    }
  };

  const isShowPercentage =
    (type === StrategyType.ROLLOUT && strategyName === 'rolloutStrategy') ||
    currentOption === StrategyType.MANUAL;

  const handleSelectExperiment = (value: number | string) => {
    if (value === 'custom') {
      let timerId: NodeJS.Timeout | null = null;
      if (timerId) clearTimeout(timerId);
      setIsCustomExperiment(true);
      timerId = setTimeout(() => inputRef.current?.focus(), 100);
      setValue(`${rootName}.${strategyName}.audience.percentage`, '');
    } else {
      setIsCustomExperiment(false);
      setValue(`${rootName}.${strategyName}.audience.percentage`, value, {
        shouldDirty: true
      });
    }
  };

  const onChangeSplitType = (value: SplitOptionType) => {
    setSplitOptionType(value);
    let timerId: NodeJS.Timeout | null = null;
    if (timerId) clearTimeout(timerId);
    timerId = setTimeout(
      () => setFocus(`${rootName}.${strategyName}.variations.0.weight`),
      100
    );
    if (value === 'equally') {
      const equallyVariations = variations.map(item => ({
        ...item,
        weight: 100 / variations.length
      }));
      setValue(`${rootName}.${strategyName}.variations`, equallyVariations, {
        shouldDirty: true,
        shouldValidate: true
      });
    } else {
      const percentVariations = variations.map((item, index) => ({
        ...item,
        weight: index === 0 ? 100 : 0
      }));
      setValue(`${rootName}.${strategyName}.variations`, percentVariations, {
        shouldDirty: true,
        shouldValidate: true
      });
    }
  };

  useEffect(() => {
    const isActiveCustomExperiment = !experimentOptions
      .map(i => i.value)
      .includes(experimentPercentage);

    const variationsTotal = variations.reduce(
      (acc, item) => acc + (item.weight || 0),
      0
    );

    if (isActiveCustomExperiment && experimentPercentage > 0) {
      setIsCustomExperiment(true);
    }
    if (variationsTotal > 0) {
      setSplitOptionType(
        isEquallyVariations(variations) ? 'equally' : 'percentage'
      );
    }
  }, []);

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
        <p className="typo-para-small text-gray-600 mt-3 uppercase min-w-fit">
          {t('feature-flags.serve')}
        </p>
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
      </div>

      {isShowPercentage && (
        <>
          <div className="typo-para-small text-gray-500 my-5">
            {t('form:experiments.define-audience.any-traffic')}
          </div>

          <div className="w-full p-4 bg-gray-additional rounded-lg">
            <div className="flex flex-col w-full gap-y-4 typo-para-small leading-[14px] text-gray-600">
              <p>{t('form:experiments.define-audience.audience-amount')}</p>
              <div className="w-full h-3 p-[1px] border border-gray-400 rounded-full bg-gray-100">
                <div
                  className={cn('h-full bg-primary-500 rounded-l-full', {
                    'rounded-r-full': experimentPercentage >= 100
                  })}
                  style={{
                    width: `${experimentPercentage > 100 ? 100 : experimentPercentage}%`
                  }}
                />
              </div>
              <div className="flex items-center w-full gap-x-4">
                <div className="flex items-center gap-x-2">
                  <div className="flex-center size-5 m-0.5 rounded bg-primary-500" />
                  <p>{`${experimentPercentage}% ${t('form:experiments.define-audience.in-this-experiment')}`}</p>
                </div>
                <div className="flex items-center gap-x-2">
                  <div className="flex-center size-5 m-0.5 border border-gray-400 rounded bg-gray-100" />
                  <p>{`${100 - experimentPercentage}% ${t('form:experiments.define-audience.not-in-experiment')}`}</p>
                </div>
              </div>
            </div>
            <Divider className="my-5 border-gray-300" />

            <Form.Field
              control={control}
              name={`${rootName}.${strategyName}.audience.percentage`}
              render={({ field }) => {
                return (
                  <Form.Item className="flex flex-col flex-1 py-0 w-full">
                    <Form.Control>
                      <div>
                        <div className="flex items-center w-full gap-x-2">
                          <p className="typo-para-medium leading-5 text-gray-600 whitespace-nowrap">
                            {t(
                              'form:experiments.define-audience.in-this-experiment'
                            )}
                            :
                          </p>
                          {experimentOptions.map((item, index) => (
                            <ExperimentSelect
                              key={index}
                              label={item.label}
                              value={item.value}
                              isActive={
                                item.value === 'custom'
                                  ? isCustomExperiment
                                  : field.value === item.value
                              }
                              onSelect={handleSelectExperiment}
                            />
                          ))}
                        </div>
                        {isCustomExperiment && (
                          <Input
                            {...field}
                            ref={inputRef}
                            type="number"
                            className="mt-5"
                            value={field.value ?? ''}
                            onWheel={e => e.currentTarget.blur()}
                            onChange={value => field.onChange(value)}
                          />
                        )}
                      </div>
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                );
              }}
            />

            <div className="flex items-center w-full gap-x-2 mt-4 typo-para-medium leading-5 text-gray-600 whitespace-nowrap">
              <Trans
                i18nKey={
                  'form:experiments.define-audience.not-in-experiment-served'
                }
                values={{
                  percent: `${100 - experimentPercentage}%`
                }}
                components={{
                  highlight: (
                    <div className="flex-center size-fit p-3 rounded-lg typo-para-medium leading-5 text-gray-700 bg-gray-additional-2" />
                  )
                }}
              />
              <div className="flex-1">
                <Form.Field
                  control={control}
                  name={`${rootName}.${strategyName}.audience.defaultVariation`}
                  render={({ field }) => {
                    const options = variationOptions.slice(0, -1);
                    const option = options.find(
                      item => item.value === field.value
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
                              {options.map((item, index) => (
                                <DropdownMenuItem
                                  {...field}
                                  key={index}
                                  label={item.label}
                                  value={item.value}
                                  icon={item?.icon}
                                  onSelectOption={() => {
                                    field.onChange(item.value);
                                  }}
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
              </div>
            </div>
            <div className="flex items-center w-full gap-x-2 mt-4 typo-para-medium leading-5 text-gray-600 whitespace-nowrap">
              <Trans
                i18nKey={
                  'form:experiments.define-audience.in-experiment-target'
                }
                values={{
                  percent: `${experimentPercentage}%`
                }}
                components={{
                  highlight: (
                    <div className="flex-center size-fit p-3 rounded-lg typo-para-medium leading-5 text-gray-700 bg-gray-additional-2" />
                  )
                }}
              />
            </div>
          </div>
          <Divider className="my-5" />
          <div className="flex items-center gap-x-2">
            <p className="typo-para-medium text-gray-700">
              {t('form:experiments.define-audience.split-experiment')}
            </p>
            <Tooltip
              align="start"
              alignOffset={-68}
              content={``}
              trigger={
                <div className="flex-center size-fit">
                  <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
                </div>
              }
              className="max-w-[400px]"
            />
          </div>
          <div className="typo-para-small text-gray-500 mt-3">
            {`${t('form:experiments.define-audience.split-experiment-desc')}`}
          </div>

          <RadioGroup
            value={splitOptionType}
            onValueChange={onChangeSplitType}
            className="flex gap-x-6 mt-5"
          >
            {startOptions.map(({ label, value }) => (
              <div key={value} className="flex items-center gap-x-2">
                <RadioGroupItem value={value} id={value} />
                <label
                  htmlFor={value}
                  className="typo-para-medium leading-4 text-gray-600 cursor-pointer"
                >
                  {label}
                </label>
              </div>
            ))}
          </RadioGroup>
          <Divider className="my-5" />
          <div className="typo-para-small text-gray-500 mt-3">
            {`Audience slplit`}
          </div>

          <Form.Field
            control={control}
            name={`${rootName}.${strategyName}.variations`}
            render={({ field }) => {
              return (
                <Form.Item className="flex flex-col w-full gap-y-2">
                  <Form.Control>
                    <>
                      {percentageValueCount > 0 && (
                        <div className="flex items-center w-full p-0.5 border border-gray-400 rounded-full">
                          {field.value.map(
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
                      <div className="flex items-center w-full my-3">
                        <div className="typo-para-small text-gray-500 uppercase flex-1">{`Name`}</div>
                        <div className="typo-para-small text-gray-500 uppercase flex-1">{`Percentage`}</div>
                      </div>
                      <div className="flex flex-col gap-y-3">
                        {field.value?.map(
                          (rollout: RuleStrategyVariation, index: number) => (
                            <PercentageInput
                              key={index}
                              isDisabled={isDisabled || !splitOptionType}
                              variationOptions={variationOptions}
                              name={`${rootName}.${strategyName}.variations.${index}.weight`}
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
                  {splitOptionType && <Form.Message />}
                </Form.Item>
              );
            }}
          />
        </>
      )}
    </div>
  );
};

export default Strategy;
