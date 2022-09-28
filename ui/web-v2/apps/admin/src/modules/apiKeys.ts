import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';
import { BoolValue } from 'google-protobuf/google/protobuf/wrappers_pb';

import * as grpc from '../grpc/apikey';
import { APIKey } from '../proto/account/api_key_pb';
import {
  EnableAPIKeyCommand,
  DisableAPIKeyCommand,
  CreateAPIKeyCommand,
  ChangeAPIKeyNameCommand,
} from '../proto/account/command_pb';
import {
  ListAPIKeysRequest,
  ListAPIKeysResponse,
  GetAPIKeyRequest,
  EnableAPIKeyRequest,
  DisableAPIKeyRequest,
  CreateAPIKeyRequest,
  ChangeAPIKeyNameRequest,
} from '../proto/account/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'apiKeys';

export const apiKeysAdapter = createEntityAdapter<APIKey.AsObject>({
  selectId: (apykey) => apykey.id,
});

export const { selectAll, selectById } = apiKeysAdapter.getSelectors();

export const listAPIKeys = createAsyncThunk<
  ListAPIKeysResponse.AsObject,
  APIKeyParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListAPIKeysRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  params.disabled != null &&
    request.setDisabled(new BoolValue().setValue(params.disabled));

  await setupAuthToken();
  const result = await grpc.listAPIKeys(request);
  return result.response.toObject();
});

export type OrderBy =
  ListAPIKeysRequest.OrderByMap[keyof ListAPIKeysRequest.OrderByMap];
export type OrderDirection =
  ListAPIKeysRequest.OrderDirectionMap[keyof ListAPIKeysRequest.OrderDirectionMap];

interface APIKeyParams {
  environmentNamespace: string;
  pageSize: number;
  cursor: string;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchKeyword: string;
  disabled?: boolean;
}

export interface GetAPIKeyParams {
  environmentNamespace: string;
  id: string;
}

export const getAPIKey = createAsyncThunk<
  APIKey.AsObject,
  GetAPIKeyParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetAPIKeyRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  await setupAuthToken();
  const result = await grpc.getAPIKey(request);
  return result.response.toObject().apiKey;
});

export interface EnableAPIKeyParams {
  environmentNamespace: string;
  id: string;
}

export const enableAPIKey = createAsyncThunk<
  void,
  EnableAPIKeyParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  const request = new EnableAPIKeyRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new EnableAPIKeyCommand());
  await setupAuthToken();
  await grpc.enableAPIKey(request);
});

export interface DisableAPIKeyParams {
  environmentNamespace: string;
  id: string;
}

export const disableAPIKey = createAsyncThunk<
  void,
  DisableAPIKeyParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  const request = new DisableAPIKeyRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new DisableAPIKeyCommand());
  await setupAuthToken();
  await grpc.disableAPIKey(request);
});

export interface CreateAPIKeyParams {
  environmentNamespace: string;
  name: string;
}

export const createAPIKey = createAsyncThunk<
  void,
  CreateAPIKeyParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/add`, async (params) => {
  const request = new CreateAPIKeyRequest();
  const cmd = new CreateAPIKeyCommand();
  cmd.setName(params.name);
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setCommand(cmd);
  await setupAuthToken();
  await grpc.createAPIKey(request);
});

export interface updateAPIKeyParams {
  environmentNamespace: string;
  id: string;
  name: string;
}

export const updateAPIKey = createAsyncThunk<
  void,
  updateAPIKeyParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new ChangeAPIKeyNameRequest();
  const cmd = new ChangeAPIKeyNameCommand();
  cmd.setName(params.name);
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(cmd);
  await setupAuthToken();
  await grpc.changeAPIKeyName(request);
});

const initialState = apiKeysAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getAPIKeyError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getAPIKeyError: null,
});

export const apiKeySlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listAPIKeys.pending, (state) => {
        state.loading = true;
      })
      .addCase(listAPIKeys.fulfilled, (state, action) => {
        apiKeysAdapter.removeAll(state);
        apiKeysAdapter.upsertMany(state, action.payload.apiKeysList);
        state.totalCount = action.payload.totalCount;
        state.loading = false;
      })
      .addCase(listAPIKeys.rejected, (state) => {
        state.loading = false;
      })
      .addCase(getAPIKey.pending, (state) => {
        state.getAPIKeyError = null;
      })
      .addCase(getAPIKey.fulfilled, (state, action) => {
        state.getAPIKeyError = null;
        if (action.payload) {
          apiKeysAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(getAPIKey.rejected, (state, action) => {
        state.getAPIKeyError = action.error;
      })
      .addCase(enableAPIKey.pending, (state) => {})
      .addCase(enableAPIKey.fulfilled, (state, action) => {})
      .addCase(enableAPIKey.rejected, (state, action) => {})
      .addCase(disableAPIKey.pending, (state) => {})
      .addCase(disableAPIKey.fulfilled, (state, action) => {})
      .addCase(disableAPIKey.rejected, (state, action) => {})
      .addCase(createAPIKey.pending, (state) => {})
      .addCase(createAPIKey.fulfilled, (state, action) => {})
      .addCase(createAPIKey.rejected, (state, action) => {})
      .addCase(updateAPIKey.pending, (state) => {})
      .addCase(updateAPIKey.fulfilled, (state, action) => {})
      .addCase(updateAPIKey.rejected, (state, action) => {});
  },
});
