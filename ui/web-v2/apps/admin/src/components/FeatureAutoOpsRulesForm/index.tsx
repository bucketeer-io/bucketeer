import { WEBHOOK_LIST_PAGE_SIZE } from '@/constants/webhook';
import { listWebhooks } from '@/modules/webhooks';
import { Webhook } from '@/proto/autoops/webhook_pb';
import { MinusCircleIcon, XIcon } from '@heroicons/react/solid';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useEffect, useCallback } from 'react';
import {
  useFormContext,
  Controller,
  useFieldArray,
  useWatch,
} from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useSelector, shallowEqual, useDispatch } from 'react-redux';
import { v4 as uuid } from 'uuid';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { selectAll as selectAllGoals } from '../../modules/goals';
import { useIsEditable, useCurrentEnvironment } from '../../modules/me';
import { selectAll } from '../../modules/webhooks';
import { OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import {
  OpsEventRateClause,
  WebhookClause,
} from '../../proto/autoops/clause_pb';
import { Goal } from '../../proto/experiment/goal_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { AppDispatch } from '../../store';
import { classNames } from '../../utils/css';
import { DatetimePicker } from '../DatetimePicker';
import { DetailSkeleton } from '../DetailSkeleton';
import { Option, Select } from '../Select';

export interface ClauseTypeMap {
  EVENT_RATE: 'bucketeer.autoops.OpsEventRateClause';
  DATETIME: 'bucketeer.autoops.DatetimeClause';
  WEBHOOK: 'bucketeer.autoops.WebhookClause';
}

export const ClauseType: ClauseTypeMap = {
  EVENT_RATE: 'bucketeer.autoops.OpsEventRateClause',
  DATETIME: 'bucketeer.autoops.DatetimeClause',
  WEBHOOK: 'bucketeer.autoops.WebhookClause',
};

interface FeatureAutoOpsRulesFormProps {
  featureId: string;
  onSubmit: () => void;
}

export const FeatureAutoOpsRulesForm: FC<FeatureAutoOpsRulesFormProps> = memo(
  ({ featureId, onSubmit }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      formState: { isDirty, isValid },
    } = methods;

    return (
      <div className="p-10 bg-gray-100">
        <form className="">
          <div className="grid grid-cols-1 gap-y-6 gap-x-4">
            <AutoOpsRulesInput featureId={featureId} />
          </div>
          {editable && (
            <div>
              <div className="flex justify-end">
                <button
                  type="button"
                  className="btn-submit"
                  disabled={!isDirty || !isValid}
                  onClick={onSubmit}
                >
                  {f(messages.button.submit)}
                </button>
              </div>
            </div>
          )}
        </form>
      </div>
    );
  }
);

export const opsTypeOptions = [
  {
    value: OpsType.ENABLE_FEATURE.toString(),
    label: intl.formatMessage(messages.autoOps.enableFeatureType),
  },
  {
    value: OpsType.DISABLE_FEATURE.toString(),
    label: intl.formatMessage(messages.autoOps.disableFeatureType),
  },
];

export const createInitialAutoOpsRule = (feature: Feature.AsObject) => {
  return {
    id: uuid(),
    featureId: feature.id,
    triggeredAt: 0,
    opsType: opsTypeOptions[0].value,
    clauses: [createInitialClause(feature)],
  };
};

export const createInitialOpsEventRateClause = (feature: Feature.AsObject) => {
  return {
    variation: feature.variationsList[0].id,
    goal: null,
    minCount: 1,
    threadsholdRate: 50,
    operator: operatorOptions[0].value,
  };
};

export const createInitialWebhookClause = () => {
  return {
    id: uuid(),
    webhookId: null,
    conditionsList: [createInitialCondition()],
  };
};

export const createInitialCondition = () => {
  return {
    id: uuid(),
    filter: null,
    operator: webhookOperatorOptions[0].value,
    value: null,
  };
};

export const createInitialDatetimeClause = () => {
  const date = new Date();
  date.setDate(date.getDate() + 1);
  return {
    time: date,
  };
};

