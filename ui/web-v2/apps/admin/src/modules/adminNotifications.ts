import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';
import { BoolValue } from 'google-protobuf/google/protobuf/wrappers_pb';

import * as grpc from '../grpc/adminSubscription';
import {
  AddSourceTypesCommand,
  CreateAdminSubscriptionCommand,
  DeleteSourceTypesCommand,
  DeleteAdminSubscriptionCommand,
  DisableAdminSubscriptionCommand,
  EnableAdminSubscriptionCommand,
  RenameAdminSubscriptionCommand,
} from '../proto/notification/command_pb';
import {
  Recipient,
  SlackChannelRecipient,
} from '../proto/notification/recipient_pb';
import {
  CreateAdminSubscriptionRequest,
  DeleteAdminSubscriptionRequest,
  DisableAdminSubscriptionRequest,
  EnableAdminSubscriptionRequest,
  GetAdminSubscriptionRequest,
  ListAdminSubscriptionsRequest,
  ListAdminSubscriptionsResponse,
  UpdateAdminSubscriptionRequest,
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
  getSubscriptionError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getSubscriptionError: null,
});

export type OrderBy =
  ListAdminSubscriptionsRequest.OrderByMap[keyof ListAdminSubscriptionsRequest.OrderByMap];
export type OrderDirection =
  ListAdminSubscriptionsRequest.OrderDirectionMap[keyof ListAdminSubscriptionsRequest.OrderDirectionMap];

interface ListAdminNotificationParams {
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  disabled?: boolean;
}

export const listNotification = createAsyncThunk<
  ListAdminSubscriptionsResponse.AsObject,
  ListAdminNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListAdminSubscriptionsRequest();
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  params.disabled != null &&
    request.setDisabled(new BoolValue().setValue(params.disabled));
  await setupAuthToken();
  const result = await grpc.listSubscriptions(request);
  return result.response.toObject();
});

export interface GetAdminNotificationParams {
  id: string;
}

export const getNotification = createAsyncThunk<
  Subscription.AsObject,
  GetAdminNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetAdminSubscriptionRequest();
  request.setId(params.id);
  await setupAuthToken();
  const result = await grpc.getSubscription(request);
  return result.response.toObject().subscription;
});

export interface CreateNotificationParams {
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
  const cmd = new CreateAdminSubscriptionCommand();
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

  const request = new CreateAdminSubscriptionRequest();
  request.setCommand(cmd);
  await setupAuthToken();
  await grpc.createSubscription(request);
});

export interface UpdateAdminNotificationParams {
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
  UpdateAdminNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new UpdateAdminSubscriptionRequest();
  request.setId(params.id);

  if (params.name) {
    const cmd = new RenameAdminSubscriptionCommand();
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
  await grpc.updateSubscription(request);
});

export interface EnableAdminNotificationParams {
  id: string;
}

export const enableNotification = createAsyncThunk<
  void,
  EnableAdminNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  const request = new EnableAdminSubscriptionRequest();
  request.setId(params.id);
  request.setCommand(new EnableAdminSubscriptionCommand());
  await setupAuthToken();
  await grpc.enableSubscription(request);
});

export interface DisableNotificationParams {
  id: string;
}

export const disableNotification = createAsyncThunk<
  void,
  DisableNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  const request = new DisableAdminSubscriptionRequest();
  request.setId(params.id);
  request.setCommand(new DisableAdminSubscriptionCommand());
  await setupAuthToken();
  await grpc.disableSubscription(request);
});

export interface DeleteAdminNotificationParams {
  id: string;
}

export const deleteNotification = createAsyncThunk<
  void,
  DeleteAdminNotificationParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeleteAdminSubscriptionRequest();
  request.setId(params.id);
  request.setCommand(new DeleteAdminSubscriptionCommand());
  await setupAuthToken();
  await grpc.deleteSubscription(request);
});

export const adminNotificationSlice = createSlice({
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
      .addCase(getNotification.pending, (state) => {
        state.getSubscriptionError = null;
      })
      .addCase(getNotification.fulfilled, (state, action) => {
        state.getSubscriptionError = null;
        if (action.payload) {
          notificationAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(getNotification.rejected, (state, action) => {
        state.getSubscriptionError = action.error;
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
