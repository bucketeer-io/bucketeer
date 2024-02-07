// package: bucketeer.batch
// file: proto/batch/service.proto

import * as jspb from 'google-protobuf';

export class BatchJobRequest extends jspb.Message {
  getJob(): BatchJobMap[keyof BatchJobMap];
  setJob(value: BatchJobMap[keyof BatchJobMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BatchJobRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: BatchJobRequest
  ): BatchJobRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: BatchJobRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): BatchJobRequest;
  static deserializeBinaryFromReader(
    message: BatchJobRequest,
    reader: jspb.BinaryReader
  ): BatchJobRequest;
}

export namespace BatchJobRequest {
  export type AsObject = {
    job: BatchJobMap[keyof BatchJobMap];
  };
}

export class BatchJobResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BatchJobResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: BatchJobResponse
  ): BatchJobResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: BatchJobResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): BatchJobResponse;
  static deserializeBinaryFromReader(
    message: BatchJobResponse,
    reader: jspb.BinaryReader
  ): BatchJobResponse;
}

export namespace BatchJobResponse {
  export type AsObject = {};
}

export class MigrationVersionRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrationVersionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MigrationVersionRequest): MigrationVersionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MigrationVersionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrationVersionRequest;
  static deserializeBinaryFromReader(message: MigrationVersionRequest, reader: jspb.BinaryReader): MigrationVersionRequest;
}

export namespace MigrationVersionRequest {
  export type AsObject = {
  }
}

export class MigrationVersionResponse extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): void;

  getDirty(): boolean;
  setDirty(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrationVersionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: MigrationVersionResponse): MigrationVersionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MigrationVersionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrationVersionResponse;
  static deserializeBinaryFromReader(message: MigrationVersionResponse, reader: jspb.BinaryReader): MigrationVersionResponse;
}

export namespace MigrationVersionResponse {
  export type AsObject = {
    version: number,
    dirty: boolean,
  }
}

export class MigrationRequest extends jspb.Message {
  getDirection(): MigrationRequest.DirectionMap[keyof MigrationRequest.DirectionMap];
  setDirection(value: MigrationRequest.DirectionMap[keyof MigrationRequest.DirectionMap]): void;

  getSteps(): number;
  setSteps(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MigrationRequest): MigrationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MigrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrationRequest;
  static deserializeBinaryFromReader(message: MigrationRequest, reader: jspb.BinaryReader): MigrationRequest;
}

export namespace MigrationRequest {
  export type AsObject = {
    direction: MigrationRequest.DirectionMap[keyof MigrationRequest.DirectionMap],
    steps: number,
  }

  export interface DirectionMap {
    UP: 0;
    DOWN: 1;
  }

  export const Direction: DirectionMap;
}

export class MigrationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: MigrationResponse): MigrationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MigrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrationResponse;
  static deserializeBinaryFromReader(message: MigrationResponse, reader: jspb.BinaryReader): MigrationResponse;
}

export namespace MigrationResponse {
  export type AsObject = {
  }
}

export interface BatchJobMap {
  EXPERIMENTSTATUSUPDATER: 0;
  EXPERIMENTRUNNINGWATCHER: 1;
  FEATURESTALEWATCHER: 2;
  MAUCOUNTWATCHER: 3;
  DATETIMEWATCHER: 4;
  EVENTCOUNTWATCHER: 5;
  DOMAINEVENTINFORMER: 6;
  REDISCOUNTERDELETER: 7;
  PROGRESSIVEROLLOUTWATCHER: 8;
  EXPERIMENTCALCULATOR: 9;
  MAUSUMMARIZER: 10;
  MAUPARTITIONDELETER: 11;
  MAUPARTITIONCREATOR: 12;
  FEATUREFLAGCACHER: 13;
  SEGMENTUSERCACHER: 14;
  APIKEYCACHER: 15;
  AUTOOPSRULESCACHER: 16;
  EXPERIMENTCACHER: 17;
}

export const BatchJob: BatchJobMap;
