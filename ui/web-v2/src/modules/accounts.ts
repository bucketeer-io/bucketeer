import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError
} from '@reduxjs/toolkit';
import {
  BoolValue,
  StringValue
} from 'google-protobuf/google/protobuf/wrappers_pb';

import * as accountGrpc from '../grpc/account';
import { AccountV2 } from '../proto/account/account_pb';
import {
  CreateAccountV2Command,
  ChangeAccountV2NameCommand,
  ChangeAccountV2EnvironmentRolesCommand,
  ChangeAccountV2OrganizationRoleCommand
} from '../proto/account/command_pb';
import {
  ListAccountsV2Request,
  ListAccountsV2Response,
  GetAccountV2Request,
  CreateAccountV2Request,
  UpdateAccountV2Request
} from '../proto/account/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'accounts';

export const accountsAdapter = createEntityAdapter<AccountV2.AsObject>({
  selectId: (segment) => segment.email
});

export const { selectAll, selectById } = accountsAdapter.getSelectors();

const initialState = accountsAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getAccountError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getAccountError: null
});

export type OrderBy =
  ListAccountsV2Request.OrderByMap[keyof ListAccountsV2Request.OrderByMap];
export type OrderDirection =
  ListAccountsV2Request.OrderDirectionMap[keyof ListAccountsV2Request.OrderDirectionMap];

interface ListAccountsParams {
  environmentId: string;
  organizationId: string;
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  role?: number;
  disabled?: boolean;
}

export const listAccounts = createAsyncThunk<
  ListAccountsV2Response.AsObject,
  ListAccountsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListAccountsV2Request();
  request.setOrganizationId(params.organizationId);
  request.setEnvironmentId(new StringValue().setValue(params.environmentId));
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  params.disabled != null &&
    request.setDisabled(new BoolValue().setValue(params.disabled));
  await setupAuthToken();
  const result = await accountGrpc.listAccounts(request);
  return result.response.toObject();
});

export interface GetAccountParams {
  organizationId: string;
  email: string;
}

export const getAccount = createAsyncThunk<
  AccountV2.AsObject,
  GetAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetAccountV2Request();
  request.setOrganizationId(params.organizationId);
  request.setEmail(params.email);
  await setupAuthToken();
  const result = await accountGrpc.getAccount(request);
  return result.response.toObject().account;
});

export interface EnableAccountParams {
  organizationId: string;
  environmentId: string;
  email: string;
}

export const enableAccount = createAsyncThunk<
  void,
  EnableAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  // TODO After migration to console3.0, we should use EnableAccountV2Command
  const request = new UpdateAccountV2Request();
  const command = new ChangeAccountV2EnvironmentRolesCommand();
  const environmentRole = new AccountV2.EnvironmentRole();
  environmentRole.setEnvironmentId(params.environmentId);
  environmentRole.setRole(AccountV2.Role.Environment.ENVIRONMENT_VIEWER);
  command.setRolesList([environmentRole]);
  command.setWriteType(
    ChangeAccountV2EnvironmentRolesCommand.WriteType.WRITETYPE_PATCH
  );
  request.setChangeEnvironmentRolesCommand(command);
  request.setOrganizationId(params.organizationId);
  request.setEmail(params.email);
  await setupAuthToken();
  await accountGrpc.updateAccount(request);
});

export interface DisableAccountParams {
  organizationId: string;
  environmentId: string;
  email: string;
}

export const disableAccount = createAsyncThunk<
  void,
  DisableAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  // TODO After migrating to the console 3.0, we should use DisableAccountV2Command
  const request = new UpdateAccountV2Request();
  const cmd = new ChangeAccountV2EnvironmentRolesCommand();
  const environmentRole = new AccountV2.EnvironmentRole();
  environmentRole.setEnvironmentId(params.environmentId);
  environmentRole.setRole(AccountV2.Role.Environment.ENVIRONMENT_UNASSIGNED);
  cmd.setRolesList([environmentRole]);
  cmd.setWriteType(
    ChangeAccountV2EnvironmentRolesCommand.WriteType.WRITETYPE_PATCH
  );
  request.setChangeEnvironmentRolesCommand(cmd);
  request.setOrganizationId(params.organizationId);
  request.setEmail(params.email);
  await setupAuthToken();
  await accountGrpc.updateAccount(request);
});

export interface CreateAccountParams {
  organizationId: string;
  name: string;
  email: string;
  organizationRole: AccountV2.Role.OrganizationMap[keyof AccountV2.Role.OrganizationMap];
  environmentRole: AccountV2.Role.EnvironmentMap[keyof AccountV2.Role.EnvironmentMap];
  environmentId: string;
}

export const createAccount = createAsyncThunk<
  void,
  CreateAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/add`, async (params) => {
  const request = new CreateAccountV2Request();
  const cmd = new CreateAccountV2Command();
  const environmentRole = new AccountV2.EnvironmentRole();
  environmentRole.setEnvironmentId(params.environmentId);
  environmentRole.setRole(params.environmentRole);
  cmd.setEnvironmentRolesList([environmentRole]);
  cmd.setName(params.name);
  cmd.setEmail(params.email);
  cmd.setOrganizationRole(params.organizationRole);
  request.setCommand(cmd);
  request.setOrganizationId(params.organizationId);
  await setupAuthToken();
  await accountGrpc.createAccount(request);
});

export interface UpdateAccountParams {
  organizationId: string;
  name: string;
  email: string;
  environmentId: string;
  environmentRole: AccountV2.Role.EnvironmentMap[keyof AccountV2.Role.EnvironmentMap];
  organizationRole: AccountV2.Role.OrganizationMap[keyof AccountV2.Role.OrganizationMap];
}

export const updateAccount = createAsyncThunk<
  void,
  UpdateAccountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new UpdateAccountV2Request();
  const changeEnvRoleCmd = new ChangeAccountV2EnvironmentRolesCommand();
  const changeOrgRoleCmd = new ChangeAccountV2OrganizationRoleCommand();
  const environmentRole = new AccountV2.EnvironmentRole();
  environmentRole.setEnvironmentId(params.environmentId);
  environmentRole.setRole(params.environmentRole);
  changeEnvRoleCmd.setRolesList([environmentRole]);
  changeEnvRoleCmd.setWriteType(
    ChangeAccountV2EnvironmentRolesCommand.WriteType.WRITETYPE_PATCH
  );
  changeOrgRoleCmd.setRole(params.organizationRole);
  request.setChangeEnvironmentRolesCommand(changeEnvRoleCmd);
  request.setChangeOrganizationRoleCommand(changeOrgRoleCmd);
  if (params.name) {
    const cmd = new ChangeAccountV2NameCommand();
    cmd.setName(params.name);
    request.setChangeNameCommand(cmd);
  }
  request.setEmail(params.email);
  request.setOrganizationId(params.organizationId);
  await setupAuthToken();
  await accountGrpc.updateAccount(request);
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
      .addCase(enableAccount.pending, () => {})
      .addCase(enableAccount.fulfilled, () => {})
      .addCase(enableAccount.rejected, () => {})
      .addCase(disableAccount.pending, () => {})
      .addCase(disableAccount.fulfilled, () => {})
      .addCase(disableAccount.rejected, () => {})
      .addCase(createAccount.pending, () => {})
      .addCase(createAccount.fulfilled, () => {})
      .addCase(createAccount.rejected, () => {})
      .addCase(updateAccount.pending, () => {})
      .addCase(updateAccount.fulfilled, () => {})
      .addCase(updateAccount.rejected, () => {});
  }
});
