// package: bucketeer.account
// file: proto/account/service.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as proto_account_account_pb from "../../proto/account/account_pb";
import * as proto_account_api_key_pb from "../../proto/account/api_key_pb";
import * as proto_account_command_pb from "../../proto/account/command_pb";
import * as proto_environment_organization_pb from "../../proto/environment/organization_pb";

export class GetMeRequest extends jspb.Message {
  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMeRequest): GetMeRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMeRequest;
  static deserializeBinaryFromReader(message: GetMeRequest, reader: jspb.BinaryReader): GetMeRequest;
}

export namespace GetMeRequest {
  export type AsObject = {
    organizationId: string,
  }
}

export class GetMeResponse extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.ConsoleAccount | undefined;
  setAccount(value?: proto_account_account_pb.ConsoleAccount): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMeResponse): GetMeResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMeResponse;
  static deserializeBinaryFromReader(message: GetMeResponse, reader: jspb.BinaryReader): GetMeResponse;
}

export namespace GetMeResponse {
  export type AsObject = {
    account?: proto_account_account_pb.ConsoleAccount.AsObject,
  }
}

export class GetMyOrganizationsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyOrganizationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyOrganizationsRequest): GetMyOrganizationsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMyOrganizationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyOrganizationsRequest;
  static deserializeBinaryFromReader(message: GetMyOrganizationsRequest, reader: jspb.BinaryReader): GetMyOrganizationsRequest;
}

export namespace GetMyOrganizationsRequest {
  export type AsObject = {
  }
}

export class GetMyOrganizationsResponse extends jspb.Message {
  clearOrganizationsList(): void;
  getOrganizationsList(): Array<proto_environment_organization_pb.Organization>;
  setOrganizationsList(value: Array<proto_environment_organization_pb.Organization>): void;
  addOrganizations(value?: proto_environment_organization_pb.Organization, index?: number): proto_environment_organization_pb.Organization;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyOrganizationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyOrganizationsResponse): GetMyOrganizationsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMyOrganizationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyOrganizationsResponse;
  static deserializeBinaryFromReader(message: GetMyOrganizationsResponse, reader: jspb.BinaryReader): GetMyOrganizationsResponse;
}

export namespace GetMyOrganizationsResponse {
  export type AsObject = {
    organizationsList: Array<proto_environment_organization_pb.Organization.AsObject>,
  }
}

export class GetMeV2Request extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMeV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: GetMeV2Request): GetMeV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMeV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMeV2Request;
  static deserializeBinaryFromReader(message: GetMeV2Request, reader: jspb.BinaryReader): GetMeV2Request;
}

export namespace GetMeV2Request {
  export type AsObject = {
  }
}

export class GetMeByEmailV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMeByEmailV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: GetMeByEmailV2Request): GetMeByEmailV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMeByEmailV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMeByEmailV2Request;
  static deserializeBinaryFromReader(message: GetMeByEmailV2Request, reader: jspb.BinaryReader): GetMeByEmailV2Request;
}

export namespace GetMeByEmailV2Request {
  export type AsObject = {
    email: string,
  }
}

export class GetMeV2Response extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getIsAdmin(): boolean;
  setIsAdmin(value: boolean): void;

  clearEnvironmentRolesList(): void;
  getEnvironmentRolesList(): Array<proto_account_account_pb.EnvironmentRoleV2>;
  setEnvironmentRolesList(value: Array<proto_account_account_pb.EnvironmentRoleV2>): void;
  addEnvironmentRoles(value?: proto_account_account_pb.EnvironmentRoleV2, index?: number): proto_account_account_pb.EnvironmentRoleV2;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMeV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: GetMeV2Response): GetMeV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMeV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMeV2Response;
  static deserializeBinaryFromReader(message: GetMeV2Response, reader: jspb.BinaryReader): GetMeV2Response;
}

export namespace GetMeV2Response {
  export type AsObject = {
    email: string,
    isAdmin: boolean,
    environmentRolesList: Array<proto_account_account_pb.EnvironmentRoleV2.AsObject>,
  }
}

export class CreateAdminAccountRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.CreateAdminAccountCommand | undefined;
  setCommand(value?: proto_account_command_pb.CreateAdminAccountCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAdminAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAdminAccountRequest): CreateAdminAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAdminAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAdminAccountRequest;
  static deserializeBinaryFromReader(message: CreateAdminAccountRequest, reader: jspb.BinaryReader): CreateAdminAccountRequest;
}

