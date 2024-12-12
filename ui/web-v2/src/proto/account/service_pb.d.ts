// package: bucketeer.account
// file: proto/account/service.proto

import * as jspb from 'google-protobuf';
import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';
import * as google_api_annotations_pb from '../../google/api/annotations_pb';
import * as google_api_field_behavior_pb from '../../google/api/field_behavior_pb';
import * as protoc_gen_openapiv2_options_annotations_pb from '../../protoc-gen-openapiv2/options/annotations_pb';
import * as proto_account_account_pb from '../../proto/account/account_pb';
import * as proto_account_api_key_pb from '../../proto/account/api_key_pb';
import * as proto_account_command_pb from '../../proto/account/command_pb';
import * as proto_environment_organization_pb from '../../proto/environment/organization_pb';

export class GetMeRequest extends jspb.Message {
  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMeRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetMeRequest
  ): GetMeRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetMeRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetMeRequest;
  static deserializeBinaryFromReader(
    message: GetMeRequest,
    reader: jspb.BinaryReader
  ): GetMeRequest;
}

export namespace GetMeRequest {
  export type AsObject = {
    organizationId: string;
  };
}

export class GetMeResponse extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.ConsoleAccount | undefined;
  setAccount(value?: proto_account_account_pb.ConsoleAccount): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMeResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetMeResponse
  ): GetMeResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetMeResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetMeResponse;
  static deserializeBinaryFromReader(
    message: GetMeResponse,
    reader: jspb.BinaryReader
  ): GetMeResponse;
}

export namespace GetMeResponse {
  export type AsObject = {
    account?: proto_account_account_pb.ConsoleAccount.AsObject;
  };
}

export class GetMyOrganizationsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyOrganizationsRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetMyOrganizationsRequest
  ): GetMyOrganizationsRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetMyOrganizationsRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetMyOrganizationsRequest;
  static deserializeBinaryFromReader(
    message: GetMyOrganizationsRequest,
    reader: jspb.BinaryReader
  ): GetMyOrganizationsRequest;
}

export namespace GetMyOrganizationsRequest {
  export type AsObject = {};
}

export class GetMyOrganizationsByEmailRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): GetMyOrganizationsByEmailRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetMyOrganizationsByEmailRequest
  ): GetMyOrganizationsByEmailRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetMyOrganizationsByEmailRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetMyOrganizationsByEmailRequest;
  static deserializeBinaryFromReader(
    message: GetMyOrganizationsByEmailRequest,
    reader: jspb.BinaryReader
  ): GetMyOrganizationsByEmailRequest;
}

export namespace GetMyOrganizationsByEmailRequest {
  export type AsObject = {
    email: string;
  };
}

export class GetMyOrganizationsResponse extends jspb.Message {
  clearOrganizationsList(): void;
  getOrganizationsList(): Array<proto_environment_organization_pb.Organization>;
  setOrganizationsList(
    value: Array<proto_environment_organization_pb.Organization>
  ): void;
  addOrganizations(
    value?: proto_environment_organization_pb.Organization,
    index?: number
  ): proto_environment_organization_pb.Organization;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyOrganizationsResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetMyOrganizationsResponse
  ): GetMyOrganizationsResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetMyOrganizationsResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetMyOrganizationsResponse;
  static deserializeBinaryFromReader(
    message: GetMyOrganizationsResponse,
    reader: jspb.BinaryReader
  ): GetMyOrganizationsResponse;
}

export namespace GetMyOrganizationsResponse {
  export type AsObject = {
    organizationsList: Array<proto_environment_organization_pb.Organization.AsObject>;
  };
}

export class CreateAccountV2Request extends jspb.Message {
  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.CreateAccountV2Command | undefined;
  setCommand(value?: proto_account_command_pb.CreateAccountV2Command): void;

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

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountV2Request.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateAccountV2Request
  ): CreateAccountV2Request.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateAccountV2Request,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountV2Request;
  static deserializeBinaryFromReader(
    message: CreateAccountV2Request,
    reader: jspb.BinaryReader
  ): CreateAccountV2Request;
}

export namespace CreateAccountV2Request {
  export type AsObject = {
    organizationId: string;
    command?: proto_account_command_pb.CreateAccountV2Command.AsObject;
    email: string;
    name: string;
    avatarImageUrl: string;
    organizationRole: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
    environmentRolesList: Array<proto_account_account_pb.AccountV2.EnvironmentRole.AsObject>;
    firstName: string;
    lastName: string;
    language: string;
  };
}

export class CreateAccountV2Response extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.AccountV2 | undefined;
  setAccount(value?: proto_account_account_pb.AccountV2): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountV2Response.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateAccountV2Response
  ): CreateAccountV2Response.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateAccountV2Response,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountV2Response;
  static deserializeBinaryFromReader(
    message: CreateAccountV2Response,
    reader: jspb.BinaryReader
  ): CreateAccountV2Response;
}

export namespace CreateAccountV2Response {
  export type AsObject = {
    account?: proto_account_account_pb.AccountV2.AsObject;
  };
}

export class EnableAccountV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.EnableAccountV2Command | undefined;
  setCommand(value?: proto_account_command_pb.EnableAccountV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAccountV2Request.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnableAccountV2Request
  ): EnableAccountV2Request.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnableAccountV2Request,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnableAccountV2Request;
  static deserializeBinaryFromReader(
    message: EnableAccountV2Request,
    reader: jspb.BinaryReader
  ): EnableAccountV2Request;
}

