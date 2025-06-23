// package: bucketeer.team
// file: proto/team/team.proto

import * as jspb from 'google-protobuf';

export class Team extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  getOrganizationName(): string;
  setOrganizationName(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Team.AsObject;
  static toObject(includeInstance: boolean, msg: Team): Team.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Team,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Team;
  static deserializeBinaryFromReader(
    message: Team,
    reader: jspb.BinaryReader
  ): Team;
}

export namespace Team {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
    organizationId: string;
    organizationName: string;
    createdAt: number;
    updatedAt: number;
  };
}
