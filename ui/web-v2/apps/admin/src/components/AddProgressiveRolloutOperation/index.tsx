import { messages } from '@/lang/messages';
import { AppState } from '@/modules';
import { selectAll as selectAllAutoOpsRules } from '@/modules/autoOpsRules';
import { useIsEditable } from '@/modules/me';
import { selectAll as selectAllProgressiveRollouts } from '@/modules/porgressiveRollout';
import { AutoOpsRule } from '@/proto/autoops/auto_ops_rule_pb';
import { ProgressiveRollout } from '@/proto/autoops/progressive_rollout_pb';
import { classNames } from '@/utils/css';
import { isArraySorted } from '@/utils/isArraySorted';
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
import { shallowEqual, useSelector } from 'react-redux';

import { DatetimePicker } from '../DatetimePicker';
import {
  ClauseType,
  ProgressiveRolloutClauseType,
  getIntervalForDayjs,
} from '../FeatureAutoOpsRulesForm';
import { ProgressiveRolloutTypeTab } from '../OperationAddUpdateForm';
import { Option, Select } from '../Select';

interface AddProgressiveRolloutOperationProps {
  featureId: string;
  variationOptions: Option[];
  isSeeDetailsSelected: boolean;
  progressiveRolloutTypeList: ProgressiveRolloutTypeTab[];
  setProgressiveRolloutTypeList: (arg) => void;
}

export const AddProgressiveRolloutOperation: FC<AddProgressiveRolloutOperationProps> =
  memo(
    ({
      featureId,
      variationOptions,
      isSeeDetailsSelected,
      progressiveRolloutTypeList,
      setProgressiveRolloutTypeList,
    }) => {
      const { formatMessage: f } = useIntl();
      const editable = useIsEditable();

      const methods = useFormContext<any>();
      const { control } = methods;

      const autoOpsRules = useSelector<AppState, AutoOpsRule.AsObject[]>(
        (state) =>
          selectAllAutoOpsRules(state.autoOpsRules).filter(
            (rule) => rule.featureId === featureId
          ),
        shallowEqual
      );

      const progressiveRolloutList = useSelector<
        AppState,
        ProgressiveRollout.AsObject[]
      >(
        (state) =>
          selectAllProgressiveRollouts(state.progressiveRollout).filter(
            (rule) => rule.featureId === featureId
          ),
        shallowEqual
      );

      const isTemplateSelected =
        progressiveRolloutTypeList.find((p) => p.selected).value ===
        ProgressiveRolloutClauseType.PROGRESSIVE_ROLLOUT_TEMPLATE_SCHEDULE;

      const findScheduleOperation = autoOpsRules.find((rule) => {
        const { typeUrl } = rule.clausesList[0].clause;
        const type = typeUrl.substring(typeUrl.lastIndexOf('/') + 1);
        return type === ClauseType.DATETIME;
      });

      if (findScheduleOperation) {
        return (
          <p className="input-error mt-2">
            There is a schedule configured in the auto operations. Please delete
            it before using the progressive rollout.
          </p>
        );
      }

      if (progressiveRolloutList.length > 0) {
        return (
          <p className="input-error mt-2">
            There is already progressive rollout configured in the auto
            operations.
          </p>
        );
      }

      return (
        <div className="mt-4 h-full flex flex-col overflow-hidden">
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
          <div className="mt-4">
            <span className="input-label">{f(messages.feature.variation)}</span>
            <Controller
              key={
                isTemplateSelected ? 'templateVariationId' : 'manualVariationId'
              }
              name={
                isTemplateSelected
                  ? 'progressiveRollout.template.variationId'
                  : 'progressiveRollout.manual.variationId'
              }
              control={control}
              render={({ field }) => (
                <Select
                  onChange={(o: Option) => field.onChange(o.value)}
                  options={variationOptions}
                  disabled={!editable || isSeeDetailsSelected}
                  value={variationOptions.find((o) => o.value === field.value)}
                />
              )}
            />
          </div>
          {isTemplateSelected ? (
            <TemplateProgressiveRollout
              isSeeDetailsSelected={isSeeDetailsSelected}
              progressiveRolloutTypeList={progressiveRolloutTypeList}
            />
          ) : (
            <ManualProgressiveRollout
              isSeeDetailsSelected={isSeeDetailsSelected}
            />
          )}
        </div>
      );
    }
  );

interface TemplateProgressiveRolloutProps {
  isSeeDetailsSelected: boolean;
  progressiveRolloutTypeList: ProgressiveRolloutTypeTab[];
}

