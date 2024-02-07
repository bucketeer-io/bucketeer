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
import {
  selectAll as selectAllProgressiveRollouts,
  deleteProgressiveRollout,
  stopProgressiveRollout,
} from '@/modules/porgressiveRollout';
import { OpsCount } from '@/proto/autoops/ops_count_pb';
import { ProgressiveRollout } from '@/proto/autoops/progressive_rollout_pb';
import { AppDispatch } from '@/store';
import { Popover } from '@headlessui/react';
import {
  PlusIcon,
  DotsHorizontalIcon,
  PencilIcon,
  TrashIcon,
  InformationCircleIcon,
  ArrowNarrowLeftIcon,
  ArrowNarrowRightIcon,
  BanIcon,
  ClockIcon,
  CalendarIcon,
} from '@heroicons/react/outline';
import { SerializedError } from '@reduxjs/toolkit';
import dayjs from 'dayjs';
import React, {
  Dispatch,
  FC,
  SetStateAction,
  memo,
  useCallback,
  useEffect,
  useState,
} from 'react';
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
import { ReactComponent as UserSvg } from '../../assets/svg/user.svg';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { useCurrentEnvironment } from '../../modules/me';
import { AutoOpsRule, OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import {
  DatetimeClause,
  OpsEventRateClause,
  ProgressiveRolloutManualScheduleClause,
  ProgressiveRolloutSchedule,
  ProgressiveRolloutTemplateScheduleClause,
} from '../../proto/autoops/clause_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { classNames } from '../../utils/css';
import { isProgressiveRolloutsRunningWaiting } from '../AddProgressiveRolloutOperation';
import { AutoOpsDeleteDialog } from '../AutoOpsDeleteDialog';
import { DetailSkeleton } from '../DetailSkeleton';
import { HoverPopover } from '../HoverPopover';
import { OperationAddUpdateForm } from '../OperationAddUpdateForm';
import { Overlay } from '../Overlay';
import { ProgressiveRolloutStopDialog } from '../ProgressiveRolloutStopDialog';
import { RelativeDateText } from '../RelativeDateText';
import { Option } from '../Select';

enum SORT_TYPE {
  ASC = 'ASC',
  DESC = 'DESC',
}

enum OperationType {
  AutoOps = 'AutoOps',
  Progressive = 'Progressive',
}
interface ProgressiveRolloutWithType extends ProgressiveRollout.AsObject {
  operationType: OperationType;
}

interface AutoOpsRuleWithType extends AutoOpsRule.AsObject {
  operationType: OperationType;
}

const extractClauseType = (typeUrl: string) => {
  return typeUrl.substring(typeUrl.lastIndexOf('/') + 1);
};

const extractDatetimeFromAutoOps = (autoOps: AutoOpsRuleWithType) => {
  const { value } = autoOps.clausesList[0].clause;
  return DatetimeClause.deserializeBinary(value as Uint8Array).toObject().time;
};

const extractDatetimeFromProgressiveRollout = (
  prOperation: ProgressiveRollout.AsObject
) => {
  const { type, clause } = prOperation;
  const { value } = clause;
  const data =
    type === ProgressiveRollout.Type.TEMPLATE_SCHEDULE
      ? ProgressiveRolloutTemplateScheduleClause.deserializeBinary(
          value as Uint8Array
        ).toObject()
      : ProgressiveRolloutManualScheduleClause.deserializeBinary(
          value as Uint8Array
        ).toObject();
  const schedulesList = data.schedulesList;
  return schedulesList[schedulesList.length - 1].executeAt;
};

const sortOperations = (
  rules: AutoOpsRuleWithType[],
  rollouts: ProgressiveRolloutWithType[],
  sortType: SORT_TYPE
) => {
  const eventRateOperations = rules.filter(
    (r) =>
      extractClauseType(r.clausesList[0].clause.typeUrl) ===
      ClauseType.EVENT_RATE
  );

  const dateTimeOperations = rules.filter(
    (r) =>
      extractClauseType(r.clausesList[0].clause.typeUrl) !==
      ClauseType.EVENT_RATE
  );

  const sortedList = [...dateTimeOperations, ...rollouts].sort((a, b) => {
    const aDatetime =
      a.operationType === OperationType.AutoOps
        ? extractDatetimeFromAutoOps(a as AutoOpsRuleWithType)
        : extractDatetimeFromProgressiveRollout(
            a as ProgressiveRolloutWithType
          );

    const bDatetime =
      b.operationType === OperationType.AutoOps
        ? extractDatetimeFromAutoOps(b as AutoOpsRuleWithType)
        : extractDatetimeFromProgressiveRollout(
            b as ProgressiveRolloutWithType
          );

    return sortType === SORT_TYPE.ASC
      ? aDatetime - bDatetime
      : bDatetime - aDatetime;
  });

  return [...eventRateOperations, ...sortedList];
};

const TabLabel = {
  ACTIVE: intl.formatMessage(messages.autoOps.active),
  FINISHED: intl.formatMessage(messages.autoOps.finished),
};

export interface ClauseTypeMap {
  EVENT_RATE: 'bucketeer.autoops.OpsEventRateClause';
  DATETIME: 'bucketeer.autoops.DatetimeClause';
  PROGRESSIVE_ROLLOUT: 'bucketeer.autoops.ProgressiveRolloutClause';
}

export const ClauseType: ClauseTypeMap = {
  EVENT_RATE: 'bucketeer.autoops.OpsEventRateClause',
  DATETIME: 'bucketeer.autoops.DatetimeClause',
  PROGRESSIVE_ROLLOUT: 'bucketeer.autoops.ProgressiveRolloutClause',
};

interface Tab {
  label: string;
  value: (AutoOpsRuleWithType | ProgressiveRolloutWithType)[];
  selected: boolean;
}

export const getIntervalForDayjs = (
  interval: ProgressiveRolloutTemplateScheduleClause.IntervalMap[keyof ProgressiveRolloutTemplateScheduleClause.IntervalMap]
) => {
  if (Number(interval) === 1) {
    return 'hour';
  } else if (Number(interval) === 2) {
    return 'day';
  } else if (Number(interval) === 3) {
    return 'week';
  }
};

export interface SelectedOperation {
  type: keyof typeof ClauseType;
  id: string;
}

interface FeatureAutoOpsRulesFormProps {
  featureId: string;
  refetchAutoOpsRules: () => void;
  refetchProgressiveRollouts: () => void;
  reset: () => void;
}

export const FeatureAutoOpsRulesForm: FC<FeatureAutoOpsRulesFormProps> = memo(
  ({ featureId, refetchAutoOpsRules, refetchProgressiveRollouts, reset }) => {
    const { operationId } = useParams<{ operationId: string }>();
    const isNew = operationId === ID_NEW;
    const dispatch = useDispatch<AppDispatch>();

    const [selectedAutoOpsRule, setSelectedAutoOpsRule] =
      useState<AutoOpsRule.AsObject | null>(null);
    const [isDeleteConfirmDialogOpen, setIsDeleteConfirmDialogOpen] =
      useState(false);
    const [isStopConfirmDialogOpen, setIsStopConfirmDialogOpen] =
      useState(false);
    const [selectedOperation, setSelectedOperation] =
      useState<SelectedOperation>(null);

    const autoOpsRules = useSelector<AppState, AutoOpsRule.AsObject[]>(
      (state) =>
        selectAllAutoOpsRules(state.autoOpsRules).filter(
          (rule) => rule.featureId === featureId
        ),
      shallowEqual
    );
    const progressiveRollout = useSelector<
      AppState,
      ProgressiveRollout.AsObject[]
    >(
      (state) =>
        selectAllProgressiveRollouts(state.progressiveRollout).filter(
          (rule) => rule.featureId === featureId
        ),
      shallowEqual
    );

    const isAutoOpsRuleLoading = useSelector<AppState, boolean>(
      (state) => state.autoOpsRules.loading,
      shallowEqual
    );
    const isProgressiveRolloutsLoading = useSelector<AppState, boolean>(
      (state) => state.progressiveRollout.loading,
      shallowEqual
    );

    const [open, setOpen] = useState(isNew);
    const [isKillSwitchSelected, setIsKillSwitchSelected] = useState(false);
    const [isProgressiveRolloutSelected, setIsProgressiveRolloutSelected] =
      useState(false);

    const history = useHistory();
    const currentEnvironment = useCurrentEnvironment();

    const methods = useFormContext();

    const { handleSubmit } = methods;

    const [tabs, setTabs] = useState<Tab[]>([]);

    useEffect(() => {
      const rules = autoOpsRules.map((r) => ({
        ...r,
        operationType: OperationType.AutoOps,
      }));
      const rollouts = progressiveRollout.map((r) => ({
        ...r,
        operationType: OperationType.Progressive,
      }));

      setTabs([
        {
          label: TabLabel.ACTIVE,
          value: sortOperations(
            rules.filter((rule) => !rule.triggeredAt),
            rollouts.filter((p) =>
              isProgressiveRolloutsRunningWaiting(p.status)
            ),
            SORT_TYPE.ASC
          ),
          selected: true,
        },
        {
          label: TabLabel.FINISHED,
          value: sortOperations(
            rules.filter((rule) => rule.triggeredAt),
            rollouts.filter(
              (p) => !isProgressiveRolloutsRunningWaiting(p.status)
            ),
            SORT_TYPE.DESC
          ),
          selected: false,
        },
      ]);
    }, [autoOpsRules, progressiveRollout, setTabs]);

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

    const handleOnSubmitProgressiveRollout = useCallback(() => {
      handleClose();
      refetchProgressiveRollouts();
    }, []);

    const handleDelete = (operation: SelectedOperation) => {
      setIsDeleteConfirmDialogOpen(true);
      setSelectedOperation(operation);
    };

    const handleDeleteConfirm = () => {
      setIsDeleteConfirmDialogOpen(false);
      if (
        selectedOperation.type === 'DATETIME' ||
        selectedOperation.type === 'EVENT_RATE'
      ) {
        dispatch(
          deleteAutoOpsRule({
            environmentNamespace: currentEnvironment.id,
            id: selectedOperation.id,
          })
        ).then(() => {
          refetchAutoOpsRules();
        });
      } else if (selectedOperation.type === 'PROGRESSIVE_ROLLOUT') {
        dispatch(
          deleteProgressiveRollout({
            environmentNamespace: currentEnvironment.id,
            id: selectedOperation.id,
          })
        ).then(() => {
          refetchProgressiveRollouts();
        });
      }
    };

    const handleStop = (operation: SelectedOperation) => {
      setIsStopConfirmDialogOpen(true);
      setSelectedOperation(operation);
    };

    const handleStopConfirm = () => {
      setIsStopConfirmDialogOpen(false);

      dispatch(
        stopProgressiveRollout({
          environmentNamespace: currentEnvironment.id,
          id: selectedOperation.id,
        })
      ).then(() => {
        refetchProgressiveRollouts();
      });
    };

    const isActiveTabSelected =
      tabs.find((tab) => tab.selected)?.label === TabLabel.ACTIVE;

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
        <ActiveCompletedTabs tabs={tabs} setTabs={setTabs} />
        <AutoOpsInfos
          openAddOperation={handleOpen}
          activateSchedule={() => {
            setIsKillSwitchSelected(false);
            setIsProgressiveRolloutSelected(false);
          }}
          activateKillSwitch={() => {
            setIsProgressiveRolloutSelected(false);
            setIsKillSwitchSelected(true);
          }}
          activateProgressiveRollout={() => {
            setIsKillSwitchSelected(false);
            setIsProgressiveRolloutSelected(true);
          }}
        />
        {isAutoOpsRuleLoading || isProgressiveRolloutsLoading ? (
          <DetailSkeleton />
        ) : (
          <div className="space-y-6 py-6">
            {tabs
              .find((tab) => tab.selected)
              ?.value.map((operation) => {
                if (operation.operationType === OperationType.AutoOps) {
                  return (
                    <Operation
                      key={operation.id}
                      rule={operation as AutoOpsRule.AsObject}
                      isActiveTabSelected={isActiveTabSelected}
                      handleOpenUpdate={handleOpenUpdate}
                      handleDelete={handleDelete}
                    />
                  );
                } else if (
                  operation.operationType === OperationType.Progressive
                ) {
                  return (
                    <ProgressiveRolloutOperation
                      key={operation.id}
                      featureId={featureId}
                      isActiveTabSelected={isActiveTabSelected}
                      progressiveRollout={
                        operation as ProgressiveRollout.AsObject
                      }
                      handleDelete={handleDelete}
                      handleStop={handleStop}
                    />
                  );
                }
              })}
          </div>
        )}
        {open && (
          <Overlay open={open} onClose={handleClose}>
            <OperationAddUpdateForm
              onSubmit={handleSubmit(handleOnSubmit)}
              onSubmitProgressiveRollout={handleOnSubmitProgressiveRollout}
              onCancel={handleClose}
              featureId={featureId}
              autoOpsRule={selectedAutoOpsRule}
              isKillSwitchSelected={isKillSwitchSelected}
              isActiveTabSelected={isActiveTabSelected}
              isProgressiveRolloutSelected={isProgressiveRolloutSelected}
            />
          </Overlay>
        )}
        {isDeleteConfirmDialogOpen && (
          <AutoOpsDeleteDialog
            open={isDeleteConfirmDialogOpen}
            onConfirm={handleDeleteConfirm}
            onClose={() => {
              setIsDeleteConfirmDialogOpen(false);
              setSelectedOperation(null);
            }}
            selectedOperation={selectedOperation}
          />
        )}
        {isStopConfirmDialogOpen && (
          <ProgressiveRolloutStopDialog
            open={isStopConfirmDialogOpen}
            onConfirm={handleStopConfirm}
            onClose={() => {
              setIsStopConfirmDialogOpen(false);
              setSelectedOperation(null);
            }}
          />
        )}
      </div>
    );
  }
);

