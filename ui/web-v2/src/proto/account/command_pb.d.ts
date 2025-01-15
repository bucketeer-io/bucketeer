// package: bucketeer.account
// file: proto/account/command.proto

import * as jspb from 'google-protobuf';
import * as proto_account_account_pb from '../../proto/account/account_pb';
import * as proto_account_api_key_pb from '../../proto/account/api_key_pb';
import * as proto_account_search_filter_pb from '../../proto/account/search_filter_pb';

export class CreateAccountV2Command extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getName(): string;
  setName(value: string): void;

  getAvatarImageUrl(): string;
  setAvatarImageUrl(value: string): void;

  getOrganizationRole(): proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
  setOrganizationRole(
    value: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap]
  ): void;

  clearEnvironmentRolesList(): void;
  getEnvironmentRolesList(): Array<proto_account_account_pb.AccountV2.EnvironmentRole>;
  setEnvironmentRolesList(
    value: Array<proto_account_account_pb.AccountV2.EnvironmentRole>
  ): void;
  addEnvironmentRoles(
    value?: proto_account_account_pb.AccountV2.EnvironmentRole,
    index?: number
  ): proto_account_account_pb.AccountV2.EnvironmentRole;

  getFirstName(): string;
  setFirstName(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  getLanguage(): string;
  setLanguage(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountV2Command.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateAccountV2Command
  ): CreateAccountV2Command.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateAccountV2Command,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountV2Command;
  static deserializeBinaryFromReader(
    message: CreateAccountV2Command,
    reader: jspb.BinaryReader
  ): CreateAccountV2Command;
}

export namespace CreateAccountV2Command {
  export type AsObject = {
    email: string;
    name: string;
    avatarImageUrl: string;
    organizationRole: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
    environmentRolesList: Array<proto_account_account_pb.AccountV2.EnvironmentRole.AsObject>;
    firstName: string;
    lastName: string;
    language: string;
    tagsList: Array<string>;
  };
}

export class ChangeAccountV2NameCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountV2NameCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAccountV2NameCommand
  ): ChangeAccountV2NameCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAccountV2NameCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountV2NameCommand;
  static deserializeBinaryFromReader(
    message: ChangeAccountV2NameCommand,
    reader: jspb.BinaryReader
  ): ChangeAccountV2NameCommand;
}

export namespace ChangeAccountV2NameCommand {
  export type AsObject = {
    name: string;
  };
}

export class ChangeAccountV2FirstNameCommand extends jspb.Message {
  getFirstName(): string;
  setFirstName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountV2FirstNameCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAccountV2FirstNameCommand
  ): ChangeAccountV2FirstNameCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAccountV2FirstNameCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountV2FirstNameCommand;
  static deserializeBinaryFromReader(
    message: ChangeAccountV2FirstNameCommand,
    reader: jspb.BinaryReader
  ): ChangeAccountV2FirstNameCommand;
}

export namespace ChangeAccountV2FirstNameCommand {
  export type AsObject = {
    firstName: string;
  };
}

export class ChangeAccountV2LastNameCommand extends jspb.Message {
  getLastName(): string;
  setLastName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountV2LastNameCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAccountV2LastNameCommand
  ): ChangeAccountV2LastNameCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAccountV2LastNameCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountV2LastNameCommand;
  static deserializeBinaryFromReader(
    message: ChangeAccountV2LastNameCommand,
    reader: jspb.BinaryReader
  ): ChangeAccountV2LastNameCommand;
}

export namespace ChangeAccountV2LastNameCommand {
  export type AsObject = {
    lastName: string;
  };
}

export class ChangeAccountV2LanguageCommand extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountV2LanguageCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAccountV2LanguageCommand
  ): ChangeAccountV2LanguageCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAccountV2LanguageCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountV2LanguageCommand;
  static deserializeBinaryFromReader(
    message: ChangeAccountV2LanguageCommand,
    reader: jspb.BinaryReader
  ): ChangeAccountV2LanguageCommand;
}