export namespace EnableAccountV2Request {
  export type AsObject = {
    email: string;
    organizationId: string;
    command?: proto_account_command_pb.EnableAccountV2Command.AsObject;
  };
}

export class EnableAccountV2Response extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.AccountV2 | undefined;
  setAccount(value?: proto_account_account_pb.AccountV2): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAccountV2Response.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnableAccountV2Response
  ): EnableAccountV2Response.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnableAccountV2Response,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnableAccountV2Response;
  static deserializeBinaryFromReader(
    message: EnableAccountV2Response,
    reader: jspb.BinaryReader
  ): EnableAccountV2Response;
}

export namespace EnableAccountV2Response {
  export type AsObject = {
    account?: proto_account_account_pb.AccountV2.AsObject;
  };
}

export class DisableAccountV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.DisableAccountV2Command | undefined;
  setCommand(value?: proto_account_command_pb.DisableAccountV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountV2Request.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DisableAccountV2Request
  ): DisableAccountV2Request.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DisableAccountV2Request,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountV2Request;
  static deserializeBinaryFromReader(
    message: DisableAccountV2Request,
    reader: jspb.BinaryReader
  ): DisableAccountV2Request;
}

export namespace DisableAccountV2Request {
  export type AsObject = {
    email: string;
    organizationId: string;
    command?: proto_account_command_pb.DisableAccountV2Command.AsObject;
  };
}

export class DisableAccountV2Response extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.AccountV2 | undefined;
  setAccount(value?: proto_account_account_pb.AccountV2): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountV2Response.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DisableAccountV2Response
  ): DisableAccountV2Response.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DisableAccountV2Response,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountV2Response;
  static deserializeBinaryFromReader(
    message: DisableAccountV2Response,
    reader: jspb.BinaryReader
  ): DisableAccountV2Response;
}

export namespace DisableAccountV2Response {
  export type AsObject = {
    account?: proto_account_account_pb.AccountV2.AsObject;
  };
}

export class DeleteAccountV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.DeleteAccountV2Command | undefined;
  setCommand(value?: proto_account_command_pb.DeleteAccountV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountV2Request.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteAccountV2Request
  ): DeleteAccountV2Request.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteAccountV2Request,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountV2Request;
  static deserializeBinaryFromReader(
    message: DeleteAccountV2Request,
    reader: jspb.BinaryReader
  ): DeleteAccountV2Request;
}

export namespace DeleteAccountV2Request {
  export type AsObject = {
    email: string;
    organizationId: string;
    command?: proto_account_command_pb.DeleteAccountV2Command.AsObject;
  };
}

export class DeleteAccountV2Response extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountV2Response.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteAccountV2Response
  ): DeleteAccountV2Response.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteAccountV2Response,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountV2Response;
  static deserializeBinaryFromReader(
    message: DeleteAccountV2Response,
    reader: jspb.BinaryReader
  ): DeleteAccountV2Response;
}

export namespace DeleteAccountV2Response {
  export type AsObject = {};
}

