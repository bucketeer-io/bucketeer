import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk
} from '@reduxjs/toolkit';

import * as autoOpsGrpc from '../grpc/autoops';
import { AutoOpsRule } from '../proto/autoops/auto_ops_rule_pb';
import {
  AddDatetimeClauseCommand,
  AddOpsEventRateClauseCommand,
  ChangeDatetimeClauseCommand,
  ChangeOpsEventRateClauseCommand,
  CreateAutoOpsRuleCommand,
  DeleteAutoOpsRuleCommand,
  DeleteClauseCommand,
  StopAutoOpsRuleCommand
} from '../proto/autoops/command_pb';
import {
  CreateAutoOpsRuleRequest,
  DeleteAutoOpsRuleRequest,
  ListAutoOpsRulesRequest,
  ListAutoOpsRulesResponse,
  StopAutoOpsRuleRequest,
  UpdateAutoOpsRuleRequest
} from '../proto/autoops/service_pb';

import { AppState } from '.';

const MODULE_NAME = 'autoOpsRules';

export const autoOpsRulesAdapter = createEntityAdapter({
  selectId: (autoOpsRule: AutoOpsRule.AsObject) => autoOpsRule.id
});

export const { selectAll, selectById } = autoOpsRulesAdapter.getSelectors();

interface CreateAutoOpsRuleParams {
  environmentId: string;
  command: CreateAutoOpsRuleCommand;
}

export const createAutoOpsRule = createAsyncThunk<
  void,
  CreateAutoOpsRuleParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateAutoOpsRuleRequest();
  request.setEnvironmentId(params.environmentId);
  request.setCommand(params.command);
  await autoOpsGrpc.createAutoOpsRule(request);
});

export interface ListAutoOpsRulesParams {
  environmentId: string;
  featureId: string;
}

export const listAutoOpsRules = createAsyncThunk<
  ListAutoOpsRulesResponse.AsObject,
  ListAutoOpsRulesParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListAutoOpsRulesRequest();
  request.setFeatureIdsList([params.featureId]);
  request.setEnvironmentId(params.environmentId);
  const result = await autoOpsGrpc.listAutoOpsRules(request);
  return result.response.toObject();
});

export interface UpdateAutoOpsRuleParams {
  environmentId: string;
  id: string;
  addOpsEventRateClauseCommands?: Array<AddOpsEventRateClauseCommand>;
  changeOpsEventRateClauseCommands?: Array<ChangeOpsEventRateClauseCommand>;
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
  request.setEnvironmentId(params.environmentId);
  request.setId(params.id);
  params.addOpsEventRateClauseCommands?.length > 0 &&
    request.setAddOpsEventRateClauseCommandsList(
      params.addOpsEventRateClauseCommands
    );
  params.changeOpsEventRateClauseCommands?.length > 0 &&
    request.setChangeOpsEventRateClauseCommandsList(
      params.changeOpsEventRateClauseCommands
    );
  params.addDatetimeClauseCommands?.length > 0 &&
    request.setAddDatetimeClauseCommandsList(params.addDatetimeClauseCommands);
  params.changeDatetimeClauseCommands?.length > 0 &&
    request.setChangeDatetimeClauseCommandsList(
      params.changeDatetimeClauseCommands
    );
  params.deleteClauseCommands?.length > 0 &&
    request.setDeleteClauseCommandsList(params.deleteClauseCommands);
  await autoOpsGrpc.updateAutoOpsRule(request);
});

export interface DeleteAutoOpsRuleParams {
  environmentId: string;
  id: string;
}

export const deleteAutoOpsRule = createAsyncThunk<
  void,
  DeleteAutoOpsRuleParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeleteAutoOpsRuleRequest();
  request.setId(params.id);
  request.setEnvironmentId(params.environmentId);
  const command = new DeleteAutoOpsRuleCommand();
  request.setCommand(command);
  await autoOpsGrpc.deleteAutoOpsRule(request);
});

export interface StopAutoOpsRuleParams {
  environmentId: string;
  id: string;
}

export const stopAutoOpsRule = createAsyncThunk<
  void,
  StopAutoOpsRuleParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/stop`, async (params) => {
  const request = new StopAutoOpsRuleRequest();
  request.setId(params.id);
  request.setEnvironmentId(params.environmentId);
  const command = new StopAutoOpsRuleCommand();
  request.setCommand(command);
  await autoOpsGrpc.stopAutoOpsRule(request);
});

const initialState = autoOpsRulesAdapter.getInitialState<{
  loading: boolean;
}>({
  loading: false
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
      .addCase(deleteAutoOpsRule.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(deleteAutoOpsRule.rejected, (state) => {
        state.loading = false;
      })
      .addCase(updateAutoOpsRule.pending, (state) => {
        state.loading = true;
      })
      .addCase(updateAutoOpsRule.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(updateAutoOpsRule.rejected, (state) => {
        state.loading = false;
      })
      .addCase(stopAutoOpsRule.pending, (state) => {
        state.loading = true;
      })
      .addCase(stopAutoOpsRule.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(stopAutoOpsRule.rejected, (state) => {
        state.loading = false;
      });
  }
});
