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
import { isEqual } from 'lodash';
import { AutoOpsRule, DatetimeClause, Rollout } from '@types';
import { isSameOrBeforeDate } from 'utils/function';
import { cn } from 'utils/style';
import {
  dateTimeClauseListSchema,
  DateTimeClauseListType
} from 'pages/feature-flag-details/operations/form-schema';
import {
  ActionTypeMap,
  OperationActionType
} from 'pages/feature-flag-details/operations/types';
import { createDatetimeClausesList } from 'pages/feature-flag-details/operations/utils';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Form from 'components/form';
import SlideModal from 'components/modal/slide';
import ScheduleList from './schedule-list';

export interface OperationModalProps {
  isFinishedTab: boolean;
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
  isFinishedTab,
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
  const { t } = useTranslation(['form', 'common', 'message']);
  const { notify, errorNotify } = useToast();

  const isCreate = useMemo(() => actionType === 'NEW', [actionType]);

  const handleCreateDefaultValues = () => {
    if (selectedData) {
      return selectedData.clauses.map(item => {
        const time = new Date(+(item.clause as DatetimeClause).time * 1000);
        return {
          id: item.id,
          actionType: item.actionType as ActionTypeMap,
          time,
          wasPassed: isSameOrBeforeDate(time)
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

  const handleCheckDateTimeClauses = useCallback(
    (datetimeClausesList: DateTimeClauseListType['datetimeClausesList']) => {
      if (selectedData) {
        const datetimeClauseChanges: ClauseUpdateType<DatetimeClause>[] = [];
        const { clauses } = selectedData;
        const clausesFormatted = clauses.map(clause => {
          const time = new Date(
            +(clause.clause as DatetimeClause)?.time * 1000
          );
          return {
            actionType: clause.actionType,
            id: clause.id,
            time,
            wasPassed: isSameOrBeforeDate(time)
          };
        });
        clausesFormatted.forEach(item => {
          const currentClause = datetimeClausesList.find(
            clause => clause?.id === item.id
          );
          if (!currentClause) {
            datetimeClauseChanges.push({
              id: item.id,
              changeType: 'DELETE'
            });
          }
        });

        datetimeClausesList.forEach(item => {
          const currentClause = clausesFormatted.find(
            clause => clause.id === item?.id
          );
          if (!currentClause) {
            datetimeClauseChanges.push({
              changeType: 'CREATE',
              clause: {
                actionType: item.actionType,
                time: Math.trunc(item.time.getTime() / 1000)?.toString()
              }
            });
          }

          if (currentClause && !isEqual(currentClause, item)) {
            datetimeClauseChanges.push({
              id: item.id || '',
              changeType: 'UPDATE',
              clause: {
                actionType: item.actionType,
                time: Math.trunc(item.time.getTime() / 1000)?.toString()
              }
            });
          }
        });
        return datetimeClauseChanges;
      }
      return [];
    },
    [selectedData]
  );

  const onSubmit = useCallback(
    async (values: DateTimeClauseListType) => {
      try {
        const { datetimeClausesList } = values;

        let resp: AutoOpsCreatorResponse | null = null;

        if (!isCreate && selectedData) {
          const datetimeClauseChanges =
            handleCheckDateTimeClauses(datetimeClausesList);

          resp = await autoOpsUpdate({
            id: selectedData.id,
            environmentId,
            datetimeClauseChanges
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
            message: t(`message:operation.${isCreate ? 'created' : 'updated'}`)
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
                    <div
                      className={cn(
                        'flex-center typo-para-small text-gray-600 px-2 py-[1px] border border-gray-400 rounded mb-[-4px]',
                        {
                          'bg-primary-500 text-white': isEnabledFlag
                        }
                      )}
                    />
                  )
                }}
              />
            </div>
            <Divider />
            <ScheduleList
              selectedData={selectedData}
              isFinishedTab={isFinishedTab}
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
                  disabled={!isValid || (isFinishedTab && !!selectedData)}
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
