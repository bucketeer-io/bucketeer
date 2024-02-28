import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from '@reduxjs/toolkit';

import * as progressiveRolloutGrpc from '../grpc/progressiveRollout';
import {
  CreateProgressiveRolloutCommand,
  DeleteProgressiveRolloutCommand,
  StopProgressiveRolloutCommand,
} from '../proto/autoops/command_pb';
import { ProgressiveRollout } from '../proto/autoops/progressive_rollout_pb';
import {
  CreateProgressiveRolloutRequest,
  ListProgressiveRolloutsRequest,
  ListProgressiveRolloutsResponse,
  DeleteProgressiveRolloutRequest,
  StopProgressiveRolloutRequest,
} from '../proto/autoops/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'progressiveRollout';

export const progressiveRolloutAdapter =
  createEntityAdapter<ProgressiveRollout.AsObject>({
    selectId: (progressiveRollout) => progressiveRollout.id,
  });

export const { selectAll, selectById } =
  progressiveRolloutAdapter.getSelectors();

interface CreateProgressiveRolloutParams {
  environmentNamespace: string;
  command: CreateProgressiveRolloutCommand;
}

export const createProgressiveRollout = createAsyncThunk<
  void,
  CreateProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateProgressiveRolloutRequest();

  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setCommand(params.command);
  await setupAuthToken();
  await progressiveRolloutGrpc.createProgressiveRollout(request);
});

export interface ListProgressiveRolloutParams {
  environmentNamespace: string;
  featureId: string;
}

export const listProgressiveRollout = createAsyncThunk<
  ListProgressiveRolloutsResponse.AsObject,
  ListProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListProgressiveRolloutsRequest();
  request.setFeatureIdsList([params.featureId]);
  request.setEnvironmentNamespace(params.environmentNamespace);
  await setupAuthToken();
  const result = await progressiveRolloutGrpc.listProgressiveRollouts(request);
  return result.response.toObject();
});

export interface DeleteProgressiveRolloutParams {
  environmentNamespace: string;
  id: string;
}

export const deleteProgressiveRollout = createAsyncThunk<
  void,
  DeleteProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeleteProgressiveRolloutRequest();
  request.setId(params.id);
  request.setEnvironmentNamespace(params.environmentNamespace);
  const command = new DeleteProgressiveRolloutCommand();
  request.setCommand(command);
  await setupAuthToken();
  await progressiveRolloutGrpc.deleteProgressiveRollout(request);
});

export interface StopProgressiveRolloutParams {
  environmentNamespace: string;
  id: string;
}

export const stopProgressiveRollout = createAsyncThunk<
  void,
  StopProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/stop`, async (params) => {
  const request = new StopProgressiveRolloutRequest();
  request.setId(params.id);
  request.setEnvironmentNamespace(params.environmentNamespace);
  const command = new StopProgressiveRolloutCommand();
  command.setStoppedBy(ProgressiveRollout.StoppedBy.USER);
  request.setCommand(command);
  await setupAuthToken();
  await progressiveRolloutGrpc.stopProgressiveRollout(request);
});

const initialState = progressiveRolloutAdapter.getInitialState<{
  loading: boolean;
}>({
  loading: false,
});

export type ProgressiveRolloutState = typeof initialState;

export const progressiveRolloutSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listProgressiveRollout.pending, (state) => {
        state.loading = true;
      })
      .addCase(listProgressiveRollout.fulfilled, (state, action) => {
        progressiveRolloutAdapter.removeAll(state);
        progressiveRolloutAdapter.upsertMany(
          state,
          action.payload.progressiveRolloutsList
        );
        state.loading = false;
      })
      .addCase(listProgressiveRollout.rejected, (state) => {
        state.loading = false;
      })
      .addCase(deleteProgressiveRollout.pending, (state) => {
        state.loading = true;
      })
      .addCase(deleteProgressiveRollout.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(deleteProgressiveRollout.rejected, (state) => {
        state.loading = false;
      });
  },
});
