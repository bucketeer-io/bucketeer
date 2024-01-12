import {
  CreateFlagTriggerCommand,
  DeleteFlagTriggerCommand,
  ChangeFlagTriggerDescriptionCommand,
  ResetFlagTriggerCommand,
  DisableFlagTriggerCommand,
  EnableFlagTriggerCommand,
} from '@/proto/feature/command_pb';
import { FlagTrigger } from '@/proto/feature/flag_trigger_pb';
import {
  CreateFlagTriggerRequest,
  CreateFlagTriggerResponse,
  ListFlagTriggersRequest,
  ListFlagTriggersResponse,
  DeleteFlagTriggerRequest,
  UpdateFlagTriggerRequest,
  ResetFlagTriggerRequest,
  DisableFlagTriggerRequest,
  EnableFlagTriggerRequest,
  ResetFlagTriggerResponse,
} from '@/proto/feature/service_pb';
import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from '@reduxjs/toolkit';

import * as flagTriggersGrpc from '../grpc/flagTriggers';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'flagTriggers';

export const flagTriggersAdapter =
  createEntityAdapter<ListFlagTriggersResponse.FlagTriggerWithUrl.AsObject>({
    selectId: (flagTriggerWithUrl) => flagTriggerWithUrl.flagTrigger.id,
  });

export const { selectAll, selectById } = flagTriggersAdapter.getSelectors();

interface CreateFlagTriggerParams {
  environmentNamespace: string;
  featureId: string;
  action: FlagTrigger.ActionMap[keyof FlagTrigger.ActionMap];
  triggerType: FlagTrigger.TypeMap[keyof FlagTrigger.TypeMap];
  description: string;
}

export const createFlagTrigger = createAsyncThunk<
  CreateFlagTriggerResponse.AsObject,
  CreateFlagTriggerParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateFlagTriggerRequest();
  const command = new CreateFlagTriggerCommand();
  request.setEnvironmentNamespace(params.environmentNamespace);

  command.setFeatureId(params.featureId);
  command.setAction(params.action);
  command.setType(params.triggerType);
  command.setDescription(params.description);

  request.setCreateFlagTriggerCommand(command);
  await setupAuthToken();
  const result = await flagTriggersGrpc.createFlagTrigger(request);
  return result.response.toObject();
});

export interface ListFlagTriggersParams {
  environmentNamespace: string;
  featureId: string;
}

export const listFlagTriggers = createAsyncThunk<
  ListFlagTriggersResponse.AsObject,
  ListFlagTriggersParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListFlagTriggersRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setFeatureId(params.featureId);
  request.setOrderDirection(ListFlagTriggersRequest.OrderDirection.DESC);
  await setupAuthToken();
  const result = await flagTriggersGrpc.listFlagTriggers(request);
  return result.response.toObject();
});

export interface UpdateFlagTriggerParams {
  environmentNamespace: string;
  id: string;
  description: string;
}

export const updateFlagTrigger = createAsyncThunk<
  void,
  UpdateFlagTriggerParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const command = new ChangeFlagTriggerDescriptionCommand();
  command.setDescription(params.description);

  const request = new UpdateFlagTriggerRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setChangeFlagTriggerDescriptionCommand(command);

  await setupAuthToken();
  await flagTriggersGrpc.updateFlagTrigger(request);
});

export interface DeleteFlagTriggerParams {
  environmentNamespace: string;
  id: string;
}

export const deleteFlagTrigger = createAsyncThunk<
  void,
  DeleteFlagTriggerParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeleteFlagTriggerRequest();
  request.setId(params.id);
  request.setEnvironmentNamespace(params.environmentNamespace);
  const command = new DeleteFlagTriggerCommand();
  request.setDeleteFlagTriggerCommand(command);
  await setupAuthToken();
  await flagTriggersGrpc.deleteFlagTrigger(request);
});

export interface ResetFlagTriggerParams {
  environmentNamespace: string;
  id: string;
}

export const resetFlagTrigger = createAsyncThunk<
  ResetFlagTriggerResponse.AsObject,
  ResetFlagTriggerParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/reset`, async (params) => {
  const request = new ResetFlagTriggerRequest();
  request.setId(params.id);
  request.setEnvironmentNamespace(params.environmentNamespace);
  const command = new ResetFlagTriggerCommand();
  request.setResetFlagTriggerCommand(command);
  await setupAuthToken();
  const result = await flagTriggersGrpc.resetFlagTrigger(request);
  return result.response.toObject();
});

export interface EnableFlagTriggerParams {
  environmentNamespace: string;
  id: string;
}

export const enableFlagTrigger = createAsyncThunk<
  void,
  EnableFlagTriggerParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  const request = new EnableFlagTriggerRequest();
  request.setId(params.id);
  request.setEnvironmentNamespace(params.environmentNamespace);
  const command = new EnableFlagTriggerCommand();
  request.setEnableFlagTriggerCommand(command);
  await setupAuthToken();
  await flagTriggersGrpc.enableFlagTrigger(request);
});

export interface DisableFlagTriggerParams {
  environmentNamespace: string;
  id: string;
}

export const disableFlagTrigger = createAsyncThunk<
  void,
  DisableFlagTriggerParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  const request = new DisableFlagTriggerRequest();
  request.setId(params.id);
  request.setEnvironmentNamespace(params.environmentNamespace);
  const command = new DisableFlagTriggerCommand();
  request.setDisableFlagTriggerCommand(command);
  await setupAuthToken();
  await flagTriggersGrpc.disableFlagTrigger(request);
});

const initialState = flagTriggersAdapter.getInitialState<{
  loading: boolean;
}>({
  loading: false,
});

export type flagTriggersState = typeof initialState;

export const flagTriggersSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listFlagTriggers.pending, (state) => {
        state.loading = true;
      })
      .addCase(listFlagTriggers.fulfilled, (state, action) => {
        flagTriggersAdapter.removeAll(state);
        flagTriggersAdapter.upsertMany(state, action.payload.flagTriggersList);
        state.loading = false;
      })
      .addCase(listFlagTriggers.rejected, (state) => {
        state.loading = false;
      })
      .addCase(createFlagTrigger.pending, (state) => {
        state.loading = true;
      })
      .addCase(createFlagTrigger.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(createFlagTrigger.rejected, (state, action) => {
        state.loading = false;
      })
      .addCase(deleteFlagTrigger.pending, (state) => {
        state.loading = true;
      })
      .addCase(deleteFlagTrigger.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(deleteFlagTrigger.rejected, (state) => {
        state.loading = false;
      })
      .addCase(updateFlagTrigger.pending, (state) => {
        state.loading = true;
      })
      .addCase(updateFlagTrigger.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(updateFlagTrigger.rejected, (state) => {
        state.loading = false;
      })
      .addCase(resetFlagTrigger.pending, (state) => {
        state.loading = true;
      })
      .addCase(resetFlagTrigger.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(resetFlagTrigger.rejected, (state) => {
        state.loading = false;
      })
      .addCase(disableFlagTrigger.pending, (state) => {
        state.loading = true;
      })
      .addCase(disableFlagTrigger.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(disableFlagTrigger.rejected, (state) => {
        state.loading = false;
      });
  },
});
