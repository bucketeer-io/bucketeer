import { useCallback } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { getCurrentEnvironment, useAuth } from 'auth';
import { Feature } from '@types';
import Form from 'components/form';
import { variationsFormSchema } from './form-schema';
import SubmitBar from './submit-bar';
import VariationsSection from './variations-section';

export interface VariationProps {
  feature: Feature;
}

const Variation = ({ feature }: VariationProps) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const form = useForm({
    resolver: yupResolver(variationsFormSchema),
    defaultValues: {
      comment: '',
      variations: feature.variations,
      variationType: feature.variationType,
      offVariation: feature.offVariation,
      onVariation: '',
      requireComment: currentEnvironment.requireComment,
      resetSampling: false
    }
  });

  const onSubmit = useCallback(() => {}, []);

  return (
    <FormProvider {...form}>
      <Form onSubmit={form.handleSubmit(onSubmit)}>
        <div className="flex flex-col w-full gap-y-6">
          <SubmitBar />
          <VariationsSection feature={feature} />
        </div>
      </Form>
    </FormProvider>
  );
};

export default Variation;
