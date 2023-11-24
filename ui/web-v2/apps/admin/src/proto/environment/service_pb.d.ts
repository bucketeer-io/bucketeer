// package: bucketeer.environment
// file: proto/environment/service.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as proto_environment_environment_pb from "../../proto/environment/environment_pb";
import * as proto_environment_project_pb from "../../proto/environment/project_pb";
import * as proto_environment_organization_pb from "../../proto/environment/organization_pb";
import * as proto_environment_command_pb from "../../proto/environment/command_pb";

export class GetEnvironmentV2Request extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEnvironmentV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: GetEnvironmentV2Request): GetEnvironmentV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetEnvironmentV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEnvironmentV2Request;
  static deserializeBinaryFromReader(message: GetEnvironmentV2Request, reader: jspb.BinaryReader): GetEnvironmentV2Request;
}

export namespace GetEnvironmentV2Request {
  export type AsObject = {
    id: string,
  }
}

export class GetEnvironmentV2Response extends jspb.Message {
  hasEnvironment(): boolean;
  clearEnvironment(): void;
  getEnvironment(): proto_environment_environment_pb.EnvironmentV2 | undefined;
  setEnvironment(value?: proto_environment_environment_pb.EnvironmentV2): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEnvironmentV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: GetEnvironmentV2Response): GetEnvironmentV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetEnvironmentV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEnvironmentV2Response;
  static deserializeBinaryFromReader(message: GetEnvironmentV2Response, reader: jspb.BinaryReader): GetEnvironmentV2Response;
}

export namespace GetEnvironmentV2Response {
  export type AsObject = {
    environment?: proto_environment_environment_pb.EnvironmentV2.AsObject,
  }
}

export class ListEnvironmentsV2Request extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListEnvironmentsV2Request.OrderByMap[keyof ListEnvironmentsV2Request.OrderByMap];
  setOrderBy(value: ListEnvironmentsV2Request.OrderByMap[keyof ListEnvironmentsV2Request.OrderByMap]): void;

  getOrderDirection(): ListEnvironmentsV2Request.OrderDirectionMap[keyof ListEnvironmentsV2Request.OrderDirectionMap];
  setOrderDirection(value: ListEnvironmentsV2Request.OrderDirectionMap[keyof ListEnvironmentsV2Request.OrderDirectionMap]): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  hasArchived(): boolean;
  clearArchived(): void;
  getArchived(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setArchived(value?: google_protobuf_wrappers_pb.BoolValue): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnvironmentsV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnvironmentsV2Request): ListEnvironmentsV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnvironmentsV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnvironmentsV2Request;
  static deserializeBinaryFromReader(message: ListEnvironmentsV2Request, reader: jspb.BinaryReader): ListEnvironmentsV2Request;
}

export namespace ListEnvironmentsV2Request {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    orderBy: ListEnvironmentsV2Request.OrderByMap[keyof ListEnvironmentsV2Request.OrderByMap],
    orderDirection: ListEnvironmentsV2Request.OrderDirectionMap[keyof ListEnvironmentsV2Request.OrderDirectionMap],
    projectId: string,
    archived?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    searchKeyword: string,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    ID: 1;
    NAME: 2;
    URL_CODE: 3;
    CREATED_AT: 4;
    UPDATED_AT: 5;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListEnvironmentsV2Response extends jspb.Message {
  clearEnvironmentsList(): void;
  getEnvironmentsList(): Array<proto_environment_environment_pb.EnvironmentV2>;
  setEnvironmentsList(value: Array<proto_environment_environment_pb.EnvironmentV2>): void;
  addEnvironments(value?: proto_environment_environment_pb.EnvironmentV2, index?: number): proto_environment_environment_pb.EnvironmentV2;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnvironmentsV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnvironmentsV2Response): ListEnvironmentsV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnvironmentsV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnvironmentsV2Response;
  static deserializeBinaryFromReader(message: ListEnvironmentsV2Response, reader: jspb.BinaryReader): ListEnvironmentsV2Response;
}

