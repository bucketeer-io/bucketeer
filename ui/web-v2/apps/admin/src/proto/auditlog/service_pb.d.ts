// package: bucketeer.auditlog
// file: proto/auditlog/service.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as proto_auditlog_auditlog_pb from "../../proto/auditlog/auditlog_pb";

export class ListAuditLogsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getOrderBy(): ListAuditLogsRequest.OrderByMap[keyof ListAuditLogsRequest.OrderByMap];
  setOrderBy(value: ListAuditLogsRequest.OrderByMap[keyof ListAuditLogsRequest.OrderByMap]): void;

  getOrderDirection(): ListAuditLogsRequest.OrderDirectionMap[keyof ListAuditLogsRequest.OrderDirectionMap];
  setOrderDirection(value: ListAuditLogsRequest.OrderDirectionMap[keyof ListAuditLogsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  getFrom(): number;
  setFrom(value: number): void;

  getTo(): number;
  setTo(value: number): void;

  hasEntityType(): boolean;
  clearEntityType(): void;
  getEntityType(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setEntityType(value?: google_protobuf_wrappers_pb.Int32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAuditLogsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAuditLogsRequest): ListAuditLogsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAuditLogsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAuditLogsRequest;
  static deserializeBinaryFromReader(message: ListAuditLogsRequest, reader: jspb.BinaryReader): ListAuditLogsRequest;
}

export namespace ListAuditLogsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    environmentNamespace: string,
    orderBy: ListAuditLogsRequest.OrderByMap[keyof ListAuditLogsRequest.OrderByMap],
    orderDirection: ListAuditLogsRequest.OrderDirectionMap[keyof ListAuditLogsRequest.OrderDirectionMap],
    searchKeyword: string,
    from: number,
    to: number,
    entityType?: google_protobuf_wrappers_pb.Int32Value.AsObject,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    TIMESTAMP: 1;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    DESC: 0;
    ASC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListAuditLogsResponse extends jspb.Message {
  clearAuditLogsList(): void;
  getAuditLogsList(): Array<proto_auditlog_auditlog_pb.AuditLog>;
  setAuditLogsList(value: Array<proto_auditlog_auditlog_pb.AuditLog>): void;
  addAuditLogs(value?: proto_auditlog_auditlog_pb.AuditLog, index?: number): proto_auditlog_auditlog_pb.AuditLog;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAuditLogsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAuditLogsResponse): ListAuditLogsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAuditLogsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAuditLogsResponse;
  static deserializeBinaryFromReader(message: ListAuditLogsResponse, reader: jspb.BinaryReader): ListAuditLogsResponse;
}

export namespace ListAuditLogsResponse {
  export type AsObject = {
    auditLogsList: Array<proto_auditlog_auditlog_pb.AuditLog.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class ListAdminAuditLogsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListAdminAuditLogsRequest.OrderByMap[keyof ListAdminAuditLogsRequest.OrderByMap];
  setOrderBy(value: ListAdminAuditLogsRequest.OrderByMap[keyof ListAdminAuditLogsRequest.OrderByMap]): void;

  getOrderDirection(): ListAdminAuditLogsRequest.OrderDirectionMap[keyof ListAdminAuditLogsRequest.OrderDirectionMap];
  setOrderDirection(value: ListAdminAuditLogsRequest.OrderDirectionMap[keyof ListAdminAuditLogsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  getFrom(): number;
  setFrom(value: number): void;

  getTo(): number;
  setTo(value: number): void;

  hasEntityType(): boolean;
  clearEntityType(): void;
  getEntityType(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setEntityType(value?: google_protobuf_wrappers_pb.Int32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAdminAuditLogsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAdminAuditLogsRequest): ListAdminAuditLogsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAdminAuditLogsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAdminAuditLogsRequest;
  static deserializeBinaryFromReader(message: ListAdminAuditLogsRequest, reader: jspb.BinaryReader): ListAdminAuditLogsRequest;
}

export namespace ListAdminAuditLogsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    orderBy: ListAdminAuditLogsRequest.OrderByMap[keyof ListAdminAuditLogsRequest.OrderByMap],
    orderDirection: ListAdminAuditLogsRequest.OrderDirectionMap[keyof ListAdminAuditLogsRequest.OrderDirectionMap],
    searchKeyword: string,
    from: number,
    to: number,
    entityType?: google_protobuf_wrappers_pb.Int32Value.AsObject,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    TIMESTAMP: 1;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    DESC: 0;
    ASC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListAdminAuditLogsResponse extends jspb.Message {
  clearAuditLogsList(): void;
  getAuditLogsList(): Array<proto_auditlog_auditlog_pb.AuditLog>;
  setAuditLogsList(value: Array<proto_auditlog_auditlog_pb.AuditLog>): void;
  addAuditLogs(value?: proto_auditlog_auditlog_pb.AuditLog, index?: number): proto_auditlog_auditlog_pb.AuditLog;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAdminAuditLogsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAdminAuditLogsResponse): ListAdminAuditLogsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAdminAuditLogsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAdminAuditLogsResponse;
  static deserializeBinaryFromReader(message: ListAdminAuditLogsResponse, reader: jspb.BinaryReader): ListAdminAuditLogsResponse;
}

export namespace ListAdminAuditLogsResponse {
  export type AsObject = {
    auditLogsList: Array<proto_auditlog_auditlog_pb.AuditLog.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class ListFeatureHistoryRequest extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getOrderBy(): ListFeatureHistoryRequest.OrderByMap[keyof ListFeatureHistoryRequest.OrderByMap];
  setOrderBy(value: ListFeatureHistoryRequest.OrderByMap[keyof ListFeatureHistoryRequest.OrderByMap]): void;

  getOrderDirection(): ListFeatureHistoryRequest.OrderDirectionMap[keyof ListFeatureHistoryRequest.OrderDirectionMap];
  setOrderDirection(value: ListFeatureHistoryRequest.OrderDirectionMap[keyof ListFeatureHistoryRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  getFrom(): number;
  setFrom(value: number): void;

  getTo(): number;
  setTo(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFeatureHistoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFeatureHistoryRequest): ListFeatureHistoryRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFeatureHistoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFeatureHistoryRequest;
  static deserializeBinaryFromReader(message: ListFeatureHistoryRequest, reader: jspb.BinaryReader): ListFeatureHistoryRequest;
}

export namespace ListFeatureHistoryRequest {
  export type AsObject = {
    featureId: string,
    pageSize: number,
    cursor: string,
    environmentNamespace: string,
    orderBy: ListFeatureHistoryRequest.OrderByMap[keyof ListFeatureHistoryRequest.OrderByMap],
    orderDirection: ListFeatureHistoryRequest.OrderDirectionMap[keyof ListFeatureHistoryRequest.OrderDirectionMap],
    searchKeyword: string,
    from: number,
    to: number,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    TIMESTAMP: 1;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    DESC: 0;
    ASC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListFeatureHistoryResponse extends jspb.Message {
  clearAuditLogsList(): void;
  getAuditLogsList(): Array<proto_auditlog_auditlog_pb.AuditLog>;
  setAuditLogsList(value: Array<proto_auditlog_auditlog_pb.AuditLog>): void;
  addAuditLogs(value?: proto_auditlog_auditlog_pb.AuditLog, index?: number): proto_auditlog_auditlog_pb.AuditLog;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFeatureHistoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFeatureHistoryResponse): ListFeatureHistoryResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFeatureHistoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFeatureHistoryResponse;
  static deserializeBinaryFromReader(message: ListFeatureHistoryResponse, reader: jspb.BinaryReader): ListFeatureHistoryResponse;
}

export namespace ListFeatureHistoryResponse {
  export type AsObject = {
    auditLogsList: Array<proto_auditlog_auditlog_pb.AuditLog.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

