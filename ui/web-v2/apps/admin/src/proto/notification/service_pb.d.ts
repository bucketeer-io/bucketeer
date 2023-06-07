// package: bucketeer.notification
// file: proto/notification/service.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as proto_notification_subscription_pb from "../../proto/notification/subscription_pb";
import * as proto_notification_command_pb from "../../proto/notification/command_pb";

export class GetAdminSubscriptionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAdminSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAdminSubscriptionRequest): GetAdminSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAdminSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAdminSubscriptionRequest;
  static deserializeBinaryFromReader(message: GetAdminSubscriptionRequest, reader: jspb.BinaryReader): GetAdminSubscriptionRequest;
}

export namespace GetAdminSubscriptionRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetAdminSubscriptionResponse extends jspb.Message {
  hasSubscription(): boolean;
  clearSubscription(): void;
  getSubscription(): proto_notification_subscription_pb.Subscription | undefined;
  setSubscription(value?: proto_notification_subscription_pb.Subscription): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAdminSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAdminSubscriptionResponse): GetAdminSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAdminSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAdminSubscriptionResponse;
  static deserializeBinaryFromReader(message: GetAdminSubscriptionResponse, reader: jspb.BinaryReader): GetAdminSubscriptionResponse;
}

export namespace GetAdminSubscriptionResponse {
  export type AsObject = {
    subscription?: proto_notification_subscription_pb.Subscription.AsObject,
  }
}

export class ListAdminSubscriptionsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearSourceTypesList(): void;
  getSourceTypesList(): Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap], index?: number): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  getOrderBy(): ListAdminSubscriptionsRequest.OrderByMap[keyof ListAdminSubscriptionsRequest.OrderByMap];
  setOrderBy(value: ListAdminSubscriptionsRequest.OrderByMap[keyof ListAdminSubscriptionsRequest.OrderByMap]): void;

  getOrderDirection(): ListAdminSubscriptionsRequest.OrderDirectionMap[keyof ListAdminSubscriptionsRequest.OrderDirectionMap];
  setOrderDirection(value: ListAdminSubscriptionsRequest.OrderDirectionMap[keyof ListAdminSubscriptionsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAdminSubscriptionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAdminSubscriptionsRequest): ListAdminSubscriptionsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAdminSubscriptionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAdminSubscriptionsRequest;
  static deserializeBinaryFromReader(message: ListAdminSubscriptionsRequest, reader: jspb.BinaryReader): ListAdminSubscriptionsRequest;
}

export namespace ListAdminSubscriptionsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    sourceTypesList: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>,
    orderBy: ListAdminSubscriptionsRequest.OrderByMap[keyof ListAdminSubscriptionsRequest.OrderByMap],
    orderDirection: ListAdminSubscriptionsRequest.OrderDirectionMap[keyof ListAdminSubscriptionsRequest.OrderDirectionMap],
    searchKeyword: string,
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
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

export class ListAdminSubscriptionsResponse extends jspb.Message {
  clearSubscriptionsList(): void;
  getSubscriptionsList(): Array<proto_notification_subscription_pb.Subscription>;
  setSubscriptionsList(value: Array<proto_notification_subscription_pb.Subscription>): void;
  addSubscriptions(value?: proto_notification_subscription_pb.Subscription, index?: number): proto_notification_subscription_pb.Subscription;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAdminSubscriptionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAdminSubscriptionsResponse): ListAdminSubscriptionsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAdminSubscriptionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAdminSubscriptionsResponse;
  static deserializeBinaryFromReader(message: ListAdminSubscriptionsResponse, reader: jspb.BinaryReader): ListAdminSubscriptionsResponse;
}

export namespace ListAdminSubscriptionsResponse {
  export type AsObject = {
    subscriptionsList: Array<proto_notification_subscription_pb.Subscription.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class ListEnabledAdminSubscriptionsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearSourceTypesList(): void;
  getSourceTypesList(): Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap], index?: number): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnabledAdminSubscriptionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnabledAdminSubscriptionsRequest): ListEnabledAdminSubscriptionsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnabledAdminSubscriptionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnabledAdminSubscriptionsRequest;
  static deserializeBinaryFromReader(message: ListEnabledAdminSubscriptionsRequest, reader: jspb.BinaryReader): ListEnabledAdminSubscriptionsRequest;
}

export namespace ListEnabledAdminSubscriptionsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    sourceTypesList: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>,
  }
}

