// package: bucketeer.team
// file: proto/team/service.proto

import * as jspb from 'google-protobuf';
import * as google_api_annotations_pb from '../../google/api/annotations_pb';
import * as google_api_field_behavior_pb from '../../google/api/field_behavior_pb';
import * as protoc_gen_openapiv2_options_annotations_pb from '../../protoc-gen-openapiv2/options/annotations_pb';
import * as proto_team_team_pb from '../../proto/team/team_pb';

export class CreateTeamRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTeamRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateTeamRequest
  ): CreateTeamRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateTeamRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateTeamRequest;
  static deserializeBinaryFromReader(
    message: CreateTeamRequest,
    reader: jspb.BinaryReader
  ): CreateTeamRequest;
}

export namespace CreateTeamRequest {
  export type AsObject = {
    name: string;
    description: string;
    organizationId: string;
  };
}

export class CreateTeamResponse extends jspb.Message {
  hasTeam(): boolean;
  clearTeam(): void;
  getTeam(): proto_team_team_pb.Team | undefined;
  setTeam(value?: proto_team_team_pb.Team): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTeamResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateTeamResponse
  ): CreateTeamResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateTeamResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateTeamResponse;
  static deserializeBinaryFromReader(
    message: CreateTeamResponse,
    reader: jspb.BinaryReader
  ): CreateTeamResponse;
}

export namespace CreateTeamResponse {
  export type AsObject = {
    team?: proto_team_team_pb.Team.AsObject;
  };
}

export class DeleteTeamRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteTeamRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteTeamRequest
  ): DeleteTeamRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteTeamRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteTeamRequest;
  static deserializeBinaryFromReader(
    message: DeleteTeamRequest,
    reader: jspb.BinaryReader
  ): DeleteTeamRequest;
}

export namespace DeleteTeamRequest {
  export type AsObject = {
    id: string;
    organizationId: string;
  };
}

export class DeleteTeamResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteTeamResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteTeamResponse
  ): DeleteTeamResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteTeamResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteTeamResponse;
  static deserializeBinaryFromReader(
    message: DeleteTeamResponse,
    reader: jspb.BinaryReader
  ): DeleteTeamResponse;
}

export namespace DeleteTeamResponse {
  export type AsObject = {};
}

export class ListTeamsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListTeamsRequest.OrderByMap[keyof ListTeamsRequest.OrderByMap];
  setOrderBy(
    value: ListTeamsRequest.OrderByMap[keyof ListTeamsRequest.OrderByMap]
  ): void;

  getOrderDirection(): ListTeamsRequest.OrderDirectionMap[keyof ListTeamsRequest.OrderDirectionMap];
  setOrderDirection(
    value: ListTeamsRequest.OrderDirectionMap[keyof ListTeamsRequest.OrderDirectionMap]
  ): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTeamsRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListTeamsRequest
  ): ListTeamsRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListTeamsRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListTeamsRequest;
  static deserializeBinaryFromReader(
    message: ListTeamsRequest,
    reader: jspb.BinaryReader
  ): ListTeamsRequest;
}

export namespace ListTeamsRequest {
  export type AsObject = {
    pageSize: number;
    cursor: string;
    orderBy: ListTeamsRequest.OrderByMap[keyof ListTeamsRequest.OrderByMap];
    orderDirection: ListTeamsRequest.OrderDirectionMap[keyof ListTeamsRequest.OrderDirectionMap];
    organizationId: string;
    searchKeyword: string;
  };

  export interface OrderByMap {
    DEFAULT: 0;
    NAME: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
    ORGANIZATION: 4;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListTeamsResponse extends jspb.Message {
  clearTeamsList(): void;
  getTeamsList(): Array<proto_team_team_pb.Team>;
  setTeamsList(value: Array<proto_team_team_pb.Team>): void;
  addTeams(
    value?: proto_team_team_pb.Team,
    index?: number
  ): proto_team_team_pb.Team;

  getNextCursor(): string;
  setNextCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTeamsResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListTeamsResponse
  ): ListTeamsResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListTeamsResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListTeamsResponse;
  static deserializeBinaryFromReader(
    message: ListTeamsResponse,
    reader: jspb.BinaryReader
  ): ListTeamsResponse;
}

export namespace ListTeamsResponse {
  export type AsObject = {
    teamsList: Array<proto_team_team_pb.Team.AsObject>;
    nextCursor: string;
    totalCount: number;
  };
}
