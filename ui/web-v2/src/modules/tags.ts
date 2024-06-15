import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';
import { Message } from 'google-protobuf';
import { Any } from 'google-protobuf/google/protobuf/any_pb';

import * as featureGrpc from '../grpc/features';
import { Command } from '../proto/feature/command_pb';
import { Tag } from '../proto/feature/feature_pb';
import { ListTagsRequest, ListTagsResponse } from '../proto/feature/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'tags';

export const tagsAdapter = createEntityAdapter<Tag.AsObject>({
  selectId: (tag) => tag.id,
});

export const { selectAll, selectById } = tagsAdapter.getSelectors();

export type OrderBy =
  ListTagsRequest.OrderByMap[keyof ListTagsRequest.OrderByMap];
export type OrderDirection =
  ListTagsRequest.OrderDirectionMap[keyof ListTagsRequest.OrderDirectionMap];

export interface ListTagsParams {
  environmentNamespace: string;
  pageSize: number;
  cursor: string;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchKeyword: string;
}

export const listTags = createAsyncThunk<
  ListTagsResponse.AsObject,
  ListTagsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/listTags`, async (params) => {
  const request = new ListTagsRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);

  await setupAuthToken();
  const result = await featureGrpc.listTags(request);
  return result.response.toObject();
});

const initialState = tagsAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getFeatureError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getFeatureError: null,
});

export const tagsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listTags.pending, (state) => {
        state.loading = true;
      })
      .addCase(listTags.fulfilled, (state, action) => {
        tagsAdapter.removeAll(state);
        tagsAdapter.upsertMany(state, action.payload.tagsList);
        state.loading = false;
      });
  },
});

export const createCommand = ({
  message,
  name,
}: {
  message: Message;
  name: string;
}): Command => {
  const command = new Command();
  const any = new Any();
  any.pack(message.serializeBinary(), `bucketeer.tags.${name}`);
  command.setCommand(any);
  return command;
};