export namespace CreateAdminAccountRequest {
  export type AsObject = {
    command?: proto_account_command_pb.CreateAdminAccountCommand.AsObject,
  }
}

export class CreateAdminAccountResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAdminAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAdminAccountResponse): CreateAdminAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAdminAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAdminAccountResponse;
  static deserializeBinaryFromReader(message: CreateAdminAccountResponse, reader: jspb.BinaryReader): CreateAdminAccountResponse;
}

export namespace CreateAdminAccountResponse {
  export type AsObject = {
  }
}

export class EnableAdminAccountRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.EnableAdminAccountCommand | undefined;
  setCommand(value?: proto_account_command_pb.EnableAdminAccountCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAdminAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAdminAccountRequest): EnableAdminAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAdminAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAdminAccountRequest;
  static deserializeBinaryFromReader(message: EnableAdminAccountRequest, reader: jspb.BinaryReader): EnableAdminAccountRequest;
}

export namespace EnableAdminAccountRequest {
  export type AsObject = {
    id: string,
    command?: proto_account_command_pb.EnableAdminAccountCommand.AsObject,
  }
}

export class EnableAdminAccountResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAdminAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAdminAccountResponse): EnableAdminAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAdminAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAdminAccountResponse;
  static deserializeBinaryFromReader(message: EnableAdminAccountResponse, reader: jspb.BinaryReader): EnableAdminAccountResponse;
}

export namespace EnableAdminAccountResponse {
  export type AsObject = {
  }
}

export class DisableAdminAccountRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.DisableAdminAccountCommand | undefined;
  setCommand(value?: proto_account_command_pb.DisableAdminAccountCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAdminAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAdminAccountRequest): DisableAdminAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAdminAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAdminAccountRequest;
  static deserializeBinaryFromReader(message: DisableAdminAccountRequest, reader: jspb.BinaryReader): DisableAdminAccountRequest;
}

export namespace DisableAdminAccountRequest {
  export type AsObject = {
    id: string,
    command?: proto_account_command_pb.DisableAdminAccountCommand.AsObject,
  }
}

export class DisableAdminAccountResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAdminAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAdminAccountResponse): DisableAdminAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAdminAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAdminAccountResponse;
  static deserializeBinaryFromReader(message: DisableAdminAccountResponse, reader: jspb.BinaryReader): DisableAdminAccountResponse;
}

export namespace DisableAdminAccountResponse {
  export type AsObject = {
  }
}

export class GetAdminAccountRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAdminAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAdminAccountRequest): GetAdminAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAdminAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAdminAccountRequest;
  static deserializeBinaryFromReader(message: GetAdminAccountRequest, reader: jspb.BinaryReader): GetAdminAccountRequest;
}

export namespace GetAdminAccountRequest {
  export type AsObject = {
    email: string,
  }
}

export class GetAdminAccountResponse extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.Account | undefined;
  setAccount(value?: proto_account_account_pb.Account): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAdminAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAdminAccountResponse): GetAdminAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAdminAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAdminAccountResponse;
  static deserializeBinaryFromReader(message: GetAdminAccountResponse, reader: jspb.BinaryReader): GetAdminAccountResponse;
}

export namespace GetAdminAccountResponse {
  export type AsObject = {
    account?: proto_account_account_pb.Account.AsObject,
  }
}

export class ListAdminAccountsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListAdminAccountsRequest.OrderByMap[keyof ListAdminAccountsRequest.OrderByMap];
  setOrderBy(value: ListAdminAccountsRequest.OrderByMap[keyof ListAdminAccountsRequest.OrderByMap]): void;

  getOrderDirection(): ListAdminAccountsRequest.OrderDirectionMap[keyof ListAdminAccountsRequest.OrderDirectionMap];
  setOrderDirection(value: ListAdminAccountsRequest.OrderDirectionMap[keyof ListAdminAccountsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAdminAccountsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAdminAccountsRequest): ListAdminAccountsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAdminAccountsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAdminAccountsRequest;
  static deserializeBinaryFromReader(message: ListAdminAccountsRequest, reader: jspb.BinaryReader): ListAdminAccountsRequest;
}

