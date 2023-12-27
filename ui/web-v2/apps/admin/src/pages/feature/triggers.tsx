import { FeatureTriggerForm } from '@/components/FeatureTriggerForm';
import { FlagTrigger } from '@/proto/feature/flag_trigger_pb';
import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useSelector } from 'react-redux';

import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { Feature } from '../../proto/feature/feature_pb';

import { triggerFormSchema } from './formSchema';

interface FeatureTriggerPageProps {
  featureId: string;
}

export const FeatureTriggerPage: FC<FeatureTriggerPageProps> = memo(
  ({ featureId }) => {
    const [feature, getFeatureError] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ]);
    const methods = useForm({
      resolver: yupResolver(triggerFormSchema),
      defaultValues: {
        triggerType: FlagTrigger.Type.TYPE_WEBHOOK.toString(),
        action: null,
        description: '',
      },
      mode: 'onChange',
    });
    return (
      <FormProvider {...methods}>
        <FeatureTriggerForm featureId={feature.id} />
      </FormProvider>
    );
  }
);
