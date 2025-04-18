import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError
} from '@reduxjs/toolkit';
import { Message } from 'google-protobuf';
import { Any } from 'google-protobuf/google/protobuf/any_pb';
import {
  BoolValue,
  Int32Value
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
  Command
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
  ListTagsResponse
} from '../proto/feature/service_pb';
import { Variation } from '../proto/feature/variation_pb';

import { AppState } from '.';

const MODULE_NAME = 'features';

export const featuresAdapter = createEntityAdapter({
  selectId: (feature: Feature.AsObject) => feature.id
});

export const { selectAll, selectById } = featuresAdapter.getSelectors();

export interface VariationParams {
  value: string;
  name: string;
  description: string;
}
export interface CreateFeatureParams {
  environmentId: string;
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

  request.setEnvironmentId(params.environmentId);
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
  await featureGrpc.createFeature(request);
  return params.id;
});

export interface CloneFeatureParams {
  environmentId: string;
  id: string;
  destinationenvironmentId: string;
}

export const cloneFeature = createAsyncThunk<
  string,
  CloneFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/clone`, async (params) => {
  const request = new CloneFeatureRequest();
  const cmd = new CloneFeatureCommand();
  cmd.setEnvironmentId(params.destinationenvironmentId);
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setCommand(cmd);
  await featureGrpc.cloneFeature(request);
  return params.id;
});

export type OrderBy =
  ListFeaturesRequest.OrderByMap[keyof ListFeaturesRequest.OrderByMap];
export type OrderDirection =
  ListFeaturesRequest.OrderDirectionMap[keyof ListFeaturesRequest.OrderDirectionMap];

export interface ListFeaturesParams {
  environmentId: string;
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
  request.setEnvironmentId(params.environmentId);
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
  const result = await featureGrpc.listFeatures(request);
  return result.response.toObject();
});

export const listTags = createAsyncThunk<
  ListTagsResponse.AsObject,
  undefined,
  { state: AppState }
>(`${MODULE_NAME}/listTags`, async () => {
  const request = new ListTagsRequest();
  const result = await featureGrpc.listTags(request);
  return result.response.toObject();
});

export interface GetFeatureParams {
  environmentId: string;
  id: string;
}

export const getFeature = createAsyncThunk<
  Feature.AsObject,
  GetFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetFeatureRequest();
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
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
  getFeatureError: null
});

export interface UpdateFeatureDetailsParams {
  environmentId: string;
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
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setComment(params.comment);

  const {
    renameCommand,
    changeDescriptionCommand,
    addTagCommands,
    removeTagCommands
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

  await featureGrpc.updateFeatureDetails(request);
});

export interface UpdateFeatureTargetingParams {
  environmentId: string;
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
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setComment(params.comment);
  request.setCommandsList(params.commands);
  request.setFrom(UpdateFeatureTargetingRequest.From.USER);
  await featureGrpc.updateFeatureTargeting(request);
});

export interface UpdateFeatureVariationsParams {
  environmentId: string;
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
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setComment(params.comment);
  request.setCommandsList(params.commands);
  await featureGrpc.updateFeatureVariations(request);
});

export interface EnableFeatureParams {
  environmentId: string;
  id: string;
  comment: string;
}

export const enableFeature = createAsyncThunk<
  void,
  EnableFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  const request = new EnableFeatureRequest();
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setComment(params.comment);
  request.setCommand(new EnableFeatureCommand());
  await featureGrpc.enableFeature(request);
});

export interface DisableFeatureParams {
  environmentId: string;
  id: string;
  comment?: string;
}

export const disableFeature = createAsyncThunk<
  void,
  DisableFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  const request = new DisableFeatureRequest();
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  if (params.comment) {
    request.setComment(params.comment);
  }
  request.setCommand(new DisableFeatureCommand());
  await featureGrpc.disableFeature(request);
});

export interface ArchiveFeatureParams {
  environmentId: string;
  id: string;
  comment: string;
}

export const archiveFeature = createAsyncThunk<
  void,
  ArchiveFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/archive`, async (params) => {
  const request = new ArchiveFeatureRequest();
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setComment(params.comment);
  request.setCommand(new ArchiveFeatureCommand());
  await featureGrpc.archiveFeature(request);
});

export interface UnarchiveFeatureParams {
  environmentId: string;
  id: string;
  comment: string;
}

export const unarchiveFeature = createAsyncThunk<
  void,
  UnarchiveFeatureParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/unarchive`, async (params) => {
  const request = new UnarchiveFeatureRequest();
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setComment(params.comment);
  request.setCommand(new UnarchiveFeatureCommand());
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
      .addCase(updateFeatureDetails.pending, () => {})
      .addCase(updateFeatureDetails.fulfilled, () => {})
      .addCase(updateFeatureDetails.rejected, () => {})
      .addCase(updateFeatureTargeting.pending, () => {})
      .addCase(updateFeatureTargeting.fulfilled, () => {})
      .addCase(updateFeatureTargeting.rejected, () => {})
      .addCase(createFeature.pending, () => {})
      .addCase(createFeature.fulfilled, () => {})
      .addCase(createFeature.rejected, () => {})
      .addCase(cloneFeature.pending, () => {})
      .addCase(cloneFeature.fulfilled, () => {})
      .addCase(cloneFeature.rejected, () => {})
      .addCase(enableFeature.pending, () => {})
      .addCase(enableFeature.fulfilled, () => {})
      .addCase(enableFeature.rejected, () => {})
      .addCase(disableFeature.pending, () => {})
      .addCase(disableFeature.fulfilled, () => {})
      .addCase(disableFeature.rejected, () => {});
  }
});

export const createCommand = ({
  message,
  name
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