export class ListEnabledAdminSubscriptionsResponse extends jspb.Message {
  clearSubscriptionsList(): void;
  getSubscriptionsList(): Array<proto_notification_subscription_pb.Subscription>;
  setSubscriptionsList(value: Array<proto_notification_subscription_pb.Subscription>): void;
  addSubscriptions(value?: proto_notification_subscription_pb.Subscription, index?: number): proto_notification_subscription_pb.Subscription;

  getCursor(): string;
  setCursor(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnabledAdminSubscriptionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnabledAdminSubscriptionsResponse): ListEnabledAdminSubscriptionsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnabledAdminSubscriptionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnabledAdminSubscriptionsResponse;
  static deserializeBinaryFromReader(message: ListEnabledAdminSubscriptionsResponse, reader: jspb.BinaryReader): ListEnabledAdminSubscriptionsResponse;
}

export namespace ListEnabledAdminSubscriptionsResponse {
  export type AsObject = {
    subscriptionsList: Array<proto_notification_subscription_pb.Subscription.AsObject>,
    cursor: string,
  }
}

export class CreateAdminSubscriptionRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_notification_command_pb.CreateAdminSubscriptionCommand | undefined;
  setCommand(value?: proto_notification_command_pb.CreateAdminSubscriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAdminSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAdminSubscriptionRequest): CreateAdminSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAdminSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAdminSubscriptionRequest;
  static deserializeBinaryFromReader(message: CreateAdminSubscriptionRequest, reader: jspb.BinaryReader): CreateAdminSubscriptionRequest;
}

export namespace CreateAdminSubscriptionRequest {
  export type AsObject = {
    command?: proto_notification_command_pb.CreateAdminSubscriptionCommand.AsObject,
  }
}

export class CreateAdminSubscriptionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAdminSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAdminSubscriptionResponse): CreateAdminSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAdminSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAdminSubscriptionResponse;
  static deserializeBinaryFromReader(message: CreateAdminSubscriptionResponse, reader: jspb.BinaryReader): CreateAdminSubscriptionResponse;
}

export namespace CreateAdminSubscriptionResponse {
  export type AsObject = {
  }
}

export class DeleteAdminSubscriptionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_notification_command_pb.DeleteAdminSubscriptionCommand | undefined;
  setCommand(value?: proto_notification_command_pb.DeleteAdminSubscriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAdminSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAdminSubscriptionRequest): DeleteAdminSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAdminSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAdminSubscriptionRequest;
  static deserializeBinaryFromReader(message: DeleteAdminSubscriptionRequest, reader: jspb.BinaryReader): DeleteAdminSubscriptionRequest;
}

export namespace DeleteAdminSubscriptionRequest {
  export type AsObject = {
    id: string,
    command?: proto_notification_command_pb.DeleteAdminSubscriptionCommand.AsObject,
  }
}

export class DeleteAdminSubscriptionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAdminSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAdminSubscriptionResponse): DeleteAdminSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAdminSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAdminSubscriptionResponse;
  static deserializeBinaryFromReader(message: DeleteAdminSubscriptionResponse, reader: jspb.BinaryReader): DeleteAdminSubscriptionResponse;
}

export namespace DeleteAdminSubscriptionResponse {
  export type AsObject = {
  }
}

export class EnableAdminSubscriptionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_notification_command_pb.EnableAdminSubscriptionCommand | undefined;
  setCommand(value?: proto_notification_command_pb.EnableAdminSubscriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAdminSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAdminSubscriptionRequest): EnableAdminSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAdminSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAdminSubscriptionRequest;
  static deserializeBinaryFromReader(message: EnableAdminSubscriptionRequest, reader: jspb.BinaryReader): EnableAdminSubscriptionRequest;
}

export namespace EnableAdminSubscriptionRequest {
  export type AsObject = {
    id: string,
    command?: proto_notification_command_pb.EnableAdminSubscriptionCommand.AsObject,
  }
}

export class EnableAdminSubscriptionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableAdminSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableAdminSubscriptionResponse): EnableAdminSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableAdminSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableAdminSubscriptionResponse;
  static deserializeBinaryFromReader(message: EnableAdminSubscriptionResponse, reader: jspb.BinaryReader): EnableAdminSubscriptionResponse;
}

export namespace EnableAdminSubscriptionResponse {
  export type AsObject = {
  }
}

export class DisableAdminSubscriptionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_notification_command_pb.DisableAdminSubscriptionCommand | undefined;
  setCommand(value?: proto_notification_command_pb.DisableAdminSubscriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAdminSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAdminSubscriptionRequest): DisableAdminSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAdminSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAdminSubscriptionRequest;
  static deserializeBinaryFromReader(message: DisableAdminSubscriptionRequest, reader: jspb.BinaryReader): DisableAdminSubscriptionRequest;
}

