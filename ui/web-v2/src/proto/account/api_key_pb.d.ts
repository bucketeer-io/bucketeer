// package: bucketeer.account
// file: proto/account/api_key.proto

import * as jspb from 'google-protobuf';
import * as proto_environment_environment_pb from '../../proto/environment/environment_pb';

export class APIKey extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getRole(): APIKey.RoleMap[keyof APIKey.RoleMap];
  setRole(value: APIKey.RoleMap[keyof APIKey.RoleMap]): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getMaintainer(): string;
  setMaintainer(value: string): void;

  getApiKey(): string;
  setApiKey(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getEnvironmentName(): string;
  setEnvironmentName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): APIKey.AsObject;
  static toObject(includeInstance: boolean, msg: APIKey): APIKey.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: APIKey,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): APIKey;
  static deserializeBinaryFromReader(
    message: APIKey,
    reader: jspb.BinaryReader
  ): APIKey;
}

export namespace APIKey {
  export type AsObject = {
    id: string;
    name: string;
    role: APIKey.RoleMap[keyof APIKey.RoleMap];
    disabled: boolean;
    createdAt: number;
    updatedAt: number;
    maintainer: string;
    apiKey: string;
    description: string;
    environmentName: string;
  };

  export interface RoleMap {
    UNKNOWN: 0;
    SDK_CLIENT: 1;
    SDK_SERVER: 2;
    PUBLIC_API_READ_ONLY: 3;
    PUBLIC_API_WRITE: 4;
    PUBLIC_API_ADMIN: 5;
  }

  export const Role: RoleMap;
}

export class EnvironmentAPIKey extends jspb.Message {
  hasApiKey(): boolean;
  clearApiKey(): void;
  getApiKey(): APIKey | undefined;
  setApiKey(value?: APIKey): void;

  getEnvironmentDisabled(): boolean;
  setEnvironmentDisabled(value: boolean): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  hasEnvironment(): boolean;
  clearEnvironment(): void;
  getEnvironment(): proto_environment_environment_pb.EnvironmentV2 | undefined;
  setEnvironment(value?: proto_environment_environment_pb.EnvironmentV2): void;

  getProjectUrlCode(): string;
  setProjectUrlCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentAPIKey.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentAPIKey
  ): EnvironmentAPIKey.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentAPIKey,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentAPIKey;
  static deserializeBinaryFromReader(
    message: EnvironmentAPIKey,
    reader: jspb.BinaryReader
  ): EnvironmentAPIKey;
}

export namespace EnvironmentAPIKey {
  export type AsObject = {
    apiKey?: APIKey.AsObject;
    environmentDisabled: boolean;
    projectId: string;
    environment?: proto_environment_environment_pb.EnvironmentV2.AsObject;
    projectUrlCode: string;
  };
}