const TemplateProgressiveRollout: FC<TemplateProgressiveRolloutProps> = memo(
  ({ isSeeDetailsSelected, progressiveRolloutTypeList }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext<any>();
    const {
      control,
      formState: { errors },
      register,
      setValue,
    } = methods;
    const editable = useIsEditable();

    const {
      template: { schedulesList, increments, interval, datetime },
    } = useWatch({
      control,
      name: 'progressiveRollout',
    });

    useEffect(() => {
      if (
        progressiveRolloutTypeList.find((p) => p.selected).value ===
          ProgressiveRolloutClauseType.PROGRESSIVE_ROLLOUT_TEMPLATE_SCHEDULE &&
        Number(increments) > 0
      ) {
        const scheduleList = Array(Math.ceil(100 / increments))
          .fill('')
          .map((_, index) => {
            // increment each schedule by {interval}
            const time = dayjs(datetime.time)
              .add(index, getIntervalForDayjs(interval))
              .toDate();

            const weight = (index + 1) * increments;

            return {
              executeAt: {
                time: time,
              },
              weight: weight > 100 ? 100 : Math.round(weight * 100) / 100,
            };
          });
        setValue('progressiveRollout.template.schedulesList', scheduleList);
      }
    }, [datetime.time, interval, increments]);

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

    return (
      <div className="mt-4 h-full flex flex-col overflow-hidden">
        <div>
          <span className="input-label">{f(messages.autoOps.startDate)}</span>
          <DatetimePicker
            name="progressiveRollout.template.datetime.time"
            disabled={isSeeDetailsSelected}
          />
          <p className="input-error">
            {errors.progressiveRollout?.template?.datetime?.time?.message && (
              <span role="alert">
                {errors.progressiveRollout?.template?.datetime?.time?.message}
              </span>
            )}
          </p>
        </div>
        <div className="flex space-x-4 mt-4">
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
              {errors.progressiveRollout?.template?.increments?.message && (
                <span role="alert">
                  {errors.progressiveRollout?.template?.increments?.message}
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
                  value={intervalOptions.find((o) => o.value === field.value)}
                />
              )}
            />
          </div>
        </div>
        <div className="mt-4 flex flex-col h-full overflow-hidden">
          <div className="flex">
            <div className="flex-1 input-label">Weight</div>
            <div className="flex-1 input-label">Execute at</div>
          </div>
          <div className="space-y-2 mt-2 overflow-y-scroll flex flex-col h-full pb-6">
            {schedulesList?.map((_, index) => (
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
      </div>
    );
  }
);

interface ManualProgressiveRolloutProps {
  isSeeDetailsSelected: boolean;
}

const ManualProgressiveRollout: FC<ManualProgressiveRolloutProps> = memo(
  ({ isSeeDetailsSelected }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext<any>();
    const {
      control,
      formState: { errors },
      register,
    } = methods;
    const editable = useIsEditable();

    const watchManualSchedulesList = useWatch({
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

    const handleAddOperation = (e) => {
      e.preventDefault();

      const lastSchedule: any =
        watchManualSchedulesList[watchManualSchedulesList.length - 1];
      const time = dayjs(lastSchedule?.executeAt.time)
        .add(5, 'minute')
        .toDate();

      let weight = lastSchedule ? Number(lastSchedule.weight) : 0;

      if (weight >= 80) {
        weight = 100;
      } else {
        weight = weight + 20;
      }

      append({
        executeAt: {
          time,
        },
        weight,
      });
    };

    const handleRemoveTrigger = useCallback(
      (idx) => {
        removeTrigger(idx);
      },
      [removeTrigger]
    );

    const isLastScheduleWeight100 =
      Number(
        watchManualSchedulesList[watchManualSchedulesList.length - 1]?.weight
      ) === 100;

    const isWeightsSorted = isArraySorted(
      watchManualSchedulesList.map((d) => Number(d.weight))
    );
    const isDatesSorted = isArraySorted(
      watchManualSchedulesList.map((d) => d.executeAt.time.getTime())
    );

    return (
      <div className="mt-4 h-full flex flex-col overflow-hidden">
        <button
          onClick={handleAddOperation}
          className={classNames(
            'text-primary space-x-2 flex items-center',
            (isLastScheduleWeight100 || !isWeightsSorted || !isDatesSorted) &&
              'opacity-50 cursor-not-allowed'
          )}
          disabled={
            isLastScheduleWeight100 || !isWeightsSorted || !isDatesSorted
          }
        >
          <PlusIcon width={16} />
          <span className="text-sm font-medium">
            {f(messages.button.addOperation)}
          </span>
        </button>
        <div className="flex mt-3">
          <div className="flex-1 input-label">Weight</div>
          <div className="flex-1 input-label">Execute at</div>
        </div>
        <div className="space-y-2 flex flex-col overflow-y-scroll h-full mt-2">
          {manualSchedulesList.map((_, index) => (
            <div key={index} className="flex space-x-4 pr-1">
              <div className="w-1/2">
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
                      errors.progressiveRollout?.manual?.schedulesList[index]
                        ?.weight?.message
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
                        errors.progressiveRollout?.manual?.schedulesList[index]
                          ?.weight?.message
                      }
                    </span>
                  )}
                </p>
              </div>
              <div className="w-1/2 flex">
                <div>
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
                  <div className="flex items-center ml-2">
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
            </div>
          ))}
          {watchManualSchedulesList.length <= 10 && (
            <ErrorMessage
              isWeightsSorted={isWeightsSorted}
              isDatesSorted={isDatesSorted}
            />
          )}
        </div>
        {watchManualSchedulesList.length > 10 && (
          <ErrorMessage
            isWeightsSorted={isWeightsSorted}
            isDatesSorted={isDatesSorted}
          />
        )}
      </div>
    );
  }
);

interface ErrorMessageProps {
  isWeightsSorted: boolean;
  isDatesSorted: boolean;
}

const ErrorMessage: FC<ErrorMessageProps> = memo(
  ({ isWeightsSorted, isDatesSorted }) => (
    <div className="flex pb-6 pt-2 space-x-2 pr-6">
      <div className="flex-1">
        {!isWeightsSorted && (
          <p className="input-error">
            <span role="alert">
              The weights need to be in increasing order.
            </span>
          </p>
        )}
      </div>
      <div className="flex-1">
        {!isDatesSorted && (
          <p className="input-error">
            <span role="alert">The dates need to be in increasing order.</span>
          </p>
        )}
      </div>
    </div>
  )
);
