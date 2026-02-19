import { useCallback, useEffect } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryExperiments } from '@queries/experiments';
import { invalidateFeature } from '@queries/feature-details';
import { invalidateFeatures } from '@queries/features';
import { useQueryRollouts } from '@queries/rollouts';
import { useCreateScheduledFlagChange } from '@queries/scheduled-flag-changes';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { SCHEDULED_FLAG_CHANGES_ENABLED } from 'configs';
import {
  PAGE_PATH_EXPERIMENTS,
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import useFormSchema from 'hooks/use-form-schema';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import isEqual from 'lodash/isEqual';
import {
  Feature,
  FeatureVariation,
  ScheduledChangePayload,
  VariationChange
} from '@types';
import { IconInfoFilled } from '@icons';
import Form from 'components/form';
import Icon from 'components/icon';
import InfoMessage from 'components/info-message';
import { CardNote } from 'elements/overview-card';
import ConfirmationRequiredModal, {
  ConfirmRequiredValues
} from '../elements/confirm-required-modal';
import { SCHEDULE_TYPE_SCHEDULE } from '../elements/confirm-required-modal/form-schema';
import ScheduledChangesBanner from '../elements/scheduled-changes-banner';
import { variationsFormSchema } from './form-schema';
import SubmitBar from './submit-bar';
import VariationsSection from './variations-section';

export interface VariationProps {
  feature: Feature;
  editable: boolean;
  isRunningExperiment?: boolean;
}

const Variation = ({ feature, editable }: VariationProps) => {
  const { t } = useTranslation(['common', 'message', 'form']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();

  const [openConfirmDialog, onOpenConfirmDialog, onCloseConfirmDialog] =
    useToggleOpen(false);

  const createScheduleMutation = useCreateScheduledFlagChange();
  const { notify, errorNotify } = useToast();

  const { data: rolloutCollection } = useQueryRollouts({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id,
      featureIds: [feature.id]
    }
  });

  const { data: experimentCollection } = useQueryExperiments({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id,
      featureId: feature.id,
      statuses: ['WAITING', 'RUNNING']
    }
  });

  const isRunningExperiment =
    !!experimentCollection && experimentCollection?.experiments?.length > 0;

  const waitingRunningRollouts =
    rolloutCollection?.progressiveRollouts?.filter(item =>
      ['WAITING', 'RUNNING'].includes(item.status)
    ) || [];

  const form = useForm({
    resolver: yupResolver(useFormSchema(variationsFormSchema)),
    defaultValues: {
      variations: feature.variations,
      variationType: feature.variationType,
      offVariation: feature.offVariation,
      onVariation: ''
    },
    mode: 'onChange'
  });

  const {
    getValues,
    formState: { isDirty, isSubmitting }
  } = form;

  const handleCheckVariations = useCallback(
    (variations: FeatureVariation[]) => {
      const { variations: featureVariations } = feature;
      const variationChanges: VariationChange[] = [];

      featureVariations.forEach(item => {
        if (!variations.find(variation => variation.id === item.id)) {
          variationChanges.push({
            changeType: 'DELETE',
            variation: item
          });
        }
      });
      variations.forEach(item => {
        const currentVariation = featureVariations.find(
          variation => variation.id === item.id
        );
        if (!currentVariation) {
          variationChanges.push({
            changeType: 'CREATE',
            variation: item
          });
        } else {
          if (!isEqual(currentVariation, item))
            variationChanges.push({
              changeType: 'UPDATE',
              variation: item
            });
        }
      });

      return {
        variationChanges
      };
    },
    [feature]
  );

  const onSubmit = useCallback(
    async (additionalValues?: ConfirmRequiredValues) => {
      if (editable) {
        try {
          const { variations, offVariation } = form.getValues();
          const { comment, resetSampling, scheduleType, scheduleAt } =
            additionalValues || {};

          const isScheduleUpdate = scheduleType === SCHEDULE_TYPE_SCHEDULE;

          if (isScheduleUpdate) {
            const { variationChanges } = handleCheckVariations(variations);
            const payload: ScheduledChangePayload = {};
            if (variationChanges.length > 0) {
              payload.variationChanges = variationChanges;
            }
            if (offVariation !== feature.offVariation) {
              payload.offVariation = offVariation;
            }
            if (resetSampling) {
              payload.resetSamplingSeed = true;
            }
            const resp = await createScheduleMutation.mutateAsync({
              environmentId: currentEnvironment.id,
              featureId: feature.id,
              scheduledAt: scheduleAt as string,
              payload,
              comment
            });
            if (resp) {
              notify({
                message: t('form:feature-flags.schedule-configured', {
                  name: feature.name
                })
              });
              onCloseConfirmDialog();
            }
          } else {
            const resp = await featureUpdater({
              id: feature.id,
              environmentId: currentEnvironment.id,
              comment,
              resetSamplingSeed: resetSampling,
              offVariation,
              ...handleCheckVariations(variations)
            });
            if (resp) {
              notify({
                message: t('message:collection-action-success', {
                  collection: t('source-type.feature-flag'),
                  action: t('updated')
                })
              });
              invalidateFeature(queryClient);
              invalidateFeatures(queryClient);
              onCloseConfirmDialog();
            }
          }
        } catch (error) {
          errorNotify(error);
        }
      }
    },
    [feature, editable]
  );

  useEffect(() => {
    form.reset({
      ...getValues(),
      variations: feature.variations
    });
  }, [feature]);
  useUnsavedLeavePage({ isShow: isDirty && !isSubmitting });
  return (
    <div className="p-6 pt-0 w-full min-w-[900px]">
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(() => onSubmit())}>
          <div className="flex flex-col w-full gap-y-6">
            {SCHEDULED_FLAG_CHANGES_ENABLED && (
              <ScheduledChangesBanner
                featureId={feature.id}
                environmentId={currentEnvironment.id}
              />
            )}
            <div className="flex flex-col gap-2">
              <SubmitBar
                editable={editable}
                feature={feature}
                onShowConfirmDialog={onOpenConfirmDialog}
              />
              {feature.variationType === 'YAML' && (
                <CardNote content={t('form:yaml-note')} className="w-fit" />
              )}
            </div>
            {waitingRunningRollouts.length > 0 && (
              <div className="flex items-center gap-x-3 p-4 rounded bg-accent-blue-50 border-l-4 border-accent-blue-500 text-accent-blue-500 typo-para-medium">
                <Icon icon={IconInfoFilled} color="accent-blue-500" size="sm" />
                <div className="flex items-center [&>a]:ml-1">
                  <Trans
                    i18nKey={'form:variation.rollout-running-message'}
                    components={{
                      comp: (
                        <Link
                          className="text-primary-500 underline"
                          to={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${feature.id}${PAGE_PATH_FEATURE_AUTOOPS}`}
                        />
                      )
                    }}
                  />
                </div>
              </div>
            )}
            {isRunningExperiment && (
              <InfoMessage
                title={t('message:validation.experiment-running-warning')}
                description={t(
                  'message:validation.experiment-running-warning-desc'
                )}
                linkElements={experimentCollection.experiments.map(
                  (item, index) => (
                    <li
                      key={index}
                      className="typo-para-small text-primary-500 underline w-fit max-w-full truncate"
                    >
                      <Link
                        to={`/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${item.id}/results`}
                      >
                        {item.name}
                      </Link>
                    </li>
                  )
                )}
              />
            )}
            <VariationsSection
              editable={editable && !waitingRunningRollouts.length}
              feature={feature}
              isRunningExperiment={isRunningExperiment}
            />
          </div>
        </Form>
      </FormProvider>
      {openConfirmDialog && (
        <ConfirmationRequiredModal
          feature={feature}
          isOpen={openConfirmDialog}
          isShowScheduleSelect={SCHEDULED_FLAG_CHANGES_ENABLED}
          onClose={onCloseConfirmDialog}
          onSubmit={additionalValues =>
            form.handleSubmit(() => onSubmit(additionalValues))()
          }
        />
      )}
    </div>
  );
};

export default Variation;
