// package: bucketeer.environment
// file: proto/environment/service.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as proto_environment_environment_pb from "../../proto/environment/environment_pb";
import * as proto_environment_project_pb from "../../proto/environment/project_pb";
import * as proto_environment_command_pb from "../../proto/environment/command_pb";

export class GetEnvironmentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEnvironmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetEnvironmentRequest): GetEnvironmentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetEnvironmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEnvironmentRequest;
  static deserializeBinaryFromReader(message: GetEnvironmentRequest, reader: jspb.BinaryReader): GetEnvironmentRequest;
}

export namespace GetEnvironmentRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetEnvironmentResponse extends jspb.Message {
  hasEnvironment(): boolean;
  clearEnvironment(): void;
  getEnvironment(): proto_environment_environment_pb.Environment | undefined;
  setEnvironment(value?: proto_environment_environment_pb.Environment): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEnvironmentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetEnvironmentResponse): GetEnvironmentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetEnvironmentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEnvironmentResponse;
  static deserializeBinaryFromReader(message: GetEnvironmentResponse, reader: jspb.BinaryReader): GetEnvironmentResponse;
}

export namespace GetEnvironmentResponse {
  export type AsObject = {
    environment?: proto_environment_environment_pb.Environment.AsObject,
  }
}

export class GetEnvironmentByNamespaceRequest extends jspb.Message {
  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEnvironmentByNamespaceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetEnvironmentByNamespaceRequest): GetEnvironmentByNamespaceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetEnvironmentByNamespaceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEnvironmentByNamespaceRequest;
  static deserializeBinaryFromReader(message: GetEnvironmentByNamespaceRequest, reader: jspb.BinaryReader): GetEnvironmentByNamespaceRequest;
}

export namespace GetEnvironmentByNamespaceRequest {
  export type AsObject = {
    namespace: string,
  }
}

export class GetEnvironmentByNamespaceResponse extends jspb.Message {
  hasEnvironment(): boolean;
  clearEnvironment(): void;
  getEnvironment(): proto_environment_environment_pb.Environment | undefined;
  setEnvironment(value?: proto_environment_environment_pb.Environment): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEnvironmentByNamespaceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetEnvironmentByNamespaceResponse): GetEnvironmentByNamespaceResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetEnvironmentByNamespaceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEnvironmentByNamespaceResponse;
  static deserializeBinaryFromReader(message: GetEnvironmentByNamespaceResponse, reader: jspb.BinaryReader): GetEnvironmentByNamespaceResponse;
}

export namespace GetEnvironmentByNamespaceResponse {
  export type AsObject = {
    environment?: proto_environment_environment_pb.Environment.AsObject,
  }
}

export class ListEnvironmentsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getOrderBy(): ListEnvironmentsRequest.OrderByMap[keyof ListEnvironmentsRequest.OrderByMap];
  setOrderBy(value: ListEnvironmentsRequest.OrderByMap[keyof ListEnvironmentsRequest.OrderByMap]): void;

  getOrderDirection(): ListEnvironmentsRequest.OrderDirectionMap[keyof ListEnvironmentsRequest.OrderDirectionMap];
  setOrderDirection(value: ListEnvironmentsRequest.OrderDirectionMap[keyof ListEnvironmentsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnvironmentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnvironmentsRequest): ListEnvironmentsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnvironmentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnvironmentsRequest;
  static deserializeBinaryFromReader(message: ListEnvironmentsRequest, reader: jspb.BinaryReader): ListEnvironmentsRequest;
}

export namespace ListEnvironmentsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    projectId: string,
    orderBy: ListEnvironmentsRequest.OrderByMap[keyof ListEnvironmentsRequest.OrderByMap],
    orderDirection: ListEnvironmentsRequest.OrderDirectionMap[keyof ListEnvironmentsRequest.OrderDirectionMap],
    searchKeyword: string,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    ID: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListEnvironmentsResponse extends jspb.Message {
  clearEnvironmentsList(): void;
  getEnvironmentsList(): Array<proto_environment_environment_pb.Environment>;
  setEnvironmentsList(value: Array<proto_environment_environment_pb.Environment>): void;
  addEnvironments(value?: proto_environment_environment_pb.Environment, index?: number): proto_environment_environment_pb.Environment;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnvironmentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnvironmentsResponse): ListEnvironmentsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnvironmentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnvironmentsResponse;
  static deserializeBinaryFromReader(message: ListEnvironmentsResponse, reader: jspb.BinaryReader): ListEnvironmentsResponse;
}