export namespace ListEnvironmentsV2Response {
  export type AsObject = {
    environmentsList: Array<proto_environment_environment_pb.EnvironmentV2.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class CreateEnvironmentV2Request extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.CreateEnvironmentV2Command | undefined;
  setCommand(value?: proto_environment_command_pb.CreateEnvironmentV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateEnvironmentV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: CreateEnvironmentV2Request): CreateEnvironmentV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateEnvironmentV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateEnvironmentV2Request;
  static deserializeBinaryFromReader(message: CreateEnvironmentV2Request, reader: jspb.BinaryReader): CreateEnvironmentV2Request;
}

export namespace CreateEnvironmentV2Request {
  export type AsObject = {
    command?: proto_environment_command_pb.CreateEnvironmentV2Command.AsObject,
  }
}

export class CreateEnvironmentV2Response extends jspb.Message {
  hasEnvironment(): boolean;
  clearEnvironment(): void;
  getEnvironment(): proto_environment_environment_pb.EnvironmentV2 | undefined;
  setEnvironment(value?: proto_environment_environment_pb.EnvironmentV2): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateEnvironmentV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: CreateEnvironmentV2Response): CreateEnvironmentV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateEnvironmentV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateEnvironmentV2Response;
  static deserializeBinaryFromReader(message: CreateEnvironmentV2Response, reader: jspb.BinaryReader): CreateEnvironmentV2Response;
}

export namespace CreateEnvironmentV2Response {
  export type AsObject = {
    environment?: proto_environment_environment_pb.EnvironmentV2.AsObject,
  }
}

export class UpdateEnvironmentV2Request extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasRenameCommand(): boolean;
  clearRenameCommand(): void;
  getRenameCommand(): proto_environment_command_pb.RenameEnvironmentV2Command | undefined;
  setRenameCommand(value?: proto_environment_command_pb.RenameEnvironmentV2Command): void;

  hasChangeDescriptionCommand(): boolean;
  clearChangeDescriptionCommand(): void;
  getChangeDescriptionCommand(): proto_environment_command_pb.ChangeDescriptionEnvironmentV2Command | undefined;
  setChangeDescriptionCommand(value?: proto_environment_command_pb.ChangeDescriptionEnvironmentV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateEnvironmentV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateEnvironmentV2Request): UpdateEnvironmentV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateEnvironmentV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateEnvironmentV2Request;
  static deserializeBinaryFromReader(message: UpdateEnvironmentV2Request, reader: jspb.BinaryReader): UpdateEnvironmentV2Request;
}

export namespace UpdateEnvironmentV2Request {
  export type AsObject = {
    id: string,
    renameCommand?: proto_environment_command_pb.RenameEnvironmentV2Command.AsObject,
    changeDescriptionCommand?: proto_environment_command_pb.ChangeDescriptionEnvironmentV2Command.AsObject,
  }
}

export class UpdateEnvironmentV2Response extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateEnvironmentV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateEnvironmentV2Response): UpdateEnvironmentV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateEnvironmentV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateEnvironmentV2Response;
  static deserializeBinaryFromReader(message: UpdateEnvironmentV2Response, reader: jspb.BinaryReader): UpdateEnvironmentV2Response;
}

export namespace UpdateEnvironmentV2Response {
  export type AsObject = {
  }
}

export class ArchiveEnvironmentV2Request extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.ArchiveEnvironmentV2Command | undefined;
  setCommand(value?: proto_environment_command_pb.ArchiveEnvironmentV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveEnvironmentV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveEnvironmentV2Request): ArchiveEnvironmentV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveEnvironmentV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveEnvironmentV2Request;
  static deserializeBinaryFromReader(message: ArchiveEnvironmentV2Request, reader: jspb.BinaryReader): ArchiveEnvironmentV2Request;
}

export namespace ArchiveEnvironmentV2Request {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.ArchiveEnvironmentV2Command.AsObject,
  }
}