export namespace DisableAdminSubscriptionRequest {
  export type AsObject = {
    id: string,
    command?: proto_notification_command_pb.DisableAdminSubscriptionCommand.AsObject,
  }
}

export class DisableAdminSubscriptionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAdminSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAdminSubscriptionResponse): DisableAdminSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAdminSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAdminSubscriptionResponse;
  static deserializeBinaryFromReader(message: DisableAdminSubscriptionResponse, reader: jspb.BinaryReader): DisableAdminSubscriptionResponse;
}

export namespace DisableAdminSubscriptionResponse {
  export type AsObject = {
  }
}

export class UpdateAdminSubscriptionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasAddSourceTypesCommand(): boolean;
  clearAddSourceTypesCommand(): void;
  getAddSourceTypesCommand(): proto_notification_command_pb.AddAdminSubscriptionSourceTypesCommand | undefined;
  setAddSourceTypesCommand(value?: proto_notification_command_pb.AddAdminSubscriptionSourceTypesCommand): void;

  hasDeleteSourceTypesCommand(): boolean;
  clearDeleteSourceTypesCommand(): void;
  getDeleteSourceTypesCommand(): proto_notification_command_pb.DeleteAdminSubscriptionSourceTypesCommand | undefined;
  setDeleteSourceTypesCommand(value?: proto_notification_command_pb.DeleteAdminSubscriptionSourceTypesCommand): void;

  hasRenameSubscriptionCommand(): boolean;
  clearRenameSubscriptionCommand(): void;
  getRenameSubscriptionCommand(): proto_notification_command_pb.RenameAdminSubscriptionCommand | undefined;
  setRenameSubscriptionCommand(value?: proto_notification_command_pb.RenameAdminSubscriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAdminSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAdminSubscriptionRequest): UpdateAdminSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateAdminSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAdminSubscriptionRequest;
  static deserializeBinaryFromReader(message: UpdateAdminSubscriptionRequest, reader: jspb.BinaryReader): UpdateAdminSubscriptionRequest;
}

export namespace UpdateAdminSubscriptionRequest {
  export type AsObject = {
    id: string,
    addSourceTypesCommand?: proto_notification_command_pb.AddAdminSubscriptionSourceTypesCommand.AsObject,
    deleteSourceTypesCommand?: proto_notification_command_pb.DeleteAdminSubscriptionSourceTypesCommand.AsObject,
    renameSubscriptionCommand?: proto_notification_command_pb.RenameAdminSubscriptionCommand.AsObject,
  }
}

export class UpdateAdminSubscriptionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAdminSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAdminSubscriptionResponse): UpdateAdminSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateAdminSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAdminSubscriptionResponse;
  static deserializeBinaryFromReader(message: UpdateAdminSubscriptionResponse, reader: jspb.BinaryReader): UpdateAdminSubscriptionResponse;
}

export namespace UpdateAdminSubscriptionResponse {
  export type AsObject = {
  }
}

export class GetSubscriptionRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSubscriptionRequest): GetSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSubscriptionRequest;
  static deserializeBinaryFromReader(message: GetSubscriptionRequest, reader: jspb.BinaryReader): GetSubscriptionRequest;
}

export namespace GetSubscriptionRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
  }
}

export class GetSubscriptionResponse extends jspb.Message {
  hasSubscription(): boolean;
  clearSubscription(): void;
  getSubscription(): proto_notification_subscription_pb.Subscription | undefined;
  setSubscription(value?: proto_notification_subscription_pb.Subscription): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSubscriptionResponse): GetSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSubscriptionResponse;
  static deserializeBinaryFromReader(message: GetSubscriptionResponse, reader: jspb.BinaryReader): GetSubscriptionResponse;
}

export namespace GetSubscriptionResponse {
  export type AsObject = {
    subscription?: proto_notification_subscription_pb.Subscription.AsObject,
  }
}

