// package: bucketeer.feature
// file: proto/feature/service.proto

import * as jspb from "google-protobuf";
import * as google_api_annotations_pb from "../../google/api/annotations_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as proto_feature_command_pb from "../../proto/feature/command_pb";
import * as proto_feature_feature_pb from "../../proto/feature/feature_pb";
import * as proto_feature_evaluation_pb from "../../proto/feature/evaluation_pb";
import * as proto_user_user_pb from "../../proto/user/user_pb";
import * as proto_feature_segment_pb from "../../proto/feature/segment_pb";
import * as proto_feature_flag_trigger_pb from "../../proto/feature/flag_trigger_pb";

export class GetFeatureRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

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
    id: string,
    environmentNamespace: string,
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

export class GetFeaturesRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  clearIdsList(): void;
  getIdsList(): Array<string>;
  setIdsList(value: Array<string>): void;
  addIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFeaturesRequest): GetFeaturesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFeaturesRequest;
  static deserializeBinaryFromReader(message: GetFeaturesRequest, reader: jspb.BinaryReader): GetFeaturesRequest;
}

export namespace GetFeaturesRequest {
  export type AsObject = {
    environmentNamespace: string,
    idsList: Array<string>,
  }
}

export class GetFeaturesResponse extends jspb.Message {
  clearFeaturesList(): void;
  getFeaturesList(): Array<proto_feature_feature_pb.Feature>;
  setFeaturesList(value: Array<proto_feature_feature_pb.Feature>): void;
  addFeatures(value?: proto_feature_feature_pb.Feature, index?: number): proto_feature_feature_pb.Feature;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFeaturesResponse): GetFeaturesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFeaturesResponse;
  static deserializeBinaryFromReader(message: GetFeaturesResponse, reader: jspb.BinaryReader): GetFeaturesResponse;
}

export namespace GetFeaturesResponse {
  export type AsObject = {
    featuresList: Array<proto_feature_feature_pb.Feature.AsObject>,
  }
}

export class ListFeaturesRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  getOrderBy(): ListFeaturesRequest.OrderByMap[keyof ListFeaturesRequest.OrderByMap];
  setOrderBy(value: ListFeaturesRequest.OrderByMap[keyof ListFeaturesRequest.OrderByMap]): void;

  getOrderDirection(): ListFeaturesRequest.OrderDirectionMap[keyof ListFeaturesRequest.OrderDirectionMap];
  setOrderDirection(value: ListFeaturesRequest.OrderDirectionMap[keyof ListFeaturesRequest.OrderDirectionMap]): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getMaintainer(): string;
  setMaintainer(value: string): void;

  hasEnabled(): boolean;
  clearEnabled(): void;
  getEnabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setEnabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasHasExperiment(): boolean;
  clearHasExperiment(): void;
  getHasExperiment(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setHasExperiment(value?: google_protobuf_wrappers_pb.BoolValue): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasArchived(): boolean;
  clearArchived(): void;
  getArchived(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setArchived(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasHasPrerequisites(): boolean;
  clearHasPrerequisites(): void;
  getHasPrerequisites(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setHasPrerequisites(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFeaturesRequest): ListFeaturesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFeaturesRequest;
  static deserializeBinaryFromReader(message: ListFeaturesRequest, reader: jspb.BinaryReader): ListFeaturesRequest;
}

export namespace ListFeaturesRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    tagsList: Array<string>,
    orderBy: ListFeaturesRequest.OrderByMap[keyof ListFeaturesRequest.OrderByMap],
    orderDirection: ListFeaturesRequest.OrderDirectionMap[keyof ListFeaturesRequest.OrderDirectionMap],
    environmentNamespace: string,
    maintainer: string,
    enabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    hasExperiment?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    searchKeyword: string,
    archived?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    hasPrerequisites?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    NAME: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
    TAGS: 4;
    ENABLED: 5;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListFeaturesResponse extends jspb.Message {
  clearFeaturesList(): void;
  getFeaturesList(): Array<proto_feature_feature_pb.Feature>;
  setFeaturesList(value: Array<proto_feature_feature_pb.Feature>): void;
  addFeatures(value?: proto_feature_feature_pb.Feature, index?: number): proto_feature_feature_pb.Feature;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFeaturesResponse): ListFeaturesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFeaturesResponse;
  static deserializeBinaryFromReader(message: ListFeaturesResponse, reader: jspb.BinaryReader): ListFeaturesResponse;
}

export namespace ListFeaturesResponse {
  export type AsObject = {
    featuresList: Array<proto_feature_feature_pb.Feature.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class ListEnabledFeaturesRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnabledFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnabledFeaturesRequest): ListEnabledFeaturesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnabledFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnabledFeaturesRequest;
  static deserializeBinaryFromReader(message: ListEnabledFeaturesRequest, reader: jspb.BinaryReader): ListEnabledFeaturesRequest;
}

export namespace ListEnabledFeaturesRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    tagsList: Array<string>,
    environmentNamespace: string,
  }
}

export class ListEnabledFeaturesResponse extends jspb.Message {
  clearFeaturesList(): void;
  getFeaturesList(): Array<proto_feature_feature_pb.Feature>;
  setFeaturesList(value: Array<proto_feature_feature_pb.Feature>): void;
  addFeatures(value?: proto_feature_feature_pb.Feature, index?: number): proto_feature_feature_pb.Feature;

  getCursor(): string;
  setCursor(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnabledFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnabledFeaturesResponse): ListEnabledFeaturesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnabledFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnabledFeaturesResponse;
  static deserializeBinaryFromReader(message: ListEnabledFeaturesResponse, reader: jspb.BinaryReader): ListEnabledFeaturesResponse;
}

export namespace ListEnabledFeaturesResponse {
  export type AsObject = {
    featuresList: Array<proto_feature_feature_pb.Feature.AsObject>,
    cursor: string,
  }
}

export class CreateFeatureRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.CreateFeatureCommand | undefined;
  setCommand(value?: proto_feature_command_pb.CreateFeatureCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateFeatureRequest): CreateFeatureRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateFeatureRequest;
  static deserializeBinaryFromReader(message: CreateFeatureRequest, reader: jspb.BinaryReader): CreateFeatureRequest;
}

