import {
  createSlice,
  createAsyncThunk,
  SerializedError,
} from '@reduxjs/toolkit';

import * as grpc from '../grpc/eventcounter';
import { Filter } from '../proto/eventcounter/filter_pb';
import {
  GetGoalCountRequest,
  GetGoalCountResponse,
} from '../proto/eventcounter/service_pb';
import { Row } from '../proto/eventcounter/table_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'goalCounts';

export interface GetGoalCountParams {
  environmentNamespace: string;
  startAt: Date;
  endAt: Date;
  featureId?: string;
  featureVersion?: number;
  reason?: string;
  goalId: string;
  filters?: Array<Filter.AsObject>;
  segments?: Array<string>;
}

export const getGoalCount = createAsyncThunk<
  GetGoalCountResponse.AsObject,
  GetGoalCountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/getGoalCount`, async (params) => {
  const request = new GetGoalCountRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setStartAt(Math.floor(params.startAt.getTime() / 1000));
  request.setEndAt(Math.floor(params.endAt.getTime() / 1000));
  request.setFeatureId(params.featureId ? params.featureId : '');
  request.setFeatureVersion(params.featureVersion ? params.featureVersion : 0);
  if (params.reason) {
    request.setReason(params.reason);
  }
  request.setGoalId(params.goalId);
  if (params.segments) {
    request.setSegmentsList(params.segments);
  }
  if (params.filters) {
    const filters = params.filters.map((filter) => {
      const f = new Filter();
      f.setKey(filter.key);
      f.setOperator(filter.operator);
      f.setValuesList(filter.valuesList);
      return f;
    });
    request.setFiltersList(filters);
  }
  await setupAuthToken();
  const result = await grpc.getGoalCount(request);
  return result.response.toObject();
});

const initialState: {
  headers?: Row.AsObject;
  rows?: Array<Row.AsObject>;
  loading: boolean;
  getGoalCountError: SerializedError | null;
} = {
  headers: null,
  rows: null,
  loading: false,
  getGoalCountError: null,
};

export const goalCountsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(getGoalCount.pending, (state) => {
        state.loading = true;
      })
      .addCase(getGoalCount.fulfilled, (state, action) => {
        state.headers = action.payload.headers;
        state.rows = action.payload.rowsList;
        state.loading = false;
      })
      .addCase(getGoalCount.rejected, (state) => {
        state.loading = false;
      });
  },
});
