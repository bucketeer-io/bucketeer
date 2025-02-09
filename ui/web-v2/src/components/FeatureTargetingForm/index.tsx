import {
  PAGE_PATH_FEATURES,
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_ROOT
} from '../../constants/routing';
import { selectAll as selectAllProgressiveRollouts } from '../../modules/porgressiveRollout';
import { ProgressiveRollout } from '../../proto/autoops/progressive_rollout_pb';
import { ListFeaturesRequest } from '../../proto/feature/service_pb';
import {
  createVariationLabel,
  getAlreadyTargetedVariation
} from '../../utils/variation';
import {
  MinusCircleIcon,
  XIcon,
  InformationCircleIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  PlusCircleIcon,
  ArrowUpIcon,
  ArrowDownIcon,
  CheckIcon
} from '@heroicons/react/solid';
import { PencilIcon, EyeIcon } from '@heroicons/react/outline';
import { FileCopyOutlined } from '@material-ui/icons';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useState, useEffect } from 'react';
import ReactDatePicker from 'react-datepicker';
import {
  useFormContext,
  Controller,
  useFieldArray,
  useWatch
} from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import 'react-datepicker/dist/react-datepicker.css';
import { Link, useHistory } from 'react-router-dom';
import { components } from 'react-select';
import ReactCreatableSelect from 'react-select/creatable';
import { v4 as uuid } from 'uuid';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectById as selectFeatureById,
  selectAll as selectAllFeatures,
  listFeatures
} from '../../modules/features';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import {
  listSegments,
  selectAll as selectAllSegments
} from '../../modules/segments';
import { Clause } from '../../proto/feature/clause_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { Strategy } from '../../proto/feature/strategy_pb';
import { AppDispatch } from '../../store';
import { classNames } from '../../utils/css';
import { isProgressiveRolloutsRunningWaiting } from '../ProgressiveRolloutAddForm';
import { CopyChip } from '../CopyChip';
import { colourStyles, CreatableSelect } from '../CreatableSelect';
import { Option, Select } from '../Select';
import { OptionFeatureFlag, SelectFeatureFlag } from '../SelectFeatureFlag';
import { Switch } from '../Switch';
import { TargetingForm, ruleClauseType } from '../../pages/feature/formSchema';
import { CancelScheduleDialog } from '../CancelScheduleDialog';
import { Dialog } from '@headlessui/react';
import ProjectSvg from '../../assets/svg/project.svg';
import { Overlay } from '../Overlay';

interface FeatureTargetingFormProps {
  featureId: string;
  onOpenConfirmDialog: () => void;
}

