// package: bucketeer.auth
// file: proto/auth/service.proto

import * as jspb from 'google-protobuf';
import * as google_api_annotations_pb from '../../google/api/annotations_pb';
import * as protoc_gen_openapiv2_options_annotations_pb from '../../protoc-gen-openapiv2/options/annotations_pb';
import * as proto_auth_token_pb from '../../proto/auth/token_pb';

export class GetAuthenticationURLRequest extends jspb.Message {
  getState(): string;
  setState(value: string): void;

  getRedirectUrl(): string;
  setRedirectUrl(value: string): void;

  getType(): AuthTypeMap[keyof AuthTypeMap];
  setType(value: AuthTypeMap[keyof AuthTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthenticationURLRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAuthenticationURLRequest
  ): GetAuthenticationURLRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAuthenticationURLRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthenticationURLRequest;
  static deserializeBinaryFromReader(
    message: GetAuthenticationURLRequest,
    reader: jspb.BinaryReader
  ): GetAuthenticationURLRequest;
}

export namespace GetAuthenticationURLRequest {
  export type AsObject = {
    state: string;
    redirectUrl: string;
    type: AuthTypeMap[keyof AuthTypeMap];
  };
}

export class GetAuthenticationURLResponse extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthenticationURLResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAuthenticationURLResponse
  ): GetAuthenticationURLResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAuthenticationURLResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthenticationURLResponse;
  static deserializeBinaryFromReader(
    message: GetAuthenticationURLResponse,
    reader: jspb.BinaryReader
  ): GetAuthenticationURLResponse;
}

export namespace GetAuthenticationURLResponse {
  export type AsObject = {
    url: string;
  };
}

export class ExchangeTokenRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  getRedirectUrl(): string;
  setRedirectUrl(value: string): void;

  getType(): AuthTypeMap[keyof AuthTypeMap];
  setType(value: AuthTypeMap[keyof AuthTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExchangeTokenRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExchangeTokenRequest
  ): ExchangeTokenRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExchangeTokenRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExchangeTokenRequest;
  static deserializeBinaryFromReader(
    message: ExchangeTokenRequest,
    reader: jspb.BinaryReader
  ): ExchangeTokenRequest;
}

export namespace ExchangeTokenRequest {
  export type AsObject = {
    code: string;
    redirectUrl: string;
    type: AuthTypeMap[keyof AuthTypeMap];
  };
}

export class ExchangeTokenResponse extends jspb.Message {
  hasToken(): boolean;
  clearToken(): void;
  getToken(): proto_auth_token_pb.Token | undefined;
  setToken(value?: proto_auth_token_pb.Token): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExchangeTokenResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExchangeTokenResponse
  ): ExchangeTokenResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExchangeTokenResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExchangeTokenResponse;
  static deserializeBinaryFromReader(
    message: ExchangeTokenResponse,
    reader: jspb.BinaryReader
  ): ExchangeTokenResponse;
}

export namespace ExchangeTokenResponse {
  export type AsObject = {
    token?: proto_auth_token_pb.Token.AsObject;
  };
}

export class RefreshTokenRequest extends jspb.Message {
  getRefreshToken(): string;
  setRefreshToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshTokenRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: RefreshTokenRequest
  ): RefreshTokenRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: RefreshTokenRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): RefreshTokenRequest;
  static deserializeBinaryFromReader(
    message: RefreshTokenRequest,
    reader: jspb.BinaryReader
  ): RefreshTokenRequest;
}

export namespace RefreshTokenRequest {
  export type AsObject = {
    refreshToken: string;
  };
}

export class RefreshTokenResponse extends jspb.Message {
  hasToken(): boolean;
  clearToken(): void;
  getToken(): proto_auth_token_pb.Token | undefined;
  setToken(value?: proto_auth_token_pb.Token): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshTokenResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: RefreshTokenResponse
  ): RefreshTokenResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: RefreshTokenResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): RefreshTokenResponse;
  static deserializeBinaryFromReader(
    message: RefreshTokenResponse,
    reader: jspb.BinaryReader
  ): RefreshTokenResponse;
}

