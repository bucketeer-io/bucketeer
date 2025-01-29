import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError
} from '@reduxjs/toolkit';

import * as grpc from '../grpc/codeRefs';
import { APIKey } from '../proto/account/api_key_pb';
import { CodeReference } from '../proto/coderef/code_reference_pb';
import {
  EnableAPIKeyCommand,
  DisableAPIKeyCommand,
  CreateAPIKeyCommand,
  ChangeAPIKeyNameCommand
} from '../proto/account/command_pb';
import {
  ListAPIKeysRequest,
  ListAPIKeysResponse,
  GetAPIKeyRequest,
  EnableAPIKeyRequest,
  DisableAPIKeyRequest,
  CreateAPIKeyRequest,
  ChangeAPIKeyNameRequest
} from '../proto/account/service_pb';

import { AppState } from '.';
import {
  GetCodeReferenceRequest,
  GetCodeReferenceResponse,
  ListCodeReferencesRequest,
  ListCodeReferencesResponse
} from '../proto/coderef/service_pb';

const MODULE_NAME = 'codeRefs';

export const codeRefsAdapter = createEntityAdapter({
  selectId: (codeRef: CodeReference.AsObject) => codeRef.id
});

export const { selectAll, selectById } = codeRefsAdapter.getSelectors();

interface ListCodeReferenceParams {
  environmentId: string;
  featureId: string;
  pageSize: number;
  fileExtension: string;
  repositoryBranch: string;
  repositoryType: CodeReference.RepositoryTypeMap[keyof CodeReference.RepositoryTypeMap];
}

export const listCodeRefs = createAsyncThunk<
  ListCodeReferencesResponse.AsObject,
  ListCodeReferenceParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListCodeReferencesRequest();

  console.log({
    params
  });

  request.setEnvironmentId(params.environmentId);
  request.setFeatureId(params.featureId);
  request.setPageSize(params.pageSize);
  params.repositoryBranch &&
    request.setRepositoryBranch(params.repositoryBranch);
  params.repositoryType && request.setRepositoryType(params.repositoryType);
  params.fileExtension && request.setFileExtension(params.fileExtension);

  const result = await grpc.listCodeRefs(request);
  return result.response.toObject();
});

interface CodeReferenceParams {
  environmentId: string;
  id: string;
}

export const getCodeRefs = createAsyncThunk<
  GetCodeReferenceResponse.AsObject,
  CodeReferenceParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetCodeReferenceRequest();

  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);

  const result = await grpc.getCodeRef(request);

  return result.response.toObject();
});

export type OrderBy =
  ListAPIKeysRequest.OrderByMap[keyof ListAPIKeysRequest.OrderByMap];
export type OrderDirection =
  ListAPIKeysRequest.OrderDirectionMap[keyof ListAPIKeysRequest.OrderDirectionMap];

export interface EnableAPIKeyParams {
  environmentId: string;
  id: string;
}

export const enableAPIKey = createAsyncThunk<
  void,
  EnableAPIKeyParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  const request = new EnableAPIKeyRequest();
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setCommand(new EnableAPIKeyCommand());
  // await grpc.enableAPIKey(request);
});

export interface DisableAPIKeyParams {
  environmentId: string;
  id: string;
}

export const disableAPIKey = createAsyncThunk<
  void,
  DisableAPIKeyParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  const request = new DisableAPIKeyRequest();
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setCommand(new DisableAPIKeyCommand());
  // await grpc.disableAPIKey(request);
});

export interface CreateAPIKeyParams {
  environmentId: string;
  name: string;
  role: APIKey.RoleMap[keyof APIKey.RoleMap];
}

export const createAPIKey = createAsyncThunk<
  void,
  CreateAPIKeyParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/add`, async (params) => {
  const request = new CreateAPIKeyRequest();
  const cmd = new CreateAPIKeyCommand();
  cmd.setName(params.name);
  cmd.setRole(params.role);
  request.setEnvironmentId(params.environmentId);
  request.setCommand(cmd);
  // await grpc.createAPIKey(request);
});

export interface updateAPIKeyParams {
  environmentId: string;
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
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setCommand(cmd);
  // await grpc.changeAPIKeyName(request);
});

const initialState = codeRefsAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getAPIKeyError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getAPIKeyError: null
});

export const codeRefsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listCodeRefs.pending, (state) => {
        state.loading = true;
      })
      .addCase(listCodeRefs.fulfilled, (state, action) => {
        codeRefsAdapter.removeAll(state);
        codeRefsAdapter.upsertMany(state, action.payload.codeReferencesList);
        state.totalCount = action.payload.totalCount;
        state.loading = false;
      })
      .addCase(listCodeRefs.rejected, (state) => {
        state.loading = false;
      })
      // .addCase(getAPIKey.pending, (state) => {
      //   state.getAPIKeyError = null;
      // })
      // .addCase(getAPIKey.fulfilled, (state, action) => {
      //   state.getAPIKeyError = null;
      //   if (action.payload) {
      //     codeRefsAdapter.upsertOne(state, action.payload);
      //   }
      // })
      // .addCase(getAPIKey.rejected, (state, action) => {
      //   state.getAPIKeyError = action.error;
      // })
      .addCase(enableAPIKey.pending, () => {})
      .addCase(enableAPIKey.fulfilled, () => {})
      .addCase(enableAPIKey.rejected, () => {})
      .addCase(disableAPIKey.pending, () => {})
      .addCase(disableAPIKey.fulfilled, () => {})
      .addCase(disableAPIKey.rejected, () => {})
      .addCase(createAPIKey.pending, () => {})
      .addCase(createAPIKey.fulfilled, () => {})
      .addCase(createAPIKey.rejected, () => {})
      .addCase(updateAPIKey.pending, () => {})
      .addCase(updateAPIKey.fulfilled, () => {})
      .addCase(updateAPIKey.rejected, () => {});
  }
});
