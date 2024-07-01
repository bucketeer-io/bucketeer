// package: bucketeer.auth
// file: proto/auth/service.proto

import * as jspb from 'google-protobuf';
import * as proto_auth_token_pb from '../../proto/auth/token_pb';

export class GetAuthCodeURLRequest extends jspb.Message {
  getState(): string;
  setState(value: string): void;

  getRedirectUrl(): string;
  setRedirectUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthCodeURLRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAuthCodeURLRequest
  ): GetAuthCodeURLRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAuthCodeURLRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthCodeURLRequest;
  static deserializeBinaryFromReader(
    message: GetAuthCodeURLRequest,
    reader: jspb.BinaryReader
  ): GetAuthCodeURLRequest;
}

export namespace GetAuthCodeURLRequest {
  export type AsObject = {
    state: string;
    redirectUrl: string;
  };
}

export class GetAuthCodeURLResponse extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthCodeURLResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAuthCodeURLResponse
  ): GetAuthCodeURLResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAuthCodeURLResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthCodeURLResponse;
  static deserializeBinaryFromReader(
    message: GetAuthCodeURLResponse,
    reader: jspb.BinaryReader
  ): GetAuthCodeURLResponse;
}

export namespace GetAuthCodeURLResponse {
  export type AsObject = {
    url: string;
  };
}

export class ExchangeTokenRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  getRedirectUrl(): string;
  setRedirectUrl(value: string): void;

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

  getRedirectUrl(): string;
  setRedirectUrl(value: string): void;

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
    redirectUrl: string;
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

export class ExchangeBucketeerTokenRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  getRedirectUrl(): string;
  setRedirectUrl(value: string): void;

  getType(): AuthTypeMap[keyof AuthTypeMap];
  setType(value: AuthTypeMap[keyof AuthTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExchangeBucketeerTokenRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExchangeBucketeerTokenRequest
  ): ExchangeBucketeerTokenRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExchangeBucketeerTokenRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExchangeBucketeerTokenRequest;
  static deserializeBinaryFromReader(
    message: ExchangeBucketeerTokenRequest,
    reader: jspb.BinaryReader
  ): ExchangeBucketeerTokenRequest;
}

export namespace ExchangeBucketeerTokenRequest {
  export type AsObject = {
    code: string;
    redirectUrl: string;
    type: AuthTypeMap[keyof AuthTypeMap];
  };
}

export class ExchangeBucketeerTokenResponse extends jspb.Message {
  hasToken(): boolean;
  clearToken(): void;
  getToken(): proto_auth_token_pb.Token | undefined;
  setToken(value?: proto_auth_token_pb.Token): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExchangeBucketeerTokenResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExchangeBucketeerTokenResponse
  ): ExchangeBucketeerTokenResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExchangeBucketeerTokenResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExchangeBucketeerTokenResponse;
  static deserializeBinaryFromReader(
    message: ExchangeBucketeerTokenResponse,
    reader: jspb.BinaryReader
  ): ExchangeBucketeerTokenResponse;
}

export namespace ExchangeBucketeerTokenResponse {
  export type AsObject = {
    token?: proto_auth_token_pb.Token.AsObject;
  };
}

export class RefreshBucketeerTokenRequest extends jspb.Message {
  getRefreshToken(): string;
  setRefreshToken(value: string): void;

  getRedirectUrl(): string;
  setRedirectUrl(value: string): void;

  getType(): AuthTypeMap[keyof AuthTypeMap];
  setType(value: AuthTypeMap[keyof AuthTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshBucketeerTokenRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: RefreshBucketeerTokenRequest
  ): RefreshBucketeerTokenRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: RefreshBucketeerTokenRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): RefreshBucketeerTokenRequest;
  static deserializeBinaryFromReader(
    message: RefreshBucketeerTokenRequest,
    reader: jspb.BinaryReader
  ): RefreshBucketeerTokenRequest;
}

export namespace RefreshBucketeerTokenRequest {
  export type AsObject = {
    refreshToken: string;
    redirectUrl: string;
    type: AuthTypeMap[keyof AuthTypeMap];
  };
}

export class RefreshBucketeerTokenResponse extends jspb.Message {
  hasToken(): boolean;
  clearToken(): void;
  getToken(): proto_auth_token_pb.Token | undefined;
  setToken(value?: proto_auth_token_pb.Token): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshBucketeerTokenResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: RefreshBucketeerTokenResponse
  ): RefreshBucketeerTokenResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: RefreshBucketeerTokenResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): RefreshBucketeerTokenResponse;
  static deserializeBinaryFromReader(
    message: RefreshBucketeerTokenResponse,
    reader: jspb.BinaryReader
  ): RefreshBucketeerTokenResponse;
}

export namespace RefreshBucketeerTokenResponse {
  export type AsObject = {
    token?: proto_auth_token_pb.Token.AsObject;
  };
}

export interface AuthTypeMap {
  AUTH_TYPE_UNSPECIFIED: 0;
  AUTH_TYPE_USER_PASSWORD: 1;
  AUTH_TYPE_GOOGLE: 2;
  AUTH_TYPE_GITHUB: 3;
}

export const AuthType: AuthTypeMap;
