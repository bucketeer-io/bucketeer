// package: bucketeer.auth
// file: proto/auth/service.proto

import * as proto_auth_service_pb from '../../proto/auth/service_pb';
import { grpc } from '@improbable-eng/grpc-web';

type AuthServiceExchangeToken = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.ExchangeTokenRequest;
  readonly responseType: typeof proto_auth_service_pb.ExchangeTokenResponse;
};

type AuthServiceGetAuthenticationURL = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.GetAuthenticationURLRequest;
  readonly responseType: typeof proto_auth_service_pb.GetAuthenticationURLResponse;
};

type AuthServiceRefreshToken = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.RefreshTokenRequest;
  readonly responseType: typeof proto_auth_service_pb.RefreshTokenResponse;
};

type AuthServiceSignIn = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.SignInRequest;
  readonly responseType: typeof proto_auth_service_pb.SignInResponse;
};

type AuthServiceSwitchOrganization = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.SwitchOrganizationRequest;
  readonly responseType: typeof proto_auth_service_pb.SwitchOrganizationResponse;
};

type AuthServiceGetDeploymentStatus = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.GetDeploymentStatusRequest;
  readonly responseType: typeof proto_auth_service_pb.GetDeploymentStatusResponse;
};

export class AuthService {
  static readonly serviceName: string;
  static readonly ExchangeToken: AuthServiceExchangeToken;
  static readonly GetAuthenticationURL: AuthServiceGetAuthenticationURL;
  static readonly RefreshToken: AuthServiceRefreshToken;
  static readonly SignIn: AuthServiceSignIn;
  static readonly SwitchOrganization: AuthServiceSwitchOrganization;
  static readonly GetDeploymentStatus: AuthServiceGetDeploymentStatus;
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
  getAuthenticationURL(
    requestMessage: proto_auth_service_pb.GetAuthenticationURLRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.GetAuthenticationURLResponse | null
    ) => void
  ): UnaryResponse;
  getAuthenticationURL(
    requestMessage: proto_auth_service_pb.GetAuthenticationURLRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.GetAuthenticationURLResponse | null
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
  signIn(
    requestMessage: proto_auth_service_pb.SignInRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.SignInResponse | null
    ) => void
  ): UnaryResponse;
  signIn(
    requestMessage: proto_auth_service_pb.SignInRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.SignInResponse | null
    ) => void
  ): UnaryResponse;
  switchOrganization(
    requestMessage: proto_auth_service_pb.SwitchOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.SwitchOrganizationResponse | null
    ) => void
  ): UnaryResponse;
  switchOrganization(
    requestMessage: proto_auth_service_pb.SwitchOrganizationRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.SwitchOrganizationResponse | null
    ) => void
  ): UnaryResponse;
  getDeploymentStatus(
    requestMessage: proto_auth_service_pb.GetDeploymentStatusRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.GetDeploymentStatusResponse | null
    ) => void
  ): UnaryResponse;
  getDeploymentStatus(
    requestMessage: proto_auth_service_pb.GetDeploymentStatusRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.GetDeploymentStatusResponse | null
    ) => void
  ): UnaryResponse;
}
