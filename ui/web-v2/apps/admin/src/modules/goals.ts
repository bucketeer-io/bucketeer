import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';
import { BoolValue } from 'google-protobuf/google/protobuf/wrappers_pb';

import * as grpc from '../grpc/experiment';
import {
  ArchiveGoalCommand,
  ChangeDescriptionGoalCommand,
  CreateGoalCommand,
  DeleteGoalCommand,
  RenameGoalCommand,
} from '../proto/experiment/command_pb';
import { Goal } from '../proto/experiment/goal_pb';
import {
  ArchiveGoalRequest,
  CreateGoalRequest,
  DeleteGoalRequest,
  GetGoalRequest,
  GetGoalResponse,
  ListGoalsRequest,
  ListGoalsResponse,
  UpdateGoalRequest,
} from '../proto/experiment/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'goals';

export const goalsAdapter = createEntityAdapter<Goal.AsObject>({
  selectId: (goal) => goal.id,
});

export const { selectAll, selectById } = goalsAdapter.getSelectors();

interface GetGoalParams {
  environmentNamespace: string;
  id: string;
}

export const getGoal = createAsyncThunk<
  GetGoalResponse.AsObject,
  GetGoalParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/getGoal`, async (params) => {
  const request = new GetGoalRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  await setupAuthToken();
  const result = await grpc.getGoal(request);
  return result.response.toObject();
});

export type OrderBy =
  ListGoalsRequest.OrderByMap[keyof ListGoalsRequest.OrderByMap];
export type OrderDirection =
  ListGoalsRequest.OrderDirectionMap[keyof ListGoalsRequest.OrderDirectionMap];

interface ListGoalsParams {
  environmentNamespace: string;
  pageSize: number;
  cursor: string;
  searchKeyword: string;
  status?: boolean;
  archived?: boolean;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

export const listGoals = createAsyncThunk<
  ListGoalsResponse.AsObject,
  ListGoalsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/listGoals`, async (params) => {
  const request = new ListGoalsRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  if (params.status != null) {
    const boolValue = new BoolValue();
    boolValue.setValue(params.status);
    request.setIsInUseStatus(boolValue);
  }
  const boolValue = new BoolValue();
  boolValue.setValue(params.archived);
  request.setArchived(boolValue);
  await setupAuthToken();
  const result = await grpc.listGoals(request);
  return result.response.toObject();
});

interface CreateGoalParams {
  environmentNamespace: string;
  id: string;
  name: string;
  description: string;
}

export const createGoal = createAsyncThunk<
  void,
  CreateGoalParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/createGoal`, async (params) => {
  const request = new CreateGoalRequest();
  const command = new CreateGoalCommand();
  request.setEnvironmentNamespace(params.environmentNamespace);
  command.setId(params.id);
  command.setName(params.name);
  command.setDescription(params.description);
  request.setCommand(command);
  await setupAuthToken();
  await grpc.createGoal(request);
});

interface DeleteGoalParams {
  environmentNamespace: string;
  id: string;
}

export const deleteGoal = createAsyncThunk<
  void,
  DeleteGoalParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/deleteGoal`, async (params) => {
  const request = new DeleteGoalRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new DeleteGoalCommand());
  await setupAuthToken();
  await grpc.deleteGoal(request);
});

interface UpdateGoalParams {
  environmentNamespace: string;
  id: string;
  name?: string;
  description?: string;
}

export const updateGoal = createAsyncThunk<
  void,
  UpdateGoalParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/updateGoal`, async (params) => {
  const request = new UpdateGoalRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  if (params.name) {
    const renameCommand = new RenameGoalCommand();
    renameCommand.setName(params.name);
    request.setRenameCommand(renameCommand);
  }
  if (params.description) {
    const changeDescCommand = new ChangeDescriptionGoalCommand();
    changeDescCommand.setDescription(params.description);
    request.setChangeDescriptionCommand(changeDescCommand);
  }
  await setupAuthToken();
  await grpc.updateGoal(request);
});

export interface ArchiveGoalParams {
  environmentNamespace: string;
  id: string;
}

export const archiveGoal = createAsyncThunk<
  void,
  ArchiveGoalParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/archive`, async (params) => {
  const request = new ArchiveGoalRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new ArchiveGoalCommand());
  await setupAuthToken();
  await grpc.archiveGoal(request);
});

const initialState = goalsAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getGoalError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getGoalError: null,
});

export type GoalsState = typeof initialState;

export const goalsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listGoals.pending, (state) => {
        state.loading = true;
      })
      .addCase(listGoals.fulfilled, (state, action) => {
        goalsAdapter.removeAll(state);
        goalsAdapter.upsertMany(state, action.payload.goalsList);
        state.totalCount = action.payload.totalCount;
        state.loading = false;
      })
      .addCase(listGoals.rejected, (state) => {
        state.loading = false;
      })
      .addCase(getGoal.pending, (state) => {
        state.getGoalError = null;
      })
      .addCase(getGoal.fulfilled, (state, action) => {
        state.getGoalError = null;
        if (action.payload) {
          goalsAdapter.upsertOne(state, action.payload.goal);
        }
      })
      .addCase(getGoal.rejected, (state, action) => {
        state.getGoalError = action.error;
      })
      .addCase(createGoal.pending, (state) => {})
      .addCase(createGoal.fulfilled, (state, action) => {})
      .addCase(createGoal.rejected, (state, action) => {})
      .addCase(updateGoal.pending, (state) => {})
      .addCase(updateGoal.fulfilled, (state, action) => {})
      .addCase(updateGoal.rejected, (state, action) => {});
  },
});
