import { useCurrentEnvironment } from '../../modules/me';
import { AutoOpsRule, OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import { AppDispatch } from '../../store';
import { PlusIcon, TrashIcon, XIcon } from '@heroicons/react/outline';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import {
  Controller,
  useFieldArray,
  useFormContext,
  useWatch
} from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useDispatch } from 'react-redux';
import { messages } from '../../lang/messages';
import { OperationForm } from '../../pages/feature/formSchema';
import { Select } from '../Select';
import { DatetimePicker } from '../DatetimePicker';
import {
  AddDatetimeClauseCommand,
  ChangeDatetimeClauseCommand,
  CreateAutoOpsRuleCommand,
  DeleteClauseCommand
} from '../../proto/autoops/command_pb';
import { ActionType, DatetimeClause } from '../../proto/autoops/clause_pb';
import {
  createAutoOpsRule,
  updateAutoOpsRule,
  UpdateAutoOpsRuleParams
} from '../../modules/autoOpsRules';
import { v4 as uuid } from 'uuid';
import dayjs from 'dayjs';
import {
  isTimestampArraySorted,
  hasDuplicateTimestamps
} from '../../utils/isArraySorted';
import { getDatetimeClause } from '../../utils/getDatetimeClause';

export const actionTypesOptions = [
  { value: ActionType.ENABLE.toString(), label: 'On' },
  { value: ActionType.DISABLE.toString(), label: 'Off' }
];

export interface ScheduleAddUpdateFormProps {
  featureId: string;
  currentFlagStatus: boolean;
  onSubmit: () => void;
  onCancel: () => void;
  autoOpsRule?: AutoOpsRule.AsObject;
  isActiveTabSelected: boolean;
}

