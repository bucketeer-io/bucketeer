import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';

import * as grpc from '../grpc/environment';
import {
  ChangeDescriptionEnvironmentCommand,
  CreateEnvironmentCommand,
} from '../proto/environment/command_pb';
import { Environment } from '../proto/environment/environment_pb';
import {
  CreateEnvironmentRequest,
  GetEnvironmentRequest,
  ListEnvironmentsRequest,
  ListEnvironmentsResponse,
  UpdateEnvironmentRequest,
} from '../proto/environment/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'environments';

export const environmentAdapter = createEntityAdapter<Environment.AsObject>({
  selectId: (e) => e.id,
});

export const { selectAll, selectById } = environmentAdapter.getSelectors();

const initialState = environmentAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getEnvironmentError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getEnvironmentError: null,
});

export type OrderBy =
  ListEnvironmentsRequest.OrderByMap[keyof ListEnvironmentsRequest.OrderByMap];
export type OrderDirection =
  ListEnvironmentsRequest.OrderDirectionMap[keyof ListEnvironmentsRequest.OrderDirectionMap];

interface ListEnvironmentsRequestParams {
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  projectId?: string;
}

export const listEnvironments = createAsyncThunk<
  ListEnvironmentsResponse.AsObject,
  ListEnvironmentsRequestParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListEnvironmentsRequest();
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  request.setProjectId(params.projectId);
  await setupAuthToken();
  const result = await grpc.listEnvironments(request);
  return result.response.toObject();
});

export interface GetEnvironmentParams {
  id: string;
}

export const getEnvironment = createAsyncThunk<
  Environment.AsObject,
  GetEnvironmentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetEnvironmentRequest();
  request.setId(params.id);
  await setupAuthToken();
  const result = await grpc.getEnvironment(request);
  return result.response.toObject().environment;
});

export interface CreateEnvironmentParams {
  id: string;
  projectId: string;
  description: string;
}

export const createEnvironment = createAsyncThunk<
  void,
  CreateEnvironmentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateEnvironmentRequest();
  const command = new CreateEnvironmentCommand();
  command.setId(params.id);
  command.setDescription(params.description);
  command.setProjectId(params.projectId);
  request.setCommand(command);
  await setupAuthToken();
  await grpc.createEnvironment(request);
});

export interface UpdateEnvironmentParams {
  id: string;
  description?: string;
}

export const updateEnvironment = createAsyncThunk<
  void,
  UpdateEnvironmentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new UpdateEnvironmentRequest();
  request.setId(params.id);
  const command = new ChangeDescriptionEnvironmentCommand();
  command.setDescription(params.description);
  request.setChangeDescriptionCommand(command);
  await setupAuthToken();
  await grpc.updateEnvironment(request);
});

export type EnvironmentsState = typeof initialState;

export const environmentsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listEnvironments.pending, (state) => {
        state.loading = true;
      })
      .addCase(listEnvironments.fulfilled, (state, action) => {
        environmentAdapter.removeAll(state);
        environmentAdapter.upsertMany(state, action.payload.environmentsList);
        state.loading = false;
        state.totalCount = action.payload.totalCount;
      })
      .addCase(listEnvironments.rejected, (state) => {
        state.loading = false;
      })
      .addCase(getEnvironment.pending, (state) => {
        state.getEnvironmentError = null;
      })
      .addCase(getEnvironment.fulfilled, (state, action) => {
        state.getEnvironmentError = null;
        if (action.payload) {
          environmentAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(getEnvironment.rejected, (state, action) => {
        state.getEnvironmentError = action.error;
      })
      .addCase(createEnvironment.pending, (state) => {})
      .addCase(createEnvironment.fulfilled, (state, action) => {})
      .addCase(createEnvironment.rejected, (state, action) => {})
      .addCase(updateEnvironment.pending, (state) => {})
      .addCase(updateEnvironment.fulfilled, (state, action) => {})
      .addCase(updateEnvironment.rejected, (state, action) => {});
  },
});
