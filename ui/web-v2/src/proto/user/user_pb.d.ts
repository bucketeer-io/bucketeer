// package: bucketeer.user
// file: proto/user/user.proto

import * as jspb from 'google-protobuf';

export class User extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDataMap(): jspb.Map<string, string>;
  clearDataMap(): void;
  getTaggedDataMap(): jspb.Map<string, User.Data>;
  clearTaggedDataMap(): void;
  getLastSeen(): number;
  setLastSeen(value: number): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: User,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(
    message: User,
    reader: jspb.BinaryReader
  ): User;
}

export namespace User {
  export type AsObject = {
    id: string;
    dataMap: Array<[string, string]>;
    taggedDataMap: Array<[string, User.Data.AsObject]>;
    lastSeen: number;
    createdAt: number;
  };

  export class Data extends jspb.Message {
    getValueMap(): jspb.Map<string, string>;
    clearValueMap(): void;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Data.AsObject;
    static toObject(includeInstance: boolean, msg: Data): Data.AsObject;
    static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
    static extensionsBinary: {
      [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
    };
    static serializeBinaryToWriter(
      message: Data,
      writer: jspb.BinaryWriter
    ): void;
    static deserializeBinary(bytes: Uint8Array): Data;
    static deserializeBinaryFromReader(
      message: Data,
      reader: jspb.BinaryReader
    ): Data;
  }

  export namespace Data {
    export type AsObject = {
      valueMap: Array<[string, string]>;
    };
  }
}

export class UserAttributes extends jspb.Message {
  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  clearUserattributesList(): void;
  getUserattributesList(): Array<UserAttribute>;
  setUserattributesList(value: Array<UserAttribute>): void;
  addUserattributes(value?: UserAttribute, index?: number): UserAttribute;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserAttributes.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UserAttributes
  ): UserAttributes.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UserAttributes,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UserAttributes;
  static deserializeBinaryFromReader(
    message: UserAttributes,
    reader: jspb.BinaryReader
  ): UserAttributes;
}

export namespace UserAttributes {
  export type AsObject = {
    environmentId: string;
    userattributesList: Array<UserAttribute.AsObject>;
  };
}

export class UserAttribute extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  clearValuesList(): void;
  getValuesList(): Array<string>;
  setValuesList(value: Array<string>): void;
  addValues(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserAttribute.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UserAttribute
  ): UserAttribute.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UserAttribute,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UserAttribute;
  static deserializeBinaryFromReader(
    message: UserAttribute,
    reader: jspb.BinaryReader
  ): UserAttribute;
}

export namespace UserAttribute {
  export type AsObject = {
    key: string;
    valuesList: Array<string>;
  };
}