export class ArchiveEnvironmentV2Response extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveEnvironmentV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveEnvironmentV2Response): ArchiveEnvironmentV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveEnvironmentV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveEnvironmentV2Response;
  static deserializeBinaryFromReader(message: ArchiveEnvironmentV2Response, reader: jspb.BinaryReader): ArchiveEnvironmentV2Response;
}

export namespace ArchiveEnvironmentV2Response {
  export type AsObject = {
  }
}

export class UnarchiveEnvironmentV2Request extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.UnarchiveEnvironmentV2Command | undefined;
  setCommand(value?: proto_environment_command_pb.UnarchiveEnvironmentV2Command): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnarchiveEnvironmentV2Request.AsObject;
  static toObject(includeInstance: boolean, msg: UnarchiveEnvironmentV2Request): UnarchiveEnvironmentV2Request.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UnarchiveEnvironmentV2Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnarchiveEnvironmentV2Request;
  static deserializeBinaryFromReader(message: UnarchiveEnvironmentV2Request, reader: jspb.BinaryReader): UnarchiveEnvironmentV2Request;
}

export namespace UnarchiveEnvironmentV2Request {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.UnarchiveEnvironmentV2Command.AsObject,
  }
}

export class UnarchiveEnvironmentV2Response extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnarchiveEnvironmentV2Response.AsObject;
  static toObject(includeInstance: boolean, msg: UnarchiveEnvironmentV2Response): UnarchiveEnvironmentV2Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UnarchiveEnvironmentV2Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnarchiveEnvironmentV2Response;
  static deserializeBinaryFromReader(message: UnarchiveEnvironmentV2Response, reader: jspb.BinaryReader): UnarchiveEnvironmentV2Response;
}

export namespace UnarchiveEnvironmentV2Response {
  export type AsObject = {
  }
}

export class GetProjectRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProjectRequest): GetProjectRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProjectRequest;
  static deserializeBinaryFromReader(message: GetProjectRequest, reader: jspb.BinaryReader): GetProjectRequest;
}

export namespace GetProjectRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetProjectResponse extends jspb.Message {
  hasProject(): boolean;
  clearProject(): void;
  getProject(): proto_environment_project_pb.Project | undefined;
  setProject(value?: proto_environment_project_pb.Project): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProjectResponse): GetProjectResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProjectResponse;
  static deserializeBinaryFromReader(message: GetProjectResponse, reader: jspb.BinaryReader): GetProjectResponse;
}

export namespace GetProjectResponse {
  export type AsObject = {
    project?: proto_environment_project_pb.Project.AsObject,
  }
}

export class ListProjectsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListProjectsRequest.OrderByMap[keyof ListProjectsRequest.OrderByMap];
  setOrderBy(value: ListProjectsRequest.OrderByMap[keyof ListProjectsRequest.OrderByMap]): void;

  getOrderDirection(): ListProjectsRequest.OrderDirectionMap[keyof ListProjectsRequest.OrderDirectionMap];
  setOrderDirection(value: ListProjectsRequest.OrderDirectionMap[keyof ListProjectsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectsRequest): ListProjectsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListProjectsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectsRequest;
  static deserializeBinaryFromReader(message: ListProjectsRequest, reader: jspb.BinaryReader): ListProjectsRequest;
}

export namespace ListProjectsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    orderBy: ListProjectsRequest.OrderByMap[keyof ListProjectsRequest.OrderByMap],
    orderDirection: ListProjectsRequest.OrderDirectionMap[keyof ListProjectsRequest.OrderDirectionMap],
    searchKeyword: string,
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    ID: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
    NAME: 4;
    URL_CODE: 5;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListProjectsResponse extends jspb.Message {
  clearProjectsList(): void;
  getProjectsList(): Array<proto_environment_project_pb.Project>;
  setProjectsList(value: Array<proto_environment_project_pb.Project>): void;
  addProjects(value?: proto_environment_project_pb.Project, index?: number): proto_environment_project_pb.Project;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectsResponse): ListProjectsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListProjectsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectsResponse;
  static deserializeBinaryFromReader(message: ListProjectsResponse, reader: jspb.BinaryReader): ListProjectsResponse;
}

