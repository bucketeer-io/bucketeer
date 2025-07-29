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

type AuthServiceGetDemoSiteStatus = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.GetDemoSiteStatusRequest;
  readonly responseType: typeof proto_auth_service_pb.GetDemoSiteStatusResponse;
};

type AuthServiceCreatePassword = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.CreatePasswordRequest;
  readonly responseType: typeof proto_auth_service_pb.CreatePasswordResponse;
};

type AuthServiceUpdatePassword = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.UpdatePasswordRequest;
  readonly responseType: typeof proto_auth_service_pb.UpdatePasswordResponse;
};

type AuthServiceInitiatePasswordReset = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.InitiatePasswordResetRequest;
  readonly responseType: typeof proto_auth_service_pb.InitiatePasswordResetResponse;
};

type AuthServiceResetPassword = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.ResetPasswordRequest;
  readonly responseType: typeof proto_auth_service_pb.ResetPasswordResponse;
};

type AuthServiceValidatePasswordResetToken = {
  readonly methodName: string;
  readonly service: typeof AuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_auth_service_pb.ValidatePasswordResetTokenRequest;
  readonly responseType: typeof proto_auth_service_pb.ValidatePasswordResetTokenResponse;
};

export class AuthService {
  static readonly serviceName: string;
  static readonly ExchangeToken: AuthServiceExchangeToken;
  static readonly GetAuthenticationURL: AuthServiceGetAuthenticationURL;
  static readonly RefreshToken: AuthServiceRefreshToken;
  static readonly SignIn: AuthServiceSignIn;
  static readonly SwitchOrganization: AuthServiceSwitchOrganization;
  static readonly GetDemoSiteStatus: AuthServiceGetDemoSiteStatus;
  static readonly CreatePassword: AuthServiceCreatePassword;
  static readonly UpdatePassword: AuthServiceUpdatePassword;
  static readonly InitiatePasswordReset: AuthServiceInitiatePasswordReset;
  static readonly ResetPassword: AuthServiceResetPassword;
  static readonly ValidatePasswordResetToken: AuthServiceValidatePasswordResetToken;
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
  getDemoSiteStatus(
    requestMessage: proto_auth_service_pb.GetDemoSiteStatusRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.GetDemoSiteStatusResponse | null
    ) => void
  ): UnaryResponse;
  getDemoSiteStatus(
    requestMessage: proto_auth_service_pb.GetDemoSiteStatusRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.GetDemoSiteStatusResponse | null
    ) => void
  ): UnaryResponse;
  createPassword(
    requestMessage: proto_auth_service_pb.CreatePasswordRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.CreatePasswordResponse | null
    ) => void
  ): UnaryResponse;
  createPassword(
    requestMessage: proto_auth_service_pb.CreatePasswordRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.CreatePasswordResponse | null
    ) => void
  ): UnaryResponse;
  updatePassword(
    requestMessage: proto_auth_service_pb.UpdatePasswordRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.UpdatePasswordResponse | null
    ) => void
  ): UnaryResponse;
  updatePassword(
    requestMessage: proto_auth_service_pb.UpdatePasswordRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.UpdatePasswordResponse | null
    ) => void
  ): UnaryResponse;
  initiatePasswordReset(
    requestMessage: proto_auth_service_pb.InitiatePasswordResetRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.InitiatePasswordResetResponse | null
    ) => void
  ): UnaryResponse;
  initiatePasswordReset(
    requestMessage: proto_auth_service_pb.InitiatePasswordResetRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.InitiatePasswordResetResponse | null
    ) => void
  ): UnaryResponse;
  resetPassword(
    requestMessage: proto_auth_service_pb.ResetPasswordRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.ResetPasswordResponse | null
    ) => void
  ): UnaryResponse;
  resetPassword(
    requestMessage: proto_auth_service_pb.ResetPasswordRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.ResetPasswordResponse | null
    ) => void
  ): UnaryResponse;
  validatePasswordResetToken(
    requestMessage: proto_auth_service_pb.ValidatePasswordResetTokenRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.ValidatePasswordResetTokenResponse | null
    ) => void
  ): UnaryResponse;
  validatePasswordResetToken(
    requestMessage: proto_auth_service_pb.ValidatePasswordResetTokenRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_auth_service_pb.ValidatePasswordResetTokenResponse | null
    ) => void
  ): UnaryResponse;
}