export namespace ListAdminAccountsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    orderBy: ListAdminAccountsRequest.OrderByMap[keyof ListAdminAccountsRequest.OrderByMap],
    orderDirection: ListAdminAccountsRequest.OrderDirectionMap[keyof ListAdminAccountsRequest.OrderDirectionMap],
    searchKeyword: string,
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    EMAIL: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListAdminAccountsResponse extends jspb.Message {
  clearAccountsList(): void;
  getAccountsList(): Array<proto_account_account_pb.Account>;
  setAccountsList(value: Array<proto_account_account_pb.Account>): void;
  addAccounts(value?: proto_account_account_pb.Account, index?: number): proto_account_account_pb.Account;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAdminAccountsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAdminAccountsResponse): ListAdminAccountsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAdminAccountsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAdminAccountsResponse;
  static deserializeBinaryFromReader(message: ListAdminAccountsResponse, reader: jspb.BinaryReader): ListAdminAccountsResponse;
}

export namespace ListAdminAccountsResponse {
  export type AsObject = {
    accountsList: Array<proto_account_account_pb.Account.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class ConvertAccountRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.ConvertAccountCommand | undefined;
  setCommand(value?: proto_account_command_pb.ConvertAccountCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConvertAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConvertAccountRequest): ConvertAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConvertAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConvertAccountRequest;
  static deserializeBinaryFromReader(message: ConvertAccountRequest, reader: jspb.BinaryReader): ConvertAccountRequest;
}

export namespace ConvertAccountRequest {
  export type AsObject = {
    id: string,
    command?: proto_account_command_pb.ConvertAccountCommand.AsObject,
  }
}

export class ConvertAccountResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConvertAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ConvertAccountResponse): ConvertAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConvertAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConvertAccountResponse;
  static deserializeBinaryFromReader(message: ConvertAccountResponse, reader: jspb.BinaryReader): ConvertAccountResponse;
}

export namespace ConvertAccountResponse {
  export type AsObject = {
  }
}

export class CreateAccountRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.CreateAccountCommand | undefined;
  setCommand(value?: proto_account_command_pb.CreateAccountCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountRequest): CreateAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountRequest;
  static deserializeBinaryFromReader(message: CreateAccountRequest, reader: jspb.BinaryReader): CreateAccountRequest;
}

export namespace CreateAccountRequest {
  export type AsObject = {
    command?: proto_account_command_pb.CreateAccountCommand.AsObject,
    environmentNamespace: string,
  }
}

export class CreateAccountResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountResponse): CreateAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountResponse;
  static deserializeBinaryFromReader(message: CreateAccountResponse, reader: jspb.BinaryReader): CreateAccountResponse;
}

export namespace CreateAccountResponse {
  export type AsObject = {
  }
}

export class EnableAccountRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.EnableAccountCommand | undefined;
  setCommand(value?: proto_account_command_pb.EnableAccountCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAccountRequest): EnableAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAccountRequest;
  static deserializeBinaryFromReader(message: EnableAccountRequest, reader: jspb.BinaryReader): EnableAccountRequest;
}

export namespace EnableAccountRequest {
  export type AsObject = {
    id: string,
    command?: proto_account_command_pb.EnableAccountCommand.AsObject,
    environmentNamespace: string,
  }
}

export class EnableAccountResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAccountResponse): EnableAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAccountResponse;
  static deserializeBinaryFromReader(message: EnableAccountResponse, reader: jspb.BinaryReader): EnableAccountResponse;
}

export namespace EnableAccountResponse {
  export type AsObject = {
  }
}

export class DisableAccountRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.DisableAccountCommand | undefined;
  setCommand(value?: proto_account_command_pb.DisableAccountCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAccountRequest): DisableAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountRequest;
  static deserializeBinaryFromReader(message: DisableAccountRequest, reader: jspb.BinaryReader): DisableAccountRequest;
}

export namespace DisableAccountRequest {
  export type AsObject = {
    id: string,
    command?: proto_account_command_pb.DisableAccountCommand.AsObject,
    environmentNamespace: string,
  }
}

export class DisableAccountResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAccountResponse): DisableAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountResponse;
  static deserializeBinaryFromReader(message: DisableAccountResponse, reader: jspb.BinaryReader): DisableAccountResponse;
}

export namespace DisableAccountResponse {
  export type AsObject = {
  }
}

export class ChangeAccountRoleRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.ChangeAccountRoleCommand | undefined;
  setCommand(value?: proto_account_command_pb.ChangeAccountRoleCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAccountRoleRequest): ChangeAccountRoleRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAccountRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountRoleRequest;
  static deserializeBinaryFromReader(message: ChangeAccountRoleRequest, reader: jspb.BinaryReader): ChangeAccountRoleRequest;
}

