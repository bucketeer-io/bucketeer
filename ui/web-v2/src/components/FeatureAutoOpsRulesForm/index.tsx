import {
  PAGE_PATH_FEATURES,
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_ROOT
} from '../../constants/routing';
import {
  deleteAutoOpsRule,
  selectAll as selectAllAutoOpsRules,
  stopAutoOpsRule
} from '../../modules/autoOpsRules';
import {
  listOpsCounts,
  selectAll as selectAllOpsCounts
} from '../../modules/opsCounts';
import {
  selectAll as selectAllProgressiveRollouts,
  deleteProgressiveRollout,
  stopProgressiveRollout
} from '../../modules/porgressiveRollout';
import { OpsCount } from '../../proto/autoops/ops_count_pb';
import { ProgressiveRollout } from '../../proto/autoops/progressive_rollout_pb';
import { AppDispatch } from '../../store';
import { Popover } from '@headlessui/react';
import {
  PlusIcon,
  DotsHorizontalIcon,
  PencilIcon,
  TrashIcon,
  InformationCircleIcon,
  BanIcon,
  ClockIcon,
  CalendarIcon,
  ChevronDownIcon
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
  useState
} from 'react';
import { useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useSelector, shallowEqual, useDispatch } from 'react-redux';
import { useParams, useHistory } from 'react-router-dom';
import ArrowTrendingUp from '../../assets/svg/arrow-trending-up.svg';
import ArrowTrendingUpGray from '../../assets/svg/arrow-trending-up-gray.svg';
import CalendarSvg from '../../assets/svg/calendar.svg';
import CrossSvg from '../../assets/svg/cross.svg';
import OpenInNewSvg from '../../assets/svg/open-new-tab.svg';
import RefreshSvg from '../../assets/svg/refresh.svg';
import RefreshPinkSvg from '../../assets/svg/refresh-pink.svg';
import SeeDetailsSvg from '../../assets/svg/see-details.svg';
import UserSvg from '../../assets/svg/user.svg';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { useCurrentEnvironment } from '../../modules/me';
import {
  AutoOpsRule,
  AutoOpsStatus,
  OpsType
} from '../../proto/autoops/auto_ops_rule_pb';
import {
  OpsEventRateClause,
  ProgressiveRolloutManualScheduleClause,
  ProgressiveRolloutSchedule,
  ProgressiveRolloutTemplateScheduleClause
} from '../../proto/autoops/clause_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { classNames } from '../../utils/css';
import { isProgressiveRolloutsRunningWaiting } from '../ProgressiveRolloutAddForm';
import { AutoOpsDeleteDialog } from '../AutoOpsDeleteDialog';
import { DetailSkeleton } from '../DetailSkeleton';
import { HoverPopover } from '../HoverPopover';
import {
  actionTypesOptions,
  ScheduleAddUpdateForm
} from '../ScheduleAddUpdateForm';
import { EventRateAddUpdateForm } from '../EventRateAddUpdateForm';
import { ProgressiveRolloutAddForm } from '../ProgressiveRolloutAddForm';
import { Overlay } from '../Overlay';
import { AutoOpsStopDialog } from '../AutoOpsStopDialog';
import { RelativeDateText } from '../RelativeDateText';
import { Option } from '../Select';
import { isLanguageJapanese } from '../../lang/getSelectedLanguage';
import { createVariationLabel } from '../../utils/variation';
import { getDatetimeClause } from '../../utils/getDatetimeClause';
import OperationPagination from '../OperationPagination';

enum SORT_TYPE {
  ASC = 'ASC',
  DESC = 'DESC'
}

export enum OperationType {
  SCHEDULE = 'schedule',
  EVENT_RATE = 'event_rate',
  PROGRESSIVE_ROLLOUT = 'progressive_rollout'
}

