// package: bucketeer.feature
// file: proto/feature/feature_last_used_info.proto

import * as jspb from "google-protobuf";

export class FeatureLastUsedInfo extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getVersion(): number;
  setVersion(value: number): void;

  getLastUsedAt(): number;
  setLastUsedAt(value: number): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getClientOldestVersion(): string;
  setClientOldestVersion(value: string): void;

  getClientLatestVersion(): string;
  setClientLatestVersion(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureLastUsedInfo.AsObject;
  static toObject(includeInstance: boolean, msg: FeatureLastUsedInfo): FeatureLastUsedInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FeatureLastUsedInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FeatureLastUsedInfo;
  static deserializeBinaryFromReader(message: FeatureLastUsedInfo, reader: jspb.BinaryReader): FeatureLastUsedInfo;
}

export namespace FeatureLastUsedInfo {
  export type AsObject = {
    featureId: string,
    version: number,
    lastUsedAt: number,
    createdAt: number,
    clientOldestVersion: string,
    clientLatestVersion: string,
  }
}

