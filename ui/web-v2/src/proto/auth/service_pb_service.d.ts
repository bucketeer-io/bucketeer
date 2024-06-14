// package: bucketeer.auth
// file: proto/auth/service.proto

import * as proto_auth_service_pb from '../../proto/auth/service_pb';
import { grpc } from '@improbable-eng/grpc-web';

type AuthServiceGetAuthCodeURL = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.GetAuthCodeURLRequest;
  readonly responseType: typeof proto_auth_service_pb.GetAuthCodeURLResponse;
};

type AuthServiceExchangeToken = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.ExchangeTokenRequest;
  readonly responseType: typeof proto_auth_service_pb.ExchangeTokenResponse;
};

type AuthServiceRefreshToken = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.RefreshTokenRequest;
  readonly responseType: typeof proto_auth_service_pb.RefreshTokenResponse;
};

type AuthServiceGetAuthenticationURL = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.GetAuthenticationURLRequest;
  readonly responseType: typeof proto_auth_service_pb.GetAuthenticationURLResponse;
};

type AuthServiceExchangeBucketeerToken = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.ExchangeBucketeerTokenRequest;
  readonly responseType: typeof proto_auth_service_pb.ExchangeBucketeerTokenResponse;
};

type AuthServiceRefreshBucketeerToken = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.RefreshBucketeerTokenRequest;
  readonly responseType: typeof proto_auth_service_pb.RefreshBucketeerTokenResponse;
};

export class AuthService {
  static readonly serviceName: string;
  static readonly GetAuthCodeURL: AuthServiceGetAuthCodeURL;
  static readonly ExchangeToken: AuthServiceExchangeToken;
  static readonly RefreshToken: AuthServiceRefreshToken;
  static readonly GetAuthenticationURL: AuthServiceGetAuthenticationURL;
  static readonly ExchangeBucketeerToken: AuthServiceExchangeBucketeerToken;
  static readonly RefreshBucketeerToken: AuthServiceRefreshBucketeerToken;
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

export class AuthServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getAuthCodeURL(
    requestMessage: proto_auth_service_pb.GetAuthCodeURLRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.GetAuthCodeURLResponse | null
    ) => void
  ): UnaryResponse;
  getAuthCodeURL(
    requestMessage: proto_auth_service_pb.GetAuthCodeURLRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.GetAuthCodeURLResponse | null
    ) => void
  ): UnaryResponse;
  exchangeToken(
    requestMessage: proto_auth_service_pb.ExchangeTokenRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.ExchangeTokenResponse | null
    ) => void
  ): UnaryResponse;
  exchangeToken(
    requestMessage: proto_auth_service_pb.ExchangeTokenRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.ExchangeTokenResponse | null
    ) => void
  ): UnaryResponse;
  refreshToken(
    requestMessage: proto_auth_service_pb.RefreshTokenRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.RefreshTokenResponse | null
    ) => void
  ): UnaryResponse;
  refreshToken(
    requestMessage: proto_auth_service_pb.RefreshTokenRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.RefreshTokenResponse | null
    ) => void
  ): UnaryResponse;
  getAuthenticationURL(
    requestMessage: proto_auth_service_pb.GetAuthenticationURLRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_auth_service_pb.GetAuthenticationURLResponse|null) => void
  ): UnaryResponse;
  getAuthenticationURL(
    requestMessage: proto_auth_service_pb.GetAuthenticationURLRequest,
    callback: (error: ServiceError|null, responseMessage: proto_auth_service_pb.GetAuthenticationURLResponse|null) => void
  ): UnaryResponse;
  exchangeBucketeerToken(
    requestMessage: proto_auth_service_pb.ExchangeBucketeerTokenRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_auth_service_pb.ExchangeBucketeerTokenResponse|null) => void
  ): UnaryResponse;
  exchangeBucketeerToken(
    requestMessage: proto_auth_service_pb.ExchangeBucketeerTokenRequest,
    callback: (error: ServiceError|null, responseMessage: proto_auth_service_pb.ExchangeBucketeerTokenResponse|null) => void
  ): UnaryResponse;
  refreshBucketeerToken(
    requestMessage: proto_auth_service_pb.RefreshBucketeerTokenRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_auth_service_pb.RefreshBucketeerTokenResponse|null) => void
  ): UnaryResponse;
  refreshBucketeerToken(
    requestMessage: proto_auth_service_pb.RefreshBucketeerTokenRequest,
    callback: (error: ServiceError|null, responseMessage: proto_auth_service_pb.RefreshBucketeerTokenResponse|null) => void
  ): UnaryResponse;
}
