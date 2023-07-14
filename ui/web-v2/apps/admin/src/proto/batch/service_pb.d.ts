/*
 * Copyright 2023 The Bucketeer Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// package: bucketeer.batch
// file: proto/batch/service.proto

import * as jspb from "google-protobuf";

export class BatchJobRequest extends jspb.Message {
  getJob(): BatchJobMap[keyof BatchJobMap];
  setJob(value: BatchJobMap[keyof BatchJobMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BatchJobRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BatchJobRequest): BatchJobRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BatchJobRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BatchJobRequest;
  static deserializeBinaryFromReader(message: BatchJobRequest, reader: jspb.BinaryReader): BatchJobRequest;
}

export namespace BatchJobRequest {
  export type AsObject = {
    job: BatchJobMap[keyof BatchJobMap],
  }
}

export class BatchJobResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BatchJobResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BatchJobResponse): BatchJobResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BatchJobResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BatchJobResponse;
  static deserializeBinaryFromReader(message: BatchJobResponse, reader: jspb.BinaryReader): BatchJobResponse;
}

export namespace BatchJobResponse {
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
}

export const BatchJob: BatchJobMap;