export class UpdateAccountV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasChangeNameCommand(): boolean;
  clearChangeNameCommand(): void;
  getChangeNameCommand():
    | proto_account_command_pb.ChangeAccountV2NameCommand
    | undefined;
  setChangeNameCommand(
    value?: proto_account_command_pb.ChangeAccountV2NameCommand
  ): void;

  hasChangeAvatarUrlCommand(): boolean;
  clearChangeAvatarUrlCommand(): void;
  getChangeAvatarUrlCommand():
    | proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand
    | undefined;
  setChangeAvatarUrlCommand(
    value?: proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand
  ): void;

  hasChangeOrganizationRoleCommand(): boolean;
  clearChangeOrganizationRoleCommand(): void;
  getChangeOrganizationRoleCommand():
    | proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand
    | undefined;
  setChangeOrganizationRoleCommand(
    value?: proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand
  ): void;

  hasChangeEnvironmentRolesCommand(): boolean;
  clearChangeEnvironmentRolesCommand(): void;
  getChangeEnvironmentRolesCommand():
    | proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand
    | undefined;
  setChangeEnvironmentRolesCommand(
    value?: proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand
  ): void;

  hasChangeFirstNameCommand(): boolean;
  clearChangeFirstNameCommand(): void;
  getChangeFirstNameCommand():
    | proto_account_command_pb.ChangeAccountV2FirstNameCommand
    | undefined;
  setChangeFirstNameCommand(
    value?: proto_account_command_pb.ChangeAccountV2FirstNameCommand
  ): void;

  hasChangeLastNameCommand(): boolean;
  clearChangeLastNameCommand(): void;
  getChangeLastNameCommand():
    | proto_account_command_pb.ChangeAccountV2LastNameCommand
    | undefined;
  setChangeLastNameCommand(
    value?: proto_account_command_pb.ChangeAccountV2LastNameCommand
  ): void;

  hasChangeLanguageCommand(): boolean;
  clearChangeLanguageCommand(): void;
  getChangeLanguageCommand():
    | proto_account_command_pb.ChangeAccountV2LanguageCommand
    | undefined;
  setChangeLanguageCommand(
    value?: proto_account_command_pb.ChangeAccountV2LanguageCommand
  ): void;

  hasChangeLastSeenCommand(): boolean;
  clearChangeLastSeenCommand(): void;
  getChangeLastSeenCommand():
    | proto_account_command_pb.ChangeAccountV2LastSeenCommand
    | undefined;
  setChangeLastSeenCommand(
    value?: proto_account_command_pb.ChangeAccountV2LastSeenCommand
  ): void;

  hasChangeAvatarCommand(): boolean;
  clearChangeAvatarCommand(): void;
  getChangeAvatarCommand():
    | proto_account_command_pb.ChangeAccountV2AvatarCommand
    | undefined;
  setChangeAvatarCommand(
    value?: proto_account_command_pb.ChangeAccountV2AvatarCommand
  ): void;

  hasName(): boolean;
  clearName(): void;
  getName(): google_protobuf_wrappers_pb.StringValue | undefined;
  setName(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasAvatarImageUrl(): boolean;
  clearAvatarImageUrl(): void;
  getAvatarImageUrl(): google_protobuf_wrappers_pb.StringValue | undefined;
  setAvatarImageUrl(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasOrganizationRole(): boolean;
  clearOrganizationRole(): void;
  getOrganizationRole():
    | UpdateAccountV2Request.OrganizationRoleValue
    | undefined;
  setOrganizationRole(
    value?: UpdateAccountV2Request.OrganizationRoleValue
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

  hasFirstName(): boolean;
  clearFirstName(): void;
  getFirstName(): google_protobuf_wrappers_pb.StringValue | undefined;
  setFirstName(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasLastName(): boolean;
  clearLastName(): void;
  getLastName(): google_protobuf_wrappers_pb.StringValue | undefined;
  setLastName(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasLanguage(): boolean;
  clearLanguage(): void;
  getLanguage(): google_protobuf_wrappers_pb.StringValue | undefined;
  setLanguage(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasLastSeen(): boolean;
  clearLastSeen(): void;
  getLastSeen(): google_protobuf_wrappers_pb.Int64Value | undefined;
  setLastSeen(value?: google_protobuf_wrappers_pb.Int64Value): void;

  hasAvatar(): boolean;
  clearAvatar(): void;
  getAvatar(): UpdateAccountV2Request.AccountV2Avatar | undefined;
  setAvatar(value?: UpdateAccountV2Request.AccountV2Avatar): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAccountV2Request.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateAccountV2Request
  ): UpdateAccountV2Request.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateAccountV2Request,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAccountV2Request;
  static deserializeBinaryFromReader(
    message: UpdateAccountV2Request,
    reader: jspb.BinaryReader
  ): UpdateAccountV2Request;
}

export namespace UpdateAccountV2Request {
  export type AsObject = {
    email: string;
    organizationId: string;
    changeNameCommand?: proto_account_command_pb.ChangeAccountV2NameCommand.AsObject;
    changeAvatarUrlCommand?: proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand.AsObject;
    changeOrganizationRoleCommand?: proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand.AsObject;
    changeEnvironmentRolesCommand?: proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand.AsObject;
    changeFirstNameCommand?: proto_account_command_pb.ChangeAccountV2FirstNameCommand.AsObject;
    changeLastNameCommand?: proto_account_command_pb.ChangeAccountV2LastNameCommand.AsObject;
    changeLanguageCommand?: proto_account_command_pb.ChangeAccountV2LanguageCommand.AsObject;
    changeLastSeenCommand?: proto_account_command_pb.ChangeAccountV2LastSeenCommand.AsObject;
    changeAvatarCommand?: proto_account_command_pb.ChangeAccountV2AvatarCommand.AsObject;
    name?: google_protobuf_wrappers_pb.StringValue.AsObject;
    avatarImageUrl?: google_protobuf_wrappers_pb.StringValue.AsObject;
    organizationRole?: UpdateAccountV2Request.OrganizationRoleValue.AsObject;
    environmentRolesList: Array<proto_account_account_pb.AccountV2.EnvironmentRole.AsObject>;
    firstName?: google_protobuf_wrappers_pb.StringValue.AsObject;
    lastName?: google_protobuf_wrappers_pb.StringValue.AsObject;
    language?: google_protobuf_wrappers_pb.StringValue.AsObject;
    lastSeen?: google_protobuf_wrappers_pb.Int64Value.AsObject;
    avatar?: UpdateAccountV2Request.AccountV2Avatar.AsObject;
  };

  export class AccountV2Avatar extends jspb.Message {
    getAvatarImage(): Uint8Array | string;
    getAvatarImage_asU8(): Uint8Array;
    getAvatarImage_asB64(): string;
    setAvatarImage(value: Uint8Array | string): void;

    getAvatarFileType(): string;
    setAvatarFileType(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AccountV2Avatar.AsObject;
    static toObject(
      includeInstance: boolean,
      msg: AccountV2Avatar
    ): AccountV2Avatar.AsObject;
    static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
    static extensionsBinary: {
      [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
    };
    static serializeBinaryToWriter(
      message: AccountV2Avatar,
      writer: jspb.BinaryWriter
    ): void;
    static deserializeBinary(bytes: Uint8Array): AccountV2Avatar;
    static deserializeBinaryFromReader(
      message: AccountV2Avatar,
      reader: jspb.BinaryReader
    ): AccountV2Avatar;
  }

  export namespace AccountV2Avatar {
    export type AsObject = {
      avatarImage: Uint8Array | string;
      avatarFileType: string;
    };
  }

  export class OrganizationRoleValue extends jspb.Message {
    getRole(): proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
    setRole(
      value: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap]
    ): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OrganizationRoleValue.AsObject;
    static toObject(
      includeInstance: boolean,
      msg: OrganizationRoleValue
    ): OrganizationRoleValue.AsObject;
    static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
    static extensionsBinary: {
      [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
    };
    static serializeBinaryToWriter(
      message: OrganizationRoleValue,
      writer: jspb.BinaryWriter
    ): void;
    static deserializeBinary(bytes: Uint8Array): OrganizationRoleValue;
    static deserializeBinaryFromReader(
      message: OrganizationRoleValue,
      reader: jspb.BinaryReader
    ): OrganizationRoleValue;
  }

  export namespace OrganizationRoleValue {
    export type AsObject = {
      role: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
    };
  }
}

export class UpdateAccountV2Response extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.AccountV2 | undefined;
  setAccount(value?: proto_account_account_pb.AccountV2): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAccountV2Response.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateAccountV2Response
  ): UpdateAccountV2Response.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateAccountV2Response,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAccountV2Response;
  static deserializeBinaryFromReader(
    message: UpdateAccountV2Response,
    reader: jspb.BinaryReader
  ): UpdateAccountV2Response;
}

export namespace UpdateAccountV2Response {
  export type AsObject = {
    account?: proto_account_account_pb.AccountV2.AsObject;
  };
}

export class GetAccountV2Request extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountV2Request.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAccountV2Request
  ): GetAccountV2Request.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAccountV2Request,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountV2Request;
  static deserializeBinaryFromReader(
    message: GetAccountV2Request,
    reader: jspb.BinaryReader
  ): GetAccountV2Request;
}

export namespace GetAccountV2Request {
  export type AsObject = {
    email: string;
    organizationId: string;
  };
}

export class GetAccountV2Response extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.AccountV2 | undefined;
  setAccount(value?: proto_account_account_pb.AccountV2): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountV2Response.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAccountV2Response
  ): GetAccountV2Response.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAccountV2Response,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountV2Response;
  static deserializeBinaryFromReader(
    message: GetAccountV2Response,
    reader: jspb.BinaryReader
  ): GetAccountV2Response;
}