interface ActiveCompletedTabsProps {
  tabs: Tab[];
  setTabs: Dispatch<SetStateAction<Tab[]>>;
}

const ActiveCompletedTabs: FC<ActiveCompletedTabsProps> = memo(
  ({ tabs, setTabs }) => {
    const handleClick = (tabLabel) => {
      setTabs(
        tabs.map((t) => ({
          ...t,
          selected: t.label === tabLabel,
        }))
      );
    };

    return (
      <div className="flex border-b border-gray-200 mt-2">
        {tabs.map((tab) => {
          return (
            <div
              key={tab.label}
              className={classNames(
                'px-4 py-3 cursor-pointer',
                tab.selected
                  ? 'text-primary border-b-2 border-primary'
                  : 'text-gray-400'
              )}
              onClick={() => handleClick(tab.label)}
            >
              {tab.label} ({tab.value.length})
            </div>
          );
        })}
      </div>
    );
  }
);

interface AutoOpsInfosProps {
  openAddOperation: () => void;
  activateSchedule: () => void;
  activateKillSwitch: () => void;
  activateProgressiveRollout: () => void;
}

const AutoOpsInfos: FC<AutoOpsInfosProps> = memo(
  ({
    openAddOperation,
    activateSchedule,
    activateKillSwitch,
    activateProgressiveRollout,
  }) => (
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
              activateSchedule();
              openAddOperation();
            },
          },
          {
            id: 2,
            title: intl.formatMessage(messages.autoOps.killSwitch),
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
              activateKillSwitch();
              openAddOperation();
            },
          },
          {
            id: 3,
            title: intl.formatMessage(messages.autoOps.progressiveRollout),
            detail: intl.formatMessage(
              messages.autoOps.infoBlocks.progressiveRolloutInfo
            ),
            bgColor: 'bg-blue-50',
            icon: <ArrowTrendingUp />,
            onClick: () => {
              activateProgressiveRollout();
              openAddOperation();
            },
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
  )
);

