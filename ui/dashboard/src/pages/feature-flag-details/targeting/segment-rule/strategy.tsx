import { useEffect, useMemo, useRef, useState } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import useOptions from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { RuleStrategyVariation, StrategyType } from '@types';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import Divider from 'components/divider';
import Dropdown, { DropdownValue } from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import { Tooltip } from 'components/tooltip';
import { StrategySchema } from '../form-schema';
import { isEquallyVariations } from '../utils';
import AudienceSelect from './audience-select';
import PercentageBar from './percentage-bar';
import PercentageInput from './percentage-input';
import { VariationOption } from './variation';

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
  const { splitExperimentOptions, audienceTrafficOptions } = useOptions();

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

  const options = useMemo(
    () =>
      variationOptions.map(opt => {
        return {
          label: (
            <div className="flex items-center gap-x-2 typo-para-medium text-gray-700">
              {opt.icon && <Icon icon={opt.icon} size={'sm'} />}
              {opt.label}
            </div>
          ),
          value: opt.value,
          type: opt.type,
          variationValue: opt?.variationValue
        };
      }),
    [variationOptions]
  );
  const handleChangeStrategy = (
    val: DropdownValue | DropdownValue[],
    onChange: (item: DropdownValue | DropdownValue[]) => void
  ) => {
    onChange(val);
    const selected = variationOptions.find(opt => opt.value === val);
    if (selected) handleSelectStrategy(selected);
  };

  const handleSelectExperiment = (value: number | string) => {
    if (value === 'custom') {
      let timerId: NodeJS.Timeout | null = null;
      if (timerId) clearTimeout(timerId);
      setIsCustomExperiment(true);
      timerId = setTimeout(() => inputRef.current?.focus(), 100);
      setValue(`${rootName}.${strategyName}.audience.percentage`, '');
    } else {
      setIsCustomExperiment(false);
      setValue(
        `${rootName}.${strategyName}.audience.defaultVariation`,
        value === 100 ? '' : rolloutStrategy.audience?.defaultVariation,
        { shouldDirty: true }
      );
      setValue(`${rootName}.${strategyName}.audience.percentage`, value, {
        shouldDirty: true,
        shouldValidate: true
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
    const isActiveCustomExperiment = !audienceTrafficOptions
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
    <div className="px-2">
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
            return (
              <Form.Item className="flex flex-col flex-1 py-0 w-full">
                <Form.Control>
                  <Dropdown
                    value={currentOption}
                    options={options}
                    onChange={val => handleChangeStrategy(val, field.onChange)}
                    disabled={isDisabled}
                    wrapTriggerStyle="flex flex-col grap-y-2 w-full"
                    className="w-full"
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            );
          }}
        />
      </div>

      {isShowPercentage && (
        <>
          <div className="flex items-center gap-x-1">
            <div className="typo-para-medium text-gray-700 my-5">
              {t('form:experiments.define-audience.audience-traffic')}
            </div>
            <Tooltip
              align="start"
              alignOffset={-68}
              content={t(
                'form:experiments.define-audience.audience-traffic-desc'
              )}
              trigger={
                <div className="flex-center size-fit mt-1">
                  <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
                </div>
              }
              className="max-w-[400px]"
            />
          </div>
          <Form.Field
            control={control}
            name={`${rootName}.${strategyName}.audience.percentage`}
            render={({ field }) => {
              return (
                <Form.Item className="flex flex-col flex-1 py-0 w-full">
                  <Form.Control>
                    <div>
                      <div className="flex items-center w-full gap-x-2">
                        {audienceTrafficOptions.map((item, index) => (
                          <AudienceSelect
                            disabled={isDisabled}
                            key={index}
                            label={item.label}
                            value={item.value}
                            isActive={
                              item.value === 'custom'
                                ? isCustomExperiment
                                : !isCustomExperiment &&
                                  field.value === item.value
                            }
                            onSelect={handleSelectExperiment}
                          />
                        ))}
                        {isCustomExperiment && (
                          <div className="flex-1 relative">
                            <InputGroup
                              addon={'%'}
                              addonSlot="right"
                              className="w-32 overflow-hidden"
                              addonClassName="top-[1px] bottom-[1px] right-[1px] translate-x-0 translate-y-0 !flex-center rounded-r-lg bg-gray-200 w-[29px] typo-para-medium text-gray-700"
                            >
                              <Input
                                {...field}
                                ref={inputRef}
                                value={field.value ?? ''}
                                onWheel={e => e.currentTarget.blur()}
                                onChange={value => field.onChange(value)}
                                type="number"
                                className="text-right pl-[5px]"
                              />
                            </InputGroup>
                          </div>
                        )}
                      </div>
                    </div>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              );
            }}
          />
          <div className="flex flex-col w-full gap-y-4 typo-para-small leading-[14px] text-gray-600 mt-6">
            <div className="w-full h-3 p-[1px] border border-gray-400 rounded-full ">
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
                <p>{`${experimentPercentage}% ${t('form:experiments.define-audience.included')}`}</p>
              </div>
              <div className="flex items-center gap-x-2">
                <div className="flex-center size-5 m-0.5 border border-gray-400 rounded bg-gray-100" />
                <p>{`${100 - experimentPercentage}% ${t('form:experiments.define-audience.excluded')}`}</p>
              </div>
            </div>
          </div>

          {experimentPercentage > 0 && Number(experimentPercentage) !== 100 && (
            <div className="flex items-center w-full gap-x-2 mt-4 typo-para-medium leading-5 text-gray-600 whitespace-nowrap">
              <Trans
                i18nKey={
                  'form:experiments.define-audience.not-included-allocation'
                }
                values={{
                  percent: `${100 - experimentPercentage}%`
                }}
                components={{
                  highlight: (
                    <div className="flex-center size-fit p-3 rounded-lg typo-para-medium leading-5 text-gray-700 bg-gray-additional-2" />
                  ),
                  select: (
                    <Form.Field
                      control={control}
                      name={`${rootName}.${strategyName}.audience.defaultVariation`}
                      render={({ field }) => (
                        <Form.Item className="flex flex-col flex-1 py-0 w-full">
                          <Form.Control>
                            <Dropdown
                              options={options.slice(0, -1)}
                              value={field.value}
                              onChange={field.onChange}
                              disabled={isDisabled}
                              className="w-full"
                              wrapTriggerStyle="flex flex-col grap-y-2 w-full"
                            />
                          </Form.Control>
                          <Form.Message />
                        </Form.Item>
                      )}
                    />
                  )
                }}
              />
            </div>
          )}
          <Divider className="my-5 border-gray-300" />
          <div className="flex items-center gap-x-1">
            <div className="typo-para-medium text-gray-700">
              {t('form:experiments.define-audience.variation-allocation')}
            </div>
            <Tooltip
              align="start"
              alignOffset={-68}
              content={t(
                'form:experiments.define-audience.variation-allocation-desc'
              )}
              trigger={
                <div className="flex-center size-fit mt-1">
                  <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
                </div>
              }
              className="max-w-[400px]"
            />
          </div>
          <RadioGroup
            disabled={isDisabled}
            value={splitOptionType}
            onValueChange={onChangeSplitType}
            className="flex gap-x-6 mt-5 px-1"
          >
            {splitExperimentOptions.map(({ label, value }) => (
              <div key={value} className="flex items-center gap-x-2">
                <RadioGroupItem
                  disabled={isDisabled}
                  value={value}
                  id={value}
                />
                <label
                  htmlFor={value}
                  className={cn(
                    'typo-para-medium leading-4 text-gray-600',
                    isDisabled
                      ? 'cursor-not-allowed pointer-events-none'
                      : 'cursor-pointer'
                  )}
                >
                  {label}
                </label>
              </div>
            ))}
          </RadioGroup>
          <div className="typo-para-small text-gray-500 mt-6">
            {t('form:experiments.define-audience.audience-split')}
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
                      <div className="flex items-center w-full mb-3 mt-5 px-2">
                        <div className="typo-para-small text-gray-500 uppercase flex-1">
                          {t('common:name')}
                        </div>
                        <div className="typo-para-small text-gray-500 uppercase flex-1">
                          {t(`common:percentage`)}
                        </div>
                      </div>
                      <div className="flex flex-col gap-y-3 px-2">
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
