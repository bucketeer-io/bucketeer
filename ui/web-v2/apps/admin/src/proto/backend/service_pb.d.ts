// package: bucketeer.backend
// file: proto/backend/service.proto

import * as jspb from "google-protobuf";
import * as google_api_annotations_pb from "../../google/api/annotations_pb";
import * as google_protobuf_field_mask_pb from "google-protobuf/google/protobuf/field_mask_pb";
import * as proto_feature_feature_pb from "../../proto/feature/feature_pb";

export class GetFeatureRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFeatureRequest): GetFeatureRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFeatureRequest;
  static deserializeBinaryFromReader(message: GetFeatureRequest, reader: jspb.BinaryReader): GetFeatureRequest;
}

export namespace GetFeatureRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
  }
}

export class GetFeatureResponse extends jspb.Message {
  hasFeature(): boolean;
  clearFeature(): void;
  getFeature(): proto_feature_feature_pb.Feature | undefined;
  setFeature(value?: proto_feature_feature_pb.Feature): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFeatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFeatureResponse): GetFeatureResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFeatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFeatureResponse;
  static deserializeBinaryFromReader(message: GetFeatureResponse, reader: jspb.BinaryReader): GetFeatureResponse;
}

export namespace GetFeatureResponse {
  export type AsObject = {
    feature?: proto_feature_feature_pb.Feature.AsObject,
  }
}

export class UpdateFeatureRequest extends jspb.Message {
  getComment(): string;
  setComment(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasFieldMask(): boolean;
  clearFieldMask(): void;
  getFieldMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setFieldMask(value?: google_protobuf_field_mask_pb.FieldMask): void;

  hasFeature(): boolean;
  clearFeature(): void;
  getFeature(): proto_feature_feature_pb.Feature | undefined;
  setFeature(value?: proto_feature_feature_pb.Feature): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateFeatureRequest): UpdateFeatureRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateFeatureRequest;
  static deserializeBinaryFromReader(message: UpdateFeatureRequest, reader: jspb.BinaryReader): UpdateFeatureRequest;
}

export namespace UpdateFeatureRequest {
  export type AsObject = {
    comment: string,
    environmentNamespace: string,
    fieldMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    feature?: proto_feature_feature_pb.Feature.AsObject,
  }
}

export class UpdateFeatureResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateFeatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateFeatureResponse): UpdateFeatureResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateFeatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateFeatureResponse;
  static deserializeBinaryFromReader(message: UpdateFeatureResponse, reader: jspb.BinaryReader): UpdateFeatureResponse;
}

export namespace UpdateFeatureResponse {
  export type AsObject = {
  }
}

