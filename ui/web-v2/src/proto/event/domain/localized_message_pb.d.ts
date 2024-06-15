// package: bucketeer.event.domain
// file: proto/event/domain/localized_message.proto

import * as jspb from "google-protobuf";

export class LocalizedMessage extends jspb.Message {
  getLocale(): string;
  setLocale(value: string): void;

  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LocalizedMessage.AsObject;
  static toObject(includeInstance: boolean, msg: LocalizedMessage): LocalizedMessage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LocalizedMessage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LocalizedMessage;
  static deserializeBinaryFromReader(message: LocalizedMessage, reader: jspb.BinaryReader): LocalizedMessage;
}

export namespace LocalizedMessage {
  export type AsObject = {
    locale: string,
    message: string,
  }
}