export namespace ChangeAccountRoleRequest {
  export type AsObject = {
    id: string,
    command?: proto_account_command_pb.ChangeAccountRoleCommand.AsObject,
    environmentNamespace: string,
  }
}

export class ChangeAccountRoleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAccountRoleResponse): ChangeAccountRoleResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAccountRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountRoleResponse;
  static deserializeBinaryFromReader(message: ChangeAccountRoleResponse, reader: jspb.BinaryReader): ChangeAccountRoleResponse;
}

export namespace ChangeAccountRoleResponse {
  export type AsObject = {
  }
}

export class GetAccountRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAccountRequest): GetAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountRequest;
  static deserializeBinaryFromReader(message: GetAccountRequest, reader: jspb.BinaryReader): GetAccountRequest;
}

export namespace GetAccountRequest {
  export type AsObject = {
    email: string,
    environmentNamespace: string,
  }
}

export class GetAccountResponse extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.Account | undefined;
  setAccount(value?: proto_account_account_pb.Account): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAccountResponse): GetAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountResponse;
  static deserializeBinaryFromReader(message: GetAccountResponse, reader: jspb.BinaryReader): GetAccountResponse;
}

export namespace GetAccountResponse {
  export type AsObject = {
    account?: proto_account_account_pb.Account.AsObject,
  }
}

export class ListAccountsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getOrderBy(): ListAccountsRequest.OrderByMap[keyof ListAccountsRequest.OrderByMap];
  setOrderBy(value: ListAccountsRequest.OrderByMap[keyof ListAccountsRequest.OrderByMap]): void;

  getOrderDirection(): ListAccountsRequest.OrderDirectionMap[keyof ListAccountsRequest.OrderDirectionMap];
  setOrderDirection(value: ListAccountsRequest.OrderDirectionMap[keyof ListAccountsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasRole(): boolean;
  clearRole(): void;
  getRole(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setRole(value?: google_protobuf_wrappers_pb.Int32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAccountsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAccountsRequest): ListAccountsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAccountsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAccountsRequest;
  static deserializeBinaryFromReader(message: ListAccountsRequest, reader: jspb.BinaryReader): ListAccountsRequest;
}

export namespace ListAccountsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    environmentNamespace: string,
    orderBy: ListAccountsRequest.OrderByMap[keyof ListAccountsRequest.OrderByMap],
    orderDirection: ListAccountsRequest.OrderDirectionMap[keyof ListAccountsRequest.OrderDirectionMap],
    searchKeyword: string,
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    role?: google_protobuf_wrappers_pb.Int32Value.AsObject,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    EMAIL: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListAccountsResponse extends jspb.Message {
  clearAccountsList(): void;
  getAccountsList(): Array<proto_account_account_pb.Account>;
  setAccountsList(value: Array<proto_account_account_pb.Account>): void;
  addAccounts(value?: proto_account_account_pb.Account, index?: number): proto_account_account_pb.Account;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAccountsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAccountsResponse): ListAccountsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAccountsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAccountsResponse;
  static deserializeBinaryFromReader(message: ListAccountsResponse, reader: jspb.BinaryReader): ListAccountsResponse;
}

export namespace ListAccountsResponse {
  export type AsObject = {
    accountsList: Array<proto_account_account_pb.Account.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class CreateAccountV2Request extends jspb.Message {
  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.CreateAccountV2Command | undefined;
  setCommand(value?: proto_account_command_pb.CreateAccountV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountV2Request): CreateAccountV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAccountV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountV2Request;
  static deserializeBinaryFromReader(message: CreateAccountV2Request, reader: jspb.BinaryReader): CreateAccountV2Request;
}

export namespace CreateAccountV2Request {
  export type AsObject = {
    organizationId: string,
    command?: proto_account_command_pb.CreateAccountV2Command.AsObject,
  }
}

export class CreateAccountV2Response extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.AccountV2 | undefined;
  setAccount(value?: proto_account_account_pb.AccountV2): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountV2Response): CreateAccountV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAccountV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountV2Response;
  static deserializeBinaryFromReader(message: CreateAccountV2Response, reader: jspb.BinaryReader): CreateAccountV2Response;
}

export namespace CreateAccountV2Response {
  export type AsObject = {
    account?: proto_account_account_pb.AccountV2.AsObject,
  }
}

export class EnableAccountV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.EnableAccountV2Command | undefined;
  setCommand(value?: proto_account_command_pb.EnableAccountV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAccountV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAccountV2Request): EnableAccountV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAccountV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAccountV2Request;
  static deserializeBinaryFromReader(message: EnableAccountV2Request, reader: jspb.BinaryReader): EnableAccountV2Request;
}

