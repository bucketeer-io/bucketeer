import { useCallback, useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { autoOpsDelete, autoOpsStop } from '@api/auto-ops';
import { rolloutDelete, rolloutStopped } from '@api/rollouts';
import { useQueryAutoOpsCount } from '@queries/auto-ops-count';
import {
  invalidateAutoOpsRules,
  useQueryAutoOpsRules
} from '@queries/auto-ops-rules';
import { invalidateRollouts, useQueryRollouts } from '@queries/rollouts';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import {
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { AutoOpsRule, Feature, Rollout, RuleStrategyVariation } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams, useSearchParams } from 'utils/search-params';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import FormLoading from 'elements/form-loading';
import CollectionLayout from './elements/collection-layout';
import OperationActions from './elements/operation-actions';
import {
  DeleteOperationModal,
  StopOperationModal
} from './elements/operation-modals/action-operation';
import EventRateOperationModal from './elements/operation-modals/event-rate';
import ProgressiveRolloutModal from './elements/operation-modals/rollout';
import RolloutCloneModal from './elements/operation-modals/rollout-clone';
import ScheduleOperationModal from './elements/operation-modals/schedule-operation';
import Overview from './elements/overview';
import { OperationActionType, OperationTab, OpsTypeMap } from './types';

export interface OperationModalState {
  operationType: OpsTypeMap | undefined;
  actionType: OperationActionType;
  selectedData?: Rollout | AutoOpsRule;
}

const Operations = ({
  feature,
  editable,
  refetchFeature
}: {
  feature: Feature;
  editable: boolean;
  refetchFeature: () => void;
}) => {
  const { t } = useTranslation(['common', 'table', 'form', 'message']);
  const navigate = useNavigate();
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { searchOptions, onChangSearchParams } = useSearchParams();

  const searchParams = stringifyParams(
    pickBy(searchOptions, v => isNotEmpty(v as string))
  );

  const getPathName = useCallback(
    (path?: string) =>
      `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${feature.id}${PAGE_PATH_FEATURE_AUTOOPS}${path}`,
    [currentEnvironment, feature]
  );

  const { id: action, onCloseActionModal } = useActionWithURL({
    closeModalPath: getPathName(searchParams ? `?${searchParams}` : '')
  });

  const [currentTab, setCurrentTab] = useState(OperationTab.ACTIVE);
  const [operationModalState, setOperationModalState] =
    useState<OperationModalState>({
      operationType: undefined,
      actionType: 'NEW',
      selectedData: undefined
    });

  const [isLoading, setIsLoading] = useState(false);

  const rolloutStrategyCount = useMemo(() => {
    return (
      feature.defaultStrategy?.rolloutStrategy?.variations.map(item => ({
        ...item,
        weight: item.weight / 1000,
        variation: feature.variations.find(
          variation => variation.id === item.variation
        )?.name
      })) || []
    );
  }, [feature]);
  const isScheduleAction = useMemo(() => action === 'schedule', [action]);
  const isEventRateAction = useMemo(() => action === 'event-rate', [action]);
  const isRolloutAction = useMemo(() => action === 'rollout', [action]);

  const isRolloutActive = useMemo(() => {
    const data = operationModalState.selectedData;
    if (!data) return false;
    // Rollout has 'status', AutoOpsRule has 'autoOpsStatus'
    const status = 'status' in data ? data.status : data.autoOpsStatus;
    return status === 'RUNNING' || status === 'WAITING';
  }, [operationModalState]);

  const isScheduleType = useMemo(
    () => operationModalState.operationType === 'SCHEDULE',
    [operationModalState]
  );
  const isRolloutType = useMemo(
    () => operationModalState.operationType === 'ROLLOUT',
    [operationModalState]
  );
  const isOpenModalAction = useMemo(
    () => ['NEW', 'UPDATE', 'DETAILS'].includes(operationModalState.actionType),
    [operationModalState]
  );
  const isStop = useMemo(
    () => operationModalState.actionType === 'STOP',
    [operationModalState]
  );
  const isDelete = useMemo(
    () => operationModalState.actionType === 'DELETE',
    [operationModalState]
  );

  const queryParams = useMemo(
    () => ({
      cursor: String(0),
      featureIds: [feature.id],
      environmentId: currentEnvironment.id
    }),
    [feature, currentEnvironment]
  );

  const { data: rolloutCollection, isLoading: isRolloutLoading } =
    useQueryRollouts({
      params: queryParams
    });

  const { data: operationCollection, isLoading: isOperationLoading } =
    useQueryAutoOpsRules({
      params: queryParams
    });

  const rollouts = rolloutCollection?.progressiveRollouts || [];
  const operations = operationCollection?.autoOpsRules || [];

  const eventRateActiveIds = operations
    ?.filter(
      item =>
        ['RUNNING', 'WAITING'].includes(item.autoOpsStatus) &&
        item.opsType === 'EVENT_RATE'
    )
    ?.map(item => item.id);

  const { data: opsCountCollection } = useQueryAutoOpsCount({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id,
      featureIds: [feature.id],
      autoOpsRuleIds: eventRateActiveIds
    },
    enabled: !!eventRateActiveIds.length
  });

  const opsCounts = opsCountCollection?.opsCounts || [];

  const onOpenOperationModal = useCallback(
    (path: string) => {
      navigate(getPathName(`${path}${searchParams ? `?${searchParams}` : ''}`));
    },
    [searchParams]
  );

  const onSubmitOperationSuccess = useCallback(() => {
    invalidateAutoOpsRules(queryClient);
    invalidateRollouts(queryClient);
    // Auto navigate to ACTIVE tab when creating operation from FINISHED tab
    if (currentTab === OperationTab.FINISHED) {
      setCurrentTab(OperationTab.ACTIVE);
      // Update search params to ACTIVE tab before closing modal
      const updatedSearchOptions = {
        ...searchOptions,
        tab: OperationTab.ACTIVE
      };
      const updatedSearchParams = stringifyParams(
        pickBy(updatedSearchOptions, v => isNotEmpty(v as string))
      );
      navigate(
        getPathName(updatedSearchParams ? `?${updatedSearchParams}` : '')
      );
    } else {
      onCloseActionModal();
    }
  }, [
    searchParams,
    currentTab,
    searchOptions,
    navigate,
    getPathName,
    queryClient
  ]);

  const onOperationActions = useCallback(
    ({
      operationType,
      actionType = 'NEW',
      selectedData
    }: OperationModalState) => {
      setOperationModalState({
        operationType,
        actionType,
        selectedData
      });

      if (!['NEW', 'UPDATE', 'DETAILS'].includes(actionType)) return;
      if (operationType === OpsTypeMap.SCHEDULE)
        return onOpenOperationModal('/schedule');

      if (operationType === OpsTypeMap.EVENT_RATE)
        return onOpenOperationModal('/event-rate');

      if (operationType === OpsTypeMap.ROLLOUT)
        return onOpenOperationModal('/rollout');
    },
    [searchParams]
  );

  const onStopOperation = useCallback(async () => {
    try {
      if (operationModalState?.selectedData) {
        setIsLoading(true);
        let resp = null;
        const isStopRollout =
          operationModalState.operationType === OpsTypeMap.ROLLOUT;
        if (isStopRollout) {
          resp = await rolloutStopped({
            environmentId: currentEnvironment.id,
            id: operationModalState?.selectedData?.id,
            stoppedBy: 'USER'
          });
        } else {
          resp = await autoOpsStop({
            environmentId: currentEnvironment.id,
            id: operationModalState?.selectedData?.id
          });
        }

        if (resp) {
          notify({
            message: t('message:collection-action-success', {
              collection: t('operation'),
              action: t('stopped')
            })
          });
          onSubmitOperationSuccess();
          setOperationModalState({
            operationType: undefined,
            actionType: 'NEW',
            selectedData: undefined
          });
        }
      }
    } catch (error) {
      errorNotify(error);
    } finally {
      setIsLoading(false);
    }
  }, [operationModalState]);

  const onDeleteOperation = useCallback(async () => {
    try {
      if (operationModalState?.selectedData) {
        setIsLoading(true);
        const isStopRollout =
          operationModalState.operationType === OpsTypeMap.ROLLOUT;
        const deleteFn = isStopRollout ? rolloutDelete : autoOpsDelete;
        const resp = await deleteFn({
          environmentId: currentEnvironment.id,
          id: operationModalState?.selectedData?.id
        });

        if (resp) {
          notify({
            message: t('message:collection-action-success', {
              collection: t('operation'),
              action: t('deleted')
            })
          });
          onSubmitOperationSuccess();
          onResetModalState();
        }
      }
    } catch (error) {
      errorNotify(error);
    } finally {
      setIsLoading(false);
    }
  }, [operationModalState]);

  const onResetModalState = useCallback(
    () =>
      setOperationModalState({
        operationType: undefined,
        actionType: 'NEW',
        selectedData: undefined
      }),
    []
  );

  useEffect(() => {
    const tab = (searchOptions?.tab || OperationTab.ACTIVE) as OperationTab;
    onChangSearchParams({
      tab
    });
    setCurrentTab(tab);
    refetchFeature();
  }, [searchOptions]);

  return (
    <div className="flex flex-col w-full gap-y-4 min-w-[900px]">
      <div className="flex flex-wrap items-center justify-between w-full gap-6 px-6">
        <p className="flex flex-1 typo-head-bold-big text-gray-800 xl:whitespace-nowrap">
          {t('table:feature-flags:operations-desc')}
        </p>
        <Filter
          action={
            <OperationActions
              disabled={!editable}
              onOperationActions={onOperationActions}
            />
          }
          className="justify-end w-fit px-0"
          link={DOCUMENTATION_LINKS.FLAG_OPERATION}
        />
      </div>
      <Overview
        disabled={!editable}
        onOperationActions={operationType =>
          onOperationActions({
            operationType,
            actionType: 'NEW'
          })
        }
      />
      {isRolloutLoading || isOperationLoading ? (
        <FormLoading />
      ) : (
        <Tabs
          className="flex-1 flex h-full flex-col"
          value={currentTab}
          onValueChange={value => {
            const tab = value as OperationTab;
            setCurrentTab(tab);
            onChangSearchParams({ tab });
          }}
        >
          <TabsList className="px-6">
            <TabsTrigger value={OperationTab.ACTIVE}>{t(`active`)}</TabsTrigger>
            <TabsTrigger value={OperationTab.FINISHED}>
              {t(`finished`)}
            </TabsTrigger>
          </TabsList>

          <TabsContent value={currentTab} className="px-6">
            <CollectionLayout
              rolloutStrategyCount={
                rolloutStrategyCount as RuleStrategyVariation[]
              }
              currentTab={currentTab}
              operations={operations}
              opsCounts={opsCounts}
              rollouts={rollouts}
              onOperationActions={onOperationActions}
            />
          </TabsContent>
        </Tabs>
      )}
      {isScheduleAction && isOpenModalAction && feature && (
        <ScheduleOperationModal
          editable={editable}
          isFinishedTab={currentTab === OperationTab.FINISHED}
          isOpen={isScheduleAction}
          featureId={feature.id}
          environmentId={currentEnvironment.id}
          isEnabledFlag={feature.enabled}
          rollouts={rollouts}
          actionType={operationModalState.actionType}
          selectedData={operationModalState?.selectedData as AutoOpsRule}
          onClose={onCloseActionModal}
          onSubmitOperationSuccess={onSubmitOperationSuccess}
        />
      )}
      {isEventRateAction && isOpenModalAction && feature && (
        <EventRateOperationModal
          editable={editable}
          isOpen={isEventRateAction}
          feature={feature}
          environmentId={currentEnvironment.id}
          actionType={operationModalState.actionType}
          isFinishedTab={currentTab === OperationTab.FINISHED}
          selectedData={operationModalState?.selectedData as AutoOpsRule}
          onClose={onCloseActionModal}
          onSubmitOperationSuccess={onSubmitOperationSuccess}
        />
      )}
      {isRolloutAction &&
        operationModalState.actionType === 'NEW' &&
        feature && (
          <ProgressiveRolloutModal
            editable={editable}
            isOpen={isRolloutAction}
            feature={feature}
            urlCode={currentEnvironment.urlCode}
            environmentId={currentEnvironment.id}
            actionType={operationModalState.actionType}
            selectedData={operationModalState?.selectedData as Rollout}
            rollouts={rollouts}
            onClose={onCloseActionModal}
            onSubmitRolloutSuccess={onSubmitOperationSuccess}
          />
        )}
      {isRolloutAction &&
        operationModalState?.selectedData &&
        operationModalState.actionType === 'DETAILS' &&
        feature && (
          <RolloutCloneModal
            isOpen={
              isRolloutAction && operationModalState.actionType === 'DETAILS'
            }
            selectedData={operationModalState?.selectedData as Rollout}
            onClose={onCloseActionModal}
          />
        )}

      {isStop && !!operationModalState?.selectedData && (
        <StopOperationModal
          environment={currentEnvironment}
          isRunning={isRolloutActive}
          feature={feature}
          editable={editable}
          loading={isLoading}
          operationType={operationModalState.operationType!}
          isOpen={isStop && !!operationModalState?.selectedData}
          refetchFeatures={refetchFeature}
          onClose={onResetModalState}
          onSubmit={onStopOperation}
        />
      )}
      {isDelete && !!operationModalState?.selectedData && (
        <DeleteOperationModal
          isRunning={isRolloutActive}
          loading={isLoading}
          isRolloutType={isRolloutType}
          isScheduleType={isScheduleType}
          editable={editable}
          feature={feature}
          environment={currentEnvironment}
          isOpen={isDelete && !!operationModalState?.selectedData}
          refetchFeature={refetchFeature}
          onClose={onResetModalState}
          onSubmit={onDeleteOperation}
        />
      )}
    </div>
  );
};

export default Operations;
