import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';
import { BoolValue } from 'google-protobuf/google/protobuf/wrappers_pb';

import * as grpc from '../grpc/adminaccount';
import { Account } from '../proto/account/account_pb';
import {
  EnableAdminAccountCommand,
  DisableAdminAccountCommand,
  CreateAdminAccountCommand,
} from '../proto/account/command_pb';
import {
  ListAdminAccountsRequest,
  ListAdminAccountsResponse,
  GetAdminAccountRequest,
  EnableAdminAccountRequest,
  DisableAdminAccountRequest,
  CreateAdminAccountRequest,
} from '../proto/account/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'accounts';

export const accountsAdapter = createEntityAdapter<Account.AsObject>({
  selectId: (segment) => segment.id,
});

export const { selectAll, selectById } = accountsAdapter.getSelectors();

const initialState = accountsAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getAccountError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getAccountError: null,
});

export type OrderBy =
  ListAdminAccountsRequest.OrderByMap[keyof ListAdminAccountsRequest.OrderByMap];
export type OrderDirection =
  ListAdminAccountsRequest.OrderDirectionMap[keyof ListAdminAccountsRequest.OrderDirectionMap];

interface ListAdminAccountsParams {
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  disabled?: boolean;
}

export const listAccounts = createAsyncThunk<
  ListAdminAccountsResponse.AsObject,
  ListAdminAccountsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListAdminAccountsRequest();
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  params.disabled != null &&
    request.setDisabled(new BoolValue().setValue(params.disabled));
  await setupAuthToken();
  const result = await grpc.listAdminAccounts(request);
  return result.response.toObject();
});

export interface GetAdminAccountParams {
  email: string;
}

export const getAccount = createAsyncThunk<
  Account.AsObject,
  GetAdminAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetAdminAccountRequest();
  request.setEmail(params.email);
  await setupAuthToken();
  const result = await grpc.getAdminAccount(request);
  return result.response.toObject().account;
});

export interface EnableAdminAccountParams {
  id: string;
}

export const enableAccount = createAsyncThunk<
  void,
  EnableAdminAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  const request = new EnableAdminAccountRequest();
  request.setId(params.id);
  request.setCommand(new EnableAdminAccountCommand());
  await setupAuthToken();
  await grpc.enableAdminAccount(request);
});

export interface DisableAdminAccountParams {
  id: string;
}

export const disableAccount = createAsyncThunk<
  void,
  DisableAdminAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  const request = new DisableAdminAccountRequest();
  request.setId(params.id);
  request.setCommand(new DisableAdminAccountCommand());
  await setupAuthToken();
  await grpc.disableAdminAccount(request);
});

export interface CreateAdminAccountParams {
  email: string;
}

export const createAccount = createAsyncThunk<
  void,
  CreateAdminAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/add`, async (params) => {
  const request = new CreateAdminAccountRequest();
  const cmd = new CreateAdminAccountCommand();
  cmd.setEmail(params.email);
  request.setCommand(cmd);
  await setupAuthToken();
  await grpc.createAdminAccount(request);
});

export type AccountsState = typeof initialState;

export const accountsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listAccounts.pending, (state) => {
        state.loading = true;
      })
      .addCase(listAccounts.fulfilled, (state, action) => {
        accountsAdapter.removeAll(state);
        accountsAdapter.upsertMany(state, action.payload.accountsList);
        state.loading = false;
        state.totalCount = action.payload.totalCount;
      })
      .addCase(listAccounts.rejected, (state) => {
        state.loading = false;
      })
      .addCase(getAccount.pending, (state) => {
        state.getAccountError = null;
      })
      .addCase(getAccount.fulfilled, (state, action) => {
        state.getAccountError = null;
        if (action.payload) {
          accountsAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(getAccount.rejected, (state, action) => {
        state.getAccountError = action.error;
      })
      .addCase(enableAccount.pending, (state) => {})
      .addCase(enableAccount.fulfilled, (state, action) => {})
      .addCase(enableAccount.rejected, (state, action) => {})
      .addCase(disableAccount.pending, (state) => {})
      .addCase(disableAccount.fulfilled, (state, action) => {})
      .addCase(disableAccount.rejected, (state, action) => {})
      .addCase(createAccount.pending, (state) => {})
      .addCase(createAccount.fulfilled, (state, action) => {})
      .addCase(createAccount.rejected, (state, action) => {});
  },
});
