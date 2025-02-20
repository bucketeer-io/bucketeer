// package: bucketeer.feature
// file: proto/feature/prerequisite.proto

import * as jspb from 'google-protobuf';

export class Prerequisite extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Prerequisite.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: Prerequisite
  ): Prerequisite.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Prerequisite,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Prerequisite;
  static deserializeBinaryFromReader(
    message: Prerequisite,
    reader: jspb.BinaryReader
  ): Prerequisite;
}

export namespace Prerequisite {
  export type AsObject = {
    featureId: string;
    variationId: string;
  };
}

export class PrerequisiteListValue extends jspb.Message {
  clearValuesList(): void;
  getValuesList(): Array<Prerequisite>;
  setValuesList(value: Array<Prerequisite>): void;
  addValues(value?: Prerequisite, index?: number): Prerequisite;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PrerequisiteListValue.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PrerequisiteListValue
  ): PrerequisiteListValue.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PrerequisiteListValue,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): PrerequisiteListValue;
  static deserializeBinaryFromReader(
    message: PrerequisiteListValue,
    reader: jspb.BinaryReader
  ): PrerequisiteListValue;
}

export namespace PrerequisiteListValue {
  export type AsObject = {
    valuesList: Array<Prerequisite.AsObject>;
  };
}
