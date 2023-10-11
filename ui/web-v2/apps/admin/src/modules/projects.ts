import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';
import { BoolValue } from 'google-protobuf/google/protobuf/wrappers_pb';

import * as grpc from '../grpc/project';
import {
  ChangeDescriptionProjectCommand,
  ConvertTrialProjectCommand,
  CreateProjectCommand,
  EnableProjectCommand,
  DisableProjectCommand,
  RenameProjectCommand,
} from '../proto/environment/command_pb';
import { Project } from '../proto/environment/project_pb';
import {
  ConvertTrialProjectRequest,
  CreateProjectRequest,
  ListProjectsRequest,
  ListProjectsResponse,
  GetProjectRequest,
  EnableProjectRequest,
  DisableProjectRequest,
  UpdateProjectRequest,
} from '../proto/environment/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'projects';

export const projectAdapter = createEntityAdapter<Project.AsObject>({
  selectId: (e) => e.id,
});

export const { selectAll, selectById } = projectAdapter.getSelectors();

const initialState = projectAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getProjectError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getProjectError: null,
});

export type OrderBy =
  ListProjectsRequest.OrderByMap[keyof ListProjectsRequest.OrderByMap];
export type OrderDirection =
  ListProjectsRequest.OrderDirectionMap[keyof ListProjectsRequest.OrderDirectionMap];

interface ListProjectsRequestParams {
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  disabled?: boolean;
}

export const listProjects = createAsyncThunk<
  ListProjectsResponse.AsObject,
  ListProjectsRequestParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListProjectsRequest();
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  params.disabled != null &&
    request.setDisabled(new BoolValue().setValue(params.disabled));
  await setupAuthToken();
  const result = await grpc.listProjects(request);
  return result.response.toObject();
});

export interface GetParams {
  id: string;
}

export const getProject = createAsyncThunk<
  Project.AsObject,
  GetParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetProjectRequest();
  request.setId(params.id);
  await setupAuthToken();
  const result = await grpc.getProject(request);
  return result.response.toObject().project;
});

export interface ConvertProjectParams {
  id: string;
}

export const convertProject = createAsyncThunk<
  void,
  ConvertProjectParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/convert`, async (params) => {
  const request = new ConvertTrialProjectRequest();
  request.setId(params.id);
  request.setCommand(new ConvertTrialProjectCommand());
  await setupAuthToken();
  await grpc.convertTrialProject(request);
});

export interface EnableProjectParams {
  id: string;
}

export const enableProject = createAsyncThunk<
  void,
  EnableProjectParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/enable`, async (params) => {
  const request = new EnableProjectRequest();
  request.setId(params.id);
  request.setCommand(new EnableProjectCommand());
  await setupAuthToken();
  await grpc.enableProject(request);
});

export interface DisableProjectParams {
  id: string;
}

export const disableProject = createAsyncThunk<
  void,
  DisableProjectParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/disable`, async (params) => {
  const request = new DisableProjectRequest();
  request.setId(params.id);
  request.setCommand(new DisableProjectCommand());
  await setupAuthToken();
  await grpc.disableProject(request);
});

export interface CreateProjectParams {
  name: string;
  urlCode: string;
  description: string;
}

export const createProject = createAsyncThunk<
  void,
  CreateProjectParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const request = new CreateProjectRequest();
  const command = new CreateProjectCommand();
  command.setName(params.name);
  command.setUrlCode(params.urlCode);
  command.setDescription(params.description);
  request.setCommand(command);
  await setupAuthToken();
  await grpc.createProject(request);
});

export interface UpdateProjectParams {
  id: string;
  name: string;
  description?: string;
}

export const updateProject = createAsyncThunk<
  void,
  UpdateProjectParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const request = new UpdateProjectRequest();
  request.setId(params.id);
  if (params.name) {
    const renameCommand = new RenameProjectCommand();
    renameCommand.setName(params.name);
    request.setRenameCommand(renameCommand);
  }
  if (params.description) {
    const changeDescCommand = new ChangeDescriptionProjectCommand();
    changeDescCommand.setDescription(params.description);
    request.setChangeDescriptionCommand(changeDescCommand);
  }
  await setupAuthToken();
  await grpc.updateProject(request);
});

export type ProjectsState = typeof initialState;

export const projectsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listProjects.pending, (state) => {
        state.loading = true;
      })
      .addCase(listProjects.fulfilled, (state, action) => {
        projectAdapter.removeAll(state);
        projectAdapter.upsertMany(state, action.payload.projectsList);
        state.loading = false;
        state.totalCount = action.payload.totalCount;
      })
      .addCase(listProjects.rejected, (state) => {
        state.loading = false;
      })
      .addCase(getProject.pending, (state) => {
        state.getProjectError = null;
      })
      .addCase(getProject.fulfilled, (state, action) => {
        state.getProjectError = null;
        if (action.payload) {
          projectAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(getProject.rejected, (state, action) => {
        state.getProjectError = action.error;
      })
      .addCase(enableProject.pending, (state) => {})
      .addCase(enableProject.fulfilled, (state, action) => {})
      .addCase(enableProject.rejected, (state, action) => {})
      .addCase(disableProject.pending, (state) => {})
      .addCase(disableProject.fulfilled, (state, action) => {})
      .addCase(disableProject.rejected, (state, action) => {})
      .addCase(createProject.pending, (state) => {})
      .addCase(createProject.fulfilled, (state, action) => {})
      .addCase(createProject.rejected, (state, action) => {});
  },
});
