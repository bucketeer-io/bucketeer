import { intl } from '@/lang';
import { AppState } from '@/modules';
import { useCurrentEnvironment, useIsEditable } from '@/modules/me';
import { addFormSchema } from '@/pages/goal/formSchema';
import { AutoOpsRule, OpsType } from '@/proto/autoops/auto_ops_rule_pb';
import {
  DatetimeClause,
  OpsEventRateClause,
  ProgressiveRolloutManualScheduleClause,
  ProgressiveRolloutTemplateScheduleClause,
  ProgressiveRolloutSchedule,
} from '@/proto/autoops/clause_pb';
import {
  ChangeAutoOpsRuleOpsTypeCommand,
  ChangeDatetimeClauseCommand,
  ChangeOpsEventRateClauseCommand,
  CreateAutoOpsRuleCommand,
  CreateProgressiveRolloutCommand,
} from '@/proto/autoops/command_pb';
import { Goal } from '@/proto/experiment/goal_pb';
import { ListGoalsRequest } from '@/proto/experiment/service_pb';
import { Feature } from '@/proto/feature/feature_pb';
import { AppDispatch } from '@/store';
import { classNames } from '@/utils/css';
import { Dialog, Transition } from '@headlessui/react';
import {
  XIcon,
  ExclamationCircleIcon,
  PlusIcon,
  MinusCircleIcon,
} from '@heroicons/react/outline';
import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import dayjs from 'dayjs';
import React, {
  FC,
  Fragment,
  memo,
  useCallback,
  useEffect,
  useState,
} from 'react';
import {
  Controller,
  useFieldArray,
  useForm,
  useFormContext,
  useWatch,
} from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import ReactSelect, { components } from 'react-select';
import { v4 as uuid } from 'uuid';

import { messages } from '../../lang/messages';
import {
  createAutoOpsRule,
  updateAutoOpsRule,
  UpdateAutoOpsRuleParams,
} from '../../modules/autoOpsRules';
import { selectById as selectFeatureById } from '../../modules/features';
import {
  createGoal,
  listGoals,
  selectAll as selectAllGoals,
} from '../../modules/goals';
import { createProgressiveRollout } from '../../modules/porgressiveRollout';
import { DatetimePicker, ReactDatePicker } from '../DatetimePicker';
import {
  ClauseType,
  createInitialDatetimeClause,
  operatorOptions,
  ProgressiveRolloutClauseType,
} from '../FeatureAutoOpsRulesForm';
import { Option, Select } from '../Select';

interface ExecuteAt {
  time: Date;
}
interface FormValues {
  executeAt: ExecuteAt;
  weight: number;
}

export interface OperationAddUpdateFormProps {
  featureId: string;
  onSubmit: () => void;
  onSubmitProgressiveRollout: () => void;
  onCancel: () => void;
  autoOpsRule?: AutoOpsRule.AsObject;
  isKillSwitchSelected: boolean;
  isActiveTabSelected: boolean;
}

const TabLabel = {
  ENABLE: intl.formatMessage(messages.autoOps.enable),
  KILL_SWITCH: intl.formatMessage(messages.autoOps.killSwitch),
};

