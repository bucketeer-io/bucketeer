// package: bucketeer.push
// file: proto/push/service.proto

import * as proto_push_service_pb from "../../proto/push/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type PushServiceListPushes = {
  readonly methodName: string;
  readonly service: typeof PushService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_push_service_pb.ListPushesRequest;
  readonly responseType: typeof proto_push_service_pb.ListPushesResponse;
};

type PushServiceCreatePush = {
  readonly methodName: string;
  readonly service: typeof PushService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_push_service_pb.CreatePushRequest;
  readonly responseType: typeof proto_push_service_pb.CreatePushResponse;
};

type PushServiceDeletePush = {
  readonly methodName: string;
  readonly service: typeof PushService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_push_service_pb.DeletePushRequest;
  readonly responseType: typeof proto_push_service_pb.DeletePushResponse;
};

type PushServiceUpdatePush = {
  readonly methodName: string;
  readonly service: typeof PushService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_push_service_pb.UpdatePushRequest;
  readonly responseType: typeof proto_push_service_pb.UpdatePushResponse;
};

export class PushService {
  static readonly serviceName: string;
  static readonly ListPushes: PushServiceListPushes;
  static readonly CreatePush: PushServiceCreatePush;
  static readonly DeletePush: PushServiceDeletePush;
  static readonly UpdatePush: PushServiceUpdatePush;
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

export class PushServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listPushes(
    requestMessage: proto_push_service_pb.ListPushesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_push_service_pb.ListPushesResponse|null) => void
  ): UnaryResponse;
  listPushes(
    requestMessage: proto_push_service_pb.ListPushesRequest,
    callback: (error: ServiceError|null, responseMessage: proto_push_service_pb.ListPushesResponse|null) => void
  ): UnaryResponse;
  createPush(
    requestMessage: proto_push_service_pb.CreatePushRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_push_service_pb.CreatePushResponse|null) => void
  ): UnaryResponse;
  createPush(
    requestMessage: proto_push_service_pb.CreatePushRequest,
    callback: (error: ServiceError|null, responseMessage: proto_push_service_pb.CreatePushResponse|null) => void
  ): UnaryResponse;
  deletePush(
    requestMessage: proto_push_service_pb.DeletePushRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_push_service_pb.DeletePushResponse|null) => void
  ): UnaryResponse;
  deletePush(
    requestMessage: proto_push_service_pb.DeletePushRequest,
    callback: (error: ServiceError|null, responseMessage: proto_push_service_pb.DeletePushResponse|null) => void
  ): UnaryResponse;
  updatePush(
    requestMessage: proto_push_service_pb.UpdatePushRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_push_service_pb.UpdatePushResponse|null) => void
  ): UnaryResponse;
  updatePush(
    requestMessage: proto_push_service_pb.UpdatePushRequest,
    callback: (error: ServiceError|null, responseMessage: proto_push_service_pb.UpdatePushResponse|null) => void
  ): UnaryResponse;
}

