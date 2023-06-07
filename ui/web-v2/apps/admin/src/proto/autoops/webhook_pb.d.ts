// package: bucketeer.autoops
// file: proto/autoops/webhook.proto

import * as jspb from "google-protobuf";

export class Webhook extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Webhook.AsObject;
  static toObject(includeInstance: boolean, msg: Webhook): Webhook.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Webhook, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Webhook;
  static deserializeBinaryFromReader(message: Webhook, reader: jspb.BinaryReader): Webhook;
}

export namespace Webhook {
  export type AsObject = {
    id: string,
    name: string,
    description: string,
    createdAt: number,
    updatedAt: number,
  }
}

