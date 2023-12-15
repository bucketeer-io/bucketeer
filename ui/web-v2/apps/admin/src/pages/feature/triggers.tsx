import { FeatureTriggerForm } from '@/components/FeatureTriggerForm';
import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useState } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useDispatch, useSelector } from 'react-redux';

import { DetailSkeleton } from '../../components/DetailSkeleton';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectById as selectFeatureById,
  UpdateDetailCommands,
  updateFeatureDetails,
  getFeature,
} from '../../modules/features';
import { useCurrentEnvironment } from '../../modules/me';
import {
  AddTagCommand,
  ChangeDescriptionCommand,
  RemoveTagCommand,
  RenameFeatureCommand,
} from '../../proto/feature/command_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { AppDispatch } from '../../store';

import { settingsFormSchema } from './formSchema';

interface FeatureTriggerPageProps {
  featureId: string;
}

export const FeatureTriggerPage: FC<FeatureTriggerPageProps> = memo(
  ({ featureId }) => {
    const { formatMessage: f } = useIntl();
    const dispatch = useDispatch<AppDispatch>();
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.features.loading
    );
    const currentEnvironment = useCurrentEnvironment();
    const [feature, getFeatureError] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ]);
    const methods = useForm({
      resolver: yupResolver(settingsFormSchema),
      defaultValues: {
        name: feature.name,
        description: feature.description,
        tags: feature.tagsList,
        comment: '',
      },
      mode: 'onChange',
    });
    const {
      handleSubmit,
      formState: { dirtyFields },
    } = methods;

    // const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);

    if (isLoading) {
      return (
        <div className="p-9 bg-gray-100">
          <DetailSkeleton />
        </div>
      );
    }

    return (
      <FormProvider {...methods}>
        <FeatureTriggerForm
        // onOpenConfirmDialog={() => setIsConfirmDialogOpen(true)}
        />
        {/* <FeatureConfirmDialog
          open={isConfirmDialogOpen}
          handleSubmit={handleSubmit(handleUpdate)}
          onClose={() => setIsConfirmDialogOpen(false)}
          title={f(messages.feature.confirm.title)}
          description={f(messages.feature.confirm.description)}
        /> */}
      </FormProvider>
    );
  }
);
