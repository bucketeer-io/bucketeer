// package: bucketeer.auditlog
// file: proto/auditlog/service.proto

import * as proto_auditlog_service_pb from "../../proto/auditlog/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type AuditLogServiceListAuditLogs = {
  readonly methodName: string;
  readonly service: typeof AuditLogService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auditlog_service_pb.ListAuditLogsRequest;
  readonly responseType: typeof proto_auditlog_service_pb.ListAuditLogsResponse;
};

type AuditLogServiceListAdminAuditLogs = {
  readonly methodName: string;
  readonly service: typeof AuditLogService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auditlog_service_pb.ListAdminAuditLogsRequest;
  readonly responseType: typeof proto_auditlog_service_pb.ListAdminAuditLogsResponse;
};

type AuditLogServiceListFeatureHistory = {
  readonly methodName: string;
  readonly service: typeof AuditLogService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auditlog_service_pb.ListFeatureHistoryRequest;
  readonly responseType: typeof proto_auditlog_service_pb.ListFeatureHistoryResponse;
};

export class AuditLogService {
  static readonly serviceName: string;
  static readonly ListAuditLogs: AuditLogServiceListAuditLogs;
  static readonly ListAdminAuditLogs: AuditLogServiceListAdminAuditLogs;
  static readonly ListFeatureHistory: AuditLogServiceListFeatureHistory;
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

export class AuditLogServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listAuditLogs(
    requestMessage: proto_auditlog_service_pb.ListAuditLogsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_auditlog_service_pb.ListAuditLogsResponse|null) => void
  ): UnaryResponse;
  listAuditLogs(
    requestMessage: proto_auditlog_service_pb.ListAuditLogsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_auditlog_service_pb.ListAuditLogsResponse|null) => void
  ): UnaryResponse;
  listAdminAuditLogs(
    requestMessage: proto_auditlog_service_pb.ListAdminAuditLogsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_auditlog_service_pb.ListAdminAuditLogsResponse|null) => void
  ): UnaryResponse;
  listAdminAuditLogs(
    requestMessage: proto_auditlog_service_pb.ListAdminAuditLogsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_auditlog_service_pb.ListAdminAuditLogsResponse|null) => void
  ): UnaryResponse;
  listFeatureHistory(
    requestMessage: proto_auditlog_service_pb.ListFeatureHistoryRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_auditlog_service_pb.ListFeatureHistoryResponse|null) => void
  ): UnaryResponse;
  listFeatureHistory(
    requestMessage: proto_auditlog_service_pb.ListFeatureHistoryRequest,
    callback: (error: ServiceError|null, responseMessage: proto_auditlog_service_pb.ListFeatureHistoryResponse|null) => void
  ): UnaryResponse;
}

