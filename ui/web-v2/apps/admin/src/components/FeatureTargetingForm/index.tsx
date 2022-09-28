import { MinusCircleIcon, XIcon } from '@heroicons/react/solid';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import ReactDatePicker from 'react-datepicker';
import {
  useFormContext,
  Controller,
  useFieldArray,
  useWatch,
} from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import 'react-datepicker/dist/react-datepicker.css';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { listSegments, selectAll } from '../../modules/segments';
import { Clause } from '../../proto/feature/clause_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { Strategy } from '../../proto/feature/strategy_pb';
import { AppDispatch } from '../../store';
import { classNames } from '../../utils/css';
import { CreatableSelect } from '../CreatableSelect';
import { Option, Select } from '../Select';
import { Switch } from '../Switch';

interface FeatureTargetingFormProps {
  featureId: string;
  onOpenConfirmDialog: () => void;
}

export const FeatureTargetingForm: FC<FeatureTargetingFormProps> = memo(
  ({ featureId, onOpenConfirmDialog }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();
    const methods = useFormContext();
    const {
      control,
      formState: { errors, isSubmitting, dirtyFields },
    } = methods;
    const { fields: targets } = useFieldArray({
      control,
      name: 'targets',
    });
    const [feature, _] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ]);
    const strategyOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: v.value,
      };
    });
    strategyOptions.push({
      value: Strategy.Type.ROLLOUT.toString(),
      label: f(messages.feature.strategy.selectRolloutPercentage),
    });
    const offVariationOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: v.value,
      };
    });
    const isValid = Object.keys(errors).length == 0;

    return (
      <div className="p-10 bg-gray-100">
        <form className="">
          <div className="grid grid-cols-1 gap-y-6 gap-x-4">
            <div className="text-sm">{`${f(
              messages.feature.targetingDescription
            )}`}</div>
            <Controller
              name="enabled"
              control={control}
              render={({ field }) => (
                <Switch
                  onChange={field.onChange}
                  size={'large'}
                  readOnly={!editable}
                  enabled={field.value}
                />
              )}
            />
            <div>
              <label className="input-section-label">
                {`${f(messages.feature.targetingUsers)}`}
              </label>
              <div className="bg-white rounded-md p-3 border">
                {targets.map((t: any, idx) => {
                  return (
                    <div key={idx} className="col-span-1">
                      <label htmlFor={`${idx}`} className="input-label">
                        {`${
                          feature.variationsList.find(
                            (v) => v.id == t.variationId
                          )?.value
                        }`}
                      </label>
                      <Controller
                        name={`targets.[${idx}].users`}
                        control={control}
                        render={({ field }) => {
                          return (
                            <CreatableSelect
                              disabled={!editable}
                              defaultValues={field.value.map((u) => {
                                return {
                                  value: u,
                                  label: u,
                                };
                              })}
                              onChange={(options: Option[]) => {
                                field.onChange(options.map((o) => o.value));
                              }}
                            />
                          );
                        }}
                      />
                    </div>
                  );
                })}
              </div>
            </div>
            <div>
              <label className="input-label">{f(messages.feature.rule)}</label>
              <div>
                <RuleInput feature={feature} />
              </div>
            </div>
            <div>
              <label className="input-label">
                {f(messages.feature.defaultStrategy)}
              </label>
              <div className="bg-white p-3 rounded-md border">
                <StrategyInput
                  feature={feature}
                  strategyName={'defaultStrategy'}
                />
                <p className="input-error">
                  {errors.defaultStrategy?.rolloutStrategy?.message && (
                    <span role="alert">
                      {errors.defaultStrategy?.rolloutStrategy?.message}
                    </span>
                  )}
                </p>
              </div>
            </div>
            <div>
              <label htmlFor="offVariation" className="input-label">
                {f(messages.feature.offVariation)}
              </label>
              <div className="bg-white p-3 rounded-md border">
                <Controller
                  name="offVariation"
                  control={control}
                  render={({ field }) => (
                    <Select
                      onChange={field.onChange}
                      options={offVariationOptions}
                      disabled={!editable}
                      value={field.value}
                      isSearchable={false}
                    />
                  )}
                />
              </div>
            </div>
          </div>
          {editable && (
            <div className="py-5">
              <div className="flex justify-end">
                <button
                  type="button"
                  className="btn-submit"
                  disabled={!Object.keys(dirtyFields).length || !isValid}
                  onClick={onOpenConfirmDialog}
                >
                  {f(messages.button.saveWithComment)}
                </button>
              </div>
            </div>
          )}
        </form>
      </div>
    );
  }
);

