// package: bucketeer.environment
// file: proto/environment/project.proto

import * as jspb from "google-protobuf";

export class Project extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getTrial(): boolean;
  setTrial(value: boolean): void;

  getCreatorEmail(): string;
  setCreatorEmail(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getName(): string;
  setName(value: string): void;

  getUrlCode(): string;
  setUrlCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Project.AsObject;
  static toObject(includeInstance: boolean, msg: Project): Project.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Project, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Project;
  static deserializeBinaryFromReader(message: Project, reader: jspb.BinaryReader): Project;
}

export namespace Project {
  export type AsObject = {
    id: string,
    description: string,
    disabled: boolean,
    trial: boolean,
    creatorEmail: string,
    createdAt: number,
    updatedAt: number,
    name: string,
    urlCode: string,
  }
}