export const createInitialClause = (feature: Feature.AsObject) => {
  return {
    id: uuid(),
    clauseType: ClauseType.DATETIME.toString(),
    datetimeClause: createInitialDatetimeClause(),
    opsEventRateClause: createInitialOpsEventRateClause(feature),
    webhookClause: createInitialWebhookClause(),
  };
};

export interface AutoOpsRulesInputProps {
  featureId: string;
}

export const AutoOpsRulesInput: FC<AutoOpsRulesInputProps> = memo(
  ({ featureId }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const [feature, getFeatureError] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ]);
    const methods = useFormContext();
    const {
      control,
      formState: { errors },
    } = methods;
    const {
      fields: rules,
      append,
      remove,
    } = useFieldArray({
      control,
      name: 'autoOpsRules',
    });

    const handleAdd = useCallback(() => {
      append(createInitialAutoOpsRule(feature));
    }, [append]);

    const handleRemove = useCallback(
      (idx) => {
        remove(idx);
      },
      [remove]
    );

    return (
      <div>
        <div className="grid grid-cols-1 gap-2">
          {rules.map((rule: any, ruleIdx) => {
            return (
              <div
                key={rule.id}
                className={classNames('bg-white p-3 rounded-md border, mb-5')}
              >
                <div className="flex text-gray-700 pb-3">
                  <div>
                    <label className={classNames('text-sm')}>{`${f(
                      messages.autoOps.operation
                    )} ${ruleIdx + 1}`}</label>
                  </div>
                  <div className="flex-grow" />
                  {editable && (
                    <div className="flex items-center">
                      <button
                        type="button"
                        className="x-icon"
                        onClick={() => handleRemove(ruleIdx)}
                      >
                        <XIcon className="w-5 h-5" aria-hidden="true" />
                      </button>
                    </div>
                  )}
                </div>
                <AutoOpsRuleInput
                  featureId={featureId}
                  ruleIdx={ruleIdx}
                  feature={feature}
                />
              </div>
            );
          })}
        </div>
        {editable && (
          <div className="flex">
            <button type="button" className="btn-submit" onClick={handleAdd}>
              {f(messages.button.addOperation)}
            </button>
          </div>
        )}
      </div>
    );
  }
);

export interface AutoOpsRuleInputProps {
  featureId: string;
  ruleIdx: number;
  feature: Feature.AsObject;
}

export const AutoOpsRuleInput: FC<AutoOpsRuleInputProps> = memo(
  ({ featureId, ruleIdx, feature }) => {
    const editable = useIsEditable();
    const ruleName = `autoOpsRules.${ruleIdx}`;
    const methods = useFormContext();
    const { control, setValue } = methods;
    const rule = useWatch({
      control,
      name: ruleName,
    });

    return (
      <>
        <Controller
          name={`${ruleName}.opsType`}
          control={control}
          render={({ field }) => (
            <Select
              onChange={(o: Option) => {
                setValue(`autoOpsRules.${ruleIdx}.clauses`, [
                  createInitialClause(feature),
                ]);
                field.onChange(o.value);
              }}
              options={opsTypeOptions}
              disabled={!editable}
              value={opsTypeOptions.find((o) => o.value == rule.opsType)}
            />
          )}
        />
        <ClausesInput featureId={featureId} ruleIdx={ruleIdx} />
      </>
    );
  }
);

export interface ClausesInputProps {
  featureId: string;
  ruleIdx: number;
}