interface ProgressiveRolloutProps {
  featureId: string;
  isActiveTabSelected: boolean;
  progressiveRollout: ProgressiveRollout.AsObject;
  handleDelete: (arg: SelectedOperation) => void;
  handleStop: (arg: SelectedOperation) => void;
}

const ProgressiveRolloutOperation: FC<ProgressiveRolloutProps> = memo(
  ({
    featureId,
    isActiveTabSelected,
    progressiveRollout,
    handleDelete,
    handleStop,
  }) => {
    const [feature] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ]);

    const variationOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: v.value,
      };
    });

    const { value } = progressiveRollout.clause;

    if (progressiveRollout.type === ProgressiveRollout.Type.TEMPLATE_SCHEDULE) {
      const data = ProgressiveRolloutTemplateScheduleClause.deserializeBinary(
        value as Uint8Array
      ).toObject();

      const { schedulesList, increments, interval, variationId } = data;
      return (
        <ProgressiveRolloutComponent
          key={progressiveRollout.id}
          variationOptions={variationOptions}
          rule={progressiveRollout}
          deleteRule={() =>
            handleDelete({
              type: 'PROGRESSIVE_ROLLOUT',
              id: progressiveRollout.id,
            })
          }
          stopRule={() =>
            handleStop({
              type: 'PROGRESSIVE_ROLLOUT',
              id: progressiveRollout.id,
            })
          }
          schedulesList={schedulesList}
          increments={increments}
          interval={interval}
          variationId={variationId}
          isActiveTabSelected={isActiveTabSelected}
        />
      );
    } else if (
      progressiveRollout.type === ProgressiveRollout.Type.MANUAL_SCHEDULE
    ) {
      const data = ProgressiveRolloutManualScheduleClause.deserializeBinary(
        value as Uint8Array
      ).toObject();

      const { schedulesList, variationId } = data;

      return (
        <ProgressiveRolloutComponent
          key={progressiveRollout.id}
          variationOptions={variationOptions}
          rule={progressiveRollout}
          deleteRule={() =>
            handleDelete({
              type: 'PROGRESSIVE_ROLLOUT',
              id: progressiveRollout.id,
            })
          }
          stopRule={() =>
            handleStop({
              type: 'PROGRESSIVE_ROLLOUT',
              id: progressiveRollout.id,
            })
          }
          schedulesList={schedulesList}
          variationId={variationId}
          isActiveTabSelected={isActiveTabSelected}
        />
      );
    }
    return null;
  }
);

