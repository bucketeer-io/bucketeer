// package: bucketeer.feature
// file: proto/feature/variation.proto

import * as jspb from "google-protobuf";

export class Variation extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Variation.AsObject;
  static toObject(includeInstance: boolean, msg: Variation): Variation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Variation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Variation;
  static deserializeBinaryFromReader(message: Variation, reader: jspb.BinaryReader): Variation;
}

export namespace Variation {
  export type AsObject = {
    id: string,
    value: string,
    name: string,
    description: string,
  }
}

