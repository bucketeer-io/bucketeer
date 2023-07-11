import {
  ID_NEW,
  PAGE_PATH_FEATURES,
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
} from '@/constants/routing';
import {
  deleteAutoOpsRule,
  selectAll as selectAllAutoOpsRules,
} from '@/modules/autoOpsRules';
import { AppDispatch } from '@/store';
import { Popover } from '@headlessui/react';
import {
  PlusIcon,
  DotsHorizontalIcon,
  PencilIcon,
  TrashIcon,
  InformationCircleIcon,
} from '@heroicons/react/outline';
import {
  XIcon,
  ArrowNarrowLeftIcon,
  ArrowNarrowRightIcon,
} from '@heroicons/react/solid';
import { SerializedError } from '@reduxjs/toolkit';
import dayjs from 'dayjs';
import React, { FC, memo, useCallback, useState } from 'react';
import {
  useFormContext,
  Controller,
  useFieldArray,
  useWatch,
} from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useSelector, shallowEqual, useDispatch } from 'react-redux';
import { useParams, useHistory } from 'react-router-dom';
import { v4 as uuid } from 'uuid';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { selectAll as selectAllGoals } from '../../modules/goals';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { AutoOpsRule, OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import {
  DatetimeClause,
  OpsEventRateClause,
} from '../../proto/autoops/clause_pb';
import { Goal } from '../../proto/experiment/goal_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { classNames } from '../../utils/css';
import { DatetimePicker } from '../DatetimePicker';
import { DetailSkeleton } from '../DetailSkeleton';
import { OperationAddUpdateForm } from '../OperationAddUpdateForm';
import { Overlay } from '../Overlay';
import { Option, Select } from '../Select';

enum TabLabel {
  ACTIVE = 'Active',
  COMPLETED = 'Completed',
}

export interface ClauseTypeMap {
  EVENT_RATE: 'bucketeer.autoops.OpsEventRateClause';
  DATETIME: 'bucketeer.autoops.DatetimeClause';
}

export const ClauseType: ClauseTypeMap = {
  EVENT_RATE: 'bucketeer.autoops.OpsEventRateClause',
  DATETIME: 'bucketeer.autoops.DatetimeClause',
};

interface FeatureAutoOpsRulesFormProps {
  featureId: string;
  refetchAutoOpsRules: () => void;
}

export const FeatureAutoOpsRulesForm: FC<FeatureAutoOpsRulesFormProps> = memo(
  ({ featureId, refetchAutoOpsRules }) => {
    const { operationId } = useParams<{ operationId: string }>();
    const isNew = operationId === ID_NEW;

    const [selectedAutoOpsRule, setSelectedAutoOpsRule] =
      useState<AutoOpsRule.AsObject | null>(null);

    const autoOpsRules = useSelector<AppState, AutoOpsRule.AsObject[]>(
      (state) =>
        selectAllAutoOpsRules(state.autoOpsRules).filter(
          (rule) => rule.featureId === featureId
        ),
      shallowEqual
    );

    const [open, setOpen] = useState(isNew);
    const [selectedPagination, setSelectedPagination] = useState(0);
    const history = useHistory();
    const currentEnvironment = useCurrentEnvironment();

    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();

    const methods = useFormContext();

    const { handleSubmit, reset } = methods;

    const [tabs, setTabs] = useState([
      {
        label: TabLabel.ACTIVE,
        value: autoOpsRules.filter((rule) => !rule.triggeredAt),
        selected: true,
      },
      {
        label: TabLabel.COMPLETED,
        value: autoOpsRules.filter((rule) => rule.triggeredAt),
        selected: false,
      },
    ]);

    const handleClose = useCallback(() => {
      reset();
      history.replace({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.id}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_AUTOOPS}`,
        search: location.search,
      });
      setOpen(false);
    }, [setOpen, history, location, reset]);

    const handleOpen = useCallback(() => {
      setOpen(true);
      history.push({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.id}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_AUTOOPS}${PAGE_PATH_NEW}`,
        search: location.search,
      });
    }, [setOpen, history, location]);

    const handleOpenUpdate = useCallback((rule: AutoOpsRule.AsObject) => {
      setSelectedAutoOpsRule(rule);
      handleOpen();
    }, []);

    const handleOnSubmit = useCallback(() => {
      handleClose();
      refetchAutoOpsRules();
    }, []);

    return (
      <div className="px-10 py-6 bg-white">
        <div className="flex justify-between">
          <div />
          <button onClick={handleOpen} className="btn-submit space-x-2">
            <PlusIcon width={18} />
            <span>New Operation</span>
          </button>
        </div>
        <div className="flex border-b border-gray-200 mt-2">
          {tabs.map((tab) => (
            <div
              key={tab.label}
              className={classNames(
                'px-4 py-3 cursor-pointer',
                tab.selected
                  ? 'text-primary border-b-2 border-primary'
                  : 'text-gray-400'
              )}
              onClick={() =>
                setTabs(
                  tabs.map((t) => ({
                    ...t,
                    selected: t.label === tab.label,
                  }))
                )
              }
            >
              {tab.label} ({tab.value.length})
            </div>
          ))}
        </div>
        <div className="space-y-6 py-6">
          {tabs
            .find((tab) => tab.selected)
            .value.map((rule) => (
              <Operation
                key={rule.id}
                rule={rule}
                isActiveSelected={
                  tabs.find((tab) => tab.selected).label === TabLabel.ACTIVE
                }
                handleOpenUpdate={handleOpenUpdate}
                refetchAutoOpsRules={refetchAutoOpsRules}
              />
            ))}
          {/* <div className="rounded-xl shadow px-6 py-4 bg-white">
            <div className="flex justify-between py-4 border-b">
              <h3 className="font-bold text-xl">Enable Operation</h3>
              <div className="flex space-x-2 items-center">
                <div className="py-[2px] px-2 bg-[#FFF7EE] rounded text-[#CE844A] text-sm">
                  Progressive Rollout
                </div>
                <Popover className="relative flex">
                  <Popover.Button>
                    <div className="pl-2 flex items-center cursor-pointer">
                      <DotsHorizontalIcon width={20} />
                    </div>
                  </Popover.Button>
                  <Popover.Panel className="absolute z-10 bg-white right-0 rounded-lg p-1 w-[166px]">
                    <button className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100">
                      <PencilIcon width={18} />
                      <span className="text-sm">Edit Operation</span>
                    </button>
                    <button className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100">
                      <TrashIcon width={18} className="text-red-500" />
                      <span className="text-red-500 text-sm">
                        Delete Operation
                      </span>
                    </button>
                  </Popover.Panel>
                </Popover>
              </div>
            </div>
            <div className="mt-4">
              <p className="font-bold text-lg">Progress Information</p>
              <div className="text-gray-400 flex items-center py-2 space-x-2">
                <div className="space-x-1 items-center flex">
                  <span className="">Increment 10%</span>
                  <InformationCircleIcon width={16} />
                </div>
                <span className="text-gray-200">/</span>
                <span>Start Date 2023-02-11</span>
                <span className="text-gray-200">/</span>
                <span>Frequency Hour</span>
              </div>
              <div className="mt-2">
                <div className="flex">
                  {Array(50)
                    .fill('')
                    .map((_, i) => {
                      const value = 54;
                      const percentage = i * 2;

                      let bgColor = 'bg-gray-200';
                      if (percentage <= value) {
                        bgColor = 'bg-pink-500';
                      } else if (percentage > 90) {
                        bgColor = 'bg-white';
                      } else if (percentage % 10 === 0) {
                        bgColor = 'bg-gray-400';
                      }

                      return (
                        <div
                          key={i}
                          className={classNames(
                            'relative h-[8px] flex-1 rounded-[60px]',
                            bgColor
                          )}
                        >
                          {i !== 0 && (
                            <div className="absolute h-[8px] w-1.5 rounded-r-full bg-white" />
                          )}
                        </div>
                      );
                    })}
                </div>
                <div className="flex mt-2">
                  {Array(10)
                    .fill('')
                    .map((_, i) => (
                      <div key={i} className="flex-1">
                        <p>{i * 10}%</p>
                        <p className="text-gray-400 text-xs">07:00</p>
                        <p className="text-gray-400 text-xs">2023-23-11</p>
                      </div>
                    ))}
                </div>
              </div>
              <div className="mt-4 flex justify-between items-center">
                <button
                  className="p-1.5 rounded border"
                  onClick={() =>
                    selectedPagination > 0 &&
                    setSelectedPagination(selectedPagination - 1)
                  }
                >
                  <ArrowNarrowLeftIcon width={16} className="text-gray-400" />
                </button>
                <div className="flex space-x-2">
                  {Array(8)
                    .fill('')
                    .map((_, i) =>
                      selectedPagination === i ? (
                        <div
                          key={i}
                          className="w-[24px] h-[8px] rounded-full bg-gray-400"
                        />
                      ) : (
                        <div
                          key={i}
                          className="w-[8px] h-[8px] rounded-full bg-gray-200"
                        />
                      )
                    )}
                </div>
                <button
                  className="p-1.5 rounded border"
                  onClick={() =>
                    selectedPagination < 8 &&
                    setSelectedPagination(selectedPagination + 1)
                  }
                >
                  <ArrowNarrowRightIcon width={16} className="text-gray-400" />
                </button>
              </div>
            </div>
          </div> */}
        </div>
        <Overlay open={open} onClose={handleClose}>
          {isNew && (
            <OperationAddUpdateForm
              onSubmit={handleSubmit(handleOnSubmit)}
              onCancel={handleClose}
              featureId={featureId}
              autoOpsRule={selectedAutoOpsRule}
            />
          )}
        </Overlay>
        <form className="">
          <div className="grid grid-cols-1 gap-y-6 gap-x-4">
            {/* <AutoOpsRulesInput featureId={featureId} /> */}
          </div>
          {/* {editable && (
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
          )} */}
        </form>
      </div>
    );
  }
);