export interface RuleInputProps {
  feature: Feature.AsObject;
}

export const RuleInput: FC<RuleInputProps> = memo(({ feature }) => {
  const { formatMessage: f } = useIntl();
  const editable = useIsEditable();
  const methods = useFormContext();
  const {
    control,
    formState: { errors },
  } = methods;
  const {
    fields: rules,
    append: appendRule,
    remove,
  } = useFieldArray({
    control,
    name: 'rules',
  });

  const newRolloutStrategy = [];
  feature.variationsList.forEach((val, _) => {
    newRolloutStrategy.push({
      id: val.id,
      percentage: 0,
    });
  });
  const handleAddRule = useCallback(() => {
    appendRule({
      strategy: {
        option: {
          value: feature.variationsList[0].id,
          label: feature.variationsList[0].value,
        },
        rolloutStrategy: newRolloutStrategy,
      },
      clauses: [
        {
          type: ClauseType.COMPARE,
          attribute: '',
          operator: Clause.Operator.EQUALS,
          values: [],
        },
      ],
    });
  }, []);

  const handleRemoveRule = useCallback(
    (idx) => {
      remove(idx);
    },
    [remove]
  );

  return (
    <div>
      <div className="grid grid-cols-1 gap-2">
        {rules.map((r: any, ruleIdx) => {
          return (
            <div
              key={ruleIdx}
              className={classNames('bg-white p-3 rounded-md border')}
            >
              <div key={ruleIdx}>
                <div className="flex mb-2">
                  <label className={classNames()}>{`${f(
                    messages.feature.rule
                  )} ${ruleIdx + 1}`}</label>
                  <div className="flex-grow" />
                  {editable && (
                    <div className="flex items-center">
                      <button
                        type="button"
                        onClick={() => handleRemoveRule(ruleIdx)}
                      >
                        <XIcon className="w-5 h-5" aria-hidden="true" />
                      </button>
                    </div>
                  )}
                </div>
                <ClausesInput ruleIdx={ruleIdx} />
              </div>
              <StrategyInput
                feature={feature}
                strategyName={`rules.${ruleIdx}.strategy`}
              />
              <p className="input-error">
                {errors.rules?.[ruleIdx]?.strategy?.rolloutStrategy
                  ?.message && (
                  <span role="alert">
                    {
                      errors.rules?.[ruleIdx]?.strategy?.rolloutStrategy
                        ?.message
                    }
                  </span>
                )}
              </p>
            </div>
          );
        })}
      </div>
      {editable && (
        <div className="py-4 flex">
          <button type="button" className="btn-submit" onClick={handleAddRule}>
            {f(messages.button.addRule)}
          </button>
        </div>
      )}
    </div>
  );
});

export const ClauseType = {
  COMPARE: 'compare',
  SEGMENT: 'segment',
  DATE: 'date',
} as const;

export type ClauseType = typeof ClauseType[keyof typeof ClauseType];

export const clauseTypeOptions: Option[] = [
  {
    value: ClauseType.COMPARE,
    label: intl.formatMessage(messages.feature.clause.type.compare),
  },
  {
    value: ClauseType.SEGMENT,
    label: intl.formatMessage(messages.feature.clause.type.segment),
  },
  {
    value: ClauseType.DATE,
    label: intl.formatMessage(messages.feature.clause.type.date),
  },
];

export const clauseCompareOperatorOptions: Option[] = [
  {
    value: Clause.Operator.EQUALS.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.equal),
  },
  {
    value: Clause.Operator.GREATER_OR_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.greaterOrEqual),
  },
  {
    value: Clause.Operator.GREATER.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.greater),
  },
  {
    value: Clause.Operator.LESS_OR_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.lessOrEqual),
  },
  {
    value: Clause.Operator.LESS.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.less),
  },
  {
    value: Clause.Operator.IN.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.in),
  },
  {
    value: Clause.Operator.STARTS_WITH.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.startWith),
  },
  {
    value: Clause.Operator.ENDS_WITH.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.endWith),
  },
];

export const clauseDateOperatorOptions: Option[] = [
  {
    value: Clause.Operator.BEFORE.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.before),
  },
  {
    value: Clause.Operator.AFTER.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.after),
  },
];

