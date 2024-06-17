import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { useCallback, useState, FC, memo, useEffect } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { DetailSkeleton } from '../../components/DetailSkeleton';
import { FeatureConfirmDialog } from '../../components/FeatureConfirmDialog';
import { FeatureVariationsForm } from '../../components/FeatureVariationsForm';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectById as selectFeatureById,
  updateFeatureVariations,
  getFeature,
  createCommand
} from '../../modules/features';
import { useCurrentEnvironment } from '../../modules/me';
import {
  AddVariationCommand,
  ChangeVariationDescriptionCommand,
  ChangeVariationNameCommand,
  ChangeVariationValueCommand,
  Command,
  RemoveVariationCommand
} from '../../proto/feature/command_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { AppDispatch } from '../../store';

import { variationsFormSchema } from './formSchema';
import { createResetSampleSeedCommand } from './targeting';

interface FeatureVariationsPageProps {
  featureId: string;
}

export const FeatureVariationsPage: FC<FeatureVariationsPageProps> = memo(
  ({ featureId }) => {
    const { formatMessage: f } = useIntl();
    const dispatch = useDispatch<AppDispatch>();
    const isFeatureLoading = useSelector<AppState, boolean>(
      (state) => state.features.loading
    );
    const isSegmentLoading = useSelector<AppState, boolean>(
      (state) => state.features.loading
    );
    const isLoading = isFeatureLoading || isSegmentLoading;
    const currentEnvironment = useCurrentEnvironment();
    const [feature, getFeatureError] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >(
      (state) => [
        selectFeatureById(state.features, featureId),
        state.features.getFeatureError
      ],
      shallowEqual
    );
    const defaultValues = (feature, requireComment: boolean) => {
      return {
        variationType: feature.variationType.toString(),
        variations: feature.variationsList,
        requireComment: requireComment,
        comment: ''
      };
    };
    const methods = useForm({
      resolver: yupResolver(variationsFormSchema),
      defaultValues: defaultValues(feature, currentEnvironment.requireComment),
      mode: 'onChange'
    });
    const {
      handleSubmit,
      formState: { dirtyFields }
    } = methods;
    const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);

    const handleUpdate = useCallback(
      async (data) => {
        const commands: Array<Command> = [];
        dirtyFields.variations &&
          commands.push(
            ...createVariationCommands(feature.variationsList, data.variations)
          );
        data.resetSampling && commands.push(createResetSampleSeedCommand());
        dispatch(
          updateFeatureVariations({
            environmentNamespace: currentEnvironment.id,
            id: feature.id,
            comment: data.comment,
            commands: commands
          })
        ).then(() => {
          setIsConfirmDialogOpen(false);
          dispatch(
            getFeature({
              environmentNamespace: currentEnvironment.id,
              id: featureId
            })
          );
        });
      },
      [feature, dispatch, dirtyFields]
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
        <FeatureVariationsForm
          featureId={featureId}
          onOpenConfirmDialog={() => setIsConfirmDialogOpen(true)}
        />
        {isConfirmDialogOpen && (
          <FeatureConfirmDialog
            open={isConfirmDialogOpen}
            handleSubmit={handleSubmit(handleUpdate)}
            onClose={() => setIsConfirmDialogOpen(false)}
            title={f(messages.feature.confirm.title)}
            description={f(messages.feature.confirm.description)}
            displayResetSampling={true}
          />
        )}
      </FormProvider>
    );
  }
);

export function createVariationCommands(
  orgVariations: any,
  valVariations: any
): Array<Command> {
  const commands: Array<Command> = [];
  const orgVariationIds = orgVariations.map((v) => v.id);
  const valVariationIds = valVariations.map((v) => v.id);
  // Intersection of org and val rules.
  const variationIds = orgVariationIds.filter((id) =>
    valVariationIds.includes(id)
  );
  orgVariations
    .filter((v) => !variationIds.includes(v.id))
    .forEach((v) => {
      const command = new RemoveVariationCommand();
      command.setId(v.id);
      commands.push(
        createCommand({ message: command, name: 'RemoveVariationCommand' })
      );
    });
  valVariations
    .filter((v) => !orgVariationIds.includes(v.id))
    .forEach((v) => {
      const command = new AddVariationCommand();
      command.setValue(v.value);
      command.setName(v.name);
      command.setDescription(v.description);
      commands.push(
        createCommand({ message: command, name: 'AddVariationCommand' })
      );
    });
  variationIds.forEach((vid) => {
    const orgVariation = orgVariations.find((v) => v.id === vid);
    const valVariation = valVariations.find((v) => v.id === vid);
    commands.push(...createVariationValueCommands(orgVariation, valVariation));
    commands.push(...createVariationNameCommands(orgVariation, valVariation));
    commands.push(
      ...createVariationDescriptionCommands(orgVariation, valVariation)
    );
  });
  return commands;
}

function createVariationValueCommands(
  orgVariation: any,
  valVariation: any
): Array<Command> {
  const commands: Array<Command> = [];
  if (orgVariation.value !== valVariation.value) {
    const command = new ChangeVariationValueCommand();
    command.setId(valVariation.id);
    command.setValue(valVariation.value);
    commands.push(
      createCommand({ message: command, name: 'ChangeVariationValueCommand' })
    );
  }
  return commands;
}

function createVariationNameCommands(
  orgVariation: any,
  valVariation: any
): Array<Command> {
  const commands: Array<Command> = [];
  if (orgVariation.name !== valVariation.name) {
    const command = new ChangeVariationNameCommand();
    command.setId(valVariation.id);
    command.setName(valVariation.name);
    commands.push(
      createCommand({ message: command, name: 'ChangeVariationNameCommand' })
    );
  }
  return commands;
}

function createVariationDescriptionCommands(
  orgVariation: any,
  valVariation: any
): Array<Command> {
  const commands: Array<Command> = [];
  if (orgVariation.description !== valVariation.description) {
    const command = new ChangeVariationDescriptionCommand();
    command.setId(valVariation.id);
    command.setDescription(valVariation.description);
    commands.push(
      createCommand({
        message: command,
        name: 'ChangeVariationDescriptionCommand'
      })
    );
  }
  return commands;
}
