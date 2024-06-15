import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from '@reduxjs/toolkit';
import { BoolValue } from 'google-protobuf/google/protobuf/wrappers_pb';

import * as subscriptionGrpc from '../grpc/subscription';
import {
  CreateSubscriptionCommand,
  DeleteSubscriptionCommand,
  EnableSubscriptionCommand,
  DisableSubscriptionCommand,
  AddSourceTypesCommand,
  DeleteSourceTypesCommand,
  RenameSubscriptionCommand,
} from '../proto/notification/command_pb';
import {
  Recipient,
  SlackChannelRecipient,
} from '../proto/notification/recipient_pb';
import {
  ListSubscriptionsRequest,
  ListSubscriptionsResponse,
  CreateSubscriptionRequest,
  UpdateSubscriptionRequest,
  DeleteSubscriptionRequest,
  EnableSubscriptionRequest,
  DisableSubscriptionRequest,
} from '../proto/notification/service_pb';
import { Subscription } from '../proto/notification/subscription_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'notifications';

export const notificationAdapter = createEntityAdapter<Subscription.AsObject>({
  selectId: (notification) => notification.id,
});

export const { selectAll, selectById } = notificationAdapter.getSelectors();

const initialState = notificationAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
}>({
  loading: false,
  totalCount: 0,
});

export type OrderBy =
  ListSubscriptionsRequest.OrderByMap[keyof ListSubscriptionsRequest.OrderByMap];
export type OrderDirection =
  ListSubscriptionsRequest.OrderDirectionMap[keyof ListSubscriptionsRequest.OrderDirectionMap];

interface ListNotificationParams {
  environmentNamespace: string;
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  disabled?: boolean;
}

export const listNotification = createAsyncThunk<
  ListSubscriptionsResponse.AsObject,
  ListNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListSubscriptionsRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  params.disabled != null &&
    request.setDisabled(new BoolValue().setValue(params.disabled));
  await setupAuthToken();
  const result = await subscriptionGrpc.listSubscriptions(request);
  return result.response.toObject();
});

export interface CreateNotificationParams {
  environmentNamespace: string;
  name: string;
  sourceTypes: Array<
    Subscription.SourceTypeMap[keyof Subscription.SourceTypeMap]
  >;
  webhookUrl: string;
}

export const createNotification = createAsyncThunk<
  void,
  CreateNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const cmd = new CreateSubscriptionCommand();
  cmd.setName(params.name);
  cmd.setSourceTypesList(params.sourceTypes);

  const recipient = new Recipient();
  recipient.setType(Recipient.Type.SLACKCHANNEL);

  const slackChannelRecipient = new SlackChannelRecipient();
  slackChannelRecipient.setWebhookUrl(params.webhookUrl);
  recipient.setSlackChannelRecipient(slackChannelRecipient);
  cmd.setRecipient(recipient);

  // TODO: We need to implement this on the admin console
  // so the user can choose the language
  recipient.setLanguage(1); // Japanese
  cmd.setRecipient(recipient);

  const request = new CreateSubscriptionRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setCommand(cmd);
  await setupAuthToken();
  await subscriptionGrpc.createSubscription(request);
});

export interface UpdateNotificationParams {
  environmentNamespace: string;
  id: string;
  name: String;
  currentSourceTypes: Array<
    Subscription.SourceTypeMap[keyof Subscription.SourceTypeMap]
  >;
  sourceTypes: Array<
    Subscription.SourceTypeMap[keyof Subscription.SourceTypeMap]
  >;
}

export const updateNotification = createAsyncThunk<
  void,
  UpdateNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new UpdateSubscriptionRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);

  if (params.name) {
    const cmd = new RenameSubscriptionCommand();
    cmd.setName(params.name.toString());
    request.setRenameSubscriptionCommand(cmd);
  }

  if (params.sourceTypes) {
    const addList = params.sourceTypes.filter(
      (type) => !params.currentSourceTypes.includes(type)
    );
    if (addList.length > 0) {
      const cmd = new AddSourceTypesCommand();
      cmd.setSourceTypesList(addList);
      request.setAddSourceTypesCommand(cmd);
    }
    const deleteList = params.currentSourceTypes.filter(
      (type) => !params.sourceTypes.includes(type)
    );
    if (deleteList.length > 0) {
      const cmd = new DeleteSourceTypesCommand();
      cmd.setSourceTypesList(deleteList);
      request.setDeleteSourceTypesCommand(cmd);
    }
  }

  await setupAuthToken();
  await subscriptionGrpc.updateSubscription(request);
});

export interface EnableNotificationParams {
  environmentNamespace: string;
  id: string;
}

export const enableNotification = createAsyncThunk<
  void,
  EnableNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  const request = new EnableSubscriptionRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new EnableSubscriptionCommand());
  await setupAuthToken();
  await subscriptionGrpc.enableSubscription(request);
});

export interface DisableNotificationParams {
  environmentNamespace: string;
  id: string;
}

export const disableNotification = createAsyncThunk<
  void,
  DisableNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  const request = new DisableSubscriptionRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new DisableSubscriptionCommand());
  await setupAuthToken();
  await subscriptionGrpc.disableSubscription(request);
});

export interface DeleteNotificationParams {
  environmentNamespace: string;
  id: string;
}

export const deleteNotification = createAsyncThunk<
  void,
  DeleteNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeleteSubscriptionRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new DeleteSubscriptionCommand());
  await setupAuthToken();
  await subscriptionGrpc.deleteSubscription(request);
});

export const notificationSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listNotification.pending, (state) => {
        state.loading = true;
      })
      .addCase(listNotification.fulfilled, (state, action) => {
        notificationAdapter.removeAll(state);
        notificationAdapter.upsertMany(state, action.payload.subscriptionsList);
        state.loading = false;
        state.totalCount = action.payload.totalCount;
      })
      .addCase(listNotification.rejected, (state) => {
        state.loading = false;
      })
      .addCase(createNotification.pending, (state) => {})
      .addCase(createNotification.fulfilled, (state, action) => {})
      .addCase(createNotification.rejected, (state, action) => {})
      .addCase(enableNotification.pending, (state) => {})
      .addCase(enableNotification.fulfilled, (state, action) => {})
      .addCase(enableNotification.rejected, (state, action) => {})
      .addCase(disableNotification.pending, (state) => {})
      .addCase(disableNotification.fulfilled, (state, action) => {})
      .addCase(disableNotification.rejected, (state, action) => {})
      .addCase(deleteNotification.pending, (state) => {})
      .addCase(deleteNotification.fulfilled, (state, action) => {})
      .addCase(deleteNotification.rejected, (state, action) => {});
  },
});