export namespace RefreshTokenResponse {
  export type AsObject = {
    token?: proto_auth_token_pb.Token.AsObject;
  };
}

export class SignInRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SignInRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SignInRequest
  ): SignInRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SignInRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SignInRequest;
  static deserializeBinaryFromReader(
    message: SignInRequest,
    reader: jspb.BinaryReader
  ): SignInRequest;
}

export namespace SignInRequest {
  export type AsObject = {
    email: string;
    password: string;
  };
}

export class SignInResponse extends jspb.Message {
  hasToken(): boolean;
  clearToken(): void;
  getToken(): proto_auth_token_pb.Token | undefined;
  setToken(value?: proto_auth_token_pb.Token): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SignInResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SignInResponse
  ): SignInResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SignInResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SignInResponse;
  static deserializeBinaryFromReader(
    message: SignInResponse,
    reader: jspb.BinaryReader
  ): SignInResponse;
}

export namespace SignInResponse {
  export type AsObject = {
    token?: proto_auth_token_pb.Token.AsObject;
  };
}

export class SwitchOrganizationRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SwitchOrganizationRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SwitchOrganizationRequest
  ): SwitchOrganizationRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SwitchOrganizationRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SwitchOrganizationRequest;
  static deserializeBinaryFromReader(
    message: SwitchOrganizationRequest,
    reader: jspb.BinaryReader
  ): SwitchOrganizationRequest;
}

export namespace SwitchOrganizationRequest {
  export type AsObject = {
    accessToken: string;
    organizationId: string;
  };
}

export class SwitchOrganizationResponse extends jspb.Message {
  hasToken(): boolean;
  clearToken(): void;
  getToken(): proto_auth_token_pb.Token | undefined;
  setToken(value?: proto_auth_token_pb.Token): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SwitchOrganizationResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SwitchOrganizationResponse
  ): SwitchOrganizationResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SwitchOrganizationResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SwitchOrganizationResponse;
  static deserializeBinaryFromReader(
    message: SwitchOrganizationResponse,
    reader: jspb.BinaryReader
  ): SwitchOrganizationResponse;
}

export namespace SwitchOrganizationResponse {
  export type AsObject = {
    token?: proto_auth_token_pb.Token.AsObject;
  };
}

export class GetDemoSiteStatusRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDemoSiteStatusRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetDemoSiteStatusRequest
  ): GetDemoSiteStatusRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetDemoSiteStatusRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetDemoSiteStatusRequest;
  static deserializeBinaryFromReader(
    message: GetDemoSiteStatusRequest,
    reader: jspb.BinaryReader
  ): GetDemoSiteStatusRequest;
}

export namespace GetDemoSiteStatusRequest {
  export type AsObject = {};
}

export class GetDemoSiteStatusResponse extends jspb.Message {
  getIsDemoSiteEnabled(): boolean;
  setIsDemoSiteEnabled(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDemoSiteStatusResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetDemoSiteStatusResponse
  ): GetDemoSiteStatusResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetDemoSiteStatusResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetDemoSiteStatusResponse;
  static deserializeBinaryFromReader(
    message: GetDemoSiteStatusResponse,
    reader: jspb.BinaryReader
  ): GetDemoSiteStatusResponse;
}

export namespace GetDemoSiteStatusResponse {
  export type AsObject = {
    isDemoSiteEnabled: boolean;
  };
}

export class CreatePasswordRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePasswordRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreatePasswordRequest
  ): CreatePasswordRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreatePasswordRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreatePasswordRequest;
  static deserializeBinaryFromReader(
    message: CreatePasswordRequest,
    reader: jspb.BinaryReader
  ): CreatePasswordRequest;
}

export namespace CreatePasswordRequest {
  export type AsObject = {
    email: string;
    password: string;
  };
}

export class CreatePasswordResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePasswordResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreatePasswordResponse
  ): CreatePasswordResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreatePasswordResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreatePasswordResponse;
  static deserializeBinaryFromReader(
    message: CreatePasswordResponse,
    reader: jspb.BinaryReader
  ): CreatePasswordResponse;
}