const extractDatetimeFromAutoOps = (autoOps: AutoOpsRule.AsObject) => {
  const { value } = autoOps.clausesList[0].clause;
  return getDatetimeClause(value).time;
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
  rules: AutoOpsRule.AsObject[],
  rollouts: ProgressiveRollout.AsObject[],
  sortType: SORT_TYPE
) => {
  const eventRateOperations = rules.filter(
    (r) => r.opsType === OpsType.EVENT_RATE
  );

  const dateTimeOperations = rules.filter(
    (r) => r.opsType === OpsType.SCHEDULE
  );

  const mergedArray = [...dateTimeOperations, ...rollouts];

  const sortedList = mergedArray.sort((a, b) => {
    const a2 = a as AutoOpsRule.AsObject & ProgressiveRollout.AsObject;
    const b2 = b as AutoOpsRule.AsObject & ProgressiveRollout.AsObject;

    const aDatetime =
      a2.opsType === OpsType.SCHEDULE || a2.opsType === OpsType.EVENT_RATE
        ? extractDatetimeFromAutoOps(a2)
        : extractDatetimeFromProgressiveRollout(a2);

    const bDatetime =
      b2.opsType === OpsType.SCHEDULE || b2.opsType === OpsType.EVENT_RATE
        ? extractDatetimeFromAutoOps(b2)
        : extractDatetimeFromProgressiveRollout(b2);

    return sortType === SORT_TYPE.ASC
      ? aDatetime - bDatetime
      : bDatetime - aDatetime;
  });

  return [...eventRateOperations, ...sortedList];
};

const TabLabel = {
  ACTIVE: intl.formatMessage(messages.autoOps.active),
  FINISHED: intl.formatMessage(messages.autoOps.finished)
};

interface Tab {
  label: string;
  value: (AutoOpsRule.AsObject | ProgressiveRollout.AsObject)[];
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
  id: string;
  type: OperationType;
}

interface FeatureAutoOpsRulesFormProps {
  featureId: string;
  refetchAutoOpsRules: () => void;
  refetchProgressiveRollouts: () => void;
  reset: () => void;
}