export namespace CreateFeatureRequest {
  export type AsObject = {
    command?: proto_feature_command_pb.CreateFeatureCommand.AsObject,
    environmentNamespace: string,
  }
}

export class CreateFeatureResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateFeatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateFeatureResponse): CreateFeatureResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateFeatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateFeatureResponse;
  static deserializeBinaryFromReader(message: CreateFeatureResponse, reader: jspb.BinaryReader): CreateFeatureResponse;
}

export namespace CreateFeatureResponse {
  export type AsObject = {
  }
}

export class EnableFeatureRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.EnableFeatureCommand | undefined;
  setCommand(value?: proto_feature_command_pb.EnableFeatureCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getComment(): string;
  setComment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableFeatureRequest): EnableFeatureRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableFeatureRequest;
  static deserializeBinaryFromReader(message: EnableFeatureRequest, reader: jspb.BinaryReader): EnableFeatureRequest;
}

export namespace EnableFeatureRequest {
  export type AsObject = {
    id: string,
    command?: proto_feature_command_pb.EnableFeatureCommand.AsObject,
    environmentNamespace: string,
    comment: string,
  }
}

export class EnableFeatureResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableFeatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableFeatureResponse): EnableFeatureResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableFeatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableFeatureResponse;
  static deserializeBinaryFromReader(message: EnableFeatureResponse, reader: jspb.BinaryReader): EnableFeatureResponse;
}

export namespace EnableFeatureResponse {
  export type AsObject = {
  }
}

export class DisableFeatureRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.DisableFeatureCommand | undefined;
  setCommand(value?: proto_feature_command_pb.DisableFeatureCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getComment(): string;
  setComment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableFeatureRequest): DisableFeatureRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableFeatureRequest;
  static deserializeBinaryFromReader(message: DisableFeatureRequest, reader: jspb.BinaryReader): DisableFeatureRequest;
}

export namespace DisableFeatureRequest {
  export type AsObject = {
    id: string,
    command?: proto_feature_command_pb.DisableFeatureCommand.AsObject,
    environmentNamespace: string,
    comment: string,
  }
}

export class DisableFeatureResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableFeatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableFeatureResponse): DisableFeatureResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableFeatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableFeatureResponse;
  static deserializeBinaryFromReader(message: DisableFeatureResponse, reader: jspb.BinaryReader): DisableFeatureResponse;
}

export namespace DisableFeatureResponse {
  export type AsObject = {
  }
}

export class ArchiveFeatureRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.ArchiveFeatureCommand | undefined;
  setCommand(value?: proto_feature_command_pb.ArchiveFeatureCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getComment(): string;
  setComment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveFeatureRequest): ArchiveFeatureRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveFeatureRequest;
  static deserializeBinaryFromReader(message: ArchiveFeatureRequest, reader: jspb.BinaryReader): ArchiveFeatureRequest;
}

export namespace ArchiveFeatureRequest {
  export type AsObject = {
    id: string,
    command?: proto_feature_command_pb.ArchiveFeatureCommand.AsObject,
    environmentNamespace: string,
    comment: string,
  }
}

export class ArchiveFeatureResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveFeatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveFeatureResponse): ArchiveFeatureResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveFeatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveFeatureResponse;
  static deserializeBinaryFromReader(message: ArchiveFeatureResponse, reader: jspb.BinaryReader): ArchiveFeatureResponse;
}

export namespace ArchiveFeatureResponse {
  export type AsObject = {
  }
}

export class UnarchiveFeatureRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.UnarchiveFeatureCommand | undefined;
  setCommand(value?: proto_feature_command_pb.UnarchiveFeatureCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getComment(): string;
  setComment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnarchiveFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UnarchiveFeatureRequest): UnarchiveFeatureRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UnarchiveFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnarchiveFeatureRequest;
  static deserializeBinaryFromReader(message: UnarchiveFeatureRequest, reader: jspb.BinaryReader): UnarchiveFeatureRequest;
}

export namespace UnarchiveFeatureRequest {
  export type AsObject = {
    id: string,
    command?: proto_feature_command_pb.UnarchiveFeatureCommand.AsObject,
    environmentNamespace: string,
    comment: string,
  }
}

export class UnarchiveFeatureResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnarchiveFeatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UnarchiveFeatureResponse): UnarchiveFeatureResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UnarchiveFeatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnarchiveFeatureResponse;
  static deserializeBinaryFromReader(message: UnarchiveFeatureResponse, reader: jspb.BinaryReader): UnarchiveFeatureResponse;
}

export namespace UnarchiveFeatureResponse {
  export type AsObject = {
  }
}

export class DeleteFeatureRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.DeleteFeatureCommand | undefined;
  setCommand(value?: proto_feature_command_pb.DeleteFeatureCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getComment(): string;
  setComment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteFeatureRequest): DeleteFeatureRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteFeatureRequest;
  static deserializeBinaryFromReader(message: DeleteFeatureRequest, reader: jspb.BinaryReader): DeleteFeatureRequest;
}

export namespace DeleteFeatureRequest {
  export type AsObject = {
    id: string,
    command?: proto_feature_command_pb.DeleteFeatureCommand.AsObject,
    environmentNamespace: string,
    comment: string,
  }
}

export class DeleteFeatureResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteFeatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteFeatureResponse): DeleteFeatureResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteFeatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteFeatureResponse;
  static deserializeBinaryFromReader(message: DeleteFeatureResponse, reader: jspb.BinaryReader): DeleteFeatureResponse;
}

export namespace DeleteFeatureResponse {
  export type AsObject = {
  }
}

export class UpdateFeatureDetailsRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasRenameFeatureCommand(): boolean;
  clearRenameFeatureCommand(): void;
  getRenameFeatureCommand(): proto_feature_command_pb.RenameFeatureCommand | undefined;
  setRenameFeatureCommand(value?: proto_feature_command_pb.RenameFeatureCommand): void;

  hasChangeDescriptionCommand(): boolean;
  clearChangeDescriptionCommand(): void;
  getChangeDescriptionCommand(): proto_feature_command_pb.ChangeDescriptionCommand | undefined;
  setChangeDescriptionCommand(value?: proto_feature_command_pb.ChangeDescriptionCommand): void;

  clearAddTagCommandsList(): void;
  getAddTagCommandsList(): Array<proto_feature_command_pb.AddTagCommand>;
  setAddTagCommandsList(value: Array<proto_feature_command_pb.AddTagCommand>): void;
  addAddTagCommands(value?: proto_feature_command_pb.AddTagCommand, index?: number): proto_feature_command_pb.AddTagCommand;

  clearRemoveTagCommandsList(): void;
  getRemoveTagCommandsList(): Array<proto_feature_command_pb.RemoveTagCommand>;
  setRemoveTagCommandsList(value: Array<proto_feature_command_pb.RemoveTagCommand>): void;
  addRemoveTagCommands(value?: proto_feature_command_pb.RemoveTagCommand, index?: number): proto_feature_command_pb.RemoveTagCommand;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getComment(): string;
  setComment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateFeatureDetailsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateFeatureDetailsRequest): UpdateFeatureDetailsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateFeatureDetailsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateFeatureDetailsRequest;
  static deserializeBinaryFromReader(message: UpdateFeatureDetailsRequest, reader: jspb.BinaryReader): UpdateFeatureDetailsRequest;
}

export namespace UpdateFeatureDetailsRequest {
  export type AsObject = {
    id: string,
    renameFeatureCommand?: proto_feature_command_pb.RenameFeatureCommand.AsObject,
    changeDescriptionCommand?: proto_feature_command_pb.ChangeDescriptionCommand.AsObject,
    addTagCommandsList: Array<proto_feature_command_pb.AddTagCommand.AsObject>,
    removeTagCommandsList: Array<proto_feature_command_pb.RemoveTagCommand.AsObject>,
    environmentNamespace: string,
    comment: string,
  }
}

export class UpdateFeatureDetailsResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateFeatureDetailsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateFeatureDetailsResponse): UpdateFeatureDetailsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateFeatureDetailsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateFeatureDetailsResponse;
  static deserializeBinaryFromReader(message: UpdateFeatureDetailsResponse, reader: jspb.BinaryReader): UpdateFeatureDetailsResponse;
}

export namespace UpdateFeatureDetailsResponse {
  export type AsObject = {
  }
}

export class UpdateFeatureVariationsRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  clearCommandsList(): void;
  getCommandsList(): Array<proto_feature_command_pb.Command>;
  setCommandsList(value: Array<proto_feature_command_pb.Command>): void;
  addCommands(value?: proto_feature_command_pb.Command, index?: number): proto_feature_command_pb.Command;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getComment(): string;
  setComment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateFeatureVariationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateFeatureVariationsRequest): UpdateFeatureVariationsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateFeatureVariationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateFeatureVariationsRequest;
  static deserializeBinaryFromReader(message: UpdateFeatureVariationsRequest, reader: jspb.BinaryReader): UpdateFeatureVariationsRequest;
}

export namespace UpdateFeatureVariationsRequest {
  export type AsObject = {
    id: string,
    commandsList: Array<proto_feature_command_pb.Command.AsObject>,
    environmentNamespace: string,
    comment: string,
  }
}

export class UpdateFeatureVariationsResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateFeatureVariationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateFeatureVariationsResponse): UpdateFeatureVariationsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateFeatureVariationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateFeatureVariationsResponse;
  static deserializeBinaryFromReader(message: UpdateFeatureVariationsResponse, reader: jspb.BinaryReader): UpdateFeatureVariationsResponse;
}

export namespace UpdateFeatureVariationsResponse {
  export type AsObject = {
  }
}

export class UpdateFeatureTargetingRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  clearCommandsList(): void;
  getCommandsList(): Array<proto_feature_command_pb.Command>;
  setCommandsList(value: Array<proto_feature_command_pb.Command>): void;
  addCommands(value?: proto_feature_command_pb.Command, index?: number): proto_feature_command_pb.Command;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getComment(): string;
  setComment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateFeatureTargetingRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateFeatureTargetingRequest): UpdateFeatureTargetingRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateFeatureTargetingRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateFeatureTargetingRequest;
  static deserializeBinaryFromReader(message: UpdateFeatureTargetingRequest, reader: jspb.BinaryReader): UpdateFeatureTargetingRequest;
}

export namespace UpdateFeatureTargetingRequest {
  export type AsObject = {
    id: string,
    commandsList: Array<proto_feature_command_pb.Command.AsObject>,
    environmentNamespace: string,
    comment: string,
  }
}

export class UpdateFeatureTargetingResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateFeatureTargetingResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateFeatureTargetingResponse): UpdateFeatureTargetingResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateFeatureTargetingResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateFeatureTargetingResponse;
  static deserializeBinaryFromReader(message: UpdateFeatureTargetingResponse, reader: jspb.BinaryReader): UpdateFeatureTargetingResponse;
}

export namespace UpdateFeatureTargetingResponse {
  export type AsObject = {
  }
}

export class CloneFeatureRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.CloneFeatureCommand | undefined;
  setCommand(value?: proto_feature_command_pb.CloneFeatureCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloneFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CloneFeatureRequest): CloneFeatureRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CloneFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloneFeatureRequest;
  static deserializeBinaryFromReader(message: CloneFeatureRequest, reader: jspb.BinaryReader): CloneFeatureRequest;
}

export namespace CloneFeatureRequest {
  export type AsObject = {
    id: string,
    command?: proto_feature_command_pb.CloneFeatureCommand.AsObject,
    environmentNamespace: string,
  }
}

export class CloneFeatureResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloneFeatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CloneFeatureResponse): CloneFeatureResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CloneFeatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloneFeatureResponse;
  static deserializeBinaryFromReader(message: CloneFeatureResponse, reader: jspb.BinaryReader): CloneFeatureResponse;
}

export namespace CloneFeatureResponse {
  export type AsObject = {
  }
}