export const OperationAddUpdateForm: FC<OperationAddUpdateFormProps> = memo(
  ({
    onSubmit,
    onSubmitProgressiveRollout,
    onCancel,
    featureId,
    autoOpsRule,
    isKillSwitchSelected,
    isActiveTabSelected,
  }) => {
    const editable = useIsEditable();
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();

    const { formatMessage: f } = useIntl();
    const [feature, getFeatureError] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ]);

    const [isAddGoalOpen, setIsAddGoalOpen] = useState(false);
    const [radioList, setRadioList] = useState([]);

    const [progressiveRolloutTypeList, setProgressiveRolloutTypeList] =
      useState([
        {
          label: 'Template',
          value:
            ProgressiveRolloutClauseType.PROGRESSIVE_ROLLOUT_TEMPLATE_SCHEDULE,
          selected: true,
        },
        {
          label: 'Manual',
          value:
            ProgressiveRolloutClauseType.PROGRESSIVE_ROLLOUT_MANUAL_SCHEDULE,
          selected: false,
        },
      ]);

    const methods = useFormContext<any>();
    const {
      handleSubmit,
      control,
      formState: { errors, isValid, isSubmitting },
      register,
      setValue,
    } = methods;

    console.log('errors', errors);
    const opsType = useWatch({
      control,
      name: 'opsType',
    });

    const clauseType = useWatch({
      control,
      name: 'clauseType',
    });

    const templateSchedulesList = useWatch({
      control,
      name: 'progressiveRollout.template.schedulesList',
    });

    const incrementsWatch = useWatch({
      control,
      name: 'progressiveRollout.template.increments',
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

    console.log('manualSchedulesList', manualSchedulesList);
    const templateIntervalWatch = useWatch({
      control,
      name: 'progressiveRollout.template.interval',
    });

    const progressiveRolloutTemplateDatetimeTimeWatch = useWatch({
      control,
      name: 'progressiveRollout.template.datetime.time',
    });

    const tabs = [
      {
        label: TabLabel.ENABLE,
        value: OpsType.ENABLE_FEATURE,
      },
      {
        label: TabLabel.KILL_SWITCH,
        value: OpsType.DISABLE_FEATURE,
      },
    ];

    const isSeeDetailsSelected = autoOpsRule && !isActiveTabSelected;

    const setEnableList = () => {
      setRadioList([
        {
          label: f(messages.autoOps.schedule),
          value: ClauseType.DATETIME,
        },
        {
          label: 'Progressive Rollout',
          value: ClauseType.PROGRESSIVE_ROLLOUT,
        },
      ]);
    };

    const setKillSwitchList = () => {
      setRadioList([
        {
          label: f(messages.autoOps.schedule),
          value: ClauseType.DATETIME,
        },
        {
          label: f(messages.autoOps.eventRate),
          value: ClauseType.EVENT_RATE,
        },
      ]);
    };

    const getInterval = (
      interval: ProgressiveRolloutTemplateScheduleClause.IntervalMap[keyof ProgressiveRolloutTemplateScheduleClause.IntervalMap]
    ) => {
      if (Number(interval) === 1) {
        return 'hour';
      } else if (Number(interval) === 2) {
        return 'day';
      } else if (Number(interval) === 3) {
        return 'week';
      }
    };

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
            const datetime = dayjs(progressiveRolloutTemplateDatetimeTimeWatch)
              .add(index, getInterval(templateIntervalWatch))
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

    useEffect(() => {
      if (autoOpsRule) {
        const typeUrl = autoOpsRule.clausesList[0].clause.typeUrl;
        const type = typeUrl.substring(typeUrl.lastIndexOf('/') + 1);

        setValue('opsType', autoOpsRule.opsType);
        setValue('clauseType', type);

        if (autoOpsRule.opsType === OpsType.ENABLE_FEATURE) {
          setEnableList();
        } else {
          setKillSwitchList();
        }

        if (type === ClauseType.DATETIME) {
          const datetime = DatetimeClause.deserializeBinary(
            autoOpsRule.clausesList[0].clause.value as Uint8Array
          ).toObject();

          setValue('datetime.time', new Date(datetime.time * 1000));
        } else if (type === ClauseType.EVENT_RATE) {
          const opsEventRateClause = OpsEventRateClause.deserializeBinary(
            autoOpsRule.clausesList[0].clause.value as Uint8Array
          ).toObject();

          setValue('eventRate.variation', opsEventRateClause.variationId);
          setValue('eventRate.goal', opsEventRateClause.goalId);
          setValue(
            'eventRate.operator',
            opsEventRateClause.operator.toString()
          );
          setValue(
            'eventRate.threadsholdRate',
            opsEventRateClause.threadsholdRate * 100
          );
          setValue('eventRate.minCount', opsEventRateClause.minCount);
        }
      } else if (isKillSwitchSelected) {
        setValue('opsType', OpsType.DISABLE_FEATURE);
        setKillSwitchList();
      } else {
        setEnableList();
      }
    }, [autoOpsRule, isKillSwitchSelected]);

    const handleOnSubmit = useCallback(
      (data) => {
        console.log('data', data);
        if (autoOpsRule) {
          const changeAutoOpsRuleOpsTypeCommand =
            new ChangeAutoOpsRuleOpsTypeCommand();

          data.opsType === OpsType.ENABLE_FEATURE.toString()
            ? changeAutoOpsRuleOpsTypeCommand.setOpsType(OpsType.ENABLE_FEATURE)
            : changeAutoOpsRuleOpsTypeCommand.setOpsType(
                OpsType.DISABLE_FEATURE
              );

          const changeDatetimeClauseCommands: ChangeDatetimeClauseCommand[] =
            [];
          const changeOpsEventRateClauseCommands: ChangeOpsEventRateClauseCommand[] =
            [];

          if (data.clauseType === ClauseType.DATETIME) {
            const clause = new DatetimeClause();
            clause.setTime(Math.round(data.datetime.time.getTime() / 1000));
            const command = new ChangeDatetimeClauseCommand();
            command.setId(autoOpsRule.clausesList[0].id);
            command.setDatetimeClause(clause);
            changeDatetimeClauseCommands.push(command);
          }
          if (data.clauseType === ClauseType.EVENT_RATE) {
            const command = new ChangeOpsEventRateClauseCommand();
            command.setId(autoOpsRule.clausesList[0].id);
            command.setOpsEventRateClause(
              createOpsEventRateClause(data.eventRate)
            );
            changeOpsEventRateClauseCommands.push(command);
          }

          const param: UpdateAutoOpsRuleParams = {
            environmentNamespace: currentEnvironment.id,
            id: autoOpsRule.id,
            changeDatetimeClauseCommands,
            changeOpsEventRateClauseCommands,
            changeAutoOpsRuleOpsTypeCommand,
          };

          dispatch(updateAutoOpsRule(param)).then(() => onSubmit());
        } else {
          if (
            data.clauseType === ClauseType.DATETIME ||
            data.clauseType === ClauseType.EVENT_RATE
          ) {
            const command = new CreateAutoOpsRuleCommand();
            command.setFeatureId(featureId);

            data.opsType === OpsType.ENABLE_FEATURE.toString()
              ? command.setOpsType(OpsType.ENABLE_FEATURE)
              : command.setOpsType(OpsType.DISABLE_FEATURE);

            if (data.clauseType === ClauseType.DATETIME) {
              const clause = new DatetimeClause();
              clause.setTime(Math.round(data.datetime.time.getTime() / 1000));
              command.setDatetimeClausesList([clause]);
            }
            if (data.clauseType === ClauseType.EVENT_RATE) {
              command.setOpsEventRateClausesList([
                createOpsEventRateClause(data.eventRate),
              ]);
            }

            dispatch(
              createAutoOpsRule({
                environmentNamespace: currentEnvironment.id,
                command: command,
              })
            ).then(() => onSubmit());
          }

          if (data.clauseType === ClauseType.PROGRESSIVE_ROLLOUT) {
            const command = new CreateProgressiveRolloutCommand();
            command.setFeatureId(featureId);

            const selectedProgressiveRolloutType =
              progressiveRolloutTypeList.find((p) => p.selected).value;

            if (
              selectedProgressiveRolloutType ===
              ProgressiveRolloutClauseType.PROGRESSIVE_ROLLOUT_TEMPLATE_SCHEDULE
            ) {
              const {
                progressiveRollout: {
                  template: {
                    increments,
                    interval,
                    schedulesList,
                    variationId,
                  },
                },
              } = data;

              const clause = new ProgressiveRolloutTemplateScheduleClause();

              clause.setIncrements(increments);
              clause.setInterval(interval);

              clause.setVariationId(variationId);

              const scheduleList = schedulesList.map((schedule) => {
                const c = new ProgressiveRolloutSchedule();

                c.setExecuteAt(
                  Math.round(schedule.executeAt.time.getTime() / 1000)
                );
                c.setWeight(schedule.weight * 1000);
                return c;
              });

              clause.setSchedulesList(scheduleList);

              command.setProgressiveRolloutTemplateScheduleClause(clause);
            } else if (
              selectedProgressiveRolloutType ===
              ProgressiveRolloutClauseType.PROGRESSIVE_ROLLOUT_MANUAL_SCHEDULE
            ) {
              const {
                progressiveRollout: {
                  manual: { schedulesList, variationId },
                },
              } = data;

              const clause = new ProgressiveRolloutManualScheduleClause();

              clause.setVariationId(variationId);

              const scheduleList = schedulesList.map((schedule) => {
                const c = new ProgressiveRolloutSchedule();

                c.setExecuteAt(
                  Math.round(schedule.executeAt.time.getTime() / 1000)
                );
                c.setWeight(schedule.weight * 1000);
                return c;
              });

              clause.setSchedulesList(scheduleList);
              command.setProgressiveRolloutManualScheduleClause(clause);
            }

            dispatch(
              createProgressiveRollout({
                environmentNamespace: currentEnvironment.id,
                command: command,
              })
            ).then(() => {
              onSubmitProgressiveRollout();
            });
          }
        }
      },
      [autoOpsRule, progressiveRolloutTypeList]
    );

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

    const variationOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: v.value,
      };
    });

    const goals = useSelector<AppState, Goal.AsObject[]>(
      (state) => selectAllGoals(state.goals),
      shallowEqual
    );

    const goalOptions = goals.map((goal) => {
      return {
        value: goal.id,
        label: goal.id,
      };
    });

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

    const title = () => {
      if (isSeeDetailsSelected) {
        return f(messages.autoOps.operationDetails);
      } else {
        return autoOpsRule
          ? f(messages.autoOps.updateAnOperation)
          : f(messages.autoOps.createAnOperation);
      }
    };

    return (
      <div className="w-[500px] h-full">
        <AddGoalModal open={isAddGoalOpen} setOpen={setIsAddGoalOpen} />
        <form className="flex flex-col h-full">
          <div className="flex-1 h-0">
            <div className="flex items-center justify-between px-4 py-5 border-b">
              <p className="text-xl font-medium">{title()}</p>
              <XIcon
                width={20}
                className="text-gray-400 cursor-pointer"
                onClick={onCancel}
              />
            </div>
            <div className="px-4 flex-1">
              <div className="flex border-b border-gray-100">
                {tabs.map((tab) => (
                  <div
                    {...register('opsType')}
                    key={tab.label}
                    className={classNames(
                      'py-3 flex-1 text-center',
                      opsType === tab.value
                        ? 'text-primary border-b-2 border-primary'
                        : 'text-gray-400',
                      isSeeDetailsSelected
                        ? 'cursor-not-allowed'
                        : 'cursor-pointer'
                    )}
                    onClick={() => {
                      if (isSeeDetailsSelected) {
                        return;
                      }
                      setValue('opsType', tab.value);
                      setValue('clauseType', ClauseType.DATETIME);

                      if (tab.value === OpsType.ENABLE_FEATURE) {
                        setEnableList();
                      } else {
                        setKillSwitchList();
                      }
                    }}
                  >
                    {tab.label}
                  </div>
                ))}
              </div>
              <div className="mt-6">
                {radioList.map((radio) => (
                  <div key={radio.label} className="mb-4 flex space-x-4">
                    <input
                      {...register('clauseType')}
                      id={radio.label}
                      type="radio"
                      value={radio.value}
                      className="h-4 w-4 text-primary focus:ring-primary border-gray-300 mt-1"
                      disabled={isSeeDetailsSelected}
                    />
                    <div className="flex-1">
                      <label htmlFor={radio.label}>{radio.label}</label>
                      {radio.value === ClauseType.DATETIME &&
                        clauseType === ClauseType.DATETIME && (
                          <div className="mt-1">
                            <span className="input-label">
                              {f(messages.autoOps.startDate)}
                            </span>
                            <DatetimePicker
                              name="datetime.time"
                              disabled={isSeeDetailsSelected}
                            />
                            <p className="input-error">
                              {errors.datetime?.time?.message && (
                                <span role="alert">
                                  {errors.datetime?.time?.message}
                                </span>
                              )}
                            </p>
                          </div>
                        )}
                      {radio.value === ClauseType.EVENT_RATE &&
                        clauseType === ClauseType.EVENT_RATE && (
                          <div className="mt-4 space-y-2">
                            <div className="px-4 py-1 text-pink-500 bg-pink-50 inline-block">
                              If
                            </div>
                            <div>
                              <span className="input-label">
                                {f(messages.feature.variation)}
                              </span>
                              <Controller
                                name="eventRate.variation"
                                control={control}
                                render={({ field }) => (
                                  <Select
                                    onChange={(o: Option) => {
                                      field.onChange(o.value);
                                    }}
                                    options={variationOptions}
                                    disabled={!editable || isSeeDetailsSelected}
                                    value={variationOptions.find(
                                      (o) => o.value === field.value
                                    )}
                                  />
                                )}
                              />
                            </div>
                            <div>
                              <span className="input-label">
                                {f(messages.autoOps.opsEventRateClause.goal)}
                              </span>
                              <Controller
                                name="eventRate.goal"
                                control={control}
                                render={({ field }) => (
                                  <AddGoalSelect
                                    onChange={(o: Option) =>
                                      field.onChange(o.value)
                                    }
                                    options={goalOptions}
                                    disabled={!editable || isSeeDetailsSelected}
                                    value={goalOptions.find(
                                      (o) => o.value === field.value
                                    )}
                                    openAddGoalModal={() =>
                                      setIsAddGoalOpen(true)
                                    }
                                  />
                                )}
                              />
                            </div>
                            <div className="grid grid-cols-3 gap-3">
                              <div>
                                <span className="input-label">
                                  {f(messages.autoOps.condition)}
                                </span>
                                <Controller
                                  name="eventRate.operator"
                                  control={control}
                                  render={({ field }) => (
                                    <Select
                                      onChange={(o: Option) =>
                                        field.onChange(o.value)
                                      }
                                      options={operatorOptions}
                                      disabled={
                                        !editable || isSeeDetailsSelected
                                      }
                                      value={operatorOptions.find(
                                        (o) => o.value === field.value
                                      )}
                                    />
                                  )}
                                />
                              </div>
                              <div>
                                <span className="input-label">
                                  {f(messages.autoOps.threshold)}
                                </span>
                                <div className="flex">
                                  <input
                                    {...register('eventRate.threadsholdRate')}
                                    type="number"
                                    min="0"
                                    max="100"
                                    className={classNames(
                                      'w-full',
                                      errors.eventRate?.threadsholdRate
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
                              </div>
                              <div>
                                <span className="input-label">
                                  {f(
                                    messages.autoOps.opsEventRateClause.minCount
                                  )}
                                </span>
                                <div>
                                  <input
                                    {...register('eventRate.minCount')}
                                    type="number"
                                    min="0"
                                    className={classNames(
                                      'w-full',
                                      errors.eventRate?.minCount
                                        ? 'input-text-error'
                                        : 'input-text'
                                    )}
                                    disabled={!editable || isSeeDetailsSelected}
                                  />
                                </div>
                              </div>
                            </div>
                            <div>
                              {errors.eventRate?.threadsholdRate?.message && (
                                <p className="input-error">
                                  <span role="alert">
                                    {errors.eventRate?.threadsholdRate?.message}
                                  </span>
                                </p>
                              )}
                              {errors.eventRate?.minCount?.message && (
                                <p className="input-error">
                                  <span role="alert">
                                    {errors.eventRate?.minCount?.message}
                                  </span>
                                </p>
                              )}
                            </div>
                          </div>
                        )}
                      {radio.value === ClauseType.PROGRESSIVE_ROLLOUT &&
                        clauseType === ClauseType.PROGRESSIVE_ROLLOUT && (
                          <div className="mt-4 space-y-4">
                            <div className="flex">
                              {progressiveRolloutTypeList.map(
                                ({ label, selected }, index) => (
                                  <div
                                    key={label}
                                    className={classNames(
                                      'flex-1 border py-2 text-center cursor-pointer',
                                      index === 0
                                        ? 'rounded-l-lg'
                                        : 'rounded-r-lg',
                                      selected
                                        ? 'bg-pink-50 border-pink-500'
                                        : ''
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
                                )
                              )}
                            </div>
                            {progressiveRolloutTypeList.find((p) => p.selected)
                              .value ===
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
                                        onChange={(o: Option) =>
                                          field.onChange(o.value)
                                        }
                                        options={variationOptions}
                                        disabled={
                                          !editable || isSeeDetailsSelected
                                        }
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
                                    {errors.datetime?.time?.message && (
                                      <span role="alert">
                                        {errors.datetime?.time?.message}
                                      </span>
                                    )}
                                  </p>
                                </div>
                                <div className="flex space-x-4">
                                  <div className="flex-1">
                                    <span className="input-label">
                                      Increment
                                    </span>
                                    <div className="flex">
                                      <input
                                        type="number"
                                        {...register(
                                          'progressiveRollout.template.increments'
                                        )}
                                        min="0"
                                        max="100"
                                        onKeyDown={(evt: any) => {
                                          if (evt.key === '.') {
                                            evt.preventDefault();
                                          }
                                        }}
                                        className={classNames(
                                          'w-full',
                                          'input-text'
                                        )}
                                        placeholder={''}
                                        required
                                        disabled={
                                          !editable || isSeeDetailsSelected
                                        }
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
                                      {errors.progressiveRollout?.template
                                        ?.increments?.message && (
                                        <span role="alert">
                                          {
                                            errors.progressiveRollout?.template
                                              ?.increments?.message
                                          }
                                        </span>
                                      )}
                                    </p>
                                  </div>
                                  <div className="flex-1">
                                    <span className="input-label">
                                      Frequency
                                    </span>
                                    <Controller
                                      name="progressiveRollout.template.interval"
                                      control={control}
                                      render={({ field }) => (
                                        <Select
                                          onChange={(o: Option) =>
                                            field.onChange(o.value)
                                          }
                                          options={intervalOptions}
                                          disabled={
                                            !editable || isSeeDetailsSelected
                                          }
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
                                    <div className="flex-1 input-label">
                                      Weight
                                    </div>
                                    <div className="flex-1 input-label">
                                      Execute at
                                    </div>
                                  </div>
                                  <div
                                    className={classNames(
                                      'space-y-2 mt-2',
                                      templateSchedulesList.length > 5 &&
                                        'max-h-[232px] overflow-y-scroll'
                                    )}
                                  >
                                    {templateSchedulesList?.map((_, index) => (
                                      <div
                                        key={index}
                                        className="flex space-x-4"
                                      >
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
                                            className={classNames(
                                              'w-full',
                                              'input-text'
                                            )}
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
                                        onChange={(o: Option) =>
                                          field.onChange(o.value)
                                        }
                                        options={variationOptions}
                                        disabled={
                                          !editable || isSeeDetailsSelected
                                        }
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
                                  // disabled={(function () {
                                  //   const lastSchedule: any =
                                  //     manualSchedulesList[
                                  //       manualSchedulesList.length - 1
                                  //     ];
                                  //   return Number(lastSchedule.weight) === 100;
                                  // })()}
                                >
                                  <PlusIcon width={16} />
                                  <span className="text-sm font-medium">
                                    Add Operation
                                  </span>
                                </button>
                                <div className="space-y-2">
                                  {manualSchedulesList.map(
                                    (schedule, index) => (
                                      <div
                                        key={index}
                                        className="flex space-x-4"
                                      >
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
                                                errors.progressiveRollout
                                                  ?.manual?.schedulesList[index]
                                                  ?.weight?.message
                                                  ? 'input-text-error'
                                                  : 'input-text'
                                              )}
                                              placeholder={''}
                                              required
                                              disabled={
                                                !editable ||
                                                isSeeDetailsSelected
                                              }
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
                                            {errors.progressiveRollout?.manual
                                              ?.schedulesList[index]?.weight
                                              ?.message && (
                                              <span role="alert">
                                                {
                                                  errors.progressiveRollout
                                                    ?.manual?.schedulesList[
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
                                            {errors.progressiveRollout?.manual
                                              ?.schedulesList[index]?.executeAt
                                              ?.time?.message && (
                                              <span role="alert">
                                                {
                                                  errors.progressiveRollout
                                                    ?.manual?.schedulesList[
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
                                              onClick={() =>
                                                handleRemoveTrigger(index)
                                              }
                                              className="minus-circle-icon"
                                              // disabled={
                                              //   index <
                                              //   manualSchedulesList.length - 1
                                              // }
                                            >
                                              <MinusCircleIcon aria-hidden="true" />
                                            </button>
                                          </div>
                                        )}
                                      </div>
                                    )
                                  )}
                                  {/* <p className="input-error">
                                    {errors.progressiveRollout?.manual
                                      ?.schedulesList?.length > 0 && (
                                      <span role="alert">
                                        {
                                          errors.progressiveRollout.manual.schedulesList.find(
                                            (s) => s?.weight
                                          )?.weight?.message
                                        }
                                      </span>
                                    )}
                                  </p> */}
                                </div>
                              </Fragment>
                            )}
                          </div>
                        )}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
          <div className="flex-shrink-0 px-4 py-4 flex justify-end border-t">
            <div className="mr-3">
              <button
                type="button"
                className="btn-cancel"
                disabled={false}
                onClick={onCancel}
              >
                {f(messages.button.cancel)}
              </button>
            </div>
            <button
              type="button"
              className="btn-submit-gradient"
              disabled={!isValid || isSubmitting || isSeeDetailsSelected}
              onClick={handleSubmit(handleOnSubmit)}
            >
              {f(messages.button.submit)}
            </button>
          </div>
        </form>
      </div>
    );
  }
);

export interface AddGoalSelectProps {
  options: Option[];
  disabled?: boolean;
  clearable?: boolean;
  isLoading?: boolean;
  isMulti?: boolean;
  isSearchable?: boolean;
  value?: Option;
  className?: string;
  onChange: ((option: Option) => void) | ((option: Option[]) => void);
  placeholder?: string;
  openAddGoalModal: () => void;
}

const AddGoalSelect: FC<AddGoalSelectProps> = memo(
  ({
    disabled,
    className,
    clearable,
    isLoading,
    isMulti,
    isSearchable,
    onChange,
    options,
    value,
    placeholder,
    openAddGoalModal,
  }) => {
    const textColor = '#3F3F46';
    const textColorDisabled = '#6B7280';
    const backgroundColor = 'white';
    const backgroundColorDisabled = '#F3F4F6';
    const borderColor = '#D1D5DB';
    const fontSize = '0.875rem';
    const lineHeight = '1.25rem';
    const minHeight = '2.5rem';
    const colourStyles = {
      control: (styles, { isDisabled }) => ({
        ...styles,
        backgroundColor: isDisabled ? backgroundColorDisabled : backgroundColor,
        borderColor: borderColor,
        '&:hover': {
          borderColor: borderColor,
        },
        fontSize: fontSize,
        lineHeight: lineHeight,
        minHeight: minHeight,
        '*': {
          boxShadow: 'none !important',
        },
      }),
      option: (styles, { isFocused, isSelected }) => {
        return {
          ...styles,
          backgroundColor: isFocused
            ? backgroundColorDisabled
            : isSelected
            ? backgroundColor
            : null,
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          whiteSpace: 'nowrap',
          color: textColor,
          ':active': {
            backgroundColor: backgroundColor,
          },
        };
      },
      menu: (base) => ({
        ...base,
        fontSize: fontSize,
        lineHeight: lineHeight,
        color: textColor,
      }),
      multiValueLabel: (base, { isDisabled }) => ({
        ...base,
        color: isDisabled ? textColorDisabled : textColor,
      }),
      singleValue: (base, { isDisabled }) => ({
        ...base,
        color: isDisabled ? textColorDisabled : textColor,
      }),
    };
    return (
      <ReactSelect
        options={options}
        className={className}
        classNamePrefix="add-goal-react-select"
        styles={colourStyles}
        components={{
          Option: (props) => (
            <CustomOption {...props} openAddGoalModal={openAddGoalModal} />
          ),
        }}
        isDisabled={isLoading || disabled}
        isClearable={clearable}
        isMulti={isMulti}
        isSearchable={isSearchable}
        isLoading={isLoading}
        placeholder={placeholder ? placeholder : ''}
        value={value}
        onChange={onChange}
        openAddGoalModal={openAddGoalModal}
      />
    );
  }
);

const CustomOption = ({ children, ...props }) => {
  const isLastOption =
    props.options[props.options.length - 1]?.value === props.data.value;

  if (isLastOption) {
    return (
      <div>
        <div
          {...props.innderRef}
          {...props.innerProps}
          className="px-3 py-2 hover:bg-[#F3F4F6]"
        >
          {children}
        </div>
        <div
          onClick={props.openAddGoalModal}
          className="text-center py-[10px] cursor-pointer border-t hover:bg-[#F3F4F6]"
        >
          Add New Goal
        </div>
      </div>
    );
  }

  return <components.Option {...props}>{children}</components.Option>;
};

interface AddGoalModalProps {
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

const AddGoalModal: FC<AddGoalModalProps> = ({ open, setOpen }) => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const methods = useFormContext();

  const { setValue } = methods;

  const {
    register,
    handleSubmit,
    formState: { errors, isValid, isSubmitting },
    reset: resetAdd,
  } = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      id: '',
      name: '',
      description: '',
    },
    mode: 'onChange',
  });

  const onSubmit = (data) => console.log(data);

  const handleCreateGoal = useCallback(
    async (data) => {
      dispatch(
        createGoal({
          environmentNamespace: currentEnvironment.id,
          id: data.id,
          name: data.name,
          description: data.description,
        })
      ).then(() => {
        setOpen(false);
        resetAdd();
        dispatch(
          listGoals({
            environmentNamespace: currentEnvironment.id,
            pageSize: 99999,
            cursor: '',
            searchKeyword: '',
            status: null,
            orderBy: ListGoalsRequest.OrderBy.NAME,
            orderDirection: ListGoalsRequest.OrderDirection.ASC,
          })
        );
        setValue('eventRate.goal', data.id);
      });
    },
    [dispatch]
  );

  return (
    <Transition.Root show={open} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={setOpen}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
        </Transition.Child>
        <form
          onSubmit={handleSubmit(onSubmit)}
          className="fixed inset-0 z-10 overflow-y-auto"
        >
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
              enterTo="opacity-100 translate-y-0 sm:scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 translate-y-0 sm:scale-100"
              leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            >
              <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all w-[542px]">
                <div className="flex items-center justify-between px-4 py-5 border-b">
                  <p className="text-xl font-medium">Create Goal</p>
                  <XIcon
                    width={20}
                    className="text-gray-400 cursor-pointer"
                    onClick={() => setOpen(false)}
                  />
                </div>
                <div className="p-4 space-y-4">
                  <p className="font-bold">General Information</p>
                  <div className="space-y-1">
                    <label
                      htmlFor="name"
                      className="flex space-x-2 items-center"
                    >
                      <span className="input-label">{f({ id: 'name' })}</span>
                      <ExclamationCircleIcon width={18} />
                    </label>
                    <input
                      {...register('name')}
                      type="text"
                      name="name"
                      id="name"
                      className="input-text w-full"
                    />
                    <p className="input-error">
                      {errors.name && (
                        <span role="alert">{errors.name.message}</span>
                      )}
                    </p>
                  </div>
                  <div className="space-y-1">
                    <label htmlFor="id" className="flex space-x-2 items-center">
                      <span className="input-label">{f({ id: 'id' })}</span>
                      <ExclamationCircleIcon width={18} />
                    </label>
                    <input
                      {...register('id')}
                      type="text"
                      name="id"
                      id="id"
                      className="input-text w-full"
                    />
                    <p className="input-error">
                      {errors.id && (
                        <span role="alert">{errors.id.message}</span>
                      )}
                    </p>
                  </div>
                  <div className="space-y-1">
                    <label htmlFor="description" className="block">
                      <span className="input-label">
                        {f(messages.description)}
                      </span>
                      <span className="input-label-optional">
                        {' '}
                        ({f(messages.input.optional)})
                      </span>
                    </label>
                    <textarea
                      {...register('description')}
                      id="description"
                      name="description"
                      rows={5}
                      className="input-text w-full break-all"
                    />
                    <p className="input-error">
                      {errors.description && (
                        <span role="alert">{errors.description.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="p-4 flex justify-end border-t space-x-4">
                  <button
                    type="button"
                    className="btn-cancel"
                    disabled={false}
                    onClick={() => setOpen(false)}
                  >
                    {f(messages.button.cancel)}
                  </button>
                  <button
                    type="button"
                    className="btn-submit"
                    disabled={!isValid || isSubmitting}
                    onClick={handleSubmit(handleCreateGoal)}
                  >
                    New Goal
                  </button>
                </div>
              </div>
            </Transition.Child>
          </div>
        </form>
      </Dialog>
    </Transition.Root>
  );
};

interface OpsEventRateClauseSchema {
  variation: string;
  goal: string;
  minCount: number;
  threadsholdRate: number;
  operator: string;
}

export function createOpsEventRateClause(
  oerc: OpsEventRateClauseSchema
): OpsEventRateClause {
  const clause = new OpsEventRateClause();
  clause.setVariationId(oerc.variation);
  clause.setGoalId(oerc.goal);
  clause.setMinCount(oerc.minCount);
  clause.setThreadsholdRate(oerc.threadsholdRate / 100);
  clause.setOperator(createOpsEventRateOperator(oerc.operator));
  return clause;
}

export function createOpsEventRateOperator(
  value: string
): OpsEventRateClause.OperatorMap[keyof OpsEventRateClause.OperatorMap] {
  if (value === OpsEventRateClause.Operator.GREATER_OR_EQUAL.toString()) {
    return OpsEventRateClause.Operator.GREATER_OR_EQUAL;
  }
  return OpsEventRateClause.Operator.LESS_OR_EQUAL;
}