export interface ClausesInputProps {
  ruleIdx: number;
}

export const ClausesInput: FC<ClausesInputProps> = memo(({ ruleIdx }) => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const editable = useIsEditable();
  const currentEnvironment = useCurrentEnvironment();
  const isSegmentLoading = useSelector<AppState, boolean>(
    (state) => state.segments.loading
  );
  const methods = useFormContext();
  const {
    register,
    control,
    formState: { errors },
  } = methods;
  const clausesName = `rules.${ruleIdx}.clauses`;
  const {
    fields: clauses,
    append,
    remove,
    update,
  } = useFieldArray({
    control,
    name: clausesName,
  });

  const segmentOptions = useSelector<AppState, Option[]>(
    (state) =>
      selectAll(state.segments).map((s) => {
        return {
          value: s.id,
          label: s.name,
        };
      }),
    shallowEqual
  );

  const handleChangeType = useCallback(
    (idx: number, type: string) => {
      switch (type) {
        case ClauseType.COMPARE: {
          update(idx, {
            type: type,
            attribute: '',
            operator: Clause.Operator.EQUALS.toString(),
            values: [],
          });
          break;
        }
        case ClauseType.SEGMENT: {
          update(idx, {
            type: type,
            attribute: '',
            operator: Clause.Operator.SEGMENT.toString(),
            values: [],
          });
          dispatch(
            listSegments({
              environmentNamespace: currentEnvironment.namespace,
              cursor: '',
            })
          );
          break;
        }
        case ClauseType.DATE: {
          const now = String(Math.round(new Date().getTime() / 1000));
          update(idx, {
            type: type,
            attribute: '',
            operator: Clause.Operator.BEFORE.toString(),
            values: [now],
          });
          break;
        }
      }
    },
    [update, dispatch, currentEnvironment]
  );

  const handleAdd = useCallback(() => {
    append({
      type: ClauseType.COMPARE,
      attribute: '',
      operator: Clause.Operator.EQUALS.toString(),
      values: [],
    });
  }, [append]);

  const handleRemove = useCallback(
    (idx) => {
      remove(idx);
    },
    [remove]
  );

  return (
    <div className="grid grid-cols-1 gap-2">
      {clauses.map((c: any, clauseIdx) => {
        const clauseName = `rules.${ruleIdx}.clauses.${clauseIdx}`;
        return (
          <div key={clauseIdx} className={classNames('flex space-x-2')}>
            <div className="w-[2rem] flex justify-center items-center">
              {clauseIdx === 0 ? (
                <div
                  className={classNames(
                    'py-1 px-2 text-xs text-white',
                    'bg-gray-400 mr-3 rounded-full'
                  )}
                >
                  IF
                </div>
              ) : (
                <div className="p-1 text-xs">AND</div>
              )}
            </div>
            <Controller
              name={`${clauseName}.type`}
              control={control}
              render={({ field }) => (
                <Select
                  onChange={(e) => {
                    handleChangeType(clauseIdx, e.value);
                    field.onChange(e.value);
                  }}
                  className={classNames('flex-none w-[200px]')}
                  options={clauseTypeOptions}
                  disabled={!editable}
                  isSearchable={false}
                  value={clauseTypeOptions.find((o) => o.value == c.type)}
                />
              )}
            />
            {c.type == ClauseType.COMPARE && (
              <div className={classNames('flex-grow grid grid-cols-4 gap-1')}>
                <div>
                  <input
                    {...register(`${clauseName}.attribute`)}
                    type="text"
                    defaultValue={c.attribute}
                    className={classNames('input-text w-full')}
                    disabled={!editable}
                  />
                  <p className="input-error">
                    {errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]?.attribute
                      ?.message && (
                      <span role="alert">
                        {
                          errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]
                            ?.attribute?.message
                        }
                      </span>
                    )}
                  </p>
                </div>
                <Controller
                  name={`${clauseName}.operator`}
                  control={control}
                  render={({ field }) => (
                    <Select
                      onChange={(e) => {
                        field.onChange(e.value);
                      }}
                      options={clauseCompareOperatorOptions}
                      disabled={!editable}
                      value={clauseCompareOperatorOptions.find(
                        (o) => o.value === c.operator
                      )}
                    />
                  )}
                />
                <div className="col-span-2">
                  <Controller
                    name={`${clauseName}.values`}
                    control={control}
                    render={({ field }) => {
                      return (
                        <CreatableSelect
                          disabled={!editable}
                          defaultValues={field.value.map((v) => {
                            return {
                              value: v,
                              label: v,
                            };
                          })}
                          onChange={(opts: Option[]) =>
                            field.onChange(opts.map((o) => o.value))
                          }
                        />
                      );
                    }}
                  />
                  <p className="input-error">
                    {errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]?.values
                      ?.message && (
                      <span role="alert">
                        {
                          errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]?.values
                            ?.message
                        }
                      </span>
                    )}
                  </p>
                </div>
              </div>
            )}
            {c.type == ClauseType.SEGMENT &&
              (segmentOptions?.length > 0 ? (
                <div className={classNames('flex-grow grid grid-cols-2 gap-1')}>
                  <div className="flex content-center">
                    <span className="inline-flex items-center text-sm text-gray-700 px-2">
                      {f(messages.feature.clause.operator.segment)}
                    </span>
                  </div>
                  {isSegmentLoading ? (
                    <div>loading</div>
                  ) : (
                    <div>
                      <Controller
                        name={`rules.${ruleIdx}.clauses.${clauseIdx}.values[0]`}
                        control={control}
                        render={({ field }) => {
                          return (
                            <Select
                              onChange={(o: Option) => {
                                field.onChange(o.value);
                              }}
                              options={segmentOptions}
                              disabled={!editable}
                              value={segmentOptions.find(
                                (o) => o.value === c.values[0]
                              )}
                              isSearchable={false}
                            />
                          );
                        }}
                      />
                      <p className="input-error">
                        {errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]?.values
                          ?.message && (
                          <span role="alert">
                            {
                              errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]
                                ?.values?.message
                            }
                          </span>
                        )}
                      </p>
                    </div>
                  )}
                </div>
              ) : (
                <div className="flex-grow flex content-center">
                  <span className="inline-flex items-center text-sm text-gray-700 px-2">
                    {f(messages.segment.select.noData.description)}
                  </span>
                </div>
              ))}
            {c.type == ClauseType.DATE && (
              <div className={classNames('flex-grow grid grid-cols-4 gap-1')}>
                <div>
                  <input
                    {...register(
                      `rules.${ruleIdx}.clauses.${clauseIdx}.attribute`
                    )}
                    type="text"
                    defaultValue={c.attribute}
                    className={classNames('input-text w-full')}
                    disabled={!editable}
                  />
                  <p className="input-error">
                    {errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]?.attribute
                      ?.message && (
                      <span role="alert">
                        {
                          errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]
                            ?.attribute?.message
                        }
                      </span>
                    )}
                  </p>
                </div>
                <Controller
                  name={`rules.${ruleIdx}.clauseIdx.${clauseIdx}.operator`}
                  control={control}
                  render={({ field }) => (
                    <Select
                      onChange={(o: Option) => field.onChange(o.value)}
                      options={clauseDateOperatorOptions}
                      disabled={!editable}
                      value={clauseDateOperatorOptions.find(
                        (o) => o.value === c.operator
                      )}
                      isSearchable={false}
                    />
                  )}
                />
                <div className="col-span-2">
                  <DatetimePicker
                    name={`rules.${ruleIdx}.clauses.${clauseIdx}.values.0`}
                  />
                  <p className="input-error">
                    {errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]?.values
                      ?.message && (
                      <span role="alert">
                        {
                          errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]?.values
                            ?.message
                        }
                      </span>
                    )}
                  </p>
                </div>
              </div>
            )}
            {editable && (
              <div className="flex items-center">
                <button
                  type="button"
                  onClick={() => handleRemove(clauseIdx)}
                  className="minus-circle-icon"
                  disabled={clauses.length <= 1}
                >
                  <MinusCircleIcon aria-hidden="true" />
                </button>
              </div>
            )}
          </div>
        );
      })}

      <div className="py-4 flex">
        {editable && (
          <button type="button" className="btn-submit" onClick={handleAdd}>
            {f(messages.button.addCondition)}
          </button>
        )}
      </div>
    </div>
  );
});

