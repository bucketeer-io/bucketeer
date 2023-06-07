// package: bucketeer.push
// file: proto/push/service.proto

import * as jspb from "google-protobuf";
import * as proto_push_push_pb from "../../proto/push/push_pb";
import * as proto_push_command_pb from "../../proto/push/command_pb";

export class CreatePushRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_push_command_pb.CreatePushCommand | undefined;
  setCommand(value?: proto_push_command_pb.CreatePushCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePushRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreatePushRequest): CreatePushRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreatePushRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreatePushRequest;
  static deserializeBinaryFromReader(message: CreatePushRequest, reader: jspb.BinaryReader): CreatePushRequest;
}

export namespace CreatePushRequest {
  export type AsObject = {
    environmentNamespace: string,
    command?: proto_push_command_pb.CreatePushCommand.AsObject,
  }
}

export class CreatePushResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePushResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreatePushResponse): CreatePushResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreatePushResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreatePushResponse;
  static deserializeBinaryFromReader(message: CreatePushResponse, reader: jspb.BinaryReader): CreatePushResponse;
}

export namespace CreatePushResponse {
  export type AsObject = {
  }
}

export class ListPushesRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListPushesRequest.OrderByMap[keyof ListPushesRequest.OrderByMap];
  setOrderBy(value: ListPushesRequest.OrderByMap[keyof ListPushesRequest.OrderByMap]): void;

  getOrderDirection(): ListPushesRequest.OrderDirectionMap[keyof ListPushesRequest.OrderDirectionMap];
  setOrderDirection(value: ListPushesRequest.OrderDirectionMap[keyof ListPushesRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPushesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListPushesRequest): ListPushesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListPushesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPushesRequest;
  static deserializeBinaryFromReader(message: ListPushesRequest, reader: jspb.BinaryReader): ListPushesRequest;
}

export namespace ListPushesRequest {
  export type AsObject = {
    environmentNamespace: string,
    pageSize: number,
    cursor: string,
    orderBy: ListPushesRequest.OrderByMap[keyof ListPushesRequest.OrderByMap],
    orderDirection: ListPushesRequest.OrderDirectionMap[keyof ListPushesRequest.OrderDirectionMap],
    searchKeyword: string,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    NAME: 1;
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

export class ListPushesResponse extends jspb.Message {
  clearPushesList(): void;
  getPushesList(): Array<proto_push_push_pb.Push>;
  setPushesList(value: Array<proto_push_push_pb.Push>): void;
  addPushes(value?: proto_push_push_pb.Push, index?: number): proto_push_push_pb.Push;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPushesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListPushesResponse): ListPushesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListPushesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPushesResponse;
  static deserializeBinaryFromReader(message: ListPushesResponse, reader: jspb.BinaryReader): ListPushesResponse;
}

export namespace ListPushesResponse {
  export type AsObject = {
    pushesList: Array<proto_push_push_pb.Push.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class DeletePushRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_push_command_pb.DeletePushCommand | undefined;
  setCommand(value?: proto_push_command_pb.DeletePushCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePushRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePushRequest): DeletePushRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeletePushRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePushRequest;
  static deserializeBinaryFromReader(message: DeletePushRequest, reader: jspb.BinaryReader): DeletePushRequest;
}

export namespace DeletePushRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    command?: proto_push_command_pb.DeletePushCommand.AsObject,
  }
}

export class DeletePushResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePushResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePushResponse): DeletePushResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeletePushResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePushResponse;
  static deserializeBinaryFromReader(message: DeletePushResponse, reader: jspb.BinaryReader): DeletePushResponse;
}

export namespace DeletePushResponse {
  export type AsObject = {
  }
}

export class UpdatePushRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasAddPushTagsCommand(): boolean;
  clearAddPushTagsCommand(): void;
  getAddPushTagsCommand(): proto_push_command_pb.AddPushTagsCommand | undefined;
  setAddPushTagsCommand(value?: proto_push_command_pb.AddPushTagsCommand): void;

  hasDeletePushTagsCommand(): boolean;
  clearDeletePushTagsCommand(): void;
  getDeletePushTagsCommand(): proto_push_command_pb.DeletePushTagsCommand | undefined;
  setDeletePushTagsCommand(value?: proto_push_command_pb.DeletePushTagsCommand): void;

  hasRenamePushCommand(): boolean;
  clearRenamePushCommand(): void;
  getRenamePushCommand(): proto_push_command_pb.RenamePushCommand | undefined;
  setRenamePushCommand(value?: proto_push_command_pb.RenamePushCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePushRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePushRequest): UpdatePushRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdatePushRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePushRequest;
  static deserializeBinaryFromReader(message: UpdatePushRequest, reader: jspb.BinaryReader): UpdatePushRequest;
}

export namespace UpdatePushRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    addPushTagsCommand?: proto_push_command_pb.AddPushTagsCommand.AsObject,
    deletePushTagsCommand?: proto_push_command_pb.DeletePushTagsCommand.AsObject,
    renamePushCommand?: proto_push_command_pb.RenamePushCommand.AsObject,
  }
}

export class UpdatePushResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePushResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePushResponse): UpdatePushResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdatePushResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePushResponse;
  static deserializeBinaryFromReader(message: UpdatePushResponse, reader: jspb.BinaryReader): UpdatePushResponse;
}

export namespace UpdatePushResponse {
  export type AsObject = {
  }
}

