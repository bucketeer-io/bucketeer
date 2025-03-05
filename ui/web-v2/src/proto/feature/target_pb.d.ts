// package: bucketeer.feature
// file: proto/feature/target.proto

import * as jspb from 'google-protobuf';

export class Target extends jspb.Message {
  getVariation(): string;
  setVariation(value: string): void;

  clearUsersList(): void;
  getUsersList(): Array<string>;
  setUsersList(value: Array<string>): void;
  addUsers(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Target.AsObject;
  static toObject(includeInstance: boolean, msg: Target): Target.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Target,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Target;
  static deserializeBinaryFromReader(
    message: Target,
    reader: jspb.BinaryReader
  ): Target;
}

export namespace Target {
  export type AsObject = {
    variation: string;
    usersList: Array<string>;
  };
}

export class TargetListValue extends jspb.Message {
  clearValuesList(): void;
  getValuesList(): Array<Target>;
  setValuesList(value: Array<Target>): void;
  addValues(value?: Target, index?: number): Target;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TargetListValue.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: TargetListValue
  ): TargetListValue.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: TargetListValue,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): TargetListValue;
  static deserializeBinaryFromReader(
    message: TargetListValue,
    reader: jspb.BinaryReader
  ): TargetListValue;
}

export namespace TargetListValue {
  export type AsObject = {
    valuesList: Array<Target.AsObject>;
  };
}
