import { useCallback, useEffect } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateFeature } from '@queries/feature-details';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
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

  const getVariationChanges = useCallback(
    (variations: FeatureVariation[]) => {
      const variationChanges: VariationChange[] = [];
      const { variations: variationsFeature } = feature;
      variationsFeature.forEach(variation => {
        const existingVariation = variations.find(v => v.id === variation.id);
        if (!existingVariation) {
          variationChanges.push({
            changeType: 'DELETE',
            variation
          });
        } else {
          const isEqualObj = isEqual(existingVariation, variation);
          if (!isEqualObj)
            variationChanges.push({
              changeType: 'UPDATE',
              variation: existingVariation
            });
        }
      });

      variations.forEach(variation => {
        const existingVariation = variationsFeature.find(
          v => v.id === variation.id
        );
        if (!existingVariation) {
          variationChanges.push({
            changeType: 'CREATE',
            variation
          });
        }
      });
      return variationChanges;
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
        variations: {
          values: variations
        },
        variationChanges: getVariationChanges(variations)
      });
      if (resp) {
        notify({
          message: (
            <span>
              <b>{feature.name}</b> {` has been successfully updated!`}
            </span>
          )
        });

        invalidateFeature(queryClient);
        onCloseConfirmDialog();
      }
    } catch (error) {
      errorNotify(
        error,
        'There is an Experiment scheduled to start or running. Please stop the Experiment before updating it.'
      );
      form.resetField('variations', {
        defaultValue: feature.variations
      });
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
    <>
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col w-full gap-y-6">
            <SubmitBar
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
              onClose={() => {
                onCloseConfirmDialog();
                form.setValue('comment', '');
                form.setValue('resetSampling', false);
                form.setValue('requireComment', false);
              }}
              onSubmit={form.handleSubmit(onSubmit)}
            />
          )}
        </Form>
      </FormProvider>
    </>
  );
};

export default Variation;