export const ScheduleAddUpdateForm: FC<ScheduleAddUpdateFormProps> = memo(
  ({
    featureId,
    currentFlagStatus,
    onSubmit,
    onCancel,
    autoOpsRule,
    isActiveTabSelected
  }) => {
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();

    const { formatMessage: f } = useIntl();

    const methods = useFormContext<OperationForm>();
    const {
      handleSubmit,
      control,
      formState: { isValid, isSubmitting },
      setValue
    } = methods;

    const isSeeDetailsSelected = autoOpsRule && !isActiveTabSelected;

    const {
      fields: datetimeClausesList,
      remove: removeDatetimeClause,
      append: appendDatetimeClause
    } = useFieldArray({
      control,
      name: 'datetimeClausesList'
    });

    const watchDatetimeClausesList = useWatch({
      control,
      name: 'datetimeClausesList'
    });

    const [disabledDatetimeClasuesListIds, setDisabledDatetimeClasuesListIds] =
      useState([]);

    useEffect(() => {
      if (autoOpsRule) {
        setDisabledDatetimeClasuesListIds(
          datetimeClausesList
            .filter((clause) => isSameOrBeforeNow(clause.time))
            .map((clause) => clause.id)
        );
      }
    }, [datetimeClausesList, autoOpsRule]);

    useEffect(() => {
      if (autoOpsRule) {
        const datetimeClausesList = autoOpsRule.clausesList.map((clause) => {
          const datetime = getDatetimeClause(clause.clause.value);
          return {
            id: clause.id,
            actionType: clause.actionType,
            time: new Date(datetime.time * 1000)
          };
        });
        setValue('datetimeClausesList', datetimeClausesList);
      }
    }, [autoOpsRule]);

    const handleOnSubmit = useCallback(
      (data) => {
        if (autoOpsRule) {
          const orgIds = autoOpsRule.clausesList.map((clause) => clause.id);
          const valIds = data.datetimeClausesList.map((clause) => clause.id);

          // delete clauses
          const toDelete = orgIds.filter((id) => !valIds.includes(id));
          const deleteClauseCommands: DeleteClauseCommand[] = [];
          toDelete.forEach((id) => {
            const command = new DeleteClauseCommand();
            command.setId(id);
            deleteClauseCommands.push(command);
          });

          const toAdd = data.datetimeClausesList.filter(
            (clause) => !orgIds.includes(clause.id)
          );

          // add clauses
          const addDatetimeClauseCommands: AddDatetimeClauseCommand[] = [];
          toAdd.forEach((datetimeClause) => {
            const command = new AddDatetimeClauseCommand();
            const clause = new DatetimeClause();
            clause.setTime(Math.round(datetimeClause.time.getTime() / 1000));
            clause.setActionType(datetimeClause.actionType);
            command.setDatetimeClause(clause);
            addDatetimeClauseCommands.push(command);
          });

          // change clauses
          const changeDatetimeClauseCommands: ChangeDatetimeClauseCommand[] =
            [];
          orgIds.forEach((orgId) => {
            const orgClause = autoOpsRule.clausesList.find(
              (clause) => clause.id === orgId
            );
            const valClause = data.datetimeClausesList.find(
              (clause) => clause.id === orgId
            );

            if (orgClause && valClause) {
              const orgDatetime = getDatetimeClause(
                orgClause.clause.value
              ).time;
              const valDatetime = Math.round(valClause.time.getTime() / 1000);

              if (
                orgDatetime !== valDatetime ||
                orgClause.actionType !== valClause.actionType
              ) {
                const command = new ChangeDatetimeClauseCommand();
                const clause = new DatetimeClause();
                clause.setTime(valDatetime);
                clause.setActionType(valClause.actionType);

                command.setId(orgId);
                command.setDatetimeClause(clause);
                changeDatetimeClauseCommands.push(command);
              }
            }
          });

          if (
            addDatetimeClauseCommands.length === 0 &&
            deleteClauseCommands.length === 0 &&
            changeDatetimeClauseCommands.length === 0
          ) {
            return onSubmit();
          }

          const param: UpdateAutoOpsRuleParams = {
            environmentId: currentEnvironment.id,
            id: autoOpsRule.id,
            addDatetimeClauseCommands,
            deleteClauseCommands,
            changeDatetimeClauseCommands
          };

          dispatch(updateAutoOpsRule(param)).then(() => onSubmit());
        } else {
          const command = new CreateAutoOpsRuleCommand();
          command.setFeatureId(featureId);
          command.setOpsType(OpsType.SCHEDULE);

          const datetimeClausesList = [];
          data.datetimeClausesList.forEach((datetimeClause) => {
            const clause = new DatetimeClause();
            clause.setTime(Math.round(datetimeClause.time.getTime() / 1000));
            clause.setActionType(datetimeClause.actionType);
            datetimeClausesList.push(clause);
          });
          command.setDatetimeClausesList(datetimeClausesList);

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

    const handleAddDatetimeClause = useCallback(() => {
      const lastDatetimeClause =
        watchDatetimeClausesList[watchDatetimeClausesList.length - 1];

      const time = dayjs(lastDatetimeClause.time).add(1, 'hour').toDate();

      appendDatetimeClause({
        id: uuid(),
        actionType: ActionType.ENABLE,
        time
      });
    }, [watchDatetimeClausesList]);

    const handleDelete = useCallback((idx) => {
      removeDatetimeClause(idx);
    }, []);

    const isSameOrBeforeNow = (time: Date) => {
      return dayjs(time).isSameOrBefore(dayjs());
    };

    const title = () => {
      if (isSeeDetailsSelected) {
        return f(messages.autoOps.operationDetails);
      } else {
        return autoOpsRule
          ? f(messages.autoOps.updateOperation)
          : f(messages.autoOps.createOperation);
      }
    };

    const isDisabled = (id: string) => {
      if (autoOpsRule) {
        return disabledDatetimeClasuesListIds.includes(id);
      }
      return false;
    };

    const isDatesSorted = isTimestampArraySorted(
      watchDatetimeClausesList.map((d) => d.time.getTime())
    );

    const hasDuplicates = hasDuplicateTimestamps(
      watchDatetimeClausesList.map((d) => d.time.getTime())
    );

    const _checkInvalidDatetime = () => {
      let list = watchDatetimeClausesList;
      if (autoOpsRule) {
        // if there are disabled clauses, we need to check only enabled clauses
        list = watchDatetimeClausesList.slice(
          disabledDatetimeClasuesListIds.length,
          watchDatetimeClausesList.length
        );
      }
      // no datetime clauses
      if (list.length === 0) {
        return true;
      }

      // check if there is a datetime clause that is before now
      return !!list.find((clause) => {
        return isSameOrBeforeNow(clause.time);
      });
    };

    return (
      <div className="w-[500px] h-full overflow-hidden">
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
              {!isSeeDetailsSelected && (
                <div className="flex py-5 border-b space-x-6">
                  <p className="font-bold">
                    {f(messages.autoOps.currentFlagState)}
                  </p>
                  <div className="px-2 py-[2px] rounded-md border border-gray-300 text-sm">
                    {currentFlagStatus ? 'On' : 'Off'}
                  </div>
                </div>
              )}
              <div className="py-6 h-full flex flex-col overflow-hidden space-y-4 px-1">
                <p className="font-bold">{f(messages.autoOps.schedule)}</p>
                <div className="h-full overflow-y-auto space-y-2">
                  {datetimeClausesList.map((datetimeClause, idx) => (
                    <div
                      key={datetimeClause.id}
                      className="flex space-x-4 pl-[2px]"
                    >
                      <div className="w-32 space-y-1">
                        <span className="input-label">
                          {f(messages.autoOps.state)}
                        </span>
                        <Controller
                          name={`datetimeClausesList.${idx}.actionType`}
                          control={control}
                          render={({ field }) => (
                            <Select
                              {...field}
                              disabled={
                                isSeeDetailsSelected ||
                                isDisabled(datetimeClause.id)
                              }
                              isSearchable={false}
                              value={actionTypesOptions.find(
                                (o) => o.value === field.value.toString()
                              )}
                              options={actionTypesOptions}
                              onChange={(option) =>
                                field.onChange(option.value)
                              }
                            />
                          )}
                        />
                      </div>
                      <div className="w-full space-y-1">
                        <span className="input-label">
                          {f(messages.autoOps.startDate)}
                        </span>
                        <DatetimePicker
                          name={`datetimeClausesList.${idx}.time`}
                          dateFormat="yyyy/MM/dd HH:mm"
                          disabled={
                            isSeeDetailsSelected ||
                            isDisabled(datetimeClause.id)
                          }
                        />
                        {!isDisabled(datetimeClause.id) &&
                          isSameOrBeforeNow(
                            watchDatetimeClausesList[idx]?.time
                          ) && (
                            <p className="input-error">
                              <span role="alert">
                                {f(
                                  messages.input.error.notLaterThanCurrentTime
                                )}
                              </span>
                            </p>
                          )}
                      </div>
                      {!isSeeDetailsSelected && (
                        <div>
                          <button
                            className="py-[11px] mt-6 text-gray-400 hover:text-gray-500 disabled:opacity-60 disabled:hover:text-gray-400 disabled:cursor-not-allowed"
                            type="button"
                            onClick={() => handleDelete(idx)}
                            disabled={
                              datetimeClausesList.length === 1 ||
                              isDisabled(datetimeClause.id)
                            }
                          >
                            <TrashIcon width={18} />
                          </button>
                        </div>
                      )}
                    </div>
                  ))}
                  {watchDatetimeClausesList.length <= 10 && (
                    <ErrorMessage
                      isDatesSorted={isDatesSorted}
                      hasDuplicates={hasDuplicates}
                    />
                  )}
                  {!isSeeDetailsSelected && (
                    <button
                      className="flex whitespace-nowrap space-x-2 text-primary max-w-min py-2 items-center"
                      type="button"
                      onClick={handleAddDatetimeClause}
                    >
                      <PlusIcon width={18} />
                      <span>{f(messages.button.addSchedule)}</span>
                    </button>
                  )}
                </div>
                {watchDatetimeClausesList.length > 10 && (
                  <ErrorMessage
                    isDatesSorted={isDatesSorted}
                    hasDuplicates={hasDuplicates}
                  />
                )}
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
              disabled={
                !isValid ||
                isSubmitting ||
                isSeeDetailsSelected ||
                _checkInvalidDatetime() ||
                !isDatesSorted ||
                hasDuplicates
              }
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

interface ErrorMessageProps {
  isDatesSorted: boolean;
  hasDuplicates: boolean;
}

const ErrorMessage: FC<ErrorMessageProps> = memo(
  ({ isDatesSorted, hasDuplicates }) => {
    const { formatMessage: f } = useIntl();

    if (isDatesSorted && !hasDuplicates) {
      return null;
    }

    return (
      <div className="flex space-x-2">
        <div className="w-32" />
        <div className="w-full">
          {hasDuplicates ? (
            <p className="input-error">
              <span role="alert">{f(messages.autoOps.duplicateDates)}</span>
            </p>
          ) : !isDatesSorted ? (
            <p className="input-error">
              <span role="alert">
                {f(messages.autoOps.dateIncreasingOrder)}
              </span>
            </p>
          ) : null}
        </div>
      </div>
    );
  }
);