export namespace GetAccountV2Response {
  export type AsObject = {
    account?: proto_account_account_pb.AccountV2.AsObject;
  };
}

export class GetAccountV2ByEnvironmentIDRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): GetAccountV2ByEnvironmentIDRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAccountV2ByEnvironmentIDRequest
  ): GetAccountV2ByEnvironmentIDRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAccountV2ByEnvironmentIDRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): GetAccountV2ByEnvironmentIDRequest;
  static deserializeBinaryFromReader(
    message: GetAccountV2ByEnvironmentIDRequest,
    reader: jspb.BinaryReader
  ): GetAccountV2ByEnvironmentIDRequest;
}

export namespace GetAccountV2ByEnvironmentIDRequest {
  export type AsObject = {
    email: string;
    environmentId: string;
  };
}

export class GetAccountV2ByEnvironmentIDResponse extends jspb.Message {
  hasAccount(): boolean;
  clearAccount(): void;
  getAccount(): proto_account_account_pb.AccountV2 | undefined;
  setAccount(value?: proto_account_account_pb.AccountV2): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): GetAccountV2ByEnvironmentIDResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAccountV2ByEnvironmentIDResponse
  ): GetAccountV2ByEnvironmentIDResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAccountV2ByEnvironmentIDResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): GetAccountV2ByEnvironmentIDResponse;
  static deserializeBinaryFromReader(
    message: GetAccountV2ByEnvironmentIDResponse,
    reader: jspb.BinaryReader
  ): GetAccountV2ByEnvironmentIDResponse;
}

export namespace GetAccountV2ByEnvironmentIDResponse {
  export type AsObject = {
    account?: proto_account_account_pb.AccountV2.AsObject;
  };
}