interface OperationProps {
  rule: AutoOpsRule.AsObject;
  isActiveSelected: any;
  handleOpenUpdate: any;
  refetchAutoOpsRules: () => void;
}

const Operation = ({
  rule,
  isActiveSelected,
  handleOpenUpdate,
  refetchAutoOpsRules,
}: OperationProps) => {
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();

  const { typeUrl, value } = rule.clausesList[0].clause;
  const type = typeUrl.substring(typeUrl.lastIndexOf('/') + 1);

  let datetime;
  if (type === ClauseType.DATETIME) {
    const datetimeClause = DatetimeClause.deserializeBinary(
      value as Uint8Array
    ).toObject();
    datetime = dayjs(new Date(datetimeClause.time * 1000)).format(
      'YYYY-MM-DD HH:mm'
    );
  }

  let goal, minCount, threadsholdRate, operator;
  if (type === ClauseType.EVENT_RATE) {
    const opsEventRateClause = OpsEventRateClause.deserializeBinary(
      value as Uint8Array
    ).toObject();

    goal = opsEventRateClause.goalId;
    minCount = opsEventRateClause.minCount;
    threadsholdRate = opsEventRateClause.threadsholdRate * 100;
    // operator = opsEventRateClause.operator.toString();
  }

  const handleDelete = (ruleId) => {
    dispatch(
      deleteAutoOpsRule({
        environmentNamespace: currentEnvironment.namespace,
        id: ruleId,
      })
    ).then(refetchAutoOpsRules);
  };

  return (
    <div className="rounded-xl shadow px-6 py-4 bg-white">
      <div className="flex justify-between py-4 border-b">
        <h3 className="font-bold text-xl text-gray-600">
          {rule.opsType === OpsType.ENABLE_FEATURE
            ? 'Enable Operation'
            : 'Kill Switch Operation'}
        </h3>
        <div className="flex space-x-2 items-center">
          <div
            className={classNames(
              'py-[2px] px-2 rounded text-sm',
              type === ClauseType.DATETIME && 'bg-[#EBF9ED] text-green-700',
              type === ClauseType.EVENT_RATE && 'bg-[#EFECF5] text-primary'
            )}
          >
            {type === ClauseType.DATETIME && 'Schedule'}
            {type === ClauseType.EVENT_RATE && 'Event Rate'}
          </div>
          <Popover className="relative flex">
            <Popover.Button>
              <div className="pl-2 flex items-center cursor-pointer">
                <DotsHorizontalIcon width={20} />
              </div>
            </Popover.Button>
            <Popover.Panel className="absolute z-10 bg-white right-0 rounded-lg p-1 w-[166px] shadow-md">
              <button
                onClick={() => handleOpenUpdate(rule)}
                className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
              >
                <PencilIcon width={18} />
                <span className="text-sm">Edit Operation</span>
              </button>
              <button
                onClick={() => handleDelete(rule.id)}
                className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100"
              >
                <TrashIcon width={18} className="text-red-500" />
                <span className="text-red-500 text-sm">Delete Operation</span>
              </button>
            </Popover.Panel>
          </Popover>
        </div>
      </div>
      <div className="mt-4">
        <p className="font-bold text-lg text-gray-600">Progress Information</p>
        {type === ClauseType.DATETIME && (
          <>
            <div
              className={classNames(
                'mt-6 h-2  flex justify-between relative mx-1',
                isActiveSelected ? 'bg-gray-200' : 'bg-pink-500'
              )}
            >
              <div className="w-[14px] h-[14px] absolute top-1/2 -translate-y-1/2 rounded-full -left-1 bg-pink-500 border border-pink-100" />
              <div
                className={classNames(
                  'w-[14px] h-[14px] absolute top-1/2 -translate-y-1/2 rounded-full -right-1 bg-gray-300 border',
                  isActiveSelected ? 'bg-gray-200' : 'bg-pink-500'
                )}
              />
            </div>
            <div className="flex justify-between mt-2">
              <span>Off</span>
              <span>On</span>
            </div>
            <div className="flex justify-between">
              <span className="text-xs text-gray-400">Today</span>
              <span className="text-xs text-gray-400">{datetime}</span>
            </div>
          </>
        )}
        {type === ClauseType.EVENT_RATE && (
          <>
            <div className="flex items-center space-x-2 mt-3">
              <span className="text-gray-400">Goal</span>
              <span className="text-gray-500">{goal}</span>
              <span className="text-gray-200">/</span>
              <span className="text-gray-400">Min Count</span>
              <span className="text-gray-500">{minCount}</span>
              <span className="text-gray-200">/</span>
              <span className="text-gray-400">Current Goal Count</span>
              <span className="text-gray-500">
                {threadsholdRate}/{minCount} (
                {Math.round((threadsholdRate / minCount) * 100) / 100}%)
              </span>
              <InformationCircleIcon width={18} />
            </div>
            <div className="mt-3">
              <div className="flex">
                {Array(50)
                  .fill('')
                  .map((_, i) => {
                    const value = 54;
                    const percentage = i * 2;

                    let bgColor = 'bg-gray-200';
                    if (percentage <= value) {
                      bgColor = 'bg-pink-500';
                    } else if (percentage > 90) {
                      bgColor = 'bg-white';
                    } else if (percentage % 10 === 0) {
                      bgColor = 'bg-gray-400';
                    }

                    return (
                      <div
                        key={i}
                        className={classNames(
                          'relative h-[8px] flex-1 rounded-[60px]',
                          bgColor
                        )}
                      >
                        {i !== 0 && (
                          <div className="absolute h-[8px] w-1.5 rounded-r-full bg-white" />
                        )}
                      </div>
                    );
                  })}
              </div>
              <div className="flex mt-2">
                {Array(10)
                  .fill('')
                  .map((_, i) => (
                    <div key={i} className="flex-1">
                      {i * 10}%
                    </div>
                  ))}
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

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
  };
};

// export interface AutoOpsRulesInputProps {
//   featureId: string;
// }

// export const AutoOpsRulesInput: FC<AutoOpsRulesInputProps> = memo(
//   ({ featureId }) => {
//     const editable = useIsEditable();
//     const { formatMessage: f } = useIntl();
//     const [feature, getFeatureError] = useSelector<
//       AppState,
//       [Feature.AsObject | undefined, SerializedError | null]
//     >((state) => [
//       selectFeatureById(state.features, featureId),
//       state.features.getFeatureError,
//     ]);
//     const methods = useFormContext();
//     const {
//       control,
//       formState: { errors },
//     } = methods;
//     const {
//       fields: rules,
//       append,
//       remove,
//     } = useFieldArray({
//       control,
//       name: 'autoOpsRules',
//     });

//     const handleAdd = useCallback(() => {
//       append(createInitialAutoOpsRule(feature));
//     }, [append]);

//     const handleRemove = useCallback(
//       (idx) => {
//         remove(idx);
//       },
//       [remove]
//     );

//     return (
//       <div>
//         <div className="grid grid-cols-1 gap-2">
//           {rules.map((rule: any, ruleIdx) => {
//             return (
//               <div
//                 key={rule.id}
//                 className={classNames('bg-white p-3 rounded-md border, mb-5')}
//               >
//                 <div className="flex text-gray-700 pb-3">
//                   <div>
//                     <label className={classNames('text-sm')}>{`${f(
//                       messages.autoOps.operation
//                     )} ${ruleIdx + 1}`}</label>
//                   </div>
//                   <div className="flex-grow" />
//                   {editable && (
//                     <div className="flex items-center">
//                       <button
//                         type="button"
//                         className="x-icon"
//                         onClick={() => handleRemove(ruleIdx)}
//                       >
//                         <XIcon className="w-5 h-5" aria-hidden="true" />
//                       </button>
//                     </div>
//                   )}
//                 </div>
//                 <AutoOpsRuleInput
//                   featureId={featureId}
//                   ruleIdx={ruleIdx}
//                   feature={feature}
//                 />
//               </div>
//             );
//           })}
//         </div>
//         {editable && (
//           <div className="flex">
//             <button type="button" className="btn-submit" onClick={handleAdd}>
//               {f(messages.button.addOperation)}
//             </button>
//           </div>
//         )}
//       </div>
//     );
//   }
// );

// export interface AutoOpsRuleInputProps {
//   featureId: string;
//   ruleIdx: number;
//   feature: Feature.AsObject;
// }

// export const AutoOpsRuleInput: FC<AutoOpsRuleInputProps> = memo(
//   ({ featureId, ruleIdx, feature }) => {
//     const editable = useIsEditable();
//     const ruleName = `autoOpsRules.${ruleIdx}`;
//     const methods = useFormContext();
//     const { control, setValue } = methods;
//     const rule = useWatch({
//       control,
//       name: ruleName,
//     });

//     return (
//       <>
//         <Controller
//           name={`${ruleName}.opsType`}
//           control={control}
//           render={({ field }) => (
//             <Select
//               onChange={(o: Option) => {
//                 setValue(`autoOpsRules.${ruleIdx}.clauses`, [
//                   createInitialClause(feature),
//                 ]);
//                 field.onChange(o.value);
//               }}
//               options={opsTypeOptions}
//               disabled={!editable}
//               value={opsTypeOptions.find((o) => o.value == rule.opsType)}
//             />
//           )}
//         />
//         <ClausesInput featureId={featureId} ruleIdx={ruleIdx} />
//       </>
//     );
//   }
// );

// export interface ClausesInputProps {
//   featureId: string;
//   ruleIdx: number;
// }

// export const ClausesInput: FC<ClausesInputProps> = ({ featureId, ruleIdx }) => {
//   const editable = useIsEditable();
//   const [feature, getFeatureError] = useSelector<
//     AppState,
//     [Feature.AsObject | undefined, SerializedError | null]
//   >((state) => [
//     selectFeatureById(state.features, featureId),
//     state.features.getFeatureError,
//   ]);
//   const { formatMessage: f } = useIntl();
//   const methods = useFormContext();
//   const {
//     control,
//     formState: { errors },
//     watch,
//   } = methods;
//   const clausesName = `autoOpsRules.${ruleIdx}.clauses`;
//   const watchClauses = useWatch({
//     control,
//     name: clausesName,
//   });
//   const {
//     fields: clauses,
//     append,
//     remove,
//   } = useFieldArray({
//     control,
//     name: clausesName,
//   });

//   const handleAdd = useCallback(() => {
//     append(createInitialClause(feature));
//   }, [append]);

//   const handleRemove = useCallback(
//     (idx) => {
//       remove(idx);
//     },
//     [remove]
//   );

//   return (
//     <div>
//       {clauses.map((clause: any, clauseIdx: number) => {
//         return (
//           <div key={clause.id}>
//             <div className={classNames('flex space-x-2')}>
//               <div className="w-[2rem] flex justify-center items-center">
//                 {clauseIdx === 0 ? (
//                   <div
//                     className={classNames(
//                       'py-1 px-2',
//                       'text-xs bg-gray-400 text-white rounded-full'
//                     )}
//                   >
//                     IF
//                   </div>
//                 ) : (
//                   <div className="p-1 text-xs">OR</div>
//                 )}
//               </div>
//               <div className="flex-grow flex mt-3 p-3 rounded-md border">
//                 <div className="flex-grow">
//                   <ClauseInput
//                     featureId={featureId}
//                     ruleIdx={ruleIdx}
//                     clauseIdx={clauseIdx}
//                   />
//                 </div>
//                 {editable && (
//                   <div className="flex items-start pl-2">
//                     <button
//                       type="button"
//                       className="x-icon"
//                       onClick={() => handleRemove(clauseIdx)}
//                     >
//                       <XIcon className="w-5 h-5" aria-hidden="true" />
//                     </button>
//                   </div>
//                 )}
//               </div>
//             </div>
//           </div>
//         );
//       })}
//       {editable && !containsDatetimeClause(watchClauses) && (
//         <div className="py-4 flex">
//           <button type="button" className="btn-submit" onClick={handleAdd}>
//             {f(messages.button.addCondition)}
//           </button>
//         </div>
//       )}
//     </div>
//   );
// };

// function containsDatetimeClause(clauses): boolean {
//   for (const clause of clauses) {
//     if (clause.clauseType === ClauseType.DATETIME.toString()) {
//       return true;
//     }
//   }
//   return false;
// }

export const clauseTypeOptionEventRate = {
  value: ClauseType.EVENT_RATE.toString(),
  label: intl.formatMessage(messages.autoOps.eventRateClauseType),
};

export const clauseTypeOptionDatetime = {
  value: ClauseType.DATETIME.toString(),
  label: intl.formatMessage(messages.autoOps.datetimeClauseType),
};

export const clauseTypeOptions = [
  clauseTypeOptionEventRate,
  clauseTypeOptionDatetime,
];

export const createClauseTypeOption = (
  clauseType: ClauseTypeMap[keyof ClauseTypeMap]
) => {
  return clauseTypeOptions.find(
    (option) => clauseType.toString() == option.value
  );
};

// export interface ClauseInputProps {
//   featureId: string;
//   ruleIdx: number;
//   clauseIdx: number;
// }

// export const ClauseInput: FC<ClauseInputProps> = ({
//   featureId,
//   ruleIdx,
//   clauseIdx,
// }) => {
//   const editable = useIsEditable();
//   const methods = useFormContext();
//   const {
//     control,
//     formState: { errors },
//   } = methods;
//   const ruleName = `autoOpsRules.${ruleIdx}`;
//   const clauseName = `${ruleName}.clauses.${clauseIdx}`;
//   const opsType = useWatch({
//     control,
//     name: `${ruleName}.opsType`,
//   });
//   const selectedClauseTypeOptions =
//     opsType === OpsType.ENABLE_FEATURE.toString()
//       ? [clauseTypeOptionDatetime]
//       : [clauseTypeOptionEventRate, clauseTypeOptionDatetime];
//   const clauseType = useWatch({
//     control,
//     name: `${clauseName}.clauseType`,
//   });

//   return (
//     <div className="grid grid-cols-1 gap-2">
//       <div className="">
//         <Controller
//           name={`${clauseName}.clauseType`}
//           control={control}
//           render={({ field }) => (
//             <Select
//               onChange={(o: Option) => field.onChange(o.value)}
//               options={selectedClauseTypeOptions}
//               disabled={!editable}
//               value={selectedClauseTypeOptions.find(
//                 (o) => o.value === clauseType
//               )}
//             />
//           )}
//         />
//         {clauseType === ClauseType.EVENT_RATE.toString() && (
//           <EventRateClauseInput
//             featureId={featureId}
//             ruleIdx={ruleIdx}
//             clauseIdx={clauseIdx}
//           />
//         )}
//         {clauseType === ClauseType.DATETIME.toString() && (
//           <DatetimeClauseInput ruleIdx={ruleIdx} clauseIdx={clauseIdx} />
//         )}
//       </div>
//     </div>
//   );
// };

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

export const createOperatorOption = (
  operator: OpsEventRateClause.OperatorMap[keyof OpsEventRateClause.OperatorMap]
) => {
  return operatorOptions.find((option) => option.value === operator.toString());
};

// export interface EventRateClauseInputProps {
//   featureId: string;
//   ruleIdx: number;
//   clauseIdx: number;
// }

// export const EventRateClauseInput: FC<EventRateClauseInputProps> = memo(
//   ({ featureId, ruleIdx, clauseIdx }) => {
//     const editable = useIsEditable();
//     const { formatMessage: f } = useIntl();
//     const isGoalLoading = useSelector<AppState, boolean>(
//       (state) => state.goals.loading,
//       shallowEqual
//     );
//     const goals = useSelector<AppState, Goal.AsObject[]>(
//       (state) => selectAllGoals(state.goals),
//       shallowEqual
//     );
//     const goalOptions = goals.map((goal) => {
//       return {
//         value: goal.id,
//         label: goal.id,
//       };
//     });
//     const methods = useFormContext();
//     const opsEventRateClauseName = `autoOpsRules.${ruleIdx}.clauses.${clauseIdx}.opsEventRateClause`;
//     const {
//       register,
//       control,
//       formState: { errors },
//       trigger,
//     } = methods;
//     const [feature, _] = useSelector<
//       AppState,
//       [Feature.AsObject | undefined, SerializedError | null]
//     >((state) => [
//       selectFeatureById(state.features, featureId),
//       state.features.getFeatureError,
//     ]);
//     const clause = useWatch({
//       control,
//       name: opsEventRateClauseName,
//     });
//     const variationOptions = feature.variationsList.map((v) => {
//       return {
//         value: v.id,
//         label: v.value,
//       };
//     });

//     if (isGoalLoading) {
//       return (
//         <div className="p-9 bg-gray-100">
//           <DetailSkeleton />
//         </div>
//       );
//     }
//     return (
//       <div className="grid grid-cols-1 gap-2">
//         <div>
//           <label htmlFor="variation" className="input-label">
//             {f(messages.feature.variation)}
//           </label>
//           <Controller
//             name={`${opsEventRateClauseName}.variation`}
//             control={control}
//             render={({ field }) => (
//               <Select
//                 onChange={(o: Option) => field.onChange(o.value)}
//                 options={variationOptions}
//                 disabled={!editable}
//                 value={variationOptions.find((o) => o.value === field.value)}
//               />
//             )}
//           />
//         </div>
//         <label htmlFor="variation" className="input-label">
//           {f(messages.autoOps.opsEventRateClause.goal)}
//         </label>
//         <div className={classNames('flex-grow grid grid-cols-4 gap-1')}>
//           <Controller
//             name={`${opsEventRateClauseName}.goal`}
//             control={control}
//             render={({ field }) => (
//               <Select
//                 onChange={(o: Option) => field.onChange(o.value)}
//                 options={goalOptions}
//                 disabled={!editable}
//                 value={goalOptions.find((o) => o.value === clause.goal)}
//               />
//             )}
//           />
//           <Controller
//             name={`${opsEventRateClauseName}.operator`}
//             control={control}
//             render={({ field }) => (
//               <Select
//                 onChange={(o: Option) => field.onChange(o.value)}
//                 options={operatorOptions}
//                 disabled={!editable}
//                 value={operatorOptions.find((o) => o.value === clause.operator)}
//               />
//             )}
//           />
//           <div className="w-36 flex">
//             <input
//               {...register(`${opsEventRateClauseName}.threadsholdRate`)}
//               type="number"
//               min="0"
//               max="100"
//               defaultValue={clause.threadsholdRate}
//               className={classNames(
//                 'flex-grow pr-0 py-1',
//                 'rounded-l border border-r-0 border-gray-300',
//                 'text-right'
//               )}
//               placeholder={''}
//               required
//               disabled={!editable}
//             />
//             <span
//               className={classNames(
//                 'px-1 py-1 inline-flex items-center bg-gray-100',
//                 'rounded-r border border-l-0 border-gray-300 text-gray-600'
//               )}
//             >
//               {'%'}
//             </span>
//           </div>
//         </div>
//         <div>
//           <p className="input-error">
//             {errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
//               ?.opsEventRateClause?.threadsholdRate?.message && (
//               <span role="alert">
//                 {
//                   errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
//                     ?.opsEventRateClause?.threadsholdRate?.message
//                 }
//               </span>
//             )}
//           </p>
//         </div>
//         <div className="w-36">
//           <label htmlFor="name">
//             <span className="input-label">
//               {f(messages.autoOps.opsEventRateClause.minCount)}
//             </span>
//           </label>
//           <div className="mt-1">
//             <input
//               {...register(`${opsEventRateClauseName}.minCount`)}
//               type="number"
//               min="0"
//               className="input-text w-full"
//               disabled={!editable}
//             />
//             <p className="input-error">
//               {errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
//                 ?.opsEventRateClause?.minCount?.message && (
//                 <span role="alert">
//                   {
//                     errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
//                       ?.opsEventRateClause?.minCount?.message
//                   }
//                 </span>
//               )}
//             </p>
//           </div>
//         </div>
//       </div>
//     );
//   }
// );

// export interface DatetimeClauseInputProps {
//   ruleIdx: number;
//   clauseIdx: number;
// }

// export const DatetimeClauseInput: FC<DatetimeClauseInputProps> = memo(
//   ({ ruleIdx, clauseIdx }) => {
//     const editable = useIsEditable();
//     const { formatMessage: f } = useIntl();
//     const methods = useFormContext();
//     const clauseName = `autoOpsRules.${ruleIdx}.clauses.${clauseIdx}`;
//     const {
//       formState: { errors },
//     } = methods;

//     return (
//       <div className="">
//         <label htmlFor="name">
//           <span className="input-label">
//             {f(messages.autoOps.datetimeClause.datetime)}
//           </span>
//         </label>
//         <DatetimePicker
//           name={`${clauseName}.datetimeClause.time`}
//           disabled={!editable}
//         />
//         <p className="input-error">
//           {errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]?.datetimeClause
//             ?.time?.message && (
//             <span role="alert">
//               {
//                 errors.autoOpsRules?.[ruleIdx]?.clauses?.[clauseIdx]
//                   ?.datetimeClause?.time?.message
//               }
//             </span>
//           )}
//         </p>
//       </div>
//     );
//   }
// );