export namespace CreatePasswordResponse {
  export type AsObject = {};
}

export class UpdatePasswordRequest extends jspb.Message {
  getCurrentPassword(): string;
  setCurrentPassword(value: string): void;

  getNewPassword(): string;
  setNewPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePasswordRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdatePasswordRequest
  ): UpdatePasswordRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdatePasswordRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePasswordRequest;
  static deserializeBinaryFromReader(
    message: UpdatePasswordRequest,
    reader: jspb.BinaryReader
  ): UpdatePasswordRequest;
}

export namespace UpdatePasswordRequest {
  export type AsObject = {
    currentPassword: string;
    newPassword: string;
  };
}

export class UpdatePasswordResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePasswordResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdatePasswordResponse
  ): UpdatePasswordResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdatePasswordResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePasswordResponse;
  static deserializeBinaryFromReader(
    message: UpdatePasswordResponse,
    reader: jspb.BinaryReader
  ): UpdatePasswordResponse;
}

export namespace UpdatePasswordResponse {
  export type AsObject = {};
}

export class InitiatePasswordResetRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitiatePasswordResetRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: InitiatePasswordResetRequest
  ): InitiatePasswordResetRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: InitiatePasswordResetRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): InitiatePasswordResetRequest;
  static deserializeBinaryFromReader(
    message: InitiatePasswordResetRequest,
    reader: jspb.BinaryReader
  ): InitiatePasswordResetRequest;
}

export namespace InitiatePasswordResetRequest {
  export type AsObject = {
    email: string;
  };
}

export class InitiatePasswordResetResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitiatePasswordResetResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: InitiatePasswordResetResponse
  ): InitiatePasswordResetResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: InitiatePasswordResetResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): InitiatePasswordResetResponse;
  static deserializeBinaryFromReader(
    message: InitiatePasswordResetResponse,
    reader: jspb.BinaryReader
  ): InitiatePasswordResetResponse;
}

export namespace InitiatePasswordResetResponse {
  export type AsObject = {
    message: string;
  };
}

export class ResetPasswordRequest extends jspb.Message {
  getResetToken(): string;
  setResetToken(value: string): void;

  getNewPassword(): string;
  setNewPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPasswordRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ResetPasswordRequest
  ): ResetPasswordRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ResetPasswordRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ResetPasswordRequest;
  static deserializeBinaryFromReader(
    message: ResetPasswordRequest,
    reader: jspb.BinaryReader
  ): ResetPasswordRequest;
}

export namespace ResetPasswordRequest {
  export type AsObject = {
    resetToken: string;
    newPassword: string;
  };
}

export class ResetPasswordResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPasswordResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ResetPasswordResponse
  ): ResetPasswordResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ResetPasswordResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ResetPasswordResponse;
  static deserializeBinaryFromReader(
    message: ResetPasswordResponse,
    reader: jspb.BinaryReader
  ): ResetPasswordResponse;
}

export namespace ResetPasswordResponse {
  export type AsObject = {};
}

export class ValidatePasswordResetTokenRequest extends jspb.Message {
  getResetToken(): string;
  setResetToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ValidatePasswordResetTokenRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ValidatePasswordResetTokenRequest
  ): ValidatePasswordResetTokenRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ValidatePasswordResetTokenRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ValidatePasswordResetTokenRequest;
  static deserializeBinaryFromReader(
    message: ValidatePasswordResetTokenRequest,
    reader: jspb.BinaryReader
  ): ValidatePasswordResetTokenRequest;
}

export namespace ValidatePasswordResetTokenRequest {
  export type AsObject = {
    resetToken: string;
  };
}

export class ValidatePasswordResetTokenResponse extends jspb.Message {
  getIsValid(): boolean;
  setIsValid(value: boolean): void;

  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ValidatePasswordResetTokenResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ValidatePasswordResetTokenResponse
  ): ValidatePasswordResetTokenResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ValidatePasswordResetTokenResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ValidatePasswordResetTokenResponse;
  static deserializeBinaryFromReader(
    message: ValidatePasswordResetTokenResponse,
    reader: jspb.BinaryReader
  ): ValidatePasswordResetTokenResponse;
}

