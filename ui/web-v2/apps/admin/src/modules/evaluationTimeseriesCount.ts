import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';

import * as grpc from '../grpc/eventcounter';
import {
  GetEvaluationTimeseriesCountRequest,
  GetEvaluationTimeseriesCountResponse,
} from '../proto/eventcounter/service_pb';
import { VariationTimeseries } from '../proto/eventcounter/timeseries_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'evaluationTimeseriesCount';

export interface EvaluationTimeseriesCount {
  loading: boolean;
  userCountsList: Array<VariationTimeseries.AsObject>;
  eventCountsList: Array<VariationTimeseries.AsObject>;
}

export interface GetEvaluationTimeseriesCountParams {
  environmentNamespace: string;
  featureId: string;
  timeRange: GetEvaluationTimeseriesCountRequest.TimeRangeMap[keyof GetEvaluationTimeseriesCountRequest.TimeRangeMap];
}

export const getEvaluationTimeseriesCount = createAsyncThunk<
  GetEvaluationTimeseriesCountResponse.AsObject,
  GetEvaluationTimeseriesCountParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/getEvaluationTimeseriesCount`, async (params) => {
  const request = new GetEvaluationTimeseriesCountRequest();

  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setFeatureId(params.featureId);
  request.setTimeRange(params.timeRange);
  await setupAuthToken();
  const result = await grpc.getEvaluationTimeseriesCount(request);
  return result.response.toObject();
});

const initialState: EvaluationTimeseriesCount = {
  loading: false,
  userCountsList: [],
  eventCountsList: [],
};

export type EvaluationTimeseriesCountState = typeof initialState;

export const evaluationTimeseriesCountSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(getEvaluationTimeseriesCount.pending, (state) => {
        state.loading = true;
      })
      .addCase(getEvaluationTimeseriesCount.fulfilled, (state, action) => {
        state.userCountsList = action.payload.userCountsList;
        state.eventCountsList = action.payload.eventCountsList;
        state.loading = false;
      })
      .addCase(getEvaluationTimeseriesCount.rejected, (state) => {
        state.loading = false;
      });
  },
});
