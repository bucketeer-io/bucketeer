// package: bucketeer.account
// file: proto/account/service.proto

import * as proto_account_service_pb from "../../proto/account/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type AccountServiceGetMeV2 = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetMeV2Request;
  readonly responseType: typeof proto_account_service_pb.GetMeV2Response;
};

type AccountServiceGetMeByEmailV2 = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetMeByEmailV2Request;
  readonly responseType: typeof proto_account_service_pb.GetMeV2Response;
};

type AccountServiceCreateAdminAccount = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.CreateAdminAccountRequest;
  readonly responseType: typeof proto_account_service_pb.CreateAdminAccountResponse;
};

type AccountServiceEnableAdminAccount = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.EnableAdminAccountRequest;
  readonly responseType: typeof proto_account_service_pb.EnableAdminAccountResponse;
};

type AccountServiceDisableAdminAccount = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.DisableAdminAccountRequest;
  readonly responseType: typeof proto_account_service_pb.DisableAdminAccountResponse;
};

type AccountServiceGetAdminAccount = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetAdminAccountRequest;
  readonly responseType: typeof proto_account_service_pb.GetAdminAccountResponse;
};

type AccountServiceListAdminAccounts = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.ListAdminAccountsRequest;
  readonly responseType: typeof proto_account_service_pb.ListAdminAccountsResponse;
};

type AccountServiceConvertAccount = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.ConvertAccountRequest;
  readonly responseType: typeof proto_account_service_pb.ConvertAccountResponse;
};

type AccountServiceCreateAccount = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.CreateAccountRequest;
  readonly responseType: typeof proto_account_service_pb.CreateAccountResponse;
};

type AccountServiceEnableAccount = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.EnableAccountRequest;
  readonly responseType: typeof proto_account_service_pb.EnableAccountResponse;
};

type AccountServiceDisableAccount = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.DisableAccountRequest;
  readonly responseType: typeof proto_account_service_pb.DisableAccountResponse;
};

type AccountServiceChangeAccountRole = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.ChangeAccountRoleRequest;
  readonly responseType: typeof proto_account_service_pb.ChangeAccountRoleResponse;
};

type AccountServiceGetAccount = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetAccountRequest;
  readonly responseType: typeof proto_account_service_pb.GetAccountResponse;
};

type AccountServiceListAccounts = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.ListAccountsRequest;
  readonly responseType: typeof proto_account_service_pb.ListAccountsResponse;
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

type AccountServiceGetAPIKeyBySearchingAllEnvironments = {
  readonly methodName: string;
  readonly service: typeof AccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_account_service_pb.GetAPIKeyBySearchingAllEnvironmentsRequest;
  readonly responseType: typeof proto_account_service_pb.GetAPIKeyBySearchingAllEnvironmentsResponse;
};

