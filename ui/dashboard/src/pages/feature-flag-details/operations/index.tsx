import { useCallback, useEffect, useMemo, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { isEmptyObject } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Form from 'components/form';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import CollectionLoader from './elements/collection-loader';
import OperationActions from './elements/operation-actions';
import NewScheduleOperationModal from './elements/operation-modals/new-schedule-operation';
import { operationFormSchema } from './form-schema';
import { OperationTab, OpsTypeMap, RolloutTypeMap } from './types';
import {
  createDatetimeClausesList,
  createEventRate,
  createProgressiveRollout
} from './utils';

export type OperationModalActionType = 'NEW' | 'UPDATE';

interface OperationModalState {
  operationType?: OpsTypeMap;
  actionType: OperationModalActionType;
}

const Operations = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['common', 'table']);
  const { searchOptions, onChangSearchParams } = useSearchParams();

  const [currentTab, setCurrentTab] = useState(OperationTab.ACTIVE);
  const [operationModalState, setOperationModalState] =
    useState<OperationModalState>({
      operationType: undefined,
      actionType: 'NEW'
    });

  const isNewAction = useMemo(
    () => operationModalState.actionType === 'NEW',
    [operationModalState]
  );

  const isOpenNewScheduleModal = useMemo(
    () =>
      operationModalState?.operationType === OpsTypeMap.SCHEDULE && isNewAction,
    [isNewAction, operationModalState]
  );

  const defaultValues = {
    opsType: OpsTypeMap.SCHEDULE,
    datetimeClausesList: [createDatetimeClausesList()],
    eventRate: createEventRate(feature),
    progressiveRolloutType: RolloutTypeMap.TEMPLATE_SCHEDULE,
    progressiveRollout: createProgressiveRollout(feature)
  };

  const form = useForm({
    resolver: yupResolver(operationFormSchema),
    defaultValues,
    mode: 'onChange'
  });

  const onSubmit = () => {};

  const onOpenOperationModal = useCallback(
    (
      operationType?: OpsTypeMap,
      actionType: OperationModalActionType = 'NEW'
    ) =>
      setOperationModalState({
        operationType,
        actionType
      }),
    []
  );

  const onCloseOperationModal = useCallback(
    () =>
      setOperationModalState({
        ...operationModalState,
        operationType: undefined
      }),
    [operationModalState]
  );

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      const tab = searchOptions?.tab;
      setCurrentTab((tab as OperationTab) || OperationTab.ACTIVE);
    }
  }, [searchOptions]);

  return (
    <FormProvider {...form}>
      <Form onSubmit={form.handleSubmit(onSubmit)}>
        <div className="flex flex-col w-full gap-y-6">
          <Filter
            searchValue=""
            isShowDocumentation={false}
            onSearchChange={() => {}}
            onOpenFilter={() => {}}
            action={
              <OperationActions onOpenOperationModal={onOpenOperationModal} />
            }
          />

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

            <TabsContent value={currentTab}>
              <CollectionLoader feature={feature} currentTab={currentTab} />
            </TabsContent>
          </Tabs>
          {isOpenNewScheduleModal && (
            <NewScheduleOperationModal
              isOpen={isOpenNewScheduleModal}
              isEnabledFlag={feature.enabled}
              onClose={() => {
                onCloseOperationModal();
                form.resetField('datetimeClausesList', {
                  defaultValue: [createDatetimeClausesList()]
                });
              }}
            />
          )}
        </div>
      </Form>
    </FormProvider>
  );
};

export default Operations;
