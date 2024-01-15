// package: bucketeer.environment
// file: proto/environment/organization.proto

import * as jspb from "google-protobuf";

export class Organization extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getUrlCode(): string;
  setUrlCode(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getArchived(): boolean;
  setArchived(value: boolean): void;

  getTrial(): boolean;
  setTrial(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getSystemAdmin(): boolean;
  setSystemAdmin(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Organization.AsObject;
  static toObject(includeInstance: boolean, msg: Organization): Organization.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Organization, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Organization;
  static deserializeBinaryFromReader(message: Organization, reader: jspb.BinaryReader): Organization;
}

export namespace Organization {
  export type AsObject = {
    id: string,
    name: string,
    urlCode: string,
    description: string,
    disabled: boolean,
    archived: boolean,
    trial: boolean,
    createdAt: number,
    updatedAt: number,
    systemAdmin: boolean,
  }
}