export namespace EnableAccountV2Request {
  export type AsObject = {
    email: string,
    organizationId: string,
    command?: proto_account_command_pb.EnableAccountV2Command.AsObject,
  }
}

export class EnableAccountV2Response extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAccountV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAccountV2Response): EnableAccountV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAccountV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAccountV2Response;
  static deserializeBinaryFromReader(message: EnableAccountV2Response, reader: jspb.BinaryReader): EnableAccountV2Response;
}

export namespace EnableAccountV2Response {
  export type AsObject = {
  }
}

export class DisableAccountV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.DisableAccountV2Command | undefined;
  setCommand(value?: proto_account_command_pb.DisableAccountV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAccountV2Request): DisableAccountV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAccountV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountV2Request;
  static deserializeBinaryFromReader(message: DisableAccountV2Request, reader: jspb.BinaryReader): DisableAccountV2Request;
}

export namespace DisableAccountV2Request {
  export type AsObject = {
    email: string,
    organizationId: string,
    command?: proto_account_command_pb.DisableAccountV2Command.AsObject,
  }
}

export class DisableAccountV2Response extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAccountV2Response): DisableAccountV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAccountV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountV2Response;
  static deserializeBinaryFromReader(message: DisableAccountV2Response, reader: jspb.BinaryReader): DisableAccountV2Response;
}

export namespace DisableAccountV2Response {
  export type AsObject = {
  }
}

export class DeleteAccountV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.DeleteAccountV2Command | undefined;
  setCommand(value?: proto_account_command_pb.DeleteAccountV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAccountV2Request): DeleteAccountV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAccountV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountV2Request;
  static deserializeBinaryFromReader(message: DeleteAccountV2Request, reader: jspb.BinaryReader): DeleteAccountV2Request;
}

export namespace DeleteAccountV2Request {
  export type AsObject = {
    email: string,
    organizationId: string,
    command?: proto_account_command_pb.DeleteAccountV2Command.AsObject,
  }
}

export class DeleteAccountV2Response extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAccountV2Response): DeleteAccountV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAccountV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountV2Response;
  static deserializeBinaryFromReader(message: DeleteAccountV2Response, reader: jspb.BinaryReader): DeleteAccountV2Response;
}

export namespace DeleteAccountV2Response {
  export type AsObject = {
  }
}

export class UpdateAccountV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasChangeNameCommand(): boolean;
  clearChangeNameCommand(): void;
  getChangeNameCommand(): proto_account_command_pb.ChangeAccountV2NameCommand | undefined;
  setChangeNameCommand(value?: proto_account_command_pb.ChangeAccountV2NameCommand): void;

  hasChangeAvatarUrlCommand(): boolean;
  clearChangeAvatarUrlCommand(): void;
  getChangeAvatarUrlCommand(): proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand | undefined;
  setChangeAvatarUrlCommand(value?: proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand): void;

  hasChangeOrganizationRoleCommand(): boolean;
  clearChangeOrganizationRoleCommand(): void;
  getChangeOrganizationRoleCommand(): proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand | undefined;
  setChangeOrganizationRoleCommand(value?: proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand): void;

  hasChangeEnvironmentRolesCommand(): boolean;
  clearChangeEnvironmentRolesCommand(): void;
  getChangeEnvironmentRolesCommand(): proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand | undefined;
  setChangeEnvironmentRolesCommand(value?: proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAccountV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAccountV2Request): UpdateAccountV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateAccountV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAccountV2Request;
  static deserializeBinaryFromReader(message: UpdateAccountV2Request, reader: jspb.BinaryReader): UpdateAccountV2Request;
}

export namespace UpdateAccountV2Request {
  export type AsObject = {
    email: string,
    organizationId: string,
    changeNameCommand?: proto_account_command_pb.ChangeAccountV2NameCommand.AsObject,
    changeAvatarUrlCommand?: proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand.AsObject,
    changeOrganizationRoleCommand?: proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand.AsObject,
    changeEnvironmentRolesCommand?: proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand.AsObject,
  }
}

export class UpdateAccountV2Response extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAccountV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAccountV2Response): UpdateAccountV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateAccountV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAccountV2Response;
  static deserializeBinaryFromReader(message: UpdateAccountV2Response, reader: jspb.BinaryReader): UpdateAccountV2Response;
}

