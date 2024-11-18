import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk
} from '@reduxjs/toolkit';

import * as progressiveRolloutGrpc from '../grpc/progressiveRollout';
import {
  CreateProgressiveRolloutCommand,
  DeleteProgressiveRolloutCommand,
  StopProgressiveRolloutCommand
} from '../proto/autoops/command_pb';
import { ProgressiveRollout } from '../proto/autoops/progressive_rollout_pb';
import {
  CreateProgressiveRolloutRequest,
  ListProgressiveRolloutsRequest,
  ListProgressiveRolloutsResponse,
  DeleteProgressiveRolloutRequest,
  StopProgressiveRolloutRequest
} from '../proto/autoops/service_pb';

import { AppState } from '.';

const MODULE_NAME = 'progressiveRollout';

export const progressiveRolloutAdapter = createEntityAdapter({
  selectId: (progressiveRollout: ProgressiveRollout.AsObject) =>
    progressiveRollout.id
});

export const { selectAll, selectById } =
  progressiveRolloutAdapter.getSelectors();

interface CreateProgressiveRolloutParams {
  environmentId: string;
  command: CreateProgressiveRolloutCommand;
}

export const createProgressiveRollout = createAsyncThunk<
  void,
  CreateProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateProgressiveRolloutRequest();

  request.setEnvironmentId(params.environmentId);
  request.setCommand(params.command);
  await progressiveRolloutGrpc.createProgressiveRollout(request);
});

export interface ListProgressiveRolloutParams {
  environmentId: string;
  featureId: string;
}

export const listProgressiveRollout = createAsyncThunk<
  ListProgressiveRolloutsResponse.AsObject,
  ListProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListProgressiveRolloutsRequest();
  request.setFeatureIdsList([params.featureId]);
  request.setEnvironmentId(params.environmentId);
  const result = await progressiveRolloutGrpc.listProgressiveRollouts(request);
  return result.response.toObject();
});

export interface DeleteProgressiveRolloutParams {
  environmentId: string;
  id: string;
}

export const deleteProgressiveRollout = createAsyncThunk<
  void,
  DeleteProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeleteProgressiveRolloutRequest();
  request.setId(params.id);
  request.setEnvironmentId(params.environmentId);
  const command = new DeleteProgressiveRolloutCommand();
  request.setCommand(command);
  await progressiveRolloutGrpc.deleteProgressiveRollout(request);
});

export interface StopProgressiveRolloutParams {
  environmentId: string;
  id: string;
}

export const stopProgressiveRollout = createAsyncThunk<
  void,
  StopProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/stop`, async (params) => {
  const request = new StopProgressiveRolloutRequest();
  request.setId(params.id);
  request.setEnvironmentId(params.environmentId);
  const command = new StopProgressiveRolloutCommand();
  command.setStoppedBy(ProgressiveRollout.StoppedBy.USER);
  request.setCommand(command);
  await progressiveRolloutGrpc.stopProgressiveRollout(request);
});

const initialState = progressiveRolloutAdapter.getInitialState<{
  loading: boolean;
}>({
  loading: false
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
      .addCase(deleteProgressiveRollout.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(deleteProgressiveRollout.rejected, (state) => {
        state.loading = false;
      });
  }
});
