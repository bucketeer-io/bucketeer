import { PAGE_PATH_FEATURES, PAGE_PATH_ROOT } from '@/constants/routing';
import { isLanguageJapanese } from '@/lang/getSelectedLanguage';
import { ListFeaturesRequest } from '@/proto/feature/service_pb';
import {
  createVariationLabel,
  getAlreadyTargetedVariation,
} from '@/utils/variation';
import {
  MinusCircleIcon,
  XIcon,
  InformationCircleIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  PlusCircleIcon,
  ArrowUpIcon,
  ArrowDownIcon,
} from '@heroicons/react/solid';
import { FileCopyOutlined } from '@material-ui/icons';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useState, useEffect } from 'react';
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
import { Link } from 'react-router-dom';
import { components } from 'react-select';
import ReactCreatableSelect from 'react-select/creatable';
import { v4 as uuid } from 'uuid';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectById as selectFeatureById,
  selectAll as selectAllFeatures,
  listFeatures,
} from '../../modules/features';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { listSegments, selectAll } from '../../modules/segments';
import { Clause } from '../../proto/feature/clause_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { Strategy } from '../../proto/feature/strategy_pb';
import { AppDispatch } from '../../store';
import { classNames } from '../../utils/css';
import { CopyChip } from '../CopyChip';
import { colourStyles, CreatableSelect } from '../CreatableSelect';
import { Option, Select } from '../Select';
import { OptionFeatureFlag, SelectFeatureFlag } from '../SelectFeatureFlag';
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
      formState: { errors, isDirty },
      watch,
    } = methods;

    const prerequisites = watch('prerequisites');

    const rules = watch('rules');
    const targets = watch('targets');

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
        label: createVariationLabel(v),
      };
    });
    strategyOptions.push({
      value: Strategy.Type.ROLLOUT.toString(),
      label: f(messages.feature.strategy.selectRolloutPercentage),
    });
    const offVariationOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: createVariationLabel(v),
      };
    });

    const checkSaveBtnDisabled = useCallback(() => {
      // check if all prerequisites fields are dirty
      const checkPrerequisites = prerequisites.every(
        (p) => p.featureId && p.variationId
      );

      // check if all rules fields are dirty
      const checkRules = rules.every((rule) =>
        rule.clauses.every((clause) => {
          if (clause.type === ClauseType.SEGMENT) {
            return clause.values.length > 0;
          }
          return clause.attribute && clause.values.length > 0;
        })
      );

      if (
        !checkPrerequisites ||
        !checkRules ||
        Object.values(errors).some(Boolean) ||
        !isDirty
      ) {
        return true;
      }
      return false;
    }, [rules, isDirty, errors, prerequisites]);

    const handleOnPaste = (e, t, field) => {
      // Stop data actually being pasted into div
      e.stopPropagation();
      e.preventDefault();

      const clipboardData = e.clipboardData;
      const pastedData: string = clipboardData.getData('Text').split(', ');

      if (pastedData) {
        const difference = t.users.filter((u) => !pastedData.includes(u));
        field.onChange([...difference, ...pastedData]);
      }
    };

    const NoOptionsMessage = ({ props }) => {
      return (
        <components.NoOptionsMessage {...props}>
          <span className="custom-css-class">
            {props.selectProps.inputValue
              ? f(messages.feature.alreadyTargeted)
              : f(messages.feature.addUserIds)}
          </span>
        </components.NoOptionsMessage>
      );
    };

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
            <FlagIsPrerequisite featureId={featureId} />
            <div>
              <label className="input-section-label">
                {`${f(messages.feature.prerequisites)}`}
              </label>
              <PrerequisiteInput feature={feature} />
            </div>
            <div>
              <label className="input-section-label">
                {`${f(messages.feature.targetingUsers)}`}
              </label>
              <div className="bg-white rounded-md p-3 border">
                {targets.map((t: any, idx) => {
                  return (
                    <div key={idx} className="col-span-1">
                      <div className="truncate">
                        <label
                          htmlFor={`${idx}`}
                          className="input-label w-full"
                        >
                          {createVariationLabel(
                            feature.variationsList.find(
                              (v) => v.id == t.variationId
                            )
                          )}
                        </label>
                      </div>
                      <div className="flex space-x-2">
                        <Controller
                          name={`targets.[${idx}].users`}
                          control={control}
                          render={({ field }) => {
                            return (
                              <div
                                className="flex-1"
                                onPaste={(e) => handleOnPaste(e, t, field)}
                              >
                                <ReactCreatableSelect
                                  isMulti
                                  placeholder={f(messages.feature.addUserIds)}
                                  classNamePrefix="react-select"
                                  styles={colourStyles}
                                  formatCreateLabel={(userInput) => {
                                    const alreadyTargetedVariaition =
                                      getAlreadyTargetedVariation(
                                        targets,
                                        t.variationId,
                                        userInput
                                      );
                                    if (alreadyTargetedVariaition) {
                                      let variationName = createVariationLabel(
                                        feature.variationsList.find(
                                          (v) =>
                                            v.id ===
                                            alreadyTargetedVariaition.variationId
                                        )
                                      );
                                      variationName =
                                        variationName.length > 50
                                          ? `${variationName.slice(0, 50)} ...`
                                          : variationName;
                                      return (
                                        <div
                                          className={
                                            'text-center text-gray-500'
                                          }
                                        >
                                          <span>
                                            {f(
                                              messages.feature
                                                .alreadyTargetedInVariation,
                                              {
                                                userId: userInput,
                                                variationName,
                                              }
                                            )}
                                          </span>
                                        </div>
                                      );
                                    }

                                    return (
                                      <div className="flex space-x-1 items-center">
                                        <PlusCircleIcon
                                          className="w-4 h-4 text-blue-400"
                                          aria-hidden="true"
                                        />

                                        <span className="text-blue-700">
                                          {f(messages.feature.addUser, {
                                            userId: userInput,
                                          })}
                                        </span>
                                      </div>
                                    );
                                  }}
                                  components={{
                                    DropdownIndicator: null,
                                    NoOptionsMessage: (props) => (
                                      <NoOptionsMessage props={props} />
                                    ),
                                  }}
                                  value={field.value.map((u) => {
                                    return {
                                      value: u,
                                      label: u,
                                    };
                                  })}
                                  onChange={(options: Option[]) => {
                                    const newOption = options.find(
                                      (o) => o['__isNew__']
                                    );

                                    const alreadyTargetedVariaition =
                                      getAlreadyTargetedVariation(
                                        targets,
                                        t.variationId,
                                        newOption?.label
                                      );

                                    if (!alreadyTargetedVariaition) {
                                      field.onChange(
                                        options.map((o) => o.value)
                                      );
                                    }
                                  }}
                                  disabled={!editable}
                                />
                              </div>
                            );
                          }}
                        />
                        <CopyChip text={t.users.join(', ')}>
                          <div className="flex items-center border border-[#D1D5DB] cursor-pointer hover:bg-gray-50 transition px-2 h-full rounded">
                            <FileCopyOutlined
                              aria-hidden="true"
                              fontSize="small"
                              className="text-gray-400"
                            />
                          </div>
                        </CopyChip>
                      </div>
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
                  disabled={checkSaveBtnDisabled()}
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

export interface FlagIsPrerequisiteProps {
  featureId: string;
}

const FlagIsPrerequisite: FC<FlagIsPrerequisiteProps> = ({ featureId }) => {
  const [isSeeMore, setSeeMore] = useState(false);
  const { formatMessage: f } = useIntl();

  const currentEnvironment = useCurrentEnvironment();

  const features = useSelector<AppState, Feature.AsObject[]>(
    (state) => selectAllFeatures(state.features),
    shallowEqual
  );

  useEffect(() => {
    listFeatures({
      environmentNamespace: currentEnvironment.id,
      pageSize: 0,
      cursor: '',
      tags: [],
      searchKeyword: null,
      maintainerId: null,
      orderBy: ListFeaturesRequest.OrderBy.DEFAULT,
      orderDirection: ListFeaturesRequest.OrderDirection.ASC,
    });
  }, []);

  const flagList = features.reduce((acc, feature) => {
    if (
      feature.prerequisitesList.find(
        (prerequisite) => prerequisite.featureId === featureId
      )
    ) {
      return [
        ...acc,
        {
          id: feature.id,
          name: feature.name,
        },
      ];
    }
    return acc;
  }, []);

  const flagListLength = flagList.length;

  if (flagListLength === 0) {
    return null;
  }

  return (
    <div className="bg-blue-50 p-4 border-l-4 border-blue-400">
      <div className="flex">
        <div className="flex-shrink-0">
          <InformationCircleIcon
            className="h-5 w-5 text-blue-400"
            aria-hidden="true"
          />
        </div>
        <div className="ml-3 flex-1">
          <p className="text-sm text-blue-700">
            {f(messages.feature.flagIsPrerequisite, {
              length: flagListLength,
            })}
          </p>
          <div
            className="inline-flex space-x-1 cursor-pointer"
            onClick={() => setSeeMore(!isSeeMore)}
          >
            <span className="text-sm font-medium text-gray-700 hover:text-gray-600">
              {isSeeMore ? f(messages.close) : f(messages.seeMore)}
            </span>
            {isSeeMore ? (
              <ChevronUpIcon className="w-5 text-gray-700" />
            ) : (
              <ChevronDownIcon className="w-5 text-gray-700" />
            )}
          </div>
          {isSeeMore && (
            <div className="pl-4 mt-2 space-y-2 text-sm">
              <p className="text-gray-600">
                {f(messages.feature.flagIsPrerequisiteDescription, {
                  length: flagListLength,
                })}
              </p>
              <ul className="list-disc pl-4">
                {flagList.map((flag) => (
                  <li key={flag.id}>
                    <Link
                      className="link text-left"
                      to={`${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${flag.id}`}
                    >
                      <p className="truncate w-96">{flag.name}</p>
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
export interface PrerequisiteInputProps {
  feature: Feature.AsObject;
}

export const PrerequisiteInput: FC<PrerequisiteInputProps> = memo(
  ({ feature }) => {
    const { formatMessage: f } = useIntl();
    const dispatch = useDispatch<AppDispatch>();
    const editable = useIsEditable();
    const methods = useFormContext();
    const currentEnvironment = useCurrentEnvironment();

    const {
      control,
      formState: { errors },
    } = methods;
    const {
      fields: prerequisites,
      append: appendPrerequisite,
      remove,
      update,
    } = useFieldArray({
      control,
      name: 'prerequisites',
      keyName: 'key',
    });

    const isFeaturesLoading = useSelector<AppState, boolean>(
      (state) => state.features.loading,
      shallowEqual
    );

    const features = useSelector<AppState, Feature.AsObject[]>(
      (state) => selectAllFeatures(state.features),
      shallowEqual
    );

    const handleAddPrerequisite = useCallback(() => {
      if (prerequisites.length === 0) {
        dispatchListFeatures().then(() => {
          setTimeout(() => {
            appendPrerequisite({
              featureId: null,
              variationId: null,
            });
          });
        });
      } else {
        appendPrerequisite({
          featureId: null,
          variationId: null,
        });
      }
    }, [prerequisites]);

    const handleRemovePrerequisite = useCallback(
      (idx) => {
        remove(idx);
      },
      [remove]
    );

    const dispatchListFeatures = () => {
      return dispatch(
        listFeatures({
          environmentNamespace: currentEnvironment.id,
          pageSize: 99999,
          cursor: '',
          tags: [],
          searchKeyword: null,
          enabled: null,
          hasExperiment: null,
          maintainerId: null,
          archived: false,
          orderBy: ListFeaturesRequest.OrderBy.DEFAULT,
          orderDirection: ListFeaturesRequest.OrderDirection.ASC,
        })
      );
    };

    useEffect(() => {
      if (prerequisites.length > 0) {
        dispatchListFeatures();
      }
    }, []);

    const disableAddPrerequisite = useCallback(() => {
      if (prerequisites.length === 0) {
        return false;
      }
      return prerequisites.length === features.length - 1;
    }, [prerequisites, features]);

    return (
      <div className="">
        {prerequisites.length > 0 && (
          <div className="bg-white rounded-md p-3 border space-y-2">
            {prerequisites.map((p: any, prerequisitesIdx) => {
              const featureIdName = `prerequisites[${prerequisitesIdx}].featureId`;
              const variationIdName = `prerequisites[${prerequisitesIdx}].variationId`;

              const variationList = features.find(
                (f) => f.id === p.featureId
              )?.variationsList;

              const variationOptions = variationList?.map((v) => ({
                label: v.value,
                value: v.id,
              }));

              const featureFlagOptions = features
                .filter((f) => f.id !== feature.id)
                .filter(
                  (f) =>
                    !prerequisites.some(
                      (p2: any) =>
                        p2.featureId === f.id && p2.featureId !== p.featureId
                    )
                )
                .map((f) => {
                  return {
                    value: f.id,
                    label: f.name,
                    enabled: f.enabled,
                  };
                });

              return (
                <div key={p.key} className="flex space-x-2">
                  <Controller
                    name={featureIdName}
                    control={control}
                    render={({ field }) => {
                      return (
                        <SelectFeatureFlag
                          placeholder={f(messages.feature.selectFlag)}
                          options={featureFlagOptions}
                          className="w-full"
                          onChange={(e: OptionFeatureFlag) => {
                            if (field.value !== e.value) {
                              field.onChange(e.value);
                              update(prerequisitesIdx, {
                                ...p,
                                featureId: e.value,
                                variationId: null,
                              });
                            }
                          }}
                          value={featureFlagOptions.find(
                            (o) => o.value === field.value
                          )}
                        />
                      );
                    }}
                  />

                  <Controller
                    name={variationIdName}
                    control={control}
                    render={({ field }) => {
                      return (
                        <Select
                          placeholder={f(messages.feature.selectVariation)}
                          options={variationOptions}
                          className="w-full"
                          onChange={(e) => {
                            field.onChange(e.value);
                            update(prerequisitesIdx, {
                              ...p,
                              variationId: e.value,
                            });
                          }}
                          value={
                            variationOptions?.find(
                              (o) => o.value === p.variationId
                            ) ?? null
                          }
                        />
                      );
                    }}
                  />
                  {editable && (
                    <div className="flex items-center">
                      <button
                        type="button"
                        onClick={() =>
                          handleRemovePrerequisite(prerequisitesIdx)
                        }
                        className="minus-circle-icon"
                      >
                        <MinusCircleIcon aria-hidden="true" />
                      </button>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
        {editable && (
          <div className="pt-4 flex">
            <button
              type="button"
              className="btn-submit"
              onClick={handleAddPrerequisite}
              disabled={disableAddPrerequisite()}
            >
              {f(messages.feature.addPrerequisites)}
            </button>
          </div>
        )}
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
    swap,
  } = useFieldArray({
    control,
    name: 'rules',
    keyName: 'key',
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
      id: uuid(),
      strategy: {
        option: {
          value: feature.variationsList[0].id,
          label: createVariationLabel(feature.variationsList[0]),
        },
        rolloutStrategy: newRolloutStrategy,
      },
      clauses: [
        {
          id: uuid(),
          type: ClauseType.COMPARE,
          attribute: '',
          operator: Clause.Operator.EQUALS.toString(),
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
              key={r.id}
              className={classNames('bg-white p-3 rounded-md border')}
            >
              <div key={ruleIdx}>
                <div className="flex justify-end space-x-2"></div>
                <div className="flex mb-2">
                  <label className={classNames()}>{`${f(
                    messages.feature.rule
                  )} ${ruleIdx + 1}`}</label>
                  <div className="flex-grow" />
                  {editable && (
                    <div className="flex py-1 space-x-4">
                      <button
                        type="button"
                        onClick={() => handleRemoveRule(ruleIdx)}
                        className="text-gray-500 hover:text-gray-800"
                      >
                        <XIcon className="w-5 h-5" aria-hidden="true" />
                      </button>
                      {ruleIdx !== 0 && (
                        <button
                          type="button"
                          onClick={() => swap(ruleIdx, ruleIdx - 1)}
                          className="text-gray-500 hover:text-gray-800"
                          disabled={ruleIdx === 0}
                        >
                          <ArrowUpIcon width={18} />
                        </button>
                      )}
                      {ruleIdx !== rules.length - 1 && (
                        <button
                          type="button"
                          onClick={() => swap(ruleIdx, ruleIdx + 1)}
                          className="text-gray-500 hover:text-gray-800"
                          disabled={ruleIdx === rules.length - 1}
                        >
                          <ArrowDownIcon width={18} />
                        </button>
                      )}
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
    keyName: 'key',
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
            id: uuid(),
            type: type,
            attribute: '',
            operator: Clause.Operator.EQUALS.toString(),
            values: [],
          });
          break;
        }
        case ClauseType.SEGMENT: {
          update(idx, {
            id: uuid(),
            type: type,
            attribute: '',
            operator: Clause.Operator.SEGMENT.toString(),
            values: [segmentOptions[0]?.value],
          });
          dispatch(
            listSegments({
              environmentNamespace: currentEnvironment.id,
              cursor: '',
            })
          );
          break;
        }
        case ClauseType.DATE: {
          const now = String(Math.round(new Date().getTime() / 1000));
          update(idx, {
            id: uuid(),
            type: type,
            attribute: '',
            operator: Clause.Operator.BEFORE.toString(),
            values: [now],
          });
          break;
        }
      }
    },
    [update, dispatch, currentEnvironment, segmentOptions]
  );

  const handleAdd = useCallback(() => {
    append({
      id: uuid(),
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
        const clauseType = `${clauseName}.type`;
        const clauseAttribute = `${clauseName}.attribute`;
        const clauseOperator = `${clauseName}.operator`;
        const clauseValues = `${clauseName}.values`;

        return (
          <div key={c.id} className={classNames('flex space-x-2')}>
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
              name={clauseType}
              control={control}
              render={({ field }) => (
                <Select
                  onChange={(e) => {
                    if (e.value === field.value) {
                      return;
                    }
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
                    {...register(clauseAttribute)}
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
                  name={clauseOperator}
                  control={control}
                  render={({ field }) => (
                    <Select
                      onChange={(e) => {
                        field.onChange(e.value);
                      }}
                      options={clauseCompareOperatorOptions}
                      disabled={!editable}
                      value={clauseCompareOperatorOptions.find(
                        (o) => o.value === field.value
                      )}
                    />
                  )}
                />
                <div className="col-span-2">
                  <Controller
                    name={clauseValues}
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
                        name={clauseValues}
                        control={control}
                        render={({ field }) => {
                          return (
                            <Select
                              onChange={(o: Option) => {
                                field.onChange([o.value]);
                              }}
                              options={segmentOptions}
                              disabled={!editable}
                              value={segmentOptions.find(
                                (o) => o.value === field.value[0]
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
                    {...register(clauseAttribute)}
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
                  name={clauseOperator}
                  control={control}
                  render={({ field }) => (
                    <Select
                      onChange={(o: Option) => field.onChange(o.value)}
                      options={clauseDateOperatorOptions}
                      disabled={!editable}
                      value={clauseDateOperatorOptions.find(
                        (o) => o.value === field.value
                      )}
                      isSearchable={false}
                    />
                  )}
                />
                <div className="col-span-2">
                  <DatetimePicker name={clauseValues} />
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
        label: createVariationLabel(v),
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
          <div className="mt-2 space-y-2">
            {rolloutStrategy.map((s: any, idx: number) => {
              return (
                <div key={s.id} className="flex items-center space-x-2">
                  <div className="flex w-36 flex-shrink-0">
                    <input
                      {...register(`${rolloutStrategyName}.${idx}.percentage`)}
                      type="number"
                      min="0"
                      max="100"
                      defaultValue={s.percentage}
                      className={classNames(
                        'flex-grow pr-0 py-1 rounded-l border border-r-0 border-gray-300 w-full',
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
                  <label className="truncate text-sm text-gray-700">
                    {createVariationLabel(
                      feature.variationsList.find((v) => v.id == s.id)
                    )}
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
      render={({ field }) => {
        return (
          <ReactDatePicker
            dateFormat="yyyy-MM-dd HH:mm"
            showTimeSelect
            timeIntervals={60}
            placeholderText=""
            className={classNames('input-text w-full')}
            wrapperClassName="w-full"
            onChange={(v) => {
              const data = [v.getTime() / 1000];
              field.onChange(data);
            }}
            selected={(() => {
              return field.value[0]
                ? new Date(Number(field.value[0]) * 1000)
                : null;
            })()}
            disabled={!editable}
          />
        );
      }}
    />
  );
});
