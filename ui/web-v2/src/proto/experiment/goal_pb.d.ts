// package: bucketeer.experiment
// file: proto/experiment/goal.proto

import * as jspb from "google-protobuf";

export class Goal extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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

  getIsInUseStatus(): boolean;
  setIsInUseStatus(value: boolean): void;

  getArchived(): boolean;
  setArchived(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Goal.AsObject;
  static toObject(includeInstance: boolean, msg: Goal): Goal.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Goal, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Goal;
  static deserializeBinaryFromReader(message: Goal, reader: jspb.BinaryReader): Goal;
}

export namespace Goal {
  export type AsObject = {
    id: string,
    name: string,
    description: string,
    deleted: boolean,
    createdAt: number,
    updatedAt: number,
    isInUseStatus: boolean,
    archived: boolean,
  }
}

