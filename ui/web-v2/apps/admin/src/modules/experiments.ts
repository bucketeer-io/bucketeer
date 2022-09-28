import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';
import {
  BoolValue,
  Int32Value,
} from 'google-protobuf/google/protobuf/wrappers_pb';

import * as grpc from '../grpc/experiment';
import {
  ArchiveExperimentCommand,
  ChangeExperimentDescriptionCommand,
  ChangeExperimentNameCommand,
  ChangeExperimentPeriodCommand,
  CreateExperimentCommand,
  StopExperimentCommand,
} from '../proto/experiment/command_pb';
import { Experiment } from '../proto/experiment/experiment_pb';
import {
  ArchiveExperimentRequest,
  CreateExperimentRequest,
  GetExperimentRequest,
  ListExperimentsRequest,
  ListExperimentsResponse,
  StopExperimentRequest,
  UpdateExperimentRequest,
} from '../proto/experiment/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'experiments';

export const experimentsAdapter = createEntityAdapter<Experiment.AsObject>({
  selectId: (experiment) => experiment.id,
});

export const { selectAll, selectById } = experimentsAdapter.getSelectors();

const initialState = experimentsAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getExperimentError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getExperimentError: null,
});

export type ExperimentsState = typeof initialState;

export interface GetExperimentParams {
  environmentNamespace: string;
  id: string;
}

export const getExperiment = createAsyncThunk<
  Experiment.AsObject,
  GetExperimentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/getExperiment`, async (params) => {
  const request = new GetExperimentRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  await setupAuthToken();
  const result = await grpc.getExperiment(request);
  return result.response.getExperiment().toObject();
});

export type OrderBy =
  ListExperimentsRequest.OrderByMap[keyof ListExperimentsRequest.OrderByMap];
export type OrderDirection =
  ListExperimentsRequest.OrderDirectionMap[keyof ListExperimentsRequest.OrderDirectionMap];

export interface ListExperimentsParams {
  environmentNamespace: string;
  pageSize: number;
  cursor: string;
  featureId?: string;
  featureVersion?: number;
  startFrom?: number;
  stopUntil?: number;
  searchKeyword: string;
  status?: number;
  archived?: boolean;
  maintainer?: string;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

export const listExperiments = createAsyncThunk<
  ListExperimentsResponse.AsObject,
  ListExperimentsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListExperimentsRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setPageSize(params.pageSize);
  request.setSearchKeyword(params.searchKeyword);
  request.setCursor(params.cursor);
  const boolValue = new BoolValue();
  boolValue.setValue(params.archived);
  request.setArchived(boolValue);
  request.setMaintainer(params.maintainer);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  if (params.status != null) {
    const int32Value = new Int32Value();
    int32Value.setValue(params.status);
    request.setStatus(int32Value);
  }
  if (params.featureId) {
    request.setFeatureId(params.featureId);
  }
  // If version is unset, fetch experiments of all versions.
  if (params.featureVersion) {
    const version = new Int32Value();
    version.setValue(params.featureVersion);
    request.setFeatureVersion(version);
  }
  if (params.startFrom) {
    request.setFrom(params.startFrom);
  }
  if (params.stopUntil) {
    request.setTo(params.stopUntil);
  }
  await setupAuthToken();
  const result = await grpc.listExperiments(request);
  return result.response.toObject();
});

export interface CreateExperimentParams {
  name: string;
  description?: string;
  environmentNamespace: string;
  featureId: string;
  baseVariationId: string;
  goalIdsList: string[];
  startAt: number;
  stopAt: number;
}

export const createExperiment = createAsyncThunk<
  Experiment.AsObject,
  CreateExperimentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateExperimentRequest();
  const command = new CreateExperimentCommand();
  command.setName(params.name);
  if (params.description) {
    command.setDescription(params.description);
  }
  command.setFeatureId(params.featureId);
  command.setBaseVariationId(params.baseVariationId);
  command.setGoalIdsList(params.goalIdsList);
  command.setStartAt(Math.floor(params.startAt));
  command.setStopAt(Math.floor(params.stopAt));
  request.setCommand(command);
  request.setEnvironmentNamespace(params.environmentNamespace);

  await setupAuthToken();
  const result = await grpc.createExperiment(request);
  return result.response.toObject().experiment;
});

export interface UpdateExperimentParams {
  environmentNamespace: string;
  id: string;
  changeNameCommand?: ChangeExperimentNameCommand;
  changeDescriptionCommand?: ChangeExperimentDescriptionCommand;
  changePeriodCommand?: ChangeExperimentPeriodCommand;
}

export const updateExperiment = createAsyncThunk<
  void,
  UpdateExperimentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new UpdateExperimentRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  if (params.changeNameCommand) {
    request.setChangeNameCommand(params.changeNameCommand);
  }
  if (params.changeDescriptionCommand) {
    request.setChangeDescriptionCommand(params.changeDescriptionCommand);
  }
  if (params.changePeriodCommand) {
    request.setChangeExperimentPeriodCommand(params.changePeriodCommand);
  }
  await setupAuthToken();
  await grpc.updateExperiment(request);
});

export interface ArchiveExperimentParams {
  environmentNamespace: string;
  id: string;
}

export const archiveExperiment = createAsyncThunk<
  void,
  ArchiveExperimentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/archive`, async (params) => {
  const request = new ArchiveExperimentRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new ArchiveExperimentCommand());
  await setupAuthToken();
  await grpc.archiveExperiment(request);
});

export interface StopExperimentParams {
  environmentNamespace: string;
  experimentId: string;
}

export const stopExperiment = createAsyncThunk<
  void,
  StopExperimentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/stop`, async (params) => {
  const request = new StopExperimentRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.experimentId);
  request.setCommand(new StopExperimentCommand());
  await setupAuthToken();
  await grpc.stopExperiment(request);
});

export const experimentsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listExperiments.pending, (state) => {
        state.loading = true;
      })
      .addCase(listExperiments.fulfilled, (state, action) => {
        experimentsAdapter.removeAll(state);
        experimentsAdapter.upsertMany(state, action.payload.experimentsList);
        state.totalCount = action.payload.totalCount;
        state.loading = false;
      })
      .addCase(listExperiments.rejected, (state) => {
        state.loading = false;
      })
      .addCase(getExperiment.pending, (state) => {
        state.getExperimentError = null;
      })
      .addCase(getExperiment.fulfilled, (state, action) => {
        state.getExperimentError = null;
        if (action.payload) {
          experimentsAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(getExperiment.rejected, (state, action) => {
        state.getExperimentError = action.error;
      })
      .addCase(createExperiment.pending, (state) => {})
      .addCase(createExperiment.fulfilled, (state, action) => {})
      .addCase(createExperiment.rejected, (state, action) => {});
  },
});
