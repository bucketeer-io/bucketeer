// package: bucketeer.eventcounter
// file: proto/eventcounter/service.proto

import * as proto_eventcounter_service_pb from "../../proto/eventcounter/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type EventCounterServiceGetExperimentEvaluationCount = {
  readonly methodName: string;
  readonly service: typeof EventCounterService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_eventcounter_service_pb.GetExperimentEvaluationCountRequest;
  readonly responseType: typeof proto_eventcounter_service_pb.GetExperimentEvaluationCountResponse;
};

type EventCounterServiceGetEvaluationTimeseriesCount = {
  readonly methodName: string;
  readonly service: typeof EventCounterService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_eventcounter_service_pb.GetEvaluationTimeseriesCountRequest;
  readonly responseType: typeof proto_eventcounter_service_pb.GetEvaluationTimeseriesCountResponse;
};

type EventCounterServiceGetExperimentResult = {
  readonly methodName: string;
  readonly service: typeof EventCounterService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_eventcounter_service_pb.GetExperimentResultRequest;
  readonly responseType: typeof proto_eventcounter_service_pb.GetExperimentResultResponse;
};

type EventCounterServiceListExperimentResults = {
  readonly methodName: string;
  readonly service: typeof EventCounterService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_eventcounter_service_pb.ListExperimentResultsRequest;
  readonly responseType: typeof proto_eventcounter_service_pb.ListExperimentResultsResponse;
};

type EventCounterServiceGetExperimentGoalCount = {
  readonly methodName: string;
  readonly service: typeof EventCounterService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_eventcounter_service_pb.GetExperimentGoalCountRequest;
  readonly responseType: typeof proto_eventcounter_service_pb.GetExperimentGoalCountResponse;
};

type EventCounterServiceGetMAUCount = {
  readonly methodName: string;
  readonly service: typeof EventCounterService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_eventcounter_service_pb.GetMAUCountRequest;
  readonly responseType: typeof proto_eventcounter_service_pb.GetMAUCountResponse;
};

type EventCounterServiceSummarizeMAUCounts = {
  readonly methodName: string;
  readonly service: typeof EventCounterService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_eventcounter_service_pb.SummarizeMAUCountsRequest;
  readonly responseType: typeof proto_eventcounter_service_pb.SummarizeMAUCountsResponse;
};

type EventCounterServiceGetOpsEvaluationUserCount = {
  readonly methodName: string;
  readonly service: typeof EventCounterService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_eventcounter_service_pb.GetOpsEvaluationUserCountRequest;
  readonly responseType: typeof proto_eventcounter_service_pb.GetOpsEvaluationUserCountResponse;
};

type EventCounterServiceGetOpsGoalUserCount = {
  readonly methodName: string;
  readonly service: typeof EventCounterService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_eventcounter_service_pb.GetOpsGoalUserCountRequest;
  readonly responseType: typeof proto_eventcounter_service_pb.GetOpsGoalUserCountResponse;
};

export class EventCounterService {
  static readonly serviceName: string;
  static readonly GetExperimentEvaluationCount: EventCounterServiceGetExperimentEvaluationCount;
  static readonly GetEvaluationTimeseriesCount: EventCounterServiceGetEvaluationTimeseriesCount;
  static readonly GetExperimentResult: EventCounterServiceGetExperimentResult;
  static readonly ListExperimentResults: EventCounterServiceListExperimentResults;
  static readonly GetExperimentGoalCount: EventCounterServiceGetExperimentGoalCount;
  static readonly GetMAUCount: EventCounterServiceGetMAUCount;
  static readonly SummarizeMAUCounts: EventCounterServiceSummarizeMAUCounts;
  static readonly GetOpsEvaluationUserCount: EventCounterServiceGetOpsEvaluationUserCount;
  static readonly GetOpsGoalUserCount: EventCounterServiceGetOpsGoalUserCount;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class EventCounterServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getExperimentEvaluationCount(
    requestMessage: proto_eventcounter_service_pb.GetExperimentEvaluationCountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetExperimentEvaluationCountResponse|null) => void
  ): UnaryResponse;
  getExperimentEvaluationCount(
    requestMessage: proto_eventcounter_service_pb.GetExperimentEvaluationCountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetExperimentEvaluationCountResponse|null) => void
  ): UnaryResponse;
  getEvaluationTimeseriesCount(
    requestMessage: proto_eventcounter_service_pb.GetEvaluationTimeseriesCountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetEvaluationTimeseriesCountResponse|null) => void
  ): UnaryResponse;
  getEvaluationTimeseriesCount(
    requestMessage: proto_eventcounter_service_pb.GetEvaluationTimeseriesCountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetEvaluationTimeseriesCountResponse|null) => void
  ): UnaryResponse;
  getExperimentResult(
    requestMessage: proto_eventcounter_service_pb.GetExperimentResultRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetExperimentResultResponse|null) => void
  ): UnaryResponse;
  getExperimentResult(
    requestMessage: proto_eventcounter_service_pb.GetExperimentResultRequest,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetExperimentResultResponse|null) => void
  ): UnaryResponse;
  listExperimentResults(
    requestMessage: proto_eventcounter_service_pb.ListExperimentResultsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.ListExperimentResultsResponse|null) => void
  ): UnaryResponse;
  listExperimentResults(
    requestMessage: proto_eventcounter_service_pb.ListExperimentResultsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.ListExperimentResultsResponse|null) => void
  ): UnaryResponse;
  getExperimentGoalCount(
    requestMessage: proto_eventcounter_service_pb.GetExperimentGoalCountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetExperimentGoalCountResponse|null) => void
  ): UnaryResponse;
  getExperimentGoalCount(
    requestMessage: proto_eventcounter_service_pb.GetExperimentGoalCountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetExperimentGoalCountResponse|null) => void
  ): UnaryResponse;
  getMAUCount(
    requestMessage: proto_eventcounter_service_pb.GetMAUCountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetMAUCountResponse|null) => void
  ): UnaryResponse;
  getMAUCount(
    requestMessage: proto_eventcounter_service_pb.GetMAUCountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetMAUCountResponse|null) => void
  ): UnaryResponse;
  summarizeMAUCounts(
    requestMessage: proto_eventcounter_service_pb.SummarizeMAUCountsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.SummarizeMAUCountsResponse|null) => void
  ): UnaryResponse;
  summarizeMAUCounts(
    requestMessage: proto_eventcounter_service_pb.SummarizeMAUCountsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.SummarizeMAUCountsResponse|null) => void
  ): UnaryResponse;
  getOpsEvaluationUserCount(
    requestMessage: proto_eventcounter_service_pb.GetOpsEvaluationUserCountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetOpsEvaluationUserCountResponse|null) => void
  ): UnaryResponse;
  getOpsEvaluationUserCount(
    requestMessage: proto_eventcounter_service_pb.GetOpsEvaluationUserCountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetOpsEvaluationUserCountResponse|null) => void
  ): UnaryResponse;
  getOpsGoalUserCount(
    requestMessage: proto_eventcounter_service_pb.GetOpsGoalUserCountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetOpsGoalUserCountResponse|null) => void
  ): UnaryResponse;
  getOpsGoalUserCount(
    requestMessage: proto_eventcounter_service_pb.GetOpsGoalUserCountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_eventcounter_service_pb.GetOpsGoalUserCountResponse|null) => void
  ): UnaryResponse;
}

