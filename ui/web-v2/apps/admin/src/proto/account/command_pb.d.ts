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

export class CreateAccountV2Command extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getName(): string;
  setName(value: string): void;

  getAvatarImageUrl(): string;
  setAvatarImageUrl(value: string): void;

  getOrganizationRole(): proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
  setOrganizationRole(value: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap]): void;

  clearEnvironmentRolesList(): void;
  getEnvironmentRolesList(): Array<proto_account_account_pb.AccountV2.EnvironmentRole>;
  setEnvironmentRolesList(value: Array<proto_account_account_pb.AccountV2.EnvironmentRole>): void;
  addEnvironmentRoles(value?: proto_account_account_pb.AccountV2.EnvironmentRole, index?: number): proto_account_account_pb.AccountV2.EnvironmentRole;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountV2Command.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountV2Command): CreateAccountV2Command.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAccountV2Command, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountV2Command;
  static deserializeBinaryFromReader(message: CreateAccountV2Command, reader: jspb.BinaryReader): CreateAccountV2Command;
}

export namespace CreateAccountV2Command {
  export type AsObject = {
    email: string,
    name: string,
    avatarImageUrl: string,
    organizationRole: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap],
    environmentRolesList: Array<proto_account_account_pb.AccountV2.EnvironmentRole.AsObject>,
  }
}

export class ChangeAccountV2NameCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountV2NameCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAccountV2NameCommand): ChangeAccountV2NameCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAccountV2NameCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountV2NameCommand;
  static deserializeBinaryFromReader(message: ChangeAccountV2NameCommand, reader: jspb.BinaryReader): ChangeAccountV2NameCommand;
}

export namespace ChangeAccountV2NameCommand {
  export type AsObject = {
    name: string,
  }
}

export class ChangeAccountV2AvatarImageUrlCommand extends jspb.Message {
  getAvatarImageUrl(): string;
  setAvatarImageUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountV2AvatarImageUrlCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAccountV2AvatarImageUrlCommand): ChangeAccountV2AvatarImageUrlCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAccountV2AvatarImageUrlCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountV2AvatarImageUrlCommand;
  static deserializeBinaryFromReader(message: ChangeAccountV2AvatarImageUrlCommand, reader: jspb.BinaryReader): ChangeAccountV2AvatarImageUrlCommand;
}

export namespace ChangeAccountV2AvatarImageUrlCommand {
  export type AsObject = {
    avatarImageUrl: string,
  }
}

export class ChangeAccountV2OrganizationRoleCommand extends jspb.Message {
  getRole(): proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
  setRole(value: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountV2OrganizationRoleCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAccountV2OrganizationRoleCommand): ChangeAccountV2OrganizationRoleCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAccountV2OrganizationRoleCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountV2OrganizationRoleCommand;
  static deserializeBinaryFromReader(message: ChangeAccountV2OrganizationRoleCommand, reader: jspb.BinaryReader): ChangeAccountV2OrganizationRoleCommand;
}

export namespace ChangeAccountV2OrganizationRoleCommand {
  export type AsObject = {
    role: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap],
  }
}

export class ChangeAccountV2EnvironmentRolesCommand extends jspb.Message {
  clearRolesList(): void;
  getRolesList(): Array<proto_account_account_pb.AccountV2.EnvironmentRole>;
  setRolesList(value: Array<proto_account_account_pb.AccountV2.EnvironmentRole>): void;
  addRoles(value?: proto_account_account_pb.AccountV2.EnvironmentRole, index?: number): proto_account_account_pb.AccountV2.EnvironmentRole;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountV2EnvironmentRolesCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAccountV2EnvironmentRolesCommand): ChangeAccountV2EnvironmentRolesCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAccountV2EnvironmentRolesCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountV2EnvironmentRolesCommand;
  static deserializeBinaryFromReader(message: ChangeAccountV2EnvironmentRolesCommand, reader: jspb.BinaryReader): ChangeAccountV2EnvironmentRolesCommand;
}

export namespace ChangeAccountV2EnvironmentRolesCommand {
  export type AsObject = {
    rolesList: Array<proto_account_account_pb.AccountV2.EnvironmentRole.AsObject>,
  }
}

export class EnableAccountV2Command extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAccountV2Command.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAccountV2Command): EnableAccountV2Command.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAccountV2Command, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAccountV2Command;
  static deserializeBinaryFromReader(message: EnableAccountV2Command, reader: jspb.BinaryReader): EnableAccountV2Command;
}

export namespace EnableAccountV2Command {
  export type AsObject = {
  }
}

export class DisableAccountV2Command extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountV2Command.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAccountV2Command): DisableAccountV2Command.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAccountV2Command, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountV2Command;
  static deserializeBinaryFromReader(message: DisableAccountV2Command, reader: jspb.BinaryReader): DisableAccountV2Command;
}

export namespace DisableAccountV2Command {
  export type AsObject = {
  }
}

export class DeleteAccountV2Command extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountV2Command.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAccountV2Command): DeleteAccountV2Command.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAccountV2Command, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountV2Command;
  static deserializeBinaryFromReader(message: DeleteAccountV2Command, reader: jspb.BinaryReader): DeleteAccountV2Command;
}

export namespace DeleteAccountV2Command {
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

