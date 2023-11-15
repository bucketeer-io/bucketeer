import { intl } from '@/lang';
import { AppState } from '@/modules';
import { useCurrentEnvironment } from '@/modules/me';
import { selectAll as selectAllProgressiveRollouts } from '@/modules/porgressiveRollout';
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
import { ProgressiveRollout } from '@/proto/autoops/progressive_rollout_pb';
import { Feature } from '@/proto/feature/feature_pb';
import { AppDispatch } from '@/store';
import { classNames } from '@/utils/css';
import { XIcon } from '@heroicons/react/outline';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useFormContext, useWatch } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import {
  createAutoOpsRule,
  updateAutoOpsRule,
  UpdateAutoOpsRuleParams,
} from '../../modules/autoOpsRules';
import { selectById as selectFeatureById } from '../../modules/features';
import { createProgressiveRollout } from '../../modules/porgressiveRollout';
import { AddProgressiveRolloutOperation } from '../AddProgressiveRolloutOperation';
import {
  AddUpdateEventRateOperation,
  createOpsEventRateClause,
} from '../AddUpdateEventRateOperation';
import { AddUpdateScheduleOperation } from '../AddUpdateScheduleOperation';
import {
  ClauseType,
  ProgressiveRolloutClauseType,
} from '../FeatureAutoOpsRulesForm';

export interface ProgressiveRolloutTypeTab {
  label: 'Template' | 'Manual';
  value: string;
  selected: boolean;
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
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();

    const { formatMessage: f } = useIntl();

    const [feature] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
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

    const [radioList, setRadioList] = useState([]);

    const methods = useFormContext<any>();
    const {
      handleSubmit,
      control,
      formState: { isValid, isSubmitting, errors },
      register,
      setValue,
    } = methods;

    const opsType = useWatch({
      control,
      name: 'opsType',
    });

    const clauseType = useWatch({
      control,
      name: 'clauseType',
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

    const [progressiveRolloutTypeList, setProgressiveRolloutTypeList] =
      useState<ProgressiveRolloutTypeTab[]>([
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

    const isSeeDetailsSelected = autoOpsRule && !isActiveTabSelected;

    const setEnableList = () => {
      if (autoOpsRule) {
        setRadioList([
          {
            label: f(messages.autoOps.schedule),
            value: ClauseType.DATETIME,
          },
        ]);
      } else {
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
      }
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
            Math.round(opsEventRateClause.threadsholdRate * 100)
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
          } else if (data.clauseType === ClauseType.PROGRESSIVE_ROLLOUT) {
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

    const variationOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: v.value,
      };
    });

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
                          <AddUpdateScheduleOperation
                            isSeeDetailsSelected={isSeeDetailsSelected}
                          />
                        )}
                      {radio.value === ClauseType.EVENT_RATE &&
                        clauseType === ClauseType.EVENT_RATE && (
                          <AddUpdateEventRateOperation
                            isSeeDetailsSelected={isSeeDetailsSelected}
                            variationOptions={variationOptions}
                          />
                        )}
                      {radio.value === ClauseType.PROGRESSIVE_ROLLOUT &&
                        clauseType === ClauseType.PROGRESSIVE_ROLLOUT && (
                          <AddProgressiveRolloutOperation
                            featureId={featureId}
                            variationOptions={variationOptions}
                            isSeeDetailsSelected={isSeeDetailsSelected}
                            progressiveRolloutTypeList={
                              progressiveRolloutTypeList
                            }
                            setProgressiveRolloutTypeList={
                              setProgressiveRolloutTypeList
                            }
                          />
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