export class ListAccountsV2Request extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  getOrderBy(): ListAccountsV2Request.OrderByMap[keyof ListAccountsV2Request.OrderByMap];
  setOrderBy(
    value: ListAccountsV2Request.OrderByMap[keyof ListAccountsV2Request.OrderByMap]
  ): void;

  getOrderDirection(): ListAccountsV2Request.OrderDirectionMap[keyof ListAccountsV2Request.OrderDirectionMap];
  setOrderDirection(
    value: ListAccountsV2Request.OrderDirectionMap[keyof ListAccountsV2Request.OrderDirectionMap]
  ): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasOrganizationRole(): boolean;
  clearOrganizationRole(): void;
  getOrganizationRole(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setOrganizationRole(value?: google_protobuf_wrappers_pb.Int32Value): void;

  hasEnvironmentId(): boolean;
  clearEnvironmentId(): void;
  getEnvironmentId(): google_protobuf_wrappers_pb.StringValue | undefined;
  setEnvironmentId(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasEnvironmentRole(): boolean;
  clearEnvironmentRole(): void;
  getEnvironmentRole(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setEnvironmentRole(value?: google_protobuf_wrappers_pb.Int32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAccountsV2Request.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListAccountsV2Request
  ): ListAccountsV2Request.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListAccountsV2Request,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListAccountsV2Request;
  static deserializeBinaryFromReader(
    message: ListAccountsV2Request,
    reader: jspb.BinaryReader
  ): ListAccountsV2Request;
}

export namespace ListAccountsV2Request {
  export type AsObject = {
    pageSize: number;
    cursor: string;
    organizationId: string;
    orderBy: ListAccountsV2Request.OrderByMap[keyof ListAccountsV2Request.OrderByMap];
    orderDirection: ListAccountsV2Request.OrderDirectionMap[keyof ListAccountsV2Request.OrderDirectionMap];
    searchKeyword: string;
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject;
    organizationRole?: google_protobuf_wrappers_pb.Int32Value.AsObject;
    environmentId?: google_protobuf_wrappers_pb.StringValue.AsObject;
    environmentRole?: google_protobuf_wrappers_pb.Int32Value.AsObject;
  };

  export interface OrderByMap {
    DEFAULT: 0;
    EMAIL: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
    ORGANIZATION_ROLE: 4;
    ENVIRONMENT_COUNT: 5;
    LAST_SEEN: 6;
    STATE: 7;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListAccountsV2Response extends jspb.Message {
  clearAccountsList(): void;
  getAccountsList(): Array<proto_account_account_pb.AccountV2>;
  setAccountsList(value: Array<proto_account_account_pb.AccountV2>): void;
  addAccounts(
    value?: proto_account_account_pb.AccountV2,
    index?: number
  ): proto_account_account_pb.AccountV2;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAccountsV2Response.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListAccountsV2Response
  ): ListAccountsV2Response.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListAccountsV2Response,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListAccountsV2Response;
  static deserializeBinaryFromReader(
    message: ListAccountsV2Response,
    reader: jspb.BinaryReader
  ): ListAccountsV2Response;
}

export namespace ListAccountsV2Response {
  export type AsObject = {
    accountsList: Array<proto_account_account_pb.AccountV2.AsObject>;
    cursor: string;
    totalCount: number;
  };
}

export class CreateAPIKeyRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.CreateAPIKeyCommand | undefined;
  setCommand(value?: proto_account_command_pb.CreateAPIKeyCommand): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getRole(): proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap];
  setRole(
    value: proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap]
  ): void;

  getMaintainer(): string;
  setMaintainer(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAPIKeyRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateAPIKeyRequest
  ): CreateAPIKeyRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateAPIKeyRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateAPIKeyRequest;
  static deserializeBinaryFromReader(
    message: CreateAPIKeyRequest,
    reader: jspb.BinaryReader
  ): CreateAPIKeyRequest;
}

export namespace CreateAPIKeyRequest {
  export type AsObject = {
    command?: proto_account_command_pb.CreateAPIKeyCommand.AsObject;
    environmentId: string;
    name: string;
    role: proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap];
    maintainer: string;
    description: string;
  };
}

export class CreateAPIKeyResponse extends jspb.Message {
  hasApiKey(): boolean;
  clearApiKey(): void;
  getApiKey(): proto_account_api_key_pb.APIKey | undefined;
  setApiKey(value?: proto_account_api_key_pb.APIKey): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAPIKeyResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateAPIKeyResponse
  ): CreateAPIKeyResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateAPIKeyResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateAPIKeyResponse;
  static deserializeBinaryFromReader(
    message: CreateAPIKeyResponse,
    reader: jspb.BinaryReader
  ): CreateAPIKeyResponse;
}

export namespace CreateAPIKeyResponse {
  export type AsObject = {
    apiKey?: proto_account_api_key_pb.APIKey.AsObject;
  };
}

export class ChangeAPIKeyNameRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.ChangeAPIKeyNameCommand | undefined;
  setCommand(value?: proto_account_command_pb.ChangeAPIKeyNameCommand): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAPIKeyNameRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAPIKeyNameRequest
  ): ChangeAPIKeyNameRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAPIKeyNameRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAPIKeyNameRequest;
  static deserializeBinaryFromReader(
    message: ChangeAPIKeyNameRequest,
    reader: jspb.BinaryReader
  ): ChangeAPIKeyNameRequest;
}

export namespace ChangeAPIKeyNameRequest {
  export type AsObject = {
    id: string;
    command?: proto_account_command_pb.ChangeAPIKeyNameCommand.AsObject;
    environmentId: string;
  };
}

export class ChangeAPIKeyNameResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAPIKeyNameResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAPIKeyNameResponse
  ): ChangeAPIKeyNameResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAPIKeyNameResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAPIKeyNameResponse;
  static deserializeBinaryFromReader(
    message: ChangeAPIKeyNameResponse,
    reader: jspb.BinaryReader
  ): ChangeAPIKeyNameResponse;
}

export namespace ChangeAPIKeyNameResponse {
  export type AsObject = {};
}