export const FeatureTargetingForm: FC<FeatureTargetingFormProps> = memo(
  ({ featureId, onOpenConfirmDialog }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();
    const methods = useFormContext<TargetingForm>();
    const {
      control,
      formState: { errors, dirtyFields },
      watch
    } = methods;
    const history = useHistory();
    const currentEnvironment = useCurrentEnvironment();
    const prerequisites = watch('prerequisites');

    const rules = watch('rules');
    const targets = watch('targets');

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
    const isProgressiveRolloutsRunning = !!progressiveRolloutList.find((p) =>
      isProgressiveRolloutsRunningWaiting(p.status)
    );

    const strategyOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: createVariationLabel(v)
      };
    });
    strategyOptions.push({
      value: Strategy.Type.ROLLOUT.toString(),
      label: f(messages.feature.strategy.selectRolloutPercentage)
    });
    const offVariationOptions: Option[] = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: createVariationLabel(v)
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
          if (clause.type === ruleClauseType.SEGMENT) {
            return clause.values.length > 0;
          }
          return clause.attribute && clause.values.length > 0;
        })
      );

      if (
        !checkPrerequisites ||
        !checkRules ||
        Object.values(errors).some(Boolean) ||
        Object.keys(dirtyFields).length === 0
      ) {
        return true;
      }
      return false;
    }, [rules, dirtyFields, errors, prerequisites]);

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

    const [open, setOpen] = useState(false);
    const [isSeeChangesModalOpen, setIsSeeChangesModalOpen] = useState(false);

    const disabled = false;

    return (
      <div className="p-10 bg-gray-100">
        <CancelScheduleDialog open={open} onClose={() => setOpen(false)} />
        <div className="bg-blue-50 p-4 border-l-4 border-blue-400 mb-7 inline-block">
          <div className="flex">
            <div className="flex-shrink-0">
              <InformationCircleIcon
                className="h-5 w-5 text-blue-400"
                aria-hidden="true"
              />
            </div>
            <div className="ml-3 flex-1 space-y-4">
              <p className="flex text-sm text-blue-700">
                Changes have been scheduled and will be applied on 01/21/2025 at
                10:00
              </p>
              <div className="flex items-center">
                <div className="flex items-center text-primary border-r border-gray-300 pr-4 space-x-2">
                  <CheckIcon className="w-5 h-5" />
                  <button className="text-sm font-normal">Apply Now</button>
                </div>
                <div className="flex items-center text-primary border-r border-gray-300 px-4 space-x-2">
                  <PencilIcon className="w-4 h-4" />
                  <button className="text-sm font-normal">Edit Schedule</button>
                </div>
                <div className="flex items-center text-primary border-r border-gray-300 px-4 space-x-2">
                  <EyeIcon className="w-5 h-5" />
                  <button
                    className="text-sm font-normal"
                    onClick={() => setIsSeeChangesModalOpen(true)}
                  >
                    See Changes
                  </button>
                </div>
                <div className="flex items-center text-red-600 px-4 space-x-2">
                  <XIcon className="w-5 h-5" />
                  <button
                    className="text-sm font-normal"
                    onClick={() => setOpen(true)}
                  >
                    Cancel Schedule
                  </button>
                </div>

                <button
                  className="text-sm text-red-600 font-medium p-2"
                  onClick={() => setOpen(true)}
                ></button>
              </div>
            </div>
          </div>
        </div>
        <Overlay
          open={isSeeChangesModalOpen}
          onClose={() => setIsSeeChangesModalOpen(false)}
        >
          <div className="w-[500px] h-full">
            <form className="flex flex-col h-full">
              <div className="h-full flex flex-col">
                <div className="flex items-center justify-between px-4 py-5 border-b">
                  <p className="text-xl font-medium">Changes Details</p>
                  <XIcon
                    width={20}
                    className="text-gray-400 cursor-pointer"
                    onClick={() => setIsSeeChangesModalOpen(false)}
                  />
                </div>
                <div className="flex-1 flex flex-col p-4">
                  <div className="space-y-4 pb-4 pt-2 border-b border-gray-200">
                    <div className="flex space-x-4 items-center">
                      <div className="bg-purple-50 p-3 rounded">
                        <ProjectSvg />
                      </div>
                      <p className="text-gray-500">
                        <strong>John Tucker</strong> updated the flag{' '}
                        <a href="#" className="text-primary underline">
                          Flag Name
                        </a>
                      </p>
                    </div>
                    <p className="text-gray-400">August 2, 2023 - 1:32PM</p>
                  </div>
                  <div className="py-4 space-y-4">
                    <p className="text-gray-400">PATCH</p>
                  </div>
                </div>
              </div>
            </form>
          </div>
        </Overlay>
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
                  readOnly={!editable || disabled}
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
                {targets.map((t, idx) => {
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
                          name={`targets.${idx}.users`}
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
                                                variationName
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
                                            userId: userInput
                                          })}
                                        </span>
                                      </div>
                                    );
                                  }}
                                  components={{
                                    DropdownIndicator: null,
                                    NoOptionsMessage: (props) => (
                                      <NoOptionsMessage props={props} />
                                    )
                                  }}
                                  value={field.value.map((u) => {
                                    return {
                                      value: u,
                                      label: u
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
                                  isDisabled={!editable || disabled}
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
                {isProgressiveRolloutsRunning && (
                  <div className="bg-blue-50 p-4 border-l-4 border-blue-400 mb-4 inline-block">
                    <div className="flex">
                      <div className="flex-shrink-0">
                        <InformationCircleIcon
                          className="h-5 w-5 text-blue-400"
                          aria-hidden="true"
                        />
                      </div>
                      <div className="ml-3 flex-1">
                        <p className="flex text-sm text-blue-700">
                          {f(
                            messages.feature.runningProgressiveRolloutMessage,
                            {
                              link: (
                                <span
                                  onClick={() => {
                                    history.push(
                                      `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_AUTOOPS}`
                                    );
                                  }}
                                  className="underline text-primary cursor-pointer ml-1"
                                >
                                  <span>
                                    {f(messages.sourceType.progressiveRollout)}
                                  </span>
                                </span>
                              )
                            }
                          )}
                        </p>
                      </div>
                    </div>
                  </div>
                )}
                <StrategyInput
                  feature={feature}
                  strategyName={'defaultStrategy'}
                  disabled={!!isProgressiveRolloutsRunning || disabled}
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
                      disabled={!editable || disabled}
                      value={offVariationOptions.find(
                        (o) => o.value === field.value.value
                      )}
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
      environmentId: currentEnvironment.id,
      pageSize: 0,
      cursor: '',
      tags: [],
      searchKeyword: null,
      maintainerId: null,
      orderBy: ListFeaturesRequest.OrderBy.DEFAULT,
      orderDirection: ListFeaturesRequest.OrderDirection.ASC
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
          name: feature.name
        }
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
              length: flagListLength
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
                  length: flagListLength
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
    const methods = useFormContext<TargetingForm>();
    const currentEnvironment = useCurrentEnvironment();

    const { control } = methods;
    const {
      fields: prerequisites,
      append: appendPrerequisite,
      remove,
      update
    } = useFieldArray({
      control,
      name: 'prerequisites',
      keyName: 'key'
    });

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
              variationId: null
            });
          });
        });
      } else {
        appendPrerequisite({
          featureId: null,
          variationId: null
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
          environmentId: currentEnvironment.id,
          pageSize: 0,
          cursor: '',
          tags: [],
          searchKeyword: null,
          enabled: null,
          hasExperiment: null,
          maintainerId: null,
          archived: false,
          orderBy: ListFeaturesRequest.OrderBy.DEFAULT,
          orderDirection: ListFeaturesRequest.OrderDirection.ASC
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

    const disabled = false;

    return (
      <div className="">
        {prerequisites.length > 0 && (
          <div className="bg-white rounded-md p-3 border space-y-2">
            {prerequisites.map((p, prerequisitesIdx) => {
              const variationList = features.find(
                (f) => f.id === p.featureId
              )?.variationsList;

              const variationOptions = variationList?.map((v) => ({
                label: v.value,
                value: v.id
              }));

              const featureFlagOptions = features
                .filter((f) => f.id !== feature.id)
                .filter(
                  (f) =>
                    !prerequisites.some(
                      (p2) =>
                        p2.featureId === f.id && p2.featureId !== p.featureId
                    )
                )
                .map((f) => {
                  return {
                    value: f.id,
                    label: f.name,
                    enabled: f.enabled
                  };
                });

              return (
                <div key={p.key} className="flex space-x-2">
                  <Controller
                    name={`prerequisites.${prerequisitesIdx}.featureId`}
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
                                variationId: null
                              });
                            }
                          }}
                          value={featureFlagOptions.find(
                            (o) => o.value === field.value
                          )}
                          disabled={disabled}
                        />
                      );
                    }}
                  />

                  <Controller
                    name={`prerequisites.${prerequisitesIdx}.variationId`}
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
                              variationId: e.value
                            });
                          }}
                          value={
                            variationOptions?.find(
                              (o) => o.value === p.variationId
                            ) ?? null
                          }
                          disabled={disabled}
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
                        disabled={disabled}
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
              disabled={disableAddPrerequisite() || disabled}
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
  const methods = useFormContext<TargetingForm>();
  const {
    control,
    formState: { errors }
  } = methods;
  const {
    fields: rules,
    append: appendRule,
    remove,
    swap
  } = useFieldArray({
    control,
    name: 'rules',
    keyName: 'key'
  });

  const newRolloutStrategy = [];
  feature.variationsList.forEach((val) => {
    newRolloutStrategy.push({
      id: val.id,
      percentage: 0
    });
  });
  const handleAddRule = useCallback(() => {
    appendRule({
      id: uuid(),
      strategy: {
        option: {
          value: feature.variationsList[0].id,
          label: createVariationLabel(feature.variationsList[0])
        },
        rolloutStrategy: newRolloutStrategy
      },
      clauses: [
        {
          id: uuid(),
          type: ruleClauseType.COMPARE,
          attribute: '',
          operator: Clause.Operator.EQUALS.toString(),
          values: []
        }
      ]
    });
  }, []);

  const handleRemoveRule = useCallback(
    (idx) => {
      remove(idx); // Remove the field
    },
    [remove]
  );

  const disabled = false;
  return (
    <div>
      <div className="grid grid-cols-1 gap-2">
        {rules.map((r, ruleIdx) => {
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
                        disabled={disabled}
                        className={classNames(
                          'text-gray-500',
                          disabled
                            ? 'cursor-not-allowed opacity-80'
                            : 'hover:text-gray-800'
                        )}
                      >
                        <XIcon className="w-5 h-5" aria-hidden="true" />
                      </button>
                      {ruleIdx !== 0 && (
                        <button
                          type="button"
                          onClick={() => swap(ruleIdx, ruleIdx - 1)}
                          className={classNames(
                            'text-gray-500',
                            disabled
                              ? 'cursor-not-allowed opacity-80'
                              : 'hover:text-gray-800'
                          )}
                          disabled={ruleIdx === 0 || disabled}
                        >
                          <ArrowUpIcon width={18} />
                        </button>
                      )}
                      {ruleIdx !== rules.length - 1 && (
                        <button
                          type="button"
                          onClick={() => swap(ruleIdx, ruleIdx + 1)}
                          className={classNames(
                            'text-gray-500',
                            disabled
                              ? 'cursor-not-allowed opacity-80'
                              : 'hover:text-gray-800'
                          )}
                          disabled={ruleIdx === rules.length - 1 || disabled}
                        >
                          <ArrowDownIcon width={18} />
                        </button>
                      )}
                    </div>
                  )}
                </div>
                <ClausesInput
                  featureId={feature.id}
                  ruleIdx={ruleIdx}
                  disabled={disabled}
                />
              </div>
              <StrategyInput
                feature={feature}
                strategyName={`rules.${ruleIdx}.strategy`}
                disabled={disabled}
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
          <button
            type="button"
            className="btn-submit"
            onClick={handleAddRule}
            disabled={disabled}
          >
            {f(messages.button.addRule)}
          </button>
        </div>
      )}
    </div>
  );
});