export class ListSubscriptionsRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearSourceTypesList(): void;
  getSourceTypesList(): Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap], index?: number): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  getOrderBy(): ListSubscriptionsRequest.OrderByMap[keyof ListSubscriptionsRequest.OrderByMap];
  setOrderBy(value: ListSubscriptionsRequest.OrderByMap[keyof ListSubscriptionsRequest.OrderByMap]): void;

  getOrderDirection(): ListSubscriptionsRequest.OrderDirectionMap[keyof ListSubscriptionsRequest.OrderDirectionMap];
  setOrderDirection(value: ListSubscriptionsRequest.OrderDirectionMap[keyof ListSubscriptionsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSubscriptionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListSubscriptionsRequest): ListSubscriptionsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListSubscriptionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSubscriptionsRequest;
  static deserializeBinaryFromReader(message: ListSubscriptionsRequest, reader: jspb.BinaryReader): ListSubscriptionsRequest;
}

export namespace ListSubscriptionsRequest {
  export type AsObject = {
    environmentNamespace: string,
    pageSize: number,
    cursor: string,
    sourceTypesList: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>,
    orderBy: ListSubscriptionsRequest.OrderByMap[keyof ListSubscriptionsRequest.OrderByMap],
    orderDirection: ListSubscriptionsRequest.OrderDirectionMap[keyof ListSubscriptionsRequest.OrderDirectionMap],
    searchKeyword: string,
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
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

export class ListSubscriptionsResponse extends jspb.Message {
  clearSubscriptionsList(): void;
  getSubscriptionsList(): Array<proto_notification_subscription_pb.Subscription>;
  setSubscriptionsList(value: Array<proto_notification_subscription_pb.Subscription>): void;
  addSubscriptions(value?: proto_notification_subscription_pb.Subscription, index?: number): proto_notification_subscription_pb.Subscription;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSubscriptionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListSubscriptionsResponse): ListSubscriptionsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListSubscriptionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSubscriptionsResponse;
  static deserializeBinaryFromReader(message: ListSubscriptionsResponse, reader: jspb.BinaryReader): ListSubscriptionsResponse;
}

export namespace ListSubscriptionsResponse {
  export type AsObject = {
    subscriptionsList: Array<proto_notification_subscription_pb.Subscription.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class ListEnabledSubscriptionsRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearSourceTypesList(): void;
  getSourceTypesList(): Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap], index?: number): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnabledSubscriptionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnabledSubscriptionsRequest): ListEnabledSubscriptionsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnabledSubscriptionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnabledSubscriptionsRequest;
  static deserializeBinaryFromReader(message: ListEnabledSubscriptionsRequest, reader: jspb.BinaryReader): ListEnabledSubscriptionsRequest;
}

export namespace ListEnabledSubscriptionsRequest {
  export type AsObject = {
    environmentNamespace: string,
    pageSize: number,
    cursor: string,
    sourceTypesList: Array<proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]>,
  }
}

export class ListEnabledSubscriptionsResponse extends jspb.Message {
  clearSubscriptionsList(): void;
  getSubscriptionsList(): Array<proto_notification_subscription_pb.Subscription>;
  setSubscriptionsList(value: Array<proto_notification_subscription_pb.Subscription>): void;
  addSubscriptions(value?: proto_notification_subscription_pb.Subscription, index?: number): proto_notification_subscription_pb.Subscription;

  getCursor(): string;
  setCursor(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnabledSubscriptionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnabledSubscriptionsResponse): ListEnabledSubscriptionsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnabledSubscriptionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnabledSubscriptionsResponse;
  static deserializeBinaryFromReader(message: ListEnabledSubscriptionsResponse, reader: jspb.BinaryReader): ListEnabledSubscriptionsResponse;
}

export namespace ListEnabledSubscriptionsResponse {
  export type AsObject = {
    subscriptionsList: Array<proto_notification_subscription_pb.Subscription.AsObject>,
    cursor: string,
  }
}

export class CreateSubscriptionRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_notification_command_pb.CreateSubscriptionCommand | undefined;
  setCommand(value?: proto_notification_command_pb.CreateSubscriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateSubscriptionRequest): CreateSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateSubscriptionRequest;
  static deserializeBinaryFromReader(message: CreateSubscriptionRequest, reader: jspb.BinaryReader): CreateSubscriptionRequest;
}

export namespace CreateSubscriptionRequest {
  export type AsObject = {
    environmentNamespace: string,
    command?: proto_notification_command_pb.CreateSubscriptionCommand.AsObject,
  }
}

export class CreateSubscriptionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateSubscriptionResponse): CreateSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateSubscriptionResponse;
  static deserializeBinaryFromReader(message: CreateSubscriptionResponse, reader: jspb.BinaryReader): CreateSubscriptionResponse;
}

export namespace CreateSubscriptionResponse {
  export type AsObject = {
  }
}