export namespace UpdateAccountV2Response {
  export type AsObject = {
  }
}

export class GetAccountV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: GetAccountV2Request): GetAccountV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAccountV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountV2Request;
  static deserializeBinaryFromReader(message: GetAccountV2Request, reader: jspb.BinaryReader): GetAccountV2Request;
}

export namespace GetAccountV2Request {
  export type AsObject = {
    email: string,
    organizationId: string,
  }
}

export class GetAccountV2Response extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.AccountV2 | undefined;
  setAccount(value?: proto_account_account_pb.AccountV2): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: GetAccountV2Response): GetAccountV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAccountV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountV2Response;
  static deserializeBinaryFromReader(message: GetAccountV2Response, reader: jspb.BinaryReader): GetAccountV2Response;
}

export namespace GetAccountV2Response {
  export type AsObject = {
    account?: proto_account_account_pb.AccountV2.AsObject,
  }
}

export class GetAccountV2ByEnvironmentIDRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountV2ByEnvironmentIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAccountV2ByEnvironmentIDRequest): GetAccountV2ByEnvironmentIDRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAccountV2ByEnvironmentIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountV2ByEnvironmentIDRequest;
  static deserializeBinaryFromReader(message: GetAccountV2ByEnvironmentIDRequest, reader: jspb.BinaryReader): GetAccountV2ByEnvironmentIDRequest;
}

export namespace GetAccountV2ByEnvironmentIDRequest {
  export type AsObject = {
    email: string,
    environmentId: string,
  }
}

export class GetAccountV2ByEnvironmentIDResponse extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.AccountV2 | undefined;
  setAccount(value?: proto_account_account_pb.AccountV2): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountV2ByEnvironmentIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAccountV2ByEnvironmentIDResponse): GetAccountV2ByEnvironmentIDResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAccountV2ByEnvironmentIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountV2ByEnvironmentIDResponse;
  static deserializeBinaryFromReader(message: GetAccountV2ByEnvironmentIDResponse, reader: jspb.BinaryReader): GetAccountV2ByEnvironmentIDResponse;
}

export namespace GetAccountV2ByEnvironmentIDResponse {
  export type AsObject = {
    account?: proto_account_account_pb.AccountV2.AsObject,
  }
}

