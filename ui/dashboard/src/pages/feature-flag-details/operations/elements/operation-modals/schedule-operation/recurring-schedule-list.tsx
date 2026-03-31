import { useCallback } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import dayjs from 'dayjs';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import { cn } from 'utils/style';
import { IconInfo, IconPlus, IconTrash, IconWatch } from '@icons';
import { ScheduleOperationFormType } from 'pages/feature-flag-details/operations/form-schema';
import {
  ActionTypeMap,
  DAY_LABELS_FULL,
  DAY_LABELS_SHORT_KEYS,
  DAYS_OF_WEEK,
  EndConditionType,
  FREQUENCY_OPTIONS
} from 'pages/feature-flag-details/operations/types';
import Button from 'components/button';
import { ReactDatePicker } from 'components/date-time-picker';
import Dropdown from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import { Tooltip } from 'components/tooltip';

const RecurringScheduleList = ({ isDisabled }: { isDisabled: boolean }) => {
  const { t } = useTranslation(['form', 'common']);
  const { control, watch, setValue, trigger } =
    useFormContext<ScheduleOperationFormType>();

  const frequency = watch('recurring.frequency');
  const daysOfWeek = watch('recurring.daysOfWeek') || [];
  const endCondition = watch('recurring.endCondition');

  const {
    fields: recurringClauses,
    append,
    remove
  } = useFieldArray({
    name: 'recurring.recurringClauses',
    control,
    keyName: 'clauseKey'
  });

  const watchClausesList = watch('recurring.recurringClauses');

  const frequencyKeyMap: Record<string, string> = {
    DAILY: 'form:daily',
    WEEKLY: 'form:weekly',
    MONTHLY: 'form:monthly'
  };

  const stateOptions = [
    {
      label: t('form:experiments.on'),
      value: ActionTypeMap.ENABLE
    },
    {
      label: t('form:experiments.off'),
      value: ActionTypeMap.DISABLE
    }
  ];

  const frequencyOptions = FREQUENCY_OPTIONS.map(freq => ({
    label: t(frequencyKeyMap[freq] ?? 'form:unknown'),
    value: freq
  }));

  const handleAddClause = useCallback(() => {
    const lastClause = watchClausesList.at(-1);
    const baseTime = lastClause?.time ?? new Date();
    const nextTime = dayjs(baseTime).add(1, 'hour');
    const normalizedTime = new Date(
      1970,
      0,
      1,
      nextTime.hour(),
      nextTime.minute(),
      0,
      0
    );
    append({
      id: uuid(),
      actionType: ActionTypeMap.ENABLE,
      time: normalizedTime
    });
  }, [append, watchClausesList]);

  const handleToggleDay = useCallback(
    (day: number) => {
      const current = daysOfWeek || [];
      const updated = current.includes(day)
        ? current.filter(d => d !== day)
        : [...current, day].sort((a, b) => a - b);
      setValue('recurring.daysOfWeek', updated, { shouldValidate: true });
    },
    [daysOfWeek, setValue]
  );

  return (
    <div className="flex flex-col gap-y-5">
      <p className="typo-head-bold-small text-gray-800">
        {t('feature-flags.schedule')}
      </p>

      <Form.Field
        control={control}
        name="recurring.startDate"
        render={({ field }) => (
          <Form.Item className="py-0">
            <Form.Label required className="relative w-fit">
              {t('feature-flags.start-date')}
              <Tooltip
                align="start"
                alignOffset={-73}
                content={t('feature-flags.recurring-tooltips.start-date')}
                trigger={
                  <button
                    type="button"
                    className="flex-center size-fit absolute top-0 -right-6"
                  >
                    <Icon icon={IconInfo} size="xs" color="gray-500" />
                  </button>
                }
                className="max-w-[300px]"
              />
            </Form.Label>
            <Form.Control>
              <ReactDatePicker
                dateFormat={'yyyy/MM/dd'}
                selected={field.value ?? null}
                showTimeSelect={false}
                className="w-[186px]"
                wrapperClassName="w-[186px]"
                disabled={isDisabled}
                onChange={date => {
                  if (date) field.onChange(date);
                }}
              />
            </Form.Control>
            <Form.Message />
          </Form.Item>
        )}
      />

      <Form.Field
        control={control}
        name="recurring.frequency"
        render={({ field }) => (
          <Form.Item className="py-0">
            <Form.Label required>{t('feature-flags.repeat')}</Form.Label>
            <Form.Control>
              <Dropdown
                value={field.value}
                options={frequencyOptions}
                onChange={field.onChange}
                className="w-[186px]"
                disabled={isDisabled}
              />
            </Form.Control>
            <Form.Message />
          </Form.Item>
        )}
      />

      {frequency === 'WEEKLY' && (
        <Form.Field
          control={control}
          name="recurring.daysOfWeek"
          render={() => (
            <Form.Item className="py-0">
              <Form.Label required>{t('feature-flags.days')}</Form.Label>
              <Form.Control>
                <div className="flex gap-x-2">
                  {DAYS_OF_WEEK.map(day => {
                    const isActive = daysOfWeek?.includes(day);
                    return (
                      <button
                        key={day}
                        type="button"
                        disabled={isDisabled}
                        onClick={() => handleToggleDay(day)}
                        aria-label={DAY_LABELS_FULL[day]}
                        aria-pressed={isActive}
                        className={cn(
                          'flex-center size-9 rounded-md border typo-para-medium transition-colors',
                          isActive
                            ? 'bg-primary-500 text-white border-primary-500'
                            : 'bg-white text-gray-700 border-gray-400 hover:border-gray-500'
                        )}
                      >
                        {t(`form:${DAY_LABELS_SHORT_KEYS[day]}`)}
                      </button>
                    );
                  })}
                </div>
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />
      )}

      {frequency === 'MONTHLY' && (
        <Form.Field
          control={control}
          name="recurring.dayOfMonth"
          render={({ field }) => (
            <Form.Item className="py-0">
              <Form.Label required>
                {t('feature-flags.day-of-month')}
              </Form.Label>
              <p className="typo-para-small text-gray-500 mb-2">
                {t('feature-flags.day-of-month-hint')}
              </p>
              <Form.Control>
                <Dropdown
                  value={String(field.value || 1)}
                  options={Array.from({ length: 31 }, (_, i) => ({
                    label: String(i + 1),
                    value: String(i + 1)
                  }))}
                  onChange={val => field.onChange(Number(val))}
                  className="w-[186px]"
                  disabled={isDisabled}
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />
      )}

      <div className="flex flex-col gap-y-3">
        <Form.Label required className="relative w-fit">
          {t('feature-flags.end')}
          <Tooltip
            align="start"
            alignOffset={-73}
            content={t('feature-flags.recurring-tooltips.end')}
            trigger={
              <button
                type="button"
                className="flex-center size-fit absolute top-0 -right-6"
              >
                <Icon icon={IconInfo} size="xs" color="gray-500" />
              </button>
            }
            className="max-w-[300px]"
          />
        </Form.Label>
        <Form.Field
          control={control}
          name="recurring.endCondition"
          render={({ field }) => (
            <Form.Item className="py-0">
              <Form.Control>
                <RadioGroup
                  value={field.value}
                  onValueChange={value => {
                    field.onChange(value);
                    if (value !== EndConditionType.ON_DATE) {
                      setValue('recurring.endDate', undefined as never, {
                        shouldValidate: true
                      });
                    }
                    if (value !== EndConditionType.AFTER) {
                      setValue('recurring.maxOccurrences', undefined as never, {
                        shouldValidate: true
                      });
                    }
                  }}
                  disabled={isDisabled}
                  className="flex flex-col gap-y-4"
                >
                  <div className="flex items-center pb-5 gap-x-2">
                    <RadioGroupItem value={EndConditionType.NEVER} />
                    <span className="typo-para-medium text-gray-700">
                      {t('feature-flags.end-never')}
                    </span>
                  </div>

                  <div className="flex items-center pb-5 gap-x-2">
                    <RadioGroupItem value={EndConditionType.ON_DATE} />
                    <span className="typo-para-medium text-gray-700">
                      {t('feature-flags.end-on-date')}
                    </span>
                    <Form.Field
                      control={control}
                      name="recurring.endDate"
                      render={({ field: dateField }) => (
                        <Form.Item className="py-0 relative">
                          <Form.Control>
                            <ReactDatePicker
                              dateFormat={'yyyy/MM/dd'}
                              selected={dateField.value ?? null}
                              showTimeSelect={false}
                              className="w-[160px]"
                              disabled={
                                isDisabled ||
                                endCondition !== EndConditionType.ON_DATE
                              }
                              onChange={date => {
                                if (date) dateField.onChange(date);
                              }}
                            />
                          </Form.Control>
                          <Form.Message className="absolute top-full left-0 whitespace-nowrap" />
                        </Form.Item>
                      )}
                    />
                  </div>

                  <div className="flex items-center pb-5 gap-x-2">
                    <RadioGroupItem value={EndConditionType.AFTER} />
                    <span className="typo-para-medium text-gray-700">
                      {t('feature-flags.end-after')}
                    </span>
                    <Form.Field
                      control={control}
                      name="recurring.maxOccurrences"
                      render={({ field: occField }) => (
                        <Form.Item className="py-0 relative">
                          <Form.Control>
                            <div className="flex items-center gap-x-2">
                              <input
                                type="number"
                                min={1}
                                value={occField.value ?? ''}
                                onChange={e =>
                                  occField.onChange(
                                    e.target.value
                                      ? Number(e.target.value)
                                      : undefined
                                  )
                                }
                                disabled={
                                  isDisabled ||
                                  endCondition !== EndConditionType.AFTER
                                }
                                className="w-[160px] py-[11px] px-4 rounded-lg border border-gray-400 typo-para-medium text-gray-700 focus:outline-none focus:border-primary-500 disabled:bg-gray-100 disabled:text-gray-400"
                              />
                              <span className="typo-para-medium text-gray-700">
                                {t('feature-flags.occurrences')}
                              </span>
                            </div>
                          </Form.Control>
                          <Form.Message className="absolute top-full left-0 whitespace-nowrap" />
                        </Form.Item>
                      )}
                    />
                  </div>
                </RadioGroup>
              </Form.Control>
            </Form.Item>
          )}
        />
      </div>

      <div className="border-t border-gray-300 my-1" />

      <Form.Field
        control={control}
        name="recurring.recurringClauses"
        render={() => (
          <Form.Item className="flex flex-col gap-y-4 py-0">
            <Form.Control>
              <div className="flex flex-col gap-y-4">
                {recurringClauses.map((item, index) => {
                  const clauseExecuted =
                    watchClausesList[index]?.wasExecuted ?? false;
                  const clauseDisabled = isDisabled || clauseExecuted;
                  return (
                    <div className="flex w-full gap-x-4" key={item.clauseKey}>
                      <Form.Field
                        name={`recurring.recurringClauses.${index}.actionType`}
                        control={control}
                        render={({ field }) => (
                          <Form.Item className="py-0">
                            <Form.Label required>
                              {t('common:state')}
                            </Form.Label>
                            <Form.Control>
                              <Dropdown
                                value={field.value}
                                options={stateOptions}
                                onChange={field.onChange}
                                className="w-[124px] uppercase"
                                disabled={clauseDisabled}
                                contentClassName="min-w-[124px]"
                              />
                            </Form.Control>
                            <Form.Message />
                          </Form.Item>
                        )}
                      />
                      <Form.Field
                        name={`recurring.recurringClauses.${index}.time`}
                        control={control}
                        render={({ field }) => (
                          <Form.Item className="py-0">
                            <Form.Label required>
                              {t('feature-flags.time')}
                            </Form.Label>
                            <Form.Control>
                              <ReactDatePicker
                                dateFormat={'HH:mm'}
                                timeFormat="HH:mm"
                                selected={field.value ?? null}
                                showTimeSelectOnly={true}
                                className="w-[124px]"
                                disabled={clauseDisabled}
                                onChange={date => {
                                  if (date) {
                                    field.onChange(date);
                                    trigger('recurring.recurringClauses');
                                  }
                                }}
                                icon={
                                  <Icon
                                    icon={IconWatch}
                                    className="flex-center"
                                  />
                                }
                              />
                            </Form.Control>
                            <Form.Message />
                          </Form.Item>
                        )}
                      />
                      <Button
                        variant={'grey'}
                        size={'icon-sm'}
                        className="self-end mb-2"
                        disabled={
                          recurringClauses.length <= 1 || clauseDisabled
                        }
                        onClick={() => remove(index)}
                      >
                        <Icon icon={IconTrash} size={'sm'} />
                      </Button>
                    </div>
                  );
                })}
              </div>
            </Form.Control>
            <Form.Message />
          </Form.Item>
        )}
      />

      <Button
        type="button"
        variant={'text'}
        className="flex items-center h-6 self-start p-0"
        disabled={isDisabled}
        onClick={handleAddClause}
      >
        <Icon icon={IconPlus} size="md" className="flex-center" />
        {t('add-schedule')}
      </Button>
    </div>
  );
};

export default RecurringScheduleList;
