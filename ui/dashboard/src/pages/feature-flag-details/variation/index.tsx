import { useCallback } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateFeature } from '@queries/feature-details';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
import { Feature } from '@types';
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
      requireComment: currentEnvironment?.requireComment,
      comment: '',
      resetSampling: false
    },
    mode: 'onChange'
  });

  const onShowConfirmDialog = useCallback(() => {
    onOpenConfirmDialog();
  }, []);

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
        }
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
      errorNotify(error);
    }
  }, []);

  return (
    <div className="flex flex-col w-full p-5 pt-0">
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col w-full gap-y-6">
            <SubmitBar onShowConfirmDialog={onShowConfirmDialog} />
            <VariationsSection feature={feature} />
          </div>
          <ConfirmationRequiredModal
            isOpen={openConfirmDialog}
            onClose={onCloseConfirmDialog}
            onSubmit={form.handleSubmit(onSubmit)}
          />
        </Form>
      </FormProvider>
    </div>
  );
};

export default Variation;
