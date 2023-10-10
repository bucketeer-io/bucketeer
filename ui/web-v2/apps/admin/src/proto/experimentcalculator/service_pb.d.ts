// package: bucketeer.experimentcalculator
// file: proto/experimentcalculator/service.proto

import * as jspb from "google-protobuf";

export class BatchCalcRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BatchCalcRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BatchCalcRequest): BatchCalcRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BatchCalcRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BatchCalcRequest;
  static deserializeBinaryFromReader(message: BatchCalcRequest, reader: jspb.BinaryReader): BatchCalcRequest;
}

export namespace BatchCalcRequest {
  export type AsObject = {
  }
}

export class BatchCalcResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BatchCalcResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BatchCalcResponse): BatchCalcResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BatchCalcResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BatchCalcResponse;
  static deserializeBinaryFromReader(message: BatchCalcResponse, reader: jspb.BinaryReader): BatchCalcResponse;
}

export namespace BatchCalcResponse {
  export type AsObject = {
  }
}

