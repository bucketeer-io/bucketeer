import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError
} from '@reduxjs/toolkit';
import * as jspb from 'google-protobuf';
import { Any } from 'google-protobuf/google/protobuf/any_pb';
import { BoolValue } from 'google-protobuf/google/protobuf/wrappers_pb';

import * as segmentGrpc from '../grpc/segments';
import {
  CreateSegmentCommand,
  BulkUploadSegmentUsersCommand,
  DeleteSegmentCommand,
  ChangeSegmentNameCommand,
  ChangeSegmentDescriptionCommand,
  Command
} from '../proto/feature/command_pb';
import { Segment } from '../proto/feature/segment_pb';
import {
  ListSegmentsRequest,
  GetSegmentRequest,
  BulkDownloadSegmentUsersRequest,
  BulkUploadSegmentUsersRequest,
  CreateSegmentRequest,
  UpdateSegmentRequest,
  DeleteSegmentRequest,
  ListSegmentsResponse
} from '../proto/feature/service_pb';

import { AppState } from '.';

const MODULE_NAME = 'segments';

export const segmentsAdapter = createEntityAdapter<Segment.AsObject>({
  selectId: (segment) => segment.id
});

export const { selectAll, selectById } = segmentsAdapter.getSelectors();

export type OrderBy =
  ListSegmentsRequest.OrderByMap[keyof ListSegmentsRequest.OrderByMap];
export type OrderDirection =
  ListSegmentsRequest.OrderDirectionMap[keyof ListSegmentsRequest.OrderDirectionMap];

export interface ListSegmentsParams {
  environmentNamespace: string;
  pageSize?: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  inUse?: boolean;
}

export const listSegments = createAsyncThunk<
  ListSegmentsResponse.AsObject,
  ListSegmentsParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/list`, async (params) => {
  const request = new ListSegmentsRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setCursor(params.cursor);
  if (params.pageSize) {
    request.setPageSize(params.pageSize);
  }
  request.setOrderBy(params.orderBy);
  request.setOrderDirection(params.orderDirection);
  request.setSearchKeyword(params.searchKeyword);
  params.inUse != null &&
    request.setIsInUseStatus(new BoolValue().setValue(params.inUse));
  const result = await segmentGrpc.listSegments(request);
  return result.response.toObject();
});

export interface GetSegmentParams {
  environmentNamespace: string;
  id: string;
}

export const getSegment = createAsyncThunk<
  Segment.AsObject,
  GetSegmentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/get`, async (params) => {
  const request = new GetSegmentRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  const result = await segmentGrpc.getSegment(request);
  return result.response.toObject().segment;
});

export interface BulkDownloadSegmentUsersParams {
  segmentId: string;
  environmentNamespace: string;
}

export const bulkDownloadSegmentUsers = createAsyncThunk<
  Uint8Array | string,
  BulkDownloadSegmentUsersParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/download`, async (params) => {
  const request = new BulkDownloadSegmentUsersRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setSegmentId(params.segmentId);
  const result = await segmentGrpc.bulkDownloadSegmentUsers(request);
  return result.response.toObject().data;
});

export interface BulkUploadSegmentUsersParams {
  environmentNamespace: string;
  segmentId: string;
  data: Uint8Array;
}

export const bulkUploadSegmentUsers = createAsyncThunk<
  void,
  BulkUploadSegmentUsersParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/upload`, async (params) => {
  const cmd = new BulkUploadSegmentUsersCommand();
  cmd.setData(params.data);
  const request = new BulkUploadSegmentUsersRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setSegmentId(params.segmentId);
  request.setCommand(cmd);
  await segmentGrpc.bulkUploadSegmentUsers(request);
});

export interface CreateSegmentParams {
  environmentNamespace: string;
  name: string;
  description: string;
}