export class AccountService {
  static readonly serviceName: string;
  static readonly GetMeV2: AccountServiceGetMeV2;
  static readonly GetMeByEmailV2: AccountServiceGetMeByEmailV2;
  static readonly CreateAdminAccount: AccountServiceCreateAdminAccount;
  static readonly EnableAdminAccount: AccountServiceEnableAdminAccount;
  static readonly DisableAdminAccount: AccountServiceDisableAdminAccount;
  static readonly GetAdminAccount: AccountServiceGetAdminAccount;
  static readonly ListAdminAccounts: AccountServiceListAdminAccounts;
  static readonly ConvertAccount: AccountServiceConvertAccount;
  static readonly CreateAccount: AccountServiceCreateAccount;
  static readonly EnableAccount: AccountServiceEnableAccount;
  static readonly DisableAccount: AccountServiceDisableAccount;
  static readonly ChangeAccountRole: AccountServiceChangeAccountRole;
  static readonly GetAccount: AccountServiceGetAccount;
  static readonly ListAccounts: AccountServiceListAccounts;
  static readonly CreateAPIKey: AccountServiceCreateAPIKey;
  static readonly ChangeAPIKeyName: AccountServiceChangeAPIKeyName;
  static readonly EnableAPIKey: AccountServiceEnableAPIKey;
  static readonly DisableAPIKey: AccountServiceDisableAPIKey;
  static readonly GetAPIKey: AccountServiceGetAPIKey;
  static readonly ListAPIKeys: AccountServiceListAPIKeys;
  static readonly GetAPIKeyBySearchingAllEnvironments: AccountServiceGetAPIKeyBySearchingAllEnvironments;
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

export class AccountServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getMeV2(
    requestMessage: proto_account_service_pb.GetMeV2Request,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetMeV2Response|null) => void
  ): UnaryResponse;
  getMeV2(
    requestMessage: proto_account_service_pb.GetMeV2Request,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetMeV2Response|null) => void
  ): UnaryResponse;
  getMeByEmailV2(
    requestMessage: proto_account_service_pb.GetMeByEmailV2Request,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetMeV2Response|null) => void
  ): UnaryResponse;
  getMeByEmailV2(
    requestMessage: proto_account_service_pb.GetMeByEmailV2Request,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetMeV2Response|null) => void
  ): UnaryResponse;
  createAdminAccount(
    requestMessage: proto_account_service_pb.CreateAdminAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.CreateAdminAccountResponse|null) => void
  ): UnaryResponse;
  createAdminAccount(
    requestMessage: proto_account_service_pb.CreateAdminAccountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.CreateAdminAccountResponse|null) => void
  ): UnaryResponse;
  enableAdminAccount(
    requestMessage: proto_account_service_pb.EnableAdminAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.EnableAdminAccountResponse|null) => void
  ): UnaryResponse;
  enableAdminAccount(
    requestMessage: proto_account_service_pb.EnableAdminAccountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.EnableAdminAccountResponse|null) => void
  ): UnaryResponse;
  disableAdminAccount(
    requestMessage: proto_account_service_pb.DisableAdminAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.DisableAdminAccountResponse|null) => void
  ): UnaryResponse;
  disableAdminAccount(
    requestMessage: proto_account_service_pb.DisableAdminAccountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.DisableAdminAccountResponse|null) => void
  ): UnaryResponse;
  getAdminAccount(
    requestMessage: proto_account_service_pb.GetAdminAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetAdminAccountResponse|null) => void
  ): UnaryResponse;
  getAdminAccount(
    requestMessage: proto_account_service_pb.GetAdminAccountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetAdminAccountResponse|null) => void
  ): UnaryResponse;
  listAdminAccounts(
    requestMessage: proto_account_service_pb.ListAdminAccountsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ListAdminAccountsResponse|null) => void
  ): UnaryResponse;
  listAdminAccounts(
    requestMessage: proto_account_service_pb.ListAdminAccountsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ListAdminAccountsResponse|null) => void
  ): UnaryResponse;
  convertAccount(
    requestMessage: proto_account_service_pb.ConvertAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ConvertAccountResponse|null) => void
  ): UnaryResponse;
  convertAccount(
    requestMessage: proto_account_service_pb.ConvertAccountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ConvertAccountResponse|null) => void
  ): UnaryResponse;
  createAccount(
    requestMessage: proto_account_service_pb.CreateAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.CreateAccountResponse|null) => void
  ): UnaryResponse;
  createAccount(
    requestMessage: proto_account_service_pb.CreateAccountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.CreateAccountResponse|null) => void
  ): UnaryResponse;
  enableAccount(
    requestMessage: proto_account_service_pb.EnableAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.EnableAccountResponse|null) => void
  ): UnaryResponse;
  enableAccount(
    requestMessage: proto_account_service_pb.EnableAccountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.EnableAccountResponse|null) => void
  ): UnaryResponse;
  disableAccount(
    requestMessage: proto_account_service_pb.DisableAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.DisableAccountResponse|null) => void
  ): UnaryResponse;
  disableAccount(
    requestMessage: proto_account_service_pb.DisableAccountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.DisableAccountResponse|null) => void
  ): UnaryResponse;
  changeAccountRole(
    requestMessage: proto_account_service_pb.ChangeAccountRoleRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ChangeAccountRoleResponse|null) => void
  ): UnaryResponse;
  changeAccountRole(
    requestMessage: proto_account_service_pb.ChangeAccountRoleRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ChangeAccountRoleResponse|null) => void
  ): UnaryResponse;
  getAccount(
    requestMessage: proto_account_service_pb.GetAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetAccountResponse|null) => void
  ): UnaryResponse;
  getAccount(
    requestMessage: proto_account_service_pb.GetAccountRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetAccountResponse|null) => void
  ): UnaryResponse;
  listAccounts(
    requestMessage: proto_account_service_pb.ListAccountsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ListAccountsResponse|null) => void
  ): UnaryResponse;
  listAccounts(
    requestMessage: proto_account_service_pb.ListAccountsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ListAccountsResponse|null) => void
  ): UnaryResponse;
  createAPIKey(
    requestMessage: proto_account_service_pb.CreateAPIKeyRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.CreateAPIKeyResponse|null) => void
  ): UnaryResponse;
  createAPIKey(
    requestMessage: proto_account_service_pb.CreateAPIKeyRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.CreateAPIKeyResponse|null) => void
  ): UnaryResponse;
  changeAPIKeyName(
    requestMessage: proto_account_service_pb.ChangeAPIKeyNameRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ChangeAPIKeyNameResponse|null) => void
  ): UnaryResponse;
  changeAPIKeyName(
    requestMessage: proto_account_service_pb.ChangeAPIKeyNameRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ChangeAPIKeyNameResponse|null) => void
  ): UnaryResponse;
  enableAPIKey(
    requestMessage: proto_account_service_pb.EnableAPIKeyRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.EnableAPIKeyResponse|null) => void
  ): UnaryResponse;
  enableAPIKey(
    requestMessage: proto_account_service_pb.EnableAPIKeyRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.EnableAPIKeyResponse|null) => void
  ): UnaryResponse;
  disableAPIKey(
    requestMessage: proto_account_service_pb.DisableAPIKeyRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.DisableAPIKeyResponse|null) => void
  ): UnaryResponse;
  disableAPIKey(
    requestMessage: proto_account_service_pb.DisableAPIKeyRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.DisableAPIKeyResponse|null) => void
  ): UnaryResponse;
  getAPIKey(
    requestMessage: proto_account_service_pb.GetAPIKeyRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetAPIKeyResponse|null) => void
  ): UnaryResponse;
  getAPIKey(
    requestMessage: proto_account_service_pb.GetAPIKeyRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetAPIKeyResponse|null) => void
  ): UnaryResponse;
  listAPIKeys(
    requestMessage: proto_account_service_pb.ListAPIKeysRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ListAPIKeysResponse|null) => void
  ): UnaryResponse;
  listAPIKeys(
    requestMessage: proto_account_service_pb.ListAPIKeysRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.ListAPIKeysResponse|null) => void
  ): UnaryResponse;
  getAPIKeyBySearchingAllEnvironments(
    requestMessage: proto_account_service_pb.GetAPIKeyBySearchingAllEnvironmentsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetAPIKeyBySearchingAllEnvironmentsResponse|null) => void
  ): UnaryResponse;
  getAPIKeyBySearchingAllEnvironments(
    requestMessage: proto_account_service_pb.GetAPIKeyBySearchingAllEnvironmentsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_account_service_pb.GetAPIKeyBySearchingAllEnvironmentsResponse|null) => void
  ): UnaryResponse;
}

