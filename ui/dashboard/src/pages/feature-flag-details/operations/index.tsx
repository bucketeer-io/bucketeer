import { useCallback, useEffect, useMemo, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useQueryAutoOpsRules } from '@queries/auto-ops-rules';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import { AutoOpsRule, Feature, OpsEventRateClause, Rollout } from '@types';
import { useSearchParams } from 'utils/search-params';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import FormLoading from 'elements/form-loading';
import TableListContainer from 'elements/table-list-container';
import { OperationActionType } from '../types';
import CollectionLayout from './elements/collection-layout';
import OperationActions from './elements/operation-actions';
import EventRateOperationModal from './elements/operation-modals/event-rate';
import ScheduleOperationModal from './elements/operation-modals/schedule-operation';
import { OperationTab, OpsTypeMap } from './types';

export interface OperationModalState {
  operationType: OpsTypeMap | undefined;
  actionType: OperationActionType;
  selectedData?: AutoOpsRule | Rollout | OpsEventRateClause;
}

const Operations = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const location = useLocation();
  const navigate = useNavigate();

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

  const isSchedule = useMemo(() => action === 'schedule', [action]);
  const isEventRate = useMemo(() => action === 'event-rate', [action]);

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

  // const defaultValues = {
  //   opsType: OpsTypeMap.SCHEDULE,
  //   datetimeClausesList: [createDatetimeClausesList()],
  //   eventRate: createEventRate(feature),
  //   progressiveRolloutType: RolloutTypeMap.TEMPLATE_SCHEDULE,
  //   progressiveRollout: createProgressiveRollout(feature)
  // };

  // const form = useForm({
  //   resolver: yupResolver(operationFormSchema),
  //   defaultValues,
  //   mode: 'onChange'
  // });

  const onSubmitOperationSuccess = useCallback(() => {
    onCloseActionModal();
    refetchAutoOpsRules();
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
      if (operationType === OpsTypeMap.SCHEDULE)
        return onOpenOperationModal('/schedule');

      if (operationType === OpsTypeMap.EVENT_RATE)
        return onOpenOperationModal('/event-rate');
    },
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
    <div className="flex flex-col w-full gap-y-6">
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
          <TabsList>
            <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
            <TabsTrigger value="COMPLETED">{t(`completed`)}</TabsTrigger>
          </TabsList>

          <TabsContent value={currentTab} className="mt-0">
            <TableListContainer>
              <CollectionLayout
                feature={feature}
                currentTab={currentTab}
                rollouts={rollouts}
                operations={operations}
                onOperationActions={onOperationActions}
              />
            </TableListContainer>
          </TabsContent>
        </Tabs>
      )}
      {isSchedule && feature && (
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
      {isEventRate && feature && (
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
    </div>
  );
};

export default Operations;