export class EnableAPIKeyRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.EnableAPIKeyCommand | undefined;
  setCommand(value?: proto_account_command_pb.EnableAPIKeyCommand): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAPIKeyRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnableAPIKeyRequest
  ): EnableAPIKeyRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnableAPIKeyRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnableAPIKeyRequest;
  static deserializeBinaryFromReader(
    message: EnableAPIKeyRequest,
    reader: jspb.BinaryReader
  ): EnableAPIKeyRequest;
}

export namespace EnableAPIKeyRequest {
  export type AsObject = {
    id: string;
    command?: proto_account_command_pb.EnableAPIKeyCommand.AsObject;
    environmentId: string;
  };
}

export class EnableAPIKeyResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAPIKeyResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnableAPIKeyResponse
  ): EnableAPIKeyResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnableAPIKeyResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnableAPIKeyResponse;
  static deserializeBinaryFromReader(
    message: EnableAPIKeyResponse,
    reader: jspb.BinaryReader
  ): EnableAPIKeyResponse;
}

export namespace EnableAPIKeyResponse {
  export type AsObject = {};
}

export class DisableAPIKeyRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.DisableAPIKeyCommand | undefined;
  setCommand(value?: proto_account_command_pb.DisableAPIKeyCommand): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAPIKeyRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DisableAPIKeyRequest
  ): DisableAPIKeyRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DisableAPIKeyRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DisableAPIKeyRequest;
  static deserializeBinaryFromReader(
    message: DisableAPIKeyRequest,
    reader: jspb.BinaryReader
  ): DisableAPIKeyRequest;
}

export namespace DisableAPIKeyRequest {
  export type AsObject = {
    id: string;
    command?: proto_account_command_pb.DisableAPIKeyCommand.AsObject;
    environmentId: string;
  };
}

export class DisableAPIKeyResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAPIKeyResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DisableAPIKeyResponse
  ): DisableAPIKeyResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DisableAPIKeyResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DisableAPIKeyResponse;
  static deserializeBinaryFromReader(
    message: DisableAPIKeyResponse,
    reader: jspb.BinaryReader
  ): DisableAPIKeyResponse;
}

export namespace DisableAPIKeyResponse {
  export type AsObject = {};
}

export class GetAPIKeyRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAPIKeyRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAPIKeyRequest
  ): GetAPIKeyRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAPIKeyRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetAPIKeyRequest;
  static deserializeBinaryFromReader(
    message: GetAPIKeyRequest,
    reader: jspb.BinaryReader
  ): GetAPIKeyRequest;
}

export namespace GetAPIKeyRequest {
  export type AsObject = {
    id: string;
    environmentId: string;
  };
}

export class GetAPIKeyResponse extends jspb.Message {
  hasApiKey(): boolean;
  clearApiKey(): void;
  getApiKey(): proto_account_api_key_pb.APIKey | undefined;
  setApiKey(value?: proto_account_api_key_pb.APIKey): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAPIKeyResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAPIKeyResponse
  ): GetAPIKeyResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAPIKeyResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetAPIKeyResponse;
  static deserializeBinaryFromReader(
    message: GetAPIKeyResponse,
    reader: jspb.BinaryReader
  ): GetAPIKeyResponse;
}

export namespace GetAPIKeyResponse {
  export type AsObject = {
    apiKey?: proto_account_api_key_pb.APIKey.AsObject;
  };
}

export class ListAPIKeysRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListAPIKeysRequest.OrderByMap[keyof ListAPIKeysRequest.OrderByMap];
  setOrderBy(
    value: ListAPIKeysRequest.OrderByMap[keyof ListAPIKeysRequest.OrderByMap]
  ): void;

  getOrderDirection(): ListAPIKeysRequest.OrderDirectionMap[keyof ListAPIKeysRequest.OrderDirectionMap];
  setOrderDirection(
    value: ListAPIKeysRequest.OrderDirectionMap[keyof ListAPIKeysRequest.OrderDirectionMap]
  ): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  clearEnvironmentIdsList(): void;
  getEnvironmentIdsList(): Array<string>;
  setEnvironmentIdsList(value: Array<string>): void;
  addEnvironmentIds(value: string, index?: number): string;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAPIKeysRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListAPIKeysRequest
  ): ListAPIKeysRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListAPIKeysRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListAPIKeysRequest;
  static deserializeBinaryFromReader(
    message: ListAPIKeysRequest,
    reader: jspb.BinaryReader
  ): ListAPIKeysRequest;
}

export namespace ListAPIKeysRequest {
  export type AsObject = {
    pageSize: number;
    cursor: string;
    orderBy: ListAPIKeysRequest.OrderByMap[keyof ListAPIKeysRequest.OrderByMap];
    orderDirection: ListAPIKeysRequest.OrderDirectionMap[keyof ListAPIKeysRequest.OrderDirectionMap];
    searchKeyword: string;
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject;
    environmentId: string;
    environmentIdsList: Array<string>;
    organizationId: string;
  };