export class DeleteSubscriptionRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_notification_command_pb.DeleteSubscriptionCommand | undefined;
  setCommand(value?: proto_notification_command_pb.DeleteSubscriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSubscriptionRequest): DeleteSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSubscriptionRequest;
  static deserializeBinaryFromReader(message: DeleteSubscriptionRequest, reader: jspb.BinaryReader): DeleteSubscriptionRequest;
}

export namespace DeleteSubscriptionRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    command?: proto_notification_command_pb.DeleteSubscriptionCommand.AsObject,
  }
}

export class DeleteSubscriptionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSubscriptionResponse): DeleteSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSubscriptionResponse;
  static deserializeBinaryFromReader(message: DeleteSubscriptionResponse, reader: jspb.BinaryReader): DeleteSubscriptionResponse;
}

export namespace DeleteSubscriptionResponse {
  export type AsObject = {
  }
}

export class EnableSubscriptionRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_notification_command_pb.EnableSubscriptionCommand | undefined;
  setCommand(value?: proto_notification_command_pb.EnableSubscriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableSubscriptionRequest): EnableSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableSubscriptionRequest;
  static deserializeBinaryFromReader(message: EnableSubscriptionRequest, reader: jspb.BinaryReader): EnableSubscriptionRequest;
}

export namespace EnableSubscriptionRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    command?: proto_notification_command_pb.EnableSubscriptionCommand.AsObject,
  }
}

export class EnableSubscriptionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableSubscriptionResponse): EnableSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableSubscriptionResponse;
  static deserializeBinaryFromReader(message: EnableSubscriptionResponse, reader: jspb.BinaryReader): EnableSubscriptionResponse;
}

export namespace EnableSubscriptionResponse {
  export type AsObject = {
  }
}

export class DisableSubscriptionRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_notification_command_pb.DisableSubscriptionCommand | undefined;
  setCommand(value?: proto_notification_command_pb.DisableSubscriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableSubscriptionRequest): DisableSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableSubscriptionRequest;
  static deserializeBinaryFromReader(message: DisableSubscriptionRequest, reader: jspb.BinaryReader): DisableSubscriptionRequest;
}

export namespace DisableSubscriptionRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    command?: proto_notification_command_pb.DisableSubscriptionCommand.AsObject,
  }
}

export class DisableSubscriptionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableSubscriptionResponse): DisableSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableSubscriptionResponse;
  static deserializeBinaryFromReader(message: DisableSubscriptionResponse, reader: jspb.BinaryReader): DisableSubscriptionResponse;
}

export namespace DisableSubscriptionResponse {
  export type AsObject = {
  }
}

export class UpdateSubscriptionRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasAddSourceTypesCommand(): boolean;
  clearAddSourceTypesCommand(): void;
  getAddSourceTypesCommand(): proto_notification_command_pb.AddSourceTypesCommand | undefined;
  setAddSourceTypesCommand(value?: proto_notification_command_pb.AddSourceTypesCommand): void;

  hasDeleteSourceTypesCommand(): boolean;
  clearDeleteSourceTypesCommand(): void;
  getDeleteSourceTypesCommand(): proto_notification_command_pb.DeleteSourceTypesCommand | undefined;
  setDeleteSourceTypesCommand(value?: proto_notification_command_pb.DeleteSourceTypesCommand): void;

  hasRenameSubscriptionCommand(): boolean;
  clearRenameSubscriptionCommand(): void;
  getRenameSubscriptionCommand(): proto_notification_command_pb.RenameSubscriptionCommand | undefined;
  setRenameSubscriptionCommand(value?: proto_notification_command_pb.RenameSubscriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSubscriptionRequest): UpdateSubscriptionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSubscriptionRequest;
  static deserializeBinaryFromReader(message: UpdateSubscriptionRequest, reader: jspb.BinaryReader): UpdateSubscriptionRequest;
}

export namespace UpdateSubscriptionRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    addSourceTypesCommand?: proto_notification_command_pb.AddSourceTypesCommand.AsObject,
    deleteSourceTypesCommand?: proto_notification_command_pb.DeleteSourceTypesCommand.AsObject,
    renameSubscriptionCommand?: proto_notification_command_pb.RenameSubscriptionCommand.AsObject,
  }
}

export class UpdateSubscriptionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSubscriptionResponse): UpdateSubscriptionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSubscriptionResponse;
  static deserializeBinaryFromReader(message: UpdateSubscriptionResponse, reader: jspb.BinaryReader): UpdateSubscriptionResponse;
}

export namespace UpdateSubscriptionResponse {
  export type AsObject = {
  }
}

