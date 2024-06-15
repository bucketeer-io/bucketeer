// package: bucketeer.auth
// file: proto/auth/token.proto

import * as jspb from "google-protobuf";

export class Token extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): void;

  getTokenType(): string;
  setTokenType(value: string): void;

  getRefreshToken(): string;
  setRefreshToken(value: string): void;

  getExpiry(): number;
  setExpiry(value: number): void;

  getIdToken(): string;
  setIdToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Token.AsObject;
  static toObject(includeInstance: boolean, msg: Token): Token.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Token, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Token;
  static deserializeBinaryFromReader(message: Token, reader: jspb.BinaryReader): Token;
}

export namespace Token {
  export type AsObject = {
    accessToken: string,
    tokenType: string,
    refreshToken: string,
    expiry: number,
    idToken: string,
  }
}

export class IDTokenSubject extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getConnId(): string;
  setConnId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDTokenSubject.AsObject;
  static toObject(includeInstance: boolean, msg: IDTokenSubject): IDTokenSubject.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: IDTokenSubject, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDTokenSubject;
  static deserializeBinaryFromReader(message: IDTokenSubject, reader: jspb.BinaryReader): IDTokenSubject;
}

export namespace IDTokenSubject {
  export type AsObject = {
    userId: string,
    connId: string,
  }
}