export namespace ListProjectsResponse {
  export type AsObject = {
    projectsList: Array<proto_environment_project_pb.Project.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class CreateProjectRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.CreateProjectCommand | undefined;
  setCommand(value?: proto_environment_command_pb.CreateProjectCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateProjectRequest): CreateProjectRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateProjectRequest;
  static deserializeBinaryFromReader(message: CreateProjectRequest, reader: jspb.BinaryReader): CreateProjectRequest;
}

export namespace CreateProjectRequest {
  export type AsObject = {
    command?: proto_environment_command_pb.CreateProjectCommand.AsObject,
  }
}

export class CreateProjectResponse extends jspb.Message {
  hasProject(): boolean;
  clearProject(): void;
  getProject(): proto_environment_project_pb.Project | undefined;
  setProject(value?: proto_environment_project_pb.Project): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateProjectResponse): CreateProjectResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateProjectResponse;
  static deserializeBinaryFromReader(message: CreateProjectResponse, reader: jspb.BinaryReader): CreateProjectResponse;
}

export namespace CreateProjectResponse {
  export type AsObject = {
    project?: proto_environment_project_pb.Project.AsObject,
  }
}

export class CreateTrialProjectRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.CreateTrialProjectCommand | undefined;
  setCommand(value?: proto_environment_command_pb.CreateTrialProjectCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTrialProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateTrialProjectRequest): CreateTrialProjectRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateTrialProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateTrialProjectRequest;
  static deserializeBinaryFromReader(message: CreateTrialProjectRequest, reader: jspb.BinaryReader): CreateTrialProjectRequest;
}

export namespace CreateTrialProjectRequest {
  export type AsObject = {
    command?: proto_environment_command_pb.CreateTrialProjectCommand.AsObject,
  }
}

export class CreateTrialProjectResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTrialProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateTrialProjectResponse): CreateTrialProjectResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateTrialProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateTrialProjectResponse;
  static deserializeBinaryFromReader(message: CreateTrialProjectResponse, reader: jspb.BinaryReader): CreateTrialProjectResponse;
}

export namespace CreateTrialProjectResponse {
  export type AsObject = {
  }
}

export class UpdateProjectRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasChangeDescriptionCommand(): boolean;
  clearChangeDescriptionCommand(): void;
  getChangeDescriptionCommand(): proto_environment_command_pb.ChangeDescriptionProjectCommand | undefined;
  setChangeDescriptionCommand(value?: proto_environment_command_pb.ChangeDescriptionProjectCommand): void;

  hasRenameCommand(): boolean;
  clearRenameCommand(): void;
  getRenameCommand(): proto_environment_command_pb.RenameProjectCommand | undefined;
  setRenameCommand(value?: proto_environment_command_pb.RenameProjectCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectRequest): UpdateProjectRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectRequest;
  static deserializeBinaryFromReader(message: UpdateProjectRequest, reader: jspb.BinaryReader): UpdateProjectRequest;
}

export namespace UpdateProjectRequest {
  export type AsObject = {
    id: string,
    changeDescriptionCommand?: proto_environment_command_pb.ChangeDescriptionProjectCommand.AsObject,
    renameCommand?: proto_environment_command_pb.RenameProjectCommand.AsObject,
  }
}

export class UpdateProjectResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectResponse): UpdateProjectResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectResponse;
  static deserializeBinaryFromReader(message: UpdateProjectResponse, reader: jspb.BinaryReader): UpdateProjectResponse;
}

export namespace UpdateProjectResponse {
  export type AsObject = {
  }
}

export class EnableProjectRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.EnableProjectCommand | undefined;
  setCommand(value?: proto_environment_command_pb.EnableProjectCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableProjectRequest): EnableProjectRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableProjectRequest;
  static deserializeBinaryFromReader(message: EnableProjectRequest, reader: jspb.BinaryReader): EnableProjectRequest;
}

export namespace EnableProjectRequest {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.EnableProjectCommand.AsObject,
  }
}

