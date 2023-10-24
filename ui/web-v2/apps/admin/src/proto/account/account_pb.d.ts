// package: bucketeer.account
// file: proto/account/account.proto

import * as jspb from "google-protobuf";
import * as proto_environment_environment_pb from "../../proto/environment/environment_pb";

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

