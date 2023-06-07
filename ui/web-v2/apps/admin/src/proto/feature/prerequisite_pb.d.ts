// package: bucketeer.feature
// file: proto/feature/prerequisite.proto

import * as jspb from "google-protobuf";

export class Prerequisite extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Prerequisite.AsObject;
  static toObject(includeInstance: boolean, msg: Prerequisite): Prerequisite.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Prerequisite, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Prerequisite;
  static deserializeBinaryFromReader(message: Prerequisite, reader: jspb.BinaryReader): Prerequisite;
}

export namespace Prerequisite {
  export type AsObject = {
    featureId: string,
    variationId: string,
  }
}

