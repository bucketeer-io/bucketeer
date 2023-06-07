// package: bucketeer.test
// file: proto/test/service.proto

import * as proto_test_service_pb from "../../proto/test/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type TestServiceTest = {
  readonly methodName: string;
  readonly service: typeof TestService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_test_service_pb.TestRequest;
  readonly responseType: typeof proto_test_service_pb.TestResponse;
};

export class TestService {
  static readonly serviceName: string;
  static readonly Test: TestServiceTest;
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

export class TestServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  test(
    requestMessage: proto_test_service_pb.TestRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_test_service_pb.TestResponse|null) => void
  ): UnaryResponse;
  test(
    requestMessage: proto_test_service_pb.TestRequest,
    callback: (error: ServiceError|null, responseMessage: proto_test_service_pb.TestResponse|null) => void
  ): UnaryResponse;
}