export namespace ListEnvironmentsResponse {
  export type AsObject = {
    environmentsList: Array<proto_environment_environment_pb.Environment.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class CreateEnvironmentRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.CreateEnvironmentCommand | undefined;
  setCommand(value?: proto_environment_command_pb.CreateEnvironmentCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateEnvironmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateEnvironmentRequest): CreateEnvironmentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateEnvironmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateEnvironmentRequest;
  static deserializeBinaryFromReader(message: CreateEnvironmentRequest, reader: jspb.BinaryReader): CreateEnvironmentRequest;
}

export namespace CreateEnvironmentRequest {
  export type AsObject = {
    command?: proto_environment_command_pb.CreateEnvironmentCommand.AsObject,
  }
}

export class CreateEnvironmentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateEnvironmentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateEnvironmentResponse): CreateEnvironmentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateEnvironmentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateEnvironmentResponse;
  static deserializeBinaryFromReader(message: CreateEnvironmentResponse, reader: jspb.BinaryReader): CreateEnvironmentResponse;
}

export namespace CreateEnvironmentResponse {
  export type AsObject = {
  }
}

export class UpdateEnvironmentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasRenameCommand(): boolean;
  clearRenameCommand(): void;
  getRenameCommand(): proto_environment_command_pb.RenameEnvironmentCommand | undefined;
  setRenameCommand(value?: proto_environment_command_pb.RenameEnvironmentCommand): void;

  hasChangeDescriptionCommand(): boolean;
  clearChangeDescriptionCommand(): void;
  getChangeDescriptionCommand(): proto_environment_command_pb.ChangeDescriptionEnvironmentCommand | undefined;
  setChangeDescriptionCommand(value?: proto_environment_command_pb.ChangeDescriptionEnvironmentCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateEnvironmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateEnvironmentRequest): UpdateEnvironmentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateEnvironmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateEnvironmentRequest;
  static deserializeBinaryFromReader(message: UpdateEnvironmentRequest, reader: jspb.BinaryReader): UpdateEnvironmentRequest;
}

export namespace UpdateEnvironmentRequest {
  export type AsObject = {
    id: string,
    renameCommand?: proto_environment_command_pb.RenameEnvironmentCommand.AsObject,
    changeDescriptionCommand?: proto_environment_command_pb.ChangeDescriptionEnvironmentCommand.AsObject,
  }
}

export class UpdateEnvironmentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateEnvironmentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateEnvironmentResponse): UpdateEnvironmentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateEnvironmentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateEnvironmentResponse;
  static deserializeBinaryFromReader(message: UpdateEnvironmentResponse, reader: jspb.BinaryReader): UpdateEnvironmentResponse;
}

export namespace UpdateEnvironmentResponse {
  export type AsObject = {
  }
}

export class DeleteEnvironmentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_environment_command_pb.DeleteEnvironmentCommand | undefined;
  setCommand(value?: proto_environment_command_pb.DeleteEnvironmentCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteEnvironmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteEnvironmentRequest): DeleteEnvironmentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteEnvironmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteEnvironmentRequest;
  static deserializeBinaryFromReader(message: DeleteEnvironmentRequest, reader: jspb.BinaryReader): DeleteEnvironmentRequest;
}

export namespace DeleteEnvironmentRequest {
  export type AsObject = {
    id: string,
    command?: proto_environment_command_pb.DeleteEnvironmentCommand.AsObject,
  }
}

export class DeleteEnvironmentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteEnvironmentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteEnvironmentResponse): DeleteEnvironmentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteEnvironmentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteEnvironmentResponse;
  static deserializeBinaryFromReader(message: DeleteEnvironmentResponse, reader: jspb.BinaryReader): DeleteEnvironmentResponse;
}

export namespace DeleteEnvironmentResponse {
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

