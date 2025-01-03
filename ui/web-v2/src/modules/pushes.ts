import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk
} from '@reduxjs/toolkit';

import * as pushGrpc from '../grpc/push';
import {
  CreatePushCommand,
  AddPushTagsCommand,
  DeletePushTagsCommand,
  RenamePushCommand,
  DeletePushCommand
} from '../proto/push/command_pb';
import { Push } from '../proto/push/push_pb';
import {
  ListPushesRequest,
  CreatePushRequest,
  UpdatePushRequest,
  DeletePushRequest,
  ListPushesResponse
} from '../proto/push/service_pb';

import { AppState } from '.';

const MODULE_NAME = 'pushes';

export const pushAdapter = createEntityAdapter({
  selectId: (push: Push.AsObject) => push.id
});

export const { selectAll, selectById } = pushAdapter.getSelectors();

const initialState = pushAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
}>({
  loading: false,
  totalCount: 0
});

export type OrderBy =
  ListPushesRequest.OrderByMap[keyof ListPushesRequest.OrderByMap];
export type OrderDirection =
  ListPushesRequest.OrderDirectionMap[keyof ListPushesRequest.OrderDirectionMap];

interface ListPushesParams {
  environmentId: string;
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
}

export const listPushes = createAsyncThunk<
  ListPushesResponse.AsObject,
  ListPushesParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListPushesRequest();
  request.setEnvironmentId(params.environmentId);
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  const result = await pushGrpc.listPushes(request);
  return result.response.toObject();
});

export interface CreatePushParams {
  environmentId: string;
  name: string;
  tags: Array<string>;
  fcmServiceAccount: Uint8Array | string;
}

export const createPush = createAsyncThunk<
  void,
  CreatePushParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const cmd = new CreatePushCommand();
  cmd.setName(params.name);
  cmd.setTagsList(params.tags);
  cmd.setFcmServiceAccount(params.fcmServiceAccount);

  const request = new CreatePushRequest();
  request.setEnvironmentId(params.environmentId);
  request.setCommand(cmd);
  await pushGrpc.createPush(request);
});

export interface UpdatePushParams {
  environmentId: string;
  id: string;
  name: string;
  currentTags: Array<string>;
  tags: Array<string>;
}

export const updatePush = createAsyncThunk<
  void,
  UpdatePushParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new UpdatePushRequest();
  if (params.name) {
    const cmd = new RenamePushCommand();
    cmd.setName(params.name.toString());
    request.setRenamePushCommand(cmd);
  }
  if (params.tags) {
    const addPushTagList = params.tags.filter(
      (tag) => !params.currentTags.includes(tag)
    );
    if (addPushTagList.length > 0) {
      const cmd = new AddPushTagsCommand();
      cmd.setTagsList(addPushTagList);
      request.setAddPushTagsCommand(cmd);
    }
    const deletePushTagList = params.currentTags.filter(
      (tag) => !params.tags.includes(tag)
    );
    if (deletePushTagList.length > 0) {
      const cmd = new DeletePushTagsCommand();
      cmd.setTagsList(deletePushTagList);
      request.setDeletePushTagsCommand(cmd);
    }
  }
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  await pushGrpc.updatePush(request);
});

export interface DeletePushParams {
  environmentId: string;
  id: string;
}

export const deletePush = createAsyncThunk<
  void,
  DeletePushParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeletePushRequest();
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  request.setCommand(new DeletePushCommand());
  await pushGrpc.deletePush(request);
});

export type PushesState = typeof initialState;

export const pushSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listPushes.pending, (state) => {
        state.loading = true;
      })
      .addCase(listPushes.fulfilled, (state, action) => {
        pushAdapter.removeAll(state);
        pushAdapter.upsertMany(state, action.payload.pushesList);
        state.loading = false;
        state.totalCount = action.payload.totalCount;
      })
      .addCase(listPushes.rejected, (state) => {
        state.loading = false;
      })
      .addCase(createPush.pending, () => {})
      .addCase(createPush.fulfilled, () => {})
      .addCase(createPush.rejected, () => {})
      .addCase(deletePush.pending, () => {})
      .addCase(deletePush.fulfilled, () => {})
      .addCase(deletePush.rejected, () => {});
  }
});
