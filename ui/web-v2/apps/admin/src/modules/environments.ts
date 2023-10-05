import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';

import * as grpc from '../grpc/environment';
import {
  ChangeDescriptionEnvironmentV2Command,
  CreateEnvironmentV2Command,
  RenameEnvironmentV2Command
} from '../proto/environment/command_pb';
import { EnvironmentV2 } from '../proto/environment/environment_pb';
import {
  CreateEnvironmentV2Request,
  GetEnvironmentV2Request,
  ListEnvironmentsV2Request,
  ListEnvironmentsV2Response,
  UpdateEnvironmentV2Request,
} from '../proto/environment/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'environments';

export const environmentAdapter = createEntityAdapter<EnvironmentV2.AsObject>({
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
  ListEnvironmentsV2Request.OrderByMap[keyof ListEnvironmentsV2Request.OrderByMap];
export type OrderDirection =
  ListEnvironmentsV2Request.OrderDirectionMap[keyof ListEnvironmentsV2Request.OrderDirectionMap];

interface ListEnvironmentsRequestParams {
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  projectId?: string;
}

export const listEnvironments = createAsyncThunk<
  ListEnvironmentsV2Response.AsObject,
  ListEnvironmentsRequestParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListEnvironmentsV2Request();
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
  EnvironmentV2.AsObject,
  GetEnvironmentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetEnvironmentV2Request();
  request.setId(params.id);
  await setupAuthToken();
  const result = await grpc.getEnvironment(request);
  return result.response.toObject().environment;
});

export interface CreateEnvironmentParams {
  name: string;
  urlCode: string;
  projectId: string;
  description: string;
}

export const createEnvironment = createAsyncThunk<
  void,
  CreateEnvironmentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateEnvironmentV2Request();
  const command = new CreateEnvironmentV2Command();
  command.setName(params.name);
  command.setUrlCode(params.urlCode);
  command.setDescription(params.description);
  command.setProjectId(params.projectId);
  request.setCommand(command);
  await setupAuthToken();
  await grpc.createEnvironment(request);
});

export interface UpdateEnvironmentParams {
  id: string;
  name?: string
  description?: string;
}

export const updateEnvironment = createAsyncThunk<
  void,
  UpdateEnvironmentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params: UpdateEnvironmentParams) => {
  const request = new UpdateEnvironmentV2Request();
  request.setId(params.id);
  if (params.name) {
    const renameCommand = new RenameEnvironmentV2Command();
    renameCommand.setName(params.name);
    request.setRenameCommand(renameCommand);
  }
  if (params.description) {
    const changeDescCommand = new ChangeDescriptionEnvironmentV2Command();
    changeDescCommand.setDescription(params.description);
    request.setChangeDescriptionCommand(changeDescCommand);
  }
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
