import { useCallback, useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import {
  autoOpsCreator,
  AutoOpsCreatorResponse,
  autoOpsUpdate
} from '@api/auto-ops';
import { yupResolver } from '@hookform/resolvers/yup';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
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
          id: uuid(),
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
        const datetimeClauses = values.datetimeClausesList.map(item => {
          const time = Math.trunc(item.time.getTime() / 1000)?.toString();
          return {
            time,
            actionType: item.actionType
          };
        });

        let resp: AutoOpsCreatorResponse | null = null;

        if (!isCreate && selectedData) {
          resp = await autoOpsUpdate({
            environmentId,
            updateDatetimeClauses: {
              id: selectedData.id,
              delete: false,
              clause: datetimeClauses
            }
          });
        } else {
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
            <ScheduleList isCreate={isCreate} rollouts={rollouts} />
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
                  disabled={!isValid}
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
