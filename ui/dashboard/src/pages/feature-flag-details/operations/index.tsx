import { useCallback, useEffect, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { autoOpsDelete, autoOpsStop } from '@api/auto-ops';
import { rolloutDelete, rolloutStopped } from '@api/rollouts';
import { useQueryAutoOpsCount } from '@queries/auto-ops-count';
import { useQueryAutoOpsRules } from '@queries/auto-ops-rules';
import { useQueryRollouts } from '@queries/rollouts';
import { useLoaderData } from '@tanstack/react-router';
import { getCurrentEnvironment, useAuth } from 'auth';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import {
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import { pickBy } from 'lodash';
import {
  Route as featureDetailsLayoutRoute,
  FeatureDetailsLoaderData
} from 'routes/_default-layout/$env/features/$featureId/_feature-details-layout';
import { AutoOpsRule, Rollout } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams, useSearchParams } from 'utils/search-params';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import ConfirmModal from 'elements/confirm-modal';
import Filter from 'elements/filter';
import FormLoading from 'elements/form-loading';
import CollectionLayout from './elements/collection-layout';
import OperationActions from './elements/operation-actions';
import EventRateOperationModal from './elements/operation-modals/event-rate';
import ProgressiveRolloutModal from './elements/operation-modals/rollout';
import RolloutCloneModal from './elements/operation-modals/rollout-clone';
import ScheduleOperationModal from './elements/operation-modals/schedule-operation';
import StopOperationModal from './elements/operation-modals/stop-operation';
import Overview from './elements/overview';
import { OperationActionType, OperationTab, OpsTypeMap } from './types';

export interface OperationModalState {
  operationType: OpsTypeMap | undefined;
  actionType: OperationActionType;
  selectedData?: AutoOpsRule | Rollout;
}

const Operations = () => {
  const { t } = useTranslation(['common', 'table', 'form', 'message']);
  const navigate = useNavigate();
  const { notify, errorNotify } = useToast();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { searchOptions, onChangSearchParams } = useSearchParams();

  const loaderData: FeatureDetailsLoaderData = useLoaderData({
    from: featureDetailsLayoutRoute.id
  });

  const feature = loaderData?.feature;

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

  const { actionType, operationType, selectedData } = operationModalState || {};

  const isScheduleAction = useMemo(() => action === 'schedule', [action]);
  const isEventRateAction = useMemo(() => action === 'event-rate', [action]);
  const isRolloutAction = useMemo(() => action === 'rollout', [action]);

  const isScheduleType = useMemo(
    () => operationType === 'SCHEDULE',
    [operationType]
  );
  const isRolloutType = useMemo(
    () => operationType === 'ROLLOUT',
    [operationType]
  );
  const isOpenModalAction = useMemo(
    () => ['NEW', 'UPDATE', 'DETAILS'].includes(actionType),
    [actionType]
  );
  const isStop = useMemo(() => actionType === 'STOP', [actionType]);
  const isDelete = useMemo(() => actionType === 'DELETE', [actionType]);

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
    onCloseActionModal();
    refetchAutoOpsRules();
  }, [searchParams]);

  const onSubmitRolloutSuccess = useCallback(() => {
    onCloseActionModal();
    refetchRollouts();
  }, [searchParams]);

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
      if (selectedData) {
        setIsLoading(true);
        let resp = null;
        const isStopRollout = operationType === OpsTypeMap.ROLLOUT;
        if (isStopRollout) {
          resp = await rolloutStopped({
            environmentId: currentEnvironment.id,
            id: selectedData?.id,
            stoppedBy: 'USER'
          });
        } else {
          resp = await autoOpsStop({
            environmentId: currentEnvironment.id,
            id: selectedData?.id
          });
        }

        if (resp) {
          notify({
            message: t('message:operation.stopped')
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
  }, [selectedData, operationType, currentEnvironment]);

  const onDeleteOperation = useCallback(async () => {
    try {
      if (selectedData) {
        setIsLoading(true);
        const isStopRollout = operationType === OpsTypeMap.ROLLOUT;
        const deleteFn = isStopRollout ? rolloutDelete : autoOpsDelete;
        const resp = await deleteFn({
          environmentId: currentEnvironment.id,
          id: selectedData?.id
        });

        if (resp) {
          notify({
            message: t('message:operation.deleted')
          });
          refetchAutoOpsRules();
          refetchRollouts();
          onResetModalState();
        }
      }
    } catch (error) {
      errorNotify(error);
    } finally {
      setIsLoading(false);
    }
  }, [selectedData, operationType, currentEnvironment]);

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
  }, [searchOptions]);

  return (
    <div className="flex flex-col w-full gap-y-4 min-w-[900px]">
      <div className="flex flex-wrap items-center justify-between w-full gap-6 px-6">
        <p className="flex flex-1 typo-head-bold-big text-gray-800 xl:whitespace-nowrap">
          {t('table:feature-flags:operations-desc')}
        </p>
        <Filter
          action={<OperationActions onOperationActions={onOperationActions} />}
          className="justify-end w-fit px-0"
          link={DOCUMENTATION_LINKS.FLAG_OPERATION}
        />
      </div>
      <Overview
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
          isFinishedTab={currentTab === OperationTab.FINISHED}
          isOpen={isScheduleAction}
          featureId={feature.id}
          environmentId={currentEnvironment.id}
          isEnabledFlag={feature.enabled}
          rollouts={rollouts}
          actionType={actionType}
          selectedData={selectedData as AutoOpsRule}
          onClose={onCloseActionModal}
          onSubmitOperationSuccess={onSubmitOperationSuccess}
        />
      )}
      {isEventRateAction && isOpenModalAction && feature && (
        <EventRateOperationModal
          isOpen={isEventRateAction}
          feature={feature}
          environmentId={currentEnvironment.id}
          actionType={actionType}
          isFinishedTab={currentTab === OperationTab.FINISHED}
          selectedData={selectedData as AutoOpsRule}
          onClose={onCloseActionModal}
          onSubmitOperationSuccess={onSubmitOperationSuccess}
        />
      )}
      {isRolloutAction && actionType === 'NEW' && feature && (
        <ProgressiveRolloutModal
          isOpen={isRolloutAction}
          feature={feature}
          urlCode={currentEnvironment.urlCode}
          environmentId={currentEnvironment.id}
          actionType={actionType}
          selectedData={selectedData as Rollout}
          rollouts={rollouts}
          onClose={onCloseActionModal}
          onSubmitRolloutSuccess={onSubmitRolloutSuccess}
        />
      )}
      {isRolloutAction &&
        selectedData &&
        actionType === 'DETAILS' &&
        feature && (
          <RolloutCloneModal
            isOpen={isRolloutAction && actionType === 'DETAILS'}
            selectedData={selectedData as Rollout}
            onClose={onCloseActionModal}
          />
        )}

      {isStop && !!selectedData && (
        <StopOperationModal
          loading={isLoading}
          operationType={operationType!}
          isOpen={isStop && !!selectedData}
          onClose={onResetModalState}
          onSubmit={onStopOperation}
        />
      )}
      {isDelete && !!selectedData && (
        <ConfirmModal
          loading={isLoading}
          isOpen={isDelete && !!selectedData}
          title={t(
            `table:popover.delete-${isRolloutType ? 'rollout' : isScheduleType ? 'operation' : 'kill-switch'}`
          )}
          description={
            <Trans
              i18nKey={'table:operations.confirm-delete-operation'}
              values={{
                type: t(
                  `form:feature-flags.${isRolloutType ? 'rollout' : isScheduleType ? 'schedule' : 'kill-switch'}`
                )
              }}
              components={{
                bold: <strong />
              }}
            />
          }
          onClose={onResetModalState}
          onSubmit={onDeleteOperation}
        />
      )}
    </div>
  );
};

export default Operations;
