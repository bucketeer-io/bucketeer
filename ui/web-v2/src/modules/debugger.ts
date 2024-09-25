import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk
} from '@reduxjs/toolkit';

import { Evaluation } from '../proto/feature/evaluation_pb';
import * as featuresGrpc from '../grpc/features';

import { AppState } from '.';
import {
  EvaluateFeaturesRequest,
  EvaluateFeaturesResponse
} from '../proto/feature/service_pb';
import { User } from '../proto/user/user_pb';

const MODULE_NAME = 'debugger';

export const debuggerAdapter = createEntityAdapter({
  selectId: (evaluation: Evaluation.AsObject) => evaluation.id
});

export const { selectAll, selectById } = debuggerAdapter.getSelectors();

interface evaluateFeaturesParams {
  environmentNamespace: string;
  flag: string;
  userId: string;
  userAttributes: Array<[string, string]>;
}

export const evaluateFeatures = createAsyncThunk<
  EvaluateFeaturesResponse.AsObject,
  evaluateFeaturesParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/evaluate`, async (params) => {
  const request = new EvaluateFeaturesRequest();
  const user = new User();
  user.setId(params.userId);
  const dataMap = user.getDataMap();
  for (const [key, value] of params.userAttributes) {
    dataMap.set(key, value);
  }

  request.setFeatureId(params.flag);
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setUser(user);
  const result = await featuresGrpc.evaluateFeatures(request);
  return result.response.toObject();
});

const initialState = debuggerAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
}>({
  loading: false,
  totalCount: 0
});

export type AuditLogsState = typeof initialState;

export const auditLogSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(evaluateFeatures.pending, (state) => {
        state.loading = true;
      })
      .addCase(evaluateFeatures.fulfilled, (state) => {
        // debuggerAdapter.removeAll(state);
        // debuggerAdapter.upsertMany(state, action.payload.userEvaluations);
        state.loading = false;
      })
      .addCase(evaluateFeatures.rejected, (state) => {
        state.loading = false;
      });
  }
});
