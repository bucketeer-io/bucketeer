import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from '@reduxjs/toolkit';

import * as grpc from '../grpc/eventcounter';
import {
  ListUserMetadataRequest,
  ListUserMetadataResponse,
} from '../proto/eventcounter/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'userMetadata';

export const userMetadataAdapter = createEntityAdapter<string>({});

export const { selectAll } = userMetadataAdapter.getSelectors();

export interface ListUserMetadataParams {
  environmentNamespace: string;
}

export const listUserMetadata = createAsyncThunk<
  ListUserMetadataResponse.AsObject,
  ListUserMetadataParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/listUserMetadata`, async (params) => {
  const request = new ListUserMetadataRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  await setupAuthToken();
  const result = await grpc.listUserMetadata(request);
  return result.response.toObject();
});

const initialState = userMetadataAdapter.getInitialState<{
  loading: boolean;
}>({
  loading: false,
});

export const userMetadataSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listUserMetadata.pending, (state) => {
        state.loading = true;
      })
      .addCase(listUserMetadata.fulfilled, (state, action) => {
        userMetadataAdapter.removeAll(state);
        userMetadataAdapter.upsertMany(state, action.payload.dataList);
        state.loading = false;
      })
      .addCase(listUserMetadata.rejected, (state) => {
        state.loading = false;
      });
  },
});
