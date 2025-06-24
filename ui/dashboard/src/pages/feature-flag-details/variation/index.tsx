import { useCallback, useEffect } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Link } from 'react-router-dom';
import { featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryExperiments } from '@queries/experiments';
import { invalidateFeature } from '@queries/feature-details';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { isEqual } from 'lodash';
import { Feature, FeatureVariation, VariationChange } from '@types';
import Form from 'components/form';
import InfoMessage from 'components/info-message';
import ConfirmationRequiredModal, {
  ConfirmRequiredValues
} from '../elements/confirm-required-modal';
import { variationsFormSchema } from './form-schema';
import SubmitBar from './submit-bar';
import VariationsSection from './variations-section';

export interface VariationProps {
  feature: Feature;
  editable: boolean;
  isRunningExperiment?: boolean;
}

const Variation = ({ feature, editable }: VariationProps) => {
  const { t } = useTranslation(['common', 'message']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();

  const [openConfirmDialog, onOpenConfirmDialog, onCloseConfirmDialog] =
    useToggleOpen(false);

  const { notify, errorNotify } = useToast();

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

  const form = useForm({
    resolver: yupResolver(variationsFormSchema),
    defaultValues: {
      variations: feature.variations,
      variationType: feature.variationType,
      offVariation: feature.offVariation,
      onVariation: ''
    },
    mode: 'onChange'
  });

  const { getValues } = form;

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
          const { comment, resetSampling } = additionalValues || {};

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
            onCloseConfirmDialog();
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

  return (
    <div className="p-6 pt-0 w-full min-w-[900px]">
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(() => onSubmit())}>
          <div className="flex flex-col w-full gap-y-6">
            <SubmitBar
              editable={editable}
              feature={feature}
              onShowConfirmDialog={onOpenConfirmDialog}
            />
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
              editable={editable}
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