export namespace ChangeAccountV2LanguageCommand {
  export type AsObject = {
    language: string;
  };
}

export class ChangeAccountV2AvatarImageUrlCommand extends jspb.Message {
  getAvatarImageUrl(): string;
  setAvatarImageUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ChangeAccountV2AvatarImageUrlCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAccountV2AvatarImageUrlCommand
  ): ChangeAccountV2AvatarImageUrlCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAccountV2AvatarImageUrlCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ChangeAccountV2AvatarImageUrlCommand;
  static deserializeBinaryFromReader(
    message: ChangeAccountV2AvatarImageUrlCommand,
    reader: jspb.BinaryReader
  ): ChangeAccountV2AvatarImageUrlCommand;
}

export namespace ChangeAccountV2AvatarImageUrlCommand {
  export type AsObject = {
    avatarImageUrl: string;
  };
}

export class ChangeAccountV2AvatarCommand extends jspb.Message {
  getAvatarImage(): Uint8Array | string;
  getAvatarImage_asU8(): Uint8Array;
  getAvatarImage_asB64(): string;
  setAvatarImage(value: Uint8Array | string): void;

  getAvatarFileType(): string;
  setAvatarFileType(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountV2AvatarCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAccountV2AvatarCommand
  ): ChangeAccountV2AvatarCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAccountV2AvatarCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountV2AvatarCommand;
  static deserializeBinaryFromReader(
    message: ChangeAccountV2AvatarCommand,
    reader: jspb.BinaryReader
  ): ChangeAccountV2AvatarCommand;
}

export namespace ChangeAccountV2AvatarCommand {
  export type AsObject = {
    avatarImage: Uint8Array | string;
    avatarFileType: string;
  };
}

export class ChangeTagsCommand extends jspb.Message {
  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeTagsCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeTagsCommand
  ): ChangeTagsCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeTagsCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeTagsCommand;
  static deserializeBinaryFromReader(
    message: ChangeTagsCommand,
    reader: jspb.BinaryReader
  ): ChangeTagsCommand;
}

export namespace ChangeTagsCommand {
  export type AsObject = {
    tagsList: Array<string>;
  };
}

export class ChangeAccountV2LastSeenCommand extends jspb.Message {
  getLastSeen(): number;
  setLastSeen(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAccountV2LastSeenCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAccountV2LastSeenCommand
  ): ChangeAccountV2LastSeenCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAccountV2LastSeenCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAccountV2LastSeenCommand;
  static deserializeBinaryFromReader(
    message: ChangeAccountV2LastSeenCommand,
    reader: jspb.BinaryReader
  ): ChangeAccountV2LastSeenCommand;
}

export namespace ChangeAccountV2LastSeenCommand {
  export type AsObject = {
    lastSeen: number;
  };
}

export class ChangeAccountV2OrganizationRoleCommand extends jspb.Message {
  getRole(): proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
  setRole(
    value: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ChangeAccountV2OrganizationRoleCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAccountV2OrganizationRoleCommand
  ): ChangeAccountV2OrganizationRoleCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAccountV2OrganizationRoleCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ChangeAccountV2OrganizationRoleCommand;
  static deserializeBinaryFromReader(
    message: ChangeAccountV2OrganizationRoleCommand,
    reader: jspb.BinaryReader
  ): ChangeAccountV2OrganizationRoleCommand;
}

export namespace ChangeAccountV2OrganizationRoleCommand {
  export type AsObject = {
    role: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
  };
}

export class ChangeAccountV2EnvironmentRolesCommand extends jspb.Message {
  clearRolesList(): void;
  getRolesList(): Array<proto_account_account_pb.AccountV2.EnvironmentRole>;
  setRolesList(
    value: Array<proto_account_account_pb.AccountV2.EnvironmentRole>
  ): void;
  addRoles(
    value?: proto_account_account_pb.AccountV2.EnvironmentRole,
    index?: number
  ): proto_account_account_pb.AccountV2.EnvironmentRole;

