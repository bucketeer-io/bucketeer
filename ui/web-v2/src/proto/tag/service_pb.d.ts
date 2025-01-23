// package: bucketeer.tag
// file: proto/tag/service.proto

import * as jspb from 'google-protobuf';
import * as google_api_annotations_pb from '../../google/api/annotations_pb';
import * as google_api_field_behavior_pb from '../../google/api/field_behavior_pb';
import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';
import * as protoc_gen_openapiv2_options_annotations_pb from '../../protoc-gen-openapiv2/options/annotations_pb';
import * as proto_tag_tag_pb from '../../proto/tag/tag_pb';

export class CreateTagRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getEntityType(): proto_tag_tag_pb.Tag.EntityTypeMap[keyof proto_tag_tag_pb.Tag.EntityTypeMap];
  setEntityType(
    value: proto_tag_tag_pb.Tag.EntityTypeMap[keyof proto_tag_tag_pb.Tag.EntityTypeMap]
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTagRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateTagRequest
  ): CreateTagRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateTagRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateTagRequest;
  static deserializeBinaryFromReader(
    message: CreateTagRequest,
    reader: jspb.BinaryReader
  ): CreateTagRequest;
}

export namespace CreateTagRequest {
  export type AsObject = {
    name: string;
    entityType: proto_tag_tag_pb.Tag.EntityTypeMap[keyof proto_tag_tag_pb.Tag.EntityTypeMap];
    environmentId: string;
  };
}

export class CreateTagResponse extends jspb.Message {
  hasTag(): boolean;
  clearTag(): void;
  getTag(): proto_tag_tag_pb.Tag | undefined;
  setTag(value?: proto_tag_tag_pb.Tag): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTagResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateTagResponse
  ): CreateTagResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateTagResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateTagResponse;
  static deserializeBinaryFromReader(
    message: CreateTagResponse,
    reader: jspb.BinaryReader
  ): CreateTagResponse;
}

export namespace CreateTagResponse {
  export type AsObject = {
    tag?: proto_tag_tag_pb.Tag.AsObject;
  };
}

export class DeleteTagRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteTagRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteTagRequest
  ): DeleteTagRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteTagRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteTagRequest;
  static deserializeBinaryFromReader(
    message: DeleteTagRequest,
    reader: jspb.BinaryReader
  ): DeleteTagRequest;
}

export namespace DeleteTagRequest {
  export type AsObject = {
    id: string;
    environmentId: string;
  };
}

export class DeleteTagResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteTagResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteTagResponse
  ): DeleteTagResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteTagResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteTagResponse;
  static deserializeBinaryFromReader(
    message: DeleteTagResponse,
    reader: jspb.BinaryReader
  ): DeleteTagResponse;
}

export namespace DeleteTagResponse {
  export type AsObject = {};
}

export class ListTagsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListTagsRequest.OrderByMap[keyof ListTagsRequest.OrderByMap];
  setOrderBy(
    value: ListTagsRequest.OrderByMap[keyof ListTagsRequest.OrderByMap]
  ): void;

  getOrderDirection(): ListTagsRequest.OrderDirectionMap[keyof ListTagsRequest.OrderDirectionMap];
  setOrderDirection(
    value: ListTagsRequest.OrderDirectionMap[keyof ListTagsRequest.OrderDirectionMap]
  ): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  getEntityType(): proto_tag_tag_pb.Tag.EntityTypeMap[keyof proto_tag_tag_pb.Tag.EntityTypeMap];
  setEntityType(
    value: proto_tag_tag_pb.Tag.EntityTypeMap[keyof proto_tag_tag_pb.Tag.EntityTypeMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTagsRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListTagsRequest
  ): ListTagsRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListTagsRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListTagsRequest;
  static deserializeBinaryFromReader(
    message: ListTagsRequest,
    reader: jspb.BinaryReader
  ): ListTagsRequest;
}

export namespace ListTagsRequest {
  export type AsObject = {
    pageSize: number;
    cursor: string;
    orderBy: ListTagsRequest.OrderByMap[keyof ListTagsRequest.OrderByMap];
    orderDirection: ListTagsRequest.OrderDirectionMap[keyof ListTagsRequest.OrderDirectionMap];
    searchKeyword: string;
    environmentId: string;
    organizationId: string;
    entityType: proto_tag_tag_pb.Tag.EntityTypeMap[keyof proto_tag_tag_pb.Tag.EntityTypeMap];
  };

  export interface OrderByMap {
    DEFAULT: 0;
    NAME: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
    ENTITY_TYPE: 4;
    ENVIRONMENT: 5;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListTagsResponse extends jspb.Message {
  clearTagsList(): void;
  getTagsList(): Array<proto_tag_tag_pb.Tag>;
  setTagsList(value: Array<proto_tag_tag_pb.Tag>): void;
  addTags(value?: proto_tag_tag_pb.Tag, index?: number): proto_tag_tag_pb.Tag;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTagsResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListTagsResponse
  ): ListTagsResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListTagsResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListTagsResponse;
  static deserializeBinaryFromReader(
    message: ListTagsResponse,
    reader: jspb.BinaryReader
  ): ListTagsResponse;
}

export namespace ListTagsResponse {
  export type AsObject = {
    tagsList: Array<proto_tag_tag_pb.Tag.AsObject>;
    cursor: string;
    totalCount: number;
  };
}