export class ListAccountsV2Request extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  getOrderBy(): ListAccountsV2Request.OrderByMap[keyof ListAccountsV2Request.OrderByMap];
  setOrderBy(value: ListAccountsV2Request.OrderByMap[keyof ListAccountsV2Request.OrderByMap]): void;

  getOrderDirection(): ListAccountsV2Request.OrderDirectionMap[keyof ListAccountsV2Request.OrderDirectionMap];
  setOrderDirection(value: ListAccountsV2Request.OrderDirectionMap[keyof ListAccountsV2Request.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasRole(): boolean;
  clearRole(): void;
  getRole(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setRole(value?: google_protobuf_wrappers_pb.Int32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAccountsV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: ListAccountsV2Request): ListAccountsV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAccountsV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAccountsV2Request;
  static deserializeBinaryFromReader(message: ListAccountsV2Request, reader: jspb.BinaryReader): ListAccountsV2Request;
}

export namespace ListAccountsV2Request {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    organizationId: string,
    orderBy: ListAccountsV2Request.OrderByMap[keyof ListAccountsV2Request.OrderByMap],
    orderDirection: ListAccountsV2Request.OrderDirectionMap[keyof ListAccountsV2Request.OrderDirectionMap],
    searchKeyword: string,
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    role?: google_protobuf_wrappers_pb.Int32Value.AsObject,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    EMAIL: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListAccountsV2Response extends jspb.Message {
  clearAccountsList(): void;
  getAccountsList(): Array<proto_account_account_pb.AccountV2>;
  setAccountsList(value: Array<proto_account_account_pb.AccountV2>): void;
  addAccounts(value?: proto_account_account_pb.AccountV2, index?: number): proto_account_account_pb.AccountV2;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAccountsV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: ListAccountsV2Response): ListAccountsV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAccountsV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAccountsV2Response;
  static deserializeBinaryFromReader(message: ListAccountsV2Response, reader: jspb.BinaryReader): ListAccountsV2Response;
}

export namespace ListAccountsV2Response {
  export type AsObject = {
    accountsList: Array<proto_account_account_pb.AccountV2.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class CreateAPIKeyRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.CreateAPIKeyCommand | undefined;
  setCommand(value?: proto_account_command_pb.CreateAPIKeyCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAPIKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAPIKeyRequest): CreateAPIKeyRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAPIKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAPIKeyRequest;
  static deserializeBinaryFromReader(message: CreateAPIKeyRequest, reader: jspb.BinaryReader): CreateAPIKeyRequest;
}

export namespace CreateAPIKeyRequest {
  export type AsObject = {
    command?: proto_account_command_pb.CreateAPIKeyCommand.AsObject,
    environmentNamespace: string,
  }
}

export class CreateAPIKeyResponse extends jspb.Message {
  hasApiKey(): boolean;
  clearApiKey(): void;
  getApiKey(): proto_account_api_key_pb.APIKey | undefined;
  setApiKey(value?: proto_account_api_key_pb.APIKey): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAPIKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAPIKeyResponse): CreateAPIKeyResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAPIKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAPIKeyResponse;
  static deserializeBinaryFromReader(message: CreateAPIKeyResponse, reader: jspb.BinaryReader): CreateAPIKeyResponse;
}

export namespace CreateAPIKeyResponse {
  export type AsObject = {
    apiKey?: proto_account_api_key_pb.APIKey.AsObject,
  }
}

export class ChangeAPIKeyNameRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.ChangeAPIKeyNameCommand | undefined;
  setCommand(value?: proto_account_command_pb.ChangeAPIKeyNameCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAPIKeyNameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAPIKeyNameRequest): ChangeAPIKeyNameRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAPIKeyNameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAPIKeyNameRequest;
  static deserializeBinaryFromReader(message: ChangeAPIKeyNameRequest, reader: jspb.BinaryReader): ChangeAPIKeyNameRequest;
}

export namespace ChangeAPIKeyNameRequest {
  export type AsObject = {
    id: string,
    command?: proto_account_command_pb.ChangeAPIKeyNameCommand.AsObject,
    environmentNamespace: string,
  }
}

export class ChangeAPIKeyNameResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAPIKeyNameResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAPIKeyNameResponse): ChangeAPIKeyNameResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAPIKeyNameResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAPIKeyNameResponse;
  static deserializeBinaryFromReader(message: ChangeAPIKeyNameResponse, reader: jspb.BinaryReader): ChangeAPIKeyNameResponse;
}

export namespace ChangeAPIKeyNameResponse {
  export type AsObject = {
  }
}

export class EnableAPIKeyRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.EnableAPIKeyCommand | undefined;
  setCommand(value?: proto_account_command_pb.EnableAPIKeyCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAPIKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAPIKeyRequest): EnableAPIKeyRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAPIKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAPIKeyRequest;
  static deserializeBinaryFromReader(message: EnableAPIKeyRequest, reader: jspb.BinaryReader): EnableAPIKeyRequest;
}

export namespace EnableAPIKeyRequest {
  export type AsObject = {
    id: string,
    command?: proto_account_command_pb.EnableAPIKeyCommand.AsObject,
    environmentNamespace: string,
  }
}

export class EnableAPIKeyResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAPIKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAPIKeyResponse): EnableAPIKeyResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAPIKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAPIKeyResponse;
  static deserializeBinaryFromReader(message: EnableAPIKeyResponse, reader: jspb.BinaryReader): EnableAPIKeyResponse;
}

export namespace EnableAPIKeyResponse {
  export type AsObject = {
  }
}

export class DisableAPIKeyRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.DisableAPIKeyCommand | undefined;
  setCommand(value?: proto_account_command_pb.DisableAPIKeyCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAPIKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAPIKeyRequest): DisableAPIKeyRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAPIKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAPIKeyRequest;
  static deserializeBinaryFromReader(message: DisableAPIKeyRequest, reader: jspb.BinaryReader): DisableAPIKeyRequest;
}

export namespace DisableAPIKeyRequest {
  export type AsObject = {
    id: string,
    command?: proto_account_command_pb.DisableAPIKeyCommand.AsObject,
    environmentNamespace: string,
  }
}

export class DisableAPIKeyResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAPIKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAPIKeyResponse): DisableAPIKeyResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAPIKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAPIKeyResponse;
  static deserializeBinaryFromReader(message: DisableAPIKeyResponse, reader: jspb.BinaryReader): DisableAPIKeyResponse;
}

export namespace DisableAPIKeyResponse {
  export type AsObject = {
  }
}

export class GetAPIKeyRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAPIKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAPIKeyRequest): GetAPIKeyRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAPIKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAPIKeyRequest;
  static deserializeBinaryFromReader(message: GetAPIKeyRequest, reader: jspb.BinaryReader): GetAPIKeyRequest;
}