export class EnableProjectResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableProjectResponse): EnableProjectResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableProjectResponse;
  static deserializeBinaryFromReader(message: EnableProjectResponse, reader: jspb.BinaryReader): EnableProjectResponse;
}

export namespace EnableProjectResponse {
  export type AsObject = {
  }
}

export class DisableProjectRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.DisableProjectCommand | undefined;
  setCommand(value?: proto_environment_command_pb.DisableProjectCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableProjectRequest): DisableProjectRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableProjectRequest;
  static deserializeBinaryFromReader(message: DisableProjectRequest, reader: jspb.BinaryReader): DisableProjectRequest;
}

export namespace DisableProjectRequest {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.DisableProjectCommand.AsObject,
  }
}

export class DisableProjectResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableProjectResponse): DisableProjectResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableProjectResponse;
  static deserializeBinaryFromReader(message: DisableProjectResponse, reader: jspb.BinaryReader): DisableProjectResponse;
}

export namespace DisableProjectResponse {
  export type AsObject = {
  }
}

export class ConvertTrialProjectRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.ConvertTrialProjectCommand | undefined;
  setCommand(value?: proto_environment_command_pb.ConvertTrialProjectCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConvertTrialProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConvertTrialProjectRequest): ConvertTrialProjectRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConvertTrialProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConvertTrialProjectRequest;
  static deserializeBinaryFromReader(message: ConvertTrialProjectRequest, reader: jspb.BinaryReader): ConvertTrialProjectRequest;
}

export namespace ConvertTrialProjectRequest {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.ConvertTrialProjectCommand.AsObject,
  }
}

export class ConvertTrialProjectResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConvertTrialProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ConvertTrialProjectResponse): ConvertTrialProjectResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConvertTrialProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConvertTrialProjectResponse;
  static deserializeBinaryFromReader(message: ConvertTrialProjectResponse, reader: jspb.BinaryReader): ConvertTrialProjectResponse;
}

export namespace ConvertTrialProjectResponse {
  export type AsObject = {
  }
}

export class GetOrganizationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationRequest): GetOrganizationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationRequest;
  static deserializeBinaryFromReader(message: GetOrganizationRequest, reader: jspb.BinaryReader): GetOrganizationRequest;
}

export namespace GetOrganizationRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetOrganizationResponse extends jspb.Message {
  hasOrganization(): boolean;
  clearOrganization(): void;
  getOrganization(): proto_environment_organization_pb.Organization | undefined;
  setOrganization(value?: proto_environment_organization_pb.Organization): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationResponse): GetOrganizationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationResponse;
  static deserializeBinaryFromReader(message: GetOrganizationResponse, reader: jspb.BinaryReader): GetOrganizationResponse;
}

export namespace GetOrganizationResponse {
  export type AsObject = {
    organization?: proto_environment_organization_pb.Organization.AsObject,
  }
}

