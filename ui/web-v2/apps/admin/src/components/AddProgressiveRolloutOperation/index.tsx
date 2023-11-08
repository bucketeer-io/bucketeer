import { messages } from '@/lang/messages';
import { useIsEditable } from '@/modules/me';
import { classNames } from '@/utils/css';
import { MinusCircleIcon, PlusIcon } from '@heroicons/react/outline';
import dayjs from 'dayjs';
import React, { memo, FC, useEffect, useCallback, Fragment } from 'react';
import {
  Controller,
  useFieldArray,
  useFormContext,
  useWatch,
} from 'react-hook-form';
import { useIntl } from 'react-intl';

import { DatetimePicker } from '../DatetimePicker';
import {
  ProgressiveRolloutClauseType,
  getIntervalForDayjs,
} from '../FeatureAutoOpsRulesForm';
import { Option, Select } from '../Select';

interface AddProgressiveRolloutOperationProps {
  variationOptions: Option[];
  isSeeDetailsSelected: boolean;
  progressiveRolloutTypeList: any;
  setProgressiveRolloutTypeList: any;
}

export const AddProgressiveRolloutOperation: FC<AddProgressiveRolloutOperationProps> =
  memo(
    ({
      variationOptions,
      isSeeDetailsSelected,
      progressiveRolloutTypeList,
      setProgressiveRolloutTypeList,
    }) => {
      const { formatMessage: f } = useIntl();
      const editable = useIsEditable();

      const methods = useFormContext<any>();
      const {
        control,
        formState: { errors },
        register,
        setValue,
      } = methods;

      const templateSchedulesList = useWatch({
        control,
        name: 'progressiveRollout.template.schedulesList',
      });

      const incrementsWatch = useWatch({
        control,
        name: 'progressiveRollout.template.increments',
      });

      const templateIntervalWatch = useWatch({
        control,
        name: 'progressiveRollout.template.interval',
      });

      const progressiveRolloutTemplateDatetimeTimeWatch = useWatch({
        control,
        name: 'progressiveRollout.template.datetime.time',
      });

      const manualSchedulesListWatch = useWatch({
        control,
        name: 'progressiveRollout.manual.schedulesList',
      });

      const {
        fields: manualSchedulesList,
        remove: removeTrigger,
        append,
      } = useFieldArray({
        control,
        name: 'progressiveRollout.manual.schedulesList',
      });

      useEffect(() => {
        if (
          progressiveRolloutTypeList.find((p) => p.selected).value ===
            ProgressiveRolloutClauseType.PROGRESSIVE_ROLLOUT_TEMPLATE_SCHEDULE &&
          Number(incrementsWatch) > 0
        ) {
          const scheduleList = Array(Math.ceil(100 / incrementsWatch))
            .fill('')
            .map((_, index) => {
              // increment each schedule by {templateIntervalWatch}
              const datetime = dayjs(
                progressiveRolloutTemplateDatetimeTimeWatch
              )
                .add(index, getIntervalForDayjs(templateIntervalWatch))
                .toDate();

              const weight = (index + 1) * incrementsWatch;

              return {
                executeAt: {
                  time: datetime,
                },
                weight: weight > 100 ? 100 : Math.round(weight * 100) / 100,
              };
            });
          setValue('progressiveRollout.template.schedulesList', scheduleList);
        }
      }, [
        progressiveRolloutTemplateDatetimeTimeWatch,
        templateIntervalWatch,
        incrementsWatch,
      ]);

      const intervalOptions = [
        {
          label: 'Hourly',
          value: '1',
        },
        {
          label: 'Daily',
          value: '2',
        },
        {
          label: 'Weekly',
          value: '3',
        },
      ];

      const handleAddOperation = (e) => {
        e.preventDefault();

        const lastSchedule: any =
          manualSchedulesListWatch[manualSchedulesListWatch.length - 1];

        const time = dayjs(lastSchedule?.executeAt.time)
          .add(5, 'minute')
          .toDate();

        const weight = lastSchedule ? Number(lastSchedule.weight) : 0;

        if (weight < 100) {
          append({
            executeAt: {
              time,
            },
            weight: weight > 80 ? 100 : weight + 20,
          });
        }
      };

      const handleRemoveTrigger = useCallback(
        (idx) => {
          removeTrigger(idx);
        },
        [removeTrigger]
      );

      return (
        <div className="mt-4 space-y-4">
          <div className="flex">
            {progressiveRolloutTypeList.map(({ label, selected }, index) => (
              <div
                key={label}
                className={classNames(
                  'flex-1 border py-2 text-center cursor-pointer',
                  index === 0 ? 'rounded-l-lg' : 'rounded-r-lg',
                  selected ? 'bg-pink-50 border-pink-500' : ''
                )}
                onClick={() => {
                  setProgressiveRolloutTypeList(
                    progressiveRolloutTypeList.map((p) => ({
                      ...p,
                      selected: p.label === label,
                    }))
                  );
                }}
              >
                {label}
              </div>
            ))}
          </div>
          {progressiveRolloutTypeList.find((p) => p.selected).value ===
          ProgressiveRolloutClauseType.PROGRESSIVE_ROLLOUT_TEMPLATE_SCHEDULE ? (
            <Fragment>
              <div className="flex-1">
                <span className="input-label">
                  {f(messages.feature.variation)}
                </span>
                <Controller
                  name="progressiveRollout.template.variationId"
                  control={control}
                  render={({ field }) => (
                    <Select
                      onChange={(o: Option) => field.onChange(o.value)}
                      options={variationOptions}
                      disabled={!editable || isSeeDetailsSelected}
                      value={variationOptions.find(
                        (o) => o.value === field.value
                      )}
                    />
                  )}
                />
              </div>
              <div className="">
                <span className="input-label">
                  {f(messages.autoOps.startDate)}
                </span>
                <DatetimePicker
                  name="progressiveRollout.template.datetime.time"
                  disabled={isSeeDetailsSelected}
                />
                <p className="input-error">
                  {errors.progressiveRollout?.template?.datetime?.time
                    ?.message && (
                    <span role="alert">
                      {
                        errors.progressiveRollout?.template?.datetime?.time
                          ?.message
                      }
                    </span>
                  )}
                </p>
              </div>
              <div className="flex space-x-4">
                <div className="flex-1">
                  <span className="input-label">Increment</span>
                  <div className="flex">
                    <input
                      type="number"
                      {...register('progressiveRollout.template.increments')}
                      min="0"
                      max="100"
                      onKeyDown={(evt: any) => {
                        if (evt.key === '.') {
                          evt.preventDefault();
                        }
                      }}
                      className={classNames('w-full', 'input-text')}
                      placeholder={''}
                      required
                      disabled={!editable || isSeeDetailsSelected}
                    />
                    <span
                      className={classNames(
                        'px-1 py-1 inline-flex items-center bg-gray-100',
                        'rounded-r border border-l-0 border-gray-300 text-gray-600'
                      )}
                    >
                      {'%'}
                    </span>
                  </div>
                  <p className="input-error">
                    {errors.progressiveRollout?.template?.increments
                      ?.message && (
                      <span role="alert">
                        {
                          errors.progressiveRollout?.template?.increments
                            ?.message
                        }
                      </span>
                    )}
                  </p>
                </div>
                <div className="flex-1">
                  <span className="input-label">Frequency</span>
                  <Controller
                    name="progressiveRollout.template.interval"
                    control={control}
                    render={({ field }) => (
                      <Select
                        onChange={(o: Option) => field.onChange(o.value)}
                        options={intervalOptions}
                        disabled={!editable || isSeeDetailsSelected}
                        value={intervalOptions.find(
                          (o) => o.value === field.value
                        )}
                      />
                    )}
                  />
                </div>
              </div>
              <div>
                <div className="flex space-x-4">
                  <div className="flex-1 input-label">Weight</div>
                  <div className="flex-1 input-label">Execute at</div>
                </div>
                <div
                  className={classNames(
                    'space-y-2 mt-2',
                    templateSchedulesList.length > 5 &&
                      'max-h-[232px] overflow-y-scroll'
                  )}
                >
                  {templateSchedulesList?.map((_, index) => (
                    <div key={index} className="flex space-x-4">
                      <div className="flex w-1/2">
                        <input
                          {...register(
                            `progressiveRollout.template.schedulesList.${index}.weight`
                          )}
                          onKeyDown={(evt: any) => {
                            if (evt.key === '.') {
                              evt.preventDefault();
                            }
                          }}
                          type="number"
                          className={classNames('w-full', 'input-text')}
                          disabled={true}
                        />
                        <span
                          className={classNames(
                            'px-1 py-1 inline-flex items-center bg-gray-100',
                            'rounded-r border border-l-0 border-gray-300 text-gray-600'
                          )}
                        >
                          {'%'}
                        </span>
                      </div>
                      <div className="w-1/2">
                        <div>
                          <DatetimePicker
                            name={`progressiveRollout.template.schedulesList.${index}.executeAt.time`}
                            disabled={true}
                          />
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </Fragment>
          ) : (
            <Fragment>
              <div className="flex-1">
                <span className="input-label">
                  {f(messages.feature.variation)}
                </span>
                <Controller
                  name="progressiveRollout.manual.variationId"
                  control={control}
                  render={({ field }) => (
                    <Select
                      onChange={(o: Option) => field.onChange(o.value)}
                      options={variationOptions}
                      disabled={!editable || isSeeDetailsSelected}
                      value={variationOptions.find(
                        (o) => o.value === field.value
                      )}
                    />
                  )}
                />
              </div>
              <button
                onClick={handleAddOperation}
                className="text-primary space-x-2 flex items-center"
              >
                <PlusIcon width={16} />
                <span className="text-sm font-medium">Add Operation</span>
              </button>
              <div className="space-y-2">
                {manualSchedulesList.map((_, index) => (
                  <div key={index} className="flex space-x-4">
                    <div className="w-full">
                      <div className="flex">
                        <input
                          {...register(
                            `progressiveRollout.manual.schedulesList.${index}.weight`
                          )}
                          type="number"
                          min="0"
                          max="100"
                          className={classNames(
                            'w-full',
                            errors.progressiveRollout?.manual?.schedulesList[
                              index
                            ]?.weight?.message
                              ? 'input-text-error'
                              : 'input-text'
                          )}
                          placeholder={''}
                          required
                          disabled={!editable || isSeeDetailsSelected}
                        />
                        <span
                          className={classNames(
                            'px-1 py-1 inline-flex items-center bg-gray-100',
                            'rounded-r border border-l-0 border-gray-300 text-gray-600'
                          )}
                        >
                          {'%'}
                        </span>
                      </div>
                      <p className="input-error">
                        {errors.progressiveRollout?.manual?.schedulesList[index]
                          ?.weight?.message && (
                          <span role="alert">
                            {
                              errors.progressiveRollout?.manual?.schedulesList[
                                index
                              ]?.weight?.message
                            }
                          </span>
                        )}
                      </p>
                    </div>
                    <div className="w-full">
                      <DatetimePicker
                        name={`progressiveRollout.manual.schedulesList.${index}.executeAt.time`}
                      />
                      <p className="input-error">
                        {errors.progressiveRollout?.manual?.schedulesList[index]
                          ?.executeAt?.time?.message && (
                          <span role="alert">
                            {
                              errors.progressiveRollout?.manual?.schedulesList[
                                index
                              ]?.executeAt?.time?.message
                            }
                          </span>
                        )}
                      </p>
                    </div>
                    {editable && (
                      <div className="flex items-center">
                        <button
                          type="button"
                          onClick={() => handleRemoveTrigger(index)}
                          className="minus-circle-icon"
                          disabled={manualSchedulesList.length === 1}
                        >
                          <MinusCircleIcon aria-hidden="true" />
                        </button>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </Fragment>
          )}
        </div>
      );
    }
  );
