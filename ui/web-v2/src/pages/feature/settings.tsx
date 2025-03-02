import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useDispatch, useSelector } from 'react-redux';

import { DetailSkeleton } from '../../components/DetailSkeleton';
import {
  FeatureConfirmDialog,
  SaveFeatureType
} from '../../components/FeatureConfirmDialog';
import { FeatureSettingsForm } from '../../components/FeatureSettingsForm';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectById as selectFeatureById,
  UpdateDetailCommands,
  updateFeatureDetails,
  getFeature,
  updateFeature
} from '../../modules/features';
import { useCurrentEnvironment } from '../../modules/me';
import {
  AddTagCommand,
  ChangeDescriptionCommand,
  RemoveTagCommand,
  RenameFeatureCommand
} from '../../proto/feature/command_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { AppDispatch } from '../../store';

import { settingsFormSchema } from './formSchema';
import { listTags } from '../../modules/tags';
import { ListTagsRequest } from '../../proto/tag/service_pb';
import { Tag } from '../../proto/tag/tag_pb';
import { ChangeType, TagChange } from '../../proto/feature/service_pb';

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
    const [feature] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError
    ]);
    const methods = useForm({
      resolver: yupResolver(
        settingsFormSchema(currentEnvironment.requireComment)
      ),
      defaultValues: {
        name: feature.name,
        description: feature.description,
        tags: feature.tagsList,
        comment: ''
      },
      mode: 'onChange'
    });
    const {
      handleSubmit,
      formState: { dirtyFields },
      reset
    } = methods;
    const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);

    useEffect(() => {
      dispatch(
        listTags({
          environmentId: currentEnvironment.id,
          pageSize: 0,
          cursor: '',
          orderBy: ListTagsRequest.OrderBy.DEFAULT,
          orderDirection: ListTagsRequest.OrderDirection.ASC,
          searchKeyword: null,
          entityType: Tag.EntityType.FEATURE_FLAG
        })
      );
    }, [dispatch]);

    const handleUpdate = useCallback(
      async (data, saveFeatureType) => {
        const prepareUpdate = async (updateAction, payload) => {
          dispatch(updateAction(payload)).then(() => {
            setIsConfirmDialogOpen(false);
            dispatch(
              getFeature({
                environmentId: currentEnvironment.id,
                id: featureId
              })
            ).then((res) => {
              const featurePayload = res.payload as Feature.AsObject;
              reset({
                name: featurePayload.name,
                description: featurePayload.description,
                tags: featurePayload.tagsList,
                comment: ''
              });
            });
          });
        };

        const tags = [];

        if (dirtyFields.tags) {
          const createTagChange = (type, tag) => {
            const tagChange = new TagChange();
            tagChange.setTag(tag);
            tagChange.setChangeType(type);
            return tagChange;
          };

          const featureTags = new Set(feature.tagsList);
          const dataTags = new Set(data.tags || []);

          dataTags.forEach((tag: string) => {
            if (!featureTags.has(tag)) {
              console.log('add tag');
              tags.push(createTagChange(ChangeType.CREATE, tag));
            }
          });

          featureTags.forEach((tag: string) => {
            if (!dataTags.has(tag)) {
              console.log('remove tag');
              tags.push(createTagChange(ChangeType.DELETE, tag));
            }
          });
        }

        if (saveFeatureType === SaveFeatureType.SCHEDULE) {
          await prepareUpdate(updateFeature, {
            environmentId: currentEnvironment.id,
            id: feature.id,
            comment: data.comment,
            name: dirtyFields.name ? data.name : undefined,
            description: dirtyFields.description ? data.description : undefined,
            tags: tags.length && tags
          });
        } else {
          const commands: UpdateDetailCommands = {};

          if (dirtyFields.name) {
            const renameCommand = new RenameFeatureCommand();
            renameCommand.setName(data.name);
            commands.renameCommand = renameCommand;
          }

          if (dirtyFields.description) {
            const descriptionCommand = new ChangeDescriptionCommand();
            descriptionCommand.setDescription(data.description);
            commands.changeDescriptionCommand = descriptionCommand;
          }

          if (dirtyFields.tags) {
            const addTags = data.tags?.filter(
              (tag) => !feature.tagsList.includes(tag)
            );
            if (addTags?.length) {
              commands.addTagCommands = addTags.map((tag) => {
                const addTagCommand = new AddTagCommand();
                addTagCommand.setTag(tag);
                return addTagCommand;
              });
            }

            const removeTags = feature.tagsList.filter(
              (tag) => !data.tags?.includes(tag)
            );
            if (removeTags?.length) {
              commands.removeTagCommands = removeTags.map((tag) => {
                const removeTagCommand = new RemoveTagCommand();
                removeTagCommand.setTag(tag);
                return removeTagCommand;
              });
            }
          }

          await prepareUpdate(updateFeatureDetails, {
            environmentId: currentEnvironment.id,
            id: feature.id,
            comment: data.comment,
            updateDetailCommands: commands
          });
        }
      },
      [dispatch, dirtyFields, currentEnvironment, featureId, feature, reset]
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
        {isConfirmDialogOpen && (
          <FeatureConfirmDialog
            open={isConfirmDialogOpen}
            handleSubmit={(arg) => {
              handleSubmit((data) => handleUpdate(data, arg))();
            }}
            onClose={() => setIsConfirmDialogOpen(false)}
            title={f(messages.feature.confirm.title)}
            description={f(messages.feature.confirm.description)}
          />
        )}
      </FormProvider>
    );
  }
);
