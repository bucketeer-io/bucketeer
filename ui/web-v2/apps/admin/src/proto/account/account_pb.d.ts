// package: bucketeer.account
// file: proto/account/account.proto

import * as jspb from "google-protobuf";
import * as proto_environment_environment_pb from "../../proto/environment/environment_pb";
import * as proto_environment_project_pb from "../../proto/environment/project_pb";
import * as proto_environment_organization_pb from "../../proto/environment/organization_pb";

export class Account extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getName(): string;
  setName(value: string): void;

  getRole(): Account.RoleMap[keyof Account.RoleMap];
  setRole(value: Account.RoleMap[keyof Account.RoleMap]): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Account.AsObject;
  static toObject(includeInstance: boolean, msg: Account): Account.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Account, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Account;
  static deserializeBinaryFromReader(message: Account, reader: jspb.BinaryReader): Account;
}

export namespace Account {
  export type AsObject = {
    id: string,
    email: string,
    name: string,
    role: Account.RoleMap[keyof Account.RoleMap],
    disabled: boolean,
    createdAt: number,
    updatedAt: number,
    deleted: boolean,
  }

  export interface RoleMap {
    VIEWER: 0;
    EDITOR: 1;
    OWNER: 2;
    UNASSIGNED: 99;
  }

  export const Role: RoleMap;
}

export class AccountV2 extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getName(): string;
  setName(value: string): void;

  getAvatarImageUrl(): string;
  setAvatarImageUrl(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  getOrganizationRole(): AccountV2.Role.OrganizationMap[keyof AccountV2.Role.OrganizationMap];
  setOrganizationRole(value: AccountV2.Role.OrganizationMap[keyof AccountV2.Role.OrganizationMap]): void;

  clearEnvironmentRolesList(): void;
  getEnvironmentRolesList(): Array<AccountV2.EnvironmentRole>;
  setEnvironmentRolesList(value: Array<AccountV2.EnvironmentRole>): void;
  addEnvironmentRoles(value?: AccountV2.EnvironmentRole, index?: number): AccountV2.EnvironmentRole;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2.AsObject;
  static toObject(includeInstance: boolean, msg: AccountV2): AccountV2.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccountV2, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2;
  static deserializeBinaryFromReader(message: AccountV2, reader: jspb.BinaryReader): AccountV2;
}

export namespace AccountV2 {
  export type AsObject = {
    email: string,
    name: string,
    avatarImageUrl: string,
    organizationId: string,
    organizationRole: AccountV2.Role.OrganizationMap[keyof AccountV2.Role.OrganizationMap],
    environmentRolesList: Array<AccountV2.EnvironmentRole.AsObject>,
    disabled: boolean,
    createdAt: number,
    updatedAt: number,
  }

  export class Role extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Role.AsObject;
    static toObject(includeInstance: boolean, msg: Role): Role.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Role, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Role;
    static deserializeBinaryFromReader(message: Role, reader: jspb.BinaryReader): Role;
  }

  export namespace Role {
    export type AsObject = {
    }

    export interface EnvironmentMap {
      ENVIRONMENT_UNASSIGNED: 0;
      ENVIRONMENT_VIEWER: 1;
      ENVIRONMENT_EDITOR: 2;
    }

    export const Environment: EnvironmentMap;

    export interface OrganizationMap {
      ORGANIZATION_UNASSIGNED: 0;
      ORGANIZATION_MEMBER: 1;
      ORGANIZATION_ADMIN: 2;
      ORGANIZATION_OWNER: 3;
    }

    export const Organization: OrganizationMap;
  }

  export class EnvironmentRole extends jspb.Message {
    getEnvironmentId(): string;
    setEnvironmentId(value: string): void;

    getRole(): AccountV2.Role.EnvironmentMap[keyof AccountV2.Role.EnvironmentMap];
    setRole(value: AccountV2.Role.EnvironmentMap[keyof AccountV2.Role.EnvironmentMap]): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): EnvironmentRole.AsObject;
    static toObject(includeInstance: boolean, msg: EnvironmentRole): EnvironmentRole.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: EnvironmentRole, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): EnvironmentRole;
    static deserializeBinaryFromReader(message: EnvironmentRole, reader: jspb.BinaryReader): EnvironmentRole;
  }

  export namespace EnvironmentRole {
    export type AsObject = {
      environmentId: string,
      role: AccountV2.Role.EnvironmentMap[keyof AccountV2.Role.EnvironmentMap],
    }
  }
}