export class ListOrganizationsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListOrganizationsRequest.OrderByMap[keyof ListOrganizationsRequest.OrderByMap];
  setOrderBy(value: ListOrganizationsRequest.OrderByMap[keyof ListOrganizationsRequest.OrderByMap]): void;

  getOrderDirection(): ListOrganizationsRequest.OrderDirectionMap[keyof ListOrganizationsRequest.OrderDirectionMap];
  setOrderDirection(value: ListOrganizationsRequest.OrderDirectionMap[keyof ListOrganizationsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasArchived(): boolean;
  clearArchived(): void;
  getArchived(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setArchived(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrganizationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrganizationsRequest): ListOrganizationsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListOrganizationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrganizationsRequest;
  static deserializeBinaryFromReader(message: ListOrganizationsRequest, reader: jspb.BinaryReader): ListOrganizationsRequest;
}

export namespace ListOrganizationsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    orderBy: ListOrganizationsRequest.OrderByMap[keyof ListOrganizationsRequest.OrderByMap],
    orderDirection: ListOrganizationsRequest.OrderDirectionMap[keyof ListOrganizationsRequest.OrderDirectionMap],
    searchKeyword: string,
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    archived?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    ID: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
    NAME: 4;
    URL_CODE: 5;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListOrganizationsResponse extends jspb.Message {
  clearOrganizationsList(): void;
  getOrganizationsList(): Array<proto_environment_organization_pb.Organization>;
  setOrganizationsList(value: Array<proto_environment_organization_pb.Organization>): void;
  addOrganizations(value?: proto_environment_organization_pb.Organization, index?: number): proto_environment_organization_pb.Organization;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrganizationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrganizationsResponse): ListOrganizationsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListOrganizationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrganizationsResponse;
  static deserializeBinaryFromReader(message: ListOrganizationsResponse, reader: jspb.BinaryReader): ListOrganizationsResponse;
}

export namespace ListOrganizationsResponse {
  export type AsObject = {
    organizationsList: Array<proto_environment_organization_pb.Organization.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class CreateOrganizationRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.CreateOrganizationCommand | undefined;
  setCommand(value?: proto_environment_command_pb.CreateOrganizationCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateOrganizationRequest): CreateOrganizationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateOrganizationRequest;
  static deserializeBinaryFromReader(message: CreateOrganizationRequest, reader: jspb.BinaryReader): CreateOrganizationRequest;
}

export namespace CreateOrganizationRequest {
  export type AsObject = {
    command?: proto_environment_command_pb.CreateOrganizationCommand.AsObject,
  }
}

export class CreateOrganizationResponse extends jspb.Message {
  hasOrganization(): boolean;
  clearOrganization(): void;
  getOrganization(): proto_environment_organization_pb.Organization | undefined;
  setOrganization(value?: proto_environment_organization_pb.Organization): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateOrganizationResponse): CreateOrganizationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateOrganizationResponse;
  static deserializeBinaryFromReader(message: CreateOrganizationResponse, reader: jspb.BinaryReader): CreateOrganizationResponse;
}

export namespace CreateOrganizationResponse {
  export type AsObject = {
    organization?: proto_environment_organization_pb.Organization.AsObject,
  }
}

export class UpdateOrganizationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasRenameCommand(): boolean;
  clearRenameCommand(): void;
  getRenameCommand(): proto_environment_command_pb.ChangeNameOrganizationCommand | undefined;
  setRenameCommand(value?: proto_environment_command_pb.ChangeNameOrganizationCommand): void;

  hasChangeDescriptionCommand(): boolean;
  clearChangeDescriptionCommand(): void;
  getChangeDescriptionCommand(): proto_environment_command_pb.ChangeDescriptionOrganizationCommand | undefined;
  setChangeDescriptionCommand(value?: proto_environment_command_pb.ChangeDescriptionOrganizationCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrganizationRequest): UpdateOrganizationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrganizationRequest;
  static deserializeBinaryFromReader(message: UpdateOrganizationRequest, reader: jspb.BinaryReader): UpdateOrganizationRequest;
}

export namespace UpdateOrganizationRequest {
  export type AsObject = {
    id: string,
    renameCommand?: proto_environment_command_pb.ChangeNameOrganizationCommand.AsObject,
    changeDescriptionCommand?: proto_environment_command_pb.ChangeDescriptionOrganizationCommand.AsObject,
  }
}

export class UpdateOrganizationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrganizationResponse): UpdateOrganizationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrganizationResponse;
  static deserializeBinaryFromReader(message: UpdateOrganizationResponse, reader: jspb.BinaryReader): UpdateOrganizationResponse;
}

export namespace UpdateOrganizationResponse {
  export type AsObject = {
  }
}

export class EnableOrganizationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.EnableOrganizationCommand | undefined;
  setCommand(value?: proto_environment_command_pb.EnableOrganizationCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableOrganizationRequest): EnableOrganizationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableOrganizationRequest;
  static deserializeBinaryFromReader(message: EnableOrganizationRequest, reader: jspb.BinaryReader): EnableOrganizationRequest;
}

export namespace EnableOrganizationRequest {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.EnableOrganizationCommand.AsObject,
  }
}