export interface StrategyInputProps {
  feature: Feature.AsObject;
  strategyName: string;
}

export const StrategyInput: FC<StrategyInputProps> = memo(
  ({ feature, strategyName }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();
    const methods = useFormContext();
    const {
      register,
      control,
      formState: { errors },
      trigger,
    } = methods;
    const optionName = `${strategyName}.option`;
    const rolloutStrategyName = `${strategyName}.rolloutStrategy`;
    const selectedOption = useWatch({
      control,
      name: optionName,
    });
    const { fields: rolloutStrategy, update } = useFieldArray({
      control,
      name: rolloutStrategyName,
      keyName: 'key', // the default keyName is "id" and it conflicts with the variation id field
    });
    const strategyOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: v.value,
      };
    });
    strategyOptions.push({
      value: String(Strategy.Type.ROLLOUT),
      label: f(messages.feature.strategy.selectRolloutPercentage),
    });
    const handleOnChange = useCallback(
      (idx: number, id: string, e: React.ChangeEvent<HTMLInputElement>) => {
        update(idx, {
          id: id,
          percentage: Number(e.target.value),
        });
        trigger(strategyName);
      },
      [update, trigger]
    );

    return (
      <div>
        <Controller
          name={optionName}
          control={control}
          render={({ field }) => (
            <Select
              options={strategyOptions}
              disabled={!editable}
              value={selectedOption}
              onChange={field.onChange}
              isSearchable={false}
            />
          )}
        />
        {selectedOption.value == Strategy.Type.ROLLOUT && (
          <div className="grid grid-cols-1 gap-2 mt-2">
            {rolloutStrategy.map((s: any, idx: number) => {
              return (
                <div key={s.id} className="flex">
                  <div className="w-36 flex">
                    <input
                      {...register(`${rolloutStrategyName}.${idx}.percentage`)}
                      type="number"
                      min="0"
                      max="100"
                      defaultValue={s.percentage}
                      className={classNames(
                        'flex-grow pr-0 py-1 rounded-l border border-r-0 border-gray-300',
                        'text-right text-sm text-gray-700'
                      )}
                      placeholder={''}
                      onChange={(e) => handleOnChange(idx, s.id, e)}
                      disabled={!editable}
                    />
                    <span
                      className={classNames(
                        'px-1 py-1 inline-flex items-center bg-gray-100',
                        'rounded-r border border-l-0 border-gray-300 text-gray-700'
                      )}
                    >
                      {'%'}
                    </span>
                  </div>
                  <label className="inline-flex items-center ml-3 text-sm text-gray-700">
                    {feature.variationsList.find((v) => v.id == s.id).value}
                  </label>
                </div>
              );
            })}
            <div className="w-36 flex">
              <span
                className={classNames(
                  'w-14 px-3 py-1 inline-flex items-center bg-gray-100',
                  'rounded-l border border-r-0 border-gray-300',
                  'text-sm text-gray-700'
                )}
              >
                {f(messages.total)}
              </span>
              <div
                className={classNames(
                  'flex-grow text-right pr-4',
                  'border border-l-0 border-r-0 border-gray-300',
                  'text-sm text-gray-700'
                )}
              >
                <span
                  className={classNames(
                    'pr-0 py-1 inline-flex items-center',
                    'text-right'
                  )}
                >
                  {rolloutStrategy
                    .map((s: any) => Number(s.percentage))
                    .reduce((previousValue, currentValue) => {
                      return previousValue + currentValue;
                    })}
                </span>
              </div>
              <span
                className={classNames(
                  'px-1 py-1 inline-flex items-center bg-gray-100',
                  'rounded-r border border-l-0 border-gray-300',
                  'text-sm text-gray-700'
                )}
              >
                {'%'}
              </span>
            </div>
          </div>
        )}
      </div>
    );
  }
);

export interface DatetimePickerProps {
  name: string;
}

export const DatetimePicker: FC<DatetimePickerProps> = memo(({ name }) => {
  const editable = useIsEditable();
  const methods = useFormContext();
  const { control } = methods;

  return (
    <Controller
      control={control}
      name={name}
      render={({ field }) => (
        <ReactDatePicker
          dateFormat="yyyy-MM-dd HH:mm"
          showTimeSelect
          timeIntervals={60}
          placeholderText=""
          className={classNames('input-text w-full')}
          wrapperClassName="w-full"
          onChange={(v) => {
            field.onChange(v.getTime() / 1000);
          }}
          selected={(() => {
            return field.value ? new Date(Number(field.value) * 1000) : null;
          })()}
          disabled={!editable}
        />
      )}
    />
  );
});
