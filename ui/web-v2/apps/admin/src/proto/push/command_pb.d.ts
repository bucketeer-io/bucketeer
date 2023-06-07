// package: bucketeer.push
// file: proto/push/command.proto

import * as jspb from "google-protobuf";

export class CreatePushCommand extends jspb.Message {
  getFcmApiKey(): string;
  setFcmApiKey(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePushCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreatePushCommand): CreatePushCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreatePushCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreatePushCommand;
  static deserializeBinaryFromReader(message: CreatePushCommand, reader: jspb.BinaryReader): CreatePushCommand;
}

export namespace CreatePushCommand {
  export type AsObject = {
    fcmApiKey: string,
    tagsList: Array<string>,
    name: string,
  }
}

export class AddPushTagsCommand extends jspb.Message {
  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddPushTagsCommand.AsObject;
  static toObject(includeInstance: boolean, msg: AddPushTagsCommand): AddPushTagsCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddPushTagsCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddPushTagsCommand;
  static deserializeBinaryFromReader(message: AddPushTagsCommand, reader: jspb.BinaryReader): AddPushTagsCommand;
}

export namespace AddPushTagsCommand {
  export type AsObject = {
    tagsList: Array<string>,
  }
}

export class DeletePushTagsCommand extends jspb.Message {
  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePushTagsCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePushTagsCommand): DeletePushTagsCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeletePushTagsCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePushTagsCommand;
  static deserializeBinaryFromReader(message: DeletePushTagsCommand, reader: jspb.BinaryReader): DeletePushTagsCommand;
}

export namespace DeletePushTagsCommand {
  export type AsObject = {
    tagsList: Array<string>,
  }
}

export class DeletePushCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePushCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePushCommand): DeletePushCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeletePushCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePushCommand;
  static deserializeBinaryFromReader(message: DeletePushCommand, reader: jspb.BinaryReader): DeletePushCommand;
}

export namespace DeletePushCommand {
  export type AsObject = {
  }
}

export class RenamePushCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenamePushCommand.AsObject;
  static toObject(includeInstance: boolean, msg: RenamePushCommand): RenamePushCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RenamePushCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenamePushCommand;
  static deserializeBinaryFromReader(message: RenamePushCommand, reader: jspb.BinaryReader): RenamePushCommand;
}

export namespace RenamePushCommand {
  export type AsObject = {
    name: string,
  }
}

