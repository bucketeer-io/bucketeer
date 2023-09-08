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
import {
  listOpsCounts,
  selectAll as selectAllOpsCounts,
} from '@/modules/opsCounts';
import { OpsCount } from '@/proto/autoops/ops_count_pb';
import { AppDispatch } from '@/store';
import { Popover } from '@headlessui/react';
import {
  PlusIcon,
  DotsHorizontalIcon,
  PencilIcon,
  TrashIcon,
  InformationCircleIcon,
} from '@heroicons/react/outline';
import dayjs from 'dayjs';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useFormContext } from 'react-hook-form';
import { useSelector, shallowEqual, useDispatch } from 'react-redux';
import { useParams, useHistory } from 'react-router-dom';
import { v4 as uuid } from 'uuid';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { useCurrentEnvironment } from '../../modules/me';
import { AutoOpsRule, OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import {
  DatetimeClause,
  OpsEventRateClause,
} from '../../proto/autoops/clause_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { classNames } from '../../utils/css';
import { HoverPopover } from '../HoverPopover';
import { OperationAddUpdateForm } from '../OperationAddUpdateForm';
import { Overlay } from '../Overlay';

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
    const dispatch = useDispatch<AppDispatch>();

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
    const history = useHistory();
    const currentEnvironment = useCurrentEnvironment();

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

    useEffect(() => {
      if (autoOpsRules?.length > 0) {
        const ids = autoOpsRules
          .filter((rule) => {
            const { typeUrl } = rule.clausesList[0].clause;
            const type = typeUrl.substring(typeUrl.lastIndexOf('/') + 1);
            return type === ClauseType.EVENT_RATE && !rule.triggeredAt;
          })
          .map((rule) => rule.id);

        if (ids.length > 0) {
          dispatch(
            listOpsCounts({
              environmentNamespace: currentEnvironment.namespace,
              ids,
            })
          );
        }
      }
    }, [autoOpsRules]);

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
          <button
            onClick={() => {
              setSelectedAutoOpsRule(null);
              handleOpen();
            }}
            className="btn-submit space-x-2"
          >
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
        </div>
        {open && (
          <Overlay open={open} onClose={handleClose}>
            <OperationAddUpdateForm
              onSubmit={handleSubmit(handleOnSubmit)}
              onCancel={handleClose}
              featureId={featureId}
              autoOpsRule={selectedAutoOpsRule}
            />
          </Overlay>
        )}
      </div>
    );
  }
);

interface OperationProps {
  rule: AutoOpsRule.AsObject;
  isActiveSelected: boolean;
  handleOpenUpdate: (arg) => void;
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

  const opsCounts = useSelector<AppState, OpsCount.AsObject[]>(
    (state) => selectAllOpsCounts(state.opsCounts),
    shallowEqual
  );

  const { typeUrl } = rule.clausesList[0].clause;
  const type = typeUrl.substring(typeUrl.lastIndexOf('/') + 1);

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
          <DateTimeOperation rule={rule} isActiveSelected={isActiveSelected} />
        )}
        {type === ClauseType.EVENT_RATE && (
          <EventRateOperation rule={rule} opsCounts={opsCounts} />
        )}
      </div>
    </div>
  );
};

interface DateTimeOperationProps {
  rule: AutoOpsRule.AsObject;
  isActiveSelected: boolean;
}

const DateTimeOperation = memo(
  ({ rule, isActiveSelected }: DateTimeOperationProps) => {
    const { value } = rule.clausesList[0].clause;

    const datetimeClause = DatetimeClause.deserializeBinary(
      value as Uint8Array
    ).toObject();

    const datetime = dayjs(new Date(datetimeClause.time * 1000)).format(
      'YYYY-MM-DD HH:mm'
    );

    return (
      <div>
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
      </div>
    );
  }
);

const getEquallyDividedArray = (maxValue: number) => {
  const totalNumbers = 9;
  const resultArray = [];
  const step = maxValue / totalNumbers;

  for (let i = 0; i < totalNumbers; i++) {
    // Calculate the next value
    let nextValue = (i + 1) * step;

    // Round the value
    nextValue = Math.round(nextValue);

    // Push the rounded value into the result array
    resultArray.push(nextValue);
  }

  return [0, ...resultArray];
};

interface EventRateOperationProps {
  rule: AutoOpsRule.AsObject;
  opsCounts: OpsCount.AsObject[];
}

const EventRateOperation = memo(
  ({ rule, opsCounts }: EventRateOperationProps) => {
    const { value } = rule.clausesList[0].clause;

    const { goalId, minCount, threadsholdRate } =
      OpsEventRateClause.deserializeBinary(value as Uint8Array).toObject();

    const opsCount = opsCounts.find(
      (opsCount) => opsCount.autoOpsRuleId === rule.id
    );

    const threadsholdPercentageRange = getEquallyDividedArray(
      threadsholdRate * 100
    );

    let currentEventRate = 0;
    if (opsCount && opsCount.opsEventCount >= minCount) {
      currentEventRate =
        Math.round(
          (opsCount.opsEventCount / opsCount.evaluationCount) * 100 * 100
        ) / 100;
    }

    return (
      <div>
        <div className="flex items-center space-x-2 mt-3">
          <span className="text-gray-400">Goal</span>
          <span className="text-gray-500 truncate max-w-[120px]">{goalId}</span>
          <span className="text-gray-200">/</span>
          <span className="text-gray-400">Min Count</span>
          <span className="text-gray-500">{minCount}</span>
          <span className="text-gray-200">/</span>
          <span className="text-gray-400">Total Count Events</span>
          <span className="text-gray-500">
            {opsCount ? opsCount.opsEventCount : 0}
          </span>
          <span className="text-gray-200">/</span>
          <span className="text-gray-400">Current Event Rate</span>
          <span className="text-gray-500">
            {opsCount ? `${currentEventRate}%` : 0}
          </span>
          <HoverPopover
            render={() => {
              return (
                <div className="shadow p-2 rounded bg-white text-sm whitespace-nowrap -ml-28 mt-[-60px]">
                  Goal count / Evaluation count * 100
                </div>
              );
            }}
          >
            <InformationCircleIcon width={18} />
          </HoverPopover>
        </div>
        <div className="mt-3">
          <div className="flex">
            {Array(50)
              .fill('')
              .map((_, i) => {
                const percentage = i * 2;

                // Calculate percentage contain by one block. There are 46 blocks in the chart.
                const oneBlockPercentage = (threadsholdRate * 100 * i) / 46;

                let bgColor = 'bg-gray-200';

                if (
                  oneBlockPercentage <= currentEventRate &&
                  currentEventRate !== 0
                ) {
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
                    {percentage === 90 && (
                      <div className="absolute -left-6 text-sm text-pink-500 bottom-5 font-semibold">
                        Threshold
                      </div>
                    )}
                    {i !== 0 && (
                      <div className="absolute h-[8px] w-1.5 rounded-r-full bg-white" />
                    )}
                  </div>
                );
              })}
          </div>
          <div className="flex mt-2">
            {threadsholdPercentageRange.map((percentage) => (
              <div key={percentage} className="flex-1">
                {percentage}%
              </div>
            ))}
          </div>
        </div>
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