export const FeatureAutoOpsRulesForm: FC<FeatureAutoOpsRulesFormProps> = memo(
  ({ featureId, refetchAutoOpsRules, refetchProgressiveRollouts, reset }) => {
    const { operationType } = useParams<{
      operationType: OperationType;
    }>();
    const dispatch = useDispatch<AppDispatch>();

    const [selectedOperationType, setSelectedOperationType] = useState<
      OperationType | undefined
    >(operationType);

    const [selectedAutoOpsRule, setSelectedAutoOpsRule] =
      useState<AutoOpsRule.AsObject | null>(null);
    const [isDeleteConfirmDialogOpen, setIsDeleteConfirmDialogOpen] =
      useState(false);
    const [isStopConfirmDialogOpen, setIsStopConfirmDialogOpen] =
      useState(false);
    const [selectedOperation, setSelectedOperation] =
      useState<SelectedOperation>(null);

    const [feature] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError
    ]);
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

    const history = useHistory();
    const currentEnvironment = useCurrentEnvironment();

    const methods = useFormContext();

    const { handleSubmit, setValue } = methods;

    const [tabs, setTabs] = useState<Tab[]>([]);
    const { formatMessage: f } = useIntl();

    useEffect(() => {
      let initialTab = TabLabel.ACTIVE;

      if (tabs.length > 0) {
        initialTab =
          tabs.find((tab) => tab.selected)?.label === TabLabel.ACTIVE
            ? TabLabel.ACTIVE
            : TabLabel.FINISHED;
      }

      setTabs([
        {
          label: TabLabel.ACTIVE,
          value: sortOperations(
            autoOpsRules.filter(
              (rule) =>
                rule.autoOpsStatus === AutoOpsStatus.RUNNING ||
                rule.autoOpsStatus === AutoOpsStatus.WAITING
            ),
            progressiveRollout.filter((p) =>
              isProgressiveRolloutsRunningWaiting(p.status)
            ),
            SORT_TYPE.ASC
          ),
          selected: initialTab === TabLabel.ACTIVE
        },
        {
          label: TabLabel.FINISHED,
          value: sortOperations(
            autoOpsRules.filter(
              (rule) =>
                rule.autoOpsStatus === AutoOpsStatus.FINISHED ||
                rule.autoOpsStatus === AutoOpsStatus.STOPPED
            ),
            progressiveRollout.filter(
              (p) => !isProgressiveRolloutsRunningWaiting(p.status)
            ),
            SORT_TYPE.DESC
          ),
          selected: initialTab === TabLabel.FINISHED
        }
      ]);
    }, [autoOpsRules, progressiveRollout, setTabs]);

    useEffect(() => {
      if (autoOpsRules?.length > 0) {
        const ids = autoOpsRules
          .filter((rule) => {
            return (
              rule.opsType === OpsType.EVENT_RATE &&
              (rule.autoOpsStatus === AutoOpsStatus.WAITING ||
                rule.autoOpsStatus === AutoOpsStatus.RUNNING)
            );
          })
          .map((rule) => rule.id);

        if (ids.length > 0) {
          dispatch(
            listOpsCounts({
              environmentId: currentEnvironment.id,
              ids
            })
          );
        }
      }
    }, [autoOpsRules]);

    const handleCloseAutoOps = useCallback(() => {
      reset();
      history.replace({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_AUTOOPS}`,
        search: location.search
      });
      setSelectedAutoOpsRule(null);
      setSelectedOperationType(undefined);
    }, [history, location, reset]);

    const handleOpenAuto = useCallback(
      (operationType: OperationType) => {
        if (operationType === OperationType.SCHEDULE) {
          setValue('opsType', OpsType.SCHEDULE);
        } else if (operationType === OperationType.EVENT_RATE) {
          setValue('opsType', OpsType.EVENT_RATE);
        } else if (operationType === OperationType.PROGRESSIVE_ROLLOUT) {
          setValue(
            'progressiveRolloutType',
            ProgressiveRollout.Type.TEMPLATE_SCHEDULE
          );
        }
        setSelectedOperationType(operationType);
        history.push({
          pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_AUTOOPS}/${operationType}`,
          search: location.search
        });
      },
      [history, location, setSelectedOperationType]
    );

    const handleOpenUpdate = useCallback((rule: AutoOpsRule.AsObject) => {
      setSelectedAutoOpsRule(rule);
      if (rule.opsType === OpsType.SCHEDULE) {
        handleOpenAuto(OperationType.SCHEDULE);
      } else if (rule.opsType === OpsType.EVENT_RATE) {
        handleOpenAuto(OperationType.EVENT_RATE);
      }
    }, []);

    const handleOnSubmit = useCallback(() => {
      handleCloseAutoOps();
      refetchAutoOpsRules();
    }, []);

    const handleOnSubmitProgressiveRollout = useCallback(() => {
      handleCloseAutoOps();
      refetchProgressiveRollouts();
    }, []);

    const handleDelete = (operation: SelectedOperation) => {
      setIsDeleteConfirmDialogOpen(true);
      setSelectedOperation(operation);
    };

    const handleDeleteConfirm = () => {
      setIsDeleteConfirmDialogOpen(false);
      if (
        selectedOperation.type === OperationType.SCHEDULE ||
        selectedOperation.type === OperationType.EVENT_RATE
      ) {
        dispatch(
          deleteAutoOpsRule({
            environmentId: currentEnvironment.id,
            id: selectedOperation.id
          })
        ).then(() => {
          refetchAutoOpsRules();
        });
      } else if (selectedOperation.type === OperationType.PROGRESSIVE_ROLLOUT) {
        dispatch(
          deleteProgressiveRollout({
            environmentId: currentEnvironment.id,
            id: selectedOperation.id
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

    const handleStopConfirm = useCallback(() => {
      setIsStopConfirmDialogOpen(false);

      if (
        selectedOperation.type === OperationType.SCHEDULE ||
        selectedOperation.type === OperationType.EVENT_RATE
      ) {
        dispatch(
          stopAutoOpsRule({
            environmentId: currentEnvironment.id,
            id: selectedOperation.id
          })
        ).then(() => {
          refetchAutoOpsRules();
        });
      } else if (selectedOperation.type === OperationType.PROGRESSIVE_ROLLOUT) {
        dispatch(
          stopProgressiveRollout({
            environmentId: currentEnvironment.id,
            id: selectedOperation.id
          })
        ).then(() => {
          refetchProgressiveRollouts();
        });
      }
    }, [selectedOperation]);

    const variationOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: createVariationLabel(v)
      };
    });

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
              {f(messages.sideMenu.documentation)}
            </span>
          </a>
          <Popover className="relative flex">
            <Popover.Button>
              <div className="btn-submit space-x-2 items-center">
                <PlusIcon width={16} />
                <span>{f(messages.autoOps.newOperation)}</span>
                <ChevronDownIcon width={16} />
              </div>
            </Popover.Button>
            <Popover.Panel className="absolute min-w-max top-10 z-10 bg-white right-0 rounded-lg p-1 whitespace-nowrap shadow-md">
              <button
                onClick={() => handleOpenAuto(OperationType.SCHEDULE)}
                className="flex w-full space-x-2 px-3 py-1.5 items-center hover:bg-gray-100"
              >
                <CalendarIcon width={16} color="#94A3B8" />
                <span className="text-sm text-gray-500">
                  {f(messages.autoOps.schedule)}
                </span>
              </button>
              <button
                onClick={() => handleOpenAuto(OperationType.EVENT_RATE)}
                className="flex w-full space-x-2 px-3 py-1.5 items-center hover:bg-gray-100"
              >
                <RefreshSvg />
                <span className="text-sm text-gray-500">
                  {f(messages.autoOps.eventRate)}
                </span>
              </button>
              <button
                onClick={() =>
                  handleOpenAuto(OperationType.PROGRESSIVE_ROLLOUT)
                }
                className="flex w-full space-x-2 px-3 py-1.5 items-center hover:bg-gray-100"
              >
                <ArrowTrendingUpGray color="#94A3B8" />
                <span className="text-sm text-gray-500">
                  {f(messages.autoOps.progressiveRollout)}
                </span>
              </button>
            </Popover.Panel>
          </Popover>
        </div>
        <ActiveCompletedTabs tabs={tabs} setTabs={setTabs} />
        <AutoOpsInfos openAddOperation={handleOpenAuto} />
        {isAutoOpsRuleLoading || isProgressiveRolloutsLoading ? (
          <DetailSkeleton />
        ) : (
          <div className="space-y-6 py-6">
            {tabs
              .find((tab) => tab.selected)
              ?.value.map((op) => {
                const operation = op as AutoOpsRule.AsObject &
                  ProgressiveRollout.AsObject;

                if (
                  operation.opsType === OpsType.SCHEDULE ||
                  operation.opsType === OpsType.EVENT_RATE
                ) {
                  return (
                    <Operation
                      key={operation.id}
                      rule={operation}
                      isActiveTabSelected={isActiveTabSelected}
                      handleOpenUpdate={handleOpenUpdate}
                      handleDelete={handleDelete}
                      handleStop={handleStop}
                    />
                  );
                } else if (
                  operation.type ===
                    ProgressiveRollout.Type.TEMPLATE_SCHEDULE ||
                  operation.type === ProgressiveRollout.Type.MANUAL_SCHEDULE
                ) {
                  return (
                    <ProgressiveRolloutOperation
                      key={operation.id}
                      featureId={featureId}
                      isActiveTabSelected={isActiveTabSelected}
                      progressiveRollout={operation}
                      handleDelete={handleDelete}
                      handleStop={handleStop}
                    />
                  );
                }
              })}
          </div>
        )}
        {selectedOperationType === OperationType.SCHEDULE && (
          <Overlay
            open={selectedOperationType === OperationType.SCHEDULE}
            onClose={handleCloseAutoOps}
          >
            <ScheduleAddUpdateForm
              featureId={featureId}
              currentFlagStatus={feature.enabled}
              onSubmit={handleSubmit(handleOnSubmit)}
              onCancel={handleCloseAutoOps}
              autoOpsRule={selectedAutoOpsRule}
              isActiveTabSelected={isActiveTabSelected}
            />
          </Overlay>
        )}
        {selectedOperationType === OperationType.EVENT_RATE && (
          <Overlay
            open={selectedOperationType === OperationType.EVENT_RATE}
            onClose={handleCloseAutoOps}
          >
            <EventRateAddUpdateForm
              onCancel={handleCloseAutoOps}
              featureId={featureId}
              autoOpsRule={selectedAutoOpsRule}
              isActiveTabSelected={isActiveTabSelected}
              variationOptions={variationOptions}
              onSubmit={handleSubmit(handleOnSubmit)}
            />
          </Overlay>
        )}
        {selectedOperationType === OperationType.PROGRESSIVE_ROLLOUT && (
          <Overlay
            open={selectedOperationType === OperationType.PROGRESSIVE_ROLLOUT}
            onClose={handleCloseAutoOps}
          >
            <ProgressiveRolloutAddForm
              featureId={featureId}
              onCancel={handleCloseAutoOps}
              autoOpsRule={selectedAutoOpsRule}
              isActiveTabSelected={isActiveTabSelected}
              variationOptions={variationOptions}
              onSubmitProgressiveRollout={handleOnSubmitProgressiveRollout}
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
          <AutoOpsStopDialog
            selectedOperation={selectedOperation}
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
          selected: t.label === tabLabel
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
  openAddOperation: (operationType: OperationType) => void;
}

const AutoOpsInfos: FC<AutoOpsInfosProps> = memo(({ openAddOperation }) => (
  <div className="py-6">
    <p className="text-xl font-bold">
      {intl.formatMessage(messages.autoOps.infoBlocks.title)}
    </p>
    <div className="flex space-x-6 mt-6">
      {[
        {
          id: 1,
          title: intl.formatMessage(messages.autoOps.schedule),
          detail: intl.formatMessage(messages.autoOps.infoBlocks.scheduleInfo),
          bgColor: 'bg-purple-50',
          icon: <CalendarSvg />,
          onClick: () => {
            openAddOperation(OperationType.SCHEDULE);
          }
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
              <RefreshPinkSvg />
              <CrossSvg className="absolute right-[2px] bottom-[1px]" />
            </div>
          ),
          onClick: () => {
            openAddOperation(OperationType.EVENT_RATE);
          }
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
            openAddOperation(OperationType.PROGRESSIVE_ROLLOUT);
          }
        }
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
));

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
    handleStop
  }) => {
    const [feature] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError
    ]);

    const variationOptions = feature.variationsList.map((v) => {
      return {
        value: v.id,
        label: v.value
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
              type: OperationType.PROGRESSIVE_ROLLOUT,
              id: progressiveRollout.id
            })
          }
          stopRule={() =>
            handleStop({
              type: OperationType.PROGRESSIVE_ROLLOUT,
              id: progressiveRollout.id
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
              type: OperationType.PROGRESSIVE_ROLLOUT,
              id: progressiveRollout.id
            })
          }
          stopRule={() =>
            handleStop({
              type: OperationType.PROGRESSIVE_ROLLOUT,
              id: progressiveRollout.id
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
  handleStop: (arg: SelectedOperation) => void;
}

const Operation: FC<OperationProps> = memo(
  ({
    rule,
    isActiveTabSelected,
    handleOpenUpdate,
    handleDelete,
    handleStop
  }) => {
    const { formatMessage: f } = useIntl();
    const opsCounts = useSelector<AppState, OpsCount.AsObject[]>(
      (state) => selectAllOpsCounts(state.opsCounts),
      shallowEqual
    );

    const { opsType } = rule;
    return (
      <div className="rounded-xl shadow px-6 py-4 bg-white">
        <div className="flex justify-between py-4 border-b">
          <h3 className="font-bold text-xl text-gray-600">
            {rule.opsType === OpsType.SCHEDULE &&
              f(messages.autoOps.scheduleOperation)}
            {rule.opsType === OpsType.EVENT_RATE &&
              f(messages.autoOps.killSwitchOperation)}
          </h3>
          <div className="flex space-x-2 items-center">
            <div
              className={classNames(
                'py-[2px] px-2 rounded text-sm',
                opsType === OpsType.SCHEDULE && 'bg-[#EBF9ED] text-green-700',
                opsType === OpsType.EVENT_RATE && 'bg-[#EFECF5] text-primary'
              )}
            >
              {opsType === OpsType.SCHEDULE && f(messages.autoOps.schedule)}
              {opsType === OpsType.EVENT_RATE && f(messages.autoOps.eventRate)}
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
                        {opsType === OpsType.SCHEDULE &&
                          f(messages.autoOps.editSchedule)}
                        {opsType === OpsType.EVENT_RATE &&
                          f(messages.autoOps.editKillSwitch)}
                      </span>
                    </button>
                    <button
                      onClick={() => {
                        if (rule.opsType === OpsType.SCHEDULE) {
                          handleStop({
                            type: OperationType.SCHEDULE,
                            id: rule.id
                          });
                        } else if (rule.opsType === OpsType.EVENT_RATE) {
                          handleStop({
                            type: OperationType.EVENT_RATE,
                            id: rule.id
                          });
                        }
                      }}
                      className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100"
                    >
                      <BanIcon width={18} className="" />
                      <span className="text-sm">
                        {rule.opsType === OpsType.SCHEDULE &&
                          f(messages.autoOps.stopSchedule)}
                        {rule.opsType === OpsType.EVENT_RATE &&
                          f(messages.autoOps.stopKillSwitch)}
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
                <button
                  onClick={() => {
                    if (opsType === OpsType.SCHEDULE) {
                      handleDelete({
                        type: OperationType.SCHEDULE,
                        id: rule.id
                      });
                    } else if (opsType === OpsType.EVENT_RATE) {
                      handleDelete({
                        type: OperationType.EVENT_RATE,
                        id: rule.id
                      });
                    }
                  }}
                  className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100"
                >
                  <TrashIcon width={18} className="text-red-500" />
                  <span className="text-red-500 text-sm">
                    {opsType === OpsType.SCHEDULE &&
                      f(messages.autoOps.deleteSchedule)}
                    {opsType === OpsType.EVENT_RATE &&
                      f(messages.autoOps.deleteKillSwitch)}
                  </span>
                </button>
              </Popover.Panel>
            </Popover>
          </div>
        </div>
        <div className="mt-4">
          <div className="flex justify-between">
            <p className="font-bold text-lg text-gray-600">
              {f(messages.autoOps.progressInformation)}
            </p>
            {rule.autoOpsStatus === AutoOpsStatus.STOPPED && (
              <div className="flex items-center text-gray-500 space-x-[6px]">
                <ClockIcon width={18} className="" />
                <span>{f(messages.experiment.status.forceStopped)}</span>
                <RelativeDateText date={new Date(rule.updatedAt * 1000)} />
              </div>
            )}
          </div>
          {rule.opsType === OpsType.SCHEDULE && (
            <DateTimeOperation rule={rule} />
          )}
          {rule.opsType === OpsType.EVENT_RATE && (
            <EventRateOperation rule={rule} opsCounts={opsCounts} />
          )}
        </div>
      </div>
    );
  }
);

interface DateTimeOperationProps {
  rule: AutoOpsRule.AsObject;
}

const DateTimeOperation = memo(({ rule }: DateTimeOperationProps) => {
  const { formatMessage: f } = useIntl();

  const [page, setPage] = useState(0);
  const count = Math.ceil(rule.clausesList.length / 10);

  const paginatedClausesList = rule.clausesList.slice(
    page * 10,
    (page + 1) * 10
  );

  const _datetimeClause = (value) => {
    const datetimeClause = getDatetimeClause(value);

    const date = dayjs(new Date(datetimeClause.time * 1000)).format(
      'YYYY/MM/DD'
    );
    const time = dayjs(new Date(datetimeClause.time * 1000)).format('HH:mm');
    return {
      date,
      time
    };
  };

  const _isSameOrBeforeCurrentDate = (date) => {
    return dayjs(date).isSameOrBefore(new Date());
  };

  let displayTime;
  let displayLabel;
  if (page === 0) {
    displayTime =
      rule.updatedAt > rule.createdAt ? rule.updatedAt : rule.createdAt;
    displayLabel =
      rule.createdAt === rule.updatedAt
        ? f(messages.autoOps.created)
        : f(messages.autoOps.updated);
  } else {
    const clause = rule.clausesList[page * 10 - 1];
    displayTime = getDatetimeClause(clause.clause.value).time;
    displayLabel = actionTypesOptions.find(
      (o) => o.value === clause.actionType.toString()
    ).label;
  }

  const handlePageChange = (page) => {
    setPage(page);
  };

  return (
    <div>
      <div className="px-12 pb-16 pt-14">
        <div className="flex relative h-[4px]">
          <div className="flex items-center">
            <div
              className={classNames(
                'w-[9px] h-[9px] relative rounded-full border',
                _isSameOrBeforeCurrentDate(new Date(displayTime * 1000))
                  ? 'bg-pink-500 border-pink-500'
                  : 'bg-white border-gray-400'
              )}
            >
              <span className="absolute -top-8 left-1/2 -translate-x-1/2 whitespace-nowrap">
                {displayLabel}
              </span>
              <div className="text-xs text-gray-400 absolute space-y-[2px] left-1/2 -translate-x-1/2 whitespace-nowrap text-center top-[18px]">
                <p>{dayjs(new Date(displayTime * 1000)).format('HH:mm')}</p>
                <p>
                  {dayjs(new Date(displayTime * 1000)).format('YYYY/MM/DD')}
                </p>
              </div>
            </div>
          </div>
          {paginatedClausesList.map((clause) => {
            const time = getDatetimeClause(clause.clause.value).time;

            return (
              <div
                key={clause.id}
                className={classNames(
                  'flex flex-1 justify-end items-center relative',
                  _isSameOrBeforeCurrentDate(time * 1000)
                    ? 'bg-pink-500 border-pink-500'
                    : 'bg-gray-200'
                )}
              >
                <div
                  key={clause.id}
                  className={classNames(
                    'w-[9px] h-[9px] relative rounded-full border',
                    _isSameOrBeforeCurrentDate(new Date(time * 1000))
                      ? 'bg-pink-500 border-pink-500'
                      : 'bg-white border-gray-400'
                  )}
                >
                  <span className="absolute -top-8 left-1/2 -translate-x-1/2">
                    {
                      actionTypesOptions.find(
                        (o) => o.value === clause.actionType.toString()
                      ).label
                    }
                  </span>
                  <div className="text-xs text-gray-400 absolute space-y-[2px] left-1/2 -translate-x-1/2 whitespace-nowrap text-center top-[18px]">
                    <p>{_datetimeClause(clause.clause.value).time}</p>
                    <p>{_datetimeClause(clause.clause.value).date}</p>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
      <OperationPagination
        page={page}
        count={count}
        onPageChange={handlePageChange}
      />
    </div>
  );
});

interface EventRateOperationProps {
  rule: AutoOpsRule.AsObject;
  opsCounts: OpsCount.AsObject[];
}

const EventRateOperation = memo(
  ({ rule, opsCounts }: EventRateOperationProps) => {
    const { value } = rule.clausesList[0].clause;
    const { formatMessage: f } = useIntl();

    const { goalId, minCount, threadsholdRate } =
      OpsEventRateClause.deserializeBinary(value as Uint8Array).toObject();

    const opsCount = opsCounts.find(
      (opsCount) => opsCount.autoOpsRuleId === rule.id
    );

    let currentEventRate = 0;
    if (opsCount && opsCount.opsEventCount >= minCount) {
      let opsEventCount = opsCount.opsEventCount;
      // If opsEventCount is unexpectedly greater than evaluationCount,
      // clamp opsEventCount to evaluationCount so that it never exceeds 100%.
      // This way, it will avoid a rate above 100%, which breaks the admin console design.
      if (opsCount.opsEventCount > opsCount.evaluationCount) {
        opsEventCount = opsCount.evaluationCount;
      }
      currentEventRate =
        Math.round((opsEventCount / opsCount.evaluationCount) * 100 * 100) /
        100;
    }

    const numberOfSteps =
      Math.round(threadsholdRate * 100) > 10
        ? 10
        : Math.round(threadsholdRate * 100);
    const step = (threadsholdRate * 100) / numberOfSteps;

    const stepArray = Array.from({ length: numberOfSteps }, (_, index) =>
      Math.round(step + index * step)
    );

    const barWidth = (currentEventRate / (threadsholdRate * 100)) * 100;

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
                width: `${barWidth}%`
              }}
            />
            <div className={classNames('flex items-center h-[4px]')}>
              <div
                className={classNames(
                  'w-[9px] h-[9px] rounded-full relative',
                  currentEventRate > 0
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
              const isActive =
                percentage <= currentEventRate && currentEventRate !== 0;

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
    isActiveTabSelected
  }: ProgressiveRolloutTemplateScheduleProps) => {
    const { formatMessage: f } = useIntl();

    const [page, setPage] = useState(0);

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

    const handlePageChange = (page) => {
      setPage(page);
    };

    const count = Math.ceil(schedulesList.length / 10);

    const paginatedScheduleList = schedulesList.slice(
      page * 10,
      (page + 1) * 10
    );

    const firstSchedule = {
      weight: page === 0 ? 0 : schedulesList[page * 10 - 1].weight / 1000,
      executeAt:
        page === 0 ? rule.createdAt : schedulesList[page * 10 - 1].executeAt
    };

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
            <Popover className="relative flex">
              <Popover.Button>
                <div className="pl-2 flex items-center cursor-pointer">
                  <DotsHorizontalIcon width={20} />
                </div>
              </Popover.Button>
              <Popover.Panel className="absolute z-10 bg-white right-0 rounded-lg p-1 whitespace-nowrap shadow-md">
                {isActiveTabSelected && (
                  <button
                    onClick={stopRule}
                    className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100"
                  >
                    <BanIcon width={18} className="" />
                    <span className="text-sm">
                      {f(messages.autoOps.stopProgressiveRollout)}
                    </span>
                  </button>
                )}
                <button
                  onClick={deleteRule}
                  className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100"
                >
                  <TrashIcon width={18} className="text-red-500" />
                  <span className="text-red-500 text-sm">
                    {f(messages.autoOps.deleteProgressiveRollout)}
                  </span>
                </button>
              </Popover.Panel>
            </Popover>
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
              <div className="text-gray-500">
                {rule.stoppedBy === ProgressiveRollout.StoppedBy.USER && (
                  <div className="flex items-center">
                    {f(messages.autoOps.stoppedByUser, {
                      relativeDate: (
                        <div
                          className={classNames(
                            !isLanguageJapanese && 'mx-[6px]'
                          )}
                        >
                          <RelativeDateText
                            date={new Date(rule.stoppedAt * 1000)}
                          />
                        </div>
                      ),
                      stoppedByIcon: <UserSvg className="mx-[6px]" />,
                      clockIcon: <ClockIcon width={18} className="mx-[6px]" />
                    })}
                  </div>
                )}
                {rule.stoppedBy ===
                  ProgressiveRollout.StoppedBy.OPS_KILL_SWITCH && (
                  <div className="flex items-center">
                    {f(messages.autoOps.stoppedByKillSwitch, {
                      relativeDate: (
                        <div
                          className={classNames(
                            !isLanguageJapanese && 'mx-[6px]'
                          )}
                        >
                          <RelativeDateText
                            date={new Date(rule.stoppedAt * 1000)}
                          />
                        </div>
                      ),
                      stoppedByIcon: (
                        <div className="relative px-[6px]">
                          <RefreshPinkSvg width={22} />
                          <CrossSvg
                            width={12}
                            className="absolute right-[6px] bottom-[3px]"
                          />
                        </div>
                      ),
                      clockIcon: <ClockIcon width={18} className="mx-[6px]" />
                    })}
                  </div>
                )}
                {rule.stoppedBy ===
                  ProgressiveRollout.StoppedBy.OPS_SCHEDULE && (
                  <div className="flex items-center">
                    {f(messages.autoOps.stoppedBySchedule, {
                      relativeDate: (
                        <div
                          className={classNames(
                            !isLanguageJapanese && 'mx-[6px]'
                          )}
                        >
                          <RelativeDateText
                            date={new Date(rule.stoppedAt * 1000)}
                          />
                        </div>
                      ),
                      stoppedByIcon: (
                        <CalendarIcon
                          width={18}
                          className="text-primary mx-[6px]"
                        />
                      ),
                      clockIcon: <ClockIcon width={18} className="mx-[6px]" />
                    })}
                  </div>
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
          <OperationPagination
            page={page}
            count={count}
            onPageChange={handlePageChange}
          />
        </div>
      </div>
    );
  }
);