  export interface OrderByMap {
    DEFAULT: 0;
    NAME: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
    ROLE: 4;
    ENVIRONMENT: 5;
    STATE: 6;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListAPIKeysResponse extends jspb.Message {
  clearApiKeysList(): void;
  getApiKeysList(): Array<proto_account_api_key_pb.APIKey>;
  setApiKeysList(value: Array<proto_account_api_key_pb.APIKey>): void;
  addApiKeys(
    value?: proto_account_api_key_pb.APIKey,
    index?: number
  ): proto_account_api_key_pb.APIKey;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAPIKeysResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListAPIKeysResponse
  ): ListAPIKeysResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListAPIKeysResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListAPIKeysResponse;
  static deserializeBinaryFromReader(
    message: ListAPIKeysResponse,
    reader: jspb.BinaryReader
  ): ListAPIKeysResponse;
}

export namespace ListAPIKeysResponse {
  export type AsObject = {
    apiKeysList: Array<proto_account_api_key_pb.APIKey.AsObject>;
    cursor: string;
    totalCount: number;
  };
}

export class GetAPIKeyBySearchingAllEnvironmentsRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getApiKey(): string;
  setApiKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): GetAPIKeyBySearchingAllEnvironmentsRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAPIKeyBySearchingAllEnvironmentsRequest
  ): GetAPIKeyBySearchingAllEnvironmentsRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAPIKeyBySearchingAllEnvironmentsRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): GetAPIKeyBySearchingAllEnvironmentsRequest;
  static deserializeBinaryFromReader(
    message: GetAPIKeyBySearchingAllEnvironmentsRequest,
    reader: jspb.BinaryReader
  ): GetAPIKeyBySearchingAllEnvironmentsRequest;
}

export namespace GetAPIKeyBySearchingAllEnvironmentsRequest {
  export type AsObject = {
    id: string;
    apiKey: string;
  };
}

export class GetAPIKeyBySearchingAllEnvironmentsResponse extends jspb.Message {
  hasEnvironmentApiKey(): boolean;
  clearEnvironmentApiKey(): void;
  getEnvironmentApiKey():
    | proto_account_api_key_pb.EnvironmentAPIKey
    | undefined;
  setEnvironmentApiKey(
    value?: proto_account_api_key_pb.EnvironmentAPIKey
  ): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): GetAPIKeyBySearchingAllEnvironmentsResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAPIKeyBySearchingAllEnvironmentsResponse
  ): GetAPIKeyBySearchingAllEnvironmentsResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAPIKeyBySearchingAllEnvironmentsResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): GetAPIKeyBySearchingAllEnvironmentsResponse;
  static deserializeBinaryFromReader(
    message: GetAPIKeyBySearchingAllEnvironmentsResponse,
    reader: jspb.BinaryReader
  ): GetAPIKeyBySearchingAllEnvironmentsResponse;
}

export namespace GetAPIKeyBySearchingAllEnvironmentsResponse {
  export type AsObject = {
    environmentApiKey?: proto_account_api_key_pb.EnvironmentAPIKey.AsObject;
  };
}

export class GetEnvironmentAPIKeyRequest extends jspb.Message {
  getApiKey(): string;
  setApiKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEnvironmentAPIKeyRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetEnvironmentAPIKeyRequest
  ): GetEnvironmentAPIKeyRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetEnvironmentAPIKeyRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetEnvironmentAPIKeyRequest;
  static deserializeBinaryFromReader(
    message: GetEnvironmentAPIKeyRequest,
    reader: jspb.BinaryReader
  ): GetEnvironmentAPIKeyRequest;
}

export namespace GetEnvironmentAPIKeyRequest {
  export type AsObject = {
    apiKey: string;
  };
}

export class GetEnvironmentAPIKeyResponse extends jspb.Message {
  hasEnvironmentApiKey(): boolean;
  clearEnvironmentApiKey(): void;
  getEnvironmentApiKey():
    | proto_account_api_key_pb.EnvironmentAPIKey
    | undefined;
  setEnvironmentApiKey(
    value?: proto_account_api_key_pb.EnvironmentAPIKey
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEnvironmentAPIKeyResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetEnvironmentAPIKeyResponse
  ): GetEnvironmentAPIKeyResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetEnvironmentAPIKeyResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetEnvironmentAPIKeyResponse;
  static deserializeBinaryFromReader(
    message: GetEnvironmentAPIKeyResponse,
    reader: jspb.BinaryReader
  ): GetEnvironmentAPIKeyResponse;
}

export namespace GetEnvironmentAPIKeyResponse {
  export type AsObject = {
    environmentApiKey?: proto_account_api_key_pb.EnvironmentAPIKey.AsObject;
  };
}

export class CreateSearchFilterRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.CreateSearchFilterCommand | undefined;
  setCommand(value?: proto_account_command_pb.CreateSearchFilterCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSearchFilterRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateSearchFilterRequest
  ): CreateSearchFilterRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateSearchFilterRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateSearchFilterRequest;
  static deserializeBinaryFromReader(
    message: CreateSearchFilterRequest,
    reader: jspb.BinaryReader
  ): CreateSearchFilterRequest;
}

export namespace CreateSearchFilterRequest {
  export type AsObject = {
    email: string;
    organizationId: string;
    environmentId: string;
    command?: proto_account_command_pb.CreateSearchFilterCommand.AsObject;
  };
}

export class CreateSearchFilterResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSearchFilterResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateSearchFilterResponse
  ): CreateSearchFilterResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateSearchFilterResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateSearchFilterResponse;
  static deserializeBinaryFromReader(
    message: CreateSearchFilterResponse,
    reader: jspb.BinaryReader
  ): CreateSearchFilterResponse;
}

export namespace CreateSearchFilterResponse {
  export type AsObject = {};
}

