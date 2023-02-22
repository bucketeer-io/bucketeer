import { ChangeWebhookClauseCommand } from './../proto/autoops/command_pb.d';
import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from '@reduxjs/toolkit';

import * as autoOpsGrpc from '../grpc/autoops';
import { AutoOpsRule } from '../proto/autoops/auto_ops_rule_pb';
import {
  AddDatetimeClauseCommand,
  AddOpsEventRateClauseCommand,
  AddWebhookClauseCommand,
  ChangeAutoOpsRuleOpsTypeCommand,
  ChangeDatetimeClauseCommand,
  ChangeOpsEventRateClauseCommand,
  CreateAutoOpsRuleCommand,
  DeleteAutoOpsRuleCommand,
  DeleteClauseCommand,
} from '../proto/autoops/command_pb';
import {
  CreateAutoOpsRuleRequest,
  DeleteAutoOpsRuleRequest,
  ListAutoOpsRulesRequest,
  ListAutoOpsRulesResponse,
  UpdateAutoOpsRuleRequest,
} from '../proto/autoops/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'autoOpsRules';

export const autoOpsRulesAdapter = createEntityAdapter<AutoOpsRule.AsObject>({
  selectId: (autoOpsRule) => autoOpsRule.id,
});

export const { selectAll, selectById } = autoOpsRulesAdapter.getSelectors();

interface CreateAutoOpsRuleParams {
  environmentNamespace: string;
  command: CreateAutoOpsRuleCommand;
}

export const createAutoOpsRule = createAsyncThunk<
  void,
  CreateAutoOpsRuleParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateAutoOpsRuleRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setCommand(params.command);
  await setupAuthToken();
  await autoOpsGrpc.createAutoOpsRule(request);
});

export interface ListAutoOpsRulesParams {
  environmentNamespace: string;
  featureId: string;
}

export const listAutoOpsRules = createAsyncThunk<
  ListAutoOpsRulesResponse.AsObject,
  ListAutoOpsRulesParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListAutoOpsRulesRequest();
  request.setFeatureIdsList([params.featureId]);
  request.setEnvironmentNamespace(params.environmentNamespace);
  await setupAuthToken();
  const result = await autoOpsGrpc.listAutoOpsRules(request);
  return result.response.toObject();
});

export interface UpdateAutoOpsRuleParams {
  environmentNamespace: string;
  id: string;
  changeAutoOpsRuleOpsTypeCommand?: ChangeAutoOpsRuleOpsTypeCommand;
  addOpsEventRateClauseCommands?: Array<AddOpsEventRateClauseCommand>;
  addWebhookClauseCommands?: Array<AddWebhookClauseCommand>;
  changeOpsEventRateClauseCommands?: Array<ChangeOpsEventRateClauseCommand>;
  changeWebhookClauseCommands?: Array<ChangeWebhookClauseCommand>;
  addDatetimeClauseCommands?: Array<AddDatetimeClauseCommand>;
  changeDatetimeClauseCommands?: Array<ChangeDatetimeClauseCommand>;
  deleteClauseCommands?: Array<DeleteClauseCommand>;
}

export const updateAutoOpsRule = createAsyncThunk<
  void,
  UpdateAutoOpsRuleParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new UpdateAutoOpsRuleRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  params.changeAutoOpsRuleOpsTypeCommand &&
    request.setChangeAutoOpsRuleOpsTypeCommand(
      params.changeAutoOpsRuleOpsTypeCommand
    );
  params.addOpsEventRateClauseCommands?.length > 0 &&
    request.setAddOpsEventRateClauseCommandsList(
      params.addOpsEventRateClauseCommands
    );
  params.addWebhookClauseCommands?.length > 0 &&
    request.setAddWebhookClauseCommandsList(params.addWebhookClauseCommands);
  params.changeOpsEventRateClauseCommands?.length > 0 &&
    request.setChangeOpsEventRateClauseCommandsList(
      params.changeOpsEventRateClauseCommands
    );
  params.changeWebhookClauseCommands?.length > 0 &&
    request.setChangeWebhookClauseCommandsList(
      params.changeWebhookClauseCommands
    );
  params.addDatetimeClauseCommands?.length > 0 &&
    request.setAddDatetimeClauseCommandsList(params.addDatetimeClauseCommands);
  params.changeDatetimeClauseCommands?.length > 0 &&
    request.setChangeDatetimeClauseCommandsList(
      params.changeDatetimeClauseCommands
    );
  params.deleteClauseCommands?.length > 0 &&
    request.setDeleteClauseCommandsList(params.deleteClauseCommands);
  await setupAuthToken();
  await autoOpsGrpc.updateAutoOpsRule(request);
});

export interface DeleteAutoOpsRuleParams {
  environmentNamespace: string;
  id: string;
}

export const deleteAutoOpsRule = createAsyncThunk<
  void,
  DeleteAutoOpsRuleParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeleteAutoOpsRuleRequest();
  request.setId(params.id);
  request.setEnvironmentNamespace(params.environmentNamespace);
  const command = new DeleteAutoOpsRuleCommand();
  request.setCommand(command);
  await setupAuthToken();
  await autoOpsGrpc.deleteAutoOpsRule(request);
});

const initialState = autoOpsRulesAdapter.getInitialState<{
  loading: boolean;
}>({
  loading: false,
});

export type AutoOpsRulesState = typeof initialState;

export const autoOpsRulesSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listAutoOpsRules.pending, (state) => {
        state.loading = true;
      })
      .addCase(listAutoOpsRules.fulfilled, (state, action) => {
        autoOpsRulesAdapter.removeAll(state);
        autoOpsRulesAdapter.upsertMany(state, action.payload.autoOpsRulesList);
        state.loading = false;
      })
      .addCase(listAutoOpsRules.rejected, (state) => {
        state.loading = false;
      })
      .addCase(deleteAutoOpsRule.pending, (state) => {
        state.loading = true;
      })
      .addCase(deleteAutoOpsRule.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(deleteAutoOpsRule.rejected, (state) => {
        state.loading = false;
      })
      .addCase(updateAutoOpsRule.pending, (state) => {
        state.loading = true;
      })
      .addCase(updateAutoOpsRule.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(updateAutoOpsRule.rejected, (state) => {
        state.loading = false;
      });
  },
});
