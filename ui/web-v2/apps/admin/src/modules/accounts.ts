import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';
import {
  BoolValue,
  Int32Value,
} from 'google-protobuf/google/protobuf/wrappers_pb';

import * as accountGrpc from '../grpc/account';
import { Account } from '../proto/account/account_pb';
import {
  EnableAccountCommand,
  DisableAccountCommand,
  CreateAccountCommand,
  ChangeAccountRoleCommand,
} from '../proto/account/command_pb';
import {
  ListAccountsRequest,
  ListAccountsResponse,
  GetAccountRequest,
  EnableAccountRequest,
  DisableAccountRequest,
  CreateAccountRequest,
  ChangeAccountRoleRequest,
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
  ListAccountsRequest.OrderByMap[keyof ListAccountsRequest.OrderByMap];
export type OrderDirection =
  ListAccountsRequest.OrderDirectionMap[keyof ListAccountsRequest.OrderDirectionMap];

interface ListAccountsParams {
  environmentNamespace: string;
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  role?: number;
  disabled?: boolean;
}

export const listAccounts = createAsyncThunk<
  ListAccountsResponse.AsObject,
  ListAccountsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListAccountsRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  params.role != null &&
    request.setRole(new Int32Value().setValue(params.role));
  params.disabled != null &&
    request.setDisabled(new BoolValue().setValue(params.disabled));
  await setupAuthToken();
  const result = await accountGrpc.listAccounts(request);
  return result.response.toObject();
});

export interface GetAccountParams {
  environmentNamespace: string;
  email: string;
}

export const getAccount = createAsyncThunk<
  Account.AsObject,
  GetAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetAccountRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setEmail(params.email);
  await setupAuthToken();
  const result = await accountGrpc.getAccount(request);
  return result.response.toObject().account;
});

export interface EnableAccountParams {
  environmentNamespace: string;
  id: string;
}

export const enableAccount = createAsyncThunk<
  void,
  EnableAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  const request = new EnableAccountRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new EnableAccountCommand());
  await setupAuthToken();
  await accountGrpc.enableAccount(request);
});

export interface DisableAccountParams {
  environmentNamespace: string;
  id: string;
}

export const disableAccount = createAsyncThunk<
  void,
  DisableAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  const request = new DisableAccountRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new DisableAccountCommand());
  await setupAuthToken();
  await accountGrpc.disableAccount(request);
});

export interface CreateAccountParams {
  environmentNamespace: string;
  email: string;
  role?: Account.RoleMap[keyof Account.RoleMap];
}

export const createAccount = createAsyncThunk<
  void,
  CreateAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/add`, async (params) => {
  const request = new CreateAccountRequest();
  const cmd = new CreateAccountCommand();
  cmd.setEmail(params.email);
  cmd.setRole(params.role);
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setCommand(cmd);
  await setupAuthToken();
  await accountGrpc.createAccount(request);
});

export interface UpdateAccountParams {
  environmentNamespace: string;
  id: string;
  role?: Account.RoleMap[keyof Account.RoleMap];
}

export const updateAccount = createAsyncThunk<
  void,
  UpdateAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new ChangeAccountRoleRequest();
  const cmd = new ChangeAccountRoleCommand();
  cmd.setRole(params.role);
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(cmd);
  await setupAuthToken();
  await accountGrpc.changeAccountRole(request);
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
      .addCase(createAccount.rejected, (state, action) => {})
      .addCase(updateAccount.pending, (state) => {})
      .addCase(updateAccount.fulfilled, (state, action) => {})
      .addCase(updateAccount.rejected, (state, action) => {});
  },
});
