import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from '@reduxjs/toolkit';

import * as progressiveRolloutGrpc from '../grpc/progressiveRollout';
import { AutoOpsRule } from '../proto/autoops/auto_ops_rule_pb';
import { ProgressiveRollout } from '../proto/autoops/progressive_rollout_pb';
import {
  AddDatetimeClauseCommand,
  AddOpsEventRateClauseCommand,
  ChangeAutoOpsRuleOpsTypeCommand,
  ChangeDatetimeClauseCommand,
  ChangeOpsEventRateClauseCommand,
  CreateAutoOpsRuleCommand,
  DeleteAutoOpsRuleCommand,
  DeleteClauseCommand,
  CreateProgressiveRolloutCommand,
  DeleteProgressiveRolloutCommand,
} from '../proto/autoops/command_pb';
import {
  CreateAutoOpsRuleRequest,
  CreateProgressiveRolloutRequest,
  DeleteAutoOpsRuleRequest,
  ListAutoOpsRulesRequest,
  ListAutoOpsRulesResponse,
  UpdateAutoOpsRuleRequest,
  ListProgressiveRolloutsRequest,
  ListProgressiveRolloutsResponse,
  DeleteProgressiveRolloutRequest,
  DeleteProgressiveRolloutResponse,
} from '../proto/autoops/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'progressiveRollout';

export const progressiveRolloutAdapter =
  createEntityAdapter<ProgressiveRollout.AsObject>({
    selectId: (progressiveRollout) => progressiveRollout.id,
  });

export const { selectAll, selectById } =
  progressiveRolloutAdapter.getSelectors();

interface CreateProgressiveRolloutParams {
  environmentNamespace: string;
  command: CreateProgressiveRolloutCommand;
}

export const createProgressiveRollout = createAsyncThunk<
  void,
  CreateProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateProgressiveRolloutRequest();

  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setCommand(params.command);
  await setupAuthToken();
  await progressiveRolloutGrpc.createProgressiveRollout(request);
});

export interface ListProgressiveRolloutParams {
  environmentNamespace: string;
  featureId: string;
}

export const listProgressiveRollout = createAsyncThunk<
  ListProgressiveRolloutsResponse.AsObject,
  ListProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListProgressiveRolloutsRequest();
  request.setFeatureIdsList([params.featureId]);
  request.setEnvironmentNamespace(params.environmentNamespace);
  await setupAuthToken();
  const result = await progressiveRolloutGrpc.listProgressiveRollouts(request);
  return result.response.toObject();
});

// export interface UpdateAutoOpsRuleParams {
//   environmentNamespace: string;
//   id: string;
//   changeAutoOpsRuleOpsTypeCommand?: ChangeAutoOpsRuleOpsTypeCommand;
//   addOpsEventRateClauseCommands?: Array<AddOpsEventRateClauseCommand>;
//   changeOpsEventRateClauseCommands?: Array<ChangeOpsEventRateClauseCommand>;
//   addDatetimeClauseCommands?: Array<AddDatetimeClauseCommand>;
//   changeDatetimeClauseCommands?: Array<ChangeDatetimeClauseCommand>;
//   deleteClauseCommands?: Array<DeleteClauseCommand>;
// }

// export const updateAutoOpsRule = createAsyncThunk<
//   void,
//   UpdateAutoOpsRuleParams | undefined,
//   { state: AppState }
// >(`${MODULE_NAME}/update`, async (params) => {
//   const request = new UpdateAutoOpsRuleRequest();
//   request.setEnvironmentNamespace(params.environmentNamespace);
//   request.setId(params.id);
//   params.changeAutoOpsRuleOpsTypeCommand &&
//     request.setChangeAutoOpsRuleOpsTypeCommand(
//       params.changeAutoOpsRuleOpsTypeCommand
//     );
//   params.addOpsEventRateClauseCommands?.length > 0 &&
//     request.setAddOpsEventRateClauseCommandsList(
//       params.addOpsEventRateClauseCommands
//     );
//   params.changeOpsEventRateClauseCommands?.length > 0 &&
//     request.setChangeOpsEventRateClauseCommandsList(
//       params.changeOpsEventRateClauseCommands
//     );
//   params.addDatetimeClauseCommands?.length > 0 &&
//     request.setAddDatetimeClauseCommandsList(params.addDatetimeClauseCommands);
//   params.changeDatetimeClauseCommands?.length > 0 &&
//     request.setChangeDatetimeClauseCommandsList(
//       params.changeDatetimeClauseCommands
//     );
//   params.deleteClauseCommands?.length > 0 &&
//     request.setDeleteClauseCommandsList(params.deleteClauseCommands);
//   await setupAuthToken();
//   await progressiveRolloutGrpc.updateAutoOpsRule(request);
// });

export interface DeleteProgressiveRolloutParams {
  environmentNamespace: string;
  id: string;
}

export const deleteProgressiveRollout = createAsyncThunk<
  void,
  DeleteProgressiveRolloutParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeleteProgressiveRolloutRequest();
  request.setId(params.id);
  request.setEnvironmentNamespace(params.environmentNamespace);
  const command = new DeleteProgressiveRolloutCommand();
  request.setCommand(command);
  await setupAuthToken();
  await progressiveRolloutGrpc.deleteProgressiveRollout(request);
});

const initialState = progressiveRolloutAdapter.getInitialState<{
  loading: boolean;
}>({
  loading: false,
});

export type ProgressiveRolloutState = typeof initialState;

export const progressiveRolloutSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listProgressiveRollout.pending, (state) => {
        state.loading = true;
      })
      .addCase(listProgressiveRollout.fulfilled, (state, action) => {
        progressiveRolloutAdapter.removeAll(state);
        progressiveRolloutAdapter.upsertMany(
          state,
          action.payload.progressiveRolloutsList
        );
        state.loading = false;
      })
      .addCase(listProgressiveRollout.rejected, (state) => {
        state.loading = false;
      });
    //   .addCase(deleteAutoOpsRule.pending, (state) => {
    //     state.loading = true;
    //   })
    //   .addCase(deleteAutoOpsRule.fulfilled, (state, action) => {
    //     state.loading = false;
    //   })
    //   .addCase(deleteAutoOpsRule.rejected, (state) => {
    //     state.loading = false;
    //   })
    //   .addCase(updateAutoOpsRule.pending, (state) => {
    //     state.loading = true;
    //   })
    //   .addCase(updateAutoOpsRule.fulfilled, (state, action) => {
    //     state.loading = false;
    //   })
    //   .addCase(updateAutoOpsRule.rejected, (state) => {
    //     state.loading = false;
    //   });
  },
});
