import { useCallback, useEffect, useMemo, useState } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import dayjs from 'dayjs';
import { useTranslation } from 'i18n';
import { AutoOpsRule, Rollout } from '@types';
import { cn } from 'utils/style';
import { IconInfoFilled, IconPlus, IconTrash, IconWatch } from '@icons';
import { DateTimeClauseListType } from 'pages/feature-flag-details/operations/form-schema';
import { ActionTypeMap } from 'pages/feature-flag-details/operations/types';
import { createDatetimeClausesList } from 'pages/feature-flag-details/operations/utils';
import Button from 'components/button';
import { ReactDatePicker } from 'components/date-time-picker';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';

const WarningMessage = ({ description }: { description: string }) => {
  return (
    <div className="flex items-center w-full p-3 gap-x-2 rounded border-l-4 border-accent-yellow-600 bg-accent-yellow-50">
      <Icon icon={IconInfoFilled} size={'xxs'} color="accent-yellow-600" />
      <p className="typo-para-medium text-accent-yellow-600">{description}</p>
    </div>
  );
};

const ScheduleList = ({
  isCreate,
  isFinishedTab,
  rollouts,
  selectedData
}: {
  isCreate: boolean;
  isFinishedTab: boolean;
  rollouts: Rollout[];
  selectedData?: AutoOpsRule;
}) => {
  const { t } = useTranslation(['form', 'common']);
  const [conflictWithRolloutIndexes, setConflictWithRolloutIndexes] = useState<
    number[]
  >([]);

  const { control, watch, trigger } = useFormContext<DateTimeClauseListType>();

  const {
    fields: scheduleData,
    append,
    remove
  } = useFieldArray({
    name: 'datetimeClausesList',
    control,
    keyName: 'scheduleOperationId'
  });
  const watchScheduleList = [...watch('datetimeClausesList')];

  const stateOptions = useMemo(
    () => [
      {
        label: t('form:experiments.on'),
        value: ActionTypeMap.ENABLE
      },
      {
        label: t('form:experiments.off'),
        value: ActionTypeMap.DISABLE
      }
    ],
    []
  );

  const handleAddDate = useCallback(() => {
    const lastTime = watchScheduleList.at(-1)?.time?.getTime();
    const dateTimeClause = createDatetimeClausesList(lastTime);
    append({
      ...dateTimeClause
    });
  }, [watchScheduleList]);

  const isDisabledField = useCallback(
    (wasPassed?: boolean) => {
      return (
        (isFinishedTab && !!selectedData) || (!isCreate ? !!wasPassed : false)
      );
    },
    [isFinishedTab, isCreate, selectedData]
  );

  useEffect(() => {
    if (rollouts.length && watchScheduleList.length) {
      const waitingRolloutItems = rollouts.filter(
        item => item.status === 'WAITING'
      );
      if (waitingRolloutItems.length) {
        const flatMapRolloutItems = waitingRolloutItems.flatMap(
          item => item.clause?.schedules
        );

        const conflictIndexes: number[] = [];
        watchScheduleList.forEach((item, index) => {
          const timeString = Math.trunc(
            dayjs(item.time)?.set('second', 0)?.valueOf() / 1000
          )?.toString();

          if (
            flatMapRolloutItems.find(item => {
              const executeAtTime = dayjs(+item?.executeAt * 1000)?.set(
                'second',
                0
              );
              const executeAtString = (
                executeAtTime.valueOf() / 1000
              )?.toString();

              return executeAtString === timeString;
            })
          ) {
            conflictIndexes.push(index);
          }
        });
        setConflictWithRolloutIndexes(conflictIndexes);
      }
    }
  }, [rollouts, watchScheduleList]);

  return (
    <>
      <p className="typo-head-bold-small text-gray-800">
        {t('feature-flags.schedule')}
      </p>
      <Form.Field
        control={control}
        name="datetimeClausesList"
        render={() => (
          <Form.Item className="flex flex-col gap-y-4 py-0">
            <Form.Control>
              <div className="flex flex-col gap-y-4">
                {scheduleData.map((item, index) => (
                  <div
                    className="flex w-full gap-x-4"
                    key={item.scheduleOperationId}
                  >
                    <Form.Field
                      name={`datetimeClausesList.${index}.actionType`}
                      control={control}
                      render={({ field }) => (
                        <Form.Item className="py-0">
                          <Form.Label required>{t('common:state')}</Form.Label>
                          <Form.Control>
                            <DropdownMenu>
                              <DropdownMenuTrigger
                                label={
                                  stateOptions.find(
                                    item => item.value === field.value
                                  )?.label
                                }
                                className="w-[124px] uppercase"
                                disabled={isDisabledField(item.wasPassed)}
                              />
                              <DropdownMenuContent
                                align="start"
                                className="min-w-[124px]"
                                {...field}
                              >
                                {stateOptions.map(({ label, value }, index) => (
                                  <DropdownMenuItem
                                    key={index}
                                    label={label}
                                    value={value}
                                    onSelectOption={value =>
                                      field.onChange(value)
                                    }
                                    className="uppercase"
                                  />
                                ))}
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </Form.Control>
                          <Form.Message />
                        </Form.Item>
                      )}
                    />
                    <Form.Field
                      name={`datetimeClausesList.${index}.time`}
                      control={control}
                      render={({ field }) => (
                        <>
                          <Form.Item className="py-0">
                            <Form.Control>
                              <div className="flex gap-x-4">
                                <div>
                                  <Form.Label required>
                                    {t('feature-flags.start-date')}
                                  </Form.Label>
                                  <ReactDatePicker
                                    dateFormat={'yyyy/MM/dd'}
                                    selected={field.value ?? null}
                                    showTimeSelect={false}
                                    className={cn('w-[186px]', {
                                      '!border !border-accent-yellow-500':
                                        conflictWithRolloutIndexes.includes(
                                          index
                                        )
                                    })}
                                    disabled={isDisabledField(item.wasPassed)}
                                    onChange={date => {
                                      if (date) {
                                        field.onChange(date, {
                                          shouldValidate: true
                                        });
                                        trigger('datetimeClausesList');
                                      }
                                    }}
                                  />
                                </div>
                                <div>
                                  <Form.Label required>
                                    {t('feature-flags.time')}
                                  </Form.Label>
                                  <ReactDatePicker
                                    dateFormat={'HH:mm'}
                                    timeFormat="HH:mm"
                                    selected={field.value ?? null}
                                    showTimeSelectOnly={true}
                                    className={cn('w-[124px]', {
                                      '!border !border-accent-yellow-500':
                                        conflictWithRolloutIndexes.includes(
                                          index
                                        )
                                    })}
                                    disabled={isDisabledField(item.wasPassed)}
                                    onChange={date => {
                                      if (date) {
                                        field.onChange(date, {
                                          shouldValidate: true
                                        });
                                        trigger('datetimeClausesList');
                                      }
                                    }}
                                    icon={
                                      <Icon
                                        icon={IconWatch}
                                        className="flex-center"
                                      />
                                    }
                                  />
                                </div>
                              </div>
                            </Form.Control>
                            <Form.Message />
                          </Form.Item>
                        </>
                      )}
                    />
                    <Button
                      variant={'grey'}
                      size={'icon-sm'}
                      className="self-end mb-2"
                      disabled={
                        scheduleData.length <= 1 ||
                        isDisabledField(item.wasPassed)
                      }
                      onClick={() => remove(index)}
                    >
                      <Icon icon={IconTrash} size={'sm'} />
                    </Button>
                  </div>
                ))}
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
        disabled={isFinishedTab && !!selectedData}
        onClick={handleAddDate}
      >
        <Icon icon={IconPlus} size="md" className="flex-center" />
        {t('add-schedule')}
      </Button>

      {conflictWithRolloutIndexes.length > 0 && (
        <WarningMessage description={t('feature-flags.conflict-rollout')} />
      )}
    </>
  );
};

export default ScheduleList;
