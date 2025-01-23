import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError
} from '@reduxjs/toolkit';
import { Message } from 'google-protobuf';
import { Any } from 'google-protobuf/google/protobuf/any_pb';

import * as tagsGrpc from '../grpc/tags';
import { Command } from '../proto/feature/command_pb';
import { Tag } from '../proto/tag/tag_pb';
import { ListTagsRequest, ListTagsResponse } from '../proto/tag/service_pb';

import { AppState } from '.';

const MODULE_NAME = 'tags';

export const tagsAdapter = createEntityAdapter({
  selectId: (tag: Tag.AsObject) => tag.id
});

export const { selectAll, selectById } = tagsAdapter.getSelectors();

export type OrderBy =
  ListTagsRequest.OrderByMap[keyof ListTagsRequest.OrderByMap];
export type OrderDirection =
  ListTagsRequest.OrderDirectionMap[keyof ListTagsRequest.OrderDirectionMap];
export type TagEntityType = Tag.EntityTypeMap[keyof Tag.EntityTypeMap];

export interface ListTagsParams {
  environmentId: string;
  pageSize: number;
  cursor: string;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchKeyword: string;
  entityType: TagEntityType;
}

export const listTags = createAsyncThunk<
  ListTagsResponse.AsObject,
  ListTagsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/listTags`, async (params) => {
  const request = new ListTagsRequest();
  request.setEnvironmentId(params.environmentId);
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  request.setEntityType(params.entityType);

  const result = await tagsGrpc.listTags(request);
  return result.response.toObject();
});

const initialState = tagsAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getFeatureError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getFeatureError: null
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
  }
});

export const createCommand = ({
  message,
  name
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