export class EnableOrganizationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableOrganizationResponse): EnableOrganizationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableOrganizationResponse;
  static deserializeBinaryFromReader(message: EnableOrganizationResponse, reader: jspb.BinaryReader): EnableOrganizationResponse;
}

export namespace EnableOrganizationResponse {
  export type AsObject = {
  }
}

export class DisableOrganizationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.DisableOrganizationCommand | undefined;
  setCommand(value?: proto_environment_command_pb.DisableOrganizationCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableOrganizationRequest): DisableOrganizationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableOrganizationRequest;
  static deserializeBinaryFromReader(message: DisableOrganizationRequest, reader: jspb.BinaryReader): DisableOrganizationRequest;
}

export namespace DisableOrganizationRequest {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.DisableOrganizationCommand.AsObject,
  }
}

export class DisableOrganizationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableOrganizationResponse): DisableOrganizationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableOrganizationResponse;
  static deserializeBinaryFromReader(message: DisableOrganizationResponse, reader: jspb.BinaryReader): DisableOrganizationResponse;
}

export namespace DisableOrganizationResponse {
  export type AsObject = {
  }
}

export class ArchiveOrganizationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.ArchiveOrganizationCommand | undefined;
  setCommand(value?: proto_environment_command_pb.ArchiveOrganizationCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveOrganizationRequest): ArchiveOrganizationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveOrganizationRequest;
  static deserializeBinaryFromReader(message: ArchiveOrganizationRequest, reader: jspb.BinaryReader): ArchiveOrganizationRequest;
}

export namespace ArchiveOrganizationRequest {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.ArchiveOrganizationCommand.AsObject,
  }
}

export class ArchiveOrganizationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveOrganizationResponse): ArchiveOrganizationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveOrganizationResponse;
  static deserializeBinaryFromReader(message: ArchiveOrganizationResponse, reader: jspb.BinaryReader): ArchiveOrganizationResponse;
}

export namespace ArchiveOrganizationResponse {
  export type AsObject = {
  }
}

export class UnarchiveOrganizationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.UnarchiveOrganizationCommand | undefined;
  setCommand(value?: proto_environment_command_pb.UnarchiveOrganizationCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnarchiveOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UnarchiveOrganizationRequest): UnarchiveOrganizationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UnarchiveOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnarchiveOrganizationRequest;
  static deserializeBinaryFromReader(message: UnarchiveOrganizationRequest, reader: jspb.BinaryReader): UnarchiveOrganizationRequest;
}

export namespace UnarchiveOrganizationRequest {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.UnarchiveOrganizationCommand.AsObject,
  }
}

export class UnarchiveOrganizationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnarchiveOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UnarchiveOrganizationResponse): UnarchiveOrganizationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UnarchiveOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnarchiveOrganizationResponse;
  static deserializeBinaryFromReader(message: UnarchiveOrganizationResponse, reader: jspb.BinaryReader): UnarchiveOrganizationResponse;
}

export namespace UnarchiveOrganizationResponse {
  export type AsObject = {
  }
}

export class ConvertTrialOrganizationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.ConvertTrialOrganizationCommand | undefined;
  setCommand(value?: proto_environment_command_pb.ConvertTrialOrganizationCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConvertTrialOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConvertTrialOrganizationRequest): ConvertTrialOrganizationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConvertTrialOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConvertTrialOrganizationRequest;
  static deserializeBinaryFromReader(message: ConvertTrialOrganizationRequest, reader: jspb.BinaryReader): ConvertTrialOrganizationRequest;
}

export namespace ConvertTrialOrganizationRequest {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.ConvertTrialOrganizationCommand.AsObject,
  }
}

export class ConvertTrialOrganizationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConvertTrialOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ConvertTrialOrganizationResponse): ConvertTrialOrganizationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConvertTrialOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConvertTrialOrganizationResponse;
  static deserializeBinaryFromReader(message: ConvertTrialOrganizationResponse, reader: jspb.BinaryReader): ConvertTrialOrganizationResponse;
}

export namespace ConvertTrialOrganizationResponse {
  export type AsObject = {
  }
}