export const clauseTypeOptions: Option[] = [
  {
    value: ruleClauseType.COMPARE,
    label: intl.formatMessage(messages.feature.clause.type.compare)
  },
  {
    value: ruleClauseType.SEGMENT,
    label: intl.formatMessage(messages.feature.clause.type.segment)
  },
  {
    value: ruleClauseType.DATE,
    label: intl.formatMessage(messages.feature.clause.type.date)
  },
  {
    value: ruleClauseType.FEATURE_FLAG,
    label: intl.formatMessage(messages.feature.clause.type.featureFlag)
  }
];

export const clauseCompareOperatorOptions: Option[] = [
  {
    value: Clause.Operator.EQUALS.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.equal)
  },
  {
    value: Clause.Operator.GREATER_OR_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.greaterOrEqual)
  },
  {
    value: Clause.Operator.GREATER.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.greater)
  },
  {
    value: Clause.Operator.LESS_OR_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.lessOrEqual)
  },
  {
    value: Clause.Operator.LESS.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.less)
  },
  {
    value: Clause.Operator.IN.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.in)
  },
  {
    value: Clause.Operator.PARTIALLY_MATCH.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.partiallyMatch)
  },
  {
    value: Clause.Operator.STARTS_WITH.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.startWith)
  },
  {
    value: Clause.Operator.ENDS_WITH.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.endWith)
  }
];

