import { useCallback, useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { rolloutCreator, RolloutCreatorParams } from '@api/rollouts';
import { yupResolver } from '@hookform/resolvers/yup';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import {
  Feature,
  IntervalType,
  Rollout,
  RolloutManualScheduleClause,
  RolloutSchedule,
  RolloutTemplateScheduleClause
} from '@types';
import { cn } from 'utils/style';
import {
  rolloutSchema,
  RolloutSchemaType
} from 'pages/feature-flag-details/operations/form-schema';
import {
  IntervalMap,
  RolloutTypeMap
} from 'pages/feature-flag-details/operations/types';
import {
  createProgressiveRollout,
  handleCreateIncrement
} from 'pages/feature-flag-details/operations/utils';
import {
  OperationActionType,
  ScheduleItem
} from 'pages/feature-flag-details/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import SlideModal from 'components/modal/slide';
import ManualSchedule from './manual-schedule';
import TemplateSchedule from './template-schedule';
import RolloutWarning from './warning';

export interface OperationModalProps {
  feature: Feature;
  environmentId: string;
  isOpen: boolean;
  actionType: OperationActionType;
  selectedData?: Rollout;
  rollouts: Rollout[];
  onClose: () => void;
  onSubmitRolloutSuccess: () => void;
}

const buttonCls =
  '!typo-para-medium flex-1 !text-gray-600 !shadow-none border border-gray-200 hover:border-gray-400';
const buttonActiveCls =
  '!text-accent-pink-500 border-accent-pink-500 hover:!text-accent-pink-500 hover:border-accent-pink-500';

const ProgressiveRolloutModal = ({
  feature,
  environmentId,
  isOpen,
  actionType,
  selectedData,
  rollouts,
  onClose,
  onSubmitRolloutSuccess
}: OperationModalProps) => {
  const { t } = useTranslation(['form', 'common', 'table', 'message']);
  const { notify, errorNotify } = useToast();

  const form = useForm({
    resolver: yupResolver(rolloutSchema),
    defaultValues: {
      progressiveRolloutType: RolloutTypeMap.TEMPLATE_SCHEDULE,
      progressiveRollout: createProgressiveRollout(feature)
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isSubmitting },
    watch,
    setValue
  } = form;

  const progressiveRolloutType = watch('progressiveRolloutType');

  const isDisableCreateRollout = useMemo(() => {
    return (
      rollouts.length > 0 &&
      !!rollouts.find(item => ['WAITING', 'RUNNING'].includes(item.status))
    );
  }, [rollouts]);

  const isTemplateRollout = useMemo(
    () => progressiveRolloutType === RolloutTypeMap.TEMPLATE_SCHEDULE,
    [progressiveRolloutType]
  );

  const variationOptions = useMemo(
    () =>
      feature.variations.map(item => ({
        label: item.name || item.value,
        value: item.id
      })),
    [feature]
  );

  const handleChangeTab = useCallback((rolloutType: RolloutTypeMap) => {
    setValue('progressiveRolloutType', rolloutType);
  }, []);

  const onSubmit = useCallback(
    async (values: RolloutSchemaType) => {
      try {
        const { manual, template } = values.progressiveRollout;
        const payload: RolloutCreatorParams = {
          environmentId,
          featureId: feature.id
        };
        if (values.progressiveRolloutType === RolloutTypeMap.MANUAL_SCHEDULE) {
          const progressiveRolloutManualScheduleClause: RolloutManualScheduleClause =
            {
              variationId: manual.variationId,
              schedules: manual.schedulesList.map(item => ({
                ...item,
                weight: item.weight * 1000,
                executeAt: Math.trunc(
                  item.executeAt?.getTime() / 1000
                )?.toString()
              }))
            };
          Object.assign(payload, { progressiveRolloutManualScheduleClause });
        } else {
          const { increments, interval, startDate, variationId } = template;
          const lastSchedule: ScheduleItem = {
            scheduleId: uuid(),
            weight: increments,
            executeAt: startDate,
            triggeredAt: '0'
          };
          const scheduleList = [lastSchedule];
          const templateInterval = interval as IntervalMap;
          const incrementType =
            templateInterval === IntervalMap.DAILY
              ? 'day'
              : templateInterval === IntervalMap.WEEKLY
                ? 'week'
                : 'hour';

          while (
            !scheduleList.find(item => item.weight === 100) &&
            lastSchedule.weight !== 100
          ) {
            const incrementItem = handleCreateIncrement({
              lastSchedule: scheduleList.at(-1)!,
              incrementType,
              increment: increments
            });
            scheduleList.push(incrementItem);
          }

          const progressiveRolloutTemplateScheduleClause: RolloutTemplateScheduleClause =
            {
              variationId,
              increments: increments.toString(),
              schedules: scheduleList.map(schedule => ({
                ...schedule,
                weight: schedule.weight * 1000,
                executeAt: Math.trunc(
                  schedule.executeAt.getTime() / 1000
                ).toString()
              })) as RolloutSchedule[],
              interval: interval as IntervalType
            };

          Object.assign(payload, { progressiveRolloutTemplateScheduleClause });
        }

        const resp = await rolloutCreator(payload);

        if (resp) {
          notify({
            message: t('message:operation.created')
          });
          onSubmitRolloutSuccess();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [actionType, selectedData]
  );

  return (
    <SlideModal
      title={t(`common:new-operation`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col gap-y-5 w-full p-5 pb-28">
            <p className="typo-head-bold-small text-gray-800">
              {t('common:source-type.progressive-rollout')}
            </p>
            {isDisableCreateRollout && <RolloutWarning />}
            <div className="flex items-center">
              <Button
                type="button"
                variant={'secondary-2'}
                size={'sm'}
                className={cn(
                  'rounded-r-none',
                  buttonCls,
                  isTemplateRollout && buttonActiveCls
                )}
                disabled={isDisableCreateRollout}
                onClick={() =>
                  handleChangeTab(RolloutTypeMap.TEMPLATE_SCHEDULE)
                }
              >
                {t(`template`)}
              </Button>
              <Button
                type="button"
                variant={'secondary-2'}
                size={'sm'}
                className={cn(
                  'rounded-l-none',
                  buttonCls,
                  !isTemplateRollout && buttonActiveCls && buttonActiveCls
                )}
                disabled={isDisableCreateRollout}
                onClick={() => handleChangeTab(RolloutTypeMap.MANUAL_SCHEDULE)}
              >
                {t(`manual`)}
              </Button>
            </div>
            {isTemplateRollout ? (
              <TemplateSchedule
                variationOptions={variationOptions}
                isDisableCreateRollout={isDisableCreateRollout}
              />
            ) : (
              <ManualSchedule
                variationOptions={variationOptions}
                isDisableCreateRollout={isDisableCreateRollout}
              />
            )}
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
                  disabled={!isValid || isDisableCreateRollout}
                >
                  {t(`feature-flags.create-operation`)}
                </Button>
              }
            />
          </div>
        </Form>
      </FormProvider>
    </SlideModal>
  );
};

export default ProgressiveRolloutModal;