export const createSegment = createAsyncThunk<
  string,
  CreateSegmentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/create`, async (params) => {
  const cmd = new CreateSegmentCommand();
  cmd.setName(params.name);
  cmd.setDescription(params.description);
  const request = new CreateSegmentRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setCommand(cmd);
  const result = await segmentGrpc.createSegment(request);
  return result.response.toObject().segment.id;
});

const convertCommandToAny = (
  command: jspb.Message,
  commandName: string
): Command => {
  const result = new Command();
  const pbAny = new Any();
  pbAny.pack(command.serializeBinary(), `bucketeer.feature.${commandName}`);

  result.setCommand(pbAny);
  return result;
};

export interface UpdateSegmentParams {
  environmentNamespace: string;
  id: string;
  name: string;
  description: string;
}

export const updateSegment = createAsyncThunk<
  void,
  UpdateSegmentParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/update`, async (params) => {
  const cmdList = [];
  if (params.name) {
    const cmd = new ChangeSegmentNameCommand();
    cmd.setName(params.name);
    cmdList.push(convertCommandToAny(cmd, 'ChangeSegmentNameCommand'));
  }
  if (params.description) {
    const cmd = new ChangeSegmentDescriptionCommand();
    cmd.setDescription(params.description.toString());
    cmdList.push(convertCommandToAny(cmd, 'ChangeSegmentDescriptionCommand'));
  }
  const request = new UpdateSegmentRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommandsList(cmdList);
  await segmentGrpc.updateSegment(request);
});

export interface DeleteSegmentUsersParams {
  id: string;
  environmentNamespace: string;
}

export const deleteSegmentUser = createAsyncThunk<
  string,
  DeleteSegmentUsersParams | undefined,
  { state: AppState }
>(`${MODULE_NAME}/delete`, async (params) => {
  const request = new DeleteSegmentRequest();
  request.setEnvironmentNamespace(params.environmentNamespace);
  request.setId(params.id);
  request.setCommand(new DeleteSegmentCommand());
  await segmentGrpc.deleteSegment(request);
  return params.id;
});

const initialState = segmentsAdapter.getInitialState<{
  loading: boolean;
  totalCount: number;
  getSegmentError: SerializedError | null;
}>({
  loading: false,
  totalCount: 0,
  getSegmentError: null
});

export type SegmentsState = typeof initialState;

export const segmentsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(listSegments.pending, (state) => {
        state.loading = true;
      })
      .addCase(listSegments.fulfilled, (state, action) => {
        segmentsAdapter.removeAll(state);
        segmentsAdapter.upsertMany(state, action.payload.segmentsList);
        state.loading = false;
        state.totalCount = action.payload.totalCount;
      })
      .addCase(listSegments.rejected, (state) => {
        state.loading = false;
      })
      .addCase(getSegment.pending, (state) => {
        state.getSegmentError = null;
      })
      .addCase(getSegment.fulfilled, (state, action) => {
        state.getSegmentError = null;
        if (action.payload) {
          segmentsAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(getSegment.rejected, (state, action) => {
        state.getSegmentError = action.error;
      })
      .addCase(bulkDownloadSegmentUsers.pending, () => {})
      .addCase(bulkDownloadSegmentUsers.fulfilled, () => {})
      .addCase(bulkDownloadSegmentUsers.rejected, () => {})
      .addCase(bulkUploadSegmentUsers.pending, () => {})
      .addCase(bulkUploadSegmentUsers.fulfilled, () => {})
      .addCase(bulkUploadSegmentUsers.rejected, () => {})
      .addCase(createSegment.pending, () => {})
      .addCase(createSegment.fulfilled, () => {})
      .addCase(createSegment.rejected, () => {})
      .addCase(updateSegment.pending, () => {})
      .addCase(updateSegment.fulfilled, () => {})
      .addCase(updateSegment.rejected, () => {})
      .addCase(deleteSegmentUser.pending, () => {})
      .addCase(deleteSegmentUser.fulfilled, (state, action) => {
        segmentsAdapter.removeOne(state, action.payload);
      })
      .addCase(deleteSegmentUser.rejected, () => {});
  }
});