export class CreateSegmentRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.CreateSegmentCommand | undefined;
  setCommand(value?: proto_feature_command_pb.CreateSegmentCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSegmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateSegmentRequest): CreateSegmentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateSegmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateSegmentRequest;
  static deserializeBinaryFromReader(message: CreateSegmentRequest, reader: jspb.BinaryReader): CreateSegmentRequest;
}

export namespace CreateSegmentRequest {
  export type AsObject = {
    command?: proto_feature_command_pb.CreateSegmentCommand.AsObject,
    environmentNamespace: string,
  }
}

export class CreateSegmentResponse extends jspb.Message {
  hasSegment(): boolean;
  clearSegment(): void;
  getSegment(): proto_feature_segment_pb.Segment | undefined;
  setSegment(value?: proto_feature_segment_pb.Segment): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSegmentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateSegmentResponse): CreateSegmentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateSegmentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateSegmentResponse;
  static deserializeBinaryFromReader(message: CreateSegmentResponse, reader: jspb.BinaryReader): CreateSegmentResponse;
}

export namespace CreateSegmentResponse {
  export type AsObject = {
    segment?: proto_feature_segment_pb.Segment.AsObject,
  }
}

export class GetSegmentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSegmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSegmentRequest): GetSegmentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSegmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSegmentRequest;
  static deserializeBinaryFromReader(message: GetSegmentRequest, reader: jspb.BinaryReader): GetSegmentRequest;
}

export namespace GetSegmentRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
  }
}

export class GetSegmentResponse extends jspb.Message {
  hasSegment(): boolean;
  clearSegment(): void;
  getSegment(): proto_feature_segment_pb.Segment | undefined;
  setSegment(value?: proto_feature_segment_pb.Segment): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSegmentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSegmentResponse): GetSegmentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSegmentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSegmentResponse;
  static deserializeBinaryFromReader(message: GetSegmentResponse, reader: jspb.BinaryReader): GetSegmentResponse;
}

export namespace GetSegmentResponse {
  export type AsObject = {
    segment?: proto_feature_segment_pb.Segment.AsObject,
  }
}

export class ListSegmentsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getOrderBy(): ListSegmentsRequest.OrderByMap[keyof ListSegmentsRequest.OrderByMap];
  setOrderBy(value: ListSegmentsRequest.OrderByMap[keyof ListSegmentsRequest.OrderByMap]): void;

  getOrderDirection(): ListSegmentsRequest.OrderDirectionMap[keyof ListSegmentsRequest.OrderDirectionMap];
  setOrderDirection(value: ListSegmentsRequest.OrderDirectionMap[keyof ListSegmentsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setStatus(value?: google_protobuf_wrappers_pb.Int32Value): void;

  hasIsInUseStatus(): boolean;
  clearIsInUseStatus(): void;
  getIsInUseStatus(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setIsInUseStatus(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSegmentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListSegmentsRequest): ListSegmentsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListSegmentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSegmentsRequest;
  static deserializeBinaryFromReader(message: ListSegmentsRequest, reader: jspb.BinaryReader): ListSegmentsRequest;
}

export namespace ListSegmentsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    environmentNamespace: string,
    orderBy: ListSegmentsRequest.OrderByMap[keyof ListSegmentsRequest.OrderByMap],
    orderDirection: ListSegmentsRequest.OrderDirectionMap[keyof ListSegmentsRequest.OrderDirectionMap],
    searchKeyword: string,
    status?: google_protobuf_wrappers_pb.Int32Value.AsObject,
    isInUseStatus?: google_protobuf_wrappers_pb.BoolValue.AsObject,
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

export class ListSegmentsResponse extends jspb.Message {
  clearSegmentsList(): void;
  getSegmentsList(): Array<proto_feature_segment_pb.Segment>;
  setSegmentsList(value: Array<proto_feature_segment_pb.Segment>): void;
  addSegments(value?: proto_feature_segment_pb.Segment, index?: number): proto_feature_segment_pb.Segment;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSegmentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListSegmentsResponse): ListSegmentsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListSegmentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSegmentsResponse;
  static deserializeBinaryFromReader(message: ListSegmentsResponse, reader: jspb.BinaryReader): ListSegmentsResponse;
}

export namespace ListSegmentsResponse {
  export type AsObject = {
    segmentsList: Array<proto_feature_segment_pb.Segment.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class DeleteSegmentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.DeleteSegmentCommand | undefined;
  setCommand(value?: proto_feature_command_pb.DeleteSegmentCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSegmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSegmentRequest): DeleteSegmentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteSegmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSegmentRequest;
  static deserializeBinaryFromReader(message: DeleteSegmentRequest, reader: jspb.BinaryReader): DeleteSegmentRequest;
}

export namespace DeleteSegmentRequest {
  export type AsObject = {
    id: string,
    command?: proto_feature_command_pb.DeleteSegmentCommand.AsObject,
    environmentNamespace: string,
  }
}

export class DeleteSegmentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSegmentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSegmentResponse): DeleteSegmentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteSegmentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSegmentResponse;
  static deserializeBinaryFromReader(message: DeleteSegmentResponse, reader: jspb.BinaryReader): DeleteSegmentResponse;
}

export namespace DeleteSegmentResponse {
  export type AsObject = {
  }
}

export class UpdateSegmentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  clearCommandsList(): void;
  getCommandsList(): Array<proto_feature_command_pb.Command>;
  setCommandsList(value: Array<proto_feature_command_pb.Command>): void;
  addCommands(value?: proto_feature_command_pb.Command, index?: number): proto_feature_command_pb.Command;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSegmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSegmentRequest): UpdateSegmentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateSegmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSegmentRequest;
  static deserializeBinaryFromReader(message: UpdateSegmentRequest, reader: jspb.BinaryReader): UpdateSegmentRequest;
}

export namespace UpdateSegmentRequest {
  export type AsObject = {
    id: string,
    commandsList: Array<proto_feature_command_pb.Command.AsObject>,
    environmentNamespace: string,
  }
}

export class UpdateSegmentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSegmentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSegmentResponse): UpdateSegmentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateSegmentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSegmentResponse;
  static deserializeBinaryFromReader(message: UpdateSegmentResponse, reader: jspb.BinaryReader): UpdateSegmentResponse;
}

export namespace UpdateSegmentResponse {
  export type AsObject = {
  }
}

