import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useState } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useDispatch, useSelector } from 'react-redux';

import { DetailSkeleton } from '../../components/DetailSkeleton';
import { FeatureConfirmDialog } from '../../components/FeatureConfirmDialog';
import { FeatureSettingsForm } from '../../components/FeatureSettingsForm';
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

interface FeatureSettingsPageProps {
  featureId: string;
}

export const FeatureSettingsPage: FC<FeatureSettingsPageProps> = memo(
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
    const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);

    const handleUpdate = useCallback(
      async (data) => {
        const commands: UpdateDetailCommands = {};
        if (dirtyFields.name) {
          commands.renameCommand = new RenameFeatureCommand();
          commands.renameCommand.setName(data.name);
        }
        if (dirtyFields.description) {
          commands.changeDescriptionCommand = new ChangeDescriptionCommand();
          commands.changeDescriptionCommand.setDescription(data.description);
        }
        if (dirtyFields.tags) {
          const addTags = data.tags?.filter(
            (tag) => !feature.tagsList.includes(tag)
          );
          if (addTags.length) {
            commands.addTagCommands = addTags.map((tag) => {
              const addTagCommand = new AddTagCommand();
              addTagCommand.setTag(tag);
              return addTagCommand;
            });
          }
          const removeTags = feature.tagsList.filter(
            (tag) => !data.tags?.includes(tag)
          );
          if (removeTags.length) {
            commands.removeTagCommands = removeTags.map((tag) => {
              const removeTagCommand = new RemoveTagCommand();
              removeTagCommand.setTag(tag);
              return removeTagCommand;
            });
          }
        }
        dispatch(
          updateFeatureDetails({
            environmentNamespace: currentEnvironment.id,
            id: feature.id,
            comment: data.comment,
            updateDetailCommands: commands,
          })
        ).then(() => {
          setIsConfirmDialogOpen(false);
          dispatch(
            getFeature({
              environmentNamespace: currentEnvironment.id,
              id: featureId,
            })
          );
        });
      },
      [dispatch, dirtyFields]
    );

    if (isLoading) {
      return (
        <div className="p-9 bg-gray-100">
          <DetailSkeleton />
        </div>
      );
    }
    return (
      <FormProvider {...methods}>
        <FeatureSettingsForm
          onOpenConfirmDialog={() => setIsConfirmDialogOpen(true)}
        />
        <FeatureConfirmDialog
          open={isConfirmDialogOpen}
          handleSubmit={handleSubmit(handleUpdate)}
          onClose={() => setIsConfirmDialogOpen(false)}
          title={f(messages.feature.confirm.title)}
          description={f(messages.feature.confirm.description)}
        />
      </FormProvider>
    );
  }
);
