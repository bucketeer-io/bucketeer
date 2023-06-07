// package: bucketeer.notification
// file: proto/notification/command.proto

import * as jspb from "google-protobuf";
import * as proto_notification_subscription_pb from "../../proto/notification/subscription_pb";
import * as proto_notification_recipient_pb from "../../proto/notification/recipient_pb";

export class CreateAdminSubscriptionCommand extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap], index?: number): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  hasRecipient(): boolean;
  clearRecipient(): void;
  getRecipient(): proto_notification_recipient_pb.Recipient | undefined;
  setRecipient(value?: proto_notification_recipient_pb.Recipient): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAdminSubscriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAdminSubscriptionCommand): CreateAdminSubscriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAdminSubscriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAdminSubscriptionCommand;
  static deserializeBinaryFromReader(message: CreateAdminSubscriptionCommand, reader: jspb.BinaryReader): CreateAdminSubscriptionCommand;
}

export namespace CreateAdminSubscriptionCommand {
  export type AsObject = {
    sourceTypesList: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>,
    recipient?: proto_notification_recipient_pb.Recipient.AsObject,
    name: string,
  }
}

export class AddAdminSubscriptionSourceTypesCommand extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap], index?: number): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAdminSubscriptionSourceTypesCommand.AsObject;
  static toObject(includeInstance: boolean, msg: AddAdminSubscriptionSourceTypesCommand): AddAdminSubscriptionSourceTypesCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddAdminSubscriptionSourceTypesCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAdminSubscriptionSourceTypesCommand;
  static deserializeBinaryFromReader(message: AddAdminSubscriptionSourceTypesCommand, reader: jspb.BinaryReader): AddAdminSubscriptionSourceTypesCommand;
}

export namespace AddAdminSubscriptionSourceTypesCommand {
  export type AsObject = {
    sourceTypesList: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>,
  }
}

export class DeleteAdminSubscriptionSourceTypesCommand extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap], index?: number): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAdminSubscriptionSourceTypesCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAdminSubscriptionSourceTypesCommand): DeleteAdminSubscriptionSourceTypesCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAdminSubscriptionSourceTypesCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAdminSubscriptionSourceTypesCommand;
  static deserializeBinaryFromReader(message: DeleteAdminSubscriptionSourceTypesCommand, reader: jspb.BinaryReader): DeleteAdminSubscriptionSourceTypesCommand;
}

export namespace DeleteAdminSubscriptionSourceTypesCommand {
  export type AsObject = {
    sourceTypesList: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>,
  }
}

export class EnableAdminSubscriptionCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAdminSubscriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAdminSubscriptionCommand): EnableAdminSubscriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAdminSubscriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAdminSubscriptionCommand;
  static deserializeBinaryFromReader(message: EnableAdminSubscriptionCommand, reader: jspb.BinaryReader): EnableAdminSubscriptionCommand;
}

export namespace EnableAdminSubscriptionCommand {
  export type AsObject = {
  }
}

export class DisableAdminSubscriptionCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAdminSubscriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAdminSubscriptionCommand): DisableAdminSubscriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAdminSubscriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAdminSubscriptionCommand;
  static deserializeBinaryFromReader(message: DisableAdminSubscriptionCommand, reader: jspb.BinaryReader): DisableAdminSubscriptionCommand;
}

export namespace DisableAdminSubscriptionCommand {
  export type AsObject = {
  }
}

export class DeleteAdminSubscriptionCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAdminSubscriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAdminSubscriptionCommand): DeleteAdminSubscriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAdminSubscriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAdminSubscriptionCommand;
  static deserializeBinaryFromReader(message: DeleteAdminSubscriptionCommand, reader: jspb.BinaryReader): DeleteAdminSubscriptionCommand;
}

export namespace DeleteAdminSubscriptionCommand {
  export type AsObject = {
  }
}

export class RenameAdminSubscriptionCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenameAdminSubscriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: RenameAdminSubscriptionCommand): RenameAdminSubscriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RenameAdminSubscriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenameAdminSubscriptionCommand;
  static deserializeBinaryFromReader(message: RenameAdminSubscriptionCommand, reader: jspb.BinaryReader): RenameAdminSubscriptionCommand;
}

export namespace RenameAdminSubscriptionCommand {
  export type AsObject = {
    name: string,
  }
}

export class CreateSubscriptionCommand extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap], index?: number): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  hasRecipient(): boolean;
  clearRecipient(): void;
  getRecipient(): proto_notification_recipient_pb.Recipient | undefined;
  setRecipient(value?: proto_notification_recipient_pb.Recipient): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSubscriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateSubscriptionCommand): CreateSubscriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateSubscriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateSubscriptionCommand;
  static deserializeBinaryFromReader(message: CreateSubscriptionCommand, reader: jspb.BinaryReader): CreateSubscriptionCommand;
}

export namespace CreateSubscriptionCommand {
  export type AsObject = {
    sourceTypesList: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>,
    recipient?: proto_notification_recipient_pb.Recipient.AsObject,
    name: string,
  }
}

export class AddSourceTypesCommand extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap], index?: number): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSourceTypesCommand.AsObject;
  static toObject(includeInstance: boolean, msg: AddSourceTypesCommand): AddSourceTypesCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddSourceTypesCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSourceTypesCommand;
  static deserializeBinaryFromReader(message: AddSourceTypesCommand, reader: jspb.BinaryReader): AddSourceTypesCommand;
}

export namespace AddSourceTypesCommand {
  export type AsObject = {
    sourceTypesList: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>,
  }
}

export class DeleteSourceTypesCommand extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap], index?: number): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSourceTypesCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSourceTypesCommand): DeleteSourceTypesCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteSourceTypesCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSourceTypesCommand;
  static deserializeBinaryFromReader(message: DeleteSourceTypesCommand, reader: jspb.BinaryReader): DeleteSourceTypesCommand;
}

export namespace DeleteSourceTypesCommand {
  export type AsObject = {
    sourceTypesList: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>,
  }
}

export class EnableSubscriptionCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableSubscriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: EnableSubscriptionCommand): EnableSubscriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableSubscriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableSubscriptionCommand;
  static deserializeBinaryFromReader(message: EnableSubscriptionCommand, reader: jspb.BinaryReader): EnableSubscriptionCommand;
}

export namespace EnableSubscriptionCommand {
  export type AsObject = {
  }
}

export class DisableSubscriptionCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableSubscriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DisableSubscriptionCommand): DisableSubscriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableSubscriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableSubscriptionCommand;
  static deserializeBinaryFromReader(message: DisableSubscriptionCommand, reader: jspb.BinaryReader): DisableSubscriptionCommand;
}

export namespace DisableSubscriptionCommand {
  export type AsObject = {
  }
}

export class DeleteSubscriptionCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSubscriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSubscriptionCommand): DeleteSubscriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteSubscriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSubscriptionCommand;
  static deserializeBinaryFromReader(message: DeleteSubscriptionCommand, reader: jspb.BinaryReader): DeleteSubscriptionCommand;
}

export namespace DeleteSubscriptionCommand {
  export type AsObject = {
  }
}

export class RenameSubscriptionCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenameSubscriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: RenameSubscriptionCommand): RenameSubscriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RenameSubscriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenameSubscriptionCommand;
  static deserializeBinaryFromReader(message: RenameSubscriptionCommand, reader: jspb.BinaryReader): RenameSubscriptionCommand;
}

export namespace RenameSubscriptionCommand {
  export type AsObject = {
    name: string,
  }
}