export namespace ValidatePasswordResetTokenResponse {
  export type AsObject = {
    isValid: boolean;
    email: string;
  };
}

export class InitiatePasswordSetupRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitiatePasswordSetupRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: InitiatePasswordSetupRequest
  ): InitiatePasswordSetupRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: InitiatePasswordSetupRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): InitiatePasswordSetupRequest;
  static deserializeBinaryFromReader(
    message: InitiatePasswordSetupRequest,
    reader: jspb.BinaryReader
  ): InitiatePasswordSetupRequest;
}

export namespace InitiatePasswordSetupRequest {
  export type AsObject = {
    email: string;
  };
}

export class InitiatePasswordSetupResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitiatePasswordSetupResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: InitiatePasswordSetupResponse
  ): InitiatePasswordSetupResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: InitiatePasswordSetupResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): InitiatePasswordSetupResponse;
  static deserializeBinaryFromReader(
    message: InitiatePasswordSetupResponse,
    reader: jspb.BinaryReader
  ): InitiatePasswordSetupResponse;
}

export namespace InitiatePasswordSetupResponse {
  export type AsObject = {
    message: string;
  };
}

export class SetupPasswordRequest extends jspb.Message {
  getSetupToken(): string;
  setSetupToken(value: string): void;

  getNewPassword(): string;
  setNewPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetupPasswordRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SetupPasswordRequest
  ): SetupPasswordRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SetupPasswordRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SetupPasswordRequest;
  static deserializeBinaryFromReader(
    message: SetupPasswordRequest,
    reader: jspb.BinaryReader
  ): SetupPasswordRequest;
}

export namespace SetupPasswordRequest {
  export type AsObject = {
    setupToken: string;
    newPassword: string;
  };
}

export class SetupPasswordResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetupPasswordResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SetupPasswordResponse
  ): SetupPasswordResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SetupPasswordResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SetupPasswordResponse;
  static deserializeBinaryFromReader(
    message: SetupPasswordResponse,
    reader: jspb.BinaryReader
  ): SetupPasswordResponse;
}

export namespace SetupPasswordResponse {
  export type AsObject = {};
}

export class ValidatePasswordSetupTokenRequest extends jspb.Message {
  getSetupToken(): string;
  setSetupToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ValidatePasswordSetupTokenRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ValidatePasswordSetupTokenRequest
  ): ValidatePasswordSetupTokenRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ValidatePasswordSetupTokenRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ValidatePasswordSetupTokenRequest;
  static deserializeBinaryFromReader(
    message: ValidatePasswordSetupTokenRequest,
    reader: jspb.BinaryReader
  ): ValidatePasswordSetupTokenRequest;
}

export namespace ValidatePasswordSetupTokenRequest {
  export type AsObject = {
    setupToken: string;
  };
}

export class ValidatePasswordSetupTokenResponse extends jspb.Message {
  getIsValid(): boolean;
  setIsValid(value: boolean): void;

  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ValidatePasswordSetupTokenResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ValidatePasswordSetupTokenResponse
  ): ValidatePasswordSetupTokenResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ValidatePasswordSetupTokenResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ValidatePasswordSetupTokenResponse;
  static deserializeBinaryFromReader(
    message: ValidatePasswordSetupTokenResponse,
    reader: jspb.BinaryReader
  ): ValidatePasswordSetupTokenResponse;
}

export namespace ValidatePasswordSetupTokenResponse {
  export type AsObject = {
    isValid: boolean;
    email: string;
  };
}

export interface AuthTypeMap {
  AUTH_TYPE_UNSPECIFIED: 0;
  AUTH_TYPE_USER_PASSWORD: 1;
  AUTH_TYPE_GOOGLE: 2;
  AUTH_TYPE_GITHUB: 3;
}

export const AuthType: AuthTypeMap;
