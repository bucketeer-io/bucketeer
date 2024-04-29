import {
  PAGE_PATH_FEATURES,
  PAGE_PATH_ROOT,
  PAGE_PATH_FEATURE_EXPERIMENTS
} from '../../constants/routing';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  listExperiments,
  selectAll as selectAllExperiment
} from '../../modules/experiments';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { selectAll as selectAllProgressiveRollouts } from '../../modules/porgressiveRollout';
import { ProgressiveRollout } from '../../proto/autoops/progressive_rollout_pb';
import { Experiment } from '../../proto/experiment/experiment_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { AppDispatch } from '../../store';
import { classNames } from '../../utils/css';
import { isArraySorted } from '../../utils/isArraySorted';
import { isIntervals5MinutesApart } from '../../utils/isIntervals5MinutesApart';
import {
  ExclamationCircleIcon,
  InformationCircleIcon,
  MinusCircleIcon,
  PlusIcon
} from '@heroicons/react/outline';
import { SerializedError } from '@reduxjs/toolkit';
import dayjs from 'dayjs';
import React, { memo, FC, useEffect, useCallback, useState } from 'react';
import {
  Controller,
  useFieldArray,
  useFormContext,
  useWatch
} from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';

import { selectById as selectFeatureById } from '../../modules/features';
import { DatetimePicker } from '../DatetimePicker';
import { isExperimentStatusWaitingRunnning } from '../ExperimentList';
import { getIntervalForDayjs } from '../FeatureAutoOpsRulesForm';
import { ProgressiveRolloutTypeTab } from '../OperationAddUpdateForm';
import { Option, Select } from '../Select';
import { OperationForm } from '../../pages/feature/formSchema';

export const isProgressiveRolloutsRunningWaiting = (
  status: ProgressiveRollout.StatusMap[keyof ProgressiveRollout.StatusMap]
) =>
  status === ProgressiveRollout.Status.RUNNING ||
  status === ProgressiveRollout.Status.WAITING;

interface isProgressiveRolloutsWarningsExists {
  feature: Feature.AsObject;
  progressiveRolloutList: ProgressiveRollout.AsObject[];
  experiments: Experiment.AsObject[];
}