  getWriteType(): ChangeAccountV2EnvironmentRolesCommand.WriteTypeMap[keyof ChangeAccountV2EnvironmentRolesCommand.WriteTypeMap];
  setWriteType(
    value: ChangeAccountV2EnvironmentRolesCommand.WriteTypeMap[keyof ChangeAccountV2EnvironmentRolesCommand.WriteTypeMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ChangeAccountV2EnvironmentRolesCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAccountV2EnvironmentRolesCommand
  ): ChangeAccountV2EnvironmentRolesCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAccountV2EnvironmentRolesCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ChangeAccountV2EnvironmentRolesCommand;
  static deserializeBinaryFromReader(
    message: ChangeAccountV2EnvironmentRolesCommand,
    reader: jspb.BinaryReader
  ): ChangeAccountV2EnvironmentRolesCommand;
}

export namespace ChangeAccountV2EnvironmentRolesCommand {
  export type AsObject = {
    rolesList: Array<proto_account_account_pb.AccountV2.EnvironmentRole.AsObject>;
    writeType: ChangeAccountV2EnvironmentRolesCommand.WriteTypeMap[keyof ChangeAccountV2EnvironmentRolesCommand.WriteTypeMap];
  };

  export interface WriteTypeMap {
    WRITETYPE_UNSPECIFIED: 0;
    WRITETYPE_OVERRIDE: 1;
    WRITETYPE_PATCH: 2;
  }

  export const WriteType: WriteTypeMap;
}

export class EnableAccountV2Command extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAccountV2Command.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnableAccountV2Command
  ): EnableAccountV2Command.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnableAccountV2Command,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnableAccountV2Command;
  static deserializeBinaryFromReader(
    message: EnableAccountV2Command,
    reader: jspb.BinaryReader
  ): EnableAccountV2Command;
}

export namespace EnableAccountV2Command {
  export type AsObject = {};
}

export class DisableAccountV2Command extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountV2Command.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DisableAccountV2Command
  ): DisableAccountV2Command.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DisableAccountV2Command,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountV2Command;
  static deserializeBinaryFromReader(
    message: DisableAccountV2Command,
    reader: jspb.BinaryReader
  ): DisableAccountV2Command;
}

export namespace DisableAccountV2Command {
  export type AsObject = {};
}

export class DeleteAccountV2Command extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountV2Command.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteAccountV2Command
  ): DeleteAccountV2Command.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteAccountV2Command,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountV2Command;
  static deserializeBinaryFromReader(
    message: DeleteAccountV2Command,
    reader: jspb.BinaryReader
  ): DeleteAccountV2Command;
}

export namespace DeleteAccountV2Command {
  export type AsObject = {};
}

export class CreateAPIKeyCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getRole(): proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap];
  setRole(
    value: proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAPIKeyCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateAPIKeyCommand
  ): CreateAPIKeyCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateAPIKeyCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateAPIKeyCommand;
  static deserializeBinaryFromReader(
    message: CreateAPIKeyCommand,
    reader: jspb.BinaryReader
  ): CreateAPIKeyCommand;
}

export namespace CreateAPIKeyCommand {
  export type AsObject = {
    name: string;
    role: proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap];
  };
}

export class ChangeAPIKeyNameCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAPIKeyNameCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAPIKeyNameCommand
  ): ChangeAPIKeyNameCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAPIKeyNameCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAPIKeyNameCommand;
  static deserializeBinaryFromReader(
    message: ChangeAPIKeyNameCommand,
    reader: jspb.BinaryReader
  ): ChangeAPIKeyNameCommand;
}

export namespace ChangeAPIKeyNameCommand {
  export type AsObject = {
    name: string;
  };
}

