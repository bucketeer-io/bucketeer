// package: bucketeer.coderef
// file: proto/coderef/service.proto

import * as proto_coderef_service_pb from '../../proto/coderef/service_pb';
import { grpc } from '@improbable-eng/grpc-web';

type CodeReferenceServiceGetCodeReference = {
  readonly methodName: string;
  readonly service: typeof CodeReferenceService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_coderef_service_pb.GetCodeReferenceRequest;
  readonly responseType: typeof proto_coderef_service_pb.GetCodeReferenceResponse;
};

type CodeReferenceServiceListCodeReferences = {
  readonly methodName: string;
  readonly service: typeof CodeReferenceService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_coderef_service_pb.ListCodeReferencesRequest;
  readonly responseType: typeof proto_coderef_service_pb.ListCodeReferencesResponse;
};

type CodeReferenceServiceCreateCodeReference = {
  readonly methodName: string;
  readonly service: typeof CodeReferenceService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_coderef_service_pb.CreateCodeReferenceRequest;
  readonly responseType: typeof proto_coderef_service_pb.CreateCodeReferenceResponse;
};

type CodeReferenceServiceUpdateCodeReference = {
  readonly methodName: string;
  readonly service: typeof CodeReferenceService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_coderef_service_pb.UpdateCodeReferenceRequest;
  readonly responseType: typeof proto_coderef_service_pb.UpdateCodeReferenceResponse;
};

type CodeReferenceServiceDeleteCodeReference = {
  readonly methodName: string;
  readonly service: typeof CodeReferenceService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_coderef_service_pb.DeleteCodeReferenceRequest;
  readonly responseType: typeof proto_coderef_service_pb.DeleteCodeReferenceResponse;
};

export class CodeReferenceService {
  static readonly serviceName: string;
  static readonly GetCodeReference: CodeReferenceServiceGetCodeReference;
  static readonly ListCodeReferences: CodeReferenceServiceListCodeReferences;
  static readonly CreateCodeReference: CodeReferenceServiceCreateCodeReference;
  static readonly UpdateCodeReference: CodeReferenceServiceUpdateCodeReference;
  static readonly DeleteCodeReference: CodeReferenceServiceDeleteCodeReference;
}

export type ServiceError = {
  message: string;
  code: number;
  metadata: grpc.Metadata;
};
export type Status = { details: string; code: number; metadata: grpc.Metadata };

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
  on(
    type: 'data',
    handler: (message: ResT) => void
  ): BidirectionalStream<ReqT, ResT>;
  on(
    type: 'end',
    handler: (status?: Status) => void
  ): BidirectionalStream<ReqT, ResT>;
  on(
    type: 'status',
    handler: (status: Status) => void
  ): BidirectionalStream<ReqT, ResT>;
}

export class CodeReferenceServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getCodeReference(
    requestMessage: proto_coderef_service_pb.GetCodeReferenceRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_coderef_service_pb.GetCodeReferenceResponse | null
    ) => void
  ): UnaryResponse;
  getCodeReference(
    requestMessage: proto_coderef_service_pb.GetCodeReferenceRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_coderef_service_pb.GetCodeReferenceResponse | null
    ) => void
  ): UnaryResponse;
  listCodeReferences(
    requestMessage: proto_coderef_service_pb.ListCodeReferencesRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_coderef_service_pb.ListCodeReferencesResponse | null
    ) => void
  ): UnaryResponse;
  listCodeReferences(
    requestMessage: proto_coderef_service_pb.ListCodeReferencesRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_coderef_service_pb.ListCodeReferencesResponse | null
    ) => void
  ): UnaryResponse;
  createCodeReference(
    requestMessage: proto_coderef_service_pb.CreateCodeReferenceRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_coderef_service_pb.CreateCodeReferenceResponse | null
    ) => void
  ): UnaryResponse;
  createCodeReference(
    requestMessage: proto_coderef_service_pb.CreateCodeReferenceRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_coderef_service_pb.CreateCodeReferenceResponse | null
    ) => void
  ): UnaryResponse;
  updateCodeReference(
    requestMessage: proto_coderef_service_pb.UpdateCodeReferenceRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_coderef_service_pb.UpdateCodeReferenceResponse | null
    ) => void
  ): UnaryResponse;
  updateCodeReference(
    requestMessage: proto_coderef_service_pb.UpdateCodeReferenceRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_coderef_service_pb.UpdateCodeReferenceResponse | null
    ) => void
  ): UnaryResponse;
  deleteCodeReference(
    requestMessage: proto_coderef_service_pb.DeleteCodeReferenceRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_coderef_service_pb.DeleteCodeReferenceResponse | null
    ) => void
  ): UnaryResponse;
  deleteCodeReference(
    requestMessage: proto_coderef_service_pb.DeleteCodeReferenceRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_coderef_service_pb.DeleteCodeReferenceResponse | null
    ) => void
  ): UnaryResponse;
}
