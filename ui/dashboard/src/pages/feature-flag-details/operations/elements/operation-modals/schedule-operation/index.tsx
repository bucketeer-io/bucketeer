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
import useFormSchema from 'hooks/use-form-schema';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import isEqual from 'lodash/isEqual';
import { v4 as uuid } from 'uuid';
import {
  AutoOpsRule,
  DatetimeClause,
  RecurrenceFrequency,
  Rollout
} from '@types';
import { isSameOrBeforeDate } from 'utils/function';
import { cn } from 'utils/style';
import {
  recurringScheduleSchema,
  ScheduleOperationFormType
} from 'pages/feature-flag-details/operations/form-schema';
import {
  ActionTypeMap,
  EndConditionType,
  OperationActionType,
  ScheduleType
} from 'pages/feature-flag-details/operations/types';
import {
  createDatetimeClausesList,
  isRecurringOperation,
  timeOfDayToSeconds
} from 'pages/feature-flag-details/operations/utils';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Form from 'components/form';
import SlideModal from 'components/modal/slide';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import RecurringScheduleList from './recurring-schedule-list';
import ScheduleList from './schedule-list';

export interface OperationModalProps {
  editable: boolean;
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
  editable,
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

  const isExistingRecurring = useMemo(
    () => selectedData && isRecurringOperation(selectedData.clauses),
    [selectedData]
  );

  const defaultScheduleType = useMemo(() => {
    if (isExistingRecurring) return ScheduleType.RECURRING;
    return ScheduleType.ONE_TIME;
  }, [isExistingRecurring]);