export class AddSegmentUserRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.AddSegmentUserCommand | undefined;
  setCommand(value?: proto_feature_command_pb.AddSegmentUserCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSegmentUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddSegmentUserRequest): AddSegmentUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddSegmentUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSegmentUserRequest;
  static deserializeBinaryFromReader(message: AddSegmentUserRequest, reader: jspb.BinaryReader): AddSegmentUserRequest;
}

export namespace AddSegmentUserRequest {
  export type AsObject = {
    id: string,
    command?: proto_feature_command_pb.AddSegmentUserCommand.AsObject,
    environmentNamespace: string,
  }
}

export class AddSegmentUserResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSegmentUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddSegmentUserResponse): AddSegmentUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddSegmentUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSegmentUserResponse;
  static deserializeBinaryFromReader(message: AddSegmentUserResponse, reader: jspb.BinaryReader): AddSegmentUserResponse;
}

export namespace AddSegmentUserResponse {
  export type AsObject = {
  }
}

export class DeleteSegmentUserRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.DeleteSegmentUserCommand | undefined;
  setCommand(value?: proto_feature_command_pb.DeleteSegmentUserCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSegmentUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSegmentUserRequest): DeleteSegmentUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteSegmentUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSegmentUserRequest;
  static deserializeBinaryFromReader(message: DeleteSegmentUserRequest, reader: jspb.BinaryReader): DeleteSegmentUserRequest;
}

export namespace DeleteSegmentUserRequest {
  export type AsObject = {
    id: string,
    command?: proto_feature_command_pb.DeleteSegmentUserCommand.AsObject,
    environmentNamespace: string,
  }
}

export class DeleteSegmentUserResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSegmentUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSegmentUserResponse): DeleteSegmentUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteSegmentUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSegmentUserResponse;
  static deserializeBinaryFromReader(message: DeleteSegmentUserResponse, reader: jspb.BinaryReader): DeleteSegmentUserResponse;
}

export namespace DeleteSegmentUserResponse {
  export type AsObject = {
  }
}

export class GetSegmentUserRequest extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getState(): proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
  setState(value: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap]): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSegmentUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSegmentUserRequest): GetSegmentUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSegmentUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSegmentUserRequest;
  static deserializeBinaryFromReader(message: GetSegmentUserRequest, reader: jspb.BinaryReader): GetSegmentUserRequest;
}

export namespace GetSegmentUserRequest {
  export type AsObject = {
    segmentId: string,
    userId: string,
    state: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap],
    environmentNamespace: string,
  }
}

export class GetSegmentUserResponse extends jspb.Message {
  hasUser(): boolean;
  clearUser(): void;
  getUser(): proto_feature_segment_pb.SegmentUser | undefined;
  setUser(value?: proto_feature_segment_pb.SegmentUser): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSegmentUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSegmentUserResponse): GetSegmentUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSegmentUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSegmentUserResponse;
  static deserializeBinaryFromReader(message: GetSegmentUserResponse, reader: jspb.BinaryReader): GetSegmentUserResponse;
}

export namespace GetSegmentUserResponse {
  export type AsObject = {
    user?: proto_feature_segment_pb.SegmentUser.AsObject,
  }
}

export class ListSegmentUsersRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getSegmentId(): string;
  setSegmentId(value: string): void;

  hasState(): boolean;
  clearState(): void;
  getState(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setState(value?: google_protobuf_wrappers_pb.Int32Value): void;

  getUserId(): string;
  setUserId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSegmentUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListSegmentUsersRequest): ListSegmentUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListSegmentUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSegmentUsersRequest;
  static deserializeBinaryFromReader(message: ListSegmentUsersRequest, reader: jspb.BinaryReader): ListSegmentUsersRequest;
}

export namespace ListSegmentUsersRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    segmentId: string,
    state?: google_protobuf_wrappers_pb.Int32Value.AsObject,
    userId: string,
    environmentNamespace: string,
  }
}

export class ListSegmentUsersResponse extends jspb.Message {
  clearUsersList(): void;
  getUsersList(): Array<proto_feature_segment_pb.SegmentUser>;
  setUsersList(value: Array<proto_feature_segment_pb.SegmentUser>): void;
  addUsers(value?: proto_feature_segment_pb.SegmentUser, index?: number): proto_feature_segment_pb.SegmentUser;

  getCursor(): string;
  setCursor(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSegmentUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListSegmentUsersResponse): ListSegmentUsersResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListSegmentUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSegmentUsersResponse;
  static deserializeBinaryFromReader(message: ListSegmentUsersResponse, reader: jspb.BinaryReader): ListSegmentUsersResponse;
}

export namespace ListSegmentUsersResponse {
  export type AsObject = {
    usersList: Array<proto_feature_segment_pb.SegmentUser.AsObject>,
    cursor: string,
  }
}

export class BulkUploadSegmentUsersRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getSegmentId(): string;
  setSegmentId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_feature_command_pb.BulkUploadSegmentUsersCommand | undefined;
  setCommand(value?: proto_feature_command_pb.BulkUploadSegmentUsersCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkUploadSegmentUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkUploadSegmentUsersRequest): BulkUploadSegmentUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BulkUploadSegmentUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkUploadSegmentUsersRequest;
  static deserializeBinaryFromReader(message: BulkUploadSegmentUsersRequest, reader: jspb.BinaryReader): BulkUploadSegmentUsersRequest;
}

export namespace BulkUploadSegmentUsersRequest {
  export type AsObject = {
    environmentNamespace: string,
    segmentId: string,
    command?: proto_feature_command_pb.BulkUploadSegmentUsersCommand.AsObject,
  }
}

export class BulkUploadSegmentUsersResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkUploadSegmentUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkUploadSegmentUsersResponse): BulkUploadSegmentUsersResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BulkUploadSegmentUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkUploadSegmentUsersResponse;
  static deserializeBinaryFromReader(message: BulkUploadSegmentUsersResponse, reader: jspb.BinaryReader): BulkUploadSegmentUsersResponse;
}

export namespace BulkUploadSegmentUsersResponse {
  export type AsObject = {
  }
}

export class BulkDownloadSegmentUsersRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getSegmentId(): string;
  setSegmentId(value: string): void;

  getState(): proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
  setState(value: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkDownloadSegmentUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkDownloadSegmentUsersRequest): BulkDownloadSegmentUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BulkDownloadSegmentUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkDownloadSegmentUsersRequest;
  static deserializeBinaryFromReader(message: BulkDownloadSegmentUsersRequest, reader: jspb.BinaryReader): BulkDownloadSegmentUsersRequest;
}

export namespace BulkDownloadSegmentUsersRequest {
  export type AsObject = {
    environmentNamespace: string,
    segmentId: string,
    state: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap],
  }
}

export class BulkDownloadSegmentUsersResponse extends jspb.Message {
  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkDownloadSegmentUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkDownloadSegmentUsersResponse): BulkDownloadSegmentUsersResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BulkDownloadSegmentUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkDownloadSegmentUsersResponse;
  static deserializeBinaryFromReader(message: BulkDownloadSegmentUsersResponse, reader: jspb.BinaryReader): BulkDownloadSegmentUsersResponse;
}

export namespace BulkDownloadSegmentUsersResponse {
  export type AsObject = {
    data: Uint8Array | string,
  }
}

export class EvaluateFeaturesRequest extends jspb.Message {
  hasUser(): boolean;
  clearUser(): void;
  getUser(): proto_user_user_pb.User | undefined;
  setUser(value?: proto_user_user_pb.User): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getTag(): string;
  setTag(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EvaluateFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EvaluateFeaturesRequest): EvaluateFeaturesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EvaluateFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EvaluateFeaturesRequest;
  static deserializeBinaryFromReader(message: EvaluateFeaturesRequest, reader: jspb.BinaryReader): EvaluateFeaturesRequest;
}

export namespace EvaluateFeaturesRequest {
  export type AsObject = {
    user?: proto_user_user_pb.User.AsObject,
    environmentNamespace: string,
    tag: string,
    featureId: string,
  }
}

export class EvaluateFeaturesResponse extends jspb.Message {
  hasUserEvaluations(): boolean;
  clearUserEvaluations(): void;
  getUserEvaluations(): proto_feature_evaluation_pb.UserEvaluations | undefined;
  setUserEvaluations(value?: proto_feature_evaluation_pb.UserEvaluations): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EvaluateFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EvaluateFeaturesResponse): EvaluateFeaturesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EvaluateFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EvaluateFeaturesResponse;
  static deserializeBinaryFromReader(message: EvaluateFeaturesResponse, reader: jspb.BinaryReader): EvaluateFeaturesResponse;
}

export namespace EvaluateFeaturesResponse {
  export type AsObject = {
    userEvaluations?: proto_feature_evaluation_pb.UserEvaluations.AsObject,
  }
}

export class ListTagsRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListTagsRequest.OrderByMap[keyof ListTagsRequest.OrderByMap];
  setOrderBy(value: ListTagsRequest.OrderByMap[keyof ListTagsRequest.OrderByMap]): void;

  getOrderDirection(): ListTagsRequest.OrderDirectionMap[keyof ListTagsRequest.OrderDirectionMap];
  setOrderDirection(value: ListTagsRequest.OrderDirectionMap[keyof ListTagsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTagsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListTagsRequest): ListTagsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListTagsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTagsRequest;
  static deserializeBinaryFromReader(message: ListTagsRequest, reader: jspb.BinaryReader): ListTagsRequest;
}

export namespace ListTagsRequest {
  export type AsObject = {
    environmentNamespace: string,
    pageSize: number,
    cursor: string,
    orderBy: ListTagsRequest.OrderByMap[keyof ListTagsRequest.OrderByMap],
    orderDirection: ListTagsRequest.OrderDirectionMap[keyof ListTagsRequest.OrderDirectionMap],
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

export class ListTagsResponse extends jspb.Message {
  clearTagsList(): void;
  getTagsList(): Array<proto_feature_feature_pb.Tag>;
  setTagsList(value: Array<proto_feature_feature_pb.Tag>): void;
  addTags(value?: proto_feature_feature_pb.Tag, index?: number): proto_feature_feature_pb.Tag;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTagsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListTagsResponse): ListTagsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListTagsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTagsResponse;
  static deserializeBinaryFromReader(message: ListTagsResponse, reader: jspb.BinaryReader): ListTagsResponse;
}

export namespace ListTagsResponse {
  export type AsObject = {
    tagsList: Array<proto_feature_feature_pb.Tag.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class CreateFlagTriggerRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasCreateFlagTriggerCommand(): boolean;
  clearCreateFlagTriggerCommand(): void;
  getCreateFlagTriggerCommand(): proto_feature_command_pb.CreateFlagTriggerCommand | undefined;
  setCreateFlagTriggerCommand(value?: proto_feature_command_pb.CreateFlagTriggerCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateFlagTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateFlagTriggerRequest): CreateFlagTriggerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateFlagTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateFlagTriggerRequest;
  static deserializeBinaryFromReader(message: CreateFlagTriggerRequest, reader: jspb.BinaryReader): CreateFlagTriggerRequest;
}

export namespace CreateFlagTriggerRequest {
  export type AsObject = {
    environmentNamespace: string,
    createFlagTriggerCommand?: proto_feature_command_pb.CreateFlagTriggerCommand.AsObject,
  }
}

export class CreateFlagTriggerResponse extends jspb.Message {
  hasFlagTrigger(): boolean;
  clearFlagTrigger(): void;
  getFlagTrigger(): proto_feature_flag_trigger_pb.FlagTrigger | undefined;
  setFlagTrigger(value?: proto_feature_flag_trigger_pb.FlagTrigger): void;

