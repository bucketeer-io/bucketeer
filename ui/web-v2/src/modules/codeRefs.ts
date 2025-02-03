import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk
} from '@reduxjs/toolkit';

import * as grpc from '../grpc/codeRefs';
import { CodeReference } from '../proto/coderef/code_reference_pb';
import { AppState } from '.';
import {
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

const initialState = codeRefsAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
}>({
  loading: false,
  totalCount: 0
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
      });
  }
});