export const ClausesInput: FC<ClausesInputProps> = ({ featureId, ruleIdx }) => {
  const editable = useIsEditable();
  const [feature, getFeatureError] = useSelector<
    AppState,
    [Feature.AsObject | undefined, SerializedError | null]
  >((state) => [
    selectFeatureById(state.features, featureId),
    state.features.getFeatureError,
  ]);

  const { formatMessage: f } = useIntl();
  const methods = useFormContext();
  const { control } = methods;
  const ruleName = `autoOpsRules.${ruleIdx}`;
  const clausesName = `${ruleName}.clauses`;
  const watchClauses = useWatch({
    control,
    name: clausesName,
  });
  const {
    append,
    remove,
    fields: clauses,
  } = useFieldArray({
    control,
    name: clausesName,
  });

  const handleAddCondition = useCallback(() => {
    append(createInitialClause(feature));
  }, [append]);

  const handleRemove = useCallback(
    (idx) => {
      remove(idx);
    },
    [remove]
  );

  return (
    <div>
      {clauses.map((clause: any, clauseIdx: number) => {
        return (
          <div key={clause.id}>
            <div className={classNames('flex space-x-2')}>
              <div className="w-[2rem] flex justify-center items-center">
                {clauseIdx === 0 ? (
                  <div
                    className={classNames(
                      'py-1 px-2',
                      'text-xs bg-gray-400 text-white rounded-full'
                    )}
                  >
                    IF
                  </div>
                ) : (
                  <div className="p-1 text-xs">OR</div>
                )}
              </div>
              <div className="flex-grow flex mt-3 p-3 rounded-md border">
                <div className="flex-grow">
                  <ClauseInput
                    featureId={featureId}
                    ruleIdx={ruleIdx}
                    clauseIdx={clauseIdx}
                    feature={feature}
                  />
                </div>
                {editable && (
                  <div className="flex items-start pl-2">
                    <button
                      type="button"
                      className="x-icon"
                      onClick={() => handleRemove(clauseIdx)}
                    >
                      <XIcon className="w-5 h-5" aria-hidden="true" />
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>
        );
      })}
      {editable && !hideAddConditionBtn(watchClauses) && (
        <div className="py-4 flex">
          <button
            type="button"
            className="btn-submit"
            onClick={handleAddCondition}
          >
            {f(messages.button.addCondition)}
          </button>
        </div>
      )}
    </div>
  );
};

function hideAddConditionBtn(clauses): boolean {
  for (const clause of clauses) {
    if (clause.clauseType === ClauseType.DATETIME.toString()) {
      return true;
    } else if (clause.clauseType === ClauseType.WEBHOOK.toString()) {
      return true;
    }
  }
  return false;
}

export const clauseTypeOptionEventRate = {
  value: ClauseType.EVENT_RATE.toString(),
  label: intl.formatMessage(messages.autoOps.eventRateClauseType),
};

export const clauseTypeOptionDatetime = {
  value: ClauseType.DATETIME.toString(),
  label: intl.formatMessage(messages.autoOps.datetimeClauseType),
};

export const clauseTypeOptionWebhook = {
  value: ClauseType.WEBHOOK.toString(),
  label: intl.formatMessage(messages.autoOps.webhookClauseType),
};

export const clauseTypeOptions = [
  clauseTypeOptionEventRate,
  clauseTypeOptionDatetime,
  clauseTypeOptionWebhook,
];

export const createClauseTypeOption = (
  clauseType: ClauseTypeMap[keyof ClauseTypeMap]
) => {
  return clauseTypeOptions.find(
    (option) => clauseType.toString() == option.value
  );
};

export interface ClauseInputProps {
  featureId: string;
  ruleIdx: number;
  clauseIdx: number;
  feature: Feature.AsObject;
}

export const ClauseInput: FC<ClauseInputProps> = ({
  featureId,
  ruleIdx,
  clauseIdx,
  feature,
}) => {
  const editable = useIsEditable();
  const methods = useFormContext();
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();

  const {
    control,
    formState: { errors },
    setValue,
  } = methods;
  const ruleName = `autoOpsRules.${ruleIdx}`;
  const clauseName = `${ruleName}.clauses.${clauseIdx}`;
  const opsType = useWatch({
    control,
    name: `${ruleName}.opsType`,
  });

  const selectedClauseTypeOptions =
    opsType === OpsType.ENABLE_FEATURE.toString()
      ? [clauseTypeOptionDatetime, clauseTypeOptionWebhook]
      : [
          clauseTypeOptionEventRate,
          clauseTypeOptionWebhook,
          clauseTypeOptionDatetime,
        ];

  const clauseType = useWatch({
    control,
    name: `${clauseName}.clauseType`,
  });
  const clause = useWatch({
    control,
    name: clauseName,
  });

  const webhookList = useSelector<AppState, Webhook.AsObject[]>(
    (state) => selectAll(state.webhook),
    shallowEqual
  );

  useEffect(() => {
    dispatch(
      listWebhooks({
        environmentNamespace: currentEnvironment.namespace,
        pageSize: WEBHOOK_LIST_PAGE_SIZE,
        cursor: String(0),
      })
    );
  }, []);

  const webhookListOptions = webhookList.map((webhook) => ({
    label: webhook.name,
    value: webhook.id,
  }));

  return (
    <div className="grid grid-cols-1 gap-2">
      <div className="">
        <Controller
          name={`${clauseName}.clauseType`}
          control={control}
          render={({ field }) => (
            <Select
              onChange={(o: Option) => {
                setValue(`${ruleName}.clauses`, [
                  {
                    id: uuid(),
                    clauseType: o.value.toString(),
                    datetimeClause: createInitialDatetimeClause(),
                    opsEventRateClause:
                      createInitialOpsEventRateClause(feature),
                    webhookClause: createInitialWebhookClause(),
                  },
                ]);
                field.onChange(o.value);
              }}
              options={selectedClauseTypeOptions}
              disabled={!editable}
              value={selectedClauseTypeOptions.find(
                (o) => o.value === clauseType
              )}
            />
          )}
        />
        {clauseType === ClauseType.EVENT_RATE.toString() && (
          <EventRateClauseInput
            featureId={featureId}
            ruleIdx={ruleIdx}
            clauseIdx={clauseIdx}
          />
        )}
        {clauseType === ClauseType.DATETIME.toString() && (
          <DatetimeClauseInput ruleIdx={ruleIdx} clauseIdx={clauseIdx} />
        )}
        {clauseType === ClauseType.WEBHOOK.toString() && (
          <div>
            <div className="pt-1">
              <label htmlFor="webhookName">
                <span className="input-label">Webhook</span>
              </label>
              <Controller
                name={`${clauseName}.webhookClause.webhookId`}
                control={control}
                render={({ field }) => (
                  <Select
                    onChange={(o: Option) => field.onChange(o.value)}
                    options={webhookListOptions}
                    disabled={!editable}
                    value={webhookListOptions.find(
                      (o) => o.value === clause.webhookClause.webhookId
                    )}
                  />
                )}
              />
            </div>
            <WebhookClauseInput ruleIdx={ruleIdx} clauseIdx={clauseIdx} />
          </div>
        )}
      </div>
    </div>
  );
};

export const operatorOptions = [
  {
    value: OpsEventRateClause.Operator.GREATER_OR_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.greaterOrEqual),
  },
  {
    value: OpsEventRateClause.Operator.LESS_OR_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.lessOrEqual),
  },
];

