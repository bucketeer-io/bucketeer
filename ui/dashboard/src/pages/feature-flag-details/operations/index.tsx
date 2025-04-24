import { useCallback, useEffect, useMemo, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { autoOpsStop } from '@api/auto-ops';
import { rolloutStopped } from '@api/rollouts';
import { useQueryAutoOpsRules } from '@queries/auto-ops-rules';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import { AutoOpsRule, Feature, Rollout } from '@types';
import { useSearchParams } from 'utils/search-params';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import FormLoading from 'elements/form-loading';
import { OperationActionType } from '../types';
import CollectionLayout from './elements/collection-layout';
import OperationActions from './elements/operation-actions';
import EventRateOperationModal from './elements/operation-modals/event-rate';
import ProgressiveRolloutModal from './elements/operation-modals/rollout';
import ScheduleOperationModal from './elements/operation-modals/schedule-operation';
import StopOperationModal from './elements/operation-modals/stop-operation';
import { OperationTab, OpsTypeMap } from './types';

export interface OperationModalState {
  operationType: OpsTypeMap | undefined;
  actionType: OperationActionType;
  selectedData?: AutoOpsRule | Rollout;
}

const Operations = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['common', 'table', 'form']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const location = useLocation();
  const navigate = useNavigate();
  const { notify, errorNotify } = useToast();
  const getPathName = useCallback(
    (path?: string) =>
      `/${currentEnvironment.urlCode}/features/${feature.id}/autoops${path}`,
    [currentEnvironment, feature]
  );

  const { id: action, onCloseActionModal } = useActionWithURL({
    closeModalPath: getPathName(location.search)
  });

  const [currentTab, setCurrentTab] = useState(OperationTab.ACTIVE);
  const [operationModalState, setOperationModalState] =
    useState<OperationModalState>({
      operationType: undefined,
      actionType: 'NEW',
      selectedData: undefined
    });

  const [isLoading, setIsLoading] = useState(false);

  const isSchedule = useMemo(() => action === 'schedule', [action]);
  const isEventRate = useMemo(() => action === 'event-rate', [action]);
  const isRollout = useMemo(() => action === 'rollout', [action]);
  const isCreateOrUpdate = useMemo(
    () => ['NEW', 'UPDATE'].includes(operationModalState.actionType),
    [operationModalState]
  );
  const isStop = useMemo(
    () => operationModalState.actionType === 'STOP',
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

  const {
    data: rolloutCollection,
    isLoading: isRolloutLoading,
    refetch: refetchRollouts
  } = useQueryRollouts({
    params: queryParams
  });

  const {
    data: operationCollection,
    isLoading: isOperationLoading,
    refetch: refetchAutoOpsRules
  } = useQueryAutoOpsRules({
    params: queryParams
  });

  const rollouts = rolloutCollection?.progressiveRollouts || [];
  const operations = operationCollection?.autoOpsRules || [];

  const onOpenOperationModal = useCallback(
    (path: string) => {
      navigate(getPathName(`${path}${location.search}`));
    },
    [location]
  );

  const onSubmitOperationSuccess = useCallback(() => {
    onCloseActionModal();
    refetchAutoOpsRules();
  }, []);

  const onSubmitRolloutSuccess = useCallback(() => {
    onCloseActionModal();
    refetchRollouts();
  }, []);

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
      if (!['NEW', 'UPDATE'].includes(actionType)) return;
      if (operationType === OpsTypeMap.SCHEDULE)
        return onOpenOperationModal('/schedule');

      if (operationType === OpsTypeMap.EVENT_RATE)
        return onOpenOperationModal('/event-rate');

      if (operationType === OpsTypeMap.ROLLOUT)
        return onOpenOperationModal('/rollout');
    },
    []
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
            message: 'Stopped operation successfully.'
          });
          refetchRollouts();
          refetchAutoOpsRules();
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

  useEffect(() => {
    const tab = (searchOptions?.tab || OperationTab.ACTIVE) as OperationTab;
    onChangSearchParams({
      tab
    });
    setCurrentTab(tab);
  }, [searchOptions]);

  return (
    <div className="flex flex-col w-full gap-y-4 min-w-[900px]">
      <Filter
        searchValue=""
        isShowDocumentation={false}
        onSearchChange={() => {}}
        onOpenFilter={() => {}}
        action={<OperationActions onOperationActions={onOperationActions} />}
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
            <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
            <TabsTrigger value="COMPLETED">{t(`completed`)}</TabsTrigger>
          </TabsList>

          <TabsContent value={currentTab} className="px-6">
            <CollectionLayout
              currentTab={currentTab}
              rollouts={rollouts}
              operations={operations}
              onOperationActions={onOperationActions}
            />
          </TabsContent>
        </Tabs>
      )}
      {isSchedule && isCreateOrUpdate && feature && (
        <ScheduleOperationModal
          isOpen={isSchedule}
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
      {isEventRate && isCreateOrUpdate && feature && (
        <EventRateOperationModal
          isOpen={isEventRate}
          feature={feature}
          environmentId={currentEnvironment.id}
          actionType={operationModalState.actionType}
          selectedData={operationModalState?.selectedData as AutoOpsRule}
          onClose={onCloseActionModal}
          onSubmitOperationSuccess={onSubmitOperationSuccess}
        />
      )}
      {isRollout && isCreateOrUpdate && feature && (
        <ProgressiveRolloutModal
          isOpen={isRollout}
          feature={feature}
          environmentId={currentEnvironment.id}
          actionType={operationModalState.actionType}
          selectedData={operationModalState?.selectedData as Rollout}
          rollouts={rollouts}
          onClose={onCloseActionModal}
          onSubmitRolloutSuccess={onSubmitRolloutSuccess}
        />
      )}
      {isStop && !!operationModalState?.selectedData && (
        <StopOperationModal
          loading={isLoading}
          operationType={operationModalState.operationType!}
          isOpen={isStop && !!operationModalState?.selectedData}
          onClose={() =>
            setOperationModalState({
              operationType: undefined,
              actionType: 'NEW',
              selectedData: undefined
            })
          }
          onSubmit={onStopOperation}
        />
      )}
    </div>
  );
};

export default Operations;
