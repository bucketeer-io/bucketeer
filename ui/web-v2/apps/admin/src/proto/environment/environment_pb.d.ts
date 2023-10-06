// package: bucketeer.environment
// file: proto/environment/environment.proto

import * as jspb from "google-protobuf";

export class EnvironmentV2 extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getUrlCode(): string;
  setUrlCode(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getArchived(): boolean;
  setArchived(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentV2.AsObject;
  static toObject(includeInstance: boolean, msg: EnvironmentV2): EnvironmentV2.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnvironmentV2, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentV2;
  static deserializeBinaryFromReader(message: EnvironmentV2, reader: jspb.BinaryReader): EnvironmentV2;
}

export namespace EnvironmentV2 {
  export type AsObject = {
    id: string,
    name: string,
    urlCode: string,
    description: string,
    projectId: string,
    archived: boolean,
    createdAt: number,
    updatedAt: number,
  }
}

