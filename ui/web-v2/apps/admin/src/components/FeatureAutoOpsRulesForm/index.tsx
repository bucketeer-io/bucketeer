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
import { useIntl } from 'react-intl';
import { useSelector, shallowEqual, useDispatch } from 'react-redux';
import { useParams, useHistory } from 'react-router-dom';
import { v4 as uuid } from 'uuid';

import { ReactComponent as ArrowTrendingUp } from '../../assets/svg/arrow-trending-up.svg';
import { ReactComponent as CalendarSvg } from '../../assets/svg/calendar.svg';
import { ReactComponent as CrossSvg } from '../../assets/svg/cross.svg';
import { ReactComponent as OpenInNewSvg } from '../../assets/svg/open-new-tab.svg';
import { ReactComponent as RefreshSvg } from '../../assets/svg/refresh.svg';
import { ReactComponent as SeeDetailsSvg } from '../../assets/svg/see-details.svg';
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

const numberOfBlocks = 51;

const TabLabel = {
  ACTIVE: intl.formatMessage(messages.autoOps.active),
  COMPLETED: intl.formatMessage(messages.autoOps.completed),
};

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
    const [isKillSwitchSelected, setIsKillSwitchSelected] = useState(false);
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
              environmentNamespace: currentEnvironment.id,
              ids,
            })
          );
        }
      }
    }, [autoOpsRules]);

    const handleClose = useCallback(() => {
      reset();
      history.replace({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_AUTOOPS}`,
        search: location.search,
      });
      setOpen(false);
      setIsKillSwitchSelected(false);
    }, [setOpen, history, location, reset]);

    const handleOpen = useCallback(() => {
      setOpen(true);
      history.push({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_AUTOOPS}${PAGE_PATH_NEW}`,
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

    const isActiveTabSelected =
      tabs.find((tab) => tab.selected).label === TabLabel.ACTIVE;

    return (
      <div className="px-10 py-6 bg-white">
        <div className="flex justify-end">
          <a
            className="space-x-2 flex items-center justify-center mr-5 text-primary cursor-pointer"
            href="https://docs.bucketeer.io/feature-flags/creating-feature-flags/auto-operation/"
            target="_blank"
            rel="noreferrer"
          >
            <OpenInNewSvg className="mt-[2px]" />
            <span className="underline">
              {intl.formatMessage(messages.sideMenu.documentation)}
            </span>
          </a>
          <button
            onClick={() => {
              setSelectedAutoOpsRule(null);
              handleOpen();
            }}
            className="btn-submit space-x-2"
          >
            <PlusIcon width={18} />
            <span>{intl.formatMessage(messages.button.add)}</span>
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
        <div className="py-6">
          <p className="text-xl font-bold">
            {intl.formatMessage(messages.autoOps.infoBlocks.title)}
          </p>
          <div className="flex space-x-6 mt-6">
            {[
              {
                id: 1,
                title: intl.formatMessage(messages.autoOps.schedule),
                detail: intl.formatMessage(
                  messages.autoOps.infoBlocks.scheduleInfo
                ),
                bgColor: 'bg-purple-50',
                icon: <CalendarSvg />,
                onClick: () => {
                  handleOpen();
                },
              },
              {
                id: 2,
                title: intl.formatMessage(
                  messages.autoOps.infoBlocks.killSwitch
                ),
                detail: intl.formatMessage(
                  messages.autoOps.infoBlocks.killSwitchInfo
                ),
                bgColor: 'bg-pink-50',
                icon: (
                  <div className="relative">
                    <RefreshSvg />
                    <CrossSvg className="absolute right-[2px] bottom-[1px]" />
                  </div>
                ),
                onClick: () => {
                  setIsKillSwitchSelected(true);
                  handleOpen();
                },
              },
              {
                id: 3,
                title: intl.formatMessage(
                  messages.autoOps.infoBlocks.progressiveRollout
                ),
                detail: intl.formatMessage(
                  messages.autoOps.infoBlocks.progressiveRolloutInfo
                ),
                bgColor: 'bg-blue-50',
                icon: <ArrowTrendingUp />,
                onclick: () => {},
              },
            ].map(({ id, title, detail, bgColor, icon, onClick }) => (
              <div
                key={id}
                className="flex flex-1 space-x-4 p-4 rounded-md shadow-md cursor-pointer"
                onClick={onClick}
              >
                <div
                  className={classNames(
                    'w-16 h-16 rounded-lg flex justify-center items-center',
                    bgColor
                  )}
                >
                  {icon}
                </div>
                <div className="flex-1">
                  <p className="text-lg font-bold">{title}</p>
                  <p className="">{detail}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
        <div className="space-y-6 py-6">
          {tabs
            .find((tab) => tab.selected)
            .value.map((rule) => (
              <Operation
                key={rule.id}
                rule={rule}
                isActiveTabSelected={isActiveTabSelected}
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
              isKillSwitchSelected={isKillSwitchSelected}
              isActiveTabSelected={isActiveTabSelected}
            />
          </Overlay>
        )}
      </div>
    );
  }
);

interface OperationProps {
  rule: AutoOpsRule.AsObject;
  isActiveTabSelected: boolean;
  handleOpenUpdate: (arg) => void;
  refetchAutoOpsRules: () => void;
}

const Operation = ({
  rule,
  isActiveTabSelected,
  handleOpenUpdate,
  refetchAutoOpsRules,
}: OperationProps) => {
  const { formatMessage: f } = useIntl();
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
        environmentNamespace: currentEnvironment.id,
        id: ruleId,
      })
    ).then(refetchAutoOpsRules);
  };

  return (
    <div className="rounded-xl shadow px-6 py-4 bg-white">
      <div className="flex justify-between py-4 border-b">
        <h3 className="font-bold text-xl text-gray-600">
          {rule.opsType === OpsType.ENABLE_FEATURE
            ? f(messages.autoOps.enableOperation)
            : f(messages.autoOps.killSwitchOperation)}
        </h3>
        <div className="flex space-x-2 items-center">
          <div
            className={classNames(
              'py-[2px] px-2 rounded text-sm',
              type === ClauseType.DATETIME && 'bg-[#EBF9ED] text-green-700',
              type === ClauseType.EVENT_RATE && 'bg-[#EFECF5] text-primary'
            )}
          >
            {type === ClauseType.DATETIME && f(messages.autoOps.schedule)}
            {type === ClauseType.EVENT_RATE && f(messages.autoOps.eventRate)}
          </div>
          <Popover className="relative flex">
            <Popover.Button>
              <div className="pl-2 flex items-center cursor-pointer">
                <DotsHorizontalIcon width={20} />
              </div>
            </Popover.Button>
            <Popover.Panel className="absolute z-10 bg-white right-0 rounded-lg p-1 whitespace-nowrap shadow-md">
              {isActiveTabSelected ? (
                <>
                  <button
                    onClick={() => handleOpenUpdate(rule)}
                    className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                  >
                    <PencilIcon width={18} />
                    <span className="text-sm">
                      {f(messages.autoOps.editOperation)}
                    </span>
                  </button>
                  <button
                    onClick={() => handleDelete(rule.id)}
                    className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100"
                  >
                    <TrashIcon width={18} className="text-red-500" />
                    <span className="text-red-500 text-sm">
                      {f(messages.autoOps.deleteOperation)}
                    </span>
                  </button>
                </>
              ) : (
                <button
                  onClick={() => handleOpenUpdate(rule)}
                  className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                >
                  <SeeDetailsSvg />
                  <span className="text-sm">
                    {f(messages.autoOps.operationDetails)}
                  </span>
                </button>
              )}
            </Popover.Panel>
          </Popover>
        </div>
      </div>
      <div className="mt-4">
        <p className="font-bold text-lg text-gray-600">
          {f(messages.autoOps.progressInformation)}
        </p>
        {type === ClauseType.DATETIME && (
          <DateTimeOperation
            rule={rule}
            isActiveTabSelected={isActiveTabSelected}
          />
        )}
        {type === ClauseType.EVENT_RATE && (
          <EventRateOperation
            rule={rule}
            opsCounts={opsCounts}
            isActiveTabSelected={isActiveTabSelected}
          />
        )}
      </div>
    </div>
  );
};

interface DateTimeOperationProps {
  rule: AutoOpsRule.AsObject;
  isActiveTabSelected: boolean;
}

const DateTimeOperation = memo(
  ({ rule, isActiveTabSelected }: DateTimeOperationProps) => {
    const { value } = rule.clausesList[0].clause;

    const datetimeClause = DatetimeClause.deserializeBinary(
      value as Uint8Array
    ).toObject();

    const datetime = dayjs(new Date(datetimeClause.time * 1000)).format(
      'YYYY-MM-DD HH:mm'
    );

    const createdAt = dayjs(new Date(rule.createdAt * 1000)).format(
      'YYYY-MM-DD HH:mm'
    );

    return (
      <div>
        <div
          className={classNames(
            'mt-6 h-2  flex justify-between relative mx-1',
            isActiveTabSelected ? 'bg-gray-200' : 'bg-pink-500'
          )}
        >
          <div className="w-[14px] h-[14px] absolute top-1/2 -translate-y-1/2 rounded-full -left-1 bg-pink-500 border border-pink-100" />
          <div
            className={classNames(
              'w-[14px] h-[14px] absolute top-1/2 -translate-y-1/2 rounded-full -right-1 bg-gray-300 border',
              isActiveTabSelected ? 'bg-gray-200' : 'bg-pink-500'
            )}
          />
        </div>
        <div className="flex justify-between mt-2">
          <span>Off</span>
          <span>On</span>
        </div>
        <div className="flex justify-between mt-1">
          <span className="text-xs text-gray-400">{createdAt}</span>
          <span className="text-xs text-gray-400">{datetime}</span>
        </div>
      </div>
    );
  }
);

interface EventRateOperationProps {
  rule: AutoOpsRule.AsObject;
  opsCounts: OpsCount.AsObject[];
  isActiveTabSelected: boolean;
}

const EventRateOperation = memo(
  ({ rule, opsCounts, isActiveTabSelected }: EventRateOperationProps) => {
    const { value } = rule.clausesList[0].clause;
    const { formatMessage: f } = useIntl();

    const { goalId, minCount, threadsholdRate } =
      OpsEventRateClause.deserializeBinary(value as Uint8Array).toObject();

    const opsCount = opsCounts.find(
      (opsCount) => opsCount.autoOpsRuleId === rule.id
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
          <span className="text-gray-400">
            {f(messages.autoOps.opsEventRateClause.goal)}
          </span>
          <span className="text-gray-500 truncate max-w-[120px]">{goalId}</span>
          <span className="text-gray-200">/</span>
          <span className="text-gray-400">
            {f(messages.autoOps.minimumGoalCount)}
          </span>
          <span className="text-gray-500">{minCount}</span>
          <span className="text-gray-200">/</span>
          <span className="text-gray-400">
            {f(messages.autoOps.totalGoalCountEvents)}
          </span>
          <span className="text-gray-500">
            {opsCount ? opsCount.opsEventCount : 0}
          </span>
          <span className="text-gray-200">/</span>
          <span className="text-gray-400">
            {f(messages.autoOps.currentEventRate)}
          </span>
          <span className="text-gray-500">
            {opsCount ? `${currentEventRate}%` : '0%'}
          </span>
          <HoverPopover
            render={() => {
              return (
                <div className="shadow p-2 rounded bg-white text-sm whitespace-nowrap -ml-28 mt-[-60px]">
                  {f(messages.autoOps.goalCount)} /{' '}
                  {f(messages.autoOps.evaluationCount)} * 100
                </div>
              );
            }}
          >
            <InformationCircleIcon width={18} />
          </HoverPopover>
        </div>
        <div className="flex pb-7 mt-4">
          {Array(numberOfBlocks)
            .fill('')
            .map((_, i) => {
              // Calculate percentage contain by one block. There are 51 blocks in the chart.
              const oneBlockPercentage =
                (threadsholdRate * 100 * i) / numberOfBlocks;

              let bgColor = 'bg-gray-200';

              let percentage;
              if (
                oneBlockPercentage <= currentEventRate &&
                currentEventRate !== 0
              ) {
                bgColor = 'bg-pink-500';
              } else if (i % 5 === 0) {
                bgColor = 'bg-gray-400';
                const step = (threadsholdRate * 100) / 10;
                percentage = Math.round((i * step) / 5);
              }

              return (
                <div
                  key={i}
                  className={classNames(
                    'relative h-[8px] flex-1 rounded-[60px]',
                    isActiveTabSelected ? bgColor : 'bg-pink-500'
                  )}
                >
                  {i === numberOfBlocks - 1 && (
                    <div className="absolute right-0 bottom-[26px] text-sm text-pink-500 font-semibold whitespace-nowrap">
                      {f(messages.autoOps.threshold)}
                    </div>
                  )}
                  {i !== 0 && (
                    <div className="absolute h-[8px] w-1.5 rounded-r-full bg-white" />
                  )}
                  {i % 5 === 0 && (
                    <div
                      className={classNames(
                        'absolute -bottom-8',
                        i !== 0 && i < numberOfBlocks - 1 && 'left-1',
                        i === numberOfBlocks - 1 && 'right-0'
                      )}
                    >
                      {percentage}%
                    </div>
                  )}
                </div>
              );
            })}
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
    minCount: 50,
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
