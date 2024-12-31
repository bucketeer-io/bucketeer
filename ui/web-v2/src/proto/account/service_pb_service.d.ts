// package: bucketeer.account
// file: proto/account/service.proto

import * as proto_account_service_pb from '../../proto/account/service_pb';
import { grpc } from '@improbable-eng/grpc-web';

type AccountServiceGetMe = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetMeRequest;
  readonly responseType: typeof proto_account_service_pb.GetMeResponse;
};

type AccountServiceGetMyOrganizations = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetMyOrganizationsRequest;
  readonly responseType: typeof proto_account_service_pb.GetMyOrganizationsResponse;
};

type AccountServiceGetMyOrganizationsByEmail = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetMyOrganizationsByEmailRequest;
  readonly responseType: typeof proto_account_service_pb.GetMyOrganizationsResponse;
};

type AccountServiceCreateAccountV2 = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.CreateAccountV2Request;
  readonly responseType: typeof proto_account_service_pb.CreateAccountV2Response;
};

type AccountServiceEnableAccountV2 = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.EnableAccountV2Request;
  readonly responseType: typeof proto_account_service_pb.EnableAccountV2Response;
};

type AccountServiceDisableAccountV2 = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.DisableAccountV2Request;
  readonly responseType: typeof proto_account_service_pb.DisableAccountV2Response;
};

type AccountServiceUpdateAccountV2 = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.UpdateAccountV2Request;
  readonly responseType: typeof proto_account_service_pb.UpdateAccountV2Response;
};

type AccountServiceDeleteAccountV2 = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.DeleteAccountV2Request;
  readonly responseType: typeof proto_account_service_pb.DeleteAccountV2Response;
};

type AccountServiceGetAccountV2 = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetAccountV2Request;
  readonly responseType: typeof proto_account_service_pb.GetAccountV2Response;
};

type AccountServiceGetAccountV2ByEnvironmentID = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetAccountV2ByEnvironmentIDRequest;
  readonly responseType: typeof proto_account_service_pb.GetAccountV2ByEnvironmentIDResponse;
};

type AccountServiceListAccountsV2 = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.ListAccountsV2Request;
  readonly responseType: typeof proto_account_service_pb.ListAccountsV2Response;
};

type AccountServiceCreateAPIKey = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.CreateAPIKeyRequest;
  readonly responseType: typeof proto_account_service_pb.CreateAPIKeyResponse;
};

type AccountServiceChangeAPIKeyName = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.ChangeAPIKeyNameRequest;
  readonly responseType: typeof proto_account_service_pb.ChangeAPIKeyNameResponse;
};

type AccountServiceEnableAPIKey = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.EnableAPIKeyRequest;
  readonly responseType: typeof proto_account_service_pb.EnableAPIKeyResponse;
};

type AccountServiceDisableAPIKey = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.DisableAPIKeyRequest;
  readonly responseType: typeof proto_account_service_pb.DisableAPIKeyResponse;
};

type AccountServiceGetAPIKey = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetAPIKeyRequest;
  readonly responseType: typeof proto_account_service_pb.GetAPIKeyResponse;
};

type AccountServiceListAPIKeys = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.ListAPIKeysRequest;
  readonly responseType: typeof proto_account_service_pb.ListAPIKeysResponse;
};

type AccountServiceGetEnvironmentAPIKey = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetEnvironmentAPIKeyRequest;
  readonly responseType: typeof proto_account_service_pb.GetEnvironmentAPIKeyResponse;
};

type AccountServiceCreateSearchFilter = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.CreateSearchFilterRequest;
  readonly responseType: typeof proto_account_service_pb.CreateSearchFilterResponse;
};

type AccountServiceUpdateSearchFilter = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.UpdateSearchFilterRequest;
  readonly responseType: typeof proto_account_service_pb.UpdateSearchFilterResponse;
};

type AccountServiceDeleteSearchFilter = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.DeleteSearchFilterRequest;
  readonly responseType: typeof proto_account_service_pb.DeleteSearchFilterResponse;
};

type AccountServiceUpdateAPIKey = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.UpdateAPIKeyRequest;
  readonly responseType: typeof proto_account_service_pb.UpdateAPIKeyResponse;
};

