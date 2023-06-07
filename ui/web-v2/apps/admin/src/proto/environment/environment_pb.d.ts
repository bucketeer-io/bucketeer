// package: bucketeer.environment
// file: proto/environment/environment.proto

import * as jspb from "google-protobuf";

export class Environment extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Environment.AsObject;
  static toObject(includeInstance: boolean, msg: Environment): Environment.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Environment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Environment;
  static deserializeBinaryFromReader(message: Environment, reader: jspb.BinaryReader): Environment;
}

export namespace Environment {
  export type AsObject = {
    id: string,
    namespace: string,
    name: string,
    description: string,
    deleted: boolean,
    createdAt: number,
    updatedAt: number,
    projectId: string,
  }
}

