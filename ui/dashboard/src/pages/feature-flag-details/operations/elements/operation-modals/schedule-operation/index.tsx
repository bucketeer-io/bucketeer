import { useCallback, useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import {
  autoOpsCreator,
  AutoOpsCreatorResponse,
  autoOpsUpdate,
  ClauseUpdateType
} from '@api/auto-ops';
import { yupResolver } from '@hookform/resolvers/yup';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { AutoOpsRule, DatetimeClause, Rollout } from '@types';
import {
  dateTimeClauseListSchema,
  DateTimeClauseListType
} from 'pages/feature-flag-details/operations/form-schema';
import { ActionTypeMap } from 'pages/feature-flag-details/operations/types';
import { createDatetimeClausesList } from 'pages/feature-flag-details/operations/utils';
import { OperationActionType } from 'pages/feature-flag-details/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Form from 'components/form';
import SlideModal from 'components/modal/slide';
import ScheduleList from './schedule-list';

export interface OperationModalProps {
  isCompletedTab: boolean;
  featureId: string;
  environmentId: string;
  isOpen: boolean;
  isEnabledFlag: boolean;
  rollouts: Rollout[];
  actionType: OperationActionType;
  selectedData?: AutoOpsRule;
  onClose: () => void;
  onSubmitOperationSuccess: () => void;
}

const ScheduleOperationModal = ({
  isCompletedTab,
  featureId,
  environmentId,
  isOpen,
  isEnabledFlag,
  rollouts,
  actionType,
  selectedData,
  onClose,
  onSubmitOperationSuccess
}: OperationModalProps) => {
  const { t } = useTranslation(['form', 'common']);
  const { notify, errorNotify } = useToast();

  const isCreate = useMemo(() => actionType === 'NEW', [actionType]);

  const handleCreateDefaultValues = () => {
    if (selectedData) {
      return selectedData.clauses.map(item => {
        const time = new Date(+(item.clause as DatetimeClause).time * 1000);
        return {
          id: item.id,
          actionType: item.actionType as ActionTypeMap,
          time
        };
      });
    }
    return [createDatetimeClausesList()];
  };

  const form = useForm({
    resolver: yupResolver(dateTimeClauseListSchema),
    defaultValues: {
      datetimeClausesList: handleCreateDefaultValues()
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isSubmitting }
  } = form;

  const onSubmit = useCallback(
    async (values: DateTimeClauseListType) => {
      try {
        const { datetimeClausesList } = values;

        let resp: AutoOpsCreatorResponse | null = null;

        if (!isCreate && selectedData) {
          const updateDatetimeClauses: ClauseUpdateType<DatetimeClause>[] = [];
          const { clauses } = selectedData;
          clauses.forEach(item => {
            if (!datetimeClausesList.find(clause => clause?.id === item.id)) {
              updateDatetimeClauses.push({
                id: item.id,
                delete: true,
                clause: {
                  actionType: item.clause.actionType,
                  time: (item.clause as DatetimeClause).time
                }
              });
            }
          });
          datetimeClausesList.forEach(item => {
            if (!clauses.find(clause => clause.id === item?.id)) {
              updateDatetimeClauses.push({
                id: item.id || '',
                delete: false,
                clause: {
                  actionType: item.actionType,
                  time: Math.trunc(item.time.getTime() / 1000)?.toString()
                }
              });
            }
          });

          resp = await autoOpsUpdate({
            id: selectedData.id,
            environmentId,
            updateDatetimeClauses
          });
        } else {
          const datetimeClauses = datetimeClausesList.map(item => {
            const time = Math.trunc(item.time.getTime() / 1000)?.toString();
            return {
              time,
              actionType: item.actionType
            };
          });
          resp = await autoOpsCreator({
            featureId,
            environmentId,
            opsType: 'SCHEDULE',
            datetimeClauses
          });
        }

        if (resp) {
          onSubmitOperationSuccess();
          notify({
            message: `Schedule operation ${isCreate ? 'created' : 'updated'} successfully`
          });
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [isCreate, actionType, selectedData]
  );

  return (
    <SlideModal
      title={t(`common:${isCreate ? 'new' : 'update'}-operation`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col gap-y-5 w-full p-5 pb-28">
            <div className="flex items-center gap-x-4 typo-head-bold-small text-gray-800">
              <Trans
                i18nKey={'form:feature-flags.current-state'}
                values={{
                  state: t(`form:experiments.${isEnabledFlag ? 'on' : 'off'}`)
                }}
                components={{
                  comp: (
                    <div className="flex-center typo-para-small text-gray-600 px-2 py-[1px] border border-gray-400 rounded mb-[-4px]" />
                  )
                }}
              />
            </div>
            <Divider />
            <ScheduleList
              isCompletedTab={isCompletedTab}
              isCreate={isCreate}
              rollouts={rollouts}
            />
          </div>
          <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
            <ButtonBar
              primaryButton={
                <Button variant="secondary" onClick={onClose}>
                  {t(`common:cancel`)}
                </Button>
              }
              secondaryButton={
                <Button
                  type="submit"
                  loading={isSubmitting}
                  disabled={!isValid || isCompletedTab}
                >
                  {t(
                    isCreate
                      ? `feature-flags.create-operation`
                      : 'common:update-operation'
                  )}
                </Button>
              }
            />
          </div>
        </Form>
      </FormProvider>
    </SlideModal>
  );
};

export default ScheduleOperationModal;
