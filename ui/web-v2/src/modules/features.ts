import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';
import { Message } from 'google-protobuf';
import { Any } from 'google-protobuf/google/protobuf/any_pb';
import {
  BoolValue,
  Int32Value,
} from 'google-protobuf/google/protobuf/wrappers_pb';

import * as featureGrpc from '../grpc/features';
import {
  RenameFeatureCommand,
  ChangeDescriptionCommand,
  AddTagCommand,
  RemoveTagCommand,
  CreateFeatureCommand,
  CloneFeatureCommand,
  EnableFeatureCommand,
  DisableFeatureCommand,
  ArchiveFeatureCommand,
  UnarchiveFeatureCommand,
  Command,
} from '../proto/feature/command_pb';
import { Feature } from '../proto/feature/feature_pb';
import {
  ArchiveFeatureRequest,
  CreateFeatureRequest,
  CloneFeatureRequest,
  DisableFeatureRequest,
  EnableFeatureRequest,
  GetFeatureRequest,
  ListFeaturesRequest,
  ListFeaturesResponse,
  UnarchiveFeatureRequest,
  UpdateFeatureDetailsRequest,
  UpdateFeatureTargetingRequest,
  UpdateFeatureVariationsRequest,
  ListTagsRequest,
  ListTagsResponse,
} from '../proto/feature/service_pb';
import { Variation } from '../proto/feature/variation_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'features';

export const featuresAdapter = createEntityAdapter<Feature.AsObject>({
  selectId: (feature) => feature.id,
});

export const { selectAll, selectById } = featuresAdapter.getSelectors();

export interface VariationParams {
  value: string;
  name: string;
  description: string;
}
export interface CreateFeatureParams {
  environmentNamespace: string;
  name: string;
  id: string;
  description: string;
  tagsList: Array<string>;
  variationType: Feature.VariationTypeMap[keyof Feature.VariationTypeMap];
  variations: VariationParams[];
  defaultOnVariationIndex: number;
  defaultOffVariationIndex: number;
}

export const createFeature = createAsyncThunk<
  string,
  CreateFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateFeatureRequest();
  const cmd = new CreateFeatureCommand();
  const variations: Variation[] = [];

  params.variations.forEach((v) => {
    const variation = new Variation();
    variation.setValue(v.value);
    variation.setName(v.name);
    variation.setDescription(v.description);
    variations.push(variation);
  });

  request.setEnvironmentNamespace(params.environmentNamespace);
  cmd.setName(params.name);
  cmd.setId(params.id);
  cmd.setDescription(params.description);
  cmd.setTagsList(params.tagsList);
  cmd.setVariationsList(variations);
  cmd.setVariationType(params.variationType);
  if (params.defaultOnVariationIndex >= 0) {
    const int32Value = new Int32Value();
    int32Value.setValue(params.defaultOnVariationIndex);
    cmd.setDefaultOnVariationIndex(int32Value);
  }
  if (params.defaultOffVariationIndex >= 0) {
    const int32Value = new Int32Value();
    int32Value.setValue(params.defaultOffVariationIndex);
    cmd.setDefaultOffVariationIndex(int32Value);
  }
  request.setCommand(cmd);
  await setupAuthToken();
  await featureGrpc.createFeature(request);
  return params.id;
});

export interface CloneFeatureParams {
  environmentNamespace: string;
  id: string;
  destinationEnvironmentNamespace: string;
}

