// package: bucketeer.batch
// file: proto/batch/service.proto

import * as jspb from 'google-protobuf';
import * as google_api_annotations_pb from '../../google/api/annotations_pb';

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
  TAGDELETER: 18;
}

export const BatchJob: BatchJobMap;
