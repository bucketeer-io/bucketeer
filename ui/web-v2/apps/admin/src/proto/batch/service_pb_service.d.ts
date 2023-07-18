// package: bucketeer.batch
// file: proto/batch/service.proto

import * as proto_batch_service_pb from "../../proto/batch/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type BatchServiceExecuteBatchJob = {
  readonly methodName: string;
  readonly service: typeof BatchService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_batch_service_pb.BatchJobRequest;
  readonly responseType: typeof proto_batch_service_pb.BatchJobResponse;
};

export class BatchService {
  static readonly serviceName: string;
  static readonly ExecuteBatchJob: BatchServiceExecuteBatchJob;
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

export class BatchServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  executeBatchJob(
    requestMessage: proto_batch_service_pb.BatchJobRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_batch_service_pb.BatchJobResponse|null) => void
  ): UnaryResponse;
  executeBatchJob(
    requestMessage: proto_batch_service_pb.BatchJobRequest,
    callback: (error: ServiceError|null, responseMessage: proto_batch_service_pb.BatchJobResponse|null) => void
  ): UnaryResponse;
}

