// package: bucketeer.migration
// file: proto/migration/mysql_service.proto

import * as proto_migration_mysql_service_pb from "../../proto/migration/mysql_service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type MigrationMySQLServiceMigrateAllMasterSchema = {
  readonly methodName: string;
  readonly service: typeof MigrationMySQLService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_migration_mysql_service_pb.MigrateAllMasterSchemaRequest;
  readonly responseType: typeof proto_migration_mysql_service_pb.MigrateAllMasterSchemaResponse;
};

type MigrationMySQLServiceRollbackMasterSchema = {
  readonly methodName: string;
  readonly service: typeof MigrationMySQLService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_migration_mysql_service_pb.RollbackMasterSchemaRequest;
  readonly responseType: typeof proto_migration_mysql_service_pb.RollbackMasterSchemaResponse;
};

export class MigrationMySQLService {
  static readonly serviceName: string;
  static readonly MigrateAllMasterSchema: MigrationMySQLServiceMigrateAllMasterSchema;
  static readonly RollbackMasterSchema: MigrationMySQLServiceRollbackMasterSchema;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class MigrationMySQLServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  migrateAllMasterSchema(
    requestMessage: proto_migration_mysql_service_pb.MigrateAllMasterSchemaRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_migration_mysql_service_pb.MigrateAllMasterSchemaResponse|null) => void
  ): UnaryResponse;
  migrateAllMasterSchema(
    requestMessage: proto_migration_mysql_service_pb.MigrateAllMasterSchemaRequest,
    callback: (error: ServiceError|null, responseMessage: proto_migration_mysql_service_pb.MigrateAllMasterSchemaResponse|null) => void
  ): UnaryResponse;
  rollbackMasterSchema(
    requestMessage: proto_migration_mysql_service_pb.RollbackMasterSchemaRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_migration_mysql_service_pb.RollbackMasterSchemaResponse|null) => void
  ): UnaryResponse;
  rollbackMasterSchema(
    requestMessage: proto_migration_mysql_service_pb.RollbackMasterSchemaRequest,
    callback: (error: ServiceError|null, responseMessage: proto_migration_mysql_service_pb.RollbackMasterSchemaResponse|null) => void
  ): UnaryResponse;
}