export class EnableAPIKeyCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAPIKeyCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnableAPIKeyCommand
  ): EnableAPIKeyCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnableAPIKeyCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnableAPIKeyCommand;
  static deserializeBinaryFromReader(
    message: EnableAPIKeyCommand,
    reader: jspb.BinaryReader
  ): EnableAPIKeyCommand;
}

export namespace EnableAPIKeyCommand {
  export type AsObject = {};
}

export class DisableAPIKeyCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAPIKeyCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DisableAPIKeyCommand
  ): DisableAPIKeyCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DisableAPIKeyCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DisableAPIKeyCommand;
  static deserializeBinaryFromReader(
    message: DisableAPIKeyCommand,
    reader: jspb.BinaryReader
  ): DisableAPIKeyCommand;
}

export namespace DisableAPIKeyCommand {
  export type AsObject = {};
}

export class CreateSearchFilterCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getFilterTargetType(): proto_account_search_filter_pb.FilterTargetTypeMap[keyof proto_account_search_filter_pb.FilterTargetTypeMap];
  setFilterTargetType(
    value: proto_account_search_filter_pb.FilterTargetTypeMap[keyof proto_account_search_filter_pb.FilterTargetTypeMap]
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getDefaultFilter(): boolean;
  setDefaultFilter(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSearchFilterCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateSearchFilterCommand
  ): CreateSearchFilterCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateSearchFilterCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateSearchFilterCommand;
  static deserializeBinaryFromReader(
    message: CreateSearchFilterCommand,
    reader: jspb.BinaryReader
  ): CreateSearchFilterCommand;
}

export namespace CreateSearchFilterCommand {
  export type AsObject = {
    name: string;
    query: string;
    filterTargetType: proto_account_search_filter_pb.FilterTargetTypeMap[keyof proto_account_search_filter_pb.FilterTargetTypeMap];
    environmentId: string;
    defaultFilter: boolean;
  };
}

export class ChangeSearchFilterNameCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeSearchFilterNameCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeSearchFilterNameCommand
  ): ChangeSearchFilterNameCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeSearchFilterNameCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeSearchFilterNameCommand;
  static deserializeBinaryFromReader(
    message: ChangeSearchFilterNameCommand,
    reader: jspb.BinaryReader
  ): ChangeSearchFilterNameCommand;
}

export namespace ChangeSearchFilterNameCommand {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class ChangeSearchFilterQueryCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeSearchFilterQueryCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeSearchFilterQueryCommand
  ): ChangeSearchFilterQueryCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeSearchFilterQueryCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeSearchFilterQueryCommand;
  static deserializeBinaryFromReader(
    message: ChangeSearchFilterQueryCommand,
    reader: jspb.BinaryReader
  ): ChangeSearchFilterQueryCommand;
}

export namespace ChangeSearchFilterQueryCommand {
  export type AsObject = {
    id: string;
    query: string;
  };
}

export class ChangeDefaultSearchFilterCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDefaultFilter(): boolean;
  setDefaultFilter(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ChangeDefaultSearchFilterCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeDefaultSearchFilterCommand
  ): ChangeDefaultSearchFilterCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeDefaultSearchFilterCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeDefaultSearchFilterCommand;
  static deserializeBinaryFromReader(
    message: ChangeDefaultSearchFilterCommand,
    reader: jspb.BinaryReader
  ): ChangeDefaultSearchFilterCommand;
}

export namespace ChangeDefaultSearchFilterCommand {
  export type AsObject = {
    id: string;
    defaultFilter: boolean;
  };
}

export class DeleteSearchFilterCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSearchFilterCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteSearchFilterCommand
  ): DeleteSearchFilterCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteSearchFilterCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSearchFilterCommand;
  static deserializeBinaryFromReader(
    message: DeleteSearchFilterCommand,
    reader: jspb.BinaryReader
  ): DeleteSearchFilterCommand;
}

export namespace DeleteSearchFilterCommand {
  export type AsObject = {
    id: string;
  };
}