interface OperationProps {
  rule: AutoOpsRule.AsObject;
  isActiveTabSelected: boolean;
  handleOpenUpdate: (arg) => void;
  handleDelete: (arg: SelectedOperation) => void;
}

const Operation: FC<OperationProps> = memo(
  ({ rule, isActiveTabSelected, handleOpenUpdate, handleDelete }) => {
    const { formatMessage: f } = useIntl();
    const opsCounts = useSelector<AppState, OpsCount.AsObject[]>(
      (state) => selectAllOpsCounts(state.opsCounts),
      shallowEqual
    );

    const { typeUrl } = rule.clausesList[0].clause;
    const type = typeUrl.substring(typeUrl.lastIndexOf('/') + 1);

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
                        {type === ClauseType.DATETIME &&
                          f(messages.autoOps.editSchedule)}
                        {type === ClauseType.EVENT_RATE &&
                          f(messages.autoOps.editKillSwitch)}
                      </span>
                    </button>
                    <button
                      onClick={() => {
                        if (type === ClauseType.DATETIME) {
                          handleDelete({
                            type: 'DATETIME',
                            id: rule.id,
                          });
                        } else if (type === ClauseType.EVENT_RATE) {
                          handleDelete({
                            type: 'EVENT_RATE',
                            id: rule.id,
                          });
                        }
                      }}
                      className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100"
                    >
                      <TrashIcon width={18} className="text-red-500" />
                      <span className="text-red-500 text-sm">
                        {type === ClauseType.DATETIME &&
                          f(messages.autoOps.deleteSchedule)}
                        {type === ClauseType.EVENT_RATE &&
                          f(messages.autoOps.deleteKillSwitch)}
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
  }
);

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
      'YYYY/MM/DD HH:mm'
    );

    const createdAt = dayjs(new Date(rule.createdAt * 1000)).format(
      'YYYY/MM/DD HH:mm'
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

    const numberOfSteps =
      Math.round(threadsholdRate * 100) > 10
        ? 10
        : Math.round(threadsholdRate * 100);
    const step = (threadsholdRate * 100) / numberOfSteps;

    const stepArray = Array.from({ length: numberOfSteps }, (_, index) =>
      Math.round(step + index * step)
    );

    const barWidth = isActiveTabSelected
      ? (currentEventRate / (threadsholdRate * 100)) * 100
      : 100;

    return (
      <div>
        <div className="flex items-center space-x-2 mt-3">
          <span className="text-gray-400">
            {f(messages.autoOps.opsEventRateClause.goal)}:
          </span>
          <span className="text-gray-500 truncate max-w-[120px]">{goalId}</span>
          <span className="text-gray-200">/</span>
          <span className="text-gray-400">
            {f(messages.autoOps.minimumGoalCount)}:
          </span>
          <span className="text-gray-500">{minCount}</span>
          <span className="text-gray-200">/</span>
          <span className="text-gray-400">
            {f(messages.autoOps.totalGoalCountEvents)}:
          </span>
          <span className="text-gray-500">
            {opsCount ? opsCount.opsEventCount : 0}
          </span>
          <span className="text-gray-200">/</span>
          <span className="text-gray-400">
            {f(messages.autoOps.currentEventRate)}:
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
        <div className="bg-gray-50 pt-16 pb-10 px-12 rounded mt-2 relative">
          <div className="absolute right-11 top-[6px] text-sm text-pink-500 font-semibold whitespace-nowrap">
            {f(messages.autoOps.threshold)}
          </div>
          <div className="flex h-[4px] bg-gray-200 relative">
            <div
              className="bg-pink-500 absolute h-[4px]"
              style={{
                width: `${barWidth}%`,
              }}
            />
            <div className={classNames('flex items-center h-[4px]')}>
              <div
                className={classNames(
                  'w-[9px] h-[9px] rounded-full relative',
                  currentEventRate > 0 || !isActiveTabSelected
                    ? 'bg-pink-500'
                    : 'border border-gray-400 bg-gray-50'
                )}
              >
                <span className="absolute -top-8 left-1/2 -translate-x-1/2">
                  0%
                </span>
              </div>
            </div>
            {stepArray.map((percentage) => {
              const isActive = isActiveTabSelected
                ? percentage <= currentEventRate && currentEventRate !== 0
                : true;

              return (
                <div
                  key={percentage}
                  className={classNames(
                    'flex justify-end flex-1 items-center h-[4px]'
                  )}
                >
                  <div
                    className={classNames(
                      'w-[9px] h-[9px] rounded-full relative',
                      isActive
                        ? 'bg-pink-500'
                        : 'border border-gray-400 bg-gray-50'
                    )}
                  >
                    <span className="absolute -top-8 left-1/2 -translate-x-1/2">
                      {percentage}%
                    </span>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    );
  }
);

interface ProgressiveRolloutTemplateScheduleProps {
  variationOptions: Option[];
  rule: ProgressiveRollout.AsObject;
  deleteRule: () => void;
  stopRule: () => void;
  schedulesList: ProgressiveRolloutSchedule.AsObject[];
  increments?: number;
  interval?: ProgressiveRolloutTemplateScheduleClause.IntervalMap[keyof ProgressiveRolloutTemplateScheduleClause.IntervalMap];
  variationId: string;
  isActiveTabSelected: boolean;
}

const ProgressiveRolloutComponent = memo(
  ({
    variationOptions,
    rule,
    deleteRule,
    stopRule,
    schedulesList,
    increments,
    interval,
    variationId,
    isActiveTabSelected,
  }: ProgressiveRolloutTemplateScheduleProps) => {
    const { formatMessage: f } = useIntl();

    const [selectedPagination, setSelectedPagination] = useState(0);

    const getFrequency = (frequency) => {
      if (frequency === 1) {
        return f(messages.autoOps.hourly);
      } else if (frequency === 2) {
        return f(messages.autoOps.daily);
      } else if (frequency === 3) {
        return f(messages.autoOps.weekly);
      } else {
        return null;
      }
    };

    const lastItemWithTriggeredAt = [...schedulesList]
      .reverse()
      .find((s) => s.triggeredAt);

    const isSameOrBeforeOfLastTriggerAt = (executeAt) => {
      if (lastItemWithTriggeredAt) {
        return dayjs(executeAt).isSameOrBefore(
          lastItemWithTriggeredAt.executeAt * 1000
        );
      }
      return false;
    };
    const totalNumberOfPages = Math.ceil(schedulesList.length / 10);

    const paginatedScheduleList = schedulesList.slice(
      selectedPagination * 10,
      (selectedPagination + 1) * 10
    );

    const firstSchedule = {
      weight:
        selectedPagination === 0
          ? 0
          : schedulesList[selectedPagination * 10 - 1].weight / 1000,
      executeAt:
        selectedPagination === 0
          ? rule.createdAt
          : schedulesList[selectedPagination * 10 - 1].executeAt,
    };

    console.log({ rule });

    return (
      <div className="rounded-xl shadow px-6 py-4 bg-white">
        <div className="flex justify-between py-4 border-b">
          <h3 className="font-bold text-xl">
            {f(messages.autoOps.enableOperation)}
          </h3>
          <div className="flex space-x-2 items-center">
            <div className="py-[2px] px-2 bg-[#FFF7EE] rounded text-[#CE844A] text-sm">
              {f(messages.autoOps.progressiveRollout)}
            </div>
            {isActiveTabSelected && (
              <Popover className="relative flex">
                <Popover.Button>
                  <div className="pl-2 flex items-center cursor-pointer">
                    <DotsHorizontalIcon width={20} />
                  </div>
                </Popover.Button>
                <Popover.Panel className="absolute z-10 bg-white right-0 rounded-lg p-1 whitespace-nowrap shadow-md">
                  <button
                    onClick={deleteRule}
                    className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100"
                  >
                    <TrashIcon width={18} className="text-red-500" />
                    <span className="text-red-500 text-sm">
                      {f(messages.autoOps.deleteProgressiveRollout)}
                    </span>
                  </button>
                  <button
                    onClick={stopRule}
                    className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100"
                  >
                    <BanIcon width={18} className="" />
                    <span className="text-sm">
                      {f(messages.autoOps.stopProgressiveRollout)}
                    </span>
                  </button>
                </Popover.Panel>
              </Popover>
            )}
          </div>
        </div>
        <div className="mt-4">
          <p className="font-bold text-lg">
            {f(messages.autoOps.progressInformation)}
          </p>
          <div className="flex py-3 items-center justify-between">
            <div className="flex items-center space-x-2">
              {increments && (
                <>
                  <div className="space-x-1 items-center flex">
                    <span className="text-gray-400">
                      {f(messages.autoOps.increment)}:
                    </span>
                    <span className="text-gray-500">{increments}%</span>
                    <InformationCircleIcon
                      width={18}
                      className="text-gray-400"
                    />
                  </div>
                  <span className="text-gray-200">/</span>
                </>
              )}
              <div className="flex space-x-1">
                <span className="text-gray-400">
                  {f(messages.autoOps.startDate)}:
                </span>
                <span className="text-gray-500">
                  {dayjs(schedulesList[0].executeAt * 1000).format(
                    'YYYY/MM/DD'
                  )}
                </span>
              </div>
              <span className="text-gray-200">/</span>
              {interval && (
                <>
                  <div className="flex space-x-1">
                    <span className="text-gray-400">
                      {f(messages.autoOps.frequency)}:
                    </span>
                    <span className="text-gray-500">
                      {getFrequency(interval)}
                    </span>
                  </div>
                  <span className="text-gray-200">/</span>
                </>
              )}
              <div className="flex space-x-1">
                <span className="text-gray-400">
                  {f(messages.feature.variation)}:
                </span>
                <span className="text-gray-500">
                  {variationOptions.find((v) => v.value === variationId)?.label}
                </span>
              </div>
            </div>
            {rule.status === ProgressiveRollout.Status.STOPPED && (
              <div className="flex space-x-[6px] text-gray-500 items-center">
                <ClockIcon width={18} />
                <span>Stopped</span>
                {<RelativeDateText date={new Date(rule.stoppedAt * 1000)} />}
                <span>by</span>
                {rule.stoppedBy ===
                  ProgressiveRollout.StoppedBy.OPS_KILL_SWITCH && (
                  <>
                    <div className="relative">
                      <RefreshSvg width={22} />
                      <CrossSvg
                        width={12}
                        className="absolute right-[1px] bottom-[3px]"
                      />
                    </div>
                    <span>Kill Switch</span>
                  </>
                )}
                {rule.stoppedBy === ProgressiveRollout.StoppedBy.USER && (
                  <>
                    <UserSvg />
                    <span>User</span>
                  </>
                )}
                {rule.stoppedBy ===
                  ProgressiveRollout.StoppedBy.OPS_SCHEDULE && (
                  <>
                    <CalendarIcon width={18} className="text-primary" />
                    <span>Schedule</span>
                  </>
                )}
              </div>
            )}
          </div>
          <div className="bg-gray-50 pt-14 pb-16 px-12 rounded mt-2">
            <div className="flex h-[4px] bg-gray-200 relative">
              <div className="h-[4px] flex items-center">
                <div
                  className={classNames(
                    'w-[9px] h-[9px] rounded-full relative',
                    isSameOrBeforeOfLastTriggerAt(firstSchedule.executeAt)
                      ? 'bg-pink-500'
                      : 'border border-gray-400 bg-gray-50'
                  )}
                >
                  <span className="absolute -top-8 left-1/2 -translate-x-1/2">
                    {firstSchedule.weight}%
                  </span>
                  <div className="absolute top-[18px] left-1/2 -translate-x-1/2 whitespace-nowrap text-center">
                    <p className="text-gray-400 text-xs">
                      {dayjs(firstSchedule.executeAt * 1000).format('HH:mm')}
                    </p>
                    <p className="text-gray-400 text-xs">
                      {dayjs(firstSchedule.executeAt * 1000).format(
                        'YYYY/MM/DD'
                      )}
                    </p>
                  </div>
                </div>
              </div>
              {paginatedScheduleList.map((schedule, i) => (
                <div
                  key={i}
                  className={classNames(
                    'flex justify-end flex-1 items-center h-[4px]',
                    isSameOrBeforeOfLastTriggerAt(schedule.executeAt * 1000) &&
                      'bg-pink-500'
                  )}
                >
                  <div
                    className={classNames(
                      'w-[9px] h-[9px] rounded-full relative',
                      isSameOrBeforeOfLastTriggerAt(schedule.executeAt * 1000)
                        ? 'bg-pink-500'
                        : 'border border-gray-400 bg-gray-50'
                    )}
                  >
                    <span className="absolute -top-8 left-1/2 -translate-x-1/2">
                      {schedule.weight / 1000}%
                    </span>
                    <div className="absolute top-[18px] left-1/2 -translate-x-1/2 whitespace-nowrap text-center">
                      <p className="text-gray-400 text-xs">
                        {dayjs(schedule.executeAt * 1000).format('HH:mm')}
                      </p>
                      <p className="text-gray-400 text-xs">
                        {dayjs(schedule.executeAt * 1000).format('YYYY/MM/DD')}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
          {totalNumberOfPages > 1 && (
            <div className="mt-4 flex justify-between items-center">
              <button
                className={classNames(
                  'p-1.5 rounded border',
                  selectedPagination === 0 && 'opacity-50 cursor-not-allowed'
                )}
                disabled={selectedPagination === 0}
                onClick={() => setSelectedPagination(selectedPagination - 1)}
              >
                <ArrowNarrowLeftIcon width={16} className="text-gray-400" />
              </button>
              <div className="flex space-x-2">
                {Array(totalNumberOfPages)
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
                className={classNames(
                  'p-1.5 rounded border',
                  selectedPagination === totalNumberOfPages - 1 &&
                    'opacity-50 cursor-not-allowed'
                )}
                disabled={selectedPagination === totalNumberOfPages - 1}
                onClick={() => setSelectedPagination(selectedPagination + 1)}
              >
                <ArrowNarrowRightIcon width={16} className="text-gray-400" />
              </button>
            </div>
          )}
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