  const handleCreateDefaultValues = () => {
    if (selectedData && !isExistingRecurring) {
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

  const handleCreateRecurringDefaults = () => {
    if (selectedData && isExistingRecurring) {
      const firstClause = selectedData.clauses[0]?.clause as DatetimeClause;
      const recurrence = firstClause?.recurrence;

      let endCondition = EndConditionType.NEVER;
      let endDate: Date | undefined;
      let maxOccurrences: number | undefined;

      if (recurrence?.endDate && Number(recurrence.endDate) > 0) {
        endCondition = EndConditionType.ON_DATE;
        endDate = new Date(Number(recurrence.endDate) * 1000);
      } else if (recurrence?.maxOccurrences && recurrence.maxOccurrences > 0) {
        endCondition = EndConditionType.AFTER;
        maxOccurrences = recurrence.maxOccurrences;
      }

      const startDate = recurrence?.startDate
        ? new Date(Number(recurrence.startDate) * 1000)
        : new Date();

      return {
        startDate,
        frequency: (recurrence?.frequency || 'WEEKLY') as RecurrenceFrequency,
        daysOfWeek: recurrence?.daysOfWeek?.map(Number) || [],
        dayOfMonth: recurrence?.dayOfMonth || 1,
        endCondition,
        endDate,
        maxOccurrences,
        recurringClauses: selectedData.clauses.map(c => {
          const dc = c.clause as DatetimeClause;
          const totalSeconds = Number(dc.time);
          const hours = Math.floor(totalSeconds / 3600);
          const minutes = Math.floor((totalSeconds % 3600) / 60);
          const timeDate = new Date();
          timeDate.setHours(hours, minutes, 0, 0);
          const wasExecuted = (dc.executionCount ?? 0) > 0;
          return {
            id: c.id,
            actionType: c.actionType as ActionTypeMap,
            time: timeDate,
            wasExecuted
          };
        })
      };
    }

    const defaultTime = new Date();
    defaultTime.setHours(defaultTime.getHours() + 1, 0, 0, 0);
    return {
      startDate: new Date(),
      frequency: 'WEEKLY' as RecurrenceFrequency,
      daysOfWeek: [1, 2, 3, 4, 5],
      dayOfMonth: 1,
      endCondition: EndConditionType.NEVER,
      endDate: undefined,
      maxOccurrences: undefined,
      recurringClauses: [
        {
          id: uuid(),
          actionType: ActionTypeMap.ENABLE,
          time: defaultTime
        }
      ]
    };
  };

  const form = useForm<ScheduleOperationFormType>({
    resolver: yupResolver(useFormSchema(recurringScheduleSchema)) as never,
    defaultValues: {
      scheduleType: defaultScheduleType,
      datetimeClausesList: handleCreateDefaultValues(),
      recurring: handleCreateRecurringDefaults()
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isDirty, isSubmitting }
  } = form;

  const scheduleType = form.watch('scheduleType');

  const handleCheckDateTimeClauses = useCallback(
    (datetimeClausesList: ScheduleOperationFormType['datetimeClausesList']) => {
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

  const buildRecurrenceRule = useCallback(
    (recurring: ScheduleOperationFormType['recurring']) => {
      const startDate = Math.trunc(
        recurring.startDate.getTime() / 1000
      ).toString();

      let endDate = '0';
      let maxOccurrences = 0;

      if (
        recurring.endCondition === EndConditionType.ON_DATE &&
        recurring.endDate
      ) {
        const endOfDay = new Date(recurring.endDate);
        endOfDay.setHours(23, 59, 59, 0);
        endDate = Math.trunc(endOfDay.getTime() / 1000).toString();
      } else if (
        recurring.endCondition === EndConditionType.AFTER &&
        recurring.maxOccurrences
      ) {
        maxOccurrences = recurring.maxOccurrences;
      }

      return {
        frequency: recurring.frequency,
        daysOfWeek:
          recurring.frequency === 'WEEKLY' ? recurring.daysOfWeek : [],
        dayOfMonth:
          recurring.frequency === 'MONTHLY' ? recurring.dayOfMonth : 0,
        startDate,
        endDate,
        maxOccurrences,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone
      };
    },
    []
  );

  const onSubmit = useCallback(
    async (values: ScheduleOperationFormType) => {
      try {
        if (editable) {
          let resp: AutoOpsCreatorResponse | null = null;

          if (values.scheduleType === ScheduleType.RECURRING) {
            const { recurring } = values;
            const recurrenceRule = buildRecurrenceRule(recurring);
            const datetimeClauses = recurring.recurringClauses.map(item => ({
              time: timeOfDayToSeconds(item.time).toString(),
              actionType: item.actionType,
              recurrence: recurrenceRule
            }));

            if (!isCreate && selectedData) {
              const datetimeClauseChanges: ClauseUpdateType<DatetimeClause>[] =
                [];

              const existingIds = new Set(selectedData.clauses.map(c => c.id));
              const formIds = new Set(
                recurring.recurringClauses.map(c => c.id).filter(Boolean)
              );

              selectedData.clauses.forEach(c => {
                if (!formIds.has(c.id)) {
                  datetimeClauseChanges.push({
                    id: c.id,
                    changeType: 'DELETE'
                  });
                }
              });

              const recurrenceRuleForUpdate = buildRecurrenceRule(recurring);
              recurring.recurringClauses.forEach(item => {
                const clause: DatetimeClause = {
                  time: timeOfDayToSeconds(item.time).toString(),
                  actionType: item.actionType,
                  recurrence: recurrenceRuleForUpdate
                };

                if (item.id && existingIds.has(item.id)) {
                  datetimeClauseChanges.push({
                    id: item.id,
                    changeType: 'UPDATE',
                    clause
                  });
                } else {
                  datetimeClauseChanges.push({
                    changeType: 'CREATE',
                    clause
                  });
                }
              });

              resp = await autoOpsUpdate({
                id: selectedData.id,
                environmentId,
                datetimeClauseChanges
              });
            } else {
              resp = await autoOpsCreator({
                featureId,
                environmentId,
                opsType: 'SCHEDULE',
                datetimeClauses
              });
            }
          } else {
            const { datetimeClausesList } = values;

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
          }

          if (resp) {
            onSubmitOperationSuccess();
            notify({
              message: t('message:collection-action-success', {
                collection: t('common:operation'),
                action: t(isCreate ? 'common:created' : 'common:updated')
              })
            });
          }
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [
      isCreate,
      actionType,
      selectedData,
      editable,
      environmentId,
      featureId,
      buildRecurrenceRule,
      handleCheckDateTimeClauses,
      onSubmitOperationSuccess,
      notify,
      errorNotify,
      t
    ]
  );

  useUnsavedLeavePage({ isShow: isDirty && !isSubmitting });

  const isDisabledMode = isFinishedTab && !!selectedData;

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

            {isCreate && (
              <Form.Field
                control={form.control}
                name="scheduleType"
                render={({ field }) => (
                  <Form.Item className="py-0">
                    <Form.Label required>
                      {t('feature-flags.schedule-type')}
                    </Form.Label>
                    <Form.Control>
                      <RadioGroup
                        value={field.value}
                        onValueChange={field.onChange}
                        className="flex gap-x-6"
                      >
                        <div className="flex items-center gap-x-2">
                          <RadioGroupItem value={ScheduleType.ONE_TIME} />
                          <span className="typo-para-medium text-gray-700">
                            {t('feature-flags.one-time')}
                          </span>
                        </div>
                        <div className="flex items-center gap-x-2">
                          <RadioGroupItem value={ScheduleType.RECURRING} />
                          <span className="typo-para-medium text-gray-700">
                            {t('feature-flags.recurring')}
                          </span>
                        </div>
                      </RadioGroup>
                    </Form.Control>
                  </Form.Item>
                )}
              />
            )}

            {scheduleType === ScheduleType.RECURRING ? (
              <RecurringScheduleList isDisabled={isDisabledMode} />
            ) : (
              <ScheduleList
                selectedData={selectedData}
                isFinishedTab={isFinishedTab}
                isCreate={isCreate}
                rollouts={rollouts}
              />
            )}
          </div>
          <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
            <ButtonBar
              primaryButton={
                <Button type="button" variant="secondary" onClick={onClose}>
                  {t(`common:cancel`)}
                </Button>
              }
              secondaryButton={
                <Button
                  type="submit"
                  loading={isSubmitting}
                  disabled={!isValid || isDisabledMode || !editable}
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