  getUrl(): string;
  setUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateFlagTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateFlagTriggerResponse): CreateFlagTriggerResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateFlagTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateFlagTriggerResponse;
  static deserializeBinaryFromReader(message: CreateFlagTriggerResponse, reader: jspb.BinaryReader): CreateFlagTriggerResponse;
}

export namespace CreateFlagTriggerResponse {
  export type AsObject = {
    flagTrigger?: proto_feature_flag_trigger_pb.FlagTrigger.AsObject,
    url: string,
  }
}

export class DeleteFlagTriggerRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasDeleteFlagTriggerCommand(): boolean;
  clearDeleteFlagTriggerCommand(): void;
  getDeleteFlagTriggerCommand(): proto_feature_command_pb.DeleteFlagTriggerCommand | undefined;
  setDeleteFlagTriggerCommand(value?: proto_feature_command_pb.DeleteFlagTriggerCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteFlagTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteFlagTriggerRequest): DeleteFlagTriggerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteFlagTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteFlagTriggerRequest;
  static deserializeBinaryFromReader(message: DeleteFlagTriggerRequest, reader: jspb.BinaryReader): DeleteFlagTriggerRequest;
}

export namespace DeleteFlagTriggerRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    deleteFlagTriggerCommand?: proto_feature_command_pb.DeleteFlagTriggerCommand.AsObject,
  }
}

export class DeleteFlagTriggerResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteFlagTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteFlagTriggerResponse): DeleteFlagTriggerResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteFlagTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteFlagTriggerResponse;
  static deserializeBinaryFromReader(message: DeleteFlagTriggerResponse, reader: jspb.BinaryReader): DeleteFlagTriggerResponse;
}

export namespace DeleteFlagTriggerResponse {
  export type AsObject = {
  }
}

export class UpdateFlagTriggerRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasChangeFlagTriggerDescriptionCommand(): boolean;
  clearChangeFlagTriggerDescriptionCommand(): void;
  getChangeFlagTriggerDescriptionCommand(): proto_feature_command_pb.ChangeFlagTriggerDescriptionCommand | undefined;
  setChangeFlagTriggerDescriptionCommand(value?: proto_feature_command_pb.ChangeFlagTriggerDescriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateFlagTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateFlagTriggerRequest): UpdateFlagTriggerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateFlagTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateFlagTriggerRequest;
  static deserializeBinaryFromReader(message: UpdateFlagTriggerRequest, reader: jspb.BinaryReader): UpdateFlagTriggerRequest;
}

export namespace UpdateFlagTriggerRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    changeFlagTriggerDescriptionCommand?: proto_feature_command_pb.ChangeFlagTriggerDescriptionCommand.AsObject,
  }
}

export class UpdateFlagTriggerResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateFlagTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateFlagTriggerResponse): UpdateFlagTriggerResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateFlagTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateFlagTriggerResponse;
  static deserializeBinaryFromReader(message: UpdateFlagTriggerResponse, reader: jspb.BinaryReader): UpdateFlagTriggerResponse;
}

export namespace UpdateFlagTriggerResponse {
  export type AsObject = {
  }
}

export class EnableFlagTriggerRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasEnableFlagTriggerCommand(): boolean;
  clearEnableFlagTriggerCommand(): void;
  getEnableFlagTriggerCommand(): proto_feature_command_pb.EnableFlagTriggerCommand | undefined;
  setEnableFlagTriggerCommand(value?: proto_feature_command_pb.EnableFlagTriggerCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableFlagTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableFlagTriggerRequest): EnableFlagTriggerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableFlagTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableFlagTriggerRequest;
  static deserializeBinaryFromReader(message: EnableFlagTriggerRequest, reader: jspb.BinaryReader): EnableFlagTriggerRequest;
}

export namespace EnableFlagTriggerRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    enableFlagTriggerCommand?: proto_feature_command_pb.EnableFlagTriggerCommand.AsObject,
  }
}

export class EnableFlagTriggerResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableFlagTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableFlagTriggerResponse): EnableFlagTriggerResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableFlagTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableFlagTriggerResponse;
  static deserializeBinaryFromReader(message: EnableFlagTriggerResponse, reader: jspb.BinaryReader): EnableFlagTriggerResponse;
}

export namespace EnableFlagTriggerResponse {
  export type AsObject = {
  }
}

export class DisableFlagTriggerRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasDisableFlagTriggerCommand(): boolean;
  clearDisableFlagTriggerCommand(): void;
  getDisableFlagTriggerCommand(): proto_feature_command_pb.DisableFlagTriggerCommand | undefined;
  setDisableFlagTriggerCommand(value?: proto_feature_command_pb.DisableFlagTriggerCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableFlagTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableFlagTriggerRequest): DisableFlagTriggerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableFlagTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableFlagTriggerRequest;
  static deserializeBinaryFromReader(message: DisableFlagTriggerRequest, reader: jspb.BinaryReader): DisableFlagTriggerRequest;
}

export namespace DisableFlagTriggerRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    disableFlagTriggerCommand?: proto_feature_command_pb.DisableFlagTriggerCommand.AsObject,
  }
}

export class DisableFlagTriggerResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableFlagTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableFlagTriggerResponse): DisableFlagTriggerResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableFlagTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableFlagTriggerResponse;
  static deserializeBinaryFromReader(message: DisableFlagTriggerResponse, reader: jspb.BinaryReader): DisableFlagTriggerResponse;
}

export namespace DisableFlagTriggerResponse {
  export type AsObject = {
  }
}

export class ResetFlagTriggerRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasResetFlagTriggerCommand(): boolean;
  clearResetFlagTriggerCommand(): void;
  getResetFlagTriggerCommand(): proto_feature_command_pb.ResetFlagTriggerCommand | undefined;
  setResetFlagTriggerCommand(value?: proto_feature_command_pb.ResetFlagTriggerCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetFlagTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetFlagTriggerRequest): ResetFlagTriggerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResetFlagTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetFlagTriggerRequest;
  static deserializeBinaryFromReader(message: ResetFlagTriggerRequest, reader: jspb.BinaryReader): ResetFlagTriggerRequest;
}

export namespace ResetFlagTriggerRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    resetFlagTriggerCommand?: proto_feature_command_pb.ResetFlagTriggerCommand.AsObject,
  }
}

export class ResetFlagTriggerResponse extends jspb.Message {
  hasFlagTrigger(): boolean;
  clearFlagTrigger(): void;
  getFlagTrigger(): proto_feature_flag_trigger_pb.FlagTrigger | undefined;
  setFlagTrigger(value?: proto_feature_flag_trigger_pb.FlagTrigger): void;