export const webhookOperatorOptions = [
  {
    value: WebhookClause.Condition.Operator.EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.equal),
  },
  {
    value: WebhookClause.Condition.Operator.NOT_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.notEqual),
  },
  {
    value: WebhookClause.Condition.Operator.MORE_THAN.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.greater),
  },
  {
    value: WebhookClause.Condition.Operator.MORE_THAN_OR_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.greaterOrEqual),
  },
  {
    value: WebhookClause.Condition.Operator.LESS_THAN.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.less),
  },
  {
    value: WebhookClause.Condition.Operator.LESS_THAN_OR_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.lessOrEqual),
  },
];

export const createOperatorOption = (
  operator: OpsEventRateClause.OperatorMap[keyof OpsEventRateClause.OperatorMap]
) => {
  return operatorOptions.find((option) => option.value === operator.toString());
};

export interface EventRateClauseInputProps {
  featureId: string;
  ruleIdx: number;
  clauseIdx: number;
}

export const EventRateClauseInput: FC<EventRateClauseInputProps> = memo(
  ({ featureId, ruleIdx, clauseIdx }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const isGoalLoading = useSelector<AppState, boolean>(
      (state) => state.goals.loading,
      shallowEqual
    );
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
    const methods = useFormContext();
    const opsEventRateClauseName = `autoOpsRules.${ruleIdx}.clauses.${clauseIdx}.opsEventRateClause`;
    const {
      register,
      control,
      formState: { errors },
      trigger,
    } = methods;
    const [feature, _] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ]);
    const clause = useWatch({
      control,
      name: opsEventRateClauseName,
    });
    const variationOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: v.value,
      };
    });

    if (isGoalLoading) {
      return (
        <div className="p-9 bg-gray-100">
          <DetailSkeleton />
        </div>
      );
    }
    return (
      <div className="grid grid-cols-1 gap-2">
        <div>
          <label htmlFor="variation" className="input-label">
            {f(messages.feature.variation)}
          </label>
          <Controller
            name={`${opsEventRateClauseName}.variation`}
            control={control}
            render={({ field }) => (
              <Select
                onChange={(o: Option) => field.onChange(o.value)}
                options={variationOptions}
                disabled={!editable}
                value={variationOptions.find((o) => o.value === field.value)}
              />
            )}
          />
        </div>
        <label htmlFor="variation" className="input-label">
          {f(messages.autoOps.opsEventRateClause.goal)}
        </label>
        <div className={classNames('flex-grow grid grid-cols-4 gap-1')}>
          <Controller
            name={`${opsEventRateClauseName}.goal`}
            control={control}
            render={({ field }) => (
              <Select
                onChange={(o: Option) => field.onChange(o.value)}
                options={goalOptions}
                disabled={!editable}
                value={goalOptions.find((o) => o.value === clause.goal)}
              />
            )}
          />
          <Controller
            name={`${opsEventRateClauseName}.operator`}
            control={control}
            render={({ field }) => (
              <Select
                onChange={(o: Option) => field.onChange(o.value)}
                options={operatorOptions}
                disabled={!editable}
                value={operatorOptions.find((o) => o.value === clause.operator)}
              />
            )}
          />
          <div className="w-36 flex">
            <input
              {...register(`${opsEventRateClauseName}.threadsholdRate`)}
              type="number"
              min="0"
              max="100"
              defaultValue={clause.threadsholdRate}
              className={classNames(
                'flex-grow pr-0 py-1',
                'rounded-l border border-r-0 border-gray-300',
                'text-right'
              )}
              placeholder={''}
              required
              disabled={!editable}
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
          <p className="input-error">
            {errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
              ?.opsEventRateClause?.threadsholdRate?.message && (
              <span role="alert">
                {
                  errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
                    ?.opsEventRateClause?.threadsholdRate?.message
                }
              </span>
            )}
          </p>
        </div>
        <div className="w-36">
          <label htmlFor="name">
            <span className="input-label">
              {f(messages.autoOps.opsEventRateClause.minCount)}
            </span>
          </label>
          <div className="mt-1">
            <input
              {...register(`${opsEventRateClauseName}.minCount`)}
              type="number"
              min="0"
              className="input-text w-full"
              disabled={!editable}
            />
            <p className="input-error">
              {errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
                ?.opsEventRateClause?.minCount?.message && (
                <span role="alert">
                  {
                    errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
                      ?.opsEventRateClause?.minCount?.message
                  }
                </span>
              )}
            </p>
          </div>
        </div>
      </div>
    );
  }
);

