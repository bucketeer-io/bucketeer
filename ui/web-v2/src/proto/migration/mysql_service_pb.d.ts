// package: bucketeer.migration
// file: proto/migration/mysql_service.proto

import * as jspb from "google-protobuf";

export class MigrateAllMasterSchemaRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrateAllMasterSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MigrateAllMasterSchemaRequest): MigrateAllMasterSchemaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MigrateAllMasterSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrateAllMasterSchemaRequest;
  static deserializeBinaryFromReader(message: MigrateAllMasterSchemaRequest, reader: jspb.BinaryReader): MigrateAllMasterSchemaRequest;
}

export namespace MigrateAllMasterSchemaRequest {
  export type AsObject = {
  }
}

export class MigrateAllMasterSchemaResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrateAllMasterSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: MigrateAllMasterSchemaResponse): MigrateAllMasterSchemaResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MigrateAllMasterSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrateAllMasterSchemaResponse;
  static deserializeBinaryFromReader(message: MigrateAllMasterSchemaResponse, reader: jspb.BinaryReader): MigrateAllMasterSchemaResponse;
}

export namespace MigrateAllMasterSchemaResponse {
  export type AsObject = {
  }
}

export class RollbackMasterSchemaRequest extends jspb.Message {
  getStep(): number;
  setStep(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RollbackMasterSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RollbackMasterSchemaRequest): RollbackMasterSchemaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RollbackMasterSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RollbackMasterSchemaRequest;
  static deserializeBinaryFromReader(message: RollbackMasterSchemaRequest, reader: jspb.BinaryReader): RollbackMasterSchemaRequest;
}

export namespace RollbackMasterSchemaRequest {
  export type AsObject = {
    step: number,
  }
}

export class RollbackMasterSchemaResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RollbackMasterSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RollbackMasterSchemaResponse): RollbackMasterSchemaResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RollbackMasterSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RollbackMasterSchemaResponse;
  static deserializeBinaryFromReader(message: RollbackMasterSchemaResponse, reader: jspb.BinaryReader): RollbackMasterSchemaResponse;
}

export namespace RollbackMasterSchemaResponse {
  export type AsObject = {
  }
}

