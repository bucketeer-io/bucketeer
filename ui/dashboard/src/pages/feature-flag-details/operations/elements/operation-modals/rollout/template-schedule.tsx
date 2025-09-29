import { useEffect, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import dayjs from 'dayjs';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import { IconInfo } from '@icons';
import { RolloutSchemaType } from 'pages/feature-flag-details/operations/form-schema';
import { IntervalMap } from 'pages/feature-flag-details/operations/types';
import { ReactDatePicker } from 'components/date-time-picker';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownOption
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import { Tooltip } from 'components/tooltip';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';

const TemplateSchedule = ({
  disabled,
  variationOptions,
  isDisableCreateRollout
}: {
  disabled: boolean;
  variationOptions: DropdownOption[];
  isDisableCreateRollout: boolean;
}) => {
  const { t } = useTranslation(['form', 'table', 'common']);
  const { control, watch, setValue } = useFormContext<RolloutSchemaType>();

  const scheduleList = watch('progressiveRollout.template.schedulesList');

  const isDisabled = useMemo(
    () => isDisableCreateRollout || disabled,
    [isDisableCreateRollout, disabled]
  );

  const intervalOptions = useMemo(
    () => [
      {
        label: t('hourly'),
        value: IntervalMap.HOURLY
      },
      {
        label: t('daily'),
        value: IntervalMap.DAILY
      },
      {
        label: t('weekly'),
        value: IntervalMap.WEEKLY
      }
    ],
    []
  );

  const handleChangeScheduleList = ({
    increments,
    interval,
    startDate
  }: {
    increments?: number;
    interval?: IntervalMap;
    startDate?: Date;
  }) => {
    const _increments =
      increments || watch('progressiveRollout.template.increments');
    const _interval = interval || watch('progressiveRollout.template.interval');
    const _startDate =
      startDate || watch('progressiveRollout.template.startDate');
    if (!_increments)
      return setValue('progressiveRollout.template.schedulesList', [], {
        shouldValidate: true
      });
    const newScheduleList = Array(Math.ceil(100 / _increments))
      .fill('')
      .map((_, index) => {
        const time = dayjs(_startDate)
          .add(
            index,
            _interval === IntervalMap.HOURLY
              ? 'hour'
              : interval === IntervalMap.DAILY
                ? 'day'
                : 'week'
          )
          .toDate();

        const weight = (index + 1) * _increments;

        return {
          scheduleId: uuid(),
          executeAt: time,
          weight: weight > 100 ? 100 : Math.round(weight * 100) / 100,
          triggeredAt: ''
        };
      });
    setValue('progressiveRollout.template.schedulesList', newScheduleList, {
      shouldValidate: true
    });
  };

  useEffect(() => {
    handleChangeScheduleList({
      increments: 10,
      interval: IntervalMap.HOURLY
    });
  }, []);

  return (
    <div className="flex flex-col w-full gap-y-4">
      <Form.Field
        control={control}
        name={`progressiveRollout.template.variationId`}
        render={({ field }) => (
          <Form.Item className="py-0">
            <Form.Label required className="relative w-fit">
              {t('flag-variation')}
              <Tooltip
                align="start"
                alignOffset={-73}
                content={t('rollout-tooltips.template.variation')}
                trigger={
                  <div className="flex-center size-fit absolute top-0 -right-6">
                    <Icon icon={IconInfo} size="xs" color="gray-500" />
                  </div>
                }
                className="max-w-[300px]"
              />
            </Form.Label>
            <Form.Control>
              <DropdownMenuWithSearch
                align="end"
                label={
                  variationOptions.find(item => item.value === field.value)
                    ?.label || ''
                }
                itemSelected={field.value}
                contentClassName="[&>div.wrapper-menu-items>div]:px-4"
                options={variationOptions}
                disabled={isDisabled}
                onSelectOption={field.onChange}
              />
            </Form.Control>
            <Form.Message />
          </Form.Item>
        )}
      />

      <Form.Field
        control={control}
        name={`progressiveRollout.template.startDate`}
        render={({ field }) => (
          <Form.Item className="py-0">
            <Form.Label required className="relative w-fit">
              {t('feature-flags.start-date')}
              <Tooltip
                align="start"
                alignOffset={-73}
                content={t('rollout-tooltips.template.start-date')}
                trigger={
                  <div className="flex-center size-fit absolute top-0 -right-6">
                    <Icon icon={IconInfo} size="xs" color="gray-500" />
                  </div>
                }
                className="max-w-[300px]"
              />
            </Form.Label>
            <Form.Control>
              <ReactDatePicker
                selected={field.value ?? null}
                disabled={isDisabled}
                onChange={date => {
                  if (date) {
                    field.onChange(date);
                    handleChangeScheduleList({
                      startDate: date
                    });
                  }
                }}
              />
            </Form.Control>
            <Form.Message />
          </Form.Item>
        )}
      />
      <div className="flex w-full gap-x-4">
        <Form.Field
          control={control}
          name={`progressiveRollout.template.increments`}
          render={({ field }) => (
            <Form.Item className="py-0 flex-1 h-full">
              <Form.Label required className="relative w-fit">
                {t('increment')}
                <Tooltip
                  align="start"
                  alignOffset={-73}
                  content={t('rollout-tooltips.template.increment')}
                  trigger={
                    <div className="flex-center size-fit absolute top-0 -right-6">
                      <Icon icon={IconInfo} size="xs" color="gray-500" />
                    </div>
                  }
                  className="max-w-[300px]"
                />
              </Form.Label>
              <Form.Control>
                <InputGroup
                  className="w-full"
                  addonSlot="right"
                  addonSize="md"
                  addon={'%'}
                >
                  <Input
                    {...field}
                    value={field.value || ''}
                    type="number"
                    className="pr-8"
                    disabled={isDisabled}
                    onWheel={e => {
                      e.currentTarget.blur();
                    }}
                    onKeyDown={e => {
                      if (e.key === 'Enter') e.preventDefault();
                    }}
                    onChange={value => {
                      field.onChange(+value);
                      handleChangeScheduleList({
                        increments: +value
                      });
                    }}
                  />
                </InputGroup>
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />
        <Form.Field
          control={control}
          name={`progressiveRollout.template.interval`}
          render={({ field }) => (
            <Form.Item className="py-0 flex-1 h-full">
              <Form.Label required className="relative w-fit">
                {t('frequency')}
                <Tooltip
                  align="start"
                  alignOffset={-150}
                  content={t('rollout-tooltips.template.frequency')}
                  trigger={
                    <div className="flex-center size-fit absolute top-0 -right-6">
                      <Icon icon={IconInfo} size="xs" color="gray-500" />
                    </div>
                  }
                  className="max-w-[300px]"
                />
              </Form.Label>
              <Form.Control>
                <DropdownMenu>
                  <DropdownMenuTrigger
                    label={
                      intervalOptions.find(item => item.value === field.value)
                        ?.label || ''
                    }
                    isExpand
                    disabled={isDisabled}
                  />
                  <DropdownMenuContent align="end" className="min-w-[243px]">
                    {intervalOptions.map((item, index) => (
                      <DropdownMenuItem
                        key={index}
                        label={item.label}
                        value={item.value}
                        isSelectedItem={item.value === field.value}
                        onSelectOption={value => {
                          field.onChange(value);
                          handleChangeScheduleList({
                            interval: value as IntervalMap
                          });
                        }}
                      />
                    ))}
                  </DropdownMenuContent>
                </DropdownMenu>
              </Form.Control>
            </Form.Item>
          )}
        />
      </div>
      <div className="flex flex-col w-full">
        <div className="flex items-center w-full gap-x-4">
          <Form.Label className="flex flex-1">{t('weight')}</Form.Label>
          <Form.Label className="flex flex-1">{t('execute-at')}</Form.Label>
        </div>
        <div className="flex flex-col w-full gap-y-3">
          {scheduleList?.map((item, index) => (
            <div key={index} className="flex items-center w-full gap-x-4">
              <div className="flex flex-1">
                <InputGroup
                  className="w-full"
                  addonSlot="right"
                  addonSize="md"
                  addon={'%'}
                >
                  <Input value={item.weight} className="pr-8" disabled />
                </InputGroup>
              </div>
              <div className="flex flex-1">
                <ReactDatePicker
                  selected={item.executeAt ?? null}
                  disabled={true}
                />
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default TemplateSchedule;
