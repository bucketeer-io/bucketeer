import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from '@reduxjs/toolkit';
import { Int32Value } from 'google-protobuf/google/protobuf/wrappers_pb';

import * as auditLogGrpc from '../grpc/auditLog';
import { AuditLog } from '../proto/auditlog/auditlog_pb';
import {
  ListAdminAuditLogsRequest,
  ListAdminAuditLogsResponse,
  ListAuditLogsRequest,
  ListAuditLogsResponse,
  ListFeatureHistoryRequest,
  ListFeatureHistoryResponse,
} from '../proto/auditlog/service_pb';

import { setupAuthToken } from './auth';

import { AppState } from '.';

const MODULE_NAME = 'auditLogs';

export const auditLogsAdapter = createEntityAdapter<AuditLog.AsObject>({
  selectId: (auditLog) => auditLog.id,
});

export const { selectAll, selectById } = auditLogsAdapter.getSelectors();

export type AdminOrderBy =
  ListAdminAuditLogsRequest.OrderByMap[keyof ListAdminAuditLogsRequest.OrderByMap];
export type AdminOrderDirection =
  ListAdminAuditLogsRequest.OrderDirectionMap[keyof ListAdminAuditLogsRequest.OrderDirectionMap];

export type OrderBy =
  ListAuditLogsRequest.OrderByMap[keyof ListAuditLogsRequest.OrderByMap];
export type OrderDirection =
  ListAuditLogsRequest.OrderDirectionMap[keyof ListAuditLogsRequest.OrderDirectionMap];

export interface ListAdminAuditLogsParams {
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  from: number;
  to: number;
  resource?: number;
}

export const listAdminAuditLogs = createAsyncThunk<
  ListAdminAuditLogsResponse.AsObject,
  ListAdminAuditLogsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/listAdminAuditLogs`, async (params) => {
  const request = new ListAdminAuditLogsRequest();
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  request.setFrom(params.from);
  request.setTo(params.to);
  params.resource != null &&
    request.setEntityType(new Int32Value().setValue(params.resource));
  await setupAuthToken();
  const result = await auditLogGrpc.listAdminAuditLogs(request);
  return result.response.toObject();
});

export interface ListAuditLogsParams {
  pageSize: number;
  cursor: string;
  environmentNamespace: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  from: number;
  to: number;
  entityType?: number;
}

export const listAuditLogs = createAsyncThunk<
  ListAuditLogsResponse.AsObject,
  ListAuditLogsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/listAuditLogs`, async (params) => {
  const request = new ListAuditLogsRequest();
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  request.setFrom(params.from);
  request.setTo(params.to);
  params.entityType != null &&
    request.setEntityType(new Int32Value().setValue(params.entityType));
  await setupAuthToken();
  const result = await auditLogGrpc.listAuditLogs(request);
  return result.response.toObject();
});

export interface ListFeatureHistoryParams {
  featureId: string;
  pageSize: number;
  cursor: string;
  environmentNamespace: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  from: number;
  to: number;
}

export const listFeatureHistory = createAsyncThunk<
  ListFeatureHistoryResponse.AsObject,
  ListFeatureHistoryParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/listFeatureHistory`, async (params) => {
  const request = new ListFeatureHistoryRequest();
  request.setFeatureId(params.featureId);
  request.setPageSize(params.pageSize);
  request.setCursor(params.cursor);
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  request.setFrom(params.from);
  request.setTo(params.to);
  await setupAuthToken();
  const result = await auditLogGrpc.listFeatureHistory(request);
  return result.response.toObject();
});

const initialState = auditLogsAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
}>({
  loading: false,
  totalCount: 0,
});

export type AuditLogsState = typeof initialState;

export const auditLogSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listAdminAuditLogs.pending, (state) => {
        state.loading = true;
      })
      .addCase(listAdminAuditLogs.fulfilled, (state, action) => {
        auditLogsAdapter.removeAll(state);
        auditLogsAdapter.upsertMany(state, action.payload.auditLogsList);
        state.loading = false;
        state.totalCount = action.payload.totalCount;
      })
      .addCase(listAdminAuditLogs.rejected, (state) => {
        state.loading = false;
      })
      .addCase(listAuditLogs.pending, (state) => {
        state.loading = true;
      })
      .addCase(listAuditLogs.fulfilled, (state, action) => {
        auditLogsAdapter.removeAll(state);
        auditLogsAdapter.upsertMany(state, action.payload.auditLogsList);
        state.loading = false;
        state.totalCount = action.payload.totalCount;
      })
      .addCase(listAuditLogs.rejected, (state) => {
        state.loading = false;
      })
      .addCase(listFeatureHistory.pending, (state) => {
        state.loading = true;
      })
      .addCase(listFeatureHistory.fulfilled, (state, action) => {
        auditLogsAdapter.removeAll(state);
        auditLogsAdapter.upsertMany(state, action.payload.auditLogsList);
        state.loading = false;
        state.totalCount = action.payload.totalCount;
      })
      .addCase(listFeatureHistory.rejected, (state) => {
        state.loading = false;
      });
  },
});
