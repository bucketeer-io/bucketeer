// package: bucketeer.backend
// file: proto/backend/service.proto

import * as proto_backend_service_pb from "../../proto/backend/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type PublicAPIServiceGetFeature = {
  readonly methodName: string;
  readonly service: typeof PublicAPIService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_backend_service_pb.GetFeatureRequest;
  readonly responseType: typeof proto_backend_service_pb.GetFeatureResponse;
};

type PublicAPIServiceUpdateFeature = {
  readonly methodName: string;
  readonly service: typeof PublicAPIService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_backend_service_pb.UpdateFeatureRequest;
  readonly responseType: typeof proto_backend_service_pb.UpdateFeatureResponse;
};

export class PublicAPIService {
  static readonly serviceName: string;
  static readonly GetFeature: PublicAPIServiceGetFeature;
  static readonly UpdateFeature: PublicAPIServiceUpdateFeature;
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

export class PublicAPIServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getFeature(
    requestMessage: proto_backend_service_pb.GetFeatureRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_backend_service_pb.GetFeatureResponse|null) => void
  ): UnaryResponse;
  getFeature(
    requestMessage: proto_backend_service_pb.GetFeatureRequest,
    callback: (error: ServiceError|null, responseMessage: proto_backend_service_pb.GetFeatureResponse|null) => void
  ): UnaryResponse;
  updateFeature(
    requestMessage: proto_backend_service_pb.UpdateFeatureRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_backend_service_pb.UpdateFeatureResponse|null) => void
  ): UnaryResponse;
  updateFeature(
    requestMessage: proto_backend_service_pb.UpdateFeatureRequest,
    callback: (error: ServiceError|null, responseMessage: proto_backend_service_pb.UpdateFeatureResponse|null) => void
  ): UnaryResponse;
}