export namespace GetAPIKeyRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
  }
}

export class GetAPIKeyResponse extends jspb.Message {
  hasApiKey(): boolean;
  clearApiKey(): void;
  getApiKey(): proto_account_api_key_pb.APIKey | undefined;
  setApiKey(value?: proto_account_api_key_pb.APIKey): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAPIKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAPIKeyResponse): GetAPIKeyResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAPIKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAPIKeyResponse;
  static deserializeBinaryFromReader(message: GetAPIKeyResponse, reader: jspb.BinaryReader): GetAPIKeyResponse;
}

export namespace GetAPIKeyResponse {
  export type AsObject = {
    apiKey?: proto_account_api_key_pb.APIKey.AsObject,
  }
}

export class ListAPIKeysRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getOrderBy(): ListAPIKeysRequest.OrderByMap[keyof ListAPIKeysRequest.OrderByMap];
  setOrderBy(value: ListAPIKeysRequest.OrderByMap[keyof ListAPIKeysRequest.OrderByMap]): void;

  getOrderDirection(): ListAPIKeysRequest.OrderDirectionMap[keyof ListAPIKeysRequest.OrderDirectionMap];
  setOrderDirection(value: ListAPIKeysRequest.OrderDirectionMap[keyof ListAPIKeysRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAPIKeysRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAPIKeysRequest): ListAPIKeysRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAPIKeysRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAPIKeysRequest;
  static deserializeBinaryFromReader(message: ListAPIKeysRequest, reader: jspb.BinaryReader): ListAPIKeysRequest;
}

export namespace ListAPIKeysRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    environmentNamespace: string,
    orderBy: ListAPIKeysRequest.OrderByMap[keyof ListAPIKeysRequest.OrderByMap],
    orderDirection: ListAPIKeysRequest.OrderDirectionMap[keyof ListAPIKeysRequest.OrderDirectionMap],
    searchKeyword: string,
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    NAME: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListAPIKeysResponse extends jspb.Message {
  clearApiKeysList(): void;
  getApiKeysList(): Array<proto_account_api_key_pb.APIKey>;
  setApiKeysList(value: Array<proto_account_api_key_pb.APIKey>): void;
  addApiKeys(value?: proto_account_api_key_pb.APIKey, index?: number): proto_account_api_key_pb.APIKey;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAPIKeysResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAPIKeysResponse): ListAPIKeysResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAPIKeysResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAPIKeysResponse;
  static deserializeBinaryFromReader(message: ListAPIKeysResponse, reader: jspb.BinaryReader): ListAPIKeysResponse;
}

export namespace ListAPIKeysResponse {
  export type AsObject = {
    apiKeysList: Array<proto_account_api_key_pb.APIKey.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class GetAPIKeyBySearchingAllEnvironmentsRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAPIKeyBySearchingAllEnvironmentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAPIKeyBySearchingAllEnvironmentsRequest): GetAPIKeyBySearchingAllEnvironmentsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAPIKeyBySearchingAllEnvironmentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAPIKeyBySearchingAllEnvironmentsRequest;
  static deserializeBinaryFromReader(message: GetAPIKeyBySearchingAllEnvironmentsRequest, reader: jspb.BinaryReader): GetAPIKeyBySearchingAllEnvironmentsRequest;
}

export namespace GetAPIKeyBySearchingAllEnvironmentsRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetAPIKeyBySearchingAllEnvironmentsResponse extends jspb.Message {
  hasEnvironmentApiKey(): boolean;
  clearEnvironmentApiKey(): void;
  getEnvironmentApiKey(): proto_account_api_key_pb.EnvironmentAPIKey | undefined;
  setEnvironmentApiKey(value?: proto_account_api_key_pb.EnvironmentAPIKey): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAPIKeyBySearchingAllEnvironmentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAPIKeyBySearchingAllEnvironmentsResponse): GetAPIKeyBySearchingAllEnvironmentsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAPIKeyBySearchingAllEnvironmentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAPIKeyBySearchingAllEnvironmentsResponse;
  static deserializeBinaryFromReader(message: GetAPIKeyBySearchingAllEnvironmentsResponse, reader: jspb.BinaryReader): GetAPIKeyBySearchingAllEnvironmentsResponse;
}

export namespace GetAPIKeyBySearchingAllEnvironmentsResponse {
  export type AsObject = {
    environmentApiKey?: proto_account_api_key_pb.EnvironmentAPIKey.AsObject,
  }
}

