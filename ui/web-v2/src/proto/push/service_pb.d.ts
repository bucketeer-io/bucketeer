// package: bucketeer.push
// file: proto/push/service.proto

import * as jspb from 'google-protobuf';
import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';
import * as google_api_annotations_pb from '../../google/api/annotations_pb';
import * as google_api_field_behavior_pb from '../../google/api/field_behavior_pb';
import * as protoc_gen_openapiv2_options_annotations_pb from '../../protoc-gen-openapiv2/options/annotations_pb';
import * as proto_push_push_pb from '../../proto/push/push_pb';
import * as proto_push_command_pb from '../../proto/push/command_pb';

export class CreatePushRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_push_command_pb.CreatePushCommand | undefined;
  setCommand(value?: proto_push_command_pb.CreatePushCommand): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  getName(): string;
  setName(value: string): void;

  getFcmServiceAccount(): Uint8Array | string;
  getFcmServiceAccount_asU8(): Uint8Array;
  getFcmServiceAccount_asB64(): string;
  setFcmServiceAccount(value: Uint8Array | string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePushRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreatePushRequest
  ): CreatePushRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreatePushRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreatePushRequest;
  static deserializeBinaryFromReader(
    message: CreatePushRequest,
    reader: jspb.BinaryReader
  ): CreatePushRequest;
}

export namespace CreatePushRequest {
  export type AsObject = {
    command?: proto_push_command_pb.CreatePushCommand.AsObject;
    tagsList: Array<string>;
    name: string;
    fcmServiceAccount: Uint8Array | string;
    environmentId: string;
  };
}

export class CreatePushResponse extends jspb.Message {
  hasPush(): boolean;
  clearPush(): void;
  getPush(): proto_push_push_pb.Push | undefined;
  setPush(value?: proto_push_push_pb.Push): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePushResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreatePushResponse
  ): CreatePushResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreatePushResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreatePushResponse;
  static deserializeBinaryFromReader(
    message: CreatePushResponse,
    reader: jspb.BinaryReader
  ): CreatePushResponse;
}

export namespace CreatePushResponse {
  export type AsObject = {
    push?: proto_push_push_pb.Push.AsObject;
  };
}

export class ListPushesRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListPushesRequest.OrderByMap[keyof ListPushesRequest.OrderByMap];
  setOrderBy(
    value: ListPushesRequest.OrderByMap[keyof ListPushesRequest.OrderByMap]
  ): void;

  getOrderDirection(): ListPushesRequest.OrderDirectionMap[keyof ListPushesRequest.OrderDirectionMap];
  setOrderDirection(
    value: ListPushesRequest.OrderDirectionMap[keyof ListPushesRequest.OrderDirectionMap]
  ): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPushesRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListPushesRequest
  ): ListPushesRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListPushesRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListPushesRequest;
  static deserializeBinaryFromReader(
    message: ListPushesRequest,
    reader: jspb.BinaryReader
  ): ListPushesRequest;
}

export namespace ListPushesRequest {
  export type AsObject = {
    pageSize: number;
    cursor: string;
    orderBy: ListPushesRequest.OrderByMap[keyof ListPushesRequest.OrderByMap];
    orderDirection: ListPushesRequest.OrderDirectionMap[keyof ListPushesRequest.OrderDirectionMap];
    searchKeyword: string;
    environmentId: string;
    organizationId: string;
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject;
  };

  export interface OrderByMap {
    DEFAULT: 0;
    NAME: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
    ENVIRONMENT: 4;
    STATE: 5;
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
  addPushes(
    value?: proto_push_push_pb.Push,
    index?: number
  ): proto_push_push_pb.Push;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPushesResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListPushesResponse
  ): ListPushesResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListPushesResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListPushesResponse;
  static deserializeBinaryFromReader(
    message: ListPushesResponse,
    reader: jspb.BinaryReader
  ): ListPushesResponse;
}

export namespace ListPushesResponse {
  export type AsObject = {
    pushesList: Array<proto_push_push_pb.Push.AsObject>;
    cursor: string;
    totalCount: number;
  };
}

