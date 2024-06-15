// package: bucketeer.notification
// file: proto/notification/recipient.proto

import * as jspb from "google-protobuf";

export class Recipient extends jspb.Message {
  getType(): Recipient.TypeMap[keyof Recipient.TypeMap];
  setType(value: Recipient.TypeMap[keyof Recipient.TypeMap]): void;

  hasSlackChannelRecipient(): boolean;
  clearSlackChannelRecipient(): void;
  getSlackChannelRecipient(): SlackChannelRecipient | undefined;
  setSlackChannelRecipient(value?: SlackChannelRecipient): void;

  getLanguage(): Recipient.LanguageMap[keyof Recipient.LanguageMap];
  setLanguage(value: Recipient.LanguageMap[keyof Recipient.LanguageMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Recipient.AsObject;
  static toObject(includeInstance: boolean, msg: Recipient): Recipient.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Recipient, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Recipient;
  static deserializeBinaryFromReader(message: Recipient, reader: jspb.BinaryReader): Recipient;
}

export namespace Recipient {
  export type AsObject = {
    type: Recipient.TypeMap[keyof Recipient.TypeMap],
    slackChannelRecipient?: SlackChannelRecipient.AsObject,
    language: Recipient.LanguageMap[keyof Recipient.LanguageMap],
  }

  export interface TypeMap {
    SLACKCHANNEL: 0;
  }

  export const Type: TypeMap;

  export interface LanguageMap {
    ENGLISH: 0;
    JAPANESE: 1;
  }

  export const Language: LanguageMap;
}

export class SlackChannelRecipient extends jspb.Message {
  getWebhookUrl(): string;
  setWebhookUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SlackChannelRecipient.AsObject;
  static toObject(includeInstance: boolean, msg: SlackChannelRecipient): SlackChannelRecipient.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SlackChannelRecipient, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SlackChannelRecipient;
  static deserializeBinaryFromReader(message: SlackChannelRecipient, reader: jspb.BinaryReader): SlackChannelRecipient;
}

export namespace SlackChannelRecipient {
  export type AsObject = {
    webhookUrl: string,
  }
}

