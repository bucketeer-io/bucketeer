import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from '@reduxjs/toolkit';

import * as webhookGrpc from '../grpc/webhook';
import {
  CreateWebhookCommand,
  DeleteWebhookCommand,
  ChangeWebhookNameCommand,
  ChangeWebhookDescriptionCommand,
} from '../proto/autoops/command_pb';
import {
  ListWebhooksRequest,
  CreateWebhookRequest,
  UpdateWebhookRequest,
  DeleteWebhookRequest,
  ListWebhooksResponse,
  GetWebhookRequest,
  GetWebhookResponse,
} from '../proto/autoops/service_pb';
import { Webhook } from '../proto/autoops/webhook_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'webhooks';

export const webhookAdapter = createEntityAdapter<Webhook.AsObject>({
  selectId: (webhook) => webhook.id,
});

export const { selectAll, selectById } = webhookAdapter.getSelectors();

const initialState = webhookAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  webhookLoading: boolean;
  webhookUrl: string;
}>({
  loading: false,
  totalCount: 0,
  webhookLoading: false,
  webhookUrl: '',
});

export type OrderBy =
  ListWebhooksRequest.OrderByMap[keyof ListWebhooksRequest.OrderByMap];
export type OrderDirection =
  ListWebhooksRequest.OrderDirectionMap[keyof ListWebhooksRequest.OrderDirectionMap];

interface ListWebhooksParams {
  environmentNamespace: string;
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
}

export const listWebhooks = createAsyncThunk<
  ListWebhooksResponse.AsObject,
  ListWebhooksParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListWebhooksRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  await setupAuthToken();
  const result = await webhookGrpc.listWebhooks(request);
  return result.response.toObject();
});

export interface CreateWebhookParams {
  environmentNamespace: string;
  name: string;
  description: string;
}

export const createWebhook = createAsyncThunk<
  void,
  CreateWebhookParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const cmd = new CreateWebhookCommand();
  cmd.setName(params.name);
  cmd.setDescription(params.description);
  const request = new CreateWebhookRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setCommand(cmd);
  await setupAuthToken();
  await webhookGrpc.createWebhook(request);
});

export interface UpdateWebhookParams {
  environmentNamespace: string;
  id: string;
  name: String;
  description: String;
}

export const updateWebhook = createAsyncThunk<
  void,
  UpdateWebhookParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new UpdateWebhookRequest();
  if (params.name) {
    const cmd = new ChangeWebhookNameCommand();
    cmd.setName(params.name.toString());
    request.setChangewebhooknamecommand(cmd);
  }
  if (params.description) {
    const cmd = new ChangeWebhookDescriptionCommand();
    cmd.setDescription(params.description.toString());
    request.setChangewebhookdescriptioncommand(cmd);
  }
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  await setupAuthToken();
  await webhookGrpc.updateWebhook(request);
});

export interface DeleteWebhookParams {
  environmentNamespace: string;
  id: string;
}

export const deleteWebhook = createAsyncThunk<
  void,
  DeleteWebhookParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeleteWebhookRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new DeleteWebhookCommand());
  await setupAuthToken();
  await webhookGrpc.deleteWebhook(request);
});

export interface GetWebhookParams {
  environmentNamespace: string;
  id: string;
}

export const getWebhook = createAsyncThunk<
  GetWebhookResponse.AsObject,
  GetWebhookParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetWebhookRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  await setupAuthToken();
  const result = await webhookGrpc.getWebhook(request);
  return result.response.toObject();
});

export type WebhooksState = typeof initialState;

export const webhookSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listWebhooks.pending, (state) => {
        state.loading = true;
      })
      .addCase(listWebhooks.fulfilled, (state, action) => {
        webhookAdapter.removeAll(state);
        webhookAdapter.upsertMany(state, action.payload.webhooksList);
        state.loading = false;
        state.totalCount = action.payload.totalCount;
      })
      .addCase(listWebhooks.rejected, (state) => {
        state.loading = false;
      })
      .addCase(getWebhook.pending, (state) => {
        state.webhookLoading = true;
      })
      .addCase(getWebhook.fulfilled, (state, action) => {
        state.webhookUrl = action.payload.url;
        state.webhookLoading = false;
      })
      .addCase(getWebhook.rejected, (state) => {
        state.webhookLoading = false;
      })
      .addCase(createWebhook.pending, (state) => {})
      .addCase(createWebhook.fulfilled, (state, action) => {})
      .addCase(createWebhook.rejected, (state, action) => {})
      .addCase(deleteWebhook.pending, (state) => {})
      .addCase(deleteWebhook.fulfilled, (state, action) => {})
      .addCase(deleteWebhook.rejected, (state, action) => {});
  },
});