export class DeletePushRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_push_command_pb.DeletePushCommand | undefined;
  setCommand(value?: proto_push_command_pb.DeletePushCommand): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePushRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeletePushRequest
  ): DeletePushRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeletePushRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeletePushRequest;
  static deserializeBinaryFromReader(
    message: DeletePushRequest,
    reader: jspb.BinaryReader
  ): DeletePushRequest;
}

export namespace DeletePushRequest {
  export type AsObject = {
    id: string;
    command?: proto_push_command_pb.DeletePushCommand.AsObject;
    environmentId: string;
  };
}

export class DeletePushResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePushResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeletePushResponse
  ): DeletePushResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeletePushResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeletePushResponse;
  static deserializeBinaryFromReader(
    message: DeletePushResponse,
    reader: jspb.BinaryReader
  ): DeletePushResponse;
}

export namespace DeletePushResponse {
  export type AsObject = {};
}

export class UpdatePushRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasAddPushTagsCommand(): boolean;
  clearAddPushTagsCommand(): void;
  getAddPushTagsCommand(): proto_push_command_pb.AddPushTagsCommand | undefined;
  setAddPushTagsCommand(value?: proto_push_command_pb.AddPushTagsCommand): void;

  hasDeletePushTagsCommand(): boolean;
  clearDeletePushTagsCommand(): void;
  getDeletePushTagsCommand():
    | proto_push_command_pb.DeletePushTagsCommand
    | undefined;
  setDeletePushTagsCommand(
    value?: proto_push_command_pb.DeletePushTagsCommand
  ): void;

  hasRenamePushCommand(): boolean;
  clearRenamePushCommand(): void;
  getRenamePushCommand(): proto_push_command_pb.RenamePushCommand | undefined;
  setRenamePushCommand(value?: proto_push_command_pb.RenamePushCommand): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  hasName(): boolean;
  clearName(): void;
  getName(): google_protobuf_wrappers_pb.StringValue | undefined;
  setName(value?: google_protobuf_wrappers_pb.StringValue): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePushRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdatePushRequest
  ): UpdatePushRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdatePushRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePushRequest;
  static deserializeBinaryFromReader(
    message: UpdatePushRequest,
    reader: jspb.BinaryReader
  ): UpdatePushRequest;
}

export namespace UpdatePushRequest {
  export type AsObject = {
    id: string;
    addPushTagsCommand?: proto_push_command_pb.AddPushTagsCommand.AsObject;
    deletePushTagsCommand?: proto_push_command_pb.DeletePushTagsCommand.AsObject;
    renamePushCommand?: proto_push_command_pb.RenamePushCommand.AsObject;
    tagsList: Array<string>;
    name?: google_protobuf_wrappers_pb.StringValue.AsObject;
    environmentId: string;
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject;
  };
}

export class UpdatePushResponse extends jspb.Message {
  hasPush(): boolean;
  clearPush(): void;
  getPush(): proto_push_push_pb.Push | undefined;
  setPush(value?: proto_push_push_pb.Push): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePushResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdatePushResponse
  ): UpdatePushResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdatePushResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePushResponse;
  static deserializeBinaryFromReader(
    message: UpdatePushResponse,
    reader: jspb.BinaryReader
  ): UpdatePushResponse;
}

export namespace UpdatePushResponse {
  export type AsObject = {
    push?: proto_push_push_pb.Push.AsObject;
  };
}

export class GetPushRequest extends jspb.Message {
  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPushRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetPushRequest
  ): GetPushRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetPushRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetPushRequest;
  static deserializeBinaryFromReader(
    message: GetPushRequest,
    reader: jspb.BinaryReader
  ): GetPushRequest;
}

export namespace GetPushRequest {
  export type AsObject = {
    environmentId: string;
    id: string;
  };
}

export class GetPushResponse extends jspb.Message {
  hasPush(): boolean;
  clearPush(): void;
  getPush(): proto_push_push_pb.Push | undefined;
  setPush(value?: proto_push_push_pb.Push): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPushResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetPushResponse
  ): GetPushResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetPushResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetPushResponse;
  static deserializeBinaryFromReader(
    message: GetPushResponse,
    reader: jspb.BinaryReader
  ): GetPushResponse;
}

export namespace GetPushResponse {
  export type AsObject = {
    push?: proto_push_push_pb.Push.AsObject;
  };
}
