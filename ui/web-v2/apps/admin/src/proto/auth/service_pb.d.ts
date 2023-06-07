// package: bucketeer.auth
// file: proto/auth/service.proto

import * as jspb from "google-protobuf";
import * as proto_auth_token_pb from "../../proto/auth/token_pb";

export class GetAuthCodeURLRequest extends jspb.Message {
  getState(): string;
  setState(value: string): void;

  getRedirectUrl(): string;
  setRedirectUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthCodeURLRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAuthCodeURLRequest): GetAuthCodeURLRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAuthCodeURLRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthCodeURLRequest;
  static deserializeBinaryFromReader(message: GetAuthCodeURLRequest, reader: jspb.BinaryReader): GetAuthCodeURLRequest;
}

export namespace GetAuthCodeURLRequest {
  export type AsObject = {
    state: string,
    redirectUrl: string,
  }
}

export class GetAuthCodeURLResponse extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthCodeURLResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAuthCodeURLResponse): GetAuthCodeURLResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAuthCodeURLResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthCodeURLResponse;
  static deserializeBinaryFromReader(message: GetAuthCodeURLResponse, reader: jspb.BinaryReader): GetAuthCodeURLResponse;
}

export namespace GetAuthCodeURLResponse {
  export type AsObject = {
    url: string,
  }
}

export class ExchangeTokenRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  getRedirectUrl(): string;
  setRedirectUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExchangeTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ExchangeTokenRequest): ExchangeTokenRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExchangeTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExchangeTokenRequest;
  static deserializeBinaryFromReader(message: ExchangeTokenRequest, reader: jspb.BinaryReader): ExchangeTokenRequest;
}

export namespace ExchangeTokenRequest {
  export type AsObject = {
    code: string,
    redirectUrl: string,
  }
}

export class ExchangeTokenResponse extends jspb.Message {
  hasToken(): boolean;
  clearToken(): void;
  getToken(): proto_auth_token_pb.Token | undefined;
  setToken(value?: proto_auth_token_pb.Token): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExchangeTokenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ExchangeTokenResponse): ExchangeTokenResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExchangeTokenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExchangeTokenResponse;
  static deserializeBinaryFromReader(message: ExchangeTokenResponse, reader: jspb.BinaryReader): ExchangeTokenResponse;
}

export namespace ExchangeTokenResponse {
  export type AsObject = {
    token?: proto_auth_token_pb.Token.AsObject,
  }
}

export class RefreshTokenRequest extends jspb.Message {
  getRefreshToken(): string;
  setRefreshToken(value: string): void;

  getRedirectUrl(): string;
  setRedirectUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RefreshTokenRequest): RefreshTokenRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RefreshTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RefreshTokenRequest;
  static deserializeBinaryFromReader(message: RefreshTokenRequest, reader: jspb.BinaryReader): RefreshTokenRequest;
}

export namespace RefreshTokenRequest {
  export type AsObject = {
    refreshToken: string,
    redirectUrl: string,
  }
}

export class RefreshTokenResponse extends jspb.Message {
  hasToken(): boolean;
  clearToken(): void;
  getToken(): proto_auth_token_pb.Token | undefined;
  setToken(value?: proto_auth_token_pb.Token): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshTokenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RefreshTokenResponse): RefreshTokenResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RefreshTokenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RefreshTokenResponse;
  static deserializeBinaryFromReader(message: RefreshTokenResponse, reader: jspb.BinaryReader): RefreshTokenResponse;
}

export namespace RefreshTokenResponse {
  export type AsObject = {
    token?: proto_auth_token_pb.Token.AsObject,
  }
}

