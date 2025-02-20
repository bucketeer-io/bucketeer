import { AppState } from '../../modules';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { AutoOpsRule, OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import { listGoals, selectAll as selectAllGoals } from '../../modules/goals';
import { AppDispatch } from '../../store';
import { XIcon } from '@heroicons/react/outline';
import React, { FC, memo, useCallback, useEffect } from 'react';
import { Controller, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import { OperationForm } from '../../pages/feature/formSchema';
import { Option, Select } from '../Select';
import { Goal } from '../../proto/experiment/goal_pb';
import { ListGoalsRequest } from '../../proto/experiment/service_pb';
import { classNames } from '../../utils/css';
import { OpsEventRateClause, ActionType } from '../../proto/autoops/clause_pb';
import {
  CreateAutoOpsRuleCommand,
  ChangeOpsEventRateClauseCommand
} from '../../proto/autoops/command_pb';
import {
  createAutoOpsRule,
  updateAutoOpsRule,
  UpdateAutoOpsRuleParams
} from '../../modules/autoOpsRules';
import { operatorOptions } from '../../pages/feature/autoops';
import { AddGoalSelect } from '../AddGoalSelect';

export interface EventRateAddUpdateFormProps {
  onCancel: () => void;
  featureId: string;
  autoOpsRule?: AutoOpsRule.AsObject;
  isActiveTabSelected: boolean;
  variationOptions: Option[];
  onSubmit: () => void;
}

export const EventRateAddUpdateForm: FC<EventRateAddUpdateFormProps> = memo(
  ({
    onCancel,
    featureId,
    autoOpsRule,
    isActiveTabSelected,
    variationOptions,
    onSubmit
  }) => {
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();

    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();

    const methods = useFormContext<OperationForm>();
    const {
      handleSubmit,
      control,
      formState: { isValid, isSubmitting, errors },
      register,
      setValue
    } = methods;

    const isSeeDetailsSelected = autoOpsRule && !isActiveTabSelected;

    const title = () => {
      if (isSeeDetailsSelected) {
        return f(messages.autoOps.operationDetails);
      } else {
        return autoOpsRule
          ? f(messages.autoOps.updateOperation)
          : f(messages.autoOps.createOperation);
      }
    };

    const goals = useSelector<AppState, Goal.AsObject[]>(
      (state) => selectAllGoals(state.goals),
      shallowEqual
    );

    useEffect(() => {
      dispatch(
        listGoals({
          environmentId: currentEnvironment.id,
          pageSize: 0,
          cursor: '',
          searchKeyword: '',
          status: null,
          orderBy: ListGoalsRequest.OrderBy.NAME,
          orderDirection: ListGoalsRequest.OrderDirection.ASC,
          connectionType: Goal.ConnectionType.OPERATION
        })
      );
    }, [dispatch, featureId, currentEnvironment]);

    useEffect(() => {
      if (autoOpsRule) {
        const opsEventRateClause = OpsEventRateClause.deserializeBinary(
          autoOpsRule.clausesList[0].clause.value as Uint8Array
        ).toObject();

        setValue('eventRate.variation', opsEventRateClause.variationId);
        setValue('eventRate.goal', opsEventRateClause.goalId);
        setValue('eventRate.operator', opsEventRateClause.operator.toString());
        setValue(
          'eventRate.threadsholdRate',
          Math.round(opsEventRateClause.threadsholdRate * 100)
        );
        setValue('eventRate.minCount', opsEventRateClause.minCount);
      }
    }, [autoOpsRule]);

    const handleOnSubmit = useCallback(
      (data) => {
        if (autoOpsRule) {
          const changeOpsEventRateClauseCommands: ChangeOpsEventRateClauseCommand[] =
            [];

          const command = new ChangeOpsEventRateClauseCommand();
          command.setId(autoOpsRule.clausesList[0].id);
          command.setOpsEventRateClause(
            createOpsEventRateClause(data.eventRate)
          );
          changeOpsEventRateClauseCommands.push(command);

          const param: UpdateAutoOpsRuleParams = {
            environmentId: currentEnvironment.id,
            id: autoOpsRule.id,
            changeOpsEventRateClauseCommands
          };

          dispatch(updateAutoOpsRule(param)).then(() => onSubmit());
        } else {
          const command = new CreateAutoOpsRuleCommand();
          command.setFeatureId(featureId);
          command.setOpsType(OpsType.EVENT_RATE);
          command.setOpsEventRateClausesList([
            createOpsEventRateClause(data.eventRate)
          ]);

          dispatch(
            createAutoOpsRule({
              environmentId: currentEnvironment.id,
              command: command
            })
          ).then(() => onSubmit());
        }
      },
      [autoOpsRule]
    );

    const goalOptions = goals.map((goal) => {
      return {
        value: goal.id,
        label: `${goal.id} (${goal.name})`
      };
    });

    return (
      <div className="w-[530px] h-full overflow-hidden">
        <form className="flex flex-col h-full overflow-hidden">
          <div className="h-full flex flex-col overflow-hidden">
            <div className="flex items-center justify-between px-4 py-5 border-b">
              <p className="text-xl font-medium">{title()}</p>
              <XIcon
                width={20}
                className="text-gray-400 cursor-pointer"
                onClick={onCancel}
              />
            </div>
            <div className="px-4 h-full flex flex-col overflow-hidden">
              <div className="flex py-5">
                <p className="font-bold">{f(messages.autoOps.eventRate)}</p>
              </div>
              <div className="h-full flex flex-col overflow-hidden space-y-4 px-1">
                <div className="flex divide-x divide-primary space-x-3">
                  <div className="px-4 py-1 text-pink-500 bg-pink-50 inline-block self-center">
                    If
                  </div>
                  <div className="space-y-3 pl-3">
                    <div>
                      <span className="input-label">
                        {f(messages.feature.variation)}
                      </span>
                      <Controller
                        name="eventRate.variation"
                        control={control}
                        render={({ field }) => (
                          <Select
                            isSearchable={false}
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
                            name="eventRate.goal"
                            onChange={(o: Option) => field.onChange(o.value)}
                            options={goalOptions}
                            disabled={!editable || isSeeDetailsSelected}
                            value={goalOptions.find(
                              (o) => o.value === field.value
                            )}
                            connectionType={Goal.ConnectionType.OPERATION}
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
                              onChange={(o: Option) => field.onChange(o.value)}
                              options={operatorOptions}
                              disabled={!editable || isSeeDetailsSelected}
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
                          {f(messages.autoOps.opsEventRateClause.minCount)}
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
                </div>
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
  clause.setActionType(ActionType.DISABLE);
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