export const isProgressiveRolloutsWarningsExists = ({
  feature,
  progressiveRolloutList,
  experiments
}: isProgressiveRolloutsWarningsExists): boolean => {
  const check =
    feature.variationsList.length !== 2 ||
    (experiments.length > 0 &&
      experiments.find((e) => isExperimentStatusWaitingRunnning(e.status))) ||
    (progressiveRolloutList.length > 0 &&
      progressiveRolloutList.find((p) =>
        isProgressiveRolloutsRunningWaiting(p.status)
      ));
  return !!check;
};

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
      setProgressiveRolloutTypeList
    }) => {
      const { formatMessage: f } = useIntl();
      const editable = useIsEditable();
      const history = useHistory();
      const dispatch = useDispatch<AppDispatch>();
      const currentEnvironment = useCurrentEnvironment();
      const [isLoading, setIsLoading] = useState(true);

      const methods = useFormContext<OperationForm>();
      const { control } = methods;

      const [feature] = useSelector<
        AppState,
        [Feature.AsObject | undefined, SerializedError | null]
      >((state) => [
        selectFeatureById(state.features, featureId),
        state.features.getFeatureError
      ]);

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

      const experiments = useSelector<AppState, Experiment.AsObject[]>(
        (state) => selectAllExperiment(state.experiments),
        shallowEqual
      );

      useEffect(() => {
        dispatch(
          listExperiments({
            featureId: featureId,
            environmentNamespace: currentEnvironment.id,
            searchKeyword: '',
            pageSize: 1000,
            cursor: ''
          })
        ).then(() => setIsLoading(false));
      }, [dispatch, featureId, currentEnvironment]);

      const isTemplateSelected =
        progressiveRolloutTypeList.find((p) => p.selected).value ===
        ProgressiveRollout.Type.TEMPLATE_SCHEDULE;

      if (isLoading) {
        return <div className="spinner mt-2 mx-auto"></div>;
      }

      if (
        isProgressiveRolloutsWarningsExists({
          feature,
          progressiveRolloutList,
          experiments
        })
      ) {
        return (
          <div className="rounded-md bg-yellow-50 p-4 mt-2">
            <div className="flex">
              <div className="flex-shrink-0">
                <ExclamationCircleIcon
                  className="h-5 w-5 text-yellow-400"
                  aria-hidden="true"
                />
              </div>
              <div className="ml-3 flex-1">
                <p className="text-sm text-yellow-700 font-semibold">
                  {f(messages.autoOps.progressiveRolloutWarningMessages.title)}
                </p>
                <div className="mt-2 text-sm text-yellow-700">
                  <ul className="list-disc space-y-1 pl-5">
                    {feature.variationsList.length !== 2 ? (
                      <li>
                        <p>
                          {f(
                            messages.autoOps.progressiveRolloutWarningMessages
                              .variations
                          )}
                        </p>
                      </li>
                    ) : null}
                    {experiments.find((e) =>
                      isExperimentStatusWaitingRunnning(e.status)
                    ) ? (
                      <li>
                        <p>
                          {f(
                            messages.autoOps.progressiveRolloutWarningMessages
                              .experimentOnProgress,
                            {
                              link: (
                                <span
                                  onClick={() => {
                                    history.push(
                                      `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_EXPERIMENTS}`
                                    );
                                  }}
                                  className="underline text-primary cursor-pointer ml-1"
                                >
                                  <span>
                                    {f(messages.sourceType.experiment)}
                                  </span>
                                </span>
                              )
                            }
                          )}
                        </p>
                      </li>
                    ) : null}
                    {progressiveRolloutList.length > 0 &&
                    progressiveRolloutList.find((p) =>
                      isProgressiveRolloutsRunningWaiting(p.status)
                    ) ? (
                      <li>
                        <p>
                          {f(
                            messages.autoOps.progressiveRolloutWarningMessages
                              .alreadyProgressiveRollout
                          )}
                        </p>
                      </li>
                    ) : null}
                  </ul>
                </div>
                <p className="text-yellow-700 text-sm mt-4">
                  {f(
                    messages.autoOps.progressiveRolloutWarningMessages
                      .moreInformation,
                    {
                      link: (
                        <a
                          href="https://docs.bucketeer.io/feature-flags/creating-feature-flags/auto-operation/rollout"
                          target="_blank"
                          rel="noreferrer"
                          className="underline text-primary"
                        >
                          {f(messages.feature.documentation)}
                        </a>
                      )
                    }
                  )}
                </p>
              </div>
            </div>
          </div>
        );
      }

      return (
        <div className="mt-4 h-full flex flex-col overflow-hidden">
          {!feature.enabled && (
            <div className="bg-blue-50 p-4 border-l-4 border-blue-400 mb-7">
              <div className="flex">
                <div className="flex-shrink-0">
                  <InformationCircleIcon
                    className="h-5 w-5 text-blue-400"
                    aria-hidden="true"
                  />
                </div>
                <div className="ml-3 flex-1">
                  <p className="text-sm text-blue-700">
                    It will enable the flag when it starts
                  </p>
                </div>
              </div>
            </div>
          )}
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
                      selected: p.label === label
                    }))
                  );
                }}
              >
                {label}
              </div>
            ))}
          </div>
          <div className="mt-4  px-[2px]">
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
    const methods = useFormContext<OperationForm>();
    const {
      control,
      formState: { errors },
      register,
      setValue
    } = methods;
    const editable = useIsEditable();

    const {
      template: { schedulesList, increments, interval, datetime }
    } = useWatch({
      control,
      name: 'progressiveRollout'
    });

    useEffect(() => {
      if (
        progressiveRolloutTypeList.find((p) => p.selected).value ===
          ProgressiveRollout.Type.TEMPLATE_SCHEDULE &&
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
                time: time
              },
              weight: weight > 100 ? 100 : Math.round(weight * 100) / 100
            };
          });
        setValue('progressiveRollout.template.schedulesList', scheduleList);
      }
    }, [datetime.time, interval, increments]);

    const intervalOptions = [
      {
        label: f(messages.autoOps.hourly),
        value: '1'
      },
      {
        label: f(messages.autoOps.daily),
        value: '2'
      },
      {
        label: f(messages.autoOps.weekly),
        value: '3'
      }
    ];

    return (
      <div className="mt-4 h-full flex flex-col overflow-hidden px-[2px]">
        <div>
          <span className="input-label">{f(messages.autoOps.startDate)}</span>
          <DatetimePicker
            name="progressiveRollout.template.datetime.time"
            dateFormat="yyyy/MM/dd HH:mm"
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
            <span className="input-label">{f(messages.autoOps.increment)}</span>
            <div className="flex">
              <input
                type="number"
                {...register('progressiveRollout.template.increments')}
                min="0"
                max="100"
                onKeyDown={(evt) => {
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
            <span className="input-label">{f(messages.autoOps.frequency)}</span>
            <Controller
              name="progressiveRollout.template.interval"
              control={control}
              render={({ field }) => (
                <Select
                  onChange={(o: Option) => field.onChange(o.value)}
                  options={intervalOptions}
                  disabled={!editable || isSeeDetailsSelected}
                  value={intervalOptions.find(
                    (o) => o.value === field.value.toString()
                  )}
                />
              )}
            />
          </div>
        </div>
        <div className="mt-4 flex flex-col h-full overflow-hidden">
          <div className="space-y-2 mt-2 overflow-y-auto flex flex-col h-full pb-6">
            {schedulesList?.map((_, index) => (
              <div key={index}>
                {index === 0 && (
                  <div className="flex space-x-4 mb-2">
                    <div className="w-1/2 input-label">
                      {f(messages.autoOps.weight)}
                    </div>
                    <div className="w-1/2 input-label">
                      {f(messages.autoOps.executeAt)}
                    </div>
                  </div>
                )}
                <div className="flex space-x-4">
                  <div className="flex w-1/2">
                    <input
                      {...register(
                        `progressiveRollout.template.schedulesList.${index}.weight`
                      )}
                      onKeyDown={(evt) => {
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
                        dateFormat="yyyy/MM/dd HH:mm"
                        disabled={true}
                      />
                    </div>
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
    const methods = useFormContext<OperationForm>();
    const {
      control,
      formState: { errors },
      register
    } = methods;
    const editable = useIsEditable();

    const watchManualSchedulesList = useWatch({
      control,
      name: 'progressiveRollout.manual.schedulesList'
    });

    const {
      fields: manualSchedulesList,
      remove: removeTrigger,
      append
    } = useFieldArray({
      control,
      name: 'progressiveRollout.manual.schedulesList'
    });

    const handleAddOperation = (e) => {
      e.preventDefault();

      const lastSchedule =
        watchManualSchedulesList[watchManualSchedulesList.length - 1];
      const time = dayjs(lastSchedule?.executeAt.time)
        .add(5, 'minute')
        .toDate();

      let weight = lastSchedule ? Number(lastSchedule.weight) : 0;

      if (weight >= 90) {
        weight = 100;
      } else {
        weight = weight + 10;
      }

      append({
        executeAt: {
          time
        },
        weight
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
    const isDatetime5MinutesApart = isIntervals5MinutesApart(
      watchManualSchedulesList.map((d) => d.executeAt.time.getTime())
    );

    return (
      <div className="mt-4 h-full flex flex-col overflow-hidden px-[2px]">
        <div className="space-y-2 flex flex-col overflow-y-auto h-full mt-2">
          {manualSchedulesList.map((schedule, index) => (
            <div key={schedule.id}>
              {index === 0 && (
                <div className="flex space-x-4 mb-2">
                  <div className="w-1/2 input-label">
                    {f(messages.autoOps.weight)}
                  </div>
                  <div className="w-1/2 input-label">
                    {f(messages.autoOps.executeAt)}
                  </div>
                </div>
              )}
              <div className="flex space-x-4 pr-1">
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
                          errors.progressiveRollout?.manual?.schedulesList[
                            index
                          ]?.weight?.message
                        }
                      </span>
                    )}
                  </p>
                </div>
                <div className="w-1/2 flex">
                  <div>
                    <DatetimePicker
                      name={`progressiveRollout.manual.schedulesList.${index}.executeAt.time`}
                      dateFormat="yyyy/MM/dd HH:mm"
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
            </div>
          ))}
          {watchManualSchedulesList.length <= 10 && (
            <ErrorMessage
              isWeightsSorted={isWeightsSorted}
              isDatesSorted={isDatesSorted}
              isDatetime5MinutesApart={isDatetime5MinutesApart}
            />
          )}
          <div className="py-3">
            <button
              onClick={handleAddOperation}
              className={classNames(
                'text-primary space-x-2 flex items-center self-start',
                (isLastScheduleWeight100 ||
                  !isWeightsSorted ||
                  !isDatesSorted) &&
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
          </div>
        </div>
        {watchManualSchedulesList.length > 10 && (
          <ErrorMessage
            isWeightsSorted={isWeightsSorted}
            isDatesSorted={isDatesSorted}
            isDatetime5MinutesApart={isDatetime5MinutesApart}
          />
        )}
      </div>
    );
  }
);

interface ErrorMessageProps {
  isWeightsSorted: boolean;
  isDatesSorted: boolean;
  isDatetime5MinutesApart: boolean;
}

const ErrorMessage: FC<ErrorMessageProps> = memo(
  ({ isWeightsSorted, isDatesSorted, isDatetime5MinutesApart }) => {
    const { formatMessage: f } = useIntl();

    if (isWeightsSorted && isDatesSorted && isDatetime5MinutesApart) {
      return null;
    }

    return (
      <div className="flex space-x-2">
        <div className="flex-1">
          {!isWeightsSorted && (
            <p className="input-error">
              <span role="alert">
                {f(messages.autoOps.weightIncreasingOrder)}
              </span>
            </p>
          )}
        </div>
        <div className="flex-1">
          {!isDatesSorted ? (
            <p className="input-error">
              <span role="alert">
                {f(messages.autoOps.dateIncreasingOrder)}
              </span>
            </p>
          ) : !isDatetime5MinutesApart ? (
            !isDatetime5MinutesApart && (
              <p className="input-error">
                <span role="alert">
                  {f(messages.autoOps.timeInterval5MinutesApart)}
                </span>
              </p>
            )
          ) : null}
        </div>
      </div>
    );
  }
);
