import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useDispatch, useSelector, shallowEqual } from 'react-redux';
import { useParams, useHistory } from 'react-router-dom';
import { v4 as uuid } from 'uuid';
import * as yup from 'yup';

import { DetailSkeleton } from '../../components/DetailSkeleton';
import {
  ClauseType,
  createInitialDatetimeClause,
  createInitialOpsEventRateClause,
  createInitialWebhookClause,
  FeatureAutoOpsRulesForm,
} from '../../components/FeatureAutoOpsRulesForm';
import {
  PAGE_PATH_FEATURES,
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_ROOT,
} from '../../constants/routing';
import { AppState } from '../../modules';
import {
  createAutoOpsRule,
  deleteAutoOpsRule,
  listAutoOpsRules,
  selectAll as selectAllAutoOpsRules,
  updateAutoOpsRule,
  UpdateAutoOpsRuleParams,
} from '../../modules/autoOpsRules';
import { selectById as selectFeatureById } from '../../modules/features';
import { listGoals } from '../../modules/goals';
import { useCurrentEnvironment } from '../../modules/me';
import { AutoOpsRule, OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import {
  DatetimeClause,
  OpsEventRateClause,
  WebhookClause,
} from '../../proto/autoops/clause_pb';
import {
  AddDatetimeClauseCommand,
  AddOpsEventRateClauseCommand,
  AddWebhookClauseCommand,
  ChangeAutoOpsRuleOpsTypeCommand,
  ChangeDatetimeClauseCommand,
  ChangeOpsEventRateClauseCommand,
  ChangeWebhookClauseCommand,
  CreateAutoOpsRuleCommand,
  DeleteAutoOpsRuleCommand,
  DeleteClauseCommand,
} from '../../proto/autoops/command_pb';
import { ListGoalsRequest } from '../../proto/experiment/service_pb';
import { AddUserToVariationCommand } from '../../proto/feature/command_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { AppDispatch } from '../../store';

import { autoOpsRulesFormSchema } from './formSchema';

interface FeatureAutoOpsPageProps {
  featureId: string;
}

export const FeatureAutoOpsPage: FC<FeatureAutoOpsPageProps> = memo(
  ({ featureId }) => {
    const dispatch = useDispatch<AppDispatch>();
    const history = useHistory();
    const currentEnvironment = useCurrentEnvironment();
    const isFeatureLoading = useSelector<AppState, boolean>(
      (state) => state.features.loading,
      shallowEqual
    );
    const isAutoOpsRuleLoading = useSelector<AppState, boolean>(
      (state) => state.autoOpsRules.loading,
      shallowEqual
    );
    const isLoading = isFeatureLoading || isAutoOpsRuleLoading;
    const autoOpsRules = useSelector<AppState, AutoOpsRule.AsObject[]>(
      (state) =>
        selectAllAutoOpsRules(state.autoOpsRules).filter(
          (rule) => rule.featureId === featureId
        ),
      shallowEqual
    );
    const [feature, getFeatureError] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ]);
    const defaultValues = {
      autoOpsRules: autoOpsRules.map((rule) => {
        return {
          id: rule.id,
          featureId: rule.featureId,
          triggeredAt: rule.triggeredAt,
          opsType: rule.opsType.toString(),
          clauses: rule.clausesList.map((clause) => {
            const typeUrl = clause.clause.typeUrl;
            const type = typeUrl.substring(typeUrl.lastIndexOf('/') + 1);
            if (type === ClauseType.EVENT_RATE) {
              const opsEventRateClause = OpsEventRateClause.deserializeBinary(
                clause.clause.value as Uint8Array
              ).toObject();
              return {
                id: clause.id,
                clauseType: ClauseType.EVENT_RATE.toString(),
                opsEventRateClause: {
                  variation: opsEventRateClause.variationId,
                  goal: opsEventRateClause.goalId,
                  minCount: opsEventRateClause.minCount,
                  threadsholdRate: opsEventRateClause.threadsholdRate * 100,
                  operator: opsEventRateClause.operator.toString(),
                },
                webhookClause: createInitialWebhookClause(),
                datetimeClause: createInitialDatetimeClause(),
              };
            }
            if (type === ClauseType.WEBHOOK) {
              const webhookClause = WebhookClause.deserializeBinary(
                clause.clause.value as Uint8Array
              ).toObject();
              return {
                id: clause.id,
                clauseType: ClauseType.WEBHOOK.toString(),
                webhookClause: {
                  webhookId: webhookClause.webhookId,
                  conditionsList: webhookClause.conditionsList.map((cond) => ({
                    ...cond,
                    id: uuid(),
                    operator: cond.operator.toString(),
                  })),
                },
                datetimeClause: createInitialDatetimeClause(),
                opsEventRateClause: createInitialOpsEventRateClause(feature),
              };
            }
            if (type === ClauseType.DATETIME) {
              const datetimeClause = DatetimeClause.deserializeBinary(
                clause.clause.value as Uint8Array
              ).toObject();
              return {
                id: clause.id,
                clauseType: ClauseType.DATETIME.toString(),
                datetimeClause: {
                  time: new Date(datetimeClause.time * 1000),
                },
                webhookClause: createInitialWebhookClause(),
                opsEventRateClause: createInitialOpsEventRateClause(feature),
              };
            }
          }),
        };
      }),
    };

    const methods = useForm({
      resolver: yupResolver(autoOpsRulesFormSchema),
      defaultValues: defaultValues,
      mode: 'onChange',
    });
    const { handleSubmit, reset, trigger } = methods;

    const handleUpdate = useCallback(
      async (data) => {
        const createAutoOpsRuleCommands = createCreateAutoOpsRuleCommands(
          defaultValues.autoOpsRules,
          data.autoOpsRules
        );
        const deleteAutoOpsRuleIds = createDeleteAutoOpsRuleIds(
          defaultValues.autoOpsRules,
          data.autoOpsRules
        );

        const updateAutoOpsRuleParams = createUpdateAutoOpsRuleParams(
          currentEnvironment.namespace,
          defaultValues.autoOpsRules,
          data.autoOpsRules
        );

        const promises = [];

        createAutoOpsRuleCommands.forEach((command) => {
          promises.push(
            new Promise((resolve) => {
              dispatch(
                createAutoOpsRule({
                  environmentNamespace: currentEnvironment.namespace,
                  command: command,
                })
              ).then((res) => {
                resolve(res);
              });
            })
          );
        });

        deleteAutoOpsRuleIds.forEach((id) => {
          promises.push(
            new Promise((resolve) => {
              dispatch(
                deleteAutoOpsRule({
                  environmentNamespace: currentEnvironment.namespace,
                  id: id,
                })
              ).then((res) => {
                resolve(res);
              });
            })
          );
        });

        updateAutoOpsRuleParams.forEach((param) => {
          promises.push(
            new Promise((resolve) => {
              dispatch(updateAutoOpsRule(param)).then((res) => {
                resolve(res);
              });
            })
          );
        });

        Promise.all(promises).then((res) => {
          dispatch(
            listAutoOpsRules({
              featureId: featureId,
              environmentNamespace: currentEnvironment.namespace,
            })
          );
        });
      },
      [dispatch, defaultValues]
    );

    useEffect(() => {
      dispatch(
        listAutoOpsRules({
          featureId: featureId,
          environmentNamespace: currentEnvironment.namespace,
        })
      );
      dispatch(
        listGoals({
          environmentNamespace: currentEnvironment.namespace,
          pageSize: 99999,
          cursor: '',
          searchKeyword: '',
          status: null,
          orderBy: ListGoalsRequest.OrderBy.NAME,
          orderDirection: ListGoalsRequest.OrderDirection.ASC,
        })
      );
    }, [dispatch, featureId, currentEnvironment]);

    useEffect(() => {
      reset(defaultValues);
    }, [autoOpsRules]);

    if (isLoading) {
      return (
        <div className="p-9 bg-gray-100">
          <DetailSkeleton />
        </div>
      );
    }
    return (
      <FormProvider {...methods}>
        <FeatureAutoOpsRulesForm
          featureId={featureId}
          onSubmit={handleSubmit(handleUpdate)}
        />
      </FormProvider>
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

interface WebhookConditionSchema {
  id: string;
  filter: string;
  operator: string;
  value: string;
}

interface WebhookClauseSchema {
  webhookId: string;
  conditionsList: WebhookConditionSchema[];
}

interface DatetimeClauseSchema {
  time: Date;
}

interface ClauseSchema {
  id: string;
  clauseType: string;
  opsEventRateClause?: OpsEventRateClauseSchema;
  datetimeClause?: DatetimeClauseSchema;
  webhookClause?: WebhookClauseSchema;
}

interface AutoOpsRuleSchema {
  id: string;
  featureId: string;
  triggeredAt: number;
  opsType: string;
  clauses: ClauseSchema[];
}

export function createCreateAutoOpsRuleCommands(
  org: AutoOpsRuleSchema[],
  val: AutoOpsRuleSchema[]
): CreateAutoOpsRuleCommand[] {
  const commands: Array<CreateAutoOpsRuleCommand> = [];
  const orgIds = org.map((r) => r.id);
  val
    .filter((r) => !orgIds.includes(r.id))
    .forEach((r) => {
      const command = new CreateAutoOpsRuleCommand();
      command.setFeatureId(r.featureId);
      r.opsType === OpsType.ENABLE_FEATURE.toString() &&
        command.setOpsType(OpsType.ENABLE_FEATURE);
      r.opsType === OpsType.DISABLE_FEATURE.toString() &&
        command.setOpsType(OpsType.DISABLE_FEATURE);
      command.setOpsEventRateClausesList(createOpsEventRateClauses(r.clauses));
      command.setDatetimeClausesList(createDatetimeClauses(r.clauses));
      command.setWebhookClausesList(createWebhookClauses(r.clauses));
      commands.push(command);
    });
  return commands;
}

export function createDeleteAutoOpsRuleIds(
  org: AutoOpsRuleSchema[],
  val: AutoOpsRuleSchema[]
): string[] {
  const ids: string[] = [];
  const valIds = val.map((r) => r.id);
  org
    .filter((r) => !valIds.includes(r.id))
    .forEach((r) => {
      ids.push(r.id);
    });
  return ids;
}

export function createOpsEventRateClauses(
  val: ClauseSchema[]
): OpsEventRateClause[] {
  const clauses: Array<OpsEventRateClause> = [];
  val.forEach((c) => {
    if (c.clauseType === ClauseType.EVENT_RATE.toString()) {
      const clause = createOpsEventRateClause(c.opsEventRateClause);
      clauses.push(clause);
    }
  });
  return clauses;
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

export function createDatetimeClauses(val: ClauseSchema[]): DatetimeClause[] {
  const clauses: Array<DatetimeClause> = [];
  val.forEach((c) => {
    if (c.clauseType === ClauseType.DATETIME.toString()) {
      const clause = createDatetimeClause(c.datetimeClause);
      clauses.push(clause);
    }
  });
  return clauses;
}

export function createDatetimeClause(
  dtc: DatetimeClauseSchema
): DatetimeClause {
  const clause = new DatetimeClause();
  clause.setTime(Math.round(dtc.time.getTime() / 1000));
  return clause;
}

export function createWebhookClauses(val: ClauseSchema[]): WebhookClause[] {
  const clauses: Array<WebhookClause> = [];
  val.forEach((c) => {
    if (c.clauseType === ClauseType.WEBHOOK.toString()) {
      const clause = createWebhookClause(c.webhookClause);
      clauses.push(clause);
    }
  });
  return clauses;
}

export function createWebhookClause(whc: WebhookClauseSchema): WebhookClause {
  const clause = new WebhookClause();
  clause.setWebhookId(whc.webhookId);
  clause.setConditionsList(createConditionsList(whc.conditionsList));
  return clause;
}

export function createConditionsList(
  conditionsList: WebhookConditionSchema[]
): Array<WebhookClause.Condition> {
  const conditions: Array<WebhookClause.Condition> = [];

  conditionsList.forEach((condition) => {
    conditions.push(createCondition(condition));
  });
  return conditions;
}

export function createCondition(
  condition: WebhookConditionSchema
): WebhookClause.Condition {
  const conditionClause = new WebhookClause.Condition();

  conditionClause.setFilter(condition.filter);
  conditionClause.setValue(condition.value);
  conditionClause.setOperator(createWebhookOperator(condition.operator));
  return conditionClause;
}

export function createWebhookOperator(
  value: string
): WebhookClause.Condition.OperatorMap[keyof WebhookClause.Condition.OperatorMap] {
  switch (value) {
    case WebhookClause.Condition.Operator.EQUAL.toString():
      return WebhookClause.Condition.Operator.EQUAL;
    case WebhookClause.Condition.Operator.NOT_EQUAL.toString():
      return WebhookClause.Condition.Operator.NOT_EQUAL;
    case WebhookClause.Condition.Operator.MORE_THAN.toString():
      return WebhookClause.Condition.Operator.MORE_THAN;
    case WebhookClause.Condition.Operator.MORE_THAN_OR_EQUAL.toString():
      return WebhookClause.Condition.Operator.MORE_THAN_OR_EQUAL;
    case WebhookClause.Condition.Operator.LESS_THAN.toString():
      return WebhookClause.Condition.Operator.LESS_THAN;
    default:
      return WebhookClause.Condition.Operator.LESS_THAN_OR_EQUAL;
  }
}

export function createChangeAutoOpsRuleOpsTypeCommand(
  org: AutoOpsRuleSchema[],
  val: AutoOpsRuleSchema[]
): ChangeAutoOpsRuleOpsTypeCommand[] {
  const commands: ChangeAutoOpsRuleOpsTypeCommand[] = [];
  const orgRuleIds = org.filter((r) => r.id).map((r) => r.id);
  const valRuleIds = val.filter((r) => r.id).map((r) => r.id);
  // Intersection of org and val rules.
  const ids = orgRuleIds.filter((id) => valRuleIds.includes(id));
  ids.forEach((rid) => {
    const orgRule = org.find((r) => r.id === rid);
    const valRule = val.find((r) => r.id === rid);
    if (orgRule.opsType !== valRule.opsType) {
      const command = new ChangeAutoOpsRuleOpsTypeCommand();
      valRule.opsType === OpsType.ENABLE_FEATURE.toString() &&
        command.setOpsType(OpsType.ENABLE_FEATURE);
      valRule.opsType === OpsType.DISABLE_FEATURE.toString() &&
        command.setOpsType(OpsType.DISABLE_FEATURE);
      commands.push(command);
    }
  });
  return commands;
}

export function createUpdateAutoOpsRuleParams(
  environmentNamespace: string,
  org: AutoOpsRuleSchema[],
  val: AutoOpsRuleSchema[]
): UpdateAutoOpsRuleParams[] {
  const params: UpdateAutoOpsRuleParams[] = [];
  const orgRuleIds = org.filter((r) => r.id).map((r) => r.id);
  const valRuleIds = val.filter((r) => r.id).map((r) => r.id);
  // Intersection of org and val rules.
  const ruleIds = orgRuleIds.filter((id) => valRuleIds.includes(id));
  ruleIds.forEach((rid) => {
    const orgRule = org.find((r) => r.id === rid);
    const valRule = val.find((r) => r.id === rid);
    const orgClauseIds = orgRule.clauses.filter((c) => c.id).map((c) => c.id);
    const valClauseIds = valRule.clauses.filter((c) => c.id).map((c) => c.id);
    // Intersection of org and val rules.
    const clauseIds = orgClauseIds.filter((id) => valClauseIds.includes(id));

    let changeAutoOpsRuleOpsTypeCommand: ChangeAutoOpsRuleOpsTypeCommand;
    if (orgRule.opsType !== valRule.opsType) {
      const command = new ChangeAutoOpsRuleOpsTypeCommand();
      valRule.opsType === OpsType.ENABLE_FEATURE.toString() &&
        command.setOpsType(OpsType.ENABLE_FEATURE);
      valRule.opsType === OpsType.DISABLE_FEATURE.toString() &&
        command.setOpsType(OpsType.DISABLE_FEATURE);
      changeAutoOpsRuleOpsTypeCommand = command;
    }

    const deleteClauseCommands: DeleteClauseCommand[] = [];
    orgRule.clauses
      .filter((c) => !clauseIds.includes(c.id))
      .forEach((c) => {
        const command = new DeleteClauseCommand();
        command.setId(c.id);
        deleteClauseCommands.push(command);
      });

    const addOpsEventRateClauseCommands: AddOpsEventRateClauseCommand[] = [];
    const addWebhookClauseCommands: AddWebhookClauseCommand[] = [];

    const addDatetimeClauseCommands: AddDatetimeClauseCommand[] = [];
    valRule.clauses
      .filter((c) => !clauseIds.includes(c.id))
      .forEach((c) => {
        if (c.clauseType === ClauseType.EVENT_RATE.toString()) {
          const clause = createOpsEventRateClause(c.opsEventRateClause);
          const command = new AddOpsEventRateClauseCommand();
          command.setOpsEventRateClause(clause);
          addOpsEventRateClauseCommands.push(command);
        }
        if (c.clauseType === ClauseType.WEBHOOK.toString()) {
          const clause = createWebhookClause(c.webhookClause);
          const command = new AddWebhookClauseCommand();
          command.setWebhookClause(clause);
          addWebhookClauseCommands.push(command);
        }
        if (c.clauseType === ClauseType.DATETIME.toString()) {
          const clause = createDatetimeClause(c.datetimeClause);
          const command = new AddDatetimeClauseCommand();
          command.setDatetimeClause(clause);
          addDatetimeClauseCommands.push(command);
        }
      });

    const changeOpsEventRateClauseCommands: ChangeOpsEventRateClauseCommand[] =
      [];
    const changeWebhookClauseCommands: ChangeWebhookClauseCommand[] = [];
    const changeDatetimeClauseCommands: ChangeDatetimeClauseCommand[] = [];
    clauseIds.forEach((cid) => {
      const orgClause = orgRule.clauses.find((c) => c.id === cid);
      const valClause = valRule.clauses.find((c) => c.id === cid);
      if (valClause.clauseType === ClauseType.EVENT_RATE.toString()) {
        const clause = createOpsEventRateClause(valClause.opsEventRateClause);
        const command = new ChangeOpsEventRateClauseCommand();
        command.setId(cid);
        command.setOpsEventRateClause(clause);
        changeOpsEventRateClauseCommands.push(command);
      }
      if (valClause.clauseType === ClauseType.WEBHOOK.toString()) {
        const clause = createWebhookClause(valClause.webhookClause);
        const command = new ChangeWebhookClauseCommand();
        command.setId(cid);
        command.setWebhookClause(clause);
        changeWebhookClauseCommands.push(command);
      }
      if (valClause.clauseType === ClauseType.DATETIME.toString()) {
        const clause = createDatetimeClause(valClause.datetimeClause);
        const command = new ChangeDatetimeClauseCommand();
        command.setId(cid);
        command.setDatetimeClause(clause);
        changeDatetimeClauseCommands.push(command);
      }
    });

    const param: UpdateAutoOpsRuleParams = {
      environmentNamespace: environmentNamespace,
      id: rid,
      changeAutoOpsRuleOpsTypeCommand: changeAutoOpsRuleOpsTypeCommand,
      addOpsEventRateClauseCommands: addOpsEventRateClauseCommands,
      changeOpsEventRateClauseCommands: changeOpsEventRateClauseCommands,
      changeWebhookClauseCommands: changeWebhookClauseCommands,
      addWebhookClauseCommands: addWebhookClauseCommands,
      addDatetimeClauseCommands: addDatetimeClauseCommands,
      changeDatetimeClauseCommands: changeDatetimeClauseCommands,
      deleteClauseCommands: deleteClauseCommands,
    };
    params.push(param);
  });
  return params;
}

export function createAddOpsEventRateClauseCommands(
  org: AutoOpsRuleSchema[],
  val: AutoOpsRuleSchema[]
): string[] {
  const ids: string[] = [];
  const valIds = val.map((r) => r.id);
  org
    .filter((r) => !valIds.includes(r.id))
    .forEach((r) => {
      ids.push(r.id);
    });
  return ids;
}

export function createChangeOpsEventRateClauseCommandsList(
  org: AutoOpsRuleSchema[],
  val: AutoOpsRuleSchema[]
): string[] {
  const ids: string[] = [];
  const valIds = val.map((r) => r.id);
  org
    .filter((r) => !valIds.includes(r.id))
    .forEach((r) => {
      ids.push(r.id);
    });
  return ids;
}

export function createAddDatetimeClauseCommandsList(
  org: AutoOpsRuleSchema[],
  val: AutoOpsRuleSchema[]
): string[] {
  const ids: string[] = [];
  const valIds = val.map((r) => r.id);
  org
    .filter((r) => !valIds.includes(r.id))
    .forEach((r) => {
      ids.push(r.id);
    });
  return ids;
}

export function createChangeDatetimeClauseCommandsList(
  org: AutoOpsRuleSchema[],
  val: AutoOpsRuleSchema[]
): string[] {
  const ids: string[] = [];
  const valIds = val.map((r) => r.id);
  org
    .filter((r) => !valIds.includes(r.id))
    .forEach((r) => {
      ids.push(r.id);
    });
  return ids;
}