export const cloneFeature = createAsyncThunk<
  string,
  CloneFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/clone`, async (params) => {
  const request = new CloneFeatureRequest();
  const cmd = new CloneFeatureCommand();
  cmd.setEnvironmentNamespace(params.destinationEnvironmentNamespace);
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(cmd);
  await setupAuthToken();
  await featureGrpc.cloneFeature(request);
  return params.id;
});

export type OrderBy =
  ListFeaturesRequest.OrderByMap[keyof ListFeaturesRequest.OrderByMap];
export type OrderDirection =
  ListFeaturesRequest.OrderDirectionMap[keyof ListFeaturesRequest.OrderDirectionMap];

export interface ListFeaturesParams {
  environmentNamespace: string;
  pageSize: number;
  cursor: string;
  tags: string[];
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchKeyword: string;
  enabled?: boolean;
  archived?: boolean;
  hasExperiment?: boolean;
  hasPrerequisites?: boolean;
  maintainerId: string;
}

export const listFeatures = createAsyncThunk<
  ListFeaturesResponse.AsObject,
  ListFeaturesParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListFeaturesRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setTagsList(params.tags);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  params.enabled != null &&
    request.setEnabled(new BoolValue().setValue(params.enabled));
  params.archived != null &&
    request.setArchived(new BoolValue().setValue(params.archived));
  params.hasExperiment != null &&
    request.setHasExperiment(new BoolValue().setValue(params.hasExperiment));
  params.hasPrerequisites != null &&
    request.setHasPrerequisites(
      new BoolValue().setValue(params.hasPrerequisites)
    );
  request.setMaintainer(params.maintainerId);
  await setupAuthToken();
  const result = await featureGrpc.listFeatures(request);
  return result.response.toObject();
});

export const listTags = createAsyncThunk<
  ListTagsResponse.AsObject,
  undefined,
  { state: AppState }
>(`${MODULE_NAME}/listTags`, async () => {
  const request = new ListTagsRequest();
  await setupAuthToken();
  const result = await featureGrpc.listTags(request);
  return result.response.toObject();
});

export interface GetFeatureParams {
  environmentNamespace: string;
  id: string;
}

export const getFeature = createAsyncThunk<
  Feature.AsObject,
  GetFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetFeatureRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  await setupAuthToken();
  const result = await featureGrpc.getFeature(request);
  return result.response.toObject().feature;
});

const initialState = featuresAdapter.getInitialState<{
  loading: boolean;
  listFeaturesLoading: boolean;
  totalCount: number;
  getFeatureError: SerializedError | null;
}>({
  loading: false,
  listFeaturesLoading: false,
  totalCount: 0,
  getFeatureError: null,
});

export interface UpdateFeatureDetailsParams {
  environmentNamespace: string;
  id: string;
  comment: string;
  updateDetailCommands: UpdateDetailCommands;
}
export interface UpdateDetailCommands {
  renameCommand?: RenameFeatureCommand;
  changeDescriptionCommand?: ChangeDescriptionCommand;
  addTagCommands?: AddTagCommand[];
  removeTagCommands?: RemoveTagCommand[];
}

export const updateFeatureDetails = createAsyncThunk<
  void,
  UpdateFeatureDetailsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/updateDetails`, async (params) => {
  const request = new UpdateFeatureDetailsRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setComment(params.comment);

  const {
    renameCommand,
    changeDescriptionCommand,
    addTagCommands,
    removeTagCommands,
  } = params.updateDetailCommands;
  if (renameCommand) {
    request.setRenameFeatureCommand(renameCommand);
  }
  if (changeDescriptionCommand) {
    request.setChangeDescriptionCommand(changeDescriptionCommand);
  }
  if (addTagCommands) {
    request.setAddTagCommandsList(addTagCommands);
  }
  if (removeTagCommands) {
    request.setRemoveTagCommandsList(removeTagCommands);
  }

  await setupAuthToken();
  await featureGrpc.updateFeatureDetails(request);
});

export interface UpdateFeatureTargetingParams {
  environmentNamespace: string;
  id: string;
  comment: string;
  commands: Command[];
}

export const updateFeatureTargeting = createAsyncThunk<
  void,
  UpdateFeatureTargetingParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/updateTargeting`, async (params) => {
  const request = new UpdateFeatureTargetingRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setComment(params.comment);
  request.setCommandsList(params.commands);
  request.setFrom(UpdateFeatureTargetingRequest.From.USER);
  await setupAuthToken();
  await featureGrpc.updateFeatureTargeting(request);
});

export interface UpdateFeatureVariationsParams {
  environmentNamespace: string;
  id: string;
  comment: string;
  commands: Command[];
}

export const updateFeatureVariations = createAsyncThunk<
  void,
  UpdateFeatureVariationsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/updateVariations`, async (params) => {
  const request = new UpdateFeatureVariationsRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setComment(params.comment);
  request.setCommandsList(params.commands);
  await setupAuthToken();
  await featureGrpc.updateFeatureVariations(request);
});

