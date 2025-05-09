import { useCallback, useEffect } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateFeature } from '@queries/feature-details';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { isEqual } from 'lodash';
import { Feature, FeatureVariation, VariationChange } from '@types';
import Form from 'components/form';
import ConfirmationRequiredModal from '../elements/confirm-required-modal';
import { variationsFormSchema } from './form-schema';
import SubmitBar from './submit-bar';
import VariationsSection from './variations-section';

export interface VariationProps {
  feature: Feature;
}

const Variation = ({ feature }: VariationProps) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();

  const [openConfirmDialog, onOpenConfirmDialog, onCloseConfirmDialog] =
    useToggleOpen(false);

  const { notify, errorNotify } = useToast();

  const form = useForm({
    resolver: yupResolver(variationsFormSchema),
    defaultValues: {
      variations: feature.variations,
      variationType: feature.variationType,
      offVariation: feature.offVariation,
      onVariation: '',
      requireComment: false,
      comment: '',
      resetSampling: false
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

  const onSubmit = useCallback(async () => {
    try {
      const { variations, offVariation, comment, resetSampling } =
        form.getValues();
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
          message: t('flag-updated')
        });

        invalidateFeature(queryClient);
        onCloseConfirmDialog();
      }
    } catch (error) {
      errorNotify(error);
    }
  }, [feature]);

  useEffect(() => {
    form.reset({
      ...getValues(),
      variations: feature.variations,
      comment: '',
      resetSampling: false,
      requireComment: false
    });
  }, [feature]);

  return (
    <div className="p-6 pt-0 w-full min-w-[900px]">
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col w-full gap-y-6">
            <SubmitBar
              feature={feature}
              onShowConfirmDialog={() => {
                form.setValue(
                  'requireComment',
                  currentEnvironment?.requireComment
                );
                onOpenConfirmDialog();
              }}
            />
            <VariationsSection feature={feature} />
          </div>
          {openConfirmDialog && (
            <ConfirmationRequiredModal
              isOpen={openConfirmDialog}
              onClose={onCloseConfirmDialog}
              onSubmit={form.handleSubmit(onSubmit)}
            />
          )}
        </Form>
      </FormProvider>
    </div>
  );
};

export default Variation;