export class AccountService {
  static readonly serviceName: string;
  static readonly GetMe: AccountServiceGetMe;
  static readonly GetMyOrganizations: AccountServiceGetMyOrganizations;
  static readonly GetMyOrganizationsByEmail: AccountServiceGetMyOrganizationsByEmail;
  static readonly CreateAccountV2: AccountServiceCreateAccountV2;
  static readonly EnableAccountV2: AccountServiceEnableAccountV2;
  static readonly DisableAccountV2: AccountServiceDisableAccountV2;
  static readonly UpdateAccountV2: AccountServiceUpdateAccountV2;
  static readonly DeleteAccountV2: AccountServiceDeleteAccountV2;
  static readonly GetAccountV2: AccountServiceGetAccountV2;
  static readonly GetAccountV2ByEnvironmentID: AccountServiceGetAccountV2ByEnvironmentID;
  static readonly ListAccountsV2: AccountServiceListAccountsV2;
  static readonly CreateAPIKey: AccountServiceCreateAPIKey;
  static readonly ChangeAPIKeyName: AccountServiceChangeAPIKeyName;
  static readonly EnableAPIKey: AccountServiceEnableAPIKey;
  static readonly DisableAPIKey: AccountServiceDisableAPIKey;
  static readonly GetAPIKey: AccountServiceGetAPIKey;
  static readonly ListAPIKeys: AccountServiceListAPIKeys;
  static readonly GetEnvironmentAPIKey: AccountServiceGetEnvironmentAPIKey;
  static readonly CreateSearchFilter: AccountServiceCreateSearchFilter;
  static readonly UpdateSearchFilter: AccountServiceUpdateSearchFilter;
  static readonly DeleteSearchFilter: AccountServiceDeleteSearchFilter;
  static readonly UpdateAPIKey: AccountServiceUpdateAPIKey;
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

export class AccountServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getMe(
    requestMessage: proto_account_service_pb.GetMeRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetMeResponse | null
    ) => void
  ): UnaryResponse;
  getMe(
    requestMessage: proto_account_service_pb.GetMeRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetMeResponse | null
    ) => void
  ): UnaryResponse;
  getMyOrganizations(
    requestMessage: proto_account_service_pb.GetMyOrganizationsRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetMyOrganizationsResponse | null
    ) => void
  ): UnaryResponse;
  getMyOrganizations(
    requestMessage: proto_account_service_pb.GetMyOrganizationsRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetMyOrganizationsResponse | null
    ) => void
  ): UnaryResponse;
  getMyOrganizationsByEmail(
    requestMessage: proto_account_service_pb.GetMyOrganizationsByEmailRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetMyOrganizationsResponse | null
    ) => void
  ): UnaryResponse;
  getMyOrganizationsByEmail(
    requestMessage: proto_account_service_pb.GetMyOrganizationsByEmailRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetMyOrganizationsResponse | null
    ) => void
  ): UnaryResponse;
  createAccountV2(
    requestMessage: proto_account_service_pb.CreateAccountV2Request,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.CreateAccountV2Response | null
    ) => void
  ): UnaryResponse;
  createAccountV2(
    requestMessage: proto_account_service_pb.CreateAccountV2Request,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.CreateAccountV2Response | null
    ) => void
  ): UnaryResponse;
  enableAccountV2(
    requestMessage: proto_account_service_pb.EnableAccountV2Request,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.EnableAccountV2Response | null
    ) => void
  ): UnaryResponse;
  enableAccountV2(
    requestMessage: proto_account_service_pb.EnableAccountV2Request,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.EnableAccountV2Response | null
    ) => void
  ): UnaryResponse;
  disableAccountV2(
    requestMessage: proto_account_service_pb.DisableAccountV2Request,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.DisableAccountV2Response | null
    ) => void
  ): UnaryResponse;
  disableAccountV2(
    requestMessage: proto_account_service_pb.DisableAccountV2Request,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.DisableAccountV2Response | null
    ) => void
  ): UnaryResponse;
  updateAccountV2(
    requestMessage: proto_account_service_pb.UpdateAccountV2Request,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.UpdateAccountV2Response | null
    ) => void
  ): UnaryResponse;
  updateAccountV2(
    requestMessage: proto_account_service_pb.UpdateAccountV2Request,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.UpdateAccountV2Response | null
    ) => void
  ): UnaryResponse;
  deleteAccountV2(
    requestMessage: proto_account_service_pb.DeleteAccountV2Request,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.DeleteAccountV2Response | null
    ) => void
  ): UnaryResponse;
  deleteAccountV2(
    requestMessage: proto_account_service_pb.DeleteAccountV2Request,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.DeleteAccountV2Response | null
    ) => void
  ): UnaryResponse;
  getAccountV2(
    requestMessage: proto_account_service_pb.GetAccountV2Request,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetAccountV2Response | null
    ) => void
  ): UnaryResponse;
  getAccountV2(
    requestMessage: proto_account_service_pb.GetAccountV2Request,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetAccountV2Response | null
    ) => void
  ): UnaryResponse;
  getAccountV2ByEnvironmentID(
    requestMessage: proto_account_service_pb.GetAccountV2ByEnvironmentIDRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetAccountV2ByEnvironmentIDResponse | null
    ) => void
  ): UnaryResponse;
  getAccountV2ByEnvironmentID(
    requestMessage: proto_account_service_pb.GetAccountV2ByEnvironmentIDRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetAccountV2ByEnvironmentIDResponse | null
    ) => void
  ): UnaryResponse;
  listAccountsV2(
    requestMessage: proto_account_service_pb.ListAccountsV2Request,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.ListAccountsV2Response | null
    ) => void
  ): UnaryResponse;
  listAccountsV2(
    requestMessage: proto_account_service_pb.ListAccountsV2Request,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.ListAccountsV2Response | null
    ) => void
  ): UnaryResponse;
  createAPIKey(
    requestMessage: proto_account_service_pb.CreateAPIKeyRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.CreateAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  createAPIKey(
    requestMessage: proto_account_service_pb.CreateAPIKeyRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.CreateAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  changeAPIKeyName(
    requestMessage: proto_account_service_pb.ChangeAPIKeyNameRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.ChangeAPIKeyNameResponse | null
    ) => void
  ): UnaryResponse;
  changeAPIKeyName(
    requestMessage: proto_account_service_pb.ChangeAPIKeyNameRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.ChangeAPIKeyNameResponse | null
    ) => void
  ): UnaryResponse;
  enableAPIKey(
    requestMessage: proto_account_service_pb.EnableAPIKeyRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.EnableAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  enableAPIKey(
    requestMessage: proto_account_service_pb.EnableAPIKeyRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.EnableAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  disableAPIKey(
    requestMessage: proto_account_service_pb.DisableAPIKeyRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.DisableAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  disableAPIKey(
    requestMessage: proto_account_service_pb.DisableAPIKeyRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.DisableAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  getAPIKey(
    requestMessage: proto_account_service_pb.GetAPIKeyRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  getAPIKey(
    requestMessage: proto_account_service_pb.GetAPIKeyRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  listAPIKeys(
    requestMessage: proto_account_service_pb.ListAPIKeysRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.ListAPIKeysResponse | null
    ) => void
  ): UnaryResponse;
  listAPIKeys(
    requestMessage: proto_account_service_pb.ListAPIKeysRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.ListAPIKeysResponse | null
    ) => void
  ): UnaryResponse;
  getEnvironmentAPIKey(
    requestMessage: proto_account_service_pb.GetEnvironmentAPIKeyRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetEnvironmentAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  getEnvironmentAPIKey(
    requestMessage: proto_account_service_pb.GetEnvironmentAPIKeyRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.GetEnvironmentAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  createSearchFilter(
    requestMessage: proto_account_service_pb.CreateSearchFilterRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.CreateSearchFilterResponse | null
    ) => void
  ): UnaryResponse;
  createSearchFilter(
    requestMessage: proto_account_service_pb.CreateSearchFilterRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.CreateSearchFilterResponse | null
    ) => void
  ): UnaryResponse;
  updateSearchFilter(
    requestMessage: proto_account_service_pb.UpdateSearchFilterRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.UpdateSearchFilterResponse | null
    ) => void
  ): UnaryResponse;
  updateSearchFilter(
    requestMessage: proto_account_service_pb.UpdateSearchFilterRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.UpdateSearchFilterResponse | null
    ) => void
  ): UnaryResponse;
  deleteSearchFilter(
    requestMessage: proto_account_service_pb.DeleteSearchFilterRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.DeleteSearchFilterResponse | null
    ) => void
  ): UnaryResponse;
  deleteSearchFilter(
    requestMessage: proto_account_service_pb.DeleteSearchFilterRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.DeleteSearchFilterResponse | null
    ) => void
  ): UnaryResponse;
  updateAPIKey(
    requestMessage: proto_account_service_pb.UpdateAPIKeyRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.UpdateAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
  updateAPIKey(
    requestMessage: proto_account_service_pb.UpdateAPIKeyRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_account_service_pb.UpdateAPIKeyResponse | null
    ) => void
  ): UnaryResponse;
}