export interface EnableFeatureParams {
  environmentNamespace: string;
  id: string;
  comment: string;
}

export const enableFeature = createAsyncThunk<
  void,
  EnableFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  const request = new EnableFeatureRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setComment(params.comment);
  request.setCommand(new EnableFeatureCommand());
  await setupAuthToken();
  await featureGrpc.enableFeature(request);
});

export interface DisableFeatureParams {
  environmentNamespace: string;
  id: string;
  comment?: string;
}

export const disableFeature = createAsyncThunk<
  void,
  DisableFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  const request = new DisableFeatureRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  if (params.comment) {
    request.setComment(params.comment);
  }
  request.setCommand(new DisableFeatureCommand());
  await setupAuthToken();
  await featureGrpc.disableFeature(request);
});

export interface ArchiveFeatureParams {
  environmentNamespace: string;
  id: string;
  comment: string;
}

export const archiveFeature = createAsyncThunk<
  void,
  ArchiveFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/archive`, async (params) => {
  const request = new ArchiveFeatureRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setComment(params.comment);
  request.setCommand(new ArchiveFeatureCommand());
  await setupAuthToken();
  await featureGrpc.archiveFeature(request);
});

export interface UnarchiveFeatureParams {
  environmentNamespace: string;
  id: string;
  comment: string;
}

export const unarchiveFeature = createAsyncThunk<
  void,
  UnarchiveFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/unarchive`, async (params) => {
  const request = new UnarchiveFeatureRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setComment(params.comment);
  request.setCommand(new UnarchiveFeatureCommand());
  await setupAuthToken();
  await featureGrpc.unarchiveFeature(request);
});

export type FeaturesState = typeof initialState;

export const featuresSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listFeatures.pending, (state) => {
        state.listFeaturesLoading = true;
      })
      .addCase(listFeatures.fulfilled, (state, action) => {
        featuresAdapter.removeAll(state);
        featuresAdapter.upsertMany(state, action.payload.featuresList);
        state.totalCount = action.payload.totalCount;
        state.listFeaturesLoading = false;
      })
      .addCase(listFeatures.rejected, (state) => {
        state.listFeaturesLoading = false;
      })
      .addCase(getFeature.pending, (state) => {
        state.getFeatureError = null;
        state.loading = true;
      })
      .addCase(getFeature.fulfilled, (state, action) => {
        state.loading = false;
        state.getFeatureError = null;
        if (action.payload) {
          featuresAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(getFeature.rejected, (state, action) => {
        state.loading = false;
        state.getFeatureError = action.error;
      })
      .addCase(updateFeatureDetails.pending, (state) => {})
      .addCase(updateFeatureDetails.fulfilled, (state, action) => {})
      .addCase(updateFeatureDetails.rejected, (state, action) => {})
      .addCase(updateFeatureTargeting.pending, (state) => {})
      .addCase(updateFeatureTargeting.fulfilled, (state, action) => {})
      .addCase(updateFeatureTargeting.rejected, (state, action) => {})
      .addCase(createFeature.pending, (state) => {})
      .addCase(createFeature.fulfilled, (state, action) => {})
      .addCase(createFeature.rejected, (state, action) => {})
      .addCase(cloneFeature.pending, (state) => {})
      .addCase(cloneFeature.fulfilled, (state, action) => {})
      .addCase(cloneFeature.rejected, (state, action) => {})
      .addCase(enableFeature.pending, (state) => {})
      .addCase(enableFeature.fulfilled, (state, action) => {})
      .addCase(enableFeature.rejected, (state, action) => {})
      .addCase(disableFeature.pending, (state) => {})
      .addCase(disableFeature.fulfilled, (state, action) => {})
      .addCase(disableFeature.rejected, (state, action) => {});
  },
});

export const createCommand = ({
  message,
  name,
}: {
  message: Message;
  name: string;
}): Command => {
  const command = new Command();
  const any = new Any();
  any.pack(message.serializeBinary(), `bucketeer.feature.${name}`);
  command.setCommand(any);
  return command;
};