  getUrl(): string;
  setUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetFlagTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetFlagTriggerResponse): ResetFlagTriggerResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResetFlagTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetFlagTriggerResponse;
  static deserializeBinaryFromReader(message: ResetFlagTriggerResponse, reader: jspb.BinaryReader): ResetFlagTriggerResponse;
}

export namespace ResetFlagTriggerResponse {
  export type AsObject = {
    flagTrigger?: proto_feature_flag_trigger_pb.FlagTrigger.AsObject,
    url: string,
  }
}

export class GetFlagTriggerRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFlagTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFlagTriggerRequest): GetFlagTriggerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFlagTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFlagTriggerRequest;
  static deserializeBinaryFromReader(message: GetFlagTriggerRequest, reader: jspb.BinaryReader): GetFlagTriggerRequest;
}

export namespace GetFlagTriggerRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
  }
}

export class GetFlagTriggerResponse extends jspb.Message {
  hasFlagTrigger(): boolean;
  clearFlagTrigger(): void;
  getFlagTrigger(): proto_feature_flag_trigger_pb.FlagTrigger | undefined;
  setFlagTrigger(value?: proto_feature_flag_trigger_pb.FlagTrigger): void;

  getUrl(): string;
  setUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFlagTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFlagTriggerResponse): GetFlagTriggerResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFlagTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFlagTriggerResponse;
  static deserializeBinaryFromReader(message: GetFlagTriggerResponse, reader: jspb.BinaryReader): GetFlagTriggerResponse;
}

export namespace GetFlagTriggerResponse {
  export type AsObject = {
    flagTrigger?: proto_feature_flag_trigger_pb.FlagTrigger.AsObject,
    url: string,
  }
}

export class ListFlagTriggersRequest extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getCursor(): string;
  setCursor(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getOrderBy(): ListFlagTriggersRequest.OrderByMap[keyof ListFlagTriggersRequest.OrderByMap];
  setOrderBy(value: ListFlagTriggersRequest.OrderByMap[keyof ListFlagTriggersRequest.OrderByMap]): void;

  getOrderDirection(): ListFlagTriggersRequest.OrderDirectionMap[keyof ListFlagTriggersRequest.OrderDirectionMap];
  setOrderDirection(value: ListFlagTriggersRequest.OrderDirectionMap[keyof ListFlagTriggersRequest.OrderDirectionMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFlagTriggersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFlagTriggersRequest): ListFlagTriggersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFlagTriggersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFlagTriggersRequest;
  static deserializeBinaryFromReader(message: ListFlagTriggersRequest, reader: jspb.BinaryReader): ListFlagTriggersRequest;
}

export namespace ListFlagTriggersRequest {
  export type AsObject = {
    featureId: string,
    environmentNamespace: string,
    cursor: string,
    pageSize: number,
    orderBy: ListFlagTriggersRequest.OrderByMap[keyof ListFlagTriggersRequest.OrderByMap],
    orderDirection: ListFlagTriggersRequest.OrderDirectionMap[keyof ListFlagTriggersRequest.OrderDirectionMap],
  }

  export interface OrderByMap {
    DEFAULT: 0;
    CREATED_AT: 1;
    UPDATED_AT: 2;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListFlagTriggersResponse extends jspb.Message {
  clearFlagTriggersList(): void;
  getFlagTriggersList(): Array<ListFlagTriggersResponse.FlagTriggerWithUrl>;
  setFlagTriggersList(value: Array<ListFlagTriggersResponse.FlagTriggerWithUrl>): void;
  addFlagTriggers(value?: ListFlagTriggersResponse.FlagTriggerWithUrl, index?: number): ListFlagTriggersResponse.FlagTriggerWithUrl;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFlagTriggersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFlagTriggersResponse): ListFlagTriggersResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFlagTriggersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFlagTriggersResponse;
  static deserializeBinaryFromReader(message: ListFlagTriggersResponse, reader: jspb.BinaryReader): ListFlagTriggersResponse;
}

export namespace ListFlagTriggersResponse {
  export type AsObject = {
    flagTriggersList: Array<ListFlagTriggersResponse.FlagTriggerWithUrl.AsObject>,
    cursor: string,
    totalCount: number,
  }

  export class FlagTriggerWithUrl extends jspb.Message {
    hasFlagTrigger(): boolean;
    clearFlagTrigger(): void;
    getFlagTrigger(): proto_feature_flag_trigger_pb.FlagTrigger | undefined;
    setFlagTrigger(value?: proto_feature_flag_trigger_pb.FlagTrigger): void;

    getUrl(): string;
    setUrl(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): FlagTriggerWithUrl.AsObject;
    static toObject(includeInstance: boolean, msg: FlagTriggerWithUrl): FlagTriggerWithUrl.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: FlagTriggerWithUrl, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): FlagTriggerWithUrl;
    static deserializeBinaryFromReader(message: FlagTriggerWithUrl, reader: jspb.BinaryReader): FlagTriggerWithUrl;
  }

  export namespace FlagTriggerWithUrl {
    export type AsObject = {
      flagTrigger?: proto_feature_flag_trigger_pb.FlagTrigger.AsObject,
      url: string,
    }
  }
}

export class FlagTriggerWebhookRequest extends jspb.Message {
  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagTriggerWebhookRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FlagTriggerWebhookRequest): FlagTriggerWebhookRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FlagTriggerWebhookRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FlagTriggerWebhookRequest;
  static deserializeBinaryFromReader(message: FlagTriggerWebhookRequest, reader: jspb.BinaryReader): FlagTriggerWebhookRequest;
}

export namespace FlagTriggerWebhookRequest {
  export type AsObject = {
    token: string,
  }
}

export class FlagTriggerWebhookResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagTriggerWebhookResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FlagTriggerWebhookResponse): FlagTriggerWebhookResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FlagTriggerWebhookResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FlagTriggerWebhookResponse;
  static deserializeBinaryFromReader(message: FlagTriggerWebhookResponse, reader: jspb.BinaryReader): FlagTriggerWebhookResponse;
}

export namespace FlagTriggerWebhookResponse {
  export type AsObject = {
  }
}

