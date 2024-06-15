// package: bucketeer.push
// file: proto/push/push.proto

import * as jspb from "google-protobuf";

export class Push extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFcmApiKey(): string;
  setFcmApiKey(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  getName(): string;
  setName(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Push.AsObject;
  static toObject(includeInstance: boolean, msg: Push): Push.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Push, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Push;
  static deserializeBinaryFromReader(message: Push, reader: jspb.BinaryReader): Push;
}

export namespace Push {
  export type AsObject = {
    id: string,
    fcmApiKey: string,
    tagsList: Array<string>,
    deleted: boolean,
    name: string,
    createdAt: number,
    updatedAt: number,
  }
}