export class UpdateSearchFilterRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  hasChangeNameCommand(): boolean;
  clearChangeNameCommand(): void;
  getChangeNameCommand():
    | proto_account_command_pb.ChangeSearchFilterNameCommand
    | undefined;
  setChangeNameCommand(
    value?: proto_account_command_pb.ChangeSearchFilterNameCommand
  ): void;

  hasChangeQueryCommand(): boolean;
  clearChangeQueryCommand(): void;
  getChangeQueryCommand():
    | proto_account_command_pb.ChangeSearchFilterQueryCommand
    | undefined;
  setChangeQueryCommand(
    value?: proto_account_command_pb.ChangeSearchFilterQueryCommand
  ): void;

  hasChangeDefaultFilterCommand(): boolean;
  clearChangeDefaultFilterCommand(): void;
  getChangeDefaultFilterCommand():
    | proto_account_command_pb.ChangeDefaultSearchFilterCommand
    | undefined;
  setChangeDefaultFilterCommand(
    value?: proto_account_command_pb.ChangeDefaultSearchFilterCommand
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSearchFilterRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateSearchFilterRequest
  ): UpdateSearchFilterRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateSearchFilterRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSearchFilterRequest;
  static deserializeBinaryFromReader(
    message: UpdateSearchFilterRequest,
    reader: jspb.BinaryReader
  ): UpdateSearchFilterRequest;
}

export namespace UpdateSearchFilterRequest {
  export type AsObject = {
    email: string;
    organizationId: string;
    environmentId: string;
    changeNameCommand?: proto_account_command_pb.ChangeSearchFilterNameCommand.AsObject;
    changeQueryCommand?: proto_account_command_pb.ChangeSearchFilterQueryCommand.AsObject;
    changeDefaultFilterCommand?: proto_account_command_pb.ChangeDefaultSearchFilterCommand.AsObject;
  };
}

export class UpdateSearchFilterResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSearchFilterResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateSearchFilterResponse
  ): UpdateSearchFilterResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateSearchFilterResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSearchFilterResponse;
  static deserializeBinaryFromReader(
    message: UpdateSearchFilterResponse,
    reader: jspb.BinaryReader
  ): UpdateSearchFilterResponse;
}

export namespace UpdateSearchFilterResponse {
  export type AsObject = {};
}

export class DeleteSearchFilterRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_account_command_pb.DeleteSearchFilterCommand | undefined;
  setCommand(value?: proto_account_command_pb.DeleteSearchFilterCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSearchFilterRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteSearchFilterRequest
  ): DeleteSearchFilterRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteSearchFilterRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSearchFilterRequest;
  static deserializeBinaryFromReader(
    message: DeleteSearchFilterRequest,
    reader: jspb.BinaryReader
  ): DeleteSearchFilterRequest;
}

export namespace DeleteSearchFilterRequest {
  export type AsObject = {
    email: string;
    organizationId: string;
    environmentId: string;
    command?: proto_account_command_pb.DeleteSearchFilterCommand.AsObject;
  };
}

export class DeleteSearchFilterResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSearchFilterResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteSearchFilterResponse
  ): DeleteSearchFilterResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteSearchFilterResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSearchFilterResponse;
  static deserializeBinaryFromReader(
    message: DeleteSearchFilterResponse,
    reader: jspb.BinaryReader
  ): DeleteSearchFilterResponse;
}

export namespace DeleteSearchFilterResponse {
  export type AsObject = {};
}

export class UpdateAPIKeyRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  hasName(): boolean;
  clearName(): void;
  getName(): google_protobuf_wrappers_pb.StringValue | undefined;
  setName(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasDescription(): boolean;
  clearDescription(): void;
  getDescription(): google_protobuf_wrappers_pb.StringValue | undefined;
  setDescription(value?: google_protobuf_wrappers_pb.StringValue): void;

  getRole(): proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap];
  setRole(
    value: proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap]
  ): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasMaintainer(): boolean;
  clearMaintainer(): void;
  getMaintainer(): google_protobuf_wrappers_pb.StringValue | undefined;
  setMaintainer(value?: google_protobuf_wrappers_pb.StringValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAPIKeyRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateAPIKeyRequest
  ): UpdateAPIKeyRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateAPIKeyRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAPIKeyRequest;
  static deserializeBinaryFromReader(
    message: UpdateAPIKeyRequest,
    reader: jspb.BinaryReader
  ): UpdateAPIKeyRequest;
}

export namespace UpdateAPIKeyRequest {
  export type AsObject = {
    id: string;
    environmentId: string;
    name?: google_protobuf_wrappers_pb.StringValue.AsObject;
    description?: google_protobuf_wrappers_pb.StringValue.AsObject;
    role: proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap];
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject;
    maintainer?: google_protobuf_wrappers_pb.StringValue.AsObject;
  };
}

export class UpdateAPIKeyResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAPIKeyResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateAPIKeyResponse
  ): UpdateAPIKeyResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateAPIKeyResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAPIKeyResponse;
  static deserializeBinaryFromReader(
    message: UpdateAPIKeyResponse,
    reader: jspb.BinaryReader
  ): UpdateAPIKeyResponse;
}

export namespace UpdateAPIKeyResponse {
  export type AsObject = {};
}
