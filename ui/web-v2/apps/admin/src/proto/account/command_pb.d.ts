// package: bucketeer.account
// file: proto/account/command.proto

import * as jspb from "google-protobuf";
import * as proto_account_account_pb from "../../proto/account/account_pb";
import * as proto_account_api_key_pb from "../../proto/account/api_key_pb";

export class CreateAdminAccountCommand extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAdminAccountCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAdminAccountCommand): CreateAdminAccountCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAdminAccountCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAdminAccountCommand;
  static deserializeBinaryFromReader(message: CreateAdminAccountCommand, reader: jspb.BinaryReader): CreateAdminAccountCommand;
}

export namespace CreateAdminAccountCommand {
  export type AsObject = {
    email: string,
  }
}

export class EnableAdminAccountCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAdminAccountCommand.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAdminAccountCommand): EnableAdminAccountCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAdminAccountCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAdminAccountCommand;
  static deserializeBinaryFromReader(message: EnableAdminAccountCommand, reader: jspb.BinaryReader): EnableAdminAccountCommand;
}

export namespace EnableAdminAccountCommand {
  export type AsObject = {
  }
}

export class DisableAdminAccountCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAdminAccountCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAdminAccountCommand): DisableAdminAccountCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAdminAccountCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAdminAccountCommand;
  static deserializeBinaryFromReader(message: DisableAdminAccountCommand, reader: jspb.BinaryReader): DisableAdminAccountCommand;
}

export namespace DisableAdminAccountCommand {
  export type AsObject = {
  }
}

export class ConvertAccountCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConvertAccountCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ConvertAccountCommand): ConvertAccountCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConvertAccountCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConvertAccountCommand;
  static deserializeBinaryFromReader(message: ConvertAccountCommand, reader: jspb.BinaryReader): ConvertAccountCommand;
}

export namespace ConvertAccountCommand {
  export type AsObject = {
  }
}

export class DeleteAccountCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAccountCommand): DeleteAccountCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAccountCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountCommand;
  static deserializeBinaryFromReader(message: DeleteAccountCommand, reader: jspb.BinaryReader): DeleteAccountCommand;
}

export namespace DeleteAccountCommand {
  export type AsObject = {
  }
}

export class CreateAccountCommand extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getRole(): proto_account_account_pb.Account.RoleMap[keyof proto_account_account_pb.Account.RoleMap];
  setRole(value: proto_account_account_pb.Account.RoleMap[keyof proto_account_account_pb.Account.RoleMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountCommand): CreateAccountCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAccountCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountCommand;
  static deserializeBinaryFromReader(message: CreateAccountCommand, reader: jspb.BinaryReader): CreateAccountCommand;
}

export namespace CreateAccountCommand {
  export type AsObject = {
    email: string,
    role: proto_account_account_pb.Account.RoleMap[keyof proto_account_account_pb.Account.RoleMap],
  }
}

export class ChangeAccountRoleCommand extends jspb.Message {
  getRole(): proto_account_account_pb.Account.RoleMap[keyof proto_account_account_pb.Account.RoleMap];
  setRole(value: proto_account_account_pb.Account.RoleMap[keyof proto_account_account_pb.Account.RoleMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountRoleCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAccountRoleCommand): ChangeAccountRoleCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAccountRoleCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountRoleCommand;
  static deserializeBinaryFromReader(message: ChangeAccountRoleCommand, reader: jspb.BinaryReader): ChangeAccountRoleCommand;
}

export namespace ChangeAccountRoleCommand {
  export type AsObject = {
    role: proto_account_account_pb.Account.RoleMap[keyof proto_account_account_pb.Account.RoleMap],
  }
}

export class EnableAccountCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAccountCommand.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAccountCommand): EnableAccountCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAccountCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAccountCommand;
  static deserializeBinaryFromReader(message: EnableAccountCommand, reader: jspb.BinaryReader): EnableAccountCommand;
}

export namespace EnableAccountCommand {
  export type AsObject = {
  }
}

export class DisableAccountCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAccountCommand): DisableAccountCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAccountCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountCommand;
  static deserializeBinaryFromReader(message: DisableAccountCommand, reader: jspb.BinaryReader): DisableAccountCommand;
}

export namespace DisableAccountCommand {
  export type AsObject = {
  }
}

export class CreateAPIKeyCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getRole(): proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap];
  setRole(value: proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAPIKeyCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAPIKeyCommand): CreateAPIKeyCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAPIKeyCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAPIKeyCommand;
  static deserializeBinaryFromReader(message: CreateAPIKeyCommand, reader: jspb.BinaryReader): CreateAPIKeyCommand;
}

export namespace CreateAPIKeyCommand {
  export type AsObject = {
    name: string,
    role: proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap],
  }
}

export class ChangeAPIKeyNameCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAPIKeyNameCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAPIKeyNameCommand): ChangeAPIKeyNameCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAPIKeyNameCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAPIKeyNameCommand;
  static deserializeBinaryFromReader(message: ChangeAPIKeyNameCommand, reader: jspb.BinaryReader): ChangeAPIKeyNameCommand;
}

export namespace ChangeAPIKeyNameCommand {
  export type AsObject = {
    name: string,
  }
}

export class EnableAPIKeyCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAPIKeyCommand.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAPIKeyCommand): EnableAPIKeyCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAPIKeyCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAPIKeyCommand;
  static deserializeBinaryFromReader(message: EnableAPIKeyCommand, reader: jspb.BinaryReader): EnableAPIKeyCommand;
}

export namespace EnableAPIKeyCommand {
  export type AsObject = {
  }
}

export class DisableAPIKeyCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAPIKeyCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAPIKeyCommand): DisableAPIKeyCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAPIKeyCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAPIKeyCommand;
  static deserializeBinaryFromReader(message: DisableAPIKeyCommand, reader: jspb.BinaryReader): DisableAPIKeyCommand;
}

export namespace DisableAPIKeyCommand {
  export type AsObject = {
  }
}