export const clauseDateOperatorOptions: Option[] = [
  {
    value: Clause.Operator.BEFORE.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.before)
  },
  {
    value: Clause.Operator.AFTER.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.after)
  }
];

interface ClausesInputProps {
  featureId: string;
  ruleIdx: number;
  disabled: boolean;
}

const ClausesInput: FC<ClausesInputProps> = memo(
  ({ featureId, ruleIdx, disabled }) => {
    const { formatMessage: f } = useIntl();
    const dispatch = useDispatch<AppDispatch>();
    const editable = useIsEditable();
    const currentEnvironment = useCurrentEnvironment();
    const isSegmentLoading = useSelector<AppState, boolean>(
      (state) => state.segments.loading
    );
    const isFeaturesLoading = useSelector<AppState, boolean>(
      (state) => state.features.loading
    );
    const methods = useFormContext<TargetingForm>();
    const {
      getValues,
      register,
      control,
      formState: { errors }
    } = methods;
    const {
      fields: clauses,
      append,
      remove,
      update
    } = useFieldArray({
      control,
      name: `rules.${ruleIdx}.clauses`,
      keyName: 'key'
    });
    const selectedFeatureIds = new Set(clauses.map((c) => c.attribute));

    const segmentOptions = useSelector<AppState, Option[]>(
      (state) =>
        selectAllSegments(state.segments).map((s) => {
          return {
            value: s.id,
            label: s.name
          };
        }),
      shallowEqual
    );

    const [featureOptions, variationOptionsMap] = useSelector<
      AppState,
      [Option[], Map<string, Option[]>]
    >((state) => {
      const features = selectAllFeatures(state.features);
      const fos = features
        .filter((f) => f.id !== featureId)
        .map((f) => {
          return {
            value: f.id,
            label: f.name
          };
        });
      const vos = new Map<string, Option[]>();
      features.map((f) => {
        vos[f.id] = f.variationsList.map((v) => {
          return {
            value: v.id,
            label: v.value
          };
        });
      });
      return [fos, vos];
    }, shallowEqual);

    const handleChangeType = useCallback(
      (idx: number, type: string) => {
        switch (type) {
          case ruleClauseType.COMPARE: {
            update(idx, {
              id: uuid(),
              type: type,
              attribute: '',
              operator: Clause.Operator.EQUALS.toString(),
              values: []
            });
            break;
          }
          case ruleClauseType.SEGMENT: {
            update(idx, {
              id: uuid(),
              type: type,
              attribute: '',
              operator: Clause.Operator.SEGMENT.toString(),
              values: [segmentOptions[0]?.value]
            });
            dispatch(
              listSegments({
                environmentId: currentEnvironment.id,
                cursor: ''
              })
            );
            break;
          }
          case ruleClauseType.DATE: {
            const now = String(Math.round(new Date().getTime() / 1000));
            update(idx, {
              id: uuid(),
              type: type,
              attribute: '',
              operator: Clause.Operator.BEFORE.toString(),
              values: [now]
            });
            break;
          }
          case ruleClauseType.FEATURE_FLAG: {
            update(idx, {
              id: uuid(),
              type: type,
              attribute: '',
              operator: Clause.Operator.FEATURE_FLAG.toString(),
              values: []
            });
            dispatch(
              listFeatures({
                environmentId: currentEnvironment.id,
                pageSize: 0,
                cursor: '',
                tags: [],
                searchKeyword: null,
                maintainerId: null,
                enabled: null,
                hasExperiment: null,
                archived: false,
                orderBy: ListFeaturesRequest.OrderBy.DEFAULT,
                orderDirection: ListFeaturesRequest.OrderDirection.ASC
              })
            );
            break;
          }
        }
      },
      [update, dispatch, currentEnvironment, segmentOptions, featureOptions]
    );

    const handleAdd = useCallback(() => {
      append({
        id: uuid(),
        type: ruleClauseType.COMPARE,
        attribute: '',
        operator: Clause.Operator.EQUALS.toString(),
        values: []
      });
    }, [append]);

    const handleRemove = useCallback(
      (idx) => {
        remove(idx);
      },
      [remove]
    );

    useEffect(() => {
      dispatch(
        listFeatures({
          environmentId: currentEnvironment.id,
          pageSize: 0,
          cursor: '',
          tags: [],
          searchKeyword: null,
          maintainerId: null,
          enabled: null,
          hasExperiment: null,
          archived: false,
          orderBy: ListFeaturesRequest.OrderBy.DEFAULT,
          orderDirection: ListFeaturesRequest.OrderDirection.ASC
        })
      );
    }, []);

    return (
      <div className="grid grid-cols-1 gap-2">
        {clauses.map((c, clauseIdx) => {
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
                name={`rules.${ruleIdx}.clauses.${clauseIdx}.type`}
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
                    disabled={!editable || disabled}
                    isSearchable={false}
                    value={clauseTypeOptions.find((o) => o.value == c.type)}
                  />
                )}
              />
              {c.type == ruleClauseType.COMPARE && (
                <div className={classNames('flex-grow grid grid-cols-4 gap-1')}>
                  <div>
                    <input
                      {...register(
                        `rules.${ruleIdx}.clauses.${clauseIdx}.attribute`
                      )}
                      type="text"
                      defaultValue={c.attribute}
                      className={classNames('input-text w-full')}
                      disabled={!editable || disabled}
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
                    name={`rules.${ruleIdx}.clauses.${clauseIdx}.operator`}
                    control={control}
                    render={({ field }) => (
                      <Select
                        onChange={(e) => {
                          field.onChange(e.value);
                        }}
                        options={clauseCompareOperatorOptions}
                        disabled={!editable || disabled}
                        value={clauseCompareOperatorOptions.find(
                          (o) => o.value === field.value
                        )}
                      />
                    )}
                  />
                  <div className="col-span-2">
                    <Controller
                      name={`rules.${ruleIdx}.clauses.${clauseIdx}.values`}
                      control={control}
                      render={({ field }) => {
                        return (
                          <CreatableSelect
                            disabled={!editable || disabled}
                            defaultValues={field.value.map((v) => {
                              return {
                                value: v,
                                label: v
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
                            errors.rules?.[ruleIdx]?.clauses?.[clauseIdx]
                              ?.values?.message
                          }
                        </span>
                      )}
                    </p>
                  </div>
                </div>
              )}
              {c.type == ruleClauseType.SEGMENT &&
                (segmentOptions?.length > 0 ? (
                  <div
                    className={classNames('flex-grow grid grid-cols-2 gap-1')}
                  >
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
                          name={`rules.${ruleIdx}.clauses.${clauseIdx}.values`}
                          control={control}
                          render={({ field }) => {
                            return (
                              <Select
                                onChange={(o: Option) => {
                                  field.onChange([o.value]);
                                }}
                                options={segmentOptions}
                                disabled={!editable || disabled}
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
              {c.type == ruleClauseType.DATE && (
                <div className={classNames('flex-grow grid grid-cols-4 gap-1')}>
                  <div>
                    <input
                      {...register(
                        `rules.${ruleIdx}.clauses.${clauseIdx}.attribute`
                      )}
                      type="text"
                      defaultValue={c.attribute}
                      className={classNames('input-text w-full')}
                      disabled={!editable || disabled}
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
                    name={`rules.${ruleIdx}.clauses.${clauseIdx}.operator`}
                    control={control}
                    render={({ field }) => (
                      <Select
                        onChange={(o: Option) => field.onChange(o.value)}
                        options={clauseDateOperatorOptions}
                        disabled={!editable || disabled}
                        value={clauseDateOperatorOptions.find(
                          (o) => o.value === field.value
                        )}
                        isSearchable={false}
                      />
                    )}
                  />
                  <div className="col-span-2">
                    <DatetimePicker
                      name={`rules.${ruleIdx}.clauses.${clauseIdx}.values`}
                      disabled={disabled}
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
                </div>
              )}
              {c.type == ruleClauseType.FEATURE_FLAG &&
                (featureOptions?.length > 0 ? (
                  <div className={classNames('flex-grow')}>
                    {isFeaturesLoading ? (
                      <div>loading</div>
                    ) : (
                      <div
                        className={classNames(
                          'flex-grow grid grid-cols-3 gap-1'
                        )}
                      >
                        <Controller
                          name={`rules.${ruleIdx}.clauses.${clauseIdx}.attribute`}
                          control={control}
                          render={({ field }) => {
                            const clauseFeatureOptions = featureOptions.filter(
                              (o) =>
                                o.value === field.value ||
                                !selectedFeatureIds.has(o.value)
                            );
                            return (
                              <Select
                                onChange={(o: Option) => {
                                  field.onChange(o.value);
                                }}
                                options={clauseFeatureOptions}
                                disabled={!editable || disabled}
                                value={featureOptions.find(
                                  (o) => o.value === field.value
                                )}
                                isSearchable={true}
                              />
                            );
                          }}
                        />
                        <span className="inline-flex items-center text-sm text-gray-700 px-2">
                          {f(messages.feature.clause.operator.equal)}
                        </span>
                        <Controller
                          name={`rules.${ruleIdx}.clauses.${clauseIdx}.values`}
                          control={control}
                          render={({ field }) => {
                            const selectedfeatureId = getValues(
                              `rules.${ruleIdx}.clauses.${clauseIdx}.attribute`
                            );
                            const variationOptions =
                              variationOptionsMap[selectedfeatureId] ?? [];
                            return (
                              <Select
                                onChange={(o: Option) => {
                                  field.onChange([o.value]);
                                }}
                                options={variationOptions}
                                disabled={!editable || disabled}
                                value={variationOptions.find(
                                  (o) => o.value === field.value[0]
                                )}
                                isSearchable={false}
                              />
                            );
                          }}
                        />
                      </div>
                    )}
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
                ) : (
                  <div className="flex-grow flex content-center">
                    <span className="inline-flex items-center text-sm text-gray-700 px-2">
                      {f(messages.feature.noSelectableFeatureFlags)}
                    </span>
                  </div>
                ))}
              {editable && (
                <div className="flex items-center">
                  <button
                    type="button"
                    onClick={() => handleRemove(clauseIdx)}
                    className="minus-circle-icon"
                    disabled={clauses.length <= 1 || disabled}
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
            <button
              type="button"
              className="btn-submit"
              onClick={handleAdd}
              disabled={disabled}
            >
              {f(messages.button.addCondition)}
            </button>
          )}
        </div>
      </div>
    );
  }
);

export interface StrategyInputProps {
  feature: Feature.AsObject;
  strategyName: `rules.${number}.strategy` | `defaultStrategy`;
  disabled?: boolean;
}

export const StrategyInput: FC<StrategyInputProps> = memo(
  ({ feature, strategyName, disabled }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();
    const methods = useFormContext<TargetingForm>();
    const { register, control, trigger } = methods;
    const selectedOption = useWatch({
      control,
      name: `${strategyName}.option`
    });
    const { fields: rolloutStrategy, update } = useFieldArray({
      control,
      name: `${strategyName}.rolloutStrategy`,
      keyName: 'key' // the default keyName is "id" and it conflicts with the variation id field
    });

    const strategyOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: createVariationLabel(v)
      };
    });
    strategyOptions.push({
      value: String(Strategy.Type.ROLLOUT),
      label: f(messages.feature.strategy.selectRolloutPercentage)
    });
    const handleOnChange = useCallback(
      (idx: number, id: string, e: React.ChangeEvent<HTMLInputElement>) => {
        update(idx, {
          id: id,
          percentage: Number(e.target.value)
        });
        trigger(strategyName);
      },
      [update, trigger]
    );

    return (
      <div>
        <Controller
          name={`${strategyName}.option`}
          control={control}
          render={({ field }) => (
            <Select
              options={strategyOptions}
              disabled={!editable || disabled}
              value={{
                label: selectedOption.label ?? '',
                value: selectedOption.value ?? ''
              }}
              onChange={field.onChange}
              isSearchable={false}
            />
          )}
        />
        {selectedOption.value == Strategy.Type.ROLLOUT.toString() && (
          <div className="mt-2 space-y-2">
            {rolloutStrategy.map((s, idx: number) => {
              return (
                <div key={s.id} className="flex items-center space-x-2">
                  <div className="flex w-36 flex-shrink-0">
                    <input
                      {...register(
                        `${strategyName}.rolloutStrategy.${idx}.percentage`
                      )}
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
                      disabled={!editable || disabled}
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
                    .map((s) => Number(s.percentage))
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
  disabled?: boolean;
}

export const DatetimePicker: FC<DatetimePickerProps> = memo(
  ({ name, disabled }) => {
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
              disabled={!editable || disabled}
            />
          );
        }}
      />
    );
  }
);