export class EnvironmentRoleV2 extends jspb.Message {
  hasEnvironment(): boolean;
  clearEnvironment(): void;
  getEnvironment(): proto_environment_environment_pb.EnvironmentV2 | undefined;
  setEnvironment(value?: proto_environment_environment_pb.EnvironmentV2): void;

  getRole(): Account.RoleMap[keyof Account.RoleMap];
  setRole(value: Account.RoleMap[keyof Account.RoleMap]): void;

  getTrialProject(): boolean;
  setTrialProject(value: boolean): void;

  getTrialStartedAt(): number;
  setTrialStartedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentRoleV2.AsObject;
  static toObject(includeInstance: boolean, msg: EnvironmentRoleV2): EnvironmentRoleV2.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnvironmentRoleV2, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentRoleV2;
  static deserializeBinaryFromReader(message: EnvironmentRoleV2, reader: jspb.BinaryReader): EnvironmentRoleV2;
}

export namespace EnvironmentRoleV2 {
  export type AsObject = {
    environment?: proto_environment_environment_pb.EnvironmentV2.AsObject,
    role: Account.RoleMap[keyof Account.RoleMap],
    trialProject: boolean,
    trialStartedAt: number,
  }
}

export class ConsoleAccount extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getName(): string;
  setName(value: string): void;

  getAvatarUrl(): string;
  setAvatarUrl(value: string): void;

  getIsSystemAdmin(): boolean;
  setIsSystemAdmin(value: boolean): void;

  hasOrganization(): boolean;
  clearOrganization(): void;
  getOrganization(): proto_environment_organization_pb.Organization | undefined;
  setOrganization(value?: proto_environment_organization_pb.Organization): void;

  getOrganizationRole(): AccountV2.Role.OrganizationMap[keyof AccountV2.Role.OrganizationMap];
  setOrganizationRole(value: AccountV2.Role.OrganizationMap[keyof AccountV2.Role.OrganizationMap]): void;

  clearEnvironmentRolesList(): void;
  getEnvironmentRolesList(): Array<ConsoleAccount.EnvironmentRole>;
  setEnvironmentRolesList(value: Array<ConsoleAccount.EnvironmentRole>): void;
  addEnvironmentRoles(value?: ConsoleAccount.EnvironmentRole, index?: number): ConsoleAccount.EnvironmentRole;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConsoleAccount.AsObject;
  static toObject(includeInstance: boolean, msg: ConsoleAccount): ConsoleAccount.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConsoleAccount, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConsoleAccount;
  static deserializeBinaryFromReader(message: ConsoleAccount, reader: jspb.BinaryReader): ConsoleAccount;
}

export namespace ConsoleAccount {
  export type AsObject = {
    email: string,
    name: string,
    avatarUrl: string,
    isSystemAdmin: boolean,
    organization?: proto_environment_organization_pb.Organization.AsObject,
    organizationRole: AccountV2.Role.OrganizationMap[keyof AccountV2.Role.OrganizationMap],
    environmentRolesList: Array<ConsoleAccount.EnvironmentRole.AsObject>,
  }

  export class EnvironmentRole extends jspb.Message {
    hasEnvironment(): boolean;
    clearEnvironment(): void;
    getEnvironment(): proto_environment_environment_pb.EnvironmentV2 | undefined;
    setEnvironment(value?: proto_environment_environment_pb.EnvironmentV2): void;

    hasProject(): boolean;
    clearProject(): void;
    getProject(): proto_environment_project_pb.Project | undefined;
    setProject(value?: proto_environment_project_pb.Project): void;

    getRole(): AccountV2.Role.EnvironmentMap[keyof AccountV2.Role.EnvironmentMap];
    setRole(value: AccountV2.Role.EnvironmentMap[keyof AccountV2.Role.EnvironmentMap]): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): EnvironmentRole.AsObject;
    static toObject(includeInstance: boolean, msg: EnvironmentRole): EnvironmentRole.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: EnvironmentRole, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): EnvironmentRole;
    static deserializeBinaryFromReader(message: EnvironmentRole, reader: jspb.BinaryReader): EnvironmentRole;
  }

  export namespace EnvironmentRole {
    export type AsObject = {
      environment?: proto_environment_environment_pb.EnvironmentV2.AsObject,
      project?: proto_environment_project_pb.Project.AsObject,
      role: AccountV2.Role.EnvironmentMap[keyof AccountV2.Role.EnvironmentMap],
    }
  }
}