export interface DatetimeClauseInputProps {
  ruleIdx: number;
  clauseIdx: number;
}

export const DatetimeClauseInput: FC<DatetimeClauseInputProps> = memo(
  ({ ruleIdx, clauseIdx }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const clauseName = `autoOpsRules.${ruleIdx}.clauses.${clauseIdx}`;
    const {
      formState: { errors },
    } = methods;

    return (
      <div className="">
        <label htmlFor="name">
          <span className="input-label">
            {f(messages.autoOps.datetimeClause.datetime)}
          </span>
        </label>
        <DatetimePicker
          name={`${clauseName}.datetimeClause.time`}
          disabled={!editable}
        />
        <p className="input-error">
          {errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]?.datetimeClause
            ?.time?.message && (
            <span role="alert">
              {
                errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
                  ?.datetimeClause?.time?.message
              }
            </span>
          )}
        </p>
      </div>
    );
  }
);

export interface WebhookClauseInputProps {
  ruleIdx: number;
  clauseIdx: number;
}

export const WebhookClauseInput: FC<WebhookClauseInputProps> = memo(
  ({ ruleIdx, clauseIdx }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const webhookClauseName = `autoOpsRules.${ruleIdx}.clauses.${clauseIdx}.webhookClause`;
    const webhookConditionsListClauseName = `${webhookClauseName}.conditionsList`;

    const {
      register,
      control,
      formState: { errors },
    } = methods;

    const {
      append,
      remove,
      fields: conditionsList,
    } = useFieldArray({
      control,
      name: webhookConditionsListClauseName,
    });

    const handleAddClauseCondition = useCallback(() => {
      append(createInitialCondition());
    }, [append]);

    const handleRemoveCondition = useCallback(
      (idx) => {
        remove(idx);
      },
      [remove]
    );

    return (
      <div className="">
        {conditionsList.map((condition: any, conditionIdx) => (
          <div key={condition.id} className="flex space-x-2 mt-2">
            <div className="w-14 self-center flex justify-center">
              <div
                className={classNames(
                  'py-1 px-2',
                  'text-xs bg-gray-400 text-white rounded-full'
                )}
              >
                {conditionIdx === 0 ? 'WHERE' : 'AND'}
              </div>
            </div>
            <div className="flex-grow flex mt-3 p-3 rounded-md border">
              <div className="flex-grow space-y-2">
                <div>
                  <input
                    {...register(
                      `${webhookConditionsListClauseName}.${conditionIdx}.filter`
                    )}
                    type="text"
                    placeholder={intl.formatMessage(
                      messages.autoOps.webhookClause.filter
                    )}
                    defaultValue={condition.filter}
                    className={classNames('input-text w-full')}
                    disabled={!editable}
                  />
                  <p className="input-error">
                    {errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
                      ?.webhookClause?.conditionsList[conditionIdx]?.filter
                      ?.message && (
                      <span role="alert">
                        {
                          errors.autoOpsRules[ruleIdx].clauses[clauseIdx]
                            .webhookClause.conditionsList[conditionIdx].filter
                            .message
                        }
                      </span>
                    )}
                  </p>
                </div>
                <Controller
                  name={`${webhookConditionsListClauseName}.${conditionIdx}.operator`}
                  control={control}
                  render={({ field }) => (
                    <Select
                      onChange={(o: Option) => field.onChange(o.value)}
                      options={webhookOperatorOptions}
                      disabled={!editable}
                      value={webhookOperatorOptions.find(
                        (o) => o.value === condition.operator
                      )}
                    />
                  )}
                />
                <div>
                  <input
                    {...register(
                      `${webhookConditionsListClauseName}.${conditionIdx}.value`
                    )}
                    type="text"
                    defaultValue={condition.value}
                    className={classNames('input-text w-full')}
                    disabled={!editable}
                    placeholder={intl.formatMessage(
                      messages.autoOps.webhookClause.value
                    )}
                  />
                  <p className="input-error">
                    {errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
                      ?.webhookClause?.conditionsList[conditionIdx]?.value
                      ?.message && (
                      <span role="alert">
                        {
                          errors.autoOpsRules[ruleIdx].clauses[clauseIdx]
                            .webhookClause.conditionsList[conditionIdx].value
                            .message
                        }
                      </span>
                    )}
                  </p>
                </div>
              </div>
              {editable && (
                <div className="flex items-start pl-2">
                  <button
                    type="button"
                    className="x-icon"
                    onClick={() => handleRemoveCondition(conditionIdx)}
                  >
                    <XIcon className="w-5 h-5" aria-hidden="true" />
                  </button>
                </div>
              )}
            </div>
          </div>
        ))}
        {editable && (
          <div className="pt-4 flex">
            <button
              type="button"
              className="btn-submit"
              onClick={handleAddClauseCondition}
            >
              {f(messages.button.addCondition)}
            </button>
          </div>
        )}
      </div>
    );
  }
);
